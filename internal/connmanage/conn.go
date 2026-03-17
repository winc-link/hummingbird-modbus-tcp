package connmanage

import (
	"fmt"
	"net"
	"time"
)

// 创建Modbus TCP连接
func createModbusTCPConnection(ip string, port int) (net.Conn, error) {
	// 构建地址
	address := fmt.Sprintf("%s:%d", ip, port)

	// 设置连接超时
	dialer := &net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	// 建立TCP连接
	conn, err := dialer.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("连接Modbus设备失败: %w", err)
	}

	// 设置读写超时
	conn.SetDeadline(time.Now().Add(10 * time.Second))
	return conn, nil
}

// 更完善的连接管理器
type ModbusTCPClient struct {
	conn    net.Conn
	timeout time.Duration
}

// 创建ModbusTCPClient
func NewModbusTCPClient(ip string, port int) (net.Conn, error) {
	conn, err := createModbusTCPConnection(ip, port)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
