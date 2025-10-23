package config

import "log"
        

func SetupLogger() {
    log.SetFlags(log.LstdFlags | log.Lshortfile)
}