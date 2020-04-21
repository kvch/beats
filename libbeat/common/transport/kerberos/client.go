package kerberos

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	krbclient "gopkg.in/jcmturner/gokrb5.v7/client"
	krbconfig "gopkg.in/jcmturner/gokrb5.v7/config"
	"gopkg.in/jcmturner/gokrb5.v7/keytab"
	"gopkg.in/jcmturner/gokrb5.v7/spnego"

	"github.com/elastic/beats/v7/libbeat/logp"
)

var (
	InvalidAuthType = errors.New("invalid authentication type")
)

type Client struct {
	spClient *spnego.Client
}

func NewClient(config *Config, httpClient *http.Client, url string) (*Client, error) {
	var krbClient *krbclient.Client
	krbConf, err := krbconfig.Load(config.ConfigPath)
	if err != nil {
		return nil, fmt.Errorf("error creating Kerberos client: %+v", err)
	}

	switch config.AuthType {
	case AUTH_KEYTAB:
		kTab, err := keytab.Load(config.KeyTabPath)
		if err != nil {
			return nil, fmt.Errorf("cannot load keytab file %s: %+v", config.KeyTabPath, err)
		}
		krbClient = krbclient.NewClientWithKeytab(config.Username, config.Realm, kTab, krbConf)
	case AUTH_PASSWORD:
		krbClient = krbclient.NewClientWithPassword(config.Username, config.Realm, config.Password, krbConf)
	default:
		return nil, InvalidAuthType
	}

	spn := fmt.Sprintf("HTTP/%s@%s", url, config.Realm)
	return &Client{
		spClient: spnego.NewClient(krbClient, httpClient, spn),
	}, nil
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	logp.Info("kerberos client do")
	return c.spClient.Do(req)
}
