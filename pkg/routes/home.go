package routes

import (
	"fmt"

	"github.com/francoganga/pagoda_bun/models"
	"github.com/francoganga/pagoda_bun/pkg/controller"

	"github.com/labstack/echo/v4"
)

type (
	home struct {
		controller.Controller
	}

	post struct {
		Title string
		Body  string
	}
)

func (c *home) Get(ctx echo.Context) error {
	page := controller.NewPage(ctx)
	page.Layout = "main"
	page.Name = "home"
	page.Metatags.Description = "Welcome to the homepage."
	page.Metatags.Keywords = []string{"Go", "MVC", "Web", "Software"}
	page.Pager = controller.NewPager(ctx, 10)

	transactions := make([]*models.Transaction, 0)

	err := c.Container.Bun.NewSelect().
		Model(&transactions).
		Scan(ctx.Request().Context())

	if err != nil {
		return err
	}


	return c.RenderPage(ctx, page)
}

func (c *home) test(ctx echo.Context) error {

    // transactions := make([]*models.Transaction, 0)

    type Month struct {
        Month string
        Year string
        My string
    }

    var months []Month;

    page := controller.NewPage(ctx)
    page.Layout = "main"
    page.Name = "months"

    err := c.Container.Bun.NewRaw("SELECT DISTINCT EXTRACT(month from date) as month, EXTRACT(year from date) as year, CONCAT(EXTRACT(month from date), ', ', EXTRACT(year from date)) as MY from transactions ORDER BY year ASC, month ASC").Scan(ctx.Request().Context(), &months)

	// err := c.Container.Bun.NewSelect().
	// 	Model(&transactions).
	//         Order("date ASC").
	// 	Scan(ctx.Request().Context())




	if err != nil {
		return err
	}

    page.Data = months

    return c.RenderPage(ctx, page)
}

func (c *home) month(ctx echo.Context) error {

    page := controller.NewPage(ctx)
    page.Layout = "main"
    page.Name = "month"

    year := ctx.QueryParam("year")
    month := ctx.QueryParam("month")

    transactions := make([]*models.Transaction, 0)

    c.Container.Bun.NewRaw("SELECT * FROM transactions WHERE EXTRACT(year from date) = ? AND EXTRACT(month from date) = ?", year, month).Scan(ctx.Request().Context(), &transactions)

    page.Data = transactions

    return c.RenderPage(ctx, page)
}

// fetchPosts is an mock example of fetching posts to illustrate how paging works
func (c *home) fetchPosts(pager *controller.Pager) []post {
	pager.SetItems(20)
	posts := make([]post, 20)

	for k := range posts {
		posts[k] = post{
			Title: fmt.Sprintf("Post example #%d", k+1),
			Body:  fmt.Sprintf("Lorem ipsum example #%d ddolor sit amet, consectetur adipiscing elit. Nam elementum vulputate tristique.", k+1),
		}
	}
	return posts[pager.GetOffset() : pager.GetOffset()+pager.ItemsPerPage]
}
