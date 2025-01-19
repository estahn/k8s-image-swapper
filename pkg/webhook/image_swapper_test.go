package webhook

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/alitto/pond"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	"github.com/estahn/k8s-image-swapper/pkg/config"
	"github.com/estahn/k8s-image-swapper/pkg/registry"
	"github.com/estahn/k8s-image-swapper/pkg/secrets"
	"github.com/estahn/k8s-image-swapper/pkg/types"
	"github.com/slok/kubewebhook/v2/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
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
	mock.Mock
	ecriface.ECRAPI
}

func (m *mockECRClient) CreateRepositoryWithContext(ctx context.Context, createRepositoryInput *ecr.CreateRepositoryInput, opts ...request.Option) (*ecr.CreateRepositoryOutput, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	m.Called(ctx, createRepositoryInput)
	return &ecr.CreateRepositoryOutput{}, nil
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	os.Exit(0)
}

func readAdmissionReviewFromFile(filename string) (*admissionv1.AdmissionReview, error) {
	data, err := os.ReadFile("../../test/requests/" + filename)
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
	expectedRepositories := []string{
		"docker.io/library/init-container",
		"docker.io/library/nginx",
		"k8s.gcr.io/ingress-nginx/controller",
		"us-central1-docker.pkg.dev/gcp-project-123/main/k8s.gcr.io/ingress-nginx/controller",
	}

	ecrClient := new(mockECRClient)

	for _, expectedRepository := range expectedRepositories {
		ecrClient.On(
			"CreateRepositoryWithContext",
			mock.AnythingOfType("*context.valueCtx"),
			&ecr.CreateRepositoryInput{
				ImageScanningConfiguration: &ecr.ImageScanningConfiguration{
					ScanOnPush: aws.Bool(true),
				},
				EncryptionConfiguration: &ecr.EncryptionConfiguration{
					EncryptionType: aws.String("AES256"),
				},
				ImageTagMutability: aws.String("MUTABLE"),
				RepositoryName:     aws.String(expectedRepository),
				RegistryId:         aws.String("123456789"),
				Tags: []*ecr.Tag{
					{
						Key:   aws.String("CreatedBy"),
						Value: aws.String("k8s-image-swapper"),
					},
					{
						Key:   aws.String("AnotherTag"),
						Value: aws.String("another-tag"),
					},
				},
			}).Return(mock.Anything)
	}

	registryClient, _ := registry.NewMockECRClient(ecrClient, "ap-southeast-2", "123456789.dkr.ecr.ap-southeast-2.amazonaws.com", "123456789", "arn:aws:iam::123456789:role/fakerole")

	admissionReview, _ := readAdmissionReviewFromFile("admissionreview-simple.json")
	admissionReviewModel := model.NewAdmissionReviewV1(admissionReview)

	copier := pond.New(1, 1)
	// TODO: test types.ImageSwapPolicyExists
	wh, err := NewImageSwapperWebhookWithOpts(
		registryClient,
		Copier(copier),
		ImageSwapPolicy(types.ImageSwapPolicyAlways),
	)

	assert.NoError(t, err, "NewImageSwapperWebhookWithOpts executed without errors")

	resp, err := wh.Review(context.Background(), admissionReviewModel)

	// TODO: think about moving "expected" into a file, e.g. admissionreview-simple-response-ecr.json
	// container with name "skip-test-gar" should be skipped, hence there is no "replace" operation for it
	expected := `[
		{"op":"replace","path":"/spec/initContainers/0/image","value":"123456789.dkr.ecr.ap-southeast-2.amazonaws.com/docker.io/library/init-container:latest"},
		{"op":"replace","path":"/spec/containers/0/image","value":"123456789.dkr.ecr.ap-southeast-2.amazonaws.com/docker.io/library/nginx:latest"},
		{"op":"replace","path":"/spec/containers/1/image","value":"123456789.dkr.ecr.ap-southeast-2.amazonaws.com/k8s.gcr.io/ingress-nginx/controller@sha256:9bba603b99bf25f6d117cf1235b6598c16033ad027b143c90fa5b3cc583c5713"},
		{"op":"replace","path":"/spec/containers/3/image","value":"123456789.dkr.ecr.ap-southeast-2.amazonaws.com/us-central1-docker.pkg.dev/gcp-project-123/main/k8s.gcr.io/ingress-nginx/controller@sha256:9bba603b99bf25f6d117cf1235b6598c16033ad027b143c90fa5b3cc583c5713"}
	]`

	assert.JSONEq(t, expected, string(resp.(*model.MutatingAdmissionResponse).JSONPatchPatch))
	assert.Nil(t, resp.(*model.MutatingAdmissionResponse).Warnings)
	assert.NoError(t, err, "Webhook executed without errors")

	// Ensure the worker pool is empty before asserting ecrClient
	copier.StopAndWait()

	ecrClient.AssertExpectations(t)
	assert.Equal(t, uint64(4), copier.SubmittedTasks())
}

