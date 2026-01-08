package main

import (
	"github.com/winc-link/hummingbird-modbus-tcp/config"
	"github.com/winc-link/hummingbird-modbus-tcp/internal/driver"
	"github.com/winc-link/hummingbird-sdk-go/service"
)

func main() {
	driverService := service.NewDriverService("official-modbus-tcp-v3.0")
	config.InitConfig(driverService)
	modbusTcpDriver := driver.NewModbusTcpProtocolDriver(driverService)
	if err := driverService.Start(modbusTcpDriver); err != nil {
		driverService.GetLogger().Error("driver service start error: %s", err)
		return
	}
}
