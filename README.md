# ipsync

`ipsync` is a service that monitors your external IP address and updates a given DNS record on `Netlify` DNS accordingly. It provides a "dynamic DNS"-like functionality, useful for when you need to connect to a system that sits behind a dynamic IP.

## Building from source

To build the `ipsync` binary you need to have  `Go` (v1.19 or higher) installed.

```bash
git clone https://github.com/gntouts/ipsync.git
cd ipsync

make # This will try to build the binary using Go. If Go is not istalled it will prompt you to install go.
sudo mv ./dist/bin/ipsync /usr/local/bin
```

## Usage

You can use `ipsync` as a Docker container or as a standalone binary.

### Docker

To run `ipsync` as a Docker container, you can use the latest [image](https://hub.docker.com/repository/docker/gntouts/ipsync) to run the container with the env variables required:

```bash
docker run -d --restart unless-stopped --env NETLIFY_TOKEN=<your_token> \
    --env DNS_TARGET=<your_dns_target> --env IPSYNC_TIMEOUT=30 --name ipsync gntouts/ipsync:latest
```

If you want to build the Docker image yourself, you can find the Dockerfile under `dist/Dockerfile`.

### Docker Compose

Another way to run `ipsync` containerized is via Docker Compose. You can edit the `dist/compose.yml` file with your env variables and run:

```bash
docker compose up -d
```

### Run as standalone binary

To run `ipsync` as a standalone binary, you can use the `ipsync` binary provided by the [building from source](#building-from-source) step and run it using your env variables:

```bash
NETLIFY_TOKEN=your_token DNS_TARGET=your_dns_target IPSYNC_TIMEOUT=30 ipsync
```

To check the logs, run:

```bash
sudo cat /var/log/syslog | grep ipsync
```
