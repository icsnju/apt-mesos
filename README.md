# apt-mesos

Apt-Mesos is a Mesos Framework for Testing which provides an easy way to build testing environment and schedule testing tasks in a multi-node cluster.

The latest version is 0.1.0 alpha(a prototype version), which support:
* Use RESTful API to submit, list, delete, kill tasks
* Run tasks with `docker container` wrapped
* Measure cluster's metrics(`cpus`, `mem`, `disk`)
Other features will be added in later versions.

## Prerequisites
* golang (nessesary)
* vagrant
* VirtualBox
* vagrant plugins
	* vagrant-omnibus `$ vagrant plugin install vagrant-omnibus`
	* vagrant-berkshelf `$ vagrant plugin install vagrant-berkshelf`
	* vaggrant-hosts `$ vagrant plugin install vagrant-hosts`

**Note:** You should build a `Mesos` environment first. We provide `Vagrantfile` and some scripts to help you build `Mesos Cluster` easily (thanks to [everpeace/vagrant-mesos](https://github.com/everpeace/vagrant-mesos)), or you can use [playa-mesos](https://github.com/mesosphere/playa-mesos) to build `Mesos Standalone` .

## Installation

### Compile source code

```
$ go get github.com/icsnju/apt-mesos
$ cd $GOPATH/src/github.com/icsnju/apt-mesos
$ go build
```

### Build Mesos cluster

```
$ cd vagrant
$ vagrant up
```

## Usage

### Start server

```
$ ./apt-mesos --master=<mesos_addr> --addr=<server_listened_addr>
```

### Submit a task:

```
$ curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" <server_listened_addr>/api/tasks -d@task.json 
```

### Task format

```
{
    "cmd": "sh /data/ping.sh",
    "cpus": "1",
    "mem": "16",
    "docker_image": "busybox",
    "volumes": [
        {
            "container_path":"/data",
            "host_path":"/vagrant"
        }
    ]
}
```

### Complete usage:

```
Usage of ./apt-mesos:
  -addr string
    	Address to listen on <ip:port> (default "127.0.0.1:3030")
  -debug
    	Run in debug mode
  -master string
    	Master to connect to <ip:port> (default "127.0.0.1:5050")
``` 

### Hack the WEBUI

We provide a simple WEBUI, welcome to fork the code [https://github.com/JetMuffin/sher-frontend](https://github.com/JetMuffin/sher-frontend) and contribute to this project!
## Contributors
[JetMuffin](https://github.com/JetMuffin)
