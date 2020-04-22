package jsonapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"reflect"
)

//NotFoundJSONHandler returns a JSON HTTP 404 Not Found response
var NotFoundJSONHandler = withJSONResponse(func(r *http.Request) (int, interface{}) {
	return http.StatusNotFound, nil
})

//withJSONResponse returns an http.Handler that writes JSON responses for the given ReturnHandlerFunc
func withJSONResponse(next ReturnHandlerFunc) http.Handler {
	type response struct {
		Code        int    `json:"code"`
		Description string `json:"description"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code, body := next(r)

		if err, ok := body.(error); ok || body == nil {
			resp := response{Code: code, Description: http.StatusText(code)}
			body = resp
			if err != nil {
				(r.Context().Value(contextKeyLogData)).(*logData).Error = err.Error()
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)

		e := json.NewEncoder(w)
		err := e.Encode(body)

		if err != nil {
			log.Println("Error writing JSON response:", err)
		}
	})
}

//ParseJSONBody decodes the http.Request body into one or more values
func ParseJSONBody(r *http.Request, vals ...interface{}) error {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return fmt.Errorf("Could not parse Content-Type: %v", err)
	}

	if mediaType != "application/json" {
		return errors.New("Content-Type not application/json")
	}

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, r.Body); err != nil {
		return fmt.Errorf("Unable to read request body: %v", err)
	}

	reader := bytes.NewReader(buf.Bytes())
	dec := json.NewDecoder(reader)

	for _, v := range vals {
		if err := dec.Decode(v); err != nil {
			return fmt.Errorf("Unable to parse request body to %s: %v", reflect.TypeOf(v).Elem().Name(), err)
		}
		if _, err = reader.Seek(0, 0); err != nil {
			return fmt.Errorf("Unable to reset byte reader: %v", err)
		}
	}

	return nil
}
