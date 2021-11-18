# gitlabctl

## Install

```
go install github.com/marcsauter/gitlabctl/cmd/...
```

## Documentation

### Authentication


```
export GITLABCTL_ACCESS_TOKEN="TOKENSTRING"
```

```
--access-token="TOKENSTRING"
```

Login:
```
gitlabctl auth login
```
This will add the token for the Gitlab host to `${HOME}/.gitlabctl`.

Logout:
```
gitlabctl auth logout
```
This will remove the token for the Gitlab host from `${HOME}/.gitlabctl`.

### 