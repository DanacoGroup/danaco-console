package stdio

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// zakonczenieWiersza oddziela wywołania na wejściu i odpowiedzi na wyjściu.
// Jeden wiersz = jedno wywołanie JSON-RPC.
const zakonczenieWiersza = '\n'

// Obsluguj prowadzi rozmowę z procesem modelu aż do wyczerpania wejścia albo
// zamknięcia kontekstu. Żaden pojedynczy komunikat nieczytelny albo błędny nie
// kończy pracy serwera przedwcześnie.
func Obsluguj(kontekst context.Context, katalog Katalog, wejscie io.Reader, wyjscie io.Writer) error {
	czytnik := bufio.NewReader(wejscie)
	pisarz := json.NewEncoder(wyjscie)
	for {
		if kontekst.Err() != nil {
			return nil
		}
		wiersz, err := czytnik.ReadBytes(zakonczenieWiersza)
		if odpowiedz, jest := rozstrzygnij(kontekst, katalog, wiersz); jest {
			if err := pisarz.Encode(odpowiedz); err != nil {
				return fmt.Errorf("stdio: zapis odpowiedzi: %w", err)
			}
		}
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("stdio: odczyt wywołania: %w", err)
		}
	}
}

// rozstrzygnij zamienia jeden wiersz wejścia na odpowiedź. Zwraca fałsz dla
// wiersza pustego i dla zawiadomienia — obie rzeczy przechodzą bez odpowiedzi.
func rozstrzygnij(kontekst context.Context, katalog Katalog, wiersz []byte) (odpowiedz, bool) {
	wiersz = bytes.TrimSpace(wiersz)
	if len(wiersz) == 0 {
		return odpowiedz{}, false
	}
	var z zadanie
	if err := json.Unmarshal(wiersz, &z); err != nil {
		// Identyfikator nieczytelny w komunikacie daje odpowiedź z identyfikatorem pustym.
		return bledem(json.RawMessage("null"), kodZlyKomunikat,
			"wywołanie niezgodne z JSON-RPC: "+err.Error()), true
	}
	if z.Metoda == "" {
		return bledem(identyfikatorAlboPusty(z.Id), kodZlyKomunikat, "wywołanie bez pola method"), true
	}
	return odpowiedzNaMetode(kontekst, katalog, z)
}

// identyfikatorAlboPusty zwraca identyfikator wywołania albo wartość pustą,
// gdy wywołanie go nie niosło.
func identyfikatorAlboPusty(id json.RawMessage) json.RawMessage {
	if len(id) == 0 {
		return json.RawMessage("null")
	}
	return id
}
