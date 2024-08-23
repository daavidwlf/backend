package main

func main() {

	//change to env later
	port := "3000"

	server := createServer(":" + port)
	connectDB()
	server.run()
}
