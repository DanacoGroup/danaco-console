// Odpowiedzialność pliku: odczyt skrzynki — foldery, wykaz nagłówków
// i pobranie jednego listu w całości wraz z załącznikami; zapis leży
// osobno.
package poczta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/emersion/go-imap/v2"
)

// dlugoscZapowiedzi to ile bajtów treści wystarcza, żeby rozpoznać sprawę,
// bez ściągania całej wiadomości z serwera.
const dlugoscZapowiedzi = 512

// sekcjaZapowiedzi opisuje część pierwszą listu, przyciętą do zapowiedzi;
// sekcja TEXT znaczy w IMAP wszystko po nagłówkach, nie treść.
func sekcjaZapowiedzi() *imap.FetchItemBodySection {
	return &imap.FetchItemBodySection{
		Part:    []int{1},
		Partial: &imap.SectionPartial{Offset: 0, Size: dlugoscZapowiedzi},
		Peek:    true,
	}
}

// Foldery oddaje wszystkie foldery skrzynki — obsługuje mail.folder.list,
// bez zawężenia po przeznaczeniu ani po zawartości.
func (k *Klient) Foldery() ([]string, error) {
	wykaz, err := k.imap.List("", "*", nil).Collect()
	if err != nil {
		return nil, fmt.Errorf("serwer poczty nie oddał wykazu folderów: %w", err)
	}
	nazwy := make([]string, 0, len(wykaz))
	for _, folder := range wykaz {
		nazwy = append(nazwy, folder.Mailbox)
	}
	sort.Strings(nazwy)
	return nazwy, nil
}

// Wykaz oddaje nagłówki listów pasujących do zawężenia, od najnowszego,
// wraz z liczbą wszystkich pasujących — obsługuje mail.message.list.
func (k *Klient) Wykaz(z Zawezenie) ([]Naglowek, int, error) {
	folder := nazwaFolderu(z.Folder, FolderOdebranych)
	if _, err := k.imap.Select(folder, &imap.SelectOptions{ReadOnly: true}).Wait(); err != nil {
		return nil, 0, fmt.Errorf("skrzynka nie ma folderu %q albo nie daje do niego dostępu: %w", folder, err)
	}

	warunki := &imap.SearchCriteria{}
	if fraza := strings.TrimSpace(z.Fraza); fraza != "" {
		// TEXT przeszukuje nagłówki i treść — dokładnie to, co obiecuje pole
		// query kontraktu.
		warunki.Text = append(warunki.Text, fraza)
	}
	if nadawca := strings.TrimSpace(z.Nadawca); nadawca != "" {
		warunki.Header = append(warunki.Header, imap.SearchCriteriaHeaderField{Key: "From", Value: nadawca})
	}
	if !z.Od.IsZero() {
		// SINCE porównuje samą datę wewnętrzną, tak stanowi RFC 3501; godziny
		// z żądania są tracone.
		warunki.Since = z.Od
	}
	if z.TylkoNieprzeczytane {
		warunki.NotFlag = append(warunki.NotFlag, imap.FlagSeen)
	}

	wynik, err := k.imap.UIDSearch(warunki, nil).Wait()
	if err != nil {
		return nil, 0, fmt.Errorf("serwer poczty odmówił szukania w folderze %q: %w", folder, err)
	}
	uidy := wynik.AllUIDs()
	if len(uidy) == 0 {
		return []Naglowek{}, 0, nil
	}
	wszystkich := len(uidy)

	// Od najnowszego: UID rośnie z czasem doręczenia, porządek malejący jest
	// porządkiem od najnowszej.
	sort.Slice(uidy, func(i, j int) bool { return uidy[i] > uidy[j] })
	if z.Granica > 0 && len(uidy) > z.Granica {
		uidy = uidy[:z.Granica]
	}

	wiadomosci, err := k.imap.Fetch(imap.UIDSetNum(uidy...), &imap.FetchOptions{
		UID:           true,
		Flags:         true,
		Envelope:      true,
		InternalDate:  true,
		BodyStructure: &imap.FetchItemBodyStructure{Extended: true},
		BodySection:   []*imap.FetchItemBodySection{sekcjaZapowiedzi()},
	}).Collect()
	if err != nil {
		return nil, 0, fmt.Errorf("serwer poczty nie oddał nagłówków z folderu %q: %w", folder, err)
	}

	naglowki := make([]Naglowek, 0, len(wiadomosci))
	for _, w := range wiadomosci {
		naglowki = append(naglowki, zlozNaglowek(folder, w))
	}
	// FETCH oddaje wiadomości w porządku numerów, nie zamówionych UID-ów —
	// porządek wraca ustawiony.
	sort.SliceStable(naglowki, func(i, j int) bool {
		return naglowki[i].Chwila.After(naglowki[j].Chwila)
	})
	return naglowki, wszystkich, nil
}

