package swu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/google/go-querystring/query"
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
	Recipient   *SWURecipient     `json:"recipient,omitempty"`
	CC          []*SWURecipient   `json:"cc,omitempty"`
	BCC         []*SWURecipient   `json:"bcc,omitempty"`
	Sender      *SWUSender        `json:"sender,omitempty"`
	EmailData   map[string]string `json:"email_data,omitempty"`
	Tags        []string          `json:"tags,omitempty"`
	Inline      *SWUAttachment    `json:"inline,omitempty"`
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

type SWULogEvent struct {
	Object  string `json:"object,omitempty"`
	Created int64  `json:"created,omitempty"`
	Type    string `json:"type,omitempty"`
	Message string `json:"message,omitempty"`
}

type SWULogQuery struct {
	Count      int   `json:"count,omitempty" url:"count,omitempty"`
	Offset     int   `json:"offset,omitempty" url:"offset,omitempty"`
	CreatedGT  int64 `json:"created_gt,omitempty" url:"created_gt,omitempty"`
	CreatedGTE int64 `json:"created_gte,omitempty" url:"created_gte,omitempty"`
	CreatedLT  int64 `json:"created_lt,omitempty" url:"created_lt,omitempty"`
	CreatedLTE int64 `json:"created_lte,omitempty" url:"created_lte,omitempty"`
}

type SWULog struct {
	SWULogEvent
	ID               string `json:"id,omitempty"`
	RecipientName    string `json:"recipient_name,omitempty"`
	RecipientAddress string `json:"recipient_address,omitempty"`
	Status           string `json:"status,omitempty"`
	EmailID          string `json:"email_id,omitempty"`
	EmailName        string `json:"email_name,omitempty"`
	EmailVersion     string `json:"email_version,omitempty"`
	EventsURL        string `json:"events_url,omitempty"`
}

type SWULogResend struct {
	Success bool   `json:"success,omitempty"`
	Status  string `json:"status,omitempty"`
	ID      string `json:"log_id,omitempty"`
	Email   struct {
		Name        string `json:"name"`
		VersionName string `json:"version_name"`
	} `json:"email"`
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

func (c *SWUClient) Send(email *SWUEmail) error {
	payload, err := json.Marshal(email)
	if err != nil {
		return err
	}
	err = c.makeRequest("POST", "/send", bytes.NewReader(payload), nil)
	return err
}

func (c *SWUClient) ActivateDripCampaign(id string, email *SWUEmail) error {
	payload, err := json.Marshal(email)
	if err != nil {
		return err
	}
	err = c.makeRequest("POST", "/drip_campaigns/"+id+"/activate", bytes.NewReader(payload), nil)
	return err
}

func (c *SWUClient) GetLogs(q *SWULogQuery) ([]*SWULog, error) {
	var parse []*SWULog
	payload, _ := query.Values(q)
	err := c.makeRequest("GET", "/logs?"+payload.Encode(), nil, &parse)
	return parse, err
}

func (c *SWUClient) GetLog(id string) (*SWULog, error) {
	var parse SWULog
	err := c.makeRequest("GET", "/logs/"+id, nil, &parse)
	return &parse, err
}

func (c *SWUClient) GetLogEvents(id string) (*SWULogEvent, error) {
	var parse SWULogEvent
	err := c.makeRequest("GET", "/logs/"+id+"/events", nil, &parse)
	return &parse, err
}

func (c *SWUClient) ResendLog(id string) (*SWULogResend, error) {
	parse := &SWULogResend{
		ID: id,
	}
	payload, _ := json.Marshal(parse)
	err := c.makeRequest("POST", "/resend", bytes.NewReader(payload), parse)
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
	if result != nil {
		return buildRespJSON(b, result)
	}
	return nil
}

func buildRespJSON(b []byte, parse interface{}) error {
	err := json.Unmarshal(b, parse)
	return err
}
