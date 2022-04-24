package triggerService

import (
	"context"
	"errors"

	"github.com/dotenx/dotenx/ao-api/models"
	"github.com/dotenx/dotenx/ao-api/services/executionService"
	"github.com/dotenx/dotenx/ao-api/services/integrationService"
	"github.com/dotenx/dotenx/ao-api/services/utopiopsService"
	"github.com/dotenx/dotenx/ao-api/stores/integrationStore"
	"github.com/dotenx/dotenx/ao-api/stores/triggerStore"
)

type TriggerService interface {
	GetTriggerTypes() (map[string][]triggerSummery, error)
	GetAllTriggers(accountId string) ([]models.EventTrigger, error)
	GetAllTriggersForPipeline(accountId, pipelineName string) ([]models.EventTrigger, error)
	GetAllTriggersForAccountByType(accountId, triggerType string) ([]models.EventTrigger, error)
	GetDefinitionForTrigger(accountId, triggerType string) (models.TriggerDefinition, error)
	AddTrigger(accountId string, trigger models.EventTrigger) error
	DeleteTrigger(accountId string, triggerName, pipeline string) error
	StartChecking(accId string, store integrationStore.IntegrationStore) error
	StartScheduller(accId string) error
	StartSchedulling(trigger models.EventTrigger) error
}

type TriggerManager struct {
	Store              triggerStore.TriggerStore
	UtopiopsService    utopiopsService.UtopiopsService
	ExecutionService   executionService.ExecutionService
	IntegrationService integrationService.IntegrationService
}

type triggerSummery struct {
	Type        string `json:"type"`
	IconUrl     string `json:"icon_url"`
	Description string `json:"description"`
}

func NewTriggerService(store triggerStore.TriggerStore, service utopiopsService.UtopiopsService, execService executionService.ExecutionService, intService integrationService.IntegrationService) TriggerService {
	return &TriggerManager{Store: store, UtopiopsService: service, ExecutionService: execService, IntegrationService: intService}
}

func (manager *TriggerManager) GetTriggerTypes() (map[string][]triggerSummery, error) {
	triggers := make(map[string][]triggerSummery)
	for _, integ := range models.AvaliableTriggers {
		if _, ok := triggers[integ.Service]; ok {
			triggers[integ.Service] = append(triggers[integ.Service], triggerSummery{Type: integ.Type, IconUrl: integ.Icon, Description: integ.Description})
		} else {
			types := make([]triggerSummery, 0)
			types = append(types, triggerSummery{Type: integ.Type, IconUrl: integ.Icon, Description: integ.Description})
			triggers[integ.Service] = types
		}

	}
	return triggers, nil
}

func (manager *TriggerManager) AddTrigger(accountId string, trigger models.EventTrigger) (err error) {
	// todo: make ready the body to be saved in table
	err = manager.Store.AddTrigger(context.Background(), accountId, trigger)
	if err == nil {
		if trigger.Type == "Schedule" {
			err = manager.StartSchedulling(trigger)
		}
	}
	return
}
func (manager *TriggerManager) DeleteTrigger(accountId string, triggerName, pipeline string) error {
	return manager.Store.DeleteTrigger(context.Background(), accountId, triggerName, pipeline)
}

func (manager *TriggerManager) GetAllTriggers(accountId string) ([]models.EventTrigger, error) {
	return manager.Store.GetAllTriggers(context.Background(), accountId)
}
func (manager *TriggerManager) GetAllTriggersForPipeline(accountId, pipelineName string) (map[]models.EventTrigger, error) {
	triggers, err := manager.Store.GetAllTriggers(context.Background(), accountId)
	if err != nil {
		return nil, err
	}
	selected := make(map[string]models.EventTrigger)
	for _, tr := range triggers {
		if tr.Pipeline == pipelineName {
			selected[tr.Name] = tr
		}
	}
	return selected, nil
}

func (manager *TriggerManager) GetAllTriggersForAccountByType(accountId, triggerType string) ([]models.EventTrigger, error) {
	return manager.Store.GetTriggersByType(context.Background(), accountId, triggerType)
}

func (manager *TriggerManager) GetDefinitionForTrigger(accountId, triggerType string) (models.TriggerDefinition, error) {
	intgType, ok := models.AvaliableTriggers[triggerType]
	if ok {
		return intgType, nil
	}
	return models.TriggerDefinition{}, errors.New("invalid trigger type")
}
