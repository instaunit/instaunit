package net

import (
	"net"
)

/**
 * Rebase an address for a host
 */
func RebaseAddr(host, addr string) (string, error) {
	_, p, err := net.SplitHostPort(addr)
	if err != nil {
		return "", err
	}
	return host + ":" + p, nil
}
