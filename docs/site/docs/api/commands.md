# Command Methods

## Overview

`Session` provides ~144 command methods, one for each MQSC command verb +
qualifier combination. Each method is a thin wrapper that calls the internal
command dispatcher with the correct verb and qualifier. Method names follow the
pattern `VerbQualifier` in PascalCase, mapping directly to MQSC commands (e.g.
`DISPLAY QUEUE` becomes `DisplayQueue()`).

## Method signature patterns

### DISPLAY commands (list return)

```go
func (session *Session) DisplayQueue(
    ctx  context.Context,
    name string,
    opts ...CommandOption,
) ([]map[string]any, error)
```

### DISPLAY commands (singleton return)

Queue manager singletons return a single `map[string]any` instead of a slice:

```go
func (session *Session) DisplayQmgr(
    ctx  context.Context,
    opts ...CommandOption,
) (map[string]any, error)
```

### Non-DISPLAY commands (error-only return)

```go
func (session *Session) DefineQlocal(
    ctx  context.Context,
    name string,
    opts ...CommandOption,
) error
```

## CommandOption

All command methods accept variadic `CommandOption` functional options:

| Option | Description |
| --- | --- |
| `WithRequestParameters(map[string]any)` | MQSC command parameters (attributes to set or filter on) |
| `WithResponseParameters([]string)` | Attribute names to include in the response (defaults to `["all"]` for DISPLAY) |
| `WithWhere(string)` | WHERE clause to filter DISPLAY command results |

```go
ctx := context.Background()

// DISPLAY with default options (all attributes, wildcard name)
queues, err := session.DisplayQueue(ctx, "*")

// DISPLAY with request parameters and response filtering
queues, err := session.DisplayQueue(ctx, "APP.*",
    mqrestadmin.WithResponseParameters([]string{"current_queue_depth", "max_queue_depth"}),
    mqrestadmin.WithWhere("current_queue_depth GT 100"),
)

// DEFINE with attributes
err = session.DefineQlocal(ctx, "MY.QUEUE",
    mqrestadmin.WithRequestParameters(map[string]any{
        "max_queue_depth": 50000,
        "description":     "Application request queue",
    }),
)

// DELETE (no options needed)
err = session.DeleteQueue(ctx, "MY.QUEUE")
```

## Return values

- **DISPLAY commands (list)**: `([]map[string]any, error)` -- one map per
  matched object. A `nil` or empty slice means no objects matched (not an error).
- **Queue manager singletons** (`DisplayQmgr`, `DisplayQmstatus`,
  `DisplayCmdserv`): `(map[string]any, error)` -- returns `nil` map if empty.
- **Non-DISPLAY commands**: `error` on failure, `nil` on success.

## DISPLAY methods

| Method | MQSC command | Qualifier mapping |
| --- | --- | --- |
| `DisplayApstatus()` | `DISPLAY APSTATUS` | apstatus |
| `DisplayArchive()` | `DISPLAY ARCHIVE` | archive |
| `DisplayAuthinfo()` | `DISPLAY AUTHINFO` | authinfo |
| `DisplayAuthrec()` | `DISPLAY AUTHREC` | authrec |
| `DisplayAuthserv()` | `DISPLAY AUTHSERV` | authserv |
| `DisplayCfstatus()` | `DISPLAY CFSTATUS` | cfstatus |
| `DisplayCfstruct()` | `DISPLAY CFSTRUCT` | cfstruct |
| `DisplayChannel()` | `DISPLAY CHANNEL` | channel |
| `DisplayChinit()` | `DISPLAY CHINIT` | chinit |
| `DisplayChlauth()` | `DISPLAY CHLAUTH` | chlauth |
| `DisplayChstatus()` | `DISPLAY CHSTATUS` | chstatus |
| `DisplayClusqmgr()` | `DISPLAY CLUSQMGR` | clusqmgr |
| `DisplayCmdserv()` | `DISPLAY CMDSERV` | cmdserv |
| `DisplayComminfo()` | `DISPLAY COMMINFO` | comminfo |
| `DisplayConn()` | `DISPLAY CONN` | conn |
| `DisplayEntauth()` | `DISPLAY ENTAUTH` | entauth |
| `DisplayGroup()` | `DISPLAY GROUP` | group |
| `DisplayListener()` | `DISPLAY LISTENER` | listener |
| `DisplayLog()` | `DISPLAY LOG` | log |
| `DisplayLsstatus()` | `DISPLAY LSSTATUS` | lsstatus |
| `DisplayMaxsmsgs()` | `DISPLAY MAXSMSGS` | maxsmsgs |
| `DisplayNamelist()` | `DISPLAY NAMELIST` | namelist |
| `DisplayPolicy()` | `DISPLAY POLICY` | policy |
| `DisplayProcess()` | `DISPLAY PROCESS` | process |
| `DisplayPubsub()` | `DISPLAY PUBSUB` | pubsub |
| `DisplayQmgr()` | `DISPLAY QMGR` | qmgr |
| `DisplayQmstatus()` | `DISPLAY QMSTATUS` | qmgr |
| `DisplayQstatus()` | `DISPLAY QSTATUS` | qstatus |
| `DisplayQueue()` | `DISPLAY QUEUE` | queue |
| `DisplaySbstatus()` | `DISPLAY SBSTATUS` | sbstatus |
| `DisplaySecurity()` | `DISPLAY SECURITY` | security |
| `DisplayService()` | `DISPLAY SERVICE` | service |
| `DisplaySmds()` | `DISPLAY SMDS` | smds |
| `DisplaySmdsconn()` | `DISPLAY SMDSCONN` | smdsconn |
| `DisplayStgclass()` | `DISPLAY STGCLASS` | stgclass |
| `DisplaySub()` | `DISPLAY SUB` | sub |
| `DisplaySvstatus()` | `DISPLAY SVSTATUS` | svstatus |
| `DisplaySystem()` | `DISPLAY SYSTEM` | system |
| `DisplayTcluster()` | `DISPLAY TCLUSTER` | tcluster |
| `DisplayThread()` | `DISPLAY THREAD` | thread |
| `DisplayTopic()` | `DISPLAY TOPIC` | topic |
| `DisplayTpstatus()` | `DISPLAY TPSTATUS` | tpstatus |
| `DisplayTrace()` | `DISPLAY TRACE` | trace |
| `DisplayUsage()` | `DISPLAY USAGE` | usage |

