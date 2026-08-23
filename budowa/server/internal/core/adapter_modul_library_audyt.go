// Moduł Library — dziennik audytu (`library.audit.list`) oraz dwa mechanizmy
// towarzyszące każdej czynności modułu: odnotowanie zdarzenia w dzienniku
// i zgłoszenie go nasłuchom zewnętrznym.
//
// Odnotowanie nie może wywrócić czynności. Zasób został przeniesiony do archiwum
// naprawdę — nieudany zapis do dziennika nie cofa tego przeniesienia, a odmowa
// oddana Operatorowi po wykonanej czynności byłaby odmową nieprawdziwą. Dziennik
// jest świadkiem czynności, nie jej warunkiem.
//
// Zgłoszenie nasłuchowi idzie w osobnym wątku z własną granicą czasu: odbiorca
// zewnętrzny bywa wolny albo martwy, a komenda repozytorium nie ma czekać na
// cudzy serwer. Adres jest w konfiguracji Operatora, więc rdzeń go nie
// weryfikuje poza wymogiem, że jest adresem HTTP.
package core

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// sprawcaBiblioteki nazywa sprawcę czynności w dzienniku.
//
// Rdzeń nie zna dziś tożsamości wołającego na poziomie komendy — bramka wiąże
// połączenie z sesją, ale adapter widzi samo żądanie. Sprawcą jest więc Operator
// i to jest prawda: komendę wywołało okno, nie moduł ani model. Gdy komenda
// przyjdzie z pętli wykonawczej modułu, sprawca zmieni się razem z drogą, którą
// przyjdzie — i wtedy będzie to zmiana jednego miejsca.
const sprawcaBiblioteki = "Operator"

// granicaZgloszeniaNasluchu jest krótka z zamysłu: zgłoszenie ma dolecieć albo
// odpaść, a nie trzymać wątku rdzenia.
const granicaZgloszeniaNasluchu = 10 * time.Second

// odnotuj dopisuje zdarzenie do dziennika audytu repozytorium.
func (a *adapterBiblioteki) odnotuj(ctx context.Context, czynnosc shared.LibraryAuditAction,
	kodPliku *string, opis string) {

	if a.repozytorium == nil {
		return
	}
	tresc := opis
	_, _ = a.repozytorium.ZapiszWpisAudytu(ctx, dane.WpisAudytuBiblioteki{
		Kod:      nowyIdentyfikator(przedrostekWpisuAudytuBiblioteki),
		PlikKod:  kodPliku,
		Czynnosc: czynnoscAudytuBazy(czynnosc),
		Sprawca:  sprawcaBiblioteki,
		Opis:     &tresc,
	})
}

// DziennikAudytu obsługuje `library.audit.list`.
func (a *adapterBiblioteki) DziennikAudytu(ctx context.Context,
	z shared.LibraryAuditListRequest) (shared.LibraryAuditListResponse, error) {

	filtr := dane.FiltrAudytuBiblioteki{PlikKod: z.FileId, Sprawca: z.Actor}
	for _, czynnosc := range z.Actions {
		filtr.Czynnosci = append(filtr.Czynnosci, czynnoscAudytuBazy(czynnosc))
	}
	if z.Since != nil {
		od := znacznikBibliotekiZChwili(*z.Since)
		filtr.Od = &od
	}
	if z.Until != nil {
		doChwili := znacznikBibliotekiZChwili(*z.Until)
		filtr.Do = &doChwili
	}
	if z.Limit != nil {
		filtr.Limit = *z.Limit
	}
	if z.Offset != nil {
		filtr.Offset = *z.Offset
	}

	wiersze, lacznie, err := a.repozytorium.WpisyAudytu(ctx, filtr)
	if err != nil {
		return shared.LibraryAuditListResponse{}, bladBiblioteki(err)
	}
	wpisy := make([]shared.LibraryAuditEntry, 0, len(wiersze))
	for _, wiersz := range wiersze {
		wpisy = append(wpisy, shared.LibraryAuditEntry{
			Id: wiersz.Kod, FileId: wiersz.PlikKod, Action: czynnoscAudytuKontraktu(wiersz.Czynnosc),
			Actor: wiersz.Sprawca, Detail: wiersz.Opis, At: chwilaBazy(wiersz.Chwila),
		})
	}
	return shared.LibraryAuditListResponse{Entries: wpisy, Total: lacznie}, nil
}

