// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package nsxtreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/nsxtreceiver"

import (
	"errors"
	"fmt"
	"net/url"

	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/receiver/scraperhelper"
	"go.uber.org/multierr"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/nsxtreceiver/internal/metadata"
)

// Config is the configuration for the NSX receiver
type Config struct {
	scraperhelper.ScraperControllerSettings `mapstructure:",squash"`
	confighttp.HTTPClientSettings           `mapstructure:",squash"`
	metadata.MetricsBuilderConfig           `mapstructure:",squash"`
	ResourceAttributes                      metadata.ResourceAttributesConfig `mapstructure:",squash"`
	Username                                string                            `mapstructure:"username"`
	Password                                string                            `mapstructure:"password"`
}

// Validate returns if the NSX configuration is valid
func (c *Config) Validate() error {
	var err error
	if c.Endpoint == "" {
		err = multierr.Append(err, errors.New("no manager endpoint was specified"))
		return err
	}

	res, err := url.Parse(c.Endpoint)
	if err != nil {
		err = multierr.Append(err, fmt.Errorf("unable to parse url %s: %w", c.Endpoint, err))
		return err
	}

	if res.Scheme != "http" && res.Scheme != "https" {
		err = multierr.Append(err, errors.New("url scheme must be http or https"))
	}

	if c.Username == "" {
		err = multierr.Append(err, errors.New("username not provided and is required"))
	}

	if c.Password == "" {
		err = multierr.Append(err, errors.New("password not provided and is required"))
	}
	return err
}