## DEFINE methods

| Method | MQSC command | Qualifier mapping |
| --- | --- | --- |
| `DefineAuthinfo()` | `DEFINE AUTHINFO` | authinfo |
| `DefineBuffpool()` | `DEFINE BUFFPOOL` | buffpool |
| `DefineCfstruct()` | `DEFINE CFSTRUCT` | cfstruct |
| `DefineChannel()` | `DEFINE CHANNEL` | channel |
| `DefineComminfo()` | `DEFINE COMMINFO` | comminfo |
| `DefineListener()` | `DEFINE LISTENER` | listener |
| `DefineLog()` | `DEFINE LOG` | log |
| `DefineMaxsmsgs()` | `DEFINE MAXSMSGS` | maxsmsgs |
| `DefineNamelist()` | `DEFINE NAMELIST` | namelist |
| `DefineProcess()` | `DEFINE PROCESS` | process |
| `DefinePsid()` | `DEFINE PSID` | psid |
| `DefineQalias()` | `DEFINE QALIAS` | queue |
| `DefineQlocal()` | `DEFINE QLOCAL` | queue |
| `DefineQmodel()` | `DEFINE QMODEL` | queue |
| `DefineQremote()` | `DEFINE QREMOTE` | queue |
| `DefineService()` | `DEFINE SERVICE` | service |
| `DefineStgclass()` | `DEFINE STGCLASS` | stgclass |
| `DefineSub()` | `DEFINE SUB` | sub |
| `DefineTopic()` | `DEFINE TOPIC` | topic |

## DELETE methods

| Method | MQSC command | Qualifier mapping |
| --- | --- | --- |
| `DeleteAuthinfo()` | `DELETE AUTHINFO` | authinfo |
| `DeleteAuthrec()` | `DELETE AUTHREC` | authrec |
| `DeleteBuffpool()` | `DELETE BUFFPOOL` | buffpool |
| `DeleteCfstruct()` | `DELETE CFSTRUCT` | cfstruct |
| `DeleteChannel()` | `DELETE CHANNEL` | channel |
| `DeleteComminfo()` | `DELETE COMMINFO` | comminfo |
| `DeleteListener()` | `DELETE LISTENER` | listener |
| `DeleteNamelist()` | `DELETE NAMELIST` | namelist |
| `DeletePolicy()` | `DELETE POLICY` | policy |
| `DeleteProcess()` | `DELETE PROCESS` | process |
| `DeletePsid()` | `DELETE PSID` | psid |
| `DeleteQueue()` | `DELETE QUEUE` | queue |
| `DeleteService()` | `DELETE SERVICE` | service |
| `DeleteStgclass()` | `DELETE STGCLASS` | stgclass |
| `DeleteSub()` | `DELETE SUB` | sub |
| `DeleteTopic()` | `DELETE TOPIC` | topic |

## ALTER methods

