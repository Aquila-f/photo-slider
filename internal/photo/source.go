package photo

import "context"

type Source interface {
	List(ctx context.Context) ([]string, error)
	Read(ctx context.Context, token string) ([]byte, error)
}

type DirSource struct {
	resolver Resolver
	reader   Reader
}

func NewDirSource(resolver Resolver, reader Reader) *DirSource {
	return &DirSource{resolver: resolver, reader: reader}
}

func (s *DirSource) List(ctx context.Context) ([]string, error) {
	return s.resolver.List(ctx)
}

func (s *DirSource) Read(ctx context.Context, token string) ([]byte, error) {
	path, err := s.resolver.Resolve(ctx, token)
	if err != nil {
		return nil, err
	}
	return s.reader.Read(ctx, path)
}
