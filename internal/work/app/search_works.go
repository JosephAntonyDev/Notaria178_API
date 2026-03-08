package app

import (
	"context"

	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/repository"
)

type SearchWorksUseCase struct {
	repo repository.WorkRepository
}

func NewSearchWorksUseCase(r repository.WorkRepository) *SearchWorksUseCase {
	return &SearchWorksUseCase{repo: r}
}

// Execute aplica aislamiento de datos basado en el rol del usuario
// antes de delegar la consulta al repositorio.
func (uc *SearchWorksUseCase) Execute(
	ctx context.Context,
	userRole string,
	userID string,
	userBranchID string,
	filters repository.WorkFilters,
) ([]WorkDTO, error) {

	// ── Data Scoping ────────────────────────────────────────────────────
	switch userRole {
	case "SUPER_ADMIN":
		// El Notario puede ver todo; si envía BranchID por URL se respeta.
	case "LOCAL_ADMIN", "DATA_ENTRY":
		// Forzar aislamiento por sucursal (ignora lo que venga por URL).
		filters.BranchID = &userBranchID
	case "DRAFTER":
		// Forzar aislamiento doble: solo expedientes donde es creador o colaborador.
		filters.ScopedUserID = &userID
	}

	works, err := uc.repo.List(ctx, filters)
	if err != nil {
		return nil, err
	}

	worksDTO := make([]WorkDTO, 0, len(works))
	for _, w := range works {
		worksDTO = append(worksDTO, ToWorkDTO(w))
	}

	return worksDTO, nil
}
