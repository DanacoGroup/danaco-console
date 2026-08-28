package protocol

import (
	"encoding/json"
	"fmt"

	"danacoconsole/shared"
)

// rodzajFragmentu rozstrzyga, co niesie fragment strumienia, wyliczeniem
// rodzajów kontraktu; warstwa protokołu żadnego rodzaju nie dopisuje.
type rodzajFragmentu = shared.ChunkKind

// Chunk jest pojedynczym fragmentem strumienia zdarzeń: treścią zdarzenia
// stream.chunk tego kontraktu API.
type Chunk = shared.StreamChunkEvent

// NowyChunk składa fragment dowolnego rodzaju z treścią nietekstową, kodując
// ją do postaci przenoszonej kopertą.
func NowyChunk(rodzaj rodzajFragmentu, okno, wiadomosc string, dane any) (Chunk, error) {
	c := Chunk{WindowId: okno, MessageId: wiadomosc, Kind: rodzaj}
	if dane == nil {
		return c, nil
	}
	surowe, err := json.Marshal(dane)
	if err != nil {
		return Chunk{}, fmt.Errorf("protocol: kodowanie fragmentu %s: %w", rodzaj, err)
	}
	c.Data = surowe
	return c, nil
}

// ChunkTekstu niesie porcję tekstu odpowiedzi modelu jako kolejny fragment
// strumienia dla wskazanego okna.
func ChunkTekstu(okno, wiadomosc, tekst string) Chunk {
	return Chunk{
		WindowId:  okno,
		MessageId: wiadomosc,
		Kind:      shared.ChunkKindText,
		Text:      wskaznikTekstu(tekst),
	}
}

// ChunkWersjiOstatecznej niesie wersję ostateczną odpowiedzi, całą jej treść,
// nie kolejną porcję do dopisania.
func ChunkWersjiOstatecznej(okno, wiadomosc, tresc string) Chunk {
	return Chunk{
		WindowId:  okno,
		MessageId: wiadomosc,
		Kind:      shared.ChunkKindFinal,
		Text:      wskaznikTekstu(tresc),
	}
}

// ChunkBledu niesie błąd techniczny wywołania. Przyczyna jedzie i jako treść
// nietekstowa, i jako tekst dla użytkownika. Fragment kończy strumień, więc
// wywołujący pakuje go kopertą z ostatni równym prawda — sesja i konto
// pozostają przy tym czynne.
func ChunkBledu(okno, wiadomosc string, b Blad) Chunk {
	c, err := NowyChunk(shared.ChunkKindError, okno, wiadomosc, b)
	if err != nil {
		c = Chunk{WindowId: okno, MessageId: wiadomosc, Kind: shared.ChunkKindError}
	}
	c.Text = wskaznikTekstu(b.Message)
	return c
}

// Tresc odczytuje treść tekstową danego fragmentu strumienia; brak tekstu
// daje pusty napis, nigdy błąd.
func Tresc(c Chunk) string {
	return wartoscTekstu(c.Text)
}

// KopertaFragmentu pakuje fragment w kopertę zdarzenia stream.chunk. Numer
// fragmentu i znacznik końca trafiają do pól strumienia koperty, a nie do
// ładunku — tak rozkłada je kontrakt.
func KopertaFragmentu(id, idSesji string, seq int, ostatni bool, c Chunk) (Koperta, error) {
	k, err := NowaKoperta(shared.EventStreamChunk, id, idSesji, c)
	if err != nil {
		return Koperta{}, err
	}
	k.Seq = wskaznik(seq)
	if ostatni {
		k.Done = wskaznik(true)
	}
	return k, nil
}

// FragmentZKoperty odczytuje fragment z odebranej koperty — działanie odwrotne
// do KopertaFragmentu. Numer i znacznik końca zostają w kopercie: czytają je
// Numer i Ostatni.
func FragmentZKoperty(k Koperta) (Chunk, error) {
	var c Chunk
	if err := LadunekDo(k, &c); err != nil {
		return Chunk{}, err
	}
	return c, nil
}
