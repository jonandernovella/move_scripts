---

# Rackham -> Dardel data transfer tool

## Overview

This command-line tool is designed to simplify the process of transferring data from one cluster to another. Specifically, it facilitates the transfer of data from Rackham to Dardel.

## Usage

```bash
$ movefiles [check|start]
```

- `check`: Runs in check mode. Data will not be transferred; instead, the tool will identify uncompressed (large) files in the specified directory.
- `start`: will ask you for some information about the data transfer and create a SLURM script so you can perform the transfer.

## Running in Check Mode

In check mode, the tool checks for uncompressed files with certain extensions and warns you if they take up a lot of space. To run in check mode:

```bash
$ movefiles check
```

## Generate the SLURM data transfer script

To create the script for performing the transfer, use the following command:

```bash
$ movefiles start
```

This will prompt you for various details required for the transfer:

- **Directory to Transfer**: Specify the directory you want to transfer.

- **Auto-Delete Option**: Choose whether to automatically delete local files after copying.

- **Keep Large Directories**: Choose whether to discard large subdirectories after packaging.

- **Target Host**: Specify the destination system (default: dardel.pdc.kth.se).

- **Target Directory**: Specify the destination directory on the target host.

- **Username on Target Host**: Provide your username on the target host.

- **Number of Connections**: Specify the number of parallel rsync connections.

- **Project ID**: Enter the UPPMAX project ID (e.g., nais2023-22-999).

The tool generates a Bash script (`transfer_<directory>.sh`) based on your inputs. This script can be edited to set the correct project ID. To execute the script, use the following command:

```bash
$ sbatch transfer_<directory>.sh
```

**Note:** Make sure to review and edit the generated script before executing it.

## Additional Information

- Large directories with more than 100,000 files will be packaged (tar) before transferring.
- Uncompressed files with certain extensions will be identified during the check mode.
