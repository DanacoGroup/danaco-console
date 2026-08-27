// Moduł Library obsługuje klasyfikację wsadową i sugestie porządkujące:
// library.classify.run, library.suggestion.list oraz library.suggestion.apply.
// Klasyfikacja wytwarza sugestie; dopiero zatwierdzenie zmienia zasób.
package core

import (
	"context"
	"encoding/json"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/models"
	"danacoconsole/shared"
)

// granicaKlasyfikacjiModelem — ilu zasobom naraz wolno zająć kanał modelu.
// Klasyfikacja całego repozytorium jednym żądaniem byłaby setkami wywołań
// w jednej komendzie.
const granicaKlasyfikacjiModelem = 25

// granicaTresciDoModelu — ile znaków treści idzie do modelu. Model rozpoznaje
// rodzaj dokumentu z jego początku; wysyłanie całości byłoby kosztem bez zysku.
const granicaTresciDoModelu = 4000

// wzorceTresciWrazliwej rozpoznają dane osobowe i finansowe w treści zasobu.
//
// Rozpoznanie jest wskazaniem do przeglądu, nie blokadą — moduł nie stawia
// blokad, a fałszywe trafienie ma kosztować Operatora jedno kliknięcie, nie
// utratę dostępu do pliku.
var wzorceTresciWrazliwej = []struct {
	nazwa   string
	wzorzec *regexp.Regexp
}{
	{"numer PESEL", regexp.MustCompile(`\b\d{11}\b`)},
	{"numer NIP", regexp.MustCompile(`\b\d{3}-\d{3}-\d{2}-\d{2}\b|\bNIP[: ]*\d{10}\b`)},
	{"numer rachunku", regexp.MustCompile(`\bPL\d{26}\b|\b\d{26}\b`)},
	{"adres poczty", regexp.MustCompile(`\b[\w.+-]+@[\w-]+\.[\w.]{2,}\b`)},
}

// ZKanalamiModelu wpina rejestr kanałów modelu — ten sam, którym jedzie okno
// rozmowy i moduł Roundtable. Bez niego klasyfikacja pracuje samym pomiarem
// rdzenia, co jest mniejszym zakresem, a nie brakiem czynności.
func (a *adapterBiblioteki) ZKanalamiModelu(kanaly *models.Rejestr) *adapterBiblioteki {
	a.kanaly = kanaly
	return a
}

