// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jaegerthrifthttpexporter

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/config/configopaque"
	"go.opentelemetry.io/collector/confmap/confmaptest"
)

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "config.yaml"))
	require.NoError(t, err)

	// URL doesn't have a default value so set it directly.
	defaultCfg := createDefaultConfig().(*Config)
	defaultCfg.Endpoint = "http://jaeger.example:14268/api/traces"

	tests := []struct {
		id       component.ID
		expected component.Config
	}{
		{
			id:       component.NewIDWithName(typeStr, ""),
			expected: defaultCfg,
		},
		{
			id: component.NewIDWithName(typeStr, "2"),
			expected: &Config{
				HTTPClientSettings: confighttp.HTTPClientSettings{
					Endpoint: "http://jaeger.example.com/api/traces",
					Headers: map[string]configopaque.String{
						"added-entry": "added value",
						"dot.test":    "test",
					},
					Timeout: 2 * time.Second,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.id.String(), func(t *testing.T) {
			factory := NewFactory()
			cfg := factory.CreateDefaultConfig()

			sub, err := cm.Sub(tt.id.String())
			require.NoError(t, err)
			require.NoError(t, component.UnmarshalConfig(sub, cfg))

			assert.NoError(t, component.ValidateConfig(cfg))
			assert.Equal(t, tt.expected, cfg)
		})
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name         string
		config       *Config
		errorMessage string
	}{
		{
			name:         "empty_url",
			config:       &Config{},
			errorMessage: "invalid \"endpoint\": parse \"\": empty url",
		},
		{
			name: "invalid_url",
			config: &Config{
				HTTPClientSettings: confighttp.HTTPClientSettings{
					Endpoint: ".example:123",
				},
			},
			errorMessage: "invalid \"endpoint\": parse \".example:123\": invalid URI for request",
		},
		{
			name: "negative_duration",
			config: &Config{
				HTTPClientSettings: confighttp.HTTPClientSettings{
					Endpoint: "example.com:123",
					Timeout:  -2 * time.Second,
				},
			},
			errorMessage: "invalid negative value for \"timeout\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.EqualError(t, tt.config.Validate(), tt.errorMessage)
		})
	}
}
