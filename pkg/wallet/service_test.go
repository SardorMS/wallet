package wallet

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/SardorMS/wallet/pkg/types"
	"github.com/google/uuid"
)

type testService struct {
	*Service
}

func newTestService() *testService {
	return &testService{Service: &Service{}}
}

type testAccount struct {
	phone    types.Phone
	balance  types.Money
	payments []struct {
		amount   types.Money
		category types.PaymentCategory
	}
}

var defaultTestAccount = testAccount{
	phone:   "+998204567898",
	balance: 500,
	payments: []struct {
		amount   types.Money
		category types.PaymentCategory
	}{
		{amount: 100, category: "auto"},
	},
}

func (s *Service) addAccount(data testAccount) (*types.Account, []*types.Payment, []*types.Favorite, error) {

	// register user
	account, err := s.RegisterAccount(data.phone)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("can't register account, error = %v", err)
	}

	//  account top up
	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("can't deposity account, error = %v", err)
	}

	// make a payment to account
	payments := make([]*types.Payment, len(data.payments))
	favorites := make([]*types.Favorite, len(data.payments))

	for i, payment := range data.payments {

		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("can't make paymnet, error = %v", err)
		}

		favorites[i], err = s.FavoritePayment(payments[i].ID, "Favorite payment #i")
		if err != nil {
			return nil, nil, nil, fmt.Errorf("can't make favorite paymnet, error = %v", err)
		}
	}
	return account, payments, favorites, nil
}

func TestService_RegisterAccount(t *testing.T) {
	s := newTestService()

	s.RegisterAccount("+1111")
	_, err := s.RegisterAccount("+1111")
	if err == nil {
		t.Error(err)
	}
}

func TestService_Deposit(t *testing.T) {
	s := newTestService()
	s.RegisterAccount("+1111")
	err := s.Deposit(0, 400)
	if err == nil {
		t.Error(err)
	}
}

func TestService_Pay(t *testing.T) {
	s := newTestService()
	s.RegisterAccount("+1111")
	s.Deposit(1, 40)
	_, err := s.Pay(1, 100, "phone")
	if err == nil {
		t.Error(err)
	}

	s.RegisterAccount("+2222")
	s.Deposit(2, 40)
	_, err = s.Pay(0, 100, "auto")
	if err == nil {
		t.Error(err)
	}
}

func TestService_FindAccountByID_success(t *testing.T) {
	s := newTestService()
	account, _, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	result, err := s.FindAccountByID(account.ID)
	if err != nil {
		t.Errorf("FindAccountByID(): error: %v ", err)
	}

	if !reflect.DeepEqual(result, account) {
		t.Errorf("FindPaymentByID(): wrong account returned = %v", err)
		return
	}

}

func TestService_FindAccountByID_notFound(t *testing.T) {
	s := newTestService()
	_, _, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	// trying to find a non-existent ID
	anotherID := s.nextAccountID + 1
	_, err = s.FindAccountByID(anotherID)
	if err == nil {
		t.Error("FindAccountByID(): must return error, returned nil")
	}

	if err != ErrAccountNotFound {
		t.Errorf("FindAccountByID(): must return ErrAccountNotFound, returned = %v", err)
		return
	}
}

func TestService_FindPaymentByID_success(t *testing.T) {
	// create a service
	s := newTestService()
	_, payments, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	// trying to find the payment
	payment := payments[0]
	result, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("FindPaymentByID(): error: %v", err)
		return
	}

	// compare the payments
	if !reflect.DeepEqual(result, payment) {
		t.Errorf("FindPaymentByID(): wrong payment returned = %v", err)
		return
	}
}

func TestService_FindPaymentByID_notFound(t *testing.T) {

	s := newTestService()
	_, _, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}
	// trying to find a non-existent payment
	paymentID := uuid.New().String()
	_, err = s.FindPaymentByID(paymentID)
	if err == nil {
		t.Error("FindPaymentByID(): must return error, returned nil")
		return
	}

	if err != ErrPaymentNotFound {
		t.Errorf("FindPaymentByID(): must return ErrPaymentNotFound, returned: %v", err)
		return
	}
}

