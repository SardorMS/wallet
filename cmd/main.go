package main

import (
	"github.com/SardorMS/wallet/pkg/wallet"
	"log"
)

func main() {
	svc := &wallet.Service{}

	// svc.RegisterAccount("+1111")
	// svc.Deposit(1, 400)

	// svc.RegisterAccount("+2222")
	// svc.Deposit(2, 500)

	// svc.RegisterAccount("+3333")
	// svc.Deposit(3, 600)

	// svc.RegisterAccount("+4444")
	// svc.Deposit(4, 700)
	// ==================================================
	svc.RegisterAccount("1111")
	svc.Deposit(1, 500)
	pay, err := svc.Pay(1, 100, "phone")
	if err != nil {
		log.Print(err)
		return
	}
	svc.FavoritePayment(pay.ID, "my_phone")
	//================================================
	svc.RegisterAccount("2222")
	svc.Deposit(2, 800)
	pay1, err := svc.Pay(2, 200, "auto")
	if err != nil {
		log.Print(err)
		return
	}
	svc.FavoritePayment(pay1.ID, "my_auto")
	//=================================================
	svc.RegisterAccount("3333")
	svc.Deposit(3, 1000)
	pay2, err := svc.Pay(3, 300, "shop")
	if err != nil {
		log.Print(err)
		return
	}
	svc.FavoritePayment(pay2.ID, "my_shop")
	//==================================================
	// svc.Export("../data")
	// svc.Import("../data")
	// svc.ExportToFile("../data/export.txt")
	// svc.ImportFromFile("../data/export.txt")
}
