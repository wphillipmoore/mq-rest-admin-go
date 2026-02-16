package mqrestadmin

import "context"

// CommandOption configures optional parameters for MQSC command methods.
type CommandOption func(*commandConfig)

type commandConfig struct {
	requestParameters  map[string]any
	responseParameters []string
	where              *string
}

func buildCommandConfig(opts []CommandOption) commandConfig {
	var config commandConfig
	for _, opt := range opts {
		opt(&config)
	}
	return config
}

// WithRequestParameters sets the MQSC command parameters (attributes to set
// or filter on).
func WithRequestParameters(params map[string]any) CommandOption {
	return func(config *commandConfig) {
		config.requestParameters = params
	}
}

// WithResponseParameters specifies which attributes to return in the response.
// For DISPLAY commands, this defaults to ["all"] if not specified.
func WithResponseParameters(params []string) CommandOption {
	return func(config *commandConfig) {
		config.responseParameters = params
	}
}

// WithWhere sets a WHERE clause to filter DISPLAY command results.
func WithWhere(clause string) CommandOption {
	return func(config *commandConfig) {
		config.where = &clause
	}
}

// ---------------------------------------------------------------------------
// DISPLAY commands — wildcard default (name defaults to "*")
// ---------------------------------------------------------------------------

// DisplayQueue displays queue attributes. Name defaults to "*" if empty.
func (session *Session) DisplayQueue(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	config := buildCommandConfig(opts)
	if name == "" {
		name = "*"
	}
	return session.mqscCommand(ctx, "DISPLAY", "QUEUE", &name, config.requestParameters, config.responseParameters, config.where, true)
}

// DisplayChannel displays channel attributes. Name defaults to "*" if empty.
func (session *Session) DisplayChannel(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	config := buildCommandConfig(opts)
	if name == "" {
		name = "*"
	}
	return session.mqscCommand(ctx, "DISPLAY", "CHANNEL", &name, config.requestParameters, config.responseParameters, config.where, true)
}

// ---------------------------------------------------------------------------
// DISPLAY commands — singleton (no name, returns single object or nil)
// ---------------------------------------------------------------------------

// DisplayQmgr displays the queue manager attributes.
func (session *Session) DisplayQmgr(ctx context.Context, opts ...CommandOption) (map[string]any, error) {
	config := buildCommandConfig(opts)
	objects, err := session.mqscCommand(ctx, "DISPLAY", "QMGR", nil, config.requestParameters, config.responseParameters, nil, true)
	if err != nil {
		return nil, err
	}
	if len(objects) == 0 {
		return nil, nil
	}
	return objects[0], nil
}

// DisplayQmstatus displays the queue manager status.
func (session *Session) DisplayQmstatus(ctx context.Context, opts ...CommandOption) (map[string]any, error) {
	config := buildCommandConfig(opts)
	objects, err := session.mqscCommand(ctx, "DISPLAY", "QMSTATUS", nil, config.requestParameters, config.responseParameters, nil, true)
	if err != nil {
		return nil, err
	}
	if len(objects) == 0 {
		return nil, nil
	}
	return objects[0], nil
}

// DisplayCmdserv displays the command server status.
func (session *Session) DisplayCmdserv(ctx context.Context, opts ...CommandOption) (map[string]any, error) {
	config := buildCommandConfig(opts)
	objects, err := session.mqscCommand(ctx, "DISPLAY", "CMDSERV", nil, config.requestParameters, config.responseParameters, nil, true)
	if err != nil {
		return nil, err
	}
	if len(objects) == 0 {
		return nil, nil
	}
	return objects[0], nil
}

// ---------------------------------------------------------------------------
// DISPLAY commands — optional name, list return
// ---------------------------------------------------------------------------

// DisplayApstatus displays application status.
func (session *Session) DisplayApstatus(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "APSTATUS", name, opts)
}

// DisplayArchive displays archive information.
func (session *Session) DisplayArchive(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "ARCHIVE", name, opts)
}

// DisplayAuthinfo displays authentication information objects.
func (session *Session) DisplayAuthinfo(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "AUTHINFO", name, opts)
}

// DisplayAuthrec displays authority records.
func (session *Session) DisplayAuthrec(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "AUTHREC", name, opts)
}

