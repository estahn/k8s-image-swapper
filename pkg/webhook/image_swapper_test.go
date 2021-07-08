package webhook

import (
	"reflect"
	"testing"

	"github.com/estahn/k8s-image-swapper/pkg/config"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
			ObjectMeta: v1.ObjectMeta{
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

func TestParseRegistryAuth(t *testing.T) {
	var testCases = []struct {
		inputConfig         []byte
		ExpectedRegistryUrl string
		ExpectedAuth        string
	}{
		{[]byte(`{"https://index.docker.io/v1/":{"username":"secretUser","password":"supersecretpass","email":"user@me.com","auth":"c2VjcmV0VXNlcjpzdXBlcnNlY3JldHBhc3M="}}`), "https://index.docker.io/v1/", "secretUser:supersecretpass"},
		{[]byte(`{"docker.io":{"username":"secretUser","password":"supersecretpass","email":"user@me.com","auth":"c2VjcmV0VXNlcjpzdXBlcnNlY3JldHBhc3M="}}`), "docker.io", "secretUser:supersecretpass"},
		{[]byte(`{"my.registry.com":{"username":"secretUser","password":"supersecretpass","email":"user@me.com","auth":"c2VjcmV0VXNlcjpzdXBlcnNlY3JldHBhc3M="}}`), "my.registry.com", "secretUser:supersecretpass"},
	}

	for _, tc := range testCases {
		t.Run("parseRegistryAuth", func(t *testing.T) {
			regisrtyUrl, auth := parseRegistryAuth(tc.inputConfig)
			if regisrtyUrl != tc.ExpectedRegistryUrl {
				t.Errorf("Got %q, Excepted %q", regisrtyUrl, tc.ExpectedRegistryUrl)
			}
			if auth != tc.ExpectedAuth {
				t.Errorf("Got %q, Excepted %q", auth, tc.ExpectedAuth)
			}
		})
	}
}

func TestAlignUrlBySrcRef(t *testing.T) {
	var testCases = []struct {
		InputRegistriesTokens    map[string][]string
		InputRegistryAliases     map[string]string
		ExpectedRegistriesTokens map[string][]string
	}{
		{map[string][]string{"a": {"0", "1", "2"}, "b": {"4", "5", "6"}, "c": {"7", "8", "9"}}, map[string]string{"a": "b"}, map[string][]string{"a": {"0", "1", "2", "4", "5", "6"}, "b": {"4", "5", "6"}, "c": {"7", "8", "9"}}},
	}

	for _, tc := range testCases {
		t.Run("alignUrlBySrcRef", func(t *testing.T) {
			registriesTokens := alignUrlBySrcRef(tc.InputRegistriesTokens, tc.InputRegistryAliases)
			if !reflect.DeepEqual(registriesTokens, tc.ExpectedRegistriesTokens) {
				t.Errorf("Got %q, expected %q", registriesTokens, tc.ExpectedRegistriesTokens)
			}
		})
	}
}

func TestAddTokenToAllRegistries(t *testing.T) {
	var testCases = []struct {
		InputRegistriesTokens    map[string][]string
		InputToken               string
		ExpectedRegistriesTokens map[string][]string
	}{
		{map[string][]string{"a": {"0", "1", "2"}, "b": {"4", "5", "6"}, "c": {"7", "8", "9"}}, "testingToken", map[string][]string{"a": {"0", "1", "2", "testingToken"}, "b": {"4", "5", "6", "testingToken"}, "c": {"7", "8", "9", "testingToken"}}},
	}

	for _, tc := range testCases {
		t.Run("addTokenToAllRegistries", func(t *testing.T) {
			registriesTokens := addTokenToAllRegistries(tc.InputRegistriesTokens, tc.InputToken)
			if !reflect.DeepEqual(registriesTokens, tc.ExpectedRegistriesTokens) {
				t.Errorf("Got %q, expected %q", registriesTokens, tc.ExpectedRegistriesTokens)
			}
		})
	}
}
