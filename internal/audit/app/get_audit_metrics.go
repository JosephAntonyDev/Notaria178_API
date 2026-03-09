package app

import (
	"context"

	"github.com/JosephAntonyDev/Notaria178_API/internal/audit/domain/repository"
)

type AuditMetricsDTO struct {
	UserActions []ActionMetricDTO `json:"user_actions"`
	WorkActions []ActionMetricDTO `json:"work_actions"`
}

type ActionMetricDTO struct {
	Action string `json:"action"`
	Count  int    `json:"count"`
}

type GetAuditMetricsUseCase struct {
	repo repository.AuditRepository
}

func NewGetAuditMetricsUseCase(r repository.AuditRepository) *GetAuditMetricsUseCase {
	return &GetAuditMetricsUseCase{repo: r}
}

func (uc *GetAuditMetricsUseCase) Execute(ctx context.Context, filters repository.AuditFilters) (*AuditMetricsDTO, error) {
	// Call User Actions
	userMetrics, err := uc.repo.GetUserActionMetrics(ctx, filters)
	if err != nil {
		return nil, err
	}

	// Call Work Actions
	workMetrics, err := uc.repo.GetWorkActionMetrics(ctx, filters)
	if err != nil {
		return nil, err
	}

	// Map User Metrics
	userActionDTOs := make([]ActionMetricDTO, len(userMetrics))
	for i, m := range userMetrics {
		userActionDTOs[i] = ActionMetricDTO{
			Action: m.Action,
			Count:  m.Count,
		}
	}

	// Map Work Metrics
	workActionDTOs := make([]ActionMetricDTO, len(workMetrics))
	for i, m := range workMetrics {
		workActionDTOs[i] = ActionMetricDTO{
			Action: m.Action,
			Count:  m.Count,
		}
	}

	return &AuditMetricsDTO{
		UserActions: userActionDTOs,
		WorkActions: workActionDTOs,
	}, nil
}
