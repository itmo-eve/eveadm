# eveadm
Manage EVE virtual machines and containers

## Building
```make```

## Testing
For testing, you need to install Ubuntu Server 18.04, esecute the setup script [tests/build_all_ubuntu_bionic_beaver.sh](https://github.com/itmo-eve/eveadm/blob/master/tests/build_all_ubuntu_bionic_beaver.sh) and run with root privileges:
* ```make test``` -- for general functional testing;
* ```make test_ext``` -- for intgration testing with 'rkt' and Xen 'xl' utilities.

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
