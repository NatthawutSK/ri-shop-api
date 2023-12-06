package entities

type Image struct {
	Id       string `json:"id" db:"id"`
	FileName string `json:"filename" db:"filename"`
	Url      string `json:"url" db:"url"`
}
