// Odpowiedzialność pliku: szkic i wysyłka — dwie komendy, które coś tworzą,
// oraz obowiązkowy ślad tej jednej, której nie da się cofnąć. Bramek
// potwierdzenia nie ma — jest ślad podwójny: wiersz w bazie i kopia w folderze
// wysłanych skrzynki Operatora.
package core

import (
	"context"
	"errors"
	"os"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/poczta"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// ZapiszSzkic odkłada szkic odpowiedzi w folderze szkiców skrzynki Operatora,
// obsługując `mail.draft.save`; nic z niego nie wychodzi w świat.
func (a *adapterPoczty) ZapiszSzkic(ctx context.Context,
	z shared.MailDraftSaveRequest) (shared.MailDraftSaveResponse, error) {

	klient, skrzynka, err := a.polacz(ctx, z.AccountId)
	if err != nil {
		return shared.MailDraftSaveResponse{}, err
	}
	defer klient.Zamknij()

	wychodzacy, err := a.zlozWychodzacy(ctx, skrzynka, wysylanyList{
		Do:            z.To,
		Kopia:         z.Cc,
		Temat:         z.Subject,
		Tresc:         z.Body,
		WOdpowiedziNa: z.InReplyTo,
		Zasoby:        z.AttachmentAssetIds,
	})
	if err != nil {
		return shared.MailDraftSaveResponse{}, err
	}

	naglowek, err := klient.ZapiszSzkic(wychodzacy, strings.TrimSpace(wartoscLubPustka(z.DraftId)))
	if err != nil {
		return shared.MailDraftSaveResponse{}, bladPoczty(err)
	}
	return shared.MailDraftSaveResponse{Draft: wiadomoscKontraktu(naglowek)}, nil
}

// Wyslij nadaje list — obsługuje `mail.send`. Oddaje trzy wartości: odpowiedź
// kontraktu, wiadomość dla zdarzenia `mail.changed` i błąd (uzasadnienie przy
// deklaracji portu w `handlers_poczta.go`).
func (a *adapterPoczty) Wyslij(ctx context.Context,
	z shared.MailSendRequest) (shared.MailSendResponse, shared.MailMessage, error) {

	klient, skrzynka, err := a.polacz(ctx, z.AccountId)
	if err != nil {
		return shared.MailSendResponse{}, shared.MailMessage{}, err
	}
	defer klient.Zamknij()

	// Szkic wskazany wyklucza się z polami treści; przyjęcie obu naraz
	// zmuszałoby do wyboru.
	if z.DraftId != nil && strings.TrimSpace(*z.DraftId) != "" {
		if len(z.To) > 0 || z.Subject != nil || z.Body != nil {
			return shared.MailSendResponse{}, shared.MailMessage{}, bladWskazaniaPoczty(
				"komenda mail.send niesie naraz szkic i treść — wskaż jedno albo drugie")
		}
		return a.wyslijSzkic(ctx, klient, skrzynka, strings.TrimSpace(*z.DraftId))
	}

	wychodzacy, err := a.zlozWychodzacy(ctx, skrzynka, wysylanyList{
		Do:            z.To,
		Kopia:         z.Cc,
		Temat:         z.Subject,
		Tresc:         z.Body,
		WOdpowiedziNa: z.InReplyTo,
		Zasoby:        z.AttachmentAssetIds,
	})
	if err != nil {
		return shared.MailSendResponse{}, shared.MailMessage{}, err
	}
	return a.nadaj(ctx, klient, skrzynka, wychodzacy)
}

// wyslijSzkic odczytuje zapisany szkic ze skrzynki i nadaje jego treść. Szkic
// po wysłaniu nie jest kasowany: kontrakt tego nie obiecuje, a skasowanie
// byłoby czynnością uboczną, której nikt nie zlecił.
func (a *adapterPoczty) wyslijSzkic(ctx context.Context, klient *poczta.Klient,
	skrzynka dane.SkrzynkaOperatora, szkic string) (shared.MailSendResponse, shared.MailMessage, error) {

	list, err := klient.Pobierz(szkic, true)
	if err != nil {
		return shared.MailSendResponse{}, shared.MailMessage{}, bladOdczytuListu(err)
	}
	return a.nadaj(ctx, klient, skrzynka, poczta.Wychodzacy{
		Od:               skrzynka.Adres,
		NazwaWyswietlana: wartoscLubPustka(skrzynka.NazwaWyswietlana),
		Do:               list.Do,
		Kopia:            list.Kopia,
		Temat:            list.Temat,
		Tresc:            list.Tresc,
		Zalaczniki:       list.Zalaczniki,
	})
}

// nadaj wykonuje wysyłkę i zawsze zostawia ślad, wiersz w bazie zapisywany
// zarówno po udanym, jak i po nieudanym nadaniu.
func (a *adapterPoczty) nadaj(ctx context.Context, klient *poczta.Klient,
	skrzynka dane.SkrzynkaOperatora, w poczta.Wychodzacy) (shared.MailSendResponse, shared.MailMessage, error) {

	identyfikator, nadano, blad := klient.Wyslij(w)
	a.zapiszSlad(ctx, skrzynka, w, blad)
	if blad != nil {
		// Kod `channel_unavailable`: niedostępność serwera bywa chwilowa,
		// ponowienie ma sens.
		return shared.MailSendResponse{}, shared.MailMessage{},
			protocolBladPoczty(shared.ErrorCodeChannelUnavailable, blad.Error())
	}

	// Wiadomość zdarzenia niesie folder wysłanych, nadawcę skrzynki i chwilę
	// nadania z serwera.
	wyslana := shared.MailMessage{
		Id:     identyfikator,
		Folder: poczta.FolderWyslanych,
		From:   skrzynka.Adres,
		To:     w.Do,
		Cc:     w.Kopia,
		Date:   nadano.UnixMilli(),
	}
	if temat := strings.TrimSpace(w.Temat); temat != "" {
		wyslana.Subject = &temat
	}
	return shared.MailSendResponse{MessageId: identyfikator, SentAt: nadano.UnixMilli()}, wyslana, nil
}

// zapiszSlad utrwala fakt nadania — udanego i nieudanego. Niepowodzenie zapisu
// śladu nie unieważnia wysyłki: list już wyszedł i nie da się go cofnąć, więc
// odmowa komendy z powodu bazy mówiłaby nieprawdę.
func (a *adapterPoczty) zapiszSlad(ctx context.Context, skrzynka dane.SkrzynkaOperatora,
	w poczta.Wychodzacy, blad error) {

	if a.skrzynki == nil {
		return
	}
	slad := dane.SladWysylki{
		SkrzynkaID: skrzynka.ID,
		Adresat:    strings.Join(append(append([]string{}, w.Do...), w.Kopia...), ", "),
		Temat:      w.Temat,
		Tresc:      w.Tresc,
		Powodzenie: blad == nil,
	}
	if odpowiedz := strings.TrimSpace(w.WOdpowiedziNa); odpowiedz != "" {
		slad.WOdpowiedziNa = &odpowiedz
	}
	if blad != nil {
		powod := blad.Error()
		slad.Blad = &powod
	}
	_ = a.skrzynki.ZapiszSladWysylki(ctx, slad)
}

// wysylanyList zbiera pola wspólne `mail.draft.save` i `mail.send`. Struktura
// pośrednia, bo obie komendy niosą te same pola treści, a składanie listu ma
// być jedno — szkic oglądany przed wysłaniem musi być tym samym
// dokumentem, który wyjdzie.
type wysylanyList struct {
	Do            []string
	Kopia         []string
	Temat         *string
	Tresc         *string
	WOdpowiedziNa *string
	Zasoby        []string
}

// zlozWychodzacy przekłada żądanie na list do nadania, dobierając bajty
// wskazanych zasobów z magazynu.
func (a *adapterPoczty) zlozWychodzacy(ctx context.Context, skrzynka dane.SkrzynkaOperatora,
	z wysylanyList) (poczta.Wychodzacy, error) {

	zalaczniki, err := a.zalacznikiZZasobow(ctx, z.Zasoby)
	if err != nil {
		return poczta.Wychodzacy{}, err
	}
	return poczta.Wychodzacy{
		Od:               skrzynka.Adres,
		NazwaWyswietlana: wartoscLubPustka(skrzynka.NazwaWyswietlana),
		Do:               z.Do,
		Kopia:            z.Kopia,
		Temat:            wartoscLubPustka(z.Temat),
		Tresc:            wartoscLubPustka(z.Tresc),
		WOdpowiedziNa:    strings.TrimSpace(wartoscLubPustka(z.WOdpowiedziNa)),
		Zalaczniki:       zalaczniki,
	}, nil
}

// zalacznikiZZasobow zamienia identyfikatory zasobów magazynu na bajty do
// dołączenia. Zasób wskazany, a nieznany jest odmową całej komendy: wysłanie
// listu bez załącznika, o który proszono, byłoby wysłaniem innego listu niż
// zamówiony.
func (a *adapterPoczty) zalacznikiZZasobow(ctx context.Context, zasoby []string) ([]poczta.Zalacznik, error) {
	if len(zasoby) == 0 {
		return nil, nil
	}
	if a.zasoby == nil || a.magazyn == nil {
		return nil, bladPoczty(errors.New(
			"magazyn zasobów nie jest wpięty — nie ma skąd wziąć bajtów załączników"))
	}
	zalaczniki := make([]poczta.Zalacznik, 0, len(zasoby))
	for _, kod := range zasoby {
		zasob, err := a.zasoby.Zasob(ctx, kod)
		if errors.Is(err, dane.ErrBrakWiersza) {
			return nil, protocolBladPoczty(shared.ErrorCodeNotFound,
				"platforma nie zna zasobu "+kod+" wskazanego jako załącznik")
		}
		if err != nil {
			return nil, bladPoczty(err)
		}
		if zasob.URI == nil || strings.TrimSpace(*zasob.URI) == "" {
			return nil, protocolBladPoczty(shared.ErrorCodeConflict,
				"zasób "+kod+" nie ma treści w magazynie rdzenia — nie ma czego dołączyć")
		}
		// Bajty pochodzą wprost ze ścieżki bloba `uri` zasobu, tą samą drogą,
		// co moduł dokumentów.
		bajty, err := os.ReadFile(*zasob.URI)
		if err != nil {
			return nil, bladPoczty(err)
		}
		zalaczniki = append(zalaczniki, poczta.Zalacznik{
			Nazwa:     nazwaZalacznika(zasob, kod),
			TypTresci: wartoscLubPustka(zasob.Format),
			Bajty:     bajty,
		})
	}
	return zalaczniki, nil
}

// nazwaZalacznika bierze nazwę zasobu, a przy jej braku — jego kod. Kod jest
// tu etykietą, nie zmyśloną nazwą pliku: odbiorca zobaczy coś, po czym da się
// ten załącznik odnaleźć w magazynie rdzenia.
func nazwaZalacznika(zasob dane.ZasobDesignu, kod string) string {
	if zasob.Nazwa != nil && strings.TrimSpace(*zasob.Nazwa) != "" {
		return *zasob.Nazwa
	}
	return kod
}

// protocolBladPoczty składa odmowę modułu z wskazanym kodem kontraktu, znacząc
// treść przedrostkiem nazwy modułu dla odróżnienia od odmów innych modułów.
func protocolBladPoczty(kod shared.ErrorCode, powod string) error {
	return protocol.JakoError(protocol.NowyBlad(kod, "moduł poczty: "+powod))
}
