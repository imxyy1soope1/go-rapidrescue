package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/imxyy1soope1/go-rapidrescue/pkg/bfs"
	. "github.com/imxyy1soope1/go-rapidrescue/pkg/json"
	"github.com/imxyy1soope1/go-rapidrescue/pkg/routing"
)

func main() {
	config, err := loadConfig()
	handleError(err)

	graph := bfs.NewGraph(config.Map, config.BrokenRoads, config.NonTuringPoints)
	data := &routing.Data{
		Carrying:       0,
		RequiredGoods:  config.RequiredGoods,
		MaterialPoints: config.MaterialPoints,
		Quarters:       config.Quarters,
	}
	routePlanner := routing.NewRoutePlanner(graph, data)
	routeResult := routePlanner.Plan().ToResult()

	err = saveResult(routeResult)
	handleError(err)

	fmt.Println("Done!")
}

func loadConfig() (Config, error) {
	if len(os.Args) != 3 {
		return Config{}, fmt.Errorf("Usage: %s [config.json] [result.json]", os.Args[0])
	}

	configFile, err := os.Open(os.Args[1])
	if err != nil {
		return Config{}, fmt.Errorf("failed to open config file: %w", err)
	}
	defer configFile.Close()

	configData, err := io.ReadAll(configFile)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	err = json.Unmarshal(configData, &config)
	if err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return config, nil
}

func saveResult(result interface{}) error {
	if _, err := os.Stat(path.Dir(os.Args[2])); err != nil {
		os.MkdirAll(path.Dir(os.Args[2]), os.ModePerm)
	}
	resultFile, err := os.Create(os.Args[2])
	if err != nil {
		return fmt.Errorf("failed to create result file: %w", err)
	}
	defer resultFile.Close()

	err = json.NewEncoder(resultFile).Encode(result)
	if err != nil {
		return fmt.Errorf("failed to encode result: %w", err)
	}

	return nil
}

func handleError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
