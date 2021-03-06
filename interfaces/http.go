package interfaces

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pkg/errors"

	"github.com/websmee/rest-service/app"
)

type httpHandler struct {
	objectProcessor *app.ObjectProcessor
}

func NewHTTPHandler(objectProcessor *app.ObjectProcessor) *httpHandler {
	return &httpHandler{objectProcessor}
}

type callbackRequest struct {
	ObjectIDs []int64 `json:"object_ids"`
}

type callbackResponse struct {
	Success bool     `json:"success"`
	Errors  []string `json:"errors"`
}

func (h httpHandler) HandleCallback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.writeResponse(w, false, []error{errors.Wrap(err, "request read failed")})
		return
	}

	var cr callbackRequest
	if err := json.Unmarshal(body, &cr); err != nil {
		h.writeResponse(w, false, []error{errors.Wrap(err, "request unmarshal failed")})
		return
	}

	errs := h.objectProcessor.Process(cr.ObjectIDs)
	if len(errs) > 0 {
		h.writeResponse(w, false, errs)
		return
	}

	log.Printf("processed %d", len(cr.ObjectIDs))
	h.writeResponse(w, true, nil)
}

func (h httpHandler) writeResponse(w http.ResponseWriter, success bool, errors []error) {
	w.Header().Set("Content-Type", "application/json")

	if success {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	errorStrings := make([]string, len(errors))
	for i := range errors {
		errorStrings[i] = errors[i].Error()
		log.Println(errors[i])
	}

	b, _ := json.Marshal(callbackResponse{
		Success: success,
		Errors:  errorStrings,
	})
	_, _ = io.WriteString(w, string(b))
}
