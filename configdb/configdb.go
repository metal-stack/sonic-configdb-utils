package configdb

import (
	"fmt"
	"slices"

	"github.com/metal-stack/sonic-configdb-utils/values"
)

type ConfigDB struct {
	ACLRules           map[string]ACLRule        `json:"ACL_RULE"`
	ACLTables          map[string]ACLTable       `json:"ACL_TABLE"`
	Breakouts          map[string]BreakoutConfig `json:"BREAKOUT_CFG"`
	DeviceMetadata     `json:"DEVICE_METADATA"`
	Features           map[string]Feature   `json:"FEATURE"`
	Interfaces         map[string]Interface `json:"INTERFACE"`
	LLDP               `json:"LLDP"`
	LoopbackInterface  map[string]struct{}       `json:"LOOPBACK_INTERFACE"`
	MCLAGDomains       map[string]MCLAGDomain    `json:"MCLAG_DOMAIN"`
	MCLAGInterfaces    map[string]MCLAGInterface `json:"MCLAG_INTERFACE"`
	MCLAGUniqueIPs     map[string]MCLAGUniqueIP  `json:"MCLAG_UNIQUE_IP"`
	MgmtInterfaces     map[string]MgmtInterface  `json:"MGMT_INTERFACE"`
	MgmtPorts          map[string]MgmtPort       `json:"MGMT_PORT"`
	MgmtVRFConfig      `json:"MGMT_VRF_CONFIG"`
	NTP                `json:"NTP"`
	NTPServers         map[string]struct{}    `json:"NTP_SERVER"`
	Ports              map[string]Port        `json:"PORT"`
	Portchannels       map[string]Portchannel `json:"PORTCHANNEL"`
	PortchannelMembers map[string]struct{}    `json:"PORTCHANNEL_MEMBERS"`
	SAG                `json:"SAG"`
	VLANs              map[string]VLAN          `json:"VLAN"`
	VLANInterfaces     map[string]VLANInterface `json:"VLAN_INTERFACE"`
	VLANMembers        map[string]VLANMember    `json:"VLAN_MEMBER"`
	VRFs               map[string]VRF           `json:"VRF"`
	VXLANEVPN          `json:"VXLAN_EVPN_NVO"`
	VXLANTunnels       map[string]VXLANTunnel      `json:"VXLAN_TUNNEL"`
	VXLANTunnelMaps    []VXLANTunnelMapWithComment `json:"VXLAN_TUNNEL_MAP"`
}

