package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"controller"
)

func main() {
	var err error
	defaultConfigPath := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	configPath := flag.String("config", defaultConfigPath, "path to kube config")
	flag.Parse()

	if *configPath, err = filepath.Abs(*configPath); err != nil {
		panic(err.Error())
	}

	controller := internfed.MakeOutOfCluster(*configPath)

	fmt.Println("Successfully created the kubernetes controller:", controller)
	fmt.Println("--------------------")

	reader := bufio.NewReader(os.Stdin)
shell:
	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		words := strings.Fields(text)

		switch words[0] {
		case "quit":
			break shell
		case "help":
			fmt.Println("Enter any kubectl command (enter 'kubectl help' to see all kubectl commands)\nEnter 'quit' to exit")
		case "create":
			if words[1] != "-f" {
				log.Fatal("Cannot create without -f")
			}
			data, err := ioutil.ReadFile(words[2])
			if err != nil {
				log.Fatal(err)
			}
			var t struct{
				Kind string
			}
			if err = yaml.Unmarshal([]byte(data), &t); err != nil {
				log.Fatal(err)
			}
			controller.Create(t.Kind, []byte(data))
		default:
		}
	}
}