# gormon

gormon is an utility for monitoring remote Redis instances via SSH.

In the beginning, as usual, run:
```
    $ go get github.com/austinov/gormon
    $ cd $GOPATH/github.com/austinov/gormon
    $ glide up
    $ go build
```

Before run *gormon*, configure it using the settings in yaml (default is ./dev.yaml).
First of all, you need to configure Redis hosts and directory contains SSH keys.
To use custom configuration file path, please flags:
```
  -cfg-dir string
    	dir with app's config (default ".")
  -cfg-name string
    	app's config base file name (default "dev")
```

By default, will be printed the following values fields from Redis **INFO**:
  - used_memory
  - used_memory_rss
  - connected_clients
  - blocked_clients
  - rejected_connections
  - keyspace_hits
  - keyspace_misses
  - used_cpu_sys
  - used_cpu_user
  - aof_last_write_status

Additionally, always printed the fields:
  - host
  - tstamp
  - error

To setup fields play with the **fields-out** setting.

To run the utility:
```
    $ ./gormon
```

##### Screenshot
![gormon](https://github.com/austinov/gormon/blob/assets/screenshot_terminal.gif)


To run the utility as a web-server and update stats via [Server-sent events](http://www.w3.org/TR/eventsource/) use:
```
    $ ./gormon -server=localhost:3000
```
##### Screenshot
![gormon](https://github.com/austinov/gormon/blob/assets/screenshot_web.png)
