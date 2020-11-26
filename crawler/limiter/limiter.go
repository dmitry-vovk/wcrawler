package limiter

var limitC chan struct{}

func SetLimit(limit int) {
	limitC = make(chan struct{}, limit)
}

func Start() {
	limitC <- struct{}{}
}

func Finish() {
	<-limitC
}
