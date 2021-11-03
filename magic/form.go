package magic

import "strings"

var abbrs = []string{
	"ACL", "API", "ASCII",
	"CPU", "CSS",
	"DNS",
	"EOF",
	"GUID",
	"HTML", "HTTP", "HTTPS",
	"ID",
	"VIP",
	"IP",
	"JSON",
	"LHS",
	"QPS",
	"RAM", "RHS", "RPC",
	"SLA", "SMTP", "SQL", "SSH",
	"TCP", "TLS", "TTL",
	"UDP", "UI", "UID", "UUID", "URI", "URL", "UTF8",
	"VM",
	"XML", "XMPP", "XSRF", "XSS",
}

var abbrMap = make(map[string]string)

func init() {
	for _, abbr := range abbrs {
		abbrMap[camelize(abbr)] = abbr
	}
}

func Standardize(s string, sep SeparatorType) string {
	if s == "" {
		return s
	}

	b := []byte{}
	words := strings.Split(s, sep)
	for _, word := range words {
		word = camelize(word)
		abbr, ok := abbrMap[word]
		if ok {
			word = abbr
		}
		b = append(b, []byte(word)...)
	}
	return string(b)
}

func camelize(s string) string {
	s = strings.ToLower(s)
	b := []byte(s)
	if b[0] >= 'a' && b[0] <= 'z' {
		b[0] -= 32
	}
	return string(b)
}
