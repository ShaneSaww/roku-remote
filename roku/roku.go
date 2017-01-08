package roku

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

const (
	Port        = 1900
	BroadcastIP = "239.255.255.250"
)

var (
	rokuUSN    = regexp.MustCompile(`uuid:roku:ecp:([\w\d]{12})`)
	searchTime = (5 * time.Second)
)

//Get the roku ip.
func GetRokuIp(searchDuration time.Duration) (string, error) {
	//First set up the listener, so when the request is made we get the response back.
	serverAddr, _ := net.ResolveUDPAddr("udp", "0.0.0.0:"+strconv.Itoa(Port))
	conn, err := net.ListenUDP("udp", serverAddr)
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	//Build the sender
	replaceMePlaceHolder := "/replacemewithstar"
	broadcastAddr, _ := net.ResolveUDPAddr("udp", BroadcastIP+":"+strconv.Itoa(Port))
	request, _ := http.NewRequest("M-SEARCH", "http://"+broadcastAddr.String()+replaceMePlaceHolder, nil)
	setSSDPHeaders(*request)

	searchBytes := make([]byte, 0, 1024)
	buffer := bytes.NewBuffer(searchBytes)
	err = request.Write(buffer)
	if err != nil {
		return "", fmt.Errorf("Fatal error writing to buffer. This should never happen (in theory).")
	}

	searchBytes = buffer.Bytes()

	searchBytes = bytes.Replace(searchBytes, []byte(replaceMePlaceHolder), []byte("*"), 1)

	_, err = conn.WriteTo(searchBytes, broadcastAddr)
	if err != nil {
		return "", err
	}
	conn.SetReadDeadline(time.Now().Add(searchDuration))

	buf := make([]byte, 1024)
	for {
		rlen, _, err := conn.ReadFromUDP(buf)
		if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
			break // duration reached, return what we've found
		}
		if err != nil {
			return "", err
		}
		reader := bufio.NewReader(bytes.NewReader(buf[:rlen]))
		request := &http.Request{} // Needed for ReadResponse but doesn't have to be real
		response, err := http.ReadResponse(reader, request)
		headers := response.Header
		if rokuUSN.MatchString(headers.Get("usn")) {
			fmt.Println("Got Roku IP.")
			return headers.Get("Location"), nil
		}
		return "", nil
	}
	return "", nil
}

func setupListener() {

}
func setSSDPHeaders(req http.Request) {
	headers := req.Header
	headers.Set("User-Agent", "")
	headers.Set("st", "roku:ecp")
	headers.Set("man", `"ssdp:discover"`)
	headers.Set("mx", "5")
}
