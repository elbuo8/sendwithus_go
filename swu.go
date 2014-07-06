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
	Versions []SWUVersion `json:"versions"`
	Name     string       `json:"name"`
}

type SWUVersion struct {
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
	var parse []SWUEmail
	err := c.makeRequest("GET", "/templates", nil, &parse)
	return parse, err
}

func (c *SWUClient) GetTemplate(id string) (SWUEmail, error) {
	var parse SWUEmail
	err := c.makeRequest("GET", "/templates/"+id, nil, &parse)
	return parse, err
}

func (c *SWUClient) GetTemplateVersion(id, version string) (SWUVersion, error) {
	var parse SWUVersion
	err := c.makeRequest("GET", "/templates/"+id+"/versions/"+version, nil, &parse)
	return parse, err
}

func (c *SWUClient) makeRequest(method, endpoint string, body io.Reader, result interface{}) error {
	r, _ := http.NewRequest(method, c.URL+endpoint, body)
	r.SetBasicAuth(c.apiKey, "")
	r.Header.Set("X-SWU-API-CLIENT", APIHeaderClient)
	res, err := c.Client.Do(r)
	if err != nil {
		return &SWUError{
			Code:    res.StatusCode,
			Message: err.Error(),
		}
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return &SWUError{
			Code:    res.StatusCode,
			Message: err.Error(),
		}
	}
	if res.StatusCode >= 300 {
		return &SWUError{
			Code:    res.StatusCode,
			Message: string(b),
		}
	}
	return buildRespJSON(b, result)
}

func buildRespJSON(b []byte, parse interface{}) error {
	err := json.Unmarshal(b, parse)
	return err
}
