package configdb

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/metal-stack/metal-lib/pkg/testcommon"
)

func Test_parseBreakout(t *testing.T) {
	tests := []struct {
		name         string
		portName     string
		breakoutMode string
		wantNumber   int
		wantSpeed    int
		wantErr      error
	}{
		{
			name:         "port name cannot be empty",
			portName:     "",
			breakoutMode: "4x24G",
			wantNumber:   0,
			wantSpeed:    0,
			wantErr:      fmt.Errorf("port name must not be empty"),
		},
		{
			name:         "port name must start with 'Ethernet'",
			portName:     "eth0",
			breakoutMode: "4x24G",
			wantNumber:   0,
			wantSpeed:    0,
			wantErr:      fmt.Errorf("port name must start with 'Ethernet' and end with a positive integer"),
		},
		{
			name:         "port name suffix must be an integer",
			portName:     "EthernetX",
			breakoutMode: "4x24G",
			wantNumber:   0,
			wantSpeed:    0,
			wantErr:      fmt.Errorf("port name must start with 'Ethernet' and end with a positive integer"),
		},
		{
			name:         "port name suffix must be a positive integer",
			portName:     "Ethernet-10",
			breakoutMode: "4x24G",
			wantNumber:   0,
			wantSpeed:    0,
			wantErr:      fmt.Errorf("port name must start with 'Ethernet' and end with a positive integer"),
		},
		{
			name:         "port number must be divisible by 4",
			portName:     "Ethernet1",
			breakoutMode: "4x24G",
			wantNumber:   0,
			wantSpeed:    0,
			wantErr:      fmt.Errorf("port number must be divisible by 4"),
		},
		{
			name:         "breakout mode must be valid",
			portName:     "Ethernet0",
			breakoutMode: "2x25G",
			wantNumber:   0,
			wantSpeed:    0,
			wantErr:      fmt.Errorf("breakout mode must be on of '1x100G[40G]', '2x50G', '4x25G', '4x10G'"),
		},
		{
			name:         "breakout 1x100G[40G]",
			portName:     "Ethernet4",
			breakoutMode: "1x100G[40G]",
			wantNumber:   1,
			wantSpeed:    100,
			wantErr:      nil,
		},
		{
			name:         "breakout 2x50G",
			portName:     "Ethernet8",
			breakoutMode: "2x50G",
			wantNumber:   2,
			wantSpeed:    50,
			wantErr:      nil,
		},
		{
			name:         "breakout 4x25G",
			portName:     "Ethernet12",
			breakoutMode: "4x25G",
			wantNumber:   4,
			wantSpeed:    25,
			wantErr:      nil,
		},
		{
			name:         "breakout 4x10G",
			portName:     "Ethernet16",
			breakoutMode: "4x10G",
			wantNumber:   4,
			wantSpeed:    10,
			wantErr:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNumber, gotSpeed, err := parseBreakout(tt.portName, tt.breakoutMode)
			if diff := cmp.Diff(err, tt.wantErr, testcommon.ErrorStringComparer()); diff != "" {
				t.Errorf("parseBreakout() error diff: %s", diff)
			}

			if gotNumber != tt.wantNumber {
				t.Errorf("parseBreakout() gotNumber = %v, want %v", gotNumber, tt.wantNumber)
			}
			if gotSpeed != tt.wantSpeed {
				t.Errorf("parseBreakout() gotSpeed = %v, want %v", gotSpeed, tt.wantSpeed)
			}
		})
	}
}
