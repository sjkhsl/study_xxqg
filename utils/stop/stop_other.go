//go:build !windows
// +build !windows

package stop

func Stop() {
	select {}
}
