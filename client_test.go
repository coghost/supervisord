package supervisord

import (
	"fmt"
	"testing"

	"github.com/coghost/xlog"
	"github.com/k0kubun/pp/v3"
	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/suite"
	"github.com/ungerik/go-dry"
)

type ClientSuite struct {
	suite.Suite
	host   string
	client *Client
}

func TestClient(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}

func (s *ClientSuite) SetupSuite() {
	xlog.InitLogDebug()
	type config struct {
		Url  string `toml:"url"`
		User string `toml:"username"`
		Pwd  string `toml:"password"`
	}

	raw, err := dry.FileGetBytes(".auth.toml")
	s.Nil(err)

	var cfg config

	toml.Unmarshal(raw, &cfg)

	s.host = fmt.Sprintf("http://%s/RPC2", cfg.Url)
	c, err := NewClient(s.host, WithAuth(cfg.User, cfg.Pwd))
	fmt.Printf("%s\n", c)
	s.Nil(err)
	s.client = c
}

func (s *ClientSuite) TearDownSuite() {
}

func (s *ClientSuite) Test_01_getversion() {
	a, b := "3.0", "4.2.5"
	v, e := s.client.GetAPIVersion()
	s.Nil(e)
	pp.Println(v)
	s.Equal(v, a)

	v1, e := s.client.GetSupervisorVersion()
	s.Nil(e)

	pp.Println(v1)
	s.Equal(v1, b)
}

func (s *ClientSuite) Test_02_getState() {
	st, err := s.client.GetState()
	s.Nil(err)

	pp.Println(st)

	pid, err := s.client.GetPID()
	s.Nil(err)
	pp.Println(pid)
}

func (s *ClientSuite) Test_getLog() {
	raw, err := s.client.ReadLog(-1024, 0)
	s.Nil(err)
	pp.Println(raw)
}

func (s *ClientSuite) Test_listMethods() {
	raw, err := s.client.ListMethods()
	s.Nil(err)
	pp.Println(raw)
}

func (s *ClientSuite) Test_03_getInfo() {
	name := "gocc-site2"
	pi, err := s.client.GetProcessInfo(name)
	s.Nil(err)
	pp.Println(pi)

	piarr, err := s.client.GetAllProcessInfo()
	s.Nil(err)

	pp.Println(piarr)

	str, err := s.client.GetAllConfigInfo()
	s.Nil(err)
	pp.Println(str)

	s.client.StartProcess(name, true)
	s.client.StopProcess(name, true)
}

func (s *ClientSuite) Test_04_tmp() {
	str, err := s.client.MethodHelp(tailProcessLog)
	s.Nil(err)
	fmt.Println(str)

	arr, err := s.client.MethodSignature(tailProcessLog)
	s.Nil(err)
	pp.Println(len(arr), arr)
}

func (s *ClientSuite) Test_05_multi() {
	c := s.client

	calls := []CmdCall{
		{
			MethodName: "supervisor.getPID",
			Params:     []interface{}{},
		},
		{
			MethodName: "system.methodHelp",
			Params:     []interface{}{c.refineCmd(methodHelp)},
		},
	}

	arr, err := c.Multicall(calls)
	s.Nil(err)

	for i, res := range arr {
		fmt.Printf("%s\n%v\n", calls[i].MethodName, res)
	}
}

func (s *ClientSuite) Test_06_readlog() {
	_, err := s.client.ReadLog(int64(0), 2024)
	s.Nil(err)
}
