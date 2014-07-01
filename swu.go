package swu

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	SWUEndpoint     = "https://api.sendwithus.com/api/v1"
	APIHeaderClient = "golang-0.0.1"
)

type SWUClient struct {
	Client *http.Client
	apiKey string
	URL    string
}

type SWUEmail struct {
	Tags     []string `json:"tags"`
	Created  int64    `json:"created"`
	Versions []struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"versions"`
}

type SWUError struct {
	Code    int
	Message string
}

func (e *SWUError) Error() string {
	return fmt.Printf("swg.go: Status code: %d, Error: %s", e.Code, e.Message)
}

func New(apiKey string) *SWUClient {
	return &SWUClient{
		Client: http.DefaultClient,
		apiKey: apiKey,
		URL:    SWUEndpoint,
	}
}

func (c *SWUClient) Emails() {

}

func (c *SWUClient) makeRequest(method, endpoint string, body io.Reader) (*http.Response, error) {
	r, _ := http.NewRequest(method, c.URL+endpoint, body)
	r.SetBasicAuth(c.apiKey, "")
	r.Header.Set("X-SWU-API-CLIENT", APIHeaderClient)
	res, err := c.Client.Do(r)
	if err != nil {
		return nil, &SWUError{
			Code:    res.StatusCode,
			Message: err,
		}
	}
}
