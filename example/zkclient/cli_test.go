package zkclient

import "testing"

var cli *Client

func init() {
	var err error
	cli, err = NewClient("http://101.251.223.190:8888")
	if err != nil {

	}
}

func TestClient_InitAccount(t *testing.T) {
	err := cli.InitAccount(10)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestClient_SendZkTx(t *testing.T) {
	resp, err := cli.SendZkTx(ZkTx{
		From:       UserAddress,
		ValueTrans: 3,
		ValueNew:   7,
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(resp)
}
