package main

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"math/rand/v2"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/seancfoley/ipaddress-go/ipaddr"
	"github.com/praserx/ipconv"
)

func decompressFile(filePath string, compression string) error {
	reader, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer reader.Close()

	uncompressedStream, err := gzip.NewReader(reader)
	if err != nil {
		panic(err)
	}

	decompressedFilePath := strings.Replace(filePath, compression, "", -1)
	decompressedFile, err := os.Create(decompressedFilePath)
	if err != nil {
		panic(err)
	}
	defer decompressedFile.Close()

	fmt.Println("decompressing: " + filePath)
	buf := make([]byte, 1024)
	for {
		_, readErr := uncompressedStream.Read(buf)
		if readErr != nil && !errors.Is(readErr, io.EOF) {
			panic(readErr)
		}

		_, writeErr := decompressedFile.Write(buf)
		if writeErr != nil {
			panic(writeErr);
		}

		if errors.Is(readErr, io.EOF) {
			break
		}
	}

	return err
}

func durationUntil(timeUntil string) time.Duration {
	updateTime	:= strings.Split(timeUntil, ":")
	hours, _	:= strconv.Atoi(updateTime[0])
	minutes, _	:= strconv.Atoi(updateTime[1])

	nsNow		:= time.Now().Nanosecond()
	secondsNow	:= time.Now().Second()
	minutesNow	:= time.Now().Minute()
	hoursNow	:= time.Now().Hour()

	adjNanoseconds := 0
	if nsNow > 0 {
		adjNanoseconds = 1000000000 - nsNow
		secondsNow++
	}

	adjSeconds := 0
	if secondsNow > 0 {
		adjSeconds = 60 - secondsNow
		minutesNow++
	}

	adjMinutes := 60 - (minutesNow - minutes)
	if minutesNow > minutes {
		hoursNow++
	}

	adjHours := 0
	if hoursNow >= hours {
		adjHours = hoursNow - hours
	} else {
		adjHours = 24 - (hoursNow - hours)
		//daysNow++
	}

	duration, err := time.ParseDuration(fmt.Sprintf("%dh%dm%ds%dns", adjHours, adjMinutes, adjSeconds, adjNanoseconds))
	if err != nil {
		panic(err);
	}

	return duration
}

func fetchIP(ipString string) (*Ip, error) {
	start := time.Now()

	ip := net.ParseIP(ipString)
	if ip == nil || ip.IsPrivate() || ip.IsLoopback() {
		return nil, errors.New("invalid IP address passed (" + ipString + "); private / loopback IP ranges are not processed")
	}

	IpResult := dbIp(ip)
	IpResult.Milliseconds = time.Now().Sub(start).Milliseconds()
	IpResult.Microseconds = time.Now().Sub(start).Microseconds()

	return IpResult, nil
}

func fetchIPJson(ipString string) ([]byte, error) {
	ipResult, err := fetchIP(ipString)
	if err != nil {
		return nil, err
	}

	jsonResult, err := json.Marshal(ipResult)
	if err != nil {
		return nil, errors.New("system error")
	}

	return jsonResult, nil
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		return false
	}

	return true
}

func fileReadSmall(filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	return string(content)
}

func fileWriteSmall(filePath string, content string) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		panic(err)
	}
}

func findIPRanges(ipRangeStart string, ipRangeEnd string) []*net.IPNet {
	ipStart	:= ipaddr.NewIPAddressString(ipRangeStart)
	ipEnd	:= ipaddr.NewIPAddressString(ipRangeEnd)

	addressStart	:= ipStart.GetAddress()
	addressEnd		:= ipEnd.GetAddress()

	ipRange := addressStart.SpanWithRange(addressEnd)
	rangeSlice := ipRange.SpanWithPrefixBlocks()

	var ipNets []*net.IPNet
	for _, val := range rangeSlice {
		_, network, err := net.ParseCIDR(val.String())
		if err != nil {
			panic(err)
		}

		ipNets = append(ipNets, network)
	}

	return ipNets
}

