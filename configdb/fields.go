package configdb

type ACLRule struct {
	EtherType    string `json:"ETHER_TYPE"`
	PacketAction `json:"PACKET_ACTION"`
	Priority     string `json:"PRIORITY"`
	SrcIP        string `json:"SRC_IP"`
}

type ACLTable struct {
	PolicyDesc string   `json:"policy_desc"`
	Ports      []string `json:"ports"`
	Services   []string `json:"services"`
	Stage      string   `json:"stage"`
	Type       string   `json:"type"`
}

type AdminStatus string

const (
	AdminStatusUp   AdminStatus = "UP"
	AdminStatusDown AdminStatus = "DOWN"
)

type AutonegMode string

const (
	AutonegModeOn  AutonegMode = "ON"
	AutonegModeOff AutonegMode = "OFF"
)

type BreakoutConfig struct {
	BreakoutMode `json:"brkout_mode"`
}

type BreakoutMode string

const (
	BreakoutMode1x100G BreakoutMode = "1x100G[40G]"
	BreakoutMode4x25G  BreakoutMode = "4x25G"
	BreakoutMode4x10G  BreakoutMode = "4x10G"
)

type DeviceMetadata struct {
	Localhost Metadata `json:"localhost"`
}

type DockerRoutingConfigMode string

const (
	DockerRoutingConfigModeSeparated DockerRoutingConfigMode = "SEPARATED"
	DockerRoutingConfigModeSplit     DockerRoutingConfigMode = "SPLIT"
	DockerRoutingConfigModeUnified   DockerRoutingConfigMode = "UNIFIED"
)

type Feature struct {
	AutoRestart FeatureMode `json:"auto_restart"`
	State       FeatureMode `json:"state"`
}

type FeatureMode string

const (
	FeatureModeEnabled  FeatureMode = "ENABLED"
	FeatureModeDisabled FeatureMode = "DISABLED"
)

type FECMode string

const (
	FECModeNone FECMode = "NONE"
	FECModeRS   FECMode = "RS"
)

type Interface struct {
	IPv6UseLinkLocalOnly IPv6UseLinkLocalOnlyMode `json:"ipv6_use_link_local_only,omitempty"`
	VRFName              string                   `json:"vrf_name,omitempty"`
}

type IPv6UseLinkLocalOnlyMode string

const (
	IPv6UseLinkLocalOnlyModeEnable IPv6UseLinkLocalOnlyMode = "ENABLE"
)

type LACPKeyMode string

const (
	LACPKeyModeAuto LACPKeyMode = "AUTO"
)

type LLDP struct {
	Global LLDPGlobal `json:"GLOBAL"`
}

type LLDPGlobal struct {
	HelloTime int `json:"hello_time"`
}

type MCLAGDomain struct {
	MCLAGSystemID string `json:"system_mac"`
	PeerIP        string `json:"peer_ip"`
	PeerLink      string `json:"peer_link"`
	SourceIP      string `json:"source_ip"`
}

type MCLAGInterface struct {
	IfType string `json:"if_type"`
}

type MCLAGUniqueIP struct {
	UniqueIP MCLAGUniqueIPMode `json:"unique_ip"`
}

type MCLAGUniqueIPMode string

const (
	MCLAGUniqueIPModeEnable MCLAGUniqueIPMode = "ENABLE"
)

type Metadata struct {
	DockerRoutingConfigMode `json:"docker_routing_config_mode"`
	FRRMgmtFrameworkConfig  bool   `json:"frr_mgmt_framework_config"`
	Hostname                string `json:"hostname"`
	RouterType              `json:"type"`
}

type MgmtInterface struct {
	GWAddr string `json:"gwaddr,omitempty"`
}

type MgmtPort struct {
	AdminStatus `json:"admin_status"`
	Alias       string `json:"alias"`
	Description string `json:"description"`
}

type MgmtVRFConfig struct {
	VRFGlobal `json:"vrf_global"`
}

type NTP struct {
	NTPGlobal `json:"global"`
}

type NTPGlobal struct {
	SrcIntf string `json:"src_intf"`
}

type PacketAction string

const (
	PacketActionDrop   PacketAction = "DROP"
	PacketActionAccept PacketAction = "ACCEPT"
)

type Port struct {
	AdminStatus `json:"admin_status"`
	Alias       string      `json:"alias"`
	Autoneg     AutonegMode `json:"autoneg"`
	FEC         FECMode     `json:"fec"`
	Index       int         `json:"index"`
	Lanes       string      `json:"lanes"`
	MTU         int         `json:"mtu"`
	ParentPort  string      `json:"parent_port"`
	Speed       int         `json:"speed"`
}

type Portchannel struct {
	AdminStatus `json:"admin_status"`
	Fallback    bool        `json:"fallback"`
	LACPKey     LACPKeyMode `json:"lacp_key"`
	MinLinks    int         `json:"min_links"`
	MixSpeed    bool        `json:"mix_speed"`
	MTU         int         `json:"mtu"`
}

type RouterType string

const (
	RouterTypeDualToR    RouterType = "DUAL_TOR"
	RouterTypeLeafRouter RouterType = "LEAF_ROUTER"
	RouterTypeToRRouter  RouterType = "TOR_ROUTER"
)

type SAG struct {
	SAGGlobal `json:"GLOBAL"`
}

type SAGGlobal struct {
	GatewayMAC string `json:"gateway_mac"`
}

type TaggingMode string

const (
	TaggingModeTagged   TaggingMode = "tagged"
	TaggingModeUntagged TaggingMode = "untagged"
)

type VLAN struct {
	DHCPServers []string `json:"dhcp_servers"`
	VLANID      string   `json:"vlanid"`
}

type VLANInterface struct {
	StaticAnycastGateway bool   `json:"static_anycast_gateway"`
	VRFName              string `json:"vrf_name"`
}

type VLANMember struct {
	TaggingMode `json:"tagging_mode"`
}

type VRF struct {
	VNI string `json:"vni"`
}

type VRFGlobal struct {
	MgmtVRFEnabled bool `json:"mgmtVrfEnabled"`
}

type VXLANEVPN struct {
	VXLANEVPNNVO `json:"nvo"`
}

type VXLANEVPNNVO struct {
	SourceVTEP string `json:"source_vtep"`
}

type VXLANTunnel struct {
	SrcIP string `json:"src_ip"`
}

type VXLANTunnelMap struct {
	VLAN string `json:"vlan"`
	VNI  string `json:"vni"`
}

type VXLANTunnelMapWithComment struct {
	Comment   string `json:"comment"`
	TunnelMap map[string]VXLANTunnelMap
}
