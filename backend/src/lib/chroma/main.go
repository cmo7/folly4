package chroma

var ansiColorMap = map[string]string{}
var ansiFormatMap = map[string]string{}

func init() {
	ansiColorMap["black"] = "30"
	ansiColorMap["red"] = "31"
	ansiColorMap["green"] = "32"
	ansiColorMap["yellow"] = "33"
	ansiColorMap["blue"] = "34"
	ansiColorMap["magenta"] = "35"
	ansiColorMap["cyan"] = "36"
	ansiColorMap["white"] = "37"

	ansiFormatMap["reset"] = "0"
	ansiFormatMap["bold"] = "1"
	ansiFormatMap["italic"] = "3"
	ansiFormatMap["underline"] = "4"
	ansiFormatMap["blink"] = "5"
	ansiFormatMap["inverse"] = "7"
	ansiFormatMap["hidden"] = "8"
}

func getColorCode(color string) string {
	code, ok := ansiColorMap[color]
	if !ok {
		code = "39"
	}
	return "\033[" + code + "m"
}

func getFormatCode(format string) string {
	code, ok := ansiFormatMap[format]
	if !ok {
		code = "0"
	}
	return "\033[" + code + "m"
}

func Format(s string, funcs ...func(string) string) string {
	for _, f := range funcs {
		s = f(s)
	}
	return s
}

func Color(color string) func(string) string {
	return func(s string) string {
		return getColorCode(color) + s + getFormatCode("reset")
	}
}

func Bold(s string) string {
	return getFormatCode("bold") + s + getFormatCode("reset")
}

func Italic(s string) string {
	return getFormatCode("italic") + s + getFormatCode("reset")
}

func Underline(s string) string {
	return getFormatCode("underline") + s + getFormatCode("reset")
}

func Blink(s string) string {
	return getFormatCode("blink") + s + getFormatCode("reset")
}

func Inverse(s string) string {
	return getFormatCode("inverse") + s + getFormatCode("reset")
}

func Hidden(s string) string {
	return getFormatCode("hidden") + s + getFormatCode("reset")
}
