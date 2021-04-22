package main

import (
	// "github.com/SardorMS/wallet/pkg/types"
	"log"

	"github.com/SardorMS/wallet/pkg/wallet"
)

func main() {
	svc := &wallet.Service{}

	svc.RegisterAccount("+999923123123")
	svc.Deposit(1, 400)
	pay, err := svc.Pay(1, 50, "phone")
	if err != nil {
		log.Print(err)
		return
	}
	svc.FavoritePayment(pay.ID, "my_phone")

	svc.RegisterAccount("+999942342342")
	svc.Deposit(2, 700)
	pay1, err := svc.Pay(2, 300, "auto")
	if err != nil {
		log.Print(err)
		return
	}
	svc.FavoritePayment(pay1.ID, "my_auto")

	svc.RegisterAccount("+999967865765")
	svc.Deposit(3, 600)
	pay2, err := svc.Pay(3, 200, "shop")
	if err != nil {
		log.Print(err)
		return
	}
	svc.FavoritePayment(pay2.ID, "my_shop")

	svc.Export("../data")

	// svc.ExportToFile("../data/export.txt")
	// svc.ImportFromFile("../data/export.txt")
}