// TestImageSwapper_MutateWithImagePullSecrets tests mutating with imagePullSecret support
func TestImageSwapper_MutateWithImagePullSecrets(t *testing.T) {
	ecrClient := new(mockECRClient)
	ecrClient.On(
		"CreateRepositoryWithContext",
		mock.AnythingOfType("*context.valueCtx"),
		&ecr.CreateRepositoryInput{
			ImageScanningConfiguration: &ecr.ImageScanningConfiguration{
				ScanOnPush: aws.Bool(true),
			},
			EncryptionConfiguration: &ecr.EncryptionConfiguration{
				EncryptionType: aws.String("AES256"),
			},
			ImageTagMutability: aws.String("MUTABLE"),
			RegistryId:         aws.String("123456789"),
			RepositoryName:     aws.String("docker.io/library/nginx"),
			Tags: []*ecr.Tag{
				{
					Key:   aws.String("CreatedBy"),
					Value: aws.String("k8s-image-swapper"),
				},
				{
					Key:   aws.String("AnotherTag"),
					Value: aws.String("another-tag"),
				},
			},
		}).Return(mock.Anything)

	registryClient, _ := registry.NewMockECRClient(ecrClient, "ap-southeast-2", "123456789.dkr.ecr.ap-southeast-2.amazonaws.com", "123456789", "arn:aws:iam::123456789:role/fakerole")

	admissionReview, _ := readAdmissionReviewFromFile("admissionreview-imagepullsecrets.json")
	admissionReviewModel := model.NewAdmissionReviewV1(admissionReview)

	clientSet := fake.NewSimpleClientset()

	svcAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-service-account",
		},
		ImagePullSecrets: []corev1.LocalObjectReference{
			{Name: "my-sa-secret"},
		},
	}
	svcAccountSecretDockerConfigJson := []byte(`{"auths":{"my-sa-secret.registry.example.com":{"username":"my-sa-secret","password":"xxxxxxxxxxx","email":"jdoe@example.com","auth":"c3R...zE2"}}}`)
	svcAccountSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-sa-secret",
		},
		Type: corev1.SecretTypeDockerConfigJson,
		Data: map[string][]byte{
			corev1.DockerConfigJsonKey: svcAccountSecretDockerConfigJson,
		},
	}
	podSecretDockerConfigJson := []byte(`{"auths":{"my-pod-secret.registry.example.com":{"username":"my-sa-secret","password":"xxxxxxxxxxx","email":"jdoe@example.com","auth":"c3R...zE2"}}}`)
	podSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-pod-secret",
		},
		Type: corev1.SecretTypeDockerConfigJson,
		Data: map[string][]byte{
			corev1.DockerConfigJsonKey: podSecretDockerConfigJson,
		},
	}

	_, _ = clientSet.CoreV1().ServiceAccounts("test-ns").Create(context.Background(), svcAccount, metav1.CreateOptions{})
	_, _ = clientSet.CoreV1().Secrets("test-ns").Create(context.Background(), svcAccountSecret, metav1.CreateOptions{})
	_, _ = clientSet.CoreV1().Secrets("test-ns").Create(context.Background(), podSecret, metav1.CreateOptions{})

	provider := secrets.NewKubernetesImagePullSecretsProvider(clientSet)

	copier := pond.New(1, 1)
	// TODO: test types.ImageSwapPolicyExists
	wh, err := NewImageSwapperWebhookWithOpts(
		registryClient,
		ImagePullSecretsProvider(provider),
		Copier(copier),
		ImageSwapPolicy(types.ImageSwapPolicyAlways),
	)

	assert.NoError(t, err, "NewImageSwapperWebhookWithOpts executed without errors")

	resp, err := wh.Review(context.Background(), admissionReviewModel)

	assert.JSONEq(t, "[{\"op\":\"replace\",\"path\":\"/spec/containers/0/image\",\"value\":\"123456789.dkr.ecr.ap-southeast-2.amazonaws.com/docker.io/library/nginx:latest\"}]", string(resp.(*model.MutatingAdmissionResponse).JSONPatchPatch))
	assert.Nil(t, resp.(*model.MutatingAdmissionResponse).Warnings)
	assert.NoError(t, err, "Webhook executed without errors")

	// Ensure the worker pool is empty before asserting ecrClient
	copier.StopAndWait()

	ecrClient.AssertExpectations(t)
	assert.Equal(t, uint64(1), copier.SubmittedTasks())
}

