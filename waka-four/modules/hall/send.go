package hall

import (
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"

	"github.com/liuhan907/waka/waka-four/conf"
	"github.com/liuhan907/waka/waka-four/database"
	"github.com/liuhan907/waka/waka-four/proto"
	"github.com/liuhan907/waka/waka/modules/supervisor/supervisor_message"
)

func (my *actorT) send(player database.Player, m proto.Message) {
	if playerData, being := my.players[player]; !being || playerData.Remote == "" {
		return
	}

	log.WithFields(logrus.Fields{
		"player":  player,
		"type":    reflect.TypeOf(m).Elem().Name(),
		"payload": m.String(),
	}).Debugln("send")

	my.supervisor.Tell(&supervisor_message.SendFromHall{uint64(player), m})
}

// ---------------------------------------------------------------------------------------------------------------------

func (my *actorT) sendPlayer(player database.Player) {
	my.send(player, my.ToPlayer(player))
}

func (my *actorT) sendPlayerSecret(player database.Player) {
	my.send(player, my.ToPlayerSecret(player))
}

func (my *actorT) sendHallEntered(player database.Player) {
	my.send(player, &four_proto.HallEntered{
		Player: my.ToPlayerSecret(player),
	})
}

func (my *actorT) sendPlayerNumber(player database.Player, number int32) {
	my.send(player, &four_proto.PlayerNumber{
		Number: number + conf.Option.Hall.MinPlayerNumber,
	})
}

func (my *actorT) sendRecover(player database.Player, is bool, name string) {
	my.send(player, &four_proto.Recover{
		Is:   is,
		Name: name,
	})
}

// ---------------------------------------------------------------------------------------------------------------------

func (my *actorT) sendFourCreateRoomFailed(player database.Player, reason int32) {
	my.send(player, &four_proto.FourCreateRoomFailed{
		Reason: reason,
	})
}

func (my *actorT) sendFourJoinRoomFailed(player database.Player, reason int32) {
	my.send(player, &four_proto.FourJoinRoomFailed{
		Reason: reason,
	})
}

func (my *actorT) sendFourCreateRoomSuccess(player database.Player) {
	my.send(player, &four_proto.FourCreateRoomSuccess{})
}

func (my *actorT) sendFourJoinRoomSuccess(player database.Player) {
	my.send(player, &four_proto.FourJoinRoomSuccess{})
}

func (my *actorT) sendFourUpdateRoom(player database.Player, room fourRoomT) {
	my.send(player, &four_proto.FourUpdateRoom{room.FourRoom2()})
}

func (my *actorT) sendFourLeftRoom(player database.Player) {
	my.send(player, &four_proto.FourLeftRoom{})
}

func (my *actorT) sendFourLeftRoomByDismiss(player database.Player) {
	my.send(player, &four_proto.FourLeftRoomByDismiss{})
}

func (my *actorT) sendFourStarted(player database.Player, number int32) {
	my.send(player, &four_proto.FourStarted{number})
}

func (my *actorT) sendFourUpdateRound(player database.Player, room fourRoomT) {
	my.send(player, &four_proto.FourUpdateRound{room.FourRoundStatus()})
}

func (my *actorT) sendFourDeal(player database.Player, pokers []string) {
	my.send(player, &four_proto.FourDeal{pokers})
}

func (my *actorT) sendFourCompare(player database.Player, room fourRoomT) {
	my.send(player, room.FourCompare())
}

func (my *actorT) sendFourSettle(player database.Player, room fourRoomT) {
	my.send(player, room.FourSettle())
}

func (my *actorT) sendFourFinallySettle(player database.Player, room fourRoomT) {
	my.send(player, room.FourFinallySettle())
}

func (my *actorT) sendFourDismissRequireVote(player, initiator database.Player) {
	my.send(player, &four_proto.FourDismissRequireVote{int32(initiator)})
}

func (my *actorT) sendFourDismissVoteCountdown(player database.Player, number int32) {
	my.send(player, &four_proto.FourDismissVoteCountdown{number})
}

func (my *actorT) sendFourGrabBankerCountdown(player database.Player, number int32) {
	my.send(player, &four_proto.FourGrabBankerCountdown{number})
}

func (my *actorT) sendFourGrabAnimationCountdown(player database.Player, number int32) {
	my.send(player, &four_proto.FourGrabAnimationCountdown{number})
}

