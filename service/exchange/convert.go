package exchange

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type apiInfo struct {
	USDTWD usdtwd `json:"USDTWD"`
}
type usdtwd struct {
	Exrate float64 `json:"Exrate"`
}

func ConvertExchange(oldMoney []float64) (newMoney []float64, err error) {
	// 先取匯率
	res, err := http.Get("https://tw.rter.info/capi.php")
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	apiInfo := &apiInfo{}

	err = json.Unmarshal(body, apiInfo)
	if err != nil {
		return nil, err
	}

	exrateStr := apiInfo.USDTWD.Exrate

	for _, v := range oldMoney {
		newMoney = append(newMoney, v*exrateStr)
	}

	return newMoney, nil
}
