package todos

import "context"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, userID int64, req CreateTodoRequest) (*Todo, error) {
	return s.repo.Create(ctx, userID, req.Title)
}

func (s *Service) ListByUser(ctx context.Context, userID int64) ([]Todo, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *Service) Update(ctx context.Context, id, userID int64, req UpdateTodoRequest) (*Todo, error) {
	return s.repo.Update(ctx, id, userID, req.Title, req.IsDone)
}

func (s *Service) Delete(ctx context.Context, id, userID int64) error {
	return s.repo.Delete(ctx, id, userID)
}
