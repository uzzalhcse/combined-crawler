package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"
)

func PluginGenerate() {
	path := "plugins/test/url_handler"
	filename := "TestHandler.go"
	functionName := "Sum"
	src := filepath.Join(path, filename)
	pluginPath := filepath.Join(path, functionName+".so")

	if err := generatePlugin(src, pluginPath); err != nil {
		fmt.Printf("failed to build plugin %s: %v\n", functionName, err)
	}

	fmt.Println("pluginPath", pluginPath)

	fnSymbol, err := loadHandler(pluginPath, functionName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Signature: %T\n", fnSymbol)
	switch v := fnSymbol.(type) {
	case func(int, int):
		fn, ok := fnSymbol.(func(a, b int))
		if !ok {
			fmt.Printf("Function %s has unexpected type in package\n", functionName)
		}
		fn(5, 6)
	default:
		fmt.Printf("Invalid Signature: %T\n", v)
	}
}

func generatePlugin(src, dest string) error {
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", dest, src)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	return nil
}
func loadHandler(pluginPath, funcName string) (plugin.Symbol, error) {
	// Load the package
	pkg, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, fmt.Errorf("error opening package: %w", err)
	}

	// Look for the function by name in the package
	fnSymbol, err := pkg.Lookup(funcName)
	if err != nil {
		return nil, fmt.Errorf("function %s not found in package %s: %w", funcName, pluginPath, err)
	}

	return fnSymbol, nil
}
