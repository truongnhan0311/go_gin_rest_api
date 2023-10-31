package services

import (
	"encoding/json"
	"fmt"
	"github.com/signintech/gopdf"
	"strings"
	"time"
	"warranty-api/forms"
	"warranty-api/models"
)

func GeneratePDF(warranty models.Warranty) string {
	var err error
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: 679, H: 206}}) //595.28, 841.89 = A4
	pdf.AddPage()
	err = pdf.AddTTFFont("arial", "./config/arial.ttf")
	if err != nil {
		panic(err)
	}
	err = pdf.SetFont("arial", "", 9)
	if err != nil {
		panic(err)
	}
	// Color the page
	pdf.SetLineWidth(0.1)
	//pdf.SetFillColor(124, 252, 0) //setup fill color
	//pdf.RectFromUpperLeftWithStyle(50, 100, 400, 600, "FD")
	//pdf.SetFillColor(0, 0, 0)

	tpl1 := pdf.ImportPage("./config/phieu-bao-hanh.pdf", 1, "/MediaBox")
	tpl2 := pdf.ImportPage("./config/phieu-bao-hanh.pdf", 2, "/MediaBox")
	pdf.UseImportedTemplate(tpl1, 0, 0, 0, 0)

	pdf.SetX(27)
	pdf.SetY(27)
	pdf.Cell(nil, warranty.Daily.Name)

	pdf.SetX(38)
	pdf.SetY(38)
	pdf.Cell(nil, warranty.Daily.Address)

	pdf.SetX(52)
	pdf.SetY(49)
	pdf.Cell(nil, warranty.Daily.Phone)

	pdf.SetX(37)
	pdf.SetY(59)
	pdf.Cell(nil, warranty.Daily.Email)

	//nhan vien
	pdf.SetX(50)
	pdf.SetY(70)
	pdf.Cell(nil, warranty.User.Name)

	pdf.SetX(30)
	pdf.SetY(112)
	pdf.Cell(nil, warranty.Name)

	pdf.SetX(52)
	pdf.SetY(123)
	pdf.Cell(nil, warranty.Phone)

	pdf.SetX(35)
	pdf.SetY(134)
	pdf.Cell(nil, warranty.Email)

	pdf.SetX(40)
	pdf.SetY(145)
	pdf.Cell(nil, warranty.LoaiXe+" "+warranty.MauXe+" "+warranty.DoiXe)

	pdf.SetX(60)
	pdf.SetY(155)
	pdf.Cell(nil, warranty.Tenxe)

	pdf.SetX(50)
	pdf.SetY(165)
	pdf.Cell(nil, warranty.Soxe)

	var ngaydan = warranty.NgayDan.Add(time.Hour * 7)
	pdf.SetX(60)
	pdf.SetY(176)
	pdf.Cell(nil, fmt.Sprintf("%v", ngaydan.Day())+"-"+fmt.Sprintf("%v", int(ngaydan.Month()))+"-"+fmt.Sprintf("%v", ngaydan.Year()))

	loaiKinh := strings.Split(warranty.LoaiKinh, ",")
	for _, kinh := range loaiKinh {

		if kinh == "clearplex" {
			pdf.SetX(181)
			pdf.SetY(191)
			pdf.Cell(nil, "x")
		}

		if kinh == "clearplexir" {
			pdf.SetX(181)
			pdf.SetY(191)
			pdf.Cell(nil, "x")
		}

		if kinh == "window_film" {
			pdf.SetX(239)
			pdf.SetY(191)
			pdf.Cell(nil, "x")
		}

		if kinh == "ppf" {
			pdf.SetX(301)
			pdf.SetY(191)
			pdf.Cell(nil, "x")
		}
	}

	var vitridan []forms.Vitri
	_ = json.Unmarshal([]byte(warranty.ViTriDan), &vitridan)

	for _, s := range vitridan {
		var value = ""
		if s.ID == "cua_so_troi" {
			err = pdf.SetFont("arial", "", 4)
			value = GetValue(s.Value)
			pdf.SetX(169)
			pdf.SetY(64)
			pdf.Cell(nil, value)
		}

		if s.ID == "kinh_suon_truoc" {
			err = pdf.SetFont("arial", "", 4)
			value = GetValue(s.Value)
			pdf.SetX(169)
			pdf.SetY(91)
			pdf.Cell(nil, value)

			pdf.SetX(290)
			pdf.SetY(91)
			pdf.Cell(nil, value)
		}

		if s.ID == "kinh_suon_sau" {
			err = pdf.SetFont("arial", "", 4)
			value = GetValue(s.Value)
			pdf.SetX(169)
			pdf.SetY(115)
			pdf.Cell(nil, value)

			pdf.SetX(290)
			pdf.SetY(115)
			pdf.Cell(nil, value)
		}

		if s.ID == "kinh_khoang_sau" {
			err = pdf.SetFont("arial", "", 4)
			value = GetValue(s.Value)
			pdf.SetX(169)
			pdf.SetY(141)
			pdf.Cell(nil, value)

			pdf.SetX(287)
			pdf.SetY(141)
			pdf.Cell(nil, value)
		}

		if s.ID == "kinh_hau" {
			err = pdf.SetFont("arial", "", 4)
			value = GetValue(s.Value)
			pdf.SetX(169)
			pdf.SetY(165)
			pdf.Cell(nil, value)
		}

		if s.ID == "kinh_lai" {
			err = pdf.SetFont("arial", "", 4)
			value = GetValue(s.Value)
			pdf.SetX(293)
			pdf.SetY(64)
			pdf.Cell(nil, value)
		}
	}

	pdf.AddPage()
	pdf.UseImportedTemplate(tpl2, 0, 0, 0, 0)

	path := "/download/" + warranty.MaSoBH + ".pdf"
	pdf.WritePdf("." + path)
	return path
}

func GetValue(loaiKinh string) string {
	if loaiKinh == "clearplexir" {
		return "IR"
	}
	if loaiKinh == "clearplex" || loaiKinh == "ppf" {
		return ""
	}
	return loaiKinh
}
