package agent

//logger settings
type Logger struct {
    //logger type, now only 'db' and 'ac' are supported
    Type string
    //list of controlles base URLs
    AgentControllers []int
    //Process which levels
    Levels []int

    //Log dir (for loggers that needs it)
    LogDir string
    //Flush interval (for loggers that needs it)
    FlushInt int
    //Flush batch size (for loggers that needs it)
    BatchSize int
}

//external cmd config
type Cmd struct {
    //binary to execute
    Binary string
    //script search path
    Cwd string
    //(optional) if set, overrides the cmd.Args['name']
    Script string
    //(optional) Env variables
    Env map[string]string
}

//main agent settings
type Settings struct {
    Main struct {
        Gid int
        Nid int
        MaxJobs int
        AgentControllers []string
        MessageIdFile string
        HistoryFile string
    }

    Cmds map[string]Cmd

    Logging map[string]Logger

    Stats struct {
        Interval int
        AgentControllers []int
    }

    Channel struct {
        Cmds []int
    }

    Security struct {
        CertificateAuthority string
        ClientCertificate string
        ClientCertificateKey string
    }
}
