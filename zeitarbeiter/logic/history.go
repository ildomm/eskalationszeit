package logic

import (
	"encoding/json"
	"github.com/go-openapi/strfmt"
	"github.com/go-redis/redis"
	"github.com/ildomm/eskalationszeit/zeitarbeiter/databases"
	"github.com/ildomm/eskalationszeit/zeitarbeiter/models"
	"log"
	"reflect"
	"sort"
	"strconv"
	"time"
)

type Pair struct {
	Source *models.Currency
	Target *models.Currency
	found  bool
}

func UpdateHistories(){

	// One against each other
	for _, currency := range convertables() {
		for _, convert := range convertables() {

			// Convert only against diffs
			if currency.Symbol != convert.Symbol {
				do( &Pair{Source: currency, Target: convert} )
			}
		}
	}
}

func do(relation *Pair){
	if h1FirstTs(relation) > 0 {
		for _, period := range relation.Source.HistoryPeriods() {
			relation.found = true
			fill(relation, period)
		}
	}
}

func key(relation *Pair, period int) string {
	return relation.Source.HistoryKey(relation.Target.Symbol, period)
}

func h1FirstTs(relation *Pair) int {
	conn := databases.Redis()
	key := key(relation,1)
	scores := conn.ZRangeWithScores(key, 0, 0)
	for _, s := range scores.Val() {
		return int(s.Score)
	}
	return 0
}

// Try to find some valid time reference 10 steps back
func h1NearTs(relation *Pair, period int, ts int) int {
	left := ts - ( 60 * (period * 10) )
	right := ts

	conn := databases.Redis()
	key := key(relation,1)
	query := &redis.ZRangeBy{ Min: strconv.Itoa(left), Max: strconv.Itoa(right) }
	scores := conn.ZRevRangeByScoreWithScores(key, *query)

	for _, s := range scores.Val() {
		log.Printf("h1NearTs %s, ts %d, len scores %d, query %v", key, ts, len(scores.Val()), query )
		return int(s.Score)
	}

	return ts - (60 * ( period * 3))
}

func lastTs(relation *Pair, period int) int {
	conn := databases.Redis()
	key := key(relation,period)

	index := int64(relation.Source.HistoryIndex(key))
	scores := conn.ZRangeWithScores(key, index - 1, index)

	for _, s := range scores.Val() {
		return int(s.Score)
	}

	return 0
}

// Calculate next time step
// When does not have reference, then first H1 is the step
func nextTs(relation *Pair, period int, _ts *int) int {
	var ts int
	if _ts == nil {
		ts = lastTs(relation, period)
	} else {
		ts = *_ts
	}

	if ts == 0 {
		ts = h1FirstTs(relation)
	} else {
		ts = ts + ( 60 * period )
	}

	return ts
}

// H"1" is the Source of all data
// When start+end range does not math, must navigate back into H1 to locate nearest
// possible score
func knFromH1(relation *Pair, period int, ts int) *models.CurrencyMilestone {
	left := ts
	right := ts + (60 * period)

	conn := databases.Redis()
	key := key(relation,1)
	query := &redis.ZRangeBy{ Min: strconv.Itoa(left), Max: strconv.Itoa(right) }
	scores := conn.ZRangeByScore(key, *query)

	// If no data, try nearest score from H1
	if len(scores.Val()) == 0 {
		left = h1NearTs(relation, period, ts)
		query := &redis.ZRangeBy{ Min: strconv.Itoa(left), Max: strconv.Itoa(right) }
		scores = conn.ZRangeByScore(key, *query)
	}

	// Load to a simple float list
	var prices []float64
	for _, s := range scores.Val() {
		var price []*float64
		if err := json.Unmarshal([]byte(s), &price); err != nil {
			log.Println(err)
		}
		prices = append(prices, *price[1])
	}

	if len(prices) == 0 {
		return nil
	}

	// Apply math
	open := prices[0]
	close := prices[len(prices) - 1]

	// Yes, that is it, reorder list to catch values from "edges"
	// Expensive/dangerous approach, but fast using small lists
	sort.Float64s(prices)

	min := prices[0]
	max := prices[len(prices) - 1]

	return &models.CurrencyMilestone{ Open: open, Max: max, Min: min, Close: close }
}

