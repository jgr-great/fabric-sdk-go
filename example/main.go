package main

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/logging"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

var logger = logging.NewLogger("example")

func main() {
	var (
		err          error
		configFile   = "./fixtures/example-org1.yaml"
		channelID    = "mychannel2"
		orgUser      = "User1"
		OrgName      = "org1"
		chaincoodeID = "sacc_2"
		setArgs      = []string{"set", "a", "b"}
	)
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

	request := channel.Request{
		ChaincodeID: chaincoodeID,
		Fcn:         setArgs[0],
		Args: [][]byte{[]byte(setArgs[1]),
			[]byte(setArgs[2])},
	}
	response, err := chClient.Execute(request)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Info("response: ", string(response.Payload))
}
