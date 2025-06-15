package internalhttp

type IDResponse struct {
	ID int `json:"id"`
}

func ToIDResponse(id int) *IDResponse {
	idRes := new(IDResponse)
	idRes.ID = id

	return idRes
}

func (i *IDResponse) Int() int {
	return i.ID
}
