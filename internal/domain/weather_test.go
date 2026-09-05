package domain

import "testing"

func TestNewWeather(t *testing.T) {
	tests := []struct {
		name        string
		code        int
		wantMeaning string
	}{
		{"0は快晴", 0, "快晴"},
		{"1は晴れ", 1, "晴れ"},
		{"2は一部曇り", 2, "一部曇り"},
		{"3は曇り", 3, "曇り"},
		{"45は霧", 45, "霧"},
		{"48は着氷性の霧", 48, "着氷性の霧"},
		{"51は弱い霧雨", 51, "弱い霧雨"},
		{"53は霧雨", 53, "霧雨"},
		{"55は強い霧雨", 55, "強い霧雨"},
		{"56は弱い着氷性の霧雨", 56, "弱い着氷性の霧雨"},
		{"57は強い着氷性の霧雨", 57, "強い着氷性の霧雨"},
		{"61は弱い雨", 61, "弱い雨"},
		{"63は雨", 63, "雨"},
		{"65は強い雨", 65, "強い雨"},
		{"66は弱い着氷性の雨", 66, "弱い着氷性の雨"},
		{"67は強い着氷性の雨", 67, "強い着氷性の雨"},
		{"71は弱い雪", 71, "弱い雪"},
		{"73は雪", 73, "雪"},
		{"75は強い雪", 75, "強い雪"},
		{"77は霧雪", 77, "霧雪"},
		{"80は弱いにわか雨", 80, "弱いにわか雨"},
		{"81はにわか雨", 81, "にわか雨"},
		{"82は激しいにわか雨", 82, "激しいにわか雨"},
		{"85は弱いにわか雪", 85, "弱いにわか雪"},
		{"86は強いにわか雪", 86, "強いにわか雪"},
		{"95は雷雨", 95, "雷雨"},
		{"96は弱い雹を伴う雷雨", 96, "弱い雹を伴う雷雨"},
		{"99は強い雹を伴う雷雨", 99, "強い雹を伴う雷雨"},
		{"表に無いコード4は不明", 4, "不明"},
		{"表に無いコード52は不明", 52, "不明"},
		{"表に無いコード100は不明", 100, "不明"},
		{"負数は不明", -1, "不明"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewWeather(tt.code)
			if got == nil {
				t.Errorf("NewWeather(%v) = nil, want non-nil", tt.code)
				return
			}
			if got.Meaning() != tt.wantMeaning {
				t.Errorf("NewWeather(%v).Meaning() = %v, want %v", tt.code, got.Meaning(), tt.wantMeaning)
			}
		})
	}
}

func TestZeroValueWeatherMeaning(t *testing.T) {
	var zero Weather
	want := "不明"
	if got := zero.Meaning(); got != want {
		t.Errorf("Weather{}.Meaning() = %v, want %v", got, want)
	}
}

func TestWeatherMeaningsCoversAllWeather(t *testing.T) {
	t.Run("WMOコード表のすべての天気が意味を持つ", func(t *testing.T) {
		unknown := weatherMeanings[weatherKindUnknown]
		for code := range weatherKindByWMOCode {
			if got := NewWeather(code).Meaning(); got == unknown || got == "" {
				t.Errorf("NewWeather(%v).Meaning() = %q, want neither %q nor empty", code, got, unknown)
			}
		}
	})

	t.Run("意味の表が全区分を網羅している", func(t *testing.T) {
		if got, want := len(weatherMeanings), int(weatherKindCount); got != want {
			t.Errorf("len(weatherMeanings) = %v, want %v", got, want)
		}
		for i := 0; i < int(weatherKindCount); i++ {
			if _, ok := weatherMeanings[weatherKind(i)]; !ok {
				t.Errorf("weatherMeanings[weatherKind(%d)] の登録 = %v, want %v", i, ok, true)
			}
		}
	})
}

func TestWeatherEquality(t *testing.T) {
	t.Run("同じ天気は値として等しい", func(t *testing.T) {
		fromCode4 := NewWeather(4)
		fromCode52 := NewWeather(52)
		if *fromCode4 != *fromCode52 {
			t.Errorf("*NewWeather(4) = %v, want %v (*NewWeather(52) と等しい)", *fromCode4, *fromCode52)
		}
	})

	t.Run("異なる天気は値として等しくない", func(t *testing.T) {
		clearSky := NewWeather(0)
		rain := NewWeather(63)
		if *clearSky == *rain {
			t.Errorf("*NewWeather(0) = %v, want not %v (*NewWeather(63) と等しくない)", *clearSky, *rain)
		}
	})
}
