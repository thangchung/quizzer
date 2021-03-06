package utils

import (
  "runtime"
  "sync"
  "time"
  "github.com/manucorporat/stats"
)

var ips = stats.New()
var messages = stats.New()
var users = stats.New()
var mutexStats sync.RWMutex
var savedStats map[string]uint64

func StatsWorker() {
  c := time.Tick(1 * time.Second)
  var lastMallocs uint64 = 0
  var lastFrees uint64 = 0
  for _ = range c {
    var stats runtime.MemStats
    runtime.ReadMemStats(&stats)

    mutexStats.Lock()
    savedStats = map[string]uint64{
      "timestamp": uint64(time.Now().Unix()),
      "HeapInuse": stats.HeapInuse,
      "StackInuse": stats.StackInuse,
      "Mallocs": (stats.Mallocs - lastMallocs),
      "Frees": (stats.Frees - lastFrees),
      "Inbound": uint64(messages.Get("inbound")),
      "Outbound": uint64(messages.Get("outbound")),
      "Connected": connectedUsers(),
    }
  }
}

func connectedUsers() uint64 {
  connected := users.Get("connected") - users.Get("disconnected")
  if connected < 0 {
    return 0
  }
  return uint64(connected)
}

func Stats() map[string]uint64 {
  mutexStats.RLock()
  defer mutexStats.RUnlock()
  return savedStats
}

