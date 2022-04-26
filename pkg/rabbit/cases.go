package rabbit

import (
	"github.com/testground/sdk-go/run"
	"github.com/testground/sdk-go/runtime"
)

// A test composed of 2 instances: one is the producer, and the other the consumer

func OneOnOne(runenv *runtime.RunEnv, initCtx *run.InitContext) error {
	return runRabbitTest(runenv, initCtx, 50)
}

// A test composed of 4 instances: one is a consumer, and the other threee are producers
func FourToOne(runenv *runtime.RunEnv, initCtx *run.InitContext) error {
	return runRabbitTest(runenv, initCtx, 25)
}

// A test with composed with 2 instances, aimed to fail (no messages will be sent by the producer)
func FailingTimeout(runenv *runtime.RunEnv, initCtx *run.InitContext) error {
	return runRabbitTest(runenv, initCtx, 0)
}
