package configdb

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/metal-stack/metal-lib/pkg/testcommon"
)

func Test_parseBreakout(t *testing.T) {
	tests := []struct {
		name          string
		portName      string
		breakoutMode  string
		wantNumber    int
		wantSpeed     int
		wantPortIndex int
		wantErr       error
	}{
		{
			name:         "port name cannot be empty",
			breakoutMode: "4x24G",
			wantErr:      fmt.Errorf("port name must not be empty"),
		},
		{
			name:         "port name must start with 'Ethernet'",
			portName:     "eth0",
			breakoutMode: "4x24G",
			wantErr:      fmt.Errorf("port name must start with 'Ethernet' and end with a positive integer"),
		},
		{
			name:         "port name suffix must be an integer",
			portName:     "EthernetX",
			breakoutMode: "4x24G",
			wantErr:      fmt.Errorf("port name must start with 'Ethernet' and end with a positive integer"),
		},
		{
			name:         "port name suffix must be a positive integer",
			portName:     "Ethernet-10",
			breakoutMode: "4x24G",
			wantErr:      fmt.Errorf("port name must start with 'Ethernet' and end with a positive integer"),
		},
		{
			name:         "port number must be divisible by 4",
			portName:     "Ethernet1",
			breakoutMode: "4x24G",
			wantErr:      fmt.Errorf("port number must be divisible by 4"),
		},
		{
			name:         "breakout mode must be valid",
			portName:     "Ethernet0",
			breakoutMode: "2x25G",
			wantErr:      fmt.Errorf("breakout mode must be one of '1x100G[40G]', '2x50G', '4x25G', '4x10G'"),
		},
		{
			name:          "breakout 1x100G[40G]",
			portName:      "Ethernet4",
			breakoutMode:  "1x100G[40G]",
			wantNumber:    1,
			wantSpeed:     100000,
			wantPortIndex: 2,
			wantErr:       nil,
		},
		{
			name:          "breakout 2x50G",
			portName:      "Ethernet8",
			breakoutMode:  "2x50G",
			wantNumber:    2,
			wantSpeed:     50000,
			wantPortIndex: 3,
			wantErr:       nil,
		},
		{
			name:          "breakout 4x25G",
			portName:      "Ethernet12",
			breakoutMode:  "4x25G",
			wantNumber:    4,
			wantSpeed:     25000,
			wantPortIndex: 4,
			wantErr:       nil,
		},
		{
			name:          "breakout 4x10G",
			portName:      "Ethernet16",
			breakoutMode:  "4x10G",
			wantNumber:    4,
			wantSpeed:     10000,
			wantPortIndex: 5,
			wantErr:       nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNumber, gotSpeed, gotPortIndex, err := parseBreakout(tt.portName, tt.breakoutMode)
			if diff := cmp.Diff(err, tt.wantErr, testcommon.ErrorStringComparer()); diff != "" {
				t.Errorf("parseBreakout() error diff = %s", diff)
			}

			if gotNumber != tt.wantNumber {
				t.Errorf("parseBreakout() gotNumber = %v, want %v", gotNumber, tt.wantNumber)
			}
			if gotSpeed != tt.wantSpeed {
				t.Errorf("parseBreakout() gotSpeed = %v, want %v", gotSpeed, tt.wantSpeed)
			}
			if gotPortIndex != tt.wantPortIndex {
				t.Errorf("parseBreakout() gotPortIndex = %v, want %v", gotPortIndex, tt.wantPortIndex)
			}
		})
	}
}

