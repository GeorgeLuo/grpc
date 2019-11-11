package models

import (
	"errors"
	"fmt"
	"sync"
)

// AliasMap is used to map an alias to a task_id, this alias can be used to
// provide a human readable or programmatically sensical string by a client
// to handle a process. Keys in the alias map 1-1 to the SyncMap, but will
// eventually be mappable to a set of processes.
type AliasMap struct {
	aliasMap map[string][]string
	mutex    *sync.Mutex
}

// NewAliasMap is used to return an empty AliasMap.
func NewAliasMap() AliasMap {
	return AliasMap{
		aliasMap: make(map[string][]string),
		mutex:    &sync.Mutex{},
	}
}

// Put is an operation insert a new mapping of alias to taskIDs.
func (rwm *AliasMap) Put(alias string, taskIDs ...string) error {
	rwm.mutex.Lock()
	defer rwm.mutex.Unlock()
	if alias == "" {
		return errors.New("illegal empty alias")
	}
	if taskID, ok := rwm.aliasMap[alias]; ok {
		return fmt.Errorf("alias [%s] already mapped to taskId %s", alias, taskID)
	}
	rwm.aliasMap[alias] = taskIDs
	return nil
}

// Get returns the task ids mapped to an alias.
func (rwm *AliasMap) Get(alias string) (prcs []string, ok bool) {
	rwm.mutex.Lock()
	defer rwm.mutex.Unlock()
	prcs, ok = rwm.aliasMap[alias]
	return
}
