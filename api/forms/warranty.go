package forms

type CreateWarrantyCommand struct {
	Name     string   `json:"name" binding:"required"`
	Phone    string   `json:"phone" binding:"required"`
	SoXe     string   `json:"soxe" binding:"required"`
	TenXe    string   `json:"tenxe" binding:"required"`
	LoaiKinh []string `json:"loai_kinh" binding:"required"`
	NgayDan  string   `json:"ngay_dan" binding:"required"`
	Agency   int      `json:"agency_id" binding:"required"`
	DaiLy    int      `json:"dai_ly" binding:"required"`
	DoiXe    string   `json:"doixe" binding:"required"`
	MauXe    string   `json:"mauxe" binding:"required"`
	LoaiXe   string   `json:"loaixe" binding:"required"`
	ViTriDan []*Vitri `json:"vi_tri_dan" binding:"required,dive"`
	Email    string   `json:"email"`
	Note     string   `json:"note"`
}

type Vitri struct {
	ID    string `json:"id" binding:"required"`
	Value string `json:"value" `
}

type UpdateWarrantyCommand struct {
	Name  string `json:"name" binding:"required"`
	Phone string `json:"phone" binding:"required"`
	DoiXe string `json:"doi_xe" binding:"required"`
	MauXe string `json:"mau_xe" binding:"required"`
	TenXe string `json:"tenxe" binding:"required"`
	Email string `json:"email"`
	Note  string `json:"note"`
	DaiLy int    `json:"dai_ly" binding:"required"`
}

type Request struct {
	Phone       string
	Soxe        string
	Maso        string
	Agency      string
	Daily       string
	Ngaydanfrom string
	Ngaydanto   string
}
