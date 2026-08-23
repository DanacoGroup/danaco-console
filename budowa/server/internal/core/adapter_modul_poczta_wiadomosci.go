// Odpowiedzialność pliku: odczyt skrzynki od strony kontraktu — foldery, wykaz
// nagłówków, pełny list wraz z wciągnięciem załączników do magazynu oraz
// oznaczenie listu. Metody stoją na `*adapterPoczty` z `adapter_modul_poczta.go`.
//
// Załącznik wchodzi do magazynu tą samą drogą, co `design.asset.upload`: bajty
// lądują w magazynie rdzenia pod sumą sha256, a obok powstaje wiersz zasobu
// w tabeli `zasob_design` — taki sam, jaki zakłada wniesienie pliku do Assets
// Panelu. Dzięki temu model dostaje w odpowiedzi `assetId`, którym woła
// `image.*`, `document.text.extract` i pozostałe narzędzia treści. Rdzeń nie ma
// drugiego magazynu bajtów, więc załącznik odłożony osobno byłby plikiem,
// którego żadne narzędzie nie widzi.
//
// Kolumna `zasob_design.okno` niesie tu odwołanie do skrzynki, nie do okna.
// Jest `NOT NULL`, bo zasób Designu należy do okna modułu, a załącznik listu do
// żadnego okna nie należy — należy do skrzynki. Wpisujemy więc „poczta:<kod>":
// to prawda o pochodzeniu zasobu, daje się odfiltrować i nie miesza się
// z zasobami okien. Zmyślenie identyfikatora istniejącego okna byłoby
// wstawieniem cudzych plików do cudzego panelu.
//
// Wciągnięcie jest warunkowe (`fetchAttachments`), bo kosztuje: list
// z wielomegabajtowym skanem odczytany po to, żeby sprawdzić datę, nie ma
// powodu zostawiać po sobie tych megabajtów w katalogu danych.
package core

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/poczta"
	"danacoconsole/shared"
)

// przedrostekZasobuZalacznika znakuje wiersze zasobów powstałe z załączników —
// odróżnia je w tabeli od `zasob-` wnoszonych i generowanych w Designie.
const przedrostekZasobuZalacznika = "zalacznik-"

// domyslnaGranicaWykazu przycina wykaz, gdy żądanie granicy nie podało. Bez
// niej pierwsze pytanie o skrzynkę Operatora ściągałoby nagłówki wszystkich
// listów, jakie w niej leżą.
const domyslnaGranicaWykazu = 50

// Foldery oddaje foldery skrzynki — obsługuje `mail.folder.list`.
func (a *adapterPoczty) Foldery(ctx context.Context,
	z shared.MailFolderListRequest) (shared.MailFolderListResponse, error) {

	klient, _, err := a.polacz(ctx, z.AccountId)
	if err != nil {
		return shared.MailFolderListResponse{}, err
	}
	defer klient.Zamknij()

	foldery, err := klient.Foldery()
	if err != nil {
		return shared.MailFolderListResponse{}, bladPoczty(err)
	}
	return shared.MailFolderListResponse{Folders: foldery}, nil
}

// Wiadomosci oddaje nagłówki listów pasujących do zawężenia — obsługuje
// `mail.message.list`, czyli krok „odnajdź sprawę, o której mówi Operator".
func (a *adapterPoczty) Wiadomosci(ctx context.Context,
	z shared.MailMessageListRequest) (shared.MailMessageListResponse, error) {

	klient, _, err := a.polacz(ctx, z.AccountId)
	if err != nil {
		return shared.MailMessageListResponse{}, err
	}
	defer klient.Zamknij()

	zawezenie := poczta.Zawezenie{
		Folder:              strings.TrimSpace(wartoscLubPustka(z.Folder)),
		Fraza:               strings.TrimSpace(wartoscLubPustka(z.Query)),
		Nadawca:             strings.TrimSpace(wartoscLubPustka(z.From)),
		TylkoNieprzeczytane: z.UnreadOnly != nil && *z.UnreadOnly,
		Granica:             domyslnaGranicaWykazu,
	}
	if z.Since != nil {
		// Kontrakt liczy czas w milisekundach epoki; IMAP porównuje samą datę.
		zawezenie.Od = time.UnixMilli(*z.Since)
	}
	if z.Limit != nil && *z.Limit > 0 {
		zawezenie.Granica = *z.Limit
	}

	naglowki, wszystkich, err := klient.Wykaz(zawezenie)
	if err != nil {
		return shared.MailMessageListResponse{}, bladPoczty(err)
	}
	wiadomosci := make([]shared.MailMessage, 0, len(naglowki))
	for _, n := range naglowki {
		wiadomosci = append(wiadomosci, wiadomoscKontraktu(n))
	}
	return shared.MailMessageListResponse{Messages: wiadomosci, Total: wszystkich}, nil
}

