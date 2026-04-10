# Agent Templates

AgentTemplate serves two purposes:

1. **Reusable base configuration for Agents**: Teams define shared settings (images, contexts, credentials) in one template. Individual users create Agents that reference it via `templateRef`.
2. **Blueprint for ephemeral tasks**: Tasks can reference a template directly via `templateRef` to run one-off, ephemeral Pods without a persistent Agent.

## Creating a Template

```yaml
apiVersion: kubeopencode.io/v1alpha1
kind: AgentTemplate
metadata:
  name: team-config
spec:
  executorImage: ghcr.io/kubeopencode/kubeopencode-agent-devbox:latest
  workspaceDir: /workspace
  serviceAccountName: kubeopencode-agent
  contexts:
    - name: coding-standards
      type: Text
      text: "Follow team coding standards..."
  credentials:
    - name: github-token
      secretRef:
        name: shared-github-creds
        key: token
      env: GITHUB_TOKEN
```

## Creating an Agent from a Template

```yaml
apiVersion: kubeopencode.io/v1alpha1
kind: Agent
metadata:
  name: my-agent
spec:
  templateRef:
    name: team-config
  profile: "My personal development agent"
  # Required fields (even with template):
  executorImage: ghcr.io/kubeopencode/kubeopencode-agent-devbox:latest
  workspaceDir: /workspace
  serviceAccountName: kubeopencode-agent
  # Instance-specific settings:
  maxConcurrentTasks: 3
```

## Running Ephemeral Tasks from a Template

Tasks can reference a template directly instead of a running Agent. This creates an ephemeral Pod that runs standalone and terminates when done — ideal for batch operations and CI/CD:

```yaml
apiVersion: kubeopencode.io/v1alpha1
kind: Task
metadata:
  name: one-off-task
spec:
  templateRef:
    name: team-config
  description: |
    Update all dependencies and run tests.
```

The Task controller creates a standalone Pod using the template's configuration. No persistent Agent is needed. Exactly one of `agentRef` or `templateRef` must be set on a Task.

## Merge Behavior

When an Agent references a template:
- **Scalar fields** (images, workspaceDir, config, etc.): Agent wins if set, otherwise template value
- **List fields** (contexts, credentials, imagePullSecrets): Agent's list **replaces** the template's (not appended)
- **Agent-only fields** (profile, port, persistence, suspend): Always from Agent

## Tracking

Agents using a template automatically get the label `kubeopencode.io/agent-template: <name>`,
enabling template-based queries:

```bash
# List all agents using a template
kubectl get agents -l kubeopencode.io/agent-template=team-config
```
