# Local Development Setup

## Quick Start

1. **Copy `.env.example` to `.env`**:
   ```bash
   cp .env.example .env
   ```

2. **Edit `.env` with your local API settings**:
   ```bash
   # Edit .env file
   INTENT_API_URL=http://localhost:8080
   INTENT_TOKEN=your-local-dev-token  # Optional if local API doesn't require auth
   ```

3. **Use the CLI** - `.env` file is automatically loaded:
   ```bash
   intent whoami           # Check connection
   intent publish dist/package.itpkg
   ```

## How It Works

The CLI automatically loads `.env` files from the project root (where `go.mod` or `.git` is located).

**Priority order:**
1. `.env` file (if present in project root)
2. Environment variables (INTENT_API_URL, INTENT_TOKEN)
3. Config file (`~/.intent/config.yaml`)
4. Defaults (production API URL)

## Getting a Token for Local API

If your local API requires authentication, you'll need to create a token:

1. **Check your API documentation** - usually there's a way to create tokens
2. **Common endpoints**:
   - `POST /v1/auth/token` - Create token
   - `POST /v1/users/me/tokens` - User tokens
   - Check your API's auth endpoints

3. **For testing, you can create a simple token**:
   ```bash
   # Example: if your API has a token endpoint
   curl -X POST http://localhost:8080/v1/auth/token \
     -H "Content-Type: application/json" \
     -d '{"username": "dev", "password": "dev"}'
   ```

4. **Or skip authentication** (if your local API allows):
   - Leave `INTENT_TOKEN=` empty in `.env`
   - The CLI will try to publish without auth (may fail if API requires it)

## Alternative: Use Environment Variables

If you prefer not to use `.env` files:

```bash
export INTENT_API_URL=http://localhost:8080
export INTENT_TOKEN=your-token
intent publish dist/package.itpkg
```

## Troubleshooting

**Check your configuration:**
```bash
intent whoami    # Shows current API URL and auth status
intent doctor    # Full diagnostics
```

**See what's being used:**
```bash
intent publish dist/package.itpkg --debug
# Shows full HTTP requests including API URL
```

