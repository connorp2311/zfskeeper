# zfsTools

Tools for working with ZFS datasets

## Installation

1. Clone the repository to your local machine using git clone:
```
git clone https://github.com/connorp2311/zfsTools.git
```

2. Install with go install
```
cd zfsTools
sudo GOBIN=/usr/local/bin/ go install ./...
```
3. Confirm the install worked
```
zfsTools --version
```

And optionally 
```
sudo zfsTools completion bash | sudo tee /etc/bash_completion.d/zfsTools
```

## Usage

Currently the main usage of this tool is for running retention on ZFS snapshots but I plan on expanding this in the future.
```
zfsTools retention <dataset> [flags]
```

Delete snapshots for dataset tank/home with the following retention policies:

* Keep intra-daily snapshots for the past 2 days
* Keep daily snapshots for the past 7 days
* Keep weekly snapshots for the past 4 weeks
* Keep monthly snapshots for the past 12 months

```
zfsTools retention tank/home --intra-daily 2 --daily 7 --weekly 4 --monthly 12
```

Perform a dry run to see what snapshots would be deleted, but do not actually delete any snapshots:

```
zfsTools retention tank/home --intra-daily 2 --daily 7 --weekly 4 --monthly 12 --dry-run
```


For a more information on usage run `zfsTools --help` or view the documentation [here](docs/zfsTools.md)

## Contributing

If you find any issues with the tool or would like to request a feature, please submit an issue on the GitHub repository.

If you would like to contribute to the project, please fork the repository and submit a pull request.