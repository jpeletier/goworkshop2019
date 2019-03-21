package timeline_test

import (
	"goworkshop/gossip/kv/account"
	"goworkshop/gossip/timeline"
	"testing"
)

func TestTimeline(t *testing.T) {

	account := account.New("someaccount")
	tm := timeline.NewMock(account)

	testMessages := []string{"My first message", "a second one", "the third message!!"}

	for _, msg := range testMessages {
		err := tm.Post(msg)
		if err != nil {
			t.Fatal(err)
		}
	}

	comments := tm.Dump(account.Addr())
	for i := len(testMessages) - 1; i > 0; i-- {
		comment, ok := <-comments
		if !ok {
			t.Fatal("channel closed prematurely")
		}
		if comment.Text != testMessages[i] {
			t.Fatalf("Expected comment text to be %s, got %s", testMessages[i], comment.Text)
		}
	}
}
