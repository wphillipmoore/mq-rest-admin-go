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

// BEGIN GENERATED MQSC METHODS

// AlterAuthinfo executes the ALTER AUTHINFO command.
func (session *Session) AlterAuthinfo(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "ALTER", "AUTHINFO", name, opts)
}

// AlterBuffpool executes the ALTER BUFFPOOL command.
func (session *Session) AlterBuffpool(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "ALTER", "BUFFPOOL", name, opts)
}

// AlterCfstruct executes the ALTER CFSTRUCT command.
func (session *Session) AlterCfstruct(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "ALTER", "CFSTRUCT", name, opts)
}

// AlterChannel executes the ALTER CHANNEL command.
func (session *Session) AlterChannel(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "ALTER", "CHANNEL", name, opts)
}

// AlterComminfo executes the ALTER COMMINFO command.
func (session *Session) AlterComminfo(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "ALTER", "COMMINFO", name, opts)
}

// AlterListener executes the ALTER LISTENER command.
func (session *Session) AlterListener(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "ALTER", "LISTENER", name, opts)
}

// AlterNamelist executes the ALTER NAMELIST command.
func (session *Session) AlterNamelist(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "ALTER", "NAMELIST", name, opts)
}

// AlterProcess executes the ALTER PROCESS command.
func (session *Session) AlterProcess(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "ALTER", "PROCESS", name, opts)
}

// AlterPsid executes the ALTER PSID command.
func (session *Session) AlterPsid(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "ALTER", "PSID", name, opts)
}

// AlterQalias executes the ALTER QALIAS command.
func (session *Session) AlterQalias(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "ALTER", "QALIAS", name, opts)
}

// AlterQlocal executes the ALTER QLOCAL command.
func (session *Session) AlterQlocal(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "ALTER", "QLOCAL", name, opts)
}

// AlterQmgr executes the ALTER QMGR command.
func (session *Session) AlterQmgr(ctx context.Context, opts ...CommandOption) error {
	return session.voidCommand(ctx, "ALTER", "QMGR", nil, opts)
}

// AlterQmodel executes the ALTER QMODEL command.
func (session *Session) AlterQmodel(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "ALTER", "QMODEL", name, opts)
}

// AlterQremote executes the ALTER QREMOTE command.
func (session *Session) AlterQremote(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "ALTER", "QREMOTE", name, opts)
}

// AlterSecurity executes the ALTER SECURITY command.
func (session *Session) AlterSecurity(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "ALTER", "SECURITY", name, opts)
}

// AlterService executes the ALTER SERVICE command.
func (session *Session) AlterService(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "ALTER", "SERVICE", name, opts)
}

// AlterSmds executes the ALTER SMDS command.
func (session *Session) AlterSmds(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "ALTER", "SMDS", name, opts)
}

// AlterStgclass executes the ALTER STGCLASS command.
func (session *Session) AlterStgclass(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "ALTER", "STGCLASS", name, opts)
}

// AlterSub executes the ALTER SUB command.
func (session *Session) AlterSub(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "ALTER", "SUB", name, opts)
}

// AlterTopic executes the ALTER TOPIC command.
func (session *Session) AlterTopic(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "ALTER", "TOPIC", name, opts)
}

// AlterTrace executes the ALTER TRACE command.
func (session *Session) AlterTrace(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "ALTER", "TRACE", name, opts)
}

// ArchiveLog executes the ARCHIVE LOG command.
func (session *Session) ArchiveLog(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "ARCHIVE", "LOG", name, opts)
}

// BackupCfstruct executes the BACKUP CFSTRUCT command.
func (session *Session) BackupCfstruct(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "BACKUP", "CFSTRUCT", name, opts)
}

// ClearQlocal executes the CLEAR QLOCAL command.
func (session *Session) ClearQlocal(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "CLEAR", "QLOCAL", name, opts)
}

// ClearTopicstr executes the CLEAR TOPICSTR command.
func (session *Session) ClearTopicstr(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "CLEAR", "TOPICSTR", name, opts)
}

// DefineAuthinfo executes the DEFINE AUTHINFO command.
func (session *Session) DefineAuthinfo(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DEFINE", "AUTHINFO", name, opts)
}

