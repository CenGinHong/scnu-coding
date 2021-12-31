package language_enum

const (
	Full = iota
	Cpp
	Java
	Python
)

func Num2LanguageString(languageEnum int) (LanguageName string) {
	switch languageEnum {
	case Full:
		LanguageName = "full"
	case Cpp:
		LanguageName = "cpp"
	case Java:
		LanguageName = "java"
	case Python:
		LanguageName = "python"
	}
	return LanguageName
}
