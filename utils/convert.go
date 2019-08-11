/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2019/5/2
   Description :
-------------------------------------------------
*/

package utils

// uint64转为bytes, 从右边开始写入
func Uint64ToBytes(v uint64) []byte {
    return []byte{
        byte(v >> 56), byte(v >> 48), byte(v >> 40), byte(v >> 32), byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v),
    }
}

// bytes转为uint64, 从右边开始读取
func BytesToUint64(b []byte) uint64 {
    return uint64(b[7]) | uint64(b[6])<<8 | uint64(b[5])<<16 | uint64(b[4])<<24 |
        uint64(b[3])<<32 | uint64(b[2])<<40 | uint64(b[1])<<48 | uint64(b[0])<<56
}

// uint32转为bytes, 从右边开始写入
func Uint32ToBytes(v uint32) []byte {
    return []byte{
        byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v),
    }
}

// bytes转为uint32, 从右边开始读取
func BytesToUint32(b []byte) uint32 {
    return uint32(b[3]) | uint32(b[2])<<8 | uint32(b[1])<<16 | uint32(b[0])<<24
}
