package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseImageSwapPolicy(t *testing.T) {
	type args struct {
		p string
	}
	tests := []struct {
		name    string
		args    args
		want    ImageSwapPolicy
		wantErr bool
	}{
		{
			name: "always",
			args: args{p: "always"},
			want: ImageSwapPolicyAlways,
		},
		{
			name: "exists",
			args: args{p: "exists"},
			want: ImageSwapPolicyExists,
		},
		{
			name:    "random-non-existent",
			args:    args{p: "random-non-existent"},
			want:    ImageSwapPolicyExists,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseImageSwapPolicy(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseImageSwapPolicy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseImageSwapPolicy() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseImageCopyPolicy(t *testing.T) {
	type args struct {
		p string
	}
	tests := []struct {
		name    string
		args    args
		want    ImageCopyPolicy
		wantErr bool
	}{
		{
			name: "delayed",
			args: args{p: "delayed"},
			want: ImageCopyPolicyDelayed,
		},
		{
			name: "immediate",
			args: args{p: "immediate"},
			want: ImageCopyPolicyImmediate,
		},
		{
			name: "force",
			args: args{p: "force"},
			want: ImageCopyPolicyForce,
		},
		{
			name: "none",
			args: args{p: "none"},
			want: ImageCopyPolicyNone,
		},
		{
			name:    "random-non-existent",
			args:    args{p: "random-non-existent"},
			want:    ImageCopyPolicyDelayed,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseImageCopyPolicy(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseImageCopyPolicy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseImageCopyPolicy() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseAwsRegistry(t *testing.T) {

	registry, err := ParseRegistry("aws")
	assert.Nil(t, err)
	assert.Equal(t, "aws", registry.String())
}

func TestParseGcpRegistry(t *testing.T) {

	registry, err := ParseRegistry("gcp")
	assert.Nil(t, err)
	assert.Equal(t, "gcp", registry.String())
}
func TestParseGenericRegistry(t *testing.T) {

	registry, err := ParseRegistry("generic")
	assert.Nil(t, err)
	assert.Equal(t, "generic", registry.String())
}

func TestParseUnknownRegistry(t *testing.T) {

	registry, err := ParseRegistry("not_known")
	assert.NotNil(t, err)
	assert.Equal(t, "unknown", registry.String())
	assert.Equal(t, "unknown target registry string: 'not_known', defaulting to unknown", err.Error())
}