// Pobierz ściąga jeden list w całości — treść i bajty załączników; obsługuje
// mail.message.get, krok przeanalizuj ten list.
func (k *Klient) Pobierz(identyfikator string, zZalacznikami bool) (List, error) {
	folder, uid, err := rozbierzIdentyfikator(identyfikator)
	if err != nil {
		return List{}, err
	}
	if _, err := k.imap.Select(folder, &imap.SelectOptions{ReadOnly: true}).Wait(); err != nil {
		return List{}, fmt.Errorf("skrzynka nie ma folderu %q albo nie daje do niego dostępu: %w", folder, err)
	}

	wiadomosci, err := k.imap.Fetch(imap.UIDSetNum(uid), &imap.FetchOptions{
		UID:          true,
		Flags:        true,
		Envelope:     true,
		InternalDate: true,
		// Cała wiadomość jednym kawałkiem; rozbiór MIME idzie w list.go, żeby
		// nie dublować odczytu sekcji.
		BodySection: []*imap.FetchItemBodySection{{Peek: true}},
	}).Collect()
	if err != nil {
		return List{}, fmt.Errorf("serwer poczty nie oddał wiadomości %s: %w", identyfikator, err)
	}
	if len(wiadomosci) == 0 {
		return List{}, fmt.Errorf("w folderze %q nie ma wiadomości o identyfikatorze %s — "+
			"list mógł zostać przeniesiony albo usunięty w kliencie poczty Operatora", folder, identyfikator)
	}

	surowy := wiadomosci[0].FindBodySection(&imap.FetchItemBodySection{Peek: true})
	if surowy == nil {
		surowy = wiadomosci[0].FindBodySection(&imap.FetchItemBodySection{})
	}
	list := List{Naglowek: zlozNaglowek(folder, wiadomosci[0])}
	if len(surowy) == 0 {
		return list, fmt.Errorf("serwer poczty oddał wiadomość %s bez treści", identyfikator)
	}

	tresc, zalaczniki, err := rozbierzList(surowy)
	if err != nil {
		return list, fmt.Errorf("treść wiadomości %s jest nieczytelna: %w", identyfikator, err)
	}
	list.Tresc = tresc
	if zZalacznikami {
		list.Zalaczniki = zalaczniki
	}
	// Nazwy załączników pochodzą z rozebranej treści, nie z BODYSTRUCTURE, bo
	// niesie ona bajty Operatora.
	list.NazwyZalacznikow = nazwyZalacznikow(zalaczniki)
	return list, nil
}

// nazwyZalacznikow wyciąga same nazwy — do pola wykazu kontraktu, bez
// treści bajtów, które niesie osobne pole załączników.
func nazwyZalacznikow(zalaczniki []Zalacznik) []string {
	if len(zalaczniki) == 0 {
		return nil
	}
	nazwy := make([]string, 0, len(zalaczniki))
	for _, z := range zalaczniki {
		nazwy = append(nazwy, z.Nazwa)
	}
	return nazwy
}

// zlozIdentyfikator skleja folder z UID-em, w postaci używanej przez
// Naglowek.Identyfikator wszędzie w tym pakiecie.
func zlozIdentyfikator(folder string, uid imap.UID) string {
	return folder + ":" + strconv.FormatUint(uint64(uid), 10)
}

// rozbierzIdentyfikator rozdziela parę folder:UID po ostatnim dwukropku,
// bo nazwa folderu ma prawo go zawierać, a UID nigdy.
func rozbierzIdentyfikator(identyfikator string) (string, imap.UID, error) {
	identyfikator = strings.TrimSpace(identyfikator)
	granica := strings.LastIndex(identyfikator, ":")
	if granica <= 0 || granica == len(identyfikator)-1 {
		return "", 0, fmt.Errorf("identyfikator wiadomości %q nie ma postaci folder:UID — "+
			"weź go z wyniku mail.message.list, rdzeń go nie zgaduje", identyfikator)
	}
	numer, err := strconv.ParseUint(identyfikator[granica+1:], 10, 32)
	if err != nil {
		return "", 0, fmt.Errorf("identyfikator wiadomości %q nie niesie liczbowego UID-u: %w", identyfikator, err)
	}
	return identyfikator[:granica], imap.UID(numer), nil
}

// nazwaFolderu bierze folder wskazany, a przy jego braku — domyślny dla
// tej rodziny czynności skrzynki pocztowej.
func nazwaFolderu(wskazany, domyslny string) string {
	if s := strings.TrimSpace(wskazany); s != "" {
		return s
	}
	return domyslny
}

// adresy zamienia koperty IMAP na adresy tekstowe; koperta bez adresu
// (początek grupy) jest pomijana, żeby nie wyglądać jak odbiorca.
func adresy(lista []imap.Address) []string {
	wynik := make([]string, 0, len(lista))
	for _, a := range lista {
		if adres := a.Addr(); adres != "" {
			wynik = append(wynik, adres)
		}
	}
	if len(wynik) == 0 {
		return nil
	}
	return wynik
}

// pierwszyAdres oddaje pierwszy adres listy albo pustkę, gdy nadawca jest
// jeden, ale koperta niesie go listą adresów.
func pierwszyAdres(lista []imap.Address) string {
	if a := adresy(lista); len(a) > 0 {
		return a[0]
	}
	return ""
}

// czyNieprzeczytana czyta brak znacznika Seen; odwrotność jest tu
// zamierzona, bo IMAP nie ma znacznika nieprzeczytana.
func czyNieprzeczytana(znaczniki []imap.Flag) bool {
	for _, z := range znaczniki {
		if z == imap.FlagSeen {
			return false
		}
	}
	return true
}

// chwilaListu bierze datę z nagłówka listu, a gdy jej nie ma — chwilę
// doręczenia; kolejność wynika z pytania kontraktu o czas nadania.
func chwilaListu(koperta *imap.Envelope, doreczono time.Time) time.Time {
	if koperta != nil && !koperta.Date.IsZero() {
		return koperta.Date
	}
	return doreczono
}
