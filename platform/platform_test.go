package platform

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/metal-stack/metal-lib/pkg/testcommon"
)

func TestPlatform_GetDefaultBreakoutConfig(t *testing.T) {
	tests := []struct {
		name string
		p    *Platform
		want BreakoutConfig
	}{
		{
			name: "unrealistic case where some interface does not have any breakout modes",
			p: &Platform{
				Interfaces: map[string]Interface{
					"Ethernet20": {},
				},
			},
			want: BreakoutConfig{},
		},
		{
			name: "fill in default breakout for all interfaces",
			p: &Platform{
				Interfaces: map[string]Interface{
					"Ethernet1": {
						BreakoutModes: map[string][]string{
							"1x1G": {},
						},
					},
					"Ethernet10": {
						BreakoutModes: map[string][]string{
							"1x100G[40G]": {},
							"2x50G":       {},
						},
					},
					"Ethernet20": {
						BreakoutModes: map[string][]string{
							"2x50G":       {},
							"1x100G[40G]": {},
						},
					},
					"Ethernet120": {
						BreakoutModes: map[string][]string{
							"1x100G[40G]": {},
							"2x50G":       {},
							"4x25G":       {},
							"4x10G":       {},
						},
					},
				},
			},
			want: BreakoutConfig{
				"Ethernet1":   "1x1G",
				"Ethernet10":  "1x100G[40G]",
				"Ethernet20":  "1x100G[40G]",
				"Ethernet120": "1x100G[40G]",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.p.GetDefaultBreakoutConfig()
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Platform.GetDefaultBreakoutConfig() diff = %v", diff)
			}
		})
	}
}

func TestPlatform_ParseBreakout(t *testing.T) {
	tests := []struct {
		name     string
		p        *Platform
		portName string
		breakout string
		want     *BreakoutPorts
		wantErr  error
	}{
		{
			name:     "port not found",
			p:        &Platform{},
			portName: "Ethernet0",
			breakout: "",
			want:     nil,
			wantErr:  fmt.Errorf("unknown port Ethernet0"),
		},
		{
			name: "invalid breakout mode",
			p: &Platform{
				Interfaces: map[string]Interface{
					"Ethernet0": {
						BreakoutModes: map[string][]string{
							"1x100G[40G]": {},
						},
					},
				},
			},
			portName: "Ethernet0",
			breakout: "1x1G",
			want:     nil,
			wantErr:  fmt.Errorf("invalid breakout mode 1x1G; must be one of [1x100G[40G]]"),
		},
		{
			name: "one port with one lane",
			p: &Platform{
				Interfaces: map[string]Interface{
					"Ethernet0": {
						BreakoutModes: map[string][]string{
							"1x1G": {
								"Eth1(Port1)",
							},
						},
						Index: "1",
						Lanes: "25",
					},
				},
			},
			portName: "Ethernet0",
			breakout: "1x1G",
			want: &BreakoutPorts{
				PortAliases: []string{"Eth1(Port1)"},
				Lanes:       []int{25},
				Index:       []int{1},
			},
			wantErr: nil,
		},
		{
			name: "one port with multiple lanes",
			p: &Platform{
				Interfaces: map[string]Interface{
					"Ethernet56": {
						BreakoutModes: map[string][]string{
							"1x100G[40G]": {
								"Eth54/1(Port54)",
								"Eth54/2(Port54)",
								"Eth54/3(Port54)",
								"Eth54/4(Port54)",
							},
						},
						Index: "54,54,54,54",
						Lanes: "69,70,71,72",
					},
				},
			},
			portName: "Ethernet56",
			breakout: "1x100G[40G]",
			want: &BreakoutPorts{
				PortAliases: []string{
					"Eth54/1(Port54)",
					"Eth54/2(Port54)",
					"Eth54/3(Port54)",
					"Eth54/4(Port54)",
				},
				Lanes: []int{69, 70, 71, 72},
				Index: []int{54, 54, 54, 54},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.ParseBreakout(tt.portName, tt.breakout)
			if diff := cmp.Diff(tt.wantErr, err, testcommon.ErrorStringComparer()); diff != "" {
				t.Errorf("Platform.ParseBreakout() error diff = %v", diff)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Platform.ParseBreakout() diff = %v", diff)
			}
		})
	}
}

func TestPlatform_ParseSpeedOptions(t *testing.T) {
	tests := []struct {
		name         string
		breakoutMode string
		want         SpeedOptions
		wantErr      bool
	}{
		{
			name:         "invalid port number",
			breakoutMode: "Onex100G[40G]",
			want:         SpeedOptions{},
			wantErr:      true,
		},
		{
			name:         "missing 'x'",
			breakoutMode: "100G[40G]",
			want:         SpeedOptions{},
			wantErr:      true,
		},
		{
			name:         "missing 'G'",
			breakoutMode: "1x100",
			want:         SpeedOptions{},
			wantErr:      true,
		},
		{
			name:         "invalid speed",
			breakoutMode: "1x-100G",
			want:         SpeedOptions{},
			wantErr:      true,
		},
		{
			name:         "invalid alt speed syntax",
			breakoutMode: "1x100G(40G)",
			want:         SpeedOptions{},
			wantErr:      true,
		},
		{
			name:         "missing 'G' for alt speed",
			breakoutMode: "1x100G[40]",
			want:         SpeedOptions{},
			wantErr:      true,
		},
		{
			name:         "invalid alt speed",
			breakoutMode: "1x100G[-40G]",
			want:         SpeedOptions{},
			wantErr:      true,
		},
		{
			name:         "negative number",
			breakoutMode: "-1x100G[40G]",
			want:         SpeedOptions{},
			wantErr:      true,
		},
		{
			name:         "only one speed option",
			breakoutMode: "1x1G",
			want:         SpeedOptions{1000},
			wantErr:      false,
		},
		{
			name:         "two speed options",
			breakoutMode: "1x100G[40G]",
			want:         SpeedOptions{100000, 40000},
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSpeedOptions(tt.breakoutMode)
			if (err != nil) != tt.wantErr {
				t.Errorf("Platform.ParseSpeedOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Platform.ParseSpeedOptions() diff = %v", diff)
			}
		})
	}
}

func Test_stringToIntSlice(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []int
		wantErr bool
	}{
		{
			name:    "parse error",
			input:   "1,2,3,a",
			want:    []int{},
			wantErr: true,
		},
		{
			name:    "return slice of ints",
			input:   "1,2,3,4",
			want:    []int{1, 2, 3, 4},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := stringToIntSlice(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("stringToIntSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("stringToIntSlice() diff = %v", diff)
			}
		})
	}
}
