package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

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
		channelID    = "mychannel"
		orgUser      = "User1"
		OrgName      = "org1"
		chaincoodeID = "sacc"
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

	//双花测试
	wait := sync.WaitGroup{}

	for i := 0; i <= 1; i++ {
		wait.Add(1)

		go func(index int) {
			defer wait.Done()
			args := []string{"set", "asdgf7sdfa657fujhg", fmt.Sprintf("sdafasd-%d", index)}
			request := channel.Request{
				ChaincodeID: chaincoodeID,
				Fcn:         args[0],
				Args: [][]byte{[]byte(args[1]),
					[]byte(args[2])},
			}
			response, err := chClient.Execute(request, channel.WithTargetEndpoints("peer0.org1.example.com"))
			if err != nil {
				logger.Error(err)
				return
			}
			logger.Info("i= ", index, " tx_id ", response.TransactionID, "response: ", string(response.Payload))
		}(i)

	}
	wait.Wait()
	//// 初始化数据
	rand.Seed(time.Now().Unix())

	//for i := 1; i < 100; i++ {
	//
	//	str := fmt.Sprintf("%d", rand.Int())
	//	setArgs := []string{"set", str, "v+" + str}
	//	request := channel.Request{
	//		ChaincodeID: chaincoodeID,
	//		Fcn:         setArgs[0],
	//		Args: [][]byte{[]byte(setArgs[1]),
	//			[]byte(setArgs[2])},
	//	}
	//	logger.Info("k: ", str, "v: ", "v+"+str)
	//	response, err := chClient.Execute(request)
	//	if err != nil {
	//		logger.Error(err)
	//		return
	//	}
	//	logger.Info("response: ", string(response.Payload))
	//
	//}

}
