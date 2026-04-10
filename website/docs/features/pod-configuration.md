# Pod Configuration

## Pod Security

KubeOpenCode applies a restricted security context by default to all agent containers, following the Kubernetes Pod Security Standards (Restricted profile).

### Default Security Context

When no `securityContext` is specified in `podSpec`, the controller applies these defaults to all containers (init containers and the worker container):

- `allowPrivilegeEscalation: false`
- `capabilities: drop: ["ALL"]`
- `seccompProfile: type: RuntimeDefault`

These defaults align with the Kubernetes Restricted Pod Security Standard and are suitable for most workloads.

### Custom Container Security Context

Override the default security context for tighter or workload-specific settings using `podSpec.securityContext`:

```yaml
apiVersion: kubeopencode.io/v1alpha1
kind: Agent
metadata:
  name: hardened-agent
spec:
  profile: "Security-hardened agent with strict container settings"
  agentImage: ghcr.io/kubeopencode/kubeopencode-agent-opencode:latest
  executorImage: ghcr.io/kubeopencode/kubeopencode-agent-devbox:latest
  workspaceDir: /workspace
  serviceAccountName: kubeopencode-agent
  podSpec:
    securityContext:
      runAsNonRoot: true
      allowPrivilegeEscalation: false
      readOnlyRootFilesystem: true
      capabilities:
        drop:
          - ALL
      seccompProfile:
        type: RuntimeDefault
```

> **Note**: When using `readOnlyRootFilesystem: true`, ensure the agent image supports it. You may need to use `emptyDir` volumes for writable paths (e.g., `/tmp`, `/home/agent`).

### Pod-Level Security Context

Use `podSpec.podSecurityContext` to configure security attributes that apply to the entire Pod (all containers):

```yaml
apiVersion: kubeopencode.io/v1alpha1
kind: Agent
metadata:
  name: uid-agent
spec:
  profile: "Agent running as specific user and group"
  agentImage: ghcr.io/kubeopencode/kubeopencode-agent-opencode:latest
  executorImage: ghcr.io/kubeopencode/kubeopencode-agent-devbox:latest
  workspaceDir: /workspace
  serviceAccountName: kubeopencode-agent
  podSpec:
    podSecurityContext:
      runAsUser: 1000
      runAsGroup: 1000
      fsGroup: 1000
```

`podSecurityContext` is useful for:
- Enforcing a specific UID/GID for all containers
- Setting `fsGroup` for shared volume permissions
- Meeting namespace-level Pod Security Admission requirements

## Advanced Pod Settings

Configure advanced Pod settings using `podSpec`:

```yaml
apiVersion: kubeopencode.io/v1alpha1
kind: Agent
metadata:
  name: advanced-agent
spec:
  profile: "Advanced agent with gVisor isolation and GPU scheduling"
  agentImage: ghcr.io/kubeopencode/kubeopencode-agent-opencode:latest
  executorImage: ghcr.io/kubeopencode/kubeopencode-agent-devbox:latest
  workspaceDir: /workspace
  serviceAccountName: kubeopencode-agent
  podSpec:
    # Labels for NetworkPolicy, monitoring, etc.
    labels:
      network-policy: agent-restricted
    # Enhanced isolation with gVisor or Kata
    runtimeClassName: gvisor
    # Scheduling configuration
    scheduling:
      nodeSelector:
        node-type: ai-workload
      tolerations:
        - key: "dedicated"
          operator: "Equal"
          value: "ai-workload"
          effect: "NoSchedule"
```
