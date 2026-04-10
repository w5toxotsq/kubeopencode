# Live Agents

Every Agent in KubeOpenCode is a persistent, running service on Kubernetes — available for interactive use anytime.

## Why Live Agents?

- **Zero cold start**: The agent is always running. No waiting for container startup when you need help.
- **Shared context**: Pre-load codebases, documentation, and organizational standards. All tasks share the same context.
- **Interactive access**: Connect via web terminal or CLI for real-time pair programming.
- **Session persistence**: Conversation history survives pod restarts (crashes, node drains, upgrades).
- **Team-shared agents**: One agent serves your entire team — consistent configuration, centralized credential management.

## Use Cases

| Use Case | Description |
|----------|-------------|
| **Team coding assistant** | A shared agent pre-loaded with your monorepo and coding standards. Team members attach via CLI to get interactive help. |
| **Slack/ChatOps bot** | An always-on agent that responds to Slack messages, creating PRs and fixing issues on demand. |
| **Code review agent** | A persistent agent that reviews PRs as they come in, leveraging shared context about your codebase. |
| **On-call assistant** | An agent with production runbooks and monitoring dashboards pre-loaded, ready to help debug incidents. |

## Setup

Creating an Agent automatically creates a persistent Deployment + Service:

```yaml
apiVersion: kubeopencode.io/v1alpha1
kind: Agent
metadata:
  name: team-agent
spec:
  profile: "Team development agent"
  executorImage: ghcr.io/kubeopencode/kubeopencode-agent-devbox:latest
  workspaceDir: /workspace
  serviceAccountName: kubeopencode-agent
  port: 4096                      # OpenCode server port (default: 4096)
  persistence:
    sessions:
      size: "2Gi"                 # Persist conversation history
  contexts:
    - name: codebase
      type: Git
      git:
        repository: https://github.com/your-org/your-repo.git
        ref: main
      mountPath: code
  credentials:
    - name: api-key
      secretRef:
        name: ai-credentials
        key: api-key
      env: OPENCODE_API_KEY
```

The controller automatically creates:
- A **Deployment** running `opencode serve` (persistent server)
- A **Service** for in-cluster access (e.g., `http://team-agent.kubeopencode-system.svc.cluster.local:4096`)

## Interacting with Live Agents

**CLI attach** (connects through kube-apiserver service proxy — no Ingress or port-forward needed):

```bash
kubeoc agent attach team-agent -n kubeopencode-system
```

**Web Terminal**: Access the agent's OpenCode TUI directly from the KubeOpenCode dashboard at `http://localhost:2746`.

**Programmatic Tasks**: Submit Tasks referencing the Agent — they run on the persistent server via `--attach` flag:

```bash
kubectl apply -f task.yaml
```

## Agent vs Template Tasks

| Aspect | `agentRef` (Agent) | `templateRef` (AgentTemplate) |
|--------|-------------------|-------------------------------|
| Lifecycle | Persistent Deployment + lightweight Pod per Task | Ephemeral Pod per Task |
| Command | `opencode run --attach <url> "task"` | `opencode run "task"` |
| Cold start | No (server already running) | Yes (container startup) |
| Context sharing | Shared across Tasks via server | Isolated per Task |
| Interaction | Web Terminal, CLI attach, API | Logs only |
| Best for | Interactive coding, team agents | Batch operations, CI/CD, one-off tasks |
| Concurrency/Quota | Enforced by Agent | Not enforced |

## Agent Status

Monitor your live agent's health:

```bash
kubectl get agent team-agent -o wide
# NAME         PROFILE                  STATUS
# team-agent   Team development agent   Ready
```

The Agent status includes deployment details:
```yaml
status:
  deploymentName: team-agent-server
  serviceName: team-agent
  url: http://team-agent.kubeopencode-system.svc.cluster.local:4096
  ready: true
  conditions:
    - type: ServerReady
      status: "True"
    - type: ServerHealthy
      status: "True"
```

See [Getting Started](../getting-started.md) for a complete walkthrough.