// KlasyfikujWsadowo obsługuje `library.classify.run`: liczy sugestie modelu
// i pomiaru rdzenia dla wskazanych zasobów i zapisuje je do przeglądu.
func (a *adapterBiblioteki) KlasyfikujWsadowo(ctx context.Context,
	z shared.LibraryClassifyRunRequest) (shared.LibraryClassifyRunResponse, error) {

	zasoby, err := a.zasobyZbioru(ctx, z.FileIds, z.CollectionId)
	if err != nil {
		return shared.LibraryClassifyRunResponse{}, err
	}
	tylkoBezEtykiet := z.UntaggedOnly != nil && *z.UntaggedOnly

	// Rozpoznanie duplikatów potrzebuje całego zbioru naraz, więc liczy się raz,
	// przed pętlą po zasobach.
	wedlugSumy := map[string][]string{}
	for _, zasob := range zasoby {
		if zasob.SumaKontrolna == nil || *zasob.SumaKontrolna == "" {
			continue
		}
		wedlugSumy[*zasob.SumaKontrolna] = append(wedlugSumy[*zasob.SumaKontrolna], zasob.Kod)
	}

	sugestie := []shared.LibrarySuggestion{}
	przetworzone := 0
	doModelu := 0

	for _, zasob := range zasoby {
		etykiety, err := a.repozytorium.Etykiety(ctx, zasob.ID)
		if err != nil {
			return shared.LibraryClassifyRunResponse{}, bladBiblioteki(err)
		}
		if tylkoBezEtykiet && len(etykiety) > 0 {
			continue
		}
		przetworzone++
		kolekcje, err := a.repozytorium.KolekcjePliku(ctx, zasob.ID)
		if err != nil {
			return shared.LibraryClassifyRunResponse{}, bladBiblioteki(err)
		}

		propozycje := []dane.SugestiaBiblioteki{}
		if len(etykiety) == 0 && len(kolekcje) == 0 {
			propozycje = append(propozycje, dane.SugestiaBiblioteki{
				PlikKod: zasob.Kod, Rodzaj: "osierocony",
				Uzasadnienie: "zasób nie ma ani jednej etykiety i nie należy do żadnej kolekcji",
			})
		}
		if zasob.SumaKontrolna != nil {
			blizniacze := wedlugSumy[*zasob.SumaKontrolna]
			if len(blizniacze) > 1 {
				oryginal := blizniacze[0]
				if oryginal != zasob.Kod {
					propozycje = append(propozycje, dane.SugestiaBiblioteki{
						PlikKod: zasob.Kod, Rodzaj: "duplikat", Wartosc: &oryginal,
						Uzasadnienie: "suma kontrolna treści jest identyczna z zasobem " + oryginal,
						Pewnosc:      wskazanieLiczbyBiblioteki(100),
					})
				}
			}
		}

		tresc := a.trescDoKlasyfikacji(zasob)
		if powod := trescWrazliwaBiblioteki(tresc); powod != "" {
			propozycje = append(propozycje, dane.SugestiaBiblioteki{
				PlikKod: zasob.Kod, Rodzaj: "wrazliwy",
				Wartosc:      wskazanieBiblioteki("poufne"),
				Uzasadnienie: "treść zawiera " + powod + " — zasób wchodzi do przeglądu",
			})
		}

		etykietyModelu := []string{}
		if doModelu < granicaKlasyfikacjiModelem && tresc != "" {
			etykietyModelu = a.etykietyOdModelu(ctx, z.ModelChannelId, zasob.Nazwa, tresc)
			if len(etykietyModelu) > 0 {
				doModelu++
			}
		}
		if len(etykietyModelu) == 0 && len(etykiety) == 0 && tresc != "" {
			// Model niedostępny nie znaczy braku propozycji: słowa najczęstsze
			// w treści też są kandydatami.
			etykietyModelu = slowaKluczoweTresciBiblioteki(tresc, 3)
		}
		for _, etykieta := range etykietyModelu {
			if zawieraTekstBiblioteki(etykiety, etykieta) {
				continue
			}
			nazwa := etykieta
			propozycje = append(propozycje, dane.SugestiaBiblioteki{
				PlikKod: zasob.Kod, Rodzaj: "etykieta", Wartosc: &nazwa,
				Uzasadnienie: "wskazanie z treści zasobu " + zasob.Nazwa,
			})
		}

		for _, propozycja := range propozycje {
			propozycja.Kod = nowyIdentyfikator(przedrostekSugestiiBiblioteki)
			zapisana, err := a.repozytorium.ZapiszSugestie(ctx, propozycja)
			if err != nil {
				return shared.LibraryClassifyRunResponse{}, bladBiblioteki(err)
			}
			sugestie = append(sugestie, sugestiaKontraktuBiblioteki(zapisana))
		}
	}

	odpowiedz := shared.LibraryClassifyRunResponse{
		Suggestions: sugestie, ProcessedCount: przetworzone,
	}
	if z.Apply != nil && *z.Apply && len(sugestie) > 0 {
		kody := make([]string, 0, len(sugestie))
		for _, sugestia := range sugestie {
			kody = append(kody, sugestia.Id)
		}
		wynik, err := a.RozstrzygnijSugestie(ctx, shared.LibrarySuggestionApplyRequest{
			SuggestionIds: kody, Accept: true,
		})
		if err != nil {
			return shared.LibraryClassifyRunResponse{}, err
		}
		odpowiedz.AppliedCount = wynik.AppliedCount
	}
	a.odnotuj(ctx, shared.LibraryAuditActionChange, nil,
		"klasyfikacja wsadowa: zasobów "+strconv.Itoa(przetworzone)+
			", sugestii "+strconv.Itoa(len(sugestie)))
	return odpowiedz, nil
}

// WykazSugestii obsługuje `library.suggestion.list`: oddaje sugestie
// zapisane dla zasobu, filtrowane rodzajem i ograniczone liczbą.
func (a *adapterBiblioteki) WykazSugestii(ctx context.Context,
	z shared.LibrarySuggestionListRequest) (shared.LibrarySuggestionListResponse, error) {

	rodzaje := make([]string, 0, len(z.Kinds))
	for _, rodzaj := range z.Kinds {
		rodzaje = append(rodzaje, rodzajSugestiiBazy(rodzaj))
	}
	granica := 0
	if z.Limit != nil {
		granica = *z.Limit
	}
	wiersze, lacznie, err := a.repozytorium.Sugestie(ctx, z.FileId, rodzaje, granica)
	if err != nil {
		return shared.LibrarySuggestionListResponse{}, bladBiblioteki(err)
	}
	sugestie := make([]shared.LibrarySuggestion, 0, len(wiersze))
	for _, wiersz := range wiersze {
		sugestie = append(sugestie, sugestiaKontraktuBiblioteki(wiersz))
	}
	return shared.LibrarySuggestionListResponse{Suggestions: sugestie, Total: lacznie}, nil
}

