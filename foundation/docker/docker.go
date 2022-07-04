// package docker support starting and stopping container for running unit tests
package docker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os/exec"
	"strings"
)

// type Container tracks information about container sstarte for the test
type Container struct {
	ID string
	Host string
}

// StartContainer start specific container for running test
func StartContainer(image string, port string, args ...string) ( *Container, error){
	arg := []string{"run", "-P", "-d"}
	arg = append(arg, args...)
	arg = append(arg, image)


	cmd := exec.Command("docker", arg...)
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("could not start container %s: %w", image, err)
	}

	id :=  out.String()[:12]

	hostIP, hostPort, err := extractIPPort(id, port)
	if err!=nil {
		return nil, fmt.Errorf("could not extract ip/port: %w", err)
	}

	c := Container{
		ID: id,
		Host: net.JoinHostPort(hostIP, hostPort),
	}

	fmt.Printf("Image:       %s\n", image)
	fmt.Printf("ContainerID: %s\n", c.ID)
	fmt.Printf("Host:        %s\n", c.Host)


	return &c, nil
}

func StopContainer(id string) error {
	if err := exec.Command("docker", "stop", id).Run(); err != nil {
		return fmt.Errorf("could not stop container: %w", err)
	}
	
	fmt.Println("Stopped:", id)

	if err :=  exec.Command("docker", "rm", id, "-v").Run(); err != nil {
		return fmt.Errorf("could not remove container: %w", err)
	}

	fmt.Println("Removed:", id)


	return nil
}


func extractIPPort(id string, port string) (hostIP string, hostPort string, err error) {
	tmpl := fmt.Sprintf("[{{range $k,$v := (index .NetworkSettings.Ports \"%s/tcp\")}}{{json $v}}{{end}}]", port)

	cmd := exec.Command("docker", "inspect", "-f", tmpl, id)
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", "", fmt.Errorf("could not inspect container %s: %w", id, err)
	}

	// When IPv6 is turned on with Docker.
	// Got  [{"HostIp":"0.0.0.0","HostPort":"49190"}{"HostIp":"::","HostPort":"49190"}]
	// Need [{"HostIp":"0.0.0.0","HostPort":"49190"},{"HostIp":"::","HostPort":"49190"}]
	data := strings.ReplaceAll(out.String(), "}{", "},{")

	var docs []struct {
		HostIP string
		HostPort string
	}

	if err := json.Unmarshal([]byte(data), &docs); err != nil {
		return "", "", fmt.Errorf("could not decode json: %w", err)
	}

	for _, doc := range docs {
		if doc.HostIP != "::" {
			return doc.HostIP, doc.HostPort, nil
		}
	}

	return "", "",  fmt.Errorf("could not locate ip/port")
}

