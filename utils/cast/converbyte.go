package cast

import (
	"bytes"
	"encoding/binary"
	"github.com/winc-link/hummingbird-modbus-tcp/constant"
)

//https://www.cnblogs.com/JohnnyLei/p/18180046

func BytesToBool(b []byte) bool {
	bytesBuffer := bytes.NewBuffer(b)
	var x bool
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return x
}

func BytesToInt16BigEndian(b []byte) int16 {
	bytesBuffer := bytes.NewBuffer(b)
	var x int16
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return x
}

func BytesToInt16LittleEndian(b []byte) int16 {
	bytesBuffer := bytes.NewBuffer(b)
	var x int16
	binary.Read(bytesBuffer, binary.LittleEndian, &x)
	return x
}

func BytesToUint16BigEndian(b []byte) uint16 {
	bytesBuffer := bytes.NewBuffer(b)
	var x uint16
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return x
}

func BytesToUint16LittleEndian(b []byte) uint16 {
	bytesBuffer := bytes.NewBuffer(b)
	var x uint16
	binary.Read(bytesBuffer, binary.LittleEndian, &x)
	return x
}

func BytesToInt32BigEndian(b []byte) int32 {
	bytesBuffer := bytes.NewBuffer(b)
	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return x
}

func BytesToInt32LittleEndian(b []byte) int32 {
	bytesBuffer := bytes.NewBuffer(b)
	var x int32
	binary.Read(bytesBuffer, binary.LittleEndian, &x)
	return x
}

// badc
func BytesToInt32BigEndianSwap(b []byte) int32 {
	if len(b) != 4 {
		return 0
	}
	newb := make([]byte, 4)
	newb[0] = b[1]
	newb[1] = b[0]
	newb[2] = b[3]
	newb[3] = b[2]
	bytesBuffer := bytes.NewBuffer(newb)
	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return x
}

// CDAB
func BytesToInt32BLittleEndianSwap(b []byte) int32 {
	if len(b) != 4 {
		return 0
	}
	newb := make([]byte, 4)
	newb[0] = b[2]
	newb[1] = b[3]
	newb[2] = b[0]
	newb[3] = b[1]
	bytesBuffer := bytes.NewBuffer(newb)
	var x int32
	binary.Read(bytesBuffer, binary.LittleEndian, &x)
	return x
}

func BytesToFloat32BigEndian(b []byte) float32 {
	bytesBuffer := bytes.NewBuffer(b)
	var x float32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return x
}

func BytesToFloat32LittleEndian(b []byte) float32 {
	bytesBuffer := bytes.NewBuffer(b)
	var x float32
	binary.Read(bytesBuffer, binary.LittleEndian, &x)
	return x
}

// BADC
func BytesToFloat32BigEndianSwap(b []byte) float32 {
	if len(b) != 4 {
		return 0
	}
	newb := make([]byte, 4)
	newb[0] = b[1]
	newb[1] = b[0]
	newb[2] = b[3]
	newb[3] = b[2]
	bytesBuffer := bytes.NewBuffer(newb)
	var x float32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return x
}

// CDAB
func BytesToFloat32BLittleEndianSwap(b []byte) float32 {
	if len(b) != 4 {
		return 0
	}
	newb := make([]byte, 4)
	newb[0] = b[2]
	newb[1] = b[3]
	newb[2] = b[0]
	newb[3] = b[1]
	bytesBuffer := bytes.NewBuffer(newb)
	var x float32
	binary.Read(bytesBuffer, binary.LittleEndian, &x)
	return x
}

func BytesToUint32BigEndian(b []byte) uint32 {
	bytesBuffer := bytes.NewBuffer(b)
	var x uint32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return x
}

func BytesToUint32LittleEndian(b []byte) uint32 {
	bytesBuffer := bytes.NewBuffer(b)
	var x uint32
	binary.Read(bytesBuffer, binary.LittleEndian, &x)
	return x
}

func BytesToUint32BigEndianSwap(b []byte) uint32 {
	if len(b) != 4 {
		return 0
	}
	newb := make([]byte, 4)
	newb[0] = b[1]
	newb[1] = b[0]
	newb[2] = b[3]
	newb[3] = b[2]
	bytesBuffer := bytes.NewBuffer(newb)
	var x uint32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return x
}

