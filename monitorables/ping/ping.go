package ping

import "github.com/monitoror/monitoror/service"

type Monitorable struct {
	store *service.Store
}

func NewMonitorable(store *service.Store) *Monitorable {
	return &Monitorable{store: store}
}

func (m *Monitorable) GetVariants() []string {
	return []string{}
}

func (m *Monitorable) IsValid(variant string) bool {
	return true
}

func (m *Monitorable) Enable(variant string) {
}
