package models

import (
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Fragment jest jednostką jednolitego strumienia odpowiedzi. Kształt
// bierze wprost z kontraktu przez pakiet protocol — pakiet models nie tworzy
// równoległej struktury komunikatu.
type Fragment = protocol.Chunk

// FragmentTekstu niesie porcję tekstu odpowiedzi modelu, złożoną w fragment strumienia gotowy do wysłania oknu rozmowy.
func FragmentTekstu(z Zapytanie, tekst string) Fragment {
	return protocol.ChunkTekstu(z.Okno(), z.Wiadomosc, tekst)
}

// TrescObrazu jest ładunkiem fragmentu obrazu — treścią pola `data` fragmentu rodzaju `image`. Bajty
// jadą albo wprost przez pole Base64, albo adresem przez pole Adres, zależnie od odpowiedzi dostawcy.
type TrescObrazu struct {
	// TypTresci nazywa rodzaj bajtów (np. „image/png").
	TypTresci string `json:"mimeType,omitempty"`
	// Base64 niesie bajty obrazu w zapisie base64.
	Base64 string `json:"contentBase64,omitempty"`
	// Adres niesie odsyłacz do obrazu, gdy dostawca oddał adres zamiast bajtów.
	Adres string `json:"uri,omitempty"`
	// Prompt niesie polecenie, z którego obraz powstał.
	Prompt string `json:"prompt,omitempty"`
}

// FragmentObrazu niesie treść wizualną odpowiedzi kanału. Pole tekstowe
// fragmentu zostaje puste świadomie: wołacze zbierające odpowiedź zawężają się
// warunkiem na rodzaju `text`, więc obraz przechodzi przez nie nietknięty.
func FragmentObrazu(z Zapytanie, t TrescObrazu) (Fragment, error) {
	return protocol.NowyChunk(shared.ChunkKindImage, z.Okno(), z.Wiadomosc, t)
}

// FragmentProwenancji niesie opis wywołania. Kanał nadaje go jako pierwszy
// fragment strumienia, przed jakąkolwiek treścią odpowiedzi.
func FragmentProwenancji(z Zapytanie, p Prowenancja) (Fragment, error) {
	return protocol.NowyChunk(RodzajProwenancji, z.Okno(), z.Wiadomosc, p)
}

// FragmentKonta niesie metadane konta użytego przez kanał — także przy zmianie
// konta w trakcie sesji, dzięki czemu rotacja jest widoczna, a nie milcząca.
func FragmentKonta(z Zapytanie, k MetadaneKonta) (Fragment, error) {
	return protocol.NowyChunk(RodzajKonta, z.Okno(), z.Wiadomosc, k)
}

// FragmentBledu niesie błąd techniczny kanału. Błąd kończy bieżące wywołanie
// i nic ponadto: sesja, konto i kolejne próby pozostają czynne.
func FragmentBledu(z Zapytanie, b protocol.Blad) Fragment {
	return protocol.ChunkBledu(z.Okno(), z.Wiadomosc, b)
}

// BladKanalu składa błąd kanału. Kod bierze z katalogu kontraktu; ponawialność
// wynika z tego katalogu, nie z decyzji pakietu. Błąd już niosący kod kontraktu
// przechodzi bez podmiany.
func BladKanalu(err error) protocol.Blad {
	return protocol.BladZeZrodla(shared.ErrorCodeChannelUnavailable, err)
}

// TrescFragmentu odczytuje treść tekstową fragmentu strumienia, niezależnie od rodzaju komunikatu, z którego fragment pochodzi.
func TrescFragmentu(f Fragment) string {
	return protocol.Tresc(f)
}
