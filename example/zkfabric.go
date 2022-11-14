package main

import (
	"encoding/json"
	"github.com/hyperledger/fabric-sdk-go/example/zkclient"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

func zkTx() {
	var (
		err          error
		configFile   = "./fixtures/example-org1.yaml"
		channelID    = "mychannel2"
		orgUser      = "User1"
		OrgName      = "org1"
		chaincoodeID = "basic"
		setArgs      = []string{"invoke", "userA", "userB", "5"}
	)
	zkCli, err := zkclient.NewClient("http://101.251.223.190:8888")
	if err != nil {
		logger.Error(err)
		return
	}
	sdk, err := fabsdk.New(config.FromFile(configFile))
	if err != nil {
		logger.Error(err)
		return
	}
	clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(orgUser), fabsdk.WithOrg(OrgName))

	chClient, err := channel.New(clientChannelContext)
	if err != nil {
		logger.Error(err)
		return
	}
	//register user
	request := channel.Request{
		ChaincodeID: chaincoodeID,
		Fcn:         "addUser",
		Args: [][]byte{[]byte("userA"),
			[]byte("100")},
	}
	response, err := chClient.Execute(request)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Info("response: ", string(response.Payload))

	request = channel.Request{
		ChaincodeID: chaincoodeID,
		Fcn:         "addUser",
		Args: [][]byte{[]byte("userB"),
			[]byte("100")},
	}
	response, err = chClient.Execute(request)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Info("response: ", string(response.Payload))

	request = channel.Request{
		ChaincodeID: chaincoodeID,
		Fcn:         setArgs[0],
		Args: [][]byte{[]byte(setArgs[1]),
			[]byte(setArgs[2]), []byte(setArgs[3])},
	}
	response, err = chClient.Execute(request)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Info("response: ", string(response.Payload))

	// init zk account for userA
	err = zkCli.InitAccount(100)
	if err != nil {
		return
	}
	// generate proof for userA
	tx, err := zkCli.SendZkTx(zkclient.ZkTx{
		From:       "userA",
		ValueTrans: 100,
		ValueNew:   10,
	})
	if err != nil {
		return
	}

	dataByte, _ := json.Marshal(tx)

	request = channel.Request{
		ChaincodeID: "sacc",
		Fcn:         "set",
		Args: [][]byte{[]byte(tx.SNSN),
			dataByte},
	}
	response, err = chClient.Execute(request)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Info("response: ", string(response.Payload))

}
