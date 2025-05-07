package platform

import (
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"sort"
	"strings"
)

type BreakoutConfig map[string]string

type BreakoutPorts struct {
	PortAliases []string
	Lanes       []int
	Index       []int
}

type Interface struct {
	BreakoutModes map[string][]string `json:"breakout_modes"`
	Index         string              `json:"index"`
	Lanes         string              `json:"lanes"`
}

type Platform struct {
	Interfaces map[string]Interface `json:"interfaces"`
}

type SpeedOptions []int

func (p *Platform) GetDefaultBreakoutConfig() BreakoutConfig {
	breakoutConfig := make(BreakoutConfig)

	for name, intf := range p.Interfaces {
		if len(intf.BreakoutModes) > 0 {
			keys := slices.Collect(maps.Keys(intf.BreakoutModes))
			sort.Strings(keys)
			breakoutConfig[name] = keys[0]
		}
	}

	return breakoutConfig
}

func (p *Platform) ParseBreakout(portName, breakout string) (*BreakoutPorts, error) {
	intf, ok := p.Interfaces[portName]
	if !ok {
		return nil, fmt.Errorf("unknown port %s", portName)
	}

	aliases, ok := intf.BreakoutModes[breakout]
	if !ok {
		return nil, fmt.Errorf("invalid breakout mode %s; must be one of %v", breakout, slices.Collect(maps.Keys(intf.BreakoutModes)))
	}

	lanes, err := stringToIntSlice(intf.Lanes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse lanes, %w", err)
	}

	index, err := stringToIntSlice(intf.Index)
	if err != nil {
		return nil, fmt.Errorf("unable to parse index, %w", err)
	}

	return &BreakoutPorts{
		PortAliases: aliases,
		Lanes:       lanes,
		Index:       index,
	}, nil
}

func ParseSpeedOptions(breakoutMode string) (SpeedOptions, error) {
	options := make(SpeedOptions, 0)
	parseError := fmt.Errorf("invalid breakout mode %s, must be of form <number>x<speed>G[<altSpeed1>G,...,<altSpeedN>G]", breakoutMode)

	altSpeedString, ok := cutBetween(breakoutMode, "[", "]")
	if !ok {
		return parseSingleSpeedOption(breakoutMode)
	}

	speedString, ok := cutBetween(breakoutMode, "x", "[")
	if !ok {
		return nil, parseError
	}

	speed, err := stringToSpeed(speedString)
	if err != nil {
		return nil, err
	}
	options = append(options, speed)

	for optionString := range strings.SplitSeq(altSpeedString, ",") {
		option, err := stringToSpeed(optionString)
		if err != nil {
			return nil, err
		}
		options = append(options, option)
	}

	return options, nil
}

func UnmarshalPlatformJSON(in []byte) (*Platform, error) {
	var platform Platform
	err := json.Unmarshal(in, &platform)
	if err != nil {
		return nil, err
	}

	return &platform, nil
}

func parseSingleSpeedOption(breakoutMode string) (SpeedOptions, error) {
	_, speedString, ok := strings.Cut(breakoutMode, "x")
	if !ok {
		return nil, fmt.Errorf("invalid breakout mode %s, must be of form <number>x<speed>G", breakoutMode)
	}

	speed, err := stringToSpeed(speedString)
	if err != nil {
		return nil, err
	}

	return SpeedOptions{speed}, nil
}
