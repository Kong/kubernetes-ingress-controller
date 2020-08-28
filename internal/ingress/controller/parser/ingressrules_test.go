package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeIngressRules(t *testing.T) {
	for _, tt := range []struct {
		name       string
		inputs     []*ingressRules
		wantOutput *ingressRules
	}{
		{
			name: "empty list",
			wantOutput: &ingressRules{
				SecretNameToSNIs:      map[string][]string{},
				ServiceNameToServices: map[string]Service{},
			},
		},
		{
			name: "nil maps",
			inputs: []*ingressRules{
				{}, {}, {},
			},
			wantOutput: &ingressRules{
				SecretNameToSNIs:      map[string][]string{},
				ServiceNameToServices: map[string]Service{},
			},
		},
		{
			name: "one input",
			inputs: []*ingressRules{
				{
					SecretNameToSNIs:      map[string][]string{"a": {"b", "c"}, "d": {"e", "f"}},
					ServiceNameToServices: map[string]Service{"1": {Namespace: "potato"}},
				},
			},
			wantOutput: &ingressRules{
				SecretNameToSNIs:      map[string][]string{"a": {"b", "c"}, "d": {"e", "f"}},
				ServiceNameToServices: map[string]Service{"1": {Namespace: "potato"}},
			},
		},
		{
			name: "three inputs",
			inputs: []*ingressRules{
				{
					SecretNameToSNIs:      map[string][]string{"a": {"b", "c"}, "d": {"e", "f"}},
					ServiceNameToServices: map[string]Service{"1": {Namespace: "potato"}},
				},
				{
					SecretNameToSNIs: map[string][]string{"g": {"h"}},
				},
				{
					ServiceNameToServices: map[string]Service{"2": {Namespace: "carrot"}},
				},
			},
			wantOutput: &ingressRules{
				SecretNameToSNIs:      map[string][]string{"a": {"b", "c"}, "d": {"e", "f"}, "g": {"h"}},
				ServiceNameToServices: map[string]Service{"1": {Namespace: "potato"}, "2": {Namespace: "carrot"}},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			gotOutput := mergeIngressRules(tt.inputs...)
			assert.Equal(t, &gotOutput, tt.wantOutput)
		})
	}
}
