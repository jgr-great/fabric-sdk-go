package main

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/common/logging"
)

var logger = logging.NewLogger("example")

func main() {
	zkTx()
}
