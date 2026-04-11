# ADR 0032: Harness Strategy — Lessons from Anthropic Managed Agents

## Status

Proposed

## Context

Anthropic launched [Claude Managed Agents](https://www.anthropic.com/engineering/managed-agents) in April 2026, a hosted agent orchestration service for enterprise API users. Its architecture solves the same class of problems as KubeOpenCode — agent loop management, state persistence, tool execution, container isolation, credential management — but through a fully managed SaaS model.

This ADR documents the architectural comparison between Managed Agents and KubeOpenCode, analyzes the emerging "closed ecosystem" trend in agent platforms, and proposes a phased harness strategy for KubeOpenCode.

### Anthropic Managed Agents Architecture

Managed Agents is built around four core abstractions:

- **Agent**: Reusable configuration (model, system prompt, tools, MCP servers, skills), version-managed.
- **Environment**: Container template (pre-installed packages, network rules).
- **Session**: Append-only event log, persisted independently of harness and container. Supports `getEvents()` for flexible context querying.
- **Events**: Message stream via SSE, durably stored.

The key architectural insight is **"decouple the brain from the hands"** — separating the system into three independently-failing components:

1. **Brain** (Harness + Claude): Stateless agent loop. Calls containers via `execute(name, input) → string`. Recoverable via `wake(sessionId)`.
2. **Hands** (Container/Sandbox): Short-lived, cattle-not-pets. Provisioned on demand, failure is a tool-call error.
3. **Session** (Event log): Durable, external to both brain and hands. Stores complete event history, making context compaction reversible.

The security boundary is structural: credentials **never enter** the sandbox. Git tokens are wired into `git remote` during initialization. MCP OAuth tokens are stored in a vault and accessed through a proxy. Even successful prompt injection cannot exfiltrate credentials or create new sessions.

**Reference:** [Scaling Managed Agents: Decoupling the brain from the hands](https://www.anthropic.com/engineering/managed-agents)

### KubeOpenCode Architecture Comparison

| Dimension | Managed Agents | KubeOpenCode |
|---|---|---|
| Deployment | Managed SaaS | Self-hosted Kubernetes |
| Core abstractions | Session + Harness + Sandbox | Task + Agent + Context |
| State management | Session log (external event stream) | etcd (CRDs) + PVC (sessions/workspace) |
| Execution isolation | Container-as-tool (`execute()`) | Pod (direct or attach to Agent Deployment) |
| Security model | Credentials never in sandbox (Vault + Proxy) | K8s Secrets + RBAC + RuntimeClass (ADR 0019) |
| Agent lifecycle | Harness is stateless, `wake(sessionId)` recovery | Agent is Deployment+Service, suspend/standby (ADR 0024) |
| Brain-hands coupling | Fully decoupled | Coupled in `agentRef` mode (OpenCode + workspace in same Pod) |
| Multi-model | Claude only | Any AI via `agentImage`/`executorImage` |
| Target user | API developers (embed into products) | Platform engineers (K8s-native operations) |

### The Closed Ecosystem Problem

Anthropic's stated design philosophy — "harnesses encode assumptions that go stale" — is technically correct from a model-capability perspective. However, it obscures a strategic reality: **each major AI vendor is building its own closed agent platform**.

| Vendor | Agent Platform | Lock-in Point |
|---|---|---|
| Anthropic | Managed Agents | Proprietary Session API, Harness, credential vault |
| OpenAI | Agents SDK + Codex | Proprietary tool calling protocol, Responses API |
| Google | ADK + Vertex AI Agent Engine | Vertex ecosystem binding |
| AWS | Bedrock Agents | Step Functions orchestration, IAM binding |

Each vendor optimizes for its own vertical integration: Model → Harness → Tool execution → State management → Billing. There is **no convergence toward a standard harness interface**. MCP is the only cross-platform protocol, but it only standardizes the tool-connection layer ("hands" interface), not:

- Agent loop / harness interface
- Session / state management
- Context engineering patterns
- Security boundary models

This is an Apple-style closed ecosystem strategy, not an open standard approach. Each vendor will build its own "best" harness and lock customers into it.

## Decision

We adopt a **phased harness strategy**: deep integration with OpenCode now, architectural readiness for harness diversity later.

### Phase 1: Deep OpenCode Integration (Current — v0.x)

Continue with OpenCode as the sole harness. Focus on delivering core platform capabilities:

- Complete MCP support (ADR 0026-mcp)
- Skills system maturation (ADR 0026-skills)
- OTel observability (ADR 0031)
- Credential isolation improvements (see Phase 2)
- Agent lifecycle polish (standby, Git sync, concurrency/quota)

**Rationale:** Trying to be harness-agnostic too early would result in poor abstraction. We need production experience with one harness before we can identify the right interface boundaries. Kubernetes itself followed this pattern — deep Docker integration first, CRI abstraction later.

### Phase 2: Security Boundary Hardening (Next — v0.x+)

Learn from Managed Agents' credential isolation model. Introduce:

1. **MCP Credential Proxy Sidecar**: When implementing MCP support (ADR 0026-mcp), store MCP OAuth/API tokens in a Kubernetes Secret accessible only to a sidecar proxy container — not to the agent container. The agent calls MCP tools via the sidecar, which injects credentials. This follows the same pattern as Managed Agents' vault-backed proxy.

2. **Git Token Isolation**: During `git-init`, wire tokens into `git remote` URL and then remove the Secret mount. The token remains in the git config inside the container, but is not exposed as an environment variable or file.

3. **Short-lived Credentials via External Secrets Operator**: Document and test integration with ESO / Vault for automatic credential rotation. This aligns with ADR 0019's long-term strategy.

**Rationale:** Credential isolation is a universally valuable security improvement, independent of harness choice. Managed Agents proved this pattern works at scale.

### Phase 3: Harness Interface Awareness (Future — v1.x)

As the market evolves, observe whether harness standards emerge. Prepare for harness diversity by:

1. **Documenting the implicit harness contract**: KubeOpenCode already has a de-facto harness interface — the two-container pattern (`agentImage` copies binary to `/tools`, `executorImage` runs the server), the `command` field, the port convention, the `OPENCODE_CONFIG` / `OPENCODE_PERMISSION` environment variables. Document this as a "Harness Provider Interface" specification.

2. **Validating with a second harness**: When a candidate emerges (e.g., a Claude Code server mode, an open-source alternative), attempt to run it on KubeOpenCode without controller changes. Success validates the interface; failure reveals what needs abstraction.

3. **Evaluating CRI-like formalization**: If multiple harnesses prove viable, consider formalizing the interface — similar to how Kubernetes moved from Docker-specific code to CRI. This is a v1.x+ concern.

**Rationale:** Premature standardization produces wrong interfaces. The right time to abstract is when we have at least two concrete implementations to generalize from. The existing `agentImage` + `executorImage` + `command` design already provides significant flexibility without formal abstraction.

### What We Explicitly Do NOT Do

- **Do not build a "meta-harness"**: Managed Agents is Anthropic's meta-harness because they control the model. We don't control any model; our value is the Kubernetes-native platform layer.
- **Do not chase harness interoperability with Managed Agents**: Their interfaces are proprietary and will change. Our API surface (Task, Agent, Context) serves different users with different needs.
- **Do not try to be "the open standard"**: Standards emerge from adoption, not from declaration. If KubeOpenCode gains adoption, its patterns may become de facto standards. But designing for this prematurely is a trap.

## Consequences

### Positive

- **Focus**: Phase 1 keeps the team focused on delivering a complete, polished experience with one harness rather than building thin support for many.
- **Security improvement**: Phase 2's credential isolation benefits all users regardless of future harness decisions.
- **Architectural optionality**: Phase 3 preserves the ability to support multiple harnesses without committing to premature abstraction.
- **Clear positioning**: KubeOpenCode is differentiated as the "self-hosted, multi-model, Kubernetes-native" alternative for enterprises that need control. Managed Agents is for enterprises that want convenience.

### Negative

- **OpenCode dependency risk**: Deep coupling with OpenCode means KubeOpenCode's capabilities are bounded by OpenCode's evolution. Mitigation: OpenCode is open-source; we can fork or contribute.
- **Late mover on harness diversity**: If a clear harness standard emerges, we may need to catch up. Mitigation: The `agentImage`/`executorImage` pattern already supports swapping; the gap is mostly documentation and testing.

### Key Architectural Lessons from Managed Agents

1. **Session as queryable event stream** (not just persistence): Worth exploring in future — a structured event log per Task/Agent that supports `getEvents()` style queries. This enables cross-Task knowledge transfer and makes context compaction reversible. Related to ADR 0031 (OTel) and ADR 0013 (token tracking).

2. **Many Brains, Many Hands**: The ability for one "brain" to operate multiple "hands" (execution environments) is a compelling pattern. In Kubernetes, this could map to an Agent orchestrating work across multiple Pods or even clusters. This is a future exploration area, not a current priority.

3. **Cattle not pets for everything**: Managed Agents makes harness, container, and session all independently recoverable. KubeOpenCode's Agent is closer to a "pet" (persistent Deployment with PVCs). The standby mechanism (ADR 0024) partially addresses this, but full cattle-like recovery (e.g., `wake(sessionId)` equivalent) is not yet supported.

## References

- [Scaling Managed Agents: Decoupling the brain from the hands](https://www.anthropic.com/engineering/managed-agents) — Anthropic Engineering Blog, April 2026
- [Building Effective Agents](https://www.anthropic.com/engineering/building-effective-agents) — Anthropic, 2025
- [Effective Harnesses for Long-Running Agents](https://www.anthropic.com/engineering/effective-harnesses-for-long-running-agents) — Anthropic, 2025
- [The History of Pets vs Cattle](https://cloudscaling.com/blog/cloud-computing/the-history-of-pets-vs-cattle/) — Randy Bias
- [The Bitter Lesson](http://www.incompleteideas.net/IncIdeas/BitterLesson.html) — Rich Sutton
- ADR 0019: Web Terminal Credential Security Strategy
- ADR 0022: Agent Always Running — Unified Execution Model
- ADR 0024: Agent Standby — Unified Suspend/Resume Lifecycle
- ADR 0026-mcp: MCP Server Support in Agent API
- ADR 0031: OpenTelemetry Observability for Tasks and Agents
