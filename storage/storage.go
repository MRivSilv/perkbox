package storage

type Entry struct {
	Service  []byte `json:"service"`
	Username []byte `json:"username"`
	Password []byte `json:"password"`
	TwoFAKey []byte `json:"twofa_key,omitempty"`
}

type Storage interface {
	LoadAll() ([]Entry, error)
	SaveAll(entries []Entry) error
	FindByService(service, masterPassword string) ([]Entry, error)
}

func GetStorage() Storage {
	return NewLocalStorage()
}
