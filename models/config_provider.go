package models

import "github.com/tyrm/supreme-robot/config"

type Provider struct {
	client *Client
}

func (c *Client) ConfigProvider() *Provider {
	return &Provider{client: c}
}

func (p *Provider) Get(k string) (string, error) {
	conf, err := p.client.ReadConfigByKey(k)
	if err != nil {
		return "", err
	}
	if conf == nil {
		return "", config.ErrorNotDefined
	}
	return conf.Value, nil
}

func (p *Provider) Set(k, v string) error {
	conf, err := p.client.ReadConfigByKey(k)
	if err != nil {
		return err
	}

	// create if it doesn't exist
	if conf == nil {
		newConf := Config{
			Key: k,
			Value: v,
		}

		err = newConf.Create(p.client)
		if err != nil {
			return err
		}
		return nil
	}

	// update value
	conf.Value = v


	return nil
}