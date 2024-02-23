package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sfvn-hbduy/common/cache"
	"sfvn-hbduy/common/model"
	"sfvn-hbduy/common/response"
	"sfvn-hbduy/common/util"
	"strconv"
	"time"
)

type IHistoriesService interface {
	GetHistories(ctx context.Context, symbol string, days int, period string) (int, any)
}

type Histories struct {
	BaseURL string
	APIKey  string
}

func NewHistories(baseUrl, apiKey string) IHistoriesService {
	return &Histories{
		BaseURL: baseUrl,
		APIKey:  apiKey,
	}
}

func (s *Histories) GetHistories(ctx context.Context, symbol string, days int, period string) (int, any) {
	histories, err := s.cacheHistories(ctx, symbol, days, period)
	if err != nil {
		return response.BadRequestMsg(err)
	}

	return response.OK(histories)
}

func (s *Histories) cacheHistories(ctx context.Context, symbol string, days int, period string) ([]model.HistoricalData, error) {
	key := fmt.Sprintf("cache_histories_%s_%d", symbol, days)
	var historicalData []model.HistoricalData
	if historiesCached, err := cache.MCache.Get(key); err != nil || historiesCached == nil {
		result, err := util.CallAPI("GET", s.BaseURL+"/coins/"+symbol+"/ohlc", map[string]string{
			"vs_currency": "usd",
			"days":        strconv.Itoa(days),
			"api_key":     s.APIKey,
		})
		if err != nil {
			return historicalData, err
		}
		var historiesResp [][]float64
		err = json.Unmarshal(result, &historiesResp)
		if err != nil {
			return historicalData, err
		}
		var historicalDataTmp []model.HistoricalData
		var orginNumber float64
		var percentageChange float64
		for _, item := range historiesResp {
			change := math.Abs(orginNumber - item[4])
			if orginNumber != 0 {
				percentageChange = (change / orginNumber) * 100
			}
			historicalDataTmp = append(historicalDataTmp, model.HistoricalData{
				High:   item[2],
				Low:    item[3],
				Open:   item[1],
				Close:  item[4],
				Time:   int64(item[0]),
				Change: percentageChange,
			})
			orginNumber = item[4]
		}
		historicalData = historicalDataTmp

		var timeTTL time.Duration
		if days == 1 {
			timeTTL = 30 * time.Minute
		} else if days >= 3 && days <= 30 {
			timeTTL = 4 * time.Hour
		} else {
			timeTTL = 4 * 24 * time.Hour
		}

		if err := cache.MCache.SetTTL(key, historicalData, timeTTL); err != nil {
			return historicalData, err
		}
	} else {
		historicalDataTmp, _ := historiesCached.([]model.HistoricalData)
		historicalData = historicalDataTmp
	}

	return historicalData, nil
}