// DisplayAuthserv displays authorization service information.
func (session *Session) DisplayAuthserv(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "AUTHSERV", name, opts)
}

// DisplayCfstatus displays CF structure status.
func (session *Session) DisplayCfstatus(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "CFSTATUS", name, opts)
}

// DisplayCfstruct displays CF structure attributes.
func (session *Session) DisplayCfstruct(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "CFSTRUCT", name, opts)
}

// DisplayChinit displays channel initiator information.
func (session *Session) DisplayChinit(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "CHINIT", name, opts)
}

// DisplayChlauth displays channel authentication records.
func (session *Session) DisplayChlauth(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "CHLAUTH", name, opts)
}

// DisplayChstatus displays channel status.
func (session *Session) DisplayChstatus(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "CHSTATUS", name, opts)
}

// DisplayClusqmgr displays cluster queue manager information.
func (session *Session) DisplayClusqmgr(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "CLUSQMGR", name, opts)
}

// DisplayComminfo displays communication information objects.
func (session *Session) DisplayComminfo(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "COMMINFO", name, opts)
}

// DisplayConn displays connection information.
func (session *Session) DisplayConn(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "CONN", name, opts)
}

// DisplayEntauth displays entity authority information.
func (session *Session) DisplayEntauth(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "ENTAUTH", name, opts)
}

// DisplayGroup displays group information.
func (session *Session) DisplayGroup(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "GROUP", name, opts)
}

// DisplayListener displays listener attributes.
func (session *Session) DisplayListener(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "LISTENER", name, opts)
}

// DisplayLog displays log information.
func (session *Session) DisplayLog(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "LOG", name, opts)
}

// DisplayLsstatus displays listener status.
func (session *Session) DisplayLsstatus(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "LSSTATUS", name, opts)
}

// DisplayMaxsmsgs displays maximum short messages information.
func (session *Session) DisplayMaxsmsgs(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "MAXSMSGS", name, opts)
}

// DisplayNamelist displays namelist attributes.
func (session *Session) DisplayNamelist(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "NAMELIST", name, opts)
}

// DisplayPolicy displays security policy information.
func (session *Session) DisplayPolicy(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "POLICY", name, opts)
}

// DisplayProcess displays process attributes.
func (session *Session) DisplayProcess(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "PROCESS", name, opts)
}

// DisplayPubsub displays pub/sub status.
func (session *Session) DisplayPubsub(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "PUBSUB", name, opts)
}

// DisplayQstatus displays queue status.
func (session *Session) DisplayQstatus(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "QSTATUS", name, opts)
}

// DisplaySbstatus displays subscription status.
func (session *Session) DisplaySbstatus(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "SBSTATUS", name, opts)
}

// DisplaySecurity displays security settings.
func (session *Session) DisplaySecurity(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "SECURITY", name, opts)
}

// DisplayService displays service attributes.
func (session *Session) DisplayService(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "SERVICE", name, opts)
}

// DisplaySmds displays shared message data set information.
func (session *Session) DisplaySmds(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "SMDS", name, opts)
}

// DisplaySmdsconn displays shared message data set connection information.
func (session *Session) DisplaySmdsconn(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "SMDSCONN", name, opts)
}

// DisplayStgclass displays storage class attributes.
func (session *Session) DisplayStgclass(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "STGCLASS", name, opts)
}

// DisplaySub displays subscription attributes.
func (session *Session) DisplaySub(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "SUB", name, opts)
}

// DisplaySvstatus displays service status.
func (session *Session) DisplaySvstatus(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "SVSTATUS", name, opts)
}

// DisplaySystem displays system information.
func (session *Session) DisplaySystem(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "SYSTEM", name, opts)
}

// DisplayTcluster displays topic cluster information.
func (session *Session) DisplayTcluster(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "TCLUSTER", name, opts)
}

// DisplayThread displays thread information.
func (session *Session) DisplayThread(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "THREAD", name, opts)
}

// DisplayTopic displays topic attributes.
func (session *Session) DisplayTopic(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "TOPIC", name, opts)
}

// DisplayTpstatus displays topic status.
func (session *Session) DisplayTpstatus(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "TPSTATUS", name, opts)
}

// DisplayTrace displays trace information.
func (session *Session) DisplayTrace(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "TRACE", name, opts)
}