func TestImageSwapper_GAR_Mutate(t *testing.T) {
	registryClient, _ := registry.NewMockGARClient(nil, "us-central1-docker.pkg.dev/gcp-project-123/main")

	admissionReview, _ := readAdmissionReviewFromFile("admissionreview-simple.json")
	admissionReviewModel := model.NewAdmissionReviewV1(admissionReview)

	copier := pond.New(1, 1)
	// TODO: test types.ImageSwapPolicyExists
	wh, err := NewImageSwapperWebhookWithOpts(
		registryClient,
		Copier(copier),
		ImageSwapPolicy(types.ImageSwapPolicyAlways),
	)

	assert.NoError(t, err, "NewImageSwapperWebhookWithOpts executed without errors")

	resp, err := wh.Review(context.TODO(), admissionReviewModel)

	// container with name "skip-test-gar" should be skipped, hence there is no "replace" operation for it
	expected := `[
		{"op":"replace","path":"/spec/initContainers/0/image","value":"us-central1-docker.pkg.dev/gcp-project-123/main/docker.io/library/init-container:latest"},
		{"op":"replace","path":"/spec/containers/0/image","value":"us-central1-docker.pkg.dev/gcp-project-123/main/docker.io/library/nginx:latest"},
		{"op":"replace","path":"/spec/containers/1/image","value":"us-central1-docker.pkg.dev/gcp-project-123/main/k8s.gcr.io/ingress-nginx/controller@sha256:9bba603b99bf25f6d117cf1235b6598c16033ad027b143c90fa5b3cc583c5713"},
		{"op":"replace","path":"/spec/containers/2/image","value":"us-central1-docker.pkg.dev/gcp-project-123/main/123456789.dkr.ecr.ap-southeast-2.amazonaws.com/k8s.gcr.io/ingress-nginx/controller@sha256:9bba603b99bf25f6d117cf1235b6598c16033ad027b143c90fa5b3cc583c5713"}
	]`

	assert.JSONEq(t, expected, string(resp.(*model.MutatingAdmissionResponse).JSONPatchPatch))
	assert.Nil(t, resp.(*model.MutatingAdmissionResponse).Warnings)
	assert.NoError(t, err, "Webhook executed without errors")
	assert.Equal(t, uint64(4), copier.SubmittedTasks())
}

