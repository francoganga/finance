package routes

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	_ "os"
	"os/exec"
	_ "path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/francoganga/pagoda_bun/models"
	"github.com/francoganga/pagoda_bun/pkg/controller"
	"github.com/francoganga/pagoda_bun/pkg/internal/parser"
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

	form, err := ctx.MultipartForm()

	if err != nil {
		return err
	}

	files := form.File["pdfs"]

	for _, file := range files {

		matches, err := getMatchesFromFile(file)

		if err != nil {
			return err
		}

        fmt.Printf("%+v\n", file.Filename)
		for _, line := range matches {
			fmt.Println(strings.Replace(line, "\n", "\\n", -1))
		}

		consumos := make([]*parser.ConsumoDto, 0)

		transactions := make([]*models.Transaction, 0)

		for _, line := range matches {
			p := parser.FromInput(line)

			consu := p.ParseConsumo()

			if len(p.Errors()) > 0 {

				msg := ""

				for _, e := range p.Errors() {
					msg += e
				}

				return fmt.Errorf("error parsing consumo: %s", msg)

			}

			pd, err := time.Parse("02/01/06", consu.Date)

			if err != nil {
				return fmt.Errorf("error in time.parse consumo: %w", err)
			}

			t := &models.Transaction{
				Date:        pd,
				Code:        consu.Code,
				Description: consu.Description,
				Amount:      consu.Amount,
				Balance:     consu.Balance,
			}

			transactions = append(transactions, t)

			consumos = append(consumos, consu)
		}

		// res := strings.Join(matches, "\n")

		_, err = c.Container.Bun.NewInsert().
			Model(&transactions).
			Exec(ctx.Request().Context())

		if err != nil {
			return err
		}

	}

	page := controller.NewPage(ctx)

	page.Layout = "main"
	page.Name = "load-pdf"
	page.Title = "Load PDF"

	return c.RenderPage(ctx, page)
	// return ctx.HTML(http.StatusOK, fmt.Sprintf("<html><body><pre>%s</pre></body></html>", "asd"))
}

func getMatchesFromFile(file *multipart.FileHeader) ([]string, error) {

	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	contents := make([]byte, file.Size)

	_, err = src.Read(contents)

	if err != nil {
		return nil, err
	}

	command := exec.Command("pdftotext", "-layout", "-f", "1", "-l", "3", "-", "-")

	stdin, err := command.StdinPipe()

	if err != nil {
		return nil, err
	}

	var outb bytes.Buffer

	command.Stdout = &outb

	if err = command.Start(); err != nil { //Use start, not run
		fmt.Println("An error occured: ", err) //replace with logger, or anything you want
	}

	_, err = io.WriteString(stdin, string(contents))

	if err != nil {
		return nil, err
	}

	stdin.Close()

	err = command.Wait()

	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`(?m)^([0-9]{2}/[0-9]{2}/[0-9]{2})\s+([0-9]+)\s+(.*?)\s{2,}(.*?)\s{2,}(.*)\n(.*)`)

	fmt.Println(outb.String())

	return re.FindAllString(outb.String(), -1), nil
}
