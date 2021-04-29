package infrastructure

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/websmee/rest-service/domain/remote"
)

const getObjectTimeout = 5 * time.Second

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
		httpClient: http.Client{Timeout: getObjectTimeout},
		baseURL:    baseURL,
	}
}

func (r remoteObjectRepository) GetByID(id int64) (*remote.Object, error) {
	responseBody, err := r.makeRequest("/objects/" + strconv.Itoa(int(id)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("GetByID request failed for id=%d", id))
	}

	var response remoteObjectResponse
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("GetByID unmarshal failed for id=%d", id))
	}

	return &remote.Object{
		ID:     response.ID,
		Online: response.Online,
	}, nil
}

func (r remoteObjectRepository) makeRequest(path string) ([]byte, error) {
	response, err := http.Get(r.baseURL + path)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("request failed code=%d, msg='%s'", response.StatusCode, responseBody))
	}

	return responseBody, nil
}
