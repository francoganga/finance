package routes

import (
	"fmt"

	"github.com/francoganga/pagoda_bun/models"
	"github.com/francoganga/pagoda_bun/pkg/controller"
	"github.com/labstack/echo/v4"
)

type transaction struct {
	controller.Controller
}

// func (c *transaction) New(ctx echo.Context) error {
//
// 	return ctx.String(200, "hello new transaction")
// }
//

func (c *transaction) Edit(ctx echo.Context) error {
	page := controller.NewPage(ctx)
	page.Layout = "main"
	page.Name = "transaction/edit"

	id := ctx.Param("id")

	var transaction models.Transaction
	var categories []models.Category

	c.Container.Bun.NewSelect().
		Model(&transaction).
		Where("id = ?", id).
		Scan(ctx.Request().Context())

	c.Container.Bun.NewSelect().Model(&categories).Scan(ctx.Request().Context())

	if page.HTMX.Request.Enabled {
		modal := controller.NewPage(ctx)
		modal.Name = "transaction/_modal"

		modal.Data = struct {
			Transaction models.Transaction
			Categories  []models.Category
		}{
			transaction,
			categories,
		}

		return c.RenderPage(ctx, modal)
	}

	return ctx.String(200, "hello edit transaction")
}

func (c *transaction) Index(ctx echo.Context) error {

	page := controller.NewPage(ctx)
	page.Layout = "main"
	page.Name = "month"

	year := ctx.QueryParam("year")
	month := ctx.QueryParam("month")

	fmt.Printf("YEAR=%v\n", year)
	fmt.Printf("MONTH=%v\n", month)

	transactions := make([]*models.Transaction, 0)

	c.Container.Bun.NewRaw("SELECT * FROM transactions WHERE EXTRACT(year from date) = ? AND EXTRACT(month from date) = ?", year, month).Scan(ctx.Request().Context(), &transactions)

	page.Data = transactions

	return c.RenderPage(ctx, page)
}
