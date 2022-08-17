# LINE WORKS CLI
CLI for LINE WORKS API

## Description
Command line tool for LINE WORKS API.
https://developers.worksmobile.com/jp/reference/introduction?lang=ja

## Config file
Configuation of this CLI is stored in home directory.

`~/.config/lineworks/`

## Usage
See help.

```bash
lineworks -h
lineworks [subcommand] -h
```

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
※ Only User Account authorization

Get redirect URL

```bash
lineworks configure get-redirect-url --profile "profile"
```

Add redirect URL to **Redirect URL** setting of App on Developer Console.

### Set Service Account setting
※ Only Service Account authorization

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

### Get Access Token (User Account authorization)

```bash
lineworks auth user-account --scopes "scopes" --profile "profile"
```

Open browser and sign-on page.

Sign on and get access token.

After request to get access token, Exec `get-access-token` command.

```bash
lineworks auth get-access-token --profile "profile"
```

### Get Access Token (Service Account authorization)

```bash
lineworks auth service-account --scopes "scopes" --profile "profile"
```

After request to get access token, Exec `get-access-token` command.

```bash
lineworks auth get-access-token --profile "profile"
```

