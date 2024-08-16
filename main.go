package main

func main() {
	server := createServer(":3000")
	//connectDB()
	server.run()
}
