package utils

import (
	"hash/crc32"
	"os"
	"strconv"
	"sync"
)

func HashStr(str string) string {
	ieee := crc32.ChecksumIEEE([]byte(str))
	return strconv.FormatInt(int64(ieee), 16)
}

func Mkdir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_ = os.MkdirAll(path, 0777)
	}
}

type KeyMutex struct {
	mutexes sync.Map
}

func (m *KeyMutex) Lock(key string) func() {
	value, _ := m.mutexes.LoadOrStore(key, &sync.Mutex{})
	mtx := value.(*sync.Mutex)
	mtx.Lock()

	return func() { mtx.Unlock() }
}
