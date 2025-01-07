package utils

import (
	"encoding/json"
	"github.com/Zhanbatyr06/PizzaGoland/utils"
	"net/http"
)

func SetError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusBadRequest)
	_, wErr := w.Write([]byte(err.Error()))
	if wErr != nil {
		utils.Logger.Error("Write error: %v", wErr)
	}
}

func SetJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	encodeErr := json.NewEncoder(w).Encode(v)
	if encodeErr != nil {
		utils.Logger.Error("Write error: %v", encodeErr)
	}
}

func SetOK(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{ok: 'true'"))
}