func Test_getPortsFromBreakout(t *testing.T) {
	tests := []struct {
		name         string
		portName     string
		breakoutMode string
		want         []Port
		wantErr      bool
	}{
		{
			name:         "add one 1x100G[40G] port",
			portName:     "Ethernet120",
			breakoutMode: "1x100G[40G]",
			want: []Port{
				{
					AdminStatus: DefaultAdminStatus,
					Alias:       "Eth31(Port31)",
					Autoneg:     DefaultAutonegMode,
					FEC:         DefaultFECMode,
					Index:       31,
					Lanes:       "121,122,123,124",
					MTU:         DefaultMTU,
					ParentPort:  "Ethernet120",
					Speed:       100000,
				},
			},
			wantErr: false,
		},
		{
			name:         "add two 2x50G ports",
			portName:     "Ethernet116",
			breakoutMode: "2x50G",
			want: []Port{
				{
					AdminStatus: DefaultAdminStatus,
					Alias:       "Eth30/1(Port30)",
					Autoneg:     DefaultAutonegMode,
					FEC:         DefaultFECMode,
					Index:       30,
					Lanes:       "117,118",
					MTU:         DefaultMTU,
					ParentPort:  "Ethernet116",
					Speed:       50000,
				},
				{
					AdminStatus: DefaultAdminStatus,
					Alias:       "Eth30/2(Port30)",
					Autoneg:     DefaultAutonegMode,
					FEC:         DefaultFECMode,
					Index:       30,
					Lanes:       "119,120",
					MTU:         DefaultMTU,
					ParentPort:  "Ethernet116",
					Speed:       50000,
				},
			},
			wantErr: false,
		},
		{
			name:         "add four 4x10G ports",
			portName:     "Ethernet8",
			breakoutMode: "4x10G",
			want: []Port{
				{
					AdminStatus: DefaultAdminStatus,
					Alias:       "Eth3/1(Port3)",
					Autoneg:     DefaultAutonegMode,
					FEC:         DefaultFECMode,
					Index:       3,
					Lanes:       "9",
					MTU:         DefaultMTU,
					ParentPort:  "Ethernet8",
					Speed:       10000,
				},
				{
					AdminStatus: DefaultAdminStatus,
					Alias:       "Eth3/2(Port3)",
					Autoneg:     DefaultAutonegMode,
					FEC:         DefaultFECMode,
					Index:       3,
					Lanes:       "10",
					MTU:         DefaultMTU,
					ParentPort:  "Ethernet8",
					Speed:       10000,
				},
				{
					AdminStatus: DefaultAdminStatus,
					Alias:       "Eth3/3(Port3)",
					Autoneg:     DefaultAutonegMode,
					FEC:         DefaultFECMode,
					Index:       3,
					Lanes:       "11",
					MTU:         DefaultMTU,
					ParentPort:  "Ethernet8",
					Speed:       10000,
				},
				{
					AdminStatus: DefaultAdminStatus,
					Alias:       "Eth3/4(Port3)",
					Autoneg:     DefaultAutonegMode,
					FEC:         DefaultFECMode,
					Index:       3,
					Lanes:       "12",
					MTU:         DefaultMTU,
					ParentPort:  "Ethernet8",
					Speed:       10000,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getPortsFromBreakout(tt.portName, tt.breakoutMode)
			if (err != nil) != tt.wantErr {
				t.Errorf("getPortsFromBreakout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("getPortsFromBreakout() diff = %s", diff)
			}
		})
	}
}

func Test_getLanesForPort(t *testing.T) {
	tests := []struct {
		name      string
		portIndex int
		number    int
		offset    int
		want      string
	}{
		{
			name:      "1x100G[40G]",
			portIndex: 31,
			number:    1,
			want:      "121,122,123,124",
		},
		{
			name:      "2x50G first",
			portIndex: 30,
			number:    2,
			want:      "117,118",
		},
		{
			name:      "2x50G second",
			portIndex: 30,
			number:    2,
			offset:    1,
			want:      "119,120",
		},
		{
			name:      "4x25G first",
			portIndex: 2,
			number:    4,
			want:      "5",
		},
		{
			name:      "4x25G second",
			portIndex: 2,
			number:    4,
			offset:    1,
			want:      "6",
		},
		{
			name:      "4x25G third",
			portIndex: 2,
			number:    4,
			offset:    2,
			want:      "7",
		},
		{
			name:      "4x25G fourth",
			portIndex: 2,
			number:    4,
			offset:    3,
			want:      "8",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getLanesForPort(tt.portIndex, tt.number, tt.offset); got != tt.want {
				t.Errorf("getLanesForPort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getPortAlias(t *testing.T) {
	tests := []struct {
		name      string
		portIndex int
		number    int
		offset    int
		want      string
	}{
		{
			name:      "1x100G[40G]",
			portIndex: 9,
			number:    1,
			offset:    3, // offset should not matter
			want:      "Eth9(Port9)",
		},
		{
			name:      "2x50G first",
			portIndex: 9,
			number:    2,
			offset:    0,
			want:      "Eth9/1(Port9)",
		},
		{
			name:      "2x50G second",
			portIndex: 9,
			number:    2,
			offset:    1,
			want:      "Eth9/2(Port9)",
		},
		{
			name:      "4x25G third",
			portIndex: 9,
			number:    4,
			offset:    2,
			want:      "Eth9/3(Port9)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPortAlias(tt.portIndex, tt.number, tt.offset); got != tt.want {
				t.Errorf("getPortAlias() = %v, want %v", got, tt.want)
			}
		})
	}
}
