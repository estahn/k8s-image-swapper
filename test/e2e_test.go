// +build e2e

package test

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/shell"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	terratesttesting "github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IsKindCluster returns true if the underlying kubernetes cluster is kind. This is determined by getting the
// associated nodes and checking if a node is named "kind-control-plane".
//func IsKindCluster(t terratestTesting.TestingT, options *k8s.KubectlOptions) (bool, error) {
//	nodes, err := k8s.GetNodesE(t, options)
//	if err != nil {
//		return false, err
//	}
//
//	// ASSUMPTION: All minikube setups will have nodes with labels that are namespaced with minikube.k8s.io
//	for _, node := range nodes {
//		if !nodeHasMinikubeLabel(node) {
//			return false, nil
//		}
//	}
//
//	// At this point we know that all the nodes in the cluster has the minikube label, so we return true.
//	return true, nil
//}

// nodeHasMinikubeLabel returns true if any of the labels on the node is namespaced with minikube.k8s.io
//func nodeHasMinikubeLabel(node corev1.Node) bool {
//	labels := node.GetLabels()
//	for key, _ := range labels {
//		if strings.HasPrefix(key, "minikube.k8s.io") {
//			return true
//		}
//	}
//	return false
//}

// This file contains examples of how to use terratest to test helm charts by deploying the chart and verifying the
// deployment by hitting the service endpoint.
func TestHelmDeployment(t *testing.T) {
	workingDir, _ := filepath.Abs("..")

	awsAccountID := os.Getenv("AWS_ACCOUNT_ID")
	awsRegion := os.Getenv("AWS_DEFAULT_REGION")
	awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	ecrRegistry := awsAccountID + ".dkr.ecr." + awsRegion + ".amazonaws.com"
	ecrRepository := "docker.io/library/nginx"

	logger.Default = logger.New(newSensitiveLogger(
		logger.Default,
		[]*regexp.Regexp{
			regexp.MustCompile(awsAccountID),
			regexp.MustCompile(awsAccessKeyID),
			regexp.MustCompile(awsSecretAccessKey),
			regexp.MustCompile(`(--docker-password=)\S+`),
		},
	))

	// To ensure we can reuse the resource config on the same cluster to test different scenarios, we setup a unique
	// namespace for the resources for this test.
	// Note that namespaces must be lowercase.
	namespaceName := fmt.Sprintf("k8s-image-swapper-%s", strings.ToLower(random.UniqueId()))
	releaseName := fmt.Sprintf("k8s-image-swapper-%s",
		strings.ToLower(random.UniqueId()))

	// Setup the kubectl config and context. Here we choose to use the defaults, which is:
	// - HOME/.kube/config for the kubectl config file
	// - Current context of the kubectl config file
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)

	// Init ECR client
	ecrClient := aws.NewECRClient(t, awsRegion)

	defer test_structure.RunTestStage(t, "cleanup_aws", func() {
		_, err := ecrClient.DeleteRepository(&ecr.DeleteRepositoryInput{
			RepositoryName: awssdk.String(ecrRepository),
			RegistryId:     awssdk.String(awsAccountID),
			Force:          awssdk.Bool(true),
		})
		require.NoError(t, err)
	})

	defer test_structure.RunTestStage(t, "cleanup_k8s", func() {
		// Return the output before cleanup - helps in debugging
		k8s.RunKubectl(t, kubectlOptions, "logs", "--selector=app.kubernetes.io/name=k8s-image-swapper", "--tail=-1")
		helm.Delete(t, &helm.Options{KubectlOptions: kubectlOptions}, releaseName, true)
		k8s.DeleteNamespace(t, kubectlOptions, namespaceName)
	})

	test_structure.RunTestStage(t, "build_and_load_docker_image", func() {
		// Generate docker image to be tested
		shell.RunCommand(t, shell.Command{
			Command:    "goreleaser",
			Args:       []string{"release", "--snapshot", "--skip-publish", "--rm-dist"},
			WorkingDir: workingDir,
		})

		// Tag with "local" to ensure kind is not pulling from the GitHub Registry
		shell.RunCommand(t, shell.Command{
			Command: "docker",
			Args:    []string{"tag", "ghcr.io/estahn/k8s-image-swapper:latest", "local/k8s-image-swapper:latest"},
		})

		// Load generated docker image into kind
		shell.RunCommand(t, shell.Command{
			Command: "kind",
			Args:    []string{"load", "docker-image", "local/k8s-image-swapper:latest"},
		})
	})

	test_structure.RunTestStage(t, "deploy_webhook", func() {
		k8s.CreateNamespace(t, kubectlOptions, namespaceName)

		// Setup permissions for kind to be able to pull from ECR
		ecrAuthToken, _ := ecrClient.GetAuthorizationToken(&ecr.GetAuthorizationTokenInput{})
		ecrDecodedAuthToken, _ := base64.StdEncoding.DecodeString(*ecrAuthToken.AuthorizationData[0].AuthorizationToken)
		ecrUsernamePassword := bytes.Split(ecrDecodedAuthToken, []byte(":"))

		secretName := awsRegion + "-ecr-registry"
		k8s.RunKubectl(t, kubectlOptions, "create", "secret", "docker-registry",
			secretName,
			"--docker-server="+*ecrAuthToken.AuthorizationData[0].ProxyEndpoint,
			"--docker-username="+string(ecrUsernamePassword[0]),
			"--docker-password="+string(ecrUsernamePassword[1]),
			"--docker-email=anymail.doesnt.matter@email.com",
		)
		k8s.RunKubectl(t, kubectlOptions, "patch", "serviceaccount", "default", "-p",
			fmt.Sprintf("{\"imagePullSecrets\":[{\"name\":\"%s\"}]}", secretName),
		)

		// Setup the args. For this test, we will set the following input values:
		options := &helm.Options{
			KubectlOptions: kubectlOptions,
			SetValues: map[string]string{
				"config.logFormat":                  "console",
				"config.logLevel":                   "debug",
				"config.dryRun":                     "false",
				"config.target.aws.accountId":       awsAccountID,
				"config.target.aws.region":          awsRegion,
				"config.imageSwapPolicy":            "always",
				"config.imageCopyPolicy":            "delayed",
				"config.source.filters[0].jmespath": "obj.metadata.name != 'nginx'",
				"awsSecretName":                     "k8s-image-swapper-aws",
				"image.repository":                  "local/k8s-image-swapper",
				"image.tag":                         "latest",
			},
		}

		k8s.RunKubectl(t, kubectlOptions, "create", "secret", "generic", "k8s-image-swapper-aws",
			fmt.Sprintf("--from-literal=aws_access_key_id=%s", awsAccessKeyID),
			fmt.Sprintf("--from-literal=aws_secret_access_key=%s", awsSecretAccessKey),
		)

		helm.AddRepo(t, options, "estahn", "https://estahn.github.io/charts/")
		helm.Install(t, options, "estahn/k8s-image-swapper", releaseName)
	})

	test_structure.RunTestStage(t, "validate", func() {
		k8s.WaitUntilNumPodsCreated(t, kubectlOptions, metav1.ListOptions{LabelSelector: "app.kubernetes.io/name=k8s-image-swapper"}, 1, 30, 10*time.Second)
		k8s.WaitUntilServiceAvailable(t, kubectlOptions, releaseName, 30, 10*time.Second)

		// Launch nginx container to verify functionality
		k8s.RunKubectl(t, kubectlOptions, "run", "nginx", "--image=nginx", "--restart=Never")
		k8s.WaitUntilPodAvailable(t, kubectlOptions, "nginx", 30, 10*time.Second)

		// Verify container is running with images from ECR.
		// Implicit proof for repository creation and images pull/push via k8s-image-swapper.
		nginxPod := k8s.GetPod(t, kubectlOptions, "nginx")

		require.Equal(t, ecrRegistry+"/"+ecrRepository+":latest", nginxPod.Spec.Containers[0].Image, "container should be prefixed with ECR address")
	})
}

type sensitiveLogger struct {
	logger   logger.TestLogger
	patterns []*regexp.Regexp
}

func newSensitiveLogger(logger *logger.Logger, patterns []*regexp.Regexp) *sensitiveLogger {
	return &sensitiveLogger{
		logger:   logger,
		patterns: patterns,
	}
}

func (l *sensitiveLogger) Logf(t terratesttesting.TestingT, format string, args ...interface{}) {
	var redactedArgs []interface{}

	obfuscateWith := "$1*******"

	redactedArgs = args

	for _, pattern := range l.patterns {
		for i, arg := range redactedArgs {
			switch arg := arg.(type) {
			case string:
				redactedArgs[i] = pattern.ReplaceAllString(arg, obfuscateWith)
			case []string:
				var result []string
				for _, s := range arg {
					result = append(result, pattern.ReplaceAllString(s, obfuscateWith))
				}
				redactedArgs[i] = result
			default:
				panic("type needs implementation")
			}
		}
	}

	l.logger.Logf(t, format, redactedArgs...)
}