func BytesToUint32BLittleEndianSwap(b []byte) uint32 {
	if len(b) != 4 {
		return 0
	}
	newb := make([]byte, 4)
	newb[0] = b[2]
	newb[1] = b[3]
	newb[2] = b[0]
	newb[3] = b[1]
	bytesBuffer := bytes.NewBuffer(newb)
	var x uint32
	binary.Read(bytesBuffer, binary.LittleEndian, &x)
	return x
}

func BytesToInt64BigEndian(b []byte) int64 {
	bytesBuffer := bytes.NewBuffer(b)
	var x int64
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return x
}

func BytesToInt64LittleEndian(b []byte) int64 {
	bytesBuffer := bytes.NewBuffer(b)
	var x int64
	binary.Read(bytesBuffer, binary.LittleEndian, &x)
	return x
}

// BADCFEHG
func BytesToInt64BigEndianSwap(b []byte) int64 {
	if len(b) != 8 {
		return 0
	}
	newb := make([]byte, 8)
	newb[0] = b[1]
	newb[1] = b[0]
	newb[2] = b[3]
	newb[3] = b[2]
	newb[5] = b[5]
	newb[6] = b[4]
	newb[7] = b[7]
	newb[8] = b[6]
	bytesBuffer := bytes.NewBuffer(newb)
	var x int64
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return x
}

// GHEFCDAB
func BytesToInt64LittleEndianSwap(b []byte) int64 {
	if len(b) != 8 {
		return 0
	}
	newb := make([]byte, 8)
	newb[0] = b[6]
	newb[1] = b[7]
	newb[2] = b[4]
	newb[3] = b[5]
	newb[5] = b[2]
	newb[6] = b[3]
	newb[7] = b[0]
	newb[8] = b[1]
	bytesBuffer := bytes.NewBuffer(newb)
	var x int64
	binary.Read(bytesBuffer, binary.LittleEndian, &x)
	return x
}

func BytesToUInt64BigEndian(b []byte) uint64 {
	bytesBuffer := bytes.NewBuffer(b)
	var x uint64
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return x
}

func BytesToUInt64LittleEndian(b []byte) uint64 {
	bytesBuffer := bytes.NewBuffer(b)
	var x uint64
	binary.Read(bytesBuffer, binary.LittleEndian, &x)
	return x
}

// BADCFEHG
func BytesToUInt64BigEndianSwap(b []byte) uint64 {
	if len(b) != 8 {
		return 0
	}
	newb := make([]byte, 8)
	newb[0] = b[1]
	newb[1] = b[0]
	newb[2] = b[3]
	newb[3] = b[2]
	newb[5] = b[5]
	newb[6] = b[4]
	newb[7] = b[7]
	newb[8] = b[6]
	bytesBuffer := bytes.NewBuffer(newb)
	var x uint64
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return x
}

// GHEFCDAB
func BytesToUInt64LittleEndianSwap(b []byte) uint64 {
	if len(b) != 8 {
		return 0
	}
	newb := make([]byte, 8)
	newb[0] = b[6]
	newb[1] = b[7]
	newb[2] = b[4]
	newb[3] = b[5]
	newb[5] = b[2]
	newb[6] = b[3]
	newb[7] = b[0]
	newb[8] = b[1]
	bytesBuffer := bytes.NewBuffer(newb)
	var x uint64
	binary.Read(bytesBuffer, binary.LittleEndian, &x)
	return x
}

func ConvertMultiplierToFloat(multiplier string) (float64, error) {
	if multiplier == "1" || multiplier == "" {
		return 1, nil
	}
	return ToFloat64(multiplier, CONVERT_ALL)
}

