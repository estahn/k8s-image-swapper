package utils

import (
	"crypto/tls"
	"path"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog/log"
)

// KeypairReloader structs holds cert path and certs
type KeypairReloader struct {
	certMu      sync.RWMutex
	cert        *tls.Certificate
	tlsCertFile string
	tlsKeyFile  string
}

// NewKeypairReloader will load certs on first run and trigger a goroutine for fsnotify watcher
func NewKeypairReloader(tlsCertFile, tlsKeyFile string) (*KeypairReloader, error) {
	result := &KeypairReloader{
		tlsCertFile: tlsCertFile,
		tlsKeyFile:  tlsKeyFile,
	}
	cert, err := tls.LoadX509KeyPair(tlsCertFile, tlsKeyFile)
	if err != nil {
		return nil, err
	}
	result.cert = &cert

	// creates a new file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			watcher.Close()
		}
	}()

	// Notify on changes to the cert directory
	if err := watcher.Add(path.Dir(tlsCertFile)); err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				// Watch for changes to the tlsCertFile
				if event.Name == tlsCertFile {
					log.Info().Msg("Reloading certs")
					if err := result.reload(); err != nil {
						log.Err(err).Msg("Could not load new certs")
					}
				}

				// watch for errors
			case err := <-watcher.Errors:
				log.Err(err).Msg("Watcher error")
			}
		}
	}()

	return result, nil
}

// reload loads updated cert and key whenever they are updated
func (kpr *KeypairReloader) reload() error {
	newCert, err := tls.LoadX509KeyPair(kpr.tlsCertFile, kpr.tlsKeyFile)
	if err != nil {
		return err
	}
	kpr.certMu.Lock()
	defer kpr.certMu.Unlock()
	kpr.cert = &newCert
	return nil
}

// GetCertificateFunc will return function which will be used as tls.Config.GetCertificate
func (kpr *KeypairReloader) GetCertificateFunc() func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
	return func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		kpr.certMu.RLock()
		defer kpr.certMu.RUnlock()
		return kpr.cert, nil
	}
}
