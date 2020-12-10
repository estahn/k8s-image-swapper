package main // import "github.com/hipages/k8s-image-swapper"

import (
	"flag"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
	whhttp "github.com/slok/kubewebhook/pkg/http"
	mutatingwh "github.com/slok/kubewebhook/pkg/webhook/mutating"
	corev1 "k8s.io/api/core/v1"
)

type config struct {
	certFile string
	keyFile  string
}

func initFlags() *config {
	cfg := &config{}

	fl := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fl.StringVar(&cfg.certFile, "tls-cert-file", "", "TLS certificate file")
	fl.StringVar(&cfg.keyFile, "tls-key-file", "", "TLS key file")

	fl.Parse(os.Args[1:])
	return cfg
}

func main() {
	flag.Parse()

	//cfg := initFlags()

	ism := NewImageSwapper(
		"***REMOVED***.dkr.ecr.ap-southeast-2.amazonaws.com",
	)

	// Create our mutator
	mt := mutatingwh.MutatorFunc(ism.Mutate)

	mcfg := mutatingwh.WebhookConfig{
		Name: "imageSwapper",
		Obj:  &corev1.Pod{},
	}
	wh, err := mutatingwh.NewWebhook(mcfg, mt, nil, nil, nil)
	if err != nil {
		log.Err(err).Msg("error creating webhook")
		os.Exit(1)
	}

	// Get the handler for our webhook.
	whHandler, err := whhttp.HandlerFor(wh)
	if err != nil {
		log.Err(err).Msg("error creating webhook handler")
		os.Exit(1)
	}

	log.Info().Msg("Listening on :8080")
	//err = http.ListenAndServeTLS(":8080", cfg.certFile, cfg.keyFile, whHandler)
	err = http.ListenAndServe(":8080", whHandler)
	if err != nil {
		log.Err(err).Msg("error serving webhook")
		os.Exit(1)
	}
}