// RozstrzygnijSugestie obsługuje `library.suggestion.apply`: przyjęcie
// wykonuje czynność, którą sugestia opisuje — etykieta zostaje nadana,
// kolekcja przypisana, duplikat zarchiwizowany.
func (a *adapterBiblioteki) RozstrzygnijSugestie(ctx context.Context,
	z shared.LibrarySuggestionApplyRequest) (shared.LibrarySuggestionApplyResponse, error) {

	if len(z.SuggestionIds) == 0 {
		return shared.LibrarySuggestionApplyResponse{}, bladWskazaniaBiblioteki(
			"decyzja bez wskazania sugestii")
	}
	if !z.Accept {
		odrzucone, err := a.repozytorium.RozstrzygnijSugestie(ctx, z.SuggestionIds, false)
		if err != nil {
			return shared.LibrarySuggestionApplyResponse{}, bladBiblioteki(err)
		}
		return shared.LibrarySuggestionApplyResponse{RejectedCount: odrzucone}, nil
	}

	zmienione := map[string]bool{}
	przyjete := 0
	for _, kod := range z.SuggestionIds {
		sugestia, err := a.repozytorium.Sugestia(ctx, kod)
		if err != nil {
			return shared.LibrarySuggestionApplyResponse{}, bladNieznanejSugestii(kod, err)
		}
		if sugestia.Stan != dane.StanSugestiiOczekujaca {
			continue
		}
		if err := a.wykonajSugestie(ctx, sugestia); err != nil {
			return shared.LibrarySuggestionApplyResponse{}, err
		}
		if _, err := a.repozytorium.RozstrzygnijSugestie(ctx, []string{kod}, true); err != nil {
			return shared.LibrarySuggestionApplyResponse{}, bladBiblioteki(err)
		}
		zmienione[sugestia.PlikKod] = true
		przyjete++
	}

	kody := make([]string, 0, len(zmienione))
	for kod := range zmienione {
		kody = append(kody, kod)
	}
	sort.Strings(kody)
	pliki := make([]shared.LibraryFile, 0, len(kody))
	for _, kod := range kody {
		wiersz, err := a.repozytorium.Plik(ctx, kod)
		if err != nil {
			continue
		}
		plik, err := a.zloz(ctx, wiersz)
		if err != nil {
			return shared.LibrarySuggestionApplyResponse{}, err
		}
		pliki = append(pliki, plik)
		a.zglosNasluchom(shared.LibraryWebhookEventFileChanged, kod)
	}
	a.odnotuj(ctx, shared.LibraryAuditActionChange, nil,
		"przyjęcie sugestii: "+strconv.Itoa(przyjete))
	return shared.LibrarySuggestionApplyResponse{AppliedCount: przyjete, Files: pliki}, nil
}

// wykonajSugestie wykonuje czynność opisaną sugestią, dobierając operację
// do jej rodzaju: etykieta, kolekcja, duplikat albo zasób osierocony.
func (a *adapterBiblioteki) wykonajSugestie(ctx context.Context, sugestia dane.SugestiaBiblioteki) error {
	zasob, err := a.plik(ctx, sugestia.PlikKod)
	if err != nil {
		return err
	}
	switch sugestia.Rodzaj {
	case "etykieta", "wrazliwy":
		if sugestia.Wartosc == nil || *sugestia.Wartosc == "" {
			return nil
		}
		etykiety, err := a.repozytorium.Etykiety(ctx, zasob.ID)
		if err != nil {
			return bladBiblioteki(err)
		}
		if zawieraTekstBiblioteki(etykiety, *sugestia.Wartosc) {
			return nil
		}
		if _, err := a.repozytorium.ZapiszEtykieteSlownika(ctx, *sugestia.Wartosc, nil); err != nil {
			return bladBiblioteki(err)
		}
		if _, err := a.repozytorium.UstawEtykiety(ctx, zasob.Kod,
			append(etykiety, *sugestia.Wartosc)); err != nil {
			return bladBiblioteki(err)
		}
	case "kolekcja":
		if sugestia.Wartosc == nil || *sugestia.Wartosc == "" {
			return nil
		}
		if _, err := a.repozytorium.PrzypiszDoKolekcji(ctx, *sugestia.Wartosc,
			[]string{zasob.Kod}); err != nil {
			return bladNieznanejKolekcji(*sugestia.Wartosc, err)
		}
	case "duplikat":
		// Przyjęcie wskazania duplikatu archiwizuje kopię — oryginał zostaje,
		// a archiwum jest odwracalne.
		if _, err := a.repozytorium.UstawStanPlikow(ctx, []string{zasob.Kod},
			dane.StanZasobuZarchiwizowany); err != nil {
			return bladBiblioteki(err)
		}
		a.odnotuj(ctx, shared.LibraryAuditActionArchive, wskazanieBiblioteki(zasob.Kod),
			"archiwizacja duplikatu po przyjęciu sugestii")
	case "osierocony":
		// Wskazanie zasobu osieroconego jest informacją, nie czynnością;
		// przyjęcie tylko zdejmuje sugestię.
		return nil
	}
	return nil
}

