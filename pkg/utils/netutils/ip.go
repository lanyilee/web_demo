package netutils

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net"
	"strconv"
	"strings"
)

type IPRangeInfo struct {
	MaskLen int
	Mask    string
	IPStart string
	IPEnd   string
	HostNum int64
}

func (i *IPRangeInfo) String() string {
	s, _ := json.Marshal(i)
	return string(s)
}

func Ipv4ToString(ip int) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		byte(ip>>24),
		byte(ip>>16),
		byte(ip>>8),
		byte(ip))
}

func NtoIP(ipAddrInt int64) string {
	var bytes [4]byte
	bytes[0] = byte(ipAddrInt & 0xFF)
	bytes[1] = byte((ipAddrInt >> 8) & 0xFF)
	bytes[2] = byte((ipAddrInt >> 16) & 0xFF)
	bytes[3] = byte((ipAddrInt >> 24) & 0xFF)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0]).String()
}

func MasklenToAddr(maskLen int) string {
	maskAddrInt := int64(1<<uint(maskLen)-1) << uint(32-maskLen)
	return NtoIP(maskAddrInt)
}

func IPton(ipAddr string) (int64, error) {
	var sum int64

	bits := strings.Split(ipAddr, ".")

	b0, err := strconv.Atoi(bits[0])
	if err != nil {
		return sum, errors.New("Invalid parameters")
	}
	b1, err := strconv.Atoi(bits[1])
	if err != nil {
		return sum, errors.New("Invalid parameters")
	}
	b2, err := strconv.Atoi(bits[2])
	if err != nil {
		return sum, errors.New("Invalid parameters")
	}
	b3, err := strconv.Atoi(bits[3])
	if err != nil {
		return sum, errors.New("Invalid parameters")
	}

	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)

	return sum, nil
}

func CalcNetworkAddr(ipAddr string, maskLen int) (string, error) {
	ipAddrInt, err := IPton(ipAddr)
	if err != nil {
		return "", err
	}

	maskBits := int64(1<<uint(maskLen)-1) << uint(32-maskLen)
	netAddrInt := ipAddrInt & maskBits

	return NtoIP(netAddrInt), nil
}

func CalcBroadcastAddr(ipAddr string, maskLen int) (string, error) {
	ipAddrInt, err := IPton(ipAddr)
	if err != nil {
		return "", err
	}

	maskBits := int64(1<<uint(maskLen)-1) << uint(32-maskLen)
	hostBits := int64(1<<uint(32-maskLen) - 1)
	netBits := ipAddrInt & maskBits
	brAddrInt := netBits | hostBits

	return NtoIP(brAddrInt), nil
}

func HostIpAddrs() []string {
	ipAddrs := make([]string, 0)

	addrs, err := net.InterfaceAddrs()

	if err != nil {
		zap.S().Error(err)
		return nil
	}

	for _, address := range addrs {

		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ipAddrs = append(ipAddrs, ipnet.IP.String())
			}

		}
	}

	return ipAddrs
}

func SubnetContainIP(network, ip string) bool {
	_, ipnet, err := net.ParseCIDR(network)
	if err == nil {
		return ipnet.Contains(net.ParseIP(ip))
	}
	return false
}

func CalcIPRangeByCidr(ipWithCidrOrMask, gateway string) (*IPRangeInfo, error) {

	ipAddr, maskLen, err := getIPAddrAndMasklen(ipWithCidrOrMask)
	if err != nil {
		return nil, err
	}

	hostNum := getHostNum(maskLen)
	ipStart, ipEnd, err := getCidrIpRange(ipAddr, maskLen)
	if err != nil {
		return nil, err
	}

	if !gatewayIsValid(ipStart, maskLen, gateway) {
		return nil, errors.New("网关配置错误")
	}

	return &IPRangeInfo{
		Mask:    MasklenToAddr(maskLen),
		MaskLen: maskLen,
		IPStart: ipStart,
		IPEnd:   ipEnd,
		HostNum: hostNum}, nil
}

func CalcIPRange(start, end, netmask, gateway string) (*IPRangeInfo, error) {

	if start == "" || end == "" || netmask == "" {
		return nil, errors.New("无效的IP区间或网络掩码")
	}

	ipStartInt, err := IPton(start)
	if err != nil {
		return nil, err
	}

	ipEndInt, err := IPton(end)
	if err != nil {
		return nil, err
	}

	if ipStartInt >= ipEndInt {
		return nil, errors.New("IP池的起始地址大于结束地址")
	}

	hostNum := ipEndInt - ipStartInt + 1

	ipMask := net.IPMask(net.ParseIP(netmask).To4())
	maskLen, _ := ipMask.Size()

	if maskLen == 0 {
		return nil, errors.New("无效的网络掩码")
	}

	ipStartNet, err := CalcNetworkAddr(start, maskLen)
	if err != nil {
		return nil, errors.New("无效的IP池起始地址")
	}

	ipEndNet, err := CalcNetworkAddr(end, maskLen)
	if err != nil {
		return nil, errors.New("无效的IP池结束地址")
	}

	if ipStartNet != ipEndNet {
		return nil, errors.New("IP池的起始地址和结束地址不在同一个地址段")
	}

	if !gatewayIsValid(start, maskLen, gateway) {
		return nil, errors.New("无效的网关地址")
	}

	if ipStartNet == start {
		hostNum -= 1
		start = NtoIP(ipStartInt + 1)
	}

	broadAddr, err := CalcBroadcastAddr(start, maskLen)
	if err != nil {
		return nil, errors.New("Invalid parameters")
	}

	if broadAddr == end {
		hostNum -= 1
		end = NtoIP(ipEndInt - 1)
	}

	if hostNum == 0 {
		return nil, errors.New("Invalid parameters")
	}

	return &IPRangeInfo{
		Mask:    MasklenToAddr(maskLen),
		MaskLen: maskLen,
		IPStart: start,
		IPEnd:   end,
		HostNum: hostNum}, nil
}

