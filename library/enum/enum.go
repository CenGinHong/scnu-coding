package enum

func Num2IsCloseString(isClose int) (isCloseString string) {
	switch isClose {
	case 0:
		isCloseString = "已结课"
	case 1:
		isCloseString = "进行中"
	}
	return isCloseString
}