// trescDoKlasyfikacji oddaje początek treści tekstowej zasobu, przycięty
// do granicy znaków przyjmowanej przez model.
func (a *adapterBiblioteki) trescDoKlasyfikacji(zasob dane.PlikBiblioteki) string {
	if zasob.TrescOdwolanie == nil || *zasob.TrescOdwolanie == "" {
		return ""
	}
	bajty, err := os.ReadFile(*zasob.TrescOdwolanie)
	if err != nil {
		return ""
	}
	if rodzajPodgladuJestObrazem(zasob) {
		return ""
	}
	tekst := string(bajty)
	runy := []rune(tekst)
	if len(runy) > granicaTresciDoModelu {
		runy = runy[:granicaTresciDoModelu]
	}
	return string(runy)
}

// etykietyOdModelu prosi model o propozycje etykiet dla zasobu i bierze
// z odpowiedzi wyłącznie to, co rozpozna: wykaz JSON albo wiersze rozdzielone
// przecinkiem.
func (a *adapterBiblioteki) etykietyOdModelu(ctx context.Context, kanal *string,
	nazwa, tresc string) []string {

	if a.kanaly == nil || kanal == nil || strings.TrimSpace(*kanal) == "" {
		return nil
	}
	var odpowiedz strings.Builder
	ujscie := models.UjscieFunkcji(func(_ context.Context, fragment models.Fragment) error {
		if fragment.Kind == shared.ChunkKindText {
			odpowiedz.WriteString(models.TrescFragmentu(fragment))
		}
		return nil
	})
	polecenie := "Zaproponuj od jednej do trzech etykiet tematycznych dla dokumentu. " +
		"Odpowiedz wylacznie wykazem JSON, na przyklad [\"umowa\",\"klient-x\"].\n" +
		"Nazwa pliku: " + nazwa + "\nTresc:\n" + tresc
	zapytanie := models.Zapytanie{
		Wiadomosc: "library.classify.run",
		Tresc:     polecenie,
		Kanal:     strings.TrimSpace(*kanal),
	}
	if err := a.kanaly.Wyslij(ctx, zapytanie, ujscie); err != nil {
		return nil
	}
	return etykietyZOdpowiedziModelu(odpowiedz.String())
}

// etykietyZOdpowiedziModelu wyciąga wykaz etykiet z odpowiedzi modelu,
// rozpoznając wykaz JSON albo wartości rozdzielone przecinkiem.
func etykietyZOdpowiedziModelu(odpowiedz string) []string {
	tekst := strings.TrimSpace(odpowiedz)
	if tekst == "" {
		return nil
	}
	if poczatek := strings.Index(tekst, "["); poczatek >= 0 {
		if koniec := strings.LastIndex(tekst, "]"); koniec > poczatek {
			var wykaz []string
			if err := json.Unmarshal([]byte(tekst[poczatek:koniec+1]), &wykaz); err == nil {
				return oczyscEtykietyBiblioteki(wykaz)
			}
		}
	}
	return oczyscEtykietyBiblioteki(strings.Split(tekst, ","))
}

// oczyscEtykietyBiblioteki sprowadza propozycje do postaci etykiety: bez spacji na
// brzegach, bez pustych, najwyżej trzy.
func oczyscEtykietyBiblioteki(propozycje []string) []string {
	wynik := make([]string, 0, 3)
	for _, propozycja := range propozycje {
		etykieta := strings.TrimSpace(strings.Trim(propozycja, "\"'.\n\r\t "))
		if etykieta == "" || len([]rune(etykieta)) > 40 {
			continue
		}
		wynik = append(wynik, etykieta)
		if len(wynik) == 3 {
			break
		}
	}
	return wynik
}

