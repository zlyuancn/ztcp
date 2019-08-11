/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2019/5/1
   Description :
-------------------------------------------------
*/

package utils

import (
    "sync/atomic"
)

type AutoId struct {
    nextId uint64
}

func (m *AutoId) Next() uint64 {
    return atomic.AddUint64(&m.nextId, 1)
}

var AutoClientID AutoId
