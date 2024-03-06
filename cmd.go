package supervisord

type CMD string

const (
	_prefix       = "supervisor."
	_prefixSystem = "system."
)

const (
	addProcessGroup      CMD = "addProcessGroup"
	clearAllProcessLogs  CMD = "clearAllProcessLogs"
	clearLog             CMD = "clearLog"
	clearProcessLog      CMD = "clearProcessLog"
	clearProcessLogs     CMD = "clearProcessLogs"
	getAPIVersion        CMD = "getAPIVersion"
	getAllConfigInfo     CMD = "getAllConfigInfo"
	getAllProcessInfo    CMD = "getAllProcessInfo"
	getIdentification    CMD = "getIdentification"
	getPID               CMD = "getPID"
	getProcessInfo       CMD = "getProcessInfo"
	getState             CMD = "getState"
	getSupervisorVersion CMD = "getSupervisorVersion"
	getVersion           CMD = "getVersion"
	readLog              CMD = "readLog"
	readMainLog          CMD = "readMainLog"
	readProcessLog       CMD = "readProcessLog"
	readProcessStderrLog CMD = "readProcessStderrLog"
	readProcessStdoutLog CMD = "readProcessStdoutLog"
	reloadConfig         CMD = "reloadConfig"
	removeProcessGroup   CMD = "removeProcessGroup"
	restart              CMD = "restart"
	sendProcessStdin     CMD = "sendProcessStdin"
	sendRemoteCommEvent  CMD = "sendRemoteCommEvent"
	shutdown             CMD = "shutdown"
	signalAllProcesses   CMD = "signalAllProcesses"
	signalProcess        CMD = "signalProcess"
	signalProcessGroup   CMD = "signalProcessGroup"
	startAllProcesses    CMD = "startAllProcesses"
	startProcess         CMD = "startProcess"
	startProcessGroup    CMD = "startProcessGroup"
	stopAllProcesses     CMD = "stopAllProcesses"
	stopProcess          CMD = "stopProcess"
	stopProcessGroup     CMD = "stopProcessGroup"
	tailProcessLog       CMD = "tailProcessLog"
	tailProcessStderrLog CMD = "tailProcessStderrLog"
	tailProcessStdoutLog CMD = "tailProcessStdoutLog"
)

const (
	listMethods     CMD = "system.listMethods"
	methodHelp      CMD = "system.methodHelp"
	methodSignature CMD = "system.methodSignature"
	multicall       CMD = "system.multicall"
)
