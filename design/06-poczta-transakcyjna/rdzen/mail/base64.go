package mail

import (
	"encoding/base64"
	"io"
)

// base64LineLen to długość wiersza zalecana przez RFC 2045. Bez łamania
// wiersz base64 przekracza limit 998 oktetów i część serwerów odrzuca
// wiadomość.
const base64LineLen = 76

// writeBase64Lines zapisuje dane w base64, łamiąc wiersze co 76 znaków.
func writeBase64Lines(w io.Writer, data []byte) error {
	encoded := base64.StdEncoding.EncodeToString(data)
	for len(encoded) > base64LineLen {
		if _, err := io.WriteString(w, encoded[:base64LineLen]+"\r\n"); err != nil {
			return err
		}
		encoded = encoded[base64LineLen:]
	}
	if len(encoded) > 0 {
		if _, err := io.WriteString(w, encoded+"\r\n"); err != nil {
			return err
		}
	}
	return nil
}
