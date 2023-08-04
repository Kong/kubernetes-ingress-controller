package sendconfig_test

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/kong/deck/file"
	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/sendconfig"
)

func TestDBLessConfigMarshalToJSON(t *testing.T) {
	dblessConfig := sendconfig.DBLessConfig{
		Content: file.Content{
			Services: []file.FService{
				{
					Service: kong.Service{
						Name: kong.String("service-id"),
					},
				},
			},
		},
		ConsumerGroupConsumerRelationships: []sendconfig.ConsumerGroupConsumerRelationship{
			{
				ConsumerGroup: "cg1",
				Consumer:      "c1",
			},
		},
		ConsumerGroupPluginRelationships: []sendconfig.ConsumerGroupPluginRelationship{
			{
				ConsumerGroup: "cg1",
				Plugin:        "p1",
			},
		},
	}

	expected := `{
  "services": [
    {
      "name": "service-id"
    }
  ],
  "consumer_group_consumers": [
    {
      "consumer_group": "cg1",
      "consumer": "c1"
    }
  ],
  "consumer_group_plugins": [
    {
      "consumer_group": "cg1",
      "plugin": "p1"
    }
  ]
}`
	b, err := json.Marshal(dblessConfig)
	require.NoError(t, err)
	require.JSONEq(t, expected, string(b))
}

