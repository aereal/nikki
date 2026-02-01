package entrypoint

import "context"

type ExitCode int

func Run(e *Entrypoint, err error) int {
	if err != nil {
		return exitCodeOf(err)
	}

	defer func() { _ = e.tp.Shutdown(context.WithoutCancel(e.ctx)) }()
	return exitCodeOf(e.server.Start(e.ctx))
}
