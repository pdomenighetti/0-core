package pm

//implement internal processes

import (
    "github.com/shirou/gopsutil/process"
)

const (
    CMD_EXECUTE = "execute"
)

type ProcessConstructor func (cmd *Cmd) Process

var CMD_MAP = map[string]ProcessConstructor {
    CMD_EXECUTE: NewExtProcess,
}


func NewProcess(cmd *Cmd) Process {
    constructor, ok := CMD_MAP[cmd.Name]
    if !ok {
        return nil
    }

    return constructor(cmd)
}


type JsScriptProcess struct {
    extps Process
    cmd *Cmd
}

//Create a constructor for external process to execute an external script
//exe: name of executor, (python, lua, bash)
//workdir: working directory of script
//name: if name != "", execute that specific script, otherwise use args[name]
func extScript(exe string, workdir string, name string) ProcessConstructor {
    //create a new execute process with python2.7 or lua as executors.
    constructor := func(cmd *Cmd) Process {
        args := cmd.Args.Clone(false)
        var script string

        if name != "" {
            script = name
        } else {
            script = cmd.Args.GetString("name")
        }

        args.Set("name", exe)
        args.Set("args", []string{script})
        args.Set("working_dir", workdir)

        extcmd := &Cmd {
            Id: cmd.Id,
            Gid: cmd.Gid,
            Nid: cmd.Nid,
            Name: CMD_EXECUTE,
            Data: cmd.Data,
            Args: args,
        }

        return &JsScriptProcess{
            extps: NewExtProcess(extcmd),
            cmd: cmd,
        }
    }

    return constructor
}


func (ps *JsScriptProcess) Run(cfg RunCfg) {
    //intercept all the messages from the 'execute' command and
    //change it to it's original value.
    extcfg := RunCfg{
        ProcessManager: cfg.ProcessManager,
        MeterHandler: func(cmd *Cmd, p *process.Process) {
            cfg.MeterHandler(ps.cmd, p)
            },
        MessageHandler: func(msg *Message) {
            msg.Cmd = ps.cmd
            cfg.MessageHandler(msg)
            },
        ResultHandler: func(result *JobResult) {
            result.Args = ps.cmd.Args
            result.Cmd = ps.cmd.Name
            result.Id = ps.cmd.Id
            result.Gid = ps.cmd.Gid
            result.Nid = ps.cmd.Nid

            cfg.ResultHandler(result)
            },
        Signal: cfg.Signal,
    }

    ps.extps.Run(extcfg)
}

func (ps *JsScriptProcess) Kill() {
    ps.extps.Kill()
}


//registers a command to the process manager.
//cmd: Command name
//exe: executing binary
//workdir: working directory
//script: script name
//if script == "", then script name will be used from cmd.Args.
func RegisterCmd(cmd string, exe string, workdir string, script string) {
    CMD_MAP[cmd] = extScript(exe, workdir, script)
}