# Deep Agent Demo

Deep Agent demo project built on the Eino framework.

[中文说明](./README.zh.md)

## Quick Start

```
git clone git@github.com:deep-agent/sandbox.git
cd sandbox


mkdir /path/to/memory
mkdir /path/to/memory/workspace
mkdir /path/to/memory/agent

# Set the host persistent directory (workspace, prompts, etc.)
export LOCAL_MEMORY="/path/to/memory"

# On Linux, container user (UID 1000) differs from host user, so adjust permissions
sudo chown -R 1000:1000 /path/to/memory
sudo chown -R 1000:1000 /path/to/memory/workspace
sudo chown -R 1000:1000 /path/to/memory/agent

# Start the sandbox container
make docker-start


# Enter deep-agent-demo and launch the web app
cd deep-agent-demo
make web
```

![homepage](./docs/images/home.png)
