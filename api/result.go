package api

import (
	"fmt"
	"net/http"
	"encoding/json"
)

const (
	MarshalError = `{"success":false,"error":"Cannot marshal json data.","result":null}`
)
// Result will be converted to json string and write to http.ResponseWriter.
type Result struct {
	Success   bool                 `json:"success"`
	Error     error                `json:"error"` // error to show to user
	Result    interface{}          `json:"result"`
}

// Write json to ResponseWriter.
func (this *Result) Response(w http.ResponseWriter) {
	
	b, err := json.Marshal(this)
	if err != nil {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, MarshalError)

		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // always return status 200 even if an error occurs
	w.Write(b)
}