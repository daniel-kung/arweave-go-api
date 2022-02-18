package e

type LocaleLanguage int

const (
	LocaleLanguageDefault LocaleLanguage = 0
	LocaleLanguageEnglish LocaleLanguage = LocaleLanguageDefault
	LocaleLanguageChinese LocaleLanguage = 1
)

var (
	gLanguage = LocaleLanguageDefault
)

func ChangeLanguage(language LocaleLanguage) {
	gLanguage = language
}

var chineseErrMsgs = map[ERRCode]string{
	ERRCodeSuccess:         "成功",
	ERRCodeInvalidParam:    "请求参数错误",
	ERRCodeChainConnection: "链节点连接失败",
}

var englishErrMsgs = map[ERRCode]string{
	ERRCodeSuccess:         "success",
	ERRCodeInvalidParam:    "invalid parameter",
	ERRCodeChainConnection: "chain connection failed",
}

// GetMsg returns error message based on Code
func GetMsg(code ERRCode) string {
	var msgs map[ERRCode]string
	switch gLanguage {
	case LocaleLanguageEnglish:
		msgs = englishErrMsgs
	case LocaleLanguageChinese:
		msgs = chineseErrMsgs
	}
	msg, ok := msgs[code]
	if ok {
		return msg
	}

	return ""
}
