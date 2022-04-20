package main

import (
	"github.com/testground/sdk-go/run"
)

var testcases = map[string]interface{}{
	"simple": run.InitializedTestCaseFn(SimpleTest),
	"4to1":   run.InitializedTestCaseFn(FourToOne),
}

func main() {
	run.InvokeMap(testcases)
}