func ConvertByteToInt(dataType, dataOrder, multiplier string, b []byte) interface{} {
	switch dataType {
	case constant.Bool:
		return BytesToBool(b)
	case constant.Hex:
	case constant.Binary:
	case constant.Int16:
		var v int16
		if dataOrder == (constant.AB) {
			v = BytesToInt16BigEndian(b)
		} else if dataOrder == (constant.BA) {
			v = BytesToInt16LittleEndian(b)
		}
		mm, err := ConvertMultiplierToFloat(multiplier)
		if err != nil {
			return v
		}
		return int16(mm * float64(v))
	case constant.UInt16:
		var v uint16
		if dataOrder == (constant.AB) {
			v = BytesToUint16BigEndian(b)
		} else if dataOrder == (constant.BA) {
			v = BytesToUint16LittleEndian(b)
		}
		mm, err := ConvertMultiplierToFloat(multiplier)
		if err != nil {
			return v
		}
		return int16(mm * float64(v))
	case constant.Int32:
		var v int32
		if dataOrder == string(constant.ABCD) {
			v = BytesToInt32BigEndian(b)
		} else if dataOrder == string(constant.DCBA) {
			v = BytesToInt32LittleEndian(b)
		} else if dataOrder == string(constant.BADC) {
			v = BytesToInt32BigEndianSwap(b)
		} else if dataOrder == string(constant.CDAB) {
			v = BytesToInt32BLittleEndianSwap(b)
		}
		mm, err := ConvertMultiplierToFloat(multiplier)
		if err != nil {
			return v
		}
		return int32(mm * float64(v))
	case constant.UInt32:
		var v uint32
		if dataOrder == string(constant.ABCD) {
			v = BytesToUint32BigEndian(b)
		} else if dataOrder == string(constant.DCBA) {
			v = BytesToUint32LittleEndian(b)
		} else if dataOrder == string(constant.BADC) {
			v = BytesToUint32BigEndianSwap(b)
		} else if dataOrder == string(constant.CDAB) {
			v = BytesToUint32BLittleEndianSwap(b)
		}
		mm, err := ConvertMultiplierToFloat(multiplier)
		if err != nil {
			return v
		}
		return uint32(mm * float64(v))
	case constant.Int64:
		var v int64
		if dataOrder == string(constant.ABCDEFGH) {
			v = BytesToInt64BigEndian(b)
		} else if dataOrder == string(constant.GHEFCDAB) {
			v = BytesToInt64LittleEndianSwap(b)
		} else if dataOrder == string(constant.BADCFEHG) {
			v = BytesToInt64BigEndianSwap(b)
		} else if dataOrder == string(constant.HGFEDCBA) {
			v = BytesToInt64LittleEndian(b)
		}
		mm, err := ConvertMultiplierToFloat(multiplier)
		if err != nil {
			return v
		}
		return int64(mm * float64(v))
	case constant.UInt64:
		var v uint64
		if dataOrder == string(constant.ABCDEFGH) {
			v = BytesToUInt64BigEndian(b)
		} else if dataOrder == string(constant.GHEFCDAB) {
			v = BytesToUInt64LittleEndianSwap(b)
		} else if dataOrder == string(constant.BADCFEHG) {
			v = BytesToUInt64BigEndianSwap(b)
		} else if dataOrder == constant.HGFEDCBA {
			v = BytesToUInt64LittleEndian(b)
		}
		mm, err := ConvertMultiplierToFloat(multiplier)
		if err != nil {
			return v
		}
		return uint64(mm * float64(v))
	case constant.Float32:
		var v float32
		if dataOrder == string(constant.ABCD) {
			v = BytesToFloat32BigEndian(b)
		} else if dataOrder == string(constant.DCBA) {
			v = BytesToFloat32LittleEndian(b)
		} else if dataOrder == string(constant.BADC) {
			v = BytesToFloat32BigEndianSwap(b)
		} else if dataOrder == string(constant.CDAB) {
			v = BytesToFloat32BLittleEndianSwap(b)
		}
		mm, err := ConvertMultiplierToFloat(multiplier)
		if err != nil {
			return v
		}
		//v = float32(mm) * v

		v = RoundFloat32(float32(mm)*v, 2)
		return v
	case constant.Float64:

	}
	return nil
}

func GetQuantityByDataType(dataType string) uint16 {
	switch dataType {
	case string(constant.Bool):
		return 1
	case string(constant.Hex):
		return 1
	case string(constant.Binary):
		return 1
	case string(constant.Int16):
		return 1
	case string(constant.UInt16):
		return 1
	case string(constant.Int32):
		return 2
	case string(constant.UInt32):
		return 2
	case string(constant.Int64):
		return 4
	case string(constant.UInt64):
		return 4
	case string(constant.Float32):
		return 2
	case string(constant.Float64):
		return 4
	}
	return 0
}

func ConvertUintToBool(v uint16) bool {
	return v == 1
}

func ConvertBoolToUint(v bool) uint16 {
	if v == true {
		return 1
	}
	return 0
}
