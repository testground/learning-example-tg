package tgsync

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/testground/learning-example-tg/pkg/types"
	"github.com/testground/learning-example-tg/pkg/util"
	"github.com/testground/learning-example/pkg/processor"
	"github.com/testground/learning-example/pkg/producer"

	"github.com/testground/learning-example/pkg/message"
	"github.com/testground/sdk-go/run"
	"github.com/testground/sdk-go/runtime"
	"github.com/testground/sdk-go/sync"
)

func RunTgSyncTest(runenv *runtime.RunEnv, initCtx *run.InitContext, messagesByNode int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var notifyChan = make(chan bool)

	go func() {
		runTgSyncTest(runenv, initCtx, ctx, messagesByNode)
		notifyChan <- true
	}()

	select {
	// timeout?
	case <-ctx.Done():
		return ctx.Err()
	// successful test
	case <-notifyChan:
		return nil
	}
}

// Runs a test between several nodes, where each producer node will create *messagesByNode*
// total messages, and the single consumer node will consume all of them
func runTgSyncTest(runenv *runtime.RunEnv, initCtx *run.InitContext, ctx context.Context, messagesByNode int) error {
	// create a bounded client to send messages between instances
	client := sync.MustBoundClient(ctx, runenv)
	defer client.Close()

	oldAddrs, err := net.InterfaceAddrs()
	if err != nil {
		return err
	}

	seq := client.MustSignalAndWait(ctx, "ip-allocation", runenv.TestInstanceCount)

	// Make sure that the IP addresses don't change unless we request it.
	if newAddrs, err := net.InterfaceAddrs(); err != nil {
		return err
	} else if !util.SameAddrs(oldAddrs, newAddrs) {
		return fmt.Errorf("interfaces changed")
	}

	runenv.RecordMessage("I am %d", seq)

	var prod producer.Producer
	var listn *TgSyncListener
	var procsr *processor.Processor
	var consumer *types.MockConsumer
	var totalMessages = (runenv.TestGroupInstanceCount - 1) * messagesByNode

	// form queue name in rabbit unique to this run, so we can avoid message conflicts
	var testQueueName = fmt.Sprintf("queue_%s", runenv.TestRun)

	st := sync.NewTopic(testQueueName, &message.DataMessage{})

	tch := make(chan *message.DataMessage)

	if seq == 1 {
		// ID 1 - consumer
		runenv.RecordMessage("Expecting %d messages by node, %d total", messagesByNode, totalMessages)
		// custom sync listener
		_, err = client.Subscribe(ctx, st, tch)
		if err != nil {
			panic(err)
		}
		listn = &TgSyncListener{ListenChannel: tch}
		// custom mock consumer
		consumer = &types.MockConsumer{TotalCount: totalMessages, DoneChannel: make(chan bool)}
		consumer.On("ConsumeMessage", mock.Anything).Return(nil)
		procsr = &processor.Processor{Producer: nil, Consumer: consumer, Listener: listn}
	} else {
		// ID 2 - producer
		prod = &TgSyncProducer{
			IdGen:  0,
			Client: client,
			Topic:  st,
		}
	}

	switch seq {
	case 1:
		runenv.RecordMessage("Listening for messages")
		go func() { procsr.StartProcessor() }()
	default:
		// runenv.RecordMessage("Doing nothing")
	}
	if err != nil {
		return err
	}

	testFunc := func() error {
		if seq == 1 {
			// Wait for done signal by the consumer
			done := <-consumer.DoneChannel
			if !done {
				return fmt.Errorf("expected all messages to be processed")
			}
		} else {
			for i := 0; i < messagesByNode; i++ {
				prod.ProduceMessage(fmt.Sprintf("Test message #%d #%d", i, seq))
			}
			runenv.RecordMessage("Finished producing messages")
		}
		return nil
	}
	err = testFunc()
	if err != nil {
		return err
	}

	return nil
}
