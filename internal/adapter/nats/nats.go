package nats

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/carfdev/microservice-go/internal/application"
	"github.com/carfdev/microservice-go/internal/domain"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

// NATS adapter for handling invoice-related messages
type InvoiceNATSAdapter struct {
	nc              *nats.Conn
	invoiceService  *application.InvoiceService
}

func NewInvoiceNATSAdapter(nc *nats.Conn, service *application.InvoiceService) *InvoiceNATSAdapter {
	return &InvoiceNATSAdapter{
		nc:             nc,
		invoiceService: service,
	}
}

func (a *InvoiceNATSAdapter) ListenForMessages() {
// Create
	a.nc.Subscribe("invoice.create", func(msg *nats.Msg) {
		var input domain.Invoice
		if err := json.Unmarshal(msg.Data, &input); err != nil {
			log.Printf("Invalid create payload: %v", err)
			sendErrorResponse(msg, http.StatusBadRequest, "Invalid request payload")
			return
		}

		if input.ID != uuid.Nil {
			log.Println("ID should not be provided during invoice creation")
			sendErrorResponse(msg, http.StatusBadRequest, "Invalid request payload")
			return
		}

		invoice, err := a.invoiceService.CreateInvoice(&input)
		if err != nil {
			log.Printf("Error creating invoice: %v", err)
			sendErrorResponse(msg, http.StatusInternalServerError, "Failed to create invoice")
			return
		}

		a.respond(msg, invoice)
	})

	// Get by ID
	a.nc.Subscribe("invoice.get", func(msg *nats.Msg) {
		id, err := parseID(msg.Data)
		if err != nil {
			log.Printf("Invalid get ID: %v", err)
			sendErrorResponse(msg, http.StatusBadRequest, "Invalid request payload")
			return
		}

		invoice, err := a.invoiceService.GetInvoiceByID(id)
		if err != nil {
			log.Printf("Invoice not found: %v", err)
			sendErrorResponse(msg, http.StatusNotFound, "Invoice not found")
			return
		}

		a.respond(msg, invoice)
	})

	// Get all
	a.nc.Subscribe("invoice.get_all", func(msg *nats.Msg) {
		invoices, err := a.invoiceService.GetAllInvoices()
		if err != nil {
			log.Printf("Error getting invoices: %v", err)
			sendErrorResponse(msg, http.StatusInternalServerError, "Failed to get invoices")
			return
		}

		a.respond(msg, invoices)
	})

	// Update
	a.nc.Subscribe("invoice.update", func(msg *nats.Msg) {
		var input domain.Invoice
		if err := json.Unmarshal(msg.Data, &input); err != nil {
			log.Printf("Invalid update payload: %v", err)
			sendErrorResponse(msg, http.StatusBadRequest, "Invalid request payload")
			return
		}
		

		// Check if the ID is empty (uuid.Nil)
		if input.ID == uuid.Nil {
			log.Println("Update requires a valid ID")
			sendErrorResponse(msg, http.StatusBadRequest, "Update requires a valid ID")
			return
		}

		invoice, err := a.invoiceService.UpdateInvoice(input.ID, &input)
		if err != nil {
			log.Printf("Error updating invoice: %v", err)
			sendErrorResponse(msg, http.StatusInternalServerError, "Failed to update invoice")
			return
		}

		a.respond(msg, invoice)
	})

	// Delete
	a.nc.Subscribe("invoice.delete", func(msg *nats.Msg) {
		id, err := parseID(msg.Data)
		if err != nil {
			log.Printf("Invalid delete ID: %v", err)
			sendErrorResponse(msg, http.StatusBadRequest, "Invalid request payload")
			return
		}

		err = a.invoiceService.DeleteInvoice(id)
		if err != nil {
			log.Printf("Error deleting invoice: %v", err)
			sendErrorResponse(msg, http.StatusInternalServerError, "Failed to delete invoice")
			return
		}

		a.respond(msg, map[string]string{"status": "deleted"})
	})

	log.Println("✔ Listening for NATS messages (CRUD operations)")
}

// Utility: respond with JSON
func (a *InvoiceNATSAdapter) respond(msg *nats.Msg, data any) {
	if msg.Reply == "" {
		return
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshalling response: %v", err)
		sendErrorResponse(msg, http.StatusInternalServerError, "Failed to marshal response")
		return
	}

	if err := a.nc.Publish(msg.Reply, bytes); err != nil {
		log.Printf("Error publishing response: %v", err)
		sendErrorResponse(msg, http.StatusInternalServerError, "Failed to publish response")
	}
}

// Utility: parse ID from payload like {"id": "c8b0a72d-afe5-464e-81d6-d24b2f92ff2d"}
func parseID(data []byte) (uuid.UUID, error) {
	var req struct {
		ID string `json:"id"`
	}

	err := json.Unmarshal(data, &req)
	if err != nil {
		return uuid.Nil, err
	}

	return uuid.Parse(req.ID)
}

// ErrorResponse structure 
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}


// Utility function to send an error response in JSON format
func sendErrorResponse(msg *nats.Msg, status int, message string) {
	errorResponse := ErrorResponse{
		Status:  status,
		Message: message,
	}

	// Convert the error response to JSON
	resp, err := json.Marshal(errorResponse)
	if err != nil {
		log.Printf("Failed to marshal error response: %v", err)
		return
	}

	// Respond to the request with the error message
	if err := msg.Respond(resp); err != nil {
		log.Printf("Failed to send error response: %v", err)
	}
}