# arkade lvlup

arkade lvlup is a Go-based CLI utility designed to effortlessly sync and manage tools in your arkade setup. It ensures your arkade toolkit is always up-to-date with your configuration, promoting consistency and reducing manual work.

## Features

- **Automatic Synchronization**: Ensures your `.arkade/bin` tools are in sync with your configuration in `lvlup.yaml`.
- **Easy Management**: Intuitive commands to add, remove, or force sync your tools.
- **User Feedback**: Provides informative feedback to ensure the user is always in the loop.

## Installation

_TODO: Include installation steps._

## Usage

### Sync tools

    arkade-lvlup -sync

### Forcefully sync tools

    arkade-lvlup -sync -f

### Install or reinstall specified tools

    arkade-lvlup -get "tool1 tool2 tool3"

### Remove specified tools

    arkade-lvlup -remove "tool1 tool2"

### Contributing
We welcome contributions! Please open an issue or submit a pull request if you would like to help improve arkade lvlup.

### License
MIT