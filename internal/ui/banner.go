package ui

import (
	"fmt"

	"github.com/Sakura-501/XSStrike-go/internal/version"
)

func Banner() string {
	return fmt.Sprintf("\n\t%s %s\n", version.AppName, version.Version)
}
