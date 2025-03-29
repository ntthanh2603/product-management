package response

const (
	ErrCodeSuccess      = 20001 // success
	ErrCodeParamInvalid = 20003 // Email is invalid
)

var msg = map[int]string{
	ErrCodeSuccess:      "success",
	ErrCodeParamInvalid: "Email is invalid",
}
