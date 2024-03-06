package supervisord

type CmdCall struct {
	MethodName string        `xmlrpc:"methodName"`
	Params     []interface{} `xmlrpc:"params"`
}

func (c *Client) ListMethods() ([]string, error) {
	return c.CallAsStrArray(listMethods)
}

func (c *Client) MethodHelp(cmd CMD) (string, error) {
	name := c.refineCmd(cmd)
	return c.CallAsStr(methodHelp, name)
}

func (c *Client) MethodSignature(cmd CMD) ([]interface{}, error) {
	name := c.refineCmd(cmd)
	return c.CallAsInterfaceArray(methodSignature, name)
}

// Multicall
// Process an array of calls, and return an array of results. Calls should be structs of the form {‘methodName’: string, ‘params’: array}. Each result will either be a single-item array containing the result value, or a struct of the form {‘faultCode’: int, ‘faultString’: string}. This is useful when you need to make lots of small calls without lots of round trips.
//
// Example:
//
//	calls := []CmdCall{
//		{
//			MethodName: "supervisor.getPID",
//			Params:     []interface{}{},
//		},
//		{
//			MethodName: "system.methodHelp",
//			Params:     []interface{}{c.refineCmd(methodHelp)},
//		},
//	}
func (c *Client) Multicall(calls []CmdCall) ([]interface{}, error) {
	return c.CallAsInterfaceArray(multicall, calls)
}
