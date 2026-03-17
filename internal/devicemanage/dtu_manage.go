package devicemanage

import (
	"log"
	"sync"
)

type DtuDeviceManager struct {
	lock     sync.RWMutex
	workers  map[string]*DtuDeviceWorker
	resultCh chan DeviceResult
}

func NewDtuDeviceManager() *DtuDeviceManager {
	return &DtuDeviceManager{
		workers:  make(map[string]*DtuDeviceWorker),
		resultCh: make(chan DeviceResult, 100),
	}
}

func (m *DtuDeviceManager) Start() {
	go func() {
		for res := range m.resultCh {
			log.Printf("[RESULT] %s data: %v", res.DeviceID, res.Data)
		}
	}()
}

func (m *DtuDeviceManager) AddOrUpdateDevice(d Device) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if w, exists := m.workers[d.ID]; exists {
		w.Stop()
		delete(m.workers, d.ID)
	}

	worker := NewDtuDeviceWorker(d, m.resultCh)
	worker.Start()
	m.workers[d.ID] = worker
}
