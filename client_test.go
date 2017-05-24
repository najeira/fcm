package fcm

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
)

type testHttpClient struct {
	req *http.Request
	res *http.Response
	err error
}

func (c *testHttpClient) Do(req *http.Request) (*http.Response, error) {
	c.req = req
	return c.res, c.err
}

func TestClientSend(t *testing.T) {
	const ApiKey = "dummy api key"
	const To = "fcm_token"
	const Priority = "high"
	const PayloadBody = "payload body"

	c := New(ApiKey)

	// for test
	tc := &testHttpClient{
		res: &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			Body: ioutil.NopCloser(bytes.NewReader([]byte(`{
				"success": 1,
				"failure": 0,
				"canonical_ids": 0,
				"results": [{"message_id": "foobar"}]
			}`))),
		},
	}
	c.SetHttpClient(tc)

	res, err := c.Send(&Message{
		To:       To,
		Priority: Priority,
		Notification: Notification{
			Body: PayloadBody,
		},
	})
	if err != nil {
		t.Error(err)
	}
	if res.Success != 1 {
		t.Errorf("Success: got %d, expect %d", res.Success, 1)
	}
	if res.Failure != 0 {
		t.Errorf("Failure: got %d, expect %d", res.Failure, 0)
	}
	if res.CanonicalIds != 0 {
		t.Errorf("CanonicalIds: got %d, expect %d", res.CanonicalIds, 0)
	}

	auth := tc.req.Header.Get("Authorization")
	if auth != "key="+ApiKey {
		t.Errorf("invalid Authoriazation %s", auth)
	}

	reqBody, err := ioutil.ReadAll(tc.req.Body)
	if err != nil {
		t.Error(err)
	}

	var ntf Message
	if err := json.Unmarshal(reqBody, &ntf); err != nil {
		t.Error(err)
	}
	if ntf.To != To {
		t.Errorf("To: got %d, expect %d", ntf.To, To)
	}
	if ntf.Priority != Priority {
		t.Errorf("Priority: got %d, expect %d", ntf.Priority, Priority)
	}
	if ntf.Notification.Body != PayloadBody {
		t.Errorf("PayloadBody: got %d, expect %d", ntf.Notification.Body, PayloadBody)
	}
}
