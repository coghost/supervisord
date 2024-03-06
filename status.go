package supervisord

type (
	StateCode int
	StateName string
)

type State struct {
	Code StateCode `xmlrpc:"statecode"`
	Name StateName `xmlrpc:"statename"`
}

const (
	StateCodeFatal      StateCode = 2  // Supervisor has experienced a serious error.
	StateCodeRunning    StateCode = 1  // Supervisor is working normally.
	StateCodeRestarting StateCode = 0  // Supervisor is in the process of restarting.
	StateCodeShutdown   StateCode = -1 // Supervisor is in the process of shutting down.

	StateNameFatal      StateName = "FATAL"      // Supervisor has experienced a serious error.
	StateNameRunning    StateName = "RUNNING"    // Supervisor is working normally.
	StateNameRestarting StateName = "RESTARTING" // Supervisor is in the process of restarting.
	StateNameShutdown   StateName = "SHUTDOWN"   // Supervisor is in the process of shutting down.
)

func (c *Client) GetAPIVersion() (string, error) {
	return c.CallAsStr(getAPIVersion)
}

func (c *Client) GetSupervisorVersion() (string, error) {
	return c.CallAsStr(getSupervisorVersion)
}

func (c *Client) GetIdentification() (string, error) {
	return c.CallAsStr(getIdentification)
}

func (c *Client) GetState() (State, error) {
	var state State
	err := c.call(getState, nil, &state)

	return state, err
}

func (c *Client) GetPID() (int, error) {
	return c.CallAsInt(getPID)
}

func (c *Client) ReadLog(offset int64, length int) (string, error) {
	arg := []interface{}{offset, length}

	return c.CallAsStr(readLog, arg) // nolint
}

func (c *Client) ClearLog() error {
	return c.CallAsBool(clearLog)
}

func (c *Client) Shutdown() error {
	return c.CallAsBool(shutdown)
}

func (c *Client) Restart() error {
	return c.CallAsBool(restart)
}
