package data

import (
	"encoding/json"
	"fmt"
	"strings"
)

//Platform defines type of platform
type Platform int

const (
	_ Platform = iota
	//IOS represents apple ios devices
	IOS
	//ANDROID represents google android devices
	ANDROID
)

var (
	platformNameToValue = map[string]Platform{
		"IOS":     IOS,
		"ANDROID": ANDROID,
	}

	platformValueToName = map[Platform]string{
		IOS:     "IOS",
		ANDROID: "ANDROID",
	}
)

//MarshalJSON for type Platform
func (p Platform) MarshalJSON() ([]byte, error) {
	if s, ok := interface{}(p).(fmt.Stringer); ok {
		return json.Marshal(s.String())
	}
	s, ok := platformValueToName[p]
	if !ok {
		return nil, fmt.Errorf("invalid Platform: %d", p)
	}
	return json.Marshal(strings.ToLower(s))
}

//UnmarshalJSON for type Platform
func (p *Platform) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("Platform should be a string, got %s", data)
	}
	v, ok := platformNameToValue[strings.ToUpper(s)]
	if !ok {
		return fmt.Errorf("invalid Platform '%s'", s)
	}
	*p = v
	return nil
}

func ParsePlatform(platform string) (Platform, error) {
	v, ok := platformNameToValue[strings.ToUpper(platform)]
	if !ok {
		return 0, fmt.Errorf("invalid Platform '%s'", platform)
	}
	return Platform(v), nil
}

func (p Platform) ValueAsInt() int {
	return int(p)
}

func PlatformToInt(platfrom Platform) int {
	return int(platfrom)
}
