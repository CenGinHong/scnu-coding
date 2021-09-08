package enum

const (
	Cpp = iota + 1
	Java
	Python
)

func Num2LanguageString(languageEnum int) (LanguageName string) {
	switch languageEnum {
	case Cpp:
		LanguageName = "Cpp"
	case Java:
		LanguageName = "Java"
	case Python:
		LanguageName = "Python"
	}
	return LanguageName
}

func Num2IsCloseString(isClose int) (isCloseString string) {
	switch isClose {
	case 0:
		isCloseString = "已结课"
	case 1:
		isCloseString = "进行中"
	}
	return isCloseString
}
