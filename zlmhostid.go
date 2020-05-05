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
	HostID []string `json:"hostid"`
	OS     string   `json:"os"`
	CPU    string   `json:"cpu"`
}

// GetOS returns the OS.
func GetOS() string {
	m := map[string]string{
		"android": "Android",
		"freebsd": "FreeBSD",
		"linux":   "Linux",
		// iOS is not supported by Go
		"darwin":  "Mac OS",
		"openbsd": "OpenBSD",
		"windows": "Windows",
	}
	return m[runtime.GOOS]
}

// GetCPU returns the CPU architecture.
func GetCPU() string {
	m := map[string]string{
		"amd64":   "x86_64",
		"386":     "x86",
		"arm64":   "aarch64",
		"arm":     "arm",
		"ppc64le": "powerpc64le",
		"ppc64":   "powerpc64",
		// powerpc32le is not supported by Go
		"ppc":      "powerpc32",
		"mipsle":   "mips32el",
		"mips":     "mips32",
		"mips64le": "mips64el",
		"mips64":   "mips64",
	}
	return m[runtime.GOARCH]
}

// Get returns the host ID.
func Get() (*HostID, error) {
	var id HostID
	var addrs []string
	ifs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, v := range ifs {
		h := v.HardwareAddr.String()
		if len(h) > 0 {
			addrs = append(addrs, h)
		}
	}
	sort.Strings(addrs) // sort host IDs
	if len(addrs) > 0 { // make host IDs unique
		id.HostID = append(id.HostID, addrs[0])
		last := addrs[0]
		for i := 1; i < len(addrs); i++ {
			if addrs[i] != last {
				id.HostID = append(id.HostID, addrs[i])
				last = addrs[i]
			}
		}
	}
	id.OS = GetOS()
	id.CPU = GetCPU()
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
