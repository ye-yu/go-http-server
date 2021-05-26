package httpserverlib

import (
	JSON "encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type JsonRequestParams struct {
	Writer  http.ResponseWriter
	Request *http.Request
	JSON    *map[string]interface{}
}
type JsonRequestHandler func(*JsonRequestParams) error

type ErrorParams struct {
	Method  string
	Pattern string
	Writer  http.ResponseWriter
	Request *http.Request
	Error   error
}
type ErrorHandler func(*ErrorParams)

func Demo(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}

	if req.Method != "GET" {
		body := ""
		for {
			size := 0
			byteBody := make([]byte, 10)
			size, _ = req.Body.Read(byteBody)
			if size == 0 {
				break
			}
			body = concat(body, string(byteBody[:size]))
		}

		fmt.Fprintf(w, "Body: %v", string(body))
	}
}

func Ping() {
	fmt.Println("I'm alive.")
}

func WithJsonOn(
	method string,
	pattern string,
	middlewares ...JsonRequestHandler) error {

	if method == "GET" {
		return errors.New("cannot parse JSON on GET request")
	}

	http.HandleFunc(pattern, func(w http.ResponseWriter, req *http.Request) {
		var reqBody = ""
		for {
			size := 0
			byteBody := make([]byte, 10)
			size, _ = req.Body.Read(byteBody)
			if size == 0 {
				break
			}
			reqBody = concat(reqBody, string(byteBody[:size]))
		}

		var jsonBody = make(map[string]interface{})
		if err := JSON.Unmarshal([]byte(reqBody), &jsonBody); err != nil {
			fmt.Printf("Cannot parse JSON body: %v\n", reqBody)
			HandleError(&ErrorParams{
				Method:  method,
				Pattern: pattern,
				Writer:  w,
				Request: req,
				Error:   err,
			})
			return
		}

		if req.Method != method {
			fmt.Printf("No %v method available for %v, body: %v\n", method, pattern, reqBody)
			return
		}

		fmt.Printf("Got JSON body: %v\n", reqBody)
		var jsonRequestParams JsonRequestParams
		jsonRequestParams.Request = req
		jsonRequestParams.Writer = w
		jsonRequestParams.JSON = &jsonBody

		for _, middleware := range middlewares {
			if err := middleware(&jsonRequestParams); err != nil {
				HandleError(&ErrorParams{
					Method:  method,
					Pattern: pattern,
					Writer:  w,
					Request: req,
					Error:   err,
				})
				break
			}
		}
	})

	return nil
}

var onError []ErrorHandler

func AddErrorHander(errorHandler ...ErrorHandler) {
	onError = append(onError, errorHandler...)
}

func HandleError(errors *ErrorParams) {
	for _, errorHandler := range onError {
		errorHandler(errors)
	}
}

func Start(port int) {
	http.ListenAndServe(concat(":", strconv.Itoa(port)), nil)
}
