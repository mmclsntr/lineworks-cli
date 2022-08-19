# LINE WORKS CLI
CLI for LINE WORKS API

## Description
Command line tool for LINE WORKS API.  
https://developers.worksmobile.com/jp/reference/introduction?lang=ja

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
todo

## Configuration
### Config file
Configuation of this CLI is stored in home directory.

On Linux, macOS,

```bash
$HOME/.config/lineworks/
```

On Windows,

```bash
%USERPROFILE%\.config\linworks\
```

If you want to change config dir path, set `$LINEWORKS_CONFIG_DIR` environment variable.

### Set OAuth client credentials

```bash
lineworks configure set-client \
    --client-id "client_id" \
    --client-secret "client_secret" \
    --profile "profile"
```

Check configure

```bash
lineworks configure get-client --profile "profile"
```

#### Set Redirect URL on Developer Console
**※ Only User Account authorization**

Get redirect URL

```bash
lineworks configure get-redirect-url --profile "profile"
```

Add it to **Redirect URL** setting of App on Developer Console.

### Set Service Account setting
**※ Only Service Account authorization**

```bash
lineworks configure set-service-account \
    --service-account-id "serivce_account_id" \
    --private-key-file "private_key_file_path" \
    --profile "profile"
```

Check configure

```bash
lineworks configure get-service-account --profile "profile"
```

## Get Access Token
### Request Access Token (User Account authorization)
Request

```bash
lineworks auth user-account --scopes "scopes" --profile "profile"
```

Automatically open browser and show sign-on page.

Sign on by user account.

After request, Exec `get-access-token` command to get access token.

```bash
lineworks auth get-access-token --profile "profile"
```

### Request Access Token (Service Account authorization)
**※ Required to set Service Account configuration before.**

Request

```bash
lineworks auth service-account --scopes "scopes" --profile "profile"
```

After request, Exec `get-access-token` command to get access token.

```bash
lineworks auth get-access-token --profile "profile"
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
