package domain

import (
	"github.com/yeencloud/lib-shared/namespace"
)

var (
	LogFieldSQL        = namespace.Namespace{Identifier: "sql"}
	LogFieldSQLRequest = namespace.Namespace{Parent: &LogFieldSQL, Identifier: "request"}

	LogFieldSQLRowsAffected = namespace.Namespace{Parent: &LogFieldSQLRequest, Identifier: "rows_affected"}
	LogFieldSQLQuery        = namespace.Namespace{Parent: &LogFieldSQLRequest, Identifier: "query"}

	LogScopeTime        = namespace.Namespace{Parent: &LogFieldSQLRequest, Identifier: "time"}
	LogFieldTimeStarted = namespace.Namespace{Parent: &LogScopeTime, Identifier: "started"}
	LogFieldTimeEnded   = namespace.Namespace{Parent: &LogScopeTime, Identifier: "ended"}
	LogFieldDuration    = namespace.Namespace{Parent: &LogScopeTime, Identifier: "duration_ms"}
)

const SQLMetricPointName = "sql"
