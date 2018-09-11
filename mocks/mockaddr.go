package mocks

type MockAddr struct{}

func (addr *MockAddr) Network() string {
	return "tcp"
}

func (addr *MockAddr) String() string {
	return "0.0.0.0"
}
