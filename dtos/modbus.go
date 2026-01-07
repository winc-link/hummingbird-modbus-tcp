package dtos

type ModbusSpec struct {
	DataOrder       string `json:"data_order"`
	DataType        string `json:"data_type"`
	Multiplier      string `json:"multiplier"`
	RegisterAddress string `json:"register_address"`
	RegisterType    string `json:"register_type"`
	SlaveId         string `json:"slave_id"`
	Unit            string `json:"unit"`
	UnitName        string `json:"unitName"`
}
