package models

import "github.com/tyrm/supreme-robot/config"

type Provider struct {
	client *Client
}

func (c *Client) ConfigProvider() *Provider {
	return &Provider{client: c}
}

func (p *Provider) Get(k string) (*string, error) {
	conf, err := p.client.ReadConfigByKey(k)
	if err != nil {
		s := ""
		return &s, err
	}
	if conf == nil {
		return config.Default(k), nil
	}
	return &conf.Value, nil
}

func (p *Provider) MGet(keys *[]string) (*map[string]*string, error) {
	confs, err := p.client.ReadConfigsByKeys(keys)
	if err != nil {
		return nil, err
	}

	confMap := make(map[string]*string)

	var (
		index int
		found bool
	)
	for _, key := range *keys {
		index = 0
		found = false
		for i, cc := range *confs {
			if key == cc.Key {
				index = i
				found = true
				break
			}
		}

		if found {
			confMap[key] = &(*confs)[index].Value
		} else {
			// use default
			confMap[key] = config.Default(key)
		}
	}

	return &confMap, nil
}

func (p *Provider) MSet(kv *map[string]string) error {
	configList := make(ConfigList, len(*kv))

	index := 0
	for k, v := range *kv {
		configList[index] = Config{
			Key:   k,
			Value: v,
		}
	}

	return configList.Upsert(p.client)
}

func (p *Provider) Set(k, v string) error {
	conf, err := p.client.ReadConfigByKey(k)
	if err != nil {
		return err
	}

	// create if it doesn't exist
	if conf == nil {
		newConf := Config{
			Key:   k,
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
	err = conf.Update(p.client)
	if err != nil {
		return err
	}
	return nil
}
