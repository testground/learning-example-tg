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

func runRabbitTest(runenv *runtime.RunEnv, initCtx *run.InitContext, messagesByNode int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	var notifyChan = make(chan bool)

	go func() {
		runTest(runenv, initCtx, ctx, messagesByNode)
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
func runTest(runenv *runtime.RunEnv, initCtx *run.InitContext, ctx context.Context, messagesByNode int) error {
	client := initCtx.SyncClient
	netclient := initCtx.NetClient

	oldAddrs, err := net.InterfaceAddrs()
	if err != nil {
		return err
	}

	config := &network.Config{
		// Control the "default" network. At the moment, this is the only network.
		Network: "default",

		// Enable this network. Setting this to false will disconnect this test
		// instance from this network. You probably don't want to do that.
		Enable: true,
		Default: network.LinkShape{
			Latency:   100 * time.Millisecond,
			Bandwidth: 1 << 20, // 1Mib
		},
		CallbackState: "network-configured",
		// Required: will not be able to connect to rabbitMQ otherwise
		RoutingPolicy: network.AllowAll,
	}

	runenv.RecordMessage("before netclient.MustConfigureNetwork")
	netclient.MustConfigureNetwork(ctx, config)

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
	var totalMessages = (runenv.TestGroupInstanceCount - 1) * messagesByNode

	// form queue name in rabbit unique to this run, so we can avoid message conflicts
	var rabbitQueueName = fmt.Sprintf("queue_%s", runenv.TestRun)
	// clean up the queue after tests
	rabbitConn := rabbit.GetConnection()
	defer rabbit.DeleteQueue(rabbitConn, rabbitQueueName)

	if seq == 1 {
		// ID 1 - consumer
		runenv.RecordMessage("Expecting %d messages by node, %d total", messagesByNode, totalMessages)
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
		runenv.RecordMessage("Doing nothing")
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
