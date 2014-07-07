package swu

import (
	"bytes"
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

type SWUTemplate struct {
	ID       string        `json:"id,omitempty"`
	Tags     []string      `json:"tags,omitempty"`
	Created  int64         `json:"created,omitempty"`
	Versions []*SWUVersion `json:"versions,omitempty"`
	Name     string        `json:"name,omitempty"`
}

type SWUVersion struct {
	Name      string `json:"name,omitempty"`
	ID        string `json:"id,omitempty"`
	Created   int64  `json:"created,omitempty"`
	HTML      string `json:"html,omitempty"`
	Text      string `json:"text,omitempty"`
	Subject   string `json:"subject,omitempty"`
	Published bool   `json:"published,omitempty"`
}

type SWUEmail struct {
	ID          string            `json:"email_id,omitempty"`
	Recipient   SWURecipient      `json:"recipient,omitempty"`
	CC          []*SWURecipient   `json:"cc,omitempty"`
	BCC         []*SWURecipient   `json:"bcc,omitempty"`
	Sender      *SWUSender        `json:"sender,omitempty"`
	EmailData   map[string]string `json:"email_data,omitempty"`
	Tags        []string          `json:"tags,omitempty"`
	Inline      SWUAttachment     `json:"inline,omitempty"`
	Files       []*SWUAttachment  `json:"files,omitempty"`
	ESPAccount  string            `json:"esp_account,omitempty"`
	VersionName string            `json:"version_name,omitempty"`
}

type SWURecipient struct {
	Address string `json:"address,omitempty"`
	Name    string `json:"name,omitempty"`
}

type SWUSender struct {
	SWURecipient
	ReplyTo string `json:"reply_to,omitempty"`
}

type SWUAttachment struct {
	ID   string `json:"id,omitempty"`
	Data string `json:"data,omitempty"`
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

func (c *SWUClient) Templates() ([]*SWUTemplate, error) {
	return c.Emails()
}

func (c *SWUClient) Emails() ([]*SWUTemplate, error) {
	var parse []*SWUTemplate
	err := c.makeRequest("GET", "/templates", nil, &parse)
	return parse, err
}

func (c *SWUClient) GetTemplate(id string) (*SWUTemplate, error) {
	var parse SWUTemplate
	err := c.makeRequest("GET", "/templates/"+id, nil, &parse)
	return &parse, err
}

func (c *SWUClient) GetTemplateVersion(id, version string) (*SWUVersion, error) {
	var parse SWUVersion
	err := c.makeRequest("GET", "/templates/"+id+"/versions/"+version, nil, &parse)
	return &parse, err
}

func (c *SWUClient) UpdateTemplateVersion(id, version string, template *SWUVersion) (*SWUVersion, error) {
	var parse SWUVersion
	payload, err := json.Marshal(template)
	if err != nil {
		return nil, err
	}
	err = c.makeRequest("PUT", "/templates/"+id+"/versions/"+version, bytes.NewReader(payload), &parse)
	return &parse, err
}

func (c *SWUClient) CreateTemplate(template *SWUVersion) (*SWUTemplate, error) {
	var parse SWUTemplate
	payload, err := json.Marshal(template)
	if err != nil {
		return nil, err
	}
	err = c.makeRequest("POST", "/templates", bytes.NewReader(payload), &parse)
	return &parse, err
}

func (c *SWUClient) CreateTemplateVersion(id string, template *SWUVersion) (*SWUTemplate, error) {
	var parse SWUTemplate
	payload, err := json.Marshal(template)
	if err != nil {
		return nil, err
	}
	err = c.makeRequest("POST", "/templates/"+id+"/versions", bytes.NewReader(payload), &parse)
	return &parse, err
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
