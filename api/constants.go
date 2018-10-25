package api

// Constants used in the api package
const (
	RouteSessionInit     = "/session/init"
	RouteLedgerState     = "/ledger/state"
	RouteTransactionList = "/transaction/list"
	RouteTransactionPoll = "/transaction/poll"
	RouteTransactionSend = "/transaction/send"
	RouteStatsReset      = "/stats/reset"
	RouteAccountLoad     = "/account/load"
	RouteAccountPoll     = "/account/poll"
	RouteServerVersion   = "/server/version"

	HeaderSessionToken      = "X-Session-Token"
	HeaderWebsocketProtocol = "Sec-Websocket-Protocol"
	HeaderUserAgent         = "User-Agent"

	MaxAllowableSessions = 50000
	MaxRequestBodySize   = 4 * 1024 * 1024
	MaxTimeOffsetInMs    = 5000
)