// DefineBuffpool executes the DEFINE BUFFPOOL command.
func (session *Session) DefineBuffpool(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DEFINE", "BUFFPOOL", name, opts)
}

// DefineCfstruct executes the DEFINE CFSTRUCT command.
func (session *Session) DefineCfstruct(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DEFINE", "CFSTRUCT", name, opts)
}

// DefineChannel executes the DEFINE CHANNEL command. Name is required.
func (session *Session) DefineChannel(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommand(ctx, "DEFINE", "CHANNEL", &name, opts)
}

// DefineComminfo executes the DEFINE COMMINFO command.
func (session *Session) DefineComminfo(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DEFINE", "COMMINFO", name, opts)
}

// DefineListener executes the DEFINE LISTENER command.
func (session *Session) DefineListener(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DEFINE", "LISTENER", name, opts)
}

// DefineLog executes the DEFINE LOG command.
func (session *Session) DefineLog(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DEFINE", "LOG", name, opts)
}

// DefineMaxsmsgs executes the DEFINE MAXSMSGS command.
func (session *Session) DefineMaxsmsgs(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DEFINE", "MAXSMSGS", name, opts)
}

// DefineNamelist executes the DEFINE NAMELIST command.
func (session *Session) DefineNamelist(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DEFINE", "NAMELIST", name, opts)
}

// DefineProcess executes the DEFINE PROCESS command.
func (session *Session) DefineProcess(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DEFINE", "PROCESS", name, opts)
}

// DefinePsid executes the DEFINE PSID command.
func (session *Session) DefinePsid(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DEFINE", "PSID", name, opts)
}

// DefineQalias executes the DEFINE QALIAS command. Name is required.
func (session *Session) DefineQalias(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommand(ctx, "DEFINE", "QALIAS", &name, opts)
}

// DefineQlocal executes the DEFINE QLOCAL command. Name is required.
func (session *Session) DefineQlocal(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommand(ctx, "DEFINE", "QLOCAL", &name, opts)
}

// DefineQmodel executes the DEFINE QMODEL command. Name is required.
func (session *Session) DefineQmodel(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommand(ctx, "DEFINE", "QMODEL", &name, opts)
}

// DefineQremote executes the DEFINE QREMOTE command. Name is required.
func (session *Session) DefineQremote(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommand(ctx, "DEFINE", "QREMOTE", &name, opts)
}

// DefineService executes the DEFINE SERVICE command.
func (session *Session) DefineService(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DEFINE", "SERVICE", name, opts)
}

// DefineStgclass executes the DEFINE STGCLASS command.
func (session *Session) DefineStgclass(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DEFINE", "STGCLASS", name, opts)
}

// DefineSub executes the DEFINE SUB command.
func (session *Session) DefineSub(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DEFINE", "SUB", name, opts)
}

// DefineTopic executes the DEFINE TOPIC command.
func (session *Session) DefineTopic(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DEFINE", "TOPIC", name, opts)
}

// DeleteAuthinfo executes the DELETE AUTHINFO command.
func (session *Session) DeleteAuthinfo(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DELETE", "AUTHINFO", name, opts)
}

// DeleteAuthrec executes the DELETE AUTHREC command.
func (session *Session) DeleteAuthrec(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DELETE", "AUTHREC", name, opts)
}

// DeleteBuffpool executes the DELETE BUFFPOOL command.
func (session *Session) DeleteBuffpool(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DELETE", "BUFFPOOL", name, opts)
}

// DeleteCfstruct executes the DELETE CFSTRUCT command.
func (session *Session) DeleteCfstruct(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DELETE", "CFSTRUCT", name, opts)
}

// DeleteChannel executes the DELETE CHANNEL command. Name is required.
func (session *Session) DeleteChannel(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommand(ctx, "DELETE", "CHANNEL", &name, opts)
}

// DeleteComminfo executes the DELETE COMMINFO command.
func (session *Session) DeleteComminfo(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DELETE", "COMMINFO", name, opts)
}

