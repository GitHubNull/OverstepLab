package vuln

import (
	"sync/atomic"
)

// SecurityMode represents the current security mode of the application.
// 0 = Vulnerable (default), 1 = Secure
type SecurityMode int32

const (
	ModeVulnerable SecurityMode = 0
	ModeSecure     SecurityMode = 1
)

var currentMode atomic.Int32

func init() {
	currentMode.Store(int32(ModeVulnerable))
}

func SetSecureMode(secure bool) {
	if secure {
		currentMode.Store(int32(ModeSecure))
	} else {
		currentMode.Store(int32(ModeVulnerable))
	}
}

func IsSecureMode() bool {
	return currentMode.Load() == int32(ModeSecure)
}

func GetCurrentMode() SecurityMode {
	return SecurityMode(currentMode.Load())
}

func GetModeString() string {
	if IsSecureMode() {
		return "secure"
	}
	return "vulnerable"
}
