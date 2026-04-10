# Persistence & Lifecycle

## Persistence

By default, Agents use ephemeral storage ([EmptyDir](https://kubernetes.io/docs/concepts/storage/volumes/#emptydir) volumes). EmptyDir is a temporary volume that Kubernetes provides automatically — no configuration required. When the server pod restarts (due to crashes, node drains, or upgrades), all EmptyDir data is lost.

Persistence replaces EmptyDir with [PersistentVolumeClaims (PVCs)](https://kubernetes.io/docs/concepts/storage/persistent-volumes/), so data survives pod restarts. Session and workspace persistence are configured independently. See [Live Agents](live-agents.md) for the full overview.

:::info Kubernetes Storage Concepts
Kubernetes storage has three layers:

| Layer | What It Is | Who Manages It |
|-------|-----------|----------------|
| **StorageClass** | Defines *how* storage is provisioned (e.g., AWS EBS, GCP Persistent Disk, local-path). | Cluster administrator |
| **PersistentVolume (PV)** | An actual block of storage in your cluster. | Created automatically by the StorageClass (dynamic provisioning) |
| **PersistentVolumeClaim (PVC)** | A request for storage — "I need 10Gi of disk". | Created automatically by the KubeOpenCode controller |

**You only need to configure `persistence` in your Agent spec.** The controller creates PVCs automatically, and Kubernetes provisions the underlying storage via your cluster's StorageClass.
:::

### Configuration

Add `persistence` to your Agent spec:

```yaml
apiVersion: kubeopencode.io/v1alpha1
kind: Agent
metadata:
  name: persistent-agent
spec:
  executorImage: ghcr.io/kubeopencode/kubeopencode-agent-devbox:latest
  workspaceDir: /workspace
  serviceAccountName: kubeopencode-agent
  port: 4096
  persistence:
    sessions:
      size: "2Gi"               # default: 1Gi
    workspace:
      size: "20Gi"              # default: 10Gi
```

Both fields are optional — you can persist sessions only, workspace only, or both.

### Specifying a StorageClass

By default, PVCs use your cluster's default StorageClass. To use a specific StorageClass (e.g., a faster SSD tier):

```yaml
spec:
  persistence:
    sessions:
      storageClassName: "fast-ssd"
      size: "2Gi"
    workspace:
      storageClassName: "fast-ssd"
      size: "50Gi"
```

To check which StorageClasses are available on your cluster:

```bash
kubectl get storageclass
```

The one marked `(default)` is used when `storageClassName` is omitted.

### How It Works

**Without persistence (default):**
- Workspace directory uses an EmptyDir volume — available immediately, no setup needed
- Session data (SQLite DB) lives inside the container's ephemeral storage
- Pod restart = data lost, git repos re-cloned by init containers

**Session persistence** (`persistence.sessions`):
- The controller creates a PVC named `{agent-name}-server-sessions` and mounts it at `/data/sessions`
- `OPENCODE_DB` env var is set to `/data/sessions/opencode.db`
- Conversation history survives pod restarts

**Workspace persistence** (`persistence.workspace`):
- The controller replaces the workspace EmptyDir with a PVC named `{agent-name}-server-workspace`
- Git-cloned repos, AI-modified files, and in-progress work survive pod restarts
- git-init skips cloning when the repository already exists on the PVC

### PVC Lifecycle

- PVCs are **created automatically** by the controller when persistence is configured — you never need to create them manually
- Each PVC has an `OwnerReference` pointing to the Agent
- When the Agent is deleted, its PVCs are automatically garbage-collected by Kubernetes
- To retain data after Agent deletion, configure the StorageClass with `reclaimPolicy: Retain`
- PVC specs are immutable after creation; to change size or storage class, delete and recreate the Agent

### Limitations

- When workspace persistence is enabled and git-init detects an existing repository, it skips cloning.
  If the Agent's git ref changes, the existing checkout is not automatically updated.
- Sessions and workspace can be configured independently (one, both, or neither)

## Suspend/Resume

Agents can be suspended to save compute resources. `spec.suspend` is the single switch — both humans and the controller operate on it.

### Manual Suspend

Set `spec.suspend: true` to scale the Deployment to 0 replicas. PVCs and Service are retained.

```yaml
apiVersion: kubeopencode.io/v1alpha1
kind: Agent
metadata:
  name: suspendable-agent
spec:
  workspaceDir: /workspace
  serviceAccountName: kubeopencode-agent
  suspend: true    # scales deployment to 0
  persistence:
    sessions:
      size: "1Gi"
```

- Tasks targeting a suspended agent enter `Queued` phase with reason `AgentSuspended`
- Set `suspend: false` to resume — queued tasks start automatically
- API: `POST /api/v1/namespaces/{ns}/agents/{name}/suspend` and `.../resume`
- Cannot suspend while tasks are running (API returns 409 Conflict)

### Standby (Automatic Suspend/Resume)

Configure `spec.standby` for automatic lifecycle management. The controller manages `spec.suspend` automatically — suspending after idle timeout and resuming when new Tasks arrive.

```yaml
apiVersion: kubeopencode.io/v1alpha1
kind: Agent
metadata:
  name: auto-scaling-agent
spec:
  workspaceDir: /workspace
  serviceAccountName: kubeopencode-agent
  standby:
    idleTimeout: "30m"    # auto-suspend after 30 minutes with no tasks
  persistence:
    sessions:
      size: "1Gi"
    workspace:
      size: "10Gi"
```

**How it works:**

1. Agent starts running normally
2. All Tasks complete and no active connections → idle timer starts (`status.idleSince` is set)
3. After 30 minutes with no new Tasks and no active connections → controller sets `spec.suspend = true` → Deployment scales to 0
4. New Task arrives → controller sets `spec.suspend = false` → Deployment scales back to 1
5. Agent becomes ready (~30-60s cold start) → queued Task executes

**Connection-aware idle detection:** The standby system considers both active Tasks and active user connections (web terminal, `kubeoc agent attach`). While a user is connected, the idle timer is deferred — your interactive session won't be interrupted by auto-suspend. Connection activity is tracked via the `kubeopencode.io/last-connection-active` annotation, updated every 60 seconds by the web terminal handler and CLI. See [ADR 0028](https://github.com/kubeopencode/kubeopencode/blob/main/docs/adr/0028-connection-aware-standby.md) for details.

**Manual override:** Even with standby configured, you can still manually suspend/resume. The controller will resume automatically when new Tasks arrive.

**Condition reasons:**
- `Suspended: True, reason: UserRequested` — suspended (no standby configured)
- `Suspended: True, reason: Standby` — suspended with standby configured
- `Suspended: False, reason: Active` — running normally
- `StandbyConfigWarning: True, reason: IdleTimeoutTooShort` — `idleTimeout` is less than 2 minutes (connection heartbeat staleness is automatically degraded)

**Best used with `persistence`**: PVCs survive restarts, so session history and workspace files don't need to be re-initialized on resume.

### UI

The Agent detail page shows a **Suspend/Resume** button, and the agents list shows a "Suspended" badge. When standby is configured, the detail page shows the standby configuration and current idle duration.
