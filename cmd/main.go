package main

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("%+v", err)
	}
}
