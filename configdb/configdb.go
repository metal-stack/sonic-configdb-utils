package configdb

import (
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"strconv"

	p "github.com/metal-stack/sonic-configdb-utils/platform"
	"github.com/metal-stack/sonic-configdb-utils/values"
)

type ConfigDB struct {
	ACLRules           map[string]ACLRule        `json:"ACL_RULE,omitempty"`
	ACLTables          map[string]ACLTable       `json:"ACL_TABLE,omitempty"`
	Breakouts          map[string]BreakoutConfig `json:"BREAKOUT_CFG,omitempty"`
	DeviceMetadata     `json:"DEVICE_METADATA"`
	Features           map[string]Feature   `json:"FEATURE,omitempty"`
	Interfaces         map[string]Interface `json:"INTERFACE,omitempty"`
	LLDP               `json:"LLDP"`
	LoopbackInterface  map[string]struct{}       `json:"LOOPBACK_INTERFACE,omitempty"`
	MCLAGDomains       map[string]MCLAGDomain    `json:"MCLAG_DOMAIN,omitempty"`
	MCLAGInterfaces    map[string]MCLAGInterface `json:"MCLAG_INTERFACE,omitempty"`
	MCLAGUniqueIPs     map[string]MCLAGUniqueIP  `json:"MCLAG_UNIQUE_IP,omitempty"`
	MgmtInterfaces     map[string]MgmtInterface  `json:"MGMT_INTERFACE,omitempty"`
	MgmtPorts          map[string]MgmtPort       `json:"MGMT_PORT,omitempty"`
	MgmtVRFConfig      `json:"MGMT_VRF_CONFIG"`
	NTP                `json:"NTP"`
	NTPServers         map[string]struct{}      `json:"NTP_SERVER,omitempty"`
	Ports              map[string]Port          `json:"PORT,omitempty"`
	PortChannels       map[string]PortChannel   `json:"PORTCHANNEL,omitempty"`
	PortChannelMembers map[string]struct{}      `json:"PORTCHANNEL_MEMBER,omitempty"`
	SAG                *SAG                     `json:"SAG,omitempty"`
	VLANs              map[string]VLAN          `json:"VLAN,omitempty"`
	VLANInterfaces     map[string]VLANInterface `json:"VLAN_INTERFACE,omitempty"`
	VLANMembers        map[string]VLANMember    `json:"VLAN_MEMBER,omitempty"`
	VRFs               map[string]VRF           `json:"VRF,omitempty"`
	VXLANEVPN          `json:"VXLAN_EVPN_NVO"`
	VXLANTunnels       map[string]VXLANTunnel `json:"VXLAN_TUNNEL,omitempty"`
	VXLANTunnelMap     `json:"VXLAN_TUNNEL_MAP,omitempty"`
}

func GenerateConfigDB(input *values.Values, platform *p.Platform, currentDeviceMetadata DeviceMetadata) (*ConfigDB, error) {
	ports, breakouts, err := getPortsAndBreakouts(input.Ports, input.Breakouts, input.PortsDefaultFEC, input.PortsDefaultMTU, platform)
	if err != nil {
		return nil, err
	}

	deviceMetadata, err := getDeviceMetadata(input, currentDeviceMetadata)
	if err != nil {
		return nil, err
	}

	rules, tables := getACLRulesAndTables(input.SSHSourceranges)

	configdb := ConfigDB{
		ACLRules:       rules,
		ACLTables:      tables,
		Breakouts:      breakouts,
		DeviceMetadata: *deviceMetadata,
		Features: map[string]Feature{
			"dhcp_relay": {
				AutoRestart: FeatureModeEnabled,
				State:       FeatureModeEnabled,
			},
		},
		Interfaces: getInterfaces(input.Ports, input.BGPPorts),
		LLDP: LLDP{
			Global: LLDPGlobal{
				HelloTime: fmt.Sprintf("%d", input.LLDPHelloTime),
			},
		},
		LoopbackInterface: map[string]struct{}{
			"Loopback0": {},
			"Loopback0|" + input.LoopbackAddress + "/32": {},
		},
		MCLAGDomains:    getMCLAGDomains(input.MCLAG),
		MCLAGInterfaces: getMCLAGInterfaces(input.MCLAG),
		MCLAGUniqueIPs:  getMCLAGUniqueIPs(input.MCLAG),
		MgmtInterfaces:  getMgmtInterfaces(input.MgmtIfIP, input.MgmtIfGateway),
		MgmtPorts: map[string]MgmtPort{
			"eth0": {
				AdminStatus: AdminStatusUp,
				Alias:       "eth0",
				Description: "Management Port",
			},
		},
		MgmtVRFConfig: MgmtVRFConfig{
			VRFGlobal: VRFGlobal{
				MgmtVRFEnabled: strconv.FormatBool(input.MgmtVRF),
			},
		},
		NTP: NTP{
			NTPGlobal: NTPGlobal{
				SrcIntf: "eth0",
			},
		},
		NTPServers:         getNTPServers(input.NTPServers),
		Ports:              ports,
		PortChannels:       getPortChannels(input.PortChannels, input.PortChannelsDefaultMTU),
		PortChannelMembers: getPortChannelMembers(input.PortChannels),
		SAG:                getSAG(input.SAG),
		VLANs:              getVLANs(input.VLANs),
		VLANInterfaces:     getVLANInterfaces(input.VLANs),
		VLANMembers:        getVLANMembers(input.VLANs, input.VLANMembers),
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
		VXLANTunnelMap: getVXLANTunnelMap(input.VTEPs),
	}
	return &configdb, nil
}

