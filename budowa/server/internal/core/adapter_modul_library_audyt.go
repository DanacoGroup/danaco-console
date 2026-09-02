// Moduł Library obsługuje dziennik audytu (`library.audit.list`) oraz dwa elementy towarzyszące każdej czynności modułu: odnotowanie zdarzenia w dzienniku i zgłoszenie go nasłuchom zewnętrznym.
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

// sprawcaBiblioteki nazywa sprawcę czynności w dzienniku; rdzeń nie zna tożsamości wołającego na poziomie komendy.
const sprawcaBiblioteki = "Operator"

// granicaZgloszeniaNasluchu: zgłoszenie ma dolecieć albo odpaść, a nie trzymać wątku rdzenia.
const granicaZgloszeniaNasluchu = 10 * time.Second

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

// zglosNasluchom rozsyła zdarzenie do nasłuchów, które zapisały to zdarzenie w wykazie — wykaz jest zgodą Operatora.
// Praca w tle niesie konto zamawiającego (decyzja 34): kontekst tła powstaje z konta żądania.
func (a *adapterBiblioteki) zglosNasluchom(ctx context.Context, zdarzenie shared.LibraryWebhookEvent,
	kodPliku string) {

	if a.repozytorium == nil {
		return
	}
	go func() {
		tlo, przerwij := context.WithTimeout(
			dane.ZKontemOperatora(context.Background(), dane.KontoOperatora(ctx)),
			granicaZgloszeniaNasluchu)
		defer przerwij()

		nasluchy, err := a.repozytorium.Webhooki(tlo, true)
		if err != nil {
			return
		}
		for _, nasluch := range nasluchy {
			if !nasluchObejmuje(nasluch.Zdarzenia, zdarzenie) {
				continue
			}
			a.wyslijZgloszenie(tlo, nasluch, zdarzenie, kodPliku)
		}
	}()
}

func nasluchObejmuje(zdarzenia []string, zdarzenie shared.LibraryWebhookEvent) bool {
	szukane := zdarzenieWebhookaBazy(zdarzenie)
	for _, zapisane := range zdarzenia {
		if zapisane == szukane {
			return true
		}
	}
	return false
}

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

	// Kolumna ostatnie_zgloszenie mówi o dolocie, nie o próbie.
	chwila := time.Now().UTC().Format(formatZnacznikaBazy)
	nasluch.OstatnieZgloszenie = &chwila
	_, _ = a.repozytorium.ZapiszWebhook(ctx, nasluch)
}

func znacznikBibliotekiZChwili(milisekundy int64) string {
	return time.UnixMilli(milisekundy).UTC().Format(formatZnacznikaBazy)
}

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

// wskazanieBiblioteki: pola opcjonalne kontraktu niosą brak, a nie pusty napis.
func wskazanieBiblioteki(wartosc string) *string {
	if strings.TrimSpace(wartosc) == "" {
		return nil
	}
	return &wartosc
}
