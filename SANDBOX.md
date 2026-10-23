# Working in a sandbox

You can do everything in this repo without installing tools on your own machine. The options below run from most isolated to least. They all end up with the same tools at the same versions, because they all use `.devcontainer/install-tools.sh`.

Reviewing a merge request needs none of this: GitLab's web UI is enough.

## Option 1: GitHub Codespaces

Nothing runs on your machine apart from a browser (or VS Code connected to the codespace).

[![Open in GitHub Codespaces](https://github.com/codespaces/badge.svg)](https://github.com/codespaces/new/BossRoxall/ts-interviews?quickstart=1)

1. Click the button above (or open [github.com/codespaces/new/BossRoxall/ts-interviews](https://github.com/codespaces/new/BossRoxall/ts-interviews?quickstart=1)) and sign in to GitHub if asked. You don't need access to the repository beyond that. If you already have a codespace for this repo, the link reopens it instead of creating another.
2. Keep the default branch (`main`) and pick a machine with at least 4 cores and 8 GB of memory, then create the codespace.
3. Wait for the terminal to print `kind cluster 'interview' is ready`. The first start takes a few minutes while the image builds.
4. Run `tilt up` in the terminal. The Ports tab forwards 10350 (Tilt UI) and 8080 (the service).

The codespace counts against your GitHub account's free monthly Codespaces allowance. A 4-core machine uses 4 core-hours per hour; GitHub Free included 120 core-hours a month at the time of writing, which is far more than an interview needs. Stop the codespace when you finish, and delete it from [github.com/codespaces](https://github.com/codespaces) when you no longer need it.

Your changes stay in your codespace. You don't need to commit or push anything, and you can't push to this repository.

To use a terminal on your own machine instead of the browser, `gh codespace ssh` connects to the codespace, and `gh codespace code` opens it in local VS Code.

Other hosted services that read `.devcontainer/devcontainer.json`, such as DevPod with a cloud provider, work the same way.

## Option 2: dev container on your machine

Your machine needs Docker (OrbStack, Docker Desktop or Podman) and either VS Code with the Dev Containers extension or the [dev container CLI](https://github.com/devcontainers/cli). Every other tool lives inside the container image and goes away when you delete it.

In VS Code, open the repo and choose "Reopen in Container". From a terminal:

```sh
devcontainer up --workspace-folder .
devcontainer exec --workspace-folder . bash
# inside the container, once the kind cluster is ready:
tilt up
```

The container runs its own Docker daemon and a one-node kind cluster, so Tilt never touches your host's Docker images or kube contexts. Running Docker inside a container needs the container to be privileged, though, so this keeps your machine tidy without being a strong security boundary. If isolation is the goal, use option 1 or option 3.

To clean up, remove the container, its Docker volume and the image. `docker ps -a`, `docker volume ls` and `docker images` show their names; they contain `vsc-` or `devcontainer`.

## Option 3: throwaway Linux VM

A local VM keeps everything off your host operating system. [Multipass](https://canonical.com/multipass) works on macOS, Windows and Linux:

```sh
multipass launch 24.04 --name interview --cpus 4 --memory 8G --disk 30G
multipass shell interview
```

OrbStack machines, Lima, UTM or a cloud VM work just as well. Inside any Debian or Ubuntu VM:

```sh
sudo apt-get update && sudo apt-get install -y docker.io git make
sudo usermod -aG docker "$USER" && newgrp docker

git clone <repo-url> rate-limiter && cd rate-limiter
sudo .devcontainer/install-tools.sh
kind create cluster --name interview

tilt up --host 0.0.0.0
```

`--host 0.0.0.0` makes the Tilt UI and the port-forward reachable from your host at the VM's address (`multipass info interview` shows it), on ports 10350 and 8080. Only use it on a VM network that is private to your machine.

To clean up: `multipass delete --purge interview`.

## What the environment contains

`.devcontainer/install-tools.sh` installs these into `/usr/local` and checks each download against a SHA-256 pinned in the script before unpacking it:

| Tool | Version |
| --- | --- |
| Go | 1.27.3, plus gopls 0.24.0 for the VS Code Go extension (checked against the Go checksum database) |
| Node | 26.13.0 |
| uv | 0.13.7 |
| Python | 3.13, installed by uv, which checks the hashes of the builds it downloads |
| kubectl | 1.37.2 |
| Helm | 4.3.1 |
| Kustomize | 5.8.3 |
| kind | 0.33.0 |
| Tilt | 0.37.8 |

Run `make doctor` inside the sandbox to confirm everything is there.
