package uuid

import (
	"fmt"
	"sync"
	"time"
)

var (
	uniqueIdMutex      sync.Mutex
	lastUniqueIdMicros int64
)

// UuidGenerator ids unicos para correlacionar entradas de log.
type UuidGenerator struct{}

func GetInstance() *UuidGenerator {
	return &UuidGenerator{}
}

// GetUniqueIdWithPrefix prefix + uniqid() de PHP: 8 hex de los segundos + 5 hex de los microsegundos.
func (u *UuidGenerator) GetUniqueIdWithPrefix(prefix string) string {
	uniqueIdMutex.Lock()
	defer uniqueIdMutex.Unlock()

	micros := time.Now().UnixMicro()
	if micros <= lastUniqueIdMicros {
		micros = lastUniqueIdMicros + 1
	}
	lastUniqueIdMicros = micros
	return fmt.Sprintf("%s%08x%05x", prefix, micros/1_000_000, micros%1_000_000)
}
