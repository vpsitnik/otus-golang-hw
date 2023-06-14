package hw06pipelineexecution

import "fmt"

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in

	for _, stage := range stages {
		outRC := make(Bi)

		outStage := stage(out)

		go func() {
			defer close(outRC)
			for {
				select {
				case <-done:
					return
				case value, ok := <-outStage:
					if !ok {
						fmt.Println("Not ok, exiting!")
						return
					}
					outRC <- value
				}
			}
		}()

		out = outRC
	}
	return out
}
