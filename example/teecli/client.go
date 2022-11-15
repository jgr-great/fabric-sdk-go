package teecli

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	timeout = time.Second * 60
)

type Client struct {
	addr string
}

func NewClient(addr string) (*Client, error) {
	_, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}
	addr = strings.TrimSuffix(addr, "/") + "/ippinte/api/scene/getall"
	return &Client{addr}, nil
}

func (tc Client) doReq(thisUrl, method string, reqBody []byte) (*KeyPair, error) {

	req, err := http.NewRequest(method, thisUrl, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to create request body for init_account")
	}
	req.Header.Add("Content-Type", "text/plain")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	httpCli := &http.Client{
		Timeout: timeout,
	}
	resp, err := httpCli.Do(req)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to request init_account")
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to read response body for init account")
	}
	respMap := KeyPair{}
	if err = json.Unmarshal(respBody, &respMap); err != nil {
		return nil, errors.WithMessagef(err, "failed to get map from bytes")
	}
	return &respMap, nil
}

type KeyPair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (tc Client) Get(key string) (*KeyPair, error) {
	reqBody := map[string]interface{}{
		"key": key,
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, errors.WithMessagef(err, "josn cannot marshal")
	}

	respMap, err := tc.doReq(tc.addr, http.MethodGet, bodyBytes)
	if err != nil {
		return nil, errors.WithMessagef(err, "josn cannot marshal")
	}
	return respMap, nil
}

func (tc Client) Put(key, value string) (*KeyPair, error) {
	reqBody := map[string]interface{}{
		"key":   key,
		"value": value,
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, errors.WithMessagef(err, "josn cannot marshal")
	}

	respMap, err := tc.doReq(tc.addr, http.MethodGet, bodyBytes)
	if err != nil {
		return nil, errors.WithMessagef(err, "josn cannot marshal")
	}
	return respMap, nil
}
