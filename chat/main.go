package chat

func Init(n int) {
	setPorts(n)
	initRoute()
	H.Run()
}
