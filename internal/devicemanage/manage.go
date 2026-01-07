package devicemanage

import (
	"errors"
	"log"
	"sync"
)

type DeviceManager struct {
	lock     sync.RWMutex
	workers  map[string]*DeviceWorker
	resultCh chan DeviceResult
}

func NewDeviceManager() *DeviceManager {
	return &DeviceManager{
		workers:  make(map[string]*DeviceWorker),
		resultCh: make(chan DeviceResult, 100),
	}
}

func (m *DeviceManager) Start() {
	go func() {
		for res := range m.resultCh {
			log.Printf("[RESULT] %s data: %v", res.DeviceID, res.Data)
		}
	}()
}

func (m *DeviceManager) AddOrUpdateDevice(d Device) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if w, exists := m.workers[d.ID]; exists {
		w.Stop()
		delete(m.workers, d.ID)
	}

	worker := NewDeviceWorker(d, m.resultCh)
	worker.Start()
	m.workers[d.ID] = worker
}

func (m *DeviceManager) RemoveDevice(id string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if w, exists := m.workers[id]; exists {
		w.Stop()
		delete(m.workers, id)
	}
}

func (m *DeviceManager) Shutdown() {
	m.lock.Lock()
	defer m.lock.Unlock()
	for id, w := range m.workers {
		w.Stop()
		delete(m.workers, id)
	}
	close(m.resultCh)
}

func (m *DeviceManager) WriteSingleCoil(deviceId string, address, value uint16) error {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if w, exists := m.workers[deviceId]; exists {
		return w.writeSingleCoil(address, value)
	}
	return errors.New("device not found")
}

func (m *DeviceManager) WriteSingleRegister(deviceId string, address, value uint16) error {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if w, exists := m.workers[deviceId]; exists {
		return w.writeSingleRegister(address, value)
	}
	return errors.New("device not found")
}

func (m *DeviceManager) WriteMultipleRegisters(deviceId string, address, quantity uint16, value []byte) error {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if w, exists := m.workers[deviceId]; exists {
		return w.writeMultipleRegisters(address, quantity, value)
	}
	return errors.New("device not found")
}
