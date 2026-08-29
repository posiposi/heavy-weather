package domain

type weatherKind int

const (
	weatherKindUnknown weatherKind = iota
	weatherKindClearSky
	weatherKindMainlyClear
	weatherKindPartlyCloudy
	weatherKindOvercast
	weatherKindFog
	weatherKindDepositingRimeFog
	weatherKindLightDrizzle
	weatherKindDrizzle
	weatherKindHeavyDrizzle
	weatherKindLightFreezingDrizzle
	weatherKindHeavyFreezingDrizzle
	weatherKindLightRain
	weatherKindRain
	weatherKindHeavyRain
	weatherKindLightFreezingRain
	weatherKindHeavyFreezingRain
	weatherKindLightSnow
	weatherKindSnow
	weatherKindHeavySnow
	weatherKindSnowGrains
	weatherKindLightRainShowers
	weatherKindRainShowers
	weatherKindViolentRainShowers
	weatherKindLightSnowShowers
	weatherKindHeavySnowShowers
	weatherKindThunderstorm
	weatherKindThunderstormWithLightHail
	weatherKindThunderstormWithHeavyHail

	weatherKindCount
)

// Weather は天気を表す値オブジェクトである。ゼロ値は天気が確定できなかったことを表す。
type Weather struct {
	kind weatherKind
}

var weatherKindByWMOCode = map[int]weatherKind{
	0:  weatherKindClearSky,
	1:  weatherKindMainlyClear,
	2:  weatherKindPartlyCloudy,
	3:  weatherKindOvercast,
	45: weatherKindFog,
	48: weatherKindDepositingRimeFog,
	51: weatherKindLightDrizzle,
	53: weatherKindDrizzle,
	55: weatherKindHeavyDrizzle,
	56: weatherKindLightFreezingDrizzle,
	57: weatherKindHeavyFreezingDrizzle,
	61: weatherKindLightRain,
	63: weatherKindRain,
	65: weatherKindHeavyRain,
	66: weatherKindLightFreezingRain,
	67: weatherKindHeavyFreezingRain,
	71: weatherKindLightSnow,
	73: weatherKindSnow,
	75: weatherKindHeavySnow,
	77: weatherKindSnowGrains,
	80: weatherKindLightRainShowers,
	81: weatherKindRainShowers,
	82: weatherKindViolentRainShowers,
	85: weatherKindLightSnowShowers,
	86: weatherKindHeavySnowShowers,
	95: weatherKindThunderstorm,
	96: weatherKindThunderstormWithLightHail,
	99: weatherKindThunderstormWithHeavyHail,
}

var weatherMeanings = map[weatherKind]string{
	weatherKindUnknown:                   "不明",
	weatherKindClearSky:                  "快晴",
	weatherKindMainlyClear:               "晴れ",
	weatherKindPartlyCloudy:              "一部曇り",
	weatherKindOvercast:                  "曇り",
	weatherKindFog:                       "霧",
	weatherKindDepositingRimeFog:         "着氷性の霧",
	weatherKindLightDrizzle:              "弱い霧雨",
	weatherKindDrizzle:                   "霧雨",
	weatherKindHeavyDrizzle:              "強い霧雨",
	weatherKindLightFreezingDrizzle:      "弱い着氷性の霧雨",
	weatherKindHeavyFreezingDrizzle:      "強い着氷性の霧雨",
	weatherKindLightRain:                 "弱い雨",
	weatherKindRain:                      "雨",
	weatherKindHeavyRain:                 "強い雨",
	weatherKindLightFreezingRain:         "弱い着氷性の雨",
	weatherKindHeavyFreezingRain:         "強い着氷性の雨",
	weatherKindLightSnow:                 "弱い雪",
	weatherKindSnow:                      "雪",
	weatherKindHeavySnow:                 "強い雪",
	weatherKindSnowGrains:                "霧雪",
	weatherKindLightRainShowers:          "弱いにわか雨",
	weatherKindRainShowers:               "にわか雨",
	weatherKindViolentRainShowers:        "激しいにわか雨",
	weatherKindLightSnowShowers:          "弱いにわか雪",
	weatherKindHeavySnowShowers:          "強いにわか雪",
	weatherKindThunderstorm:              "雷雨",
	weatherKindThunderstormWithLightHail: "弱い雹を伴う雷雨",
	weatherKindThunderstormWithHeavyHail: "強い雹を伴う雷雨",
}

// NewWeather は WMO Weather interpretation code から Weather を生成する。
// 表に無いコードは未確定として扱うため失敗せず、常に非 nil を返す。
func NewWeather(code int) *Weather {
	return &Weather{kind: weatherKindByWMOCode[code]}
}

// Meaning は天気の意味を日本語で返す。
func (w Weather) Meaning() string {
	return weatherMeanings[w.kind]
}
