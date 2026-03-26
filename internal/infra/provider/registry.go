package provider

import (
	"github.com/wisnuaga/flight-api/internal/config"
	"github.com/wisnuaga/flight-api/internal/infra/provider/airasia"
	"github.com/wisnuaga/flight-api/internal/infra/provider/batikair"
	"github.com/wisnuaga/flight-api/internal/infra/provider/garuda"
	"github.com/wisnuaga/flight-api/internal/infra/provider/lionair"
	"github.com/wisnuaga/flight-api/internal/port"
)

// Registry builds and returns the list of enabled FlightProvider adapters.
type Registry struct {
	cfg *config.Config
}

func NewRegistry(cfg *config.Config) *Registry {
	return &Registry{cfg: cfg}
}

func (r *Registry) GetProviders() []port.FlightProvider {
	providers := []port.FlightProvider{}
	for _, name := range r.cfg.Providers {
		if p := r.getProvider(name); p != nil {
			providers = append(providers, p)
		}
	}
	return providers
}

func (r *Registry) getProvider(name string) port.FlightProvider {
	switch name {
	case "garuda":
		return garuda.NewClient(r.cfg.GarudaConfig.MockPath)
	case "lionair":
		return lionair.NewClient(r.cfg.LionAirConfig.MockPath)
	case "batikair":
		return batikair.NewClient(r.cfg.BatikAirConfig.MockPath)
	case "airasia":
		return airasia.NewClient(r.cfg.AirAsiaConfig.MockPath)
	default:
		return nil
	}
}
