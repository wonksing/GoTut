package main

import (
	"fmt"

	"github.com/spf13/viper"
)

func main() {
	fmt.Println("hello")
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // path to look for the config file in
	// viper.AddConfigPath("$HOME/.config") // call multiple times to add many search paths
	// viper.AddConfigPath(".")              // optionally look for config in the working directory
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	fmt.Printf("jacket is %s \n", viper.GetString("clothing.jacket"))
}
