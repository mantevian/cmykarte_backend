package game

import (
	"log"
	"slices"
	"time"

	"mantevian.xyz/cmykarte_backend/common/types"
)

var QUEUE_INTERVAL = 1 * time.Second
var CLEANUP_INTERVAL = 1 * time.Minute

var Queue struct {
	Entries map[types.Uuid]QueuedPlayer
} = struct {
	Entries map[types.Uuid]QueuedPlayer
}{Entries: make(map[types.Uuid]QueuedPlayer)}

type QueuedPlayer struct {
	UserId   types.Uuid
	Username string
	JoinedAt time.Time
}

func StartLookingForMatches() {
	go func() {
		ticker := time.NewTicker(QUEUE_INTERVAL)
		defer ticker.Stop()

		for range ticker.C {
			log.Printf("looking for match, currently have %d entries", len(Queue.Entries))
			startQueuedMatch()
		}
	}()

	go func() {
		ticker := time.NewTicker(CLEANUP_INTERVAL)
		defer ticker.Stop()

		for range ticker.C {
			cleared := Cleanup()
			log.Printf("done cleanup, cleared %d entries", cleared)
			startQueuedMatch()
		}
	}()
}

func startQueuedMatch() {
	found, player1, player2 := TakeMatch()
	if found {
		StartMatch(player1, player2)
	}
}

func AddQueueEntry(user types.Uuid, username string) bool {
	if _, ok := Queue.Entries[user]; ok {
		return false
	}

	Queue.Entries[user] = QueuedPlayer{
		UserId:   user,
		Username: username,
		JoinedAt: time.Now(),
	}

	return true
}

func RemoveQueueEntry(user types.Uuid) bool {
	if _, ok := Queue.Entries[user]; !ok {
		return true
	}

	delete(Queue.Entries, user)
	return true
}

func TakeMatch() (found bool, player1 QueuedPlayer, player2 QueuedPlayer) {
	if len(Queue.Entries) < 2 {
		return false, QueuedPlayer{}, QueuedPlayer{}
	}

	keys := make([]string, 0, len(Queue.Entries))
	for k := range Queue.Entries {
		keys = append(keys, k)
	}
	slices.SortFunc(keys, func(a string, b string) int {
		return Queue.Entries[a].JoinedAt.Compare(Queue.Entries[b].JoinedAt)
	})

	player1 = Queue.Entries[keys[0]]
	player2 = Queue.Entries[keys[1]]
	found = true

	delete(Queue.Entries, player1.UserId)
	delete(Queue.Entries, player2.UserId)

	return
}

func Cleanup() (cleared int) {
	cleared = 0
	now := time.Now()

	for userId, player := range Queue.Entries {
		if now.Sub(player.JoinedAt) > time.Minute {
			delete(Queue.Entries, userId)
			cleared += 1
		}
	}

	return
}
