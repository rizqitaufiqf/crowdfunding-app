package payment

import (
	"crowdfunding/user"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"os"
)

type service struct {
}

type Service interface {
	GetPaymentURL(transaction Transaction, user user.User) (string, error)
}

func NewService() *service {
	return &service{}
}

func (s *service) GetPaymentURL(transaction Transaction, user user.User) (string, error) {
	var sn = snap.Client{}
	sn.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

	req := &snap.Request{
		CustomerDetail: &midtrans.CustomerDetails{
			Email: user.Email,
			FName: user.Name,
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  transaction.ID,
			GrossAmt: int64(transaction.Amount),
		},
	}

	snapResp, err := sn.CreateTransaction(req)
	if err != nil {
		return "", err
	}

	return snapResp.RedirectURL, nil
}
