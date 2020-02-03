# eveadm
Manage EVE virtual machines and containers

## Building
```make```

## Testing
This command should run with 'root' privileges:

```make test```

## Running
Supports the command hierarchy described by the command:
```eveadm help```

## Modules
Currently implemented modules:
* xen -- for working with Xen VMs;
* rkt -- for working with RKT containers.

Each module supports actions:
* create
* delete
* info
* list
* start
* stop
* update
