package neta

import (
	"testing"
)

func TestIsCompatibleVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    bool
	}{
		{
			name:    "same version",
			version: "1.0",
			want:    true,
		},
		{
			name:    "same major, different minor",
			version: "1.1",
			want:    true,
		},
		{
			name:    "same major, higher minor",
			version: "1.5",
			want:    true,
		},
		{
			name:    "different major",
			version: "2.0",
			want:    false,
		},
		{
			name:    "different major higher",
			version: "3.0",
			want:    false,
		},
		{
			name:    "empty version",
			version: "",
			want:    false,
		},
		{
			name:    "invalid format - no dot",
			version: "abc",
			want:    false,
		},
		{
			name:    "invalid format - letters",
			version: "a.b",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isCompatibleVersion(tt.version)
			if got != tt.want {
				t.Errorf("isCompatibleVersion(%q) = %v, want %v", tt.version, got, tt.want)
			}
		})
	}
}

func TestValidateVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		wantErr bool
	}{
		{
			name:    "valid version 1.0",
			version: "1.0",
			wantErr: false,
		},
		{
			name:    "valid version 1.1",
			version: "1.1",
			wantErr: false,
		},
		{
			name:    "missing version",
			version: "",
			wantErr: true,
		},
		{
			name:    "incompatible version 2.0",
			version: "2.0",
			wantErr: true,
		},
		{
			name:    "incompatible version 0.9",
			version: "0.9",
			wantErr: true,
		},
		{
			name:    "invalid format",
			version: "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateVersion(tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateVersion(%q) error = %v, wantErr %v", tt.version, err, tt.wantErr)
			}
		})
	}
}

func TestDefinition_IsVersionCompatible(t *testing.T) {
	tests := []struct {
		name string
		def  Definition
		want bool
	}{
		{
			name: "compatible version 1.0",
			def:  Definition{Version: "1.0", Type: "http"},
			want: true,
		},
		{
			name: "compatible version 1.1",
			def:  Definition{Version: "1.1", Type: "http"},
			want: true,
		},
		{
			name: "incompatible version 2.0",
			def:  Definition{Version: "2.0", Type: "http"},
			want: false,
		},
		{
			name: "empty version",
			def:  Definition{Version: "", Type: "http"},
			want: false,
		},
		{
			name: "invalid version",
			def:  Definition{Version: "bad", Type: "http"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.def.IsVersionCompatible()
			if got != tt.want {
				t.Errorf("Definition.IsVersionCompatible() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseMajorVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    int
		wantErr bool
	}{
		{
			name:    "version 1.0",
			version: "1.0",
			want:    1,
			wantErr: false,
		},
		{
			name:    "version 2.5",
			version: "2.5",
			want:    2,
			wantErr: false,
		},
		{
			name:    "version with patch",
			version: "1.2.3",
			want:    1,
			wantErr: false,
		},
		{
			name:    "empty version",
			version: "",
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid format",
			version: "abc",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseMajorVersion(tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseMajorVersion(%q) error = %v, wantErr %v", tt.version, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseMajorVersion(%q) = %v, want %v", tt.version, got, tt.want)
			}
		})
	}
}
