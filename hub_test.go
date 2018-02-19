package main

import (
	"fmt"
	"testing"

	"github.com/rs/xid"
)

type mockClient struct {
	id string
}

func (mc *mockClient) getID() string {
	return mc.id
}

func (mc *mockClient) getDesc() string {
	return fmt.Sprintf("mock %s", mc.id)
}

func (mc *mockClient) joined(h *hub) {

}

func (mc *mockClient) close() {

}

func (mc *mockClient) sendMessage(m interface{}) error {
	return nil
}

func TestHub(t *testing.T) {
	chatHub := &hub{}
	chatHub.start()

	for i := 0; i < 10; i++ {
		id := xid.New()
		c := &mockClient{id: id.String()}
		chatHub.connect(c)
	}
}
