# wallet 
Wallet - is a service for registration, deposit, payment, search and data processing in payment systems.

![screenshot](./img/wallet_logo.png)

## Install
Install last version - v1.9.2.

```
go get github.com/SardorMS/wallet
```

## Methods
Below are a few methods you could use:

```go
// RegisterAccount - authentication processes method performing.
func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
  ...}

// Deposit -  replenish the user's account.
func (s *Service) Deposit(accountID int64, amount types.Money) error {
  ...}

// Pay - payments method.
func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
  ...}

// FindAccountByID - method that find account by ID.
func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
  ...}

// FindPaymentByID - method that find payment by ID.
func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
  ...}

// Reject - method that returns payment in a accident of error.
func (s *Service) Reject(paymentID string) error {
  ...}

// Repeat - repeats payment.
func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
  ...}

// FavoritePayment - makes a favorite from a specific payment.
func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
  ...}

// FindFavoriteByID - method that find favorite payment by ID.
func (s *Service) FindFavoriteByID(favoriteID string) (*types.Favorite, error) {
  ...}

// PayFromFavorites - makes a payment from a specific favorite one.
func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
  ...}

// Export - writes accounts, payments, favorites to a dump file(full_version).
func (s *Service) Export(dir string) error {
  ...}

// Import - import(reads) from dump file to accounts, payments and favorites(full_version).
func (s *Service) Import(dir string) error {
  ...}

// ExportAccountHistory - pulls out payments of a specific account.
func (s *Service) ExportAccountHistory(accountID int64) ([]types.Payment, error) {
...}

// HistoryToFiles - save all data(information about the payments) to files.
func (s *Service) HistoryToFiles(payments []types.Payment, dir string, records int) error {
  ...}

// SumPayments - summarizes payments using goroutines.
func (s *Service) SumPayments(goroutines int) types.Money {
  ...}

// FilterPayments - filters out payments by accountID using goroutines.
func (s *Service) FilterPayments(accountID int64, goroutines int) ([]types.Payment, error) {
  ...}

// FilterPaymentsByFn - filters out payments by any function.
func (s *Service) FilterPaymentsByFn(
	filter func(payment types.Payment) bool, goroutines int) ([]types.Payment, error) {
    ...}

// SumPaymentsWithProgress - summarizes payments by adding them to the channel.
func (s *Service) SumPaymentsWithProgress() <-chan Progress {
  ...}

// Merge - creates a channel, in which messages from all channels appear,
// which are sent in a slice.
func Merge(channels []<-chan Progress) <-chan Progress {
  ...}

```
 
## Usage

