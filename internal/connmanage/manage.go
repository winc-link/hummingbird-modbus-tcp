package connmanage

import (
	"fmt"
	"net"
	"sync"
	"time"
)

var DtuClientManage *ConnectManager

type ConnectManager struct {
	lock sync.RWMutex
	conn map[string]*ModbusRtuClient
}

func init() {
	DtuClientManage = &ConnectManager{
		conn: make(map[string]*ModbusRtuClient),
	}
}

// GetDtuClientByImei 获取 Modbus 客户端，并尝试锁定
func (c *ConnectManager) GetDtuClientByAddress(address string, timeout time.Duration) (*ModbusRtuClient, error) {
	c.lock.RLock()
	client, ok := c.conn[address]
	c.lock.RUnlock()

	if !ok {
		return nil, fmt.Errorf("Address %s 不存在", address)
	}

	if !client.TryLock(timeout) {
		return nil, fmt.Errorf("获取客户端 %s 超时", address)
	}

	if !client.HealthCheck() {
		//c.lock.Lock()
		//client.conn.Close()
		//delete(c.conn, imei)
		//c.lock.Unlock()
		//return nil, fmt.Errorf("连接[%s]已失效", imei)
	}
	client.lastUsed = time.Now()
	return client, nil
}

// WithClient 自动处理锁/解锁的封装方式
func (c *ConnectManager) WithClient(address string, timeout time.Duration, fn func(*ModbusRtuClient) error) error {
	client, err := c.GetDtuClientByAddress(address, timeout)
	if err != nil {
		return err
	}

	defer client.Unlock()
	return fn(client)
}

// AddOrUpdateDtuConnectManager 添加或更新客户端连接
func (c *ConnectManager) AddOrUpdateDtuConnectManager(address string, conn net.Conn) {
	c.lock.Lock()
	defer c.lock.Unlock()

	//if client, ok := c.conn[address]; ok {
	//	client.conn.Close()
	//}

	if _, ok := c.conn[address]; ok {
		return
	}
	c.conn[address] = NewModbusRtuClient(conn)
}

// DeleteDtuConnectManager 删除连接
func (c *ConnectManager) DeleteDtuConnectManager(address string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if client, ok := c.conn[address]; ok {
		client.conn.Close()
		delete(c.conn, address)
	}
}

// CloseAll 关闭所有连接
func (c *ConnectManager) CloseAll() {
	c.lock.Lock()
	defer c.lock.Unlock()

	for imei, client := range c.conn {
		client.conn.Close()
		delete(c.conn, imei)
	}
}
