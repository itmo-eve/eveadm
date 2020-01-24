# eveadm
Manage EVE virtual machines and containers

## Building
```make```

## Running
Supports the command hierarchy described by the command:
```./eveadm help```

## Modules
Currently implemented modules:
* test -- for testing purposes only, each action simply launches a shell command from the arguments specified for this action;
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
