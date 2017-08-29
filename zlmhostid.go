// zlmhostid is a minimal host ID tool for ZLM2.
package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"runtime"
	"sort"
)

// HostID defines a host ID.
type HostID struct {
	HostID []string
	OS     string
	Arch   string
}

// Get returns the host ID.
func Get() (*HostID, error) {
	var id HostID
	ifs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, v := range ifs {
		h := v.HardwareAddr.String()
		if len(h) > 0 {
			id.HostID = append(id.HostID, h)
		}
	}
	sort.Strings(id.HostID) // sort host IDs
	id.OS = runtime.GOOS
	id.Arch = runtime.GOARCH
	return &id, nil
}

// MarshalIndent the host ID as a JSON string (indented).
func (id *HostID) MarshalIndent() (string, error) {
	buf, err := json.MarshalIndent(id, "", "  ")
	if err != nil {
		return "", err
	}
	return string(buf), nil
}
func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err)
	os.Exit(1)
}

func main() {
	hostID, err := Get()
	if err != nil {
		fatal(err)
	}
	jsn, err := hostID.MarshalIndent()
	if err != nil {
		fatal(err)
	}
	fmt.Println(jsn)
}