| Method | MQSC command | Qualifier mapping |
| --- | --- | --- |
| `AlterAuthinfo()` | `ALTER AUTHINFO` | authinfo |
| `AlterBuffpool()` | `ALTER BUFFPOOL` | buffpool |
| `AlterCfstruct()` | `ALTER CFSTRUCT` | cfstruct |
| `AlterChannel()` | `ALTER CHANNEL` | channel |
| `AlterComminfo()` | `ALTER COMMINFO` | comminfo |
| `AlterListener()` | `ALTER LISTENER` | listener |
| `AlterNamelist()` | `ALTER NAMELIST` | namelist |
| `AlterProcess()` | `ALTER PROCESS` | process |
| `AlterPsid()` | `ALTER PSID` | psid |
| `AlterQmgr()` | `ALTER QMGR` | qmgr |
| `AlterSecurity()` | `ALTER SECURITY` | security |
| `AlterService()` | `ALTER SERVICE` | service |
| `AlterSmds()` | `ALTER SMDS` | smds |
| `AlterStgclass()` | `ALTER STGCLASS` | stgclass |
| `AlterSub()` | `ALTER SUB` | sub |
| `AlterTopic()` | `ALTER TOPIC` | topic |
| `AlterTrace()` | `ALTER TRACE` | trace |

## SET methods

| Method | MQSC command | Qualifier mapping |
| --- | --- | --- |
| `SetArchive()` | `SET ARCHIVE` | archive |
| `SetAuthrec()` | `SET AUTHREC` | authrec |
| `SetChlauth()` | `SET CHLAUTH` | chlauth |
| `SetLog()` | `SET LOG` | log |
| `SetPolicy()` | `SET POLICY` | policy |
| `SetSystem()` | `SET SYSTEM` | system |

## START methods

| Method | MQSC command | Qualifier mapping |
| --- | --- | --- |
| `StartChannel()` | `START CHANNEL` | channel |
| `StartChinit()` | `START CHINIT` | chinit |
| `StartCmdserv()` | `START CMDSERV` | cmdserv |
| `StartListener()` | `START LISTENER` | listener |
| `StartQmgr()` | `START QMGR` | qmgr |
| `StartService()` | `START SERVICE` | service |
| `StartSmdsconn()` | `START SMDSCONN` | smdsconn |
| `StartTrace()` | `START TRACE` | trace |

## STOP methods

| Method | MQSC command | Qualifier mapping |
| --- | --- | --- |
| `StopChannel()` | `STOP CHANNEL` | channel |
| `StopChinit()` | `STOP CHINIT` | chinit |
| `StopCmdserv()` | `STOP CMDSERV` | cmdserv |
| `StopConn()` | `STOP CONN` | conn |
| `StopListener()` | `STOP LISTENER` | listener |
| `StopQmgr()` | `STOP QMGR` | qmgr |
| `StopService()` | `STOP SERVICE` | service |
| `StopSmdsconn()` | `STOP SMDSCONN` | smdsconn |
| `StopTrace()` | `STOP TRACE` | trace |

## Other methods

| Method | MQSC command | Qualifier mapping |
| --- | --- | --- |
| `ArchiveLog()` | `ARCHIVE LOG` | log |
| `BackupCfstruct()` | `BACKUP CFSTRUCT` | cfstruct |
| `ClearQlocal()` | `CLEAR QLOCAL` | queue |
| `ClearTopicstr()` | `CLEAR TOPICSTR` | topicstr |
| `MoveQlocal()` | `MOVE QLOCAL` | queue |
| `PingChannel()` | `PING CHANNEL` | channel |
| `PingQmgr()` | `PING QMGR` | qmgr |
| `PurgeChannel()` | `PURGE CHANNEL` | channel |
| `RecoverBsds()` | `RECOVER BSDS` | bsds |
| `RecoverCfstruct()` | `RECOVER CFSTRUCT` | cfstruct |
| `RefreshCluster()` | `REFRESH CLUSTER` | cluster |
| `RefreshQmgr()` | `REFRESH QMGR` | qmgr |
| `RefreshSecurity()` | `REFRESH SECURITY` | security |
| `ResetCfstruct()` | `RESET CFSTRUCT` | cfstruct |
| `ResetChannel()` | `RESET CHANNEL` | channel |
| `ResetCluster()` | `RESET CLUSTER` | cluster |
| `ResetQmgr()` | `RESET QMGR` | qmgr |
| `ResetQstats()` | `RESET QSTATS` | queue |
| `ResetSmds()` | `RESET SMDS` | smds |
| `ResetTpipe()` | `RESET TPIPE` | tpipe |
| `ResolveChannel()` | `RESOLVE CHANNEL` | channel |
| `ResolveIndoubt()` | `RESOLVE INDOUBT` | indoubt |
| `ResumeQmgr()` | `RESUME QMGR` | qmgr |
| `RverifySecurity()` | `RVERIFY SECURITY` | security |
| `SuspendQmgr()` | `SUSPEND QMGR` | qmgr |

!!! note
    The full list of command methods is generated from the mapping data.
    See the [Qualifier Mapping Reference](../mappings/index.md) for per-qualifier
    details including attribute names and value mappings for each object type.
