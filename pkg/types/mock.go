package types

import (
	"github.com/stretchr/testify/mock"
	"github.com/testground/learning-example/pkg/message"
)

// A mock consumer consumes all incoming messages and sends a "done" signal
// on its DoneChannel once the expected number of messages has been consumed
type MockConsumer struct {
	// add a Mock object instance
	mock.Mock
	// other fields go here as normal
	TotalCount  int
	DoneChannel chan bool
}

func (test *MockConsumer) ConsumeMessage(msg *message.DataMessage) {
	test.Called(msg)
	test.TotalCount--
	if test.TotalCount <= 0 {
		test.DoneChannel <- true
	}
}
