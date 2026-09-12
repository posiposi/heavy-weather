package domain

import (
	"errors"
	"fmt"
	"time"
)

// ForecastDate は予報の対象日を表す値オブジェクトである。月と日の整合が取れた日付のみを取り、時刻とタイムゾーンを持たない。
type ForecastDate struct {
	value time.Time
}

// ErrForecastDateNotExist は暦上実在しない日付が渡されたことを表す。
var ErrForecastDateNotExist = errors.New("forecast date does not exist")

// NewForecastDate は年・月・日から ForecastDate を生成する。
func NewForecastDate(year int, month time.Month, day int) (*ForecastDate, error) {
	value := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	if value.Year() != year || value.Month() != month || value.Day() != day {
		return nil, fmt.Errorf("%w: %04d-%02d-%02d", ErrForecastDateNotExist, year, month, day)
	}
	return &ForecastDate{value: value}, nil
}
