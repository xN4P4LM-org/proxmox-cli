# Proxmox CLI

This is a CLI tool to help manage the use of a template VM in Proxmox. It allows you to list all VMs, list server resources, and clone a VM from a template.

This is a personal project and is not officially supported by Proxmox, also it is not intended to be used in a production environment.

## Table of Contents
- [Installation](#installation)
- [Usage](#usage)
  - [Commands](#commands)
- [Required Environment Variables](#required-environment-variables)
- [Environments file](#environments-file)


## Installation

1. Clone this repository
2. Ensure you have the following dependencies installed:
   - `make`
   - `go`
3. Run `make build` to build the binary
4. Run `make install` to install the binary to `/usr/local/bin`

## Usage

```bash
proxmox-cli [command] [flags]
```

### Commands

- `list` - List all VMs
```bash
proxmox-cli list
```

Output:
```
Node    VMID    Name                    Status  Tags
----    ----    ----                    ------  ----
node1   101     vm1                     running
node1   200     vm2                     running Template-source
node2   100     vm3                     stopped Template-source
node3   1000    vm4                     stopped Template-source
```

- `resources` - List server resources
```bash
proxmox-cli resources
```

Output:
```
Node    CPU (%) Memory Used (%) VMs
----    ---     ------          -------
node1   0.00    9.24            3
node2   0.00    17.35           2
node3   0.00    33.73           0
node4   0.00    2.74            4
node5   0.00    2.14            3
```

- `clone` - Clones a VM from a template
```bash
proxmox-cli clone -s 100 -e dev -n test-vm -p
```

flags:
```
-n, --name <name> - (required) - Name of the VM"
-s, --source <sourceVM> - (required) - Source VM to clone"
-e, --environment <environment> - (required) - Environment to use"
-p, --poweron - (optional) - Power on the VM after creation"
```

## Required Environment Variables
```bash
set -xg PROXMOX_HOST proxmox.example.com
set -xg PROXMOX_URL https://$PROXMOX_HOST:8006/api2/json
set -xg PROXMOX_TOKEN_ID (token_id)
set -xg PROXMOX_TOKEN_SECRET (token_secret)
set -xg PROXMOX_IGNORE_NODE {nodes to ignore} # comma separated, and optional
set -xg PROXMOX_IGNORE_STORAGE {storage to ignore} # comma separated, and optional
```


## Environments file

The environments file by default is located at `/etc/proxmox-cli/environments.json`. 

This file defines the vmid block for each environment, and will be used to generate a vmid for the `clone` command.

Additionally, the program searches the following locations for the environments file:
- `$HOME/.config/proxmox-cli/environments.json`

The file should be in the following format:
```json
{
    "0": {
        "enviroment": "default",
        "vmid": "100"
    },
    "1": {
        "enviroment":"dev",
        "vmid": "200"
    },
    "2": {
        "enviroment": "staging",
        "vmid": "300"
    },
    "3": {
        "enviroment": "prod",
        "vmid": "400"
    }
}
```
