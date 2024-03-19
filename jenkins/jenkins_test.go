package jenkins

import (
	"testing"
)

var (
	user        string   = "zhangwanfeng"
	pass        string   = "zhangwanfeng"
	allServices []string = []string{
		"game-service",
		"money-service",
		"quest-service",
		"gametree-service",
		"account-service",
		"stage-service",
		"antiaddiction-service",
		"item-service",
		"team-service",
		"box-service",
		"login-service",
		"casualgame-service",
		"match-service",
		"user-service",
	}
	allGateway = []string{
		"sapi-gw",
		"papi-gw",
		"casualgame-gw",
	}
)

func TestJenKins_CreateSngService(t *testing.T) {
	NewJenkins(user, pass).getJob("http://pink-jenkins.kaiqitech.com/job/gateway/job/sapi-gw/job/test/")
	// if err := NewJenkins(user, pass).CreateSngService("account-service", "user-service", "quest-service"); err != nil {
	// 	t.Errorf(err.Error())
	// }
}

func TestJenKins_UpdateSngService(t *testing.T) {
	NewJenkins(user, pass).UpdateSngService(allServices...)
}

func TestJenKins_BuildSngServiceDev(t *testing.T) {
	NewJenkins(user, pass).BuildSngServiceDev(allServices...)
}

func TestJenKins_BuildSngServiceTest(t *testing.T) {
	NewJenkins(user, pass).BuildSngServiceTest(allServices...)
}

func TestJenKins_CreateSngGateway(t *testing.T) {
	if err := NewJenkins(user, pass).CreateSngGateway(allGateway...); err != nil {
		t.Errorf(err.Error())
	}
}

func TestJenKins_UpdateSngGateway(t *testing.T) {
	NewJenkins(user, pass).UpdateSngGateway(allGateway...)
}