// Load very last db entry
func knLast(relation *Pair, period int, ts int) *models.CurrencyMilestone {
	conn := databases.Redis()
	key := key(relation,period)

	index := int64(relation.Source.HistoryIndex(key))
	scores := conn.ZRange(key, index - 1, index)

	for _, s := range scores.Val() {
		var price []*float64
		if err := json.Unmarshal([]byte(s), &price); err != nil {
			log.Println(err)
		}

		return &models.CurrencyMilestone{ Open: *price[1], Max: *price[2], Min: *price[3], Close: *price[4] }
	}

	// In case of no data at all
	return nil
}

// Construct milestone based on H1
// or just load last older entry
func getPoint(relation *Pair, period int, ts int) *models.CurrencyMilestone {
	point := knFromH1(relation, period, ts)

	if point == nil {
		point = knLast(relation, period, ts)
		relation.found = false
	} else {
		relation.found = true
	}

	return point
}

// Include new, without repetition
func appendPoint(relation *Pair, period int, ts int, last *models.CurrencyMilestone, nearest *models.CurrencyMilestone, forceUpdate bool) *models.CurrencyMilestone {
	key := key(relation,period)
	point := getPoint(relation, period, ts)

	if point == nil {
		return nil
	}

	// Not allowed repetition
	if allowPoint(point, last, nearest, forceUpdate) {
		info(forceUpdate, key, ts, point)
		relation.Source.UpdateHistory(relation.Target.Symbol, float64(ts), period, *point)
	}
	return point
}

func info(forceUpdate bool, key string, ts int, point *models.CurrencyMilestone) {
	// Apply format
	at_ := time.Unix(int64(ts), 0)
	at, _ := strfmt.ParseDateTime(at_.Format(time.RFC3339))
	point.At = at

	if forceUpdate {
		log.Printf("Update %s, ts %d, point %+v", key, ts, point)
	} else {
		log.Printf("Append %s, ts %d, point %+v", key, ts, point)
	}
}

func allowPoint(point *models.CurrencyMilestone, last *models.CurrencyMilestone, nearest *models.CurrencyMilestone, forceUpdate bool) bool {
	if forceUpdate {
		return true
	}

	if last == nil {
		return true
	}

	if nearest == nil {
		nearest = last
	}

	// Must compare with older and recent entry
	candidate := []float64{point.Open,point.Max,point.Min,point.Close}
	last_     := []float64{last.Open,last.Max,last.Min,last.Close}
	nearest_  := []float64{nearest.Open,nearest.Max,nearest.Min,nearest.Close}

	// Not allowed repetition
	if reflect.DeepEqual(candidate, last_) || reflect.DeepEqual(candidate, nearest_)  {
		return false
	}
	return true
}

func updatePoint(relation *Pair, period int, ts int, lastPoint *models.CurrencyMilestone) {
	key := key(relation,period)

	conn := databases.Redis()
	conn.ZRemRangeByScore(key, strconv.Itoa(ts), strconv.Itoa(ts))

	appendPoint(relation, period, ts, lastPoint, lastPoint,true)
}

func fill(relation *Pair, period int) {
	var _lastTs int
	var lastPoint *models.CurrencyMilestone
	var nearestPoint *models.CurrencyMilestone

	// Calculate next time step
	ts := nextTs(relation, period, nil)

	// Only plays if there are possible new data inside H"1"
	for ts > 0 && ( int64(ts) <= time.Now().Unix() ) {

		// Save two db connections during every period loop
		if lastPoint == nil {
			_lastTs = lastTs(relation, period)
			lastPoint = knLast(relation, period, _lastTs)
		}

		// Include new, without repetition
		// Must compare with older and recent entry
		nearestPoint = appendPoint(relation, period, ts, lastPoint, nearestPoint, false)

		// To avoid populate next steps with same values when H"1" is down
		if relation.found == false {
			log.Printf("Key %s: H1 data not found for ts %d", key(relation,period), ts)
		}

		if nearestPoint == nil {
			break
		}

		ts = nextTs(relation, period, &ts)
	}

    // Has to refresh OLDER point based in new recent data
    // but once, and only once
	if lastPoint != nil {
		updatePoint(relation, period, _lastTs, lastPoint)
	}
}