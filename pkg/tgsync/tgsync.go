package tgsync

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"net"
	"time"

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
	runenv.RecordMessage("Test plan started...")
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

	// determine instance role based on run params, seq ID, etc.
	var role = getInstanceRole(runenv, seq)

	var prod producer.Producer
	var listn *TgSyncListener
	var procsr *processor.Processor
	var consumer TgSyncConsumer
	var totalMessages = (runenv.TestInstanceCount - 1) * messagesByNode

	// form queue name unique to this run, so we can avoid message conflicts
	var testQueueName = fmt.Sprintf("queue_%s", runenv.TestRun)

	st := sync.NewTopic(testQueueName, &message.DataMessage{
		Id: uuid.New().String(),
	})

	tch := make(chan *message.DataMessage)

	if role == ConsumerRole {
		// consumer
		runenv.RecordMessage("Expecting %d messages by node, %d total", messagesByNode, totalMessages)
		// custom sync listener
		_, err = client.Subscribe(ctx, st, tch)
		if err != nil {
			return err
		}
		listn = &TgSyncListener{ListenChannel: tch}
		// custom mock consumer
		consumer = TgSyncConsumer{TotalCount: totalMessages, DoneChannel: make(chan bool), IdGen: int32(seq)}
		procsr = &processor.Processor{Producer: nil, Consumer: &consumer, Listener: listn}

		runenv.RecordMessage("Listening for messages")
		go func() { procsr.StartProcessor() }()
	} else {
		// producer
		prod = &TgSyncProducer{
			IdGen:  0,
			Client: client,
			Topic:  st,
		}
	}

	testFunc := func() error {
		if role == ConsumerRole {
			// Wait for done signal by the consumer
			done := <-consumer.DoneChannel
			if !done {
				return fmt.Errorf("expected all messages to be processed")
			}
		} else {
			for i := 0; i < messagesByNode; i++ {
				msg := fmt.Sprintf("Test message #%d #%d", i, seq)
				prod.ProduceMessage(msg)
				runenv.RecordMessage("Producing message:" + msg)
			}
			runenv.RecordMessage("Finished producing messages")
		}
		return nil
	}
	err = testFunc()
	if err != nil {
		return err
	}

	filePath, err := runenv.CreateRandomFile("./", 2*1024) // 2MB
	if err != nil {
		return err
	}
	fmt.Println("Output file created on location: " + filePath)
	return nil
}
