// Odpowiedzialność pliku: układ sekcji panelu okna — `panel.sections.get`
// i `panel.sections.set`. Rola okna prowadzi osobny plik
// `handlers_role_wykaz.go`: to dwie odpowiedzialności.
package core

import (
	"context"
	"sort"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// SekcjePaneli jest portem układu sekcji panelu okna, obsługującym rodzinę
// komend `panel.sections.*` niezależnie od portu ról.
type SekcjePaneli interface {
	UkladPanelu(ctx context.Context, z shared.PanelSectionsGetRequest) (shared.PanelSectionsGetResponse, error)
	UstawUklad(ctx context.Context, z shared.PanelSectionsSetRequest) (shared.PanelSectionsSetResponse, error)
}

// Rozjazd portu z adapterem zatrzymuje kompilację tutaj, nie na martwej
// komendzie odkrytej dopiero w czasie działania rdzenia.
var _ SekcjePaneli = (*adapterSekcjiPaneli)(nil)

// zarejestrujSekcjePaneli wpina `panel.sections.get` i `panel.sections.set`.
// Obie komendy oddają układ po operacji, więc zmiana nie wymaga osobnego
// zdarzenia.
func zarejestrujSekcjePaneli(r *Rejestr, sp SekcjePaneli) {
	if r == nil || sp == nil {
		return
	}

	r.Zarejestruj(shared.CommandPanelSectionsGet,
		obsluz(func(ctx context.Context, z shared.PanelSectionsGetRequest) (shared.PanelSectionsGetResponse, error) {
			return sp.UkladPanelu(ctx, z)
		}))

	r.Zarejestruj(shared.CommandPanelSectionsSet,
		obsluz(func(ctx context.Context, z shared.PanelSectionsSetRequest) (shared.PanelSectionsSetResponse, error) {
			return sp.UstawUklad(ctx, z)
		}))
}

// adapterSekcjiPaneli wypełnia port SekcjePaneli. Nie rozszerza adaptera okien:
// układ panelu nie jest polem okna, tylko wierszem pary (okno, panel).
type adapterSekcjiPaneli struct {
	uklady dane.RepozytoriumSekcjiPaneli
}

func nowyAdapterSekcjiPaneli(uklady dane.RepozytoriumSekcjiPaneli) *adapterSekcjiPaneli {
	return &adapterSekcjiPaneli{uklady: uklady}
}

// UkladPanelu oddaje układ sekcji panelu w kolejności widoku. Panel nigdy
// nieustawiany oddaje wykaz pusty, nie odmowę: układ domyślny należy do widoku,
// który sekcje rysuje, a nie do rdzenia.
func (a *adapterSekcjiPaneli) UkladPanelu(ctx context.Context,
	z shared.PanelSectionsGetRequest) (shared.PanelSectionsGetResponse, error) {

	okno, panel, err := adresPanelu(z.WindowId, z.PanelId)
	if err != nil {
		return shared.PanelSectionsGetResponse{}, err
	}
	if a.uklady == nil {
		return shared.PanelSectionsGetResponse{}, bladNosnikaPaneli()
	}
	uklad, err := a.uklady.SekcjePanelu(ctx, okno, panel)
	if err != nil {
		return shared.PanelSectionsGetResponse{}, bladOdczytuPanelu(okno, panel, err)
	}
	return shared.PanelSectionsGetResponse{
		WindowId: okno, PanelId: panel, Sections: sekcjeKontraktu(uklad),
	}, nil
}

// UstawUklad zapisuje układ sekcji panelu w całości: podane sekcje wyznaczają
// układ, a czego w żądaniu nie ma, tego po zapisie nie ma.
func (a *adapterSekcjiPaneli) UstawUklad(ctx context.Context,
	z shared.PanelSectionsSetRequest) (shared.PanelSectionsSetResponse, error) {

	okno, panel, err := adresPanelu(z.WindowId, z.PanelId)
	if err != nil {
		return shared.PanelSectionsSetResponse{}, err
	}
	if a.uklady == nil {
		return shared.PanelSectionsSetResponse{}, bladNosnikaPaneli()
	}
	// Pole `sections` jest wymagane: brak pola to brak, nie polecenie.

	// Wykaz pusty (`"sections": []`) zdejmuje układ własny, jest poprawny.
	if z.Sections == nil {
		return shared.PanelSectionsSetResponse{}, bladSekcjiPanelu("żądanie zapisu układu " +
			"panelu " + panel + " okna " + okno + " bez pola `sections`; Operator poda sekcje " +
			"w kolejności docelowej — wykaz pusty zdejmuje układ własny, brak pola nie znaczy nic")
	}
	zadane, err := sekcjeDanych(z.Sections)
	if err != nil {
		return shared.PanelSectionsSetResponse{}, err
	}
	uklad, err := a.uklady.ZapiszSekcje(ctx, okno, panel, zadane)
	if err != nil {
		return shared.PanelSectionsSetResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeInternalError, "sekcje panelu: nie można zapisać układu panelu "+panel+
				" okna "+okno+": "+err.Error()+"; Operator powtórzy zapis, a przy nawrocie sprawdzi "+
				"dziennik serwera — żądanie było poprawne, zawiódł zapis"))
	}
	return shared.PanelSectionsSetResponse{
		WindowId: okno, PanelId: panel, Sections: sekcjeKontraktu(uklad),
	}, nil
}

