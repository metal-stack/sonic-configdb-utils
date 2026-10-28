package configdb

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/metal-stack/metal-lib/pkg/testcommon"
	p "github.com/metal-stack/sonic-configdb-utils/platform"
	"github.com/metal-stack/sonic-configdb-utils/values"
	v "github.com/metal-stack/sonic-configdb-utils/version"
)

func Test_getInterfaces(t *testing.T) {
	tests := []struct {
		name          string
		ports         *values.Ports
		bgpPorts      []string
		interconnects map[string]values.Interconnect
		interfaces    map[string]values.Interface
		want          map[string]Interface
	}{
		{
			name: "empty ports",
			ports: &values.Ports{
				List: []values.Port{},
			},
			want: map[string]Interface{},
		},
		{
			name: "port not in bgp ports, with no vrf and no ips",
			ports: &values.Ports{
				List: []values.Port{
					{
						Name: "Ethernet0",
					},
				},
			},
			want: map[string]Interface{},
		},
		{
			name: "port in bgp ports",
			ports: &values.Ports{
				List: []values.Port{
					{
						Name: "Ethernet0",
					},
				},
			},
			bgpPorts: []string{"Ethernet0"},
			want: map[string]Interface{
				"Ethernet0": {
					IPv6UseLinkLocalOnly: IPv6UseLinkLocalOnlyModeEnable,
				},
			},
		},
		{
			name: "port with vrf",
			ports: &values.Ports{
				List: []values.Port{
					{
						Name: "Ethernet0",
						VRF:  "Vrf40",
					},
				},
			},
			want: map[string]Interface{
				"Ethernet0": {
					VRFName: "Vrf40",
				},
			},
		},
		{
			name: "port in bgp ports with vrf",
			ports: &values.Ports{
				List: []values.Port{
					{
						Name: "Ethernet0",
						VRF:  "Vrf40",
					},
				},
			},
			bgpPorts: []string{"Ethernet0"},
			want: map[string]Interface{
				"Ethernet0": {
					IPv6UseLinkLocalOnly: IPv6UseLinkLocalOnlyModeEnable,
					VRFName:              "Vrf40",
				},
			},
		},
		{
			name: "port not in bgp ports without vrf but with ip",
			ports: &values.Ports{
				List: []values.Port{
					{
						Name: "Ethernet0",
						IPs:  []string{"10.1.1.1"},
					},
				},
			},
			want: map[string]Interface{
				"Ethernet0":          {},
				"Ethernet0|10.1.1.1": {},
			},
		},
		{
			name: "port in bgp ports with vrf and ips",
			ports: &values.Ports{
				List: []values.Port{
					{
						Name: "Ethernet0",
						VRF:  "Vrf40",
						IPs:  []string{"10.1.1.1"},
					},
				},
			},
			bgpPorts: []string{"Ethernet0"},
			want: map[string]Interface{
				"Ethernet0": {
					IPv6UseLinkLocalOnly: IPv6UseLinkLocalOnlyModeEnable,
					VRFName:              "Vrf40",
				},
				"Ethernet0|10.1.1.1": {},
			},
		},
		{
			name:     "port in bgp ports but not in ports",
			bgpPorts: []string{"Ethernet0"},
			want: map[string]Interface{
				"Ethernet0": {
					IPv6UseLinkLocalOnly: IPv6UseLinkLocalOnlyModeEnable,
				},
			},
		},
		{
			name: "interconnect without unnumbered interfaces",
			interconnects: map[string]values.Interconnect{
				"internet": {
					UnnumberedInterfaces: []string{},
					VNI:                  "104000",
					VRF:                  "VrfInternet",
				},
			},
			want: map[string]Interface{},
		},
		{
			name: "interconnect with unnumbered interfaces",
			interconnects: map[string]values.Interconnect{
				"internet": {
					UnnumberedInterfaces: []string{
						"Ethernet0",
						"Ethernet1",
					},
					VNI: "104000",
					VRF: "VrfInternet",
				},
			},
			want: map[string]Interface{
				"Ethernet0": {
					IPv6UseLinkLocalOnly: IPv6UseLinkLocalOnlyModeEnable,
					VRFName:              "VrfInternet",
				},
				"Ethernet1": {
					IPv6UseLinkLocalOnly: IPv6UseLinkLocalOnlyModeEnable,
					VRFName:              "VrfInternet",
				},
			},
		},
		{
			name: "interconnect with unnumbered interfaces but without vrf",
			interconnects: map[string]values.Interconnect{
				"internet": {
					UnnumberedInterfaces: []string{
						"Ethernet0",
						"Ethernet1",
					},
				},
			},
			want: map[string]Interface{
				"Ethernet0": {
					IPv6UseLinkLocalOnly: IPv6UseLinkLocalOnlyModeEnable,
				},
				"Ethernet1": {
					IPv6UseLinkLocalOnly: IPv6UseLinkLocalOnlyModeEnable,
				},
			},
		},
		{
			name: "port with ip in interconnects",
			ports: &values.Ports{
				List: []values.Port{
					{
						Name: "Ethernet0",
						IPs:  []string{"10.0.0.1/32"},
					},
				},
			},
			interconnects: map[string]values.Interconnect{
				"internal": {
					UnnumberedInterfaces: []string{"Ethernet0"},
					VRF:                  "internal",
				},
			},
			want: map[string]Interface{
				"Ethernet0": {
					IPv6UseLinkLocalOnly: IPv6UseLinkLocalOnlyModeEnable,
					VRFName:              "internal",
				},
				"Ethernet0|10.0.0.1/32": {},
			},
		},
		{
			name: "additional interfaces with ips",
			interfaces: map[string]values.Interface{
				"Ethernet50.10": {
					IPs: []string{
						"172.16.0.1/32",
						"10.1.1.1/32",
					},
				},
			},
			want: map[string]Interface{
				"Ethernet50.10":               {},
				"Ethernet50.10|172.16.0.1/32": {},
				"Ethernet50.10|10.1.1.1/32":   {},
			},
		},
		{
			name: "additional ips for existing interface entry",
			ports: &values.Ports{
				List: []values.Port{
					{
						Name: "Ethernet0",
						IPs:  []string{"10.0.0.1/32"},
					},
				},
			},
			interconnects: map[string]values.Interconnect{
				"internal": {
					UnnumberedInterfaces: []string{"Ethernet0"},
					VRF:                  "internal",
				},
			},
			interfaces: map[string]values.Interface{
				"Ethernet0": {
					IPs: []string{"172.16.0.1/32"},
				},
			},
			want: map[string]Interface{
				"Ethernet0": {
					IPv6UseLinkLocalOnly: IPv6UseLinkLocalOnlyModeEnable,
					VRFName:              "internal",
				},
				"Ethernet0|10.0.0.1/32":   {},
				"Ethernet0|172.16.0.1/32": {},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getInterfaces(tt.ports, tt.bgpPorts, tt.interconnects, tt.interfaces)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("getInterfaces() %v", diff)
			}
		})
	}
}

