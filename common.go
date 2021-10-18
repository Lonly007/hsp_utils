package hsp_utils

import (
	"sync"
)

type HspRWLock struct {
	sync.RWMutex
}
