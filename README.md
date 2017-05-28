# configfs

Config files as a fuse filesystem.

## Why

Make config management easier across micro services while complying with production deployment constraints.

## Getting started

### `configfs`'s config

`configfs`'s config lives at `~/.configfs/.config`. Only one parameter is required: the env.

```sh
echo 'env = "local"' > ~/.configfs/.config
```

### Example

Config values are read from `~/.configfs/<env>.toml`, ie : `~/.configfs/local.toml`.

Say we have a config file template `~/my-project/config.json.tmpl` :
```json
{
	"server": {
		"port": LISTEN_PORT,
		"addr": "LISTEN_ADDR"
	}
}
```
We then populate `~/.configfs/local.toml` like so :
```toml
LISTEN_PORT = 3000
LISTEN_ADDR = "127.0.0.1"
```

```sh
$ touch ~/my-project/config.json
$ configfs -source ~/my-project/config.json.tmpl -dest ~/my-project/config.json
```

```json
$ cat ~/my-project/config.json
{
	"server": {
		"port": 3000,
		"addr": "127.0.0.1"
	}
}
```

There it is.

## Envs

You can define as many envs as you like and switch between them by editing `~/.configfs/.config`

## Switch env only for some vars

You can define go-style regexps to override the env only for some vars in `~/.configfs/.config` :

```toml
env = "local"

[env_patterns]
"^MYSQL_" = "staging"
```

## Project override

That all great and all but what if I use `LISTEN_PORT` in many projects and want different values for them ?
`configfs` has an additional flag : `-project`

```sh
$ cat local.toml
MYSQL_PORT = 3306

LISTEN_PORT = 8080 # used when no `-project` flag is given or if project does not match

[project-1]
LISTEN_PORT = 3000 # used when `-project project-1` is given

[project-2]
LISTEN_PORT = 4000 # used when `-project project-1` is given
```
