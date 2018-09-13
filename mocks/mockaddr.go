package mocks

type MockAddr struct {
	network string
	address string
}

func NewMockAddr(network, address string) *MockAddr {
	return &MockAddr{network, address}
}

func (addr *MockAddr) Network() string {
	return addr.network
}

func (addr *MockAddr) String() string {
	return addr.address
}