// zglosNasluchom rozsyła zdarzenie repozytorium do nasłuchów zewnętrznych.
//
// Nasłuch wybiera się po zdarzeniu: nasłuch bez tego zdarzenia w wykazie nie
// dostaje zgłoszenia, bo zapisany wykaz zdarzeń jest zgodą Operatora na to, co
// wychodzi na zewnątrz.
func (a *adapterBiblioteki) zglosNasluchom(zdarzenie shared.LibraryWebhookEvent, kodPliku string) {
	if a.repozytorium == nil {
		return
	}
	// Odczyt idzie w tle razem z wysyłką: wykaz nasłuchów jest zwykle pusty,
	// a gdy nie jest — komenda nie ma czekać ani na bazę, ani na cudzy serwer.
	go func() {
		ctx, przerwij := context.WithTimeout(context.Background(), granicaZgloszeniaNasluchu)
		defer przerwij()

		nasluchy, err := a.repozytorium.Webhooki(ctx, true)
		if err != nil {
			return
		}
		for _, nasluch := range nasluchy {
			if !nasluchObejmuje(nasluch.Zdarzenia, zdarzenie) {
				continue
			}
			a.wyslijZgloszenie(ctx, nasluch, zdarzenie, kodPliku)
		}
	}()
}

// nasluchObejmuje mówi, czy nasłuch prosił o to zdarzenie.
func nasluchObejmuje(zdarzenia []string, zdarzenie shared.LibraryWebhookEvent) bool {
	szukane := zdarzenieWebhookaBazy(zdarzenie)
	for _, zapisane := range zdarzenia {
		if zapisane == szukane {
			return true
		}
	}
	return false
}

