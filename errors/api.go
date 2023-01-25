package errors

import (
	"github.com/joomcode/errorx"
)

var ApiErrors = errorx.NewNamespace("api")
var InvalidWhereClauseError = ApiErrors.NewType("invalid_where_clause")
var CannotStopApiServiceError = ApiErrors.NewType("cannot_stop_api_service")
var CannotUpgradeWebsocket = ApiErrors.NewType("cannot_upgrade_websocket", fatalTrait)
var CannotWriteWebsocket = ApiErrors.NewType("cannot_write_websocket")
