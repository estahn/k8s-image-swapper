package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func TestImageSwapperMutator(t *testing.T) {
	tests := []struct {
		name   string
		pod    *corev1.Pod
		labels map[string]string
		expPod *corev1.Pod
		expErr bool
	}{
		{
			name: "Prefix docker hub images with host docker.io.",
			pod: &corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image: "nginx:latest",
						},
					},
				},
			},
			expPod: &corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image: "foobar.com/docker.io/nginx:latest",
						},
					},
				},
			},
		},
		{
			name: "Don't mutate if targetRegistry host is target targetRegistry.",
			pod: &corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image: "foobar.com/docker.io/nginx:latest",
						},
					},
				},
			},
			expPod: &corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image: "foobar.com/docker.io/nginx:latest",
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			pl := NewImageSwapper("foobar.com")

			gotPod := test.pod
			_, err := pl.Mutate(context.TODO(), gotPod)

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				assert.Equal(test.expPod, gotPod)
			}
		})
	}

}
//
//func TestAnnotatePodMutator2(t *testing.T) {
//	tests := []struct {
//		name   string
//		pod    *corev1.Pod
//		labels map[string]string
//		expPod *corev1.Pod
//		expErr bool
//	}{
//		{
//			name: "Mutating a pod without labels should set the labels correctly.",
//			pod: &corev1.Pod{
//				ObjectMeta: metav1.ObjectMeta{
//					Name: "test",
//				},
//			},
//			labels: map[string]string{"bruce": "wayne", "peter": "parker"},
//			expPod: &corev1.Pod{
//				ObjectMeta: metav1.ObjectMeta{
//					Name:   "test",
//					Labels: map[string]string{"bruce": "wayne", "peter": "parker"},
//				},
//			},
//		},
//		{
//			name: "Mutating a pod with labels should aggregate and replace the labels with the existing ones.",
//			pod: &corev1.Pod{
//				ObjectMeta: metav1.ObjectMeta{
//					Name:   "test",
//					Labels: map[string]string{"bruce": "banner", "tony": "stark"},
//				},
//			},
//			labels: map[string]string{"bruce": "wayne", "peter": "parker"},
//			expPod: &corev1.Pod{
//				ObjectMeta: metav1.ObjectMeta{
//					Name:   "test",
//					Labels: map[string]string{"bruce": "wayne", "peter": "parker", "tony": "stark"},
//				},
//			},
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			assert := assert.New(t)
//
//			pl := mutatortesting.NewPodLabeler(test.labels)
//			gotPod := test.pod
//			_, err := pl.Mutate(context.TODO(), gotPod)
//
//			if test.expErr {
//				assert.Error(err)
//			} else if assert.NoError(err) {
//				// Check the expected pod.
//				assert.Equal(test.expPod, gotPod)
//			}
//		})
//	}
//
//}

//func TestRegistryHost(t *testing.T) {
//	assert.Equal(t, "", registryDomain("nginx:latest"))
//	assert.Equal(t, "docker.io", registryDomain("docker.io/nginx:latest"))
//}