// slowaKluczoweTresciBiblioteki wybiera najczęstsze słowa treści jako
// kandydatów na etykiety, gdy model nie jest dostępny albo nie odpowiedział.
func slowaKluczoweTresciBiblioteki(tresc string, ile int) []string {
	licznik := map[string]int{}
	for slowo := range odciskTekstuBiblioteki([]byte(tresc)) {
		licznik[slowo] = 0
	}
	// Zbiór odcisku nie liczy wystąpień, więc zliczanie idzie po samej treści.
	for _, slowo := range strings.FieldsFunc(strings.ToLower(tresc), func(znak rune) bool {
		return !(znak >= 'a' && znak <= 'z' || znak >= '0' && znak <= '9' || znak > 127)
	}) {
		if _, znane := licznik[slowo]; znane {
			licznik[slowo]++
		}
	}
	slowa := make([]string, 0, len(licznik))
	for slowo := range licznik {
		slowa = append(slowa, slowo)
	}
	sort.SliceStable(slowa, func(pierwsze, drugie int) bool {
		if licznik[slowa[pierwsze]] != licznik[slowa[drugie]] {
			return licznik[slowa[pierwsze]] > licznik[slowa[drugie]]
		}
		return slowa[pierwsze] < slowa[drugie]
	})
	if len(slowa) > ile {
		slowa = slowa[:ile]
	}
	return slowa
}

// trescWrazliwaBiblioteki nazywa rodzaj danych wrażliwych rozpoznanych
// w treści, łącząc nazwy wszystkich dopasowanych wzorców.
func trescWrazliwaBiblioteki(tresc string) string {
	if tresc == "" {
		return ""
	}
	trafienia := []string{}
	for _, wzorzec := range wzorceTresciWrazliwej {
		if wzorzec.wzorzec.MatchString(tresc) {
			trafienia = append(trafienia, wzorzec.nazwa)
		}
	}
	return strings.Join(trafienia, ", ")
}

// zawieraTekstBiblioteki mówi, czy wykaz niesie już tę wartość, porównując
// teksty bez rozróżniania wielkości liter.
func zawieraTekstBiblioteki(wykaz []string, szukany string) bool {
	for _, wartosc := range wykaz {
		if strings.EqualFold(wartosc, szukany) {
			return true
		}
	}
	return false
}

// wskazanieLiczbyBiblioteki oddaje wskaźnik na liczbę — pola opcjonalne
// kontraktu przyjmują wskaźnik, a nie wartość.
func wskazanieLiczbyBiblioteki(wartosc int64) *int64 {
	return &wartosc
}

// sugestiaKontraktuBiblioteki przenosi wiersz sugestii na kontrakt,
// ustawiając pole pewności tylko wtedy, gdy wiersz je niesie.
func sugestiaKontraktuBiblioteki(wiersz dane.SugestiaBiblioteki) shared.LibrarySuggestion {
	sugestia := shared.LibrarySuggestion{
		Id: wiersz.Kod, FileId: wiersz.PlikKod, Kind: rodzajSugestiiKontraktu(wiersz.Rodzaj),
		Value: wiersz.Wartosc, Reason: wiersz.Uzasadnienie,
		CreatedAt: chwilaBazy(wiersz.Utworzono),
	}
	if wiersz.Pewnosc != nil {
		pewnosc := int(*wiersz.Pewnosc)
		sugestia.Confidence = &pewnosc
	}
	return sugestia
}

// rodzajSugestiiBazy i rodzajSugestiiKontraktu przekładają rodzaj sugestii
// (odwzorowanie: `sugestia_biblioteki.rodzaj`).
func rodzajSugestiiBazy(rodzaj shared.LibrarySuggestionKind) string {
	switch rodzaj {
	case shared.LibrarySuggestionKindCollection:
		return "kolekcja"
	case shared.LibrarySuggestionKindDuplicate:
		return "duplikat"
	case shared.LibrarySuggestionKindOrphan:
		return "osierocony"
	case shared.LibrarySuggestionKindSensitive:
		return "wrazliwy"
	default:
		return "etykieta"
	}
}

func rodzajSugestiiKontraktu(rodzaj string) shared.LibrarySuggestionKind {
	switch rodzaj {
	case "kolekcja":
		return shared.LibrarySuggestionKindCollection
	case "duplikat":
		return shared.LibrarySuggestionKindDuplicate
	case "osierocony":
		return shared.LibrarySuggestionKindOrphan
	case "wrazliwy":
		return shared.LibrarySuggestionKindSensitive
	default:
		return shared.LibrarySuggestionKindTag
	}
}
