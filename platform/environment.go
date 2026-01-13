package platform

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Environment struct {
	HWSKU    string
	MAC      string
	Platform string
}

func GetEnvironment(envFile string) (*Environment, error) {
	var (
		hwsku    string
		mac      string
		platform string
	)

	file, err := os.Open(envFile)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Printf("failed to close file:%v", err)
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Trim(line, "\n\t ")
		if p, ok := strings.CutPrefix(line, "PLATFORM="); ok {
			platform = p
		}
		if h, ok := strings.CutPrefix(line, "HWSKU="); ok {
			hwsku = h
		}
	}

	addrBytes, err := os.ReadFile("/sys/class/net/eth0/address")
	if err != nil {
		return nil, fmt.Errorf("failed to get mac address:%w", err)
	}
	mac = strings.TrimSpace(string(addrBytes))

	return &Environment{
		HWSKU:    hwsku,
		MAC:      mac,
		Platform: platform,
	}, nil
}
