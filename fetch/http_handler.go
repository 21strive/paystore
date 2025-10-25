package fetch

// HTTPCommandHandler
/*
	- FetchBalance
	- FetchTranscations
	- FetchPayment
	- FetchWithdraw
*/
type HTTPFetcherHandler struct {
	paystoreFetcher *PaystoreFetcher
}
