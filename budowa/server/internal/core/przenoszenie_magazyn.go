package core

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/shared"
)

// magazynKontekstu przechowuje komplet kontekstu okna pod jednym kluczem konfiguracji na poziomie okna, z pamięcią procesu jako buforem przejmującym zapis i odczyt przy awarii trwałości.
type magazynKontekstu struct {
	mu           sync.RWMutex
	okna         map[string]shared.ContextBundle
	zdegradowane map[string]struct{}

	zycie        context.Context
	repozytorium dane.RepozytoriumKonfiguracji
	dziennik     *log.Logger
}

// kluczKompletuKontekstu jest kluczem ustawienia konfiguracji okna, pod którym leży zapisany komplet kontekstu tego okna.
const kluczKompletuKontekstu = "kontekst.komplet"

// nowyMagazynKontekstu zakłada magazyn nad repozytorium konfiguracji. Puste
// repozytorium jest dopuszczalne — magazyn pracuje wtedy wyłącznie na buforze.
func nowyMagazynKontekstu(zycie context.Context, repozytorium dane.RepozytoriumKonfiguracji,
	dziennik *log.Logger) *magazynKontekstu {

	if zycie == nil {
		zycie = context.Background()
	}
	return &magazynKontekstu{
		okna:         map[string]shared.ContextBundle{},
		zdegradowane: map[string]struct{}{},
		zycie:        zycie,
		repozytorium: repozytorium,
		dziennik:     dziennik,
	}
}

// Zapisz kładzie komplet kontekstu okna do bufora pamięci procesu oraz, gdy repozytorium jest dostępne, do trwałości konfiguracji.
func (m *magazynKontekstu) Zapisz(idOkna string, komplet shared.ContextBundle) {
	if m == nil || idOkna == "" {
		return
	}
	m.mu.Lock()
	m.okna[idOkna] = komplet
	m.mu.Unlock()

	if m.repozytorium == nil || m.czyZdegradowane(idOkna) {
		return
	}
	tresc, err := json.Marshal(komplet)
	if err != nil {
		m.zdegraduj(idOkna, "kodowanie kompletu kontekstu", err)
		return
	}
	wartosc := string(tresc)
	ustawienie := dane.Ustawienie{
		Poziom: shared.ConfigScopeWindow, KluczZasiegu: idOkna,
		Klucz: kluczKompletuKontekstu, Wartosc: &wartosc,
		RodzajWartosci: string(konfig.RodzajJSON),
	}
	if err := m.repozytorium.Ustaw(m.zycie, ustawienie); err != nil {
		m.zdegraduj(idOkna, "zapis kompletu kontekstu", err)
	}
}

// Odczytaj zwraca komplet kontekstu okna. Okno bez zapisanego kompletu daje
// komplet pusty, a nie błąd: przenoszenie z okna, które jeszcze niczego nie
// dostało, jest zwyczajnym przypadkiem.
func (m *magazynKontekstu) Odczytaj(idOkna string) shared.ContextBundle {
	if m == nil || idOkna == "" {
		return shared.ContextBundle{}
	}
	if m.repozytorium != nil && !m.czyZdegradowane(idOkna) {
		ustawienie, jest, err := m.repozytorium.Odczytaj(m.zycie,
			shared.ConfigScopeWindow, idOkna, kluczKompletuKontekstu)
		switch {
		case err != nil:
			m.zdegraduj(idOkna, "odczyt kompletu kontekstu", err)
		case !jest:
			return m.bufor(idOkna)
		default:
			return kompletZUstawienia(ustawienie, m.bufor(idOkna))
		}
	}
	return m.bufor(idOkna)
}

// kompletZUstawienia rozpakowuje komplet z wiersza ustawienia. Wiersz
// nieczytelny ustępuje buforowi zamiast zerować kontekst okna.
func kompletZUstawienia(u dane.Ustawienie, zapasowy shared.ContextBundle) shared.ContextBundle {
	if u.Wartosc == nil || *u.Wartosc == "" {
		return zapasowy
	}
	var komplet shared.ContextBundle
	if err := json.Unmarshal([]byte(*u.Wartosc), &komplet); err != nil {
		return zapasowy
	}
	return komplet
}

// bufor zwraca komplet kontekstu okna trzymany w pamięci procesu, niezależnie od stanu trwałości konfiguracji.
func (m *magazynKontekstu) bufor(idOkna string) shared.ContextBundle {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.okna[idOkna]
}

// czyZdegradowane mówi, czy okno pracuje już wyłącznie na buforze pamięci procesu, bez zapisu do trwałości.
func (m *magazynKontekstu) czyZdegradowane(idOkna string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, zdegradowane := m.zdegradowane[idOkna]
	return zdegradowane
}

// zdegraduj przenosi okno na bufor pamięci procesu i zgłasza to raz, do dziennika procesu, zamiast przy każdym zapisie.
func (m *magazynKontekstu) zdegraduj(idOkna, czynnosc string, przyczyna error) {
	m.mu.Lock()
	_, juz := m.zdegradowane[idOkna]
	m.zdegradowane[idOkna] = struct{}{}
	m.mu.Unlock()
	if juz || m.dziennik == nil {
		return
	}
	m.dziennik.Printf("kontekst okna %s schodzi na bufor pamięci (%s): %v", idOkna, czynnosc, przyczyna)
}
