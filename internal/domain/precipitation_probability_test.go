package domain

import (
	"errors"
	"math"
	"testing"
)

func TestNewPrecipitationProbability(t *testing.T) {
	tests := []struct {
		name  string
		value int
		want  int
	}{
		{"下限の0", 0, 0},
		{"中間の50", 50, 50},
		{"上限の100", 100, 100},
		{"1", 1, 1},
		{"30", 30, 30},
		{"99", 99, 99},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPrecipitationProbability(tt.value)
			if err != nil {
				t.Fatalf("NewPrecipitationProbability(%v) のエラー = %v, want %v", tt.value, err, nil)
			}
			if got.value != tt.want {
				t.Errorf("NewPrecipitationProbability(%v) = %v, want %v", tt.value, got.value, tt.want)
			}
		})
	}
}

func TestNewPrecipitationProbabilityOutOfRange(t *testing.T) {
	tests := []struct {
		name  string
		value int
	}{
		{"下限より1小さい-1", -1},
		{"上限より1大きい101", 101},
		{"極端に小さい値", math.MinInt},
		{"極端に大きい値", math.MaxInt},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPrecipitationProbability(tt.value)
			if err == nil {
				t.Errorf("NewPrecipitationProbability(%v) のエラー = %v, want not %v", tt.value, err, nil)
			}
			if ok := errors.Is(err, ErrPrecipitationProbabilityOutOfRange); !ok {
				t.Errorf("errors.Is(NewPrecipitationProbability(%v) のエラー(%v), ErrPrecipitationProbabilityOutOfRange) = %v, want %v", tt.value, err, ok, true)
			}
			if got != nil {
				t.Errorf("NewPrecipitationProbability(%v) = %v, want %v", tt.value, got, nil)
			}
		})
	}
}

func TestPrecipitationProbabilityEquality(t *testing.T) {
	tests := []struct {
		name  string
		left  int
		right int
		want  bool
	}{
		{"同じ値同士は等しい", 30, 30, true},
		{"異なる値同士は等しくない", 30, 31, false},
		{"下限と上限は等しくない", 0, 100, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			left, err := NewPrecipitationProbability(tt.left)
			if err != nil {
				t.Fatalf("NewPrecipitationProbability(%v) のエラー = %v, want %v", tt.left, err, nil)
			}
			right, err := NewPrecipitationProbability(tt.right)
			if err != nil {
				t.Fatalf("NewPrecipitationProbability(%v) のエラー = %v, want %v", tt.right, err, nil)
			}
			if got := *left == *right; got != tt.want {
				t.Errorf("*NewPrecipitationProbability(%v) == *NewPrecipitationProbability(%v) = %v, want %v", tt.left, tt.right, got, tt.want)
			}
		})
	}
}

func TestPrecipitationProbabilityAtLeast(t *testing.T) {
	tests := []struct {
		name      string
		value     int
		threshold int
		want      bool
	}{
		{"閾値ちょうどはtrue", 50, 50, true},
		{"閾値を上回る場合はtrue", 51, 50, true},
		{"閾値を下回る場合はfalse", 49, 50, false},
		{"閾値0は下限の0でもtrue", 0, 0, true},
		{"閾値100は上限の100でtrue", 100, 100, true},
		{"閾値100は99でfalse", 99, 100, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewPrecipitationProbability(tt.value)
			if err != nil {
				t.Fatalf("NewPrecipitationProbability(%v) のエラー = %v, want %v", tt.value, err, nil)
			}
			threshold, err := NewPrecipitationProbability(tt.threshold)
			if err != nil {
				t.Fatalf("NewPrecipitationProbability(%v) のエラー = %v, want %v", tt.threshold, err, nil)
			}
			if got := p.AtLeast(threshold); got != tt.want {
				t.Errorf("NewPrecipitationProbability(%v).AtLeast(%v) = %v, want %v", tt.value, tt.threshold, got, tt.want)
			}
		})
	}
}
