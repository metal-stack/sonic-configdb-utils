package configdb

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/metal-stack/metal-lib/pkg/testcommon"
	p "github.com/metal-stack/sonic-configdb-utils/platform"
	"github.com/metal-stack/sonic-configdb-utils/values"
)

func Test_incrementPortNameSuffix(t *testing.T) {
	tests := []struct {
		name      string
		portName  string
		increment int
		want      string
		wantErr   bool
	}{
		{
			name:      "invalid prefix",
			portName:  "eth1",
			increment: 0,
			want:      "",
			wantErr:   true,
		},
		{
			name:      "increment suffix",
			portName:  "Ethernet56",
			increment: 1,
			want:      "Ethernet57",
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := incrementPortNameSuffix(tt.portName, tt.increment)
			if (err != nil) != tt.wantErr {
				t.Errorf("incrementPortNameSuffix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("incrementPortNameSuffix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getPortsFromBreakout(t *testing.T) {
	tests := []struct {
		name               string
		portName           string
		breakoutMode       string
		defaultPortFECMode values.FECMode
		defaultMTU         int
		platform           *p.Platform
		want               map[string]Port
		wantErr            bool
	}{
		{
			name:         "add port with only one lane",
			portName:     "Ethernet1",
			breakoutMode: "1x1G",
			platform: &p.Platform{
				Interfaces: map[string]p.Interface{
					"Ethernet1": {
						BreakoutModes: map[string][]string{
							"1x1G": {"Eth2(Port2)"},
						},
						Index: "2",
						Lanes: "25",
					},
				},
			},
			want: map[string]Port{
				"Ethernet1": {
					AdminStatus:    defaultAdminStatus,
					Alias:          "Eth2(Port2)",
					Autoneg:        defaultAutonegMode,
					Index:          "2",
					Lanes:          "25",
					MTU:            fmt.Sprintf("%d", defaultMTU),
					parentBreakout: "1x1G",
					Speed:          "1000",
				},
			},
			wantErr: false,
		},
		{
			name:               "add one 1x100G[40G] port with different defaults",
			portName:           "Ethernet120",
			breakoutMode:       "1x100G[40G]",
			defaultPortFECMode: values.FECModeRS,
			defaultMTU:         1500,
			platform: &p.Platform{
				Interfaces: map[string]p.Interface{
					"Ethernet120": {
						BreakoutModes: map[string][]string{
							"1x100G[40G]": {"Eth31(Port31)"},
						},
						Index: "31,31,31,31",
						Lanes: "121,122,123,124",
					},
				},
			},
			want: map[string]Port{
				"Ethernet120": {
					AdminStatus:    defaultAdminStatus,
					Alias:          "Eth31(Port31)",
					Autoneg:        defaultAutonegMode,
					FEC:            FECModeRS,
					Index:          "31",
					Lanes:          "121,122,123,124",
					MTU:            "1500",
					parentBreakout: "1x100G[40G]",
					Speed:          "100000",
				},
			},
			wantErr: false,
		},
		{
			name:         "add two 2x50G ports",
			portName:     "Ethernet116",
			breakoutMode: "2x50G",
			platform: &p.Platform{
				Interfaces: map[string]p.Interface{
					"Ethernet116": {
						BreakoutModes: map[string][]string{
							"2x50G": {"Eth30/1(Port30)", "Eth30/2(Port30)"},
						},
						Index: "30,30,30,30",
						Lanes: "117,118,119,120",
					},
				},
			},
			want: map[string]Port{
				"Ethernet116": {
					AdminStatus:    defaultAdminStatus,
					Alias:          "Eth30/1(Port30)",
					Autoneg:        defaultAutonegMode,
					Index:          "30",
					Lanes:          "117,118",
					MTU:            fmt.Sprintf("%d", defaultMTU),
					parentBreakout: "2x50G",
					Speed:          "50000",
				},
				"Ethernet118": {
					AdminStatus:    defaultAdminStatus,
					Alias:          "Eth30/2(Port30)",
					Autoneg:        defaultAutonegMode,
					Index:          "30",
					Lanes:          "119,120",
					MTU:            fmt.Sprintf("%d", defaultMTU),
					parentBreakout: "2x50G",
					Speed:          "50000",
				},
			},
			wantErr: false,
		},
		{
			name:         "add four 4x10G ports",
			portName:     "Ethernet8",
			breakoutMode: "4x10G",
			platform: &p.Platform{
				Interfaces: map[string]p.Interface{
					"Ethernet8": {
						BreakoutModes: map[string][]string{
							"4x10G": {"Eth3/1(Port3)", "Eth3/2(Port3)", "Eth3/3(Port3)", "Eth3/4(Port3)"},
						},
						Index: "3,3,3,3",
						Lanes: "9,10,11,12",
					},
				},
			},
			want: map[string]Port{
				"Ethernet8": {
					AdminStatus:    defaultAdminStatus,
					Alias:          "Eth3/1(Port3)",
					Autoneg:        defaultAutonegMode,
					Index:          "3",
					Lanes:          "9",
					MTU:            fmt.Sprintf("%d", defaultMTU),
					parentBreakout: "4x10G",
					Speed:          "10000",
				},
				"Ethernet9": {
					AdminStatus:    defaultAdminStatus,
					Alias:          "Eth3/2(Port3)",
					Autoneg:        defaultAutonegMode,
					Index:          "3",
					Lanes:          "10",
					MTU:            fmt.Sprintf("%d", defaultMTU),
					parentBreakout: "4x10G",
					Speed:          "10000",
				},
				"Ethernet10": {
					AdminStatus:    defaultAdminStatus,
					Alias:          "Eth3/3(Port3)",
					Autoneg:        defaultAutonegMode,
					Index:          "3",
					Lanes:          "11",
					MTU:            fmt.Sprintf("%d", defaultMTU),
					parentBreakout: "4x10G",
					Speed:          "10000",
				},
				"Ethernet11": {
					AdminStatus:    defaultAdminStatus,
					Alias:          "Eth3/4(Port3)",
					Autoneg:        defaultAutonegMode,
					Index:          "3",
					Lanes:          "12",
					MTU:            fmt.Sprintf("%d", defaultMTU),
					parentBreakout: "4x10G",
					Speed:          "10000",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getPortsFromBreakout(tt.portName, tt.breakoutMode, tt.defaultPortFECMode, tt.defaultMTU, tt.platform)
			if (err != nil) != tt.wantErr {
				t.Errorf("getPortsFromBreakout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got, testcommon.IgnoreUnexported()); diff != "" {
				t.Errorf("getPortsFromBreakout() diff = %s", diff)
			}
		})
	}
}
