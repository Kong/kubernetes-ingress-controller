package sendconfig

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"testing"

	"github.com/kong/deck/file"
	"github.com/kong/go-kong/kong"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

type configServiceMock struct {
	lock   sync.RWMutex
	config map[string]any
}

func (m *configServiceMock) ReloadDeclarativeRawConfig(
	ctx context.Context,
	config io.Reader,
	checkHash bool,
	flattenErrors bool,
) ([]byte, error) {
	var mockRetBody = []byte("{}")
	rawConfig, err := io.ReadAll(config)
	if err != nil {
		return mockRetBody, err
	}

	cfg := map[string]any{}
	err = json.Unmarshal(rawConfig, &cfg)
	if err != nil {
		return mockRetBody, err
	}

	m.lock.Lock()
	defer m.lock.Unlock()
	m.config = cfg
	return mockRetBody, nil
}

func (m *configServiceMock) getCurrentState() (*file.Content, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	buf, err := json.Marshal(m.config)
	if err != nil {
		return nil, err
	}

	c := &file.Content{}
	if err = json.Unmarshal(buf, c); err != nil {
		return nil, err
	}
	return c, nil
}

func extractPluginInFileContentByID(c *file.Content, id string) (*file.FPlugin, error) {
	for i, plugin := range c.Plugins {
		if plugin.ID != nil && *(plugin.ID) == id {
			return &c.Plugins[i], nil
		}
	}
	return nil, fmt.Errorf("plugin not found")
}

func TestUpdateInMemoryPlugin(t *testing.T) {
	m := &configServiceMock{}
	s := NewUpdateStrategyInMemory(m, false, logrus.New())
	ctx := context.Background()

	testCases := []struct {
		name   string
		plugin file.FPlugin
	}{
		{
			name: "plugin with nulls",
			plugin: file.FPlugin{
				Plugin: kong.Plugin{
					ID:   kong.String("plugin-test"),
					Name: kong.String("plugin-test"),
					Config: kong.Configuration{
						"key1": "value1",
						"key2": nil,
					},
				},
				ConfigSource: nil,
			},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			targetState := &file.Content{
				Plugins: []file.FPlugin{tc.plugin},
			}
			_, _, err := s.Update(ctx, targetState)
			require.NoError(t, err)

			currentState, err := m.getCurrentState()
			require.NoError(t, err)
			plugin, err := extractPluginInFileContentByID(currentState, *tc.plugin.ID)
			require.NoError(t, err)
			require.Equal(t, tc.plugin.Config, plugin.Config)
		})
	}
}
