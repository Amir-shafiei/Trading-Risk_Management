package checklist

import (
	"Trading-Risk_Management/models"
	"Trading-Risk_Management/repository"
	"encoding/json"
	"errors"
)

type ChecklistService interface {
	Create(userID, tradeID uint, items []string) error
	GetByTrade(userID, tradeID uint) (interface{}, error)
	UpdateCheck(userID, tradeID uint, itemIndex int, checked bool) error
	GetDefaults(userID uint) ([]string, error)
	SetDefaults(userID uint, items []string) error
}

type ChecklistServiceImpl struct {
	checklistRepo repository.ChecklistRepository
}

func NewChecklistService(checklistRepo repository.ChecklistRepository) ChecklistService {
	return &ChecklistServiceImpl{checklistRepo: checklistRepo}
}

func (s *ChecklistServiceImpl) Create(userID, tradeID uint, items []string) error {
	jsonItems, err := json.Marshal(items)
	if err != nil {
		return err
	}
	checklist := &models.PreTradeChecklist{
		UserID:  userID,
		TradeID: tradeID,
		Items:   string(jsonItems),
	}
	return s.checklistRepo.Create(checklist)
}

func (s *ChecklistServiceImpl) GetByTrade(userID, tradeID uint) (interface{}, error) {
	return s.checklistRepo.GetByTrade(userID, tradeID)
}

func (s *ChecklistServiceImpl) UpdateCheck(userID, tradeID uint, itemIndex int, checked bool) error {
	checklist, err := s.checklistRepo.GetByTrade(userID, tradeID)
	if err != nil {
		return err
	}

	var items []map[string]interface{}
	if err := json.Unmarshal([]byte(checklist.Items), &items); err != nil {
		return err
	}

	if itemIndex < 0 || itemIndex >= len(items) {
		return errors.New("invalid item index")
	}

	items[itemIndex]["checked"] = checked

	allChecked := true
	for _, item := range items {
		if c, ok := item["checked"].(bool); !ok || !c {
			allChecked = false
			break
		}
	}

	jsonItems, _ := json.Marshal(items)
	checklist.Items = string(jsonItems)
	checklist.AllMet = allChecked

	return s.checklistRepo.Update(checklist)
}

func (s *ChecklistServiceImpl) GetDefaults(userID uint) ([]string, error) {
	defaults, err := s.checklistRepo.GetDefaults(userID)
	if err != nil {
		return []string{}, nil
	}
	var items []string
	json.Unmarshal([]byte(defaults), &items)
	return items, nil
}

func (s *ChecklistServiceImpl) SetDefaults(userID uint, items []string) error {
	jsonItems, err := json.Marshal(items)
	if err != nil {
		return err
	}
	return s.checklistRepo.SetDefaults(userID, string(jsonItems))
}
