package util

import (
	"encoding/binary"
	"hash/fnv"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

func GenerateSpanID(addr string) string {
	ip := extractIPFromAddr(addr)
	ipLong, err := Ip2Long(ip)
	if err != nil {
		// 如果IP解析失败，使用默认值或者基于地址字符串的哈希
		ipLong = hashStringToUint32(addr)
	}
	times := uint64(time.Now().UnixNano())
	rand.Seed(time.Now().UnixNano())
	spanId := ((times ^ uint64(ipLong)) << 32) | uint64(rand.Int31())
	return strconv.FormatUint(spanId, 16)
}

// extractIPFromAddr 从地址字符串中提取IP部分
// 支持格式: "192.168.1.1:8080", "[::1]:8080", "::1", "192.168.1.1"
func extractIPFromAddr(addr string) string {
	// 处理IPv6格式 [::1]:port
	if strings.HasPrefix(addr, "[") {
		if idx := strings.Index(addr, "]"); idx != -1 {
			return addr[1:idx]
		}
	}

	// 处理IPv4格式 192.168.1.1:port 或者纯IP
	if idx := strings.LastIndex(addr, ":"); idx != -1 {
		// 检查是否是IPv6地址（包含多个冒号）
		if strings.Count(addr, ":") > 1 {
			// 这是一个IPv6地址，直接返回
			return addr
		}
		// IPv4地址，返回IP部分
		return addr[:idx]
	}

	// 没有端口，直接返回
	return addr
}

// hashStringToUint32 将字符串哈希为uint32，用作IP解析失败时的回退方案
func hashStringToUint32(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

// normalizeIP 标准化IP地址，将IPv6回环地址转换为IPv4
func normalizeIP(ip string) string {
	// 解析IP地址
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return ip
	}

	// 如果是IPv6回环地址，转换为IPv4
	if parsedIP.IsLoopback() {
		if parsedIP.To4() == nil {
			// IPv6回环地址，转换为IPv4回环地址
			return "127.0.0.1"
		}
	}

	// 如果是IPv4映射的IPv6地址，转换为IPv4
	if ipv4 := parsedIP.To4(); ipv4 != nil {
		return ipv4.String()
	}

	return parsedIP.String()
}

func Ip2Long(ip string) (uint32, error) {
	// 首先标准化IP地址
	normalizedIP := normalizeIP(ip)

	// 解析IP地址
	parsedIP := net.ParseIP(normalizedIP)
	if parsedIP == nil {
		return 0, net.InvalidAddrError("invalid IP address: " + ip)
	}

	// 转换为IPv4
	ipv4 := parsedIP.To4()
	if ipv4 == nil {
		// 如果不能转换为IPv4，使用哈希值
		return hashStringToUint32(ip), nil
	}

	return binary.BigEndian.Uint32(ipv4), nil
}
