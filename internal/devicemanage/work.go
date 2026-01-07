package devicemanage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/winc-link/hummingbird-modbus-tcp/constant"
	"github.com/winc-link/hummingbird-modbus-tcp/dtos"
	"github.com/winc-link/hummingbird-modbus-tcp/utils/cast"
	"github.com/winc-link/hummingbird-sdk-go/model"
	"github.com/winc-link/hummingbird-sdk-go/service"
	"strings"

	"strconv"
	"sync"
	"time"

	"github.com/goburrow/modbus"
)

var SDKDriver *service.DriverService

type Device struct {
	ID        string
	IP        string
	Port      int
	ProductId string
	Period    time.Duration // 采集周期
	SlaveId   int
}

type DeviceStatus struct {
	LastSuccess time.Time
	LastError   string
	IsOnline    bool
	Lock        sync.RWMutex
}

type DeviceResult struct {
	DeviceID string
	Data     []byte
	Err      error
	Time     time.Time
}

type DeviceWorker struct {
	device   Device
	status   *DeviceStatus
	ctx      context.Context
	cancel   context.CancelFunc
	handler  *modbus.TCPClientHandler
	client   modbus.Client
	resultCh chan<- DeviceResult
}

func NewDeviceWorker(device Device, resultCh chan<- DeviceResult) *DeviceWorker {
	ctx, cancel := context.WithCancel(context.Background())
	return &DeviceWorker{
		device:   device,
		status:   &DeviceStatus{},
		ctx:      ctx,
		cancel:   cancel,
		resultCh: resultCh,
	}
}

func (w *DeviceWorker) connect() {
	address := fmt.Sprintf("%s:%d", w.device.IP, w.device.Port)
	h := modbus.NewTCPClientHandler(address)
	h.Timeout = 3 * time.Second
	h.SlaveId = byte(w.device.SlaveId)
	if err := h.Connect(); err != nil {
		w.updateStatus(true, err.Error())
	}
	w.updateStatus(true, "")
	w.handler = h
	w.client = modbus.NewClient(h)
}

func (w *DeviceWorker) Start() {

	go func() {
		w.connect()

		ticker := time.NewTicker(w.device.Period)
		defer ticker.Stop()

		//address := fmt.Sprintf("%s:%d", w.device.IP, w.device.Port)
		for {
			select {
			case <-ticker.C:
				if w.handler == nil || w.client == nil {
					w.connect()
				}
				//if w.handler == nil || w.client == nil {
				//	h := modbus.NewTCPClientHandler(address)
				//	h.Timeout = 3 * time.Second
				//	h.SlaveId = byte(w.device.SlaveId)
				//	if err := h.Connect(); err != nil {
				//		w.updateStatus(false, err.Error())
				//		continue
				//	}
				//	SDKDriver.Online(w.device.ID)
				//	w.handler = h
				//	w.client = modbus.NewClient(h)
				//}
				//采集设备数据然后上报。
				w.collectDataAndReport()
			case <-w.ctx.Done():
				if w.handler != nil {
					w.handler.Close()
				}
				SDKDriver.GetLogger().Infof("[%s] worker stopped", w.device.ID)
				return
			}
		}
	}()
}

func (w *DeviceWorker) collectDataAndReport() {
	product, ok := SDKDriver.GetProductById(w.device.ProductId)
	if !ok {
		return
	}

	reportData := make(map[string]interface{})
	for _, property := range product.Properties {
		var modbusSpec dtos.ModbusSpec

		err := json.Unmarshal([]byte(property.TypeSpec.Specs), &modbusSpec)
		if err != nil {
			SDKDriver.GetLogger().Errorf("unmarshal property %s error: %v", property.TypeSpec.Specs, err)
			continue
		}
		address, err := strconv.ParseUint(modbusSpec.RegisterAddress, 10, 16)
		if err != nil {
			SDKDriver.GetLogger().Errorf("parse register address %s error: %v", modbusSpec.RegisterAddress, err)
			continue
		}
		var (
			value []byte
		)
		switch modbusSpec.RegisterType {
		case string(constant.CoilStatus):
			value, err = w.client.ReadCoils(uint16(address), cast.GetQuantityByDataType(modbusSpec.DataType))
		case string(constant.InputStatus):
			value, err = w.client.ReadDiscreteInputs(uint16(address), cast.GetQuantityByDataType(modbusSpec.DataType))
		case string(constant.HoldingRegister):
			value, err = w.client.ReadHoldingRegisters(uint16(address), cast.GetQuantityByDataType(modbusSpec.DataType))
		case string(constant.InputRegisters):
			value, err = w.client.ReadInputRegisters(uint16(address), cast.GetQuantityByDataType(modbusSpec.DataType))
		default:
			err = fmt.Errorf("unknown register")
		}

		if err != nil {
			SDKDriver.GetLogger().Errorf("read register %s error: %v", modbusSpec.RegisterAddress, err)
			if strings.Contains(err.Error(), "i/o timeout") {
				w.updateStatus(false, err.Error())
			}
			continue
		}
		v := cast.ConvertByteToInt(modbusSpec.DataType, modbusSpec.DataOrder, modbusSpec.Multiplier, value)

		reportData[property.Code] = v

	}
	if len(reportData) == 0 {
		return
	}
	_, err := SDKDriver.PropertyReport(w.device.ID, model.NewPropertyReport(model.NewDefaultCommonRequest(), reportData))

	if err != nil {
		SDKDriver.GetLogger().Errorf("property report %s error: %v", w.device.ID, err)
	}
}

func (w *DeviceWorker) updateStatus(online bool, errMsg string) {
	w.status.Lock.Lock()
	defer w.status.Lock.Unlock()

	if online {
		if !w.status.IsOnline {
			SDKDriver.Online(w.device.ID)
		}
	} else {
		if w.status.IsOnline {
			SDKDriver.Offline(w.device.ID)
		}
	}
	w.status.LastError = errMsg
	w.status.IsOnline = online

}

func (w *DeviceWorker) writeSingleCoil(address, value uint16) error {
	SDKDriver.GetLogger().Infof("writeSingleCoil address: %v value %v", address, value)
	var v uint16
	if value == 0 {
		v = 0x0000
	} else if value == 1 {
		v = 0xFF00
	}
	if w.client == nil {
		return fmt.Errorf("client is nil")
	}
	_, err := w.client.WriteSingleCoil(address, v)
	if err != nil {
		SDKDriver.GetLogger().Errorf("writeSingleCoil [device Id %v]  error: %v", w.device.ID, err)
		return err
	}
	return nil

}

func (w *DeviceWorker) writeSingleRegister(address, value uint16) error {
	SDKDriver.GetLogger().Infof("writeSingleRegister address: %v value %v", address, value)
	_, err := w.client.WriteSingleRegister(address, value)
	if err != nil {
		SDKDriver.GetLogger().Errorf("writeSingleRegister [device Id %v] error: %v", w.device.ID, err)
		return err
	}
	return nil
}

func (w *DeviceWorker) writeMultipleRegisters(address, quantity uint16, value []byte) error {
	SDKDriver.GetLogger().Infof("WriteMultipleRegisters address: %v value %v", address, value)
	_, err := w.client.WriteMultipleRegisters(address, quantity, value)
	if err != nil {
		SDKDriver.GetLogger().Errorf("WriteMultipleRegisters [device Id %v] error: %v", w.device.ID, err)
		return err
	}
	return nil
}

func (w *DeviceWorker) Stop() {
	w.cancel()
}
