package util

var (
	Version      string
	ShortHash    string
	Architecture string
)

func GetVersion() string {
	return stringOrDefault(Version, "dev")
}

func GetShortHash() string {
	return stringOrDefault(ShortHash, "none")
}

func GetArchitecture() string {
	return stringOrDefault(Architecture, "unspecified")
}

func stringOrDefault(s, d string) string {
	if s == "" {
		return d
	} else {
		return s
	}
}
