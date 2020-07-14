package step

import (
	"deploy/temp"
	"deploy/utils"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"sync"
)

/*
* 并发初始化服务器
 */
func SysInit(config *viper.Viper) {
	hosts := make([]interface{}, 0)
	for key, value := range config.Get("server").(map[string]interface{}) {
		if key == "app" || key == "lab" {
			for _, a := range []string{"etcd", "master", "node"} {
				i := value.(map[string]interface{})[a].(map[string]interface{})["ip"]
				hosts = append(hosts, i.([]interface{})...)
			}
			continue
		}
		h := value.(map[string]interface{})["ip"]
		hosts = append(hosts, h.([]interface{})...)
	}
	newHosts := RemoveRepeatElement(hosts)

	// 并发初始化服务器
	var wg sync.WaitGroup
	for _, h := range newHosts {
		wg.Add(1)
		go initHost(config, h, &wg)
	}
	wg.Wait()
}

/*
* 初始化服务器行为
 */
func initHost(c *viper.Viper, host interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	h := host.(string)
	log.Printf("*** Host Init %s", h)
	info := GetHostInfo(c, h)
	//fmt.Println(info.User, info.Port, info.Pass)
	utils.SFTPut(h, info, temp.Sysctl, "/etc/sysctl.conf")
	utils.SFTPut(h, info, temp.EPEL, "/etc/yum.repos.d/epel.repo")
	utils.SSHExec(h, info, "systemctl stop firewalld && "+
		"systemctl disable firewalld && "+
		"sed -i 's/SELINUX=.*/SELINUX=disabled/g' /etc/selinux/config && "+
		"mkdir -p /ddhome && "+
		"sysctl -p")
	domain := fmt.Sprintf("%s %s", c.Get("server.harbor").(map[string]interface{})["ip"].([]interface{})[0].(string),
		c.Get("registry.local").(map[string]interface{})["domain"].(string))
	utils.SSHExec(h, info,"echo "+domain+" >> /etc/hosts")
}

/*
* 去重
 */
func RemoveRepeatElement(slice []interface{}) (newSlice []interface{}) {
	mapArr := make(map[interface{}]int, 0)

	for i := 0; i < len(slice); i++ {
		if _, ok := mapArr[slice[i]]; ok {
			continue
		}
		mapArr[slice[i]] = i
	}

	for k, _ := range mapArr {
		newSlice = append(newSlice, k)
	}
	return
}

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