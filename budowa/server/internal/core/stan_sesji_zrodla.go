package core

import (
	"context"

	"danacoconsole/server/internal/dane"
)

// Źródła żywego stanu sesji: osadzenie sesji w środowisku platformy i nadania dostępu okien rozmowy.

// zrodloOsadzenia mówi, w którym środowisku platformy stoi karta sesji, jako wąski port bez dostępu do bazy.
type zrodloOsadzenia interface {
	KodSrodowiska(ctx context.Context, idSesji string) string
}

// zrodloNadan wylicza identyfikatory nadań dostępu okna rozmowy w kolejności
// nadanej przez Operatora. Kolejność i oznaczenie głównego niosą znaczenie,
// więc wykaz nie jest zbiorem nieuporządkowanym.
type zrodloNadan interface {
	NadaniaOkna(ctx context.Context, idOkna string) []string
}

// zrodloKontaSesji odnajduje kontekst konta, do którego należy sesja; rozgłoszenie
// bez żądania nie zna konta z góry i idzie po kontach po kolei (decyzja 34).
type zrodloKontaSesji interface {
	KontekstSesji(ctx context.Context, idSesji string) context.Context
}

// zrodlaObecnosciZBazy wypełnia porty repozytoriami warstwy danych, albo zwraca porty puste dla zestawu pustego.
func zrodlaObecnosciZBazy(zestaw *dane.Zestaw) (zrodloOsadzenia, zrodloNadan, zrodloKontaSesji) {
	if zestaw == nil {
		return nil, nil, nil
	}
	return osadzenieZBazy{zestaw: zestaw}, nadaniaZBazy{zestaw: zestaw}, kontaSesjiZBazy{zestaw: zestaw}
}

// kontaSesjiZBazy szuka sesji pod każdym kontem właściciela; kontekst bez wskazania wraca, gdy żadne konto sesji nie zna.
type kontaSesjiZBazy struct {
	zestaw *dane.Zestaw
}

func (k kontaSesjiZBazy) KontekstSesji(ctx context.Context, idSesji string) context.Context {
	konteksty, err := kontekstyKont(ctx, k.zestaw.KontoWlasciciela)
	if err != nil {
		return ctx
	}
	for _, kontekstKonta := range konteksty {
		if _, err := k.zestaw.Sesje.PoIdentyfikatorze(kontekstKonta, idSesji); err == nil {
			return kontekstKonta
		}
	}
	return ctx
}

// osadzenieZBazy czyta łańcuch srodowisko-karta_sesji-sesja, implementując port zrodloOsadzenia bazą danych.
type osadzenieZBazy struct {
	zestaw *dane.Zestaw
}

// KodSrodowiska idzie od wiersza sesji do karty, a od karty do środowiska, bo karta sesji nie ma odczytu po identyfikatorze. Niepowodzenie odczytu daje kod pusty — sesja trwa niezależnie od tego, czy rdzeń umie ją w tej chwili osadzić w nawigacji.
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
		karty, err := o.zestaw.KartySesji.Lista(ctx, srodowisko.ID, dane.KontoOperatora(ctx))
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

// nadaniaZBazy czyta nadania dostępu jednego okna rozmowy, implementując port zrodloNadan bazą danych.
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
