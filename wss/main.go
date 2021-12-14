package wss

func Init(n int) {
	setPorts(n)
	setRoute()
	H.Run()
}
