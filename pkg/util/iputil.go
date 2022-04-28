package util

import (
	"context"
	"net"
	"time"

	"github.com/testground/sdk-go/network"
)

// Configures the instance's network network, using the passed netClient and routingPolicy
// This can affect behavior in certain tests (e.g. denying all routing will effectively block the instance from the network)
func ConfigureNetwork(netClient *network.Client, routingPolicy network.RoutingPolicyType, ctx context.Context) {
	config := &network.Config{
		// Control the "default" network. At the moment, this is the only network.
		Network: "default",

		// Enable this network. Setting this to false will disconnect this test
		// instance from this network. You probably don't want to do that.
		Enable: true,
		Default: network.LinkShape{
			Latency:   100 * time.Millisecond,
			Bandwidth: 1 << 20, // 1Mib
		},
		CallbackState: "network-configured",
		RoutingPolicy: routingPolicy,
	}

	netClient.MustConfigureNetwork(ctx, config)
}

// Returns true if the two given addresses are the same
func SameAddrs(a, b []net.Addr) bool {
	if len(a) != len(b) {
		return false
	}
	aset := make(map[string]bool, len(a))
	for _, addr := range a {
		aset[addr.String()] = true
	}
	for _, addr := range b {
		if !aset[addr.String()] {
			return false
		}
	}
	return true
}
