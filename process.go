package supervisord

import (
	"syscall"
)

type ProcessState int

type ProcessInfo struct {
	Name          string       `xmlrpc:"name"`           // Name of the process
	Group         string       `xmlrpc:"group"`          // Name of the process’ group
	Start         int          `xmlrpc:"start"`          // UNIX timestamp of when the process was started
	Stop          int          `xmlrpc:"stop"`           // UNIX timestamp of when the process last ended, or 0 if the process has never been stopped
	Now           int          `xmlrpc:"now"`            // UNIX timestamp of the current time, which can be used to calculate process up-time.
	State         ProcessState `xmlrpc:"state"`          // State code, see ProcessState.
	StateName     StateName    `xmlrpc:"statename"`      // String description of state
	SpawnErr      string       `xmlrpc:"spawnerr"`       // Description of error that occurred during spawn, or empty string if none
	ExitStatus    int          `xmlrpc:"exitstatus"`     // Exit status (errorlevel) of process, or 0 if the process is still running
	StdoutLogfile string       `xmlrpc:"stdout_logfile"` // Absolute path and filename to the STDOUT logfile
	StderrLogfile string       `xmlrpc:"stderr_logfile"` // Absolute path and filename to the STDOUT logfile
	Pid           int          `xmlrpc:"pid"`            // UNIX process ID (PID) of the process, or 0 if the process is not running
}

const (
	StateStopped  ProcessState = 0    // The process has been stopped due to a stop request or has never been started
	StateStarting ProcessState = 10   // The process is starting due to a start request
	StateRunning  ProcessState = 20   // The process is running
	StateBackoff  ProcessState = 30   // The process entered the StateStarting state but subsequently exited too quickly to move to the StateRunning state
	StateStopping ProcessState = 40   // The process is stopping due to a stop request
	StateExited   ProcessState = 100  // The process exited from the StateRunning state (expectedly or unexpectedly)
	StateFatal    ProcessState = 200  // The process could not be started successfully
	StateUnknown  ProcessState = 1000 // The process is in an unknown state (supervisord programming error)
)

func (c *Client) HandleAllProcesses(name CMD, args ...any) ([]ProcessInfo, error) {
	var piArr []ProcessInfo

	err := c.call(name, args, &piArr)

	return piArr, err
}

func (c *Client) GetProcessInfo(name string) (*ProcessInfo, error) {
	var processinfo ProcessInfo
	err := c.call(getProcessInfo, name, &processinfo)

	return &processinfo, err
}

func (c *Client) GetAllProcessInfo() ([]ProcessInfo, error) {
	return c.HandleAllProcesses(getAllProcessInfo)
}

type ProcessConfig struct {
	Autostart             bool   `xmlrpc:"autostart"`
	Command               string `xmlrpc:"command"`
	Directory             string `xmlrpc:"directory"`
	Exitcodes             []int  `xmlrpc:"exitcodes"`
	Group                 string `xmlrpc:"group"`
	GroupPrio             int    `xmlrpc:"group_prio"`
	Inuse                 bool   `xmlrpc:"inuse"`
	Killasgroup           bool   `xmlrpc:"killasgroup"`
	Name                  string `xmlrpc:"name"`
	ProcessPrio           int    `xmlrpc:"process_prio"`
	RedirectStderr        bool   `xmlrpc:"redirect_stderr"`
	Serverurl             string `xmlrpc:"serverurl"`
	Startretries          int    `xmlrpc:"startretries"`
	Startsecs             int    `xmlrpc:"startsecs"`
	StderrCaptureMaxbytes int    `xmlrpc:"stderr_capture_maxbytes"`
	StderrEventsEnabled   bool   `xmlrpc:"stderr_events_enabled"`
	StderrLogfile         string `xmlrpc:"stderr_logfile"`
	StderrLogfileBackups  int    `xmlrpc:"stderr_logfile_backups"`
	StderrLogfileMaxbytes int    `xmlrpc:"stderr_logfile_maxbytes"`
	StderrSyslog          bool   `xmlrpc:"stderr_syslog"`
	StdoutCaptureMaxbytes int    `xmlrpc:"stdout_capture_maxbytes"`
	StdoutEventsEnabled   bool   `xmlrpc:"stdout_events_enabled"`
	StdoutLogfile         string `xmlrpc:"stdout_logfile"`
	StdoutLogfileBackups  int    `xmlrpc:"stdout_logfile_backups"`
	StdoutLogfileMaxbytes int    `xmlrpc:"stdout_logfile_maxbytes"`
	StdoutSyslog          bool   `xmlrpc:"stdout_syslog"`
	Stopsignal            int    `xmlrpc:"stopsignal"`
	Stopwaitsecs          int    `xmlrpc:"stopwaitsecs"`
	UID                   int    `xmlrpc:"uid"`
}

func (c *Client) GetAllConfigInfo() ([]ProcessConfig, error) {
	var piArr []ProcessConfig
	err := c.call(getAllConfigInfo, nil, &piArr)

	return piArr, err
}

func (c *Client) StartProcess(name string, wait bool) error {
	return c.CallAsBool(startProcess, name, wait)
}

func (c *Client) StartAllProcesses(wait bool) ([]ProcessInfo, error) {
	return c.HandleAllProcesses(startAllProcesses, wait)
}

func (c *Client) StartProcessGroup(name string, wait bool) ([]ProcessInfo, error) {
	return c.HandleAllProcesses(startProcessGroup, name, wait)
}

func (c *Client) StopProcess(name string, wait bool) error {
	return c.CallAsBool(stopProcess, name, wait)
}

func (c *Client) StopProcessGroup(name string, wait bool) ([]ProcessInfo, error) {
	return c.HandleAllProcesses(stopProcessGroup, name, wait)
}

func (c *Client) StopAllProcesses(wait bool) ([]ProcessInfo, error) {
	return c.HandleAllProcesses(stopAllProcesses, wait)
}

func (c *Client) SignalProcess(name string, signal syscall.Signal) error {
	return c.CallAsBool(signalProcess, name, int(signal))
}

func (c *Client) SignalProcessGroup(signal syscall.Signal) ([]ProcessInfo, error) {
	return c.HandleAllProcesses(signalProcessGroup, int(signal))
}

func (c *Client) SignalAllProcesses(signal syscall.Signal) ([]ProcessInfo, error) {
	return c.HandleAllProcesses(signalAllProcesses, int(signal))
}

func (c *Client) SendProcessStdin(name string, chars string) error {
	return c.CallAsBool(sendProcessStdin, name, chars)
}

func (c *Client) SendRemoteCommEvent(eventHeader string, eventBody string) error {
	return c.CallAsBool(sendRemoteCommEvent, eventHeader, eventBody)
}

func (c *Client) ReloadConfig() ([]interface{}, error) {
	return c.CallAsInterfaceArray(reloadConfig)
}

func (c *Client) AddProcessGroup(name string) error {
	return c.CallAsBool(addProcessGroup, name)
}

func (c *Client) RemoveProcessGroup(name string) error {
	return c.CallAsBool(removeProcessGroup, name)
}
