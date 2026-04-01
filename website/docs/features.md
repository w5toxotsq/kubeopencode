---
sidebar_position: 2
title: Features
description: Live Agents for interactive work, AgentTemplates for automated workflows, and enterprise-grade platform capabilities
---

# Features

KubeOpenCode supports two primary scenarios for running AI agents on Kubernetes. Choose based on your use case:

| | **Live Agent** | **AgentTemplate Task** |
|---|---|---|
| **When to use** | Complex tasks needing human interaction | Stable, repeatable workflows at scale |
| **Lifecycle** | Persistent (Deployment + Service) | Ephemeral (Pod per Task) |
| **Interaction** | Web terminal, CLI, or programmatic Tasks | Submit-and-forget via `kubectl apply` |
| **State** | Session history persists across restarts | Stateless — each Task starts fresh |
| **Example** | Pair programming, incident debugging, code review | CI/CD pipelines, batch updates, dependency upgrades |

---

## Scenario 1: Live Agent — Human-in-the-Loop

A Live Agent is a persistent AI coding service running on Kubernetes. It is always on, shared by your team, and supports real-time interaction.

### Real-time Access

Connect to a running Agent at any time:

**Web Terminal** — Open the KubeOpenCode dashboard at `http://localhost:2746`, navigate to the Agent, and click "Terminal" for a full OpenCode TUI in your browser.

**CLI** — Attach directly through the kube-apiserver service proxy (no Ingress or port-forward needed):

```bash
kubeoc agent attach team-agent -n kubeopencode-system
```

**Programmatic Tasks** — Submit Tasks as YAML. They run on the persistent Agent's server:

```yaml
apiVersion: kubeopencode.io/v1alpha1
kind: Task
metadata:
  name: fix-bug-123
spec:
  agentRef:
    name: team-agent
  description: |
    Fix the null pointer exception in UserService.java.
```

### Resource Management

Agents consume cluster resources. KubeOpenCode provides controls to manage this:

**Idle Timeout** — Automatically suspend Agents after a period of inactivity. When a new Task arrives, the Agent resumes automatically:

```yaml
spec:
  idleTimeout: "30m"   # Auto-suspend after 30 minutes idle
```

**Manual Suspend / Resume** — Scale an Agent to zero while retaining all persistent data:

```bash
# Suspend (Deployment scales to 0, PVCs retained)
kubectl patch agent team-agent --type=merge -p '{"spec":{"suspend":true}}'

# Resume
kubectl patch agent team-agent --type=merge -p '{"spec":{"suspend":false}}'
```

Tasks submitted to a suspended Agent enter `Queued` phase and run automatically when the Agent resumes.

### Quick Setup via AgentTemplate

Use an AgentTemplate to define shared configuration (images, credentials, workspace). Agents inherit from it via `templateRef`:

```yaml
# Template — shared base configuration
apiVersion: kubeopencode.io/v1alpha1
kind: AgentTemplate
metadata:
  name: org-base
spec:
  agentImage: quay.io/kubeopencode/kubeopencode-agent-opencode:latest
  executorImage: quay.io/kubeopencode/kubeopencode-agent-devbox:latest
  workspaceDir: /workspace
  serviceAccountName: kubeopencode-agent
  credentials:
    - name: api-key
      secretRef:
        name: ai-credentials
        key: api-key
      env: OPENCODE_API_KEY
---
# Agent — inherits from template, adds its own config
apiVersion: kubeopencode.io/v1alpha1
kind: Agent
metadata:
  name: team-agent
spec:
  templateRef:
    name: org-base
  profile: "Team development agent"
  executorImage: quay.io/kubeopencode/kubeopencode-agent-devbox:latest
  workspaceDir: /workspace
  serviceAccountName: kubeopencode-agent
  port: 4096
  persistence:
    sessions:
      size: "2Gi"
  idleTimeout: "30m"
  maxConcurrentTasks: 3
```

### Session Persistence

Conversation history survives pod restarts (crashes, node drains, upgrades):

```yaml
spec:
  persistence:
    sessions:
      size: "2Gi"               # PVC for OpenCode SQLite database
      storageClassName: "gp3"   # Optional, uses cluster default
    workspace:
      size: "10Gi"              # PVC for workspace files (git repos, modified files)
```

When configured, the controller creates PVCs with OwnerReferences — they are garbage-collected when the Agent is deleted.

### Concurrency Control

Limit simultaneous Tasks to respect AI provider rate limits:

```yaml
spec:
  maxConcurrentTasks: 3   # Only 3 Tasks run at once; extras enter Queued phase
```

---

## Scenario 2: AgentTemplate Tasks — Automated Workflows at Scale

For CI/CD pipelines, batch operations, and one-off tasks, use an **AgentTemplate** directly. Each Task runs in an ephemeral Pod that is created and destroyed per execution.

### Create a Task with templateRef

