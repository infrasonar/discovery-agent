[![CI](https://github.com/infrasonar/discovery-agent/workflows/CI/badge.svg)](https://github.com/infrasonar/discovery-agent/actions)
[![Release Version](https://img.shields.io/github/release/infrasonar/discovery-agent)](https://github.com/infrasonar/discovery-agent/releases)

# InfraSonar Discovery Agent

Documentation: https://docs.infrasonar.com/collectors/agents/discovery/

## Environment variables

Environment                 | Default                               | Description
----------------------------|---------------------------------------|-------------------
`NETWORK`                   | _required_                            | Network to scan. For example `192.168.0.1/24`.
`DAEMON`                    | `0`                                   | Start the agent as daemon when set to `1`. If not (`0`), the agent will run only once.
`CONFIG_PATH`       		| `/etc/infrasonar` 			        | Path where configuration files are loaded and stored _(note: for a user, the `$HOME` path will be used instead of `/etc`)_
`TOKEN`                     | _required_                            | Token used for authentication _(This MUST be a container token)_.
`ASSET_NAME`                | _none_                                | Initial Asset Name. This will only be used at the announce. Once the asset is created, `ASSET_NAME` will be ignored.
`ASSET_ID`                  | _none_                                | Asset Id _(If not given, the asset Id will be stored and loaded from file)_.
`API_URI`                   | https://api.infrasonar.com            | InfraSonar API.
`SKIP_VERIFY`				| _none_						        | Set to `1` or something else to skip certificate validation.
`CHECK_NMAP_INTERVAL`       | `14400`                               | Interval in seconds for `namp` check or `0` to disable the check.

## Download

- [Linux (amd64)](https://github.com/infrasonar/discovery-agent/releases/download/v1.0.0/discovery-agent-linux-amd64-1.0.0.tar.gz)
- [Windows (amd64)](https://github.com/infrasonar/discovery-agent/releases/download/v1.0.0/discovery-agent-windows-amd64-1.0.0.tar.gz)
- [Darwin (amd64)](https://github.com/infrasonar/discovery-agent/releases/download/v1.0.0/discovery-agent-darwin-amd64-1.0.0.tar.gz)
- [Solaris (amd64)](https://github.com/infrasonar/discovery-agent/releases/download/v1.0.0/discovery-agent-solaris-amd64-1.0.0.tar.gz)

> If your platform is not listed above, refer to the [build from source](#build-from-source) section for instructions.

## Docker

**Latest stable release**
```
docker pull ghcr.io/infrasonar/discovery-agent:latest
```

**Unstable release**
```
docker pull ghcr.io/infrasonar/discovery-agent:unstable
```

## Build from source
```
CGO_ENABLED=0 go build -trimpath -o discovery-agent
```

## Schedule

Download the latest release from here:

https://github.com/infrasonar/discovery-agent/releases/
```

Ensure the binary is executable:
```
chmod +x discovery-agent
```

Copy the binary to `/usr/sbin/infrasonar-discovery-agent`

```
sudo cp discovery-agent /usr/sbin/infrasonar-discovery-agent
```

### Using Systemd

```bash
sudo touch /etc/systemd/system/infrasonar-discovery-agent.service
sudo chmod 664 /etc/systemd/system/infrasonar-discovery-agent.service
```

**1. Using you favorite editor, add the content below to the file created:**

```
[Unit]
Description=InfraSonar Discovery Agent
Wants=network.target

[Service]
EnvironmentFile=/etc/infrasonar/discovery-agent.env
ExecStart=/usr/sbin/infrasonar-discovery-agent

[Install]
WantedBy=multi-user.target
```

**2. Create the directory `/etc/infrasonar`**

```bash
sudo mkdir /etc/infrasonar
```

**3. Create the file `/etc/infrasonar/discovery-agent.env` with at least:**

```
TOKEN=<YOUR TOKEN HERE>
```

Optionaly, add environment variable to the `discovery-agent.env` file for settings like `ASSET_ID` or `CONFIG_PATH` _(see all [environment variables](#environment-variables) in the table above)_.

**4. Reload systemd:**

```bash
sudo systemctl daemon-reload
```

**5. Install the service:**

```bash
sudo systemctl enable infrasonar-discovery-agent
```

**Finally, you may want to start/stop or view the status:**
```bash
sudo systemctl start infrasonar-discovery-agent
sudo systemctl stop infrasonar-discovery-agent
sudo systemctl status infrasonar-discovery-agent
```

**View logging:**
```bash
journalctl -u infrasonar-discovery-agent
```
