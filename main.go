package main

func main() {
	server := NewFromEnv(".env")
	server.Start()
}
