package main

import (
	"fmt"
	"sync"
)

/*
SINGLETON
*/
type Config struct {
	AppName string
}

var instance *Config
var once sync.Once

func GetInstance() *Config {
	once.Do(func() {
		instance = &Config{
			AppName: "MyApp",
		}
	})
	return instance
}

func main() {
	c1 := GetInstance()
	c2 := GetInstance()

	fmt.Println(c1 == c2) // true
	fmt.Println(c1.AppName)
}
