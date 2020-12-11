package net2

import (
	"gosml"
	"math"
	"strings"
)

type Ipv4Range struct {
	Start int64
	End   int64
}
type Ipv4Ranges struct {
	ips []*Ipv4Range
}

func (this *Ipv4Range) Contains(ipstr string) bool {
	var ipLong int64 = IpToInt64(ipstr)
	return this.ContainsIpLong(ipLong)
}
func (this *Ipv4Range) ContainsIpLong(ipLong int64) bool {
	return ipLong >= this.Start && ipLong <= this.End
}
func (this *Ipv4Ranges) Contains(ipstr string) bool {
	for _, v := range this.ips {
		if v.Contains(ipstr) {
			return true
		}
	}
	return false
}
func NewIpv4Ranges(ipstr string) *Ipv4Ranges {
	ipv4s := make([]*Ipv4Range, 0)
	for _, s := range strings.Split(ipstr, ",") {
		ipv4s = append(ipv4s, NewIpv4Range(s))
	}
	return &Ipv4Ranges{ipv4s}
}
func NewIpv4Range(ipstr string) *Ipv4Range {
	var ss, se int64
	if strings.Contains(ipstr, "-") {
		ips := strings.Split(ipstr, "-")
		ss = gosml.ConvertToInt(ips[0])
		se = gosml.ConvertToInt(ips[1])
	} else if strings.Contains(ipstr, "*") {
		ips := strings.SplitN(ipstr, ".", 4)
		ip0 := ipMM(ips[0])
		ip1 := ipMM(ips[1])
		ip2 := ipMM(ips[2])
		ip3 := ipMM(ips[3])
		ss = IpToInt64(gosml.ConvertToString(ip0[0]) + "." + gosml.ConvertToString(ip1[0]) + "." + gosml.ConvertToString(ip2[0]) + "." + gosml.ConvertToString(ip3[0]))
		se = IpToInt64(gosml.ConvertToString(ip0[1]) + "." + gosml.ConvertToString(ip1[1]) + "." + gosml.ConvertToString(ip2[1]) + "." + gosml.ConvertToString(ip3[1]))
	} else if strings.Contains(ipstr, "/") {
		ips := strings.SplitN(ipstr, "/", 2)
		ss = IpToInt64(ips[0])
		se = ss + gosml.ConvertToInt(math.Pow(2, gosml.ConvertToFloat(32-gosml.ConvertToInt(ips[1])))) - 1
	} else {
		ss = IpToInt64(ipstr)
		se = ss
	}
	return &Ipv4Range{ss, se}
}
func ipMM(v string) []int64 {
	var rs []int64 = make([]int64, 2)
	if !strings.Contains(v, "*") {
		rs[0] = gosml.ConvertToInt(v)
		rs[1] = rs[0]
	} else if strings.EqualFold(v, "*") {
		rs[0] = 0
		rs[1] = 255
	} else if strings.HasPrefix(v, "*") && strings.HasSuffix(v, "*") {
		rs[0] = gosml.ConvertToInt(gosml.SubStr(v, 1, 2))
		if rs[0] < 5 {
			rs[1] = gosml.ConvertToInt("2" + gosml.ConvertToString(rs[0]) + "9")
		} else if rs[0] == 5 {
			rs[1] = 255
		} else {
			rs[1] = gosml.ConvertToInt("1" + gosml.ConvertToString(rs[0]) + "9")
		}
	} else if strings.HasPrefix(v, "*") {
		switch len(v) {
		case 2:
			{
				rs[0] = gosml.ConvertToInt(gosml.SubStr(v, 1, 2))
				pres := "25"
				if rs[0] > 5 {
					pres = "24"
				}
				rs[1] = gosml.ConvertToInt(pres + gosml.SubStr(v, 1, 2))
			}
		case 3:
			{
				rs[0] = gosml.ConvertToInt(gosml.SubStr(v, 1, 3))
				pres := "2"
				if rs[0] > 55 {
					pres = "1"
				}
				rs[1] = gosml.ConvertToInt(pres + gosml.SubStr(v, 1, 3))
			}
		}
	} else if strings.HasSuffix(v, "*") {
		switch len(v) {
		case 2:
			{
				rs[0] = gosml.ConvertToInt(gosml.SubStr(v, 0, 1))
				sufs := "99"
				if rs[0] == 2 {
					sufs = "55"
				} else if rs[0] > 2 {
					sufs = "9"
				}
				rs[1] = gosml.ConvertToInt(gosml.SubStr(v, 0, 1) + sufs)
			}
		case 3:
			{
				rs[0] = gosml.ConvertToInt(gosml.SubStr(v, 0, 2))
				sufs := "9"
				if rs[0] > 25 {
					sufs = ""
				}
				rs[1] = gosml.ConvertToInt(gosml.SubStr(v, 0, 2) + sufs)
			}
		}
	} else if strings.Contains(v, "*") {
		rs[0] = gosml.ConvertToInt(strings.ReplaceAll(v, "*", "0"))
		rs[1] = gosml.ConvertToInt(strings.ReplaceAll(v, "*", "9"))
	}
	return rs
}

func IpToInt64(ip string) int64 {
	ips := strings.SplitN(ip, ".", 4)
	return gosml.ConvertToInt(ips[0])<<24 | gosml.ConvertToInt(ips[1])<<16 | gosml.ConvertToInt(ips[2])<<8 | gosml.ConvertToInt(ips[3])<<0
}
func Int64ToIp(l int64) string {
	return gosml.ConvertToString(l>>24&0xff) + "." + gosml.ConvertToString(l>>16&0xff) + "." + gosml.ConvertToString(l>>8&0xff) + "." + gosml.ConvertToString(l&0xff)
}
