package logic

import (
	"encoding/json"
	"github.com/ildomm/eskalationszeit/zeitarbeiter/config"
	"github.com/ildomm/eskalationszeit/zeitarbeiter/models"
	"io/ioutil"
	"log"
	"net/http"
)

var Printf = log.Printf
func UpdatePrices(){

	// One against each other
	for _, currency := range convertables() {
		for _, convert := range convertables() {

			// Convert only against diffs
			if currency.Symbol != convert.Symbol {

				price := price(currency.Symbol, convert.Symbol)
				currency.UpdatePrice(convert.Symbol, float64(price) )
			}
		}
	}
}

func price( symbol, convert string ) float32 {
	resp, err := http.Get(config.App.Runtime.GeneratorUrl)
	if err != nil {
		log.Println("There was an error:", err)
		return 0
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	var price models.Price
	err = json.Unmarshal(body, &price)
	if err != nil {
		log.Println("There was an error:", err)
		return 0
	}

	return price.Price
}

func convertables() []*models.Currency {
	var entries []*models.Currency

	entries = append(entries, &models.Currency{Symbol: models.CurrencySymbolUSD} )
	entries = append(entries, &models.Currency{Symbol: models.CurrencySymbolBTC} )
	entries = append(entries, &models.Currency{Symbol: models.CurrencySymbolNZD} )

	return entries
}