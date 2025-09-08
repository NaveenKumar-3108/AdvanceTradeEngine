package processor

import (
	"AdvanceTradeEngine/models"
	"sync"
	"time"
)

type CandleAggregator struct {
	Interval time.Duration
	Candles  map[string]*models.Candle
	Mu       sync.RWMutex
}

func NewCandleAggregator(lInterval time.Duration) *CandleAggregator {
	return &CandleAggregator{
		Interval: lInterval,
		Candles:  make(map[string]*models.Candle),
	}
}

func (Agg *CandleAggregator) AddTrade(lTrade models.Trade) {
	Agg.Mu.Lock()
	defer Agg.Mu.Unlock()

	lKey := lTrade.Pair
	lTs := lTrade.Timestamp.Truncate(Agg.Interval)

	lCandle, lExists := Agg.Candles[lKey]
	if !lExists || !lCandle.Start.Equal(lTs) {
		lCandle = &models.Candle{
			Open:   lTrade.Price,
			High:   lTrade.Price,
			Low:    lTrade.Price,
			Close:  lTrade.Price,
			Volume: lTrade.Quantity,
			Start:  lTs,
		}
		Agg.Candles[lKey] = lCandle
	} else {
		if lTrade.Price > lCandle.High {
			lCandle.High = lTrade.Price
		}
		if lTrade.Price < lCandle.Low {
			lCandle.Low = lTrade.Price
		}
		lCandle.Close = lTrade.Price
		lCandle.Volume += lTrade.Quantity
	}
}
