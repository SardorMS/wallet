package wallet

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/SardorMS/wallet/pkg/types"
	"github.com/google/uuid"
)

//Error variables.
var (
	ErrPhoneRegistered      = errors.New("phnone number already registered")
	ErrAmountMustBePositive = errors.New("amount must be greater than zero")
	ErrAccountNotFound      = errors.New("account not found")
	ErrNotEnoughBalance     = errors.New("not enough balance")
	ErrPaymentNotFound      = errors.New("payment not found")
	ErrFavoriteNotFound     = errors.New("favorite not found")
)

//Service - service struct.
type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
	favorites     []*types.Favorite
}

//RegisterAccount - authentication processes method performing.
func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistered
		}
	}

	s.nextAccountID++
	account := &types.Account{
		ID:      s.nextAccountID,
		Phone:   phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts, account)

	return account, nil
}

//FindAccountByID - method that find account by ID.
func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.ID == accountID {
			return account, nil
		}
	}
	return nil, ErrAccountNotFound
}

// Deposit -  replenish the user's account.
func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount <= 0 {
		return ErrAmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}

	if account == nil {
		return ErrAccountNotFound
	}

	account.Balance += amount
	return nil
}

//Pay - payments method.
func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}
	if account == nil {
		return nil, ErrAccountNotFound
	}

	if account.Balance < amount {
		return nil, ErrNotEnoughBalance
	}

	account.Balance -= amount
	paymentID := uuid.New().String()
	payment := &types.Payment{
		ID:        paymentID,
		AccountID: accountID,
		Amount:    amount,
		Category:  category,
		Status:    types.PaymentStatusInProgress,
	}
	s.payments = append(s.payments, payment)
	return payment, nil
}

//FindPaymentByID - method that find payment by ID.
func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	for _, payment := range s.payments {
		if payment.ID == paymentID {
			return payment, nil
		}
	}
	return nil, ErrPaymentNotFound
}

//Reject - method that returns payment in a accident of error.
func (s *Service) Reject(paymentID string) error {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return err
	}
	account, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		return err
	}

	payment.Status = types.PaymentStatusFail
	account.Balance += payment.Amount
	return nil
}

//Repeat - repeats payment.
func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	return s.Pay(payment.AccountID, payment.Amount, payment.Category)
}

//FavoritePayment - makes a favorite from a specific payment.
func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	favoriteID := uuid.New().String()
	favorite := &types.Favorite{
		ID:        favoriteID,
		AccountID: payment.AccountID,
		Amount:    payment.Amount,
		Name:      name,
		Category:  payment.Category,
	}

	s.favorites = append(s.favorites, favorite)
	return favorite, nil
}

//FindFavoriteByID - method that find favorite payment by ID.
func (s *Service) FindFavoriteByID(favoriteID string) (*types.Favorite, error) {
	for _, favorite := range s.favorites {
		if favorite.ID == favoriteID {
			return favorite, nil
		}
	}

	return nil, ErrFavoriteNotFound
}

//PayFromFavorites - makes a payment from a specific favorite one.
func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
	favorite, err := s.FindFavoriteByID(favoriteID)
	if err != nil {
		return nil, err
	}

	return s.Pay(favorite.AccountID, favorite.Amount, favorite.Category)
}

//ExportToFile - writes accounts to a file.
func (s *Service) ExportToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		log.Print(err)
		return err
	}

	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Print(err)
		}
	}()

	data := make([]byte, 0)
	lastStr := ""
	for _, account := range s.accounts {
		text := []byte(
			strconv.FormatInt(int64(account.ID), 10) + string(";") +
				string(account.Phone) + string(";") +
				strconv.FormatInt(int64(account.Balance), 10) + string("|"))

		data = append(data, text...)
		str := string(data)
		lastStr = strings.TrimSuffix(str, "|")
	}

	_, err = file.Write([]byte(lastStr))
	if err != nil {
		log.Print(err)
		return err
	}
	log.Printf("%#v", file)
	return nil
}

