package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/sebsvt/organisation-service/services"
)

type organisationHandler struct {
	organisation_srv services.OrganisationService
}

func NewOrganisationHandler(organisation_srv services.OrganisationService) organisationHandler {
	return organisationHandler{organisation_srv: organisation_srv}
}

func (h organisationHandler) CreateNewOrganisation(c *fiber.Ctx) error {
	user_id := c.Locals("user_id").(string)

	var new_organisation services.OrgansationCreated
	if err := c.BodyParser(&new_organisation); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	org_id, err := h.organisation_srv.CreateNewOrganisation(user_id, new_organisation)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"organisation_id": org_id,
	})
}

func (h organisationHandler) GetOrganisationFromDomain(c *fiber.Ctx) error {
	domain := c.Params("domain")
	organisation, err := h.organisation_srv.GetOrganisationByDomain(domain)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(organisation)
}

func (h organisationHandler) GetAllOrganisationMember(c *fiber.Ctx) error {
	org_id := c.Params("organisation_id")
	user_id := c.Locals("user_id").(string)

	members, err := h.organisation_srv.GetAllOrganisationMember(org_id, user_id)
	if err != nil {
		return c.JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(members)
}

func (h organisationHandler) GetAllOrganisationsFromUserID(c *fiber.Ctx) error {
	user_id := c.Locals("user_id").(string)
	organisations, err := h.organisation_srv.GetAllOrganisationFromUserID(user_id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(organisations)
}
