package infrastructure

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/pkg/errors"

	"github.com/websmee/rest-service/domain/remote"
)

type remoteObjectRepository struct {
	httpClient http.Client
	baseURL    string
}

type remoteObjectResponse struct {
	ID     int64 `json:"id"`
	Online bool  `json:"online"`
}

func NewRemoteObjectRepository(baseURL string) remote.Repository {
	return &remoteObjectRepository{
		httpClient: http.Client{},
		baseURL:    baseURL,
	}
}

func (r remoteObjectRepository) GetByID(ctx context.Context, id int64) (*remote.Object, error) {
	respChan := make(chan []byte)
	errorChan := make(chan error)

	go r.makeRequest(http.MethodGet, "/objects/"+strconv.Itoa(int(id)), nil, respChan, errorChan)
	select {
	case <-ctx.Done():
		return nil, errors.New("GetByID context canceled")
	case err := <-errorChan:
		return nil, err
	case responseBody := <-respChan:
		var response remoteObjectResponse
		if err := json.Unmarshal(responseBody, &response); err != nil {
			return nil, errors.Wrap(err, "GetByID unmarshal failed")
		}

		return &remote.Object{
			ID:     response.ID,
			Online: response.Online,
		}, nil
	}
}

func (r remoteObjectRepository) makeRequest(method, path string, requestBody io.Reader, respChan chan<- []byte, errorChan chan<- error) {
	url := r.baseURL + path
	request, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		errorChan <- errors.Wrap(err, "makeRequest failed")
		return
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")

	response, err := r.httpClient.Do(request)
	if err != nil {
		errorChan <- errors.Wrap(err, "makeRequest do failed")
		return
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		errorChan <- errors.Wrap(err, "makeRequest read failed")
		return
	}

	respChan <- responseBody
}
