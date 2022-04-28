package main

import (
	"github.com/testground/learning-example-tg/pkg/rabbit"
	"github.com/testground/learning-example-tg/pkg/tgsync"
	"github.com/testground/sdk-go/run"
)

var testcases = map[string]interface{}{
	"rabbit-1to1":            run.InitializedTestCaseFn(rabbit.OneOnOne),
	"rabbit-4to1":            run.InitializedTestCaseFn(rabbit.FourToOne),
	"rabbit-failing-timeout": run.InitializedTestCaseFn(rabbit.FailingTimeout),
	"rabbit-failing-policy":  run.InitializedTestCaseFn(rabbit.FailingPolicy),
	"tg-sync-1to1":           run.InitializedTestCaseFn(tgsync.OneOnOne),
	"tg-sync-manyto1":        run.InitializedTestCaseFn(tgsync.ManyToOne),
}

func main() {
	run.InvokeMap(testcases)
}
