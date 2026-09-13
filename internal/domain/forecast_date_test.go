package domain

import (
	"errors"
	"testing"
	"time"
)

func TestNewForecastDate(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  time.Time
	}{
		{"実在する日付の2026-09-05", "2026-09-05", time.Date(2026, time.September, 5, 0, 0, 0, 0, time.UTC)},
		{"年始の2026-01-01", "2026-01-01", time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)},
		{"年末の2026-12-31", "2026-12-31", time.Date(2026, time.December, 31, 0, 0, 0, 0, time.UTC)},
		{"非閏年の2月末日の2026-02-28", "2026-02-28", time.Date(2026, time.February, 28, 0, 0, 0, 0, time.UTC)},
		{"閏年の2月末日の2028-02-29", "2028-02-29", time.Date(2028, time.February, 29, 0, 0, 0, 0, time.UTC)},
		{"30日までの月の末日の2026-04-30", "2026-04-30", time.Date(2026, time.April, 30, 0, 0, 0, 0, time.UTC)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewForecastDate(tt.value)
			if err != nil {
				t.Fatalf("NewForecastDate(%v) のエラー = %v, want %v", tt.value, err, nil)
			}
			if got.value != tt.want {
				t.Errorf("NewForecastDate(%v) = %v, want %v", tt.value, got.value, tt.want)
			}
		})
	}
}

func TestNewForecastDateInvalid(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{"非閏年の2月29日", "2026-02-29"},
		{"2月30日", "2026-02-30"},
		{"月が13", "2026-13-05"},
		{"ゼロ埋めされていない", "2026-9-5"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewForecastDate(tt.value)
			if err == nil {
				t.Errorf("NewForecastDate(%v) のエラー = %v, want not %v", tt.value, err, nil)
			}
			if ok := errors.Is(err, ErrForecastDateInvalid); !ok {
				t.Errorf("errors.Is(NewForecastDate(%v) のエラー(%v), ErrForecastDateInvalid) = %v, want %v", tt.value, err, ok, true)
			}
			if got != nil {
				t.Errorf("NewForecastDate(%v) = %v, want %v", tt.value, got, nil)
			}
		})
	}
}

func TestForecastDateEquality(t *testing.T) {
	tests := []struct {
		name  string
		left  string
		right string
		want  bool
	}{
		{"同じ日付同士は等しい", "2026-09-05", "2026-09-05", true},
		{"日だけが異なる日付同士は等しくない", "2026-09-05", "2026-09-06", false},
		{"月だけが異なる日付同士は等しくない", "2026-09-05", "2026-10-05", false},
		{"年だけが異なる日付同士は等しくない", "2026-09-05", "2027-09-05", false},
		{"月と日を入れ替えた日付同士は等しくない", "2026-09-05", "2026-05-09", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			left, err := NewForecastDate(tt.left)
			if err != nil {
				t.Fatalf("NewForecastDate(%v) のエラー = %v, want %v", tt.left, err, nil)
			}
			right, err := NewForecastDate(tt.right)
			if err != nil {
				t.Fatalf("NewForecastDate(%v) のエラー = %v, want %v", tt.right, err, nil)
			}
			if got := *left == *right; got != tt.want {
				t.Errorf("*NewForecastDate(%v) == *NewForecastDate(%v) = %v, want %v", tt.left, tt.right, got, tt.want)
			}
		})
	}
}
