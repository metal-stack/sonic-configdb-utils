package platform

import (
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"sort"
	"strconv"
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

type PortAliases []string

type SFP struct {
	Name string `json:"name"`
}

type SpeedOptions [2]int

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
	parseError := fmt.Errorf("error parsing breakout mode; must be of form <number>x<speed>G or <number>x<speed>G[<alternative-speed>G]")
	numString, suffix, ok := strings.Cut(breakoutMode, "x")
	if !ok {
		return SpeedOptions{}, parseError
	}

	num, err := strconv.Atoi(numString)
	if err != nil || num < 0 {
		return SpeedOptions{}, parseError
	}

	speedString, suffix, ok := strings.Cut(suffix, "G")
	if !ok {
		return SpeedOptions{}, parseError
	}

	speed, err := strconv.Atoi(speedString)
	if err != nil || speed < 0 {
		return SpeedOptions{}, parseError
	}

	if breakoutMode == fmt.Sprintf("%dx%dG", num, speed) {
		return SpeedOptions{speed * 1000, 0}, nil
	}

	_, suffix, ok = strings.Cut(suffix, "[")
	if !ok {
		return SpeedOptions{}, parseError
	}

	altSpeedString, _, ok := strings.Cut(suffix, "G]")
	if !ok {
		return SpeedOptions{}, parseError
	}

	altSpeed, err := strconv.Atoi(altSpeedString)
	if err != nil || altSpeed < 0 {
		return SpeedOptions{}, parseError
	}

	return SpeedOptions{speed * 1000, altSpeed * 1000}, nil
}

func UnmarshalPlatformJSON(in []byte) (*Platform, error) {
	var platform Platform
	err := json.Unmarshal(in, &platform)
	if err != nil {
		return nil, err
	}

	return &platform, nil
}

func stringToIntSlice(input string) ([]int, error) {
	ints := make([]int, 0)

	for n := range strings.SplitSeq(input, ",") {
		number, err := strconv.Atoi(n)
		if err != nil {
			return []int{}, err
		}
		ints = append(ints, number)
	}

	return ints, nil
}