//ImportFromFile - import(reads) from file to accounts.
func (s *Service) ImportFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		log.Print(err)
		return err
	}

	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()

	content := make([]byte, 0)
	buf := make([]byte, 4)
	for {
		read, err := file.Read(buf)
		if err == io.EOF {
			content = append(content, buf[:read]...)
			break
		}

		if err != nil {		
			log.Print(err)
			return err
		}
		content = append(content, buf[:read]...)
	}

	data := string(content)
	log.Println("data: ", data)

	acc := strings.Split(data, "|")
	log.Println("acc: ", acc)

	// account := strings.TrimSuffix(data, "|")
	// log.Println("account: ", account)

	// acc := make([]string, len(account))
	// log.Println("accounts until split:", accounts)

	// accounts = accounts[:len(accounts) -1]
	// log.Println("accounts after split: ", accounts)

	for _, operation := range acc {

		strAcc := strings.Split(operation, ";")
		log.Println("strAcc:", strAcc)

		id, err := strconv.ParseInt(strAcc[0], 10, 64)
		if err != nil {
			log.Print(err)
			return err
		}

		phone := types.Phone(strAcc[1])

		balance, err := strconv.ParseInt(strAcc[2], 10, 64)
		if err != nil {
			log.Print(err)
			return err
		}

		account := &types.Account{
			ID:      id,
			Phone:   phone,
			Balance: types.Money(balance),
		}

		s.accounts = append(s.accounts, account)
		log.Print(account)
	}
	return nil
}


//Export - ...
func (s *Service) Export(dir string) error {

	if s.accounts != nil && len(s.accounts) > 0 {

		accDir, err := filepath.Abs(dir)
		if err != nil {
			log.Print(err)
			return err
		}

		data := make([]byte, 0)
		for _, account := range s.accounts {
			text := []byte(
				strconv.FormatInt(int64(account.ID), 10) + ";" +
					string(account.Phone) + ";" +
					strconv.FormatInt(int64(account.Balance), 10) + "\n")

			data = append(data, text...)
		}

		err = os.WriteFile(accDir+"/"+"accounts.dump", data, 0666)
		if err != nil {
			log.Print(err)
			return err
		}
	}

	if s.payments != nil && len(s.payments) > 0 {

		payDir, err := filepath.Abs(dir)
		if err != nil {
			log.Print(err)
			return err
		}

		data := make([]byte, 0)
		for _, payment := range s.payments {
			text := []byte(
				string(payment.ID) + ";" +
					strconv.FormatInt(int64(payment.AccountID), 10) + ";" +
					strconv.FormatInt(int64(payment.Amount), 10) + ";" +
					string(payment.Category) + ";" +
					string(payment.Status) + "\n")

			data = append(data, text...)
		}

		err = os.WriteFile(payDir+"/"+"payments.dump", data, 0666)
		if err != nil {
			log.Print(err)
			return err
		}
	}

	if s.favorites != nil && len(s.favorites) > 0 {

		favDir, err := filepath.Abs(dir)
		if err != nil {
			log.Print(err)
			return err
		}

		data := make([]byte, 0)
		for _, favorite := range s.favorites {
			text := []byte(
				string(favorite.ID) + ";" +
					strconv.FormatInt(int64(favorite.AccountID), 10) + ";" +
					string(favorite.Name) + ";" +
					strconv.FormatInt(int64(favorite.Amount), 10) + ";" +
					string(favorite.Category) + "\n")

			data = append(data, text...)
		}

		err = os.WriteFile(favDir+"/"+"favorites.dump", data, 0666)
		if err != nil {
			log.Print(err)
			return err
		}
	}

	return nil
}

//Import - ...
func (s * Service) Import(dir string) error  {
	accDir, err := filepath.Abs(dir)
		if err != nil {
			log.Print(err)
			return err
		}

	file, err := os.ReadFile(accDir + "")
	if err != nil {
		log.Print(err)
		return err
	}



	return nil
}