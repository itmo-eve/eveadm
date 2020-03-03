package pillar

import "github.com/lf-edge/eve/pkg/pillar/pubsub"
import "github.com/itmo-eve/eveadm/cmd"

func Run(ps *pubsub.PubSub) {
	cmd.Execute()
}
