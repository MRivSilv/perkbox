package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Entry representa una contraseña guardada
type Entry struct {
	Service  string `json:"service"`
	Username string `json:"username"`
	Password []byte `json:"password"` // guardada encriptada
}

// getStoragePath retorna dónde guardar el archivo
func getStoragePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".perkbox.json")
}

// LoadAll carga todas las entradas del archivo
func LoadAll() ([]Entry, error) {
	path := getStoragePath()

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return []Entry{}, nil // si no existe, retornamos lista vacía
	}
	if err != nil {
		return nil, err
	}

	var entries []Entry
	err = json.Unmarshal(data, &entries)
	return entries, err
}

// SaveAll guarda todas las entradas en el archivo
func SaveAll(entries []Entry) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(getStoragePath(), data, 0600)
	// 0600 = solo el dueño puede leer/escribir
}

// FindByService busca entradas por nombre de servicio
func FindByService(service string) ([]Entry, error) {
	all, err := LoadAll()
	if err != nil {
		return nil, err
	}

	var found []Entry
	for _, entry := range all {
		if entry.Service == service {
			found = append(found, entry)
		}
	}
	return found, nil
}
