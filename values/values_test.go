package values

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_UnmarshalValues_acls(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []ACL
		wantErr string
	}{
		{
			name: "known fields",
			input: `
acls:
  - name: PXE_IN
    description: pxe
    ports: [Vlan4000]
    rules:
      - action: drop
        dst: 10.0.0.0/16
      - action: forward
`,
			want: []ACL{{
				Name:        "PXE_IN",
				Description: "pxe",
				Ports:       []string{"Vlan4000"},
				Rules:       []ACLRule{{Action: "drop", Dst: "10.0.0.0/16"}, {Action: "forward"}},
			}},
		},
		{
			name: "unknown rule field",
			input: `
acls:
  - name: PXE_IN
    rules:
      - action: drop
        dst_ip: 10.0.0.0/16
`,
			wantErr: `line 6: field dst_ip not found in type values.ACLRule`,
		},
		{
			name: "fields from a merge key",
			input: `
acls:
  - &pxe
    name: PXE_IN
    ports: [Vlan4000]
    rules:
      - action: drop
  - <<: *pxe
    name: PXE_OUT
    stage: egress
`,
			want: []ACL{
				{
					Name:  "PXE_IN",
					Ports: []string{"Vlan4000"},
					Rules: []ACLRule{{Action: "drop"}},
				},
				{
					Name:  "PXE_OUT",
					Ports: []string{"Vlan4000"},
					Rules: []ACLRule{{Action: "drop"}},
					Stage: "egress",
				},
			},
		},
		{
			name: "helper key for an anchor",
			input: `
common: &common
  ports: [Vlan4000]
acls:
  - <<: *common
    name: PXE_IN
`,
			wantErr: `line 2: field common not found in type values.Values`,
		},
		{
			name: "unknown acl field",
			input: `
acls:
  - name: PXE_IN
    port: [Vlan4000]
`,
			wantErr: `line 4: field port not found in type values.ACL`,
		},
		{
			name:  "empty input",
			input: ``,
		},
		{
			name: "misspelled acls key",
			input: `
acl:
  - name: PXE_IN
    ports: [Vlan4000]
    rules:
      - action: drop
`,
			wantErr: `line 2: field acl not found`,
		},
		{
			name: "frr keys that metal-roles passes through are accepted",
			input: `
interconnects:
  external:
    vrf: Vrf46
    vni: 46
    peer_group: EXTERNAL
    neighbor_ip: 192.168.255.2
    remote_as: 4200000000
    keep_private_as: true
    announcements:
      - redistribute connected
    routemap_in:
      name: ALLOW-EXTERNAL-IN
      entries:
        - match ip address prefix-list EXTERNAL_IN
    prefixlists:
      - ip prefix-list EXTERNAL_IN seq 10 permit 10.64.64.0/22 ge 32
vlans:
  - id: "4000"
    vrrp:
      group: "1"
      priority: 66
      ip: 10.1.1.1/24
`,
		},
		{
			name: "top-level keys from a merge key",
			input: `
<<: {hostname: leaf01}
`,
		},
		{
			name: "misspelled top-level key from a merge key",
			input: `
<<: {hostnme: leaf01}
`,
			wantErr: `line 2: field hostnme not found in type values.Values`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalValues([]byte(tt.input))
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("UnmarshalValues() error = %v, want it to contain %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("UnmarshalValues() unexpected error = %v", err)
			}
			if diff := cmp.Diff(tt.want, got.ACLs); diff != "" {
				t.Errorf("UnmarshalValues() acls diff = %s", diff)
			}
		})
	}
}

func Test_UnmarshalValues_typeErrorHasNoUnknownFieldHint(t *testing.T) {
	_, err := UnmarshalValues([]byte("lldp_hello_time: ten\n"))
	if err == nil || !strings.Contains(err.Error(), "cannot unmarshal") {
		t.Fatalf("UnmarshalValues() error = %v, want a type error", err)
	}
	if strings.Contains(err.Error(), "fix or remove the fields") {
		t.Errorf("UnmarshalValues() error = %v, a type error must not ask to remove fields", err)
	}
}
