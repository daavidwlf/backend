package main

func main() {
	server := createServer(":3000")
	server.run()
}
