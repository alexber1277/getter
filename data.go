package getter

const (
	dumpFileList = "net_{EXCH}_{SYM}.dtx"
)

type Kline struct {
	Index    int     `json:"i"`
	Open     float64 `json:"o"`
	Close    float64 `json:"c"`
	High     float64 `json:"h"`
	Low      float64 `json:"l"`
	Volume   float64 `json:"v"`
	Trades   int64   `json:"t"`
	OpenTime int64   `json:"ot"`
}

var Klines []*Kline
