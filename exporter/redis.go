package exporter

import (
	"crypto/tls"
	"strings"

	"github.com/gomodule/redigo/redis"
	//log "github.com/sirupsen/logrus"
	log "github.com/golang/glog"
)

func (e *Exporter) connectToRedis() (redis.Conn, error) {
	options := []redis.DialOption{
		redis.DialConnectTimeout(e.options.ConnectionTimeouts),
		redis.DialReadTimeout(e.options.ConnectionTimeouts),
		redis.DialWriteTimeout(e.options.ConnectionTimeouts),

		redis.DialTLSConfig(&tls.Config{
			InsecureSkipVerify: e.options.SkipTLSVerification,
			Certificates:       e.options.ClientCertificates,
			RootCAs:            e.options.CaCertificates,
		}),
	}

	if e.options.User != "" {
		options = append(options, redis.DialUsername(e.options.User))
	}

	if e.options.Password != "" {
		options = append(options, redis.DialPassword(e.options.Password))
	}

	uri := e.redisAddr
	if !strings.Contains(uri, "://") {
		uri = "redis://" + uri
	}

	if e.options.PasswordMap[uri] != "" {
		options = append(options, redis.DialPassword(e.options.PasswordMap[uri]))
	}

	log.V(5).Infof("Trying DialURL(): %s", uri)
	c, err := redis.DialURL(uri, options...)
	if err != nil {
		log.V(5).Infof("DialURL() failed, err: %s", err)
		if frags := strings.Split(e.redisAddr, "://"); len(frags) == 2 {
			log.V(5).Infof("Trying: Dial(): %s %s", frags[0], frags[1])
			c, err = redis.Dial(frags[0], frags[1], options...)
		} else {
			log.V(5).Infof("Trying: Dial(): tcp %s", e.redisAddr)
			c, err = redis.Dial("tcp", e.redisAddr, options...)
		}
	}
	return c, err
}

func doRedisCmd(c redis.Conn, cmd string, args ...interface{}) (interface{}, error) {
	log.V(5).Infof("c.Do() - running command: %s %s", cmd, args)
	res, err := c.Do(cmd, args...)
	if err != nil {
		log.V(5).Infof("c.Do() - err: %s", err)
	}
	log.V(5).Infof("c.Do() - done")
	return res, err
}
