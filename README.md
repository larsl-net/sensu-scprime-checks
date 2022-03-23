# sensu-scprime-checks

## Table of Contents
- [sensu-scprime-checks](#sensu-scprime-checks)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Files](#files)
  - [Usage examples](#usage-examples)
    - [Help output](#help-output)
  - [Configuration](#configuration)
    - [Asset registration](#asset-registration)
    - [Check definition](#check-definition)
  - [Installation from source](#installation-from-source)
  - [Additional notes](#additional-notes)
  - [Contributing](#contributing)

## Overview

The sensu-scprime-checks is a [Sensu Check][6] that checks the state of a ScPrime Node.

## Files
* scprime-wallet-check

## Usage examples
### Help output

```
Check if ScPrime Wallet is unlocked.

Usage:
  scprime-wallet-check [flags]
  scprime-wallet-check [command]

Available Commands:
  help        Help about any command
  version     Print the version number of this plugin

Flags:
  -h, --help            help for scprime-wallet-check
  -p, --port int        ScPrime server port to check (default 4280)
  -s, --server string   ScPrime server to check (default "localhost")

Use "scprime-wallet-check [command] --help" for more information about a command.
```

## Configuration

### Asset registration

[Sensu Assets][10] are the best way to make use of this plugin. If you're not using an asset, please
consider doing so! If you're using sensuctl 5.13 with Sensu Backend 5.13 or later, you can use the
following command to add the asset:

```
sensuctl asset add larsl-net/sensu-scprime-checks
```

If you're using an earlier version of sensuctl, you can find the asset on the [Bonsai Asset Index][https://bonsai.sensu.io/assets/larsl-net/sensu-scprime-checks].

### Check definition

```yml
---
type: CheckConfig
api_version: core/v2
metadata:
  name: sensu-scprime-check
  namespace: default
spec:
  command: scprime-wallet-check
  subscriptions:
  - system
  runtime_assets:
  - larsl-net/sensu-scprime-checks
```

## Installation from source

The preferred way of installing and deploying this plugin is to use it as an Asset. If you would
like to compile and install the plugin from source or contribute to it, download the latest version
or create an executable script from this source.

From the local path of the sensu-scprime-check repository:

```
cd src/scprime-wallet-check/
go build
```

## Additional notes

## Contributing

For more information about contributing to this plugin, see [Contributing][1].

[1]: https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md
[6]: https://docs.sensu.io/sensu-go/latest/reference/checks/
[10]: https://docs.sensu.io/sensu-go/latest/reference/assets/
