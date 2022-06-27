package util

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	Model "salbackend/model"
	s "strings"

	wkhtml "github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

func GetHTMLTemplateForEvent(data Model.EmailDataForEvent, filepath string) string {
	var templateBuffer bytes.Buffer

	// You can bind custom data here as per requirements.

	htmlData, err := ioutil.ReadFile(filepath)

	if err != nil {
		log.Fatal(err)
		return ""
	}

	htmlTemplate := template.Must(template.New("email.html").Parse(string(htmlData)))

	err = htmlTemplate.ExecuteTemplate(&templateBuffer, "email.html", data)

	if err != nil {
		log.Fatal(err)
		return ""
	}

	return templateBuffer.String()
}

func GetHTMLTemplateForProfile(data Model.EmailDataForCounsellorProfile, filepath string) string {
	var templateBuffer bytes.Buffer

	// You can bind custom data here as per requirements.

	htmlData, err := ioutil.ReadFile(filepath)

	if err != nil {
		log.Fatal(err)
		return ""
	}

	htmlTemplate := template.Must(template.New("email.html").Parse(string(htmlData)))

	err = htmlTemplate.ExecuteTemplate(&templateBuffer, "email.html", data)

	if err != nil {
		log.Fatal(err)
		return ""
	}

	return templateBuffer.String()
}

func GetHTMLTemplateForCounsellorProfileText(data Model.EmailBodyMessageModel, filepath string) string {
	var templateBuffer bytes.Buffer

	// You can bind custom data here as per requirements.

	htmlData, err := ioutil.ReadFile(filepath)

	if err != nil {
		log.Fatal(err)
		return ""
	}

	htmlTemplate := template.Must(template.New("email.html").Parse(string(htmlData)))

	err = htmlTemplate.ExecuteTemplate(&templateBuffer, "email.html", data)

	if err != nil {
		log.Fatal(err)
		return ""
	}

	return templateBuffer.String()
}

func GetHTMLTemplateForClientAppointmentConfirmation(data Model.ClientAppointmentConfirmation, filepath string) string {
	var templateBuffer bytes.Buffer

	// You can bind custom data here as per requirements.

	htmlData, err := ioutil.ReadFile(filepath)

	if err != nil {
		log.Fatal(err)
		return ""
	}

	htmlTemplate := template.Must(template.New("email.html").Parse(string(htmlData)))

	err = htmlTemplate.ExecuteTemplate(&templateBuffer, "email.html", data)

	if err != nil {
		log.Fatal(err)
		return ""
	}

	return templateBuffer.String()
}

func GetHTMLTemplateForReceipt(data Model.EmailDataForPaymentReceipt, filepath string) (string, bool) {
	var templateBuffer bytes.Buffer

	// You can bind custom data here as per requirements.

	htmlData, err := ioutil.ReadFile(filepath)

	if err != nil {
		fmt.Println("file is not read")
		return "", false
	}

	htmlTemplate := template.Must(template.New("email.html").Parse(string(htmlData)))

	err = htmlTemplate.ExecuteTemplate(&templateBuffer, "email.html", data)

	if err != nil {
		fmt.Println("Data not pass in html")
		return "", false
	}

	return templateBuffer.String(), true
}

func GetHTMLTemplateForAssessmentAIS(data Model.AssessmentDownloadAIS, filepath string) (string, bool) {
	var templateBuffer bytes.Buffer

	// You can bind custom data here as per requirements.

	htmlData, err := ioutil.ReadFile(filepath)

	if err != nil {
		fmt.Println("file is not read")
		return "", false
	}

	htmlTemplate := template.Must(template.New("email.html").Parse(string(htmlData)))

	err = htmlTemplate.ExecuteTemplate(&templateBuffer, "email.html", data)

	if err != nil {
		fmt.Println("Data not pass in html")
		return "", false
	}

	return templateBuffer.String(), true
}

func GetHTMLTemplateForAssessmentBDI(data Model.AssessmentDownloadBDIModel, filepath string) (string, bool) {
	var templateBuffer bytes.Buffer

	// You can bind custom data here as per requirements.

	htmlData, err := ioutil.ReadFile(filepath)

	if err != nil {
		fmt.Println("file is not read")
		return "", false
	}

	htmlTemplate := template.Must(template.New("email.html").Parse(string(htmlData)))

	err = htmlTemplate.ExecuteTemplate(&templateBuffer, "email.html", data)

	if err != nil {
		fmt.Println("Data not pass in html")
		return "", false
	}

	return templateBuffer.String(), true
}

func GetHTMLTemplateForAssessmentSRS(data Model.AssessmentDownloadSRSModel, filepath string) (string, bool) {
	var templateBuffer bytes.Buffer

	// You can bind custom data here as per requirements.

	htmlData, err := ioutil.ReadFile(filepath)

	if err != nil {
		fmt.Println("file is not read")
		return "", false
	}

	htmlTemplate := template.Must(template.New("email.html").Parse(string(htmlData)))

	err = htmlTemplate.ExecuteTemplate(&templateBuffer, "email.html", data)

	if err != nil {
		fmt.Println("Data not pass in html")
		return "", false
	}

	return templateBuffer.String(), true
}

func GetHTMLTemplateForAssessmentGAD7(data Model.AssessmentDownloadGAD7Model, filepath string) (string, bool) {
	var templateBuffer bytes.Buffer

	// You can bind custom data here as per requirements.

	htmlData, err := ioutil.ReadFile(filepath)

	if err != nil {
		fmt.Println("file is not read")
		return "", false
	}

	htmlTemplate := template.Must(template.New("email.html").Parse(string(htmlData)))

	err = htmlTemplate.ExecuteTemplate(&templateBuffer, "email.html", data)

	if err != nil {
		fmt.Println("Data not pass in html")
		return "", false
	}

	return templateBuffer.String(), true
}

func GeneratePdf(htmlfile, filepath string) ([]byte, bool) { // ([]byte
	pdfg, err := wkhtml.NewPDFGenerator()
	if err != nil {
		fmt.Println("New not pdf sesssion")
		return nil, false
	}
	pdfg.AddPage(wkhtml.NewPageReader(s.NewReader(htmlfile)))

	// Create PDF document in internal buffer
	err = pdfg.Create()
	if err != nil {
		return nil, false
	}

	//Your Pdf Name
	return pdfg.Bytes(), true

}

func GeneratePdfHeaderAndFooterFixted(htmlfile, filepath string) ([]byte, bool) { // ([]byte
	pdfg, err := wkhtml.NewPDFGenerator()
	if err != nil {
		fmt.Println("New not pdf sesssion")
		return nil, false
	}
	page := wkhtml.NewPageReader(s.NewReader(htmlfile))
	page.DisableExternalLinks.Set(false)
	// page.PageOptions.HeaderHTML.Set(`<img src='https://sal-prod.s3.ap-south-1.amazonaws.com/miscellaneous/assessment_header.png' alt='Assessment Header' height='80%'' width='100%'/>`)
	page.FooterFontSize.Set(7)
	page.FooterRight.Set("[page]")
	page.EnableLocalFileAccess.Set(true)
	page.HeaderHTML.Set(`D:\Git\SALBackend\htmlfile\AssessmentHeader.html`)
	page.FooterHTML.Set("htmlfile/AssessmentFooter.html")
	pdfg.AddPage(page)

	// Create PDF document in internal buffer
	err = pdfg.Create()
	if err != nil {
		return nil, false
	}

	// err = pdfg.WriteFile(`D:\TestHFHTML.pdf`)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	//Your Pdf Name
	return pdfg.Bytes(), true

}
