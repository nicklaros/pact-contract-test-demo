package errors

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
