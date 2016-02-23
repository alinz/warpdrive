package data

import (
	"encoding/json"
	"fmt"
	"strings"
)

type FileType int

const (
	//JS represents file type as javascript
	JS FileType = iota
)

var (
	fileTypeNameToValue = map[string]FileType{
		"JS": JS,
	}

	fileTypeValueToName = map[FileType]string{
		JS: "JS",
	}
)

//MarshalJSON for type FileType
func (ft FileType) MarshalJSON() ([]byte, error) {
	if s, ok := interface{}(ft).(fmt.Stringer); ok {
		return json.Marshal(s.String())
	}
	s, ok := fileTypeValueToName[ft]
	if !ok {
		return nil, fmt.Errorf("invalid FileType: %d", ft)
	}
	return json.Marshal(strings.ToLower(s))
}

//UnmarshalJSON for type FileType
func (ft *FileType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("FileType should be a string, got %s", data)
	}
	v, ok := fileTypeNameToValue[strings.ToUpper(s)]
	if !ok {
		return fmt.Errorf("invalid FileType '%s'", s)
	}
	*ft = v
	return nil
}

func ParseFileType(fileType string) (FileType, error) {
	v, ok := fileTypeNameToValue[strings.ToUpper(fileType)]
	if !ok {
		return 0, fmt.Errorf("invalid FileType '%s'", fileType)
	}
	return FileType(v), nil
}
