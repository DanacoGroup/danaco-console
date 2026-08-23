// Moduł Library — reguły repozytorium: `library.rule.set`, `library.rule.list`,
// `library.rule.remove`.
//
// Reguła jest warunkiem wraz z tym, co ma się stać z zasobem, który go spełnia.
// Trzy rodzaje różnią się chwilą zastosowania, nie kształtem: kolekcja
// inteligentna przelicza zawartość na żądanie, reguła napływu stosuje się przy
// wejściu zasobu do repozytorium, folder obserwowany wciąga pliki spod ścieżki.
//
// Przeliczenie jest tutaj, a nie w warstwie danych, bo to rdzeń zna znaczenie
// członów warunku. Warunek ma postać:
//
//	{"sourceModuleId": "studio", "tags": ["poufne"], "mimeType": "application/pdf",
//	 "from": "2026-01-01", "to": "2026-12-31", "query": "umowa"}
//
// Każdy człon jest opcjonalny, a człony łączą się koniunkcyjnie — tak samo jak
// filtry fasetowe Library Explorera, bo to ten sam sposób zawężania widziany
// z drugiej strony.
//
// Przeliczenie NIE opróżnia kolekcji: dokłada zasoby spełniające warunek ze
// znacznikiem pochodzenia `regula`. Zdjęcie tych zasobów ma własne żądanie
// (`library.rule.remove` z `detachFiles`) — regułą wolno dodać, a zabrać
// wyłącznie na wyraźne polecenie.
package core

