package main

import (
	_ "embed"

	"github.com/Wikiwoop/wikiwoop-chain/command/root"
	"github.com/Wikiwoop/wikiwoop-chain/licenses"
)

var (
	//go:embed LICENSE
	license string
)

func main() {
	licenses.SetLicense(license)

	root.NewRootCommand().Execute()
}
