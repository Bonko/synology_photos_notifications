# Synology Photos Notifications

A Go tool that monitors a Synology Photos directory and notifies users when new photos are uploaded by others.

## Overview

This tool scans a Synology Photos directory, identifies files by their owner (using file system UID), tracks the number of files per user, and notifies configured users when new files are uploaded by someone else.

## Features

- **Automatic file tracking**: Scans the entire directory structure and tracks files by owner
- **Change detection**: Compares current file counts with previously stored counts to detect new uploads
- **User notifications**: Notifies configured users when new files are detected
- **Persistent state**: Stores file counts per owner to track changes over time
- **Configurable**: Uses YAML configuration for easy setup

## Installation

### Prerequisites

- Go 1.25.4 or later
- Access to a Synology Photos directory

### Build

```bash
go build -o synology_photo_notifications
```

### Run

```bash
./synology_photo_notifications
```

## Configuration

Create a `config.yml` file in the same directory as the executable:

```yaml
rootpath: /path/to/synology/photos
users:
  - name: alice
    email: alice@example.com
  - name: bob
    email: bob@example.com
```

### Configuration Fields

- **rootpath**: The root directory path of your Synology Photos folder
- **users**: List of users to notify
  - **name**: Username (must match the system username of file owners)
  - **email**: Email address for notifications (currently logged, email sending not yet implemented)

## How It Works

1. **File Scanning**: The tool walks through the entire directory tree starting from `rootpath`
2. **Owner Detection**: For each file, it determines the owner using the file's UID and looks up the corresponding username
3. **File Tracking**: Files are grouped by owner in memory
4. **Change Detection**: For each owner, the tool:
   - Counts current files
   - Reads the previously stored count from `last_num_files_<owner>` file
   - Calculates the difference to determine new files
5. **State Persistence**: Updates the `last_num_files_<owner>` file with the current count
6. **Notifications**: Logs notification messages for configured users when new files are detected

## State Files

The tool creates state files in the root directory to track file counts:
- Format: `last_num_files_<owner>`
- Location: Same as `rootpath`
- Content: A single integer representing the last known file count for that owner

## Development

### Running Tests

```bash
go test ./...
```
