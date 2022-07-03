package vatsim

import "github.com/lcnssantos/online-activity/internal/domain"

type count struct {
	data map[string]*domain.Activity
}

func newCount() *count {
	return &count{
		data: make(map[string]*domain.Activity),
	}
}

func (c *count) increment(id string, connType string) {
	if _, ok := c.data[id]; !ok {
		c.data[id] = &domain.Activity{
			Pilot: 0,
			ATC:   0,
		}
	}

	if connType == "pilot" {
		c.data[id] = &domain.Activity{
			Pilot: c.data[id].Pilot + 1,
			ATC:   c.data[id].ATC,
		}
	} else if connType == "atc" {
		c.data[id] = &domain.Activity{
			Pilot: c.data[id].Pilot,
			ATC:   c.data[id].ATC + 1,
		}
	}
}

func (c *count) get() map[string]*domain.Activity {
	return c.data
}
