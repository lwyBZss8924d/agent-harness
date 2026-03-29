package devopsstatus

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/config"
)

func TestCollectRegistry(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	data := `{
	  "credsStore": "osxkeychain",
	  "credHelpers": {"ghcr.io": "ghcr-helper"},
	  "auths": {
	    "ghcr.io": {"auth": "abc"},
	    "docker.io": {"identitytoken": "tok"}
	  }
	}`
	if err := os.WriteFile(configPath, []byte(data), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	cfg := config.Config{
		DevOps: config.DevOpsConfig{
			DockerConfigPath: configPath,
		},
	}
	status := CollectRegistry(cfg)
	if !status.ConfigExists {
		t.Fatal("expected config exists")
	}
	if status.CredsStore != "osxkeychain" {
		t.Fatalf("CredsStore = %q", status.CredsStore)
	}
	if len(status.AuthEntries) != 2 {
		t.Fatalf("len(AuthEntries) = %d", len(status.AuthEntries))
	}
}

func TestListDockerComposeProjects(t *testing.T) {
	oldLookPath := lookPath
	oldRunCmd := runCmd
	t.Cleanup(func() {
		lookPath = oldLookPath
		runCmd = oldRunCmd
	})
	lookPath = func(name string) (string, error) {
		if name == "docker" {
			return "/usr/local/bin/docker", nil
		}
		return "", errors.New("missing")
	}
	runCmd = func(name string, args ...string) string {
		return `[{"Name":"app","Status":"running(2)","ConfigFiles":"/tmp/docker-compose.yml"}]`
	}
	projects := ListDockerComposeProjects()
	if len(projects) != 1 || projects[0].Name != "app" {
		t.Fatalf("projects = %#v", projects)
	}
}

func TestListKubernetesContextsAndNodes(t *testing.T) {
	oldLookPath := lookPath
	oldRunCmd := runCmd
	t.Cleanup(func() {
		lookPath = oldLookPath
		runCmd = oldRunCmd
	})
	lookPath = func(name string) (string, error) {
		if name == "kubectl" {
			return "/usr/local/bin/kubectl", nil
		}
		return "", errors.New("missing")
	}
	runCmd = func(name string, args ...string) string {
		joined := name + " " + stringJoin(args, " ")
		switch joined {
		case "kubectl config current-context":
			return "orbstack"
		case "kubectl config get-contexts -o name":
			return "orbstack\nother"
		case "kubectl get nodes -o json":
			return `{"items":[{"metadata":{"name":"orbstack"},"status":{"nodeInfo":{"kubeletVersion":"v1.33.5"},"conditions":[{"type":"Ready","status":"True"}]}}]}`
		default:
			return ""
		}
	}
	contexts := ListKubernetesContexts()
	if len(contexts) != 2 || !contexts[0].Current {
		t.Fatalf("contexts = %#v", contexts)
	}
	nodes := ListKubernetesNodes()
	if len(nodes) != 1 || nodes[0].Status != "Ready" {
		t.Fatalf("nodes = %#v", nodes)
	}
}

func TestDockerLogsAndInspect(t *testing.T) {
	oldRunCmd := runCmd
	t.Cleanup(func() {
		runCmd = oldRunCmd
	})
	runCmd = func(name string, args ...string) string {
		if len(args) > 0 && args[0] == "logs" {
			return "hello"
		}
		if len(args) > 0 && args[0] == "inspect" {
			return `[{"Name":"demo"}]`
		}
		return ""
	}
	logs := DockerLogs("abc", 3)
	if logs.Logs != "hello" {
		t.Fatalf("logs = %#v", logs)
	}
	inspect := DockerInspect("abc")
	if inspect.Data["Name"] != "demo" {
		t.Fatalf("inspect = %#v", inspect)
	}
}

func TestListDockerImages(t *testing.T) {
	oldLookPath := lookPath
	oldRunCmd := runCmd
	t.Cleanup(func() {
		lookPath = oldLookPath
		runCmd = oldRunCmd
	})
	lookPath = func(name string) (string, error) {
		if name == "docker" {
			return "/usr/local/bin/docker", nil
		}
		return "", errors.New("missing")
	}
	runCmd = func(name string, args ...string) string {
		return `{"Repository":"postgres","Tag":"17","ID":"abc","CreatedSince":"1 day ago","Size":"100MB"}`
	}
	images := ListDockerImages()
	if len(images) != 1 || images[0].Repository != "postgres" {
		t.Fatalf("images = %#v", images)
	}
}

func TestListDockerVolumes(t *testing.T) {
	oldLookPath := lookPath
	oldRunCmd := runCmd
	t.Cleanup(func() {
		lookPath = oldLookPath
		runCmd = oldRunCmd
	})
	lookPath = func(name string) (string, error) {
		if name == "docker" {
			return "/usr/local/bin/docker", nil
		}
		return "", errors.New("missing")
	}
	runCmd = func(name string, args ...string) string {
		return `{"Name":"vol1","Driver":"local","Mountpoint":"/var/lib/docker/volumes/vol1/_data"}`
	}
	volumes := ListDockerVolumes()
	if len(volumes) != 1 || volumes[0].Name != "vol1" {
		t.Fatalf("volumes = %#v", volumes)
	}
}

func TestListDockerComposeProjectsAndContainers(t *testing.T) {
	oldLookPath := lookPath
	oldRunCmd := runCmd
	t.Cleanup(func() {
		lookPath = oldLookPath
		runCmd = oldRunCmd
	})
	lookPath = func(name string) (string, error) {
		if name == "docker" {
			return "/usr/local/bin/docker", nil
		}
		return "", errors.New("missing")
	}
	runCmd = func(name string, args ...string) string {
		joined := name + " " + stringJoin(args, " ")
		switch joined {
		case "docker compose ls --format json":
			return `[{"Name":"demo","Status":"running(1)","ConfigFiles":"/tmp/compose.yml"}]`
		case "docker ps --format {{json .}}":
			return `{"ID":"abc","Image":"postgres:17","Command":"postgres","RunningFor":"1 day","Status":"Up","Names":"pg","Ports":"5432/tcp"}`
		default:
			return ""
		}
	}
	projects := ListDockerComposeProjects()
	if len(projects) != 1 || projects[0].Name != "demo" {
		t.Fatalf("projects = %#v", projects)
	}
	containers := ListDockerContainers()
	if len(containers) != 1 || containers[0].Names != "pg" {
		t.Fatalf("containers = %#v", containers)
	}
}

func TestListKubernetesPodsNamespacesServices(t *testing.T) {
	oldLookPath := lookPath
	oldRunCmd := runCmd
	t.Cleanup(func() {
		lookPath = oldLookPath
		runCmd = oldRunCmd
	})
	lookPath = func(name string) (string, error) {
		if name == "kubectl" {
			return "/usr/local/bin/kubectl", nil
		}
		return "", errors.New("missing")
	}
	runCmd = func(name string, args ...string) string {
		joined := name + " " + stringJoin(args, " ")
		switch joined {
		case "kubectl get pods -A -o json":
			return `{"items":[{"metadata":{"name":"pod1","namespace":"ns1"},"status":{"phase":"Running"},"spec":{"nodeName":"node1"}}]}`
		case "kubectl get namespaces -o json":
			return `{"items":[{"metadata":{"name":"ns1"},"status":{"phase":"Active"}}]}`
		case "kubectl get services -A -o json":
			return `{"items":[{"metadata":{"name":"svc1","namespace":"ns1"},"spec":{"type":"ClusterIP","clusterIP":"10.0.0.1","ports":[{"port":80}]}}]}`
		case "kubectl get deployments -A -o json":
			return `{"items":[{"metadata":{"name":"deploy1","namespace":"ns1"},"status":{"readyReplicas":1,"replicas":2,"availableReplicas":1}}]}`
		default:
			return ""
		}
	}
	if got := ListKubernetesPods(""); len(got) != 1 || got[0].Name != "pod1" {
		t.Fatalf("pods = %#v", got)
	}
	if got := ListKubernetesNamespaces(); len(got) != 1 || got[0].Name != "ns1" {
		t.Fatalf("namespaces = %#v", got)
	}
	if got := ListKubernetesServices(""); len(got) != 1 || got[0].Name != "svc1" {
		t.Fatalf("services = %#v", got)
	}
	if got := ListKubernetesDeployments(""); len(got) != 1 || got[0].Name != "deploy1" {
		t.Fatalf("deployments = %#v", got)
	}
}

func TestListKubernetesEvents(t *testing.T) {
	oldLookPath := lookPath
	oldRunCmd := runCmd
	t.Cleanup(func() {
		lookPath = oldLookPath
		runCmd = oldRunCmd
	})
	lookPath = func(name string) (string, error) {
		if name == "kubectl" {
			return "/usr/local/bin/kubectl", nil
		}
		return "", errors.New("missing")
	}
	runCmd = func(name string, args ...string) string {
		return `{"items":[{"metadata":{"namespace":"ns1"},"regarding":{"kind":"Pod","name":"pod1"},"reason":"Started","type":"Normal","note":"ok"}]}`
	}
	events := ListKubernetesEvents("")
	if len(events) != 1 || events[0].Object != "Pod/pod1" {
		t.Fatalf("events = %#v", events)
	}
}

func TestListKubernetesNamespacesContextsPodsDeploymentsServices(t *testing.T) {
	oldLookPath := lookPath
	oldRunCmd := runCmd
	t.Cleanup(func() {
		lookPath = oldLookPath
		runCmd = oldRunCmd
	})
	lookPath = func(name string) (string, error) {
		if name == "kubectl" {
			return "/usr/local/bin/kubectl", nil
		}
		return "", errors.New("missing")
	}
	runCmd = func(name string, args ...string) string {
		joined := name + " " + stringJoin(args, " ")
		switch joined {
		case "kubectl config current-context":
			return "ctx1"
		case "kubectl config get-contexts -o name":
			return "ctx1\nctx2"
		case "kubectl get namespaces -o json":
			return `{"items":[{"metadata":{"name":"ns1"},"status":{"phase":"Active"}}]}`
		case "kubectl get pods -A -o json":
			return `{"items":[{"metadata":{"name":"pod1","namespace":"ns1"},"status":{"phase":"Running"},"spec":{"nodeName":"node1"}}]}`
		case "kubectl get deployments -A -o json":
			return `{"items":[{"metadata":{"name":"dep1","namespace":"ns1"},"status":{"readyReplicas":1,"replicas":1,"availableReplicas":1}}]}`
		case "kubectl get services -A -o json":
			return `{"items":[{"metadata":{"name":"svc1","namespace":"ns1"},"spec":{"type":"ClusterIP","clusterIP":"10.0.0.1","ports":[{"port":80}]}}]}`
		default:
			return ""
		}
	}
	if got := ListKubernetesContexts(); len(got) != 2 || !got[0].Current {
		t.Fatalf("contexts = %#v", got)
	}
	if got := ListKubernetesNamespaces(); len(got) != 1 || got[0].Name != "ns1" {
		t.Fatalf("namespaces = %#v", got)
	}
	if got := ListKubernetesPods(""); len(got) != 1 || got[0].Name != "pod1" {
		t.Fatalf("pods = %#v", got)
	}
	if got := ListKubernetesDeployments(""); len(got) != 1 || got[0].Name != "dep1" {
		t.Fatalf("deployments = %#v", got)
	}
	if got := ListKubernetesServices(""); len(got) != 1 || got[0].Name != "svc1" {
		t.Fatalf("services = %#v", got)
	}
}

func TestKubernetesLogsAndNodes(t *testing.T) {
	oldLookPath := lookPath
	oldRunCmd := runCmd
	t.Cleanup(func() {
		lookPath = oldLookPath
		runCmd = oldRunCmd
	})
	lookPath = func(name string) (string, error) {
		if name == "kubectl" {
			return "/usr/local/bin/kubectl", nil
		}
		return "", errors.New("missing")
	}
	runCmd = func(name string, args ...string) string {
		joined := name + " " + stringJoin(args, " ")
		switch joined {
		case "kubectl logs pod1 -n ns1 --tail 5":
			return "hello"
		case "kubectl get nodes -o json":
			return `{"items":[{"metadata":{"name":"n1"},"status":{"nodeInfo":{"kubeletVersion":"v1"},"conditions":[{"type":"Ready","status":"False"}]}}]}`
		default:
			return ""
		}
	}
	logs := KubernetesLogs("pod1", "ns1", "", 5)
	if logs.Logs != "hello" {
		t.Fatalf("logs = %#v", logs)
	}
	nodes := ListKubernetesNodes()
	if len(nodes) != 1 || nodes[0].Status != "NotReady" {
		t.Fatalf("nodes = %#v", nodes)
	}
}

func stringJoin(items []string, sep string) string {
	if len(items) == 0 {
		return ""
	}
	result := items[0]
	for _, item := range items[1:] {
		result += sep + item
	}
	return result
}
