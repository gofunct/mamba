package config

import (
	"bytes"
	"io"
	"reflect"
	"github.com/prometheus/common/log"
)

type defaultRemoteProvider struct {
	provider      string
	endpoint      string
	path          string
	secretKeyring string
}

func (rp defaultRemoteProvider) Provider() string {
	return rp.provider
}

func (rp defaultRemoteProvider) Endpoint() string {
	return rp.endpoint
}

func (rp defaultRemoteProvider) Path() string {
	return rp.path
}

func (rp defaultRemoteProvider) SecretKeyring() string {
	return rp.secretKeyring
}

// RemoteProvider stores the configuration necessary
// to connect to a remote key/value store.
// Optional secretKeyring to unencrypt encrypted values
// can be provided.
type RemoteProvider interface {
	Provider() string
	Endpoint() string
	Path() string
	SecretKeyring() string
}

func (v *config) AddRemoteProvider(provider, endpoint, path string) error {
	if !stringInSlice(provider, SupportedRemoteProviders) {
		return UnsupportedRemoteProviderError(provider)
	}
	if provider != "" && endpoint != "" {
		rp := &defaultRemoteProvider{
			endpoint: endpoint,
			provider: provider,
			path:     path,
		}
		if !v.providerPathExists(rp) {
			v.remoteProviders = append(v.remoteProviders, rp)
		}
	}
	return nil
}

func (v *config) AddSecureRemoteProvider(provider, endpoint, path, secretkeyring string) error {
	if !stringInSlice(provider, SupportedRemoteProviders) {
		return UnsupportedRemoteProviderError(provider)
	}
	if provider != "" && endpoint != "" {
		rp := &defaultRemoteProvider{
			endpoint:      endpoint,
			provider:      provider,
			path:          path,
			secretKeyring: secretkeyring,
		}
		if !v.providerPathExists(rp) {
			v.remoteProviders = append(v.remoteProviders, rp)
		}
	}
	return nil
}

func (v *config) providerPathExists(p *defaultRemoteProvider) bool {
	for _, y := range v.remoteProviders {
		if reflect.DeepEqual(y, p) {
			return true
		}
	}
	return false
}

// Remoteconfig is optional, see the remote package
var Remoteconfig remoteconfigFactory

type remoteconfigFactory interface {
	Get(rp RemoteProvider) (io.Reader, error)
	Watch(rp RemoteProvider) (io.Reader, error)
	WatchChannel(rp RemoteProvider) (<-chan *RemoteResponse, chan bool)
}

type RemoteResponse struct {
	Value []byte
	Error error
}

func (v *config) watchRemoteconfig(provider RemoteProvider) (map[string]interface{}, error) {
	log.Debug( "event", "watching remote config...")

	reader, err := Remoteconfig.Watch(provider)
	if err != nil {
		return nil, err
	}
	err = v.unmarshalReader(reader, v.kvstore)
	return v.kvstore, err
}

func (v *config) getRemoteconfig(provider RemoteProvider) (map[string]interface{}, error) {
	log.Debug( "event", "getting remote config...")

	reader, err := Remoteconfig.Get(provider)
	if err != nil {
		return nil, err
	}
	err = v.unmarshalReader(reader, v.kvstore)
	return v.kvstore, err
}

// Retrieve the first found remote configuration.
func (v *config) watchKeyValueconfigOnChannel() error {
	log.Debug( "event", "watching kv config on channel...")

	for _, rp := range v.remoteProviders {
		respc, _ := Remoteconfig.WatchChannel(rp)
		//Todo: Add quit channel
		go func(rc <-chan *RemoteResponse) {
			for {
				b := <-rc
				reader := bytes.NewReader(b.Value)
				v.unmarshalReader(reader, v.kvstore)
			}
		}(respc)
		return nil
	}
	return RemoteconfigError("No Files Found")
}

// Retrieve the first found remote configuration.
func (v *config) watchKeyValueconfig() error {
	log.Debug( "event", "watching kv config...")
	for _, rp := range v.remoteProviders {
		val, err := v.watchRemoteconfig(rp)
		if err != nil {
			continue
		}
		v.kvstore = val
		return nil
	}
	return RemoteconfigError("No Files Found")
}

// Retrieve the first found remote configuration.
func (v *config) getKeyValueconfig() error {
	if Remoteconfig == nil {
		return RemoteconfigError("Enable the remote features by doing a blank import of the config/remote package: '_ github.com/spf13/config/remote'")
	}

	for _, rp := range v.remoteProviders {
		val, err := v.getRemoteconfig(rp)
		if err != nil {
			continue
		}
		v.kvstore = val
		return nil
	}
	return RemoteconfigError("No Files Found")
}

func (v *config) ReadRemoteconfig() error {
	return v.getKeyValueconfig()
}

func (v *config) WatchRemoteconfig() error {
	return v.watchKeyValueconfig()
}

func (v *config) WatchRemoteconfigOnChannel() error {
	return v.watchKeyValueconfigOnChannel()
}
