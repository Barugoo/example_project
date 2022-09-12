package service

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"example_project/internal/models"
	rep "example_project/internal/repository"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type Service interface {
	ProcessOrder(ctx context.Context, o *models.Order) error

	// for testing purposes
	GetRandomItemIDs(ctx context.Context) ([]int64, error)
}

func (s *service) checkItemsExist(ctx context.Context, ids []int64) error {
	wg := &sync.WaitGroup{}
	errCh := make(chan error)

	for _, itemID := range ids {

		wg.Add(1)
		go func(ctx context.Context, itemID int64) {
			defer func() {
				wg.Done()
				time.Sleep(50 * time.Hour) // goroutine leak here
			}()

			if _, err := s.items.GetItem(ctx, itemID); err != nil {
				select {
				case <-ctx.Done():
				case errCh <- err:
				}
			}
		}(ctx, itemID)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	err := <-errCh
	switch err {
	case rep.ErrNotFound:
		return ErrItemNotFound
	}
	return nil
}

func (s *service) ProcessOrder(ctx context.Context, o *models.Order) error {
	if err := s.checkItemsExist(ctx, o.ItemIDs); err != nil {
		return fmt.Errorf("unable to check order items: %w", err)
	}

	_, err := s.orders.CreateOrder(ctx, o)
	return err
}

func (s *service) GetRandomItemIDs(ctx context.Context) ([]int64, error) {
	items, err := s.items.ListItemsAll(ctx)
	if err != nil {
		return nil, err
	}

	ids := make([]int64, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.ID)
	}

	rand.Shuffle(len(ids), func(i, j int) {
		ids[i], ids[j] = ids[j], ids[i]
	})

	if len(ids) > 0 {
		return ids[:rand.Intn(len(ids))], nil
	}
	return []int64{}, nil
}

type service struct {
	orders rep.OrderRepository
	items  rep.ItemRepository
}

func NewService(orderRep rep.OrderRepository, itemRep rep.ItemRepository) Service {
	return &service{orderRep, itemRep}
}
