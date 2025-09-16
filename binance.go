package getter

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	urlKlines     = "https://api1.binance.com/api/v3/klines?symbol={SYM}&interval=1h&limit=1000&endTime={TM}"
	urlKlinesLast = "https://api1.binance.com/api/v3/klines?symbol={SYM}&interval=1h&limit=500"
)

var localList []*Kline

func GetBinanceDataLast(symbol string) error {
	list, err := reqBinanceLast(symbol)
	if err != nil {
		return err
	}
	Klines = formedList(list)
	return nil
}

func GetBinanceData(symbol string) {
	if loadDump("binance", symbol) {
		return
	}
	tm := time.Now().Unix() * 1000
	for i := 0; true; i++ {
		list := reqBinance(symbol, tm)
		tm = list[0].OpenTime
		localList = append(localList, list...)
		log.Println("iter: ", i, len(list), list[0].OpenTime)
		if len(list) < 1000 {
			break
		}
	}
	localList = formedList(localList)
	saveDump("binance", symbol, localList)
	Klines = localList
}

func formedList(list []*Kline) []*Kline {
	var newList []*Kline
	mpData := make(map[int64]struct{})
	for _, l := range list {
		if l.Close == l.Open {
			continue
		}
		if _, ok := mpData[l.OpenTime]; !ok {
			mpData[l.OpenTime] = struct{}{}
			newList = append(newList, l)
		}
	}
	sort.Slice(newList, func(i, j int) bool {
		return newList[i].OpenTime < newList[j].OpenTime
	})
	for i, l := range newList {
		l.Index = i
	}
	return newList
}

func reqBinance(sym string, offset int64) []*Kline {
	var list []*Kline
	var inters [][]interface{}
	sym = strings.ToUpper(sym)
	offs := strconv.FormatInt(offset, 10)
	urlReq := strings.ReplaceAll(urlKlines, "{SYM}", sym)
	urlReq = strings.ReplaceAll(urlReq, "{TM}", offs)
	resp, err := http.Get(urlReq)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(bts, &inters); err != nil {
		log.Fatal(err)
	}
	// =================================================
	for _, v := range inters {
		kl := Kline{}
		kl.OpenTime = int64(v[0].(float64))
		kl.Open = strToFloat(v[1].(string))
		kl.High = strToFloat(v[2].(string))
		kl.Low = strToFloat(v[3].(string))
		kl.Close = strToFloat(v[4].(string))
		kl.Volume = strToFloat(v[5].(string))
		kl.Trades = int64(v[8].(float64))
		list = append(list, &kl)
	}
	// =================================================
	return list
}

func reqBinanceLast(sym string) ([]*Kline, error) {
	var list []*Kline
	var inters [][]interface{}
	sym = strings.ToUpper(sym)
	urlReq := strings.ReplaceAll(urlKlinesLast, "{SYM}", sym)
	resp, err := http.Get(urlReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(bts, &inters); err != nil {
		return nil, err
	}
	// =================================================
	for _, v := range inters {
		kl := Kline{}
		kl.OpenTime = int64(v[0].(float64))
		kl.Open = strToFloat(v[1].(string))
		kl.High = strToFloat(v[2].(string))
		kl.Low = strToFloat(v[3].(string))
		kl.Close = strToFloat(v[4].(string))
		kl.Volume = strToFloat(v[5].(string))
		kl.Trades = int64(v[8].(float64))
		list = append(list, &kl)
	}
	// =================================================
	return list, nil
}
