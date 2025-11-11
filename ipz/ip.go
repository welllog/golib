package ipz

import (
	"errors"
	"net"
	"net/http"
	"strings"
)

// GetRemoteIp try to get the client's real IP address from the HTTP request.
func GetRemoteIp(r *http.Request) string {
	// try to get from X-Forwarded-For
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		// get the first IP and validate
		firstIpStr := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
		ip := net.ParseIP(firstIpStr)
		if ip != nil {
			return ip.String()
		}
	}

	// try to get from X-Real-Ip
	xRealIp := r.Header.Get("X-Real-Ip")
	if xRealIp != "" {
		ip := net.ParseIP(strings.TrimSpace(xRealIp))
		if ip != nil {
			return ip.String()
		}
	}

	// fall back to RemoteAddr
	if ipStr, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		ip := net.ParseIP(ipStr)
		if ip != nil {
			return ip.String()
		}
	}

	return ""
}

func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, address := range addrs {
		// Check the address type and loopback status
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", errors.New("no local ip found")
}

func GetLocalIPByName(ifaceName string) (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range interfaces {
		if iface.Name == ifaceName && iface.Flags&net.FlagUp != 0 {
			addrs, err := iface.Addrs()
			if err != nil {
				return "", err
			}
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String(), nil
					}
				}
			}
		}
	}
	return "", errors.New("no ip found for interface " + ifaceName)
}

// GetOutboundIP get the preferred outbound IP address of this machine
func GetOutboundIP() (string, error) {
	// Connect to a public IP address (Google DNS) to determine the outbound IP
	// port 80 is usually open, and we don't need to send any data
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String(), nil
}
