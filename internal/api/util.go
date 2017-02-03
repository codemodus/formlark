package api

import (
	"encoding/json"
	"io"
	"net/http"
)

func closeBody(r *http.Request) {
	_ = r.Body.Close()
}

func decodeBody(r *http.Request, i interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(i); err != nil {
		return err
	}
	defer closeBody(r)

	return nil
}

func encodeBody(w io.Writer, i interface{}) error {
	b, err := json.Marshal(i)
	if err != nil {
		return err
	}

	if _, err = w.Write(b); err != nil {
		return err
	}

	return nil
}

func httpError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
