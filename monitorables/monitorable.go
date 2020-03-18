package monitorables

import (
	"github.com/monitoror/monitoror/service"
)

type Monitorable interface {
	GetVariants() []string

	IsValid(variant string) bool
	Enable(variant string)
}

type (
	manager struct {
		store *service.Store

		monitorables []Monitorable
	}
)

func NewMonitorableManager(store *service.Store) *manager {
	manager := &manager{store: store}

	// Register all monitorables
	manager.init()

	return manager
}

func (m *manager) register(monitorable Monitorable) {
	m.monitorables = append(m.monitorables, monitorable)
}

func (m *manager) EnableMonitorables() {
	//TODO: LOGS

	for _, monitorable := range m.monitorables {
		for _, variant := range monitorable.GetVariants() {
			monitorable.Enable(variant)
		}
	}
}
