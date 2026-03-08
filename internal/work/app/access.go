package app

import "github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/entities"

// ─── Helpers de acceso ──────────────────────────────────────────────────────

// canAccessWork verifica si el usuario puede VER el expediente
func canAccessWork(work *entities.Work, reqCtx RequestContext, isCollaborator bool) bool {
	if reqCtx.UserRole == "SUPER_ADMIN" {
		return true
	}
	if reqCtx.UserRole == "LOCAL_ADMIN" && work.BranchID.String() == reqCtx.BranchID {
		return true
	}
	if reqCtx.UserRole == "DATA_ENTRY" && work.BranchID.String() == reqCtx.BranchID {
		return true
	}
	if work.MainDrafterID != nil && work.MainDrafterID.String() == reqCtx.UserID {
		return true
	}
	return isCollaborator
}

// canModifyWork verifica si el usuario puede EDITAR el expediente (info, actos, colaboradores)
func canModifyWork(work *entities.Work, reqCtx RequestContext) bool {
	if reqCtx.UserRole == "SUPER_ADMIN" {
		return true
	}
	if reqCtx.UserRole == "LOCAL_ADMIN" && work.BranchID.String() == reqCtx.BranchID {
		return true
	}
	if work.MainDrafterID != nil && work.MainDrafterID.String() == reqCtx.UserID {
		return true
	}
	return false
}

// isAdminOrNotario verifica si el usuario tiene rol administrativo sobre este expediente
func isAdminOrNotario(reqCtx RequestContext, work *entities.Work) bool {
	if reqCtx.UserRole == "SUPER_ADMIN" {
		return true
	}
	return reqCtx.UserRole == "LOCAL_ADMIN" && work.BranchID.String() == reqCtx.BranchID
}

// isDrafterOrCollaborator verifica si el usuario es proyectista principal o colaborador
func isDrafterOrCollaborator(work *entities.Work, reqCtx RequestContext, isCollaborator bool) bool {
	if work.MainDrafterID != nil && work.MainDrafterID.String() == reqCtx.UserID {
		return true
	}
	return isCollaborator
}
