package admission

const (
	ErrTextAdminAPIUnavailable                = "could not talk to Kong admin API"
	ErrTextConsumerCredentialSecretNotFound   = "consumer referenced non-existent credentials secret"
	ErrTextConsumerCredentialValidationFailed = "consumer credential failed validation"
	ErrTextConsumerExists                     = "consumer already exists"
	ErrTextConsumerUnretrievable              = "failed to fetch consumer from kong"
	ErrTextConsumerGroupUnsupported           = "consumer group support requires Kong Enterprise"
	ErrTextConsumerGroupUnlicensed            = "consumer group support requires a valid Kong Enterprise license"
	ErrTextConsumerGroupUnexpected            = "unexpected error during checking support for consumer group"
	ErrTextCustomEntityFieldsUnmarshalFailed  = "failed to unmarshal fields of custom entity: %v"
	ErrTextCustomEntityGetSchemaFailed        = "failed to get schema of Kong entity type '%s': %v"
	ErrTextFailedToRetrieveSecret             = "could not retrieve secrets from the kubernetes API" //nolint:revive,gosec
	ErrTextPluginConfigInvalid                = "could not parse plugin configuration"
	ErrTextPluginConfigValidationFailed       = "unable to validate plugin schema"
	ErrTextPluginConfigViolatesSchema         = "plugin failed schema validation: %s"
	ErrTextPluginSecretConfigUnretrievable    = "could not load secret plugin configuration"
	ErrTextVaultConfigUnmarshalFailed         = "failed to unmarshal vault configuration: %v"
	ErrTextVaultUnableToValidate              = "unable to validate vault on Kong gateway"
	ErrTextVaultConfigValidationResultInvalid = "vault configuration in invalid: %s"
)

const (
	ErrTextCantRetrieveGatewayClass    = "gatewayclass for this gateway could not be retrieved"
	ErrTextInvalidGatewayConfiguration = "gateway metadata and/or spec are invalid"
)
