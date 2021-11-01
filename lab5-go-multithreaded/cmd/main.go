package main

import "github.com/borodun/nsu-nets/lab5/socks5"

func main() {
	conf := &socks5.Config{}
	server, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}

	if err := server.ListenAndServe("tcp", "127.0.0.1:11111"); err != nil {
		panic(err)
	}
}
