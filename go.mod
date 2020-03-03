module github.com/itmo-eve/eveadm

go 1.13

require (
	github.com/hako/durafmt v0.0.0-20191009132224-3f39dc1ed9f4
	github.com/lf-edge/eve/pkg/pillar v0.0.0-20200303001835-43cf0e139d28 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.6.1
)

replace github.com/lf-edge/eve/api/go => github.com/lf-edge/eve/api/go v0.0.0-20200301202154-704247b2b305
