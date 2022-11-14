package main

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/example/teecli"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

func TeeTx() {
	var (
		err          error
		configFile   = "./fixtures/example-org1.yaml"
		channelID    = "mychannel"
		orgUser      = "User1"
		OrgName      = "org1"
		chaincoodeID = "sacc"
		key, value   = "a", "b"
		setArgs      = []string{"set", key, value}
	)
	sdk, err := fabsdk.New(config.FromFile(configFile))
	if err != nil {
		logger.Error(err)
		return
	}

	teeCli, err := teecli.NewClient("http://175.5.45.65:8888")
	if err != nil {
		logger.Error(err)
		return
	}

	put, err := teeCli.Put(key, value)
	if err != nil {
		return
	}
	fmt.Println("tee post: ", put)

	clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(orgUser), fabsdk.WithOrg(OrgName))

	chClient, err := channel.New(clientChannelContext)
	if err != nil {
		logger.Error(err)
		return
	}

	request := channel.Request{
		ChaincodeID: chaincoodeID,
		Fcn:         setArgs[0],
		Args: [][]byte{[]byte(key),
			[]byte("")},
	}
	response, err := chClient.Execute(request)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Info(" sacc response: ", string(response.Payload))

	res, err := teeCli.Get(key)
	if err != nil {
		return
	}
	fmt.Println("get from tee", res)
}