// DeleteListener executes the DELETE LISTENER command.
func (session *Session) DeleteListener(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DELETE", "LISTENER", name, opts)
}

// DeleteNamelist executes the DELETE NAMELIST command.
func (session *Session) DeleteNamelist(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DELETE", "NAMELIST", name, opts)
}

// DeletePolicy executes the DELETE POLICY command.
func (session *Session) DeletePolicy(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DELETE", "POLICY", name, opts)
}

// DeleteProcess executes the DELETE PROCESS command.
func (session *Session) DeleteProcess(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DELETE", "PROCESS", name, opts)
}

// DeletePsid executes the DELETE PSID command.
func (session *Session) DeletePsid(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DELETE", "PSID", name, opts)
}

// DeleteQalias executes the DELETE QALIAS command. Name is required.
func (session *Session) DeleteQalias(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommand(ctx, "DELETE", "QALIAS", &name, opts)
}

// DeleteQlocal executes the DELETE QLOCAL command. Name is required.
func (session *Session) DeleteQlocal(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommand(ctx, "DELETE", "QLOCAL", &name, opts)
}

// DeleteQmodel executes the DELETE QMODEL command. Name is required.
func (session *Session) DeleteQmodel(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommand(ctx, "DELETE", "QMODEL", &name, opts)
}

// DeleteQremote executes the DELETE QREMOTE command. Name is required.
func (session *Session) DeleteQremote(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommand(ctx, "DELETE", "QREMOTE", &name, opts)
}

// DeleteQueue executes the DELETE QUEUE command. Name is required.
func (session *Session) DeleteQueue(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommand(ctx, "DELETE", "QUEUE", &name, opts)
}

// DeleteService executes the DELETE SERVICE command.
func (session *Session) DeleteService(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DELETE", "SERVICE", name, opts)
}

// DeleteStgclass executes the DELETE STGCLASS command.
func (session *Session) DeleteStgclass(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DELETE", "STGCLASS", name, opts)
}

// DeleteSub executes the DELETE SUB command.
func (session *Session) DeleteSub(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DELETE", "SUB", name, opts)
}

// DeleteTopic executes the DELETE TOPIC command.
func (session *Session) DeleteTopic(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "DELETE", "TOPIC", name, opts)
}

// DisplayApstatus executes the DISPLAY APSTATUS command.
func (session *Session) DisplayApstatus(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "APSTATUS", name, opts)
}

// DisplayArchive executes the DISPLAY ARCHIVE command.
func (session *Session) DisplayArchive(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "ARCHIVE", name, opts)
}

// DisplayAuthinfo executes the DISPLAY AUTHINFO command.
func (session *Session) DisplayAuthinfo(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "AUTHINFO", name, opts)
}

// DisplayAuthrec executes the DISPLAY AUTHREC command.
func (session *Session) DisplayAuthrec(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "AUTHREC", name, opts)
}

// DisplayAuthserv executes the DISPLAY AUTHSERV command.
func (session *Session) DisplayAuthserv(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "AUTHSERV", name, opts)
}

// DisplayCfstatus executes the DISPLAY CFSTATUS command.
func (session *Session) DisplayCfstatus(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "CFSTATUS", name, opts)
}

// DisplayCfstruct executes the DISPLAY CFSTRUCT command.
func (session *Session) DisplayCfstruct(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "CFSTRUCT", name, opts)
}

// DisplayChannel executes the DISPLAY CHANNEL command. Name defaults to "*" if empty.
func (session *Session) DisplayChannel(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	config := buildCommandConfig(opts)
	displayName := name
	if displayName == "" {
		displayName = "*"
	}
	return session.mqscCommand(ctx, "DISPLAY", "CHANNEL", &displayName, config.requestParameters, config.responseParameters, config.where, true)
}

// DisplayChinit executes the DISPLAY CHINIT command.
func (session *Session) DisplayChinit(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "CHINIT", name, opts)
}

// DisplayChlauth executes the DISPLAY CHLAUTH command.
func (session *Session) DisplayChlauth(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "CHLAUTH", name, opts)
}

// DisplayChstatus executes the DISPLAY CHSTATUS command.
func (session *Session) DisplayChstatus(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "CHSTATUS", name, opts)
}

