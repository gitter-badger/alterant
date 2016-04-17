
# Alterant

## Features
* Automatic dependency resolution
* OpenPGP encryption
* Symlinking
* Single and multiline commands

## Example
### Preparation
Alterant uses git as a backend for storage. The structure is simple, each
`branch` within a repository is considered to be a `machine` and is orphaned.
Within each `machine`, a YAML file named after the `machine` is used to describe
how a `machine` is to be provisioned. To generate a new `machine`:
````
alterant new test_machine
````

Let's edit `test_machine.yaml` to be:
````yaml
tasks:
  test1:
     dependencies:
      - "test2"
    links:
      -
        target: "test1"
        destination: "test1"
    commands:
        - "echo Test1"
  test2:
    dependencies:
      - "test3"
    links:
      -
        target:  "test2"
        destination: "test2"
    commands:
      - |
        #!/usr/bin/env bash

        # This is a multiline command that is useful for scripts
        echo Test2
  test3:
    links:
      -
        target:  "test3"
        destination: "test3"
        encrypted: true
````
We can see from of our `test_machine.yaml` that a `task` has three fields.
* `dependencies`
* `links`
* `commands`

#### Dependencies
A task can be dependent on multiple other tasks in order for it to successfully
finish. Alterant will automatically resolve the order that a `machine` should
be provisioned. In this case our `tasks` will be executed in the order:
`test3` -> `test2` -> `test1`

#### Links
A `link` is specified assuming that the `target` is relative to the repository
root, and that the `destination` is relative to `$HOME`. Alterant uses OpenPGP
for encryption of sensitive data. A file can be encrypted using the
`encrypted: true` flag for the `target` file. To generate a key-pair used for
encryption/decryption:
````
$ alterant gen-key "NAME" "COMMENT" "EMAIL"
````

The keys are generated and placed at `~/.alterant/{pubring.gpg,secring.gpg}`.
Be sure to add the unencrypted file to your `.gitignore` (test3 in the case of
our example).

#### Commands
A `command` can be any script/command you would run from a shell. They can be
formatted as single line, or multiline.

### Provisioning
With our machine ready for deployment we push the changes upstream. Alterant can
now use the machine for provisioning:
````
$ alterant --verbose provision --clobber --parents https:/path/to/dotfiles.git test_machine
````

