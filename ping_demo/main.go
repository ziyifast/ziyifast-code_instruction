package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"time"
)

var (
	timeout int64 //ping请求超时时间
	num     int   //发送请求包的个数
	size    int64 //每个包的大小
	stop    bool  //是否一直ping
	icmp    ICMP
)

// ICMP ICMP包头
type ICMP struct {
	Type        uint8
	Code        uint8
	Checksum    uint16
	Identifier  uint16
	SequenceNum uint16
}

func main() {
	ParseArgs()
	args := os.Args
	if len(args) < 2 {
		Usage()
	}
	desIp := args[len(args)-1]

	conn, err := net.DialTimeout("ip:icmp", desIp, time.Duration(timeout)*time.Millisecond)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	//icmp头部填充
	icmp.Type = 8 //表示为icmp请求 ping请求
	icmp.Code = 0
	icmp.Checksum = 0
	icmp.Identifier = 1
	icmp.SequenceNum = 1

	fmt.Printf("\n正在 ping %s 具有 %d 字节的数据:\n", desIp, size)

	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, icmp) // 以大端模式写入（低位对应高地址）
	data := make([]byte, size)
	//写入icmp包头及空数据
	buffer.Write(data)
	data = buffer.Bytes()

	var SuccessTimes int // 成功次数
	var FailTimes int    // 失败次数
	var minTime = math.MaxInt32
	var maxTime int
	var totalTime int
	for i := 0; i < num; i++ {
		icmp.SequenceNum = uint16(1)
		// 检验和设为0
		data[2] = byte(0)
		data[3] = byte(0)

		data[6] = byte(icmp.SequenceNum >> 8)
		data[7] = byte(icmp.SequenceNum)
		//设置checksum
		icmp.Checksum = CheckSum(data)
		data[2] = byte(icmp.Checksum >> 8)
		data[3] = byte(icmp.Checksum)

		// 开始时间
		t1 := time.Now()
		conn.SetDeadline(t1.Add(time.Duration(timeout) * time.Millisecond))
		//设置icmp包checksum 校验和
		n, err := conn.Write(data)
		if err != nil {
			log.Fatal(err)
		}
		buf := make([]byte, 65535)
		n, err = conn.Read(buf)
		if err != nil {
			fmt.Println("请求超时。")
			FailTimes++
			continue
		}
		//time.Now()转换为毫秒
		et := int(time.Since(t1) / 1000000)
		if minTime > et {
			minTime = et
		}
		if maxTime < et {
			maxTime = et
		}
		totalTime += et
		fmt.Printf("来自 %s 的回复: 字节=%d 时间=%dms TTL=%d\n", desIp, len(buf[28:n]), et, buf[8])
		SuccessTimes++
		time.Sleep(1 * time.Second)
	}
	fmt.Printf("\n%s 的 Ping 统计信息:\n", desIp)
	fmt.Printf("    数据包: 已发送 = %d，已接收 = %d，丢失 = %d (%.2f%% 丢失)，\n", SuccessTimes+FailTimes, SuccessTimes, FailTimes, float64(FailTimes*100)/float64(SuccessTimes+FailTimes))
	if maxTime != 0 && minTime != math.MaxInt32 {
		fmt.Printf("往返行程的估计时间(以毫秒为单位):\n")
		fmt.Printf("    最短 = %dms，最长 = %dms，平均 = %dms\n", minTime, maxTime, totalTime/SuccessTimes)
	}
}

func ParseArgs() {
	flag.Int64Var(&timeout, "w", 10000, "超时时间(毫秒)")
	flag.IntVar(&num, "n", 4, "发送的回显请求数")
	flag.Int64Var(&size, "l", 32, "发送缓冲区大小")
	flag.BoolVar(&stop, "t", false, "Ping 指定的主机，直到停止")
	flag.Parse()
}

func Usage() {
	argNum := len(os.Args)
	if argNum < 2 {
		fmt.Print(
			`
用法: ping [-t] [-a] [-n count] [-l size] [-f] [-i TTL] [-v TOS]
            [-r count] [-s count] [[-j host-list] | [-k host-list]]
            [-w timeout] [-R] [-S srcaddr] [-c compartment] [-p]
            [-4] [-6] target_name
选项:
    -t             Ping 指定的主机，直到停止。
                   若要查看统计信息并继续操作，请键入 Ctrl+Break；
                   若要停止，请键入 Ctrl+C。
    -a             将地址解析为主机名。
    -n count       要发送的回显请求数。
    -l size        发送缓冲区大小。
    -f             在数据包中设置“不分段”标记(仅适用于 IPv4)。
    -i TTL         生存时间。
    -v TOS         服务类型(仅适用于 IPv4。该设置已被弃用，
                   对 IP 标头中的服务类型字段没有任何
                   影响)。
    -r count       记录计数跃点的路由(仅适用于 IPv4)。
    -s count       计数跃点的时间戳(仅适用于 IPv4)。
    -j host-list   与主机列表一起使用的松散源路由(仅适用于 IPv4)。
    -k host-list    与主机列表一起使用的严格源路由(仅适用于 IPv4)。
    -w timeout     等待每次回复的超时时间(毫秒)。
    -R             同样使用路由标头测试反向路由(仅适用于 IPv6)。
                   根据 RFC 5095，已弃用此路由标头。
                   如果使用此标头，某些系统可能丢弃
                   回显请求。
    -S srcaddr     要使用的源地址。
    -c compartment 路由隔离舱标识符。
    -p             Ping Hyper-V 网络虚拟化提供程序地址。
    -4             强制使用 IPv4。
    -6             强制使用 IPv6。
`)
	}
}

// CheckSum 计算校验和
func CheckSum(data []byte) uint16 {
	var sum uint32
	var length = len(data)
	var index int

	for length > 1 { // 溢出部分直接去除
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length == 1 {
		sum += uint32(data[index])
	}
	sum = uint32(uint16(sum>>16) + uint16(sum))
	sum = uint32(uint16(sum>>16) + uint16(sum))
	return uint16(^sum)
}
