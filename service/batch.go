package service

// Credit for a lot of the code goes to dkasabovn:
// https://github.com/googleapis/google-api-go-client/issues/179#issuecomment-1641490906
import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"

	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/googleapi"
)

const (
	basePath = "https://gmail.googleapis.com/"
)

type Service struct {
	Regular *gmail.Service
	client  http.Client
	Batch   *BatchEmailService
}

type BatchEmailService struct {
	s *Service
}

func NewBatchEmailService(s *Service) *BatchEmailService {
	rs := &BatchEmailService{s: s}
	return rs
}

func NewService(ctx context.Context) (*Service, error) {
	gsvc, client, err := GetService(ctx)
	if err != nil {
		return nil, err
	}
	svc := &Service{
		Regular: gsvc,
		client:  *client,
	}
	batch := NewBatchEmailService(svc)
	svc.Batch = batch
	return svc, nil
}

type BatchGetEmailsRequest struct {
	// A list of UsersMessagesGetCall requests to send to the batch api
	IDS []string `json:"-"`
}

type BatchEmailsCall struct {
	s                     *Service
	userId                string
	batchgetemailsrequest *BatchGetEmailsRequest
	urlParams_            url.Values
	ctx_                  context.Context
	header_               http.Header
}

func (r *BatchEmailService) Get(userId string, IDS []string) *BatchEmailsCall {
	c := &BatchEmailsCall{s: r.s, urlParams_: make(url.Values), batchgetemailsrequest: &BatchGetEmailsRequest{IDS: IDS}}
	c.userId = userId
	return c
}

// "full" "raw" "metadata"
func (c *BatchEmailsCall) Format(format string) *BatchEmailsCall {
	c.urlParams_.Set("format", format)
	return c
}

func (c *BatchEmailsCall) Fields(s ...googleapi.Field) *BatchEmailsCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

func (c *BatchEmailsCall) Context(ctx context.Context) *BatchEmailsCall {
	c.ctx_ = ctx
	return c
}

func (c *BatchEmailsCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *BatchEmailsCall) doRequest() (*http.Response, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	defer writer.Close()

	for _, id := range c.batchgetemailsrequest.IDS {
		part, err := writer.CreatePart(textproto.MIMEHeader{
			"Content-Type": {"application/http"},
			"Content-ID":   {id},
		})
		if err != nil {
			return nil, err
		}

		/* --- making single messages get request --- */
		reqHeaders := make(http.Header)
		for k, v := range c.header_ {
			reqHeaders[k] = v
		}
		reqHeaders.Set("User-Agent", c.s.Regular.UserAgent)

		urls := googleapi.ResolveRelative(c.s.Regular.BasePath, "gmail/v1/users/{userId}/messages/{id}")
		urls += "?" + c.urlParams_.Encode()
		req, err := http.NewRequest("GET", urls, nil)
		if err != nil {
			return nil, err
		}
		req.Header = reqHeaders
		googleapi.Expand(req.URL, map[string]string{
			"userId": c.userId,
			"id":     id,
		})
		/* --- making single messages get request: Output, req --- */

		// hopefully this just works; would be pog
		req.Write(part)
	}

	writer.Close()

	batchUrl := googleapi.ResolveRelative(c.s.Regular.BasePath, "batch/gmail/v1")
	req, err := http.NewRequest("POST", batchUrl, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", fmt.Sprintf("multipart/mixed; boundary=%s", writer.Boundary()))

	return c.s.client.Do(req.WithContext(c.ctx_))
}

func (c *BatchEmailsCall) Do() ([]*gmail.Message, error) {
	res, err := c.doRequest()
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}

	// b, _ := httputil.DumpResponse(res, true)
	// log.Println(string(b))

	_, params, err := mime.ParseMediaType(res.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}

	mpr := multipart.NewReader(res.Body, params["boundary"])

	ret := make([]*gmail.Message, 0)

	for part, err := mpr.NextPart(); err != io.EOF; part, err = mpr.NextPart() {
		buf := bufio.NewReader(part)
		resp, err := http.ReadResponse(buf, nil)
		if err != nil {
			log.Printf("error reading response from part: %v", err)
			continue
		}

		var message gmail.Message
		if err := json.NewDecoder(resp.Body).Decode(&message); err != nil {
			log.Printf("error reading gmail message from part response: %v", err)
			continue
		}

		ret = append(ret, &message)
	}

	return ret, nil
}
