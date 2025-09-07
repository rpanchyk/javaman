package models

import "fmt"

type Os int

const (
	Linux Os = iota
	Macos
	Windows
)

func (o Os) String() string {
	switch o {
	case Linux:
		return "linux"
	case Macos:
		return "macos"
	case Windows:
		return "windows"
	default:
		return fmt.Sprintf("Unknown Os %d", o)
	}
}

func (o *Os) MarshalJSON() ([]byte, error) {
	return []byte("\"" + o.String() + "\""), nil
}

func (o *Os) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case "linux":
		*o = Linux
	case "macos":
		*o = Macos
	case "windows":
		*o = Windows
	}
	return nil
}
