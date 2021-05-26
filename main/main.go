package main

import (
	JSON "encoding/json"
	"fmt"
	"net/http"

	httpserverlib "yeyu.local/http-server-lib"
)

func strStartsWith(str, test string) bool {
	switch {
	case len(str) < len(test):
		return false
	default:
		return str[:len(test)] == test
	}
}

func beautify(jsonRequestBody *httpserverlib.JsonRequestParams) error {
	rw := jsonRequestBody.Writer
	m := jsonRequestBody.JSON
	jsonBody := *m

	if bytes, err := JSON.MarshalIndent(jsonBody, "", "  "); err != nil {
		return err
	} else {
		fmt.Fprintf(rw, "%v", string(bytes))
	}
	return nil
}

func handleError(ep *httpserverlib.ErrorParams) {
	rw := ep.Writer
	err := ep.Error.Error()

	if strStartsWith(err, "json: cannot unmarshal array") {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(rw, "Cannot parse JSON array!")
		return
	}

	if strStartsWith(err, "json: cannot unmarshal object") {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(rw, "Cannot parse JSON object!")
		return
	}

	rw.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(rw, err)
}

func main() {
	httpserverlib.Ping()
	httpserverlib.WithJsonOn("POST", "/beautify", beautify)
	httpserverlib.AddErrorHander(handleError)
	httpserverlib.Start(4000)
}