func Test_getPortsAndBreakouts(t *testing.T) {
	tests := []struct {
		name          string
		ports         *values.Ports
		breakouts     map[string]string
		platform      *p.Platform
		wantPorts     map[string]Port
		wantBreakouts map[string]BreakoutConfig
		wantErr       error
	}{
		{
			name: "only breakouts defined",
			ports: &values.Ports{
				List: []values.Port{},
			},
			breakouts: map[string]string{
				"Ethernet0": "1x100G[40G]",
			},
			platform: &p.Platform{
				Interfaces: map[string]p.Interface{
					"Ethernet0": {
						BreakoutModes: map[string][]string{
							"1x100G[40G]": {"Eth1(Port1)"},
						},
						Index: "1,1,1,1",
						Lanes: "1,2,3,4",
					},
				},
			},
			wantPorts: map[string]Port{
				"Ethernet0": {
					AdminStatus:    defaultAdminStatus,
					Alias:          "Eth1(Port1)",
					Index:          "1",
					Lanes:          "1,2,3,4",
					MTU:            fmt.Sprintf("%d", defaultMTU),
					parentBreakout: "1x100G[40G]",
					Speed:          "100000",
				},
			},
			wantBreakouts: map[string]BreakoutConfig{
				"Ethernet0": {
					BreakoutMode: "1x100G[40G]",
				},
			},
			wantErr: nil,
		},
		{
			name: "only ports defined",
			ports: &values.Ports{
				List: []values.Port{
					{
						Name:    "Ethernet0",
						Speed:   40000,
						Autoneg: values.AutonegModeOn,
					},
					{
						Name:  "Ethernet4",
						Speed: 40000,
					},
				},
			},
			platform: &p.Platform{
				Interfaces: map[string]p.Interface{
					"Ethernet0": {
						BreakoutModes: map[string][]string{
							"1x100G[40G]": {"Eth1(Port1)"},
							"2x50G":       {"Eth1/1(Port1)", "Eth1/2(Port1)"},
						},
						Index: "1,1,1,1",
						Lanes: "1,2,3,4",
					},
					"Ethernet4": {
						BreakoutModes: map[string][]string{
							"1x100G[40G]": {"Eth2(Port2)"},
							"2x50G":       {"Eth2/1(Port2)", "Eth2/2(Port2)"},
						},
						Index: "2,2,2,2",
						Lanes: "5,6,7,8",
					},
					"Ethernet8": {
						BreakoutModes: map[string][]string{
							"1x100G[40G]": {"Eth3(Port3)"},
							"2x50G":       {"Eth3/1(Port3)", "Eth3/2(Port3)"},
						},
						Index: "3,3,3,3",
						Lanes: "9,10,11,12",
					},
				},
			},
			wantPorts: map[string]Port{
				"Ethernet0": {
					AdminStatus:    AdminStatusUp,
					Alias:          "Eth1(Port1)",
					Autoneg:        AutonegModeOn,
					Index:          "1",
					Lanes:          "1,2,3,4",
					MTU:            fmt.Sprintf("%d", defaultMTU),
					parentBreakout: "1x100G[40G]",
					Speed:          "40000",
				},
				"Ethernet4": {
					AdminStatus:    AdminStatusUp,
					Alias:          "Eth2(Port2)",
					Index:          "2",
					Lanes:          "5,6,7,8",
					MTU:            fmt.Sprintf("%d", defaultMTU),
					parentBreakout: "1x100G[40G]",
					Speed:          "40000",
				},
				"Ethernet8": {
					AdminStatus:    AdminStatusUp,
					Alias:          "Eth3(Port3)",
					Index:          "3",
					Lanes:          "9,10,11,12",
					MTU:            fmt.Sprintf("%d", defaultMTU),
					parentBreakout: "1x100G[40G]",
					Speed:          "100000",
				},
			},
			wantBreakouts: map[string]BreakoutConfig{
				"Ethernet0": {
					BreakoutMode: "1x100G[40G]",
				},
				"Ethernet4": {
					BreakoutMode: "1x100G[40G]",
				},
				"Ethernet8": {
					BreakoutMode: "1x100G[40G]",
				},
			},
			wantErr: nil,
		},
		{
			name: "child-port is not present if breakout config 'absorbs' it",
			breakouts: map[string]string{
				"Ethernet4": "1x100G[40G]",
			},
			ports: &values.Ports{
				List: []values.Port{
					{
						Name: "Ethernet5",
					},
				},
			},
			platform: &p.Platform{
				Interfaces: map[string]p.Interface{
					"Ethernet4": {
						BreakoutModes: map[string][]string{
							"1x100G[40G]": {"Eth2(Port2)"},
							"4x25G":       {"Eth2/1(Port2)", "Eth2/2(Port2)", "Eth2/3(Port2)", "Eth2/4(Port2)"},
						},
						Index: "2,2,2,2",
						Lanes: "5,6,7,8",
					},
				},
			},
			wantErr: fmt.Errorf("invalid port name Ethernet5; if you think it should be available please check your breakout configuration"),
		},
		{
			name: "port speed 0 is allowed",
			breakouts: map[string]string{
				"Ethernet4": "1x100G[40G]",
			},
			ports: &values.Ports{
				DefaultFEC: values.FECModeRS,
				DefaultMTU: 1500,
				List: []values.Port{
					{
						Name:  "Ethernet4",
						Speed: 0,
					},
				}},
			platform: &p.Platform{
				Interfaces: map[string]p.Interface{
					"Ethernet4": {
						BreakoutModes: map[string][]string{
							"1x100G[40G]": {"Eth2(Port2)"},
							"4x25G":       {"Eth2/1(Port2)", "Eth2/2(Port2)", "Eth2/3(Port2)", "Eth2/4(Port2)"},
						},
						Index: "2,2,2,2",
						Lanes: "5,6,7,8",
					},
				},
			},
			wantBreakouts: map[string]BreakoutConfig{
				"Ethernet4": {
					BreakoutMode: "1x100G[40G]",
				},
			},
			wantPorts: map[string]Port{
				"Ethernet4": {
					AdminStatus:    AdminStatusUp,
					Alias:          "Eth2(Port2)",
					FEC:            FECModeRS,
					Index:          "2",
					Lanes:          "5,6,7,8",
					MTU:            "1500",
					parentBreakout: "1x100G[40G]",
					Speed:          "100000",
				},
			},
			wantErr: nil,
		},
		{
			name: "fec, mtu and speed can be overriden",
			breakouts: map[string]string{
				"Ethernet4": "1x100G[40G]",
			},
			ports: &values.Ports{
				List: []values.Port{
					{
						FECMode: values.FECModeRS,
						MTU:     1500,
						Name:    "Ethernet4",
						Speed:   40000,
					},
				},
			},
			platform: &p.Platform{
				Interfaces: map[string]p.Interface{
					"Ethernet4": {
						BreakoutModes: map[string][]string{
							"1x100G[40G]": {"Eth2(Port2)"},
							"4x25G":       {"Eth2/1(Port2)", "Eth2/2(Port2)", "Eth2/3(Port2)", "Eth2/4(Port2)"},
						},
						Index: "2,2,2,2",
						Lanes: "5,6,7,8",
					},
				},
			},
			wantBreakouts: map[string]BreakoutConfig{
				"Ethernet4": {
					BreakoutMode: "1x100G[40G]",
				},
			},
			wantPorts: map[string]Port{
				"Ethernet4": {
					AdminStatus:    AdminStatusUp,
					Alias:          "Eth2(Port2)",
					FEC:            FECModeRS,
					Index:          "2",
					Lanes:          "5,6,7,8",
					MTU:            "1500",
					parentBreakout: "1x100G[40G]",
					Speed:          "40000",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPorts, gotBreakouts, err := getPortsAndBreakouts(tt.ports, tt.breakouts, tt.platform)
			if diff := cmp.Diff(tt.wantErr, err, testcommon.ErrorStringComparer()); diff != "" {
				t.Errorf("getPortsAndBreakouts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.wantPorts, gotPorts, testcommon.IgnoreUnexported()); diff != "" {
				t.Errorf("getPortsAndBreakouts() diff = %v", diff)
			}
			if diff := cmp.Diff(tt.wantBreakouts, gotBreakouts); diff != "" {
				t.Errorf("getPortsAndBreakouts() diff = %v", diff)
			}
		})
	}
}

func Test_getVRFs(t *testing.T) {
	tests := []struct {
		name          string
		interconnects map[string]values.Interconnect
		ports         *values.Ports
		vlans         []values.VLAN
		want          map[string]VRF
	}{
		{
			name: "no ports or vlans to add",
			interconnects: map[string]values.Interconnect{
				"mpls": {
					VNI: "1",
					VRF: "Vrf40",
				},
			},
			ports: &values.Ports{
				List: []values.Port{},
			},
			vlans: []values.VLAN{},
			want: map[string]VRF{
				"Vrf40": {
					VNI: "1",
				},
			},
		},
		{
			name: "duplicates are not added",
			interconnects: map[string]values.Interconnect{
				"mpls": {
					VNI: "1",
					VRF: "Vrf40",
				},
			},
			ports: &values.Ports{
				List: []values.Port{
					{
						VRF: "Vrf40",
					},
				},
			},
			vlans: []values.VLAN{
				{
					VRF: "Vrf40",
				},
			},
			want: map[string]VRF{
				"Vrf40": {
					VNI: "1",
				},
			},
		},
		{
			name: "new vrfs are added",
			interconnects: map[string]values.Interconnect{
				"mpls": {
					VNI: "1",
					VRF: "Vrf40",
				},
			},
			ports: &values.Ports{
				List: []values.Port{
					{
						VRF: "Vrf41",
					},
				},
			},
			vlans: []values.VLAN{
				{
					VRF: "Vrf42",
				},
			},
			want: map[string]VRF{
				"Vrf40": {
					VNI: "1",
				},
				"Vrf41": {},
				"Vrf42": {},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getVRFs(tt.interconnects, tt.ports, tt.vlans)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("getVRFs() diff = %s", diff)
			}
		})
	}
}

func Test_getSAG(t *testing.T) {
	tests := []struct {
		name    string
		sag     *values.SAG
		version v.Branch
		want    *SAG
		wantErr bool
	}{
		{
			name:    "wrong version",
			sag:     &values.SAG{},
			version: v.Branch202111,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty mac",
			sag:     &values.SAG{},
			version: v.Branch202211,
			want:    nil,
			wantErr: false,
		},
		{
			name: "valid",
			sag: &values.SAG{
				MAC: "11:11:11:11:11:11",
			},
			version: v.Branch202211,
			want: &SAG{
				SAGGlobal: SAGGlobal{
					GatewayMAC: "11:11:11:11:11:11",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getSAG(tt.sag, tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("getSAG() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("getSAG() diff = %s", diff)
			}
		})
	}
}

func Test_getACLRulesAndTables(t *testing.T) {
	pxeTable := map[string]values.ACLTable{
		"PXE_IN": {PolicyDesc: "pxe", Ports: []string{"Vlan4000"}, Stage: "ingress", Type: "L3"},
	}

	tests := []struct {
		name         string
		sourceRanges []string
		acl          values.ACL
		wantRules    map[string]ACLRule
		wantTables   map[string]ACLTable
		wantErr      string
	}{
		{
			name: "nothing given",
		},
		{
			name: "custom table and rule without ssh sourceranges",
			acl: values.ACL{
				Tables: pxeTable,
				Rules: map[string]values.ACLRule{
					"PXE_IN|DENY_1": {DstIP: "10.0.0.0/16", PacketAction: "DROP", Priority: "7999"},
				},
			},
			wantRules: map[string]ACLRule{
				"PXE_IN|DENY_1": {DstIP: "10.0.0.0/16", PacketAction: PacketActionDrop, Priority: "7999"},
			},
			wantTables: map[string]ACLTable{
				"PXE_IN": {PolicyDesc: "pxe", Ports: []string{"Vlan4000"}, Stage: "ingress", Type: "L3"},
			},
		},
		{
			name:         "custom acls are added next to the ssh acls",
			sourceRanges: []string{"10.1.23.1/30"},
			acl: values.ACL{
				Tables: pxeTable,
				Rules: map[string]values.ACLRule{
					"PXE_IN|PERMIT_REST": {IPType: "IPV4ANY", PacketAction: "FORWARD", Priority: "1000"},
				},
			},
			wantRules: map[string]ACLRule{
				"ALLOW_SSH|DEFAULT_RULE": {EtherType: "2048", PacketAction: PacketActionDrop, Priority: "1"},
				"ALLOW_NTP|DEFAULT_RULE": {EtherType: "2048", PacketAction: PacketActionDrop, Priority: "1"},
				"ALLOW_NTP|RULE_1":       {PacketAction: PacketActionAccept, Priority: "99", SrcIP: "0.0.0.0/0"},
				"ALLOW_SSH|RULE_1":       {PacketAction: PacketActionAccept, Priority: "91", SrcIP: "10.1.23.1/30"},
				"PXE_IN|PERMIT_REST":     {IPType: "IPV4ANY", PacketAction: PacketActionForward, Priority: "1000"},
			},
			wantTables: map[string]ACLTable{
				"ALLOW_SSH": {PolicyDesc: "Allow SSH access", Services: []string{"SSH"}, Stage: "ingress", Type: "CTRLPLANE"},
				"ALLOW_NTP": {PolicyDesc: "Allow NTP", Services: []string{"NTP"}, Stage: "ingress", Type: "CTRLPLANE"},
				"PXE_IN":    {PolicyDesc: "pxe", Ports: []string{"Vlan4000"}, Stage: "ingress", Type: "L3"},
			},
		},
		{
			name:         "table clashes with a generated ssh table",
			sourceRanges: []string{"10.1.23.1/30"},
			acl: values.ACL{
				Tables: map[string]values.ACLTable{"ALLOW_SSH": {Stage: "ingress", Type: "CTRLPLANE"}},
			},
			wantErr: `acl table "ALLOW_SSH" is generated from ssh_sourceranges, choose another name`,
		},
		{
			name:         "rule clashes with a generated ssh rule",
			sourceRanges: []string{"10.1.23.1/30"},
			acl: values.ACL{
				Rules: map[string]values.ACLRule{"ALLOW_SSH|RULE_1": {PacketAction: "ACCEPT", Priority: "50"}},
			},
			wantErr: `acl rule "ALLOW_SSH|RULE_1" is generated from ssh_sourceranges, choose another name`,
		},
		{
			name: "rule name without table",
			acl: values.ACL{
				Tables: pxeTable,
				Rules:  map[string]values.ACLRule{"DENY_1": {PacketAction: "DROP", Priority: "10"}},
			},
			wantErr: `acl rule "DENY_1" must be named TABLE|RULE`,
		},
		{
			name: "rule for an unknown table",
			acl: values.ACL{
				Tables: pxeTable,
				Rules:  map[string]values.ACLRule{"PXE|DENY_1": {PacketAction: "DROP", Priority: "10"}},
			},
			wantErr: `acl rule "PXE|DENY_1" refers to table "PXE", which is not defined in ACL_TABLE`,
		},
		{
			name: "unknown packet action",
			acl: values.ACL{
				Tables: pxeTable,
				Rules:  map[string]values.ACLRule{"PXE_IN|DENY_1": {PacketAction: "REJECT", Priority: "10"}},
			},
			wantErr: `acl rule "PXE_IN|DENY_1" has PACKET_ACTION "REJECT", use one of ACCEPT, DROP, FORWARD`,
		},
		{
			name: "rule without priority",
			acl: values.ACL{
				Tables: pxeTable,
				Rules:  map[string]values.ACLRule{"PXE_IN|DENY_1": {PacketAction: "DROP"}},
			},
			wantErr: `acl rule "PXE_IN|DENY_1" has no PRIORITY, set one`,
		},
		{
			name: "table without type",
			acl: values.ACL{
				Tables: map[string]values.ACLTable{"PXE_IN": {Stage: "ingress"}},
			},
			wantErr: `acl table "PXE_IN" has no type, set one such as L3 or CTRLPLANE`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRules, gotTables, err := getACLRulesAndTables(tt.sourceRanges, tt.acl)
			if diff := cmp.Diff(tt.wantErr, errorString(err)); diff != "" {
				t.Fatalf("getACLRulesAndTables() error diff = %s", diff)
			}
			if diff := cmp.Diff(tt.wantRules, gotRules); diff != "" {
				t.Errorf("getACLRulesAndTables() rules diff = %s", diff)
			}
			if diff := cmp.Diff(tt.wantTables, gotTables); diff != "" {
				t.Errorf("getACLRulesAndTables() tables diff = %s", diff)
			}
		})
	}
}

func errorString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
