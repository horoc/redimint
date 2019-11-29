package code

const (
	CodeTypeOK               = 0
	CodeTypeInvalidSign      = 1
	CodeTypeInvalidValidator = 2
	CodeTypeEncodingError    = 3
)


func InfoWithDetail(p int, msg string) string {
	return Info(p) + ": " + msg
}

func Info(p int) string {
	switch (p) {
	case CodeTypeOK:
		return "Success!"
	case CodeTypeInvalidSign:
		return "Invalid Node Signature"
	case CodeTypeInvalidValidator:
		return "Invalid validator"
	case CodeTypeEncodingError:
		return "Encoding error"
	default:
		return "UNKNOWN"
	}
}
