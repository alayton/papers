package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func readJSON(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(buf.Bytes())
}
