Running GitHub Actions Locally
---
---

You can locally test the GitHub Actions defined in this your project using [nektos act](https://github.com/nektos/act).



## Setup

The first time you run [act](https://github.com/nektos/act), it can take a long time (with no output) to download
the various Docker goodies. Give it time before deciding it is stuck.

Simply run:
```bash
act
```
### Rancher Desktop on Mac

If you are using Rancher Desktop on Mac, you may see errors like this:

```bash
%  act -j test 
INFO[0000] Using docker host 'unix:///var/run/docker.sock', and daemon socket 'unix:///var/run/docker.sock' 
[test/test] 🚀  Start image=ghcr.io/catthehacker/ubuntu:act-latest
[test/test]   🐳  docker pull image=ghcr.io/catthehacker/ubuntu:act-latest platform= username= forcePull=true
Error: Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?
```
You can print out your Docker contexts with the following command:
```bash
% docker context ls
NAME                DESCRIPTION                               DOCKER ENDPOINT                         ERROR
default             Current DOCKER_HOST based configuration   unix:///var/run/docker.sock             
rancher-desktop *   Rancher Desktop moby context              unix:///Users/bhamail/.rd/docker.sock   
```
You can set the DOCKER_HOST environment variable as shown below which will allow `act` to use the correct Docker context:
```bash
export DOCKER_HOST=$(docker context inspect | jq -r '.[0].Endpoints.docker.Host')
```
See: [slight modification that helps for rancher desktop](https://github.com/nektos/act/issues/1051#issuecomment-1732542268)

### Apple Silcon

If running on Apple Silicon (ARM), launch `act` with this flag:
```bash
   act --container-architecture linux/amd64
```

Without this flag, I saw this warning:
```bash
WARN  ⚠ You are using Apple M-series chip and you have not specified container architecture, you might encounter issues while running act. If so, try running it with '--container-architecture linux/amd64'. ⚠
```

and this error:
```bash
...   🐳  docker exec cmd=[bash --noprofile --norc -e -o pipefail /var/run/act/workflow/1] user= workdir=
| docker compose build
[+] Building 0.0s (0/0)                                                         
| permission denied while trying to connect to the Docker daemon socket at unix:///var/run/docker.sock: Get "http://%2Fvar%2Frun%2Fdocker.sock/_ping": dial unix /var/run/docker.sock: connect: permission denied
...
```

## Usage

To get a list of available jobs, run:
```bash
    $ act -l
```

To run a specific job, use the `-j` flag:
```bash
    $ act -j <job-name>
```

For example, to run the `build` job from the `ci.yml` file, use this command:
```bash
    $ act --workflows .github/workflows/ci.yml -j build
```
