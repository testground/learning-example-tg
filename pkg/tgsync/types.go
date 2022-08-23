package tgsync

import (
	"context"
	"fmt"

	"github.com/testground/learning-example/pkg/message"
	"github.com/testground/sdk-go/sync"
)

// A message listener that uses the Testgorund's sync service channels
// to receive messages
type TgSyncListener struct {
	// The channel to listen at
	ListenChannel <-chan *message.DataMessage
	// The channel to send messages to (basically the consumer)
	NotifyChannel chan<- *message.DataMessage
}

func (listener *TgSyncListener) ListenForMessages() {
	forever := make(chan bool)

	go func() {
		for chanMsg := range listener.ListenChannel {
			listener.NotifyChannel <- chanMsg
		}
	}()

	<-forever
}

func (listener *TgSyncListener) SetNotifyChannel(channel chan<- *message.DataMessage) {
	listener.NotifyChannel = channel
}

// A message producer that uses Testground's sync channels to produce messages
// for listeners/consumers
type TgSyncProducer struct {
	IdGen  int32
	Client *sync.DefaultClient
	Topic  *sync.Topic
}

func (prod *TgSyncProducer) ProduceMessage(data string) {
	prod.IdGen++
	msg := &message.DataMessage{
		Id:   fmt.Sprint(prod.IdGen),
		Data: data,
	}

	prod.Client.Publish(context.TODO(), prod.Topic, msg)
}

type TgSyncConsumer struct {
	IdGen int32
	// other fields go here as normal
	TotalCount  int
	DoneChannel chan bool
}

func (cons *TgSyncConsumer) ConsumeMessage(msg *message.DataMessage) {
	cons.TotalCount--
	if cons.TotalCount <= 0 {
		cons.DoneChannel <- true
	}
	fmt.Printf("Consumed message: %s data :%s\n", msg.Id, msg.Data)
}
