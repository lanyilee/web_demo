package netutils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	C_MAC_QEMU_PREFIX = "52:54:00"
)

func RandomMac() (mac string) {

	mac = C_MAC_QEMU_PREFIX

	seed := rand.NewSource(time.Now().UnixNano())
	rander := rand.New(seed)

	mac = fmt.Sprintf("%s:%02x:%02x:%02x",
		mac,
		rander.Intn(255),
		rander.Intn(255),
		rander.Intn(255))
	return
}

func IsVirnetMac(mac string) bool {
	return strings.HasPrefix(mac, C_MAC_QEMU_PREFIX)
}
