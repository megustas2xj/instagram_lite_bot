// Please let author have a drink, usdt trc20: TEpSxaE3kexE4e5igqmCZRMJNoDiQeWx29
// tg: @fuckins996
package google

import "net/http"

type gmsManger struct {
}

func (this *gmsManger) GetClient(email string) *GmsClient {
	return &GmsClient{
		host: "http://192.168.123.241:10089",
		h:    &http.Client{},
	}
}

var GmsManger gmsManger
