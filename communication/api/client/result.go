package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// MarshalError
const (
	MarshalError = `{"success":false,"error":"Cannot marshal json data.","result":null}`
)

// Response write json to ResponseWriter.
func writeResponse(w http.ResponseWriter, code int, content interface{}) {
	data := struct {
		Code    int         `json:"code"`
		Message interface{} `json:"message"`
	}{
		code,
		content,
	}

	message, err := json.Marshal(data)
	if err != nil {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, MarshalError)

		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(message)
}

func writeError(w http.ResponseWriter, err error) {
	writeResponse(w, http.StatusInternalServerError, err.Error())
}
