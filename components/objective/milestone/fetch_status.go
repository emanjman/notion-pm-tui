package milestone

type FetchStatus int

const (
	FetchIdle FetchStatus = iota
	FetchPending
	FetchSuccess
	FetchFailed
)
