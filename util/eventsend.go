package util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"text/template"
)

type HtmlEventdatastruct struct {
	Counsellor_id string
	Type          string
	Title         string
	Description   string
	Photo         string
	Topic_id      string
	Date          string
	Time          string
	Duration      string
	Price         string
}

type HtmlCounsellorProfileForApproval struct {
	First_Name  string
	Last_Name   string
	Gender      string
	Phone       string
	Photo       string
	Email       string
	Education   string
	Experience  string
	About       string
	Resume      string
	Certificate string
	Aadhar      string
	Linkedin    string
	Status      string
}

func GetHTMLTemplateForEvent(orderdetails []map[string]string) string {
	var templateBuffer bytes.Buffer

	htmlData, err := ioutil.ReadFile("htmlfile/Event.html")
	if err != nil {
		return ""
	}

	neworderdetails := HtmlEventdatastruct{
		Counsellor_id: orderdetails[0]["counsellor_id"],
		Type:          orderdetails[0]["type"],
		Title:         orderdetails[0]["title"],
		Description:   orderdetails[0]["description"],
		Photo:         orderdetails[0]["photo"],
		Topic_id:      orderdetails[0]["topic_id"],
		Date:          orderdetails[0]["date"],
		Time:          orderdetails[0]["time"],
		Duration:      orderdetails[0]["duration"],
		Price:         orderdetails[0]["price"],
	}

	htmlTemplate := template.Must(template.New("Event.html").Parse(string(htmlData)))

	err = htmlTemplate.ExecuteTemplate(&templateBuffer, "Event.html", neworderdetails)

	if err != nil {
		fmt.Println("Data is not fetch")
		return ""

	}

	return templateBuffer.String()
}

func GetHTMLTemplateForCounsellor(orderdetails []map[string]string) string {
	var templateBuffer bytes.Buffer

	htmlData, err := ioutil.ReadFile("htmlfile/CounsellorProfile.html")
	if err != nil {
		return ""
	}

	neworderdetails := HtmlCounsellorProfileForApproval{
		First_Name:  orderdetails[0]["first_name"],
		Last_Name:   orderdetails[0]["last_name"],
		Gender:      orderdetails[0]["gender"],
		Phone:       orderdetails[0]["phone"],
		Photo:       orderdetails[0]["photo"],
		Email:       orderdetails[0]["email"],
		Education:   orderdetails[0]["education"],
		Experience:  orderdetails[0]["experience"],
		About:       orderdetails[0]["about"],
		Resume:      orderdetails[0]["resume"],
		Certificate: orderdetails[0]["certificate"],
		Aadhar:      orderdetails[0]["aadhar"],
		Linkedin:    orderdetails[0]["linkedin"],
		Status:      orderdetails[0]["status"],
	}

	htmlTemplate := template.Must(template.New("CounsellorProfile.html").Parse(string(htmlData)))

	err = htmlTemplate.ExecuteTemplate(&templateBuffer, "CounsellorProfile.html", neworderdetails)

	if err != nil {
		fmt.Println("Data is not fetch")
		return ""

	}

	return templateBuffer.String()
}
