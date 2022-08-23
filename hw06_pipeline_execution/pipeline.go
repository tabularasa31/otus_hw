package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

//TODO: need to finish
func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in
	select {
	case <-done:
		return out
	case <-in:
		for _, stage := range stages {
			out = stage(out)
		}
	}

	return out
}
