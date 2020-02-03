# eveadm
Manage EVE virtual machines and containers

## Building
```make```

## Testing
This command should run with 'root' privileges.
```make test```

## Running
Supports the command hierarchy described by the command:
```eveadm help```

## Modules
Currently implemented modules:
* xen -- for working with Xen VMs (now only the 'test' functionality);
* rkt -- for working with RKT containers (now only the 'test' functionality) .

Each module supports actions:
* create
* delete
* info
* list
* start
* stop
* update