func GenerateConfigDB(input *values.Values) (*ConfigDB, error) {
	ports, breakouts, err := getPortsAndBreakouts(input.Ports, input.Breakouts, input.PortsDefaultFEC, input.PortsDefaultMTU)
	if err != nil {
		return nil, err
	}

	rules, tables := getACLRulesAndTables(input.SSHSourceranges)

	configdb := ConfigDB{
		ACLRules:  rules,
		ACLTables: tables,
		Breakouts: breakouts,
		DeviceMetadata: DeviceMetadata{
			Localhost: Metadata{
				DockerRoutingConfigMode: DockerRoutingConfigMode(input.DockerRoutingConfigMode),
				FRRMgmtFrameworkConfig:  input.FRRMgmtFrameworkConfig,
				Hostname:                input.Hostname,
				RouterType:              "LeafRouter",
			},
		},
		Features: map[string]Feature{
			"dhcp_relay": {
				AutoRestart: FeatureModeEnabled,
				State:       FeatureModeEnabled,
			},
		},
		Interfaces: getInterfaces(input.Ports, input.BGPPorts),
		LLDP: LLDP{
			Global: LLDPGlobal{
				HelloTime: input.LLDPHelloTimer,
			},
		},
		LoopbackInterface: map[string]struct{}{
			"Loopback0": {},
			"Loopback0|" + input.LoopbackAddress + "/32": {},
		},
		MCLAGDomains: map[string]MCLAGDomain{
			"1": {
				MCLAGSystemID: input.SystemMAC,
				PeerIP:        input.PeerIP,
				PeerLink:      input.PeerLink,
				SourceIP:      input.SourceIP,
			},
		},
		MCLAGInterfaces: getMCLAGInterfaces(input.MemberPortChannels),
		MCLAGUniqueIPs: map[string]MCLAGUniqueIP{
			"Vlan" + input.KeepaliveVLAN: {
				UniqueIP: MCLAGUniqueIPModeEnable,
			},
		},
		MgmtInterfaces: getMgmtInterfaces(input.MgmtIfIP, input.MgmtIfGateway),
		MgmtPorts: map[string]MgmtPort{
			"eth0": {
				AdminStatus: AdminStatusUp,
				Alias:       "eth0",
				Description: "Management Port",
			},
		},
		MgmtVRFConfig: MgmtVRFConfig{
			VRFGlobal: VRFGlobal{
				MgmtVRFEnabled: input.MgmtVRF,
			},
		},
		NTP: NTP{
			NTPGlobal: NTPGlobal{
				SrcIntf: "eth0",
			},
		},
		NTPServers:         getNTPServers(input.NTPServers),
		Ports:              ports,
		Portchannels:       getPortchannels(input.Portchannels, input.PortchannelsDefaultMTU),
		PortchannelMembers: getPortchannelMembers(input.Portchannels),
		SAG:                getSAG(input.SAG),
		VLANs:              getVLANs(input.VLANs),
		VLANInterfaces:     getVLANInterfaces(input.VLANs),
		VLANMembers:        getVLANMembers(input.VLANs),
		VRFs:               getVRFs(input.Interconnects, input.Ports, input.VLANs),
		VXLANEVPN: VXLANEVPN{
			VXLANEVPNNVO: VXLANEVPNNVO{
				SourceVTEP: "vtep",
			},
		},
		VXLANTunnels: map[string]VXLANTunnel{
			"vtep": {
				SrcIP: input.LoopbackAddress,
			},
		},
		VXLANTunnelMaps: getVXLANTunnelMapWithComment(input.VTEPs),
	}
	return &configdb, nil
}

func getACLRulesAndTables(sourceRanges []string) (map[string]ACLRule, map[string]ACLTable) {
	if len(sourceRanges) == 0 {
		return nil, nil
	}

	rules := map[string]ACLRule{
		"ALLOW_SSH|DEFAULT_RULE": {
			EtherType:    "2048",
			PacketAction: PacketActionDrop,
			Priority:     "1",
		},
		"ALLOW_NTP|DEFAULT_RULE": {
			EtherType:    "2048",
			PacketAction: PacketActionDrop,
			Priority:     "1",
		},
		"ALLOW_NTP|RULE_1": {
			PacketAction: PacketActionAccept,
			Priority:     "99",
			SrcIP:        "0.0.0.0/0",
		},
	}

	tables := map[string]ACLTable{
		"ALLOW_SSH": {
			PolicyDesc: "Allow SSH access",
			Ports:      []string{},
			Services:   []string{"SSH"},
			Stage:      "ingress",
			Type:       "CTRLPLANE",
		},
		"ALLOW_NTP": {
			PolicyDesc: "Allow NTP",
			Ports:      []string{},
			Services:   []string{"NTP"},
			Stage:      "ingress",
			Type:       "CTRLPLANE",
		},
	}

	for i, cidr := range sourceRanges {
		rules[fmt.Sprintf("ALLOW_SSH|RULE_%d", i+1)] = ACLRule{
			PacketAction: PacketActionAccept,
			Priority:     fmt.Sprintf("9%d", i+1),
			SrcIP:        cidr,
		}
	}

	return rules, tables
}

