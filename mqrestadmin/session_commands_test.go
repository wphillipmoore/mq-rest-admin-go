package mqrestadmin

import (
	"context"
	"testing"
)

// displayListEntry defines a table entry for testing displayList-based commands.
type displayListEntry struct {
	name      string
	qualifier string
	call      func(*Session, context.Context, string) ([]map[string]any, error)
}

func TestDisplayListCommands(t *testing.T) {
	entries := []displayListEntry{
		{"DisplayApstatus", "APSTATUS", func(s *Session, ctx context.Context, name string) ([]map[string]any, error) {
			return s.DisplayApstatus(ctx, name)
		}},
		{"DisplayArchive", "ARCHIVE", func(s *Session, ctx context.Context, name string) ([]map[string]any, error) {
			return s.DisplayArchive(ctx, name)
		}},
		{"DisplayAuthinfo", "AUTHINFO", func(s *Session, ctx context.Context, name string) ([]map[string]any, error) {
			return s.DisplayAuthinfo(ctx, name)
		}},
		{"DisplayAuthrec", "AUTHREC", func(s *Session, ctx context.Context, name string) ([]map[string]any, error) {
			return s.DisplayAuthrec(ctx, name)
		}},
		{"DisplayAuthserv", "AUTHSERV", func(s *Session, ctx context.Context, name string) ([]map[string]any, error) {
			return s.DisplayAuthserv(ctx, name)
		}},
		{"DisplayCfstatus", "CFSTATUS", func(s *Session, ctx context.Context, name string) ([]map[string]any, error) {
			return s.DisplayCfstatus(ctx, name)
		}},
		{"DisplayCfstruct", "CFSTRUCT", func(s *Session, ctx context.Context, name string) ([]map[string]any, error) {
			return s.DisplayCfstruct(ctx, name)
		}},
		{"DisplayChinit", "CHINIT", func(s *Session, ctx context.Context, name string) ([]map[string]any, error) {
			return s.DisplayChinit(ctx, name)
		}},
		{"DisplayChlauth", "CHLAUTH", func(s *Session, ctx context.Context, name string) ([]map[string]any, error) {
			return s.DisplayChlauth(ctx, name)
		}},
		{"DisplayChstatus", "CHSTATUS", func(s *Session, ctx context.Context, name string) ([]map[string]any, error) {
			return s.DisplayChstatus(ctx, name)
		}},
		{"DisplayClusqmgr", "CLUSQMGR", func(s *Session, ctx context.Context, name string) ([]map[string]any, error) {
			return s.DisplayClusqmgr(ctx, name)
		}},
		{"DisplayComminfo", "COMMINFO", func(s *Session, ctx context.Context, name string) ([]map[string]any, error) {
			return s.DisplayComminfo(ctx, name)
		}},
		{"DisplayConn", "CONN", func(s *Session, ctx context.Context, name string) ([]map[string]any, error) {
			return s.DisplayConn(ctx, name)
		}},
		{"DisplayEntauth", "ENTAUTH", func(s *Session, ctx context.Context, name string) ([]map[string]any, error) {
			return s.DisplayEntauth(ctx, name)
		}},
		{"DisplayGroup", "GROUP", func(s *Session, ctx context.Context, name string) ([]map[string]any, error) {
			return s.DisplayGroup(ctx, name)
		}},
		{"DisplayListener", "LISTENER", func(s *Session, ctx context.Context, name string) ([]map[string]any, error) {
			return s.DisplayListener(ctx, name)
		}},
		{"DisplayLog", "LOG", func(s *Session, ctx context.Context, name string) ([]map[string]any, error) {
			return s.DisplayLog(ctx, name)
		}},
		{"DisplayLsstatus", "LSSTATUS", func(s *Session, ctx context.Context, name string) ([]map[string]any, error) {
			return s.DisplayLsstatus(ctx, name)
		}},
		{"DisplayMaxsmsgs", "MAXSMSGS", func(s *Session, ctx context.Context, name string) ([]map[string]any, error) {
			return s.DisplayMaxsmsgs(ctx, name)
		}},
		{"DisplayNamelist", "NAMELIST", func(s *Session, ctx context.Context, name string) ([]map[string]any, error) {
			return s.DisplayNamelist(ctx, name)
		}},
		{"DisplayPolicy", "POLICY", func(s *Session, ctx context.Context, name string) ([]map[string]any, error) {
			return s.DisplayPolicy(ctx, name)
		}},
		{"DisplayProcess", "PROCESS", func(s *Session, ctx context.Context, name string) ([]map[string]any, error) {
			return s.DisplayProcess(ctx, name)
		}},
		{"DisplayPubsub", "PUBSUB", func(s *Session, ctx context.Context, name string) ([]map[string]any, error) {
			return s.DisplayPubsub(ctx, name)
		}},
		{"DisplayQstatus", "QSTATUS", func(s *Session, ctx context.Context, name string) ([]map[string]any, error) {
			return s.DisplayQstatus(ctx, name)
		}},
		{"DisplaySbstatus", "SBSTATUS", func(s *Session, ctx context.Context, name string) ([]map[string]any, error) {
			return s.DisplaySbstatus(ctx, name)
		}},
		{"DisplaySecurity", "SECURITY", func(s *Session, ctx context.Context, name string) ([]map[string]any, error) {
			return s.DisplaySecurity(ctx, name)
		}},
		{"DisplayService", "SERVICE", func(s *Session, ctx context.Context, name string) ([]map[string]any, error) {
			return s.DisplayService(ctx, name)
		}},
		{"DisplaySmds", "SMDS", func(s *Session, ctx context.Context, name string) ([]map[string]any, error) {
			return s.DisplaySmds(ctx, name)
		}},
		{"DisplaySmdsconn", "SMDSCONN", func(s *Session, ctx context.Context, name string) ([]map[string]any, error) {
			return s.DisplaySmdsconn(ctx, name)
		}},
		{"DisplayStgclass", "STGCLASS", func(s *Session, ctx context.Context, name string) ([]map[string]any, error) {
			return s.DisplayStgclass(ctx, name)
		}},
		{"DisplaySub", "SUB", func(s *Session, ctx context.Context, name string) ([]map[string]any, error) {
			return s.DisplaySub(ctx, name)
		}},
		{"DisplaySvstatus", "SVSTATUS", func(s *Session, ctx context.Context, name string) ([]map[string]any, error) {
			return s.DisplaySvstatus(ctx, name)
		}},
		{"DisplaySystem", "SYSTEM", func(s *Session, ctx context.Context, name string) ([]map[string]any, error) {
			return s.DisplaySystem(ctx, name)
		}},
		{"DisplayTcluster", "TCLUSTER", func(s *Session, ctx context.Context, name string) ([]map[string]any, error) {
			return s.DisplayTcluster(ctx, name)
		}},
		{"DisplayThread", "THREAD", func(s *Session, ctx context.Context, name string) ([]map[string]any, error) {
			return s.DisplayThread(ctx, name)
		}},
		{"DisplayTopic", "TOPIC", func(s *Session, ctx context.Context, name string) ([]map[string]any, error) {
			return s.DisplayTopic(ctx, name)
		}},
		{"DisplayTpstatus", "TPSTATUS", func(s *Session, ctx context.Context, name string) ([]map[string]any, error) {
			return s.DisplayTpstatus(ctx, name)
		}},
		{"DisplayTrace", "TRACE", func(s *Session, ctx context.Context, name string) ([]map[string]any, error) {
			return s.DisplayTrace(ctx, name)
		}},
		{"DisplayUsage", "USAGE", func(s *Session, ctx context.Context, name string) ([]map[string]any, error) {
			return s.DisplayUsage(ctx, name)
		}},
	}

	for _, entry := range entries {
		t.Run(entry.name, func(t *testing.T) {
			transport := newMockTransport()
			transport.addSuccessResponse(map[string]any{"NAME": "OBJ1"})
			session := newTestSession(transport)

			result, err := entry.call(session, context.Background(), "OBJ1")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(result) != 1 {
				t.Fatalf("expected 1 result, got %d", len(result))
			}

			payload := transport.lastCall().Payload
			if payload["command"] != "DISPLAY" {
				t.Errorf("command = %v, want DISPLAY", payload["command"])
			}
			if payload["qualifier"] != entry.qualifier {
				t.Errorf("qualifier = %v, want %s", payload["qualifier"], entry.qualifier)
			}
		})
	}
}

