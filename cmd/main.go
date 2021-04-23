package main

import (
	"github.com/SardorMS/wallet/pkg/wallet"
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
	pay, _ := svc.Pay(1, 100, "phone")
	t, _ := svc.FavoritePayment(pay.ID, "my_phone")
	svc.PayFromFavorite(t.ID)
	//================================================
	svc.RegisterAccount("2222")
	svc.Deposit(2, 1000)
	pay1, _ := svc.Pay(2, 200, "auto")
	t1, _ := svc.FavoritePayment(pay1.ID, "my_auto")
	svc.PayFromFavorite(t1.ID)
	//=================================================
	svc.RegisterAccount("3333")
	svc.Deposit(3, 12000)
	pay2, _ := svc.Pay(3, 300, "shop")
	t2, _ := svc.FavoritePayment(pay2.ID, "my_shop")
	svc.PayFromFavorite(t2.ID)

	//==================================================
	svc.Export("../data")
	svc.Import("../data")
	// svc.ExportToFile("../data/export.txt")
	// svc.ImportFromFile("../data/export.txt")
}
