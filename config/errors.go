package config

import 	"github.com/fatih/color"

// UnsupportedRemoteProviderError denotes encountering an unsupported remote
// provider. Currently only etcd and Consul are supported.
type UnsupportedRemoteProviderError string

// Error returns the formatted remote provider error.
func (str UnsupportedRemoteProviderError) Error() string {
	return color.RedString("Unsupported Remote Provider Type %q", string(str))
}

// UnsupportedconfigError denotes encountering an unsupported
// configuration filetype.
type UnsupportedconfigError string

// Error returns the formatted configuration error.
func (str UnsupportedconfigError) Error() string {
	return color.RedString("Unsupported config Type %q", string(str))
}

// configMarshalError happens when failing to marshal the configuration.
type configMarshalError struct {
	err error
}

// Error returns the formatted configuration error.
func (e configMarshalError) Error() string {
	return color.RedString("While marshaling config: %s", e.err.Error())
}

// RemoteconfigError denotes encountering an error while trying to
// pull the configuration from the remote provider.
type RemoteconfigError string

// Error returns the formatted remote provider error
func (rce RemoteconfigError) Error() string {
	return color.RedString("Remote configurations Error: %s", string(rce))
}

// configFileNotFoundError denotes failing to find configuration file.
type configFileNotFoundError struct {
	name, locations string
}

// Error returns the formatted configuration error.
func (fnfe configFileNotFoundError) Error() string {
	return color.RedString("config File %q Not Found in %q", fnfe.name, fnfe.locations)
}