// DisplayUsage displays usage information.
func (session *Session) DisplayUsage(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "USAGE", name, opts)
}

// displayList is the shared implementation for optional-name DISPLAY commands
// that return a list.
func (session *Session) displayList(ctx context.Context, qualifier, name string, opts []CommandOption) ([]map[string]any, error) {
	config := buildCommandConfig(opts)
	var namePtr *string
	if name != "" {
		namePtr = &name
	}
	return session.mqscCommand(ctx, "DISPLAY", qualifier, namePtr, config.requestParameters, config.responseParameters, config.where, true)
}

// ---------------------------------------------------------------------------
// DEFINE commands
// ---------------------------------------------------------------------------

// DefineQlocal defines a local queue. Name is required.
func (session *Session) DefineQlocal(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommand(ctx, "DEFINE", "QLOCAL", &name, opts)
}

// DefineQremote defines a remote queue. Name is required.
func (session *Session) DefineQremote(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommand(ctx, "DEFINE", "QREMOTE", &name, opts)
}

// DefineQalias defines an alias queue. Name is required.
func (session *Session) DefineQalias(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommand(ctx, "DEFINE", "QALIAS", &name, opts)
}

// DefineQmodel defines a model queue. Name is required.
func (session *Session) DefineQmodel(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommand(ctx, "DEFINE", "QMODEL", &name, opts)
}

// DefineChannel defines a channel. Name is required.
func (session *Session) DefineChannel(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommand(ctx, "DEFINE", "CHANNEL", &name, opts)
}

// DefineAuthinfo defines an authentication information object.
func (session *Session) DefineAuthinfo(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DEFINE", "AUTHINFO", name, opts)
}

// DefineBuffpool defines a buffer pool.
func (session *Session) DefineBuffpool(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DEFINE", "BUFFPOOL", name, opts)
}

// DefineCfstruct defines a CF structure.
func (session *Session) DefineCfstruct(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DEFINE", "CFSTRUCT", name, opts)
}

// DefineComminfo defines a communication information object.
func (session *Session) DefineComminfo(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DEFINE", "COMMINFO", name, opts)
}

// DefineListener defines a listener.
func (session *Session) DefineListener(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DEFINE", "LISTENER", name, opts)
}

// DefineLog defines a log.
func (session *Session) DefineLog(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DEFINE", "LOG", name, opts)
}

// DefineMaxsmsgs defines maximum short messages.
func (session *Session) DefineMaxsmsgs(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DEFINE", "MAXSMSGS", name, opts)
}

// DefineNamelist defines a namelist.
func (session *Session) DefineNamelist(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DEFINE", "NAMELIST", name, opts)
}

// DefineProcess defines a process.
func (session *Session) DefineProcess(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DEFINE", "PROCESS", name, opts)
}

// DefinePsid defines a page set ID.
func (session *Session) DefinePsid(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DEFINE", "PSID", name, opts)
}

// DefineService defines a service.
func (session *Session) DefineService(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DEFINE", "SERVICE", name, opts)
}

// DefineStgclass defines a storage class.
func (session *Session) DefineStgclass(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DEFINE", "STGCLASS", name, opts)
}

// DefineSub defines a subscription.
func (session *Session) DefineSub(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DEFINE", "SUB", name, opts)
}

// DefineTopic defines a topic.
func (session *Session) DefineTopic(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DEFINE", "TOPIC", name, opts)
}

// ---------------------------------------------------------------------------
// ALTER commands
// ---------------------------------------------------------------------------

// AlterQmgr alters queue manager attributes.
func (session *Session) AlterQmgr(ctx context.Context, opts ...CommandOption) error {
	return session.voidCommand(ctx, "ALTER", "QMGR", nil, opts)
}

// AlterAuthinfo alters an authentication information object.
func (session *Session) AlterAuthinfo(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "ALTER", "AUTHINFO", name, opts)
}

// AlterBuffpool alters a buffer pool.
func (session *Session) AlterBuffpool(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "ALTER", "BUFFPOOL", name, opts)
}

// AlterCfstruct alters a CF structure.
func (session *Session) AlterCfstruct(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "ALTER", "CFSTRUCT", name, opts)
}

// AlterChannel alters a channel.
func (session *Session) AlterChannel(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "ALTER", "CHANNEL", name, opts)
}

