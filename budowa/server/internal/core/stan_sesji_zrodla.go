package core

import (
	"context"

	"danacoconsole/server/internal/dane"
)

// Źródła żywego stanu sesji, których nie zna ani rejestr nadzorcy, ani pętla:
// osadzenie sesji w środowisku platformy oraz nadania dostępu okien rozmowy.
//
// Oba są wąskimi portami, a nie całym zestawem repozytoriów. Kontrolka powrotu
// do sesji potrzebuje dwóch odpowiedzi — „w którym środowisku ta sesja stoi"
// i „co to okno widzi" — a nie dostępu do bazy. Dzięki temu składacz obecności
// da się sprawdzić bez SQL, a warstwa danych pozostaje po swojej stronie
// granicy.
//
// Brak źródła nie unieważnia odpisu: sesja bez rozpoznanego środowiska i okno
// bez nadań są poprawnym stanem, w którym kontrolka pokazuje mniej.

// zrodloOsadzenia mówi, w którym środowisku platformy stoi karta sesji.
type zrodloOsadzenia interface {
	KodSrodowiska(ctx context.Context, idSesji string) string
}

// zrodloNadan wylicza identyfikatory nadań dostępu okna rozmowy w kolejności
// nadanej przez Operatora. Kolejność i oznaczenie głównego niosą znaczenie,
// więc wykaz nie jest zbiorem nieuporządkowanym.
type zrodloNadan interface {
	NadaniaOkna(ctx context.Context, idOkna string) []string
}

// zrodlaObecnosciZBazy wypełnia oba porty repozytoriami warstwy danych.
func zrodlaObecnosciZBazy(zestaw *dane.Zestaw) (zrodloOsadzenia, zrodloNadan) {
	if zestaw == nil {
		return nil, nil
	}
	return osadzenieZBazy{zestaw: zestaw}, nadaniaZBazy{zestaw: zestaw}
}

// osadzenieZBazy czyta łańcuch `srodowisko → karta_sesji → sesja`.
type osadzenieZBazy struct {
	zestaw *dane.Zestaw
}

// KodSrodowiska idzie od wiersza sesji do karty, a od karty do środowiska.
// Karta sesji nie ma odczytu po identyfikatorze, więc przejście prowadzi przez
// słownik środowisk — jest ich tyle, ile profili widoczności modułów, czyli
// garść wierszy zasianych migracją, a nie zbiór rosnący z pracą Operatora.
//
// Niepowodzenie odczytu daje kod pusty. Sesja trwa niezależnie od tego, czy
// rdzeń umie ją w tej chwili osadzić w nawigacji.
func (o osadzenieZBazy) KodSrodowiska(ctx context.Context, idSesji string) string {
	wiersz, err := o.zestaw.Sesje.PoIdentyfikatorze(ctx, idSesji)
	if err != nil {
		return ""
	}
	srodowiska, err := o.zestaw.Srodowiska.Lista(ctx)
	if err != nil {
		return ""
	}
	for _, srodowisko := range srodowiska {
		karty, err := o.zestaw.KartySesji.Lista(ctx, srodowisko.ID)
		if err != nil {
			continue
		}
		for _, karta := range karty {
			if karta.ID == wiersz.KartaSesjiID {
				return srodowisko.Kod
			}
		}
	}
	return ""
}

// nadaniaZBazy czyta nadania dostępu jednego okna rozmowy.
type nadaniaZBazy struct {
	zestaw *dane.Zestaw
}

// NadaniaOkna zwraca identyfikatory nadań czynnych. Okno bez wiersza i okno
// bez nadań dają wykaz pusty — okno pracuje dalej, tylko niczego nie widzi.
func (n nadaniaZBazy) NadaniaOkna(ctx context.Context, idOkna string) []string {
	wiersz, err := n.zestaw.Okna.PoIdentyfikatorze(ctx, idOkna)
	if err != nil {
		return nil
	}
	nadania, err := n.zestaw.Nadania.ListaOkna(ctx, wiersz.ID, true)
	if err != nil {
		return nil
	}
	identyfikatory := make([]string, 0, len(nadania))
	for _, nadanie := range nadania {
		identyfikatory = append(identyfikatory,
			identyfikatorWiersza(nadanie.IdentyfikatorZewnetrzny, nadanie.ID))
	}
	return identyfikatory
}
