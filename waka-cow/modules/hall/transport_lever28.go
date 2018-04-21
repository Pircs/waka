package hall

import (
	"sync/atomic"

	"github.com/liuhan907/waka/waka-cow/proto"
	"github.com/liuhan907/waka/waka/modules/supervisor/supervisor_message"
	"github.com/sirupsen/logrus"
)

func (my *actorT) playerTransportedLever28(player *playerT, ev *supervisor_message.PlayerTransported) bool {
	switch evd := ev.Payload.(type) {
	case *waka.Lever28CreateRedPaperBag:
		my.Lever28CreateRedPaperBag(player, evd)
	case *waka.Lever28Grab:
		my.Lever28Grab(player, evd)
	case *waka.Lever28Leave:
		my.Lever28Leave(player, evd)
	default:
		return false
	}
	return true
}

func (my *actorT) Lever28CreateRedPaperBag(player *playerT, ev *waka.Lever28CreateRedPaperBag) {
	ev.GetOption().Money *= 100

	id := atomic.AddInt32(&my.lever28IdPool, 1)

	bag := new(lever28BagT)
	bag.Create(my, id, ev.GetOption(), player.Player)
}

func (my *actorT) Lever28Grab(player *playerT, ev *waka.Lever28Grab) {
	bag, being := my.lever28Bags[ev.GetId()]
	if !being {
		log.WithFields(logrus.Fields{
			"player": player.Player,
			"id":     ev.GetId(),
		}).Warnln("grab lever28 but not found")
		my.sendRedGrabFailed(player.Player, 1)
		return
	}

	bag.Grab(player)
}

func (my *actorT) Lever28Leave(player *playerT, ev *waka.Lever28Leave) {
	player.InsideLever28 = 0
}
