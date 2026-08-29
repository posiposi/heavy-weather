package domain

import (
	"errors"
	"fmt"
)

// PrecipitationProbability は降水確率を表す値オブジェクトである。単位はパーセントで、0 以上 100 以下の整数を取る。
type PrecipitationProbability int

// ErrPrecipitationProbabilityOutOfRange は降水確率として取り得ない値が渡されたことを表す。
var ErrPrecipitationProbabilityOutOfRange = errors.New("precipitation probability out of range")

const (
	minPrecipitationProbability = 0
	maxPrecipitationProbability = 100
)

// NewPrecipitationProbability はパーセント表記の整数から PrecipitationProbability を生成する。
// 0 以上 100 以下でない場合はゼロ値と ErrPrecipitationProbabilityOutOfRange を返す。
func NewPrecipitationProbability(percent int) (PrecipitationProbability, error) {
	if percent < minPrecipitationProbability || percent > maxPrecipitationProbability {
		return 0, fmt.Errorf("%w: %d", ErrPrecipitationProbabilityOutOfRange, percent)
	}
	return PrecipitationProbability(percent), nil
}

// AtLeast は降水確率が threshold 以上かを返す。
func (p PrecipitationProbability) AtLeast(threshold PrecipitationProbability) bool {
	return p >= threshold
}

// String は "30%" の形式で降水確率を返す。
func (p PrecipitationProbability) String() string {
	return fmt.Sprintf("%d%%", int(p))
}
