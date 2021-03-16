package ledger

import (
	"account/pkg/accountinfo/dto"
	"account/pkg/ledger/model"
	"account/pkg/repository"
	"context"
	"fmt"
	"sort"
	"time"
)

type Service interface {
	CreateLedgerEntry(ctx context.Context, info *model.Ledger) error
	GetEntries(ctx context.Context, query *dto.LogQuery) ([]*model.Ledger, error)
	ExpireCredits(ctx context.Context) error
	AddDebitEntry(ctx context.Context, debitEntry *model.Ledger) error
}

type ledgerService struct {
	repository repository.LedgerRepository
}

func (ls *ledgerService) CreateLedgerEntry(ctx context.Context, info *model.Ledger) error {
	if info.Activity == dto.Credit {
		err := ls.repository.CreateLedgerEntries(ctx, []*model.Ledger{info})
		if err != nil {
			return fmt.Errorf("Service.CreateLedgerEntry failed. Error: %w", err)
		}
	} else if info.Activity == dto.Debit {

	}

	return nil
}

func (ls *ledgerService) GetEntries(ctx context.Context, query *dto.LogQuery) ([]*model.Ledger, error) {
	entries, err := ls.repository.GetEntries(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("Service.CreateLedgerEntry failed. Error: %w", err)
	}

	return entries, nil
}

func (ls *ledgerService) ExpireCredits(ctx context.Context) error {
	return nil
}
func (ls *ledgerService) AddDebitEntry(ctx context.Context, debitEntryRequest *model.Ledger) error {
	aggregateEntries, err := ls.repository.GetEntriesByPriority(ctx)
	if err != nil {
		return err
	}
	groupedAggregateEntries := groupByPriorityAndType(aggregateEntries)
	sortedPriorities := getSortedKeys(groupedAggregateEntries)
	totalCredit := make(map[int64]int64, len(sortedPriorities))
	for _, priority := range sortedPriorities {
		totalCredit[priority] = 0
	}
	// read the debit, credit and expired
	// credit minus expired grouped by priority
	for _, priority := range sortedPriorities {
		// Assumption is that no negative credits are there per priority.
		// The current logic is ensuring that.
		if priorityEntries, ok := groupedAggregateEntries[priority]; ok {
			if creditEntry, ok := priorityEntries[dto.Credit]; ok {
				totalCredit[creditEntry.Priority] += creditEntry.Amount
			}
			if debitEntry, ok := priorityEntries[dto.Debit]; ok {
				totalCredit[debitEntry.Priority] -= debitEntry.Amount
			}
			if expiredEntry, ok := priorityEntries[dto.Expiration]; ok {
				totalCredit[expiredEntry.Priority] -= expiredEntry.Amount
			}
		}
	}
	fmt.Println("555555")
	targetDebit := debitEntryRequest.Amount
	debitEntries := make([]*model.Ledger, 0)

	for priority, credit := range totalCredit {
		if targetDebit <= 0 {
			break
		}
		var debitEntry *model.Ledger
		if credit <= targetDebit {
			debitEntry = &model.Ledger{
				AccountID: debitEntryRequest.AccountID,
				Amount:    credit,
				Priority:  priority,
				Activity:  dto.Debit,
				Expiry:    time.Now(),
				CreatedAt: time.Now(),
			}
			targetDebit -= credit
		} else {
			debitEntry = &model.Ledger{
				AccountID: debitEntryRequest.AccountID,
				Amount:    credit - targetDebit,
				Priority:  priority,
				Activity:  dto.Debit,
				Expiry:    time.Now(),
				CreatedAt: time.Now(),
			}
			targetDebit = 0
		}
		debitEntries = append(debitEntries, debitEntry)
	}
	if targetDebit != 0 {
		return fmt.Errorf("Service.DebitRequest failed. Not enough credits. Error: %w", err)
	}
	err = ls.repository.CreateLedgerEntries(ctx, debitEntries)
	if err != nil {
		return fmt.Errorf("Service.CreateLedgerEntry failed. Error: %w", err)
	}

	return nil
}

func getSortedKeys(entries map[int64]map[string]*model.AggregateEntry) []int64 {
	priorities := make([]int64, 0)
	for priority, _ := range entries {
		priorities = append(priorities, priority)
	}
	int64AsIntValues := make([]int, len(priorities))

	for i, val := range priorities {
		int64AsIntValues[i] = int(val)
	}
	sort.Ints(int64AsIntValues)

	for i, val := range int64AsIntValues {
		priorities[i] = int64(val)
	}
	return priorities
}

func groupByPriorityAndType(entries []*model.AggregateEntry) map[int64]map[string]*model.AggregateEntry {
	groupedEntries := make(map[int64]map[string]*model.AggregateEntry, 0)
	for _, entry := range entries {
		if prioritisedEntries, ok := groupedEntries[entry.Priority]; !ok {
			newPrioritisedEntries := make(map[string]*model.AggregateEntry, 0)
			newPrioritisedEntries[entry.Activity] = entry
			groupedEntries[entry.Priority] = newPrioritisedEntries
		} else {
			prioritisedEntries[entry.Activity] = entry
		}
	}

	return groupedEntries
}
func NewLedgerService(repository repository.LedgerRepository) Service {
	return &ledgerService{
		repository: repository,
	}
}
