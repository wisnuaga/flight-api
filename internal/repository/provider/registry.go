package provider

import "github.com/wisnuaga/flight-api/internal/repository/provider/garuda"

type RegistryConfig struct {
	EnabledProviders []string
	MockPath         map[string]string
}

type Registry struct {
	cfg RegistryConfig
}

func NewRegistry(cfg RegistryConfig) *Registry {
	return &Registry{cfg: cfg}
}

func (r *Registry) GetProviders() []FlightProvider {
	providers := []FlightProvider{}
	for _, name := range r.cfg.EnabledProviders {
		providers = append(providers, r.getProvider(name))
	}
	return providers
}

func (r *Registry) getProvider(name string) FlightProvider {
	switch name {
	case "garuda":
		return garuda.NewClient(r.cfg.MockPath["garuda"])
	default:
		return nil
	}
}
