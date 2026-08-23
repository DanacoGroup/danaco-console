// Odpowiedzialność pliku: wpięcie dziesięciu komend obszaru `mail.*` i
// rozgłoszenie `mail.changed` po tych z nich, które skrzynkę zmieniają.
//
// Siedem komend kontrakt wystawia modelowi jako narzędzia: wykaz skrzynek,
// wykaz folderów, odnalezienie listu, odczytanie go w całości, zapisanie szkicu,
// wysłanie i oznaczenie. Podpięcie, rozpoznanie i odpięcie skrzynki narzędziami
// nie są, bo rozstrzygają, do czego platforma ma dostęp. Wpinają się tak samo
// jak reszta — różnica leży w tym, że nie ma ich w wykazie narzędzi kontraktu.
//
// `mail.changed` idzie po trzech komendach, nie po dziesięciu: opisuje zmianę
// w skrzynce — zapisany szkic, wysłany list, zmianę oznaczenia. Wykazy niczego
// nie zmieniają, a podpięcie i odpięcie skrzynki zmieniają katalog skrzynek,
// dla którego kontrakt osobnego zdarzenia nie ma; rozgłaszanie ich zdarzeniem
// o wiadomości donosiłoby oknu o zmianie bytu, który się nie zmienił.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Poczta jest portem rodziny `mail.*`.
//
// Port stoi osobno, a nie jako rozszerzenie portu `Asystent`: Asystent prowadzi
// zlecenie, a skrzynka jest zasobem urządzenia, po który model sięga w dowolnym
// oknie — tak samo jak po narzędzia obrazu czy dokumentów. Wtopienie poczty
// w Asystenta odebrałoby ją każdemu innemu oknu.
type Poczta interface {
	Skrzynki(ctx context.Context, z shared.MailAccountListRequest) (shared.MailAccountListResponse, error)
	Foldery(ctx context.Context, z shared.MailFolderListRequest) (shared.MailFolderListResponse, error)
	Wiadomosci(ctx context.Context, z shared.MailMessageListRequest) (shared.MailMessageListResponse, error)
	Wiadomosc(ctx context.Context, z shared.MailMessageGetRequest) (shared.MailMessageGetResponse, error)
	ZapiszSzkic(ctx context.Context, z shared.MailDraftSaveRequest) (shared.MailDraftSaveResponse, error)
	// Wyslij oddaje trzy wartości, jako jedyna w tym porcie. Odpowiedź kontraktu
	// niesie sam identyfikator i chwilę nadania, a zdarzenie `mail.changed` musi
	// nieść wiadomość, żeby okno pokazało, co wyszło. Adapter podaje ją obok
	// odpowiedzi, bo rozgłaszanie z wnętrza modułu byłoby drugą drogą do szyny
	// zdarzeń.
	Wyslij(ctx context.Context, z shared.MailSendRequest) (shared.MailSendResponse, shared.MailMessage, error)
	Oznacz(ctx context.Context, z shared.MailMessageFlagRequest) (shared.MailMessageFlagResponse, error)
	Podepnij(ctx context.Context, z shared.MailAccountAddRequest) (shared.MailAccountAddResponse, error)
	Rozpoznaj(ctx context.Context, z shared.MailAccountDiscoverRequest) (shared.MailAccountDiscoverResponse, error)
	Odepnij(ctx context.Context, z shared.MailAccountRemoveRequest) (shared.MailAccountRemoveResponse, error)
}

// zarejestrujPoczte wpina dziesięć komend rodziny `mail.*`.
func zarejestrujPoczte(r *Rejestr, m Poczta, e *emiter) {
	if r == nil || m == nil {
		return
	}

	// Cztery odczyty — bez opakowania rozgłaszającego, bo niczego nie zmieniają.
	r.Zarejestruj(shared.CommandMailAccountList, obsluz(m.Skrzynki))
	r.Zarejestruj(shared.CommandMailFolderList, obsluz(m.Foldery))
	r.Zarejestruj(shared.CommandMailMessageList, obsluz(m.Wiadomosci))
	r.Zarejestruj(shared.CommandMailMessageGet, obsluz(m.Wiadomosc))

	// Trzy czynności nad katalogiem skrzynek — patrz nagłówek pliku.
	r.Zarejestruj(shared.CommandMailAccountAdd, obsluz(m.Podepnij))
	r.Zarejestruj(shared.CommandMailAccountDiscover, obsluz(m.Rozpoznaj))
	r.Zarejestruj(shared.CommandMailAccountRemove, obsluz(m.Odepnij))

	// Szkic rozgłasza się zawsze jako `created`, także przy poprawianiu: IMAP nie
	// zna poprawiania w miejscu, więc zmieniony szkic to nowa wiadomość w folderze
	// szkiców, ze swoim nowym identyfikatorem (`poczta/imap_zapis.go`).
	// Rozgłoszenie `updated` kazałoby oknu odświeżyć pozycję, której już
	// w skrzynce nie ma.
	r.Zarejestruj(shared.CommandMailDraftSave,
		obsluz(func(ctx context.Context, z shared.MailDraftSaveRequest) (shared.MailDraftSaveResponse, error) {
			odpowiedz, err := m.ZapiszSzkic(ctx, z)
			if err == nil {
				e.poczta(ctx, shared.ChangeKindCreated, wskazanaSkrzynka(z.AccountId), &odpowiedz.Draft)
			}
			return odpowiedz, err
		}))

	// Wysłanie rozgłasza `created`, bo powstał nowy byt — kopia w folderze
	// wysłanych. Zdarzenie idzie wyłącznie po nadaniu udanym: rozgłoszenie po
	// odmowie donosiłoby o liście, który nie wyszedł. Ślad w bazie zapisuje się
	// w obu przypadkach i leży w adapterze.
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

// wskazanaSkrzynka oddaje skrzynkę z żądania albo pustkę. Pustka znaczy
// „domyślna" — tak samo, jak rozumie ją adapter. Podstawianie tu nazwy
// skrzynki domyślnej wymagałoby zapytania bazy z wnętrza rejestru, który bazy
// nie zna i znać nie ma.
func wskazanaSkrzynka(id *string) string {
	if id == nil {
		return ""
	}
	return *id
}

// poczta rozgłasza `mail.changed`. Skrzynka jest bytem platformy, nie karty
// sesji — jedna poczta obsługuje wszystkie sesje konta — więc zdarzenie idzie
// bez wskazania sesji.
func (e *emiter) poczta(ctx context.Context, zmiana shared.ChangeKind,
	skrzynka string, wiadomosc *shared.MailMessage) {

	zdarzenie := shared.MailChangedEvent{Change: zmiana, AccountId: skrzynka, Message: wiadomosc}
	zdarzenie.Actor, zdarzenie.ActorClientId = sprawca(ctx)
	e.wyslij(shared.EventMailChanged, "", zdarzenie)
}
