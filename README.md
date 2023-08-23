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

### Contributing
We welcome contributions! Please open an issue or submit a pull request if you would like to help improve arkade lvlup.

### License
This project is licensed under the MIT License.
