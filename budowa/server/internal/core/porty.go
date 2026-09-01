// Pakiet core jest rdzeniem dyspozycji: rozstrzyga komendę po nazwie z kontraktu, kieruje ją do obsługiwacza właściwej domeny przez interfejsy tego pliku i zwraca kopertę odpowiedzi.
package core

import (
	"context"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Sesje obsługuje sesję — byt wspólny dla plików, pamięci, projektu i agentów.
// Okna komunikacji sesji są bytem osobnym, obsługiwanym przez Okna.
type Sesje interface {
	Utworz(ctx context.Context, z shared.SessionCreateRequest) (shared.SessionCreateResponse, error)
	Wykaz(ctx context.Context, z shared.SessionListRequest) (shared.SessionListResponse, error)
	Otworz(ctx context.Context, z shared.SessionOpenRequest) (shared.SessionOpenResponse, error)
	Wydaj(ctx context.Context, z shared.SessionExportRequest) (shared.SessionExportResponse, error)
	Zamknij(ctx context.Context, z shared.SessionCloseRequest) (shared.SessionCloseResponse, error)
	Usun(ctx context.Context, z shared.SessionDeleteRequest) (shared.SessionDeleteResponse, error)
	// Czynności Operatora na wykazie sesji; wszystkie są odwracalne poza usunięciem.
	ZmienNazwe(ctx context.Context, z shared.SessionRenameRequest) (shared.SessionRenameResponse, error)
	Archiwizuj(ctx context.Context, z shared.SessionArchiveRequest) (shared.SessionArchiveResponse, error)
	Przywroc(ctx context.Context, z shared.SessionRestoreRequest) (shared.SessionRestoreResponse, error)
	WykazArchiwum(ctx context.Context, z shared.SessionArchiveListRequest) (shared.SessionArchiveListResponse, error)
	PrzypiszProjekt(ctx context.Context, z shared.SessionProjectSetRequest) (shared.SessionProjectSetResponse, error)
	OdepnijProjekt(ctx context.Context, z shared.SessionProjectClearRequest) (shared.SessionProjectClearResponse, error)
	Kopiuj(ctx context.Context, z shared.SessionCopyRequest) (shared.SessionCopyResponse, error)
	Wznow(ctx context.Context, z shared.SessionResumeRequest) (shared.SessionResumeResponse, error)
	Zatrzymaj(ctx context.Context, z shared.SessionStopRequest) (shared.SessionStopResponse, error)
}

// Okna obsługuje okno komunikacji — byt pośredni między sesją a wiadomością.
// Jedna sesja ma wiele okien, każde z własnym modułem, kanałem modelu, listą
// katalogów, zasięgiem wykonania, trybem uprawnień i rolą.
type Okna interface {
	Utworz(ctx context.Context, z shared.WindowCreateRequest) (shared.WindowCreateResponse, error)
	Wykaz(ctx context.Context, z shared.WindowListRequest) (shared.WindowListResponse, error)
	Zmien(ctx context.Context, z shared.WindowUpdateRequest) (shared.WindowUpdateResponse, error)
	Zamknij(ctx context.Context, z shared.WindowCloseRequest) (shared.WindowCloseResponse, error)
}

// Rozmowa obsługuje wiadomości okna komunikacji. Fragmenty odpowiedzi modelu
// nie wracają tą drogą — idą strumieniem stream.chunk przez Nadajnik,
// bo strumień trwa dłużej niż wykonanie komendy.
type Rozmowa interface {
	Wyslij(ctx context.Context, z shared.MessageSendRequest) (shared.MessageSendResponse, error)
	Zatrzymaj(ctx context.Context, z shared.MessageStopRequest) (shared.MessageStopResponse, error)
	Wykaz(ctx context.Context, z shared.MessageListRequest) (shared.MessageListResponse, error)
}

// Ustawienia obsługuje konfigurację warstwową ośmiu poziomów zasięgu, w tym samym rejestrze co konfiguracja obszarów sesji; brak ustawienia jest wartością domyślną, nigdy blokadą.
type Ustawienia interface {
	Odczytaj(ctx context.Context, z shared.ConfigGetRequest) (shared.ConfigGetResponse, error)
	Zapisz(ctx context.Context, z shared.ConfigSetRequest) (shared.ConfigSetResponse, error)
	Przywroc(ctx context.Context, z shared.ConfigResetRequest) (shared.ConfigResetResponse, error)
	KonfiguracjaSesji
}

// Kanaly obsługuje rejestr kanałów modelu sterowany danymi:
// nowy kanał to nowy wiersz rejestru, nie nowy typ w kodzie.
type Kanaly interface {
	Dodaj(ctx context.Context, z shared.ChannelAddRequest) (shared.ChannelAddResponse, error)
	Zmien(ctx context.Context, z shared.ChannelUpdateRequest) (shared.ChannelUpdateResponse, error)
	Usun(ctx context.Context, z shared.ChannelRemoveRequest) (shared.ChannelRemoveResponse, error)
	Wykaz(ctx context.Context, z shared.ChannelListRequest) (shared.ChannelListResponse, error)
	// Sprawdz obsluguje `channel.check` — prawdziwe wywolanie kanalu, nie
	// ogladanie wiersza rejestru.
	Sprawdz(ctx context.Context, z shared.ChannelCheckRequest) (shared.ChannelCheckResponse, error)
	// StanPoswiadczenia obsluguje `channel.credential.status` — STAN, nigdy tresc.
	StanPoswiadczenia(ctx context.Context, z shared.ChannelCredentialStatusRequest) (shared.ChannelCredentialStatusResponse, error)
}

// Kolejki obsługuje jeden silnik pętli koordynator–wykonawca, wspólny dla pętli
// sesyjnej i MultitaskingAI. Bieg naprawczy nie ma limitu obiegów
// — rdzeń nie dokłada własnego progu.
type Kolejki interface {
	Utworz(ctx context.Context, z shared.QueueCreateRequest) (shared.QueueCreateResponse, error)
	Wykonaj(ctx context.Context, z shared.QueueActionRequest) (shared.QueueActionResponse, error)
}

// Przenoszenie obsługuje przekazanie kompletu kontekstu między modułami jedną
// komendą: polecenie, dokumenty, projekt, agenci, historia, źródła
// wiedzy i parametry wykonania.
type Przenoszenie interface {
	Przenies(ctx context.Context, z shared.ContextTransferRequest) (shared.ContextTransferResponse, error)
}

// Nadajnik oddaje warstwie transportu komunikat wychodzący spoza pary żądanie–odpowiedź: zdarzenie zmiany rozgłaszane do połączeń konta oraz fragment strumienia odpowiedzi modelu. Konto puste znaczy wszystkie połączenia rdzenia.
type Nadajnik interface {
	Rozglos(konto string, k protocol.Koperta)
}

// RozlaczanieSesji zrywa gniazda, których sesja bramki przestała nadawać, od razu po unieważnieniu — bez niego gniazdo z unieważnioną sesją żyje do swojej następnej komendy. Konto puste znaczy wszystkie konta; gniazdo wołającego zostaje, bo odpowiedź na komendę ma do niego dojść.
type RozlaczanieSesji interface {
	RozlaczPoUniewaznieniu(ctx context.Context, konto, powod string) int
}

// Nasluch jest warstwą transportu widzianą przez rdzeń. Sluchaj pracuje do
// zamknięcia kontekstu i dopiero wtedy zwraca sterowanie.
type Nasluch interface {
	Sluchaj(ctx context.Context) error
}
