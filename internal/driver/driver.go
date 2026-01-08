/*******************************************************************************
 * Copyright 2017.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/

package driver

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/winc-link/hummingbird-modbus-tcp/constant"
	"github.com/winc-link/hummingbird-modbus-tcp/dtos"
	"github.com/winc-link/hummingbird-modbus-tcp/internal/devicemanage"
	"github.com/winc-link/hummingbird-modbus-tcp/utils/cast"
	"strconv"
	"time"

	"github.com/winc-link/hummingbird-sdk-go/commons"
	"github.com/winc-link/hummingbird-sdk-go/model"
	"github.com/winc-link/hummingbird-sdk-go/service"
)

type ModbusTcpProtocolDriver struct {
	sd      *service.DriverService
	manager *devicemanage.DeviceManager
}

func (dr ModbusTcpProtocolDriver) HandlePropertyReportDebug(ctx context.Context, deviceId string, data model.PropertyReport) error {
	return nil
}

func (dr ModbusTcpProtocolDriver) HandleEventReportDebug(ctx context.Context, deviceId string, data model.EventReport) error {
	return nil
}

// DeviceNotify 设备添加/修改/删除通知
func (dr ModbusTcpProtocolDriver) DeviceNotify(ctx context.Context, t commons.DeviceNotifyType, deviceId string, device model.Device) error {
	switch t {
	case commons.DeviceAddNotify, commons.DeviceUpdateNotify:
		device.Id = deviceId
		deviceM, err := tranDeviceToMangeDeviceModel(dr.sd, device)
		if err != nil {
			return err
		}
		dr.manager.AddOrUpdateDevice(deviceM)
	case commons.DeviceDeleteNotify:
		dr.manager.RemoveDevice(deviceId)
	}
	return nil
}

// ProductNotify 产品添加/修改/删除通知
func (dr ModbusTcpProtocolDriver) ProductNotify(ctx context.Context, t commons.ProductNotifyType, productId string, product model.Product) error {
	return nil
}

// Stop 驱动退出通知。
func (dr ModbusTcpProtocolDriver) Stop(ctx context.Context) error {
	dr.manager.Shutdown()
	return nil
}

// HandlePropertySet 设备属性设置
func (dr ModbusTcpProtocolDriver) HandlePropertySet(ctx context.Context, deviceId string, data model.PropertySet) error {
	dr.sd.GetLogger().Infof("HandlePropertySet deviceId:%s param:%v", deviceId, data)
	device, ok := dr.sd.GetDeviceById(deviceId)
	if ok != true {
		_ = dr.sd.PropertySetResponse(deviceId, model.PropertySetResponse{
			MsgId: data.MsgId,
			Data: model.PropertySetResponseData{
				Success:      false,
				Code:         uint32(constant.DeviceNotFound),
				ErrorMessage: string(constant.ErrorCodeMsgMap[constant.DeviceNotFound]),
			},
		})
		return nil
	}

	_, ok = dr.sd.GetProductById(device.ProductId)
	if ok != true {
		_ = dr.sd.PropertySetResponse(deviceId, model.PropertySetResponse{
			MsgId: data.MsgId,
			Data: model.PropertySetResponseData{
				Success:      false,
				Code:         uint32(constant.ProductNotFound),
				ErrorMessage: string(constant.ErrorCodeMsgMap[constant.ProductNotFound]),
			},
		})
		return nil
	}

	var (
		code  string
		value interface{}
	)
	for k, v := range data.Data {
		code = k
		value = v
	}
	dr.sd.GetLogger().Infof("code:%s value:%v", code, value)

	property := data.Spec[code]
	var modbusSpec dtos.ModbusSpec
	err := json.Unmarshal([]byte(property.TypeSpec.Specs), &modbusSpec)
	if err != nil {
		dr.sd.GetLogger().Errorf("unmarshal property %s error: %v", property.TypeSpec.Specs, err)
		_ = dr.sd.PropertySetResponse(deviceId, model.PropertySetResponse{
			MsgId: data.MsgId,
			Data: model.PropertySetResponseData{
				Success: false,
			},
		})
		return nil
	}
	address, err := strconv.Atoi(modbusSpec.RegisterAddress)
	if err != nil {
		dr.sd.GetLogger().Errorf("convert modbus register address to int error: %v", err)
		_ = dr.sd.PropertySetResponse(deviceId, model.PropertySetResponse{
			MsgId: data.MsgId,
			Data: model.PropertySetResponseData{
				Success: false,
			},
		})
		return nil
	}

	switch modbusSpec.RegisterType {
	case string(constant.CoilStatus):
		v, err := cast.ToBool(value, cast.CONVERT_ALL)
		if err != nil {
			dr.sd.GetLogger().Errorf("modbus register value error: %v", value)
			_ = dr.sd.PropertySetResponse(deviceId, model.PropertySetResponse{
				MsgId: data.MsgId,
				Data: model.PropertySetResponseData{
					Success: false,
				},
			})
			return nil
		}
		err = dr.manager.WriteSingleCoil(deviceId, uint16(address), cast.ConvertBoolToUint(v))
		if err != nil {
			_ = dr.sd.PropertySetResponse(deviceId, model.PropertySetResponse{
				MsgId: data.MsgId,
				Data: model.PropertySetResponseData{
					Success: false,
				},
			})
			return nil
		}

		_ = dr.sd.PropertySetResponse(deviceId, model.PropertySetResponse{
			MsgId: data.MsgId,
			Data: model.PropertySetResponseData{
				Success: true,
			},
		})
		return nil
	case string(constant.HoldingRegister):
		quantity, v, err := cast.ConvertToRegisters(modbusSpec.DataType, modbusSpec.DataOrder, value)
		if err != nil {
			_ = dr.sd.PropertySetResponse(deviceId, model.PropertySetResponse{
				MsgId: data.MsgId,
				Data: model.PropertySetResponseData{
					Success: false,
				},
			})
			return nil
		}

		err = dr.manager.WriteMultipleRegisters(deviceId, uint16(address), quantity, v)
		if err != nil {
			_ = dr.sd.PropertySetResponse(deviceId, model.PropertySetResponse{
				MsgId: data.MsgId,
				Data: model.PropertySetResponseData{
					Success: false,
				},
			})
			return nil
		}

		_ = dr.sd.PropertySetResponse(deviceId, model.PropertySetResponse{
			MsgId: data.MsgId,
			Data: model.PropertySetResponseData{
				Success: true,
			},
		})
		return nil
	}

	return nil
}

// HandlePropertyGet 设备属性查询
func (dr ModbusTcpProtocolDriver) HandlePropertyGet(ctx context.Context, deviceId string, data model.PropertyGet) error {
	_, ok := dr.sd.GetDeviceById(deviceId)
	if ok != true {
		_ = dr.sd.PropertyGetResponse(deviceId, model.PropertyGetResponse{
			MsgId: data.MsgId,
			Data:  []model.PropertyGetResponseData{},
		})
		return nil
	}

	return nil
}

// HandleServiceExecute 设备服务调用
func (dr ModbusTcpProtocolDriver) HandleServiceExecute(ctx context.Context, deviceId string, data model.ServiceExecuteRequest) error {
	_, ok := dr.sd.GetDeviceById(deviceId)
	if ok != true {
		_ = dr.sd.ServiceExecuteResponse(deviceId, model.ServiceExecuteResponse{
			MsgId: data.MsgId,
			Data:  model.ServiceDataOut{},
		})
		return nil
	}

	return nil
}

func tranDeviceToMangeDeviceModel(sd *service.DriverService, dev model.Device) (devicemanage.Device, error) {
	if dev.Ip == "" {
		sd.GetLogger().Errorf("failed to convert deviceName [%s] device ip is empty %s", dev.Name)
		return devicemanage.Device{}, errors.New("device ip is empty")
	}
	if !cast.IsValidIP(dev.Ip) {
		sd.GetLogger().Errorf("failed to convert deviceName [%s] device ip is invalid", dev.Name)
		return devicemanage.Device{}, errors.New("device ip is invalid")
	}

	port, err := strconv.Atoi(dev.Port)
	if err != nil {
		sd.GetLogger().Errorf("failed to convert deviceName [%s] port to int: %s", dev.Name, err.Error())
		return devicemanage.Device{}, errors.New("failed to convert deviceName to int")
	}
	period, err := strconv.Atoi(dev.Period)
	if err != nil {
		sd.GetLogger().Errorf("failed to convert deviceName [%s] period to int: %s", dev.Name, err.Error())
	}
	slaveId, err := strconv.Atoi(dev.SlaveId)
	if err != nil {
		sd.GetLogger().Errorf("failed to convert deviceName [%s] slaveId to int: %s", dev.Name, err.Error())
	}

	var periodTime time.Duration
	if period <= 0 {
		sd.GetLogger().Errorf("failed to convert deviceName [%s] period %d", dev.Name, period)
		periodTime = 24 * time.Hour * 365
	} else {
		periodTime = time.Duration(period) * time.Second
	}

	return devicemanage.Device{
		ID:        dev.Id,
		IP:        dev.Ip,
		Port:      port,
		ProductId: dev.ProductId,
		Period:    periodTime,
		SlaveId:   slaveId,
	}, nil
}

// NewModbusTcpProtocolDriver 协议驱动
func NewModbusTcpProtocolDriver(sd *service.DriverService) *ModbusTcpProtocolDriver {
	manager := devicemanage.NewDeviceManager()
	manager.Start()

	for _, device := range sd.GetDeviceList() {
		deviceM, err := tranDeviceToMangeDeviceModel(sd, device)
		if err != nil {
			continue
		}
		manager.AddOrUpdateDevice(deviceM)
	}
	devicemanage.SDKDriver = sd
	return &ModbusTcpProtocolDriver{
		sd:      sd,
		manager: manager,
	}
}
