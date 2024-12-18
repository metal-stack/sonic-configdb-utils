package configdb

type ACLRule struct {
	EtherType    string `json:"ETHER_TYPE,omitempty"`
	PacketAction `json:"PACKET_ACTION,omitempty"`
	Priority     string `json:"PRIORITY,omitempty"`
	SrcIP        string `json:"SRC_IP,omitempty"`
}

type ACLTable struct {
	PolicyDesc string   `json:"policy_desc,omitempty"`
	Ports      []string `json:"ports,omitempty"`
	Services   []string `json:"services,omitempty"`
	Stage      string   `json:"stage,omitempty"`
	Type       string   `json:"type,omitempty"`
}

type AdminStatus string

const (
	AdminStatusUp   AdminStatus = "up"
	AdminStatusDown AdminStatus = "down"
)

type AutonegMode string

const (
	AutonegModeOn  AutonegMode = "on"
	AutonegModeOff AutonegMode = "off"
)

type BreakoutConfig struct {
	BreakoutMode `json:"brkout_mode,omitempty"`
}

type BreakoutMode string

const (
	BreakoutMode1x100G BreakoutMode = "1x100G[40G]"
	BreakoutMode4x25G  BreakoutMode = "4x25G"
	BreakoutMode4x10G  BreakoutMode = "4x10G"
)

type DeviceMetadata struct {
	Localhost Metadata `json:"localhost,omitempty"`
}

type DockerRoutingConfigMode string

const (
	DockerRoutingConfigModeSeparated DockerRoutingConfigMode = "separated"
	DockerRoutingConfigModeSplit     DockerRoutingConfigMode = "split"
	DockerRoutingConfigModeUnified   DockerRoutingConfigMode = "unified"
)

type Feature struct {
	AutoRestart FeatureMode `json:"auto_restart,omitempty"`
	State       FeatureMode `json:"state,omitempty"`
}

type FeatureMode string

const (
	FeatureModeEnabled  FeatureMode = "enabled"
	FeatureModeDisabled FeatureMode = "disabled"
)

type FECMode string

const (
	FECModeNone FECMode = "none"
	FECModeRS   FECMode = "rs"
)

type Interface struct {
	IPv6UseLinkLocalOnly IPv6UseLinkLocalOnlyMode `json:"ipv6_use_link_local_only,omitempty"`
	VRFName              string                   `json:"vrf_name,omitempty"`
}

type IPv6UseLinkLocalOnlyMode string

const (
	IPv6UseLinkLocalOnlyModeEnable IPv6UseLinkLocalOnlyMode = "enable"
)

type LACPKeyMode string

const (
	LACPKeyModeAuto LACPKeyMode = "auto"
)

type LLDP struct {
	Global LLDPGlobal `json:"GLOBAL,omitempty"`
}

type LLDPGlobal struct {
	HelloTime int `json:"hello_time,omitempty"`
}

type MCLAGDomain struct {
	MCLAGSystemID string `json:"system_mac,omitempty"`
	PeerIP        string `json:"peer_ip,omitempty"`
	PeerLink      string `json:"peer_link,omitempty"`
	SourceIP      string `json:"source_ip,omitempty"`
}

type MCLAGInterface struct {
	IfType string `json:"if_type,omitempty"`
}

type MCLAGUniqueIP struct {
	UniqueIP MCLAGUniqueIPMode `json:"unique_ip,omitempty"`
}

type MCLAGUniqueIPMode string

const (
	MCLAGUniqueIPModeEnable MCLAGUniqueIPMode = "enable"
)

type Metadata struct {
	DockerRoutingConfigMode `json:"docker_routing_config_mode,omitempty"`
	FRRMgmtFrameworkConfig  bool   `json:"frr_mgmt_framework_config,omitempty"`
	Hostname                string `json:"hostname,omitempty"`
	RouterType              `json:"type,omitempty"`
}

type MgmtInterface struct {
	GWAddr string `json:"gwaddr,omitempty"`
}

type MgmtPort struct {
	AdminStatus `json:"admin_status,omitempty"`
	Alias       string `json:"alias,omitempty"`
	Description string `json:"description,omitempty"`
}

type MgmtVRFConfig struct {
	VRFGlobal `json:"vrf_global,omitempty"`
}

type NTP struct {
	NTPGlobal `json:"global,omitempty"`
}

type NTPGlobal struct {
	SrcIntf string `json:"src_intf,omitempty"`
}

type PacketAction string

const (
	PacketActionDrop   PacketAction = "drop"
	PacketActionAccept PacketAction = "accept"
)

type Port struct {
	AdminStatus `json:"admin_status,omitempty"`
	Alias       string      `json:"alias,omitempty"`
	Autoneg     AutonegMode `json:"autoneg,omitempty"`
	FEC         FECMode     `json:"fec,omitempty"`
	Index       int         `json:"index,omitempty"`
	Lanes       string      `json:"lanes,omitempty"`
	MTU         int         `json:"mtu,omitempty"`
	ParentPort  string      `json:"parent_port,omitempty"`
	Speed       int         `json:"speed,omitempty"`
}

type Portchannel struct {
	AdminStatus `json:"admin_status,omitempty"`
	Fallback    bool        `json:"fallback,omitempty"`
	LACPKey     LACPKeyMode `json:"lacp_key,omitempty"`
	MinLinks    int         `json:"min_links,omitempty"`
	MixSpeed    bool        `json:"mix_speed,omitempty"`
	MTU         int         `json:"mtu,omitempty"`
}

type RouterType string

const (
	RouterTypeDualToR    RouterType = "DualToR"
	RouterTypeLeafRouter RouterType = "LeafRouter"
	RouterTypeToRRouter  RouterType = "ToRRouter"
)

type SAG struct {
	SAGGlobal `json:"GLOBAL,omitempty"`
}

type SAGGlobal struct {
	GatewayMAC string `json:"gateway_mac,omitempty"`
}

type TaggingMode string

const (
	TaggingModeTagged   TaggingMode = "tagged"
	TaggingModeUntagged TaggingMode = "untagged"
)

type VLAN struct {
	DHCPServers []string `json:"dhcp_servers,omitempty"`
	VLANID      string   `json:"vlanid,omitempty"`
}

type VLANInterface struct {
	StaticAnycastGateway bool   `json:"static_anycast_gateway,omitempty"`
	VRFName              string `json:"vrf_name,omitempty"`
}

type VLANMember struct {
	TaggingMode `json:"tagging_mode,omitempty"`
}

type VRF struct {
	VNI string `json:"vni,omitempty"`
}

type VRFGlobal struct {
	MgmtVRFEnabled bool `json:"mgmtVrfEnabled,omitempty"`
}

type VXLANEVPN struct {
	VXLANEVPNNVO `json:"nvo,omitempty"`
}

type VXLANEVPNNVO struct {
	SourceVTEP string `json:"source_vtep,omitempty"`
}

type VXLANTunnel struct {
	SrcIP string `json:"src_ip,omitempty"`
}

type VXLANTunnelMap struct {
	VLAN string `json:"vlan,omitempty"`
	VNI  string `json:"vni,omitempty"`
}

type VXLANTunnelMapWithComment struct {
	Comment   string `json:"comment,omitempty"`
	TunnelMap map[string]VXLANTunnelMap
}
