package sendemail

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

func generateMessageIDFromSender(sender string) string {
	// attempt to extract domain
	parts := strings.Split(sender, "@")
	domain := "local"
	if len(parts) == 2 && parts[1] != "" {
		domain = parts[1]
	}
	return fmt.Sprintf("%s@%s", uuid.NewString(), domain)
}

func formatDateRFC5322(t time.Time) string {
	return t.Format(time.RFC1123Z)
}
