/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2019/5/2
   Description :
-------------------------------------------------
*/

package utils

import (
    "fmt"
    "github.com/zlyuancn/zassert"
    "github.com/zlyuancn/ztcp/config"
    "net"
)
var DataBuffSize = config.DefaultDataBuffSize
var DataPackageBuffSize = config.DefaultDataPackageBuffSize

// 等待一次指定长度的数据(已连接的conn, 数据总长度, 单次数据缓存大小)
func WaitConnData(conn net.Conn, length int) ([] byte, error) {
    fullbuff := make([]byte, length)

    var index int
    var size int
    for index < length {
        size = length - index
        if size > DataBuffSize {
            size = DataBuffSize
        }

        buff := fullbuff[index : index+size]
        le, err := conn.Read(buff)
        index += le

        if err != nil {
            return fullbuff[:index], err
        }
    }
    return fullbuff, nil
}

// 等待一个完整的数据
func WaitConnFullData(conn net.Conn) ([] byte, error) {
    dataHeader, err := WaitConnData(conn, config.DataHeaderLength)
    if err != nil {
        return dataHeader, err
    }

    dataSize := int(BytesToUint32(dataHeader))
    if dataSize == 0 {
        return dataHeader, nil
    }
    if dataSize >= DataPackageBuffSize {
        return dataHeader, zassert.AssertError{fmt.Sprintf("数据长度超过设置的 %d bytes", DataPackageBuffSize)}
    }

    return WaitConnData(conn, dataSize)
}
