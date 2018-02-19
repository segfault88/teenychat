package main

import "testing"

type mockClient struct {
}

func (mc *mockClient) getID() string {
	return "mock"
}

func (mc *mockClient) getDesc() string {
	return "mock"
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
		c := &mockClient{}
		chatHub.connect(c)
	}
}
