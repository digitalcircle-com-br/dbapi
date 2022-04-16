package types

type DBIn struct {
	T      string        `json:"t"`
	Q      string        `json:"q"`
	Params []interface{} `json:"params"`
}

type DBOut struct {
	Data interface{} `json:"data"`
	Err  string      `json:"err"`
}
