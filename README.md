# xray-unikraft

Generic xray image for Unikraft Cloud. The image contains no secrets: the
xray config is rendered at runtime from environment variables by a small
launcher (itself a static PIE binary, as Unikraft Cloud's loader requires
position-independent executables).

- `UUID` (required): VLESS client UUID
- `WS_PATH` (default `/ws`): WebSocket path
- `PORT` (default `8080`): inbound listen port

Build: GitHub Actions compiles xray-core from source with
`CGO_ENABLED=1 -buildmode=pie -linkmode=external -extldflags=-static-pie`
and pushes `ghcr.io/future0001/xray-unikraft:latest` (linux/amd64).

Deploy example:

```
unikraft run --metro sfo -m 256MiB --restart=always \
  -p 443:8080/http+tls \
  -e UUID=<uuid> -e WS_PATH=/ws -e PORT=8080 \
  --image ghcr.io/future0001/xray-unikraft:latest@sha256:<digest>
```