func TestService_Reject_success(t *testing.T) {
	s := newTestService()
	_, payments, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	// trying to cancel the payment
	payment := payments[0]
	err = s.Reject(payment.ID)
	if err != nil {
		t.Errorf("Reject(): error = %v", err)
		return
	}

	savedPayment, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("Reject(): can't find payment by id, error = %v", err)
		return
	}

	if savedPayment.Status != types.PaymentStatusFail {
		t.Errorf("Reject(): status didn't change, payment = %v", savedPayment)
		return
	}

	savedAccount, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		t.Errorf("Reject(): can't find account by id, error = %v", savedAccount)
		return
	}
	if savedAccount.Balance != defaultTestAccount.balance {
		t.Errorf("Reject(): balance didn't change, payment = %v", savedPayment)
		return
	}

}

func TestService_Reject_notFound1(t *testing.T) {
	s := newTestService()
	
	_, payments, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	// trying to cancel the payment
	payment := payments[0]
	changedPayment := payment.ID + "2"
	err = s.Reject(changedPayment)
	if err == nil {
		t.Errorf("Reject(): error = %v", err)
		return
	}

	if err != ErrPaymentNotFound {
		t.Errorf("Reject(): must return ErrPaymentNotFound, returned: %v", err)
		return
	}
}

func TestService_Repeat_success(t *testing.T) {
	s := newTestService()
	_, payments, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	// trying to cancel the payment
	payment := payments[0]

	// trying to repeat the payment
	newPayment, err := s.Repeat(payment.ID)
	if err != nil {
		t.Errorf("Repeat(): error = %v", err)
		return
	}

	if newPayment.AccountID != payment.AccountID {
		t.Errorf("Repeat(): account ID's difference,\n Repeated payment = %v,\n Rejected payment = %v", newPayment, payment)
		return
	}

	if newPayment.Amount != payment.Amount {
		t.Errorf("Repeat(): amount of payments difference,\n Repeated payment = %v,\n Rejected payment = %v", newPayment, payment)
		return
	}

	if newPayment.Category != payment.Category {
		t.Errorf("Repeat(): category of payments difference,\n Repeated payment = %v,\n Rejected payment = %v", newPayment, payment)
		return
	}

	if newPayment.Status != payment.Status {
		t.Errorf("Repeat(): status of payments difference,\n Repeated payment = %v,\n Rejected payment = %v", newPayment, payment)
		return
	}
}

func TestService_Repeat_notFound(t *testing.T) {
	s := newTestService()
	_, _, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}
	payment := uuid.New().String()
	_, err = s.Repeat(payment)
	if err == nil {
		t.Errorf("Repeat(): must return error, returned nil")
		return
	}
	if err != ErrPaymentNotFound {
		t.Errorf("Repeat(): must return ErrPaymentNotFound, returned: %v", err)
		return
	}

}

func TestService_FavoritePayment_success(t *testing.T) {
	s := newTestService()
	_, payments, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	//create a favorite payment
	payment := payments[0]
	_, err = s.FavoritePayment(payment.ID, "my favorite payment")
	if err != nil {
		t.Errorf("FavoritePayment(): error: %v", err)
		return
	}
}

func TestService_FavoritePayment_notFound(t *testing.T) {
	s := newTestService()
	_, _, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	favorite := uuid.New().String()
	_, err = s.FavoritePayment(favorite, "my favorite payment")
	if err == nil {
		t.Error("FavoritePayment(): must return error, returned nil")
		return
	}

	if err != ErrPaymentNotFound {
		t.Errorf("FavoritePayment(): must return ErrPaymentNotFound, returned: %v", err)
		return
	}
}

