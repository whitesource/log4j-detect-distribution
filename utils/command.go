package utils

import (
	"os/exec"
)

func IsInstalled(binary string) bool {
	if _, err := exec.LookPath(binary); err != nil {
		return false
	}

	return true
}
