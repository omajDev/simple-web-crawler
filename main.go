package main

import (
	"fmt"
	"sync"
)

type (
	Observable interface {
		Add(observer Observer)
		Notify(event interface{})
		Remove(event interface{})
	}

	Observer interface {
		NotifyCallback(event interface{})
	}

	WatchTower struct {
		observer sync.Map
	}

	Soldier struct {
		id   int
		zone string
	}
)

func (wt *WatchTower) Add(observer Observer) {
	wt.observer.Store(observer, struct{}{})
}

func (wt *WatchTower) Remove(observer Observer) {
	wt.observer.Delete(observer)
}

func (wt *WatchTower) Notify(event interface{}) {
	wt.observer.Range(func(key, value interface{}) bool {
		if key == nil {
			return false
		}

		key.(Observer).NotifyCallback(event)
		return true
	})
}

func (s Soldier) NotifyCallback(event interface{}) {
	if event.(string) == s.zone {
		fmt.Printf("Soldier %d, seen an enemy on zone %s\n", s.id, event)
	}
}

func main() {
	watchTower := WatchTower{}
	soldier_1 := Soldier{id: 1, zone: "B"}
	soldier_2 := Soldier{id: 2, zone: "A"}

	watchTower.Add(soldier_1)
	watchTower.Add(soldier_2)

	// Notify Zone A
	watchTower.Notify("A")

	// Notify Zone B
	watchTower.Notify("B")

	// Remove soldier 1 (No soldier on Zome B anymore)
	watchTower.Remove(soldier_1)

	// Notify Zone B (Enemy is free to pass here)
	watchTower.Notify("B")
}
