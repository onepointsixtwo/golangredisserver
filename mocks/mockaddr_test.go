package mocks

import (
	"net"
	"testing"
)

func TestMockAddrCanBeUsedAsNetAddr(t *testing.T) {
	var addr net.Addr
	addr = NewMockAddr("tcp", "0.0.0.0")

	if addr.Network() != "tcp" || addr.String() != "0.0.0.0" {
		t.Error("MockAddr should be useable as a net.Addr interface")
	}
}
