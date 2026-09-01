// Plik wpina dziesięć komend obszaru mail.* i rozgłasza mail.changed po tych, które skrzynkę zmieniają: zapisany szkic, wysłany list, zmianę oznaczenia.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Poczta jest portem rodziny mail.*, stojącym osobno od portu Asystenta: skrzynka jest zasobem urządzenia, po który model sięga w dowolnym oknie, tak samo jak po narzędzia obrazu czy dokumentów.
type Poczta interface {
	Skrzynki(ctx context.Context, z shared.MailAccountListRequest) (shared.MailAccountListResponse, error)
	Foldery(ctx context.Context, z shared.MailFolderListRequest) (shared.MailFolderListResponse, error)
	Wiadomosci(ctx context.Context, z shared.MailMessageListRequest) (shared.MailMessageListResponse, error)
	Wiadomosc(ctx context.Context, z shared.MailMessageGetRequest) (shared.MailMessageGetResponse, error)
	ZapiszSzkic(ctx context.Context, z shared.MailDraftSaveRequest) (shared.MailDraftSaveResponse, error)
	// Wyslij oddaje trzy wartości, bo zdarzenie mail.changed niesie wiadomość, której odpowiedź nie ma.
	Wyslij(ctx context.Context, z shared.MailSendRequest) (shared.MailSendResponse, shared.MailMessage, error)
	Oznacz(ctx context.Context, z shared.MailMessageFlagRequest) (shared.MailMessageFlagResponse, error)
	Podepnij(ctx context.Context, z shared.MailAccountAddRequest) (shared.MailAccountAddResponse, error)
	Rozpoznaj(ctx context.Context, z shared.MailAccountDiscoverRequest) (shared.MailAccountDiscoverResponse, error)
	Odepnij(ctx context.Context, z shared.MailAccountRemoveRequest) (shared.MailAccountRemoveResponse, error)
}

// zarejestrujPoczte wpina dziesięć komend rodziny mail.* obsługujących skrzynki i wiadomości poczty użytkownika.
func zarejestrujPoczte(r *Rejestr, m Poczta, e *emiter) {
	if r == nil || m == nil {
		return
	}

	// Cztery odczyty — bez opakowania rozgłaszającego, bo niczego nie zmieniają.
	r.Zarejestruj(shared.CommandMailAccountList, obsluz(m.Skrzynki))
	r.Zarejestruj(shared.CommandMailFolderList, obsluz(m.Foldery))
	r.Zarejestruj(shared.CommandMailMessageList, obsluz(m.Wiadomosci))
	r.Zarejestruj(shared.CommandMailMessageGet, obsluz(m.Wiadomosc))

	// Trzy czynności nad katalogiem skrzynek: podpięcie, rozpoznanie i odpięcie.
	r.Zarejestruj(shared.CommandMailAccountAdd, obsluz(m.Podepnij))
	r.Zarejestruj(shared.CommandMailAccountDiscover, obsluz(m.Rozpoznaj))
	r.Zarejestruj(shared.CommandMailAccountRemove, obsluz(m.Odepnij))

	// Szkic rozgłasza się zawsze jako created, bo IMAP nie zna poprawiania w miejscu.
	r.Zarejestruj(shared.CommandMailDraftSave,
		obsluz(func(ctx context.Context, z shared.MailDraftSaveRequest) (shared.MailDraftSaveResponse, error) {
			odpowiedz, err := m.ZapiszSzkic(ctx, z)
			if err == nil {
				e.poczta(ctx, shared.ChangeKindCreated, wskazanaSkrzynka(z.AccountId), &odpowiedz.Draft)
			}
			return odpowiedz, err
		}))

	// Wysłanie rozgłasza created wyłącznie po nadaniu udanym, nie po liście, który nie wyszedł.
	r.Zarejestruj(shared.CommandMailSend,
		obsluz(func(ctx context.Context, z shared.MailSendRequest) (shared.MailSendResponse, error) {
			odpowiedz, wyslany, err := m.Wyslij(ctx, z)
			if err == nil {
				e.poczta(ctx, shared.ChangeKindCreated, wskazanaSkrzynka(z.AccountId), &wyslany)
			}
			return odpowiedz, err
		}))

	// Oznaczenie zmienia list zastany — `updated`.
	r.Zarejestruj(shared.CommandMailMessageFlag,
		obsluz(func(ctx context.Context, z shared.MailMessageFlagRequest) (shared.MailMessageFlagResponse, error) {
			odpowiedz, err := m.Oznacz(ctx, z)
			if err == nil {
				e.poczta(ctx, shared.ChangeKindUpdated, wskazanaSkrzynka(z.AccountId), &odpowiedz.Message)
			}
			return odpowiedz, err
		}))
}

// wskazanaSkrzynka oddaje skrzynkę z żądania albo pustkę. Pustka znaczy „domyślna" — tak samo, jak rozumie ją adapter. Podstawianie tu nazwy skrzynki domyślnej wymagałoby zapytania bazy z wnętrza rejestru, który bazy nie zna i znać nie ma.
func wskazanaSkrzynka(id *string) string {
	if id == nil {
		return ""
	}
	return *id
}

// poczta rozgłasza `mail.changed`. Skrzynka jest bytem platformy, nie karty sesji — jedna poczta obsługuje wszystkie sesje konta — więc zdarzenie idzie bez wskazania sesji.
func (e *emiter) poczta(ctx context.Context, zmiana shared.ChangeKind,
	skrzynka string, wiadomosc *shared.MailMessage) {

	zdarzenie := shared.MailChangedEvent{Change: zmiana, AccountId: skrzynka, Message: wiadomosc}
	zdarzenie.Actor, zdarzenie.ActorClientId = sprawca(ctx)
	e.wyslijDoKonta(ctx, shared.EventMailChanged, "", zdarzenie)
}
