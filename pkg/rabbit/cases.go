package rabbit

import (
	"fmt"
	"strconv"

	"github.com/testground/sdk-go/network"
	"github.com/testground/sdk-go/run"
	"github.com/testground/sdk-go/runtime"
)

// A test composed of 2 instances: one is the producer, and the other the consumer
func OneOnOne(runenv *runtime.RunEnv, initCtx *run.InitContext) error {
	// We use the given "messages" parameter passed either through the command line, or using a default
	// value, specified in the manifest file
	var messagesParam = runenv.RunParams.StringParam("messages")
	messages, err := strconv.Atoi(messagesParam)

	if err != nil {
		return fmt.Errorf("invalid param 'messages': %s", messagesParam)
	}

	params := &RabbitTestParams{
		MessagesByNode: messages,
		RoutingPolicy:  network.AllowAll,
	}

	return runRabbitTest(runenv, initCtx, params)
}

// A test composed of 4 instances: one is a consumer, and the other threee are producers
func FourToOne(runenv *runtime.RunEnv, initCtx *run.InitContext) error {
	params := &RabbitTestParams{
		MessagesByNode: 50,
		RoutingPolicy:  network.AllowAll,
	}
	return runRabbitTest(runenv, initCtx, params)
}

// A test with composed with 2 instances, aimed to fail (no messages will be sent by the producer)
func FailingTimeout(runenv *runtime.RunEnv, initCtx *run.InitContext) error {
	params := &RabbitTestParams{
		MessagesByNode: 0,
		RoutingPolicy:  network.AllowAll,
	}
	return runRabbitTest(runenv, initCtx, params)
}

// A test with composed with 2 instances, aimed to fail (routing policy will block connection to rabbit broker)
func FailingPolicy(runenv *runtime.RunEnv, initCtx *run.InitContext) error {
	params := &RabbitTestParams{
		MessagesByNode: 50,
		RoutingPolicy:  network.DenyAll,
	}
	return runRabbitTest(runenv, initCtx, params)
}
