# Guacinator

[![License](https://img.shields.io/github/license/CowDogMoo/guacinator?label=License&style=flat&color=blue&logo=github)](https://github.com/CowDogMoo/guacinator/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/cowdogmoo/guacinator)](https://goreportcard.com/report/github.com/cowdogmoo/guacinator)
[![🚨 CodeQL Analysis](https://github.com/CowDogMoo/guacinator/actions/workflows/codeql-analysis.yaml/badge.svg)](https://github.com/CowDogMoo/guacinator/actions/workflows/codeql-analysis.yaml)
[![🚨 Semgrep Analysis](https://github.com/CowDogMoo/guacinator/actions/workflows/semgrep.yaml/badge.svg)](https://github.com/CowDogMoo/guacinator/actions/workflows/semgrep.yaml)
[![Pre-Commit](https://github.com/CowDogMoo/guacinator/actions/workflows/pre-commit.yaml/badge.svg)](https://github.com/CowDogMoo/guacinator/actions/workflows/pre-commit.yaml)

<img src="docs/images/guacinator-logo.jpeg" alt="Guacinator Logo" width="100%">

Guacinator is a command line utility to interact programmatically with [Apache Guacamole](https://guacamole.apache.org/).

---

## Table of Contents

- [Getting Started](#getting-started)
- [Dependencies](#dependencies)
- [Usage](#usage)
- [For Contributors and Developers](#for-contributors-and-developers)

---

## Getting Started

1. Download and install the [gh cli tool](https://cli.github.com/).

1. Clone the repo:

   ```bash
   gh repo clone CowDogMoo/guacinator
   cd guacinator
   ```

1. Get latest guacinator release:

   ```bash
   OS="$(uname | python3 -c 'print(open(0).read().lower().strip())')"
   ARCH="$(uname -a | awk '{ print $NF }')"
   gh release download -p "*${OS}_${ARCH}.tar.gz"
   tar -xvf *tar.gz
   ```

---

## Dependencies

- [Install asdf](https://asdf-vm.com/):

  ```bash
  git clone https://github.com/asdf-vm/asdf.git ~/.asdf
  ```

- Install and use asdf plugins to manage go, python, and ruby for this project:

  ```bash
  source .asdf
  ```

  Alternatively, you can pick and choose which plugins to install:

  ```bash
  # Employ asdf for this project's python:
  source .asdf python
  ```

- [Install pre-commit](https://pre-commit.com/):

  ```bash
  python3 -m pip install --upgrade pip
  python3 -m pip install pre-commit
  ```

- [Install Mage](https://magefile.org/):

  ```bash
  go install github.com/magefile/mage@latest
  ```

---

## Usage

- Compile guacinator (if not using downloaded release):

  ```bash
  go build
  ```

- Create a new VNC connection in Guacamole:

  ```bash
  GUAC_URL=https://guacamole.techvomit.xyz
  CONNECTION_NAME=test-connection
  GUAC_USER=guacadmin
  GUAC_PW=guacadmin
  VNC_IP="$(kubectl get service -o wide | grep ubuntu-vnc | awk '{print $3}')"
  VNC_PW="$(kubectl exec -it deployments/ubuntu -- zsh -c 'vncpwd \
    /home/ubuntu/.vnc/passwd' | awk -F ' ' '{print $2}')"

  ./guacinator guacamole -u "${GUAC_USER}" -p "${GUAC_PW}" -l "${GUAC_URL}" \
    --connection "${CONNECTION_NAME}" --vnc-ip "${VNC_IP}" --vnc-pw "${VNC_PW}"
  ```

- Update the `guacadmin` user's password in Guacamole:

  ```bash
  GUAC_URL=https://guacamole.techvomit.xyz
  GUAC_USER=guacadmin
  # Default unless changed
  GUAC_PW=guacadmin
  NEW_GUAC_PW=s1cknewpassword

  ./guacinator guacamole -u "${GUAC_USER}" -p "${GUAC_PW}" -l "${GUAC_URL}" \
    --guacadmin-pw "${NEW_GUAC_PW}"
  ```

- Create a new Guacamole admin user:

  ```bash
  GUAC_URL=https://guacamole.techvomit.xyz
  GUAC_USER=guacadmin
  # Default password for the new account will be this:
  GUAC_PW=guacadmin
  NEW_GUAC_ADMIN=guacadmindos

  ./guacinator guacamole -u "${GUAC_USER}" -p "${GUAC_PW}" -l "${GUAC_URL}" \
    --new-admin "${NEW_GUAC_ADMIN}"
  ```

- Delete a Guacamole user:

  ```bash
  GUAC_URL=https://guacamole.techvomit.xyz
  GUAC_USER=guacadmin
  GUAC_PW=guacadmin
  USER_TO_DELETE=someuser

  ./guacinator guacamole -u "${GUAC_USER}" -p "${GUAC_PW}" -l "${GUAC_URL}" \
  --delete-user "${USER_TO_DELETE}"
  ```

---

## For Contributors and Developers

If you're a developer or contributor looking to work on
Guacinator, follow these additional setup steps:

1. Clone the repository and navigate to the project root:

   ```bash
   gh repo clone CowDogMoo/guacinator
   cd guacinator
   ```

1. Install the dependencies:

   ```bash
   mage InstallDependencies
   ```