func TestDisplaySingletonCommands(t *testing.T) {
	type singletonEntry struct {
		name      string
		qualifier string
		call      func(*Session, context.Context) (map[string]any, error)
	}

	entries := []singletonEntry{
		{"DisplayQmstatus", "QMSTATUS", func(s *Session, ctx context.Context) (map[string]any, error) { return s.DisplayQmstatus(ctx) }},
		{"DisplayCmdserv", "CMDSERV", func(s *Session, ctx context.Context) (map[string]any, error) { return s.DisplayCmdserv(ctx) }},
	}

	for _, entry := range entries {
		t.Run(entry.name+"_success", func(t *testing.T) {
			transport := newMockTransport()
			transport.addSuccessResponse(map[string]any{"STATUS": "RUNNING"})
			session := newTestSession(transport)

			result, err := entry.call(session, context.Background())
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result == nil {
				t.Fatal("expected non-nil result")
			}
			if result["STATUS"] != "RUNNING" {
				t.Errorf("STATUS = %v, want RUNNING", result["STATUS"])
			}

			payload := transport.lastCall().Payload
			if payload["qualifier"] != entry.qualifier {
				t.Errorf("qualifier = %v, want %s", payload["qualifier"], entry.qualifier)
			}
		})

		t.Run(entry.name+"_empty", func(t *testing.T) {
			transport := newMockTransport()
			transport.addSuccessResponse()
			session := newTestSession(transport)

			result, err := entry.call(session, context.Background())
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != nil {
				t.Errorf("expected nil result, got %v", result)
			}
		})
	}
}

