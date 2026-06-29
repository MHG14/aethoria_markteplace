package scheduler

import (
	"context"
	"log"
	"time"
)

type AuctionCloser interface {
	CloseExpiredAuctions(ctx context.Context) error
}

type Scheduler struct {
	closer   AuctionCloser
	interval time.Duration
	stop     chan struct{}
}

func New(closer AuctionCloser, interval time.Duration) *Scheduler {
	return &Scheduler{
		closer:   closer,
		interval: interval,
		stop:     make(chan struct{}),
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()

		log.Printf("scheduler started — closing expired auctions every %s", s.interval)

		for {
			select {
			case <-ticker.C:
				if err := s.closer.CloseExpiredAuctions(ctx); err != nil {
					log.Printf("scheduler: close expired auctions: %v", err)
				}
			case <-s.stop:
				log.Println("scheduler stopped")
				return
			case <-ctx.Done():
				log.Println("scheduler context cancelled")
				return
			}
		}
	}()
}

func (s *Scheduler) Stop() {
	close(s.stop)
}
