/*
Copyright 2023 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package environment

import (
	"context"
)

type EventsHubConfig struct {
	TLS struct {
		Enabled bool
	}
}

type eventsHubConfigKey struct{}

func withEventsHubConfig(ctx context.Context, config *EventsHubConfig) context.Context {
	return context.WithValue(ctx, eventsHubConfigKey{}, config)
}

// GetEventsHubConfig returns the configured EventsHubConfig
func GetEventsHubConfig(ctx context.Context) *EventsHubConfig {
	config := ctx.Value(eventsHubConfigKey{})
	if config == nil {
		return &EventsHubConfig{}
	}
	return config.(*EventsHubConfig)
}

func initEventsHubFlags() ConfigurationOption {
	return func(configuration Configuration) Configuration {
		fs := configuration.Flags.Get(configuration.Context)

		cfg := &EventsHubConfig{}

		fs.BoolVar(&cfg.TLS.Enabled, "eventshub.tls.enabled", false, "Enable testing with EventsHub TLS enabled")

		configuration.Context = withEventsHubConfig(configuration.Context, cfg)

		return configuration
	}
}
