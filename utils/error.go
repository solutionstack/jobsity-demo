package utils

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
)

func HandlerError(w http.ResponseWriter, err error, code int) {
	errMsg := ErrorResponse{Message: errors.Wrap(err, "error in request").Error()}
	jsonErr, _ := json.Marshal(errMsg)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(jsonErr)
}

type ErrorResponse struct {
	Message string `json:"message"`
}
