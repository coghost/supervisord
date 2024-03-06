package supervisord

import (
	"errors"
	"fmt"
	"strings"

	"github.com/kolo/xmlrpc"
	"github.com/rs/zerolog/log"
)

type Client struct {
	*xmlrpc.Client

	host     string
	username string
	password string

	must  bool
	debug bool
}

var ErrorReturnedFalse = errors.New("Call returned false")

type ClientOptions func(*Client)

func bindOptions(opt *Client, opts ...ClientOptions) {
	for _, f := range opts {
		f(opt)
	}
}

func WithAuth(username, password string) ClientOptions {
	return func(o *Client) {
		o.username = username
		o.password = password
	}
}

func WithMust(b bool) ClientOptions {
	return func(o *Client) {
		o.must = b
	}
}

func NewClient(url string, opts ...ClientOptions) (*Client, error) {
	opt := &Client{}
	bindOptions(opt, opts...)

	tr := newBasicAuth(opt.username, opt.password)

	rpc, err := xmlrpc.NewClient(url, tr)
	if err != nil {
		return nil, err
	}

	return &Client{
		Client:   rpc,
		host:     url,
		username: opt.username,
		password: opt.password,
		must:     opt.must,
		debug:    true,
	}, nil
}

func (c *Client) String() string {
	return fmt.Sprintf("[%s] %s, %t, %t", c.host, c.username, c.must, c.debug)
}

func (c *Client) CallAsStr(serviceMethod CMD, args ...any) (string, error) {
	var reply string
	err := c.call(serviceMethod, c.refineArgs(args...), &reply)

	return reply, err
}

func (c *Client) CallAsStrArray(serviceMethod CMD, args ...any) ([]string, error) {
	var reply []string
	err := c.call(serviceMethod, nil, &reply)

	return reply, err
}

func (c *Client) CallAsInt(serviceMethod CMD, args ...any) (int, error) {
	var reply int
	err := c.call(serviceMethod, c.refineArgs(args...), &reply)

	return reply, err
}

func (c *Client) CallAsBool(serviceMethod CMD, args ...any) error {
	var reply bool
	err := c.call(serviceMethod, c.refineArgs(args...), &reply)

	if !reply {
		return ErrorReturnedFalse
	}

	return err
}

func (c *Client) CallAsInterface(cmdIn CMD, args ...any) (interface{}, error) {
	var reply interface{}
	err := c.call(cmdIn, args, &reply)

	return reply, err
}

func (c *Client) CallAsInterfaceArray(cmdIn CMD, args ...any) ([]interface{}, error) {
	var arr []interface{}
	err := c.call(cmdIn, args, &arr)

	return arr, err
}

func (c *Client) call(serviceMethod CMD, args any, reply any) error {
	err := c.Call(c.refineCmd(serviceMethod), args, reply)
	return c.pie(err)
}

func (c *Client) pie(err error) error {
	if !c.must {
		return err
	}

	if err != nil {
		panic(err)
	}

	return nil
}

func (c *Client) refineArgs(args ...any) any {
	if len(args) == 0 {
		return nil
	}

	if c.debug {
		log.Debug().Interface("args", args).Msg("got args")
	}

	return args
}

func (c *Client) refineCmd(rawCmd CMD) string {
	cmd := string(rawCmd)
	if strings.HasPrefix(cmd, _prefixSystem) {
		return cmd
	}

	if !strings.HasPrefix(cmd, _prefix) {
		cmd = _prefix + cmd
	}

	return cmd
}

func genCmd(cmdIn CMD, _typ string) {
	raw := `
	func (c *Client) %s() (%s, error) {
		return c.<>(%s)
	}
	`

	cmd := string(cmdIn)
	cmd1 := strings.ToUpper(cmd[0:1]) + cmd[1:]
	fmt.Printf(raw, cmd1, _typ, string(cmd)) // nolint
}
