package main

import (
	"github.com/winc-link/hummingbird-modbus-tcp/config"
	"github.com/winc-link/hummingbird-modbus-tcp/internal/driver"
	"github.com/winc-link/hummingbird-sdk-go/service"
)

func main() {
	driverService := service.NewDriverService("official-modbus-tcp-v2.7")
	config.InitConfig(driverService)
	tcpDriver := driver.NewModbusTcpProtocolDriver(driverService)
	if err := driverService.Start(tcpDriver); err != nil {
		driverService.GetLogger().Error("driver service start error: %s", err)
		return
	}
}
