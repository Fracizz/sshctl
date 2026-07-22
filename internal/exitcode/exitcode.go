package exitcode

// Process exit codes for sshctl.
const (
	OK          = 0 // success
	Usage       = 2 // invalid usage / local config error
	ExecFailed  = 1 // local runtime failure (dial, IO, decrypt, …)
	// Remote command failures: sshctl exec exits with the remote process status when available.
)
