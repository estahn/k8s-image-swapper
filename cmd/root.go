/*
Copyright © 2020 Enrico Stahn <enrico.stahn@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/estahn/k8s-image-swapper/pkg/config"
	"github.com/estahn/k8s-image-swapper/pkg/registry"
	"github.com/estahn/k8s-image-swapper/pkg/secrets"
	"github.com/estahn/k8s-image-swapper/pkg/types"
	"github.com/estahn/k8s-image-swapper/pkg/webhook"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	kwhhttp "github.com/slok/kubewebhook/v2/pkg/http"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var cfgFile string
var cfg *config.Config = &config.Config{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "k8s-image-swapper",
	Short: "Mirror images into your own registry and swap image references automatically.",
	Long: `Mirror images into your own registry and swap image references automatically.

A mutating webhook for Kubernetes, pointing the images to a new location.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		//promReg := prometheus.NewRegistry()
		//metricsRec := metrics.NewPrometheus(promReg)
		log.Trace().Interface("config", cfg).Msg("config")

		// Create registry clients for source registries
		sourceRegistryClients := []registry.Client{}
		for _, reg := range cfg.Source.Registries {
			sourceRegistryClient, err := registry.NewClient(reg)
			if err != nil {
				log.Err(err).Msgf("error connecting to source registry at %s", reg.Domain())
				os.Exit(1)
			}
			sourceRegistryClients = append(sourceRegistryClients, sourceRegistryClient)
		}

		// Create a registry client for private target registry
		targetRegistryClient, err := registry.NewClient(cfg.Target)
		if err != nil {
			log.Err(err).Msgf("error connecting to target registry at %s", cfg.Target.Domain())
			os.Exit(1)
		}

		imageSwapPolicy, err := types.ParseImageSwapPolicy(cfg.ImageSwapPolicy)
		if err != nil {
			log.Err(err).Str("policy", cfg.ImageSwapPolicy).Msg("parsing image swap policy failed")
		}

		imageCopyPolicy, err := types.ParseImageCopyPolicy(cfg.ImageCopyPolicy)
		if err != nil {
			log.Err(err).Str("policy", cfg.ImageCopyPolicy).Msg("parsing image copy policy failed")
		}

		imageCopyDeadline := config.DefaultImageCopyDeadline
		if cfg.ImageCopyDeadline != 0 {
			imageCopyDeadline = cfg.ImageCopyDeadline
		}

		imagePullSecretProvider := setupImagePullSecretsProvider()

		// Inform secret provider about managed private source registries
		imagePullSecretProvider.SetAuthenticatedRegistries(sourceRegistryClients)

		wh, err := webhook.NewImageSwapperWebhookWithOpts(
			targetRegistryClient,
			webhook.Filters(cfg.Source.Filters),
			webhook.ImagePullSecretsProvider(imagePullSecretProvider),
			webhook.ImageSwapPolicy(imageSwapPolicy),
			webhook.ImageCopyPolicy(imageCopyPolicy),
			webhook.ImageCopyDeadline(imageCopyDeadline),
			webhook.ImageCopySkipRegistries(cfg.ImageCopySkipRegistries),
		)
		if err != nil {
			log.Err(err).Msg("error creating webhook")
			os.Exit(1)
		}

		// Get the handler for our webhook.
		whHandler, err := kwhhttp.HandlerFor(kwhhttp.HandlerConfig{Webhook: wh})
		if err != nil {
			log.Err(err).Msg("error creating webhook handler")
			os.Exit(1)
		}

		handler := http.NewServeMux()
		handler.Handle("/webhook", whHandler)
		handler.Handle("/metrics", promhttp.Handler())
		handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte(`<html>
			 <head><title>k8s-image-webhook</title></head>
			 <body>
			 <h1>k8s-image-webhook</h1>
			 <ul><li><a href='/metrics'>Metrics</a></li><li><a href='/webhook'>Webhook</a></li></ul>
			 </body>
			 </html>`))

			if err != nil {
				log.Error()
			}
		})

		srv := &http.Server{
			Addr: cfg.ListenAddress,
			// Good practice to set timeouts to avoid Slowloris attacks.
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60,
			Handler:      handler,
		}

		go func() {
			log.Info().Msgf("Listening on %v", cfg.ListenAddress)
			//err = http.ListenAndServeTLS(":8080", cfg.certFile, cfg.keyFile, whHandler)
			if cfg.TLSCertFile != "" && cfg.TLSKeyFile != "" {
				if err := srv.ListenAndServeTLS(cfg.TLSCertFile, cfg.TLSKeyFile); err != nil {
					log.Err(err).Msg("error serving webhook")
					os.Exit(1)
				}
			} else {
				if err := srv.ListenAndServe(); err != nil {
					log.Err(err).Msg("error serving webhook")
					os.Exit(1)
				}
			}
		}()

		c := make(chan os.Signal, 1)
		// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C) or SIGTERM
		// SIGKILL, SIGQUIT will not be caught.
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		// Block until we receive our signal.
		<-c

		// Create a deadline to wait for.
		var wait time.Duration
		ctx, cancel := context.WithTimeout(context.Background(), wait)
		defer cancel()
		// Doesn't block if no connections, but will otherwise wait
		// until the timeout deadline.
		if err := srv.Shutdown(ctx); err != nil {
			log.Err(err).Msg("Error during shutdown")
		}
		// Optionally, you could run srv.Shutdown in a goroutine and block on
		// <-ctx.Done() if your application should wait for other services
		// to finalize based on context cancellation.
		log.Info().Msg("Shutting down")
		os.Exit(0)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initLogger)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.k8s-image-swapper.yaml)")
	rootCmd.PersistentFlags().StringVar(&cfg.LogLevel, "log-level", "info", "Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, fatal]")
	rootCmd.PersistentFlags().StringVar(&cfg.LogFormat, "log-format", "json", "Format of the log messages. Valid levels: [json, console]")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringVar(&cfg.ListenAddress, "listen-address", ":8443", "Address on which to expose the webhook")
	rootCmd.Flags().StringVar(&cfg.TLSCertFile, "tls-cert-file", "", "File containing the TLS certificate")
	rootCmd.Flags().StringVar(&cfg.TLSKeyFile, "tls-key-file", "", "File containing the TLS private key")
	rootCmd.Flags().BoolVar(&cfg.DryRun, "dry-run", true, "If true, print the action taken without taking it")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Default to aws target registry type if none are defined
	config.SetViperDefaults(viper.GetViper())

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".k8s-image-swapper" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".k8s-image-swapper")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Info().Str("file", viper.ConfigFileUsed()).Msg("using config file")
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Err(err).Msg("failed to unmarshal the config file")
	}

	//validate := validator.New()
	//if err := validate.Struct(cfg); err != nil {
	//	validationErrors := err.(validator.ValidationErrors)
	//	log.Err(validationErrors).Msg("validation errors for config file")
	//}
}

// initLogger configures the log level
func initLogger() {
	if cfg.LogFormat == "console" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	lvl, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		lvl = zerolog.InfoLevel
		log.Err(err).Msgf("Could not set log level to '%v'.", cfg.LogLevel)
	}

	zerolog.SetGlobalLevel(lvl)

	// add file and line number to log if level is trace
	if lvl == zerolog.TraceLevel {
		log.Logger = log.With().Caller().Logger()
	}
}

// setupImagePullSecretsProvider configures the provider handling secrets
func setupImagePullSecretsProvider() secrets.ImagePullSecretsProvider {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Warn().Err(err).Msg("failed to configure Kubernetes client, will continue without reading secrets")
		return secrets.NewDummyImagePullSecretsProvider()
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Warn().Err(err).Msg("failed to configure Kubernetes client, will continue without reading secrets")
		return secrets.NewDummyImagePullSecretsProvider()
	}

	return secrets.NewKubernetesImagePullSecretsProvider(clientset)
}
