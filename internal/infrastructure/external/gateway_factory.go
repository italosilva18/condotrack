package external

import (
	"fmt"

	"github.com/condotrack/api/internal/domain/gateway"
)

// GatewayFactory manages payment gateway instances and selects the active one.
type GatewayFactory struct {
	gateways      map[string]gateway.PaymentGateway
	activeGateway string
}

// NewGatewayFactory creates a new gateway factory.
func NewGatewayFactory() *GatewayFactory {
	return &GatewayFactory{
		gateways: make(map[string]gateway.PaymentGateway),
	}
}

// Register adds a gateway implementation to the factory.
func (f *GatewayFactory) Register(gw gateway.PaymentGateway) {
	f.gateways[gw.Name()] = gw
}

// SetActive sets which gateway is the default for new payments.
func (f *GatewayFactory) SetActive(name string) error {
	if _, ok := f.gateways[name]; !ok {
		return fmt.Errorf("gateway %q not registered", name)
	}
	f.activeGateway = name
	return nil
}

// GetActive returns the currently active gateway.
func (f *GatewayFactory) GetActive() gateway.PaymentGateway {
	if gw, ok := f.gateways[f.activeGateway]; ok {
		return gw
	}
	// Fallback: return any registered gateway
	for _, gw := range f.gateways {
		return gw
	}
	return nil
}

// Get returns a specific gateway by name.
func (f *GatewayFactory) Get(name string) (gateway.PaymentGateway, error) {
	gw, ok := f.gateways[name]
	if !ok {
		return nil, fmt.Errorf("gateway %q not registered", name)
	}
	return gw, nil
}

// ListRegistered returns the names of all registered gateways.
func (f *GatewayFactory) ListRegistered() []string {
	names := make([]string, 0, len(f.gateways))
	for name := range f.gateways {
		names = append(names, name)
	}
	return names
}
