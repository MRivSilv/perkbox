package storage

type Entry struct {
	Service  string `json:"service"`
	Username string `json:"username"`
	Password []byte `json:"password"`
	TwoFAKey []byte `json:"twofa_key,omitempty"`
}

type Storage interface {
	LoadAll() ([]Entry, error)
	SaveAll(entries []Entry) error
	FindByService(service string) ([]Entry, error)
}

func GetStorage() Storage {
	return NewLocalStorage()
}
