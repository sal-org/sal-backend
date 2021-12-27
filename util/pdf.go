package util

import (
	s "strings"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

func GeneratePdf(pdfPath string, htmlfile string) error {
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return err
	}

	pdfg.AddPage(wkhtmltopdf.NewPageReader(s.NewReader(htmlfile)))

	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)

	pdfg.Dpi.Set(300)

	err = pdfg.Create()
	if err != nil {
		return err
	}

	err = pdfg.WriteFile(pdfPath)
	if err != nil {
		return err
	}

	return nil
}
