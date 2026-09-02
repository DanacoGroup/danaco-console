// Odpowiedzialność pliku: odsłuch debaty syntezą mowy —
// `roundtable.speech.synthesize` (panel akcji Debate Panel). Mowa idzie
// z programu serwerowego `espeak-ng`.
package core

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

// narzedzieSyntezyMowy — silnik mowy arsenału serwera, deklarowany w
// sondzie zależności `zaleznosci_zewnetrzne.go`.
var narzedzieSyntezyMowy = zewnetrzne.Narzedzie{
	Nazwa: "eSpeak NG", Program: "espeak-ng", Pakiet: "espeak-ng",
}

const (
	// granicaSyntezyMowy — odczyt całej debaty bywa długi, ale nie godzinny;
	// przekroczenie granicy znaczy proces, który utknął.
	granicaSyntezyMowy = 10 * time.Minute
	// glosDomyslnyOdsluchu — głos silnika eSpeak NG użyty tam, gdzie Operator
	// nie wskazał głosu własnego dla uczestnika.
	glosDomyslnyOdsluchu = "pl"
	// czestotliwoscOdsluchu — próbkowanie, w którym silnik eSpeak NG zapisuje
	// nagranie wyjściowe formatu WAV.
	czestotliwoscOdsluchu = 22050
)

// Odsluch składa nagranie z przebiegu debaty i wydaje je jako artefakt
// odsłuchu wskazanej debaty Roundtable.
func (a *adapterDebaty) Odsluch(ctx context.Context,
	z shared.RoundtableSpeechSynthesizeRequest) (shared.RoundtableSpeechSynthesizeResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	if okno == "" {
		return shared.RoundtableSpeechSynthesizeResponse{},
			bladWskazaniaDebaty("komenda bez wskazania okna debaty")
	}
	if a.uruchamiacz == nil {
		return shared.RoundtableSpeechSynthesizeResponse{}, odmowaBrakuUruchamiacza()
	}

	glosy := map[string]string{}
	if len(z.VoiceByParticipant) > 0 {
		if err := json.Unmarshal(z.VoiceByParticipant, &glosy); err != nil {
			return shared.RoundtableSpeechSynthesizeResponse{}, odmowaNieczytelnychGlosow(err)
		}
	}

	uczestnicy, err := a.repozytorium.Uczestnicy(ctx, okno)
	if err != nil {
		return shared.RoundtableSpeechSynthesizeResponse{}, bladDebaty(err)
	}
	podpisy := podpisyUczestnikow(uczestnicy)

	wypowiedzi, err := a.wypowiedziZakresu(ctx, okno, strings.TrimSpace(wartoscTekstu(z.TurnId)))
	if err != nil {
		return shared.RoundtableSpeechSynthesizeResponse{}, err
	}
	zawezenie := strings.TrimSpace(wartoscTekstu(z.ParticipantId))

	katalog, err := os.MkdirTemp("", "roundtable-odsluch-")
	if err != nil {
		return shared.RoundtableSpeechSynthesizeResponse{}, bladDebaty(err)
	}
	defer os.RemoveAll(katalog)

	oknoZasiegu, zasady, obszar := a.zasiegNarzedziDebaty(ctx)
	probki := make([]byte, 0, 1<<20)
	odczytanych := 0
	for numer, wypowiedz := range wypowiedzi {
		if strings.TrimSpace(wypowiedz.Tresc) == "" {
			continue
		}
		if zawezenie != "" && wypowiedz.Uczestnik != zawezenie {
			continue
		}
		glos := glosDomyslnyOdsluchu
		if wskazany := strings.TrimSpace(glosy[wypowiedz.Uczestnik]); wskazany != "" {
			glos = wskazany
		}
		// Podpis mówcy wchodzi do odczytu: odsłuch bez wskazania, kto mówi,
		// jest ciągiem zdań bez autora.
		tresc := podpis(podpisy, wypowiedz.Uczestnik) + ". " + strings.TrimSpace(wypowiedz.Tresc)

		plik := filepath.Join(katalog, "glos-"+itoa(numer)+".wav")
		if _, err := zewnetrzne.Wolaj(ctx, a.uruchamiacz, oknoZasiegu, zasady, obszar,
			narzedzieSyntezyMowy, []string{"-v", glos, "-w", plik, tresc}, katalog,
			granicaSyntezyMowy); err != nil {
			return shared.RoundtableSpeechSynthesizeResponse{}, odmowaSyntezyMowy(err)
		}
		bajty, err := os.ReadFile(plik)
		if err != nil {
			return shared.RoundtableSpeechSynthesizeResponse{}, bladDebaty(err)
		}
		probki = append(probki, probkiZPliku(bajty)...)
		odczytanych++
	}
	if odczytanych == 0 {
		return shared.RoundtableSpeechSynthesizeResponse{}, odmowaOdsluchuBezZapisu(okno)
	}

	nagranie := nagranieZProbek(probki)
	artefakt, err := a.wydajArtefaktDebaty(ctx, okno, rodzajArtefaktuOdsluchu, "wav",
		nagranie, dlugoscNagraniaMs(len(probki)))
	if err != nil {
		return shared.RoundtableSpeechSynthesizeResponse{}, err
	}
	dlugosc := artefakt.DlugoscMs
	return shared.RoundtableSpeechSynthesizeResponse{
		ArtifactId: artefakt.Kod, Uri: odwolanieArtefaktu(artefakt), DurationMs: &dlugosc,
	}, nil
}

