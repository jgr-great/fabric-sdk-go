package teecli

import "testing"

func TestNewClient(t *testing.T) {

	teeCli, err := NewClient("http://175.5.59.81:8888")
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("put data into tee")
	put, err := teeCli.Put("aaa", "bbb")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("tee post: ", put)
}
