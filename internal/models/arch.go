package models

type Arch int

const (
	UNKNOWN Arch = iota
	X64
	ARM
)

func (a Arch) String() string {
	switch a {
	case X64:
		return "x64"
	case ARM:
		return "arm"
	default:
		return "unknown"
	}
}

func (a Arch) GoArch() string {
	switch a {
	case X64:
		return "amd64"
	case ARM:
		return "arm64"
	default:
		return "unknown"
	}
}

func (a *Arch) MarshalJSON() ([]byte, error) {
	return []byte("\"" + a.String() + "\""), nil
}

func (a *Arch) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case "x64":
		*a = X64
	case "arm":
		*a = ARM
	}
	return nil
}
