package beacon

import (
	"net"
	"time"
)

func sendUDP(command string, ip string, port int, timeout time.Duration) ([]byte, error) {
	// Connect.
	conn, err := net.DialUDP("udp4", nil, &net.UDPAddr{IP: net.ParseIP(ip), Port: port})
	if err != nil {
		return []byte{}, err
	}
	defer conn.Close()

	// Send the command.
	conn.SetWriteDeadline(time.Now().Add(timeout))
	if _, err = conn.Write([]byte(command)); err != nil {
		return nil, err
	}

	// Read the response.
	buf := make([]byte, beaconBufferSize)
	conn.SetReadDeadline(time.Now().Add(timeout))
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}

	return buf[:n], nil
}
