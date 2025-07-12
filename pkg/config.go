package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Addr string `yaml:"address" env:"HTTP_SERVER_ADDRESS" env-required:"true"`
}

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath string `yaml:"stroage_path" env:"STORAGE_PATH" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

func MustLoad() *Config {
	// It means we are getting value from the environment variable
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		// flag is basically use to take the value from the command line. Means maybe user is passing the config file path as a command line argument
		// Note Basically flag is gonna be a pointer.
		flags := flag.String("config", "", "Path to the configuration file")
		// Then parse each and every flag provided by the user.
		flag.Parse()

		configPath = *flags // Then destructure and get the value of the flag

		if configPath == "" {
			log.Fatal("Config Path is not set.") // It will log and exit the os in the cmd.
		}
	}

	//Stat returns a [FileInfo] describing the named file. If there is an error, it will be of type [*PathError].
	//It is saying that if the file does not exist, then it will return an error of type [*PathError].
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist at path: %s", configPath)
	}

	var cfg Config

	// cleanenv is a library that helps to read configuration files and environment variables.
	// It will read the config file and store it in the cfg variable.
	// It will return the error if there is any issue in reading the config file and we are storing that.
	err := cleanenv.ReadConfig(configPath, &cfg)

	if err != nil {
		log.Fatalf("Failed to read config file: %s", err)
	}

	return &cfg

}