func TestService_FindFavoriteByID_success(t *testing.T) {
	s := newTestService()
	_, _, favorites, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	// trying to find favorite the payment
	favorite := favorites[0]
	result, err := s.FindFavoriteByID(favorite.ID)
	if err != nil {
		t.Errorf("FindFavoriteByID(): error: %v", err)
		return
	}

	// compare the payments
	if !reflect.DeepEqual(result, favorite) {
		t.Errorf("FindFavoriteByID(): wrong payment returned = %v", err)
		return
	}
}

func TestService_FindFavoriteByID_notFound(t *testing.T) {

	s := newTestService()
	_, _, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}
	// trying to find a non-existent payment
	favoriteID := uuid.New().String()
	_, err = s.FindFavoriteByID(favoriteID)
	if err == nil {
		t.Error("FindFavoriteByID(): must return error, returned nil")
		return
	}

	if err != ErrFavoriteNotFound {
		t.Errorf("FindFavoriteByID(): must return ErrFavoriteNotFound, returned: %v", err)
		return
	}
}

func TestService_PayFromFavorite_success(t *testing.T) {
	s := newTestService()
	_, _, favorites, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	favorite := favorites[0]
	payment, err := s.PayFromFavorite(favorite.ID)
	if err != nil {
		t.Errorf("PayFromFavorite(): error: %v", err)
		return
	}

	if payment.AccountID != favorite.AccountID {
		t.Errorf("PayFromFavorite(): account ID's difference,\n Current payment = %v,\n favorite payment = %v", payment, favorite)
		return
	}

	if payment.Amount != favorite.Amount {
		t.Errorf("PayFromFavorite(): amount of payments difference,\n Current payment = %v,\n favorite payment = %v", payment, favorite)
		return
	}

	if payment.Category != favorite.Category {
		t.Errorf("PayFromFavorite(): category of payments difference,\n Current payment = %v,\n favorite payment = %v", payment, favorite)
		return
	}
}

func TestService_PayFromFavorite_notFound(t *testing.T) {
	s := newTestService()
	_, _, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	favoriteID := uuid.New().String()
	_, err = s.PayFromFavorite(favoriteID)
	if err == nil {
		t.Errorf("PayFromFavorite(): must return error, returned nil")
		return
	}
	if err != ErrFavoriteNotFound {
		t.Errorf("PayFromFavorite(): must return ErrFavoriteNotFound, returned: %v", err)
		return
	}

}

