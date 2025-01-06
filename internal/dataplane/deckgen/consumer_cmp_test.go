package deckgen

import (
	"sort"
	"testing"

	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/assert"
)

func TestConsumerCmp(t *testing.T) {
	testCases := []struct {
		name     string
		input    []file.FConsumer
		expected []file.FConsumer
	}{
		{
			name: "sort by username",
			input: []file.FConsumer{
				{
					Consumer: kong.Consumer{
						Username: kong.String("b"),
					},
				},
				{
					Consumer: kong.Consumer{
						Username: kong.String("a"),
					},
				},
			},
			expected: []file.FConsumer{
				{
					Consumer: kong.Consumer{
						Username: kong.String("a"),
					},
				},
				{
					Consumer: kong.Consumer{
						Username: kong.String("b"),
					},
				},
			},
		},
		{
			name: "sort by custom_id",
			input: []file.FConsumer{
				{
					Consumer: kong.Consumer{
						CustomID: kong.String("b"),
					},
				},
				{
					Consumer: kong.Consumer{
						CustomID: kong.String("a"),
					},
				},
			},
			expected: []file.FConsumer{
				{
					Consumer: kong.Consumer{
						CustomID: kong.String("a"),
					},
				},
				{
					Consumer: kong.Consumer{
						CustomID: kong.String("b"),
					},
				},
			},
		},
		{
			name: "sort by username and custom_id",
			input: []file.FConsumer{
				{
					Consumer: kong.Consumer{
						Username: kong.String("b"),
						CustomID: kong.String("b"),
					},
				},
				{
					Consumer: kong.Consumer{
						Username: kong.String("a"),
						CustomID: kong.String("a"),
					},
				},
			},
			expected: []file.FConsumer{
				{
					Consumer: kong.Consumer{
						Username: kong.String("a"),
						CustomID: kong.String("a"),
					},
				},
				{
					Consumer: kong.Consumer{
						Username: kong.String("b"),
						CustomID: kong.String("b"),
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sort.Sort(fConsumerByUsernameAndCustomID(tc.input))
			assert.Equal(t, tc.expected, tc.input)
		})
	}
}
