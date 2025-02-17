package service

type PaymentService interface {
	CreatePaymentIntent(amount int64, currency, stripeAccountID string) (string, error)
}