// wyslijZgloszenie wysyła jedno zgłoszenie i odnotowuje jego czas przy nasłuchu.
//
// Podpis idzie nagłówkiem HMAC-SHA256 po treści zgłoszenia, gdy nasłuch ma
// sekret: odbiorca ma móc rozstrzygnąć, że zgłoszenie pochodzi z tego rdzenia,
// a nie od kogokolwiek, kto zna adres. Kryptografia jest ze standardowej
// biblioteki Go — żadnego programu z zewnątrz.
func (a *adapterBiblioteki) wyslijZgloszenie(ctx context.Context, nasluch dane.WebhookBiblioteki,
	zdarzenie shared.LibraryWebhookEvent, kodPliku string) {

	tresc, err := json.Marshal(map[string]string{
		"event":  string(zdarzenie),
		"fileId": kodPliku,
		"at":     time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		return
	}
	zadanie, err := http.NewRequestWithContext(ctx, http.MethodPost, nasluch.Adres,
		bytes.NewReader(tresc))
	if err != nil {
		return
	}
	zadanie.Header.Set("Content-Type", "application/json")
	if nasluch.Sekret != nil && *nasluch.Sekret != "" {
		podpis := hmac.New(sha256.New, []byte(*nasluch.Sekret))
		podpis.Write(tresc)
		zadanie.Header.Set("X-Danaco-Podpis", hex.EncodeToString(podpis.Sum(nil)))
	}
	odpowiedz, err := http.DefaultClient.Do(zadanie)
	if err != nil {
		return
	}
	_ = odpowiedz.Body.Close()

	// Czas ostatniego zgłoszenia zapisuje się po locie udanym: kolumna ma mówić
	// „dolatuje", a nie „próbowaliśmy".
	chwila := time.Now().UTC().Format(formatZnacznikaBazy)
	nasluch.OstatnieZgloszenie = &chwila
	_, _ = a.repozytorium.ZapiszWebhook(ctx, nasluch)
}

// znacznikBibliotekiZChwili przekłada milisekundy epoki kontraktu na znacznik
// czasu schematu — odwrotność `chwilaBazy`, potrzebna zawężeniom czasu
// w dzienniku audytu.
func znacznikBibliotekiZChwili(milisekundy int64) string {
	return time.UnixMilli(milisekundy).UTC().Format(formatZnacznikaBazy)
}

// czynnoscAudytuBazy i czynnoscAudytuKontraktu przekładają wyliczenie czynności
// w obie strony (odwzorowanie: `wpis_audytu_biblioteki.czynnosc`).
func czynnoscAudytuBazy(czynnosc shared.LibraryAuditAction) string {
	switch czynnosc {
	case shared.LibraryAuditActionChange:
		return "zmiana"
	case shared.LibraryAuditActionArchive:
		return "archiwizacja"
	case shared.LibraryAuditActionRestore:
		return "przywrocenie"
	case shared.LibraryAuditActionExport:
		return "eksport"
	case shared.LibraryAuditActionDelete:
		return "usuniecie"
	case shared.LibraryAuditActionPreserve:
		return "utrwalenie"
	default:
		return "dostep"
	}
}

func czynnoscAudytuKontraktu(czynnosc string) shared.LibraryAuditAction {
	switch czynnosc {
	case "zmiana":
		return shared.LibraryAuditActionChange
	case "archiwizacja":
		return shared.LibraryAuditActionArchive
	case "przywrocenie":
		return shared.LibraryAuditActionRestore
	case "eksport":
		return shared.LibraryAuditActionExport
	case "usuniecie":
		return shared.LibraryAuditActionDelete
	case "utrwalenie":
		return shared.LibraryAuditActionPreserve
	default:
		return shared.LibraryAuditActionAccess
	}
}

// zdarzenieWebhookaBazy i zdarzenieWebhookaKontraktu przekładają zdarzenie
// nasłuchu w obie strony (odwzorowanie: `zdarzenie_webhooka_biblioteki.zdarzenie`).
func zdarzenieWebhookaBazy(zdarzenie shared.LibraryWebhookEvent) string {
	switch zdarzenie {
	case shared.LibraryWebhookEventFileChanged:
		return "plik_zmieniony"
	case shared.LibraryWebhookEventFileArchived:
		return "plik_zarchiwizowany"
	case shared.LibraryWebhookEventFileRestored:
		return "plik_przywrocony"
	case shared.LibraryWebhookEventVersionAdded:
		return "wersja_dolozona"
	case shared.LibraryWebhookEventRuleFired:
		return "regula_zadzialala"
	default:
		return "plik_dodany"
	}
}

func zdarzenieWebhookaKontraktu(zdarzenie string) shared.LibraryWebhookEvent {
	switch zdarzenie {
	case "plik_zmieniony":
		return shared.LibraryWebhookEventFileChanged
	case "plik_zarchiwizowany":
		return shared.LibraryWebhookEventFileArchived
	case "plik_przywrocony":
		return shared.LibraryWebhookEventFileRestored
	case "wersja_dolozona":
		return shared.LibraryWebhookEventVersionAdded
	case "regula_zadzialala":
		return shared.LibraryWebhookEventRuleFired
	default:
		return shared.LibraryWebhookEventFileAdded
	}
}

// wskazanieBiblioteki oddaje wskaźnik na łańcuch albo nic dla łańcucha pustego —
// pola opcjonalne kontraktu mają nieść brak, a nie pusty napis.
func wskazanieBiblioteki(wartosc string) *string {
	if strings.TrimSpace(wartosc) == "" {
		return nil
	}
	return &wartosc
}