func (my *actorT) sendFourSetMultipleCountdown(player database.Player, number int32) {
	my.send(player, &four_proto.FourSetMultipleCountdown{number})
}

func (my *actorT) sendFourRequireGrabBanker(player database.Player) {
	my.send(player, &four_proto.FourRequireGrabBanker{})
}

func (my *actorT) sendFourGrabAnimation(player database.Player, room fourRoomT) {
	my.send(player, room.FourGrabAnimation())
}

func (my *actorT) sendFourRequireSetMultiple(player database.Player) {
	my.send(player, &four_proto.FourRequireSetMultiple{})
}

func (my *actorT) sendFourSetMultipleSuccess(player database.Player, operator database.Player, multiple int32) {
	my.send(player, &four_proto.FourSetMultipleSuccess{
		PlayerId: int32(operator),
		Multiple: multiple,
	})
}

func (my *actorT) sendFourReceivedMessage(player database.Player, sender database.Player, messageType int32, text string) {
	my.send(player, &four_proto.FourReceivedMessage{int32(sender), &four_proto.FourMessage{messageType, text}})
}

func (my *actorT) sendFourUpdateDismissVoteStatus(player database.Player, room fourRoomT) {
	payload, _, _ := room.FourUpdateDismissVoteStatus()
	my.send(player, payload)
}

func (my *actorT) sendFourUpdateContinueWithStatus(player database.Player, room fourRoomT) {
	my.send(player, room.FourUpdateContinueWithStatus())
}

func (my *actorT) sendFourRequireCut(player database.Player, is bool) {
	my.send(player, &four_proto.FourRequireCut{is})
}

func (my *actorT) sendFourRequireCutAnimation(player database.Player, pos int32) {
	my.send(player, &four_proto.FourRequireCutAnimation{pos})
}

func (my *actorT) sendFourDismissFinally(player database.Player, dismiss bool, r fourRoomT) {
	my.send(player, &four_proto.FourDismissFinally{dismiss, r.FourFinallySettle()})
}

// --------------------------------------------------------

func (my *actorT) sendFourSetMultipleSuccessForAll(room fourRoomT, operator database.Player, multiple int32) {
	for _, player := range room.GetPlayers() {
		my.sendFourSetMultipleSuccess(player, operator, multiple)
	}
}

func (my *actorT) sendFourUpdateRoomForAll(room fourRoomT) {
	for _, player := range room.GetPlayers() {
		my.sendFourUpdateRoom(player, room)
	}
}

func (my *actorT) sendFourStartedForAll(room fourRoomT, number int32) {
	for _, player := range room.GetPlayers() {
		my.sendFourStarted(player, number)
	}
}

func (my *actorT) sendFourDismissVoteCountdownForAll(room fourRoomT, number int32) {
	for _, player := range room.GetPlayers() {
		my.sendFourDismissVoteCountdown(player, number)
	}
}

func (my *actorT) sendFourGrabAnimationCountdownForAll(room fourRoomT, number int32) {
	for _, player := range room.GetPlayers() {
		my.sendFourGrabAnimationCountdown(player, number)
	}
}

func (my *actorT) sendFourSetMultipleCountdownForAll(room fourRoomT, number int32) {
	for _, player := range room.GetPlayers() {
		my.sendFourSetMultipleCountdown(player, number)
	}
}

func (my *actorT) sendFourGrabBankerCountdownForAll(room fourRoomT, number int32) {
	for _, player := range room.GetPlayers() {
		my.sendFourGrabBankerCountdown(player, number)
	}
}

func (my *actorT) sendFourUpdateRoundForAll(room fourRoomT) {
	for _, player := range room.GetPlayers() {
		my.sendFourUpdateRound(player, room)
	}
}

func (my *actorT) sendFourUpdateDismissVoteStatusForAll(room fourRoomT) {
	for _, player := range room.GetPlayers() {
		my.sendFourUpdateDismissVoteStatus(player, room)
	}
}

func (my *actorT) sendFourUpdateContinueWithStatusForAll(room fourRoomT) {
	for _, player := range room.GetPlayers() {
		my.sendFourUpdateContinueWithStatus(player, room)
	}
}

func (my *actorT) sendFourDismissFinallyForAll(room fourRoomT, dismiss bool) {
	for _, player := range room.GetPlayers() {
		my.sendFourDismissFinally(player, dismiss, room)
	}
}

// ---------------------------------------------------------------------------------------------------------------------
