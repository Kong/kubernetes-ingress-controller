package sdk

import (
	sdkkonnectgo "github.com/Kong/sdk-konnect-go"
	sdkkonnectcomp "github.com/Kong/sdk-konnect-go/models/components"
)

// SDK returns a new SDK instance.
func SDK(token string, opts ...sdkkonnectgo.SDKOption) *sdkkonnectgo.SDK {
	sdkOpts := []sdkkonnectgo.SDKOption{
		sdkkonnectgo.WithSecurity(
			sdkkonnectcomp.Security{
				PersonalAccessToken: sdkkonnectgo.String(token),
			},
		),
	}
	sdkOpts = append(sdkOpts, opts...)

	return sdkkonnectgo.New(sdkOpts...)
}
