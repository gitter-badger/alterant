# Alterant


Alterant is a self-contained dotfile manager and lightweight provisioning tool that supports encryption of sensitive data, multiline command execution, and a single file configuration for multiple machines.

## Usage
Alterant assumes
* the current directory contains a file named `alter.yaml`
* all link targets are relative to `$PWD`
* and that all link destinations are relative to `$HOME`.

Here is a basic `alter.yaml`:
```` yaml
machines:
  mymachine:
    environment:
      ENV_VAR: "an environment variable available to link and command declarations."
    requests:
      - "common"
      - "provision_mymachine"
tasks:
  common:
    links:
      $MACHINE/bash_profile: ".bash_profile"
    commands:
      - |
        echo "This is multiline command using YAML's literal block feature"
        echo "$ENV_VAR"
  provision_mymachine:
    links:
      mymachine/ssh/config: ".ssh/config"
      mymachine/ssh/id_rsa: ".ssh/id_rsa"
      mymachine/ssh/id_rsa.pub: ".ssh/id_rsa.pub"
    commands:
      - "echo Hello, alterant!"
encrypted:
  - "ssh/id_rsa"
````
## Features
### Environments
Notice that link targets can contain environment variables. The same is true for the destination. Alterant also exports the name of the machine currently being provisioned under the environment variable `$MACHINE` and can be used in any link path or command.

### Commands
Lightweight scripting can be achieved via YAML's literal block feature. This makes commands that require preparation and/or cleaning up easier and ensures that the same environment is available throughout the execution.

### Encyption
It is important to note that in order to encrypt the items listed in the `encrypted` section Alterant's `encrypt` command must be invoked. In preparation for adding encrypted items to a repository invoke the `encrypt` command with the `--remove` flag to remove the unencrypted versions of items listed in the `encrypted` section.

### Tasks and Requests
Tasks can be defined and are available to any machine. To consume a task a request must be indicated in the `requests` section of a machine. If one tasks depends on another we can ensure the dependency is executed first by indicating it before the target task in the `requests` section. In other words, `tasks` are executed in the order they are listed in `requests`.
