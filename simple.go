package main

import (
	"github.com/testground/sdk-go/run"
	"github.com/testground/sdk-go/runtime"
)

// A test composed of 2 instances: one is the producer, and the other the consumer

func SimpleTest(runenv *runtime.RunEnv, initCtx *run.InitContext) error {
	return RunProcessingTest(runenv, initCtx, 50)
}

// A test composed of 4 instances: one is a consumer, and the other threee are producers
func FourToOne(runenv *runtime.RunEnv, initCtx *run.InitContext) error {
	return RunProcessingTest(runenv, initCtx, 25)
}
