package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if in == nil {
		out := make(Bi)
		close(out)
		return out
	}

	for _, stage := range stages {
		in = execStage(in, done, stage)
	}

	return in
}

func execStage(in In, done In, stage Stage) Bi {
	out := make(Bi)

	go func() {
		defer close(out)

		in = stage(in)

		for {
			select {
			case <-done:
				return
			case val, ok := <-in:
				if !ok {
					return
				}
				out <- val
			}
		}
	}()

	return out
}
