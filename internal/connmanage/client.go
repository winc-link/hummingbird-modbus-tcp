package connmanage

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type ModbusRtuClient struct {
	conn  net.Conn
	inUse bool
	mutex sync.Mutex
	//cond     *sync.Cond
	timeout  time.Duration
	lastUsed time.Time
}

func NewModbusRtuClient(conn net.Conn) *ModbusRtuClient {
	client := &ModbusRtuClient{
		conn:    conn,
		timeout: 5 * time.Second,
	}
	return client
}

func (m *ModbusRtuClient) Unlock() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.inUse = false
	//m.cond.Broadcast()
}

func (m *ModbusRtuClient) TryLock(timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		m.mutex.Lock()
		if !m.inUse {
			m.inUse = true
			m.mutex.Unlock()
			return true
		}
		m.mutex.Unlock()
		time.Sleep(10 * time.Millisecond)
	}
	return false
}

func (m *ModbusRtuClient) HealthCheck() bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.conn.SetReadDeadline(time.Now().Add(3000 * time.Millisecond))
	_, err := m.conn.Read(make([]byte, 1))

	if netErr, ok := err.(net.Error); ok {
		return netErr.Timeout() // 超时表示活跃
	}
	return err == nil
}

// 基本 Modbus 操作接口封装

func (m *ModbusRtuClient) ReadCoils(slaveID byte, address, quantity uint16) ([]byte, error) {
	return m.buildRequestAndParseResponse(slaveID, 0x01, address, quantity)
}

func (m *ModbusRtuClient) ReadDiscreteInputs(slaveID byte, address, quantity uint16) ([]byte, error) {
	return m.buildRequestAndParseResponse(slaveID, 0x02, address, quantity)
}

func (m *ModbusRtuClient) ReadHoldingRegisters(slaveID byte, address, quantity uint16) ([]byte, error) {
	return m.buildRequestAndParseResponse(slaveID, 0x03, address, quantity)
}

func (m *ModbusRtuClient) ReadInputRegisters(slaveID byte, address, quantity uint16) ([]byte, error) {
	return m.buildRequestAndParseResponse(slaveID, 0x04, address, quantity)
}

func (m *ModbusRtuClient) WriteSingleCoil(slaveID byte, address, value uint16) ([]byte, error) {
	if value == 1 {
		value = 0xFF00
	} else {
		value = 0x0000
	}

	data := make([]byte, 2)
	binary.BigEndian.PutUint16(data, value)
	request := buildModbusRequest(slaveID, 0x05, address, data)

	_, err := m.conn.Write(request)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 256)
	m.conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	n, err := m.conn.Read(buf)
	if err != nil {
		return nil, err
	}
	return parseModbusResponse(buf[:n])
}

func (m *ModbusRtuClient) WriteMultipleRegisters(slaveID byte, address, quantity uint16, values []byte) ([]byte, error) {
	request := buildModbusRequest(slaveID, 0x10, address, values)

	_, err := m.conn.Write(request)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 256)
	m.conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	n, err := m.conn.Read(buf)
	if err != nil {
		return nil, err
	}

	return parseModbusResponse(buf[:n])
}

func (m *ModbusRtuClient) buildRequestAndParseResponse(slaveID byte, functionCode byte, addr uint16, quantity uint16) ([]byte, error) {
	data := make([]byte, 2)
	binary.BigEndian.PutUint16(data, quantity)
	request := buildModbusRequest(slaveID, functionCode, addr, data)

	_ = m.conn.SetReadDeadline(time.Now().Add(50 * time.Millisecond))
	tmp := make([]byte, 256)
	m.conn.Read(tmp)
	_, err := m.conn.Write(request)
	if err != nil {
		return nil, err
	}
	// 等待设备响应（485 半双工）
	time.Sleep(100 * time.Millisecond)
	// 读 RTU 响应
	resp, err := readRTUFrame(m.conn, 5*time.Second)
	if err != nil {
		return nil, err
	}
	// 过滤回显
	resp = filterEcho(resp, request)
	return extractModbusData(resp)
}

