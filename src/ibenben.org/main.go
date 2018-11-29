package main

import "ibenben.org/blog/csdn"

func main() {
	csdn.Pages().QueryAndParse()
	//csdn.NewAcquirer("https://blog.csdn.net/jrainbow/article/details/8912273").Parse()
	//csdn.NewAcquirer("https://blog.csdn.net/jrainbow/article/details/51980036").Parse()
	//csdn.NewAcquirer("https://blog.csdn.net/jrainbow/article/details/83339972").Parse()
}

