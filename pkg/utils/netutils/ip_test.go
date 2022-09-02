package netutils

import "testing"

func TestAllocate(t *testing.T) {
	ipstr, err := Allocate("188.8.5.0/24", []string{"188.8.5.254"})
	t.Log(ipstr, err)
}
