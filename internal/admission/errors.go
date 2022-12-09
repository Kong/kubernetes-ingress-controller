package admission

const (
	ErrTextConsumerCredentialSecretNotFound   = "consumer referenced non-existent credentials secret"
	ErrTextConsumerCredentialValidationFailed = "consumer credential failed validation"
	ErrTextConsumerExists                     = "consumer already exists"
	ErrTextConsumerUsernameEmpty              = "username cannot be empty"
	ErrTextPluginConfigInvalid                = "could not parse plugin configuration"
	ErrTextPluginConfigValidationFailed       = "unable to validate plugin schema"
	ErrTextPluginConfigViolatesSchema         = "plugin failed schema validation: %s"
	ErrTextPluginNameEmpty                    = "plugin name cannot be empty"
	ErrTextPluginSecretConfigUnretrievable    = "could not load secret plugin configuration"
	ErrTextPluginUsesBothConfigTypes          = "plugin cannot use both Config and ConfigFrom"
)

const (
	ErrTextCantRetrieveGatewayClass    = "gatewayclass for this gateway could not be retrieved"
	ErrTextInvalidGatewayConfiguration = "gateway metadata and/or spec are invalid"
)