func extractModbusData(resp []byte) ([]byte, error) {
	if len(resp) < 5 {
		return nil, fmt.Errorf("response too short: %x", resp)
	}

	slave := resp[0]
	function := resp[1]

	// 异常响应
	if function&0x80 != 0 {
		if len(resp) < 5 {
			return nil, fmt.Errorf("exception response too short: %x", resp)
		}
		return nil, fmt.Errorf(
			"modbus exception, slave=%d func=0x%x code=0x%x",
			slave, function, resp[2],
		)
	}

	switch function {
	case 0x01, 0x02, 0x03, 0x04:
		byteCount := int(resp[2])

		expectedLen := 3 + byteCount + 2
		if len(resp) < expectedLen {
			return nil, fmt.Errorf(
				"incomplete response, want %d, got %d",
				expectedLen, len(resp),
			)
		}

		return resp[3 : 3+byteCount], nil

	case 0x05, 0x06, 0x0F, 0x10:
		// 写类命令：无数据区，返回地址和数量
		return nil, nil

	default:
		return nil, fmt.Errorf("unsupported function code: 0x%x", function)
	}
}

func readRTUFrame(conn net.Conn, timeout time.Duration) ([]byte, error) {
	deadline := time.Now().Add(timeout)
	_ = conn.SetReadDeadline(deadline)

	buf := make([]byte, 0, 256)
	tmp := make([]byte, 128)

	for {
		n, err := conn.Read(tmp)
		if n > 0 {
			buf = append(buf, tmp[:n]...)

			if len(buf) >= 5 && isValidCRC(buf) {
				return buf, nil
			}
		}

		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				return nil, fmt.Errorf("modbus read timeout, buf=%x", buf)
			}
			return nil, err
		}
	}
}

func filterEcho(resp, req []byte) []byte {
	if len(resp) >= len(req) && string(resp[:len(req)]) == string(req) {
		return resp[len(req):]
	}
	return resp
}

func isValidCRC(frame []byte) bool {
	if len(frame) < 3 {
		return false
	}
	data := frame[:len(frame)-2]
	recv := binary.LittleEndian.Uint16(frame[len(frame)-2:])
	return crc16(data) == recv
}

// ------------------ 帧构造、解析与 CRC ----------------------

func buildModbusRequest(slaveID, functionCode byte, addr uint16, data []byte) []byte {
	buf := new(bytes.Buffer)
	buf.WriteByte(slaveID)
	buf.WriteByte(functionCode)

	switch functionCode {
	case 0x01, 0x02, 0x03, 0x04:
		binary.Write(buf, binary.BigEndian, addr)
		binary.Write(buf, binary.BigEndian, binary.BigEndian.Uint16(data))
	case 0x05, 0x06:
		binary.Write(buf, binary.BigEndian, addr)
		binary.Write(buf, binary.BigEndian, binary.BigEndian.Uint16(data))
	case 0x10:
		binary.Write(buf, binary.BigEndian, addr)
		quantity := len(data) / 2
		buf.WriteByte(byte(quantity >> 8))
		buf.WriteByte(byte(quantity & 0xFF))
		buf.WriteByte(byte(len(data)))
		buf.Write(data)
	default:
		log.Fatalf("不支持的功能码: 0x%02X", functionCode)
	}

	crc := crc16(buf.Bytes())
	binary.Write(buf, binary.LittleEndian, crc)
	return buf.Bytes()
}

func parseModbusResponse(response []byte) ([]byte, error) {
	length := len(response)
	if length < 5 {
		return nil, fmt.Errorf("响应长度过短")
	}

	crcData := response[:length-2]
	crcRecv := binary.LittleEndian.Uint16(response[length-2:])
	crcCalc := crc16(crcData)
	if crcRecv != crcCalc {
		return nil, fmt.Errorf("CRC 校验失败：收到 0x%04X，计算得 0x%04X", crcRecv, crcCalc)
	}

	functionCode := response[1]
	if functionCode&0x80 != 0 {
		return nil, fmt.Errorf("异常响应，异常码: 0x%02X", response[2])
	}

	byteCount := int(response[2])
	if byteCount%2 != 0 || length != 3+byteCount+2 {
		return nil, fmt.Errorf("响应长度异常")
	}

	data := response[3 : 3+byteCount]
	return data, nil
}

func crc16(data []byte) uint16 {
	crc := uint16(0xFFFF)
	for _, b := range data {
		crc ^= uint16(b)
		for i := 0; i < 8; i++ {
			if (crc & 0x0001) != 0 {
				crc = (crc >> 1) ^ 0xA001
			} else {
				crc >>= 1
			}
		}
	}
	return crc
}