// Wiadomosc pobiera list w całości i — na życzenie — wciąga jego załączniki do
// magazynu. Obsługuje `mail.message.get`.
func (a *adapterPoczty) Wiadomosc(ctx context.Context,
	z shared.MailMessageGetRequest) (shared.MailMessageGetResponse, error) {

	if strings.TrimSpace(z.MessageId) == "" {
		return shared.MailMessageGetResponse{}, bladWskazaniaPoczty(
			"komenda mail.message.get bez wskazania wiadomości")
	}
	klient, skrzynka, err := a.polacz(ctx, z.AccountId)
	if err != nil {
		return shared.MailMessageGetResponse{}, err
	}
	defer klient.Zamknij()

	zZalacznikami := z.FetchAttachments == nil || *z.FetchAttachments
	list, err := klient.Pobierz(z.MessageId, zZalacznikami)
	if err != nil {
		return shared.MailMessageGetResponse{}, bladOdczytuListu(err)
	}

	wiadomosc := wiadomoscKontraktu(list.Naglowek)
	wiadomosc.Body = &list.Tresc

	if !zZalacznikami || len(list.Zalaczniki) == 0 {
		return shared.MailMessageGetResponse{Message: wiadomosc}, nil
	}
	zasoby, err := a.wciagnijZalaczniki(ctx, skrzynka.Kod, list.Zalaczniki)
	if err != nil {
		return shared.MailMessageGetResponse{}, err
	}
	return shared.MailMessageGetResponse{Message: wiadomosc, AttachmentAssetIds: zasoby}, nil
}

// Oznacz zmienia oznaczenia listu — obsługuje `mail.message.flag`.
//
// Żądanie bez ani jednego oznaczenia jest odmową. Kontrakt daje oba pola jako
// opcjonalne, ale komenda bez żadnego z nich niczego nie zmienia, a odpowiedź
// `ok` z niezmienioną wiadomością wyglądałaby jak wykonana czynność.
func (a *adapterPoczty) Oznacz(ctx context.Context,
	z shared.MailMessageFlagRequest) (shared.MailMessageFlagResponse, error) {

	if strings.TrimSpace(z.MessageId) == "" {
		return shared.MailMessageFlagResponse{}, bladWskazaniaPoczty(
			"komenda mail.message.flag bez wskazania wiadomości")
	}
	if z.Unread == nil && z.Flagged == nil {
		return shared.MailMessageFlagResponse{}, bladWskazaniaPoczty(
			"komenda mail.message.flag bez ani jednego oznaczenia — wskaż unread albo flagged")
	}
	klient, _, err := a.polacz(ctx, z.AccountId)
	if err != nil {
		return shared.MailMessageFlagResponse{}, err
	}
	defer klient.Zamknij()

	naglowek, err := klient.Oznacz(z.MessageId, z.Unread, z.Flagged)
	if err != nil {
		return shared.MailMessageFlagResponse{}, bladOdczytuListu(err)
	}
	return shared.MailMessageFlagResponse{Message: wiadomoscKontraktu(naglowek)}, nil
}

// wciagnijZalaczniki odkłada bajty w magazynie i zakłada wiersze zasobów,
// oddając ich identyfikatory kontraktowe — patrz nagłówek pliku.
//
// Kolejność jest zamierzona: najpierw bajty, potem wiersz — ten sam porządek,
// co w `design.asset.upload`. Wiersz wskazujący blob, którego nie ma, byłby
// zasobem, po który model sięgnie i niczego nie znajdzie.
func (a *adapterPoczty) wciagnijZalaczniki(ctx context.Context, kodSkrzynki string,
	zalaczniki []poczta.Zalacznik) ([]string, error) {

	if a.magazyn == nil {
		return nil, bladPoczty(errors.New(
			"magazyn załączników nie jest wpięty — nie ma gdzie odłożyć bajtów"))
	}
	if a.zasoby == nil {
		return nil, bladPoczty(errors.New(
			"magazyn zasobów nie jest wpięty — bajty załącznika nie miałyby jak trafić " +
				"do narzędzi obrazu i dokumentów"))
	}

	identyfikatory := make([]string, 0, len(zalaczniki))
	for _, zalacznik := range zalaczniki {
		if len(zalacznik.Bajty) == 0 {
			// Załącznik pusty pomijamy w ciszy: nie ma czego odkładać, a wiersz
			// zasobu bez bajtów jest dokładnie tym, czego moduł Design się pozbył.
			continue
		}
		suma := sha256.Sum256(zalacznik.Bajty)
		odwolanie, err := a.magazyn.Zapisz(zalacznik.Bajty, hex.EncodeToString(suma[:]))
		if err != nil {
			return nil, bladPoczty(err)
		}
		nazwa := zalacznik.Nazwa
		format := zalacznik.TypTresci
		zapisany, err := a.zasoby.ZapiszZasob(ctx, dane.ZasobDesignu{
			Kod:  nowyIdentyfikator(przedrostekZasobuZalacznika),
			Okno: podkatalogPoczty + ":" + kodSkrzynki,
			// Rodzaj `document`, nie `image`: załącznik listu bywa PDF-em,
			// arkuszem albo obrazem, a zgadywanie po rozszerzeniu byłoby cechą
			// zmyśloną. Model i tak rozstrzyga po `format`, który niesie
			// typ MIME wprost z nagłówka części listu.
			Rodzaj: "document",
			Nazwa:  &nazwa,
			Format: &format,
			URI:    &odwolanie,
		})
		if err != nil {
			return nil, bladPoczty(err)
		}
		identyfikatory = append(identyfikatory, zapisany.Kod)
	}
	return identyfikatory, nil
}

// bladOdczytuListu nazywa niepowodzenie czynności na konkretnym liście.
//
// To brak wiadomości, odróżniony od braku skrzynki (`wybierzSkrzynke`)
// i od braku łączności (`polacz`). List
// przeniesiony albo skasowany w kliencie poczty Operatora między wykazem
// a odczytem jest normalnym stanem cudzej skrzynki, więc kod jest `not_found`,
// a nie awarią rdzenia.
func bladOdczytuListu(przyczyna error) error {
	return protocolBladPoczty(shared.ErrorCodeNotFound, przyczyna.Error())
}
