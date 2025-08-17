package getter

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func d(in interface{}) {
	bts, err := json.Marshal(in)
	if err != nil {
		log.Fatal(err)
	}
	println(string(bts))
	os.Exit(0)
}

func strToFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func loadDump(exch, sym string) bool {
	exch = strings.ToLower(exch)
	sym = strings.ToLower(sym)
	fl := strings.ReplaceAll(dumpFileList, "{EXCH}", exch)
	fl = strings.ReplaceAll(fl, "{SYM}", sym)
	bts, err := ioutil.ReadFile(fl)
	if err != nil {
		return false
	}
	if err := json.Unmarshal(bts, &Klines); err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

func saveDump(exch, sym string, data []*Kline) {
	exch = strings.ToLower(exch)
	sym = strings.ToLower(sym)
	fl := strings.ReplaceAll(dumpFileList, "{EXCH}", exch)
	fl = strings.ReplaceAll(fl, "{SYM}", sym)
	bts, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile(fl, bts, 0644); err != nil {
		log.Fatal(err)
	}
}
