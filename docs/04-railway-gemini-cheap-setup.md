# Ultra-Cheap Railway Deployment & Gemini Spark Integration Guide

This guide details how to deploy the Slack MCP Server on **Railway** and connect it with **Gemini Spark** (or any HTTP/SSE MCP client) for **less than $0.50/month** (well within Railway's $5 Hobby Plan credits).

---

## Why Is This Setup Ultra-Cheap?

By default, in-memory caching for thousands of users and channels can consume 100MB+ RAM in large Slack workspaces. By configuring low-memory flags and Go runtime constraints:
- **RAM usage drops to < 15 MB**
- **CPU utilization drops to < 0.01 vCPU**
- Total monthly compute cost on Railway Hobby plan is estimated at **$0.10 to $0.40/month**, allowing you to run this server continuously virtually for free!

---

## Step 1: Obtain Slack Credentials

Select one of the following authentication methods (see [Authentication Setup](./01-authentication-setup.md) for full details):

1. **User OAuth Token (`xoxp-...`) [Recommended]**:
   - Create a Slack App in your workspace.
   - Add necessary user scopes (e.g. `channels:history`, `channels:read`, `chat:write`, `search:read`).
   - Copy the User OAuth Token (`xoxp-...`).

2. **Bot Token (`xoxb-...`)**:
   - Copy Bot User OAuth Token from your Slack App. (Note: Search API is not supported with bot tokens).

3. **Stealth Browser Tokens (`xoxc-...` + `xoxd-...`)**:
   - Extract `xoxc` token and `d` cookie (`xoxd`) from browser devtools while logged into Slack Web.

---

## Step 2: Deploy on Railway

1. **Push or Fork Repo to GitHub**:
   - Ensure your repo includes `railway.json` and `Dockerfile`.

2. **Create New Project in Railway**:
   - Go to [Railway Dashboard](https://railway.app/).
   - Click **+ New Project** -> **Deploy from GitHub repo**.
   - Select your repository.

3. **Configure Environment Variables in Railway**:
   Navigate to **Variables** tab in your Railway service and add:

   | Variable | Value | Purpose |
   |---|---|---|
   | `SLACK_MCP_XOXP_TOKEN` | `xoxp-...` | Your Slack User Token |
   | `SLACK_MCP_API_KEY` | `your-secret-api-key` | Security key for Gemini Spark authorization |
   | `GOMEMLIMIT` | `64MiB` | Forces Go GC to keep memory usage under 64MB |
   | `GOMAXPROCS` | `1` | Restricts Go runtime to 1 CPU core for lower cost |
   | `SLACK_MCP_LOG_LEVEL` | `warn` | Reduces logging overhead |

4. **Configure Ultra-Cheap Startup Command**:
   In Railway, navigate to **Settings** -> **Deploy** -> **Start Command** (or set `CMD` in Dockerfile/variables):
   ```bash
   mcp-server --transport sse --no-cache
   ```
   > **Note on `--no-cache`:** Skipping initial channel/user cache loading cuts startup memory usage down to sub-15MB. When using `--no-cache`, pass channel/user IDs directly (e.g., `C0123456789`) rather than `#channel-name`.

5. **Expose Public Domain**:
   - Under Railway service **Networking**, click **Generate Domain** (e.g. `slack-mcp-production.up.railway.app`).

---

## Step 3: Connect to Gemini Spark / MCP Clients

Gemini Spark and modern MCP clients support **SSE (Server-Sent Events)** or **Streamable HTTP** transports.

### SSE Endpoint:
- **URL**: `https://<your-railway-domain>.up.railway.app/sse`
- **Headers**:
  ```http
  Authorization: Bearer your-secret-api-key
  ```

### HTTP Endpoint (if starting with `--transport http`):
- **URL**: `https://<your-railway-domain>.up.railway.app/mcp`
- **Headers**:
  ```http
  Authorization: Bearer your-secret-api-key
  ```

---

## Step 4: Verification & Monitoring

1. Check your Railway container logs. You should see:
   ```text
   SSE server listening on 0.0.0.0:XXXX/sse
   ```
2. Verify RAM usage in Railway's **Metrics** tab — it should hover between **10 MB - 20 MB**.
3. Test a tool call from Gemini Spark (e.g. `conversations_history` with a `channel_id`).

---

## Summary Checklist for Ultra-Cheap Usage

- [x] `--no-cache` enabled in start command.
- [x] `GOMEMLIMIT=64MiB` & `GOMAXPROCS=1` env vars configured.
- [x] Public HTTPS domain generated on Railway.
- [x] `SLACK_MCP_API_KEY` set to protect public endpoint.
