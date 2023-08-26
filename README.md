# arkade-lvlup

`arkade-lvlup` is a CLI tool designed to help users synchronize, install, and manage their tools using `arkade`.

## Features

1. Synchronize tools based on a configuration file.
2. Install or reinstall specific tools.
3. Remove tools.
4. Update shell configuration to include `arkade-lvlup` in the PATH.

## Installation

_TODO: Include installation steps._

## Usage

Here are the available flags and their explanations:

- `-f`: Force sync. This is only valid with `-sync` or `--sync`.
- `-p`: Display `arkade` outputs.
- `-sync` or `-s`: Synchronize tools based on the configuration.
- `-get` or `-g [tool1,tool2,...]`: Install or reinstall specified tools.
- `-remove` or `-r [tool1,tool2,...]`: Remove specified tools.
- `-config-shell` or `-c`: Update the shell configuration to include `arkade-lvlup` in the PATH.

### Examples:

1. **Synchronizing tools based on a configuration**:

    arkade-lvlup -s


2. **Forcefully synchronizing tools**:

    arkade-lvlup -s -f


3. **Installing a tool**:

    arkade-lvlup -g [tool_name1,tool_name2,tool_name3,...]


4. **Removing a tool**:

    arkade-lvlup -r [tool_name1,tool_name2,tool_name3,...]


5. **Updating shell configuration**:

    arkade-lvlup -c


## Configuration

The configuration is stored in a YAML file and specifies which tools should be installed. Here is an example structure:

```yaml
tools:
- tool1
- tool2
- tool3
```

When you synchronize using arkade-lvlup, it will ensure that the tools listed in the configuration are installed and any tools not in the list are removed.

### Directory Tree

    ./arkade-lvlup
    ├── LICENSE
    ├── README.md
    ├── arkade-lvlup
    ├── cmd
    │   ├── get.go
    │   ├── remove.go
    │   ├── shellconfig.go
    │   └── sync.go
    ├── config
    │   ├── flags.go
    │   ├── model.go
    │   └── reader.go
    ├── go.mod
    ├── go.sum
    ├── handlers
    │   ├── gethandler.go
    │   ├── removehandler.go
    │   ├── shellconfighandler.go
    │   └── synchandler.go
    ├── main.go
    ├── shell
    │   └── configurer.go
    └── tools
        ├── config_utils.go
        ├── file_ops.go
        ├── general_utils.go
        ├── installer.go
        ├── remover.go
        └── syncer.go

6 directories, 24 files

## Continuous Integration and Continuous Deployment (CI/CD)
Our project uses GitHub Actions to automate the processes of building, testing, and releasing versions of the application. Here's a brief overview:

### Triggers
- **Pushes to main:** Every push to the **main** branch will trigger the build and test jobs.
- **Pull Requests to main:** Any PR opened against the **main** branch will also initiate the build and test workflows.
- **Commit Message Triggers:** Including the **[FORCE BUILD]** keyword in your commit message will force a build regardless of the branch you're working on.
- **Tagging Releases:** When a new semantic version tag (e.g., 1.2.3) is pushed to the repository, this will initiate the release process.

### Workflows
1. **Build & Test:**

- Checkout the code from the repository.
- Setup the desired Go environment.
- Cache and download Go module dependencies.
- Run all tests in the project using go test ./....

2. **Release:**

- Triggered after a successful build when a new tag is pushed.
- Creates a new release on GitHub using the tag name.
- Builds a binary named arkade-lvlup.
- Attaches the binary to the release as an asset.

### Versioning
We follow the versioning style of arkade, which is semantic versioning without the 'v' prefix. When you're ready to create a new release:

1. Tag the commit: **git tag 1.0.1**.
2. Push the tag: **git push --tags**.

The CI/CD pipeline will then automatically create a release for that version.

### Contributing
We welcome contributions! Please open an issue or submit a pull request if you would like to help improve arkade lvlup.

### License
This project is licensed under the MIT License.