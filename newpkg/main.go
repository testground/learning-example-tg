package main

import (
	"github.com/testground/sdk-go/run"
	"github.com/testground/sdk-go/runtime"
)

var testcases = map[string]interface{}{
	"basic": run.InitializedTestCaseFn(overrideBuilderConfiguration),
}

func main() {
	run.InvokeMap(testcases)
}

// The basic test used here to demonstrate the use of different builder options in compositions
func overrideBuilderConfiguration(runenv *runtime.RunEnv, initCtx *run.InitContext) error {
	runenv.RecordMessage("Test finished successfully!")
	return nil
}