// naglowekWav — długość nagłówka RIFF/WAVE dla zapisu bez dodatkowych
// bloków, w bajtach jednego nagłówka.
const naglowekWav = 44

// probkiZPliku wyjmuje z pliku WAV same próbki dźwięku, zdejmując nagłówek
// RIFF/WAVE z każdej jego części.
func probkiZPliku(bajty []byte) []byte {
	if len(bajty) <= naglowekWav {
		return nil
	}
	return bajty[naglowekWav:]
}

// nagranieZProbek zakłada nagłówek RIFF/WAVE nad złożonymi próbkami.
// Zapis jest szesnastobitowy, jednokanałowy — tak zapisuje silnik mowy.
func nagranieZProbek(probki []byte) []byte {
	const bitowNaProbke = 16
	const kanalow = 1
	bajtowNaSekunde := czestotliwoscOdsluchu * kanalow * bitowNaProbke / 8

	nagranie := make([]byte, naglowekWav+len(probki))
	copy(nagranie[0:4], "RIFF")
	binary.LittleEndian.PutUint32(nagranie[4:8], uint32(36+len(probki)))
	copy(nagranie[8:12], "WAVE")
	copy(nagranie[12:16], "fmt ")
	binary.LittleEndian.PutUint32(nagranie[16:20], 16) // długość bloku formatu
	binary.LittleEndian.PutUint16(nagranie[20:22], 1)  // zapis nieskompresowany
	binary.LittleEndian.PutUint16(nagranie[22:24], kanalow)
	binary.LittleEndian.PutUint32(nagranie[24:28], czestotliwoscOdsluchu)
	binary.LittleEndian.PutUint32(nagranie[28:32], uint32(bajtowNaSekunde))
	binary.LittleEndian.PutUint16(nagranie[32:34], kanalow*bitowNaProbke/8)
	binary.LittleEndian.PutUint16(nagranie[34:36], bitowNaProbke)
	copy(nagranie[36:40], "data")
	binary.LittleEndian.PutUint32(nagranie[40:44], uint32(len(probki)))
	copy(nagranie[naglowekWav:], probki)
	return nagranie
}

// dlugoscNagraniaMs przelicza liczbę bajtów próbek na czas odsłuchu
// w milisekundach, wedle formatu WAV.
func dlugoscNagraniaMs(bajtowProbek int) int {
	const bajtowNaProbke = 2 // szesnaście bitów, jeden kanał
	if bajtowProbek <= 0 {
		return 0
	}
	return bajtowProbek * 1000 / (bajtowNaProbke * czestotliwoscOdsluchu)
}