// DisplayClusqmgr executes the DISPLAY CLUSQMGR command.
func (session *Session) DisplayClusqmgr(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "CLUSQMGR", name, opts)
}

// DisplayCmdserv executes the DISPLAY CMDSERV command.
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

// DisplayComminfo executes the DISPLAY COMMINFO command.
func (session *Session) DisplayComminfo(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "COMMINFO", name, opts)
}

// DisplayConn executes the DISPLAY CONN command.
func (session *Session) DisplayConn(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "CONN", name, opts)
}

// DisplayEntauth executes the DISPLAY ENTAUTH command.
func (session *Session) DisplayEntauth(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "ENTAUTH", name, opts)
}

// DisplayGroup executes the DISPLAY GROUP command.
func (session *Session) DisplayGroup(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "GROUP", name, opts)
}

// DisplayListener executes the DISPLAY LISTENER command.
func (session *Session) DisplayListener(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "LISTENER", name, opts)
}

// DisplayLog executes the DISPLAY LOG command.
func (session *Session) DisplayLog(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "LOG", name, opts)
}

// DisplayLsstatus executes the DISPLAY LSSTATUS command.
func (session *Session) DisplayLsstatus(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "LSSTATUS", name, opts)
}

// DisplayMaxsmsgs executes the DISPLAY MAXSMSGS command.
func (session *Session) DisplayMaxsmsgs(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "MAXSMSGS", name, opts)
}

// DisplayNamelist executes the DISPLAY NAMELIST command.
func (session *Session) DisplayNamelist(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "NAMELIST", name, opts)
}

// DisplayPolicy executes the DISPLAY POLICY command.
func (session *Session) DisplayPolicy(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "POLICY", name, opts)
}

// DisplayProcess executes the DISPLAY PROCESS command.
func (session *Session) DisplayProcess(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "PROCESS", name, opts)
}

// DisplayPubsub executes the DISPLAY PUBSUB command.
func (session *Session) DisplayPubsub(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "PUBSUB", name, opts)
}

// DisplayQmgr executes the DISPLAY QMGR command.
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

// DisplayQmstatus executes the DISPLAY QMSTATUS command.
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

// DisplayQstatus executes the DISPLAY QSTATUS command.
func (session *Session) DisplayQstatus(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "QSTATUS", name, opts)
}

// DisplayQueue executes the DISPLAY QUEUE command. Name defaults to "*" if empty.
func (session *Session) DisplayQueue(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	config := buildCommandConfig(opts)
	displayName := name
	if displayName == "" {
		displayName = "*"
	}
	return session.mqscCommand(ctx, "DISPLAY", "QUEUE", &displayName, config.requestParameters, config.responseParameters, config.where, true)
}

// DisplaySbstatus executes the DISPLAY SBSTATUS command.
func (session *Session) DisplaySbstatus(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "SBSTATUS", name, opts)
}

// DisplaySecurity executes the DISPLAY SECURITY command.
func (session *Session) DisplaySecurity(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "SECURITY", name, opts)
}

// DisplayService executes the DISPLAY SERVICE command.
func (session *Session) DisplayService(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "SERVICE", name, opts)
}

// DisplaySmds executes the DISPLAY SMDS command.
func (session *Session) DisplaySmds(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "SMDS", name, opts)
}

// DisplaySmdsconn executes the DISPLAY SMDSCONN command.
func (session *Session) DisplaySmdsconn(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "SMDSCONN", name, opts)
}

// DisplayStgclass executes the DISPLAY STGCLASS command.
func (session *Session) DisplayStgclass(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "STGCLASS", name, opts)
}

// DisplaySub executes the DISPLAY SUB command.
func (session *Session) DisplaySub(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "SUB", name, opts)
}

// DisplaySvstatus executes the DISPLAY SVSTATUS command.
func (session *Session) DisplaySvstatus(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "SVSTATUS", name, opts)
}

// DisplaySystem executes the DISPLAY SYSTEM command.
func (session *Session) DisplaySystem(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "SYSTEM", name, opts)
}

