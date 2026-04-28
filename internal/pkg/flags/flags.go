package flags

import "flag"

type Flags struct {
    Path string
}

func Parse() *Flags {
    configPath := flag.String("config", "./config/config.yaml", "path to config file")
    help := flag.Bool("help", false, "show help message")
    
    flag.Parse()
    
    if *help {
        flag.Usage()
    }
    
    return &Flags{
        Path: *configPath,
    }
}