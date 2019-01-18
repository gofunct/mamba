package config

import (
	"github.com/fsnotify/fsnotify"
	"path/filepath"
	"github.com/prometheus/common/log"
	"sync"
)

func (v *config) OnconfigChange(run func(in fsnotify.Event)) {
	v.onconfigChange = run
}

func (v *config) Watchconfig() {
	initWG := sync.WaitGroup{}
	initWG.Add(1)
	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()
		// we have to watch the entire directory to pick up renames/atomic saves in a cross-platform way
		filename, err := v.getconfigFile()
		if err != nil {
			log.Debug("error: %v\n", err.Error())
			return
		}

		configFile := filepath.Clean(filename)
		configDir, _ := filepath.Split(configFile)
		realconfigFile, _ := filepath.EvalSymlinks(filename)

		eventsWG := sync.WaitGroup{}
		eventsWG.Add(1)
		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok { // 'Events' channel is closed
						eventsWG.Done()
						return
					}
					currentconfigFile, _ := filepath.EvalSymlinks(filename)
					// we only care about the config file with the following cases:
					// 1 - if the config file was modified or created
					// 2 - if the real path to the config file changed (eg: k8s configMap replacement)
					const writeOrCreateMask = fsnotify.Write | fsnotify.Create
					if (filepath.Clean(event.Name) == configFile &&
						event.Op&writeOrCreateMask != 0) ||
						(currentconfigFile != "" && currentconfigFile != realconfigFile) {
						realconfigFile = currentconfigFile
						err := v.ReadInconfig()
						if err != nil {
							log.Debug("error reading config file: %v\n", err)
						}
						if v.onconfigChange != nil {
							v.onconfigChange(event)
						}
					} else if filepath.Clean(event.Name) == configFile &&
						event.Op&fsnotify.Remove&fsnotify.Remove != 0 {
						eventsWG.Done()
						return
					}

				case err, ok := <-watcher.Errors:
					if ok { // 'Errors' channel is not closed
						log.Debug("watcher error: %v\n", err)
					}
					eventsWG.Done()
					return
				}
			}
		}()
		watcher.Add(configDir)
		initWG.Done()   // done initalizing the watch in this go routine, so the parent routine can move on...
		eventsWG.Wait() // now, wait for event loop to end in this go-routine...
	}()
	initWG.Wait() // make sure that the go routine above fully ended before returning
}
