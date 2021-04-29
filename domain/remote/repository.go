package remote

type Repository interface {
	GetByID(id int64) (*Object, error)
}
