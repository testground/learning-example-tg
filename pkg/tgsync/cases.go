package tgsync

import (
	"github.com/testground/sdk-go/run"
	"github.com/testground/sdk-go/runtime"
)

func OneOnOne(runenv *runtime.RunEnv, initCtx *run.InitContext) error {
	return RunTgSyncTest(runenv, initCtx, 10)
}

func ManyToOne(runenv *runtime.RunEnv, initCtx *run.InitContext) error {
	return RunTgSyncTest(runenv, initCtx, 10)
}