// DisplayTcluster executes the DISPLAY TCLUSTER command.
func (session *Session) DisplayTcluster(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "TCLUSTER", name, opts)
}

// DisplayThread executes the DISPLAY THREAD command.
func (session *Session) DisplayThread(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "THREAD", name, opts)
}

// DisplayTopic executes the DISPLAY TOPIC command.
func (session *Session) DisplayTopic(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "TOPIC", name, opts)
}

// DisplayTpstatus executes the DISPLAY TPSTATUS command.
func (session *Session) DisplayTpstatus(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "TPSTATUS", name, opts)
}

// DisplayTrace executes the DISPLAY TRACE command.
func (session *Session) DisplayTrace(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "TRACE", name, opts)
}

// DisplayUsage executes the DISPLAY USAGE command.
func (session *Session) DisplayUsage(ctx context.Context, name string, opts ...CommandOption) ([]map[string]any, error) {
	return session.displayList(ctx, "USAGE", name, opts)
}

// MoveQlocal executes the MOVE QLOCAL command.
func (session *Session) MoveQlocal(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "MOVE", "QLOCAL", name, opts)
}

// PingChannel executes the PING CHANNEL command.
func (session *Session) PingChannel(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "PING", "CHANNEL", name, opts)
}

// PingQmgr executes the PING QMGR command.
func (session *Session) PingQmgr(ctx context.Context, opts ...CommandOption) error {
	return session.voidCommand(ctx, "PING", "QMGR", nil, opts)
}

// PurgeChannel executes the PURGE CHANNEL command.
func (session *Session) PurgeChannel(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "PURGE", "CHANNEL", name, opts)
}

// RecoverBsds executes the RECOVER BSDS command.
func (session *Session) RecoverBsds(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "RECOVER", "BSDS", name, opts)
}

// RecoverCfstruct executes the RECOVER CFSTRUCT command.
func (session *Session) RecoverCfstruct(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "RECOVER", "CFSTRUCT", name, opts)
}

// RefreshCluster executes the REFRESH CLUSTER command.
func (session *Session) RefreshCluster(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "REFRESH", "CLUSTER", name, opts)
}

// RefreshQmgr executes the REFRESH QMGR command.
func (session *Session) RefreshQmgr(ctx context.Context, opts ...CommandOption) error {
	return session.voidCommand(ctx, "REFRESH", "QMGR", nil, opts)
}

// RefreshSecurity executes the REFRESH SECURITY command.
func (session *Session) RefreshSecurity(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "REFRESH", "SECURITY", name, opts)
}

// ResetCfstruct executes the RESET CFSTRUCT command.
func (session *Session) ResetCfstruct(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "RESET", "CFSTRUCT", name, opts)
}

// ResetChannel executes the RESET CHANNEL command.
func (session *Session) ResetChannel(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "RESET", "CHANNEL", name, opts)
}

// ResetCluster executes the RESET CLUSTER command.
func (session *Session) ResetCluster(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "RESET", "CLUSTER", name, opts)
}

// ResetQmgr executes the RESET QMGR command.
func (session *Session) ResetQmgr(ctx context.Context, opts ...CommandOption) error {
	return session.voidCommand(ctx, "RESET", "QMGR", nil, opts)
}

// ResetQstats executes the RESET QSTATS command.
func (session *Session) ResetQstats(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "RESET", "QSTATS", name, opts)
}

// ResetSmds executes the RESET SMDS command.
func (session *Session) ResetSmds(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "RESET", "SMDS", name, opts)
}

// ResetTpipe executes the RESET TPIPE command.
func (session *Session) ResetTpipe(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "RESET", "TPIPE", name, opts)
}

// ResolveChannel executes the RESOLVE CHANNEL command.
func (session *Session) ResolveChannel(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "RESOLVE", "CHANNEL", name, opts)
}

// ResolveIndoubt executes the RESOLVE INDOUBT command.
func (session *Session) ResolveIndoubt(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "RESOLVE", "INDOUBT", name, opts)
}

// ResumeQmgr executes the RESUME QMGR command.
func (session *Session) ResumeQmgr(ctx context.Context, opts ...CommandOption) error {
	return session.voidCommand(ctx, "RESUME", "QMGR", nil, opts)
}

