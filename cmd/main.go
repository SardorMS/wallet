package main

import (
	"github.com/SardorMS/wallet/pkg/wallet"
)

func main() {
	svc := &wallet.Service{}
	svc.RegisterAccount("+999923123123")
	svc.RegisterAccount("+999942342342")
	svc.RegisterAccount("+999964564566")
	svc.ExportToFile("../data/export.txt")
	// svc.ImportToFile("../data/export.txt")
}
