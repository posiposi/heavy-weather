package domain

import (
	"errors"
	"fmt"
)

// PrecipitationProbability は降水確率を表す値オブジェクトである。単位はパーセントで、0 以上 100 以下の整数を取る。
type PrecipitationProbability struct {
	value int
}

// ErrPrecipitationProbabilityOutOfRange は降水確率として取り得ない値が渡されたことを表す。
var ErrPrecipitationProbabilityOutOfRange = errors.New("precipitation probability out of range")

const (
	minPrecipitationProbability = 0
	maxPrecipitationProbability = 100
)

// NewPrecipitationProbability はパーセント表記の整数から PrecipitationProbability を生成する。
func NewPrecipitationProbability(value int) (*PrecipitationProbability, error) {
	if value < minPrecipitationProbability || value > maxPrecipitationProbability {
		return nil, fmt.Errorf("%w: %d", ErrPrecipitationProbabilityOutOfRange, value)
	}
	return &PrecipitationProbability{value: value}, nil
}
