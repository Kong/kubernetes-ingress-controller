package admission

const (
	ErrTextConsumerUsernameEmpty           = "username cannot be empty"
	ErrTextConsumerUnretrievable           = "failed to fetch consumer from kong"
	ErrTextConsumerExists                  = "consumer already exists"
	ErrTextPluginNameEmpty                 = "plugin name cannot be empty"
	ErrTextPluginConfigInvalid             = "could not parse plugin configuration"
	ErrTextPluginUsesBothConfigTypes       = "plugin cannot use both Config and ConfigFrom"
	ErrTextPluginConfigViolatesSchema      = "plugin failed schema validation"
	ErrTextPluginSecretConfigUnretrievable = "could not load secret plugin configuration"
)
