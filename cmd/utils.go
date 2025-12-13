package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"sync"
)

var (
	ipRegex    = regexp.MustCompile(`\d+$`)
	dnsCache   sync.Map
	bufferPool = sync.Pool{
		New: func() interface{} {
			return make([]byte, 4096)
		},
	}
)

func ResolveIP(ctx context.Context, host string) (string, error) {
	// If it's already an IP, return it
	if ipRegex.MatchString(host) {
		return host, nil
	}

	// Check cache first
	if val, ok := dnsCache.Load(host); ok {
		return val.(string), nil
	}

	// Lookup
	ips, err := net.DefaultResolver.LookupIP(ctx, "ip4", host)
	if err != nil {
		return "", err
	}
	if len(ips) == 0 {
		return "", fmt.Errorf("no IP found for host: %s", host)
	}

	ipStr := ips[0].String()
	dnsCache.Store(host, ipStr)
	return ipStr, nil
}

func ReadFile(filename string) ([]string, error) {
	var reader io.Reader

	if filename == "" || filename == "-" {
		stat, err := os.Stdin.Stat()
		if err != nil {
			return nil, fmt.Errorf("error checking stdin: %w", err)
		}

		if (stat.Mode() & os.ModeCharDevice) == 0 {
			reader = os.Stdin
		} else {
			return nil, fmt.Errorf("no input provided: use -f flag or pipe data via stdin")
		}
	} else {
		file, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		reader = file
	}

	var lines []string
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func ipInc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func IPsFromCIDR(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for currentIP := ip.Mask(ipnet.Mask); ipnet.Contains(currentIP); ipInc(currentIP) {
		ips = append(ips, currentIP.String())
	}
	if len(ips) <= 1 {
		return ips, nil
	}

	return ips[1 : len(ips)-1], nil
}

func fatal(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}