func getInterfaces(ports []values.Port, bgpPorts []string) map[string]Interface {
	interfaces := make(map[string]Interface)

	for _, port := range ports {
		if len(port.IPs) == 0 && port.VRF == "" && !slices.Contains(bgpPorts, port.Name) {
			continue
		}

		intf := Interface{}
		if slices.Contains(bgpPorts, port.Name) {
			intf.IPv6UseLinkLocalOnly = IPv6UseLinkLocalOnlyModeEnable
		}
		if port.VRF != "" {
			intf.VRFName = port.VRF
		}
		interfaces[port.Name] = intf

		for _, ip := range port.IPs {
			intf = Interface{}
			interfaces[port.Name+"|"+ip] = intf
		}
	}

	return interfaces
}

func getMCLAGInterfaces(memberPortChannels []string) map[string]MCLAGInterface {
	mclagInterfaces := make(map[string]MCLAGInterface)

	for _, channel := range memberPortChannels {
		mclagInterfaces["1|Portchannel"+channel] = MCLAGInterface{
			IfType: "PortChannel",
		}
	}

	return mclagInterfaces
}

func getMgmtInterfaces(mgmtIfIP, mgmtIfGateway string) map[string]MgmtInterface {
	if mgmtIfIP == "" {
		return nil
	}

	mgmtInterfaces := make(map[string]MgmtInterface)

	eth0 := MgmtInterface{}
	if mgmtIfGateway != "" {
		eth0.GWAddr = mgmtIfGateway
	}

	mgmtInterfaces["eth0|"+mgmtIfIP] = eth0

	return mgmtInterfaces
}

func getNTPServers(servers []string) map[string]struct{} {
	ntpServers := make(map[string]struct{})

	for _, server := range servers {
		ntpServers[server] = struct{}{}
	}

	return ntpServers
}

func getPortchannels(portchannels []values.Portchannel, defaultPortchannelMTU int) map[string]Portchannel {
	configPortchannels := make(map[string]Portchannel)

	for _, pc := range portchannels {
		mtu := defaultMTU
		if defaultPortchannelMTU != 0 {
			mtu = defaultPortchannelMTU
		}
		if pc.MTU != 0 {
			mtu = pc.MTU
		}

		configPortchannels["PortChannel"+pc.Number] = Portchannel{
			AdminStatus: defaultAdminStatus,
			Fallback:    true,
			LACPKey:     LACPKeyModeAuto,
			MinLinks:    1,
			MixSpeed:    false,
			MTU:         mtu,
		}
	}

	return configPortchannels
}

func getPortchannelMembers(portchannels []values.Portchannel) map[string]struct{} {
	portchannelMembers := make(map[string]struct{})

	for _, pc := range portchannels {
		for _, member := range pc.Members {
			portchannelMembers["PortChannel"+pc.Number+"|"+member] = struct{}{}
		}
	}

	return portchannelMembers
}

func getPortsAndBreakouts(ports []values.Port, breakouts map[string]string, defaultFECMode values.FECMode, defaultMTU int) (map[string]Port, map[string]BreakoutConfig, error) {
	configPorts := make(map[string]Port)
	configBreakouts := make(map[string]BreakoutConfig)

	for portName, breakout := range breakouts {
		breakoutPorts, err := getPortsFromBreakout(portName, breakout, defaultFECMode, defaultMTU)
		if err != nil {
			return nil, nil, err
		}
		for name, port := range breakoutPorts {
			configPorts[name] = port
		}
		configBreakouts[portName] = BreakoutConfig{
			BreakoutMode: BreakoutMode(breakout),
		}
	}

	for _, port := range ports {
		configPort, ok := configPorts[port.Name]
		if !ok {
			return nil, nil, fmt.Errorf("no breakout configuration found for port %s", port.Name)
		}

		switch configPort.Speed {
		case 100000:
			if port.Speed == 100000 || port.Speed == 40000 {
				configPort.Speed = port.Speed
			} else {
				return nil, nil, fmt.Errorf("invalid speed %d for port %s; current breakout configuration only allows values 100000 or 40000", port.Speed, port.Name)
			}
		default:
			if port.Speed != configPort.Speed {
				return nil, nil, fmt.Errorf("invalid speed %d for port %s; check breakout configuration", port.Speed, port.Name)
			}
		}

		if string(port.FECMode) != string(configPort.FEC) {
			configPort.FEC = FECMode(port.FECMode)
		}
		if port.MTU != configPort.MTU {
			configPort.MTU = port.MTU
		}
		configPorts[port.Name] = configPort
	}

	return configPorts, configBreakouts, nil
}

