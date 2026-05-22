package cli

import (
	"fmt"
	"os"
	"syscall"

	"github.com/atotto/clipboard"

	"perkbox/crypto"
	"perkbox/storage"

	"golang.org/x/term"
)

// Run despacha el comando correcto
func Run(args []string) {
	switch args[0] {
	case "add":
		cmdAdd()
	case "get":
		if len(args) < 2 {
			fmt.Println("Uso: perkbox get <servicio>")
			os.Exit(1)
		}
		cmdGet(args[1])
	case "list":
		cmdList()
	case "delete":
		if len(args) < 2 {
			fmt.Println("Uso: perkbox delete <servicio>")
			os.Exit(1)
		}
		cmdDelete(args[1])
	default:
		fmt.Printf("Comando desconocido: %s\n", args[0])
		os.Exit(1)
	}
}

func readPassword(prompt string) string {
	fmt.Print(prompt)
	password, _ := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	return string(password)
}

func cmdAdd() {
	var service, username string

	fmt.Print("Servicio (ej: github.com): ")
	fmt.Scanln(&service)

	fmt.Print("Usuario: ")
	fmt.Scanln(&username)

	password := readPassword("Contraseña a guardar: ")
	masterPwd := readPassword("Tu master password: ")

	encrypted, err := crypto.Encrypt(password, masterPwd)
	if err != nil {
		fmt.Println("Error encriptando:", err)
		os.Exit(1)
	}

	entries, err := storage.LoadAll()
	if err != nil {
		fmt.Println("Error cargando datos:", err)
		os.Exit(1)
	}

	entries = append(entries, storage.Entry{
		Service:  service,
		Username: username,
		Password: encrypted,
	})

	if err := storage.SaveAll(entries); err != nil {
		fmt.Println("Error guardando:", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Contraseña para %s guardada!\n", service)
}

func cmdGet(service string) {
	masterPwd := readPassword("Master password: ")

	entries, err := storage.FindByService(service)
	if err != nil || len(entries) == 0 {
		fmt.Printf("No se encontró nada para: %s\n", service)
		return
	}

	for _, e := range entries {
		pwd, err := crypto.Decrypt(e.Password, masterPwd)
		if err != nil {
			fmt.Println("Error: master password incorrecta")
			return
		}
		fmt.Printf("\nServicio:  %s\nUsuario:   %s\n", e.Service, e.Username)
		copiedPass := clipboard.WriteAll(pwd)
		if copiedPass != nil {
			fmt.Printf("No se pudo copiar su contraseña")
			return
		}
		fmt.Printf("Password copiada en su clipboard")
	}
}

func cmdList() {
	entries, err := storage.LoadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if len(entries) == 0 {
		fmt.Println("No hay contraseñas guardadas.")
		return
	}

	fmt.Println("\n=== Tus servicios ===")
	for _, e := range entries {
		fmt.Printf("• %s (%s)\n", e.Service, e.Username)
	}
}

func cmdDelete(service string) {
	entries, err := storage.LoadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var filtered []storage.Entry
	deleted := 0
	for _, e := range entries {
		if e.Service != service {
			filtered = append(filtered, e)
		} else {
			deleted++
		}
	}

	if deleted == 0 {
		fmt.Printf("No se encontró: %s\n", service)
		return
	}

	storage.SaveAll(filtered)
	fmt.Printf("✓ Eliminado %s (%d entrada/s)\n", service, deleted)
}