// voidEntry defines a table entry for testing void (non-DISPLAY) commands.
type voidEntry struct {
	name      string
	command   string
	qualifier string
	call      func(*Session, context.Context) error
}

func TestVoidCommands(t *testing.T) {
	entries := []voidEntry{
		// DEFINE commands
		{"DefineQlocal", "DEFINE", "QLOCAL", func(s *Session, ctx context.Context) error { return s.DefineQlocal(ctx, "OBJ") }},
		{"DefineQremote", "DEFINE", "QREMOTE", func(s *Session, ctx context.Context) error { return s.DefineQremote(ctx, "OBJ") }},
		{"DefineQalias", "DEFINE", "QALIAS", func(s *Session, ctx context.Context) error { return s.DefineQalias(ctx, "OBJ") }},
		{"DefineQmodel", "DEFINE", "QMODEL", func(s *Session, ctx context.Context) error { return s.DefineQmodel(ctx, "OBJ") }},
		{"DefineChannel", "DEFINE", "CHANNEL", func(s *Session, ctx context.Context) error { return s.DefineChannel(ctx, "OBJ") }},
		{"DefineAuthinfo", "DEFINE", "AUTHINFO", func(s *Session, ctx context.Context) error { return s.DefineAuthinfo(ctx, "OBJ") }},
		{"DefineBuffpool", "DEFINE", "BUFFPOOL", func(s *Session, ctx context.Context) error { return s.DefineBuffpool(ctx, "OBJ") }},
		{"DefineCfstruct", "DEFINE", "CFSTRUCT", func(s *Session, ctx context.Context) error { return s.DefineCfstruct(ctx, "OBJ") }},
		{"DefineComminfo", "DEFINE", "COMMINFO", func(s *Session, ctx context.Context) error { return s.DefineComminfo(ctx, "OBJ") }},
		{"DefineListener", "DEFINE", "LISTENER", func(s *Session, ctx context.Context) error { return s.DefineListener(ctx, "OBJ") }},
		{"DefineLog", "DEFINE", "LOG", func(s *Session, ctx context.Context) error { return s.DefineLog(ctx, "OBJ") }},
		{"DefineMaxsmsgs", "DEFINE", "MAXSMSGS", func(s *Session, ctx context.Context) error { return s.DefineMaxsmsgs(ctx, "OBJ") }},
		{"DefineNamelist", "DEFINE", "NAMELIST", func(s *Session, ctx context.Context) error { return s.DefineNamelist(ctx, "OBJ") }},
		{"DefineProcess", "DEFINE", "PROCESS", func(s *Session, ctx context.Context) error { return s.DefineProcess(ctx, "OBJ") }},
		{"DefinePsid", "DEFINE", "PSID", func(s *Session, ctx context.Context) error { return s.DefinePsid(ctx, "OBJ") }},
		{"DefineService", "DEFINE", "SERVICE", func(s *Session, ctx context.Context) error { return s.DefineService(ctx, "OBJ") }},
		{"DefineStgclass", "DEFINE", "STGCLASS", func(s *Session, ctx context.Context) error { return s.DefineStgclass(ctx, "OBJ") }},
		{"DefineSub", "DEFINE", "SUB", func(s *Session, ctx context.Context) error { return s.DefineSub(ctx, "OBJ") }},
		{"DefineTopic", "DEFINE", "TOPIC", func(s *Session, ctx context.Context) error { return s.DefineTopic(ctx, "OBJ") }},

		// ALTER commands
		{"AlterQmgr", "ALTER", "QMGR", func(s *Session, ctx context.Context) error { return s.AlterQmgr(ctx) }},
		{"AlterAuthinfo", "ALTER", "AUTHINFO", func(s *Session, ctx context.Context) error { return s.AlterAuthinfo(ctx, "OBJ") }},
		{"AlterBuffpool", "ALTER", "BUFFPOOL", func(s *Session, ctx context.Context) error { return s.AlterBuffpool(ctx, "OBJ") }},
		{"AlterCfstruct", "ALTER", "CFSTRUCT", func(s *Session, ctx context.Context) error { return s.AlterCfstruct(ctx, "OBJ") }},
		{"AlterChannel", "ALTER", "CHANNEL", func(s *Session, ctx context.Context) error { return s.AlterChannel(ctx, "OBJ") }},
		{"AlterComminfo", "ALTER", "COMMINFO", func(s *Session, ctx context.Context) error { return s.AlterComminfo(ctx, "OBJ") }},
		{"AlterListener", "ALTER", "LISTENER", func(s *Session, ctx context.Context) error { return s.AlterListener(ctx, "OBJ") }},
		{"AlterNamelist", "ALTER", "NAMELIST", func(s *Session, ctx context.Context) error { return s.AlterNamelist(ctx, "OBJ") }},
		{"AlterProcess", "ALTER", "PROCESS", func(s *Session, ctx context.Context) error { return s.AlterProcess(ctx, "OBJ") }},
		{"AlterPsid", "ALTER", "PSID", func(s *Session, ctx context.Context) error { return s.AlterPsid(ctx, "OBJ") }},
		{"AlterSecurity", "ALTER", "SECURITY", func(s *Session, ctx context.Context) error { return s.AlterSecurity(ctx, "OBJ") }},
		{"AlterService", "ALTER", "SERVICE", func(s *Session, ctx context.Context) error { return s.AlterService(ctx, "OBJ") }},
		{"AlterSmds", "ALTER", "SMDS", func(s *Session, ctx context.Context) error { return s.AlterSmds(ctx, "OBJ") }},
		{"AlterStgclass", "ALTER", "STGCLASS", func(s *Session, ctx context.Context) error { return s.AlterStgclass(ctx, "OBJ") }},
		{"AlterSub", "ALTER", "SUB", func(s *Session, ctx context.Context) error { return s.AlterSub(ctx, "OBJ") }},
		{"AlterTopic", "ALTER", "TOPIC", func(s *Session, ctx context.Context) error { return s.AlterTopic(ctx, "OBJ") }},
		{"AlterTrace", "ALTER", "TRACE", func(s *Session, ctx context.Context) error { return s.AlterTrace(ctx, "OBJ") }},

		// DELETE commands
		{"DeleteQueue", "DELETE", "QUEUE", func(s *Session, ctx context.Context) error { return s.DeleteQueue(ctx, "OBJ") }},
		{"DeleteChannel", "DELETE", "CHANNEL", func(s *Session, ctx context.Context) error { return s.DeleteChannel(ctx, "OBJ") }},
		{"DeleteAuthinfo", "DELETE", "AUTHINFO", func(s *Session, ctx context.Context) error { return s.DeleteAuthinfo(ctx, "OBJ") }},
		{"DeleteAuthrec", "DELETE", "AUTHREC", func(s *Session, ctx context.Context) error { return s.DeleteAuthrec(ctx, "OBJ") }},
		{"DeleteBuffpool", "DELETE", "BUFFPOOL", func(s *Session, ctx context.Context) error { return s.DeleteBuffpool(ctx, "OBJ") }},
		{"DeleteCfstruct", "DELETE", "CFSTRUCT", func(s *Session, ctx context.Context) error { return s.DeleteCfstruct(ctx, "OBJ") }},
		{"DeleteComminfo", "DELETE", "COMMINFO", func(s *Session, ctx context.Context) error { return s.DeleteComminfo(ctx, "OBJ") }},
		{"DeleteListener", "DELETE", "LISTENER", func(s *Session, ctx context.Context) error { return s.DeleteListener(ctx, "OBJ") }},
		{"DeleteNamelist", "DELETE", "NAMELIST", func(s *Session, ctx context.Context) error { return s.DeleteNamelist(ctx, "OBJ") }},
		{"DeletePolicy", "DELETE", "POLICY", func(s *Session, ctx context.Context) error { return s.DeletePolicy(ctx, "OBJ") }},
		{"DeleteProcess", "DELETE", "PROCESS", func(s *Session, ctx context.Context) error { return s.DeleteProcess(ctx, "OBJ") }},
		{"DeletePsid", "DELETE", "PSID", func(s *Session, ctx context.Context) error { return s.DeletePsid(ctx, "OBJ") }},
		{"DeleteService", "DELETE", "SERVICE", func(s *Session, ctx context.Context) error { return s.DeleteService(ctx, "OBJ") }},
		{"DeleteStgclass", "DELETE", "STGCLASS", func(s *Session, ctx context.Context) error { return s.DeleteStgclass(ctx, "OBJ") }},
		{"DeleteSub", "DELETE", "SUB", func(s *Session, ctx context.Context) error { return s.DeleteSub(ctx, "OBJ") }},
		{"DeleteTopic", "DELETE", "TOPIC", func(s *Session, ctx context.Context) error { return s.DeleteTopic(ctx, "OBJ") }},

		// START commands
		{"StartQmgr", "START", "QMGR", func(s *Session, ctx context.Context) error { return s.StartQmgr(ctx) }},
		{"StartCmdserv", "START", "CMDSERV", func(s *Session, ctx context.Context) error { return s.StartCmdserv(ctx) }},
		{"StartChannel", "START", "CHANNEL", func(s *Session, ctx context.Context) error { return s.StartChannel(ctx, "OBJ") }},
		{"StartChinit", "START", "CHINIT", func(s *Session, ctx context.Context) error { return s.StartChinit(ctx, "OBJ") }},
		{"StartListener", "START", "LISTENER", func(s *Session, ctx context.Context) error { return s.StartListener(ctx, "OBJ") }},
		{"StartService", "START", "SERVICE", func(s *Session, ctx context.Context) error { return s.StartService(ctx, "OBJ") }},
		{"StartSmdsconn", "START", "SMDSCONN", func(s *Session, ctx context.Context) error { return s.StartSmdsconn(ctx, "OBJ") }},
		{"StartTrace", "START", "TRACE", func(s *Session, ctx context.Context) error { return s.StartTrace(ctx, "OBJ") }},

		// STOP commands
		{"StopQmgr", "STOP", "QMGR", func(s *Session, ctx context.Context) error { return s.StopQmgr(ctx) }},
		{"StopCmdserv", "STOP", "CMDSERV", func(s *Session, ctx context.Context) error { return s.StopCmdserv(ctx) }},
		{"StopChannel", "STOP", "CHANNEL", func(s *Session, ctx context.Context) error { return s.StopChannel(ctx, "OBJ") }},
		{"StopChinit", "STOP", "CHINIT", func(s *Session, ctx context.Context) error { return s.StopChinit(ctx, "OBJ") }},
		{"StopConn", "STOP", "CONN", func(s *Session, ctx context.Context) error { return s.StopConn(ctx, "OBJ") }},
		{"StopListener", "STOP", "LISTENER", func(s *Session, ctx context.Context) error { return s.StopListener(ctx, "OBJ") }},
		{"StopService", "STOP", "SERVICE", func(s *Session, ctx context.Context) error { return s.StopService(ctx, "OBJ") }},
		{"StopSmdsconn", "STOP", "SMDSCONN", func(s *Session, ctx context.Context) error { return s.StopSmdsconn(ctx, "OBJ") }},
		{"StopTrace", "STOP", "TRACE", func(s *Session, ctx context.Context) error { return s.StopTrace(ctx, "OBJ") }},

		// PING commands
		{"PingQmgr", "PING", "QMGR", func(s *Session, ctx context.Context) error { return s.PingQmgr(ctx) }},
		{"PingChannel", "PING", "CHANNEL", func(s *Session, ctx context.Context) error { return s.PingChannel(ctx, "OBJ") }},

		// CLEAR commands
		{"ClearQlocal", "CLEAR", "QLOCAL", func(s *Session, ctx context.Context) error { return s.ClearQlocal(ctx, "OBJ") }},
		{"ClearTopicstr", "CLEAR", "TOPICSTR", func(s *Session, ctx context.Context) error { return s.ClearTopicstr(ctx, "OBJ") }},

		// REFRESH commands
		{"RefreshQmgr", "REFRESH", "QMGR", func(s *Session, ctx context.Context) error { return s.RefreshQmgr(ctx) }},
		{"RefreshCluster", "REFRESH", "CLUSTER", func(s *Session, ctx context.Context) error { return s.RefreshCluster(ctx, "OBJ") }},
		{"RefreshSecurity", "REFRESH", "SECURITY", func(s *Session, ctx context.Context) error { return s.RefreshSecurity(ctx, "OBJ") }},

		// RESET commands
		{"ResetQmgr", "RESET", "QMGR", func(s *Session, ctx context.Context) error { return s.ResetQmgr(ctx) }},
		{"ResetCfstruct", "RESET", "CFSTRUCT", func(s *Session, ctx context.Context) error { return s.ResetCfstruct(ctx, "OBJ") }},
		{"ResetChannel", "RESET", "CHANNEL", func(s *Session, ctx context.Context) error { return s.ResetChannel(ctx, "OBJ") }},
		{"ResetCluster", "RESET", "CLUSTER", func(s *Session, ctx context.Context) error { return s.ResetCluster(ctx, "OBJ") }},
		{"ResetQstats", "RESET", "QSTATS", func(s *Session, ctx context.Context) error { return s.ResetQstats(ctx, "OBJ") }},
		{"ResetSmds", "RESET", "SMDS", func(s *Session, ctx context.Context) error { return s.ResetSmds(ctx, "OBJ") }},
		{"ResetTpipe", "RESET", "TPIPE", func(s *Session, ctx context.Context) error { return s.ResetTpipe(ctx, "OBJ") }},

		// RESOLVE commands
		{"ResolveChannel", "RESOLVE", "CHANNEL", func(s *Session, ctx context.Context) error { return s.ResolveChannel(ctx, "OBJ") }},
		{"ResolveIndoubt", "RESOLVE", "INDOUBT", func(s *Session, ctx context.Context) error { return s.ResolveIndoubt(ctx, "OBJ") }},

		// RESUME / SUSPEND commands
		{"ResumeQmgr", "RESUME", "QMGR", func(s *Session, ctx context.Context) error { return s.ResumeQmgr(ctx) }},
		{"SuspendQmgr", "SUSPEND", "QMGR", func(s *Session, ctx context.Context) error { return s.SuspendQmgr(ctx) }},

		// SET commands
		{"SetArchive", "SET", "ARCHIVE", func(s *Session, ctx context.Context) error { return s.SetArchive(ctx, "OBJ") }},
		{"SetAuthrec", "SET", "AUTHREC", func(s *Session, ctx context.Context) error { return s.SetAuthrec(ctx, "OBJ") }},
		{"SetChlauth", "SET", "CHLAUTH", func(s *Session, ctx context.Context) error { return s.SetChlauth(ctx, "OBJ") }},
		{"SetLog", "SET", "LOG", func(s *Session, ctx context.Context) error { return s.SetLog(ctx, "OBJ") }},
		{"SetPolicy", "SET", "POLICY", func(s *Session, ctx context.Context) error { return s.SetPolicy(ctx, "OBJ") }},
		{"SetSystem", "SET", "SYSTEM", func(s *Session, ctx context.Context) error { return s.SetSystem(ctx, "OBJ") }},

		// Miscellaneous commands
		{"ArchiveLog", "ARCHIVE", "LOG", func(s *Session, ctx context.Context) error { return s.ArchiveLog(ctx, "OBJ") }},
		{"BackupCfstruct", "BACKUP", "CFSTRUCT", func(s *Session, ctx context.Context) error { return s.BackupCfstruct(ctx, "OBJ") }},
		{"RecoverBsds", "RECOVER", "BSDS", func(s *Session, ctx context.Context) error { return s.RecoverBsds(ctx, "OBJ") }},
		{"RecoverCfstruct", "RECOVER", "CFSTRUCT", func(s *Session, ctx context.Context) error { return s.RecoverCfstruct(ctx, "OBJ") }},
		{"PurgeChannel", "PURGE", "CHANNEL", func(s *Session, ctx context.Context) error { return s.PurgeChannel(ctx, "OBJ") }},
		{"MoveQlocal", "MOVE", "QLOCAL", func(s *Session, ctx context.Context) error { return s.MoveQlocal(ctx, "OBJ") }},
		{"RverifySecurity", "RVERIFY", "SECURITY", func(s *Session, ctx context.Context) error { return s.RverifySecurity(ctx, "OBJ") }},
	}

	for _, entry := range entries {
		t.Run(entry.name, func(t *testing.T) {
			transport := newMockTransport()
			transport.addSuccessResponse()
			session := newTestSession(transport)

			err := entry.call(session, context.Background())
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			payload := transport.lastCall().Payload
			if payload["command"] != entry.command {
				t.Errorf("command = %v, want %s", payload["command"], entry.command)
			}
			if payload["qualifier"] != entry.qualifier {
				t.Errorf("qualifier = %v, want %s", payload["qualifier"], entry.qualifier)
			}
		})
	}
}
