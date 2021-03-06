package pm

import "fmt"

//implement internal processes

/*
Global command ProcessConstructor registery
*/
var factories = map[string]ProcessFactory{
	CommandSystem: NewSystemProcess,
}

/*
NewProcess creates a new process from a command
*/
func GetProcessFactory(cmd *Command) ProcessFactory {
	return factories[cmd.Command]
}

func Register(name string, factory ProcessFactory) {
	if _, ok := factories[name]; ok {
		panic(fmt.Sprintf("command registered with same name: %s", name))
	}
	factories[name] = factory
}

/*
RegisterExtension registers a new command (extension) so it can be executed via commands
*/
func RegisterExtension(cmd string, exe string, workdir string, cmdargs []string, env map[string]string) error {
	if _, ok := factories[cmd]; ok {
		return fmt.Errorf("job factory with the same name already registered: %s", cmd)
	}

	Register(cmd, extensionProcessFactory(exe, workdir, cmdargs, env))
	return nil
}

func RegisterBuiltIn(name string, runnable Runnable) {
	Register(name, internalProcessFactory(runnable))
}

func RegisterBuiltInWithCtx(name string, runnable RunnableWithCtx) {
	Register(name, internalProcessFactoryWithCtx(runnable))
}
