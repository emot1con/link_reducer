package utils

import (
	"encoding/json"
	"net/http"
)

func ReadJSON(r http.Request, data any) error {
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return err
	}
	return nil
}

func WriteJSON(w http.ResponseWriter, data any) error {
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return err
	}
	return nil
}

func WriteError(w http.ResponseWriter, err error) error {
	if err := WriteJSON(w, map[string]any{"error": err.Error()}); err != nil {
		return err
	}
	return nil
}