func IPInsideRange(ipAddr, rangeStart, rangeEnd string) bool {
	ipAddrInt, _ := IPton(ipAddr)
	rangeStartInt, _ := IPton(rangeStart)
	rangeEndInt, _ := IPton(rangeEnd)

	if ipAddrInt >= rangeStartInt && ipAddrInt <= rangeEndInt {
		return true
	}

	return false
}

func GetIfaceIpAddrs(ifaceName string) (addrs []string, err error) {
	rawIface, netErr := net.InterfaceByName(ifaceName)
	if netErr != nil {
		err = netErr
		return
	}

	ifaceAddrs, netErr := rawIface.Addrs()
	if netErr != nil {
		err = netErr
		return
	}

	for _, addr := range ifaceAddrs {
		if ipnet, ok := addr.(*net.IPNet); ok {
			ipv4Addr := ipnet.IP.To4()
			if ipv4Addr != nil {
				addrs = append(addrs, ipv4Addr.String())
			}
		}
	}

	return
}

func gatewayIsValid(ipAddr string, maskLen int, gateway string) bool {

	_, ipNet, err := net.ParseCIDR(fmt.Sprintf("%s/%d", ipAddr, maskLen))
	if err != nil {
		return false
	}

	ipv4Addr := ipNet.IP.To4()
	if ipv4Addr != nil {
		if ipv4Addr.String() == gateway {
			return false
		}
	} else {
		return false
	}

	return ipNet.Contains(net.ParseIP(gateway))
}

func getHostNum(maskLen int) int64 {
	cidrIpNum := uint32(0)
	var i uint = uint(32 - maskLen - 1)
	for ; i >= 1; i-- {
		cidrIpNum += 1 << i
	}
	return int64(cidrIpNum)
}

func getCidrIpRange(ip string, maskLen int) (minIP, maxIP string, err error) {
	ipSegs := strings.Split(ip, ".")

	minIPSegs := make([]string, 0)
	maxIPSegs := make([]string, 0)
	for idx, ipSeg := range ipSegs {
		if (idx+1)*8 <= maskLen {
			minIPSegs = append(minIPSegs, ipSeg)
			maxIPSegs = append(maxIPSegs, ipSeg)
			continue
		}

		ipSegInt, err := strconv.ParseUint(ipSeg, 10, 8)
		if err != nil {
			return "", "", errors.New("Invalid parameters")
		}

		segMaskLen := (idx+1)*8 - maskLen

		segMin, segMax := getIpSegRange(idx+1, uint8(ipSegInt), uint8(segMaskLen))
		minIPSegs = append(minIPSegs, strconv.Itoa(segMin))
		maxIPSegs = append(maxIPSegs, strconv.Itoa(segMax))
	}

	minIP = strings.Join(minIPSegs, ".")
	maxIP = strings.Join(maxIPSegs, ".")
	return
}

func getIpSegRange(segIdx int, userSegIp, offset uint8) (int, int) {
	var ipSegMax uint8 = 255
	netSegIp := ipSegMax << offset
	segMinIp := netSegIp & userSegIp
	segMaxIp := userSegIp&(255<<offset) | ^(255 << offset)

	if segIdx == 4 {
		if segMinIp == 0 {
			segMinIp += 1
		}

		if segMaxIp == 255 {
			segMaxIp -= 1
		}
	}

	return int(segMinIp), int(segMaxIp)
}

func getIPAddrAndMasklen(ipWithCidrOrMask string) (ipAddr string, maskLen int, err error) {
	ipArray := strings.Split(ipWithCidrOrMask, "/")
	if len(ipArray) != 2 {
		err = errors.New("网络地址格式错误")
		return
	}

	ipAddr = ipArray[0]
	maskLen, err = strconv.Atoi(ipArray[1])
	if err != nil {
		netmask := net.IPMask(net.ParseIP(ipArray[1]).To4())
		maskLen, _ = netmask.Size()
		err = nil
	}

	if maskLen > 24 || maskLen <= 0 {
		err = errors.New("无效的网络掩码")
	}

	return
}

func Allocate(cidr string, useds []string) (string, error) {
	ipAddr, maskLen, err := getIPAddrAndMasklen(cidr)
	if err != nil {
		return "", err
	}
	ipStart, ipEnd, err := getCidrIpRange(ipAddr, maskLen)
	if err != nil {
		return "", err
	}
	startipInt, err := IPton(ipStart)
	if err != nil {
		return "", err
	}
	endipInt, err := IPton(ipEnd)
	if err != nil {
		return "", err
	}
	sum := endipInt - startipInt + 1
	fmt.Println(ipStart, ipEnd)
	ipPools := make([]int64, sum)
	for _, u := range useds {
		uIP, _ := IPton(u)
		ipPools[int(uIP-startipInt)] = uIP
	}
	for i := 0; i < int(sum); i++ {
		if ipPools[int(i)] == 0 {
			return NtoIP(int64(i) + startipInt), nil
		}
	}
	return "", errors.New("not found")
}
