package models

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// zbudujZadanie składa żądanie HTTP z wiersza rejestru i zapytania okna.
// Kształt ciała jest domyślny dla większości punktów końcowych rozmowy, a każdy
// dostawca dostraja go parametrem cialo_dodatkowe, bez zmiany w kodzie.
func (k *KanalAPI) zbudujZadanie(ctx context.Context, z Zapytanie) (*http.Request, error) {
	cialo, err := k.zbudujCialo(z)
	if err != nil {
		return nil, err
	}
	zadanie, err := http.NewRequestWithContext(ctx, http.MethodPost,
		k.def.Parametr("base_url"), bytes.NewReader(cialo))
	if err != nil {
		return nil, fmt.Errorf("models: kanał %s: budowa żądania: %w", k.def.Kod, err)
	}
	zadanie.Header.Set("Content-Type", "application/json")
	if k.strumieniowy() {
		zadanie.Header.Set("Accept", "text/event-stream")
	}
	for nazwa, wartosc := range naglowkiWiersza(k.def) {
		zadanie.Header.Set(nazwa, wartosc)
	}
	if err := dolozKlucz(k.def, zadanie); err != nil {
		return nil, err
	}
	return zadanie, nil
}

// zbudujCialo składa treść żądania: model, wiadomości i znacznik strumienia,
// scalone z obiektem cialo_dodatkowe wiersza.
func (k *KanalAPI) zbudujCialo(z Zapytanie) ([]byte, error) {
	wiadomosci := make([]map[string]string, 0, len(z.Historia)+2)
	if prompt := z.Nakladka.PromptSystemowy(); prompt != "" {
		wiadomosci = append(wiadomosci, map[string]string{"role": "system", "content": prompt})
	}
	// Pamięć rozmowy idzie między promptem systemowym a wypowiedzią bieżącą;
	// puste wypowiedzi są pomijane.
	for _, wpis := range z.Historia {
		if strings.TrimSpace(wpis.Tresc) == "" {
			continue
		}
		wiadomosci = append(wiadomosci, map[string]string{"role": wpis.RolaLub(), "content": wpis.Tresc})
	}
	wiadomosci = append(wiadomosci, map[string]string{"role": "user", "content": z.Tresc})

	cialo := map[string]any{
		"model":    z.WybranyModel(k.def),
		"messages": wiadomosci,
		"stream":   k.strumieniowy(),
	}
	for klucz, wartosc := range dodatkiCiala(k.def) {
		cialo[klucz] = wartosc
	}
	surowe, err := json.Marshal(cialo)
	if err != nil {
		return nil, fmt.Errorf("models: kanał %s: kodowanie żądania: %w", k.def.Kod, err)
	}
	return surowe, nil
}

// przedrostekSejfu znakuje odwołanie, które wskazuje wpis sejfu poświadczeń,
// a nie zmienną środowiskową; wartość jest tym samym napisem, którym sejf
// znakuje własny zapis w pakiecie dane.
const przedrostekSejfu = "sejf:"

// czytnikSejfu jest wąskim portem odczytu sejfu poświadczeń, którym kanał API
// rozwiązuje odwołanie sejf:<byt>; typ sejfu mieszka w pakiecie dane, więc
// pakiet models bierze go strukturalnie, przez tę jedną metodę.
type czytnikSejfu interface {
	// Odczytaj zwraca poświadczenie bytu i znacznik, czy wpis istnieje.
	Odczytaj(ctx context.Context, byt string) (string, bool)
}

// sejfPoswiadczen trzyma sejf wpięty przez montaż, jeden na cały proces,
// ponieważ kanały API budowane fabryką nie mają innej wspólnej drogi do jego
// uchwytu. Do wpięcia jest nil, a odwołanie sejf: wtedy odmawia jawnie.
var sejfPoswiadczen czytnikSejfu

// UstawSejfPoswiadczen wpina sejf poświadczeń do kanału API. Woła go montaż raz,
// przy składaniu portów, tym samym sejfem, którym jadą adaptery kont i punktów
// dostępu.
func UstawSejfPoswiadczen(c czytnikSejfu) {
	sejfPoswiadczen = c
}

// dolozKlucz wstawia dane dostępowe do nagłówka żądania. Odwołanie w wierszu
// rejestru jest referencją sejfu albo nazwą zmiennej środowiskowej; brak
// odwołania oznacza punkt końcowy bez uwierzytelnienia.
func dolozKlucz(d Definicja, zadanie *http.Request) error {
	odwolanie := strings.TrimSpace(d.PoswiadczenieOdwolanie)
	if odwolanie == "" {
		return nil
	}
	klucz, err := poswiadczenieZOdwolania(zadanie.Context(), d, odwolanie)
	if err != nil {
		return err
	}
	nazwa := d.ParametrLub("naglowek_klucza", "Authorization")
	przedrostek := "Bearer "
	if wskazany, jest := d.ParametrJest("przedrostek_klucza"); jest {
		przedrostek = wskazany
	}
	zadanie.Header.Set(nazwa, przedrostek+klucz)
	return nil
}

// poswiadczenieZOdwolania rozbiera odwołanie na źródło i pobiera z niego sekret.
// Referencja sejfu prowadzi do sejfu wpiętego przy montażu, każda inna postać
// jest nazwą zmiennej środowiskowej.
func poswiadczenieZOdwolania(ctx context.Context, d Definicja, odwolanie string) (string, error) {
	if byt, jestReferencja := strings.CutPrefix(odwolanie, przedrostekSejfu); jestReferencja {
		sejf := sejfPoswiadczen
		if sejf == nil {
			return "", fmt.Errorf("models: kanał %s: odwołanie %s wskazuje sejf, a sejfu nie wpięto", d.Kod, odwolanie)
		}
		klucz, jest := sejf.Odczytaj(ctx, strings.TrimSpace(byt))
		if !jest || strings.TrimSpace(klucz) == "" {
			return "", fmt.Errorf("models: kanał %s: brak poświadczenia w sejfie pod odwołaniem %s", d.Kod, odwolanie)
		}
		return klucz, nil
	}
	klucz := strings.TrimSpace(os.Getenv(odwolanie))
	if klucz == "" {
		return "", fmt.Errorf("models: kanał %s: brak wartości pod odwołaniem %s", d.Kod, odwolanie)
	}
	return klucz, nil
}

// naglowkiWiersza odczytuje dodatkowe nagłówki HTTP zapisane w wierszu rejestru
// kanału i przygotowane do dołożenia do żądania.
func naglowkiWiersza(d Definicja) map[string]string {
	naglowki := map[string]string{}
	obiekt, jest := d.Parametry["naglowki"].(map[string]any)
	if !jest {
		return naglowki
	}
	for nazwa, wartosc := range obiekt {
		if tekst, jest := wartosc.(string); jest {
			naglowki[nazwa] = tekst
		}
	}
	return naglowki
}

// dodatkiCiala odczytuje obiekt scalany z ciałem żądania, zapisany w wierszu
// rejestru kanału jako parametr cialo_dodatkowe.
func dodatkiCiala(d Definicja) map[string]any {
	obiekt, jest := d.Parametry["cialo_dodatkowe"].(map[string]any)
	if !jest {
		return map[string]any{}
	}
	return obiekt
}
