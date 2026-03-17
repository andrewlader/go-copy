# go-copy

| Category | Badges |
| :--- | :--- |
| __Project Information__ | ![GitHub Release](https://img.shields.io/github/v/release/andrewlader/go-copy?style=plastic) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/andrewlader/go-copy?style=plastic) |
| __Testing__ | ![Scrutinizer coverage (GitHub/Bitbucket) with branch](https://img.shields.io/scrutinizer/coverage/g/andrewlader/go-copy/main?style=plastic) |
| __Branch__ | ![GitHub branch check runs](https://img.shields.io/github/check-runs/andrewlader/go-copy/main?style=plastic) |
| __License__ | ![GitHub License](https://img.shields.io/github/license/andrewlader/go-copy?style=plastic) |
| __Purpose__ | ![Static Badge](https://img.shields.io/badge/CLI-Copy_Tool-blue?style=plastic) |

## Description
This command-line tool can copy files from a source directory to one or more destinations. It is easily configured, and can run on Windows, Linux and Mac (Intel and RISC). It's simple to install, and easy to use.

## Installation
There are a couple of ways to install `go-copy`.

### 1. If you have Go installed
```bash
go get https://github.com/andrewlader/go-copy.git
```

### 2. Download and Install the Executable
Download the desired binary from the [releases page](https://github.com/andrewlader/go-copy/releases)

## Configuration File

A configuration file is required for `go-copy` to define which files or folders to copy, where to copy them, and how to handle replacements. The config file is written in YAML format and can be placed in either the `configs/` directory or your user directory (e.g., `C:\Users\<username>\.go-copy\go-copy-config.yaml`).

### Format
Each operation is defined as a YAML key, with the following structure:

```yaml
<operation_name>:
  name: <Display Name>
  source: <Source Path>
  destinations:
    - <Destination Path 1>
    - <Destination Path 2>
  replace: <replace mode>
```

#### Example

```yaml
borderlands3:
  name: Borderlands 3
  source: C:\Users\john\Documents\My Games\Borderlands 3\Saved
  destinations:
    - D:\Game Saves\Borderlands 3 Backup
    - E:\More Game Saves\Borderlands 3 Backup
  replace: skip

oblivion:
  name: Oblivion Remastered
  source: C:\Users\john\Documents\My Games\Oblivion Remastered\Saved\SaveGames
  destinations:
    - D:\Game Saves\Oblivion Remastered Backup
    - E:\More Game Saves\Oblivion Remastered Backup
  replace: always
```

### Instructions for New Users
1. Create a file named `go-copy-config.yaml` in either the `configs/` directory or your user directory (e.g., `C:\Users\<username>\.go-copy`) if you are using Windows. The best location for the configuration file is dependent on the OS. Consult documentation for the most appropriate location.
2. For each backup operation, add a section as shown above.
	 - `name`: A friendly name for the operation.
	 - `source`: The folder or file to copy.
	 - `destinations`: One or more backup locations.
	 - `replace`: How to handle existing files (`never`, `skip`, `always`).
3. Replace modes:
    - `never` - copy over new files, but never replace existing files
    - `skip` - skip files that match the date and size of the backed up file
    - `always` - always copy, replacing the existing backup files if they exist
3. Save the file and run `go-copy --operation <operation_name>` to execute the copy.
4. You can add multiple operations for different games, projects, or folders.

For more details, see the sample config in your user directory or `configs/go-copy-config.yaml`.

## Usage
```bash
go-copy --operation <operation-name> ...<other options>
```

## Building a New Release
1. Push new branch
2. Merge branch
3. If tests succeed, then create a new tage (`v#.#.#`)
4. After creating new tag, click on `New Release`

## Contributing
Currently not accepting contributions.

## Authors and acknowledgment
Show your appreciation to those who have contributed to the project.