func TestImageSwapper_skipRegistryGcrIo_Mutate(t *testing.T) {
	registryClient, _ := registry.NewMockGARClient(nil, "us-central1-docker.pkg.dev/gcp-project-123/main")

	admissionReview, _ := readAdmissionReviewFromFile("admissionreview-simple.json")
	admissionReviewModel := model.NewAdmissionReviewV1(admissionReview)

	copier := pond.New(1, 1)
	wh, err := NewImageSwapperWebhookWithOpts(
		registryClient,
		Copier(copier),
		ImageSwapPolicy(types.ImageSwapPolicyAlways),
		ImageCopySkipRegistries([]string{"k8s.gcr.io"}),
	)

	assert.NoError(t, err, "NewImageSwapperWebhookWithOpts executed without errors")

	resp, err := wh.Review(context.TODO(), admissionReviewModel)

	// container with name "skip-test-gar" should be skipped, hence there is no "replace" operation for it
	expected := `[
		{"op":"replace","path":"/spec/initContainers/0/image","value":"us-central1-docker.pkg.dev/gcp-project-123/main/docker.io/library/init-container:latest"},
		{"op":"replace","path":"/spec/containers/0/image","value":"us-central1-docker.pkg.dev/gcp-project-123/main/docker.io/library/nginx:latest"},
		{"op":"replace","path":"/spec/containers/1/image","value":"us-central1-docker.pkg.dev/gcp-project-123/main/k8s.gcr.io/ingress-nginx/controller@sha256:9bba603b99bf25f6d117cf1235b6598c16033ad027b143c90fa5b3cc583c5713"},
		{"op":"replace","path":"/spec/containers/2/image","value":"us-central1-docker.pkg.dev/gcp-project-123/main/123456789.dkr.ecr.ap-southeast-2.amazonaws.com/k8s.gcr.io/ingress-nginx/controller@sha256:9bba603b99bf25f6d117cf1235b6598c16033ad027b143c90fa5b3cc583c5713"}
	]`

	assert.JSONEq(t, expected, string(resp.(*model.MutatingAdmissionResponse).JSONPatchPatch))
	assert.Nil(t, resp.(*model.MutatingAdmissionResponse).Warnings)
	assert.NoError(t, err, "Webhook executed without errors")

	// check the amount of submitted tasks for the copier
	assert.Equal(t, uint64(3), copier.SubmittedTasks())
}

func TestImageSwapper_skipRegistryDockerIo_Mutate(t *testing.T) {
	registryClient, _ := registry.NewMockGARClient(nil, "us-central1-docker.pkg.dev/gcp-project-123/main")

	admissionReview, _ := readAdmissionReviewFromFile("admissionreview-simple.json")
	admissionReviewModel := model.NewAdmissionReviewV1(admissionReview)

	copier := pond.New(1, 1)
	wh, err := NewImageSwapperWebhookWithOpts(
		registryClient,
		Copier(copier),
		ImageSwapPolicy(types.ImageSwapPolicyAlways),
		ImageCopySkipRegistries([]string{"docker.io"}),
	)

	assert.NoError(t, err, "NewImageSwapperWebhookWithOpts executed without errors")

	resp, err := wh.Review(context.TODO(), admissionReviewModel)

	// container with name "skip-test-gar" should be skipped, hence there is no "replace" operation for it
	expected := `[
		{"op":"replace","path":"/spec/initContainers/0/image","value":"us-central1-docker.pkg.dev/gcp-project-123/main/docker.io/library/init-container:latest"},
		{"op":"replace","path":"/spec/containers/0/image","value":"us-central1-docker.pkg.dev/gcp-project-123/main/docker.io/library/nginx:latest"},
		{"op":"replace","path":"/spec/containers/1/image","value":"us-central1-docker.pkg.dev/gcp-project-123/main/k8s.gcr.io/ingress-nginx/controller@sha256:9bba603b99bf25f6d117cf1235b6598c16033ad027b143c90fa5b3cc583c5713"},
		{"op":"replace","path":"/spec/containers/2/image","value":"us-central1-docker.pkg.dev/gcp-project-123/main/123456789.dkr.ecr.ap-southeast-2.amazonaws.com/k8s.gcr.io/ingress-nginx/controller@sha256:9bba603b99bf25f6d117cf1235b6598c16033ad027b143c90fa5b3cc583c5713"}
	]`

	assert.JSONEq(t, expected, string(resp.(*model.MutatingAdmissionResponse).JSONPatchPatch))
	assert.Nil(t, resp.(*model.MutatingAdmissionResponse).Warnings)
	assert.NoError(t, err, "Webhook executed without errors")

	// check the amount of submitted tasks for the copier
	assert.Equal(t, uint64(2), copier.SubmittedTasks())
}
