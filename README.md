# Alterant


Alterant is a self-contained dotfile manager and lightweight provisioning tool that supports encryption of sensitive data using the OpenPGP standard.

## Usage
#### Setting Up a New Machine
First you will need a remote repository to host your files. After creating the repository and cloning it, `cd` into the repository and create a new configuration file for the machine you would like to provision. In this case we will use the machine name `machine`.
````
$ alterant new machine
````
After creating the new machine we can see that Alterant automatically created an orphaned branch, checked it out, added a blank `machine.yaml`, and committed it
````
$ git branch
* machine
  master
$ git log --stat
commit 881cd3dfd8d2c485d29e6422f459150a0ab1490d
Author: Alterant <https://github.com/autonomy/alterant>
Date:   ****

    Add machine: machine

 README.md    | 2 ++
 machine.yaml | 0
 2 files changed, 2 insertions(+)
````

The YAML file, `machine.yaml`, contains all `tasks` pertaining to the machine. Here is a simple example:
```` yaml
order:
  - "etc"
  - "bash"
  - "ssh"
  - "git"
tasks:
  etc:
    commands:
      - |
        echo "Hello, Alterant!"
  bash:
    links:
      -
        target:  "bash_profile"
        destination: ".bash_profile"
  git:
    links:
      -
        target: "gitconfig"
        destination: ".gitconfig"
      -
        target: "gitignore"
        destination: ".gitignore"
  ssh:
    links:
      -
        target: "ssh/config"
        destination: ".ssh/config"
      -
        target: "/ssh/id_rsa"
        destination: ".ssh/id_rsa"
        encrypted: true
````
It is important to note the following asssumptions that Alterant makes
* a YAML file named after the machine is in `$PWD`
* all link `target`s are relative to `$PWD`
* all link `destination`s are relative to `$HOME`

This is what the directory structure should be according to our `tasks` in  `machine.yaml`
````
$ tree .
.
├── bash_profile
├── gitconfig
├── gitignore
├── machine.yaml
└── ssh
    ├── config
    └── id_rsa

1 directory, 6 files
````
Notice that the `ssh` task has a file, `id_rsa`, that should be `encrypted`. In order to use the encryption functionality, we first need a private/public key-pair. To generate a key-pair
````
$ alterant gen-key "First Last" "Alterant encryption key" "alterant@autonomy.io"
````
The arguments provided are your name, a comment, and your email address. The key-pair are generated with these and placed in
````
$ ls ~/.alterant/
pubring.gpg secring.gpg
````
Before we commit any files with the `encrypted` option into our repository, we should encrypt them and add the unencrypted filename to the `gitignore` of your dotfiles repository
````
$ alterant encrypt
$ echo ssh/id_rsa >>.gitignore
````
If we look at the directory structure once more, we can see an encrypted version of `id_rsa`
````
$ tree .
.
├── bash_profile
├── gitconfig
├── gitignore
├── machine.yaml
└── ssh
    ├── config
    ├── id_rsa
    └── id_rsa.encrypted

1 directory, 7 files
````
With `id_rsa` ignored and `id_rsa.encrypted` created, we can now safely commit, and push our changes to the repository.

#### Provisioning a Machine
````
  alterant provision https:://github.com/user/dotfiles.git machine
````
