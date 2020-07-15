package temp

import "github.com/spf13/viper"

/*
*  根据IP地址获取用户名密码信息
 */
type HostInfo struct {
	User string
	Pass string
	Port int64
}

func GetHostInfo(config *viper.Viper, host string) HostInfo {
	hostinfo := HostInfo{}
	for k, h := range config.Get("server").(map[string]interface{}) {
		if k == "app" || k == "lab" {
			for _, a := range []string{"etcd", "master", "node"} {
				i := h.(map[string]interface{})[a].(map[string]interface{})["ip"]
				for _, j := range i.([]interface{}) {
					if j.(string) == host {
						hostinfo.User = h.(map[string]interface{})[a].(map[string]interface{})["username"].(string)
						hostinfo.Pass = h.(map[string]interface{})[a].(map[string]interface{})["password"].(string)
						hostinfo.Port = h.(map[string]interface{})[a].(map[string]interface{})["port"].(int64)
						return hostinfo
					}
				}
			}
			continue
		}

		i := h.(map[string]interface{})["ip"]
		for _, j := range i.([]interface{}) {
			if j.(string) == host {
				hostinfo.User = h.(map[string]interface{})["username"].(string)
				hostinfo.Pass = h.(map[string]interface{})["password"].(string)
				hostinfo.Port = h.(map[string]interface{})["port"].(int64)
				return hostinfo
			}
		}

	}
	return hostinfo
}
