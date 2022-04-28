package rabbit

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/testground/learning-example-tg/pkg/types"
	"github.com/testground/learning-example-tg/pkg/util"
	"github.com/testground/learning-example/pkg/listener"
	"github.com/testground/learning-example/pkg/processor"
	"github.com/testground/learning-example/pkg/producer"
	"github.com/testground/learning-example/pkg/rabbit"

	"github.com/testground/sdk-go/network"
	"github.com/testground/sdk-go/run"
	"github.com/testground/sdk-go/runtime"
)

type RabbitTestParams struct {
	MessagesByNode int
	RoutingPolicy  network.RoutingPolicyType
}

// Wraps a test in a context that will timeout after a set amount of time
func runRabbitTest(runenv *runtime.RunEnv, initCtx *run.InitContext, testParams *RabbitTestParams) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var notifyChan = make(chan bool)

	go func() {
		runTest(runenv, initCtx, ctx, testParams)
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
func runTest(runenv *runtime.RunEnv, initCtx *run.InitContext, ctx context.Context, testParams *RabbitTestParams) error {
	client := initCtx.SyncClient
	netclient := initCtx.NetClient

	oldAddrs, err := net.InterfaceAddrs()
	if err != nil {
		return err
	}

	// Configure network according to desired routing policy
	util.ConfigureNetwork(netclient, testParams.RoutingPolicy, ctx)

	seq := client.MustSignalAndWait(ctx, "ip-allocation", runenv.TestInstanceCount)

	// Make sure that the IP addresses don't change unless we request it.
	if newAddrs, err := net.InterfaceAddrs(); err != nil {
		return err
	} else if !util.SameAddrs(oldAddrs, newAddrs) {
		return fmt.Errorf("interfaces changed")
	}

	runenv.RecordMessage("I am %d", seq)

	var prod producer.Producer
	var listn listener.Listener
	var procsr *processor.Processor
	var consumer *types.MockConsumer
	var totalMessages = (runenv.TestGroupInstanceCount - 1) * testParams.MessagesByNode

	// form queue name in rabbit unique to this run, so we can avoid message conflicts
	var rabbitQueueName = fmt.Sprintf("queue_%s", runenv.TestRun)
	// clean up the queue after tests
	rabbitConn := rabbit.GetConnection()
	defer rabbit.DeleteQueue(rabbitConn, rabbitQueueName)

	if seq == 1 {
		// ID 1 - consumer
		runenv.RecordMessage("Expecting %d messages by node, %d total", testParams.MessagesByNode, totalMessages)
		listn = &listener.RabbitListener{QueueName: rabbitQueueName}

		consumer = &types.MockConsumer{TotalCount: totalMessages, DoneChannel: make(chan bool)}
		consumer.On("ConsumeMessage", mock.Anything).Return(nil)
		procsr = &processor.Processor{Producer: nil, Consumer: consumer, Listener: listn}
	} else {
		// ID 2 - producer
		prod = &producer.RabbitProducer{
			IdGen:     0,
			QueueName: rabbitQueueName,
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
			done := <-consumer.DoneChannel
			if !done {
				return fmt.Errorf("expected all messages to be processed")
			}
		} else {
			for i := 0; i < testParams.MessagesByNode; i++ {
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
