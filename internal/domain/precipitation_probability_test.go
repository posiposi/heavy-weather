package domain

import (
	"errors"
	"math"
	"testing"
)

func TestNewPrecipitationProbability(t *testing.T) {
	tests := []struct {
		name       string
		percent    int
		want       int
		wantString string
	}{
		{"下限の0", 0, 0, "0%"},
		{"中間の50", 50, 50, "50%"},
		{"上限の100", 100, 100, "100%"},
		{"1", 1, 1, "1%"},
		{"30", 30, 30, "30%"},
		{"99", 99, 99, "99%"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPrecipitationProbability(tt.percent)
			if err != nil {
				t.Fatalf("NewPrecipitationProbability(%v) のエラー = %v, want %v", tt.percent, err, nil)
			}
			if got.percent != tt.want {
				t.Errorf("NewPrecipitationProbability(%v) = %v, want %v", tt.percent, got.percent, tt.want)
			}
			if got.String() != tt.wantString {
				t.Errorf("NewPrecipitationProbability(%v).String() = %v, want %v", tt.percent, got.String(), tt.wantString)
			}
		})
	}
}

func TestNewPrecipitationProbabilityOutOfRange(t *testing.T) {
	tests := []struct {
		name    string
		percent int
	}{
		{"下限より1小さい-1", -1},
		{"上限より1大きい101", 101},
		{"極端に小さい値", math.MinInt},
		{"極端に大きい値", math.MaxInt},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPrecipitationProbability(tt.percent)
			if err == nil {
				t.Errorf("NewPrecipitationProbability(%v) のエラー = %v, want not %v", tt.percent, err, nil)
			}
			if ok := errors.Is(err, ErrPrecipitationProbabilityOutOfRange); !ok {
				t.Errorf("errors.Is(NewPrecipitationProbability(%v) のエラー(%v), ErrPrecipitationProbabilityOutOfRange) = %v, want %v", tt.percent, err, ok, true)
			}
			var zero PrecipitationProbability
			if got != zero {
				t.Errorf("NewPrecipitationProbability(%v) = %v, want %v", tt.percent, got, zero)
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
			if got := left == right; got != tt.want {
				t.Errorf("NewPrecipitationProbability(%v) == NewPrecipitationProbability(%v) = %v, want %v", tt.left, tt.right, got, tt.want)
			}
		})
	}
}

func TestPrecipitationProbabilityAtLeast(t *testing.T) {
	tests := []struct {
		name      string
		percent   int
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
			p, err := NewPrecipitationProbability(tt.percent)
			if err != nil {
				t.Fatalf("NewPrecipitationProbability(%v) のエラー = %v, want %v", tt.percent, err, nil)
			}
			threshold, err := NewPrecipitationProbability(tt.threshold)
			if err != nil {
				t.Fatalf("NewPrecipitationProbability(%v) のエラー = %v, want %v", tt.threshold, err, nil)
			}
			if got := p.AtLeast(threshold); got != tt.want {
				t.Errorf("NewPrecipitationProbability(%v).AtLeast(%v) = %v, want %v", tt.percent, tt.threshold, got, tt.want)
			}
		})
	}
}
