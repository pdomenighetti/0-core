package settings

import (
	"fmt"
	"github.com/g8os/core/agent/lib/utils"
	"io/ioutil"
	"path"
)

const (
	//Init happens before handshake
	AfterInit = "init"

	//Core happens with core is up and running (also networking)
	AfterNet = "net"

	//Default for startup commands that doesn't specify dependency
	AfterBoot = "boot"
)

var (
	CyclicDependency = fmt.Errorf("cyclic dependency")

	Priority = map[string]int64{
		AfterInit: 1,
		AfterNet:  1000,
		AfterBoot: 1000000,
	}
)

type IncludedSettings struct {
	Extension map[string]Extension
	Startup   map[string]Startup
}

//GetPartialSettings loads partial settings according to main configurations
func (s *AppSettings) GetIncludedSettings() (partial *IncludedSettings, errors []error) {
	errors = make([]error, 0)

	partial = &IncludedSettings{
		Extension: make(map[string]Extension),
		Startup:   make(map[string]Startup),
	}

	if s.Main.Include == "" {
		return
	}

	infos, err := ioutil.ReadDir(s.Main.Include)
	if err != nil {
		errors = append(errors, fmt.Errorf("failed to read dir %s: %s", s.Main.Include, err))
		return
	}

	for _, info := range infos {
		if info.IsDir() {
			continue
		}
		name := info.Name()
		if len(name) <= len(ConfigSuffix) {
			//file name too short to be a config file (shorter than the extension)
			continue
		}
		if name[len(name)-len(ConfigSuffix):] != ConfigSuffix {
			continue
		}

		partialCfg := IncludedSettings{}
		partialPath := path.Join(s.Main.Include, name)

		err := utils.LoadTomlFile(partialPath, &partialCfg)
		if err != nil {
			errors = append(errors,
				fmt.Errorf("failed to load file %s: %s", partialPath, err))
			continue
		}

		//merge into settings
		for key, ext := range partialCfg.Extension {
			_, m := s.Extensions[key]
			_, p := partial.Extension[key]
			if m || p {
				errors = append(errors,
					fmt.Errorf("extension override in '%s' name '%s'", partialPath, key))
				continue
			}

			ext.key = key
			partial.Extension[key] = ext
		}

		for key, startup := range partialCfg.Startup {
			startup.key = key
			partial.Startup[key] = startup
		}
	}

	return
}
