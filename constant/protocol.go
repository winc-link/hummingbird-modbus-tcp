package constant

const (
	ModbusTcpZh        = "Modbus TCP"
	ModbusTcpEh        = "Modbus TCP"
	ModbusRtuZh        = "Modbus RTU"
	ModbusRtuEh        = "Modbus RTU"
	ModbusRtuOverTcpZh = "Modbus RTU Over TCP"
	ModbusRtuOverTcpEh = "Modbus RTU Over TCP"
	SiemensS7Zh        = "西门子 S7"
	SiemensS7Eh        = "Siemens S7"
	SiemensPPIZh       = "西门子 PPI"
	SiemensPPIEh       = "Siemens PPI"
	CJT188Zh           = "CJ-T188"
	CJT188Eh           = "CJ-T188"
)

type Protocol string

const (
	ModbusTcpClient        Protocol = "ModbusTcp客户端"
	ModbusRtuClient        Protocol = "ModbusRtu客户端"
	ModbusRtuOverTcpClient Protocol = "ModbusRtuOverTcp客户端"
	DLT645Client           Protocol = "DLT645客户端"
	DLT698Client           Protocol = "DLT698客户端"
	SiemensClient          Protocol = "西门子PLC协议"
	BACnetIPClient         Protocol = "BACnetIP客户端"
	IEC104Client           Protocol = "IEC104客户端"
	OPCUAClient            Protocol = "OPC UA客户端"
	Serial                 Protocol = "通用串口"
)

// "无校验", "奇校验", "偶校验"
// Parity: N - None, E - Even, O - Odd (default E)
const (
	ModbusParityN = "NONE"
	ModbusParityE = "EVEN"
	ModbusParityO = "ODD"
)

type RegisterType string

const (
	CoilStatus      RegisterType = "01 读线圈状态"
	InputStatus     RegisterType = "02 读离散输入状态"
	HoldingRegister RegisterType = "03 读保持寄存器"
	InputRegisters  RegisterType = "04 读输入寄存器"
)

const (
	NULL     string = "无"
	AB       string = "AB"
	BA       string = "BA"
	ABCD     string = "ABCD"
	DCBA     string = "DCBA"
	BADC     string = "BADC"
	CDAB     string = "CDAB"
	ABCDEFGH string = "ABCDEFGH"
	GHEFCDAB string = "GHEFCDAB"
	BADCFEHG string = "BADCFEHG"
	HGFEDCBA string = "HGFEDCBA"
)

const (
	//1字节
	Bool string = "Bool(1 Bytes)"
	//2字节
	Hex    string = "Hex(2 Bytes)"
	Binary string = "Binary(2 Bytes)"
	//2字节
	Int16  string = "Int16(2 Bytes)"
	UInt16 string = "UInt16(2 Bytes)"
	//4字节
	Int32  string = "Int32(4 Bytes)"
	UInt32 string = "UInt32(4 Bytes)"
	//8字节
	Int64  string = "Int64(8 Bytes)"
	UInt64 string = "UInt64(8 Bytes)"
	//4字节
	Float32 string = "Float32(4 Bytes)"
	//8字节
	Float64 string = "Float64(8 Bytes)"
)
