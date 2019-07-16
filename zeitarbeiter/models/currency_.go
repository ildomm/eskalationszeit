package models

import (
	"encoding/json"
	"github.com/go-openapi/strfmt"
	"github.com/go-redis/redis"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
	"github.com/ildomm/eskalationszeit/zeitarbeiter/databases"
)

func (m *Currency) Price(convert *string) float32 {
	return m.LoadPrice(convert)
}

func (m *Currency) LoadPrice(otherCoinSymbol *string) float32 {
	redis := databases.Redis()
	value := redis.Get(m.PriceKey(*otherCoinSymbol))
	if value != nil {
		price, _ := strconv.ParseFloat(value.Val(), 32)
		return float32(price)
	}
	return 0
}

func (m *Currency) UpdatePrice(convert string, price float64) {

	// Do not accept weird values
	if price == math.Inf(+1) ||
		price == math.Inf(-1) ||
		price == math.Copysign(0, -1) ||
		price == math.NaN() {
		return
	}

	if price > 0 {
		log.Printf("Price convert %s to %s = %03.8f", m.Symbol, convert, price)
		conn := databases.Redis()
		conn.Set(m.PriceKey(convert), price, (time.Duration(math.MaxInt16) * time.Hour))

		score := float64( time.Now().Unix() )

		milestone := &CurrencyMilestone{
			Open: price,
			Max: 0,
			Min: 0,
			Close: 0,
		}

		m.UpdateHistory(convert, score, 1, *milestone)
	}
}

func (m *Currency) PriceKey(otherCoinSymbol string) string {
	return "zeitar:currency:" + strings.ToLower(m.Symbol) + "_" + strings.ToLower(otherCoinSymbol) + ":now"
}

func (m *Currency) UpdateHistory(otherCoinSymbol string, score float64, period int, milestone CurrencyMilestone) {
	if milestone.Open > 0 {
		conn := databases.Redis()
		key := m.HistoryKey(otherCoinSymbol, period)

		index := m.HistoryIndex(key) + 1

		zValue := new(redis.Z)
		if period == 1 {
			value := []*float64{&index, &milestone.Open}
			jValue, err := json.Marshal(value)
			if err != nil {
				log.Fatal(err)
			}

			zValue = &redis.Z{
				Score:  score,
				Member: jValue }

		} else {
			value := []*float64{&index, &milestone.Open, &milestone.Max, &milestone.Min, &milestone.Close}
			jValue, err := json.Marshal(value)
			if err != nil {
				log.Fatal(err)
			}
			zValue = &redis.Z{
				Score:  score,
				Member: string(jValue) }
		}

		err := conn.ZAdd(m.HistoryKey(otherCoinSymbol, period), *zValue)

		if err != nil {
			if err.Err() != nil {
				log.Fatal(err)
			}
		}
	}
}

/* Returns the last index value of sorted list */
func (m *Currency) HistoryIndex(key string) float64 {
	conn := databases.Redis()
	_index := int64(0)
	_maxIndex := conn.ZCard(key)
	if _maxIndex == nil { _index = 0 }

	_index = _maxIndex.Val()
	if _index == 0 { _index = 0 }

	return float64( _index)
}

func (m *Currency) HistoryMilestones(pair string, period int, start int64, end int64) []*CurrencyMilestone {
	var milestones []*CurrencyMilestone
	conn := databases.Redis()

	key := m.HistoryKey(pair, period)
	query := &redis.ZRangeBy{ Min: strconv.FormatInt(start, 10), Max: strconv.FormatInt(end, 10) }
	scores := conn.ZRangeByScoreWithScores(key, *query)

	for _, s := range scores.Val() {

		m := (s.Member).(string)
		var p []*float64
		if err := json.Unmarshal([]byte(m), &p); err != nil {
			log.Println(err)
		}

		at_ := time.Unix(int64(s.Score), 0)
		at, _ := strfmt.ParseDateTime(at_.Format(time.RFC3339))

		milestone := &CurrencyMilestone{ At: at, Open: *p[1], Max: *p[2], Min: *p[3], Close: *p[4] }
		milestones = append(milestones, milestone)
	}

	return milestones
}

func (m *Currency) HistoryKey(otherCoinSymbol string, period int) string {
	return "zeitar:currency:" + strings.ToLower(m.Symbol) + "_" + strings.ToLower(otherCoinSymbol) + ":history:" + strconv.Itoa(period)
}

func (m *Currency) HistoryPeriods() []int {
	return []int{5, 15, 30, 60, 120, 720, 1440, 43800, 262800}
}

func (m *Currency) ConvertHistoryPeriod(period string) int {
	var periods = struct{M map[string]int }{M: map[string]int{
		"5": 5,
		"15": 15,
		"30": 30,
		"60": 60,
		"120": 120,
		"720": 720,
		"1D": 1440,
		"1M": 43800,
		"6M": 262800,
	}}
	return periods.M[period]
}