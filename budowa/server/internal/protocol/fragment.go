package protocol

import (
	"encoding/json"
	"fmt"

	"danacoconsole/shared"
)

// rodzajFragmentu rozstrzyga, co niesie fragment strumienia. Jeden strumień
// przenosi wszystkie rodzaje treści z każdego kanału, więc odbiorca
// kieruje się wyłącznie tym polem. Katalog rodzajów jest wyliczeniem ChunkKind
// kontraktu — warstwa protokołu żadnego rodzaju nie dopisuje.
type rodzajFragmentu = shared.ChunkKind

// Chunk jest pojedynczym fragmentem strumienia — treścią zdarzenia
// stream.chunk w kształcie kontraktu.
//
// Numeru fragmentu ani znacznika końca tutaj nie ma: kontrakt umieszcza je
// w kopercie (pola seq i done). Wstawia je KopertaFragmentu, odczytują Numer
// i Ostatni.
type Chunk = shared.StreamChunkEvent

// NowyChunk składa fragment dowolnego rodzaju z treścią nietekstową.
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

// ChunkTekstu niesie porcję tekstu odpowiedzi modelu.
func ChunkTekstu(okno, wiadomosc, tekst string) Chunk {
	return Chunk{
		WindowId:  okno,
		MessageId: wiadomosc,
		Kind:      shared.ChunkKindText,
		Text:      wskaznikTekstu(tekst),
	}
}

// ChunkWersjiOstatecznej niesie wersję ostateczną odpowiedzi — całą jej treść,
// nie kolejną porcję.
//
// Odbiorca zastępuje nią tekst złożony z fragmentów tekstowych, zamiast dopisywać
// ją na końcu. Rodzaj `final` odróżnia całość od ciągu dalszego.
//
// Fragment pakuje się kopertą ze znacznikiem końca — jest zdarzeniem domykającym
// strumień.
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

// Tresc odczytuje treść tekstową fragmentu; brak tekstu daje pusty napis.
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
