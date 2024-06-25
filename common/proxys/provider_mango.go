// Please let author have a drink, usdt trc20: TEpSxaE3kexE4e5igqmCZRMJNoDiQeWx29
// tg: @fuckins996
package proxys

import (
	"CentralizedControl/common/log"
	"CentralizedControl/common/utils"
	"fmt"
)

type MangoProxy struct {
	username string
	password string
	dns      string
	dnsPort  string
	time     string
}

func CreateMangoProxy() *MangoProxy {
	return &MangoProxy{
		username: "-zone-custom",
		password: "",
		dns:      "43.152.113.55",
		dnsPort:  "2333",
		time:     "15",
	}
}

func (this *MangoProxy) GetProxy(region string, asn string) Proxy {
	proxyUrl := fmt.Sprintf("socks5://a6a9a0f0331d1f9598335-zone--sto469323-region-%s-asn-%s-session-%s-sessTime-5:@43.153.237.55:2333",
		region, asn, utils.GenString(utils.CharSet_abc, 9))
	log.Info("%s", proxyUrl)
	proxy, err := CreateSocks5Proxy(proxyUrl)
	if err != nil {
		log.Warn("create s5 error: %v", err)
		return nil
	}
	return proxy
}
