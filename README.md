# zfskeeper

zfskeeper is a command-line tool for working with ZFS datasets, including retention policies for snapshots.

## Installation

To install zfskeeper, follow these steps:

1. Clone the repository to your local machine using git clone:
```
git clone https://github.com/connorp2311/zfskeeper.git
```

2. Install with `go install`. This will compile the source code and install the binary to `/usr/local/bin`:
```
cd zfskeeper
sudo GOBIN=/usr/local/bin/ go install ./...
```

3. Confirm the install worked, You should see the version of zfskeeper that you just installed:
```
zfskeeper --version
```

4. And optionally, setup tab completion for zfskeeper commands in your bash shell:
```
sudo zfskeeper completion bash | sudo tee /etc/bash_completion.d/zfskeeper
```

## Usage

Currently, zfskeeper provides commands for managing retention policies on ZFS snapshots. Here are some examples:
```
zfskeeper retention <dataset> [flags]
```


This command deletes snapshots for a specified dataset according to the given retention policies.

For example, to keep intra-daily snapshots for the past 2 days, daily snapshots for the past 7 days, weekly snapshots for the past 4 weeks, and monthly snapshots for the past 12 months:

* Keep intra-daily snapshots for the past 2 days
* Keep daily snapshots for the past 7 days
* Keep weekly snapshots for the past 4 weeks
* Keep monthly snapshots for the past 12 months

```
zfskeeper retention tank/home --intra-daily 2 --daily 7 --weekly 4 --monthly 12
```

You can perform a dry run to see what snapshots would be deleted, but not actually delete anything:

```
zfskeeper retention tank/home --intra-daily 2 --daily 7 --weekly 4 --monthly 12 --dry-run
```


For more information on usage, run `zfskeeper --help` or view the [documentation](docs/zfskeeper.md).

## Contributing

If you find any issues with the tool or would like to request a feature, please submit an issue on the [GitHub repository](https://github.com/connorp2311/zfskeeper/issues).

If you would like to contribute to the project, please fork the repository and submit a pull request.