func getEtag(url string) string {
	resp, err := http.Head(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	etag, ok := resp.Header["Etag"]
	if ok {
		return etag[0];
	}

	return ""
}

func getIpVersion(ipString string) int {
	ipVersion := 4
	if strings.Contains(ipString, ":") {
		ipVersion = 6
	}

	return ipVersion
}

func getLogFrequency() int {
	loadLogFrequency := os.Getenv("LOAD_LOG_FREQ")
	if len(loadLogFrequency) > 0 {
		loadLogFrequencyInt, err := strconv.Atoi(loadLogFrequency)
		if err != nil {
			panic(err)
		}

		return loadLogFrequencyInt
	}

	return 1000
}

func hasASNDatabase() bool {
	return len(os.Getenv("ASN")) > 0
}

func hasCityDatabase() bool {
	return len(os.Getenv("CITY")) > 0
}

func hasCountryDatabase() bool {
	return len(os.Getenv("COUNTRY")) > 0
}

func ipv4ToNumber(ipString string) int64 {
	ip, ipVersion, err := ipconv.ParseIP(ipString);
	if err == nil && ipVersion == 4 {
		number, err := ipconv.IPv4ToInt(ip)
		if err == nil {
			return int64(number)
		}
		panic(err)
	}

	return 0
}

func ipv6ToNumber(ipString string) string {
	ip, ipVersion, err := ipconv.ParseIP(ipString);
	if err == nil && ipVersion == 16 {
		arr, err := ipconv.IPv6ToInt(ip)
		if err == nil {
			bigInt := big.NewInt(0)
			for _, val := range arr {
				newNigInt := new(big.Int).SetUint64(val)
				bigInt.Add(bigInt, newNigInt)
			}
			return bigInt.String()
		}
		panic(err)
	}

	return ""
}

func isIpv4Reserved(ip string) bool {
	return strings.HasPrefix(ip, "0.") || strings.HasPrefix(ip, "127.") || strings.HasPrefix(ip, "10.") ||
	strings.HasPrefix(ip, "100.64.") || strings.HasPrefix(ip, "169.254.") || strings.HasPrefix(ip, "172.16.") ||
	strings.HasPrefix(ip, "172.17.") || strings.HasPrefix(ip, "172.18.") || strings.HasPrefix(ip, "172.19.") ||
	strings.HasPrefix(ip, "172.20.") || strings.HasPrefix(ip, "172.21.") || strings.HasPrefix(ip, "172.22.") ||
	strings.HasPrefix(ip, "172.23.") || strings.HasPrefix(ip, "172.24.") || strings.HasPrefix(ip, "172.25.") ||
	strings.HasPrefix(ip, "172.26.") || strings.HasPrefix(ip, "172.27.") || strings.HasPrefix(ip, "172.28.") ||
	strings.HasPrefix(ip, "172.29.") || strings.HasPrefix(ip, "172.30.") || strings.HasPrefix(ip, "172.31.") ||
	strings.HasPrefix(ip, "192.0.0.") || strings.HasPrefix(ip, "192.0.2.") || strings.HasPrefix(ip, "192.88.") ||
	strings.HasPrefix(ip, "192.168.") || strings.HasPrefix(ip, "198.18.") || strings.HasPrefix(ip, "198.19.") ||
	strings.HasPrefix(ip, "198.51.100") || strings.HasPrefix(ip, "203.0.113") || strings.HasPrefix(ip, "224.") ||
	strings.HasPrefix(ip, "225.") || strings.HasPrefix(ip, "226.") || strings.HasPrefix(ip, "227.") ||
	strings.HasPrefix(ip, "228.") || strings.HasPrefix(ip, "229.") || strings.HasPrefix(ip, "230.") ||
	strings.HasPrefix(ip, "231.") || strings.HasPrefix(ip, "232.") || strings.HasPrefix(ip, "233.") ||
	strings.HasPrefix(ip, "234.") || strings.HasPrefix(ip, "235.") || strings.HasPrefix(ip, "236.") ||
	strings.HasPrefix(ip, "237.") || strings.HasPrefix(ip, "238.") || strings.HasPrefix(ip, "239.") ||
	strings.HasPrefix(ip, "240.") || strings.HasPrefix(ip, "241.") || strings.HasPrefix(ip, "242.") ||
	strings.HasPrefix(ip, "243.") || strings.HasPrefix(ip, "244.") || strings.HasPrefix(ip, "245.") ||
	strings.HasPrefix(ip, "246.") || strings.HasPrefix(ip, "247.") || strings.HasPrefix(ip, "248.") ||
	strings.HasPrefix(ip, "249.") || strings.HasPrefix(ip, "250.") || strings.HasPrefix(ip, "251.") ||
	strings.HasPrefix(ip, "252.") || strings.HasPrefix(ip, "253.") || strings.HasPrefix(ip, "254.") ||
	strings.HasPrefix(ip, "255.")
}

func randomIpv4() string {
	numbers := []int{ randomNumber(0, 255), randomNumber(0, 255), randomNumber(0, 255), randomNumber(0, 255) }
	var parts []string
	for _, number := range numbers {
		parts = append(parts, strconv.Itoa(number))
	}

	ip := strings.Join(parts, ".")

	if isIpv4Reserved(ip) {
		return randomIpv4()
	}

	return ip
}

func randomIpv6() string {
	var parts []string

	first := []string{ "2001", "2002", "2003", "2400", "2401", "2402", "2403", "2404", "2405", "2406", "2407", "2408",
	"2409", "240a", "2600", "2601", "2602", "2603", "2604", "2605", "2606", "2607", "2608", "2609", "2610", "2620",
	"2800", "2801", "2802", "2803", "2804", "2806", "2a00", "2a01", "2a2", "2a03", "2a04", "2a05", "2a06", "2a07",
	"2a08", "2a09", "2a0a", "2a0b", "2a0c", "2a0d", "2a0e", "2a0f", "2a10", "2a11", "2a12", "2a13", "2a14", "2c0e",
	"2c0f" }
	pick := randomNumber(0, len(first))

	parts = append(parts, first[pick])

	// Not strictly accurate, but good enough
	for i := 0; i < 7; i++ {
		parts = append(parts, fmt.Sprintf("%02x", randomNumber(0, 255)) + fmt.Sprintf("%02x", randomNumber(0, 255)))
	}

	return strings.Join(parts, ":")
}

func randomNumber(min, max int) int {
	return rand.IntN(max-min) + min
}

func validApiKey(request *http.Request, enforceKey bool) bool {
	if len(os.Getenv("API_KEY")) > 0 || enforceKey {
		if len(os.Getenv("API_KEY")) == 0 || request.Header.Get("API-KEY") != os.Getenv("API_KEY") {
			return false
		}
	}

	return true
}