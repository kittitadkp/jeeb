# Auth Debug Prompt Template

Use this when the issue involves login flow, redirects, callback URLs, reverse proxy behavior, Keycloak, cookies, or host-based auth routing.

## Template

### Role

You are a senior authentication, reverse proxy, and web debugging engineer. Investigate the issue conservatively and prefer the smallest safe fix.

### Current issue

Problem: [describe the auth or redirect issue in one or two sentences]  
Where: [URL, route, environment, browser, or service]  
Expected: [what the login or redirect flow should do]  
Actual: [what happens instead]  
Evidence: [exact browser error, status code, log lines, screenshots, or observed behavior]

### Context needed

- Which host or URL is failing
- Whether the failure happens before Keycloak loads or after callback
- Relevant ingress, reverse proxy, frontend base URL, backend auth config, and Keycloak client settings
- Browser behavior such as redirect loop, DNS failure, connection refused, mixed content, or cookie issue
- Any recent changes to local DNS, hosts file, ingress, TLS, callback URI, or realm config

### Investigation checklist

- Confirm whether the hostname resolves and the browser can reach the target host
- Check whether the failure is DNS, ingress, reverse proxy, frontend routing, backend routing, or Keycloak callback config
- Verify Keycloak client redirect URIs, web origins, realm URL, and frontend auth base URL
- Inspect ingress or proxy host rules, service targets, ports, and TLS behavior
- Check browser network behavior and whether any request reaches the app before redirect
- Rank the most likely causes and explain how to verify each one
- Propose the smallest safe fix and the validation steps after the change

### Constraints

- Prefer minimal, safe changes
- Do not suggest secret rotation unless evidence points to it
- Do not rewrite the auth flow if the problem is likely DNS, ingress, or host config
- Keep the advice actionable for local debugging work

### Expected result

Produce a practical debugging plan that identifies the most likely failure layer, the exact checks to run, and the minimal change most likely to fix the issue.

### Output format

1. Short problem framing
2. Ranked likely causes
3. Checks to run now
4. Minimal fix options
5. Validation after fix
6. Any assumptions or missing inputs

## Complete example

### Role

You are a senior authentication, reverse proxy, and web debugging engineer. Investigate the issue conservatively and prefer the smallest safe fix.

### Current issue

Problem: Fix the redirect issue for `http://jeeb-dev.local/` with Keycloak.  
Where: Local development environment in the browser.  
Expected: Opening `http://jeeb-dev.local/` should load the app and redirect unauthenticated users to the Keycloak login page.  
Actual: The browser shows `This site can’t be reached` before any Keycloak page appears.  
Evidence: The failure happens at the initial site URL, not after a visible login redirect.

### Context needed

- Local hostname resolution for `jeeb-dev.local`
- Browser request and any HTTP response, if one exists
- Frontend base URL and auth redirect configuration
- Ingress or reverse proxy config for `jeeb-dev.local`
- Keycloak client redirect URIs and web origins
- Any recent changes to hosts file, ingress, TLS, or local dev networking

### Investigation checklist

- Confirm whether `jeeb-dev.local` resolves locally and points to the expected IP
- Check whether the browser can reach the ingress, reverse proxy, or local frontend on that host
- Determine whether the failure is DNS, connection refused, timeout, ingress host mismatch, or bad local routing
- Verify ingress or proxy host rules for `jeeb-dev.local`, including backend service name and port
- Verify frontend auth settings and Keycloak client redirect URIs for the same host
- Confirm whether any request reaches the app before a redirect to Keycloak is attempted
- Rank the likely causes and recommend the smallest safe fix first

### Constraints

- Prefer minimal, safe changes
- Do not redesign the auth flow
- Focus first on hostname resolution, ingress or proxy routing, and URL configuration before deeper Keycloak changes
- Keep the debugging steps practical for a local developer workstation

### Expected result

Produce a focused debugging prompt that helps identify whether the issue is caused by DNS, ingress, reverse proxy, frontend URL config, or Keycloak client config, then recommends the smallest safe fix.

### Output format

1. Short explanation of the most likely failure layer
2. Ranked root-cause hypotheses
3. Exact checks or commands to run
4. Minimal fix recommendation
5. Validation steps after the fix
6. Missing information to confirm if needed
