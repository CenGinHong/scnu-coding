package language_enum

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
