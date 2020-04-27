package session

type Session interface {
	ID() string
	Get(k interface{}) interface{}
	Set(k interface{}, v interface{})
}

type session struct {
	id   string
	data map[interface{}]interface{}
}

func NewSession() Session {
	id := NewID()
	return &session{id: id, data: make(map[interface{}]interface{})}
}

func (s *session) ID() string {
	return s.id
}
func (s *session) Get(k interface{}) interface{} {
	if item, ok := s.data[k]; ok {
		return item
	}

	return nil
}

func (s *session) Set(k interface{}, v interface{}) {
	s.data[k] = v
}
