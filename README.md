# LINE WORKS CLI
CLI for LINE WORKS API

## Description
Command line tool for LINE WORKS API.  
https://developers.worksmobile.com/jp/reference/introduction?lang=ja

## Article
Japanese article : 

## Feature
- Get Access Token the way you choose.
    - [User Account authorization](https://developers.worksmobile.com/jp/reference/authorization-auth?lang=ja)
    - [Service Account authorization](https://developers.worksmobile.com/jp/reference/authorization-sa?lang=ja)

## Usage
See help.

```bash
lineworks -h
lineworks [subcommand] -h
```

## Installation
Download binary files from [Releases](https://github.com/mmclsntr/lineworks-cli/releases)

Choose a download link according to your environment.

- macOS (Intel) : lineworks-cli_x.x.x_Darwin_x86_64.tar.gz
- macOS (M1/M2) : lineworks-cli_x.x.x_Darwin_arm64.tar.gz
- Windows (64bit) : lineworks-cli_x.x.x_Windows_x86_64.tar.gz
- Windows (Arm) : lineworks-cli_x.x.x_Windows_arm64.tar.gz
- Linux (Arm) : lineworks-cli_x.x.x_Linux_arm64.tar.gz
- Linux (Intel 64bit) : lineworks-cli_x.x.x_Linux_x86_64.tar.gz
- Linux (Intel 32bit) : lineworks-cli_x.x.x_Linux_i386.tar.gz

Binary file name

- Linux/macOS : `lineworks`
- Windows : `lineworks.exe`

## Configuration
### Config file
Configuation of this CLI is stored in home directory.

On Linux, macOS,

```bash
$HOME/.config/lineworks/
```

On Windows,

```powershell
%USERPROFILE%\.config\linworks\
```

If you want to change config dir path, set `$LINEWORKS_CONFIG_DIR` environment variable.

### Profile
The setting can be set for each profile, and the setting can be switched by the `profile` parameter specified when executing the command.

Refer to the list of configured profiles with the following command.

On Linux, macOS,

```bash
./lineworks list-profiles
```

On Windows,

```powershell
.\lineworks.exe list-profiles
```

### Set OAuth client credentials

On Linux, macOS,

```bash
./lineworks configure set-client \
    --client-id "client_id" \
    --client-secret "client_secret" \
    --profile "profile"
```

On Windows,

```powershell
.\lineworks.exe configure set-client `
    --client-id "client_id" `
    --client-secret "client_secret" `
    --profile "profile"
```

You can Check configured settings.

On Linux, macOS,

```bash
./lineworks configure get-client --profile "profile"
```

On Windows,

```powershell
.\lineworks.exe configure get-client --profile "profile"
```

#### Set Redirect URL on Developer Console
**※ Only User Account authorization**

Get redirect URL

On Linux, macOS,

```bash
./lineworks configure get-redirect-url --profile "profile"
```

On Windows,

```powershell
.\lineworks.exe configure get-redirect-url --profile "profile"
```

Add it to **Redirect URL** setting of App on Developer Console.

### Set Service Account setting
**※ Only Service Account authorization**

On Linux, macOS,

```bash
./lineworks configure set-service-account \
    --service-account-id "serivce_account_id" \
    --private-key-file "private_key_file_path" \
    --profile "profile"
```

On Windows,

```powershell
.\lineworks.exe configure set-service-account `
    --service-account-id "serivce_account_id" `
    --private-key-file "private_key_file_path" `
    --profile "profile"
```

You can Check configured settings.

On Linux, macOS,

```bash
./lineworks configure get-service-account --profile "profile"
```

On Windows,

```powershell
.\lineworks.exe configure get-service-account --profile "profile"
```

## Get Access Token
### Request Access Token (User Account authorization)
Request

On Linux, macOS,

```bash
./lineworks auth user-account --scopes "scopes" --profile "profile"
```

On Windows,

```powershell
.\lineworks.exe auth user-account --scopes "scopes" --profile "profile"
```

Automatically open browser and show sign-on page.

Sign on by user account.

### Request Access Token (Service Account authorization)
**※ Required to set Service Account configuration before.**

Request

On Linux, macOS,

```bash
./lineworks auth service-account --scopes "scopes" --profile "profile"
```

On Windows,

```powershell
.\lineworks.exe auth service-account --scopes "scopes" --profile "profile"
```

### Refer Access Token
On Linux, macOS,

```bash
./lineworks auth get-access-token --profile "profile"
```

On Windows,

```powershell
.\lineworks auth get-access-token --profile "profile"
```

## Contribution

1. Fork ([https://github.com/mmclsntr/lineworks-cli](https://github.com/mmclsntr/lineworks-cli))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Create a new Pull Request

## Authors
[mmclsntr](https://github.com/mmclsntr)

## License
[MIT](LICENSE.txt)