// adresPanelu sprawdza parę (okno, panel) — klucz układu. Pusta połowa klucza
// wskazywałaby byt, którego tabela sekcji paneli nie zna.
func adresPanelu(idOkna, idPanelu string) (string, string, error) {
	okno := strings.TrimSpace(idOkna)
	panel := strings.TrimSpace(idPanelu)
	if okno == "" {
		return "", "", bladSekcjiPanelu("żądanie układu panelu bez wskazania okna; " +
			"Operator poda `windowId` — układ jest zapisany dla pary (okno, panel)")
	}
	if panel == "" {
		return "", "", bladSekcjiPanelu("żądanie układu bez wskazania panelu; " +
			"Operator poda `panelId` — układ jest zapisany dla pary (okno, panel)")
	}
	return okno, panel, nil
}

// sekcjeDanych przekłada sekcje kontraktu na sekcje danych, ustalając porządek.
// Sortowanie jest stabilne: sekcje o tym samym `order` zachowują kolejność
// z żądania — remis rozstrzygnięty losowo dawałby dwa układy z jednego zapisu.
func sekcjeDanych(sekcje []shared.PanelSection) ([]dane.SekcjaPanelu, error) {
	uporzadkowane := make([]shared.PanelSection, len(sekcje))
	copy(uporzadkowane, sekcje)
	sort.SliceStable(uporzadkowane, func(i, j int) bool {
		return uporzadkowane[i].Order < uporzadkowane[j].Order
	})

	widziane := make(map[string]bool, len(uporzadkowane))
	wynik := make([]dane.SekcjaPanelu, 0, len(uporzadkowane))
	for _, sekcja := range uporzadkowane {
		id := strings.TrimSpace(sekcja.Id)
		if id == "" {
			return nil, bladSekcjiPanelu("sekcja bez identyfikatora; Operator poda " +
				"`id` każdej sekcji — bez niego nie ma czego ustawić w układzie")
		}
		if widziane[id] {
			return nil, bladSekcjiPanelu("sekcja " + id + " podana dwa razy; Operator zostawi " +
				"w `sections` jeden wpis na sekcję — sekcja stoi na widoku w jednym miejscu")
		}
		widziane[id] = true
		wynik = append(wynik, dane.SekcjaPanelu{
			Id:       id,
			Zwinieta: sekcja.Collapsed,
			Zdjeta:   sekcja.Hidden != nil && *sekcja.Hidden,
		})
	}
	return wynik, nil
}

// sekcjeKontraktu przekłada układ warstwy danych na sekcje kontraktu. Pole
// `hidden` wychodzi wyłącznie przy sekcji zdjętej: jest opcjonalne, a „hidden:
// false" przy każdej sekcji opisywałoby regułę, nie wyjątek.
func sekcjeKontraktu(uklad []dane.SekcjaPanelu) []shared.PanelSection {
	sekcje := make([]shared.PanelSection, 0, len(uklad))
	for _, sekcja := range uklad {
		wpis := shared.PanelSection{
			Id:        sekcja.Id,
			Order:     sekcja.Kolejnosc,
			Collapsed: sekcja.Zwinieta,
		}
		if sekcja.Zdjeta {
			zdjeta := true
			wpis.Hidden = &zdjeta
		}
		sekcje = append(sekcje, wpis)
	}
	return sekcje
}

// Odmowy: każda mówi co, dlaczego i czym czytelnik to zmieni.

// bladSekcjiPanelu składa odmowę żądania niezgodnego z kontraktem rodziny,
// wspólną dla całego pliku, wraz z powodem i czynnością naprawy.
func bladSekcjiPanelu(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"sekcje panelu: "+powod))
}

// bladNosnikaPaneli nazywa brak trwałości układu — milczenie oddałoby pusty
// panel jako stan zastany, zamiast nazwać brak repozytorium.
func bladNosnikaPaneli() error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"sekcje panelu: trwałość układu paneli niewpięta, migracja 105 nie ma nośnika; "+
			"Operator uruchomi serwer z bazą danych — układ paneli nie żyje w pamięci"))
}

// bladOdczytuPanelu niesie usterkę odczytu wraz z adresem, którego dotyczy,
// i odsyła Operatora do dziennika rdzenia przy powtórnym nawrocie.
func bladOdczytuPanelu(okno, panel string, err error) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"sekcje panelu: nie można odczytać układu panelu "+panel+" okna "+okno+": "+
			err.Error()+"; Operator powtórzy odczyt, a przy nawrocie sprawdzi dziennik serwera"))
}
