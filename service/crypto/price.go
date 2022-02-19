package crypto

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type (
	apiInfo struct {
		Ethereum usd `json:"ethereum"`
	}

	usd struct {
		Price float64 `json:"usd"`
	}
)

func GetPrice() (float64, error) {
	res, err := http.Get("https://api.coingecko.com/api/v3/simple/price?ids=ethereum&vs_currencies=usd")
	if err != nil {
		return 0, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	apiInfo := &apiInfo{
		Ethereum: usd{},
	}

	err = json.Unmarshal(body, apiInfo)
	if err != nil {
		return 0, err
	}

	return apiInfo.Ethereum.Price, nil
}
