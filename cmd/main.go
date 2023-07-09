package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/imxyy1soope1/go-rapidrescue/pkg/bfs"
	. "github.com/imxyy1soope1/go-rapidrescue/pkg/json"
	"github.com/imxyy1soope1/go-rapidrescue/pkg/routing"
)

func main() {
	var (
		configFile string
		config     Config
		errHandler = func(err error) {
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}
	)

	if len(os.Args) > 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [config.json]\n", os.Args[0])
		os.Exit(1)
	} else if len(os.Args) == 2 {
		configFile = os.Args[1]
	} else {
		configFile = "./config.json"
	}

	f, err := os.Open(configFile)
	errHandler(err)
	defer f.Close()

	configData, err := io.ReadAll(f)
	errHandler(err)

	err = json.Unmarshal(configData, &config)
	errHandler(err)

	graph := bfs.NewGraph(config.Map, config.BrokenRoads, config.NonTuringPoints)
	data := &routing.Data{
		Carrying:       0,
		RequiredGoods:  config.RequiredGoods,
		MaterialPoints: config.MaterialPoints,
		Quarters:       config.Quarters,
	}
	routePlanner := routing.NewRoutePlanner(graph, data)
	routeResult := routePlanner.Plan().ToResult()

	resultFile, err := os.Create("result.json")
	errHandler(err)
	defer resultFile.Close()

	err = json.NewEncoder(resultFile).Encode(routeResult)
	errHandler(err)

	fmt.Println("Done!")
}
