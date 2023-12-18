---

# darsync: Rackham -> Dardel data transfer tool

## Overview

This command-line tool is designed to simplify the process of transferring data from one cluster to another. Specifically, it facilitates the transfer of data from Rackham to Dardel.

## Building the tool

After cloning the repository and changing directory, you may run:

```bash
$ go build
```

## Usage

```bash
$ darsync [check|gen]
```

- `check`: will search for uncompressed files in the specified directory and warn you if they take up a lot of space.
- `gen`: will ask you for some information about the data transfer and create a SLURM script so you can perform the transfer.

## Generate the SLURM data transfer script

To create the script for performing the transfer, use the following command:

```bash
$ darsync gen
```

This will prompt you for various details required for the transfer:

- **Directory to Transfer**: Specify the directory you want to transfer.

- **Target Host**: Specify the destination system (default: dardel.pdc.kth.se).

- **Target Directory**: Specify the destination directory on the target host.

- **Username on Target Host**: Provide your username on the target host.

- **Project ID**: Enter the UPPMAX project ID (e.g., nais2023-22-999).

- **Private key**: Enter the path to the private key in order to connect to the Target Host (Dardel).

The command generates a Bash script (`transfer_<directory>.sh`) based on your inputs. This script can be edited to set the correct project ID. To execute the script, use the following command:

```bash
$ sbatch transfer_<directory>.sh
```

**Note:** Make sure to review and edit the generated script before executing it.

