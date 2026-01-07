package cast

import (
	"encoding/binary"
	"errors"
	"github.com/winc-link/hummingbird-modbus-tcp/constant"
	"math"
)

// Int16ToBytesBigEndian 大端
func Int16ToBytesBigEndian(n int16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(n))
	return b
}

// Int16ToBytesLittleEndian 小端
func Int16ToBytesLittleEndian(n int16) []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, uint16(n))
	return b
}

// Uint16ToBytesBigEndian 大端
func Uint16ToBytesBigEndian(n uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, n)
	return b
}

// Uint16ToBytesLittleEndian 小端
func Uint16ToBytesLittleEndian(n uint16) []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, n)
	return b
}

// Int32ToBytesBigEndian 大端
func Int32ToBytesBigEndian(n int32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(n))
	return b
}

// Int32ToBytesLittleEndian 小端
func Int32ToBytesLittleEndian(n int32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(n))
	return b
}

// Int32ToBytesBigEndianBigEndianSwap BADC
func Int32ToBytesBigEndianBigEndianSwap(n int32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(n))
	newb := make([]byte, 4)
	newb[0] = b[1]
	newb[1] = b[0]
	newb[2] = b[3]
	newb[3] = b[2]
	return newb
}

// Int32ToBytesBigEndianLittleEndianSwap CDAB
func Int32ToBytesBigEndianLittleEndianSwap(n int32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(n))
	newb := make([]byte, 4)
	newb[0] = b[2]
	newb[1] = b[3]
	newb[2] = b[0]
	newb[3] = b[1]
	return newb
}

// Int32ToBytesBigEndianLittleEndianSwap
func Uint32ToBytesBigEndian(n uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, n)
	return b
}

// Uint32ToBytesLittleEndian
func Uint32ToBytesLittleEndian(n uint32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, n)
	return b
}

// Uint32ToBytesBigEndianBigEndianSwap
func Uint32ToBytesBigEndianBigEndianSwap(n uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(n))
	newb := make([]byte, 4)
	newb[0] = b[1]
	newb[1] = b[0]
	newb[2] = b[3]
	newb[3] = b[2]
	return newb
}

// Uint32ToBytesBigEndianLittleEndianSwap
func Uint32ToBytesBigEndianLittleEndianSwap(n uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(n))
	newb := make([]byte, 4)
	newb[0] = b[2]
	newb[1] = b[3]
	newb[2] = b[0]
	newb[3] = b[1]
	return newb
}

func Int64ToBytesBigEndian(n int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(n))
	return b
}

func Int64ToBytesLittleEndian(n int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(n))
	return b
}

// BADCFEHG
func Int64ToBytesBigEndianSwap(n int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(n))
	newb := make([]byte, 8)
	newb[0] = b[1]
	newb[1] = b[0]
	newb[2] = b[3]
	newb[3] = b[2]
	newb[5] = b[5]
	newb[6] = b[4]
	newb[7] = b[7]
	newb[8] = b[6]
	return newb
}

// GHEFCDAB
func Int64ToBytesLittleEndianSwap(n int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(n))
	newb := make([]byte, 8)
	newb[0] = b[6]
	newb[1] = b[7]
	newb[2] = b[4]
	newb[3] = b[5]
	newb[5] = b[2]
	newb[6] = b[3]
	newb[7] = b[0]
	newb[8] = b[1]
	return newb
}

func Uint64ToBytesBigEndian(n uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, n)
	return b
}

func Uint64ToBytesLittleEndian(n uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, n)
	return b
}

// BADCFEHG
func Uint64ToBytesBigEndianSwap(n uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(n))
	newb := make([]byte, 8)
	newb[0] = b[1]
	newb[1] = b[0]
	newb[2] = b[3]
	newb[3] = b[2]
	newb[5] = b[5]
	newb[6] = b[4]
	newb[7] = b[7]
	newb[8] = b[6]
	return newb
}

// GHEFCDAB
func Uint64ToBytesLittleEndianSwap(n uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(n))
	newb := make([]byte, 8)
	newb[0] = b[6]
	newb[1] = b[7]
	newb[2] = b[4]
	newb[3] = b[5]
	newb[5] = b[2]
	newb[6] = b[3]
	newb[7] = b[0]
	newb[8] = b[1]
	return newb
}

func Float32ToBytesBigEndian(f float32) []byte {
	bits := math.Float32bits(f)
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, bits)
	return b
}

func Float32ToBytesLittleEndian(f float32) []byte {
	bits := math.Float32bits(f)
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, bits)
	return b
}

func Float32ToBytesBigEndianSwap(f float32) []byte {
	bits := math.Float32bits(f)
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, bits)

	newb := make([]byte, 4)
	newb[0] = b[1]
	newb[1] = b[0]
	newb[2] = b[3]
	newb[3] = b[2]
	return newb
}

