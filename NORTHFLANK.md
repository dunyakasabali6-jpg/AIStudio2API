# Northflank Free deployment profile

This repository is a lightweight deployment fork of [Mag1cFall/AIStudio2API](https://github.com/Mag1cFall/AIStudio2API). The standard API and web admin UI are retained. The container uses the CLI build (not the Wails desktop build) and the SpiderMonkey WAA backend. Browser-based account login, local Chrome/CDP discovery, and Camoufox downloads are removed from this server profile; add accounts by pasting AI Studio cookies in the web admin UI or through the bootstrap environment variables. Legacy accounts that depended on Chrome OAuth/DBSC refresh must be re-imported with fresh cookies.

## Container settings

- Build: Dockerfile at repository root
- Exposed/public port: `2048` (HTTP)
- Start command: default Dockerfile command
- Persistent storage: not required for the bootstrap configuration described below; runtime databases/logs are ephemeral unless a volume is attached
- Health check: `/health` on port `2048`
- Updates: redeploy the service with a newer image; container mode does not poll GitHub releases or replace the running binary.

The app keeps its isolated HTML preview shim on `127.0.0.1` using an OS-assigned ephemeral port. There is no `PREVIEW_LISTEN_ADDR` setting and no second public port; this separate loopback origin preserves the browser sandbox boundary. The API and admin UI share public port 2048.

## Required environment variables

Set these in Northflank's service environment. Do not commit real cookies, API keys, or passwords.

```dotenv
LISTEN_ADDR=0.0.0.0:2048
WAA_BACKEND=spidermonkey
AUTO_START=true
WARM_WORKER_LIMIT=1
MAX_ACTIVE_WORKERS=2
WARM_STARTUP_CONCURRENCY=1
PER_ACCOUNT_CONCURRENCY=1
LOG_FILE=false
ADMIN_AUTH_ENABLED=true
ADMIN_USERNAME=admin
ADMIN_PASSWORD=CHANGE_THIS_TO_A_LONG_RANDOM_PASSWORD
PROXY_API_KEY=CHANGE_THIS_TO_A_LONG_RANDOM_API_KEY
```

Use the API key as the bearer key for OpenAI-compatible API requests. Keep the admin password and API key separate.

## Account persistence warning

The current upstream account store reads account state from `AISTUDIO_AUTH_STATES` (default: `auth`). A container without a persistent volume loses files written inside its filesystem when it is recreated. Environment cookie bootstrap is supported. Set `BOOTSTRAP_EMAIL` and `BOOTSTRAP_COOKIE` together for one account, or set `BOOTSTRAP_COOKIES_JSON` to a JSON array of `{"email":"you@gmail.com","cookie":"Cookie: ..."}` entries. At startup, the application validates the cookie format and required AI Studio signing cookies, then writes the account into the normal account directory. The environment secret itself is not logged. Without a persistent volume, runtime account changes are lost on container replacement; the bootstrap secret recreates the configured accounts on the next start. Cookie validity and Google session acceptance still depend on the source account. Bootstrap is create-if-missing: an account ID already present in the configured auth directory is skipped rather than overwritten. If you use a persistent volume and need to replace an existing account's cookies, remove/re-import that account through the admin UI.

## Memory notes

The worker limits above are conservative starting values, not a guarantee that the service will fit every 512 MB plan under load. SpiderMonkey, account count, concurrent requests, generated media, and Go's memory overhead all affect usage. Start with one account and monitor the service's memory before increasing concurrency.

## Security

- Keep `ADMIN_AUTH_ENABLED=true`.
- Set strong, unique values for `ADMIN_PASSWORD` and `PROXY_API_KEY`.
- Never expose the admin UI without authentication.
- Do not paste account cookies into GitHub files or logs.
- Use HTTPS at the platform ingress.
