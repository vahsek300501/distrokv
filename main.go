package main

import registry "github.com/Vahsek/distrokv/internal/registry"

func main() {
	registry.StartRegistryServer(":8080")
}
