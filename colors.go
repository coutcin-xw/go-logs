package logs

//颜色部分:

func Black(s string) string {
	return "\033[0;30m" + s + "\033[0m"
}

func BlackBold(s string) string {
	return "\033[1;30m" + s + "\033[0m"
}

func BlackLine(s string) string {
	return "\033[4;30m" + s + "\033[0m"
}

func Red(s string) string {
	return "\033[0;31m" + s + "\033[0m"
}

func RedBold(s string) string {
	return "\033[1;31m" + s + "\033[0m"
}

func RedLine(s string) string {
	return "\033[4;31m" + s + "\033[0m"
}

func Green(s string) string {
	return "\033[0;32m" + s + "\033[0m"
}

func GreenBold(s string) string {
	return "\033[1;32m" + s + "\033[0m"
}

func GreenLine(s string) string {
	return "\033[4;32m" + s + "\033[0m"
}

func Yellow(s string) string {
	return "\033[0;33m" + s + "\033[0m"
}

func YellowBold(s string) string {
	return "\033[1;33m" + s + "\033[0m"
}

func YellowLine(s string) string {
	return "\033[4;33m" + s + "\033[0m"
}

func Blue(s string) string {
	return "\033[0;34m" + s + "\033[0m"
}

func BlueBold(s string) string {
	return "\033[1;34m" + s + "\033[0m"
}

func BlueLine(s string) string {
	return "\033[4;34m" + s + "\033[0m"
}

func Purple(s string) string {
	return "\033[0;35m" + s + "\033[0m"
}

func PurpleBold(s string) string {
	return "\033[1;35m" + s + "\033[0m"
}

func PurpleLine(s string) string {
	return "\033[4;35m" + s + "\033[0m"
}

func Cyan(s string) string {
	return "\033[0;36m" + s + "\033[0m"
}

func CyanBold(s string) string {
	return "\033[1;36m" + s + "\033[0m"
}

func CyanLine(s string) string {
	return "\033[4;36m" + s + "\033[0m"
}

func White(s string) string {
	return "\033[0;37m" + s + "\033[0m"
}

func WhiteBold(s string) string {
	return "\033[1;37m" + s + "\033[0m"
}

func WhiteLine(s string) string {
	return "\033[4;37m" + s + "\033[0m"
}
