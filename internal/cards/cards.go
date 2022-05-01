package cards

import (
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
)

type Card struct {
	Secret   string
	Key      string
	Currency string
}

type Transaction struct {
	TransactionStatusID int
	Amount              int
	Currency            string
	LastFour            string
	BankReturnCode      string
}

func (c *Card) Charge(currency string, amount int) (*stripe.PaymentIntent, string, error) {
	return c.CreatePaymentIntent(currency, amount)
}

func (c *Card) CreatePaymentIntent(currency string, amount int) (*stripe.PaymentIntent, string, error) {
	stripe.Key = c.Secret

	// create a payment intent
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(amount)),
		Currency: stripe.String(currency),
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		msg := ""
		if stripeErr, ok := err.(*stripe.Error); ok {
			msg = cardErrorMessage(stripeErr.Code)
		}
		return nil, msg, err
	}

	return pi, "", nil
}

func cardErrorMessage(code stripe.ErrorCode) string {
	var msg = ""
	switch code {
	case stripe.ErrorCodeCardDeclined:
		msg = "このカードが使えませんでした"
	case stripe.ErrorCodeExpiredCard:
		msg = "カードの有効期間が超えています"
	case stripe.ErrorCodeIncorrectCVC:
		msg = "CVCコードが正しくありません"
	case stripe.ErrorCodeIncorrectZip:
		msg = "郵便番号が正しくありません"
	case stripe.ErrorCodeAmountTooLarge:
		msg = "カードの支払上限金額を超えています"
	case stripe.ErrorCodeAmountTooSmall:
		msg = "指定された金額以下は処理できません"
	case stripe.ErrorCodeBalanceInsufficient:
		msg = "残高が足りませんでした"
	case stripe.ErrorCodePostalCodeInvalid:
		msg = "郵便番号が正しくありません"
	default:
		msg = "このカードが使えませんでした"
	}
	return msg
}
