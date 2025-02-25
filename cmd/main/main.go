package main

func main() {
	setupCli()

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func setupCli() {
	rootCmd.setup()

	rootCmd.AddCommand(versionCmd)
}