func TestDefaultContentToDBLessConfigConverter(t *testing.T) {
	converter := sendconfig.DefaultContentToDBLessConfigConverter{}

	testCases := []struct {
		name                 string
		content              *file.Content
		expectedDBLessConfig sendconfig.DBLessConfig
	}{
		{
			name:    "empty content",
			content: &file.Content{},
			expectedDBLessConfig: sendconfig.DBLessConfig{
				Content: file.Content{},
			},
		},
		{
			name: "content with info",
			content: &file.Content{
				Info: &file.Info{
					SelectorTags: []string{"tag1", "tag2"},
				},
			},
			expectedDBLessConfig: sendconfig.DBLessConfig{
				Content: file.Content{},
			},
		},
		{
			name: "content with consumer group consumers and plugins",
			content: &file.Content{
				ConsumerGroups: []file.FConsumerGroupObject{
					{
						ConsumerGroup: kong.ConsumerGroup{
							Name: kong.String("cg1"),
						},
						Consumers: []*kong.Consumer{{Username: kong.String("c1")}},
						Plugins:   []*kong.ConsumerGroupPlugin{{Name: kong.String("p1")}},
					},
				},
				Consumers: []file.FConsumer{
					{
						Consumer: kong.Consumer{
							Username: kong.String("c1"),
						},
						Groups: []*kong.ConsumerGroup{{ID: kong.String("cg1"), Name: kong.String("cg1")}},
					},
				},
				Plugins: []file.FPlugin{
					{
						Plugin: kong.Plugin{
							Name:          kong.String("p1"),
							ConsumerGroup: &kong.ConsumerGroup{ID: kong.String("cg1"), Name: kong.String("cg1")},
						},
					},
				},
			},
			expectedDBLessConfig: sendconfig.DBLessConfig{
				Content: file.Content{
					ConsumerGroups: []file.FConsumerGroupObject{
						{
							ConsumerGroup: kong.ConsumerGroup{
								Name: kong.String("cg1"),
							},
						},
					},
					Consumers: []file.FConsumer{
						{
							Consumer: kong.Consumer{
								Username: kong.String("c1"),
							},
						},
					},
					Plugins: []file.FPlugin{
						{
							Plugin: kong.Plugin{
								Name: kong.String("p1"),
							},
						},
					},
				},
				ConsumerGroupConsumerRelationships: []sendconfig.ConsumerGroupConsumerRelationship{
					{
						ConsumerGroup: "cg1",
						Consumer:      "c1",
					},
				},
				ConsumerGroupPluginRelationships: []sendconfig.ConsumerGroupPluginRelationship{
					{
						ConsumerGroup: "cg1",
						Plugin:        "p1",
					},
				},
			},
		},
		{
			name: "content with plugin config nulls",
			content: &file.Content{
				Plugins: []file.FPlugin{
					{
						Plugin: kong.Plugin{
							Name: kong.String("p1"),
							Config: kong.Configuration{
								"config1": nil,
								"config2": "value2",
							},
						},
					},
				},
				Consumers: []file.FConsumer{
					{
						Consumer: kong.Consumer{
							Username: kong.String("c1"),
						},
						Plugins: []*file.FPlugin{
							{
								Plugin: kong.Plugin{
									Name: kong.String("p1"),
									Config: kong.Configuration{
										"config1": nil,
										"config2": "value2",
									},
								},
							},
						},
					},
				},
				Routes: []file.FRoute{
					{
						Route: kong.Route{
							Name: kong.String("r1"),
						},
						Plugins: []*file.FPlugin{
							{
								Plugin: kong.Plugin{
									Name: kong.String("p1"),
									Config: kong.Configuration{
										"config1": nil,
										"config2": "value2",
									},
								},
							},
						},
					},
				},
				Services: []file.FService{
					{
						Service: kong.Service{
							Name: kong.String("s1"),
						},
						Plugins: []*file.FPlugin{
							{
								Plugin: kong.Plugin{
									Name: kong.String("p1"),
									Config: kong.Configuration{
										"config1": nil,
										"config2": "value2",
									},
								},
							},
						},
					},
				},
			},
			expectedDBLessConfig: sendconfig.DBLessConfig{
				Content: file.Content{
					Plugins: []file.FPlugin{
						{
							Plugin: kong.Plugin{
								Name: kong.String("p1"),
								Config: kong.Configuration{
									"config2": "value2",
								},
							},
						},
					},
					Consumers: []file.FConsumer{
						{
							Consumer: kong.Consumer{
								Username: kong.String("c1"),
							},
							Plugins: []*file.FPlugin{
								{
									Plugin: kong.Plugin{
										Name: kong.String("p1"),
										Config: kong.Configuration{
											"config2": "value2",
										},
									},
								},
							},
						},
					},
					Routes: []file.FRoute{
						{
							Route: kong.Route{
								Name: kong.String("r1"),
							},
							Plugins: []*file.FPlugin{
								{
									Plugin: kong.Plugin{
										Name: kong.String("p1"),
										Config: kong.Configuration{
											"config2": "value2",
										},
									},
								},
							},
						},
					},
					Services: []file.FService{
						{
							Service: kong.Service{
								Name: kong.String("s1"),
							},
							Plugins: []*file.FPlugin{
								{
									Plugin: kong.Plugin{
										Name: kong.String("p1"),
										Config: kong.Configuration{
											"config2": "value2",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	// deepCopyContent makes a deep copy of the given file.Content. We use gob to do this, instead of e.g. json, because
	// various types inside *file.Content implement custom json marshalling/unmarshalling, which would break the test
	// because of some fields being dropped in this process (e.g. FPlugin.MarshalJSON drops all its fields but ID).
	deepCopyContent := func(t *testing.T, content *file.Content) *file.Content {
		var originalContentBlob bytes.Buffer
		enc := gob.NewEncoder(&originalContentBlob)
		err := enc.Encode(content)
		require.NoError(t, err)

		copiedContent := &file.Content{}
		err = gob.NewDecoder(&originalContentBlob).Decode(copiedContent)
		require.NoError(t, err)
		return copiedContent
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Preserve original content to make sure it's not modified after conversion.
			originalContent := deepCopyContent(t, tc.content)

			dblessConfig := converter.Convert(tc.content)
			require.Equal(t, tc.expectedDBLessConfig, dblessConfig)

			// We're fine with Config being modified in the original *file.Content - we're just dropping null values.
			ignorePluginConfig := cmpopts.IgnoreFields(kong.Plugin{}, "Config")
			require.Empty(t, cmp.Diff(originalContent, tc.content, ignorePluginConfig),
				"passed *file.Content should not be modified",
			)
		})
	}
}

func BenchmarkDefaultContentToDBLessConfigConverter_Convert(b *testing.B) {
	content := &file.Content{
		Info: &file.Info{
			SelectorTags: []string{"tag1", "tag2"},
		},
		ConsumerGroups: []file.FConsumerGroupObject{
			{
				ConsumerGroup: kong.ConsumerGroup{
					Name: kong.String("cg1"),
				},
				Consumers: []*kong.Consumer{{Username: kong.String("c1")}},
				Plugins:   []*kong.ConsumerGroupPlugin{{Name: kong.String("p1")}},
			},
		},
		Consumers: []file.FConsumer{
			{
				Consumer: kong.Consumer{
					Username: kong.String("c1"),
				},
				Groups: []*kong.ConsumerGroup{{Name: kong.String("cg1")}},
			},
		},
		Plugins: []file.FPlugin{
			{
				Plugin: kong.Plugin{
					Name:          kong.String("p1"),
					ConsumerGroup: &kong.ConsumerGroup{Name: kong.String("cg1")},
					Config:        kong.Configuration{"config1": nil},
				},
			},
			{
				Plugin: kong.Plugin{
					Name:          kong.String("p2"),
					ConsumerGroup: &kong.ConsumerGroup{Name: kong.String("cg1")},
					Config:        kong.Configuration{"config1": nil},
				},
			},
			{
				Plugin: kong.Plugin{
					Name:          kong.String("p3"),
					ConsumerGroup: &kong.ConsumerGroup{Name: kong.String("cg1")},
					Config:        kong.Configuration{"config1": nil},
				},
			},
		},
	}

	converter := sendconfig.DefaultContentToDBLessConfigConverter{}
	for i := 0; i < b.N; i++ {
		dblessConfig := converter.Convert(content)
		_ = dblessConfig
	}
}
