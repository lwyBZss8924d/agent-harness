package devopsstatus

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/config"
)

var (
	lookPath = exec.LookPath
	runCmd   = commandOutput
)

type ToolStatus struct {
	Name      string `json:"name"`
	Required  bool   `json:"required"`
	Binary    string `json:"binary,omitempty"`
	Available bool   `json:"available"`
	Version   string `json:"version,omitempty"`
}

type DockerStatus struct {
	Available                bool   `json:"available"`
	Binary                   string `json:"binary,omitempty"`
	Version                  string `json:"version,omitempty"`
	ComposeVersion           string `json:"compose_version,omitempty"`
	Context                  string `json:"context,omitempty"`
	ContextHost              string `json:"context_host,omitempty"`
	ContextHostKind          string `json:"context_host_kind,omitempty"`
	ContextSocketPath        string `json:"context_socket_path,omitempty"`
	ContextSocketExists      bool   `json:"context_socket_exists"`
	ContextHostInForeignHome bool   `json:"context_host_in_foreign_home"`
	ServerVersion            string `json:"server_version,omitempty"`
	ServerReachable          bool   `json:"server_reachable"`
}

type KubernetesStatus struct {
	Available        bool   `json:"available"`
	Binary           string `json:"binary,omitempty"`
	ClientVersion    string `json:"client_version,omitempty"`
	CurrentContext   string `json:"current_context,omitempty"`
	ClusterInfo      string `json:"cluster_info,omitempty"`
	ClusterReachable bool   `json:"cluster_reachable"`
}

type Status struct {
	GeneratedAtUTC string           `json:"generated_at_utc"`
	Tools          []ToolStatus     `json:"tools"`
	Docker         DockerStatus     `json:"docker"`
	Kubernetes     KubernetesStatus `json:"kubernetes"`
	Registry       RegistryStatus   `json:"registry"`
}

type DockerContainer struct {
	ID      string `json:"id,omitempty"`
	Image   string `json:"image,omitempty"`
	Command string `json:"command,omitempty"`
	Created string `json:"created,omitempty"`
	Status  string `json:"status,omitempty"`
	Names   string `json:"names,omitempty"`
	Ports   string `json:"ports,omitempty"`
}

type DockerImage struct {
	Repository string `json:"repository,omitempty"`
	Tag        string `json:"tag,omitempty"`
	ID         string `json:"id,omitempty"`
	Created    string `json:"created,omitempty"`
	Size       string `json:"size,omitempty"`
}

type DockerVolume struct {
	Name       string `json:"name"`
	Driver     string `json:"driver,omitempty"`
	Mountpoint string `json:"mountpoint,omitempty"`
}

type DockerInspectResult struct {
	Target string         `json:"target"`
	Data   map[string]any `json:"data,omitempty"`
}

type KubernetesNamespace struct {
	Name   string `json:"name"`
	Status string `json:"status,omitempty"`
}

type DockerComposeProject struct {
	Name        string `json:"name"`
	Status      string `json:"status,omitempty"`
	ConfigFiles string `json:"config_files,omitempty"`
}

type KubernetesContext struct {
	Name    string `json:"name"`
	Current bool   `json:"current"`
}

type KubernetesPod struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Phase     string `json:"phase,omitempty"`
	Node      string `json:"node,omitempty"`
}

type KubernetesDeployment struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Ready     string `json:"ready,omitempty"`
	Available string `json:"available,omitempty"`
}

type KubernetesService struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Type      string `json:"type,omitempty"`
	ClusterIP string `json:"cluster_ip,omitempty"`
	Ports     string `json:"ports,omitempty"`
}

type KubernetesNode struct {
	Name    string `json:"name"`
	Status  string `json:"status,omitempty"`
	Version string `json:"version,omitempty"`
}

type KubernetesEvent struct {
	Namespace string `json:"namespace"`
	Object    string `json:"object,omitempty"`
	Reason    string `json:"reason,omitempty"`
	Type      string `json:"type,omitempty"`
	Message   string `json:"message,omitempty"`
}

