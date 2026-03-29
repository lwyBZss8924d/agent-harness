package servicestatus

import (
	"context"
	"net"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/config"
)

type Listener struct {
	Command string `json:"command"`
	PID     int    `json:"pid,omitempty"`
	Address string `json:"address,omitempty"`
	Port    int    `json:"port"`
}

type ProbeResult struct {
	Name      string `json:"name"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Required  bool   `json:"required"`
	Reachable bool   `json:"reachable"`
}

type Status struct {
	GeneratedAtUTC string        `json:"generated_at_utc"`
	Listeners      []Listener    `json:"listeners"`
	Probes         []ProbeResult `json:"probes"`
}

func Collect(cfg config.Config) Status {
	listeners := listeningTCP()
	probes := make([]ProbeResult, 0, len(cfg.Service.PortProbes))
	for _, probe := range cfg.Service.PortProbes {
		host := probe.Host
		if host == "" {
			host = "127.0.0.1"
		}
		probes = append(probes, ProbeResult{
			Name:      probe.Name,
			Host:      host,
			Port:      probe.Port,
			Required:  probe.Required,
			Reachable: portListening(host, probe.Port),
		})
	}
	return Status{
		GeneratedAtUTC: time.Now().UTC().Format(time.RFC3339),
		Listeners:      listeners,
		Probes:         probes,
	}
}

func listeningTCP() []Listener {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	output, err := exec.CommandContext(ctx, "lsof", "-nP", "-iTCP", "-sTCP:LISTEN").CombinedOutput()
	if err != nil {
		return []Listener{}
	}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) <= 1 {
		return []Listener{}
	}
	result := make([]Listener, 0, len(lines)-1)
	seen := map[string]struct{}{}
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) < 9 {
			continue
		}
		command := fields[0]
		pid, _ := strconv.Atoi(fields[1])
		nameField := fields[len(fields)-2]
		address, port := parseAddressPort(nameField)
		if port == 0 {
			continue
		}
		key := command + "|" + strconv.Itoa(pid) + "|" + strconv.Itoa(port)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, Listener{
			Command: command,
			PID:     pid,
			Address: address,
			Port:    port,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Port == result[j].Port {
			return result[i].Command < result[j].Command
		}
		return result[i].Port < result[j].Port
	})
	return result
}

func parseAddressPort(value string) (string, int) {
	idx := strings.LastIndex(value, ":")
	if idx < 0 || idx+1 >= len(value) {
		return "", 0
	}
	port, err := strconv.Atoi(value[idx+1:])
	if err != nil {
		return "", 0
	}
	address := strings.Trim(value[:idx], "[]")
	return address, port
}

func portListening(host string, port int) bool {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, strconv.Itoa(port)), 500*time.Millisecond)
	if err == nil {
		_ = conn.Close()
		return true
	}
	return false
}