// RverifySecurity executes the RVERIFY SECURITY command.
func (session *Session) RverifySecurity(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "RVERIFY", "SECURITY", name, opts)
}

// SetArchive executes the SET ARCHIVE command.
func (session *Session) SetArchive(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "SET", "ARCHIVE", name, opts)
}

// SetAuthrec executes the SET AUTHREC command.
func (session *Session) SetAuthrec(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "SET", "AUTHREC", name, opts)
}

// SetChlauth executes the SET CHLAUTH command.
func (session *Session) SetChlauth(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "SET", "CHLAUTH", name, opts)
}

// SetLog executes the SET LOG command.
func (session *Session) SetLog(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "SET", "LOG", name, opts)
}

// SetPolicy executes the SET POLICY command.
func (session *Session) SetPolicy(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "SET", "POLICY", name, opts)
}

// SetSystem executes the SET SYSTEM command.
func (session *Session) SetSystem(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "SET", "SYSTEM", name, opts)
}

// StartChannel executes the START CHANNEL command.
func (session *Session) StartChannel(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "START", "CHANNEL", name, opts)
}

// StartChinit executes the START CHINIT command.
func (session *Session) StartChinit(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "START", "CHINIT", name, opts)
}

// StartCmdserv executes the START CMDSERV command.
func (session *Session) StartCmdserv(ctx context.Context, opts ...CommandOption) error {
	return session.voidCommand(ctx, "START", "CMDSERV", nil, opts)
}

// StartListener executes the START LISTENER command.
func (session *Session) StartListener(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "START", "LISTENER", name, opts)
}

// StartQmgr executes the START QMGR command.
func (session *Session) StartQmgr(ctx context.Context, opts ...CommandOption) error {
	return session.voidCommand(ctx, "START", "QMGR", nil, opts)
}

// StartService executes the START SERVICE command.
func (session *Session) StartService(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "START", "SERVICE", name, opts)
}

// StartSmdsconn executes the START SMDSCONN command.
func (session *Session) StartSmdsconn(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "START", "SMDSCONN", name, opts)
}

// StartTrace executes the START TRACE command.
func (session *Session) StartTrace(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "START", "TRACE", name, opts)
}

// StopChannel executes the STOP CHANNEL command.
func (session *Session) StopChannel(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "STOP", "CHANNEL", name, opts)
}

// StopChinit executes the STOP CHINIT command.
func (session *Session) StopChinit(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "STOP", "CHINIT", name, opts)
}

// StopCmdserv executes the STOP CMDSERV command.
func (session *Session) StopCmdserv(ctx context.Context, opts ...CommandOption) error {
	return session.voidCommand(ctx, "STOP", "CMDSERV", nil, opts)
}

// StopConn executes the STOP CONN command.
func (session *Session) StopConn(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "STOP", "CONN", name, opts)
}

// StopListener executes the STOP LISTENER command.
func (session *Session) StopListener(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "STOP", "LISTENER", name, opts)
}

// StopQmgr executes the STOP QMGR command.
func (session *Session) StopQmgr(ctx context.Context, opts ...CommandOption) error {
	return session.voidCommand(ctx, "STOP", "QMGR", nil, opts)
}

// StopService executes the STOP SERVICE command.
func (session *Session) StopService(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "STOP", "SERVICE", name, opts)
}

// StopSmdsconn executes the STOP SMDSCONN command.
func (session *Session) StopSmdsconn(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "STOP", "SMDSCONN", name, opts)
}

// StopTrace executes the STOP TRACE command.
func (session *Session) StopTrace(ctx context.Context, name string, opts ...CommandOption) error {
	return session.voidCommandOptionalName(ctx, "STOP", "TRACE", name, opts)
}

// SuspendQmgr executes the SUSPEND QMGR command.
func (session *Session) SuspendQmgr(ctx context.Context, opts ...CommandOption) error {
	return session.voidCommand(ctx, "SUSPEND", "QMGR", nil, opts)
}

// END GENERATED MQSC METHODS

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

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