type LogResult struct {
	Target string `json:"target"`
	Tail   int    `json:"tail"`
	Logs   string `json:"logs"`
}

type RegistryStatus struct {
	ConfigPath   string              `json:"config_path"`
	ConfigExists bool                `json:"config_exists"`
	CredsStore   string              `json:"creds_store,omitempty"`
	CredHelpers  map[string]string   `json:"cred_helpers,omitempty"`
	AuthEntries  []RegistryAuthEntry `json:"auth_entries,omitempty"`
}

type RegistryAuthEntry struct {
	Registry      string `json:"registry"`
	HasAuth       bool   `json:"has_auth"`
	CredHelper    string `json:"cred_helper,omitempty"`
	CredsStore    string `json:"creds_store,omitempty"`
	IdentityToken bool   `json:"identity_token"`
}

func Collect(cfg config.Config) Status {
	return Status{
		GeneratedAtUTC: time.Now().UTC().Format(time.RFC3339),
		Tools:          CollectTools(cfg),
		Docker:         CollectDocker(),
		Kubernetes:     CollectKubernetes(),
		Registry:       CollectRegistry(cfg),
	}
}

func CollectTools(cfg config.Config) []ToolStatus {
	tools := make([]ToolStatus, 0, len(cfg.DevOps.Tools))
	for _, item := range cfg.DevOps.Tools {
		path, err := lookPath(item.Binary)
		available := err == nil
		version := ""
		if available && len(item.VersionArgs) > 0 {
			version = runCmd(item.Binary, item.VersionArgs...)
		}
		tools = append(tools, ToolStatus{
			Name:      item.Name,
			Required:  item.Required,
			Binary:    path,
			Available: available,
			Version:   version,
		})
	}
	return tools
}

func CollectDocker() DockerStatus {
	path, err := lookPath("docker")
	if err != nil {
		return DockerStatus{}
	}
	contextName := runCmd("docker", "context", "show")
	host := ""
	hostKind := ""
	socketPath := ""
	socketExists := false
	hostInForeignHome := false

	if contextName != "" {
		if output := runCmd("docker", "context", "inspect", contextName); output != "" {
			var payload []map[string]any
			if json.Unmarshal([]byte(output), &payload) == nil && len(payload) > 0 {
				endpoint, _ := payload[0]["Endpoints"].(map[string]any)
				dockerEndpoint, _ := endpoint["docker"].(map[string]any)
				if rawHost, ok := dockerEndpoint["Host"].(string); ok {
					host = rawHost
					switch {
					case strings.HasPrefix(rawHost, "unix://"):
						hostKind = "unix"
						socketPath = strings.TrimPrefix(rawHost, "unix://")
						if _, statErr := exec.Command("test", "-S", socketPath).CombinedOutput(); statErr == nil {
							socketExists = true
						}
						homePrefix := filepath.ToSlash(filepath.Dir(filepath.Dir(socketPath)))
						hostInForeignHome = strings.HasPrefix(socketPath, "/Users/") &&
							!strings.Contains(homePrefix, "/Users/"+currentUser())
					case strings.HasPrefix(rawHost, "tcp://"):
						hostKind = "tcp"
					default:
						hostKind = "other"
					}
				}
			}
		}
	}

	serverVersion := runCmd("docker", "version", "--format", "{{.Server.Version}}")
	return DockerStatus{
		Available:                true,
		Binary:                   path,
		Version:                  runCmd("docker", "--version"),
		ComposeVersion:           runCmd("docker", "compose", "version"),
		Context:                  contextName,
		ContextHost:              host,
		ContextHostKind:          hostKind,
		ContextSocketPath:        socketPath,
		ContextSocketExists:      socketExists,
		ContextHostInForeignHome: hostInForeignHome,
		ServerVersion:            serverVersion,
		ServerReachable:          serverVersion != "",
	}
}

