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

func (s *Service) addAccount(data testAccount) (*types.Account, []*types.Payment, error) {

	// if balance < 1 {
	// 	return nil, nil, errors.New("Please give balance with positive number")
	// }

	// register user
	account, err := s.RegisterAccount(data.phone)
	if err != nil {
		return nil, nil, fmt.Errorf("can't register account, error = %v", err)
	}

	//  account top up
	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, fmt.Errorf("can't deposity account, error = %v", err)
	}

	// make a payment to account
	payments := make([]*types.Payment, len(data.payments))
	for i, payment := range data.payments {
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, fmt.Errorf("can't make paymnet, error = %v", err)
		}
	}
	return account, payments, nil
}

func TestService_FindAccountByID_success(t *testing.T) {
	s := newTestService()
	account, _, err := s.addAccount(defaultTestAccount)
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
	_, _, err := s.addAccount(defaultTestAccount)
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
	_, payments, err := s.addAccount(defaultTestAccount)
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
	_, _, err := s.addAccount(defaultTestAccount)
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
	_, payments, err := s.addAccount(defaultTestAccount)
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

func TestService_Reject_notFound(t *testing.T) {
	s := newTestService()

	_, payments, err := s.addAccount(defaultTestAccount)
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
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	// trying to cancel the payment
	payment := payments[0]
	// err = s.Reject(payment.ID)
	// if err != nil {
	// 	t.Errorf("Reject(): can't reject payment error=%v", err)
	// 	return
	// }

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
	_, _, err := s.addAccount(defaultTestAccount)
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
