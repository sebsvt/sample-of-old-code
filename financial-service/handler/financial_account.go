package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sebsvt/financial-service/service"
)

type financialAccountHandler struct {
	fnc_acc_srv service.FinancialAccountService
}

func NewFinancialAccountHandler(fnc_acc_srv service.FinancialAccountService) financialAccountHandler {
	return financialAccountHandler{fnc_acc_srv: fnc_acc_srv}
}

func (h financialAccountHandler) GetFinancialFromOrganisation(c *fiber.Ctx) error {
	organisation_id := c.Query("organisation_id")
	if organisation_id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": fiber.ErrBadRequest.Message,
		})
	}
	res, err := h.fnc_acc_srv.GetFinancialFromOrganisation(organisation_id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(res)
}

func (h financialAccountHandler) SetupAccountHandler(c *fiber.Ctx) error {
	organisationID := c.Params("organisation_id") // Assuming you pass the org ID as a URL parameter

	stripeAccountID, err := h.fnc_acc_srv.SetUpFinancialAccount(organisationID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"url": stripeAccountID,
	})
}
