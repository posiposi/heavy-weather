package domain

import (
	"errors"
	"testing"
	"time"
)

func TestNewForecastDate(t *testing.T) {
	tests := []struct {
		name string
		want time.Time
	}{
		{"実在する日付の2026-09-05", time.Date(2026, time.September, 5, 0, 0, 0, 0, time.UTC)},
		{"年始の2026-01-01", time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)},
		{"年末の2026-12-31", time.Date(2026, time.December, 31, 0, 0, 0, 0, time.UTC)},
		{"非閏年の2月末日の2026-02-28", time.Date(2026, time.February, 28, 0, 0, 0, 0, time.UTC)},
		{"閏年の2月末日の2028-02-29", time.Date(2028, time.February, 29, 0, 0, 0, 0, time.UTC)},
		{"30日までの月の末日の2026-04-30", time.Date(2026, time.April, 30, 0, 0, 0, 0, time.UTC)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			year, month, day := tt.want.Date()
			got, err := NewForecastDate(year, month, day)
			if err != nil {
				t.Fatalf("NewForecastDate(%v, %v, %v) のエラー = %v, want %v", year, month, day, err, nil)
			}
			if got.value != tt.want {
				t.Errorf("NewForecastDate(%v, %v, %v) = %v, want %v", year, month, day, got.value, tt.want)
			}
		})
	}
}

func TestNewForecastDateNotExist(t *testing.T) {
	tests := []struct {
		name  string
		year  int
		month time.Month
		day   int
	}{
		{"非閏年の2月29日", 2026, time.February, 29},
		{"2月30日は3月2日へ正規化されない", 2026, time.February, 30},
		{"月が13", 2026, time.Month(13), 5},
		{"日が32", 2026, time.September, 32},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewForecastDate(tt.year, tt.month, tt.day)
			if err == nil {
				t.Errorf("NewForecastDate(%v, %v, %v) のエラー = %v, want not %v", tt.year, tt.month, tt.day, err, nil)
			}
			if ok := errors.Is(err, ErrForecastDateNotExist); !ok {
				t.Errorf("errors.Is(NewForecastDate(%v, %v, %v) のエラー(%v), ErrForecastDateNotExist) = %v, want %v", tt.year, tt.month, tt.day, err, ok, true)
			}
			if got != nil {
				t.Errorf("NewForecastDate(%v, %v, %v) = %v, want %v", tt.year, tt.month, tt.day, got, nil)
			}
		})
	}
}

func TestForecastDateEquality(t *testing.T) {
	type date struct {
		year  int
		month time.Month
		day   int
	}
	tests := []struct {
		name  string
		left  date
		right date
		want  bool
	}{
		{"同じ日付同士は等しい", date{2026, time.September, 5}, date{2026, time.September, 5}, true},
		{"日だけが異なる日付同士は等しくない", date{2026, time.September, 5}, date{2026, time.September, 6}, false},
		{"月だけが異なる日付同士は等しくない", date{2026, time.September, 5}, date{2026, time.October, 5}, false},
		{"年だけが異なる日付同士は等しくない", date{2026, time.September, 5}, date{2027, time.September, 5}, false},
		{"月と日を入れ替えた日付同士は等しくない", date{2026, time.September, 5}, date{2026, time.May, 9}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			left, err := NewForecastDate(tt.left.year, tt.left.month, tt.left.day)
			if err != nil {
				t.Fatalf("NewForecastDate(%v, %v, %v) のエラー = %v, want %v", tt.left.year, tt.left.month, tt.left.day, err, nil)
			}
			right, err := NewForecastDate(tt.right.year, tt.right.month, tt.right.day)
			if err != nil {
				t.Fatalf("NewForecastDate(%v, %v, %v) のエラー = %v, want %v", tt.right.year, tt.right.month, tt.right.day, err, nil)
			}
			if got := *left == *right; got != tt.want {
				t.Errorf("*NewForecastDate(%v, %v, %v) == *NewForecastDate(%v, %v, %v) = %v, want %v", tt.left.year, tt.left.month, tt.left.day, tt.right.year, tt.right.month, tt.right.day, got, tt.want)
			}
		})
	}
}
