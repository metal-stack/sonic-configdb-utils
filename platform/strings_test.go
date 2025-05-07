package platform

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/metal-stack/metal-lib/pkg/testcommon"
)

func Test_cutBetween(t *testing.T) {
	tests := []struct {
		name        string
		s           string
		leftDelim   string
		rightDelim  string
		wantBetween string
		wantFound   bool
	}{
		{
			name:        "no occurence found",
			s:           "(x)",
			leftDelim:   "[",
			rightDelim:  "]",
			wantBetween: "(x)",
			wantFound:   false,
		},
		{
			name:        "only left delimiter found",
			s:           "[x)",
			leftDelim:   "[",
			rightDelim:  "]",
			wantBetween: "[x)",
			wantFound:   false,
		},
		{
			name:        "only right delimiter found",
			s:           "(x]",
			leftDelim:   "[",
			rightDelim:  "]",
			wantBetween: "(x]",
			wantFound:   false,
		},
		{
			name:        "one occurence found",
			s:           "adsf[x]qewr",
			leftDelim:   "[",
			rightDelim:  "]",
			wantBetween: "x",
			wantFound:   true,
		},
		{
			name:        "multiple occurences found",
			s:           "adsf[x]qewr[y]erty",
			leftDelim:   "[",
			rightDelim:  "]",
			wantBetween: "x",
			wantFound:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBetween, gotFound := cutBetween(tt.s, tt.leftDelim, tt.rightDelim)
			if gotBetween != tt.wantBetween {
				t.Errorf("cutBetween() gotBetween = %v, want %v", gotBetween, tt.wantBetween)
			}
			if gotFound != tt.wantFound {
				t.Errorf("cutBetween() gotOk = %v, want %v", gotFound, tt.wantFound)
			}
		})
	}
}

func Test_stringToSpeed(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    int
		wantErr error
	}{
		{
			name:    "no 'G' suffix",
			s:       "100",
			want:    0,
			wantErr: fmt.Errorf("speed must be of format xG, where x is a positive integer"),
		},
		{
			name:    "prefix is not an number",
			s:       "hundredG",
			want:    0,
			wantErr: fmt.Errorf("speed must be of format xG, where x is a positive integer"),
		},
		{
			name:    "speed is negative",
			s:       "-100G",
			want:    0,
			wantErr: fmt.Errorf("speed must be of format xG, where x is a positive integer"),
		},
		{
			name:    "speed is not an integer",
			s:       "10.5G",
			want:    0,
			wantErr: fmt.Errorf("speed must be of format xG, where x is a positive integer"),
		},
		{
			name:    "speed is valid",
			s:       "1G",
			want:    1000,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := stringToSpeed(tt.s)
			if diff := cmp.Diff(tt.wantErr, err, testcommon.ErrorStringComparer()); diff != "" {
				t.Errorf("stringToSpeed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("stringToSpeed() = %v, want %v", got, tt.want)
			}
		})
	}
}
