package main

import (
	"bufio"
	"fmt"
	"os"
	"quest/internal"
)

func main() {
	player := internal.InitGame()
	fmt.Println("Добро пожаловать в текстовую игру! \nОсмотритесь.")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		command := scanner.Text()
		result := internal.HandleCommand(player, command)
		fmt.Println(result)
		fmt.Print("> ")
	}
}
