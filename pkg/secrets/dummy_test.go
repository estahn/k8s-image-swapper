package secrets

import (
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestDummyImagePullSecretsProvider_GetImagePullSecrets(t *testing.T) {
	type args struct {
		pod *corev1.Pod
	}
	tests := []struct {
		name    string
		args    args
		want    *ImagePullSecretsResult
		wantErr bool
	}{
		{
			name: "default",
			args: args{
				pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "test-ns",
						Name:      "my-pod",
					},
					Spec: corev1.PodSpec{
						ServiceAccountName: "my-service-account",
						ImagePullSecrets: []corev1.LocalObjectReference{
							{Name: "my-pod-secret"},
						},
					},
				},
			},
			want:    NewImagePullSecretsResult(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &DummyImagePullSecretsProvider{}
			got, err := p.GetImagePullSecrets(tt.args.pod)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetImagePullSecrets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetImagePullSecrets() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDummyImagePullSecretsProvider(t *testing.T) {
	tests := []struct {
		name string
		want ImagePullSecretsProvider
	}{
		{
			name: "default",
			want: &DummyImagePullSecretsProvider{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDummyImagePullSecretsProvider(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDummyImagePullSecretsProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}
