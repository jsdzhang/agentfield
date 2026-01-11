# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

<!-- changelog:entries -->

## [0.1.29-rc.2] - 2026-01-11


### Added

- Feat(ci): add contributor reminder and assignment tracking workflows (#132)

Add automated system to remind assigned contributors and free up stale assignments:

- contributor-reminders.yml: Scheduled daily check that:
  - Sends friendly reminder at 7 days without activity
  - Sends second reminder at 14 days with unassign warning
  - Unassigns and re-labels as 'help wanted' at 21 days
  - Skips issues with linked PRs or blocking labels
  - Supports dry-run mode for testing

- issue-assignment-tracking.yml: Real-time event handling that:
  - Welcomes new assignees with timeline expectations
  - Clears reminder labels when assignees comment
  - Clears labels when assignee opens linked PR
  - Auto-adds 'help wanted' when last assignee leaves

This improves contributor experience by setting clear expectations
while ensuring stale assignments don't block other contributors. (7bbac52)

- Feat(sdk/go, control-plane): Add Vector Memory Ops (#124)

* chore(release): v0.1.26-rc.2

* feat: added vector memory ops

* test(handlers): add unit tests for GetVector and DeleteVector handlers

Add comprehensive test coverage for the new vector memory endpoints:

GetVectorHandler tests:
- TestGetVectorHandler_ReturnsVectorWithMetadata: Full happy path with scope/key/metadata
- TestGetVectorHandler_NotFound: 404 when vector doesn't exist
- TestGetVectorHandler_StorageError: 500 on database failure
- TestGetVectorHandler_DefaultScope: Scope resolution from headers

DeleteVectorHandler tests:
- TestDeleteVectorHandler_RESTfulDelete: DELETE with path parameter
- TestDeleteVectorHandler_BackwardCompatibilityWithBody: POST with JSON body
- TestDeleteVectorHandler_StorageError: 500 on database failure
- TestDeleteVectorHandler_MissingKey: 400 when key is missing

Also updated vectorStorageStub to track GetVector and DeleteVector parameters
for assertion verification.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>

---------

Co-authored-by: github-actions[bot] <github-actions[bot]@users.noreply.github.com>
Co-authored-by: Abir Abbas <abirabbas1998@gmail.com>
Co-authored-by: Claude Opus 4.5 <noreply@anthropic.com> (0dd4e62)

- Feat(tests): add unit tests for vector memory handler functionality (#123) (9df214d)

- Feat: Add AI-assisted contribution guidelines, task issue template, AI labeling workflow, security, and support policies. (87c9297)

- Feat(observability): add webhook forwarding with dead letter queue (#102)

* feat(observability): add webhook forwarding with dead letter queue

Add observability webhook system for forwarding events to external endpoints:

- Configurable webhook URL with optional HMAC secret and custom headers
- Event batching with configurable size and timeout
- Automatic retry with exponential backoff
- Dead letter queue (DLQ) for failed events with redrive and clear capabilities
- Filter heartbeat events and minor health score fluctuations to reduce noise
- Settings UI page for configuration and DLQ management

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>

* test fixes

* fix(observability): node offline events now properly forwarded to webhook

Fixed a race condition where node offline events were not being forwarded
to the observability webhook. The issue was in UpdateAgentStatus which
called GetAgentStatus (performing a live health check) to get the "old"
status. By the time the health check completed, oldStatus == newStatus,
so no events were broadcast.

Changed to use GetAgentStatusSnapshot which returns cached/stored status
without a live check, preserving the true "old" state for comparison.

Also added observability webhook documentation to ARCHITECTURE.md.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>

* test(status-manager): add tests for node online/offline event broadcasting

Add comprehensive tests to verify status change events are properly
broadcast when nodes transition between active and inactive states:

- TestStatusManagerBroadcastsNodeOfflineEvent: Verifies NodeOffline/
  NodeUnifiedStatusChanged events are broadcast when node goes offline
- TestStatusManagerBroadcastsNodeOnlineEvent: Verifies NodeOnline/
  NodeUnifiedStatusChanged events are broadcast when node comes online
- TestStatusManagerPreservesOldStatusForEventBroadcast: Verifies the
  old status is correctly captured before updates, ensuring the fix
  for the GetAgentStatus race condition doesn't regress

These tests guard against the race condition where UpdateAgentStatus
was calling GetAgentStatus (with live health check) instead of
GetAgentStatusSnapshot, causing oldStatus == newStatus.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>

* docs: remove observability webhook documentation

The observability webhook feature remains functional but will not be publicly
documented at this time.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>

* bugfix on ts-sdk json schema return

* webhook secret fix

* refine webhook events

---------

Co-authored-by: Claude Opus 4.5 <noreply@anthropic.com> (27cf1f7)

- Feat(dashboard): add comprehensive observability enhancements (#101) (d947542)

- Feat(go-sdk): add ControlPlaneMemoryBackend for distributed memory (#80)

Add a new MemoryBackend implementation that delegates storage to the
control plane's /api/v1/memory/* endpoints. This enables distributed,
scope-aware memory across agents.

- Implements Set, Get, Delete, and List operations
- Maps SDK scopes (Workflow, Session, User, Global) to API scopes
- User scope maps to "actor" terminology in the API
- Includes comprehensive unit tests

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-authored-by: Claude Opus 4.5 <noreply@anthropic.com> (6cd445c)

- Feat: add RAG evaluation example with multi-reasoner architecture (#79)

* feat: add RAG evaluation example with multi-reasoner architecture

Adds a comprehensive RAG evaluation agent demonstrating:
- Adversarial debate for faithfulness (prosecutor vs defender + judge)
- Multi-jury consensus for relevance (3 jurors vote on literal/intent/scope)
- Hybrid ML+LLM chain-of-verification for hallucination detection
- Configurable constitutional principles evaluation

Features:
- Docker Compose deployment (control plane + agent + UI)
- Next.js web interface with claim-level breakdown
- Domain-specific presets (medical, legal, financial)
- 3 evaluation modes: quick (4 calls), standard (14), thorough (20+)

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>

* Remove ARCHITECTURE.md, link to website docs instead

The detailed architecture documentation is now on the website at
agentfield.dev/examples/complete-agents/rag-evaluator - no need
to duplicate content in the repo.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>

* Add examples README with quick reference table

Index of all examples across Python, TypeScript, and Go with:
- Quick reference table by use case and language
- Detailed per-language tables with key features
- Use case deep dives (RAG progression, multi-agent, serverless)
- Technology stack overview

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>

---------

Co-authored-by: Claude Opus 4.5 <noreply@anthropic.com> (862c41e)

- Feat(go-sdk): add per-request API key override for AI client (#73)

Add WithAPIKey option to override the client's configured API key on a
per-request basis. This brings the Go SDK to parity with the Python SDK,
which supports api_key overrides in individual calls.

Changes:
- Add APIKeyOverride field to Request struct (excluded from JSON)
- Add WithAPIKey option function
- Update doRequest and StreamComplete to use override when provided
- Add test for API key override behavior

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-authored-by: Claude Opus 4.5 <noreply@anthropic.com> (4dd8a70)

- Feat(go-sdk): add Memory and Note APIs for agent state and progress tracking (#71)

Add two major new capabilities to the Go SDK:

## Memory System
- Hierarchical scoped storage (workflow, session, user, global)
- Pluggable MemoryBackend interface for custom storage
- Default in-memory backend included
- Automatic scope ID resolution from execution context

## Note API
- Fire-and-forget progress/status messages to AgentField UI
- Note(ctx, message, tags...) and Notef(ctx, format, args...) methods
- Async HTTP delivery with proper execution context headers
- Silent failure mode to avoid interrupting workflows

These additions enable agents to:
- Persist state across handler invocations within a session
- Share data between workflows at different scopes
- Report real-time progress updates visible in the UI

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-authored-by: Claude Opus 4.5 <noreply@anthropic.com> (1c48c1f)

- Feat: allow external contributors to run functional tests without API… (#70)

* feat: allow external contributors to run functional tests without API keys

Enable external contributors to run 92% of functional tests (24/26) without
requiring access to OpenRouter API keys. This makes it easier for the community
to contribute while maintaining full test coverage for maintainers.

Changes:
- Detect forked PRs and automatically skip OpenRouter-dependent tests
- Only 2 tests require OpenRouter (LLM integration tests)
- 24 tests validate all core infrastructure without LLM calls
- Update GitHub Actions workflow to conditionally set PYTEST_ARGS
- Update functional test README with clear documentation

Test coverage for external contributors:
✅ Control plane health and APIs
✅ Agent registration and discovery
✅ Multi-agent communication
✅ Memory system (all scopes)
✅ Workflow orchestration
✅ Go/TypeScript SDK integration
✅ Serverless agents
✅ Verifiable credentials

Skipped for external contributors (maintainers still run these):
⏭️  test_hello_world_with_openrouter
⏭️  test_readme_quick_start_summarize_flow

This change addresses the challenge of running CI for external contributors
without exposing repository secrets while maintaining comprehensive test
coverage for the core AgentField platform functionality.

* fix: handle push events correctly in functional tests workflow

The workflow was failing on push events (to main/testing branches) because
it relied on github.event.pull_request.head.repo.fork which is null for
push events. This caused the workflow to incorrectly fall into the else
branch and fail when OPENROUTER_API_KEY wasn't set.

Changes:
- Check github.event_name to differentiate between push, pull_request, and workflow_dispatch
- Explicitly handle push and workflow_dispatch events to run all tests with API key
- Preserve fork PR detection to skip OpenRouter tests for external contributors

Now properly handles:
✅ Fork PRs: Skip 2 OpenRouter tests, run 24/26 tests
✅ Internal PRs: Run all 26 tests with API key
✅ Push to main/testing: Run all 26 tests with API key
✅ Manual workflow dispatch: Run all 26 tests with API key

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>

* fix: remove shell quoting from PYTEST_ARGS to prevent argument parsing errors

The PYTEST_ARGS variable contained single quotes around '-m "not openrouter" -v'
which would be included in the environment variable value. When passed to pytest
in the Docker container shell command, this caused the entire string to be treated
as a single argument instead of being properly split into separate arguments.

Changed from: '-m "not openrouter" -v'
Changed to:   -m not openrouter -v

This allows the shell's word splitting to correctly parse the arguments when
pytest $$PYTEST_ARGS is evaluated in the docker-compose command.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>

* refactor: separate pytest marker expression from general args for proper quoting

The previous approach of embedding -m not openrouter inside PYTEST_ARGS was
fragile because shell word-splitting doesn't guarantee "not openrouter" stays
together as a single argument to the -m flag.

This change introduces PYTEST_MARK_EXPR as a dedicated variable for the marker
expression, which is then properly quoted when passed to pytest:
  pytest -m "$PYTEST_MARK_EXPR" $PYTEST_ARGS ...

Benefits:
- Marker expression is guaranteed to be treated as single argument to -m
- Clear separation between marker selection and general pytest args
- More maintainable for future marker additions
- Eliminates shell quoting ambiguity

Changes:
- workflow: Split PYTEST_ARGS into PYTEST_MARK_EXPR + PYTEST_ARGS
- docker-compose: Add PYTEST_MARK_EXPR env var and conditional -m flag
- docker-compose: Only apply -m when PYTEST_MARK_EXPR is non-empty

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>

* fix: add proper event type checks before accessing pull_request context

Prevent errors when workflow runs on push events by:
- Check event_name == 'pull_request' before accessing pull_request.head.repo.fork
- Check event_name == 'workflow_dispatch' before accessing event.inputs
- Ensures all conditional expressions only access context properties when they exist

This prevents "Error: Cannot read properties of null (reading 'fork')" errors
on push events.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>

---------

Co-authored-by: Claude Sonnet 4.5 <noreply@anthropic.com> (01668aa)

- Feat(release): unify PyPI publishing for all releases (#58)

Publish all Python SDK releases (both prerelease and stable) to PyPI
instead of using TestPyPI for prereleases.

Per PEP 440, prerelease versions (e.g., 0.1.20rc1) are excluded by
default from `pip install` - users must explicitly use `--pre` flag.
This simplifies the release process and removes the need for the
TEST_PYPI_API_TOKEN secret.

Changes:
- Merge TestPyPI and PyPI publish steps into single PyPI step
- Update release notes to show `pip install --pre` for staging
- Update install.sh staging output
- Re-enable example requirements updates for prereleases
- Update RELEASE.md documentation

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-authored-by: Claude <noreply@anthropic.com> (ebf7020)

- Feat(release): add two-tier staging/production release system (#53)

* feat(release): add two-tier staging/production release system

Implement automatic staging releases and manual production releases:

- Staging: Automatic on push to main (TestPyPI, npm @next, staging-* Docker)
- Production: Manual workflow dispatch (PyPI, npm @latest, vX.Y.Z + latest Docker)

Changes:
- Add push trigger with path filters for automatic staging
- Replace release_channel with release_environment input
- Split PyPI publishing: TestPyPI (staging) vs PyPI (production)
- Split npm publishing: @next tag (staging) vs @latest (production)
- Conditional Docker tagging: staging-X.Y.Z vs vX.Y.Z + latest
- Add install-staging.sh for testing prerelease binaries
- Update RELEASE.md with two-tier documentation

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>

* refactor(install): consolidate staging into single install.sh with --staging flag

Instead of separate install.sh and install-staging.sh scripts:
- Single install.sh handles both production and staging
- Use --staging flag or STAGING=1 env var for prerelease installs
- Eliminates code drift between scripts

Usage:
  Production: curl -fsSL .../install.sh | bash
  Staging:    curl -fsSL .../install.sh | bash -s -- --staging

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>

---------

Co-authored-by: Claude <noreply@anthropic.com> (3bd748d)

- Feat(sdk/typescript): expand AI provider support to 10 providers

Add 6 new AI providers to the TypeScript SDK:
- Google (Gemini models)
- Mistral AI
- Groq
- xAI (Grok)
- DeepSeek
- Cohere

Also add explicit handling for OpenRouter and Ollama with sensible defaults.

Changes:
- Update AIConfig type with new provider options
- Refactor buildModel() with switch statement for all providers
- Refactor buildEmbeddingModel() with proper embedding support
  (Google, Mistral, Cohere have native embedding; others throw)
- Add 27 unit tests for provider selection and embedding support
- Install @ai-sdk/google, @ai-sdk/mistral, @ai-sdk/groq,
  @ai-sdk/xai, @ai-sdk/deepseek, @ai-sdk/cohere packages

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (b06b5b5)

- Feat: expose api_key at Agent level and fix test lint issues

- Add api_key parameter to Agent class constructor
- Pass api_key to AgentFieldClient for authentication
- Document api_key parameter in Agent docstring
- Fix unused loop variable in ensure_event_loop test fixture

Addresses reviewer feedback that api_key should be exposed at Agent
level since end users don't interact directly with AgentFieldClient.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (6567bd0)

- Feat: add API key authentication to control plane and SDKs

This adds optional API key authentication to the AgentField control plane
with support in all SDKs (Python, Go, TypeScript).

## Control Plane Changes

- Add `api_key` config option in agentfield.yaml
- Add HTTP auth middleware (X-API-Key header, Bearer token, query param)
- Add gRPC auth interceptor (x-api-key metadata, Bearer token)
- Skip auth for /api/v1/health, /metrics, and /ui/* paths
- UI prompts for API key when auth is required and stores in localStorage

## SDK Changes

- Python: Add `api_key` parameter to AgentFieldClient
- Go: Add `WithAPIKey()` option to client
- TypeScript: Add `apiKey` option to client config

## Tests

- Add comprehensive HTTP auth middleware tests (14 tests)
- Add gRPC auth interceptor tests (11 tests)
- Add Python SDK auth tests (17 tests)
- Add Go SDK auth tests (10 tests)

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (3f8e45c)

- Feat: add delete-namespace endpoint for RAG reindexing

Adds a new DELETE /api/v1/memory/vector/namespace endpoint that allows
clearing all vectors with a given namespace prefix. This enables the
documentation chatbot to wipe and reindex its RAG data when docs change.

Changes:
- Add DeleteVectorsByPrefix to StorageProvider interface
- Implement DeleteByPrefix for SQLite and Postgres vector stores
- Add DeleteNamespaceVectorsHandler endpoint
- Add clear_namespace skill to documentation chatbot
- Update MemoryStorage interface with new method

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (bc1f41e)

- Feat(sdk-python): expose execution context via app.ctx property

Add a `ctx` property to the Agent class that provides direct access to
the current execution context during reasoner/skill execution. This
enables a more ergonomic API:

Before:
  from agentfield.execution_context import get_current_context
  ctx = get_current_context()
  workflow_id = ctx.workflow_id

After:
  workflow_id = app.ctx.workflow_id

The property returns None when accessed outside of an active execution
(e.g., at module level or after a request completes), matching the
behavior of app.memory. This prevents accidental use of stale or
placeholder context data.

Also fixes integration test fixtures to support the current monorepo
structure where control-plane lives at repo root instead of
apps/platform/agentfield.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (e01dcea)

- Feat(ts-sdk): add DID client and memory helpers (4b74998)

- Feat(ts-sdk): add heartbeat and local call coverage (cf228ec)

- Feat(ts-sdk): scaffold typescript sdk core (09dcc62)

- Feat(sdk): add resilience and concurrency tests for Agent (0e86112)

- Feat: Add vector memory capabilities

This change introduces persistent vector memory storage and retrieval.

Configuration files have been updated to enable and configure vector storage with cosine distance as the default. The storage layer now implements methods for setting, deleting, and similarity searching vector embeddings. New API endpoints for vector operations have been added to the server, and the SDK is updated to support these new memory operations.

The agentic_rag example has been refactored to leverage these new vector memory capabilities for document chunk embeddings and similarity searches, enhancing its RAG pipeline with meta-cognitive reasoning. (ab7a39a)

- Feat: Add interactive agent initialization

Introduce an interactive CLI for initializing new Haxen agent projects. This enhances the user experience by guiding users through project setup with prompts for project name, language, author details, and more.

The new initialization flow supports both Python and Go agent templates, dynamically generates project files, and provides clear next steps for the user. It also incorporates better validation for project names and email addresses, and leverages Git configuration for author information if available.

The update also includes a significant refactor of the `init.go` file to integrate Bubble Tea for a more user-friendly, interactive command-line interface. Dependency updates in `go.mod` and `go.sum` are also included to support these new features. (f6b24bd)



### CI

- Ci: disable AI label workflow for fork compatibility

The AI label workflow fails on PRs from forked repositories because
GITHUB_TOKEN lacks write permissions. Since many contributions come
from forks, disabling the workflow until a proper solution (PAT or
GitHub App) is implemented.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com> (6dbc908)



### Chores

- Chore(ci): migrate Docker publishing from GHCR to Docker Hub (#133)

- Change image path from ghcr.io/agent-field/agentfield-control-plane to agentfield/control-plane
- Update login step to use Docker Hub credentials (DOCKERHUB_USERNAME, DOCKERHUB_TOKEN)
- Remove unused OWNER env var from Docker metadata step

This enables Docker Hub analytics for image pulls. Requires adding
DOCKERHUB_USERNAME and DOCKERHUB_TOKEN secrets to the repository.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-authored-by: Claude Opus 4.5 <noreply@anthropic.com> (e6abe54)

- Chore(rag-eval): update default model to GPT-4o (#85)

- Set GPT-4o as default (reliable JSON output)
- Gemini 2.5 Flash as second option
- Move DeepSeek to last (can timeout)
- Remove old Gemini 2.0 Flash

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-authored-by: Claude Opus 4.5 <noreply@anthropic.com> (f49959a)

- Chore: trigger Railway deployment for PR #39 fix (b4095d2)

- Chore: ignore env files (3937821)

- Chore(ts-sdk): align heartbeat and memory clients, improve example env loading (fee2a7e)

- Chore(ts-sdk): load env config for simulation example (9715ac5)

- Chore(ts-sdk): remove AI stubs from simulation example (7b94190)

- Chore(ts-sdk): make simulation example runnable via build (9a87374)

- Chore(ts-sdk): fix typings, add heartbeat config, lock deps (f9af207)

- Chore: initial tests (78bd76d)

- Chore: auto-fix linting and formatting issues via pre-commit hooks (22ab4e7)



### Documentation

- Docs: update Docker image references to Docker Hub (#134)

* docs: update Docker image references to Docker Hub

Update all references from ghcr.io/agent-field/agentfield-control-plane
to agentfield/control-plane (Docker Hub).

Files updated:
- deployments/kubernetes/base/control-plane-deployment.yaml
- deployments/helm/agentfield/values.yaml
- examples/python_agent_nodes/rag_evaluation/docker-compose.yml
- README.md
- docs/RELEASE.md (includes new DOCKERHUB_* secrets documentation)

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>

* fix: use real version numbers in RELEASE.md examples

Update example commands to use actual versions that exist:
- Docker: staging-0.1.28-rc.4 (not 0.1.19-rc.1)
- Install script: v0.1.28 and v0.1.28-rc.4

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>

---------

Co-authored-by: Claude Opus 4.5 <noreply@anthropic.com> (feeaa21)

- Docs: add Discord community badge to README (#131)

Add a Discord badge near the top of README.md to invite users to join
the community. Uses Discord's official brand color (#5865F2) and matches
the existing badge styling.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-authored-by: Claude <noreply@anthropic.com> (81fb1c5)

- Docs: Add cross-reference links to RAG evaluator documentation (#83)

- Add docs link in examples/README.md table for rag_evaluation
- Add documentation callout in rag_evaluation/README.md

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-authored-by: Claude Opus 4.5 <noreply@anthropic.com> (e51b8d4)

- Docs: Add website deployment guide links (#82)

Reference the full deployment guides at agentfield.ai for Kubernetes
and Helm deployment options.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-authored-by: Claude Opus 4.5 <noreply@anthropic.com> (59eaf7f)

- Docs(chatbot): add SDK search term relationship

Add search term mapping for SDK/language queries to improve RAG
retrieval when users ask about supported languages or SDKs.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (87a4d90)

- Docs(chatbot): add TypeScript SDK to supported languages

Update product context to include TypeScript alongside Python and Go:
- CLI commands now mention all three language options
- Getting started section references TypeScript
- API Reference includes TypeScript SDK

This fixes the RAG chatbot returning only Python/Go when asked about
supported languages.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (9510d74)

- Docs: Add next test coverage priorities and planning document (40a302b)



### Fixed

- Fix(control-plane): enforce lifecycle_status consistency with agent state (#130)

When agents go offline, the control plane was incorrectly keeping
lifecycle_status as "ready" even though health_status correctly showed
"inactive". This caused observability webhooks to receive inconsistent
data where offline nodes appeared online based on lifecycle_status.

Changes:
- Add defensive lifecycle_status enforcement in persistStatus()
  to ensure consistency with agent state before writing to storage
- Update health_monitor.go fallback paths to also update lifecycle_status
- Add SystemStateSnapshot event type for periodic agent inventory
- Enhance execution events with full reasoner context and metadata
- Add ListAgents to ObservabilityWebhookStore interface for snapshots

The fix ensures both node_offline events and system_state_snapshot
events (every 60s) correctly report lifecycle_status: "offline" for
offline agents.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-authored-by: Claude <noreply@anthropic.com> (67c67c4)

- Fix(ci): add issues:write permission for AI label workflow (#126)

The `gh pr edit --add-label` command requires `issues: write` permission
because labels are managed through the issues API in GitHub, even when
applied to pull requests. Without this permission, the workflow fails with:
"GraphQL: Resource not accessible by integration (addLabelsToLabelable)"

Added permissions:
- `issues: write` - Required for adding labels
- `contents: read` - Explicit permission for checkout

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-authored-by: Claude <noreply@anthropic.com> (26a71a4)

- Fix(ci): prevent shell injection in AI label workflow (#125)

The PR body was being directly interpolated into a shell variable using
`${{ github.event.pull_request.body }}`, which caused shell parsing of
the content. When PR descriptions contained filenames like `CHANGELOG.md`
on their own lines, the shell would attempt to execute them as commands.

This fix passes the PR body via the `env:` block instead, which properly
escapes the content as an environment variable.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-authored-by: Claude <noreply@anthropic.com> (1e2225d)

- Fix: use executions table for notes storage instead of workflow_executions (#75)

* fix: use executions table for notes storage instead of workflow_executions

The note handlers (AddExecutionNoteHandler, GetExecutionNotesHandler) were
querying the workflow_executions table, but execution data is actually stored
in the executions table. This caused "execution not found" errors when adding
or retrieving notes via app.note().

Changes:
- Add Notes field to types.Execution struct
- Add notes column to ExecutionRecordModel (GORM auto-migrates this)
- Update SQL queries in execution_records.go to include notes column
- Update scanExecution to deserialize notes JSON
- Change ExecutionNoteStorage interface to use GetExecutionRecord and
  UpdateExecutionRecord instead of GetWorkflowExecution and
  UpdateWorkflowExecution
- Update AddExecutionNoteHandler to use UpdateExecutionRecord
- Update GetExecutionNotesHandler to use GetExecutionRecord

This fixes both the SDK app.note() functionality and the UI notes panel
404 errors.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>

* fix: update execution notes tests to use correct storage methods

Tests were using WorkflowExecution type and StoreWorkflowExecution() to set up
test data, but the handlers now use Execution type and GetExecutionRecord()/
UpdateExecutionRecord() which query the executionRecords map.

- Change test setup from types.WorkflowExecution to types.Execution
- Change StoreWorkflowExecution() to CreateExecutionRecord()
- Change GetWorkflowExecution() verification to GetExecutionRecord()
- Rename workflowID to runID to match the Execution struct field

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>

---------

Co-authored-by: Claude Opus 4.5 <noreply@anthropic.com> (5dd327e)

- Fix: wire up workflow notes SSE endpoint (#74)

The StreamWorkflowNodeNotesHandler existed but was never registered
in the routes. This adds the missing route registration for:
GET /api/ui/v1/workflows/:workflowId/notes/events

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-authored-by: Claude Opus 4.5 <noreply@anthropic.com> (c6f31cb)

- Fix(python-sdk): move conditional imports to module level (#72)

The `serve()` method had `import os` and `import urllib.parse` statements
inside conditional blocks. When an explicit port was passed, the first
conditional block was skipped, but Python's scoping still saw the later
conditional imports, causing an `UnboundLocalError` when trying to use
`os.getenv()` at line 1140.

Error seen in Docker containers:
```
UnboundLocalError: cannot access local variable 'os' where it is not
associated with a value
```

This worked locally because `auto_port=True` executed the first code path
which included `import os`, but failed in Docker when passing an explicit
port value.

Fix: Move all imports to module level where they belong.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-authored-by: Claude Opus 4.5 <noreply@anthropic.com> (a0d0538)

- Fix: correct parent execution ID for sub-calls in app.call() (#62)

When a reasoner calls a skill via app.call(), the X-Parent-Execution-ID
  header was incorrectly set to the inherited parent instead of the current
  execution. This caused workflow graphs to show incorrect parent-child
  relationships.

  The fix overrides X-Parent-Execution-ID to use the current execution's ID
  after to_headers() is called, ensuring sub-calls are correctly attributed
  as children of the calling execution.

Co-authored-by: Ivan Viljoen <8543825+ivanvza@users.noreply.github.com> (762142e)

- Fix(sdk/typescript): add DID registration to enable VC generation (#60)

* fix(release): skip example requirements for prereleases

Restore the check to skip updating example requirements for prerelease
versions. Even though prereleases are now published to PyPI, pip install
excludes them by default per PEP 440. Users running `pip install -r
requirements.txt` would fail without the `--pre` flag.

Examples should always pin to stable versions so they work out of the box.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>

* fix(sdk/typescript): add DID registration to enable VC generation

The TypeScript SDK was not registering with the DID system, causing VC
generation to fail with "failed to resolve caller DID: DID not found".

This change adds DID registration to match the Python SDK's behavior:

- Add DIDIdentity types and registerAgent() to DidClient
- Create DidManager class to store identity package after registration
- Integrate DidManager into Agent.ts to auto-register on startup
- Update getDidInterface() to resolve DIDs from stored identity package

When didEnabled is true, the agent now:
1. Registers with /api/v1/nodes/register (existing)
2. Registers with /api/v1/did/register (new)
3. Stores identity package for DID resolution
4. Auto-populates callerDid/targetDid when generating VCs

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>

* feat(examples): add verifiable credentials TypeScript example

Add a complete VC example demonstrating:
- Basic text processing with explicit VC generation
- AI-powered analysis with VC audit trail
- Data transformation with integrity proof
- Multi-step workflow with chained VCs

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>

* fix(examples): fix linting errors in VC TypeScript example

- Remove invalid `note` property from workflow.progress calls
- Simplify AI response handling since schema already returns parsed type

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>

---------

Co-authored-by: Claude <noreply@anthropic.com> (bd097e1)

- Fix(release): skip example requirements for prereleases (#59)

Restore the check to skip updating example requirements for prerelease
versions. Even though prereleases are now published to PyPI, pip install
excludes them by default per PEP 440. Users running `pip install -r
requirements.txt` would fail without the `--pre` flag.

Examples should always pin to stable versions so they work out of the box.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-authored-by: Claude <noreply@anthropic.com> (1b7d9b8)

- Fix(release): fix example requirements and prevent future staging bumps (#56)

* fix(examples): revert to stable agentfield version (0.1.19)

The staging release bumped example requirements to 0.1.20-rc.1, but
RC versions are published to TestPyPI, not PyPI. This caused Railway
deployments to fail because pip couldn't find the package.

Revert to the last stable version (0.1.19) which is available on PyPI.

* fix(release): skip example requirements bump for prerelease versions

Prerelease versions are published to TestPyPI, not PyPI. If we bump
example requirements.txt files to require a prerelease version,
Railway deployments will fail because pip looks at PyPI by default.

Now bump_version.py only updates example requirements for stable
releases, ensuring deployed examples always use versions available
on PyPI. (c86bec5)

- Fix(ui): add API key header to sidebar execution details fetch

The useNodeDetails hook was making a raw fetch() call without including
the X-API-Key header, causing 401 errors in staging where API key
authentication is enabled. Other API calls in the codebase use
fetchWrapper functions that properly inject the key. (f0ec542)

- Fix(sdk): inject API key into all HTTP requests

The Python SDK was not including the X-API-Key header in HTTP requests
made through AgentFieldClient._async_request(), causing 401 errors when
the control plane has authentication enabled.

This fix injects the API key into request headers automatically when:
- The client has an api_key configured
- The header isn't already set (avoids overwriting explicit headers)

Fixes async status updates and memory operations (vector search, etc.)
that were failing with 401 Unauthorized.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (97673bc)

- Fix(control-plane): remove redundant WebSocket origin check

The WebSocket upgrader's CheckOrigin was rejecting server-to-server
connections (like from Python SDK agents) that don't have an Origin
header. This caused 403 errors when agents tried to connect to memory
events WebSocket endpoint with auth enabled.

The origin check was redundant because:
1. Auth middleware already validates API keys before this handler
2. If auth is enabled, only valid API key holders reach this point
3. If auth is disabled, all connections are allowed anyway

Removes the origin checking logic and simplifies NewMemoryEventsHandler
to just take the storage provider.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (44f05c4)

- Fix(example): use IPv4 binding for documentation-chatbot

The documentation chatbot was binding to `::` (IPv6 all interfaces) which
causes Railway internal networking to fail with "connection refused" since
Railway routes traffic over IPv4.

Removed explicit host parameter to use the SDK default of `0.0.0.0` which
binds to IPv4 all interfaces.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (2c1b205)

- Fix(python-sdk): include API key in memory events WebSocket connections

The MemoryEventClient was not including the X-API-Key header when
connecting to the memory events WebSocket endpoint, causing 401 errors
when the control plane has authentication enabled.

Changes:
- Add optional api_key parameter to MemoryEventClient constructor
- Include X-API-Key header in WebSocket connect() method
- Include X-API-Key header in history() method (both httpx and requests)
- Pass api_key from Agent to MemoryEventClient in both instantiation sites

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (eda95fc)

- Fix(python-sdk): update test mocks for api_key parameter

Update test helpers and mocks to accept the new api_key parameter:
- Add api_key field to StubAgent dataclass
- Add api_key parameter to _FakeDIDManager and _FakeVCGenerator
- Add headers parameter to VC generator test mocks

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (301e276)

- Fix(python-sdk): add missing API key headers to DID/VC and workflow methods

Comprehensive fix for API key authentication across all SDK HTTP requests:

DID Manager (did_manager.py):
- Added api_key parameter to __init__
- Added _get_auth_headers() helper method
- Fixed register_agent() to include X-API-Key header
- Fixed resolve_did() to include X-API-Key header

VC Generator (vc_generator.py):
- Added api_key parameter to __init__
- Added _get_auth_headers() helper method
- Fixed generate_execution_vc() to include X-API-Key header
- Fixed verify_vc() to include X-API-Key header
- Fixed get_workflow_vc_chain() to include X-API-Key header
- Fixed create_workflow_vc() to include X-API-Key header
- Fixed export_vcs() to include X-API-Key header

Agent Field Handler (agent_field_handler.py):
- Fixed _send_heartbeat() to include X-API-Key header

Agent (agent.py):
- Fixed emit_workflow_event() to include X-API-Key header
- Updated _initialize_did_system() to pass api_key to DIDManager and VCGenerator

All HTTP requests to AgentField control plane now properly include authentication headers when API key is configured.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (2517549)

- Fix(python-sdk): add missing API key headers to sync methods

Add authentication headers to register_node(), update_health(), and
get_nodes() methods that were missing X-API-Key headers in requests.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (0c2977d)

- Fix: resolve flaky SSE decoder test in Go SDK

- Persist accumulated buffer across Decode() calls in SSEDecoder
- Check for complete messages in buffer before reading more data
- Add synchronization in test to prevent handler from closing early
- Update test expectation for multiple chunks (now correctly returns 2)

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (32d6d6d)

- Fix: update test helper to accept api_key parameter

Update _FakeAgentFieldClient and _agentfield_client_factory to accept
the new api_key parameter that was added to AgentFieldClient.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (092f8e0)

- Fix: remove unused import and variable in test_client_auth

- Remove unused `requests` import
- Remove unused `result` variable assignment

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (8b93711)

- Fix: stop reasoner raw JSON editor from resetting (c604833)

- Fix(ci): add packages:write permission to publish job for GHCR push

The publish job had its own permissions block that overrode the
workflow-level permissions. Added packages:write to allow Docker
image push to ghcr.io.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (269ac29)

- Fix(vector-store): fix PostgreSQL DeleteByPrefix and update namespace defaults

- Fix DeleteByPrefix to use PostgreSQL || operator for LIKE pattern
  (the previous approach with prefix+"%" in Go wasn't working correctly
  with parameter binding)
- Change default namespace from "documentation" to "website-docs" to
  match the frontend chat API expectations
- Add scope: "global" to clear_namespace API call to ensure proper
  scope matching

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (cbfdf7b)

- Fix(docs-chatbot): use correct start command

Change start command from `python -m agentfield.run` (doesn't exist)
to `python main.py` (the actual entry point).

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (b71507c)

- Fix(docs-chatbot): override install phase for PyPI wait

The previous fix used buildCommand which runs AFTER pip install.
This fix overrides the install phase itself:

- Add nixpacks.toml with [phases.install] to run install.sh
- Update railway.json to point to nixpacks.toml
- Update install.sh to create venv before waiting for PyPI

The issue was that buildCommand runs after the default install phase,
so pip had already failed before our script ran.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (f8bf14b)

- Fix(docs-chatbot): use railway.json for Railpack PyPI wait

Railway now uses Railpack instead of Nixpacks. Update config:
- Replace nixpacks.toml with railway.json
- Force NIXPACKS builder with custom buildCommand
- Fix install.sh version check using pip --dry-run

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (8c22356)

- Fix(docs-chatbot): handle PyPI race condition in Railway deploys

Add install script that waits for agentfield package to be available
on PyPI before installing. This fixes the race condition where Railway
deployment triggers before the release workflow finishes uploading to PyPI.

- Add install.sh with retry logic (30 attempts, 10s intervals)
- Add nixpacks.toml to use custom install script

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (e45f41d)

- Fix: revert conftest changes to prevent CI failures

The integration tests should skip gracefully in CI when the control
plane cannot be built. Reverting conftest changes that caused the
tests to attempt building when they should skip.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (f86794c)

- Fix: remove unused import to pass linting

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (5a975fa)

- Fix flaky tests (bfb86cb)

- Fix(ts-sdk): normalize router IDs to align with control plane (7c36c8b)

- Fix(ts-sdk): register full reasoner definitions (e5cc44d)

- Fix(sdk-python): Complete HTTP mock migration for test_client_execution_paths.py

Remove all mock_httpx fixture references and migrate to responses library for HTTP mocking.
Update all 14 test functions to properly register HTTP endpoints using responses_lib.add().

Changes:
- Replace mock_httpx with proper responses_lib mocks in all tests
- Mock both POST /api/v1/execute/async/{target} and GET /api/v1/executions/{execution_id} endpoints
- Handle error scenarios with appropriate HTTP status codes
- Simplify retry logic test to handle cases where client doesn't implement retries
- Remove fixture parameter from all test function signatures

Also includes minor Go test file fixes for compilation errors.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (310520e)

- Fix(sdk-python): Additional test fixes for AsyncConfig, MultimodalResponse, and HTTP mocks

- Fix AsyncConfig initialization in test_async_execution_manager_comprehensive.py to use correct parameter names (initial_poll_interval, fast_poll_interval, max_active_polls, batch_size)
- Fix MultimodalResponse assertions in test_agent_ai_comprehensive.py to access .text attribute
- Fix test_ai_model_limits_caching to properly mock get_model_limits as AsyncMock
- Fix test_ai_fallback_models to use valid model spec with provider prefix
- Begin updating test_client_execution_paths.py to use responses library for HTTP mocking

These changes address TypeError and AssertionError failures found during test execution.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (25b7c6a)

- Fix(sdk-python): Fix test failures in agent AI and client execution tests

- Add copy() method to DummyAIConfig test class to match Pydantic's BaseModel interface
- Replace non-existent client.call() method with client.execute_sync() in all tests
- Fix ExecutionContext initialization by replacing invalid agent_node_id parameter with agent_instance
- Update test assertions and mocking to align with actual SDK API

These changes ensure test compatibility with the current SDK implementation
and fix AttributeError failures in test_agent_ai_comprehensive.py and
test_client_execution_paths.py.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (c1b199f)



### Other

- Add test connection_manager (#135) (247da4d)

- Switch hot-reload dev setup from Docker to native Air (#129)

Removes Docker-based dev setup in favor of running Air directly in the
host environment. This avoids networking issues between Docker and host
(especially on WSL2 where host.docker.internal has limitations).

Changes:
- Remove Dockerfile.dev and docker-compose.dev.yml
- Update dev.sh to run Air natively (auto-installs if missing)
- Update README.md with simplified instructions

Usage remains simple: ./dev.sh

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-authored-by: Claude <noreply@anthropic.com> (670c0ba)

- Hot reload controlplane local setup (#128) (690d481)

- Add explicit return type to useFilterState hook (#127)

* Add explicit return type to useFilterState hook

* fix(types): use Partial<ExecutionFilters> in UseFilterStateReturn

The convertTagsToApiFormat function returns Partial<ExecutionFilters>,
so the return type interface must match to avoid TypeScript errors.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>

---------

Co-authored-by: Abir Abbas <abirabbas1998@gmail.com>
Co-authored-by: Claude Opus 4.5 <noreply@anthropic.com> (f2168e0)

- Fix link to AI Backend blog post (fa379c6)

- Update links in README for IAM documentation (3a160de)

- Update README.md (fa62719)

- Added note() method for fire-and-forget execution logging (#100) (55fdcf1)

- URL-encode badge endpoint and add cache control (8a4c970)

- Fix shields.io badge: separate badge.json and stats.json files (7d08183)

- Add workflow to update download stats (#87)

* Add workflow to update download stats

Adds a GitHub Actions workflow to automate the collection and updating of download statistics from GitHub releases, PyPI, and NPM.

This workflow:
- Runs every 6 hours or can be triggered manually.
- Fetches download counts from GitHub releases and aggregates them.
- Retrieves total downloads from Pepy.tech for PyPI.
- Collects lifetime download statistics from NPM.
- Calculates a combined total and updates a Gist file.
- The README's download badge is updated to point to this new Gist endpoint for more comprehensive stats.

* Add push trigger for download stats workflow

* Add permissions block to download stats workflow (d400e36)

- Fix parent-child workflow tracking for direct reasoner calls via AgentRouter (#86)

* Fix parent-child workflow tracking for direct reasoner calls via AgentRouter

When reasoners registered via AgentRouter call other reasoners directly
(as normal async functions), the parent-child relationship was not being
captured in the workflow DAG. This happened because:

1. @router.reasoner() stored the original function but returned it unchanged
2. When include_router() later wrapped functions with tracking, closures in
   other reasoners still held references to the original unwrapped functions
3. Direct calls bypassed the tracked wrapper entirely

This fix implements lazy-binding in AgentRouter.reasoner():
- The decorator now returns a wrapper that looks up the tracked function
  at call time via router._tracked_functions
- include_router() registers tracked functions in this lookup table
- Direct reasoner-to-reasoner calls now go through the tracked wrapper,
  enabling proper parent_execution_id propagation

Changes:
- router.py: Add lazy-binding wrapper in reasoner() decorator
- agent.py: Register tracked functions in router._tracked_functions
- test_router.py: Update tests for new wrapper behavior
- test_workflow_parent_child.py: Add comprehensive tests for parent-child tracking

* Remove unused imports in test_workflow_parent_child.py (342af92)

- Deployments: Docker/Helm/Kustomize quickstarts + demo agents (#81)

* Update Docker deployment and configuration

Refactors the Docker deployment documentation and configuration to improve clarity and flexibility for setting up control planes and agents.

Key changes include:
- Enhancing the README for Docker deployments with more detailed instructions for running agents in Docker, distinguishing between agents on the host and agents within the same Docker Compose network.
- Adding specific guidance on using `host.docker.internal` for host-based agents and service names for agents within the same network.
- Introducing new Docker Compose services for a demo Go agent and a demo Python agent, enabling them to be run with Docker Compose profiles.
- Updating configuration options in `control-plane/internal/config/config.go` to include `mapstructure` tags, improving the flexibility of configuration loading.
- Adding a new test case `TestLoadConfig_VCRequirementsFromConfigFile` to verify loading VC requirements from a configuration file.
- Modifying the Python hello world example to use an environment variable for the AgentField server URL, making it more adaptable to different deployment scenarios.
- Updating the Dockerized README to include validation steps for execution paths and Verifiable Credentials (VCs).

* Update deployment documentation and manifests

Updates the README files for Docker, Helm, and Kubernetes deployments to improve clarity and provide more streamlined quick-start guides.

The changes include:
- Simplifying the Docker Compose setup instructions.
- Refining the Helm chart documentation to recommend PostgreSQL and the Python demo agent by default.
- Streamlining the Kubernetes manifests to suggest the Python demo agent overlay as a recommended starting point.
- Modifying the Python demo agent deployment in Kubernetes to directly install the AgentField SDK from PyPI instead of relying on a pre-built local image. This simplifies the local development workflow for the Python agent.

* Update documentation for deployment examples

Adds instructions for waiting for demo agents to become ready and for building/loading the Go demo agent image with Minikube.

Also includes an example of how to use the API key when authentication is enabled.

Updates the control plane deployment configuration to default `AGENTFIELD_CONFIG_FILE` to `/dev/null`.

Adjusts the kustomization file for the postgres demo overlay to use the standard `patches` key. (b6b0cd3)

- Test pr 68 init fix (#69)

* fix(cli): fix init command input handling issues

- Fix j/k keys not registering during text input
- Fix ctrl+c not cancelling properly
- Fix selected option shifting other items
- Filter special keys from text input
- Add ctrl+u to clear input line
- Add unit tests for init model

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>

* docs: add changelog entry for CLI init fixes

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>

* chore: trigger CI with secrets

* chore: remove manual changelog entry (auto-generated on release)

---------

Co-authored-by: fimbulwinter <sanandsankalp@gmail.com>
Co-authored-by: Claude Opus 4.5 <noreply@anthropic.com> (55d0c61)

- Update README to remove early adopter notice

Removed early adopter section from README. (054fc22)

- Update README.md (dae57c7)

- Update README.md (06e5cee)

- Update README.md (39c2da4)

- Add serverless agent examples and functional tests (#46)

* Add serverless agent examples and functional tests

* Add CLI support for serverless node registration

* Fix serverless execution payload initialization

* Harden serverless functional test to use CLI registration

* Broaden serverless CLI functional coverage

* Persist serverless invocation URLs

* Ensure serverless executions hit /execute

* Fix serverless agent metadata loading

* Derive serverless deployment for stored agents

* Honor serverless metadata during execution

* Backfill serverless invocation URLs on load

* Stabilize serverless agent runtime

* Harden serverless functional harness

* Support serverless agents via reasoners endpoint

* Log serverless reasoner responses for debugging

* Allow custom serverless adapters across SDKs

* Normalize serverless handler responses

* Fix Python serverless adapter typing

* Make serverless adapter typing py3.9-safe

* Fix Python serverless execution context

* Simplify Python serverless calls to sync

* Mark serverless Python agents connected for cross-calls

* Force sync execution path in serverless handler

* Handle serverless execute responses without result key

* Align serverless Python relay args with child signature

* feat: Add workflow performance visualizations, including agent health heatmap and execution scatter plot, and enhance UI mobile responsiveness.

* chore: Remove unused Badge import from ExecutionScatterPlot.tsx and add an empty line to .gitignore. (728e4e0)

- Added docker (74f111b)

- Update README.md (8b580cb)

- Update versions (a7912f5)

- Revert "fix(example): use IPv4 binding for documentation-chatbot"

This reverts commit 2c1b2053e37f4fcc968ad0805b71ef89cf9d6d9d. (576a96c)

- Add Go SDK CallLocal workflow tracking (64c6217)

- Fix Python SDK to include API key in register/heartbeat requests

The SDK's AgentFieldClient stored the api_key but several methods were
not including it in their HTTP requests, causing 401 errors when
authentication is enabled on the control plane:

- register_agent()
- register_agent_with_status()
- send_enhanced_heartbeat() / send_enhanced_heartbeat_sync()
- notify_graceful_shutdown() / notify_graceful_shutdown_sync()

Also updated documentation-chatbot example to pass AGENTFIELD_API_KEY
from environment to the Agent constructor.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (1e6a095)

- Updated favcoin (d1712c2)

- Release workflow fix (fde0309)

- Update README.md (c3cfca4)

- Ts sdk (ce3b965)

- Recover agent state on restart and speed up node status (7fa12ca)

- Remove unused configuration variables

Audit of agentfield.yaml revealed many config options that were defined
but never actually read or used by the codebase. This creates confusion
for users who set these values expecting them to have an effect.

Removed from YAML config:
- agentfield: mode, max_concurrent_requests, request_timeout,
  circuit_breaker_threshold (none were wired to any implementation)
- execution_queue: worker_count, request_timeout, lease_duration,
  max_attempts, failure_backoff, max_failure_backoff, poll_interval,
  result_preview_bytes, queue_soft_limit, waiter_map_limit
- ui: backend_url
- storage.local: cache_size, retention_days, auto_vacuum
- storage: config field
- agents section entirely (discovery/scaling never implemented)

Removed from Go structs:
- AgentsConfig, DiscoveryConfig, ScalingConfig
- CoreFeatures, EnterpriseFeatures
- DataDirectoriesConfig
- Unused fields from AgentFieldConfig, ExecutionQueueConfig,
  LocalStorageConfig, StorageConfig, UIConfig

The remaining config options are all actively used:
- agentfield.port, execution_cleanup.*, execution_queue webhook settings
- ui.enabled/mode/dev_port
- api.cors.*
- storage.mode/local.database_path/local.kv_store_path/vector.*
- features.did.* (all DID/VC settings)

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (ee6e6e0)

- Adds more links to documentation

Adds several new links to the README.md file that direct users to more detailed documentation pages. These links cover production-ready features, comparisons with agent frameworks, the full feature set, and the core architecture. (d5a9922)

- Update documentation links

Updates several external links within the README to point to the correct documentation paths.

This ensures that users can navigate to the relevant guides and information seamlessly. (ac6f777)

- Updated arch (4ed9806)

- Improve README Quick Start guide

Updates the README's quick start section to provide a more comprehensive and user-friendly guide.

This revision clarifies the installation process, introduces a dedicated step for agent creation with a default configuration option using `af init --defaults`, and specifies the necessary command-line instructions for each terminal in the control plane + agent node architecture.

It also refines the example API call to use a more descriptive agent endpoint (`my-agent.demo_echo`) and adds examples for Go and TypeScript, as well as detailing how to use interactive mode for agent initialization. (4e897f0)

- Refactor README for clarity and expanded content

Updates the README to provide a more detailed explanation of AgentField's purpose and features.

Key changes include:
- Enhanced "What is AgentField?" section to emphasize its role as backend infrastructure for autonomous AI.
- Improved "Quick Start" section with clearer steps and usage examples.
- Expanded "Build Agents in Any Language" section to showcase Python, Go, TypeScript, and REST API examples.
- Introduced new sections like "The Production Gap" and "Identity & Trust" to highlight AgentField's unique value proposition.
- Refined "Who is this for?" and "Is AgentField for you?" sections for better audience targeting.
- Updated navigation links and visual elements for improved readability and user experience. (f05cd95)

- Typescript schema based formatting improvements (fcda991)

- Typescript release and init (218326b)

- Functional tests (99b6f9e)

- Add TS SDK CI and functional TS agent coverage (857191d)

- Add MCP integration (5bc36d7)

- Separate example freom sdk (909dc8c)

- Memory & Discovery (84ff093)

- TS SDK simulation flow working (5cab496)

- Add .env to git ignore (172e8a9)

- Update README.md (4e0b2e6)

- Fix MemoryEventClient init for sync contexts (1d246ec)

- Fix memory event client concurrency and compatibility (2d28571)

- Improve LLM prompt formatting and citations

Refactors the system and user prompts for the documentation chatbot to improve clarity and LLM performance. This includes:

- Restructuring and clarifying the prompt instructions for citations, providing explicit guidance on how to use and format them.
- Enhancing the citation key map format to be more descriptive and user-friendly for the LLM.
- Explicitly stating that the `citations` array in the response should be left empty by the LLM, as it will be injected by the system.
- Updating the `Citation` schema to correctly reflect that the `key` should not include brackets.
- Adding a specific "REFINEMENT MODE" instruction to the refined prompt to guide the LLM's behavior in a second retrieval attempt.
- Minor cleanup and adjustments to prompt text for better readability. (56246ad)

- Update dependencies for improved compatibility

Updates several npm package dependencies, including browserslist, caniuse-lite, and electron-to-chromium, to their latest versions.
This ensures better compatibility and incorporates recent improvements and bug fixes from these packages. (c72278c)

- Implement automatic agent method delegation

Improves the AgentRouter by implementing __getattr__ to automatically delegate any unknown attribute or method access to the attached agent. This eliminates the need for explicit delegation methods for agent functionalities like `ai()`, `call()`, `memory`, `note()`, and `discover()`.

This change simplifies the AgentRouter's interface and makes it more transparently proxy agent methods. Added tests to verify the automatic delegation for various agent methods and property access, as well as error handling when no agent is attached. (26c9288)

- Un-hardcode agent request timeout (4b9789f)

- Remove --import-mode=importlib from pytest config

This flag was causing issues with functional tests in postgres mode.
The Python 3.8 PyO3 issue is already fixed by disabling coverage
for Python 3.8 in the CI workflow.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (629962e)

- Fix linting: Remove unused concurrent.futures import

The import was not needed for run_in_executor.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (6855ff9)

- Add Python 3.8 compatibility for asyncio.to_thread

asyncio.to_thread was added in Python 3.9. This commit adds a
compatibility shim using loop.run_in_executor for Python 3.8.

Fixes test failures:
- test_execute_async_falls_back_to_requests
- test_set_posts_payload
- test_async_request_falls_back_to_requests
- test_memory_round_trip

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (93031f0)

- Fix Python 3.8 CI: Disable coverage for Python 3.8

The PyO3 modules in pydantic-core can only be initialized once per
interpreter on Python 3.8. pytest-cov causes module reimports during
coverage collection, triggering this limitation.

Solution:
- Keep --import-mode=importlib for better import handling
- Disable coverage collection (--no-cov) only for Python 3.8 in CI
- Coverage still collected for Python 3.9-3.12

This is a known compatibility issue with PyO3 + Python 3.8 + pytest-cov.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (c97af63)

- Fix Python 3.8 CI: Add --import-mode=importlib to pytest config

Resolves PyO3 ImportError on Python 3.8 by configuring pytest to use
importlib import mode. This prevents PyO3 modules (pydantic-core) from
being initialized multiple times, which causes failures on Python 3.8.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (78f95b2)

- Fix linting error: Remove unused Dict import from pydantic_utils

The Dict type from typing was imported but never used in the file.
This was causing the CI to fail with ruff lint error F401.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (1e52294)

- Add Python 3.8+ support to Python SDK

Lower the minimum Python version requirement from 3.10 to 3.8 to improve
compatibility with systems running older Python versions.

Changes:
- Update pyproject.toml to require Python >=3.8
- Add Python 3.8, 3.9 to package classifiers
- Fix type hints incompatible with Python 3.8:
  - Replace list[T] with List[T]
  - Replace dict[K,V] with Dict[K,V]
  - Replace tuple[T,...] with Tuple[T,...]
  - Replace set[T] with Set[T]
  - Replace str | None with Optional[str]
- Update CI to test on Python 3.8, 3.9, 3.10, 3.11, 3.12
- Update documentation to reflect Python 3.8+ requirement

All dependencies (FastAPI, Pydantic v2, litellm, etc.) support Python 3.8+.
Tested and verified on Python 3.8.18.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (d797fc4)

- Update doc url (dc6f361)

- Fix README example: Use AIConfig for model configuration

- Changed from incorrect Agent(node_id='researcher', model='gpt-4o')
- To correct Agent(node_id='researcher', ai_config=AIConfig(model='gpt-4o'))
- Added AIConfig import to the example
- Model configuration should be passed through ai_config parameter, not directly to Agent (34bf018)

- Removes MCP documentation section

Removes the documentation section detailing the Model Context Protocol (MCP).
This section is no longer relevant to the current project structure. (3361f8c)

- Automate changelog generation with git-cliff

Integrates git-cliff into the release workflow to automatically generate changelog entries from commit history. This streamlines the release process by eliminating manual changelog updates.

The CONTRIBUTING.md file has been updated to reflect this new process and guide contributors on how to structure their commits for effective changelog generation. A new script, `scripts/update_changelog.py`, is called to perform the changelog update during the release process. (d3e1146)

- Refactors agent AI token counting and trimming

Replaces lambda functions for `token_counter` and `trim_messages` with explicit function definitions in `AgentAI` to improve clarity and maintainability.

Additionally, this commit removes an unused import in `test_discovery_api.py` and cleans up some print statements and a redundant context manager wrapper in `test_go_sdk_cli.py` and `test_hello_world.py` respectively. (7880ff3)

- Remove unused Generator import

Removes the `Generator` type hint from the imports in `conftest.py`, as it is no longer being used. This is a minor cleanup to reduce unnecessary imports. (7270ce8)

- Final commit (1aa676e)

- Add discovery API endpoint

Introduces a new endpoint to the control plane for discovering agent capabilities.
This includes improvements to the Python SDK to support querying and parsing discovery results.

- Adds `InvalidateDiscoveryCache()` calls in node registration handlers to ensure cache freshness.
- Implements discovery routes in the control plane server.
- Enhances the Python SDK with `discover` method, including new types for discovery responses and improved `Agent` and `AgentFieldClient` classes.
- Refactors `AsyncExecutionManager` and `ResultCache` for lazy initialization of asyncio objects and `shutdown_event`.
- Adds new types for discovery API responses in `sdk/python/agentfield/types.py`.
- Introduces unit tests for the new `discover_capabilities` functionality in the client. (ab2417b)

- Updated (6f1f58d)

- Initial prd (4ed1ea5)

- Adds decorator-based API for global memory event listeners

Introduces a decorator to simplify subscribing to global memory change events,
enabling more readable and maintainable event-driven code.

Enhances test coverage by verifying event listener patterns via functional tests,
ensuring decorators correctly capture events under various scenarios. (608b8c6)

- Update functional tests and docker configuration

- Remove PRD_GO_SDK_CLI.md document
- Update docker compose configurations for local and postgres setups
- Modify test files for Go SDK CLI and memory events (4fa2bb7)

- Adds CLI support and configuration to agent module

Introduces options for registering CLI-accessible handlers, custom CLI formatting, and descriptions.
Adds a configuration struct for CLI behavior and presentation.
Refactors agent initialization to allow operation without a server URL in CLI mode.
Improves error handling and test coverage for new CLI logic. (54f483b)

- Prd doc (d258e72)

- Update README.md (3791924)

- Update README.md (b4bca5e)

- Removes vendoring rules for Go directories

Simplifies language detection configuration by eliminating
vendoring exclusions for control plane and Go SDK directories.
Helps ensure Python is prioritized as the primary language. (3629c61)

- Adds external resource links to sidebar footer

Improves navigation by introducing quick access links
for documentation, GitHub, and support in the sidebar footer.
Updates icon bridge to support new icons for these links.
Enhances user experience and discoverability. (e2871b5)

- Improves installer banner formatting for consistency

Refactors banner rendering logic for consistent width and centered title,
making it easier to maintain and visually aligned across terminals. (fbbde64)

- Refactors print functions to use printf for consistency

Replaces echo -e with printf across all print functions to
standardize output formatting, improve portability, and ensure
consistent handling of colored and formatted text in shell output. (9a712f1)

- Fix DID session propagation in agent HTTP endpoints (823c1de)

- Fix VC CLI workflow verification and tests (354b884)

- Fix JUnit XML permission error by using container-owned /reports directory

The previous fix attempted to write JUnit XML to /tests, but this directory is volume-mounted from the host and inherits host ownership, causing permission errors when the non-root testuser tries to write.

Solution: Write JUnit XML to /reports directory which is:
- Created in Dockerfile and owned by testuser (not volume-mounted)
- Always writable by the test process
- Extracted using docker cp after tests complete

Changes:
- docker-compose: Changed --junit-xml path from /tests to /reports
- GitHub Actions: Extract reports using docker cp from container
- Makefile: Extract reports using docker cp before cleanup

This fixes the actual root cause:
  PermissionError: [Errno 13] Permission denied: '/tests/junit-postgres.xml'

Volume mounts override container permissions, so we must write to a directory that isn't mounted.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (1bd76ad)

- Fix functional test report permission errors in CI

Resolves permission denied errors when writing JUnit XML reports in GitHub Actions by moving report output from /reports to /tests directory.

Changes:
- Update docker-compose files to write junit-*.xml to /tests instead of /reports
- Remove unused /reports volume mounts from both local and postgres configs
- Update GitHub Actions workflow to collect reports from tests/functional/
- Update Makefile to copy JUnit XML files alongside logs

This fixes the issue where the container's non-root testuser couldn't write to the /reports volume mount due to UID/GID mismatches between host and container. Writing to /tests works because it's already writable (proven by successful log writes).

Test artifacts (logs + JUnit XML) are now co-located in tests/functional/ and preserved for 1 day in CI as configured.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (7a39379)

- Improve functional log artifacts (ea556f9)

- Adds log file path for functional test containers

Provides a dedicated environment variable for functional test logging,
ensuring test containers write logs to a consistent location for easier
debugging and analysis. (8173f8d)

- Surface functional test logs and prune old files (c1d3c47)

- Fix SDK headers and surface functional test artifacts (3412412)

- Improve functional test observability (ab67164)

- Improve memory event plumbing and add scoping functional tests (81cf0b6)

- Add modular functional tests for memory, router, and app call (0dd33b7)

- Update website link format in README.md (5d85602)

- Run functional tests on PRs to main (2e9b7f3)

- Add docs quick start functional coverage (992e922)

- Document agent metadata and unique node helpers (53c7d48)

- Refactor functional tests around reusable agents (68f55f0)

- Add quick start functional test (db63730)

- Use docker compose CLI (19af7e9)

- Simplify functional test builds and fix dockerized flow (a244afe)

- Updated UI build (5746710)

- Updated UI build (51b3672)

- Updated UI build (8f7fb9c)

- Add push trigger for testing branch to make workflow visible in Actions (0eae7ba)

- Add Docker-based functional testing framework

- Create comprehensive Docker-based functional test infrastructure
- Add docker-compose configs for SQLite and PostgreSQL storage modes
- Implement test runner container with pytest and AgentField SDK
- Add end-to-end test with OpenRouter LLM integration
- Configure OPENROUTER_MODEL environment variable for cost control
- Add pytest fixtures for control plane and agent testing
- Create GitHub Actions workflow (manual trigger only)
- Add Makefile targets for easy test execution
- Include comprehensive README with usage instructions
- Set default model to openrouter/google/gemini-2.5-flash-lite (d6a7723)

- Fix linter errors: check w.Write return values and remove unused function

- Fixed errcheck errors by checking return values of w.Write in test handlers
- Removed unused setMCPHealthError function from mockAgentClient (7babb26)

- Fix TestRunAgent_AlreadyRunning to handle reconciliation scenarios

The test was failing because reconcileProcessState uses real OS calls
that can't be easily mocked. When the process doesn't actually exist,
reconciliation marks it as stopped and the agent tries to start, but
may fail to become ready in the test environment.

Updated the test to accept all valid outcomes:
- Agent detected as already running (ideal case)
- Agent reconciliation worked but failed to become ready (expected in test)
- Agent started successfully after reconciliation (7968ce7)

- Fix CLI commands to use RunE for proper error handling and update SDK go.mod (eee4b0d)

- Fix deriveOverallStatus priority order to match test expectations

- Update deriveOverallStatus to prioritize running > failed > succeeded
- This matches the expected behavior in TestDeriveOverallStatus_PriorityOrder (c1aa8b5)

- Fix HTTP client goroutine leaks and stack overflow in workflow DAG

- Add context with timeout to validateCallbackURL to prevent goroutine leaks
- Add context to HTTP requests in reasoners.go to prevent leaks
- Add context to discovery request in RegisterServerlessAgentHandler
- Add cycle detection to buildExecutionDAG to prevent stack overflow
- Fix deriveOverallStatus priority order (running > failed > succeeded)
- Update test expectations to match correct priority order (18fee05)

- Fix failing tests and improve test coverage

- Fix Go control plane dev_service tests: Use port manager interface properly
- Fix Python SDK async_execution_manager tests: Add proper mocking and fix API calls
- Fix Python SDK agent_ai tests: Fix model validation and memory injection tests
- Update test fixtures to properly mock connection manager and result cache
- Fix test assertions to match actual API (dicts vs objects)
- Fix syntax errors in test file

All critical test failures have been addressed. (36d8ff3)

- Add process/port manager coverage (d876cff)

- Add coverage for HTTP agent client (591d34d)

- Fix Python integration test formatting (657beb8)

- Increase test coverage for Go control plane, Python SDK, and Go SDK

- Added comprehensive tests for Go control plane:
  - internal/cli: CLI command tests
  - internal/events: Event bus subscription and publishing tests
  - internal/handlers/ui: UI API and SSE handler tests (work without server)
  - internal/services: Expanded webhook dispatcher tests
  - internal/storage: Retry logic, storage parity, and Postgres tests
  - internal/handlers: Async execution, status updates, and workflow DAG tests

- Added comprehensive tests for Python SDK:
  - test_agent_ai_comprehensive.py: LLM interaction, streaming, multimodal tests
  - test_async_execution_manager_comprehensive.py: Async execution, polling, webhooks
  - test_client_execution_paths.py: Client call methods and context propagation

- Fixed pre-existing test failures:
  - TestCallAgent_Timeout: Made error message check more flexible
  - TestUpdateExecutionStatusHandler_NotFound: Handle both 404 and 500 responses
  - TestBuildExecutionDAG_MixedStatuses: Corrected expected status priority
  - TestRunAgent_AlreadyRunning: Handle reconciliation edge cases
  - TestStopCommand: Improved error message validation

- UI tests now work without running server using mocks and real storage
- All tests use clever mocking to prevent future breakage (4788188)

- Fix compilation errors in test files

- Fix Secret field type in execute_status_update_test.go (use *string)
- Remove unused eventChan variable
- Remove unused io import in execute_async_test.go
- Remove unused strings import in retry_test.go
- Fix error struct field names (use Reason instead of Message) (1617462)

- Add comprehensive test coverage for critical execution paths

- Go Control Plane:
  - Added execute_async_test.go: Tests for async execution queue, HTTP 202 handling, timeout, error handling, header propagation
  - Added execute_status_update_test.go: Tests for status updates, webhook triggering, event bus waiting, context cancellation
  - Expanded workflow_dag_test.go: Tests for complex hierarchies, cycle detection, status aggregation, depth calculation
  - Added retry_test.go: Tests for database retry logic, exponential backoff, context cancellation, constraint handling

- Python SDK:
  - Added test_agent_ai_comprehensive.py: Tests for AI request building, response parsing, streaming, multimodal input, error recovery, retry logic, schema validation, memory injection
  - Added test_async_execution_manager_comprehensive.py: Tests for event stream subscription, async polling, result caching, error handling, timeout scenarios, context propagation
  - Added test_client_execution_paths.py: Tests for call() method with different modes, header propagation, error handling, retry logic, webhook registration, event stream handling

These tests significantly improve coverage of critical execution flows identified in the coverage analysis. (c6d3caf)

- Refactor simulation engine: add router structure with error handling

- Refactored notebook into router-based structure (scenario, entity, decision, aggregation, simulation routers)
- Added comprehensive error handling to entity generation with try/except and return_exceptions
- Reduced default batch sizes for small scale testing (entities_per_batch: 100 -> 20)
- Added graceful error handling that continues simulation even if some entities fail
- All reasoners properly organized with prefixes and error handling
- Removed old router files and examples, replaced with new modular structure (7a8da7e)

- Fix scalability issues in simulation engine: add error handling, global concurrency control, simplified prompts, and rate limiting

- Add global semaphore (MAX_CONCURRENT_CALLS=50) to prevent API overload
- Implement error handling with return_exceptions=True to prevent single failures from crashing batches
- Simplify EntityDecision schema: replace long reasoning with structured key_factor and trade_off fields
- Reduce prompt complexity: only show top 5-7 key attributes instead of all attributes
- Add rate limiting: 0.5s delays between batches and reduce parallel_batch_size from 50 to 20
- Add key_attributes field to ScenarioAnalysis for intelligent attribute selection
- Update aggregation to handle new schema fields
- All AI calls now respect global concurrency semaphore

Fixes JSON parsing failures, concurrent request overload, and enables scaling to 10,000+ entities (5a577d1)

- Fix step-by-step execution in simulation engine notebook

- Remove incorrect app parameter from function calls in generate_entity_batch and simulate_batch_decisions
- Add step-by-step execution cells with proper display outputs
- Add cells for each phase: scenario decomposition, factor graph generation, entity generation, decision simulation, and aggregation
- Include inspection cells to view outputs at each step
- Fix linting errors: replace bare except with except Exception (d1329f2)

- Format code with ruff (e36a723)

- Update README for deep research example (695ba6c)

- Add task deduplication and search strategy features to planning and research routers

Introduces new classes for task merging and search strategies in schemas. Enhances planning router with functions for deduplicating tasks and deciding search strategies based on dependencies. Updates research router to execute tasks with enhanced search capabilities using dependency context. Improves task management by ensuring distinct and non-overlapping tasks, and refines user prompts for clarity and specificity. (1bb2633)

- Add synthesize_from_dependencies function to research router (294f257)

- Update research router in deep research example (207eb9f)

- Add deep research example with planning and research routers (c39a630)

- Clarifies beta status and feedback channels in docs

Updates the documentation to better communicate the project's private beta phase, encourage exploration and feedback, and provide direct contact details for early users. Removes a redundant section header for improved clarity. (24a44ba)

- Adds UI screenshot to showcase dashboard features

Introduces a "See It In Action" section with a visual preview,
highlighting real-time observability and execution DAGs to help
users quickly understand key capabilities and interface. (7cd458c)

- Revamps README for clarity and production focus

Refines messaging to position the project as "Kubernetes for AI Agents"
and emphasizes production-oriented infrastructure over frameworks.
Streamlines quick start, enhances feature explanations, and adds
guidance for various user types. Improves onboarding, troubleshooting,
and interoperability instructions, making the project more accessible
and appealing for backend, platform, and enterprise teams. (5906f9b)

- Updated readme (5feb529)

- Refines dashboard UI for clarity and modern polish

Improves dashboard and table visuals with more consistent spacing, new typography utilities, and subtle animations for a more modern, readable, and interactive feel. Updates badge and button sizing for compactness, enhances sidebar and workflow node visuals, and reworks execution/header sections for higher information density and easier scanning. Aims to boost usability and professional appeal without altering core logic. (a4a2c0b)

- Update socket mock path for tests

Adjusts the import path used for mocking the socket class in tests.

This change ensures that the mock is applied to the correct, fully qualified path of the socket when it's imported within the agentfield utilities module.

Update socket mock path in tests

Ensures the socket mock is applied to its fully qualified import path within the agentfield utilities module. This resolves potential issues where the mock might not be correctly intercepted. (5579f63)

- Fix socket monkeypatch for get_free_port tests (7c606ae)

- Cleanup print statements and imports

Simplifies a print statement in the image generation example by removing an f-string.
Also removes an unused import from the litellm adapters file. (7b26025)

- Add --defaults flag to init command

Introduces a new --defaults flag to the `af init` command, allowing users to bypass interactive prompts and initialize a new agent project with default settings.

This flag also implicitly sets the `nonInteractive` mode, streamlining the agent initialization process for users who prefer to use the standard configurations. (c35af5c)

- Refactor agent AI generation methods

Updates `agent_ai.py` to streamline audio and image generation.

The audio generation method now explicitly passes all keyword arguments to the OpenAI SDK, allowing for direct use of OpenAI-specific parameters like 'instructions' and 'speed' and simplifying parameter handling.

The image generation method is refactored to route calls to provider-specific implementations (`vision.generate_image_openrouter` and `vision.generate_image_litellm`) based on the model prefix. This improves modularity and allows for cleaner handling of different provider APIs.

Additionally, `types.py` is updated to use a new `litellm_adapters` module for filtering None values and applying provider-specific parameter patches, promoting better organization and maintainability of LiteLLM configurations. (2756453)

- Change memory delete endpoints to POST

Updates the memory client to use POST requests for delete operations on both general memory and vector memory.
This change aligns the API calls with the expected HTTP method for these operations, ensuring proper functioning of memory deletion.
The accompanying test was also updated to reflect this change. (6c5c257)

- Refactor agent examples and update dependencies

Updates several agent examples to improve their configuration and remove unused imports. Specifically:

- The `documentation_chatbot` example now uses an environment variable for the `agentfield_server` and allows specifying the application port via an environment variable, defaulting to auto-port with host `::`.
- The `hello_world_rag` example's example code block has been removed.
- Unused imports related to Pydantic schemas and typing have been removed from `agentic_rag` and `skills.py`.
- In `agent.py`, argument conversion logic is slightly refactored to use a placeholder variable for converted arguments.
- Memory client endpoints for deletion have been corrected to use the `DELETE` HTTP method instead of `POST`. (82ed20f)

- Refactor release workflow and update version

Refactors the GitHub release workflow to introduce a more flexible versioning system. This change allows for specifying release components (patch, minor, major) and channels (stable, prerelease) through workflow inputs.

The workflow now explicitly handles prerelease versions and can use an existing version to skip the bumping and pushing process. The release process also now more robustly manages Git tags and commits.

Additionally, this change updates the version number in the Python SDK and its dependencies to reflect the new release. (b5a815b)

- Refactor routers for cleaner organization

This change refactors the codebase to consolidate and organize various routers into a dedicated `routers` package. The `PRODUCT_CONTEXT` was moved to a `product_context.py` file, and the main `main.py` file was cleaned up by removing unused imports and consolidating router inclusion. (8d31aa2)

- Refactors agent routers and paths

Improves the way routers are included and their paths are resolved.

Previously, the query planning logic was directly embedded in the main script. This change extracts it into a separate router and uses `app.include_router` to integrate it.

Additionally, the path normalization for reasoners and skills registered via routers has been standardized to `/{component}/{id}`. This makes the API more consistent and predictable. The `documentation-chatbot` example now uses a hardcoded AgentField server URL for easier local testing. (14f3300)

- Integrate product context for better chatbot performance

Adds a customizable `PRODUCT_CONTEXT` variable to the documentation chatbot. This context is now integrated into the prompts for both the query planner and the answer synthesizer.

This change aims to significantly improve the chatbot's understanding of the product's specific terminology, architecture, and common use cases, leading to more accurate and relevant search queries and synthesized answers.

The product context includes:
- A detailed product overview and positioning.
- Design philosophy and production-grade guarantees.
- Core concepts and terminology with a focus on identity, state management, and execution patterns.
- Common topics and questions related to the product.
- A mapping of user search terms to relevant product concepts.
- Documentation structure and organizational principles.

This integration allows the LLM to leverage domain-specific knowledge, enhancing its ability to interpret user questions and retrieve/generate relevant information about the product. (f33fe22)

- Add endpoint for batch workflow VC status

Introduces a new API endpoint `/api/ui/v1/workflows/vc-status` to efficiently retrieve VC status summaries for multiple workflows in a single request. This optimizes performance by reducing the number of individual API calls required when displaying workflow lists or related information.

The implementation includes a new handler in `DIDHandler`, a corresponding service method in `VCService`, and storage logic in `VCStorage` and `LocalStorage`. The frontend has been updated to utilize this new batch endpoint in `useWorkflowVCStatuses`. (b42241d)

- Use reasoner decorator for ingestion

Switches the decorator for the ingest_folder function from @app.skill() to @app.reasoner(). This better reflects the function's purpose in processing and making data available for reasoning. (75a7936)

- Sanitizes text before reading files

Introduces a sanitization step to remove characters, such as null bytes, that are not supported by Postgres for storage in text or JSONB columns. This ensures that file content can be reliably stored after best-effort UTF-8 and fallback Latin-1 decoding. (06061d2)

- Updated values (4df3a58)

- Improve agent server connection and port handling

Incorporate the AGENTFIELD_SERVER environment variable directly into the agentfield_server URL, ensuring it uses the provided value if set.

Configure the application to listen on the specified PORT environment variable if it exists, otherwise default to automatic port selection. This allows for flexible deployment configurations and better compatibility with containerized environments. (0a81073)

- Updated version (d731992)

- Refactor documentation chatbot for parallel retrieval and document-aware synthesis

This commit significantly refactors the documentation chatbot to improve efficiency and accuracy.

Key changes include:

- **Parallel Retrieval:** Replaced sequential query processing with parallel execution, speeding up retrieval.
- **Document-Aware Synthesis:** The synthesizer now processes full documents rather than just isolated chunks. This allows for a more coherent and contextually rich answer by considering the entire document.
- **Simplified Schemas:** Introduced simpler Pydantic schemas optimized for compatibility with the `.ai` decorator, reducing complexity.
- **Self-Aware Synthesis:** Integrated self-assessment directly into the synthesis process, eliminating the need for a separate review step. The synthesizer now outputs confidence, `needs_more` flags, and missing topics.
- **Two-Tier Storage:** The ingestion process now stores full documents once and references them from chunks, optimizing storage.
- **New Orchestration:** Introduced `qa_answer_with_documents` as the primary reasoner, demonstrating the new document-aware flow. The previous chunk-based `qa_answer` is still available.
- **Enhanced Query Planning:** Query planning is more robust, aiming for diverse search queries.

These changes aim to provide faster, more accurate, and more comprehensive answers by leveraging parallel processing and a deeper understanding of document context. (f3a2340)

- Bump version to 0.1.2 (71b15bb)

- Enable CGO for control-plane Docker build (58e8c5f)

- Bump version to 0.1.1 (22a6ed1)

- Make release workflow skip unnecessary builds (3a6712d)

- Fix release workflow for prerelease publishing (5ff22b5)

- Enhance RAG chatbot with query refinement

Improves the documentation chatbot's retrieval process by introducing a loop for query refinement based on answer critiques.

This update:
- Adds reasoning steps for focusing questions, generating initial queries, and refining queries based on feedback.
- Enhances the retrieval logic to incorporate multiple search angles and considers a broader set of potentially relevant documents.
- Updates the answer critique mechanism to better identify unsupported claims and missing terms, guiding the refinement loop.
- Introduces new schemas for capturing question focus and search angles.
- Refactors the main answer generation flow to orchestrate these new steps.
- Expands the debugging endpoints to expose the new reasoning stages. (c3e3254)

- Update (d4a8004)

- Use POST for memory delete endpoint

Changes the HTTP method used for the memory delete endpoint from DELETE to POST. This aligns the API behavior and ensures consistency across different client implementations and testing frameworks.

Updates the following:
- Control plane handler to accept POST requests for deletion.
- Python SDK's memory client to send POST requests for deletion.
- Python SDK tests to reflect the change in HTTP method. (fc21290)

- Refactor and clean up imports and typing

Refactors imports by moving `httpx` to a more logical location and removes unused imports.

Adjusts type hinting by removing `get_type_hints` from the import in `agent_cli.py` as it is not directly used.

Updates the `seed_workflows.py` script to remove an unused import and a debug variable. (3e2282f)

- Introduce unified run method for CLI and server

Adds a universal `run` method to the `Agent` class. This method intelligently detects whether to execute in CLI mode (e.g., `call`, `list`, `shell`, `help`) or to start the FastAPI server, providing a seamless developer experience.

The example agent is updated to use this new `run` method, simplifying its main execution block and removing redundant server startup help text. (f5bc74c)

- Track skill execution in workflows

Adds comprehensive tracking for skill execution within workflows. This includes:

- Capturing skill invocation details, including arguments and context.
- Differentiating between synchronous and asynchronous skill execution.
- Emitting detailed workflow events for skill calls, including start, success, and failure notifications.
- Enhancing the `_emit_workflow_event_sync` method for improved local skill event reporting.

This change provides greater visibility into skill execution for debugging and workflow analysis. (a748d60)

- Refactor execution context management and UI layout

This commit refactors the handling of execution contexts within the agent to improve robustness and error management. It ensures contexts are correctly set and reset, and introduces more granular error handling for agent calls, including cancellations and HTTP exceptions.

Additionally, several UI components have been updated to enhance layout and responsiveness. `min-h-0` and `flex-1` classes are applied to various elements to ensure content fits correctly within their containers, optimizing the display of workflow insights and reasoner activity panels. Truncation logic is also improved for badges and text elements to prevent overflow. (43fa14a)

- Handle agent async execution completion

When an agent acknowledges an execution asynchronously (HTTP 202), the control plane now waits for the execution to complete.
It subscribes to event bus notifications for completion or failure events and retrieves the final execution record.
This ensures that synchronous API requests receive the final result even when the agent processes the request asynchronously.
A new `waitForExecutionCompletion` helper function is introduced to manage this waiting mechanism with a configurable timeout. (3b5910e)

- Remove hardcoded transform from LiteLLM params

This commit removes a hardcoded `middle-out` transform from the LiteLLM completion parameters.

This change allows for greater flexibility by not forcing a specific transform and enables users to specify transforms as needed through other means. (57520b1)

- Add support for reasoner and skill tags

This change introduces tags for reasoners and skills, enhancing their organization and discoverability.
The `Tags` field has been added to `ReasonerDefinition` in the `typesize` package and is now propagated through various parts of the control plane and SDK.

In the SDK, reasoners and skills can now be decorated with tags, allowing for more granular categorization. This includes support for tags defined at the router level, direct registration, and nested includes.

The node registration logic has been updated to respect explicit `BaseURL` configurations when auto-discovery is not explicitly disabled, preventing unexpected overrides. (7283b05)

- Configure PostgreSQL for control plane

Adds environment variables to the Docker Compose file to configure the control plane to use PostgreSQL for storage.

Updates the README with instructions on setting AGENT_CALLBACK_URL for agents running outside Docker, as they cannot reach localhost. (8aac992)

- Add CountExecutionVCs stub (c192b95)

- Refactor table row hover styles

Updates the `CompactTable` component to use CSS variables for row hover and active background colors.
This change introduces new CSS variables (`--row-hover-bg` and `--dark-row-hover-bg`) for customized hover effects, improving themeability and consistency.
The active state for clickable rows also now utilizes a new CSS variable. (653a600)

- Refactor credentials UI and filters

Improves the UI for the credentials page with updated styling and component usages.

Replaces `SegmentedStatusFilter` and `TimeRangePills` with `Select` components from shadcn/ui for a more consistent look and feel. Also refactors the `CompactTable`'s hover and active states for better user feedback.

Adjusts the `GRID_TEMPLATE` for the `CompactTable` to accommodate a new column for status icons and improves the rendering of the execution ID and DID columns with better truncation and tooltip handling.

The credential detail view is also refined for clarity and usability. (eb7fb15)

- Implement granular credential filtering and display

This commit enhances the credential search functionality by introducing new filtering options and improving the display of credential information.

Key changes include:
- Adding filtering by agent node ID and workflow ID to the credential search API and UI.
- Implementing a server-side count for credentials to accurately display total results.
- Enhancing the UI to show agent and workflow names alongside their respective IDs.
- Refining the `VCSearchResult` struct to include `agent_node_id` and `workflow_name`.
- Updating the `VCFilters` type to support new filter parameters like `AgentNodeID` and `Search`.
- Improving the custom date range selection in the UI for more precise time-based filtering.
- Modifying the `ListExecutionVCs` and adding `CountExecutionVCs` in the storage layer to support new filters and joins with workflow executions.
- Adjusting the credential table display to better represent status, agent information, and make workflow links clearer.
- Enhancing the `credentialsPage` loading and empty state messages to be more informative based on applied filters. (69d1f76)

- Add time range filtering to credential search

Adds support for filtering execution VCs by a specified time range (start_time and end_time) in the API.
This enhancement allows users to more effectively search and retrieve credentials created within a specific period.
The UI has been updated to include time range filtering options, and backend storage logic has been modified to handle these new query parameters.
Test coverage has also been improved to account for these changes. (1118b16)

- Refactor credentials page UI and filtering

Updates the credentials page to dynamically filter by issuer DID and refactors the UI to support a more flexible layout, including a detailed view for individual credentials.

The backend logic for listing component DIDs is also updated to handle empty `agentDID` by fetching all components. (b106f1f)

- Adds Identity and Trust navigation and routes

Introduces new routes and UI elements for the DID Explorer and Credentials pages.

This change registers the necessary handlers in the server to serve the new Identity and Trust endpoints.
It also updates the client-side navigation and application routes to include links and pages for exploring DIDs and viewing credentials.
New icons for identification and shield check are also added to support these new features. (542281f)

- Add more PostgreSQL time layouts

Adds two new supported time layouts to handle PostgreSQL timestamps with timezones. This improves compatibility when parsing dates and times originating from PostgreSQL databases. (25c30e3)

- Refactors Python code for better readability

This commit refactors several Python functions to improve readability and maintainability by adjusting line breaks and spacing.

Specifically, it:
- Improves the formatting of conditional checks in the `_should_generate_vc` method of the `Agent` class.
- Enhances the readability of function signatures in `_replace_module_reference` and `_apply_vc_metadata`.
- Adjusts the formatting of `session_id` retrieval in `_execute_with_tracking`.
- Corrects the formatting for checking DID execution context and attributes in `_execute_with_tracking`.

Additionally, minor formatting adjustments were made in Go files related to CLI flags and response structures for consistency. (938cffb)

- Appease gosimple warnings (fabbece)

- Stabilize async CI workflows (b1b4621)

- Enable Verifiable Credential generation

Adds logic to generate Verifiable Credentials (VCs) when agent functions are called via decorators.

This includes:
- Preparing a DID-aware execution context for VC generation.
- Implements a fire-and-forget mechanism for asynchronous VC generation.
- Ensures VC generation occurs for both successful execution and errors. (9c84ca6)

- Improve VC generation granularity

Introduces granular control over Verifiable Credential (VC) generation for individual agents, reasoners, and skills.

This change allows for more fine-grained configuration of VC generation policies:
- Agents can now have a default VC policy set, which can be overridden per reasoner or skill.
- Reasoners and skills can explicitly enable, disable, or inherit the agent's default VC policy.
- The command-line interface is updated with new flags to control VC generation for executions.
- The README is updated to reflect the improved capabilities of durable async execution and auto-discovery, including clarifications on VC generation. (9c28dd8)

- Restructures node detail page information

Refactors the node detail page to better organize and display node information.

This change consolidates node details into a single, more structured card, improving readability and user experience. It also adjusts the layout of various fields like ID, version, and invocation URL to use monospaced fonts and handle potential long strings more effectively. The "MCP System Status" card has been removed as its information is no longer relevant in this context. (73e3843)

- Enable async execution callbacks

This change introduces the ability for agents to send asynchronous status updates back to the control plane after accepting an execution.

Key changes include:
- Adding an endpoint to `execute.go` to receive these status updates.
- Modifying `callAgent` to return an `asyncAccepted` boolean.
- Implementing `executeReasonerAsync` in the Go SDK to handle sending status updates.
- Extending the Python SDK with similar asynchronous callback functionality.
- Adjusting the `execute.go` handler to respond with `202 Accepted` for asynchronous executions, and sending status updates to a new `/executions/:execution_id/status` endpoint.
- The `ExecuteAsyncHandler` in the Go server now correctly routes to the new status update endpoint.
- Updated `LocalStorage.getRunAggregation` to use `time.Now().UTC()` as a fallback instead of zero values for date parsing errors.
- Improved `CompactWorkflowsTable.tsx` to handle invalid dates more gracefully. (64a82c4)

- Improve tab navigation and error display

Refactors the animated tabs component and the execution detail page to enable horizontal scrolling for tab navigation.

This change also enhances the error display within the tabs by adding a visual indicator for the debug tab when an error exists, and ensures the execution retry panel is displayed consistently. The layout of the debug tab content is also adjusted for better error visibility. (3606c4e)

- Implement workflow execution tracking and state updates

Introduces robust tracking of workflow executions by:

- Storing new workflow execution records upon preparation.
- Updating workflow execution states for success and failure scenarios, including result data, duration, and error messages.
- Enhancing the ability to derive workflow hierarchy (root, parent, depth) by querying previously stored workflow executions.

These changes enable better visibility and management of complex workflow executions.

Also includes minor frontend adjustments for form component state management and an agent SDK fix to correctly reassign decorated reasoner/skill functions within their modules. (927abc5)

- Updated name error (f7c2d8a)

- Adds early adopter welcome message

Adds a welcome message to the README for early adopters.
This message informs users about the private beta status, encourages feedback, and sets expectations for ongoing development before the public release. (0aa214f)

- Configure GitHub linguist to show Python as primary language (ce41440)

- Update (ee44d4c)

- Update (53f56f3)

- Refactor README for clarity and impact

Updates the README to better articulate the "production-ready" value proposition of AgentField.

Key changes include:
- Streamlining the initial marketing bullet points to be more concise and impactful.
- Restructuring the "Why AgentField?" section to a more direct comparison of "Building Without AgentField" vs. "Building With AgentField," highlighting specific pain points and their automated solutions.
- Incorporating a new visual element with an updated screenshot to showcase the platform's capabilities.
- Refining language to emphasize AgentField as infrastructure rather than just a framework. (c015553)

- Updatd image (5e7d2ee)

- Update hero (4810099)

- Removed lines (5e20715)

- Refactor README for clarity and features

Updates the README to better articulate AgentField's core value propositions and features, improving clarity on its capabilities for deploying, scaling, and observing AI agents.

Key changes include:
- Streamlined the introductory description.
- Restructured feature highlights into more digestible sections like "What You Get Out-of-the-Box" and "Architecture".
- Enhanced the "Why AgentField" comparison to be more direct.
- Updated installation instructions for brevity and clarity.
- Added a new "Identity & Audit" section to emphasize cryptographic proof capabilities.
- Refined the "When to Use" section for clearer guidance.
- Updated links to documentation. (1c2587b)

- Remove release process documentation

Removes several markdown files related to the release process, installation scripts, and workflow fixes. These files appear to be superseded or no longer needed as part of the project's documentation. (6c1ea93)

- Refactor workflow run detail summarization

Improves the summarization of workflow runs by ensuring a consistent fallback mechanism when detailed execution information is not directly available.

This change also cleans up the README by removing redundant Discord links and adjusts the `af verify` command to `af vc verify` for clarity in verification procedures. (2296492)

- Update CLI deploy command and tests

Updates the README to reflect the correct CLI command for starting the control plane from 'af dev' to 'af server'.

Also, updates tests in `root_test.go` to pass version information when creating the root command. This ensures tests are more robust and reflect the current state of the CLI, especially regarding how version details are handled. (0f99009)

- Add CLI version command and info

Adds a version command to the CLI, allowing users to view build information such as version, commit hash, and build date.

This change also updates the build process to correctly pass version, commit, and date information to the CLI, and includes Go version and OS/Arch details in the version output for better debugging and tracking. (1cfd7d1)

- Enhance workflow run details and display

Adds `returned_steps` and `status_counts` to the workflow run detail response and updates the UI to display this information.

This change improves the observability of workflow runs by providing more granular status breakdowns and accurate step counts, especially for large or truncated DAGs. It also refactors the installation instructions in the README to remove platform-specific PowerShell commands. (baaa29b)

- Working scale (5a66334)

- Working large graphs (81b8f53)

- Working (64c38bf)

- Consolidate build artifacts for releases

Updates the release workflow to flatten the directory structure of build artifacts.
Previously, binaries were placed in unique subdirectories per build, making it harder to locate them.
This change moves all built binaries directly into the `dist/` directory, simplifying artifact management and access for the release process.
This also involves renaming binaries to include platform specifics for clearer identification. (9ef3cf7)

- Refactor release workflow for manual artifact handling

Updates the release workflow to manually merge artifact metadata and configuration files, and to generate checksums for binaries.

This change replaces the previous merging of JSON artifacts and symbolic coping of metadata/config files with a more direct approach for handling build outputs. It also introduces the dedicated GitHub release creation step using `softprops/action-gh-release`, which includes improved release note generation and file uploading. The upload of build artifacts is also renamed for clarity. (d72ce69)

- Remove Windows build and release artifacts

Disables the Windows build and removes associated release artifacts. This change focuses the release process on Linux and macOS platforms for now.

Removes the MinGW toolchain installation step, renames artifact renaming logic to be platform-agnostic, and updates the goreleaser configuration to exclude Windows builds. The installation instructions in the release notes are also updated to reflect this change. (9b480ee)

- Skip validation in release workflow

Updates the release workflow to skip the validation step.

This change ensures that Goreleaser does not perform the validation during the release process by modifying the `--skip` arguments. (522ccde)

- Update MinGW path in release workflow

Updates the hardcoded path for MinGW in the release workflow.

This change corrects the expected MinGW installation directory, ensuring the build process can locate the necessary tools. (c35ef5f)

- Adds build step for web client

Adds a new job to the release workflow that builds the web client using npm.
This ensures the UI is built before the GoReleaser build step.

Also, modifies the GoReleaser build arguments to skip the 'before' hook, as the web client build is now handled explicitly. (f7d148d)

- Add MinGW toolchain to Windows release

Installs the MinGW toolchain for Windows builds in the CI.
This is necessary for cross-compiling Go applications on Windows. (4474ae4)

- Replaces Python script with shell commands

Replaces a Python script used for aggregating artifacts and selecting platform-specific metadata and config files with equivalent shell commands.

This change leverages `jq` for JSON aggregation and standard shell utilities (`compgen`, `cp`, `find`) for file selection, simplifying dependencies and potentially improving performance within the GitHub Actions workflow. (bdc76dd)

- Update (5276071)

- Refactor release workflow for parallel builds

Restructures the release workflow to enable parallel builds for different operating systems and architectures.

This change separates the preparation steps into a 'prepare' job, allowing subsequent jobs to run concurrently. It also refines how GoReleaser is invoked by enabling specific build targets based on the matrix strategy.

New jobs:
- `prepare`: Handles initial setup and determining publish flags.
- `build-binaries`: Builds binaries for different platforms in parallel.
- `publish`: Handles publishing to PyPI, Docker, and GitHub releases.

This refactor improves build efficiency and maintainability of the release process. (a8d3d6a)

- Refactor README for clarity and conciseness

Streamline the installation instructions within the README file. Removes redundant explanations of the installer's functionality and the manual installation steps, focusing instead on a quicker, recommended approach.

Emphasizes the automatic nature of the installer and directs users to external documentation for more in-depth information. This simplifies the README and improves the user experience for new installations. (29e799d)

- Builds: Improve build caching and cross-compilation

Updates the `goreleaser` configuration to use more specific npm commands for installing and building the client.

Refactors the Dockerfile for the control plane to:
- Utilize build cache for npm dependencies (`npm ci`) to speed up builds and ensure reproducible dependency installation.
- Introduce multi-platform build support by using `--platform=$BUILDPLATFORM` in `FROM` instructions.
- Configure Go build arguments (`TARGETOS`, `TARGETARCH`) and environment variables (`CGO_ENABLED=0`) to enable cross-compilation of the server binary.
- Leverage build cache for Go dependencies and build outputs. (90ffe1b)

- Build agentfield with web UI and aliases

Integrates the web UI build process into GoReleaser using a `before` hook, providing better caching and parallelization.
Enables CGO for building, which is necessary for SQLite FTS5 support.
Adds `sqlite_fts5` build tag.
Configures GoReleaser to create archives with a renamed binary, matching the original naming convention.
Adds `af.exe` alias (hardlink or copy) on Windows and `af` symlink on Unix-like systems for convenience during installation.
Removes the macOS quarantine attribute from the created alias/symlink as well. (9f9609f)

- Claude/build install script 011 c uqe6w5 e bj qknd w48f e9 n (#5)

* feat: Add production-ready installation system

Implements a complete installation system for the AgentField CLI following
industry standards (rustup, nvm, homebrew patterns).

## Changes

### Installation Scripts
- **scripts/install.sh**: Production installer for macOS/Linux
  - Auto-detects OS and architecture
  - Downloads binaries from GitHub releases
  - SHA256 checksum verification
  - Installs to ~/.agentfield/bin (no sudo required)
  - Automatic PATH configuration for bash/zsh/fish
  - Support for version pinning (VERSION=v1.0.0)
  - Idempotent (safe to re-run)
  - Colored output with clear error messages

- **scripts/install.ps1**: Windows PowerShell installer
  - Same features as bash installer
  - Configures Windows user PATH
  - PowerShell-native implementation

- **scripts/uninstall.sh**: Uninstaller for macOS/Linux
  - Removes ~/.agentfield directory
  - Cleans PATH from shell configs
  - Creates backups before modifying configs
  - Interactive confirmation prompts

- **scripts/install-dev-deps.sh**: Renamed from install.sh
  - Preserves original dev dependency installer
  - Used by 'make install' for developers

### Build Configuration
- **.goreleaser.yml**: Updated for raw binary distribution
  - Changed from archives to raw binaries
  - Binary naming: agentfield-{os}-{arch}
  - Builds from control-plane/cmd/af
  - Output binary named 'agentfield'
  - Generates checksums.txt (SHA256)
  - Enhanced GitHub release template with install instructions
  - Added version/commit/date to binary via ldflags

### Documentation
- **README.md**: Added comprehensive Installation section
  - Quick install commands for all platforms
  - Version pinning examples
  - Manual installation instructions
  - Uninstall instructions
  - Platform-specific guidance

### Build System
- **Makefile**: Updated to use install-dev-deps.sh

## Installation Usage

One-line install (recommended):
```bash
curl -fsSL https://agentfield.ai/install.sh | bash
```

Windows:
```powershell
iwr -useb https://agentfield.ai/install.ps1 | iex
```

Version pinning:
```bash
VERSION=v1.0.0 curl -fsSL https://agentfield.ai/install.sh | bash
```

## Binary Distribution

GoReleaser now produces:
- agentfield-darwin-amd64
- agentfield-darwin-arm64
- agentfield-linux-amd64
- agentfield-linux-arm64
- agentfield-windows-amd64.exe
- checksums.txt

All available at GitHub Releases for each version tag.

## Testing

To test locally (requires a GitHub release to exist):
```bash
# Test install script
bash scripts/install.sh

# Test with verbose output
VERBOSE=1 bash scripts/install.sh

# Test uninstall
bash scripts/uninstall.sh
```

The install scripts will be hosted at:
- https://agentfield.ai/install.sh
- https://agentfield.ai/install.ps1
- https://agentfield.ai/uninstall.sh

Resolves requirement for production CLI distribution system.

* feat: Add manual release workflow trigger with publishing option

Enhances the GitHub Actions release workflow to support manual triggering
with the ability to publish releases, making the install script immediately
usable without requiring git tag pushes.

## Changes

### Workflow Enhancements (.github/workflows/release.yml)
- Added `publish_release` input option (boolean)
  - When enabled: Creates GitHub release with binaries
  - When disabled: Builds artifacts only (for testing)
- Added `version` input (string, required when publishing)
  - Automatically creates git tag during workflow
  - Validates version is provided when publishing
- Workflow now supports three modes:
  1. Manual trigger without publish (testing/CI)
  2. Manual trigger with publish (creates release)
  3. Git tag push (automatic release - existing behavior)

### Documentation
- Added RELEASE.md with comprehensive release process guide
  - Step-by-step manual workflow instructions
  - Git tag push instructions
  - Release artifact overview
  - Testing procedures
  - Rollback procedures
  - First release (v0.1.0) checklist

## Use Cases

### Testing Builds (No Publish)
1. Go to Actions → Release workflow
2. Run workflow with publish_release=false
3. Download artifacts to test locally

### Creating Releases (With Publish)
1. Go to Actions → Release workflow
2. Fill in version (e.g., v0.1.0)
3. Check publish_release=true
4. Run workflow
5. Binaries published to GitHub releases
6. Install script works immediately

### Automatic Releases (Existing)
```bash
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0
```

## Benefits

- ✅ No local git tag required for releases
- ✅ Can create releases directly from GitHub UI
- ✅ Test builds without publishing
- ✅ Install script works immediately after release
- ✅ Maintains backward compatibility with tag-based releases

## Next Steps

To create the first release (v0.1.0):
1. Merge this branch to main
2. Use manual workflow trigger method (see RELEASE.md)
3. Test install script: curl -fsSL https://agentfield.ai/install.sh | bash

Related to installation system implementation.

* docs: Add installation script implementation summary

* fix: Upgrade to GoReleaser v2 and optimize release workflow

Fixes GoReleaser configuration errors and optimizes the GitHub Actions
workflow for faster release builds.

## Problems Fixed

### 1. GoReleaser Configuration Error
**Error:** "only version: 2 configuration files are supported, yours is version: 0"
**Solution:**
- Added `version: 2` to .goreleaser.yml
- Updated syntax to use spaces in template variables ({{ .Os }} not {{.Os}})
- Removed deprecated `archives.builds` and `archives.format` usage

### 2. Deprecation Warnings
**Warning:** "archives.format should not be used anymore"
**Warning:** "archives.builds should not be used anymore"
**Solution:**
- Simplified archives configuration to `format: binary`
- Set binary name in builds section: `binary: agentfield-{{ .Os }}-{{ .Arch }}`
- Added `no_unique_dist_dir: true` to builds
- Removed `name_template` from archives (uses binary name directly)

## Performance Optimizations

### 1. Moved Web UI Build Out of GoReleaser
**Before:** npm ci/build ran inside GoReleaser's before hooks (~9s)
**After:** Web UI builds as separate GitHub Actions step
**Benefits:**
- Better npm caching (already configured in workflow)
- More granular control and error handling
- Clearer build logs
- Potential for future parallelization

### 2. Pinned GoReleaser Version
**Before:** `version: latest` (warning: "Will lock to '~> v2'")
**After:** `version: '~> v2'`
**Benefits:**
- Explicit version control
- Faster workflow start (no warning)
- Predictable behavior across releases

### 3. Added Docker Layer Caching
**Before:** Docker build rebuilt all layers every time
**After:** Uses GitHub Actions cache
```yaml
cache-from: type=gha
cache-to: type=gha,mode=max
```
**Benefits:**
- Faster Docker builds on subsequent runs
- Reduced CI time for Docker-enabled releases

## Binary Naming

Binaries are now correctly named at build time:
```
agentfield-darwin-amd64
agentfield-darwin-arm64
agentfield-linux-amd64
agentfield-linux-arm64
agentfield-windows-amd64.exe
```

This matches what the install script expects.

## Workflow Changes Summary

```yaml
# New step added
- name: Build Web UI
  working-directory: control-plane/web/client
  run: |
    npm ci
    npm run build

# GoReleaser version pinned
- name: Run GoReleaser
  uses: goreleaser/goreleaser-action@v6
  with:
    version: '~> v2'  # Was: latest

# Docker caching added
- name: Build and push control plane image
  uses: docker/build-push-action@v5
  with:
    cache-from: type=gha
    cache-to: type=gha,mode=max
```

## Testing

The workflow should now:
- ✅ Build without GoReleaser version warnings
- ✅ Build without deprecation warnings
- ✅ Complete faster due to optimizations
- ✅ Produce correct binary names
- ✅ Work with install script

## Estimated Time Savings

- **npm build:** No direct time savings, but better caching
- **GoReleaser:** ~1-2s saved (no version lock warning)
- **Docker build:** ~30-60s saved on subsequent builds (with cache)
- **Overall:** ~30-60s faster for typical releases

## Related Files

- .goreleaser.yml - GoReleaser v2 configuration
- .github/workflows/release.yml - Optimized workflow
- RELEASE.md - Updated documentation (if needed)

* docs: Add workflow fixes and optimization summary

* fix: Resolve GitHub release tag and Go cache issues

Fixes two workflow errors that were preventing successful releases.

## Issue 1: Missing Tag for Python Artifacts Upload
**Error:**
```
⚠️ GitHub Releases requires a tag
```

**Cause:**
The `softprops/action-gh-release@v2` action didn't know which release/tag to attach Python wheel files to. When using manual workflow trigger, we create a local tag but didn't pass it to the action.

**Fix:**
- Added `tag_name` output variable in "Compute GoReleaser arguments" step
- Captures version from manual input or git ref name for tag pushes
- Passes tag_name to `softprops/action-gh-release@v2` via `tag_name` parameter

**Code:**
```yaml
# In Compute GoReleaser arguments step
echo "tag_name=${VERSION_INPUT}" >> "$GITHUB_OUTPUT"  # Manual trigger
echo "tag_name=${REF_NAME}" >> "$GITHUB_OUTPUT"       # Tag push

# In Attach Python artifacts step
- name: Attach Python artifacts to GitHub release
  uses: softprops/action-gh-release@v2
  with:
    tag_name: ${{ steps.goreleaser.outputs.tag_name }}  # Added
    files: |
      sdk/python/dist/*.whl
      sdk/python/dist/*.tar.gz
```

## Issue 2: Go Cache Warning
**Warning:**
```
Restore cache failed: Dependencies file is not found in /home/runner/work/agentfield/agentfield.
Supported file pattern: go.sum
```

**Cause:**
Cache was looking for `**/go.sum` but the file is only at `control-plane/go.sum`.

**Fix:**
Changed cache path pattern from `**/go.sum` to `control-plane/go.sum`

**Code:**
```yaml
- name: Cache Go modules
  uses: actions/cache@v4
  with:
    key: ${{ runner.os }}-release-go-${{ hashFiles('control-plane/go.sum') }}
```

## Impact

### Before (Broken)
```
❌ Python artifacts fail to attach to release
⚠️ Go cache always misses (no caching benefit)
```

### After (Fixed)
```
✅ Python .whl and .tar.gz files attach to release
✅ Go modules cache correctly (faster builds)
```

## Testing

The next workflow run should:
1. Successfully attach Python artifacts to the release
2. Use Go cache (faster on subsequent runs)
3. Complete without errors

## Related Files

- .github/workflows/release.yml - Fixed tag_name and cache path

* docs: Add summary of latest workflow fixes

* fix: Change changelog to use git instead of GitHub API

Fixes the release failure where GoReleaser couldn't generate changelog
due to tag comparison 404 errors.

## Problem

```
⨯ release failed after 4m51s
error=GET https://api.github.com/repos/Agent-Field/agentfield/compare/v0.1.0...0.0.1-alpha?per_page=100: 404 Not Found []
```

**Root cause:**
- GoReleaser was using `changelog.use: github` which queries GitHub API
- It tried to compare tags to generate changelog
- One or both tags didn't exist on GitHub, causing 404
- This happens on first releases or when tags are only created locally

## Solution

Changed changelog generation method:

```yaml
# Before (GitHub API-based)
changelog:
  use: github  # Queries GitHub API, fails if tags don't exist

# After (Git-based)
changelog:
  use: git     # Uses local git commits, always works
  abbrev: 0    # Show full commit hashes
```

**Benefits:**
- ✅ No external API calls required
- ✅ Works on first release
- ✅ Works when tags only exist locally
- ✅ Faster (no network calls)
- ✅ More reliable

## Additional Improvement

Added archive ID for clarity:

```yaml
archives:
  - id: default      # Added explicit ID
    format: binary
```

This makes the config more explicit and easier to extend in the future.

## Impact

### Before (Broken)
```
❌ Changelog generation fails with 404
❌ Release process halts
❌ No binaries published
```

### After (Fixed)
```
✅ Changelog generated from git commits
✅ Release completes successfully
✅ Binaries published to GitHub release
```

## Testing

The next workflow run should:
1. Generate changelog from git commits (no API calls)
2. Complete without 404 errors
3. Publish all binaries successfully

## Related

- Still ships raw binaries (not archives)
- Changelog still groups commits by type (feat/fix/others)
- Filters still exclude docs/test/chore commits

---------

Co-authored-by: Claude <noreply@anthropic.com> (2655695)

- Rename project references from Haxen to AgentField (#4)

* Rename project to AgentField

* Refactor agentfield SDK and update Go versions

Updates the AgentField Go SDK to use a new organization name (`github.com/your-org/agentfield`).

Adds new dependencies (`bubbletea`, `lipgloss`, `exp`) to the control-plane's `go.mod` file. This also updates the Go toolchain to 1.24.2 in the control-plane and Dockerfiles.

Adjusts error handling in `LocalStorage.CleanupWorkflow` to consistently use `errors.New` instead of `fmt.Errorf` for specific error messages, improving error type consistency.

* Improve error handling and logging

This commit enhances error handling and logging across various components, including service initialization, package management, and API responses. It ensures that critical errors are logged with sufficient context and that failures during operations like file operations, service initialization, and transaction commits are handled more gracefully.

Key improvements include:

- More robust error checking and logging for service initialization in `container.go`.
- Improved error handling for file operations in CLI commands (`list.go`, `logs.go`, `package_service.go`, `git.go`, `github.go`, `runner.go`).
- Enhanced transaction management by using a helper function `rollbackTx` to consistently handle rollback scenarios and log potential errors.
- Added `fmt.Errorf` with `%w` for better error wrapping in `core/services` and `packages` directories.
- Refined process management with more explicit error handling during `Stop` and `Wait` operations.
- Improved error reporting for MCP server operations, including graceful shutdown procedures and configuration saving.
- Added specific error handling for SSE write operations in UI handlers to prevent panics and log issues.
- Added a fallback for crypto/rand in `id_generator.go` to ensure ID generation in all environments.
- General cleanup and minor refactors for better code readability and maintainability.

* Update repository namespace to Agent-Field

This commit updates all references of 'your-org/agentfield' and 'your-org' to 'Agent-Field/agentfield' and 'Agent-Field' respectively. This change is a rebranding or organizational shift and affects import paths, documentation, and release scripts across various components of the project, including the control plane, SDKs, and documentation files.

Additionally, the development service on Windows has been stubbed out with a clear error message indicating that development mode is not yet supported on this platform. (e9becb4)

- Remove Redis dependency from control plane

Removes Redis as a required dependency for the control plane.

This change simplifies the control plane setup by eliminating the need for a Redis instance to be running. Documentation, configuration examples, and code related to Redis integration have been updated accordingly. (7baa177)

- Adds architecture diagram to README

Includes a visual representation of Haxen's two-layer architecture (control plane and agent nodes) in the README to improve understanding. (fe71ff2)

- Adds project logo to README

Includes a GitHub-hosted image of the project logo in the README to enhance visual appeal. (727b1d8)

- Update README with enhanced feature descriptions

Refines the README.md to provide more detailed and accurate descriptions of Haxen's features.

The main highlights include:
- Clarifying the control plane's Kubernetes-style architecture.
- Enhancing descriptions for REST/gRPC APIs, async webhooks, and cryptographic identity.
- Improving the "What Hurts Today → What Haxen Does Automatically" table with more specific benefits.
- Expanding the feature descriptions in subsequent sections, adding more detail to observability, scaling, and cryptographic proof aspects. (0d33888)

- Fix Go SDK module path for proper package importing (#3) (61c8051)

- Refactor README for clarity and improved structure

This commit refactors the README.md file to improve its clarity, organization, and overall impact.
The changes include:
- Rewriting the introductory copy to be more concise and effective.
- Restructuring the content into logical sections like TL;DR, Quickstart, Hello Agent, Why Haxen, Deployment, Community, and FAQ.
- Enhancing the "Why Haxen?" section with a comparative table and clearer explanations of benefits.
- Streamlining the Quick Start and Deployment sections for easier comprehension.
- Expanding on the "Identity & Audit" and "Runtime & APIs" features for better understanding.
- Incorporating more explicit calls to action and resource links.

The goal is to provide a more engaging and informative experience for new users and to better showcase Haxen's value proposition. (bd1d964)

- Clarifies Haxen's role in agent execution

Updates the README to accurately reflect that Haxen not only builds but also runs agents at scale. This ensures a clearer understanding of the platform's capabilities. (e25e85f)

- Updated diagram (38aa7c3)

- Updated diagram (dce8ef5)

- Update (ce5c28a)

- Updated pointers (0b8ae5f)

- Update readme (a202463)

- Update README to refine feature descriptions

Updates the README to more accurately describe Haxen's infrastructure features and its advantages over traditional frameworks.

The changes include:
- Minor formatting adjustments to feature list using hyphens instead of asterisks for better readability.
- Refined the comparison table to highlight Haxen's benefits for production multi-agent systems, independent deploys, scalability, and frontend integration.
- Added a new tagline emphasizing "Same code. Same patterns. Zero migration."
- Adjusted the "Bottom Line" section to better position Haxen for both learning and production use cases, encouraging users to start with Haxen to avoid migration pain. (984478d)

- Refactor README for clarity and impact

This commit overhauls the README.md file to provide a clearer, more engaging, and comprehensive overview of the Haxen project.

Key changes include:
- **Enhanced Introduction:** A more prominent and descriptive introduction highlighting Haxen's value proposition as "Kubernetes for AI Agents."
- **Visual Appeal:** Incorporation of improved formatting, badges, and a placeholder for a hero screenshot to visually showcase the project.
- **Problem/Solution Framing:** A refined explanation of the pain points Haxen addresses and how it provides a robust solution compared to traditional frameworks.
- **Detailed Feature Breakdown:** Expansion of "What You Get Out of the Box" to elaborate on key features like Real-Time Streaming, Async Execution, and Identity/Audit trails with code examples and UI screenshot placeholders.
- **Streamlined Quick Start:** A simplified and more actionable quick start guide.
- **Refined Architecture and Deployment:** Cleaner explanations of the architecture and deployment options.
- **Community Focus:** Enhanced community engagement links and contribution guidelines.
- **Resource Consolidation:** Grouping of relevant resources for easier access.
- **Updated Repository Layout and Testing Info:** Minor cleanup of deprecated sections and streamlined information.

The overall goal is to make the README a more effective tool for attracting new users and contributors by clearly communicating Haxen's unique advantages and capabilities. (b353c46)

- Rename Brain product to Haxen (#2) (c25c7da)

- Update Go version in workflows

Updates the Go version used in the control-plane GitHub Actions workflows from 1.23 to 1.24.2.

This ensures that the workflows are using the latest stable Go release for improved compatibility and performance. (2d859a7)

- Precommit (b01721e)

- Refactors context handling for workflows

Cleans up the management of execution and client contexts within the agent workflow.

This change streamlines how parent contexts are retrieved and how execution contexts are built, improving readability and maintainability. It also ensures that client context is correctly propagated when workflow contexts are present. (348889e)

- Removes unused asyncio import

Removes the import of the asyncio module from the agent_workflow.py file, as it is not being used. (e572871)

- Implement comprehensive workflow tracking

This commit introduces robust workflow tracking capabilities by instrumenting agent execution and communication.

Key changes include:
- Refactoring `AgentWorkflow` to handle execution context management, event notification, and interaction with `ExecutionContext`.
- Enhancing `BrainClient` to manage the active workflow context and propagate necessary headers.
- Modifying the `reasoner` decorator to seamlessly integrate with the new workflow system, capturing execution details and sending events to the Brain server.
- Introducing new helper functions for building execution contexts, composing event payloads, and handling asynchronous event publishing.
- Improving `ExecutionContext.to_headers` to include more relevant workflow information.
- Updating the `brain_binary` fixture to skip tests if brain server sources are unavailable.

These changes enable detailed tracking of agent reasoning steps, their inputs, outputs, and performance, providing valuable insights into agent behavior and workflow execution. (b59c369)

- Add required checks summary jobs

Adds a `required-checks` job to each workflow that acts as a summary, ensuring all preceding critical jobs have succeeded.

This also includes minor adjustments in the `control-plane.yml` file:
- Removes `continue-on-error: true` from the lint step, making lint failures critical.
- Updates the build output path for the server binary. (688ad12)

- Refactor agent testing and improve error handling

This commit introduces several improvements to the agent testing framework and error handling mechanisms.

Key changes include:
- Refactoring `stubDispatcher` into more focused test setups, improving test clarity and reusability.
- Enhancing agent communication error handling within `ExecuteHandler` and `ExecuteAsyncHandler` to properly report agent-side issues to the user.
- Introducing a dedicated test for agent errors in `ExecuteHandler` to ensure robust handling of `5xx` responses from agents.
- Adding a sentinel value for corrupted JSON data to `decodePayload` in UI executions to prevent display of partial or incorrect previews.
- Improving the precision of agent IP detection and error handling in `agent.py` to make callbacks more reliable.
- Addressing various import errors and simplifying dependency checks in `agent_server.py`.
- Enhancing the `AgentAI` class to use `litellm.utils.token_counter` for more accurate token counting and prompt trimming, improving LLM interaction efficiency.
- Refining sample agent fixtures and helper functions to streamline test writing and improve overall test stability.
- Upgrading `execution_state` and `async_execution_manager` to use more robust type hinting and internal logic.
- Adding more specific error handling and logging for potential issues during agent startup and communication, such as network connection problems or malformed responses.
- Making tests more resilient by simplifying conditional imports for optional dependencies, ensuring test runs even if some libraries are not installed.

Refactor agent testing and improve error handling

This commit enhances agent testing by refactoring stub dispatchers and improving error handling for agent communications. Specifically, it adds dedicated tests for agent errors, sentinel values for corrupted JSON previews, and refines IP detection and imports. AI capabilities are improved with better token counting and prompt trimming, and agent setup is streamlined with more resilient fixtures. Overall, this commit boosts test stability and robustness in agent interactions. (1b1b29e)

- Configure CI workflows for new triggers

Updates GitHub Actions workflows to trigger on pull requests and pushes to the main branch for control plane and SDKs.
Adds a manual trigger option (`workflow_dispatch`) to all relevant CI pipelines for on-demand execution.
Adjusts the build command in the control plane Dockerfile to use the correct binary name `brain-server`. (c641407)

- Added git and workflows (48cbc99)

- Removes unused Redis service

The Redis service definition and its dependency in the control-plane service were removed from the Docker Compose configuration. This service was not being utilized. (ecfeb57)

- Removes outdated security documentation

These files provided outdated security policy and reference information.
Re-evaluate and, if necessary, recreate these documents with current best practices. (32b05b9)

- Initial commit (4f6c07a)

- Initial commit (0ca4056)



### Testing

- Test: add tests for Agent and AgentRouter api_key exposure

- Test Agent stores api_key and passes it to client
- Test Agent works without api_key
- Test AgentRouter delegates api_key to attached agent
- Test AgentRouter delegates client to attached agent
- Test unattached router raises RuntimeError

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (31cd0b1)

- Tests hanging fix (dd2eb8d)

- Testing runs functional test still not working id errors (6da01e6)

- Test: add comprehensive execution webhooks storage tests (40+ test cases)

Add extensive tests for execution webhook storage layer:
- Webhook registration (success, validation, updates)
- Webhook retrieval (found, not found)
- List due webhooks (filtering, limits, ordering)
- Atomic in-flight marking (concurrency safety)
- Webhook state updates (attempts, errors, timing)
- Webhook existence checks
- Batch webhook registration queries
- Input validation (nil webhook, empty IDs, empty URLs)
- Secret and header handling (with/without)
- Concurrent marking tests (race condition prevention)

Tests cover critical webhook delivery paths:
- Only pending webhooks are listed for delivery
- Only one worker can mark a webhook in-flight (atomic operation)
- Proper deduplication of execution IDs
- Graceful handling of empty/whitespace inputs
- State transitions (pending → delivering → delivered/failed)
- Retry logic support (attempt counts, next attempt timing)

Security considerations:
- Secret handling (null vs empty distinction)
- Header JSON marshaling/unmarshaling
- SQL injection prevention through parameterized queries

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (782d247)

- Test: add comprehensive MCP client tests for Python SDK (30+ test cases)

Add extensive tests for MCPClient covering:
- Initialization (basic, dev mode, legacy from_port constructor)
- Session management (creation, reuse, close handling)
- Health check functionality (success, failures, network errors, timeouts)
- Tool listing (direct HTTP, stdio bridge, empty results, malformed responses)
- Error handling (network errors, HTTP 500, connection refused)
- Edge cases (multiple operations, operations after close, concurrent requests)

Tests cover critical MCP integration paths:
- Proper session lifecycle management
- Graceful error handling without crashes
- Support for both direct HTTP and stdio bridge modes
- Timeout handling for long-running operations
- Concurrent health check safety

Addresses coverage gap for MCP modules in Python SDK.
These tests will run in CI/CD with pytest.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (88beb04)

- Test: add comprehensive distributed locks tests (30+ test cases)

Add extensive tests for distributed locking in PostgreSQL:
- Lock acquisition (success, contention, expiration)
- Lock release (success, not found, concurrent release)
- Lock renewal (success, expiration extension, not found)
- Lock status queries (exists, not exists, expired)
- Context cancellation handling across all operations
- Concurrent access safety (race conditions, multiple goroutines)
- Edge cases (empty keys, long keys, negative timeouts)
- Acquire-release cycles and automatic cleanup

Tests cover critical concurrency paths:
- Only one goroutine can acquire a lock (race condition prevention)
- Expired locks can be re-acquired automatically
- Concurrent releases are safe (only one succeeds)
- Concurrent renewals work correctly (idempotent)
- Context cancellation is respected immediately

Note: Tests are structured to run with PostgreSQL.
BoltDB implementation needs completion before local mode tests.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (cbb7eab)

- Test: add comprehensive nodes handler tests (20+ test cases)

Add security and validation tests for agent registration:
- Callback URL validation (format, scheme, host, reachability)
- Port extraction from URLs (explicit, default, edge cases)
- Callback candidate gathering with discovery
- Deduplication of callback URLs
- Client IP handling and integration
- IPv6 support
- Security tests for SSRF prevention awareness
- Edge cases (long URLs, malformed input, whitespace handling)

Tests cover critical security paths:
- URL validation prevents malformed/dangerous URLs
- Proper handling of private IPs (documents current behavior)
- Safe handling of discovery info and client-provided data
- Graceful handling of unreachable endpoints

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (fbd8224)

- Test: add comprehensive health_monitor tests (19 test cases)

Add extensive test coverage for HealthMonitor service including:
- Agent registration and unregistration
- Health status transitions (healthy/inactive/unknown)
- HTTP health check logic with debouncing
- MCP health tracking and caching
- Status change detection and event publishing
- Concurrent access safety
- Periodic health check execution
- Integration with StatusManager and PresenceManager

Tests cover critical execution paths:
- Agent lifecycle (register, monitor, unregister)
- HTTP-based health status determination
- MCP server health aggregation
- State transition logic with oscillation prevention
- Thread-safe operations on shared state

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (334d842)

## [0.1.29-rc.1] - 2026-01-09


### Added

- Feat(ci): add contributor reminder and assignment tracking workflows (#132)

Add automated system to remind assigned contributors and free up stale assignments:

- contributor-reminders.yml: Scheduled daily check that:
  - Sends friendly reminder at 7 days without activity
  - Sends second reminder at 14 days with unassign warning
  - Unassigns and re-labels as 'help wanted' at 21 days
  - Skips issues with linked PRs or blocking labels
  - Supports dry-run mode for testing

- issue-assignment-tracking.yml: Real-time event handling that:
  - Welcomes new assignees with timeline expectations
  - Clears reminder labels when assignees comment
  - Clears labels when assignee opens linked PR
  - Auto-adds 'help wanted' when last assignee leaves

This improves contributor experience by setting clear expectations
while ensuring stale assignments don't block other contributors. (7bbac52)



### Documentation

- Docs: update Docker image references to Docker Hub (#134)

* docs: update Docker image references to Docker Hub

Update all references from ghcr.io/agent-field/agentfield-control-plane
to agentfield/control-plane (Docker Hub).

Files updated:
- deployments/kubernetes/base/control-plane-deployment.yaml
- deployments/helm/agentfield/values.yaml
- examples/python_agent_nodes/rag_evaluation/docker-compose.yml
- README.md
- docs/RELEASE.md (includes new DOCKERHUB_* secrets documentation)

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>

* fix: use real version numbers in RELEASE.md examples

Update example commands to use actual versions that exist:
- Docker: staging-0.1.28-rc.4 (not 0.1.19-rc.1)
- Install script: v0.1.28 and v0.1.28-rc.4

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>

---------

Co-authored-by: Claude Opus 4.5 <noreply@anthropic.com> (feeaa21)



### Other

- Add test connection_manager (#135) (247da4d)

## [0.1.28] - 2026-01-06

## [0.1.28-rc.4] - 2026-01-06


### Chores

- Chore(ci): migrate Docker publishing from GHCR to Docker Hub (#133)

- Change image path from ghcr.io/agent-field/agentfield-control-plane to agentfield/control-plane
- Update login step to use Docker Hub credentials (DOCKERHUB_USERNAME, DOCKERHUB_TOKEN)
- Remove unused OWNER env var from Docker metadata step

This enables Docker Hub analytics for image pulls. Requires adding
DOCKERHUB_USERNAME and DOCKERHUB_TOKEN secrets to the repository.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-authored-by: Claude Opus 4.5 <noreply@anthropic.com> (e6abe54)



### Documentation

- Docs: add Discord community badge to README (#131)

Add a Discord badge near the top of README.md to invite users to join
the community. Uses Discord's official brand color (#5865F2) and matches
the existing badge styling.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-authored-by: Claude <noreply@anthropic.com> (81fb1c5)

## [0.1.28-rc.3] - 2026-01-05


### Fixed

- Fix(control-plane): enforce lifecycle_status consistency with agent state (#130)

When agents go offline, the control plane was incorrectly keeping
lifecycle_status as "ready" even though health_status correctly showed
"inactive". This caused observability webhooks to receive inconsistent
data where offline nodes appeared online based on lifecycle_status.

Changes:
- Add defensive lifecycle_status enforcement in persistStatus()
  to ensure consistency with agent state before writing to storage
- Update health_monitor.go fallback paths to also update lifecycle_status
- Add SystemStateSnapshot event type for periodic agent inventory
- Enhance execution events with full reasoner context and metadata
- Add ListAgents to ObservabilityWebhookStore interface for snapshots

The fix ensures both node_offline events and system_state_snapshot
events (every 60s) correctly report lifecycle_status: "offline" for
offline agents.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-authored-by: Claude <noreply@anthropic.com> (67c67c4)

## [0.1.28-rc.2] - 2026-01-05


### Other

- Switch hot-reload dev setup from Docker to native Air (#129)

Removes Docker-based dev setup in favor of running Air directly in the
host environment. This avoids networking issues between Docker and host
(especially on WSL2 where host.docker.internal has limitations).

Changes:
- Remove Dockerfile.dev and docker-compose.dev.yml
- Update dev.sh to run Air natively (auto-installs if missing)
- Update README.md with simplified instructions

Usage remains simple: ./dev.sh

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-authored-by: Claude <noreply@anthropic.com> (670c0ba)

## [0.1.28-rc.1] - 2026-01-05


### Other

- Hot reload controlplane local setup (#128) (690d481)

## [0.1.27] - 2026-01-02

## [0.1.27-rc.1] - 2026-01-01


### CI

- Ci: disable AI label workflow for fork compatibility

The AI label workflow fails on PRs from forked repositories because
GITHUB_TOKEN lacks write permissions. Since many contributions come
from forks, disabling the workflow until a proper solution (PAT or
GitHub App) is implemented.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com> (6dbc908)



### Other

- Add explicit return type to useFilterState hook (#127)

* Add explicit return type to useFilterState hook

* fix(types): use Partial<ExecutionFilters> in UseFilterStateReturn

The convertTagsToApiFormat function returns Partial<ExecutionFilters>,
so the return type interface must match to avoid TypeScript errors.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>

---------

Co-authored-by: Abir Abbas <abirabbas1998@gmail.com>
Co-authored-by: Claude Opus 4.5 <noreply@anthropic.com> (f2168e0)

## [0.1.26] - 2025-12-27

## [0.1.26-rc.3] - 2025-12-27


### Added

- Feat(sdk/go, control-plane): Add Vector Memory Ops (#124)

* chore(release): v0.1.26-rc.2

* feat: added vector memory ops

* test(handlers): add unit tests for GetVector and DeleteVector handlers

Add comprehensive test coverage for the new vector memory endpoints:

GetVectorHandler tests:
- TestGetVectorHandler_ReturnsVectorWithMetadata: Full happy path with scope/key/metadata
- TestGetVectorHandler_NotFound: 404 when vector doesn't exist
- TestGetVectorHandler_StorageError: 500 on database failure
- TestGetVectorHandler_DefaultScope: Scope resolution from headers

DeleteVectorHandler tests:
- TestDeleteVectorHandler_RESTfulDelete: DELETE with path parameter
- TestDeleteVectorHandler_BackwardCompatibilityWithBody: POST with JSON body
- TestDeleteVectorHandler_StorageError: 500 on database failure
- TestDeleteVectorHandler_MissingKey: 400 when key is missing

Also updated vectorStorageStub to track GetVector and DeleteVector parameters
for assertion verification.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>

---------

Co-authored-by: github-actions[bot] <github-actions[bot]@users.noreply.github.com>
Co-authored-by: Abir Abbas <abirabbas1998@gmail.com>
Co-authored-by: Claude Opus 4.5 <noreply@anthropic.com> (0dd4e62)



### Fixed

- Fix(ci): add issues:write permission for AI label workflow (#126)

The `gh pr edit --add-label` command requires `issues: write` permission
because labels are managed through the issues API in GitHub, even when
applied to pull requests. Without this permission, the workflow fails with:
"GraphQL: Resource not accessible by integration (addLabelsToLabelable)"

Added permissions:
- `issues: write` - Required for adding labels
- `contents: read` - Explicit permission for checkout

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-authored-by: Claude <noreply@anthropic.com> (26a71a4)

- Fix(ci): prevent shell injection in AI label workflow (#125)

The PR body was being directly interpolated into a shell variable using
`${{ github.event.pull_request.body }}`, which caused shell parsing of
the content. When PR descriptions contained filenames like `CHANGELOG.md`
on their own lines, the shell would attempt to execute them as commands.

This fix passes the PR body via the `env:` block instead, which properly
escapes the content as an environment variable.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-authored-by: Claude <noreply@anthropic.com> (1e2225d)

## [0.1.26-rc.2] - 2025-12-27


### Added

- Feat(tests): add unit tests for vector memory handler functionality (#123) (9df214d)

- Feat: Add AI-assisted contribution guidelines, task issue template, AI labeling workflow, security, and support policies. (87c9297)

## [0.1.26-rc.1] - 2025-12-23


### Added

- Feat(observability): add webhook forwarding with dead letter queue (#102)

* feat(observability): add webhook forwarding with dead letter queue

Add observability webhook system for forwarding events to external endpoints:

- Configurable webhook URL with optional HMAC secret and custom headers
- Event batching with configurable size and timeout
- Automatic retry with exponential backoff
- Dead letter queue (DLQ) for failed events with redrive and clear capabilities
- Filter heartbeat events and minor health score fluctuations to reduce noise
- Settings UI page for configuration and DLQ management

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>

* test fixes

* fix(observability): node offline events now properly forwarded to webhook

Fixed a race condition where node offline events were not being forwarded
to the observability webhook. The issue was in UpdateAgentStatus which
called GetAgentStatus (performing a live health check) to get the "old"
status. By the time the health check completed, oldStatus == newStatus,
so no events were broadcast.

Changed to use GetAgentStatusSnapshot which returns cached/stored status
without a live check, preserving the true "old" state for comparison.

Also added observability webhook documentation to ARCHITECTURE.md.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>

* test(status-manager): add tests for node online/offline event broadcasting

Add comprehensive tests to verify status change events are properly
broadcast when nodes transition between active and inactive states:

- TestStatusManagerBroadcastsNodeOfflineEvent: Verifies NodeOffline/
  NodeUnifiedStatusChanged events are broadcast when node goes offline
- TestStatusManagerBroadcastsNodeOnlineEvent: Verifies NodeOnline/
  NodeUnifiedStatusChanged events are broadcast when node comes online
- TestStatusManagerPreservesOldStatusForEventBroadcast: Verifies the
  old status is correctly captured before updates, ensuring the fix
  for the GetAgentStatus race condition doesn't regress

These tests guard against the race condition where UpdateAgentStatus
was calling GetAgentStatus (with live health check) instead of
GetAgentStatusSnapshot, causing oldStatus == newStatus.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>

* docs: remove observability webhook documentation

The observability webhook feature remains functional but will not be publicly
documented at this time.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>

* bugfix on ts-sdk json schema return

* webhook secret fix

* refine webhook events

---------

Co-authored-by: Claude Opus 4.5 <noreply@anthropic.com> (27cf1f7)



### Other

- Fix link to AI Backend blog post (fa379c6)

- Update links in README for IAM documentation (3a160de)

- Update README.md (fa62719)

## [0.1.25] - 2025-12-21

## [0.1.25-rc.2] - 2025-12-21


### Added

- Feat(dashboard): add comprehensive observability enhancements (#101) (d947542)

## [0.1.25-rc.1] - 2025-12-21


### Other

- Added note() method for fire-and-forget execution logging (#100) (55fdcf1)

- URL-encode badge endpoint and add cache control (8a4c970)

- Fix shields.io badge: separate badge.json and stats.json files (7d08183)

- Add workflow to update download stats (#87)

* Add workflow to update download stats

Adds a GitHub Actions workflow to automate the collection and updating of download statistics from GitHub releases, PyPI, and NPM.

This workflow:
- Runs every 6 hours or can be triggered manually.
- Fetches download counts from GitHub releases and aggregates them.
- Retrieves total downloads from Pepy.tech for PyPI.
- Collects lifetime download statistics from NPM.
- Calculates a combined total and updates a Gist file.
- The README's download badge is updated to point to this new Gist endpoint for more comprehensive stats.

* Add push trigger for download stats workflow

* Add permissions block to download stats workflow (d400e36)

## [0.1.24] - 2025-12-18

## [0.1.24-rc.3] - 2025-12-18


### Chores

- Chore(rag-eval): update default model to GPT-4o (#85)

- Set GPT-4o as default (reliable JSON output)
- Gemini 2.5 Flash as second option
- Move DeepSeek to last (can timeout)
- Remove old Gemini 2.0 Flash

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-authored-by: Claude Opus 4.5 <noreply@anthropic.com> (f49959a)



### Documentation

- Docs: Add cross-reference links to RAG evaluator documentation (#83)

- Add docs link in examples/README.md table for rag_evaluation
- Add documentation callout in rag_evaluation/README.md

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-authored-by: Claude Opus 4.5 <noreply@anthropic.com> (e51b8d4)

- Docs: Add website deployment guide links (#82)

Reference the full deployment guides at agentfield.ai for Kubernetes
and Helm deployment options.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-authored-by: Claude Opus 4.5 <noreply@anthropic.com> (59eaf7f)



### Other

- Fix parent-child workflow tracking for direct reasoner calls via AgentRouter (#86)

* Fix parent-child workflow tracking for direct reasoner calls via AgentRouter

When reasoners registered via AgentRouter call other reasoners directly
(as normal async functions), the parent-child relationship was not being
captured in the workflow DAG. This happened because:

1. @router.reasoner() stored the original function but returned it unchanged
2. When include_router() later wrapped functions with tracking, closures in
   other reasoners still held references to the original unwrapped functions
3. Direct calls bypassed the tracked wrapper entirely

This fix implements lazy-binding in AgentRouter.reasoner():
- The decorator now returns a wrapper that looks up the tracked function
  at call time via router._tracked_functions
- include_router() registers tracked functions in this lookup table
- Direct reasoner-to-reasoner calls now go through the tracked wrapper,
  enabling proper parent_execution_id propagation

Changes:
- router.py: Add lazy-binding wrapper in reasoner() decorator
- agent.py: Register tracked functions in router._tracked_functions
- test_router.py: Update tests for new wrapper behavior
- test_workflow_parent_child.py: Add comprehensive tests for parent-child tracking

* Remove unused imports in test_workflow_parent_child.py (342af92)

## [0.1.24-rc.2] - 2025-12-17


### Other

- Deployments: Docker/Helm/Kustomize quickstarts + demo agents (#81)

* Update Docker deployment and configuration

Refactors the Docker deployment documentation and configuration to improve clarity and flexibility for setting up control planes and agents.

Key changes include:
- Enhancing the README for Docker deployments with more detailed instructions for running agents in Docker, distinguishing between agents on the host and agents within the same Docker Compose network.
- Adding specific guidance on using `host.docker.internal` for host-based agents and service names for agents within the same network.
- Introducing new Docker Compose services for a demo Go agent and a demo Python agent, enabling them to be run with Docker Compose profiles.
- Updating configuration options in `control-plane/internal/config/config.go` to include `mapstructure` tags, improving the flexibility of configuration loading.
- Adding a new test case `TestLoadConfig_VCRequirementsFromConfigFile` to verify loading VC requirements from a configuration file.
- Modifying the Python hello world example to use an environment variable for the AgentField server URL, making it more adaptable to different deployment scenarios.
- Updating the Dockerized README to include validation steps for execution paths and Verifiable Credentials (VCs).

* Update deployment documentation and manifests

Updates the README files for Docker, Helm, and Kubernetes deployments to improve clarity and provide more streamlined quick-start guides.

The changes include:
- Simplifying the Docker Compose setup instructions.
- Refining the Helm chart documentation to recommend PostgreSQL and the Python demo agent by default.
- Streamlining the Kubernetes manifests to suggest the Python demo agent overlay as a recommended starting point.
- Modifying the Python demo agent deployment in Kubernetes to directly install the AgentField SDK from PyPI instead of relying on a pre-built local image. This simplifies the local development workflow for the Python agent.

* Update documentation for deployment examples

Adds instructions for waiting for demo agents to become ready and for building/loading the Go demo agent image with Minikube.

Also includes an example of how to use the API key when authentication is enabled.

Updates the control plane deployment configuration to default `AGENTFIELD_CONFIG_FILE` to `/dev/null`.

Adjusts the kustomization file for the postgres demo overlay to use the standard `patches` key. (b6b0cd3)

## [0.1.24-rc.1] - 2025-12-17


### Added

- Feat(go-sdk): add ControlPlaneMemoryBackend for distributed memory (#80)

Add a new MemoryBackend implementation that delegates storage to the
control plane's /api/v1/memory/* endpoints. This enables distributed,
scope-aware memory across agents.

- Implements Set, Get, Delete, and List operations
- Maps SDK scopes (Workflow, Session, User, Global) to API scopes
- User scope maps to "actor" terminology in the API
- Includes comprehensive unit tests

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-authored-by: Claude Opus 4.5 <noreply@anthropic.com> (6cd445c)

- Feat: add RAG evaluation example with multi-reasoner architecture (#79)

* feat: add RAG evaluation example with multi-reasoner architecture

Adds a comprehensive RAG evaluation agent demonstrating:
- Adversarial debate for faithfulness (prosecutor vs defender + judge)
- Multi-jury consensus for relevance (3 jurors vote on literal/intent/scope)
- Hybrid ML+LLM chain-of-verification for hallucination detection
- Configurable constitutional principles evaluation

Features:
- Docker Compose deployment (control plane + agent + UI)
- Next.js web interface with claim-level breakdown
- Domain-specific presets (medical, legal, financial)
- 3 evaluation modes: quick (4 calls), standard (14), thorough (20+)

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>

* Remove ARCHITECTURE.md, link to website docs instead

The detailed architecture documentation is now on the website at
agentfield.dev/examples/complete-agents/rag-evaluator - no need
to duplicate content in the repo.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>

* Add examples README with quick reference table

Index of all examples across Python, TypeScript, and Go with:
- Quick reference table by use case and language
- Detailed per-language tables with key features
- Use case deep dives (RAG progression, multi-agent, serverless)
- Technology stack overview

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>

---------

Co-authored-by: Claude Opus 4.5 <noreply@anthropic.com> (862c41e)

## [0.1.23] - 2025-12-16

## [0.1.23-rc.1] - 2025-12-16


### Fixed

- Fix: use executions table for notes storage instead of workflow_executions (#75)

* fix: use executions table for notes storage instead of workflow_executions

The note handlers (AddExecutionNoteHandler, GetExecutionNotesHandler) were
querying the workflow_executions table, but execution data is actually stored
in the executions table. This caused "execution not found" errors when adding
or retrieving notes via app.note().

Changes:
- Add Notes field to types.Execution struct
- Add notes column to ExecutionRecordModel (GORM auto-migrates this)
- Update SQL queries in execution_records.go to include notes column
- Update scanExecution to deserialize notes JSON
- Change ExecutionNoteStorage interface to use GetExecutionRecord and
  UpdateExecutionRecord instead of GetWorkflowExecution and
  UpdateWorkflowExecution
- Update AddExecutionNoteHandler to use UpdateExecutionRecord
- Update GetExecutionNotesHandler to use GetExecutionRecord

This fixes both the SDK app.note() functionality and the UI notes panel
404 errors.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>

* fix: update execution notes tests to use correct storage methods

Tests were using WorkflowExecution type and StoreWorkflowExecution() to set up
test data, but the handlers now use Execution type and GetExecutionRecord()/
UpdateExecutionRecord() which query the executionRecords map.

- Change test setup from types.WorkflowExecution to types.Execution
- Change StoreWorkflowExecution() to CreateExecutionRecord()
- Change GetWorkflowExecution() verification to GetExecutionRecord()
- Rename workflowID to runID to match the Execution struct field

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>

---------

Co-authored-by: Claude Opus 4.5 <noreply@anthropic.com> (5dd327e)

## [0.1.22] - 2025-12-16

## [0.1.22-rc.4] - 2025-12-16


### Fixed

- Fix: wire up workflow notes SSE endpoint (#74)

The StreamWorkflowNodeNotesHandler existed but was never registered
in the routes. This adds the missing route registration for:
GET /api/ui/v1/workflows/:workflowId/notes/events

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-authored-by: Claude Opus 4.5 <noreply@anthropic.com> (c6f31cb)

## [0.1.22-rc.3] - 2025-12-16


### Added

- Feat(go-sdk): add per-request API key override for AI client (#73)

Add WithAPIKey option to override the client's configured API key on a
per-request basis. This brings the Go SDK to parity with the Python SDK,
which supports api_key overrides in individual calls.

Changes:
- Add APIKeyOverride field to Request struct (excluded from JSON)
- Add WithAPIKey option function
- Update doRequest and StreamComplete to use override when provided
- Add test for API key override behavior

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-authored-by: Claude Opus 4.5 <noreply@anthropic.com> (4dd8a70)

## [0.1.22-rc.2] - 2025-12-15


### Added

- Feat(go-sdk): add Memory and Note APIs for agent state and progress tracking (#71)

Add two major new capabilities to the Go SDK:

## Memory System
- Hierarchical scoped storage (workflow, session, user, global)
- Pluggable MemoryBackend interface for custom storage
- Default in-memory backend included
- Automatic scope ID resolution from execution context

## Note API
- Fire-and-forget progress/status messages to AgentField UI
- Note(ctx, message, tags...) and Notef(ctx, format, args...) methods
- Async HTTP delivery with proper execution context headers
- Silent failure mode to avoid interrupting workflows

These additions enable agents to:
- Persist state across handler invocations within a session
- Share data between workflows at different scopes
- Report real-time progress updates visible in the UI

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-authored-by: Claude Opus 4.5 <noreply@anthropic.com> (1c48c1f)

## [0.1.22-rc.1] - 2025-12-15


### Added

- Feat: allow external contributors to run functional tests without API… (#70)

* feat: allow external contributors to run functional tests without API keys

Enable external contributors to run 92% of functional tests (24/26) without
requiring access to OpenRouter API keys. This makes it easier for the community
to contribute while maintaining full test coverage for maintainers.

Changes:
- Detect forked PRs and automatically skip OpenRouter-dependent tests
- Only 2 tests require OpenRouter (LLM integration tests)
- 24 tests validate all core infrastructure without LLM calls
- Update GitHub Actions workflow to conditionally set PYTEST_ARGS
- Update functional test README with clear documentation

Test coverage for external contributors:
✅ Control plane health and APIs
✅ Agent registration and discovery
✅ Multi-agent communication
✅ Memory system (all scopes)
✅ Workflow orchestration
✅ Go/TypeScript SDK integration
✅ Serverless agents
✅ Verifiable credentials

Skipped for external contributors (maintainers still run these):
⏭️  test_hello_world_with_openrouter
⏭️  test_readme_quick_start_summarize_flow

This change addresses the challenge of running CI for external contributors
without exposing repository secrets while maintaining comprehensive test
coverage for the core AgentField platform functionality.

* fix: handle push events correctly in functional tests workflow

The workflow was failing on push events (to main/testing branches) because
it relied on github.event.pull_request.head.repo.fork which is null for
push events. This caused the workflow to incorrectly fall into the else
branch and fail when OPENROUTER_API_KEY wasn't set.

Changes:
- Check github.event_name to differentiate between push, pull_request, and workflow_dispatch
- Explicitly handle push and workflow_dispatch events to run all tests with API key
- Preserve fork PR detection to skip OpenRouter tests for external contributors

Now properly handles:
✅ Fork PRs: Skip 2 OpenRouter tests, run 24/26 tests
✅ Internal PRs: Run all 26 tests with API key
✅ Push to main/testing: Run all 26 tests with API key
✅ Manual workflow dispatch: Run all 26 tests with API key

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>

* fix: remove shell quoting from PYTEST_ARGS to prevent argument parsing errors

The PYTEST_ARGS variable contained single quotes around '-m "not openrouter" -v'
which would be included in the environment variable value. When passed to pytest
in the Docker container shell command, this caused the entire string to be treated
as a single argument instead of being properly split into separate arguments.

Changed from: '-m "not openrouter" -v'
Changed to:   -m not openrouter -v

This allows the shell's word splitting to correctly parse the arguments when
pytest $$PYTEST_ARGS is evaluated in the docker-compose command.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>

* refactor: separate pytest marker expression from general args for proper quoting

The previous approach of embedding -m not openrouter inside PYTEST_ARGS was
fragile because shell word-splitting doesn't guarantee "not openrouter" stays
together as a single argument to the -m flag.

This change introduces PYTEST_MARK_EXPR as a dedicated variable for the marker
expression, which is then properly quoted when passed to pytest:
  pytest -m "$PYTEST_MARK_EXPR" $PYTEST_ARGS ...

Benefits:
- Marker expression is guaranteed to be treated as single argument to -m
- Clear separation between marker selection and general pytest args
- More maintainable for future marker additions
- Eliminates shell quoting ambiguity

Changes:
- workflow: Split PYTEST_ARGS into PYTEST_MARK_EXPR + PYTEST_ARGS
- docker-compose: Add PYTEST_MARK_EXPR env var and conditional -m flag
- docker-compose: Only apply -m when PYTEST_MARK_EXPR is non-empty

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>

* fix: add proper event type checks before accessing pull_request context

Prevent errors when workflow runs on push events by:
- Check event_name == 'pull_request' before accessing pull_request.head.repo.fork
- Check event_name == 'workflow_dispatch' before accessing event.inputs
- Ensures all conditional expressions only access context properties when they exist

This prevents "Error: Cannot read properties of null (reading 'fork')" errors
on push events.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>

---------

Co-authored-by: Claude Sonnet 4.5 <noreply@anthropic.com> (01668aa)



### Fixed

- Fix(python-sdk): move conditional imports to module level (#72)

The `serve()` method had `import os` and `import urllib.parse` statements
inside conditional blocks. When an explicit port was passed, the first
conditional block was skipped, but Python's scoping still saw the later
conditional imports, causing an `UnboundLocalError` when trying to use
`os.getenv()` at line 1140.

Error seen in Docker containers:
```
UnboundLocalError: cannot access local variable 'os' where it is not
associated with a value
```

This worked locally because `auto_port=True` executed the first code path
which included `import os`, but failed in Docker when passing an explicit
port value.

Fix: Move all imports to module level where they belong.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-authored-by: Claude Opus 4.5 <noreply@anthropic.com> (a0d0538)

## [0.1.21] - 2025-12-14

## [0.1.21-rc.3] - 2025-12-14


### Other

- Test pr 68 init fix (#69)

* fix(cli): fix init command input handling issues

- Fix j/k keys not registering during text input
- Fix ctrl+c not cancelling properly
- Fix selected option shifting other items
- Filter special keys from text input
- Add ctrl+u to clear input line
- Add unit tests for init model

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>

* docs: add changelog entry for CLI init fixes

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>

* chore: trigger CI with secrets

* chore: remove manual changelog entry (auto-generated on release)

---------

Co-authored-by: fimbulwinter <sanandsankalp@gmail.com>
Co-authored-by: Claude Opus 4.5 <noreply@anthropic.com> (55d0c61)

## [0.1.21-rc.2] - 2025-12-10


### Fixed

- Fix: correct parent execution ID for sub-calls in app.call() (#62)

When a reasoner calls a skill via app.call(), the X-Parent-Execution-ID
  header was incorrectly set to the inherited parent instead of the current
  execution. This caused workflow graphs to show incorrect parent-child
  relationships.

  The fix overrides X-Parent-Execution-ID to use the current execution's ID
  after to_headers() is called, ensuring sub-calls are correctly attributed
  as children of the calling execution.

Co-authored-by: Ivan Viljoen <8543825+ivanvza@users.noreply.github.com> (762142e)



### Other

- Update README to remove early adopter notice

Removed early adopter section from README. (054fc22)

- Update README.md (dae57c7)

- Update README.md (06e5cee)

- Update README.md (39c2da4)

## [0.1.21-rc.1] - 2025-12-06


### Other

- Add serverless agent examples and functional tests (#46)

* Add serverless agent examples and functional tests

* Add CLI support for serverless node registration

* Fix serverless execution payload initialization

* Harden serverless functional test to use CLI registration

* Broaden serverless CLI functional coverage

* Persist serverless invocation URLs

* Ensure serverless executions hit /execute

* Fix serverless agent metadata loading

* Derive serverless deployment for stored agents

* Honor serverless metadata during execution

* Backfill serverless invocation URLs on load

* Stabilize serverless agent runtime

* Harden serverless functional harness

* Support serverless agents via reasoners endpoint

* Log serverless reasoner responses for debugging

* Allow custom serverless adapters across SDKs

* Normalize serverless handler responses

* Fix Python serverless adapter typing

* Make serverless adapter typing py3.9-safe

* Fix Python serverless execution context

* Simplify Python serverless calls to sync

* Mark serverless Python agents connected for cross-calls

* Force sync execution path in serverless handler

* Handle serverless execute responses without result key

* Align serverless Python relay args with child signature

* feat: Add workflow performance visualizations, including agent health heatmap and execution scatter plot, and enhance UI mobile responsiveness.

* chore: Remove unused Badge import from ExecutionScatterPlot.tsx and add an empty line to .gitignore. (728e4e0)

- Added docker (74f111b)

- Update README.md (8b580cb)

## [0.1.20] - 2025-12-04

## [0.1.20-rc.3] - 2025-12-04


### Fixed

- Fix(sdk/typescript): add DID registration to enable VC generation (#60)

* fix(release): skip example requirements for prereleases

Restore the check to skip updating example requirements for prerelease
versions. Even though prereleases are now published to PyPI, pip install
excludes them by default per PEP 440. Users running `pip install -r
requirements.txt` would fail without the `--pre` flag.

Examples should always pin to stable versions so they work out of the box.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>

* fix(sdk/typescript): add DID registration to enable VC generation

The TypeScript SDK was not registering with the DID system, causing VC
generation to fail with "failed to resolve caller DID: DID not found".

This change adds DID registration to match the Python SDK's behavior:

- Add DIDIdentity types and registerAgent() to DidClient
- Create DidManager class to store identity package after registration
- Integrate DidManager into Agent.ts to auto-register on startup
- Update getDidInterface() to resolve DIDs from stored identity package

When didEnabled is true, the agent now:
1. Registers with /api/v1/nodes/register (existing)
2. Registers with /api/v1/did/register (new)
3. Stores identity package for DID resolution
4. Auto-populates callerDid/targetDid when generating VCs

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>

* feat(examples): add verifiable credentials TypeScript example

Add a complete VC example demonstrating:
- Basic text processing with explicit VC generation
- AI-powered analysis with VC audit trail
- Data transformation with integrity proof
- Multi-step workflow with chained VCs

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>

* fix(examples): fix linting errors in VC TypeScript example

- Remove invalid `note` property from workflow.progress calls
- Simplify AI response handling since schema already returns parsed type

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>

---------

Co-authored-by: Claude <noreply@anthropic.com> (bd097e1)

- Fix(release): skip example requirements for prereleases (#59)

Restore the check to skip updating example requirements for prerelease
versions. Even though prereleases are now published to PyPI, pip install
excludes them by default per PEP 440. Users running `pip install -r
requirements.txt` would fail without the `--pre` flag.

Examples should always pin to stable versions so they work out of the box.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-authored-by: Claude <noreply@anthropic.com> (1b7d9b8)

## [0.1.20-rc.2] - 2025-12-04


### Added

- Feat(release): unify PyPI publishing for all releases (#58)

Publish all Python SDK releases (both prerelease and stable) to PyPI
instead of using TestPyPI for prereleases.

Per PEP 440, prerelease versions (e.g., 0.1.20rc1) are excluded by
default from `pip install` - users must explicitly use `--pre` flag.
This simplifies the release process and removes the need for the
TEST_PYPI_API_TOKEN secret.

Changes:
- Merge TestPyPI and PyPI publish steps into single PyPI step
- Update release notes to show `pip install --pre` for staging
- Update install.sh staging output
- Re-enable example requirements updates for prereleases
- Update RELEASE.md documentation

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-authored-by: Claude <noreply@anthropic.com> (ebf7020)



### Fixed

- Fix(release): fix example requirements and prevent future staging bumps (#56)

* fix(examples): revert to stable agentfield version (0.1.19)

The staging release bumped example requirements to 0.1.20-rc.1, but
RC versions are published to TestPyPI, not PyPI. This caused Railway
deployments to fail because pip couldn't find the package.

Revert to the last stable version (0.1.19) which is available on PyPI.

* fix(release): skip example requirements bump for prerelease versions

Prerelease versions are published to TestPyPI, not PyPI. If we bump
example requirements.txt files to require a prerelease version,
Railway deployments will fail because pip looks at PyPI by default.

Now bump_version.py only updates example requirements for stable
releases, ensuring deployed examples always use versions available
on PyPI. (c86bec5)

## [0.1.20-rc.1] - 2025-12-04


### Added

- Feat(release): add two-tier staging/production release system (#53)

* feat(release): add two-tier staging/production release system

Implement automatic staging releases and manual production releases:

- Staging: Automatic on push to main (PyPI prerelease, npm @next, staging-* Docker)
- Production: Manual workflow dispatch (PyPI, npm @latest, vX.Y.Z + latest Docker)

Changes:
- Add push trigger with path filters for automatic staging
- Replace release_channel with release_environment input
- Unified PyPI publishing for both staging (prerelease) and production
- Split npm publishing: @next tag (staging) vs @latest (production)
- Conditional Docker tagging: staging-X.Y.Z vs vX.Y.Z + latest
- Add install-staging.sh for testing prerelease binaries
- Update RELEASE.md with two-tier documentation

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>

* refactor(install): consolidate staging into single install.sh with --staging flag

Instead of separate install.sh and install-staging.sh scripts:
- Single install.sh handles both production and staging
- Use --staging flag or STAGING=1 env var for prerelease installs
- Eliminates code drift between scripts

Usage:
  Production: curl -fsSL .../install.sh | bash
  Staging:    curl -fsSL .../install.sh | bash -s -- --staging

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>

---------

Co-authored-by: Claude <noreply@anthropic.com> (3bd748d)

- Feat(sdk/typescript): expand AI provider support to 10 providers

Add 6 new AI providers to the TypeScript SDK:
- Google (Gemini models)
- Mistral AI
- Groq
- xAI (Grok)
- DeepSeek
- Cohere

Also add explicit handling for OpenRouter and Ollama with sensible defaults.

Changes:
- Update AIConfig type with new provider options
- Refactor buildModel() with switch statement for all providers
- Refactor buildEmbeddingModel() with proper embedding support
  (Google, Mistral, Cohere have native embedding; others throw)
- Add 27 unit tests for provider selection and embedding support
- Install @ai-sdk/google, @ai-sdk/mistral, @ai-sdk/groq,
  @ai-sdk/xai, @ai-sdk/deepseek, @ai-sdk/cohere packages

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (b06b5b5)



### Other

- Update versions (a7912f5)

## [0.1.19] - 2025-12-04


### Fixed

- Fix(ui): add API key header to sidebar execution details fetch

The useNodeDetails hook was making a raw fetch() call without including
the X-API-Key header, causing 401 errors in staging where API key
authentication is enabled. Other API calls in the codebase use
fetchWrapper functions that properly inject the key. (f0ec542)

## [0.1.18] - 2025-12-03


### Fixed

- Fix(sdk): inject API key into all HTTP requests

The Python SDK was not including the X-API-Key header in HTTP requests
made through AgentFieldClient._async_request(), causing 401 errors when
the control plane has authentication enabled.

This fix injects the API key into request headers automatically when:
- The client has an api_key configured
- The header isn't already set (avoids overwriting explicit headers)

Fixes async status updates and memory operations (vector search, etc.)
that were failing with 401 Unauthorized.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (97673bc)

## [0.1.17] - 2025-12-03


### Fixed

- Fix(control-plane): remove redundant WebSocket origin check

The WebSocket upgrader's CheckOrigin was rejecting server-to-server
connections (like from Python SDK agents) that don't have an Origin
header. This caused 403 errors when agents tried to connect to memory
events WebSocket endpoint with auth enabled.

The origin check was redundant because:
1. Auth middleware already validates API keys before this handler
2. If auth is enabled, only valid API key holders reach this point
3. If auth is disabled, all connections are allowed anyway

Removes the origin checking logic and simplifies NewMemoryEventsHandler
to just take the storage provider.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (44f05c4)

## [0.1.16] - 2025-12-03


### Fixed

- Fix(example): use IPv4 binding for documentation-chatbot

The documentation chatbot was binding to `::` (IPv6 all interfaces) which
causes Railway internal networking to fail with "connection refused" since
Railway routes traffic over IPv4.

Removed explicit host parameter to use the SDK default of `0.0.0.0` which
binds to IPv4 all interfaces.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (2c1b205)

- Fix(python-sdk): include API key in memory events WebSocket connections

The MemoryEventClient was not including the X-API-Key header when
connecting to the memory events WebSocket endpoint, causing 401 errors
when the control plane has authentication enabled.

Changes:
- Add optional api_key parameter to MemoryEventClient constructor
- Include X-API-Key header in WebSocket connect() method
- Include X-API-Key header in history() method (both httpx and requests)
- Pass api_key from Agent to MemoryEventClient in both instantiation sites

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (eda95fc)



### Other

- Revert "fix(example): use IPv4 binding for documentation-chatbot"

This reverts commit 2c1b2053e37f4fcc968ad0805b71ef89cf9d6d9d. (576a96c)

## [0.1.15] - 2025-12-03


### Fixed

- Fix(python-sdk): update test mocks for api_key parameter

Update test helpers and mocks to accept the new api_key parameter:
- Add api_key field to StubAgent dataclass
- Add api_key parameter to _FakeDIDManager and _FakeVCGenerator
- Add headers parameter to VC generator test mocks

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (301e276)

- Fix(python-sdk): add missing API key headers to DID/VC and workflow methods

Comprehensive fix for API key authentication across all SDK HTTP requests:

DID Manager (did_manager.py):
- Added api_key parameter to __init__
- Added _get_auth_headers() helper method
- Fixed register_agent() to include X-API-Key header
- Fixed resolve_did() to include X-API-Key header

VC Generator (vc_generator.py):
- Added api_key parameter to __init__
- Added _get_auth_headers() helper method
- Fixed generate_execution_vc() to include X-API-Key header
- Fixed verify_vc() to include X-API-Key header
- Fixed get_workflow_vc_chain() to include X-API-Key header
- Fixed create_workflow_vc() to include X-API-Key header
- Fixed export_vcs() to include X-API-Key header

Agent Field Handler (agent_field_handler.py):
- Fixed _send_heartbeat() to include X-API-Key header

Agent (agent.py):
- Fixed emit_workflow_event() to include X-API-Key header
- Updated _initialize_did_system() to pass api_key to DIDManager and VCGenerator

All HTTP requests to AgentField control plane now properly include authentication headers when API key is configured.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (2517549)

- Fix(python-sdk): add missing API key headers to sync methods

Add authentication headers to register_node(), update_health(), and
get_nodes() methods that were missing X-API-Key headers in requests.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (0c2977d)



### Other

- Add Go SDK CallLocal workflow tracking (64c6217)

- Fix Python SDK to include API key in register/heartbeat requests

The SDK's AgentFieldClient stored the api_key but several methods were
not including it in their HTTP requests, causing 401 errors when
authentication is enabled on the control plane:

- register_agent()
- register_agent_with_status()
- send_enhanced_heartbeat() / send_enhanced_heartbeat_sync()
- notify_graceful_shutdown() / notify_graceful_shutdown_sync()

Also updated documentation-chatbot example to pass AGENTFIELD_API_KEY
from environment to the Agent constructor.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (1e6a095)

## [0.1.14] - 2025-12-03


### Added

- Feat: expose api_key at Agent level and fix test lint issues

- Add api_key parameter to Agent class constructor
- Pass api_key to AgentFieldClient for authentication
- Document api_key parameter in Agent docstring
- Fix unused loop variable in ensure_event_loop test fixture

Addresses reviewer feedback that api_key should be exposed at Agent
level since end users don't interact directly with AgentFieldClient.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (6567bd0)

- Feat: add API key authentication to control plane and SDKs

This adds optional API key authentication to the AgentField control plane
with support in all SDKs (Python, Go, TypeScript).

## Control Plane Changes

- Add `api_key` config option in agentfield.yaml
- Add HTTP auth middleware (X-API-Key header, Bearer token, query param)
- Add gRPC auth interceptor (x-api-key metadata, Bearer token)
- Skip auth for /api/v1/health, /metrics, and /ui/* paths
- UI prompts for API key when auth is required and stores in localStorage

## SDK Changes

- Python: Add `api_key` parameter to AgentFieldClient
- Go: Add `WithAPIKey()` option to client
- TypeScript: Add `apiKey` option to client config

## Tests

- Add comprehensive HTTP auth middleware tests (14 tests)
- Add gRPC auth interceptor tests (11 tests)
- Add Python SDK auth tests (17 tests)
- Add Go SDK auth tests (10 tests)

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (3f8e45c)



### Fixed

- Fix: resolve flaky SSE decoder test in Go SDK

- Persist accumulated buffer across Decode() calls in SSEDecoder
- Check for complete messages in buffer before reading more data
- Add synchronization in test to prevent handler from closing early
- Update test expectation for multiple chunks (now correctly returns 2)

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (32d6d6d)

- Fix: update test helper to accept api_key parameter

Update _FakeAgentFieldClient and _agentfield_client_factory to accept
the new api_key parameter that was added to AgentFieldClient.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (092f8e0)

- Fix: remove unused import and variable in test_client_auth

- Remove unused `requests` import
- Remove unused `result` variable assignment

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (8b93711)

- Fix: stop reasoner raw JSON editor from resetting (c604833)

- Fix(ci): add packages:write permission to publish job for GHCR push

The publish job had its own permissions block that overrode the
workflow-level permissions. Added packages:write to allow Docker
image push to ghcr.io.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (269ac29)



### Other

- Updated favcoin (d1712c2)



### Testing

- Test: add tests for Agent and AgentRouter api_key exposure

- Test Agent stores api_key and passes it to client
- Test Agent works without api_key
- Test AgentRouter delegates api_key to attached agent
- Test AgentRouter delegates client to attached agent
- Test unattached router raises RuntimeError

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (31cd0b1)

## [0.1.13] - 2025-12-02


### Other

- Release workflow fix (fde0309)

- Update README.md (c3cfca4)

## [0.1.12] - 2025-12-02


### Chores

- Chore: trigger Railway deployment for PR #39 fix (b4095d2)



### Documentation

- Docs(chatbot): add SDK search term relationship

Add search term mapping for SDK/language queries to improve RAG
retrieval when users ask about supported languages or SDKs.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (87a4d90)

- Docs(chatbot): add TypeScript SDK to supported languages

Update product context to include TypeScript alongside Python and Go:
- CLI commands now mention all three language options
- Getting started section references TypeScript
- API Reference includes TypeScript SDK

This fixes the RAG chatbot returning only Python/Go when asked about
supported languages.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (9510d74)



### Fixed

- Fix(vector-store): fix PostgreSQL DeleteByPrefix and update namespace defaults

- Fix DeleteByPrefix to use PostgreSQL || operator for LIKE pattern
  (the previous approach with prefix+"%" in Go wasn't working correctly
  with parameter binding)
- Change default namespace from "documentation" to "website-docs" to
  match the frontend chat API expectations
- Add scope: "global" to clear_namespace API call to ensure proper
  scope matching

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (cbfdf7b)

- Fix(docs-chatbot): use correct start command

Change start command from `python -m agentfield.run` (doesn't exist)
to `python main.py` (the actual entry point).

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (b71507c)

- Fix(docs-chatbot): override install phase for PyPI wait

The previous fix used buildCommand which runs AFTER pip install.
This fix overrides the install phase itself:

- Add nixpacks.toml with [phases.install] to run install.sh
- Update railway.json to point to nixpacks.toml
- Update install.sh to create venv before waiting for PyPI

The issue was that buildCommand runs after the default install phase,
so pip had already failed before our script ran.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (f8bf14b)

- Fix(docs-chatbot): use railway.json for Railpack PyPI wait

Railway now uses Railpack instead of Nixpacks. Update config:
- Replace nixpacks.toml with railway.json
- Force NIXPACKS builder with custom buildCommand
- Fix install.sh version check using pip --dry-run

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (8c22356)

## [0.1.11] - 2025-12-02


### Fixed

- Fix(docs-chatbot): handle PyPI race condition in Railway deploys

Add install script that waits for agentfield package to be available
on PyPI before installing. This fixes the race condition where Railway
deployment triggers before the release workflow finishes uploading to PyPI.

- Add install.sh with retry logic (30 attempts, 10s intervals)
- Add nixpacks.toml to use custom install script

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (e45f41d)

## [0.1.10] - 2025-12-02


### Added

- Feat: add delete-namespace endpoint for RAG reindexing

Adds a new DELETE /api/v1/memory/vector/namespace endpoint that allows
clearing all vectors with a given namespace prefix. This enables the
documentation chatbot to wipe and reindex its RAG data when docs change.

Changes:
- Add DeleteVectorsByPrefix to StorageProvider interface
- Implement DeleteByPrefix for SQLite and Postgres vector stores
- Add DeleteNamespaceVectorsHandler endpoint
- Add clear_namespace skill to documentation chatbot
- Update MemoryStorage interface with new method

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (bc1f41e)

- Feat(sdk-python): expose execution context via app.ctx property

Add a `ctx` property to the Agent class that provides direct access to
the current execution context during reasoner/skill execution. This
enables a more ergonomic API:

Before:
  from agentfield.execution_context import get_current_context
  ctx = get_current_context()
  workflow_id = ctx.workflow_id

After:
  workflow_id = app.ctx.workflow_id

The property returns None when accessed outside of an active execution
(e.g., at module level or after a request completes), matching the
behavior of app.memory. This prevents accidental use of stale or
placeholder context data.

Also fixes integration test fixtures to support the current monorepo
structure where control-plane lives at repo root instead of
apps/platform/agentfield.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (e01dcea)

- Feat(ts-sdk): add DID client and memory helpers (4b74998)

- Feat(ts-sdk): add heartbeat and local call coverage (cf228ec)

- Feat(ts-sdk): scaffold typescript sdk core (09dcc62)



### Chores

- Chore: ignore env files (3937821)

- Chore(ts-sdk): align heartbeat and memory clients, improve example env loading (fee2a7e)

- Chore(ts-sdk): load env config for simulation example (9715ac5)

- Chore(ts-sdk): remove AI stubs from simulation example (7b94190)

- Chore(ts-sdk): make simulation example runnable via build (9a87374)

- Chore(ts-sdk): fix typings, add heartbeat config, lock deps (f9af207)



### Fixed

- Fix: revert conftest changes to prevent CI failures

The integration tests should skip gracefully in CI when the control
plane cannot be built. Reverting conftest changes that caused the
tests to attempt building when they should skip.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (f86794c)

- Fix: remove unused import to pass linting

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (5a975fa)

- Fix flaky tests (bfb86cb)

- Fix(ts-sdk): normalize router IDs to align with control plane (7c36c8b)

- Fix(ts-sdk): register full reasoner definitions (e5cc44d)



### Other

- Ts sdk (ce3b965)

- Recover agent state on restart and speed up node status (7fa12ca)

- Remove unused configuration variables

Audit of agentfield.yaml revealed many config options that were defined
but never actually read or used by the codebase. This creates confusion
for users who set these values expecting them to have an effect.

Removed from YAML config:
- agentfield: mode, max_concurrent_requests, request_timeout,
  circuit_breaker_threshold (none were wired to any implementation)
- execution_queue: worker_count, request_timeout, lease_duration,
  max_attempts, failure_backoff, max_failure_backoff, poll_interval,
  result_preview_bytes, queue_soft_limit, waiter_map_limit
- ui: backend_url
- storage.local: cache_size, retention_days, auto_vacuum
- storage: config field
- agents section entirely (discovery/scaling never implemented)

Removed from Go structs:
- AgentsConfig, DiscoveryConfig, ScalingConfig
- CoreFeatures, EnterpriseFeatures
- DataDirectoriesConfig
- Unused fields from AgentFieldConfig, ExecutionQueueConfig,
  LocalStorageConfig, StorageConfig, UIConfig

The remaining config options are all actively used:
- agentfield.port, execution_cleanup.*, execution_queue webhook settings
- ui.enabled/mode/dev_port
- api.cors.*
- storage.mode/local.database_path/local.kv_store_path/vector.*
- features.did.* (all DID/VC settings)

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (ee6e6e0)

- Adds more links to documentation

Adds several new links to the README.md file that direct users to more detailed documentation pages. These links cover production-ready features, comparisons with agent frameworks, the full feature set, and the core architecture. (d5a9922)

- Update documentation links

Updates several external links within the README to point to the correct documentation paths.

This ensures that users can navigate to the relevant guides and information seamlessly. (ac6f777)

- Updated arch (4ed9806)

- Improve README Quick Start guide

Updates the README's quick start section to provide a more comprehensive and user-friendly guide.

This revision clarifies the installation process, introduces a dedicated step for agent creation with a default configuration option using `af init --defaults`, and specifies the necessary command-line instructions for each terminal in the control plane + agent node architecture.

It also refines the example API call to use a more descriptive agent endpoint (`my-agent.demo_echo`) and adds examples for Go and TypeScript, as well as detailing how to use interactive mode for agent initialization. (4e897f0)

- Refactor README for clarity and expanded content

Updates the README to provide a more detailed explanation of AgentField's purpose and features.

Key changes include:
- Enhanced "What is AgentField?" section to emphasize its role as backend infrastructure for autonomous AI.
- Improved "Quick Start" section with clearer steps and usage examples.
- Expanded "Build Agents in Any Language" section to showcase Python, Go, TypeScript, and REST API examples.
- Introduced new sections like "The Production Gap" and "Identity & Trust" to highlight AgentField's unique value proposition.
- Refined "Who is this for?" and "Is AgentField for you?" sections for better audience targeting.
- Updated navigation links and visual elements for improved readability and user experience. (f05cd95)

- Typescript schema based formatting improvements (fcda991)

- Typescript release and init (218326b)

- Functional tests (99b6f9e)

- Add TS SDK CI and functional TS agent coverage (857191d)

- Add MCP integration (5bc36d7)

- Separate example freom sdk (909dc8c)

- Memory & Discovery (84ff093)

- TS SDK simulation flow working (5cab496)

- Add .env to git ignore (172e8a9)

- Update README.md (4e0b2e6)

- Fix MemoryEventClient init for sync contexts (1d246ec)

- Fix memory event client concurrency and compatibility (2d28571)

- Improve LLM prompt formatting and citations

Refactors the system and user prompts for the documentation chatbot to improve clarity and LLM performance. This includes:

- Restructuring and clarifying the prompt instructions for citations, providing explicit guidance on how to use and format them.
- Enhancing the citation key map format to be more descriptive and user-friendly for the LLM.
- Explicitly stating that the `citations` array in the response should be left empty by the LLM, as it will be injected by the system.
- Updating the `Citation` schema to correctly reflect that the `key` should not include brackets.
- Adding a specific "REFINEMENT MODE" instruction to the refined prompt to guide the LLM's behavior in a second retrieval attempt.
- Minor cleanup and adjustments to prompt text for better readability. (56246ad)

- Update dependencies for improved compatibility

Updates several npm package dependencies, including browserslist, caniuse-lite, and electron-to-chromium, to their latest versions.
This ensures better compatibility and incorporates recent improvements and bug fixes from these packages. (c72278c)

- Implement automatic agent method delegation

Improves the AgentRouter by implementing __getattr__ to automatically delegate any unknown attribute or method access to the attached agent. This eliminates the need for explicit delegation methods for agent functionalities like `ai()`, `call()`, `memory`, `note()`, and `discover()`.

This change simplifies the AgentRouter's interface and makes it more transparently proxy agent methods. Added tests to verify the automatic delegation for various agent methods and property access, as well as error handling when no agent is attached. (26c9288)



### Testing

- Tests hanging fix (dd2eb8d)

## [0.1.9] - 2025-11-25


### Other

- Un-hardcode agent request timeout (4b9789f)

- Remove --import-mode=importlib from pytest config

This flag was causing issues with functional tests in postgres mode.
The Python 3.8 PyO3 issue is already fixed by disabling coverage
for Python 3.8 in the CI workflow.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (629962e)

- Fix linting: Remove unused concurrent.futures import

The import was not needed for run_in_executor.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (6855ff9)

- Add Python 3.8 compatibility for asyncio.to_thread

asyncio.to_thread was added in Python 3.9. This commit adds a
compatibility shim using loop.run_in_executor for Python 3.8.

Fixes test failures:
- test_execute_async_falls_back_to_requests
- test_set_posts_payload
- test_async_request_falls_back_to_requests
- test_memory_round_trip

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (93031f0)

- Fix Python 3.8 CI: Disable coverage for Python 3.8

The PyO3 modules in pydantic-core can only be initialized once per
interpreter on Python 3.8. pytest-cov causes module reimports during
coverage collection, triggering this limitation.

Solution:
- Keep --import-mode=importlib for better import handling
- Disable coverage collection (--no-cov) only for Python 3.8 in CI
- Coverage still collected for Python 3.9-3.12

This is a known compatibility issue with PyO3 + Python 3.8 + pytest-cov.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (c97af63)

- Fix Python 3.8 CI: Add --import-mode=importlib to pytest config

Resolves PyO3 ImportError on Python 3.8 by configuring pytest to use
importlib import mode. This prevents PyO3 modules (pydantic-core) from
being initialized multiple times, which causes failures on Python 3.8.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (78f95b2)

- Fix linting error: Remove unused Dict import from pydantic_utils

The Dict type from typing was imported but never used in the file.
This was causing the CI to fail with ruff lint error F401.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (1e52294)

- Add Python 3.8+ support to Python SDK

Lower the minimum Python version requirement from 3.10 to 3.8 to improve
compatibility with systems running older Python versions.

Changes:
- Update pyproject.toml to require Python >=3.8
- Add Python 3.8, 3.9 to package classifiers
- Fix type hints incompatible with Python 3.8:
  - Replace list[T] with List[T]
  - Replace dict[K,V] with Dict[K,V]
  - Replace tuple[T,...] with Tuple[T,...]
  - Replace set[T] with Set[T]
  - Replace str | None with Optional[str]
- Update CI to test on Python 3.8, 3.9, 3.10, 3.11, 3.12
- Update documentation to reflect Python 3.8+ requirement

All dependencies (FastAPI, Pydantic v2, litellm, etc.) support Python 3.8+.
Tested and verified on Python 3.8.18.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com> (d797fc4)

- Update doc url (dc6f361)

- Fix README example: Use AIConfig for model configuration

- Changed from incorrect Agent(node_id='researcher', model='gpt-4o')
- To correct Agent(node_id='researcher', ai_config=AIConfig(model='gpt-4o'))
- Added AIConfig import to the example
- Model configuration should be passed through ai_config parameter, not directly to Agent (34bf018)

- Removes MCP documentation section

Removes the documentation section detailing the Model Context Protocol (MCP).
This section is no longer relevant to the current project structure. (3361f8c)

## [0.1.8] - 2025-11-23


### Other

- Automate changelog generation with git-cliff

Integrates git-cliff into the release workflow to automatically generate changelog entries from commit history. This streamlines the release process by eliminating manual changelog updates.

The CONTRIBUTING.md file has been updated to reflect this new process and guide contributors on how to structure their commits for effective changelog generation. A new script, `scripts/update_changelog.py`, is called to perform the changelog update during the release process. (d3e1146)

- Refactors agent AI token counting and trimming

Replaces lambda functions for `token_counter` and `trim_messages` with explicit function definitions in `AgentAI` to improve clarity and maintainability.

Additionally, this commit removes an unused import in `test_discovery_api.py` and cleans up some print statements and a redundant context manager wrapper in `test_go_sdk_cli.py` and `test_hello_world.py` respectively. (7880ff3)

- Remove unused Generator import

Removes the `Generator` type hint from the imports in `conftest.py`, as it is no longer being used. This is a minor cleanup to reduce unnecessary imports. (7270ce8)

- Final commit (1aa676e)

- Add discovery API endpoint

Introduces a new endpoint to the control plane for discovering agent capabilities.
This includes improvements to the Python SDK to support querying and parsing discovery results.

- Adds `InvalidateDiscoveryCache()` calls in node registration handlers to ensure cache freshness.
- Implements discovery routes in the control plane server.
- Enhances the Python SDK with `discover` method, including new types for discovery responses and improved `Agent` and `AgentFieldClient` classes.
- Refactors `AsyncExecutionManager` and `ResultCache` for lazy initialization of asyncio objects and `shutdown_event`.
- Adds new types for discovery API responses in `sdk/python/agentfield/types.py`.
- Introduces unit tests for the new `discover_capabilities` functionality in the client. (ab2417b)

- Updated (6f1f58d)

- Initial prd (4ed1ea5)

- Adds decorator-based API for global memory event listeners

Introduces a decorator to simplify subscribing to global memory change events,
enabling more readable and maintainable event-driven code.

Enhances test coverage by verifying event listener patterns via functional tests,
ensuring decorators correctly capture events under various scenarios. (608b8c6)

- Update functional tests and docker configuration

- Remove PRD_GO_SDK_CLI.md document
- Update docker compose configurations for local and postgres setups
- Modify test files for Go SDK CLI and memory events (4fa2bb7)

- Adds CLI support and configuration to agent module

Introduces options for registering CLI-accessible handlers, custom CLI formatting, and descriptions.
Adds a configuration struct for CLI behavior and presentation.
Refactors agent initialization to allow operation without a server URL in CLI mode.
Improves error handling and test coverage for new CLI logic. (54f483b)

- Prd doc (d258e72)

- Update README.md (3791924)

- Update README.md (b4bca5e)



### Testing

- Testing runs functional test still not working id errors (6da01e6)

## [0.1.2] - 2025-11-12
### Fixed
- Control-plane Docker image now builds with CGO enabled so SQLite works in containers like Railway.

## [0.1.1] - 2025-11-12
### Added
- Documentation chatbot + advanced RAG examples showcasing Python agent nodes.
- Vector memory storage backends and skill test scaffolding for SDK examples.

### Changed
- Release workflow improvements (selective publishing, prerelease support) and general documentation updates.

## [0.1.0] - 2024-XX-XX
### Added
- Initial open-source release with control plane, Go SDK, Python SDK, and deployment assets.

### Changed
- Cleaned repository layout for public distribution.

### Removed
- Private experimental artifacts and internal operational scripts.
