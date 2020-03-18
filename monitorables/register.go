package monitorables

import "github.com/monitoror/monitoror/monitorables/ping"

func (m *manager) init() {
	m.register(ping.NewMonitorable(m.store))
}
