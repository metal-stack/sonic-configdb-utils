package configdb

import (
	"fmt"
	"strconv"
	"strings"
)

func parseBreakout(portName string, breakoutMode string) (number, speed int, err error) {
	if portName == "" {
		return 0, 0, fmt.Errorf("port name must not be empty")
	}

	invalidPortName := fmt.Errorf("port name must start with 'Ethernet' and end with a positive integer")

	prefix, suffix, has := strings.Cut(portName, "Ethernet")
	if !has || prefix != "" {
		return 0, 0, invalidPortName
	}
	portNumber, err := strconv.Atoi(suffix)
	if err != nil || portNumber < 0 {
		return 0, 0, invalidPortName
	}
	if portNumber%4 != 0 {
		return 0, 0, fmt.Errorf("port number must be divisible by 4")
	}

	switch breakoutMode {
	case "1x100G[40G]":
		return 1, 100, nil
	case "2x50G":
		return 2, 50, nil
	case "4x25G":
		return 4, 25, nil
	case "4x10G":
		return 4, 10, nil
	default:
		return 0, 0, fmt.Errorf("breakout mode must be on of '1x100G[40G]', '2x50G', '4x25G', '4x10G'")
	}
}
