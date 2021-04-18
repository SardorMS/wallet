package main

import (
	"github.com/SardorMS/wallet/pkg/wallet"
)

func main() {

	s := &wallet.Service{}
	s.RegisterAccount("+9989956645321")
	s.RegisterAccount("+9891244367564")
	s.RegisterAccount("+9949852352353")
	s.ExportToFile("../data/export.txt")
}
