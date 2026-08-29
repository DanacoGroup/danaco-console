package mail

import (
	"html"
	"strings"
	"unicode/utf8"
)

// maxUntrustedLen to długość, do której skracana jest wartość spoza rdzenia.
const maxUntrustedLen = 120

// cleanUntrusted usuwa znaki sterujące i skraca wartość. Wartość trafia także
// do tematu, a złamanie wiersza w nagłówku pozwala dopisać własne nagłówki.
func cleanUntrusted(v string) string {
	var b strings.Builder
	b.Grow(len(v))
	spacja := false
	for _, r := range v {
		switch {
		case r == '\r' || r == '\n' || r == '\t' || r == ' ':
			// Ciąg białych znaków zwijamy do jednej spacji — CRLF dałby dwie.
			spacja = true
		case r < 0x20 || r == 0x7f || r == utf8.RuneError:
		default:
			if spacja && b.Len() > 0 {
				b.WriteByte(' ')
			}
			spacja = false
			b.WriteRune(r)
		}
	}
	return truncateRunes(b.String(), maxUntrustedLen)
}

// escapeForHTML dokłada ucieczkę znaczników. Wyłącznie dla części HTML —
// w części tekstowej dałaby odbiorcy widoczne encje &amp; i &quot;.
func escapeForHTML(v string) string {
	return html.EscapeString(v)
}

// truncateRunes skraca do n znaków, nie bajtów — cięcie po bajtach rozbiłoby
// polski znak diakrytyczny na końcu.
func truncateRunes(v string, n int) string {
	if utf8.RuneCountInString(v) <= n {
		return v
	}
	count := 0
	for i := range v {
		if count == n {
			return v[:i]
		}
		count++
	}
	return v
}
