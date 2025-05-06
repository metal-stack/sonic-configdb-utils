package configdb

import (
	"fmt"
	"strconv"
	"strings"

	p "github.com/metal-stack/sonic-configdb-utils/platform"
)

const (
	defaultAdminStatus = AdminStatusUp
	defaultAutonegMode = AutonegModeOff
	defaultMTU         = 9000
)

type portDefaults struct {
	autoneg AutonegMode
	fec     FECMode
	mtu     int
}

func getPortsFromBreakout(portName, breakoutMode string, defaults portDefaults, platform *p.Platform) (map[string]Port, error) {
	ports := make(map[string]Port)

	breakoutPorts, err := platform.ParseBreakout(portName, breakoutMode)
	if err != nil {
		return nil, err
	}

	speedOptions, err := p.ParseSpeedOptions(breakoutMode)
	if err != nil {
		return nil, err
	}

	for i, alias := range breakoutPorts.PortAliases {
		numAliases := len(breakoutPorts.PortAliases)
		numLanes := len(breakoutPorts.Lanes)
		if numLanes < 1 {
			return nil, fmt.Errorf("no lanes given for port %s", portName)
		}
		lanesPerPort := numLanes / numAliases
		begin := i * lanesPerPort
		end := begin + lanesPerPort
		lanes := breakoutPorts.Lanes[begin:end]

		lanesString := fmt.Sprintf("%d", lanes[0])
		for i, lane := range lanes {
			if i == 0 {
				continue
			}
			lanesString += fmt.Sprintf(",%d", lane)
		}

		port := Port{
			AdminStatus:    defaultAdminStatus,
			Alias:          alias,
			Autoneg:        defaultAutonegMode,
			Index:          fmt.Sprintf("%d", breakoutPorts.Index[i]),
			Lanes:          lanesString,
			MTU:            fmt.Sprintf("%d", defaultMTU),
			parentBreakout: breakoutMode,
			Speed:          fmt.Sprintf("%d", speedOptions[0]),
		}

		if defaults.fec != "" {
			port.FEC = defaults.fec
		}
		if defaults.mtu != 0 {
			port.MTU = fmt.Sprintf("%d", defaults.mtu)
		}
		if defaults.autoneg != "" {
			port.Autoneg = defaults.autoneg
		}

		name, err := incrementPortNameSuffix(portName, i*lanesPerPort)
		if err != nil {
			return nil, err
		}
		ports[name] = port
	}

	return ports, nil
}

func incrementPortNameSuffix(portName string, increment int) (string, error) {
	parseError := fmt.Errorf("invalid port name %s; must be of form EthernetX, where X is a positive number", portName)

	suffix, ok := strings.CutPrefix(portName, "Ethernet")
	if !ok {
		return "", parseError
	}

	number, err := strconv.Atoi(suffix)
	if err != nil || number < 0 {
		return "", parseError
	}

	return fmt.Sprintf("Ethernet%d", number+increment), nil
}
