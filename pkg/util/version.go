package util

var Version string
var ShortHash string

func GetVersion() string {
	if Version == "" {
		return "dev"
	} else {
		return Version
	}
}

func GetShortHash() string {
	if ShortHash == "" {
		return "none"
	} else {
		return ShortHash
	}
}