// AlterComminfo alters a communication information object.
func (session *Session) AlterComminfo(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "ALTER", "COMMINFO", name, opts)
}

// AlterListener alters a listener.
func (session *Session) AlterListener(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "ALTER", "LISTENER", name, opts)
}

// AlterNamelist alters a namelist.
func (session *Session) AlterNamelist(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "ALTER", "NAMELIST", name, opts)
}

// AlterProcess alters a process.
func (session *Session) AlterProcess(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "ALTER", "PROCESS", name, opts)
}

// AlterPsid alters a page set ID.
func (session *Session) AlterPsid(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "ALTER", "PSID", name, opts)
}

// AlterSecurity alters security settings.
func (session *Session) AlterSecurity(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "ALTER", "SECURITY", name, opts)
}

// AlterService alters a service.
func (session *Session) AlterService(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "ALTER", "SERVICE", name, opts)
}

// AlterSmds alters a shared message data set.
func (session *Session) AlterSmds(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "ALTER", "SMDS", name, opts)
}

// AlterStgclass alters a storage class.
func (session *Session) AlterStgclass(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "ALTER", "STGCLASS", name, opts)
}

// AlterSub alters a subscription.
func (session *Session) AlterSub(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "ALTER", "SUB", name, opts)
}

// AlterTopic alters a topic.
func (session *Session) AlterTopic(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "ALTER", "TOPIC", name, opts)
}

// AlterTrace alters trace settings.
func (session *Session) AlterTrace(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "ALTER", "TRACE", name, opts)
}

// ---------------------------------------------------------------------------
// DELETE commands
// ---------------------------------------------------------------------------

// DeleteQueue deletes a queue. Name is required.
func (session *Session) DeleteQueue(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommand(ctx, "DELETE", "QUEUE", &name, opts)
}

// DeleteChannel deletes a channel. Name is required.
func (session *Session) DeleteChannel(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommand(ctx, "DELETE", "CHANNEL", &name, opts)
}

// DeleteAuthinfo deletes an authentication information object.
func (session *Session) DeleteAuthinfo(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DELETE", "AUTHINFO", name, opts)
}

// DeleteAuthrec deletes an authority record.
func (session *Session) DeleteAuthrec(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DELETE", "AUTHREC", name, opts)
}

// DeleteBuffpool deletes a buffer pool.
func (session *Session) DeleteBuffpool(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DELETE", "BUFFPOOL", name, opts)
}

// DeleteCfstruct deletes a CF structure.
func (session *Session) DeleteCfstruct(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DELETE", "CFSTRUCT", name, opts)
}

// DeleteComminfo deletes a communication information object.
func (session *Session) DeleteComminfo(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DELETE", "COMMINFO", name, opts)
}

// DeleteListener deletes a listener.
func (session *Session) DeleteListener(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DELETE", "LISTENER", name, opts)
}

// DeleteNamelist deletes a namelist.
func (session *Session) DeleteNamelist(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DELETE", "NAMELIST", name, opts)
}

// DeletePolicy deletes a security policy.
func (session *Session) DeletePolicy(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DELETE", "POLICY", name, opts)
}

// DeleteProcess deletes a process.
func (session *Session) DeleteProcess(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DELETE", "PROCESS", name, opts)
}

// DeletePsid deletes a page set ID.
func (session *Session) DeletePsid(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DELETE", "PSID", name, opts)
}

// DeleteService deletes a service.
func (session *Session) DeleteService(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DELETE", "SERVICE", name, opts)
}

// DeleteStgclass deletes a storage class.
func (session *Session) DeleteStgclass(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DELETE", "STGCLASS", name, opts)
}

// DeleteSub deletes a subscription.
func (session *Session) DeleteSub(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DELETE", "SUB", name, opts)
}

// DeleteTopic deletes a topic.
func (session *Session) DeleteTopic(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DELETE", "TOPIC", name, opts)
}

// ---------------------------------------------------------------------------
// START commands
// ---------------------------------------------------------------------------

// StartQmgr starts the queue manager.
func (session *Session) StartQmgr(ctx context.Context, opts ...CommandOption) error {
	return session.voidCommand(ctx, "START", "QMGR", nil, opts)
}

