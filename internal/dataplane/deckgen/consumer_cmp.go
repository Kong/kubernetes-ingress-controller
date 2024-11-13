package deckgen

import (
	"strings"

	"github.com/kong/go-database-reconciler/pkg/file"
)

type fConsumerByUsernameAndCustomID []file.FConsumer

func (f fConsumerByUsernameAndCustomID) Len() int      { return len(f) }
func (f fConsumerByUsernameAndCustomID) Swap(i, j int) { f[i], f[j] = f[j], f[i] }
func (f fConsumerByUsernameAndCustomID) Less(i, j int) bool {
	if f[i].Username == nil && f[j].Username != nil {
		return true
	}
	if f[i].Username != nil && f[j].Username == nil {
		return false
	}
	if f[i].Username != nil && f[j].Username != nil {
		switch cmp := strings.Compare(*f[i].Username, *f[j].Username); cmp {
		case -1:
			return true
		case 1:
			return false
		case 0:
			break
		}
	}

	// Both usernames are equal, compare custom_id.
	if f[i].CustomID == nil && f[j].CustomID != nil {
		return true
	}
	if f[i].CustomID != nil && f[j].CustomID == nil {
		return false
	}
	if f[i].CustomID != nil && f[j].CustomID != nil {
		switch cmp := strings.Compare(*f[i].CustomID, *f[j].CustomID); cmp {
		case -1:
			return true
		case 1:
			return false
		case 0:
			break
		}
	}

	return false
}