func UnmarshalConfigDB(in []byte) (*ConfigDB, error) {
	var configDB ConfigDB
	err := json.Unmarshal(in, &configDB)
	if err != nil {
		return nil, err
	}

	return &configDB, nil
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
			Services:   []string{"SSH"},
			Stage:      "ingress",
			Type:       "CTRLPLANE",
		},
		"ALLOW_NTP": {
			PolicyDesc: "Allow NTP",
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

func getDeviceMetadata(input *values.Values, currentMetadata DeviceMetadata) (*DeviceMetadata, error) {
	hint := "remove current config_db.json and run `config-setup factory` to generate an initial config_db.json with all the necessary information"

	if currentMetadata.Localhost.Platform == "" {
		return nil, fmt.Errorf("missing platform from current device metadata\n%s", hint)
	}

	if currentMetadata.Localhost.HWSKU == "" {
		return nil, fmt.Errorf("missing hwsku from current device metadata\n%s", hint)
	}

	if currentMetadata.Localhost.MAC == "" {
		return nil, fmt.Errorf("missing mac from current device metadata\n%s", hint)
	}

	return &DeviceMetadata{
		Localhost: Metadata{
			DockerRoutingConfigMode: DockerRoutingConfigMode(input.DockerRoutingConfigMode),
			FRRMgmtFrameworkConfig:  strconv.FormatBool(input.FRRMgmtFrameworkConfig),
			Hostname:                input.Hostname,
			HWSKU:                   currentMetadata.Localhost.HWSKU,
			MAC:                     currentMetadata.Localhost.MAC,
			Platform:                currentMetadata.Localhost.Platform,
			RouterType:              "LeafRouter",
		},
	}, nil
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

func getMCLAGDomains(mclag *values.MCLAG) map[string]MCLAGDomain {
	if mclag == nil {
		return nil
	}

	return map[string]MCLAGDomain{
		"1": {
			MCLAGSystemID: mclag.SystemMAC,
			PeerIP:        mclag.PeerIP,
			PeerLink:      mclag.PeerLink,
			SourceIP:      mclag.SourceIP,
		},
	}
}

func getMCLAGInterfaces(mclag *values.MCLAG) map[string]MCLAGInterface {
	if mclag == nil {
		return nil
	}

	mclagInterfaces := make(map[string]MCLAGInterface)

	for _, channel := range mclag.MemberPortChannels {
		mclagInterfaces["1|PortChannel"+channel] = MCLAGInterface{
			IfType: "PortChannel",
		}
	}

	return mclagInterfaces
}

func getMCLAGUniqueIPs(mclag *values.MCLAG) map[string]MCLAGUniqueIP {
	if mclag == nil {
		return nil
	}

	return map[string]MCLAGUniqueIP{
		"Vlan" + mclag.KeepaliveVLAN: {
			UniqueIP: MCLAGUniqueIPModeEnable,
		},
	}
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

func getPortChannels(portChannels []values.PortChannel, defaultPortChannelMTU int) map[string]PortChannel {
	configPortChannels := make(map[string]PortChannel)

	for _, pc := range portChannels {
		mtu := defaultMTU
		if defaultPortChannelMTU != 0 {
			mtu = defaultPortChannelMTU
		}
		if pc.MTU != 0 {
			mtu = pc.MTU
		}

		configPortChannels["PortChannel"+pc.Number] = PortChannel{
			AdminStatus: defaultAdminStatus,
			Fallback:    "true",
			FastRate:    "false",
			LACPKey:     LACPKeyModeAuto,
			MinLinks:    "1",
			MixSpeed:    "false",
			MTU:         fmt.Sprintf("%d", mtu),
		}
	}

	return configPortChannels
}

func getPortChannelMembers(portchannels []values.PortChannel) map[string]struct{} {
	portchannelMembers := make(map[string]struct{})

	for _, pc := range portchannels {
		for _, member := range pc.Members {
			portchannelMembers["PortChannel"+pc.Number+"|"+member] = struct{}{}
		}
	}

	return portchannelMembers
}

func getPortsAndBreakouts(ports []values.Port, breakouts map[string]string, defaultFECMode values.FECMode, defaultMTU int, platform *p.Platform) (map[string]Port, map[string]BreakoutConfig, error) {
	configPorts := make(map[string]Port)
	configBreakouts := make(map[string]BreakoutConfig)

	defaultBreakouts := platform.GetDefaultBreakoutConfig()
	for portName, breakout := range defaultBreakouts {
		breakoutPorts, err := getPortsFromBreakout(portName, breakout, defaultFECMode, defaultMTU, platform)
		if err != nil {
			return nil, nil, err
		}
		maps.Copy(configPorts, breakoutPorts)

		configBreakouts[portName] = BreakoutConfig{
			BreakoutMode: breakout,
		}
	}

	for portName, breakout := range breakouts {
		breakoutPorts, err := getPortsFromBreakout(portName, breakout, defaultFECMode, defaultMTU, platform)
		if err != nil {
			return nil, nil, err
		}
		maps.Copy(configPorts, breakoutPorts)

		configBreakouts[portName] = BreakoutConfig{
			BreakoutMode: breakout,
		}
	}

	for _, port := range ports {
		configPort, ok := configPorts[port.Name]
		if !ok {
			return nil, nil, fmt.Errorf("invalid port name %s; if you think it should be available please check your breakout configuration", port.Name)
		}

		// error is ignored because the breakout gets parsed before and any error would have made the function return before reaching this line
		speedOptions, _ := p.ParseSpeedOptions(configPort.parentBreakout)

		if port.Speed != 0 && !slices.Contains(speedOptions[:], port.Speed) {
			return nil, nil, fmt.Errorf("invalid speed %d for port %s; current breakout configuration %s only allows speed options %v", port.Speed, port.Name, configPort.parentBreakout, speedOptions)
		}

		if port.Speed != 0 {
			configPort.Speed = fmt.Sprintf("%d", port.Speed)
		}

		if port.FECMode != "" && string(port.FECMode) != string(configPort.FEC) {
			configPort.FEC = FECMode(port.FECMode)
		}
		if port.MTU != 0 && fmt.Sprintf("%d", port.MTU) != configPort.MTU {
			configPort.MTU = fmt.Sprintf("%d", port.MTU)
		}
		configPorts[port.Name] = configPort
	}

	return configPorts, configBreakouts, nil
}

func getSAG(sag values.SAG) *SAG {
	if sag.MAC == "" {
		return nil
	}

	return &SAG{
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
				StaticAnycastGateway: strconv.FormatBool(vlan.SAG),
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

func getVLANMembers(vlans []values.VLAN, addVlanMembers bool) map[string]VLANMember {
	if !addVlanMembers {
		return nil
	}

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

func getVXLANTunnelMap(vteps []values.VTEP) VXLANTunnelMap {
	vxlanTunnelMap := make(VXLANTunnelMap)

	for _, vtep := range vteps {
		vxlanTunnelMap["vtep|map_"+vtep.VNI+"_"+vtep.VLAN] = VXLANTunnelMapEntry{
			VLAN: vtep.VLAN,
			VNI:  vtep.VNI,
		}
	}

	return vxlanTunnelMap
}
