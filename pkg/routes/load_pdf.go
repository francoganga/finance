package routes

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	_ "time"

	_ "github.com/francoganga/pagoda_bun/models"
	"github.com/francoganga/pagoda_bun/pkg/controller"
	_ "github.com/francoganga/pagoda_bun/pkg/internal/parser"
	"github.com/labstack/echo/v4"
)

type (
	loadPdf struct {
		controller.Controller
	}
)

func (c *loadPdf) Get(ctx echo.Context) error {

	page := controller.NewPage(ctx)

	page.Layout = "main"
	page.Name = "load-pdf"
	page.Title = "Load PDF"

	return c.RenderPage(ctx, page)
}

func (c *loadPdf) Post(ctx echo.Context) error {

	file, err := ctx.FormFile("pdf")

	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	path := filepath.Join("/tmp", file.Filename)

	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	out, err := exec.Command("pdftotext", "-layout", "-f", "1", "-l", "3", path, "-").Output()

	if err != nil {
		return err
	}

	re, err := regexp.Compile("([0-9]{2}\\/[0-9]{2}\\/[0-9]{2})\\s+([0-9]+)\\s+(.*?)\\s{2,}(.*?)\\s{2,}(.*)\n(.*)")

	if err != nil {
		return err
	}
	matches := re.FindAllString(string(out), -1)

	for _, line := range matches {
		fmt.Println(strings.Replace(line, "\n", "\\n", -1))
	}

	// consumos := make([]*parser.ConsumoDto, 0)

	// transactions := make([]*models.Transaction, 0)

	// for _, line := range matches {
	// 	p := parser.FromInput(line)
	//
	// 	consu, err := p.ParseConsumo()
	//
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	pd, err := time.Parse("02/01/06", consu.Date)
	//
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	t := &models.Transaction{
	// 		Date:        pd,
	// 		Code:        consu.Code,
	// 		Description: consu.Description,
	// 		Amount:      consu.Amount,
	// 		Balance:     consu.Balance,
	// 	}
	//
	// 	transactions = append(transactions, t)
	//
	// 	consumos = append(consumos, consu)
	// }

	res := strings.Join(matches, "\n")

	// _, err = c.Container.Bun.NewInsert().
	// 	Model(&transactions).
	// 	Exec(ctx.Request().Context())
	//
	// if err != nil {
	// 	return err
	// }

	return ctx.HTML(http.StatusOK, fmt.Sprintf("<html><body><pre>%s</pre></body></html>", res))
}
