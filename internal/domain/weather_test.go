package domain

import "testing"

func TestNewWeather(t *testing.T) {
	tests := []struct {
		name        string
		code        int
		want        Weather
		wantMeaning string
	}{
		{"0は快晴", 0, WeatherClearSky, "快晴"},
		{"1は晴れ", 1, WeatherMainlyClear, "晴れ"},
		{"2は一部曇り", 2, WeatherPartlyCloudy, "一部曇り"},
		{"3は曇り", 3, WeatherOvercast, "曇り"},
		{"45は霧", 45, WeatherFog, "霧"},
		{"48は着氷性の霧", 48, WeatherDepositingRimeFog, "着氷性の霧"},
		{"51は弱い霧雨", 51, WeatherLightDrizzle, "弱い霧雨"},
		{"53は霧雨", 53, WeatherDrizzle, "霧雨"},
		{"55は強い霧雨", 55, WeatherHeavyDrizzle, "強い霧雨"},
		{"56は弱い着氷性の霧雨", 56, WeatherLightFreezingDrizzle, "弱い着氷性の霧雨"},
		{"57は強い着氷性の霧雨", 57, WeatherHeavyFreezingDrizzle, "強い着氷性の霧雨"},
		{"61は弱い雨", 61, WeatherLightRain, "弱い雨"},
		{"63は雨", 63, WeatherRain, "雨"},
		{"65は強い雨", 65, WeatherHeavyRain, "強い雨"},
		{"66は弱い着氷性の雨", 66, WeatherLightFreezingRain, "弱い着氷性の雨"},
		{"67は強い着氷性の雨", 67, WeatherHeavyFreezingRain, "強い着氷性の雨"},
		{"71は弱い雪", 71, WeatherLightSnow, "弱い雪"},
		{"73は雪", 73, WeatherSnow, "雪"},
		{"75は強い雪", 75, WeatherHeavySnow, "強い雪"},
		{"77は霧雪", 77, WeatherSnowGrains, "霧雪"},
		{"80は弱いにわか雨", 80, WeatherLightRainShowers, "弱いにわか雨"},
		{"81はにわか雨", 81, WeatherRainShowers, "にわか雨"},
		{"82は激しいにわか雨", 82, WeatherViolentRainShowers, "激しいにわか雨"},
		{"85は弱いにわか雪", 85, WeatherLightSnowShowers, "弱いにわか雪"},
		{"86は強いにわか雪", 86, WeatherHeavySnowShowers, "強いにわか雪"},
		{"95は雷雨", 95, WeatherThunderstorm, "雷雨"},
		{"96は弱い雹を伴う雷雨", 96, WeatherThunderstormWithLightHail, "弱い雹を伴う雷雨"},
		{"99は強い雹を伴う雷雨", 99, WeatherThunderstormWithHeavyHail, "強い雹を伴う雷雨"},
		{"表に無いコード4は不明", 4, WeatherUnknown, "不明"},
		{"表に無いコード52は不明", 52, WeatherUnknown, "不明"},
		{"表に無いコード100は不明", 100, WeatherUnknown, "不明"},
		{"負数は不明", -1, WeatherUnknown, "不明"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewWeather(tt.code)
			if got != tt.want {
				t.Errorf("NewWeather(%v) = %d, want %d", tt.code, int(got), int(tt.want))
			}
			if got.Meaning() != tt.wantMeaning {
				t.Errorf("NewWeather(%v).Meaning() = %v, want %v", tt.code, got.Meaning(), tt.wantMeaning)
			}
		})
	}
}

func TestMeaningForUndefinedValue(t *testing.T) {
	undefined := Weather(weatherCount)
	want := "不明"
	if got := undefined.Meaning(); got != want {
		t.Errorf("Weather(%v).Meaning() = %v, want %v", int(undefined), got, want)
	}
}

func TestWeatherMeaningsCoversAllWeather(t *testing.T) {
	t.Run("WMOコード表のすべての天気が意味を持つ", func(t *testing.T) {
		unknown := WeatherUnknown.Meaning()
		for code, w := range weatherByWMOCode {
			if got := w.Meaning(); got == unknown {
				t.Errorf("NewWeather(%v).Meaning() = %v, want not %v", code, got, unknown)
			}
		}
	})

	t.Run("意味の表が全定数を網羅している", func(t *testing.T) {
		if got, want := len(weatherMeanings), int(weatherCount); got != want {
			t.Errorf("len(weatherMeanings) = %v, want %v", got, want)
		}
		for i := 0; i < int(weatherCount); i++ {
			w := Weather(i)
			if _, ok := weatherMeanings[w]; !ok {
				t.Errorf("weatherMeanings[Weather(%d)] の登録 = %v, want %v", i, ok, true)
			}
		}
	})
}

func TestWeatherEquality(t *testing.T) {
	fromCode4 := NewWeather(4)
	fromCode52 := NewWeather(52)
	if fromCode4 != fromCode52 {
		t.Errorf("NewWeather(4) = %v(%d), want %v(%d) (NewWeather(52) と等しい)", fromCode4, int(fromCode4), fromCode52, int(fromCode52))
	}
}
