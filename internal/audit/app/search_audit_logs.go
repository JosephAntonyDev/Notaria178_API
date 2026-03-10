package app

import (
	"context"

	"github.com/JosephAntonyDev/Notaria178_API/internal/audit/domain/repository"
)

type SearchAuditLogsResult struct {
	Logs  []AuditLogDTO
	Total int
}

type SearchAuditLogsUseCase struct {
	repo repository.AuditRepository
}

func NewSearchAuditLogsUseCase(r repository.AuditRepository) *SearchAuditLogsUseCase {
	return &SearchAuditLogsUseCase{repo: r}
}

func (uc *SearchAuditLogsUseCase) Execute(ctx context.Context, filters repository.AuditFilters) (*SearchAuditLogsResult, error) {
	logs, err := uc.repo.List(ctx, filters)
	if err != nil {
		return nil, err
	}

	total, err := uc.repo.Count(ctx, filters)
	if err != nil {
		return nil, err
	}

	result := make([]AuditLogDTO, 0, len(logs))
	for _, l := range logs {
		result = append(result, ToAuditLogDTO(l))
	}

	return &SearchAuditLogsResult{Logs: result, Total: total}, nil
}
