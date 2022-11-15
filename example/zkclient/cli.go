package zkclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/logging"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

var (
	logger = logging.NewLogger("zk")

	UserAddress = "0x71c7656ec7ab88b098defb751b7401b5f6d8976f"
	timeout     = time.Second * 120
)

type Client struct {
	addr string
}

func NewClient(addr string) (*Client, error) {
	_, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}
	addr = strings.TrimSuffix(addr, "/")
	return &Client{addr}, nil
}

func (zk Client) doReq(thisUrl string, reqBody []byte) (map[string]interface{}, error) {

	req, err := http.NewRequest(http.MethodPost, thisUrl, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to create request body for init_account")
	}
	req.Header.Set("Content-Type", "application/json")

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
	respMap := map[string]interface{}{}
	if err = json.Unmarshal(respBody, &respMap); err != nil {
		return nil, errors.WithMessagef(err, "failed to get map from bytes")
	}
	return respMap, nil
}

func (zk Client) InitAccount(initValue int) error {
	reqBody := map[string]interface{}{
		"initvalue": initValue,
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return errors.WithMessagef(err, "josn cannot marshal")
	}
	thisUrl := fmt.Sprintf("%s/init_account", zk.addr)

	respMap, err := zk.doReq(thisUrl, bodyBytes)
	if err != nil {
		return errors.WithMessagef(err, "josn cannot marshal")
	}

	v, ok := respMap["initres"]
	switch {
	case !ok:
		return errors.Errorf("unexpected response, %v", respMap)
	case fmt.Sprintf("%v", v) != "1":
		return errors.Errorf("wrong response, %v", respMap)
	}

	return nil
}

type ZkTx struct {
	From       string `json:"from"`
	ValueTrans int    `json:"valuetrans"`
	ValueNew   int    `json:"valuenew"`
}

func (zk Client) SendZkTx(tx ZkTx) (*VerifiedDoc, error) {
	reqBody, err := json.Marshal(tx)
	if err != nil {
		return nil, errors.WithMessagef(err, "josn cannot marshal zk tx")
	}
	thisUrl := fmt.Sprintf("%s/send_zktx", zk.addr)

	respMap, err := zk.doReq(thisUrl, reqBody)
	if err != nil {
		return nil, errors.WithMessagef(err, "josn cannot marshal")
	}
	logger.Info("get zk response", respMap)
	if v, ok := respMap["status"]; ok && fmt.Sprintf("%v", v) == "0" {
		return nil, errors.Errorf("%s,failed to send zk tx", respMap["msg"])
	}
	//if v, ok := respMap["SNSN"]; !ok || fmt.Sprintf("%v", v) == "" {
	//	return nil, errors.Errorf("not contain SNSN in response")
	//}
	if v, ok := respMap["SNCMTOld"]; !ok || fmt.Sprintf("%v", v) == "" {
		return nil, errors.Errorf("not contain SNCMTOld in response")
	}
	if v, ok := respMap["CMTs"]; !ok || fmt.Sprintf("%v", v) == "" {
		return nil, errors.Errorf("not contain CMTs in response")
	}
	if v, ok := respMap["zkProof"]; !ok || fmt.Sprintf("%v", v) == "" {
		return nil, errors.Errorf("not contain zkProof in response")
	}
	if v, ok := respMap["CMTNew"]; !ok || fmt.Sprintf("%v", v) == "" {
		return nil, errors.Errorf("not contain CMTNew in response")
	}
	respBytes, _ := json.Marshal(respMap)
	txResp := VerifiedDoc{}
	_ = json.Unmarshal(respBytes, &txResp)
	return &txResp, nil
}

type VerifiedDoc struct {
	SNSN     string `json:"SNSN"`
	SNCMTOld string `json:"SNCMTOld"`
	CMTs     string `json:"CMTs"`
	ZkProof  string `json:"zkProof"`
	CMTNew   string `json:"CMTNew"`
}

func (zk Client) VerifyZkTx(vd VerifiedDoc) (bool, error) {
	reqBody, err := json.Marshal(vd)
	if err != nil {
		return false, errors.WithMessagef(err, "josn cannot marshal zk tx")
	}
	thisUrl := fmt.Sprintf("%s/verify_zktx", zk.addr)

	respMap, err := zk.doReq(thisUrl, reqBody)
	if err != nil {
		return false, errors.WithMessagef(err, "josn cannot marshal")
	}

	if v, ok := respMap["status"]; !ok || fmt.Sprintf("%v", v) == "1" {
		return false, errors.Errorf("response status is %v,expected 1", v)
	}

	return true, nil
}
