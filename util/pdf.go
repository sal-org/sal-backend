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
