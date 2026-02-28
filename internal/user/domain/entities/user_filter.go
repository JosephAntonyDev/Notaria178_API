package entities

type UserFilters struct {
	Limit    int
	Offset   int
	Search   *string
	Status   *string
	Role     *string
	BranchID *string
}