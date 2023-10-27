package admission

const (
	ErrTextAdminAPIUnavailable                = "Could not talk to Kong admin API"
	ErrTextConsumerCredentialSecretNotFound   = "Consumer referenced non-existent credentials secret"
	ErrTextConsumerCredentialValidationFailed = "Consumer credential failed validation"
	ErrTextConsumerExists                     = "Consumer already exists"
	ErrTextConsumerUnretrievable              = "Failed to fetch consumer from kong"
	ErrTextConsumerGroupUnsupported           = "Consumer group support requires Kong Enterprise"
	ErrTextConsumerGroupUnlicensed            = "Consumer group support requires a valid Kong Enterprise license"
	ErrTextConsumerGroupUnexpected            = "Unexpected error during checking support for consumer group"
	ErrTextConsumerUsernameEmpty              = "Username cannot be empty"
	//nolint:revive
	ErrTextFailedToRetrieveSecret       = "Could not retrieve secrets from the kubernetes API" //nolint:gosec // Ignore G101 error
	ErrTextPluginConfigInvalid          = "Could not parse plugin configuration"
	ErrTextPluginConfigValidationFailed = "Unable to validate plugin schema"
	ErrTextPluginConfigViolatesSchema   = "Plugin failed schema validation: %s"
	ErrTextPluginNameEmpty              = "Plugin name cannot be empty"
	//nolint:revive
	ErrTextPluginSecretConfigUnretrievable = "Could not load secret plugin configuration" //nolint:gosec // Ignore G101 error
	ErrTextPluginUsesBothConfigTypes       = "Plugin cannot use both Config and ConfigFrom"
)

const (
	ErrTextCantRetrieveGatewayClass    = "Gatewayclass for this gateway could not be retrieved"
	ErrTextInvalidGatewayConfiguration = "Gateway metadata and/or spec are invalid"
)