import (
	"context"
	"encoding/json"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// warunekRegulyBiblioteki jest kopertą warunku reguły — dokładnie tymi polami,
// które rdzeń umie zastosować. Pole nieznane w żądaniu jest pomijane, a nie
// odmawiane: warunek rośnie razem z modułem, a reguła zapisana wcześniej ma
// dalej działać.
type warunekRegulyBiblioteki struct {
	ModulZrodlowy string   `json:"sourceModuleId"`
	Etykiety      []string `json:"tags"`
	MimeType      string   `json:"mimeType"`
	Od            string   `json:"from"`
	Do            string   `json:"to"`
	Fraza         string   `json:"query"`
	ProjektID     string   `json:"projectId"`
}

// UstawRegule obsługuje `library.rule.set`.
func (a *adapterBiblioteki) UstawRegule(ctx context.Context,
	z shared.LibraryRuleSetRequest) (shared.LibraryRuleSetResponse, error) {

	if strings.TrimSpace(z.Rule.Name) == "" {
		return shared.LibraryRuleSetResponse{}, bladWskazaniaBiblioteki("reguła bez nazwy")
	}
	rodzaj := rodzajRegulyBazy(z.Rule.Kind)
	if rodzaj == "obserwacja" && strings.TrimSpace(wartoscTekstu(z.Rule.WatchPath)) == "" {
		return shared.LibraryRuleSetResponse{}, bladWskazaniaBiblioteki(
			"folder obserwowany bez ścieżki nie ma czego obserwować")
	}
	if rodzaj != "obserwacja" && strings.TrimSpace(wartoscTekstu(z.Rule.TargetCollectionId)) == "" {
		return shared.LibraryRuleSetResponse{}, bladWskazaniaBiblioteki(
			"reguła rodzaju „" + string(z.Rule.Kind) + "” bez kolekcji docelowej nie ma dokąd " +
				"przypisać zasobu")
	}
	warunek, err := trescWarunkuBiblioteki(z.Rule.Condition)
	if err != nil {
		return shared.LibraryRuleSetResponse{}, err
	}

	kod := strings.TrimSpace(z.Rule.Id)
	if kod == "" {
		kod = nowyIdentyfikator(przedrostekRegulyBiblioteki)
	}
	zapisana, err := a.repozytorium.ZapiszRegule(ctx, dane.RegulaBiblioteki{
		Kod: kod, Rodzaj: rodzaj, Nazwa: z.Rule.Name, Warunek: warunek,
		KolekcjaDocelowaKod: z.Rule.TargetCollectionId, SciezkaObserwowana: z.Rule.WatchPath,
		Czynna: z.Rule.Enabled,
	})
	if err != nil {
		return shared.LibraryRuleSetResponse{}, bladBiblioteki(err)
	}

	odpowiedz := shared.LibraryRuleSetResponse{Rule: regulaKontraktuBiblioteki(zapisana)}
	if z.RunNow != nil && *z.RunNow {
		trafione, err := a.przeliczRegule(ctx, zapisana)
		if err != nil {
			return shared.LibraryRuleSetResponse{}, err
		}
		odpowiedz.MatchedFiles = &trafione
		// Ponowny odczyt niesie czas przeliczenia zapisany przy regule — bez
		// niego odpowiedź mówiłaby o regule sprzed własnego przebiegu.
		odswiezona, err := a.repozytorium.Regula(ctx, zapisana.Kod)
		if err == nil {
			odpowiedz.Rule = regulaKontraktuBiblioteki(odswiezona)
		}
	}
	a.odnotuj(ctx, shared.LibraryAuditActionChange, nil, "zapis reguły "+z.Rule.Name)
	return odpowiedz, nil
}

// WykazRegul obsługuje `library.rule.list`.
func (a *adapterBiblioteki) WykazRegul(ctx context.Context,
	z shared.LibraryRuleListRequest) (shared.LibraryRuleListResponse, error) {

	var rodzaj *string
	if z.Kind != nil {
		wartosc := rodzajRegulyBazy(*z.Kind)
		rodzaj = &wartosc
	}
	tylkoCzynne := z.EnabledOnly != nil && *z.EnabledOnly
	wiersze, err := a.repozytorium.Reguly(ctx, rodzaj, tylkoCzynne)
	if err != nil {
		return shared.LibraryRuleListResponse{}, bladBiblioteki(err)
	}
	reguly := make([]shared.LibraryRule, 0, len(wiersze))
	for _, wiersz := range wiersze {
		reguly = append(reguly, regulaKontraktuBiblioteki(wiersz))
	}
	return shared.LibraryRuleListResponse{Rules: reguly, Total: len(reguly)}, nil
}

// UsunRegule obsługuje `library.rule.remove`.
//
// Zasoby przypisane regułą zostają w kolekcji, dopóki żądanie nie powie inaczej:
// usunięcie reguły opróżniające kolekcję po cichu zabrałoby Operatorowi
// zawartość, o którą nie prosił.
func (a *adapterBiblioteki) UsunRegule(ctx context.Context,
	z shared.LibraryRuleRemoveRequest) (shared.LibraryRuleRemoveResponse, error) {

	kod := strings.TrimSpace(z.RuleId)
	if kod == "" {
		return shared.LibraryRuleRemoveResponse{}, bladWskazaniaBiblioteki("komenda bez wskazania reguły")
	}
	regula, err := a.repozytorium.Regula(ctx, kod)
	if err != nil {
		return shared.LibraryRuleRemoveResponse{}, bladNieznanejReguly(kod, err)
	}

	zdjete := 0
	if z.DetachFiles != nil && *z.DetachFiles && regula.KolekcjaDocelowaKod != nil {
		liczba, err := a.repozytorium.OdepnijPrzypisaniaReguly(ctx, *regula.KolekcjaDocelowaKod)
		if err != nil {
			return shared.LibraryRuleRemoveResponse{}, bladBiblioteki(err)
		}
		zdjete = liczba
	}
	usunieta, err := a.repozytorium.UsunRegule(ctx, kod)
	if err != nil {
		return shared.LibraryRuleRemoveResponse{}, bladBiblioteki(err)
	}
	a.odnotuj(ctx, shared.LibraryAuditActionChange, nil, "usunięcie reguły "+regula.Nazwa)
	return shared.LibraryRuleRemoveResponse{Removed: usunieta, DetachedFiles: zdjete}, nil
}

// przeliczRegule stosuje warunek reguły i przypisuje trafione zasoby do
// kolekcji docelowej. Oddaje liczbę zasobów spełniających warunek.
//
// Folder obserwowany nie jest tu przeliczany: wciąganie plików spod ścieżki jest
// napływem z zewnątrz, a nie przeglądem repozytorium. Reguła obserwacji zapisuje
// się i czeka na napływ, więc jej przeliczenie oddaje zero i tak jest uczciwie —
// zamiast udawać, że coś zrobiono.
func (a *adapterBiblioteki) przeliczRegule(ctx context.Context, regula dane.RegulaBiblioteki) (int, error) {
	if regula.Rodzaj == "obserwacja" || regula.KolekcjaDocelowaKod == nil {
		return 0, nil
	}
	var warunek warunekRegulyBiblioteki
	if err := json.Unmarshal([]byte(regula.Warunek), &warunek); err != nil {
		return 0, bladWskazaniaBiblioteki("warunek reguły " + regula.Nazwa +
			" nie jest czytelnym zapisem JSON: " + err.Error())
	}

	stan := dane.StanZasobuCzynny
	filtr := dane.FiltrPlikow{Etykiety: warunek.Etykiety, Stan: &stan}
	if warunek.Fraza != "" {
		fraza := warunek.Fraza
		filtr.Fraza = &fraza
	}
	if warunek.ProjektID != "" {
		projekt := warunek.ProjektID
		filtr.ProjektID = &projekt
	}
	wiersze, _, err := a.repozytorium.Pliki(ctx, filtr)
	if err != nil {
		return 0, bladBiblioteki(err)
	}

	trafione := make([]string, 0, len(wiersze))
	for _, wiersz := range wiersze {
		if !zasobSpelniaWarunek(wiersz, warunek) {
			continue
		}
		trafione = append(trafione, wiersz.Kod)
	}
	if len(trafione) > 0 {
		if _, err := a.repozytorium.PrzypiszRegula(ctx, *regula.KolekcjaDocelowaKod, trafione); err != nil {
			return 0, bladBiblioteki(err)
		}
		a.zglosNasluchom(shared.LibraryWebhookEventRuleFired, trafione[0])
	}
	// Czas przeliczenia zapisuje się przy regule, żeby wykaz reguł mówił, kiedy
	// reguła ostatnio pracowała — inaczej „czynna" znaczyłoby tylko „włączona".
	regula.OstatniePrzeliczenie = wskazanieBiblioteki(terazZnacznikBiblioteki())
	if _, err := a.repozytorium.ZapiszRegule(ctx, regula); err != nil {
		return 0, bladBiblioteki(err)
	}
	return len(trafione), nil
}

// zasobSpelniaWarunek sprawdza człony warunku, których nie umie zawęzić filtr
// warstwy danych: rodzaj treści, moduł wytwórcy i zakres dat.
func zasobSpelniaWarunek(zasob dane.PlikBiblioteki, warunek warunekRegulyBiblioteki) bool {
	if warunek.MimeType != "" {
		if zasob.MimeType == nil || !strings.EqualFold(*zasob.MimeType, warunek.MimeType) {
			return false
		}
	}
	if warunek.ModulZrodlowy != "" {
		if zasob.ModulZrodlowyID == nil || *zasob.ModulZrodlowyID != warunek.ModulZrodlowy {
			return false
		}
	}
	// Zakres dat porównuje się po znaczniku bazy: zapis ISO 8601 jest
	// porównywalny leksykograficznie, więc data podana samym dniem („2026-01-01")
	// działa jako granica bez przekładu na czas.
	if warunek.Od != "" && zasob.Utworzono < warunek.Od {
		return false
	}
	if warunek.Do != "" {
		// Granica górna obejmuje cały wskazany dzień: znacznik zasobu przycina
		// się do długości granicy, więc „2026-12-31" nie odcina zasobu
		// powstałego tego dnia o dwunastej.
		granica := zasob.Utworzono
		if len(granica) > len(warunek.Do) {
			granica = granica[:len(warunek.Do)]
		}
		if granica > warunek.Do {
			return false
		}
	}
	return true
}

// trescWarunkuBiblioteki sprawdza czytelność warunku i oddaje go w postaci do
// zapisu. Warunek pusty jest wartością poprawną — reguła bez zawężenia obejmuje
// wszystko i to jest wybór Operatora, nie usterka.
func trescWarunkuBiblioteki(warunek json.RawMessage) (string, error) {
	if len(warunek) == 0 || string(warunek) == "null" {
		return "{}", nil
	}
	var probny warunekRegulyBiblioteki
	if err := json.Unmarshal(warunek, &probny); err != nil {
		return "", bladWskazaniaBiblioteki(
			"warunek reguły nie jest czytelnym zapisem JSON: " + err.Error())
	}
	return string(warunek), nil
}

// regulaKontraktuBiblioteki przenosi wiersz reguły na kontrakt.
func regulaKontraktuBiblioteki(wiersz dane.RegulaBiblioteki) shared.LibraryRule {
	regula := shared.LibraryRule{
		Id: wiersz.Kod, Kind: rodzajRegulyKontraktu(wiersz.Rodzaj), Name: wiersz.Nazwa,
		Condition: json.RawMessage(wiersz.Warunek), TargetCollectionId: wiersz.KolekcjaDocelowaKod,
		WatchPath: wiersz.SciezkaObserwowana, Enabled: wiersz.Czynna,
		CreatedAt: chwilaBazy(wiersz.Utworzono),
	}
	if wiersz.OstatniePrzeliczenie != nil && *wiersz.OstatniePrzeliczenie != "" {
		chwila := chwilaBazy(*wiersz.OstatniePrzeliczenie)
		regula.LastRunAt = &chwila
	}
	return regula
}

// rodzajRegulyBazy i rodzajRegulyKontraktu przekładają wyliczenie rodzaju reguły
// (odwzorowanie: `regula_biblioteki.rodzaj`).
func rodzajRegulyBazy(rodzaj shared.LibraryRuleKind) string {
	switch rodzaj {
	case shared.LibraryRuleKindIngest:
		return "naplyw"
	case shared.LibraryRuleKindWatch:
		return "obserwacja"
	default:
		return "kolekcja"
	}
}

func rodzajRegulyKontraktu(rodzaj string) shared.LibraryRuleKind {
	switch rodzaj {
	case "naplyw":
		return shared.LibraryRuleKindIngest
	case "obserwacja":
		return shared.LibraryRuleKindWatch
	default:
		return shared.LibraryRuleKindCollection
	}
}
