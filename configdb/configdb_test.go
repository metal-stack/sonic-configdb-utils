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
	pxe := values.ACL{
		Name:        "PXE_IN",
		Description: "pxe",
		Ports:       []string{"Vlan4000"},
		Rules: []values.ACLRule{
			{Action: "forward", Dst: "10.254.254.130/32"},
			{Action: "drop", Dst: "10.254.0.0/16"},
			{Action: "forward"},
		},
	}
	pxeRules := map[string]ACLRule{
		"PXE_IN|RULE_1": {DstIP: "10.254.254.130/32", PacketAction: PacketActionForward, Priority: "30"},
		"PXE_IN|RULE_2": {DstIP: "10.254.0.0/16", PacketAction: PacketActionDrop, Priority: "20"},
		"PXE_IN|RULE_3": {IPType: "IPV4ANY", PacketAction: PacketActionForward, Priority: "10"},
	}
	pxeTable := ACLTable{PolicyDesc: "pxe", Ports: []string{"Vlan4000"}, Stage: "ingress", Type: "L3"}
	interfaces := aclInterfaces{
		defined:        map[string]bool{"Ethernet0": true, "Vlan4000": true, "Vlan4001": true},
		portsGenerated: true,
	}

	tests := []struct {
		name         string
		sourceRanges []string
		acls         []values.ACL
		wantRules    map[string]ACLRule
		wantTables   map[string]ACLTable
		wantErr      string
	}{
		{
			name: "nothing given",
		},
		{
			name:       "rules get names and descending priorities in list order",
			acls:       []values.ACL{pxe},
			wantRules:  pxeRules,
			wantTables: map[string]ACLTable{"PXE_IN": pxeTable},
		},
		{
			name: "src and egress stage",
			acls: []values.ACL{{
				Name:  "OUT",
				Ports: []string{"Ethernet0"},
				Stage: "egress",
				Rules: []values.ACLRule{{Action: "drop", Src: "10.0.0.0/8", Dst: "192.168.0.0/16"}},
			}},
			wantRules: map[string]ACLRule{
				"OUT|RULE_1": {DstIP: "192.168.0.0/16", PacketAction: PacketActionDrop, Priority: "10", SrcIP: "10.0.0.0/8"},
			},
			wantTables: map[string]ACLTable{"OUT": {Ports: []string{"Ethernet0"}, Stage: "egress", Type: "L3"}},
		},
		{
			name:         "custom acls are added next to the ssh acls",
			sourceRanges: []string{"10.1.23.1/30"},
			acls:         []values.ACL{pxe},
			wantRules: map[string]ACLRule{
				"ALLOW_SSH|DEFAULT_RULE": {EtherType: "2048", PacketAction: PacketActionDrop, Priority: "1"},
				"ALLOW_NTP|DEFAULT_RULE": {EtherType: "2048", PacketAction: PacketActionDrop, Priority: "1"},
				"ALLOW_NTP|RULE_1":       {PacketAction: PacketActionAccept, Priority: "99", SrcIP: "0.0.0.0/0"},
				"ALLOW_SSH|RULE_1":       {PacketAction: PacketActionAccept, Priority: "91", SrcIP: "10.1.23.1/30"},
				"PXE_IN|RULE_1":          pxeRules["PXE_IN|RULE_1"],
				"PXE_IN|RULE_2":          pxeRules["PXE_IN|RULE_2"],
				"PXE_IN|RULE_3":          pxeRules["PXE_IN|RULE_3"],
			},
			wantTables: map[string]ACLTable{
				"ALLOW_SSH": {PolicyDesc: "Allow SSH access", Services: []string{"SSH"}, Stage: "ingress", Type: "CTRLPLANE"},
				"ALLOW_NTP": {PolicyDesc: "Allow NTP", Services: []string{"NTP"}, Stage: "ingress", Type: "CTRLPLANE"},
				"PXE_IN":    pxeTable,
			},
		},
		{
			name:    "name reserved for the ssh acls even without ssh sourceranges",
			acls:    []values.ACL{{Name: "ALLOW_SSH", Ports: []string{"Vlan4000"}, Rules: []values.ACLRule{{Action: "drop"}}}},
			wantErr: `acl "ALLOW_SSH" is reserved for the acls generated from ssh_sourceranges, choose another name`,
		},
		{
			name:    "missing name",
			acls:    []values.ACL{{Ports: []string{"Vlan4000"}, Rules: []values.ACLRule{{Action: "drop"}}}},
			wantErr: "acl number 1 has no name, set one",
		},
		{
			name:    "duplicate name",
			acls:    []values.ACL{pxe, pxe},
			wantErr: `acl "PXE_IN" is defined twice, merge the rules into one entry`,
		},
		{
			name:    "no ports",
			acls:    []values.ACL{{Name: "PXE_IN", Rules: []values.ACLRule{{Action: "drop"}}}},
			wantErr: `acl "PXE_IN" has no ports, list the ports or vlans it applies to`,
		},
		{
			name:    "no rules",
			acls:    []values.ACL{{Name: "PXE_IN", Ports: []string{"Vlan4000"}}},
			wantErr: `acl "PXE_IN" has no rules, add at least one`,
		},
		{
			name:    "unknown stage",
			acls:    []values.ACL{{Name: "PXE_IN", Ports: []string{"Vlan4000"}, Stage: "in", Rules: []values.ACLRule{{Action: "drop"}}}},
			wantErr: `acl "PXE_IN" has stage "in", use ingress or egress`,
		},
		{
			name:    "unknown action",
			acls:    []values.ACL{{Name: "PXE_IN", Ports: []string{"Vlan4000"}, Rules: []values.ACLRule{{Action: "accept"}}}},
			wantErr: `acl "PXE_IN" rule 1 has action "accept", use forward or drop`,
		},
		{
			name:    "invalid dst",
			acls:    []values.ACL{{Name: "PXE_IN", Ports: []string{"Vlan4000"}, Rules: []values.ACLRule{{Action: "drop", Dst: "10.254.0.0"}}}},
			wantErr: `acl "PXE_IN" rule 1 has dst "10.254.0.0", use an ipv4 cidr such as 10.0.0.0/16`,
		},
		{
			name:    "name with the configdb key separator",
			acls:    []values.ACL{{Name: "PXE|IN", Ports: []string{"Vlan4000"}, Rules: []values.ACLRule{{Action: "drop"}}}},
			wantErr: `acl "PXE|IN" has an invalid name, use only letters, digits, '_' and '-'`,
		},
		{
			name:    "undefined vlan",
			acls:    []values.ACL{{Name: "PXE_IN", Ports: []string{"Vlan400"}, Rules: []values.ACLRule{{Action: "drop"}}}},
			wantErr: `acl "PXE_IN" is bound to "Vlan400", which is not defined, add it to vlans, portchannels or ports first`,
		},
		{
			name:    "undefined ethernet port",
			acls:    []values.ACL{{Name: "PXE_IN", Ports: []string{"Ethernet9"}, Rules: []values.ACLRule{{Action: "drop"}}}},
			wantErr: `acl "PXE_IN" is bound to "Ethernet9", which is not defined, add it to vlans, portchannels or ports first`,
		},
		{
			name: "two acls on the same port and stage",
			acls: []values.ACL{
				pxe,
				{Name: "PXE_IN_2", Ports: []string{"Vlan4001", "Vlan4000"}, Rules: []values.ACLRule{{Action: "drop"}}},
			},
			wantErr: `acls "PXE_IN" and "PXE_IN_2" are both bound to "Vlan4000" at stage ingress, merge their rules into one acl`,
		},
		{
			name:    "same port listed twice in one acl",
			acls:    []values.ACL{{Name: "PXE_IN", Ports: []string{"Vlan4000", "Vlan4000"}, Rules: []values.ACLRule{{Action: "drop"}}}},
			wantErr: `acl "PXE_IN" lists "Vlan4000" twice, remove one`,
		},
		{
			name: "same port with different stages",
			acls: []values.ACL{
				{Name: "IN", Ports: []string{"Vlan4000"}, Rules: []values.ACLRule{{Action: "drop"}}},
				{Name: "OUT", Ports: []string{"Vlan4000"}, Stage: "egress", Rules: []values.ACLRule{{Action: "drop"}}},
			},
			wantRules: map[string]ACLRule{
				"IN|RULE_1":  {IPType: "IPV4ANY", PacketAction: PacketActionDrop, Priority: "10"},
				"OUT|RULE_1": {IPType: "IPV4ANY", PacketAction: PacketActionDrop, Priority: "10"},
			},
			wantTables: map[string]ACLTable{
				"IN":  {Ports: []string{"Vlan4000"}, Stage: "ingress", Type: "L3"},
				"OUT": {Ports: []string{"Vlan4000"}, Stage: "egress", Type: "L3"},
			},
		},
		{
			name:    "host bits set",
			acls:    []values.ACL{{Name: "PXE_IN", Ports: []string{"Vlan4000"}, Rules: []values.ACLRule{{Action: "forward", Dst: "10.254.254.130/16"}}}},
			wantErr: `acl "PXE_IN" rule 1 has dst "10.254.254.130/16" with host bits set, use 10.254.0.0/16`,
		},
		{
			name:    "dst is reported before src",
			acls:    []values.ACL{{Name: "PXE_IN", Ports: []string{"Vlan4000"}, Rules: []values.ACLRule{{Action: "drop", Dst: "x", Src: "y"}}}},
			wantErr: `acl "PXE_IN" rule 1 has dst "x", use an ipv4 cidr such as 10.0.0.0/16`,
		},
		{
			name:    "ipv6 src",
			acls:    []values.ACL{{Name: "PXE_IN", Ports: []string{"Vlan4000"}, Rules: []values.ACLRule{{Action: "drop", Src: "fd00::/8"}}}},
			wantErr: `acl "PXE_IN" rule 1 has src "fd00::/8", use an ipv4 cidr such as 10.0.0.0/16`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRules, gotTables, err := getACLRulesAndTables(tt.sourceRanges, tt.acls, interfaces)
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

func Test_getACLRulesAndTables_portChannelMember(t *testing.T) {
	ports := map[string]Port{"Ethernet0": {}, "Ethernet4": {}}
	portChannels := []values.PortChannel{{Number: "1", Members: []string{"Ethernet4"}}}
	interfaces := getACLInterfaces(ports, nil, portChannels)
	acls := []values.ACL{{Name: "IN", Ports: []string{"Ethernet0", "Ethernet4"}, Rules: []values.ACLRule{{Action: "drop"}}}}

	_, _, err := getACLRulesAndTables(nil, acls, interfaces)

	want := `acl "IN" is bound to "Ethernet4", which is a member of PortChannel1, bind it to PortChannel1 instead`
	if diff := cmp.Diff(want, errorString(err)); diff != "" {
		t.Errorf("getACLRulesAndTables() error diff = %s", diff)
	}
}

func Test_aclInterfaces_canBind(t *testing.T) {
	ports := map[string]Port{"Ethernet0": {}}
	vlans := []values.VLAN{{ID: "4000"}}
	portChannels := []values.PortChannel{{Number: "1"}}

	tests := []struct {
		name  string
		ports map[string]Port
		iface string
		want  bool
	}{
		{name: "generated port", ports: ports, iface: "Ethernet0", want: true},
		{name: "port missing from the generated ports", ports: ports, iface: "Ethernet1", want: false},
		{name: "ethernet ports are not checked without a ports section", ports: nil, iface: "Ethernet1", want: true},
		{name: "defined vlan", ports: nil, iface: "Vlan4000", want: true},
		{name: "undefined vlan without a ports section", ports: nil, iface: "Vlan4001", want: false},
		{name: "defined portchannel", ports: ports, iface: "PortChannel1", want: true},
		{name: "undefined portchannel", ports: ports, iface: "PortChannel2", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getACLInterfaces(tt.ports, vlans, portChannels).canBind(tt.iface)
			if got != tt.want {
				t.Errorf("canBind(%q) = %v, want %v", tt.iface, got, tt.want)
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
