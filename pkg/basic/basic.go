package basic

import (
	"context"
	"time"

	"github.com/testground/sdk-go/run"
	"github.com/testground/sdk-go/runtime"
)

func BasicTest(runenv *runtime.RunEnv, initCtx *run.InitContext) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client := initCtx.SyncClient
	seq := client.MustSignalAndWait(ctx, "ip-allocation", runenv.TestInstanceCount)

	runenv.RecordMessage("I am %d, and basic test is executed successfully!", seq)
	return nil
}
