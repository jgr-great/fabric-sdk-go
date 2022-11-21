package main

import (
	"encoding/json"
	"github.com/hyperledger/fabric-sdk-go/example/zkclient"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

func ZkTx() {
	var (
		err          error
		configFile   = "./fixtures/example-org1.yaml"
		channelID    = "mychannel"
		orgUser      = "User1"
		userAAdress  = "0x71c7656ec7ab88b098defb751b7401b5f6d8976f"
		OrgName      = "org1"
		chaincoodeID = "abstore"
		setArgs      = []string{"invoke", "userA", "userB", "5"}
	)
	logger.Info("init zk client")
	zkCli, err := zkclient.NewClient("http://101.251.223.190:8888")
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Info("new fabric client")
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
	logger.Info("register userA")
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
	logger.Info("response: ", string(response.TransactionID))
	logger.Info("register userB")
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
	txID := response.TransactionID
	logger.Info("fabric response: ", string(txID))

	logger.Info("userA transfer 5 token to userB")
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
	logger.Info("fabric response: ", string(response.TransactionID))

	// init zk account for userA
	logger.Info("init zk account for userA")
	err = zkCli.InitAccount(95)
	if err != nil {
		logger.Error("zk InitAccount error", err)
		return
	}
	// generate proof for userA
	logger.Info("send zk tx for userA")
	tx, err := zkCli.SendZkTx(zkclient.ZkTx{
		From:       userAAdress,
		ValueTrans: 50,
		ValueNew:   45,
	})
	if err != nil {
		logger.Error("zk SendZkTx error", err)
		return
	}
	logger.Info("zk response", tx)
	dataByte, _ := json.Marshal(tx)

	logger.Info("save zk proof for userA, record tx id and zk proof")
	request = channel.Request{
		ChaincodeID: "sacc",
		Fcn:         "set",
		Args: [][]byte{[]byte(response.TransactionID),
			dataByte},
	}
	response, err = chClient.Execute(request)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Info("response: ", string(response.TransactionID))

}
