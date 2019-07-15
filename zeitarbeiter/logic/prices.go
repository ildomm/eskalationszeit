package logic

import (
	"github.com/ildomm/eskalationszeit/zeitarbeiter/models"
	"log"
)

var Printf = log.Printf
func UpdatePrices(){

	for _, c := range convertables() {
		log.Println(c)
	}

}

func convertables() []*models.Currency {
	var entries []*models.Currency

	entries = append(entries, &models.Currency{Symbol: models.CurrencySymbolUSD} )
	entries = append(entries, &models.Currency{Symbol: models.CurrencySymbolBTC} )
	entries = append(entries, &models.Currency{Symbol: models.CurrencySymbolNZD} )

	return entries
}