func Float32ToBytesLittleEndianSwap(f float32) []byte {
	bits := math.Float32bits(f)
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, bits)
	newb := make([]byte, 4)
	newb[0] = b[2]
	newb[1] = b[3]
	newb[2] = b[0]
	newb[3] = b[1]
	return newb
}

func ConvertToRegisters(dataType, dataOrder string, v interface{}) (uint16, []byte, error) {
	var (
		quantity uint16
		value    []byte
		err      error
	)

	switch dataType {
	case constant.Int16:
		quantity = 1
		int16v, err := ToInt16(v, CONVERT_ALL)
		if err != nil {
			return 0, nil, err
		}
		if dataOrder == constant.AB {
			value = Int16ToBytesBigEndian(int16v)
		} else if dataOrder == constant.BA {
			value = Int16ToBytesLittleEndian(int16v)
		} else {
			err = errors.New("invalid data order")
		}
	case constant.UInt16:
		quantity = 1

		uint16v, err := ToUint16(v, CONVERT_ALL)
		if err != nil {
			return 0, nil, err
		}
		if dataOrder == constant.AB {
			value = Uint16ToBytesBigEndian(uint16v)
		} else if dataOrder == constant.BA {
			value = Uint16ToBytesLittleEndian(uint16v)
		} else {
			err = errors.New("invalid data order")
		}

	case constant.Int32:
		quantity = 2
		int32v, err := ToInt32(v, CONVERT_ALL)
		if err != nil {
			return 0, nil, err
		}
		if dataOrder == constant.ABCD {
			value = Int32ToBytesBigEndian(int32v)
		} else if dataOrder == constant.DCBA {
			value = Int32ToBytesLittleEndian(int32v)
		} else if dataOrder == constant.BADC {
			value = Int32ToBytesBigEndianBigEndianSwap(int32v)
		} else if dataOrder == constant.CDAB {
			value = Int32ToBytesBigEndianLittleEndianSwap(int32v)
		} else {
			err = errors.New("invalid data order")
		}

	case constant.UInt32:
		quantity = 2
		uint32v, err := ToUint32(v, CONVERT_ALL)
		if err != nil {
			return 0, nil, err
		}
		if dataOrder == constant.ABCD {
			value = Uint32ToBytesBigEndian(uint32v)
		} else if dataOrder == constant.DCBA {
			value = Uint32ToBytesLittleEndian(uint32v)
		} else if dataOrder == constant.BADC {
			value = Uint32ToBytesBigEndianBigEndianSwap(uint32v)
		} else if dataOrder == constant.CDAB {
			value = Uint32ToBytesBigEndianLittleEndianSwap(uint32v)
		} else {
			err = errors.New("invalid data order")
		}
	case constant.Int64:
		quantity = 4
		int64v, err := ToInt64(v, CONVERT_ALL)
		if err != nil {
			return 0, nil, err
		}
		if dataOrder == constant.ABCDEFGH {
			value = Int64ToBytesBigEndian(int64v)
		} else if dataOrder == constant.HGFEDCBA {
			value = Int64ToBytesLittleEndian(int64v)
		} else if dataOrder == constant.BADCFEHG {
			value = Int64ToBytesBigEndianSwap(int64v)
		} else if dataOrder == constant.GHEFCDAB {
			value = Int64ToBytesLittleEndianSwap(int64v)
		} else {
			err = errors.New("invalid data order")
		}

	case constant.UInt64:
		quantity = 4

		uint64v, err := ToUint64(v, CONVERT_ALL)
		if err != nil {
			return 0, nil, err
		}
		if dataOrder == constant.ABCDEFGH {
			value = Uint64ToBytesBigEndian(uint64v)
		} else if dataOrder == constant.HGFEDCBA {
			value = Uint64ToBytesLittleEndian(uint64v)
		} else if dataOrder == constant.BADCFEHG {
			value = Uint64ToBytesBigEndianSwap(uint64v)
		} else if dataOrder == constant.GHEFCDAB {
			value = Uint64ToBytesLittleEndianSwap(uint64v)
		} else {
			err = errors.New("invalid data order")
		}

	case constant.Float32:
		quantity = 2

		float32v, err := ToFloat32(v, CONVERT_ALL)
		if err != nil {
			return 0, nil, err
		}
		if dataOrder == constant.ABCD {
			value = Float32ToBytesBigEndian(float32v)
		} else if dataOrder == constant.DCBA {
			value = Float32ToBytesLittleEndian(float32v)
		} else if dataOrder == constant.BADC {
			value = Float32ToBytesBigEndianSwap(float32v)
		} else if dataOrder == constant.CDAB {
			value = Float32ToBytesLittleEndianSwap(float32v)
		} else {
			err = errors.New("invalid data order")
		}

	default:
		err = errors.New("unsupported data type")
	}

	return quantity, value, err
}
