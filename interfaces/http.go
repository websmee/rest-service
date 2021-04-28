package interfaces

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pkg/errors"

	"github.com/websmee/rest-service/app"
)

type httpHandler struct {
	ctx     context.Context
	service *app.ObjectProcessor
}

func NewHTTPHandler(ctx context.Context, service *app.ObjectProcessor) *httpHandler {
	return &httpHandler{ctx, service}
}

type callbackRequest struct {
	ObjectIDs []int64 `json:"object_ids"`
}

type callbackResponse struct {
	Success bool     `json:"success"`
	Errors  []string `json:"errors"`
}

func (h httpHandler) HandleCallback(w http.ResponseWriter, r *http.Request) {
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

	errs := h.service.Process(h.ctx, cr.ObjectIDs)
	if len(errs) > 0 {
		h.writeResponse(w, false, errs)
		return
	}

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
