package domain

import (
	"errors"
	"fmt"
	"time"
)

// ForecastDate は予報の対象日を表す値オブジェクトである。YYYY-MM-DD 形式で表された日付のみを取り、時刻とタイムゾーンを持たない。
type ForecastDate struct {
	value time.Time
}

// ErrForecastDateInvalid は YYYY-MM-DD 形式の日付として解釈できない値が渡されたことを表す。
var ErrForecastDateInvalid = errors.New("forecast date is invalid")

// NewForecastDate は YYYY-MM-DD 形式の文字列から ForecastDate を生成する。
func NewForecastDate(value string) (*ForecastDate, error) {
	parsed, err := time.Parse(time.DateOnly, value)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrForecastDateInvalid, value)
	}
	return &ForecastDate{value: parsed}, nil
}
