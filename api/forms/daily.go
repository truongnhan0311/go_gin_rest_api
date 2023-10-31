package forms

type CreateDailyCommand struct {
	Name     string `json:"name" binding:"required"`
	Address  string `json:"address" binding:"required"`
	MaSoDL   string `json:"ma_so_dl" binding:"required"`
	AgencyID int    `json:"agency_id" binding:"required"`
	Province string `json:"province" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
}
