package service

import (
	"math/rand"
	"time"
)

var Comm = commService{}

type commService struct{}

func (s *commService) ParseEnvs(envs map[string]interface{}) map[string]string {
	var new_envs map[string]string = make(map[string]string)
	for k, v := range envs {
		switch v := v.(type) {
		case string:
			new_envs[k] = v
		case []interface{}:
			new_v := ""
			for _, u := range v {
				if new_v == "" {
					new_v = u.(string)
					continue
				}
				new_v = new_v + "," + u.(string)
			}
			new_envs[k] = new_v
		}
	}
	return new_envs
}

func (s *commService) RandSeq(randlen int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, randlen)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// func (s *commService) GetScriptBody(script_name string) string {
// 	script_body, err := dao.CicdScript.Fields("script_body").Where("script_name=", script_name).Value()
// 	if err != nil {
// 		g.Log().Error(err)
// 	}
// 	return script_body.String()
// }
