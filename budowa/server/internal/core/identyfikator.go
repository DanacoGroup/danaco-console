package core

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"sync/atomic"
)

// przedrostekZdarzenia odróżnia identyfikator nadany przez rdzeń od
// identyfikatora żądania nadanego przez klienta. Zdarzenie nie odpowiada na
// żadne żądanie, więc niesie własny identyfikator.
const przedrostekZdarzenia = "ev-"

// licznikZdarzen zapewnia rozróżnialność identyfikatorów wytworzonych w tej
// samej milisekundzie przez równoległe wywołania.
var licznikZdarzen atomic.Uint64

// przedrostekWiadomosci znakuje identyfikator wiadomości nadany przez rdzeń,
// gdy klient wysyła treść bez własnego identyfikatora.
const przedrostekWiadomosci = "msg-"

// identyfikatorZdarzenia nadaje identyfikator komunikatu wychodzącego od
// rdzenia.
func identyfikatorZdarzenia() string {
	return nowyIdentyfikator(przedrostekZdarzenia)
}

// identyfikatorWiadomosci nadaje identyfikator wiadomości okna.
func identyfikatorWiadomosci() string {
	return nowyIdentyfikator(przedrostekWiadomosci)
}

// nowyIdentyfikator składa identyfikator z licznika i części losowej. Losowa
// część pochodzi ze źródła kryptograficznego; gdy źródło zawiedzie, wystarcza
// część licznikowa — brak losowości nie ma prawa wstrzymać pracy.
func nowyIdentyfikator(przedrostek string) string {
	kolejny := licznikZdarzen.Add(1)
	losowe := make([]byte, 8)
	if _, err := rand.Read(losowe); err != nil {
		return przedrostek + strconv.FormatUint(kolejny, 36)
	}
	return przedrostek + strconv.FormatUint(kolejny, 36) + "-" + hex.EncodeToString(losowe)
}
