# U-Service Subscription

Manage user subscriptions and tiers-based authorizations.

## Requirements

- [Go](https://go.dev/doc/install)
- [Docker](https://docker.com/products/docker-desktop/):
  - On Windows, you may run into this message while running Docker: `Windows Hypervisor is not present`.
    To fix this, open a PowerShell with admin privileges and run:
    ```bash
    bcdedit /set hypervisorlaunchtype auto
    ```
- Make:
  - macOS: 
    ```bash
    brew install make
    ```
  - Ubuntu:
    ```bash
    sudo apt-get install make
    ```
  - Windows: Install [chocolatey](https://chocolatey.org/install) (from a PowerShell with admin privileges), then run:
    ```bash
    choco install make
    ```

## Installation

Install dependencies:

```bash
go mod download
```

## Usage

Run the server:

```bash
make run
```

## Test

Install [gotestsum](https://github.com/gotestyourself/gotestsum):

```bash
go install gotest.tools/gotestsum@latest
```

Install [mockery](https://vektra.github.io/mockery/latest/installation/):

```bash
go install github.com/vektra/mockery/v2@v2.43.2
```

Run tests:

```bash
make test
```
