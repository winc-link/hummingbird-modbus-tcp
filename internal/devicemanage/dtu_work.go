package devicemanage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/winc-link/hummingbird-modbus-tcp/constant"
	"github.com/winc-link/hummingbird-modbus-tcp/dtos"
	"github.com/winc-link/hummingbird-modbus-tcp/internal/connmanage"
	"github.com/winc-link/hummingbird-modbus-tcp/utils/cast"
	"github.com/winc-link/hummingbird-sdk-go/model"
	"strconv"
	"strings"
	"time"
)

const maxRetry = 2

type DtuDeviceWorker struct {
	device   Device
	status   *DeviceStatus
	ctx      context.Context
	cancel   context.CancelFunc
	resultCh chan<- DeviceResult
}

func NewDtuDeviceWorker(device Device, resultCh chan<- DeviceResult) *DtuDeviceWorker {
	ctx, cancel := context.WithCancel(context.Background())
	return &DtuDeviceWorker{
		device:   device,
		status:   &DeviceStatus{},
		ctx:      ctx,
		cancel:   cancel,
		resultCh: resultCh,
	}
}

func (w *DtuDeviceWorker) Start() {

	go func() {
		//w.connect()
		ticker := time.NewTicker(w.device.Period)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				err := connmanage.DtuClientManage.WithClient(w.device.IP+":"+strconv.Itoa(w.device.Port), 10*time.Second, func(client *connmanage.ModbusRtuClient) error {
					w.collectDataAndReport(client)
					return nil
				})
				if err != nil {
					SDKDriver.GetLogger().Errorf("WithClient error: %v", err)
				}
			case <-w.ctx.Done():
				//if w.handler != nil {
				//	w.handler.Close()
				//}
				SDKDriver.GetLogger().Infof("[%s] worker stopped", w.device.ID)
				return
			}
		}
	}()
}

func (w *DtuDeviceWorker) collectDataAndReport(conn *connmanage.ModbusRtuClient) {
	SDKDriver.GetLogger().Infof("[%s] start collect data and report", w.device.ID)
	product, ok := SDKDriver.GetProductById(w.device.ProductId)
	if !ok {
		return
	}
	SDKDriver.GetLogger().Infof("[%s] get product info", w.device.ID)

	reportData := make(map[string]interface{})
	for _, property := range product.Properties {
		time.Sleep(200 * time.Millisecond)
		var modbusSpec dtos.ModbusSpec

		err := json.Unmarshal([]byte(property.TypeSpec.Specs), &modbusSpec)
		if err != nil {
			SDKDriver.GetLogger().Errorf("unmarshal property %s error: %v", property.TypeSpec.Specs, err)
			continue
		}
		SDKDriver.GetLogger().Infof("modbus spec: %v", modbusSpec)
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
			for i := 0; i < maxRetry; i++ {
				value, err = conn.ReadCoils(byte(w.device.SlaveId), uint16(address), cast.GetQuantityByDataType(modbusSpec.DataType))

				if err == nil {
					break
				}

				// 判断是否是超时错误（根据你实际库的错误类型调整）
				if strings.Contains(err.Error(), "timeout") {
					if i < maxRetry-1 {
						time.Sleep(100 * time.Millisecond)
						continue
					}
				}
			}
			//value, err = conn.ReadCoils(byte(w.device.SlaveId), uint16(address), cast.GetQuantityByDataType(modbusSpec.DataType))
		case string(constant.InputStatus):
			for i := 0; i < maxRetry; i++ {
				value, err = conn.ReadDiscreteInputs(byte(w.device.SlaveId), uint16(address), cast.GetQuantityByDataType(modbusSpec.DataType))

				if err == nil {
					break
				}

				// 判断是否是超时错误（根据你实际库的错误类型调整）
				if strings.Contains(err.Error(), "timeout") {
					if i < maxRetry-1 {
						time.Sleep(100 * time.Millisecond)
						continue
					}
				}
			}

		case string(constant.HoldingRegister):
			for i := 0; i < maxRetry; i++ {
				value, err = conn.ReadHoldingRegisters(byte(w.device.SlaveId), uint16(address), cast.GetQuantityByDataType(modbusSpec.DataType))
				if err == nil {
					break
				}

				// 判断是否是超时错误（根据你实际库的错误类型调整）
				if strings.Contains(err.Error(), "timeout") {
					if i < maxRetry-1 {
						time.Sleep(100 * time.Millisecond)
						continue
					}
				}
			}

		case string(constant.InputRegisters):
			for i := 0; i < maxRetry; i++ {
				value, err = conn.ReadInputRegisters(byte(w.device.SlaveId), uint16(address), cast.GetQuantityByDataType(modbusSpec.DataType))
				if err == nil {
					break
				}

				// 判断是否是超时错误（根据你实际库的错误类型调整）
				if strings.Contains(err.Error(), "timeout") {
					if i < maxRetry-1 {
						time.Sleep(100 * time.Millisecond)
						continue
					}
				}
			}

		default:
			err = fmt.Errorf("unknown register")
		}

		if err != nil {
			SDKDriver.GetLogger().Errorf("deviceId %s read register %s error: %v", w.device.ID, modbusSpec.RegisterAddress, err)
			w.updateStatus(false, err.Error())
			return
		}
		v := cast.ConvertByteToInt(modbusSpec.DataType, modbusSpec.DataOrder, modbusSpec.Multiplier, value)

		reportData[property.Code] = v

	}
	if len(reportData) == 0 {
		return
	}

	r, _ := json.Marshal(reportData)
	SDKDriver.GetLogger().Infof("device ID %s report data: %v", w.device.ID, string(r))

	resp, err := SDKDriver.PropertyReport(w.device.ID, model.PropertyReport{
		CommonRequest: model.NewDefaultCommonRequest(),
		Data:          reportData,
	})

	SDKDriver.GetLogger().Infof("resp: %v", resp)
	if err != nil {
		SDKDriver.GetLogger().Errorf("property report %s error: %v", w.device.ID, err)
	}
}

func (w *DtuDeviceWorker) updateStatus(online bool, errMsg string) {
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

func (w *DtuDeviceWorker) Stop() {
	w.cancel()
}
