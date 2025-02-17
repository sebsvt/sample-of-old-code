package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sebsvt/financial-service/service"
)

type paymentHandler struct {
	payment_srv service.PaymentService
}

func NewPaymentHandler(paymentService service.PaymentService) paymentHandler {
	return paymentHandler{payment_srv: paymentService}
}

func (h paymentHandler) CreatePaymentIntent(c *fiber.Ctx) error {
	amountStr := c.Query("amount")
	currency := c.Query("currency")
	stripeAccountID := c.Query("stripe_account_id")

	// Convert the amount to int64
	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid amount"})
	}

	// Create the payment intent
	intent, err := h.payment_srv.CreatePaymentIntent(amount, currency, stripeAccountID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create payment intent"})
	}

	// Respond with the payment intent
	return c.JSON(fiber.Map{
		"client_secret": intent,
	})
}

//
