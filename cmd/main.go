package main

import (
	"github.com/SardorMS/wallet/pkg/wallet"
)

func main() {
	svc := &wallet.Service{}

	svc.RegisterAccount("+999923123123")
	svc.Deposit(1, 100)

	svc.RegisterAccount("+999942342342")
	svc.Deposit(2, 500)

	svc.RegisterAccount("+999964564566")
	svc.Deposit(3, 400)

	// svc.ExportToFile("../data/export.txt")
	// svc.ImportFromFile("../data/export.txt")
	svc.Export("../data")
}
