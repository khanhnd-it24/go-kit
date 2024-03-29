// File generated by fault/cmd/gen.go. DO NOT EDIT.

package fault

var (
	KeyAuthInvalidToken = "auth_invalid_token"
)

var (
	keyFactory = map[string]string{
		KeyAuthInvalidToken: "invalid token",
	}
)

var (
	codeFactory = map[string]int64{
		KeyAuthInvalidToken: 1000,
	}
)

func getDescriptionFromKey(code string) string {
	if des, ok := keyFactory[code]; ok {
		return des
	}
	return ""
}

const (
	CodeUnknown = -1
)

func getCodeFromKey(code string) int64 {
	if des, ok := codeFactory[code]; ok {
		return des
	}
	return CodeUnknown
}
