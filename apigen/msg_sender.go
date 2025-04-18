package main

import (
	"strings"

	"github.com/Primital/go-dota2/protocol"
	"github.com/fatih/camelcase"
)

// MsgSender is the sender type of a message.
type MsgSender uint32

//go:generate stringer -type=MsgSender
const (
	MsgSenderNone MsgSender = iota
	MsgSenderGC
	MsgSenderClient
)

// GetMessageSender determines which party (GC, Client, None) would send this message in a session.
func GetMessageSender(msg protocol.EDOTAGCMsg) MsgSender {
	if !IsValidMsg(msg) {
		return MsgSenderNone
	}

	if sender, ok := msgSenderOverrides[msg]; ok {
		return sender
	}

	msgName := msg.String()
	msgName = msgName[6:]
	msgName = strings.TrimPrefix(msgName, "DOTA")

	switch {
	case strings.HasPrefix(msgName, "SQL"):
		return MsgSenderNone
	case strings.Contains(msgName, "ClientToGC"):
		if strings.HasSuffix(msgName, "Response") {
			return MsgSenderGC
		}

		return MsgSenderClient
	case strings.Contains(msgName, "GCResponseTo"):
		return MsgSenderNone
	case strings.Contains(msgName, "GCRequestTo"):
		return MsgSenderNone
	case strings.Contains(msgName, "GCToGC"):
		return MsgSenderNone
	case strings.HasPrefix(msgName, "Server"):
		return MsgSenderNone
	case strings.HasPrefix(msgName, "Gameserver"):
		return MsgSenderNone
	case strings.Contains(msgName, "ServerToGC"):
		return MsgSenderNone
	case strings.Contains(msgName, "GCToServer"):
		return MsgSenderNone
	case strings.Contains(msgName, "GC_GameServer"):
		return MsgSenderNone
	case strings.Contains(msgName, "GCToClient"):
		return MsgSenderGC
	case strings.Contains(msgName, "SignOut"):
		return MsgSenderNone
	case strings.HasSuffix(msgName, "Request"):
		return MsgSenderClient
	case strings.HasSuffix(msgName, "Response"):
		return MsgSenderGC
	case strings.HasPrefix(msgName, "GC"):
		words := camelcase.Split(msgName[2:])
		for _, w := range words {
			if IsWordVerb(w) {
				return MsgSenderClient
			}
		}

		return MsgSenderGC
	}

	return MsgSenderClient
}
