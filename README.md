# Warp (Cli)

token contains the following:

```json
{
  "type": "admin|user",
  "appId": "",
  "releaseId": ""
}
```

this will create `warpdrive.admin.json`

- login to server
- lookup for app's name get the id of the app or create a new one

```bash
  # create warpdrive.admin.json
  # this command should be called once per project. call this command multiple
  # will result the same thing
  warp init --server "127.0.0.1:8080" --user "admin" --pass "admin" --app "my awesome app"
```

```json
  // this is the content of warpdrive.admin.json file created by above command
  // it is the best practice to include this file into .gitignore.
  // if you do want to include it into your git, then make sure that you don't include
  // this on bundle.
  {
    "addr": "localhost:8080",
    "token": "......",
    "certificate": "...."
  }
```

```bash
  # switch to release's name or create a new one
  # behined the scene this command creates warpdrive.json. use this command to
  # switch or create new release cycles.
  # for best practice, before publish any build, call this command to make sure
  # you will target the right release.
  warp release setup --name "prod"
```

```json
  // this is the content of warpdrive.json. this file must be included in bundle
  // and it is the best practice to include it git as well.
  {
    "addr": "...",
    "token": "...",
    "certificate": "..."
  }
```

```bash
  # this will list all of the release cycles. and show you which one is actually selected
  warp release list
```

the output of the above command

```bash
- prod [currently selected]
- alpha
- beta
```

```bash
  warp publish --ios --version "0.0.1"
```


