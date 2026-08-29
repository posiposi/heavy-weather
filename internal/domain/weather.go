package domain

// Weather は天気を表す値オブジェクトである。
type Weather int

// Weather の取り得る値。ゼロ値は WeatherUnknown であり、天気が確定できなかったことを表す。
const (
	WeatherUnknown Weather = iota
	WeatherClearSky
	WeatherMainlyClear
	WeatherPartlyCloudy
	WeatherOvercast
	WeatherFog
	WeatherDepositingRimeFog
	WeatherLightDrizzle
	WeatherDrizzle
	WeatherHeavyDrizzle
	WeatherLightFreezingDrizzle
	WeatherHeavyFreezingDrizzle
	WeatherLightRain
	WeatherRain
	WeatherHeavyRain
	WeatherLightFreezingRain
	WeatherHeavyFreezingRain
	WeatherLightSnow
	WeatherSnow
	WeatherHeavySnow
	WeatherSnowGrains
	WeatherLightRainShowers
	WeatherRainShowers
	WeatherViolentRainShowers
	WeatherLightSnowShowers
	WeatherHeavySnowShowers
	WeatherThunderstorm
	WeatherThunderstormWithLightHail
	WeatherThunderstormWithHeavyHail

	weatherCount
)

var weatherByWMOCode = map[int]Weather{
	0:  WeatherClearSky,
	1:  WeatherMainlyClear,
	2:  WeatherPartlyCloudy,
	3:  WeatherOvercast,
	45: WeatherFog,
	48: WeatherDepositingRimeFog,
	51: WeatherLightDrizzle,
	53: WeatherDrizzle,
	55: WeatherHeavyDrizzle,
	56: WeatherLightFreezingDrizzle,
	57: WeatherHeavyFreezingDrizzle,
	61: WeatherLightRain,
	63: WeatherRain,
	65: WeatherHeavyRain,
	66: WeatherLightFreezingRain,
	67: WeatherHeavyFreezingRain,
	71: WeatherLightSnow,
	73: WeatherSnow,
	75: WeatherHeavySnow,
	77: WeatherSnowGrains,
	80: WeatherLightRainShowers,
	81: WeatherRainShowers,
	82: WeatherViolentRainShowers,
	85: WeatherLightSnowShowers,
	86: WeatherHeavySnowShowers,
	95: WeatherThunderstorm,
	96: WeatherThunderstormWithLightHail,
	99: WeatherThunderstormWithHeavyHail,
}

var weatherMeanings = map[Weather]string{
	WeatherUnknown:                   "不明",
	WeatherClearSky:                  "快晴",
	WeatherMainlyClear:               "晴れ",
	WeatherPartlyCloudy:              "一部曇り",
	WeatherOvercast:                  "曇り",
	WeatherFog:                       "霧",
	WeatherDepositingRimeFog:         "着氷性の霧",
	WeatherLightDrizzle:              "弱い霧雨",
	WeatherDrizzle:                   "霧雨",
	WeatherHeavyDrizzle:              "強い霧雨",
	WeatherLightFreezingDrizzle:      "弱い着氷性の霧雨",
	WeatherHeavyFreezingDrizzle:      "強い着氷性の霧雨",
	WeatherLightRain:                 "弱い雨",
	WeatherRain:                      "雨",
	WeatherHeavyRain:                 "強い雨",
	WeatherLightFreezingRain:         "弱い着氷性の雨",
	WeatherHeavyFreezingRain:         "強い着氷性の雨",
	WeatherLightSnow:                 "弱い雪",
	WeatherSnow:                      "雪",
	WeatherHeavySnow:                 "強い雪",
	WeatherSnowGrains:                "霧雪",
	WeatherLightRainShowers:          "弱いにわか雨",
	WeatherRainShowers:               "にわか雨",
	WeatherViolentRainShowers:        "激しいにわか雨",
	WeatherLightSnowShowers:          "弱いにわか雪",
	WeatherHeavySnowShowers:          "強いにわか雪",
	WeatherThunderstorm:              "雷雨",
	WeatherThunderstormWithLightHail: "弱い雹を伴う雷雨",
	WeatherThunderstormWithHeavyHail: "強い雹を伴う雷雨",
}

// NewWeather は WMO Weather interpretation code から Weather を生成する。
// 表に無いコードは WeatherUnknown となる。
func NewWeather(code int) Weather {
	w, ok := weatherByWMOCode[code]
	if !ok {
		return WeatherUnknown
	}
	return w
}

// Meaning は天気の意味を日本語で返す。
func (w Weather) Meaning() string {
	m, ok := weatherMeanings[w]
	if !ok {
		return weatherMeanings[WeatherUnknown]
	}
	return m
}
