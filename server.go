package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	_ "github.com/denisenkom/go-mssqldb"
)

type Message struct {
	Sender  string `json:"sender"`
	Content string `json:"content"`
}

var (
	clients []chan Message
	mutex   sync.Mutex
	db      *sql.DB
)

func main() {
	// Conexión a la base de datos
	var err error
	connString := "sqlserver://@LESLIE:1433?database=CHAT"

	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatalf("Error conectando a la base de datos: %v", err)
	}
	defer db.Close()

	// Verificar conexión
	if err = db.Ping(); err != nil {
		log.Fatalf("Error en la conexión: %v", err)
	}

	http.HandleFunc("/send", handleSendMessage)
	http.HandleFunc("/receive", handleReceiveMessages)

	log.Println("Servidor corriendo en :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleSendMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var msg Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Error procesando el mensaje", http.StatusBadRequest)
		return
	}

	// Validar que el mensaje tenga un remitente
	if msg.Sender == "" {
		http.Error(w, "El remitente no puede estar vacío", http.StatusBadRequest)
		return
	}

	// Guardar mensaje en la base de datos
	usuarioID, err := getOrCreateUser(msg.Sender)
	if err != nil {
		http.Error(w, "Error al guardar el usuario", http.StatusInternalServerError)
		return
	}

	if err = saveMessage(usuarioID, msg.Content); err != nil {
		http.Error(w, "Error al guardar el mensaje", http.StatusInternalServerError)
		return
	}

	distributeMessage(msg)
	w.WriteHeader(http.StatusOK)
}

func handleReceiveMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	ch := make(chan Message, 10)
	registerClient(ch)
	defer unregisterClient(ch)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Transfer-Encoding", "chunked")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "El servidor no soporta streaming", http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	for msg := range ch {
		if err := encoder.Encode(msg); err != nil {
			break
		}
		flusher.Flush()
	}
}

func distributeMessage(msg Message) {
	mutex.Lock()
	defer mutex.Unlock()

	for _, ch := range clients {
		ch <- msg // Enviar el mensaje a todos los clientes registrados
	}
}

func registerClient(ch chan Message) {
	mutex.Lock()
	clients = append(clients, ch)
	mutex.Unlock()
}

func unregisterClient(ch chan Message) {
	mutex.Lock()
	defer mutex.Unlock()

	for i, client := range clients {
		if client == ch {
			clients = append(clients[:i], clients[i+1:]...)
			close(ch)
			break
		}
	}
}

func getOrCreateUser(APODO string) (int, error) {
	var ID int
	err := db.QueryRow("SELECT ID FROM Usuarios WHERE APODO = @p1", APODO).Scan(&ID)
	if err == sql.ErrNoRows {
		// Crear nuevo usuario si no existe
		result, err := db.Exec("INSERT INTO Usuarios (APODO) VALUES (@p1)", APODO)
		if err != nil {
			return 0, err
		}
		insertedID, _ := result.LastInsertId()
		ID = int(insertedID)
	} else if err != nil {
		return 0, err
	}
	return ID, nil
}

func saveMessage(USUARIO_ID int, CONTENIDO string) error {
	_, err := db.Exec("INSERT INTO Mensajes (USUARIO_ID, CONTENIDO) VALUES (@p1, @p2)", USUARIO_ID, CONTENIDO)
	return err
}
