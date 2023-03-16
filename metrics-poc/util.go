package main

import (
	"encoding/json"
	"io"
)

func parseRequest(body io.Reader, target any) error {
	bb, err := io.ReadAll(body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bb, target)
	if err != nil {
		return err
	}
	return nil
}

func writeResponse(w io.Writer, body any) error {
	resp, err := json.Marshal(body)
	if err != nil {
		return err
	}

	_, err = w.Write(resp)
	if err != nil {
		return err
	}
	return nil
}