// StartCmdserv starts the command server.
func (session *Session) StartCmdserv(ctx context.Context, opts ...CommandOption) error {
	return session.voidCommand(ctx, "START", "CMDSERV", nil, opts)
}

// StartChannel starts a channel.
func (session *Session) StartChannel(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "START", "CHANNEL", name, opts)
}

// StartChinit starts the channel initiator.
func (session *Session) StartChinit(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "START", "CHINIT", name, opts)
}

// StartListener starts a listener.
func (session *Session) StartListener(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "START", "LISTENER", name, opts)
}

// StartService starts a service.
func (session *Session) StartService(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "START", "SERVICE", name, opts)
}

// StartSmdsconn starts a shared message data set connection.
func (session *Session) StartSmdsconn(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "START", "SMDSCONN", name, opts)
}

// StartTrace starts tracing.
func (session *Session) StartTrace(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "START", "TRACE", name, opts)
}

// ---------------------------------------------------------------------------
// STOP commands
// ---------------------------------------------------------------------------

// StopQmgr stops the queue manager.
func (session *Session) StopQmgr(ctx context.Context, opts ...CommandOption) error {
	return session.voidCommand(ctx, "STOP", "QMGR", nil, opts)
}

// StopCmdserv stops the command server.
func (session *Session) StopCmdserv(ctx context.Context, opts ...CommandOption) error {
	return session.voidCommand(ctx, "STOP", "CMDSERV", nil, opts)
}

// StopChannel stops a channel.
func (session *Session) StopChannel(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "STOP", "CHANNEL", name, opts)
}

// StopChinit stops the channel initiator.
func (session *Session) StopChinit(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "STOP", "CHINIT", name, opts)
}

// StopConn stops a connection.
func (session *Session) StopConn(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "STOP", "CONN", name, opts)
}

// StopListener stops a listener.
func (session *Session) StopListener(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "STOP", "LISTENER", name, opts)
}

// StopService stops a service.
func (session *Session) StopService(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "STOP", "SERVICE", name, opts)
}

// StopSmdsconn stops a shared message data set connection.
func (session *Session) StopSmdsconn(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "STOP", "SMDSCONN", name, opts)
}

// StopTrace stops tracing.
func (session *Session) StopTrace(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "STOP", "TRACE", name, opts)
}

// ---------------------------------------------------------------------------
// PING commands
// ---------------------------------------------------------------------------

// PingQmgr pings the queue manager.
func (session *Session) PingQmgr(ctx context.Context, opts ...CommandOption) error {
	return session.voidCommand(ctx, "PING", "QMGR", nil, opts)
}

// PingChannel pings a channel.
func (session *Session) PingChannel(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "PING", "CHANNEL", name, opts)
}

// ---------------------------------------------------------------------------
// CLEAR commands
// ---------------------------------------------------------------------------

// ClearQlocal clears a local queue.
func (session *Session) ClearQlocal(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "CLEAR", "QLOCAL", name, opts)
}

// ClearTopicstr clears a topic string.
func (session *Session) ClearTopicstr(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "CLEAR", "TOPICSTR", name, opts)
}

// ---------------------------------------------------------------------------
// REFRESH commands
// ---------------------------------------------------------------------------

// RefreshQmgr refreshes the queue manager.
func (session *Session) RefreshQmgr(ctx context.Context, opts ...CommandOption) error {
	return session.voidCommand(ctx, "REFRESH", "QMGR", nil, opts)
}

// RefreshCluster refreshes cluster information.
func (session *Session) RefreshCluster(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "REFRESH", "CLUSTER", name, opts)
}

// RefreshSecurity refreshes security settings.
func (session *Session) RefreshSecurity(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "REFRESH", "SECURITY", name, opts)
}

// ---------------------------------------------------------------------------
// RESET commands
// ---------------------------------------------------------------------------

// ResetQmgr resets the queue manager.
func (session *Session) ResetQmgr(ctx context.Context, opts ...CommandOption) error {
	return session.voidCommand(ctx, "RESET", "QMGR", nil, opts)
}

// ResetCfstruct resets a CF structure.
func (session *Session) ResetCfstruct(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "RESET", "CFSTRUCT", name, opts)
}

// ResetChannel resets a channel.
func (session *Session) ResetChannel(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "RESET", "CHANNEL", name, opts)
}