1. Test metricks: 
```sh
$ go test -v ./...
?       github.com/SardorMS/wallet/cmd  [no test files]
?       github.com/SardorMS/wallet/pkg/types    [no test files]
=== RUN   TestService_RegisterAccount
--- PASS: TestService_RegisterAccount (0.00s)
=== RUN   TestService_Deposit
--- PASS: TestService_Deposit (0.00s)
=== RUN   TestService_Pay
--- PASS: TestService_Pay (0.00s)
=== RUN   TestService_FindAccountByID_success
--- PASS: TestService_FindAccountByID_success (0.00s)
=== RUN   TestService_FindAccountByID_notFound
--- PASS: TestService_FindAccountByID_notFound (0.00s)
=== RUN   TestService_FindPaymentByID_success
--- PASS: TestService_FindPaymentByID_success (0.00s)
=== RUN   TestService_FindPaymentByID_notFound
--- PASS: TestService_FindPaymentByID_notFound (0.00s)
=== RUN   TestService_Reject_success
--- PASS: TestService_Reject_success (0.00s)
=== RUN   TestService_Reject_notFound1
--- PASS: TestService_Reject_notFound1 (0.00s)
=== RUN   TestService_Repeat_success
--- PASS: TestService_Repeat_success (0.00s)
=== RUN   TestService_Repeat_notFound
--- PASS: TestService_Repeat_notFound (0.00s)
=== RUN   TestService_FavoritePayment_success
--- PASS: TestService_FavoritePayment_success (0.00s)
=== RUN   TestService_FavoritePayment_notFound
--- PASS: TestService_FavoritePayment_notFound (0.00s)
=== RUN   TestService_FindFavoriteByID_success
--- PASS: TestService_FindFavoriteByID_success (0.00s)
=== RUN   TestService_FindFavoriteByID_notFound
--- PASS: TestService_FindFavoriteByID_notFound (0.00s)
=== RUN   TestService_PayFromFavorite_success
--- PASS: TestService_PayFromFavorite_success (0.00s)
=== RUN   TestService_PayFromFavorite_notFound
--- PASS: TestService_PayFromFavorite_notFound (0.00s)
=== RUN   TestService_ExportToFile_success
2021/06/29 16:45:33 &os.File{file:(*os.file)(0xc0000cac80)}
--- PASS: TestService_ExportToFile_success (0.05s)
=== RUN   TestService_ExportToFile_notFound
2021/06/29 16:45:33 open : The system cannot find the file specified.
--- PASS: TestService_ExportToFile_notFound (0.00s)
=== RUN   TestService_ImportFromFile_success
2021/06/29 16:45:33 data:  1;+998204567898;400
2021/06/29 16:45:33 acc:  [1;+998204567898;400]
2021/06/29 16:45:33 strAcc: [1 +998204567898 400]
2021/06/29 16:45:33 &{1 +998204567898 400}
--- PASS: TestService_ImportFromFile_success (0.00s)
=== RUN   TestService_ImportFromFile_notFound
2021/06/29 16:45:33 open : The system cannot find the file specified.
--- PASS: TestService_ImportFromFile_notFound (0.00s)
=== RUN   TestService_Export
--- PASS: TestService_Export (0.00s)
=== RUN   TestService_Import_success
2021/06/29 16:45:33 accounts : [1;+998204567898;400]
2021/06/29 16:45:33 accStr: [1 +998204567898 400]
2021/06/29 16:45:33 &{1 +998204567898 400}
2021/06/29 16:45:33 paySlice : [bc28cb95-0914-4fc4-baa7-3ddbd831bccb;1;100;auto;INPROGRESS]
2021/06/29 16:45:33 payStr: [bc28cb95-0914-4fc4-baa7-3ddbd831bccb 1 100 auto INPROGRESS]
2021/06/29 16:45:33 &{bc28cb95-0914-4fc4-baa7-3ddbd831bccb 1 100 auto INPROGRESS}
2021/06/29 16:45:33 favSlice : [25baf3da-d054-4de2-87f1-7c5f54a6bfcf;1;Favorite payment #i;100;auto d42dd726-c0b7-4bcf-9903-43a7d29180c5;1;my favorite payment;100;auto]
2021/06/29 16:45:33 favStr: [25baf3da-d054-4de2-87f1-7c5f54a6bfcf 1 Favorite payment #i 100 auto]
2021/06/29 16:45:33 &{25baf3da-d054-4de2-87f1-7c5f54a6bfcf 1 Favorite payment #i 100 auto}
2021/06/29 16:45:33 favStr: [d42dd726-c0b7-4bcf-9903-43a7d29180c5 1 my favorite payment 100 auto]
2021/06/29 16:45:33 &{d42dd726-c0b7-4bcf-9903-43a7d29180c5 1 my favorite payment 100 auto}
--- PASS: TestService_Import_success (0.00s)
=== RUN   TestService_Import_notFound1
2021/06/29 16:45:33 open /accounts.dump: The system cannot find the file specified.
2021/06/29 16:45:33 open /payments.dump: The system cannot find the file specified.
2021/06/29 16:45:33 open /favorites.dump: The system cannot find the file specified.
--- PASS: TestService_Import_notFound1 (0.00s)
=== RUN   TestService_Import_notFound2
2021/06/29 16:45:33 open /accounts.dump: The system cannot find the file specified.
2021/06/29 16:45:33 open /payments.dump: The system cannot find the file specified.
2021/06/29 16:45:33 open /favorites.dump: The system cannot find the file specified.
--- PASS: TestService_Import_notFound2 (0.00s)
=== RUN   TestService_Import_Error
2021/06/29 16:45:33 accounts : [1;+998204567898;400]
2021/06/29 16:45:33 accStr: [1 +998204567898 400]
2021/06/29 16:45:33 paySlice : [bc28cb95-0914-4fc4-baa7-3ddbd831bccb;1;100;auto;INPROGRESS]
2021/06/29 16:45:33 payStr: [bc28cb95-0914-4fc4-baa7-3ddbd831bccb 1 100 auto INPROGRESS]
2021/06/29 16:45:33 &{bc28cb95-0914-4fc4-baa7-3ddbd831bccb 1 100 auto INPROGRESS}
2021/06/29 16:45:33 favSlice : [25baf3da-d054-4de2-87f1-7c5f54a6bfcf;1;Favorite payment #i;100;auto d42dd726-c0b7-4bcf-9903-43a7d29180c5;1;my favorite payment;100;auto]
2021/06/29 16:45:33 favStr: [25baf3da-d054-4de2-87f1-7c5f54a6bfcf 1 Favorite payment #i 100 auto]
2021/06/29 16:45:33 &{25baf3da-d054-4de2-87f1-7c5f54a6bfcf 1 Favorite payment #i 100 auto}
2021/06/29 16:45:33 favStr: [d42dd726-c0b7-4bcf-9903-43a7d29180c5 1 my favorite payment 100 auto]
2021/06/29 16:45:33 &{d42dd726-c0b7-4bcf-9903-43a7d29180c5 1 my favorite payment 100 auto}
--- PASS: TestService_Import_Error (0.00s)
=== RUN   TestService_Import_emptyFiles
2021/06/29 16:45:33 accounts : []
2021/06/29 16:45:33 paySlice : []
2021/06/29 16:45:33 favSlice : []
--- PASS: TestService_Import_emptyFiles (0.00s)
=== RUN   TestService_ExportAccountHistory_success
--- PASS: TestService_ExportAccountHistory_success (0.00s)
=== RUN   TestService_ExportAccountHistory_notSuccess1
--- PASS: TestService_ExportAccountHistory_notSuccess1 (0.00s)
=== RUN   TestService_ExportAccountHistory_notSuccess2
--- PASS: TestService_ExportAccountHistory_notSuccess2 (0.00s)
=== RUN   TestService_HistoryToFiles_success1
--- PASS: TestService_HistoryToFiles_success1 (0.00s)
=== RUN   TestService_HistoryToFiles_success2
--- PASS: TestService_HistoryToFiles_success2 (0.00s)
=== RUN   TestService_HistoryToFiles_notSuccess1
--- PASS: TestService_HistoryToFiles_notSuccess1 (0.00s)
=== RUN   TestService_HistoryToFiles_notSuccess2
--- PASS: TestService_HistoryToFiles_notSuccess2 (0.00s)
=== RUN   TestService_SumPayments
--- PASS: TestService_SumPayments (0.00s)
=== RUN   TestService_FilterPayments_success
--- PASS: TestService_FilterPayments_success (0.00s)
=== RUN   TestService_FilterPayments_not_Success
--- PASS: TestService_FilterPayments_not_Success (0.00s)
=== RUN   TestService_FilterPaymentsByFn
--- PASS: TestService_FilterPaymentsByFn (0.00s)
=== RUN   TestService_SumPaymentsWithProgress_success1
--- PASS: TestService_SumPaymentsWithProgress_success1 (0.36s)
=== RUN   TestService_SumPaymentsWithProgress_success2
--- PASS: TestService_SumPaymentsWithProgress_success2 (0.00s)
PASS
ok      github.com/SardorMS/wallet/pkg/wallet   0.470s
```
2. Code coverage metrics:
```sh
$ go test -cover ./...
?       github.com/SardorMS/wallet/cmd  [no test files]
?       github.com/SardorMS/wallet/pkg/types    [no test files]
ok      github.com/SardorMS/wallet/pkg/wallet   0.459s  coverage: 93.1% of statements
```
