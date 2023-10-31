package forms

type CreateAgencyCommand struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`
	MaSoDL  string `json:"ma_so_dl" binding:"required"`
	Phone   string `json:"phone" binding:"required"`
}
