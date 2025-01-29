package configdb

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/metal-stack/sonic-configdb-utils/values"
)

const (
	defaultAdminStatus = AdminStatusUp
	defaultAutonegMode = AutonegModeOff
	defaultFECMode     = FECModeNone
	defaultMTU         = 9000
)

func getPortAlias(portIndex, number, offset int) string {
	if number <= 1 {
		return fmt.Sprintf("Eth%d(Port%d)", portIndex, portIndex)
	}
	return fmt.Sprintf("Eth%d/%d(Port%d)", portIndex, offset+1, portIndex)
}

func getLanesForPort(portIndex, number, offset int) string {
	firstLaneIndex := (portIndex-1)*4 + 1

	switch number {
	case 1:
		return fmt.Sprintf("%d,%d,%d,%d", firstLaneIndex, firstLaneIndex+1, firstLaneIndex+2, firstLaneIndex+3)
	case 2:
		return fmt.Sprintf("%d,%d", firstLaneIndex+2*offset, firstLaneIndex+2*offset+1)
	case 4:
		return fmt.Sprintf("%d", firstLaneIndex+offset)
	default:
		return ""
	}
}

func getPortsFromBreakout(portName, breakoutMode string, defaultPortFECMode values.FECMode, defaultPortMTU int) (map[string]Port, error) {
	number, speed, portIndex, err := parseBreakout(portName, breakoutMode)
	if err != nil {
		return nil, err
	}

	ports := make(map[string]Port)

	for i := 0; i < number; i++ {
		port := Port{
			AdminStatus: defaultAdminStatus,
			Alias:       getPortAlias(portIndex, number, i),
			Autoneg:     defaultAutonegMode,
			FEC:         defaultFECMode,
			Index:       fmt.Sprintf("%d", portIndex),
			Lanes:       getLanesForPort(portIndex, number, i),
			MTU:         fmt.Sprintf("%d", defaultMTU),
			Speed:       fmt.Sprintf("%d", speed),
		}

		if defaultPortFECMode != "" {
			port.FEC = FECMode(defaultPortFECMode)
		}
		if defaultPortMTU != 0 {
			port.MTU = fmt.Sprintf("%d", defaultPortMTU)
		}

		nameSuffix := (portIndex-1)*4 + 4/number*i // works because i > 0 for number = 1 never occurs
		name := fmt.Sprintf("Ethernet%d", nameSuffix)
		ports[name] = port
	}

	return ports, nil
}

func parseBreakout(portName, breakoutMode string) (number, speed, portIndex int, err error) {
	if portName == "" {
		return 0, 0, 0, fmt.Errorf("port name must not be empty")
	}

	invalidPortName := fmt.Errorf("port name must start with 'Ethernet' and end with a positive integer")

	prefix, suffix, has := strings.Cut(portName, "Ethernet")
	if !has || prefix != "" {
		return 0, 0, 0, invalidPortName
	}
	portNumber, err := strconv.Atoi(suffix)
	if err != nil || portNumber < 0 {
		return 0, 0, 0, invalidPortName
	}
	if portNumber%4 != 0 {
		return 0, 0, 0, fmt.Errorf("port number must be divisible by 4")
	}

	switch breakoutMode {
	case "1x100G[40G]":
		return 1, 100 * 1000, portNumber/4 + 1, nil
	case "2x50G":
		return 2, 50 * 1000, portNumber/4 + 1, nil
	case "4x25G":
		return 4, 25 * 1000, portNumber/4 + 1, nil
	case "4x10G":
		return 4, 10 * 1000, portNumber/4 + 1, nil
	default:
		return 0, 0, 0, fmt.Errorf("breakout mode must be one of '1x100G[40G]', '2x50G', '4x25G', '4x10G'")
	}
}
