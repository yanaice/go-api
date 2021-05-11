package auth

import (
	"io/ioutil"
	"net/http"
	"time"

	"go-starter-project/pkg/config"
	"go-starter-project/pkg/derror"
	"go-starter-project/pkg/dhttp"
)

type AuthTokenAPI interface {
	CheckToken(token string) error
}

type authTokenAPIImpl struct {
	base   string
	client dhttp.Client
}

func AuthTokenAPIInit() AuthTokenAPI {
	return &authTokenAPIImpl{
		base: config.Conf.AuthTokenAPI.BaseURL,
		client: &dhttp.DiancaiClient{
			Client: &http.Client{
				Timeout: 2 * time.Second,
			},
		},
	}
}

func (api *authTokenAPIImpl) CheckToken(token string) error {
	url := api.base + "/token/check"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return derror.E(err).SetDebug("check token via API failed: cannot build request").Log()
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := api.client.Do(req)
	if err != nil {
		return derror.E(err).SetDebug("check token via API failed: cannot do request").Log()
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return derror.E(err).SetDebug("check token via API failed: cannot read response body").Log()
	}

	if resp.StatusCode != http.StatusOK {
		derr, err := derror.UnmarshalDerrorReply(body)
		if err != nil {
			return derror.E(err).SetDebug("check token via API failed: cannot unmarshal error").Log()
		}

		return derror.WrapReply(derr).SetDebug("check token via API failed: non-OK response").Log()
	}

	return nil
}