```yaml
apiVersion: kubeopencode.io/v1alpha1
kind: AgentTemplate
metadata:
  name: ci-runner
spec:
  agentImage: quay.io/kubeopencode/kubeopencode-agent-opencode:latest
  executorImage: quay.io/kubeopencode/kubeopencode-agent-devbox:latest
  workspaceDir: /workspace
  serviceAccountName: kubeopencode-agent
  credentials:
    - name: opencode-key
      secretRef:
        name: ai-credentials
        key: opencode-key
      env: OPENCODE_API_KEY
---
apiVersion: kubeopencode.io/v1alpha1
kind: Task
metadata:
  name: update-deps
spec:
  templateRef:
    name: ci-runner
  description: |
    Update all dependencies to latest versions.
    Run tests and create a pull request.
  contexts:
    - type: Git
      git:
        repository: https://github.com/your-org/your-repo.git
        ref: main
      mountPath: code
```

### Batch Operations with Helm

Run the same task across multiple targets:

```yaml
# values.yaml
tasks:
  - name: update-service-a
    repo: service-a
  - name: update-service-b
    repo: service-b

# templates/tasks.yaml
{{- range .Values.tasks }}
---
apiVersion: kubeopencode.io/v1alpha1
kind: Task
metadata:
  name: {{ .name }}
spec:
  templateRef:
    name: ci-runner
  description: "Update dependencies for {{ .repo }}"
{{- end }}
```

### Quota (Rate Limiting)

In addition to `maxConcurrentTasks`, use a sliding time window to control the rate of new Task starts:

```yaml
spec:
  quota:
    maxTaskStarts: 10     # Maximum 10 task starts
    windowSeconds: 3600   # Per hour (sliding window)
```

| Feature | `maxConcurrentTasks` | `quota` |
|---------|----------------------|---------|
| What it limits | Simultaneous running Tasks | Rate of new Task starts |
| Time component | No (instant check) | Yes (sliding window) |
| Use case | Limit resource usage | API rate limiting |

---

## Platform Capabilities

These capabilities apply to both Live Agents and AgentTemplate Tasks.

### Context System

Provide additional context to AI agents via inline **ContextItem**:

| Type | Description |
|------|-------------|
| **Text** | Inline text content (coding standards, instructions) |
| **ConfigMap** | Content from a Kubernetes ConfigMap |
| **Git** | Clone a Git repository (public or private with secretRef) |
| **Runtime** | KubeOpenCode platform awareness system prompt |
| **URL** | Fetch content from remote HTTP/HTTPS URL |

Contexts without `mountPath` are aggregated into `.kubeopencode/context.md`. See [Architecture](architecture.md) for full field reference and examples.

### OpenCode Configuration

Specify model and OpenCode settings via inline JSON:

```yaml
spec:
  config: |
    {
      "$schema": "https://opencode.ai/config.json",
      "model": "opencode/big-pickle",
      "small_model": "opencode/big-pickle"
    }
```

OpenCode supports [75+ AI providers](https://opencode.ai) — Anthropic, OpenAI, Google, AWS Bedrock, Azure OpenAI, and more. Change the `model` field and provide the corresponding API key via `credentials`.

### Task Stop

Stop a running Task gracefully:

```bash
kubectl annotate task my-task kubeopencode.io/stop=true
```

The controller deletes the Pod, sets status to `Completed` with a `Stopped` condition. Logs are lost when stopped.

### Enterprise Features

KubeOpenCode is designed for enterprise environments:

| Feature | Description | Configuration |
|---------|-------------|---------------|
| **Custom CA Certificates** | Trust private CAs for internal HTTPS/Git servers | `spec.caBundle.configMapRef` or `spec.caBundle.secretRef` |
| **HTTP/HTTPS Proxy** | Route traffic through corporate proxy | `spec.proxy` (Agent-level) or `KubeOpenCodeConfig.spec.proxy` (cluster-level) |
| **Private Registries** | Pull images from authenticated registries | `spec.imagePullSecrets` |
| **Pod Security** | Restricted security context by default (Kubernetes PSS Restricted) | `spec.podSpec.securityContext` and `spec.podSpec.podSecurityContext` |
| **Pod Configuration** | Node selectors, tolerations, runtime classes | `spec.podSpec.scheduling`, `spec.podSpec.runtimeClassName` |
| **RBAC** | Kubernetes-native access control | `spec.serviceAccountName` + Role/RoleBinding |

See [Security](security.md) for credential management and best practices. See [Architecture](architecture.md) for detailed YAML examples of each feature.

---

## Next Steps

- [Getting Started](getting-started.md) — Quick local setup with Kind
- [Agent Images](agent-images.md) — Build custom agent images
- [Security](security.md) — RBAC, credential management, and best practices
- [Architecture](architecture.md) — System design, API reference, and detailed configuration examples
