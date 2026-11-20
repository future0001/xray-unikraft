# syntax=docker/dockerfile:1
# Binaries are built in GitHub Actions as static PIE executables
# (Unikraft Cloud's loader requires position-independent binaries).
FROM scratch
COPY xray /xray
COPY launcher-bin /launcher
EXPOSE 8080
ENTRYPOINT ["/launcher"]