func getSAG(sag values.SAG) SAG {
	if sag.MAC == "" {
		return SAG{}
	}

	return SAG{
		SAGGlobal: SAGGlobal{
			GatewayMAC: sag.MAC,
		},
	}
}

func getVLANs(vlans []values.VLAN) map[string]VLAN {
	configVLANs := make(map[string]VLAN)

	for _, vlan := range vlans {
		configVLANs["Vlan"+vlan.ID] = VLAN{
			DHCPServers: vlan.DHCPServers,
			VLANID:      vlan.ID,
		}
	}

	return configVLANs
}

func getVLANInterfaces(vlans []values.VLAN) map[string]VLANInterface {
	vlanInterfaces := make(map[string]VLANInterface)

	for _, vlan := range vlans {
		var vlanInterface VLANInterface

		if vlan.VRF != "" {
			vlanInterface = VLANInterface{
				StaticAnycastGateway: vlan.SAG,
				VRFName:              vlan.VRF,
			}
		}

		vlanInterfaces["Vlan"+vlan.ID] = vlanInterface

		if vlan.IP == "" {
			continue
		}

		vlanInterfaces["Vlan"+vlan.ID+"|"+vlan.IP] = VLANInterface{}
	}

	return vlanInterfaces
}

func getVLANMembers(vlans []values.VLAN) map[string]VLANMember {
	vlanMembers := make(map[string]VLANMember)

	for _, vlan := range vlans {
		for _, untagged := range vlan.UntaggedPorts {
			vlanMembers["Vlan"+vlan.ID+"|"+untagged] = VLANMember{
				TaggingMode: TaggingModeUntagged,
			}
		}
		for _, tagged := range vlan.TaggedPorts {
			vlanMembers["Vlan"+vlan.ID+"|"+tagged] = VLANMember{
				TaggingMode: TaggingModeTagged,
			}
		}
	}

	return vlanMembers
}

func getVRFs(interconnects map[string]values.Interconnect, ports []values.Port, vlans []values.VLAN) map[string]VRF {
	vrfs := make(map[string]VRF)

	for _, interconnect := range interconnects {
		vrfs[interconnect.VRF] = VRF{
			VNI: interconnect.VNI,
		}
	}

	for _, port := range ports {
		if port.VRF == "" {
			continue
		}
		if _, ok := vrfs[port.VRF]; ok {
			continue
		}
		vrfs[port.VRF] = VRF{}
	}

	for _, vlan := range vlans {
		if vlan.VRF == "" {
			continue
		}
		if _, ok := vrfs[vlan.VRF]; ok {
			continue
		}
		vrfs[vlan.VRF] = VRF{}
	}

	return vrfs
}

func getVXLANTunnelMapWithComment(vteps []values.VTEP) []VXLANTunnelMapWithComment {
	vxlanTunnelMaps := make([]VXLANTunnelMapWithComment, 0)

	for _, vtep := range vteps {
		mapWithComment := VXLANTunnelMapWithComment{
			Comment: "#" + vtep.Comment,
			TunnelMap: map[string]VXLANTunnelMap{
				"vtep|map_" + vtep.VNI + "_" + vtep.VLAN: {
					VLAN: vtep.VLAN,
					VNI:  vtep.VNI,
				},
			},
		}
		vxlanTunnelMaps = append(vxlanTunnelMaps, mapWithComment)
	}

	return vxlanTunnelMaps
}
