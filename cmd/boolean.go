package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gadumitrachioaiei/boolean"

	_ "net/http/pprof"
)

func main() {
	text := flag.String("text", "", "boolean expression")
	values := flag.String("values", "", "values of used boolean operands")
	flag.Parse()
	if *text == "" || *values == "" {
		log.Fatalln("text and values are mandatory")
	}
	data := make(map[string]bool)
	bools := strings.Fields(*values)
	for i := 0; i < len(bools)-1; i = i + 2 {
		value, err := strconv.ParseBool(bools[i+1])
		if err != nil {
			log.Fatalf("%s has to be true or false but it is: %s\n", bools[i], bools[i+1])
		}
		data[bools[i]] = value
	}
	tree, err := boolean.New().Parse(*text)
	if err != nil {
		log.Fatalln(err)
	}
	result, err := tree.Execute(data)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(result)
}