// ResetCluster resets cluster information.
func (session *Session) ResetCluster(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "RESET", "CLUSTER", name, opts)
}

// ResetQstats resets queue statistics.
func (session *Session) ResetQstats(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "RESET", "QSTATS", name, opts)
}

// ResetSmds resets a shared message data set.
func (session *Session) ResetSmds(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "RESET", "SMDS", name, opts)
}

// ResetTpipe resets a transaction pipe.
func (session *Session) ResetTpipe(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "RESET", "TPIPE", name, opts)
}

// ---------------------------------------------------------------------------
// RESOLVE commands
// ---------------------------------------------------------------------------

// ResolveChannel resolves an in-doubt channel.
func (session *Session) ResolveChannel(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "RESOLVE", "CHANNEL", name, opts)
}

// ResolveIndoubt resolves in-doubt transactions.
func (session *Session) ResolveIndoubt(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "RESOLVE", "INDOUBT", name, opts)
}

// ---------------------------------------------------------------------------
// RESUME / SUSPEND commands
// ---------------------------------------------------------------------------

// ResumeQmgr resumes the queue manager.
func (session *Session) ResumeQmgr(ctx context.Context, opts ...CommandOption) error {
	return session.voidCommand(ctx, "RESUME", "QMGR", nil, opts)
}

// SuspendQmgr suspends the queue manager.
func (session *Session) SuspendQmgr(ctx context.Context, opts ...CommandOption) error {
	return session.voidCommand(ctx, "SUSPEND", "QMGR", nil, opts)
}

// ---------------------------------------------------------------------------
// SET commands
// ---------------------------------------------------------------------------

// SetArchive sets archive parameters.
func (session *Session) SetArchive(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "SET", "ARCHIVE", name, opts)
}

// SetAuthrec sets an authority record.
func (session *Session) SetAuthrec(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "SET", "AUTHREC", name, opts)
}

// SetChlauth sets a channel authentication record.
func (session *Session) SetChlauth(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "SET", "CHLAUTH", name, opts)
}

// SetLog sets log parameters.
func (session *Session) SetLog(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "SET", "LOG", name, opts)
}

// SetPolicy sets a security policy.
func (session *Session) SetPolicy(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "SET", "POLICY", name, opts)
}

// SetSystem sets system parameters.
func (session *Session) SetSystem(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "SET", "SYSTEM", name, opts)
}

// ---------------------------------------------------------------------------
// Miscellaneous commands
// ---------------------------------------------------------------------------

// ArchiveLog archives the log.
func (session *Session) ArchiveLog(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "ARCHIVE", "LOG", name, opts)
}

// BackupCfstruct backs up a CF structure.
func (session *Session) BackupCfstruct(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "BACKUP", "CFSTRUCT", name, opts)
}

// RecoverBsds recovers a bootstrap data set.
func (session *Session) RecoverBsds(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "RECOVER", "BSDS", name, opts)
}

// RecoverCfstruct recovers a CF structure.
func (session *Session) RecoverCfstruct(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "RECOVER", "CFSTRUCT", name, opts)
}

// PurgeChannel purges a channel.
func (session *Session) PurgeChannel(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "PURGE", "CHANNEL", name, opts)
}

// MoveQlocal moves messages from one local queue to another.
func (session *Session) MoveQlocal(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "MOVE", "QLOCAL", name, opts)
}

// RverifySecurity reverifies security for a user.
func (session *Session) RverifySecurity(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "RVERIFY", "SECURITY", name, opts)
}

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

// voidCommand dispatches a non-DISPLAY MQSC command and discards the result.
func (session *Session) voidCommand(ctx context.Context, command, qualifier string, name *string, opts []CommandOption) error {
	config := buildCommandConfig(opts)
	_, err := session.mqscCommand(ctx, command, qualifier, name, config.requestParameters, config.responseParameters, nil, false)
	return err
}

// voidCommandOptionalName dispatches a non-DISPLAY MQSC command with an
// optional name parameter.
func (session *Session) voidCommandOptionalName(ctx context.Context, command, qualifier, name string, opts []CommandOption) error {
	var namePtr *string
	if name != "" {
		namePtr = &name
	}
	return session.voidCommand(ctx, command, qualifier, namePtr, opts)
}
