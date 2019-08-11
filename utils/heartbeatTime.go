/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2019/5/3
   Description :
-------------------------------------------------
*/

package utils

import (
    "sync/atomic"
    "time"
)

type heartbeatEvent func(timer *HeartbeatTime)

//心跳时间类型,不应该用值类型去做拷贝，而应该用指针
type HeartbeatTime struct {
    isrun int32
    reset int32
}

func NewHeartbeatTime(target time.Duration, precision time.Duration, callback heartbeatEvent) *HeartbeatTime {
    timer := &HeartbeatTime{
        isrun: 1,
    }
    go func() {
        var clock time.Duration
        for timer.IsRun() {
            time.Sleep(precision)

            // 要求重置计时
            if atomic.CompareAndSwapInt32(&timer.reset, 1, 0) {
                clock = 0
            }

            clock += precision
            if clock >= target {
                clock -= target
                callback(timer)
            }
        }
    }()
    return timer
}

func (m *HeartbeatTime) RefHeartbeat() {
    atomic.StoreInt32(&m.reset, 1)
}

func (m *HeartbeatTime) Stop() {
    atomic.StoreInt32(&m.isrun, 0)
}

func (m *HeartbeatTime) IsRun() bool {
    return atomic.LoadInt32(&m.isrun) == 1
}
