package webhook

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	"github.com/davecgh/go-spew/spew"
	"github.com/estahn/k8s-image-swapper/pkg/config"
	"github.com/estahn/k8s-image-swapper/pkg/registry"
	"github.com/estahn/k8s-image-swapper/pkg/secrets"
	"github.com/estahn/k8s-image-swapper/pkg/types"
	"github.com/slok/kubewebhook/v2/pkg/model"
	"github.com/stretchr/testify/assert"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//func TestImageSwapperMutator(t *testing.T) {
//	tests := []struct {
//		name   string
//		pod    *corev1.Pod
//		labels map[string]string
//		expPod *corev1.Pod
//		expErr bool
//	}{
//		{
//			name: "Prefix docker hub images with host docker.io.",
//			pod: &corev1.Pod{
//				Spec: corev1.PodSpec{
//					Containers: []corev1.Container{
//						{
//							Image: "nginx:latest",
//						},
//					},
//				},
//			},
//			expPod: &corev1.Pod{
//				Spec: corev1.PodSpec{
//					Containers: []corev1.Container{
//						{
//							Image: "foobar.com/docker.io/nginx:latest",
//						},
//					},
//				},
//			},
//		},
//		{
//			name: "Don't mutate if targetRegistry host is target targetRegistry.",
//			pod: &corev1.Pod{
//				Spec: corev1.PodSpec{
//					Containers: []corev1.Container{
//						{
//							Image: "foobar.com/docker.io/nginx:latest",
//						},
//					},
//				},
//			},
//			expPod: &corev1.Pod{
//				Spec: corev1.PodSpec{
//					Containers: []corev1.Container{
//						{
//							Image: "foobar.com/docker.io/nginx:latest",
//						},
//					},
//				},
//			},
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			assert := assert.New(t)
//
//			pl := NewImageSwapper("foobar.com")
//
//			gotPod := test.pod
//			_, err := pl.Mutate(context.TODO(), gotPod)
//
//			if test.expErr {
//				assert.Error(err)
//			} else if assert.NoError(err) {
//				assert.Equal(test.expPod, gotPod)
//			}
//		})
//	}
//
//}
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

func TestFilterMatch(t *testing.T) {
	filterContext := FilterContext{
		Obj: &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "kube-system",
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  "nginx",
						Image: "nginx:latest",
					},
				},
			},
		},
		Container: corev1.Container{
			Name:  "nginx",
			Image: "nginx:latest",
		},
	}

	assert.True(t, filterMatch(filterContext, []config.JMESPathFilter{{JMESPath: "obj.metadata.namespace == 'kube-system'"}}))
	assert.False(t, filterMatch(filterContext, []config.JMESPathFilter{{JMESPath: "obj.metadata.namespace != 'kube-system'"}}))
	assert.False(t, filterMatch(filterContext, []config.JMESPathFilter{{JMESPath: "obj"}}))
	assert.True(t, filterMatch(filterContext, []config.JMESPathFilter{{JMESPath: "container.name == 'nginx'"}}))
	// false syntax test
	assert.False(t, filterMatch(filterContext, []config.JMESPathFilter{{JMESPath: "."}}))
	// non-boolean value
	assert.False(t, filterMatch(filterContext, []config.JMESPathFilter{{JMESPath: "obj"}}))
	assert.False(t, filterMatch(filterContext, []config.JMESPathFilter{{JMESPath: "contains(container.image, '.dkr.ecr.') && contains(container.image, '.amazonaws.com')"}}))
}

type mockECRClient struct {
	ecriface.ECRAPI
}

func (m *mockECRClient) CreateRepository(*ecr.CreateRepositoryInput) (*ecr.CreateRepositoryOutput, error) {
	return &ecr.CreateRepositoryOutput{}, nil
}

func readAdmissionReviewFromFile(filename string) (*admissionv1.AdmissionReview, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	ar := &admissionv1.AdmissionReview{}
	if err := json.Unmarshal(data, ar); err != nil {
		return nil, err
	}

	return ar, nil
}

func TestImageSwapper_Mutate(t *testing.T) {
	ecrClient := &mockECRClient{}
	registryClient, _ := registry.NewMockECRClient(ecrClient, "", "")
	imagePullSecretProvider := secrets.NewDummyImagePullSecretsProvider()

	admissionReview, _ := readAdmissionReviewFromFile("../../test/requests/admissionreview.json")
	arm := model.NewAdmissionReviewV1(admissionReview)

	wh, _ := NewImageSwapperWebhook(registryClient, imagePullSecretProvider, []config.JMESPathFilter{}, types.ImageSwapPolicyAlways, types.ImageSwapPolicyAlways)
	resp, err := wh.Review(context.TODO(), arm)
	spew.Dump(resp)
	spew.Dump(err)
}
