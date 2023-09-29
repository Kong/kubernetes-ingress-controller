package admission

import (
	"context"
	"fmt"

	"github.com/blang/semver/v4"
	"github.com/kong/go-kong/kong"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/annotations"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/versions"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1beta1"
)

func (validator KongHTTPValidator) KongConsumerGroup() CustomValidatorAdapter {
	return CustomValidatorAdapter{
		validateCreate: func(ctx context.Context, obj runtime.Object) (bool, string, error) {
			consumerGroup, ok := obj.(*kongv1beta1.KongConsumerGroup)
			if !ok {
				return false, "", fmt.Errorf("unexpected type, expected *kongv1beta1.KongConsumerGroup, got %T", obj)
			}
			return validator.ValidateConsumerGroup(ctx, *consumerGroup)
		},
	}
}

func (validator KongHTTPValidator) ValidateConsumerGroup(
	ctx context.Context,
	consumerGroup kongv1beta1.KongConsumerGroup,
) (bool, string, error) {
	// Ignore ConsumerGroups that are being managed by another controller.
	if !validator.ingressClassMatcher(&consumerGroup.ObjectMeta, annotations.IngressClassKey, annotations.ExactClassMatch) {
		return true, "", nil
	}

	// Consumer groups work only for Kong Enterprise >=3.4.
	infoSvc, ok := validator.AdminAPIServicesProvider.GetInfoService()
	if !ok {
		return true, "", nil
	}
	info, err := infoSvc.Get(ctx)
	if err != nil {
		validator.Logger.Debugf("failed to fetch Kong info: %v", err)
		return false, ErrTextAdminAPIUnavailable, nil
	}
	version, err := kong.NewVersion(info.Version)
	if err != nil {
		validator.Logger.Debugf("failed to parse Kong version: %v", err)
	} else {
		kongVer := semver.Version{Major: version.Major(), Minor: version.Minor()}
		if !version.IsKongGatewayEnterprise() || !kongVer.GTE(versions.ConsumerGroupsVersionCutoff) {
			return false, ErrTextConsumerGroupUnsupported, nil
		}
	}

	cgs, ok := validator.AdminAPIServicesProvider.GetConsumerGroupsService()
	if !ok {
		return true, "", nil
	}
	// This check forbids consumer group creation if the license is invalid or missing.
	// There is no other way to robustly check the validity of a license than actually trying an enterprise feature.
	if _, _, err := cgs.List(ctx, &kong.ListOpt{Size: 0}); err != nil {
		switch {
		case kong.IsNotFoundErr(err):
			// This is the case when consumer group is not supported (Kong OSS) and previous version
			// check (if !version.IsKongGatewayEnterprise()) has been omitted due to a parsing error.
			return false, ErrTextConsumerGroupUnsupported, nil
		case kong.IsForbiddenErr(err):
			return false, ErrTextConsumerGroupUnlicensed, nil
		default:
			return false, fmt.Sprintf("%s: %s", ErrTextConsumerGroupUnexpected, err), nil
		}
	}
	return true, "", nil
}
