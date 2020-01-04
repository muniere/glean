package pubsub

import (
	"github.com/muniere/glean/internal/app/server/pubsub/manager"
)

type Manager = manager.Manager
type Config = manager.Config

var NewManager = manager.NewManager
