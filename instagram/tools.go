// Please let author have a drink, usdt trc20: TEpSxaE3kexE4e5igqmCZRMJNoDiQeWx29
// tg: @fuckins996
package instagram

import url2 "net/url"
import "strings"

func DecodeShareUrl(url string) string {
	parse, err := url2.Parse(url)
	if err != nil {
		return ""
	}
	username := ""
	if strings.Index(parse.Path, "/") == 0 {
		username = parse.Path[1:]
	}
	return username
}