func CollectKubernetes() KubernetesStatus {
	path, err := lookPath("kubectl")
	if err != nil {
		return KubernetesStatus{}
	}
	clusterInfo := runCmd("kubectl", "cluster-info")
	return KubernetesStatus{
		Available:        true,
		Binary:           path,
		ClientVersion:    runCmd("kubectl", "version", "--client=true"),
		CurrentContext:   runCmd("kubectl", "config", "current-context"),
		ClusterInfo:      clusterInfo,
		ClusterReachable: clusterInfo != "",
	}
}

func CollectRegistry(cfg config.Config) RegistryStatus {
	result := RegistryStatus{
		ConfigPath: cfg.DevOps.DockerConfigPath,
	}
	if cfg.DevOps.DockerConfigPath == "" {
		return result
	}
	data, err := os.ReadFile(cfg.DevOps.DockerConfigPath)
	if err != nil {
		return result
	}
	result.ConfigExists = true

	var payload map[string]any
	if err := json.Unmarshal(data, &payload); err != nil {
		return result
	}

	if credsStore, ok := payload["credsStore"].(string); ok {
		result.CredsStore = credsStore
	}
	if rawHelpers, ok := payload["credHelpers"].(map[string]any); ok {
		result.CredHelpers = map[string]string{}
		for registry, value := range rawHelpers {
			if helper, ok := value.(string); ok {
				result.CredHelpers[registry] = helper
			}
		}
	}

	if rawAuths, ok := payload["auths"].(map[string]any); ok {
		entries := make([]RegistryAuthEntry, 0, len(rawAuths))
		for registry, rawValue := range rawAuths {
			entry := RegistryAuthEntry{
				Registry:   registry,
				CredsStore: result.CredsStore,
			}
			if result.CredHelpers != nil {
				entry.CredHelper = result.CredHelpers[registry]
			}
			if authMap, ok := rawValue.(map[string]any); ok {
				if authValue, ok := authMap["auth"].(string); ok && authValue != "" {
					entry.HasAuth = true
				}
				if identityToken, ok := authMap["identitytoken"].(string); ok && identityToken != "" {
					entry.IdentityToken = true
				}
			}
			entries = append(entries, entry)
		}
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Registry < entries[j].Registry
		})
		result.AuthEntries = entries
	}

	return result
}

func ListDockerContainers() []DockerContainer {
	if _, err := lookPath("docker"); err != nil {
		return []DockerContainer{}
	}
	output := runCmd("docker", "ps", "--format", "{{json .}}")
	if output == "" {
		return []DockerContainer{}
	}
	lines := strings.Split(output, "\n")
	result := make([]DockerContainer, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var payload map[string]string
		if err := json.Unmarshal([]byte(line), &payload); err != nil {
			continue
		}
		result = append(result, DockerContainer{
			ID:      payload["ID"],
			Image:   payload["Image"],
			Command: payload["Command"],
			Created: payload["RunningFor"],
			Status:  payload["Status"],
			Names:   payload["Names"],
			Ports:   payload["Ports"],
		})
	}
	return result
}

