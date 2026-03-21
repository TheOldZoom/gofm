package utils

import "strconv"

func Commas(n int) string {
	sign := ""
	if n < 0 {
		sign = "-"
		n = -n
	}

	s := strconv.Itoa(n)
	if len(s) <= 3 {
		return sign + s
	}

	rem := len(s) % 3
	if rem == 0 {
		rem = 3
	}

	out := make([]byte, 0, len(s)+len(s)/3)
	out = append(out, sign...)
	out = append(out, s[:rem]...)

	for i := rem; i < len(s); i += 3 {
		out = append(out, ',')
		out = append(out, s[i:i+3]...)
	}

	return string(out)
}
