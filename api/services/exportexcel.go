package services

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"log"
	"strconv"
	"time"
	"warranty-api/forms"
	"warranty-api/models"
)

func GenerateExcel(warrantyList []models.Warranty) string {

	f := excelize.NewFile()
	id := uuid.New()

	style, _ := f.NewStyle(`{"alignment":{"horizontal":"left"}, "font":{"bold":true,"italic":false}}`)
	f.SetColWidth("Sheet1", "A", "Z", 25)
	f.SetColWidth("Sheet1", "J", "J", 40)
	f.SetCellStyle("Sheet1", MapRowCol(0, 1), MapRowCol(20, 1), style)
	f.SetCellValue("Sheet1", MapRowCol(0, 1), "Họ Tên ")
	f.SetCellValue("Sheet1", MapRowCol(1, 1), "Số Điện Thoại ")
	f.SetCellValue("Sheet1", MapRowCol(2, 1), "Email")
	f.SetCellValue("Sheet1", MapRowCol(3, 1), "Số Xe")
	f.SetCellValue("Sheet1", MapRowCol(4, 1), "Tên Xe ")
	f.SetCellValue("Sheet1", MapRowCol(5, 1), "Loại Kính ")
	f.SetCellValue("Sheet1", MapRowCol(6, 1), "Ngày Dán ")
	f.SetCellValue("Sheet1", MapRowCol(7, 1), "Nhà Phân Phối ")
	f.SetCellValue("Sheet1", MapRowCol(8, 1), "Đại Lý")
	f.SetCellValue("Sheet1", MapRowCol(9, 1), "Mã Số Bảo Hành")
	f.SetCellValue("Sheet1", MapRowCol(10, 1), "Đời Xe")
	f.SetCellValue("Sheet1", MapRowCol(11, 1), "Màu Xe")
	f.SetCellValue("Sheet1", MapRowCol(12, 1), "Loại Xe")
	f.SetCellValue("Sheet1", MapRowCol(13, 1), "Cửa sổ trời")
	f.SetCellValue("Sheet1", MapRowCol(14, 1), "Kính sườn trước")
	f.SetCellValue("Sheet1", MapRowCol(15, 1), "Kính sườn sau")
	f.SetCellValue("Sheet1", MapRowCol(16, 1), "Kính khoang sau")
	f.SetCellValue("Sheet1", MapRowCol(17, 1), "Kính Hậu")
	f.SetCellValue("Sheet1", MapRowCol(18, 1), "Kính Lái")
	f.SetCellValue("Sheet1", MapRowCol(19, 1), "Ghi Chú")

	for i, warranty := range warrantyList {
		i = i + 2
		var ngaydan = warranty.NgayDan.Add(time.Hour * 7)
		f.SetCellValue("Sheet1", MapRowCol(0, i), warranty.Name)
		f.SetCellValue("Sheet1", MapRowCol(1, i), warranty.Phone)
		f.SetCellValue("Sheet1", MapRowCol(2, i), warranty.Email)
		f.SetCellValue("Sheet1", MapRowCol(3, i), warranty.Soxe)
		f.SetCellValue("Sheet1", MapRowCol(4, i), warranty.Tenxe)
		f.SetCellValue("Sheet1", MapRowCol(5, i), warranty.LoaiKinh)
		f.SetCellValue("Sheet1", MapRowCol(6, i), fmt.Sprintf("%v", ngaydan.Day())+"-"+fmt.Sprintf("%v", int(ngaydan.Month()))+"-"+fmt.Sprintf("%v", ngaydan.Year()))
		f.SetCellValue("Sheet1", MapRowCol(7, i), warranty.Agency.Name)
		f.SetCellValue("Sheet1", MapRowCol(8, i), warranty.Daily.Name)
		f.SetCellValue("Sheet1", MapRowCol(9, i), warranty.MaSo)
		f.SetCellValue("Sheet1", MapRowCol(10, i), warranty.DoiXe)
		f.SetCellValue("Sheet1", MapRowCol(11, i), warranty.MauXe)
		f.SetCellValue("Sheet1", MapRowCol(12, i), warranty.LoaiXe)

		var vitridan []forms.Vitri
		if json.Valid([]byte(warranty.ViTriDan)) {
			_ = json.Unmarshal([]byte(warranty.ViTriDan), &vitridan)
			for _, s := range vitridan {
				var value = ""
				if s.ID == "cua_so_troi" {
					value = GetValue(s.Value)
					f.SetCellValue("Sheet1", MapRowCol(13, i), value)

				}

				if s.ID == "kinh_suon_truoc" {
					value = GetValue(s.Value)
					f.SetCellValue("Sheet1", MapRowCol(14, i), value)

				}

				if s.ID == "kinh_suon_sau" {
					value = GetValue(s.Value)
					f.SetCellValue("Sheet1", MapRowCol(15, i), value)

				}

				if s.ID == "kinh_khoang_sau" {
					value = GetValue(s.Value)
					f.SetCellValue("Sheet1", MapRowCol(16, i), value)

				}

				if s.ID == "kinh_hau" {
					value = GetValue(s.Value)
					f.SetCellValue("Sheet1", MapRowCol(17, i), value)

				}

				if s.ID == "kinh_lai" {
					value = GetValue(s.Value)
					f.SetCellValue("Sheet1", MapRowCol(18, i), value)
				}
			}
		} else {
			f.SetCellValue("Sheet1", MapRowCol(20, i), warranty.ViTriDan)
		}

		f.SetCellValue("Sheet1", MapRowCol(19, i), warranty.Note)
	}

	if err := f.SaveAs("download/export_" + id.String() + ".xlsx"); err != nil {
		log.Fatal(err)
	}

	path := "/download/export_" + id.String() + ".xlsx"
	return path
}

func MapRowCol(col int, row int) string {
	array := [21]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V"}
	return array[col] + strconv.Itoa(row)
}
