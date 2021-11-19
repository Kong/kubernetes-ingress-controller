package admission

const (
	ErrTextConsumerCredentialSecretNotFound   = "consumer referenced non-existent credentials secret"
	ErrTextConsumerCredentialValidationFailed = "consumer credential failed validation"
	ErrTextConsumerExists                     = "consumer already exists"
	ErrTextConsumerUnretrievable              = "failed to fetch consumer from kong"
	ErrTextConsumerUsernameEmpty              = "username cannot be empty"
	ErrTextFailedToRetrieveSecret             = "could not retrieve secrets from the kubernets API" //nolint:gosec
	ErrTextPluginConfigInvalid                = "could not parse plugin configuration"
	ErrTextPluginConfigValidationFailed       = "unable to validate plugin schema"
	ErrTextPluginConfigViolatesSchema         = "plugin failed schema validation: %s"
	ErrTextPluginNameEmpty                    = "plugin name cannot be empty"
	ErrTextPluginSecretConfigUnretrievable    = "could not load secret plugin configuration"
	ErrTextPluginUsesBothConfigTypes          = "plugin cannot use both Config and ConfigFrom"
)
