package exporter

import (
	"encoding/json"
	"io/ioutil"

	//log "github.com/sirupsen/logrus"
	log "github.com/golang/glog"
)

// LoadPwdFile reads the redis password file and returns the password map
func LoadPwdFile(passwordFile string) (map[string]string, error) {
	passwordMap := make(map[string]string)

	log.V(5).Infof("start load password file: %s", passwordFile)
	bytes, err := ioutil.ReadFile(passwordFile)
	if err != nil {
		log.Warningf("load password file failed: %s", err)
		return nil, err
	}
	err = json.Unmarshal(bytes, &passwordMap)
	if err != nil {
		log.Warningf("password file format error: %s", err)
		return nil, err
	}
	return passwordMap, nil
}