func TestService_ExportToFile_success(t *testing.T) {
	s := newTestService()
	_, _, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	err = s.ExportToFile("hello.txt")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestService_ExportToFile_notFound(t *testing.T) {
	s := newTestService()
	_, _, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	err = s.ExportToFile("")
	if err == nil {
		t.Error(err)
		return
	}
}

func TestService_ImportFromFile_success(t *testing.T) {
	s := newTestService()
	s.RegisterAccount("1111")
	s.Deposit(1, 500)
	pay, _ := s.Pay(1, 100, "phone")
	s.FavoritePayment(pay.ID, "my_phone")

	err := s.ImportFromFile("hello.txt")
	if err != nil {
		t.Error(err)
		return
	}
}
func TestService_ImportFromFile_noSuccess(t *testing.T) {
	s := newTestService()

	err := s.ImportFromFile("")
	if err == nil {
		t.Error(err)
		return
	}
}

func Transactions(s *testService) {
	s.RegisterAccount("1111")
	s.Deposit(1, 500)
	s.Pay(1, 10, "food")
	s.Pay(1, 10, "phone")
	s.Pay(1, 15, "cafe")
	s.Pay(1, 25, "auto")
	s.Pay(1, 30, "restaurant")
	s.Pay(1, 50, "auto")
	s.Pay(1, 60, "bank")
	s.Pay(1, 50, "bank")

	s.RegisterAccount("2222")
	s.Deposit(2, 200)
	s.Pay(2, 40, "phone")

	s.RegisterAccount("3333")
	s.Deposit(3, 300)
	s.Pay(3, 36, "auto")
	s.Pay(3, 12, "food")
	s.Pay(3, 25, "phone")
}

func TestService_ExportAccountHistory_success(t *testing.T) {
	s := newTestService()
	Transactions(s)
	_, err := s.ExportAccountHistory(1)
	if err != nil {
		t.Error(err)
	}
}

func TestService_ExportAccountHistory_notSuccess1(t *testing.T) {
	s := newTestService()
	s.RegisterAccount("")
	s.Deposit(0, 0)
	s.Pay(0, 0, "")
	_, err := s.ExportAccountHistory(1)
	if err == nil {
		t.Error(err)
	}
}
func TestService_ExportAccountHistory_notSuccess2(t *testing.T) {
	s := newTestService()
	_, _, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	anotherID := s.nextAccountID + 1
	_, err = s.FindAccountByID(anotherID)
	if err == nil {
		t.Error("ExportAccountHistory(): must return error, returned nil")
	}

	_, err = s.ExportAccountHistory(3)
	if err == nil {
		t.Error(err)
	}
	if err != ErrAccountNotFound {
		t.Errorf("ExportAccountHistory(): must return ErrAccountNotFound, returned = %v", err)
		return
	}

}

func TestService_HistoryToFiles_success1(t *testing.T) {
	s := newTestService()
	Transactions(s)

	payments, err := s.ExportAccountHistory(1)
	if err != nil {
		t.Error(err)
	}
	err = s.HistoryToFiles(payments, "data", 3)
	if err != nil {
		t.Error(err)
	}
}
func TestService_HistoryToFiles_Success2(t *testing.T) {
	s := newTestService()
	Transactions(s)

	payments, err := s.ExportAccountHistory(1)
	if err != nil {
		t.Error(err)
	}
	err = s.HistoryToFiles(payments, "data", 12)
	if err != nil {
		t.Error(err)
	}
}

func TestService_HistoryToFiles_notSuccess1(t *testing.T) {
	s := newTestService()
	Transactions(s)

	payment := []types.Payment{}
	err := s.HistoryToFiles(payment, "data", 12)
	if err != nil {
		t.Error(err)
	}
}
func TestService_HistoryToFiles_notSuccess2(t *testing.T) {
	s := newTestService()
	Transactions(s)

	payment := []types.Payment{}
	err := s.HistoryToFiles(payment, "", 0)
	if err == nil {
		t.Error(err)
	}
}

func TestService_SumPayments(t *testing.T) {
	s := newTestService()
	Transactions(s)
	sum := s.SumPayments(0)
	if sum != 363 {
		t.Errorf("TestService_SumPayments(): sum=%v", sum)
	}

}

func BenchmarkSumPayments(b *testing.B) {
	s := newTestService()
	Transactions(s)
	want := types.Money(363)
	for i := 0; i < b.N; i++ {
		result := s.SumPayments(3)
		if result != want {
			b.Fatalf("INVALID: result_we_got %v, result_we_want %v", result, want)
		}
	}
}

func TestService_FilterPayments_success(t *testing.T) {
	s := newTestService()
	Transactions(s)

	paymnet, err := s.FilterPayments(1, 0)
	if err != nil {
		t.Error(err)
	}

	want := 8
	result := len(paymnet)
	if !reflect.DeepEqual(result, want) {
		t.Errorf("INVALID: result_we_got %v, result_we_want %v", result, want)
		return
	}
}
func TestService_FilterPayments_not_Success(t *testing.T) {
	s := newTestService()
	Transactions(s)

	_, err := s.FilterPayments(0, 0)
	if err == nil {
		t.Error(err)
	}
}

func BenchmarkFilterPayments(b *testing.B) {
	s := newTestService()
	Transactions(s)

	for i := 0; i < b.N; i++ {
		paymnet, err := s.FilterPayments(1, 3)
		if err != nil {
			b.Error(err)
		}

		want := 8
		result := len(paymnet)
		if !reflect.DeepEqual(result, want) {
			b.Fatalf("INVALID: result_we_got %v, result_we_want %v", result, want)
			return
		}
	}
}


