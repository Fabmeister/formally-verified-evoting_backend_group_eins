package configuration

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Server        ServerConfig
	Database      DatabaseConfig
	Messaging     MessagingConfig
	Use_mock_data bool
}

func (this *Config) SetStandard() {
	this.Server.SetStandard()
	this.Database.SetStandard()
	this.Messaging.SetStandard()
	this.Use_mock_data = false
}

const std_config_path string = "./config.json"

var GlobalConfig Config

// This has to be called, otherwise GlobalConfig is uninitialized!
func Read_config(path string) Config {

	bytes, err := os.ReadFile(path)
	var conf Config

	if err != nil {
		log.Printf("Configpfad nicht lesbar: %v,versuche standard (./config.json)...", err)

		bytes, err = os.ReadFile(std_config_path)
		if err != nil {
			log.Printf("Konnte Standard Config nicht öffen! Grund: %v", err)

			conf.SetStandard()

			bytes, err = json.Marshal(conf)
			if err != nil {
				log.Fatalf("Fehler beim json string erstellen! %v", err)
			}

			err = os.WriteFile(std_config_path, bytes, 0770)

			if err != nil {
				log.Fatalf("Fehler beim schreiben der Beispiels config! %v", err)
			}

			log.Fatal("Beende Programm, bitte Config schreiben!")
		}
	}

	log.Println("Config geöffnet!")

	err = json.Unmarshal(bytes, &conf)

	if err != nil {
		log.Fatalf("Fehler beim parsen der Config! %v", err)
	}

	GlobalConfig = conf
	return conf
}
