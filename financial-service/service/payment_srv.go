package service

import (
	"fmt"
	"log"

	stripe "github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/paymentintent"
)

type paymentService struct{}

func NewPaymentService() PaymentService {
	return &paymentService{}
}

// CreatePaymentIntent implements PaymentService.
func (srv *paymentService) CreatePaymentIntent(amount int64, currency string, stripeAccountID string) (string, error) {
	stripe.Key = "testingkeywaithingforrealkeyfromenv"
	params := &stripe.PaymentIntentParams{
		Amount:               stripe.Int64(amount),
		Currency:             stripe.String(currency),
		ApplicationFeeAmount: stripe.Int64(int64(float64(amount) * 0.04)),
	}
	if stripeAccountID != "" {
		params.StripeAccount = stripe.String(stripeAccountID)
	}
	intent, err := paymentintent.New(params)
	if err != nil {
		log.Printf("Failed to create payment intent: %v", err)
		return "", err
	}
	fmt.Println(intent)
	return intent.ClientSecret, nil
}
