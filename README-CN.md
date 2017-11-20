# drone-rsync

[![Build Status](http://beta.drone.io/api/badges/drone-plugins/drone-rsync/status.svg)](http://beta.drone.io/drone-plugins/drone-rsync)
[![Coverage Status](https://aircover.co/badges/drone-plugins/drone-rsync/coverage.svg)](https://aircover.co/drone-plugins/drone-rsync)
[![](https://badge.imagelayers.io/plugins/drone-rsync:latest.svg)](https://imagelayers.io/?images=plugins/drone-rsync:latest 'Get your own badge on imagelayers.io')

Drone plugin to deploy or update a project via Rsync. For the usage information and a listing of the available options please take a look at [the docs](DOCS.md).

# feature

1. 多主机并行同步                   默认开启，通过 `sync: true` 关闭并发
2. 远程主机同步完成后修改用户和用户组  默认 root：root 通过 `chown: 33:33`
3. 同步完成后,远程主机修改权限         `chmod: 0755`
4. 支持所有bash脚本

### Example

```sh
./drone-rsync <<EOF
{
    "repo": {
        "clone_url": "git://github.com/drone/drone",
        "owner": "drone",
        "name": "drone",
        "full_name": "drone/drone"
    },
    "system": {
        "link_url": "https://beta.drone.io"
    },
    "build": {
        "number": 22,
        "status": "success",
        "started_at": 1421029603,
        "finished_at": 1421029813,
        "message": "Update the Readme",
        "author": "johnsmith",
        "author_email": "john.smith@gmail.com"
        "event": "push",
        "branch": "master",
        "commit": "436b7a6e2abaddfd35740527353e78a227ddcb2c",
        "ref": "refs/heads/master"
    },
    "workspace": {
        "root": "/drone/src",
        "path": "/drone/src/github.com/drone/drone"
    },
    "vargs": {
        "user": "root",
        "host": "localhost",
        "port": 22,
        "source": "dist/",
        "target": "/path/on/server",
        "delete": false,
        "recursive": false
    }
}
EOF
```

## Docker

Build the container using `make`:

```
make deps docker
```

### Example

```sh
docker run -i plugins/drone-rsync <<EOF
{
    "repo": {
        "clone_url": "git://github.com/drone/drone",
        "owner": "drone",
        "name": "drone",
        "full_name": "drone/drone"
    },
    "system": {
        "link_url": "https://beta.drone.io"
    },
    "build": {
        "number": 22,
        "status": "success",
        "started_at": 1421029603,
        "finished_at": 1421029813,
        "message": "Update the Readme",
        "author": "johnsmith",
        "author_email": "john.smith@gmail.com"
        "event": "push",
        "branch": "master",
        "commit": "436b7a6e2abaddfd35740527353e78a227ddcb2c",
        "ref": "refs/heads/master"
    },
    "workspace": {
        "root": "/drone/src",
        "path": "/drone/src/github.com/drone/drone"
    },
    "vargs": {
        "user": "root",
        "host":" test.drone.io",
        "port": 22,
        "source": "dist/",
        "target": "/path/on/server",
        "delete": false,
        "recursive": false
    }
}
EOF
```
