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

	if stages == nil {
		return in
	}

	res := in

	for _, stage := range stages {
		res = stage(chanStage(res, done))
	}

	return res
}

func chanStage(in In, done In) Bi {
	chanStage := make(Bi)

	go func() {
		defer func() {
			close(chanStage)
			//nolint:revive
			for range in {
			}
		}()

		for {
			select {
			case v, ok := <-in:
				if !ok {
					return
				}
				select {
				case chanStage <- v:
				case <-done:
					return
				}
			case <-done:
				return
			}
		}
	}()

	return chanStage
}
