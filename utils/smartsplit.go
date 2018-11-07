package utils

const (
	stateParsing = false
	stateEscaped = true
)

// Split strings
// input := "a,b\,c"
// del := ","
// esc := "\\"
//output: ["a","b,c"]
func Split(input, escape, delimiter string) []string {
	state := stateParsing
	found := make([]string, 0)
	parsed := ""

	for _, c := range input {
		c := string(c)
		if state == stateParsing {
			if c == delimiter {
				found = append(found, parsed)
				parsed = ""
			} else if c == escape {
				state = stateEscaped
			} else {
				parsed += c
			}
		} else {
			parsed += c
			state = stateParsing
		}
	}

	if parsed != "" {
		found = append(found, parsed)
	}

	return found
}

// Replace ("aa,bb", "\\", ",", "\n")
func Replace(input, escape, delimiter, new string) string {
	out := ""
	state := stateParsing

	for _, c := range input {
		c := string(c)
		// 1. 默认
		if state == stateParsing {
			if c == delimiter { // "," 替换
				out += new
			} else if c == escape { //转义
				state = stateEscaped
			} else {
				out += c
			}
		} else {
			//启动转义
			if c == delimiter {
				out += c
			} else {
				out += escape
				out += c
			}
			state = stateParsing
		}
	}

	return out
}
