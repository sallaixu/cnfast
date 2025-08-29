package help

import "fmt"

func PrintHelp() {
	fmt.Println("Usage: cnfast <command> [arguments]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  clone <repository>   Clone a GitHub repository with acceleration")
	fmt.Println("  pull <image>         Pull a Docker image with acceleration")
	fmt.Println("  version              Show version information")
	fmt.Println("  help                 Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  cnfast clone https://github.com/user/repo.git")
	fmt.Println("  cnfast pull nginx:latest")
	fmt.Println("  cnfast version")
}
