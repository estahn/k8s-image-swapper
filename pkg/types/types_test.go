package types

import "testing"

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