func ListKubernetesNamespaces() []KubernetesNamespace {
	if _, err := lookPath("kubectl"); err != nil {
		return []KubernetesNamespace{}
	}
	output := runCmd("kubectl", "get", "namespaces", "-o", "json")
	if output == "" {
		return []KubernetesNamespace{}
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(output), &payload); err != nil {
		return []KubernetesNamespace{}
	}
	items, _ := payload["items"].([]any)
	result := make([]KubernetesNamespace, 0, len(items))
	for _, item := range items {
		obj, ok := item.(map[string]any)
		if !ok {
			continue
		}
		metadata, _ := obj["metadata"].(map[string]any)
		status, _ := obj["status"].(map[string]any)
		name, _ := metadata["name"].(string)
		phase, _ := status["phase"].(string)
		if name == "" {
			continue
		}
		result = append(result, KubernetesNamespace{
			Name:   name,
			Status: phase,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result
}

func ListDockerImages() []DockerImage {
	if _, err := lookPath("docker"); err != nil {
		return []DockerImage{}
	}
	output := runCmd("docker", "images", "--format", "{{json .}}")
	if output == "" {
		return []DockerImage{}
	}
	lines := strings.Split(output, "\n")
	result := make([]DockerImage, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var payload map[string]string
		if err := json.Unmarshal([]byte(line), &payload); err != nil {
			continue
		}
		result = append(result, DockerImage{
			Repository: payload["Repository"],
			Tag:        payload["Tag"],
			ID:         payload["ID"],
			Created:    payload["CreatedSince"],
			Size:       payload["Size"],
		})
	}
	return result
}

func ListDockerVolumes() []DockerVolume {
	if _, err := lookPath("docker"); err != nil {
		return []DockerVolume{}
	}
	output := runCmd("docker", "volume", "ls", "--format", "{{json .}}")
	if output == "" {
		return []DockerVolume{}
	}
	lines := strings.Split(output, "\n")
	result := make([]DockerVolume, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var payload map[string]string
		if err := json.Unmarshal([]byte(line), &payload); err != nil {
			continue
		}
		result = append(result, DockerVolume{
			Name:       payload["Name"],
			Driver:     payload["Driver"],
			Mountpoint: payload["Mountpoint"],
		})
	}
	return result
}

func DockerInspect(target string) DockerInspectResult {
	if target == "" {
		return DockerInspectResult{}
	}
	output := runCmd("docker", "inspect", target)
	if output == "" {
		return DockerInspectResult{Target: target}
	}
	var payload []map[string]any
	if err := json.Unmarshal([]byte(output), &payload); err != nil || len(payload) == 0 {
		return DockerInspectResult{Target: target}
	}
	return DockerInspectResult{
		Target: target,
		Data:   payload[0],
	}
}

func ListDockerComposeProjects() []DockerComposeProject {
	if _, err := lookPath("docker"); err != nil {
		return []DockerComposeProject{}
	}
	output := runCmd("docker", "compose", "ls", "--format", "json")
	if output == "" {
		return []DockerComposeProject{}
	}
	var payload []map[string]any
	if err := json.Unmarshal([]byte(output), &payload); err != nil {
		return []DockerComposeProject{}
	}
	result := make([]DockerComposeProject, 0, len(payload))
	for _, item := range payload {
		name, _ := item["Name"].(string)
		if name == "" {
			continue
		}
		status, _ := item["Status"].(string)
		configFiles, _ := item["ConfigFiles"].(string)
		result = append(result, DockerComposeProject{
			Name:        name,
			Status:      status,
			ConfigFiles: configFiles,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result
}

func ListKubernetesContexts() []KubernetesContext {
	if _, err := lookPath("kubectl"); err != nil {
		return []KubernetesContext{}
	}
	current := runCmd("kubectl", "config", "current-context")
	output := runCmd("kubectl", "config", "get-contexts", "-o", "name")
	if output == "" {
		return []KubernetesContext{}
	}
	lines := strings.Split(output, "\n")
	result := make([]KubernetesContext, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		result = append(result, KubernetesContext{
			Name:    line,
			Current: line == current,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Current == result[j].Current {
			return result[i].Name < result[j].Name
		}
		return result[i].Current
	})
	return result
}

func ListKubernetesPods(namespace string) []KubernetesPod {
	if _, err := lookPath("kubectl"); err != nil {
		return []KubernetesPod{}
	}
	args := []string{"get", "pods"}
	if namespace != "" {
		args = append(args, "-n", namespace)
	} else {
		args = append(args, "-A")
	}
	args = append(args, "-o", "json")
	output := runCmd("kubectl", args...)
	if output == "" {
		return []KubernetesPod{}
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(output), &payload); err != nil {
		return []KubernetesPod{}
	}
	items, _ := payload["items"].([]any)
	result := make([]KubernetesPod, 0, len(items))
	for _, item := range items {
		obj, ok := item.(map[string]any)
		if !ok {
			continue
		}
		metadata, _ := obj["metadata"].(map[string]any)
		status, _ := obj["status"].(map[string]any)
		spec, _ := obj["spec"].(map[string]any)
		name, _ := metadata["name"].(string)
		ns, _ := metadata["namespace"].(string)
		phase, _ := status["phase"].(string)
		node, _ := spec["nodeName"].(string)
		if name == "" {
			continue
		}
		result = append(result, KubernetesPod{
			Namespace: ns,
			Name:      name,
			Phase:     phase,
			Node:      node,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Namespace == result[j].Namespace {
			return result[i].Name < result[j].Name
		}
		return result[i].Namespace < result[j].Namespace
	})
	return result
}

func ListKubernetesNodes() []KubernetesNode {
	if _, err := lookPath("kubectl"); err != nil {
		return []KubernetesNode{}
	}
	output := runCmd("kubectl", "get", "nodes", "-o", "json")
	if output == "" {
		return []KubernetesNode{}
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(output), &payload); err != nil {
		return []KubernetesNode{}
	}
	items, _ := payload["items"].([]any)
	result := make([]KubernetesNode, 0, len(items))
	for _, item := range items {
		obj, ok := item.(map[string]any)
		if !ok {
			continue
		}
		metadata, _ := obj["metadata"].(map[string]any)
		status, _ := obj["status"].(map[string]any)
		nodeInfo, _ := status["nodeInfo"].(map[string]any)
		name, _ := metadata["name"].(string)
		version, _ := nodeInfo["kubeletVersion"].(string)
		phase := "Unknown"
		if conditions, ok := status["conditions"].([]any); ok {
			for _, raw := range conditions {
				cond, ok := raw.(map[string]any)
				if !ok {
					continue
				}
				t, _ := cond["type"].(string)
				s, _ := cond["status"].(string)
				if t == "Ready" {
					if s == "True" {
						phase = "Ready"
					} else {
						phase = "NotReady"
					}
				}
			}
		}
		if name == "" {
			continue
		}
		result = append(result, KubernetesNode{
			Name:    name,
			Status:  phase,
			Version: version,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result
}

func ListKubernetesEvents(namespace string) []KubernetesEvent {
	if _, err := lookPath("kubectl"); err != nil {
		return []KubernetesEvent{}
	}
	args := []string{"get", "events"}
	if namespace != "" {
		args = append(args, "-n", namespace)
	} else {
		args = append(args, "-A")
	}
	args = append(args, "-o", "json")
	output := runCmd("kubectl", args...)
	if output == "" {
		return []KubernetesEvent{}
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(output), &payload); err != nil {
		return []KubernetesEvent{}
	}
	items, _ := payload["items"].([]any)
	result := make([]KubernetesEvent, 0, len(items))
	for _, item := range items {
		obj, ok := item.(map[string]any)
		if !ok {
			continue
		}
		metadata, _ := obj["metadata"].(map[string]any)
		regarding, _ := obj["regarding"].(map[string]any)
		namespaceValue, _ := metadata["namespace"].(string)
		reason, _ := obj["reason"].(string)
		eventType, _ := obj["type"].(string)
		message, _ := obj["note"].(string)
		if message == "" {
			message, _ = obj["message"].(string)
		}
		objectName, _ := regarding["name"].(string)
		objectKind, _ := regarding["kind"].(string)
		target := objectName
		if objectKind != "" {
			target = objectKind + "/" + objectName
		}
		result = append(result, KubernetesEvent{
			Namespace: namespaceValue,
			Object:    target,
			Reason:    reason,
			Type:      eventType,
			Message:   message,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Namespace == result[j].Namespace {
			if result[i].Object == result[j].Object {
				return result[i].Reason < result[j].Reason
			}
			return result[i].Object < result[j].Object
		}
		return result[i].Namespace < result[j].Namespace
	})
	return result
}

func ListKubernetesDeployments(namespace string) []KubernetesDeployment {
	if _, err := lookPath("kubectl"); err != nil {
		return []KubernetesDeployment{}
	}
	args := []string{"get", "deployments"}
	if namespace != "" {
		args = append(args, "-n", namespace)
	} else {
		args = append(args, "-A")
	}
	args = append(args, "-o", "json")
	output := runCmd("kubectl", args...)
	if output == "" {
		return []KubernetesDeployment{}
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(output), &payload); err != nil {
		return []KubernetesDeployment{}
	}
	items, _ := payload["items"].([]any)
	result := make([]KubernetesDeployment, 0, len(items))
	for _, item := range items {
		obj, ok := item.(map[string]any)
		if !ok {
			continue
		}
		metadata, _ := obj["metadata"].(map[string]any)
		status, _ := obj["status"].(map[string]any)
		name, _ := metadata["name"].(string)
		ns, _ := metadata["namespace"].(string)
		ready := intString(intValue(status["readyReplicas"])) + "/" + intString(intValue(status["replicas"]))
		available := intString(intValue(status["availableReplicas"]))
		if name == "" {
			continue
		}
		result = append(result, KubernetesDeployment{
			Namespace: ns,
			Name:      name,
			Ready:     ready,
			Available: available,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Namespace == result[j].Namespace {
			return result[i].Name < result[j].Name
		}
		return result[i].Namespace < result[j].Namespace
	})
	return result
}

func ListKubernetesServices(namespace string) []KubernetesService {
	if _, err := lookPath("kubectl"); err != nil {
		return []KubernetesService{}
	}
	args := []string{"get", "services"}
	if namespace != "" {
		args = append(args, "-n", namespace)
	} else {
		args = append(args, "-A")
	}
	args = append(args, "-o", "json")
	output := runCmd("kubectl", args...)
	if output == "" {
		return []KubernetesService{}
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(output), &payload); err != nil {
		return []KubernetesService{}
	}
	items, _ := payload["items"].([]any)
	result := make([]KubernetesService, 0, len(items))
	for _, item := range items {
		obj, ok := item.(map[string]any)
		if !ok {
			continue
		}
		metadata, _ := obj["metadata"].(map[string]any)
		spec, _ := obj["spec"].(map[string]any)
		name, _ := metadata["name"].(string)
		ns, _ := metadata["namespace"].(string)
		serviceType, _ := spec["type"].(string)
		clusterIP, _ := spec["clusterIP"].(string)
		portsValue := ""
		if ports, ok := spec["ports"].([]any); ok {
			portStrings := make([]string, 0, len(ports))
			for _, raw := range ports {
				portMap, ok := raw.(map[string]any)
				if !ok {
					continue
				}
				portStrings = append(portStrings, intString(intValue(portMap["port"])))
			}
			portsValue = strings.Join(portStrings, ",")
		}
		if name == "" {
			continue
		}
		result = append(result, KubernetesService{
			Namespace: ns,
			Name:      name,
			Type:      serviceType,
			ClusterIP: clusterIP,
			Ports:     portsValue,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Namespace == result[j].Namespace {
			return result[i].Name < result[j].Name
		}
		return result[i].Namespace < result[j].Namespace
	})
	return result
}

func DockerLogs(container string, tail int) LogResult {
	if tail <= 0 {
		tail = 200
	}
	if container == "" {
		return LogResult{Tail: tail}
	}
	output := runCmd("docker", "logs", "--tail", intString(tail), container)
	return LogResult{
		Target: container,
		Tail:   tail,
		Logs:   output,
	}
}

func KubernetesLogs(pod string, namespace string, container string, tail int) LogResult {
	if tail <= 0 {
		tail = 200
	}
	if pod == "" {
		return LogResult{Tail: tail}
	}
	args := []string{"logs", pod}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	if container != "" {
		args = append(args, "-c", container)
	}
	args = append(args, "--tail", intString(tail))
	output := runCmd("kubectl", args...)
	target := pod
	if namespace != "" {
		target = namespace + "/" + pod
	}
	return LogResult{
		Target: target,
		Tail:   tail,
		Logs:   output,
	}
}

func commandOutput(name string, args ...string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	output, err := exec.CommandContext(ctx, name, args...).CombinedOutput()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func currentUser() string {
	if out := commandOutput("whoami"); out != "" {
		return out
	}
	return ""
}

func intString(value int) string {
	return strconv.Itoa(value)
}

func intValue(value any) int {
	switch typed := value.(type) {
	case float64:
		return int(typed)
	case int:
		return typed
	default:
		return 0
	}
}
