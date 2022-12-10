# ipsync
`ipsync` is a service that monitors your external IP address and updates a given DNS record on `Netlify` DNS accordingly. It provides a "dynamic DNS"-like functionality, useful for when you need to connect to a system that sits behind a dynamic IP.

> Note: To be able to access your system, you need to forward the appropriate port from you router (eg port 2222 if you want to use `ssh`).

## Building from source

To build the `ipsync` binary you need to have  `Go` (v1.16 or higher) installed.

```bash
git clone https://github.com/gntouts/ipsync.git
cd ipsync

make # This will try to build the binary using Go. If Go is not istalled it will prompt you tinstall go.

# You can also use the `go build` command if you want to customize build flags etc
go build -o ./ipsync
```

## Usage

You can use `ipsync` as a Docker container, as a service, or as a standalone binary.

### Run as service

To use `ipsync` as a service, we will first need to make the binary accessible to the system, by moving it to a directory that is in `PATH` (eg `/usr/local/bin`).
Next we will create a `Bash` script where we will define our env vars needed to run `ipsync`.
Create a new file called `run-ipsync` inside a directory that is in `PATH` and add the following code:

```bash
#!/bin/bash
NETLIFY_TOKEN=<your-token>
DNS_TARGET=<your-dns-target> # eg. ipsync.gntouts.gr
IPSYNC_TIMEOUT=30

NETLIFY_TOKEN=$NETLIFY_TOKEN DNS_TARGET=$DNS_TARGET IPSYNC_TIMEOUT=$IPSYNC_TIMEOUT ipsync
```

Once the file is created, we need to make it executable.

```bash
chmod +x run-ipsync
```

Now we need to create a service file inside `/lib/systemd/system/` and name it `ipsync.service` or anything else you like.
Add the following code to the service file:

```bash
[Unit]
Description=Monitor IP address and update DNS record on Netlify DNS

[Service]
ExecStart=run-ipsync

[Install]
WantedBy=multi-user.target
```

Enable the service by running:

```bash
sudo systemctl daemon-reload
sudo systemctl enable ipsync.service
sudo systemctl start ipsync.service
```

You can check that everything runs smoothly by running:

```bash
sudo systemctl status ipsync.service
```

To check the logs, run:

```bash
sudo journalctl -u ipsync.service
# or
sudo cat /var/log/syslog | grep ipsync
```

> Note: If you want to make any changes, you can edit the `run-ipsync` file to change the variables and then restart the service.

### Run with Docker

To run `ipsync` as a Docker container, you can use the latest [image](https://hub.docker.com/repository/docker/gntouts/ipsync) to run the container with the env variables required:

```bash
docker run -d --restart unless-stopped --env NETLIFY_TOKEN=<your_token> \
    --env DNS_TARGET=<your_dns_target> --env IPSYNC_TIMEOUT=30 --name ipsync gntouts/ipsync:latest
```

Or you can create a file to store the variables:

```bash
echo "NETLIFY_TOKEN=<your_token>" > env.list && \
    echo "DNS_TARGET=<your_dns_target>" >> env.list && \
    echo "IPSYNC_TIMEOUT=30" >> env.list

docker run -d --restart unless-stopped --env-file env.list --name ipsync gntouts/ipsync:latest
```

Or you can you can use the `Dockerfile` provided to build your own image:

```bash
docker build -t ipsync:latest .
```

To check the logs, run:

```bash
 docker logs --details ipsync
 ```

### Run as standalone binary

To run `ipsync` as a standalone binary, you can use the `ipsync` binary provided by the [building from source](#building-from-source) step and run it using your env variables:

```bash
NETLIFY_TOKEN=your_token DNS_TARGET=your_dns_target IPSYNC_TIMEOUT=30 ./bin/ipsync
```

To check the logs, run:

```bash
sudo cat /var/log/syslog | grep ipsync
```
