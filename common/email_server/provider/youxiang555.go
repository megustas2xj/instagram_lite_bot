// Please let author have a drink, usdt trc20: TEpSxaE3kexE4e5igqmCZRMJNoDiQeWx29
// tg: @fuckins996
package provider

import (
	"CentralizedControl/common"
	"CentralizedControl/common/email_server/base"
	"CentralizedControl/common/http_helper"
	"CentralizedControl/common/log"
	"net/http"
	"strings"
)

//https://www.youxiang555.com/api/pub/email.html?type=33&token=f822649b9102bb58cfc3d0506e55a58c&count=1
//tgrftw32323
//tgrftw32323

type ProviderYouXiang555 struct {
	Provider
	client *http.Client
	config *base.Config
	url    string
}

func (this *ProviderYouXiang555) GetType() string {
	return ProviderYouXiang555Name
}

func (this *ProviderYouXiang555) GetEmail() (base.EmailInterface, error) {
	var resp string
	var err error
	for retry := 0; retry < 3; retry++ {
		resp, err = http_helper.HttpDo(this.client, &http_helper.RequestOpt{
			IsPost: false,
			ReqUrl: this.url,
		})
		if err != nil {
			log.Error("yx1024 get email error: %v", err)
			continue
		} else {
			err = nil
			break
		}
	}
	if err != nil {
		return nil, err
	}

	resp = strings.ReplaceAll(resp, "<br>", "")
	sp := strings.Split(resp, "----")
	if len(sp) != 2 {
		return nil, common.NerError(resp)
	}
	log.Info("youxiang555 get email: %s, passwd: %s", sp[0], sp[1])
	return base.CreateImapEmail(this.config, &base.Account{
		Username: sp[0],
		Password: sp[1],
	}), nil
}

func CreateEmailProviderYouXiang555(link string, retry *base.RetryConfig) Provider {
	return &ProviderYouXiang555{
		url:    link,
		client: http_helper.CreateGoHttpClient(http_helper.HttpTimeout(20)),
		config: &base.Config{
			Server:      "outlook.office365.com",
			Port:        "993",
			RetryConfig: retry,
		},
	}
}
