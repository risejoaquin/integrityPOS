package hardware

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"integritypos/internal/models"
)

func formatLine(left, right string, maxWidth int) string {
	spaceCount := maxWidth - len(left) - len(right)
	if spaceCount < 1 {
		spaceCount = 1
	}
	if len(left) > (maxWidth - len(right) - 1) {
		left = left[:maxWidth-len(right)-1]
		spaceCount = 1
	}
	return left + strings.Repeat(" ", spaceCount) + right
}

func PrintTicket(ticketData []models.OrderItem, total float64, orderID int) error {
	var b bytes.Buffer

	// ESC/POS Commands
	const (
		escInit       = "\x1B\x40"
		escAlignLeft  = "\x1B\x61\x00"
		escAlignCenter = "\x1B\x61\x01"
		escAlignRight = "\x1B\x61\x02"
		escBoldOn     = "\x1B\x45\x01"
		escBoldOff    = "\x1B\x45\x00"
		escCut        = "\x1D\x56\x41\x10"
		escDrawer     = "\x1B\x70\x00\x19\xFA"
		maxWidth      = 32
	)

	// Inicializar impresora (\x1B\x40)
	b.WriteString(escInit)

	// Centrar texto, modo negrita para el encabezado
	b.WriteString(escAlignCenter)
	b.WriteString(escBoldOn)
	b.WriteString("INTEGRITY POS\n")
	b.WriteString("SOLIDBIT CORE\n")
	b.WriteString(escBoldOff)
	b.WriteString(fmt.Sprintf("ORDER: %d\n", orderID))

	// Separador de guiones
	b.WriteString(strings.Repeat("-", maxWidth) + "\n")

	// Iterar sobre ticketData
	b.WriteString(escAlignLeft)
	for _, item := range ticketData {
		left := fmt.Sprintf("%dx PRD #%d", item.Quantity, item.ProductID)
		right := fmt.Sprintf("$%.2f", item.Subtotal)
		b.WriteString(formatLine(left, right, maxWidth) + "\n")
	}

	// Separador de guiones
	b.WriteString(strings.Repeat("-", maxWidth) + "\n")

	// Total alineado a la derecha
	b.WriteString(escAlignRight)
	b.WriteString(escBoldOn)
	b.WriteString(fmt.Sprintf("TOTAL: $%.2f\n", total))
	b.WriteString(escBoldOff)

	// Comando para abrir cajón de dinero y corte de papel
	b.WriteString(escDrawer)
	b.WriteString(escCut)

	// Impresión temporal para depuración en consola
	log.Println("INFO: --- TICKET ESC/POS GENERADO ---")
	log.Printf("%q", b.String())

	return nil
}
