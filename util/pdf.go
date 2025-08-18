package util

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	Model "salbackend/model"
	CONFIG "salbackend/config"
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

func GetHTMLTemplateForCounsellorCancellation(data Model.EmailDataForCounsellorCancellation, filepath string) string {
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

func GetHTMLTemplateForAppFeedBack(data Model.EmailDataForFeedback, filepath string) string {
	var templateBuffer bytes.Buffer

	// You can bind custom data here as per requirements.

	htmlData, err := ioutil.ReadFile(filepath)
{}
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

func GetHTMLTemplateForWebB2CClientProfile(data Model.EmailDataForWebClientB2CProfile, filepath string) string {
	var templateBuffer bytes.Buffer

	// You can bind custom data here as per requirements.

	htmlData, err := ioutil.ReadFile(filepath)

	if err != nil {
		log.Fatal(err)
		return ""
	}

	htmlTemplate := template.Must(template.New("email1.html").Parse(string(htmlData)))

	err = htmlTemplate.ExecuteTemplate(&templateBuffer, "email1.html", data)

	if err != nil {
		log.Fatal(err)
		return ""
	}

	return templateBuffer.String()
}

func GetHTMLTemplateForCounsellorRecord(data Model.EmailDataForCounsellorRecord, filepath string) string {
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

func GetHTMLTemplateForCounsellorVisit(data Model.EmailDataForCounsellorVisit, filepath string) string {
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

func GetHTMLTemplateForClientConfirmationWithAccessCodeText(data Model.EmailBodyWithAccessCodeMessageModel, filepath string) string {
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

func GetHTMLTemplateForWithDocument(data Model.EmailBodyMessageModelWithDocu, filepath string) string {
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

func GetHTMLTemplateForAssessmentSelfEsteem(data Model.AssessmentDownloadSelfEsteemModel, filepath string) (string, bool) {
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

func GetHTMLTemplateForAssessmentBurnOut(data Model.AssessmentDownloadBurnOutModel, filepath string) (string, bool) {
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

func GetHTMLTemplateForAssessmentGWB(data Model.AssessmentDownloadGWBModel, filepath string) (string, bool) {
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

// htmlToPDF converts HTML content to a PDF file using chromedp
func HtmlToPDFInvoice(htmlContent string) []byte {

	// Create pdf using api for that html content
	var pdfBuffer []byte

	body := Model.PostRequestForHTMLToPDF{
		HtmlContent: htmlContent,
	}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(body)

	client := &http.Client{}
	req, err := http.NewRequest("POST", CONFIG.PDFURL, payloadBuf)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	// Add Authorization header
	req.Header.Add("Content-Type", "application/json")

	// Set timeout for the request
	fmt.Println(" for the request")
	// Send HTTP request
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer res.Body.Close()

	// Read the response body
	fmt.Println("Reading response body")
	if res.StatusCode != http.StatusOK {
		fmt.Println("Error: received non-200 response status:", res.StatusCode)
		return nil
	}

	bodyy, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	mapp := make(map[string]string)

	err = json.Unmarshal(bodyy, &mapp)
	if err != nil {
		return nil
	}

	// Check if the response status is 200
	fmt.Println("Response status:", mapp)

	if mapp["status"] != "200" {
		fmt.Println("Error in response:", mapp["message"])
		return nil
	}

	// Decode base64 string back to bytes
	pdfBuffer, err = base64.StdEncoding.DecodeString(mapp["body"])
	if err != nil {
		fmt.Println("Error decoding base64 string:", err)
		return nil
	}

	return pdfBuffer

}

// htmlToPDF converts HTML content to a PDF file using chromedp
func HtmlToPDFAssessment(htmlContent string) []byte {

	// // Header HTML

	// Get base64 encoded image for header
	// headerImageBase64, err := getBase64Image("https://sal-prod.s3.ap-south-1.amazonaws.com/miscellaneous/assessment_header.png")
	// if err != nil {
	// 	fmt.Println("Error getting base64 image:", err)
	// 	return nil
	// }

	// Create a context

	customCSS := `
	<style>
		@media print {
			/* Hide default headers/footers */
			@page {
				margin: 0.75in 0.4in 0.75in 0.4in;
			}

			/* First page header/footer */
			.first-page-header, .first-page-footer {
				display: none;
			}
			.page:first-of-type .first-page-header,
			.page:first-of-type .first-page-footer {
				display: block;
			}

			/* Last page header/footer */
			.last-page-header, .last-page-footer {
				display: none;
			}
			.page:last-of-type .last-page-header,
			.page:last-of-type .last-page-footer {
				display: block;
			}

			/* Page structure */
			.page {
				page-break-after: always;
				position: relative;
				padding-top: 180px;  /* Space for header */
				padding-bottom: 150px; /* Space for footer */
			}

			/* Header positioning */
			.first-page-header, .last-page-header {
				position: absolute;
				top: 0;
				left: 0;
				right: 0;
				height: 40px;
			}

			/* Footer positioning */
			.first-page-footer, .last-page-footer {
				position: absolute;
				bottom: 0;
				left: 0;
				right: 0;
				height: 130px;
			}
		}
	</style>
	`

	//Create first page header with the image
	firstPageHeader := `
	<div class="first-page-header">
		<div style="text-align: center; line-height: 1.0;">
			<img src="https://sal-prod.s3.ap-south-1.amazonaws.com/miscellaneous/assessment_header.png" alt="Assessment Header" style="max-width: 80%%; height: auto;"/>
		</div>
	</div>
	`

	// Create first page footer (if needed)
	firstPageFooter := `
	<div class="first-page-footer">
		<div style="width: 100%; font-size: 10pt; line-height: 1.4; font-family: Arial, sans-serif; color: #333333; padding: 10px 20px;">
			<hr style="border: none; height: 1px; background-color: #999999; margin-bottom: 12px;">
			<div style="font-size: 10pt; font-weight: bold; margin-bottom: 8px;">Salubrium Private Limited</div>
			<div style="font-size: 8pt; margin-bottom: 8px;">If you or someone you know is experiencing suicidal thoughts, please reach out to the nearest hospital that offers emergency services. Clove does not offer these services at this point in time.</div>
			<div style="font-size: 8pt; font-weight: bold;">2025 Salubrium Private Limited. All Rights Reserved.</div>
		</div>
	</div>
	`

	// and adding our custom headers/footers
	enhancedHTMLContent := customCSS

	// Split content by page breaks or add page wrapper based on your needs
	// This is a simplistic approach - you may need to modify based on your actual HTML structure
	enhancedHTMLContent += `<div class="page">`
	enhancedHTMLContent += firstPageHeader
	enhancedHTMLContent += htmlContent
	enhancedHTMLContent += firstPageFooter
	enhancedHTMLContent += `</div>`

	// Create pdf using api for that html content
	var pdfBuffer []byte

	body := Model.PostRequestForHTMLToPDF{
		HtmlContent: enhancedHTMLContent,
	}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(body)

	client := &http.Client{}
	req, err := http.NewRequest("POST", CONFIG.PDFURL, payloadBuf)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	// Add Authorization header
	req.Header.Add("Content-Type", "application/json")

	// Set timeout for the request
	fmt.Println(" for the request")
	// Send HTTP request
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer res.Body.Close()

	// Read the response body
	fmt.Println("Reading response body")
	if res.StatusCode != http.StatusOK {
		fmt.Println("Error: received non-200 response status:", res.StatusCode)
		return nil
	}

	bodyy, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	mapp := make(map[string]string)

	err = json.Unmarshal(bodyy, &mapp)
	if err != nil {
		return nil
	}

	// Check if the response status is 200
	fmt.Println("Response status:", mapp)

	if mapp["status"] != "200" {
		fmt.Println("Error in response:", mapp["message"])
		return nil
	}

	// Decode base64 string back to bytes
	pdfBuffer, err = base64.StdEncoding.DecodeString(mapp["body"])
	if err != nil {
		fmt.Println("Error decoding base64 string:", err)
		return nil
	}

	return pdfBuffer

}
