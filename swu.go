package swu

import (
	"encoding/json"
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
	ID       string       `json:"id"`
	Tags     []string     `json:"tags"`
	Created  int64        `json:"created"`
	Versions []SWEVersion `json:"versions"`
	Name     string       `json:"name"`
}

type SWEVersion struct {
	Name      string `json:"name"`
	ID        string `json:"id"`
	Created   int64  `json:"created"`
	HTML      string `json:"html"`
	Text      string `json:"text"`
	Subject   string `json:"subject"`
	Published bool   `json:"published"`
}

type SWUError struct {
	Code    int
	Message string
}

func (e *SWUError) Error() string {
	return fmt.Sprintf("swu.go: Status code: %d, Error: %s", e.Code, e.Message)
}

func New(apiKey string) *SWUClient {
	return &SWUClient{
		Client: http.DefaultClient,
		apiKey: apiKey,
		URL:    SWUEndpoint,
	}
}

func (c *SWUClient) Templates() ([]SWUEmail, error) {
	return c.Emails()
}

func (c *SWUClient) Emails() ([]SWUEmail, error) {
	res, err := c.makeRequest("GET", "/templates", nil)
	if err != nil {
		return nil, err
	}
	var parse []SWUEmail
	err = buildRespJSON(res.Body, &parse)
	if err != nil {
		return nil, err
	}
	return parse, err
}

func (c *SWUClient) makeRequest(method, endpoint string, body io.Reader) (*http.Response, error) {
	r, _ := http.NewRequest(method, c.URL+endpoint, body)
	r.SetBasicAuth(c.apiKey, "")
	r.Header.Set("X-SWU-API-CLIENT", APIHeaderClient)
	res, err := c.Client.Do(r)
	if err != nil {
		return nil, &SWUError{
			Code:    res.StatusCode,
			Message: err.Error(),
		}
	}
	if res.StatusCode >= 300 {
		return nil, buildError(res)
	}
	return res, nil
}

func buildError(resp *http.Response) error {
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &SWUError{
			Code:    resp.StatusCode,
			Message: err.Error(),
		}
	}
	return &SWUError{
		Code:    resp.StatusCode,
		Message: string(b),
	}
}

func buildRespJSON(body io.ReadCloser, parse interface{}) error {
	defer body.Close()
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, parse)
	if err != nil {
		return err
	}
	return nil
}
