package core

import (
	"context"
	"testing"

	"danacoconsole/shared"
)

// zalozSesjeSprawdzianu i zalozOknoSprawdzianu składają kartę sesji oraz okno
// komunikacji na wskazanym kanale — najkrótszą drogą, jaką idzie klient.
func zalozSesjeSprawdzianu(t *testing.T, zmontowany *Zmontowany, zycie context.Context) string {
	t.Helper()

	var sesja shared.SessionCreateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandSessionCreate,
		shared.SessionCreateRequest{Title: wskaznik("sprawdzian kontekstu")}, &sesja)
	return sesja.Session.Id
}

func zalozOknoSprawdzianu(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	sesja, kanal string) string {
	t.Helper()

	var okno shared.WindowCreateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandWindowCreate, shared.WindowCreateRequest{
		SessionId:      sesja,
		ModuleId:       "communication",
		ModelChannelId: kanal,
		ExecutionEnv:   shared.ExecutionEnvCore,
		PermissionMode: shared.PermissionModeAuto,
		WindowRole:     shared.WindowRoleStandalone,
	}, &okno)
	return okno.Window.Id
}

// Sprawdziany tego pliku biorą kanał echo, adapter bez sieci wkompilowany w rdzeń.

// kanalEchoSprawdzianu zakłada kanał bez sieci i oddaje jego identyfikator.
// Parametry niosą wielkość okna kontekstu, bo bez niej pomiar zajętości nie ma
// się do czego odnieść — i mówi o tym wprost zamiast zgadywać.
func kanalEchoSprawdzianu(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	oknoKontekstu int) string {
	t.Helper()

	var wynik shared.ChannelAddResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandChannelAdd, shared.ChannelAddRequest{
		Name:  "kanał echa sprawdzianu",
		Kind:  "api",
		Model: wskaznik("gpt-4o-sprawdzian"),
		Config: jsonSurowy(t, map[string]any{
			"adapter":            "echo",
			"porcja":             "64",
			parametrGranicyOkna:  oknoKontekstu,
			kluczOdwolaniaKanalu: "DANACO_KLUCZ_SPRAWDZIANU",
		}),
	}, &wynik)
	if wynik.Channel.Id == "" {
		t.Fatal("channel.add nie oddał identyfikatora kanału")
	}
	return wynik.Channel.Id
}

// TestSprawdzenieKanaluWolaKanalNaprawde wykazuje, że `channel.check` jest
// wywołaniem, a nie oglądaniem wiersza rejestru: kanał echa odpowiada treścią,
// więc sprawdzenie melduje osiągalność wraz ze zmierzonym czasem.
func TestSprawdzenieKanaluWolaKanalNaprawde(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	kanal := kanalEchoSprawdzianu(t, zmontowany, zycie, 8192)

	var sprawdzenie shared.ChannelCheckResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandChannelCheck,
		shared.ChannelCheckRequest{ChannelId: kanal}, &sprawdzenie)

	if !sprawdzenie.Reachable {
		szczegol := ""
		if sprawdzenie.Detail != nil {
			szczegol = *sprawdzenie.Detail
		}
		t.Fatalf("kanał echa nie odpowiedział: %s", szczegol)
	}
	if sprawdzenie.CheckedAt <= 0 {
		t.Fatal("sprawdzenie nie niesie chwili wykonania")
	}
	if sprawdzenie.LatencyMs == nil {
		t.Fatal("sprawdzenie nie zmierzyło czasu odpowiedzi — a miało wywołać kanał")
	}
}

// TestSprawdzenieNieznanegoKanaluNieUdajeOsiagalnosci wykazuje sprawdzian
// przeciwny: kanał, którego nie ma, nie jest osiągalny, a odpowiedź mówi
// dlaczego, zamiast odmawiać.
func TestSprawdzenieNieznanegoKanaluNieUdajeOsiagalnosci(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	var sprawdzenie shared.ChannelCheckResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandChannelCheck,
		shared.ChannelCheckRequest{ChannelId: "kanal-ktorego-nie-ma"}, &sprawdzenie)

	if sprawdzenie.Reachable {
		t.Fatal("rdzeń zameldował osiągalność kanału, którego nie ma w rejestrze")
	}
	if sprawdzenie.Detail == nil || *sprawdzenie.Detail == "" {
		t.Fatal("odpowiedź nie mówi, dlaczego kanał nie odpowiedział")
	}
}

// TestStanPoswiadczeniaNieOddajeTresci wykazuje granicę tej komendy: mówi, CZY
// poświadczenie jest, i nie niesie ani jednego znaku sekretu.
func TestStanPoswiadczeniaNieOddajeTresci(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	kanal := kanalEchoSprawdzianu(t, zmontowany, zycie, 8192)

	var stan shared.ChannelCredentialStatusResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandChannelCredentialStatus,
		shared.ChannelCredentialStatusRequest{ChannelId: kanal}, &stan)

	if stan.Status.ChannelId != kanal {
		t.Fatalf("odpowiedź dotyczy kanału %q, a pytano o %q", stan.Status.ChannelId, kanal)
	}
	if stan.Status.Kind == nil || *stan.Status.Kind == "" {
		t.Fatal("stan poświadczenia nie mówi, jakiego jest rodzaju")
	}
	// Poświadczenia nikt nie zapisał, więc present ma być fałszem.
	if stan.Status.Present {
		t.Fatal("rdzeń melduje ustawione poświadczenie, którego nikt nie zapisał")
	}
	if stan.Status.ManagedBy == nil || *stan.Status.ManagedBy == "" {
		t.Fatal("stan nie mówi, gdzie poświadczenie miałoby mieszkać")
	}
}

// TestZajetoscKontekstuLiczySieTokenizatorem wykazuje, że pomiar jest
// pomiarem: liczba żetonów rośnie wraz z treścią i nie jest zerem, a odpowiedź
// mówi, którym słownikiem policzono.
func TestZajetoscKontekstuLiczySieTokenizatorem(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianu(t, katalog)
	kanal := kanalEchoSprawdzianu(t, zmontowany, zycie, 8192)

	sesja := zalozSesjeSprawdzianu(t, zmontowany, zycie)
	okno := zalozOknoSprawdzianu(t, zmontowany, zycie, sesja, kanal)

	// Historia rozmowy wnoszona wprost do tabeli: mierzymy tokenizator, a nie
	// drogę tury przez kanał.
	if _, err := baza.Exec(
		`INSERT INTO wiadomosc (okno_komunikacji_id, rola, rodzaj_tresci, stan, tresc, kolejnosc)
		 SELECT o.id, 'uzytkownik', 'tekst', 'zakonczona',
		        'Przygotuj notatkę ze spotkania zespołu i zapisz ją w bibliotece projektu.', 1
		   FROM okno_komunikacji o WHERE o.identyfikator_zewnetrzny = ?`, okno); err != nil {
		t.Fatalf("nie można dołożyć wiadomości do historii: %v", err)
	}

	var pomiar shared.ContextUsageGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandContextUsageGet,
		shared.ContextUsageGetRequest{WindowId: okno}, &pomiar)

	if !pomiar.Available {
		powod := ""
		if pomiar.Reason != nil {
			powod = *pomiar.Reason
		}
		t.Fatalf("pomiar zajętości nie doszedł do skutku: %s", powod)
	}
	if pomiar.Usage.LimitTokens != 8192 {
		t.Fatalf("granica okna wyszła %d, a kanał deklaruje 8192", pomiar.Usage.LimitTokens)
	}
	if pomiar.Usage.Tokenizer == "" {
		t.Fatal("odpowiedź nie mówi, czym policzono — liczba bez świadka")
	}
	if pomiar.Usage.HistoryTokens == nil || *pomiar.Usage.HistoryTokens <= 0 {
		t.Fatal("historia rozmowy została policzona na zero żetonów, choć w tabeli leży zdanie")
	}
	if pomiar.Usage.UsedTokens < *pomiar.Usage.HistoryTokens {
		t.Fatal("suma żetonów jest mniejsza od samej historii — części nie sumują się do całości")
	}
	if pomiar.Usage.MeasuredAt <= 0 {
		t.Fatal("pomiar nie niesie chwili wykonania")
	}

	// Sprawdzian rozstrzygający: dłuższa historia daje więcej żetonów, nie stałą.
	poprzednie := *pomiar.Usage.HistoryTokens
	if _, err := baza.Exec(
		`INSERT INTO wiadomosc (okno_komunikacji_id, rola, rodzaj_tresci, stan, tresc, kolejnosc)
		 SELECT o.id, 'model', 'tekst', 'zakonczona',
		        'Notatka gotowa. Zapisałem ją w bibliotece projektu pod nazwą „Spotkanie zespołu”.', 2
		   FROM okno_komunikacji o WHERE o.identyfikator_zewnetrzny = ?`, okno); err != nil {
		t.Fatalf("nie można dołożyć drugiej wiadomości: %v", err)
	}

	var drugi shared.ContextUsageGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandContextUsageGet,
		shared.ContextUsageGetRequest{WindowId: okno}, &drugi)
	if drugi.Usage.HistoryTokens == nil || *drugi.Usage.HistoryTokens <= poprzednie {
		t.Fatalf("po dołożeniu wypowiedzi historia liczy %v żetonów, a wcześniej %d — "+
			"licznik nie reaguje na treść", drugi.Usage.HistoryTokens, poprzednie)
	}
}

// TestZajetoscBezGranicyOknaNieZmyslaLiczby wykazuje uczciwe „nie wiem": kanał
// bez zadeklarowanej wielkości okna daje odpowiedź niedostępną wraz z powodem,
// a nie pasek wobec granicy wziętej z głowy.
func TestZajetoscBezGranicyOknaNieZmyslaLiczby(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	var kanal shared.ChannelAddResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandChannelAdd, shared.ChannelAddRequest{
		Name: "kanał bez zadeklarowanego okna", Kind: "api",
		Model:  wskaznik("model-bez-okna"),
		Config: jsonSurowy(t, map[string]any{"adapter": "echo"}),
	}, &kanal)

	sesja := zalozSesjeSprawdzianu(t, zmontowany, zycie)
	okno := zalozOknoSprawdzianu(t, zmontowany, zycie, sesja, kanal.Channel.Id)

	var pomiar shared.ContextUsageGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandContextUsageGet,
		shared.ContextUsageGetRequest{WindowId: okno}, &pomiar)

	if pomiar.Available {
		t.Fatalf("rdzeń zmierzył zajętość wobec granicy %d, której nikt nie ustalił",
			pomiar.Usage.LimitTokens)
	}
	if pomiar.Reason == nil || *pomiar.Reason == "" {
		t.Fatal("odpowiedź nie mówi, dlaczego pomiaru nie ma")
	}
}

// TestPowtorzenieWywolaniaZostawiaSlad wykazuje, że powtórzenie jest NOWYM
// wywołaniem: w tabeli śladu leży drugi wiersz wskazujący pierwowzór jako
// rodzica.
func TestPowtorzenieWywolaniaZostawiaSlad(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianu(t, katalog)
	kanal := kanalEchoSprawdzianu(t, zmontowany, zycie, 8192)

	// Pierwowzór wnoszony wprost do śladu: mierzymy powtórzenie, a nie drogę
	// tury przez pętlę sesyjną.
	if _, err := baza.Exec(
		`INSERT INTO prowenancja_wywolanie
		     (kod, kanal_kod, model, stan, poczatek,
		      opoznienie_ms, tresc_zapisana, prompt, odpowiedz)
		 VALUES ('wywolanie-pierwowzor', ?, 'gpt-4o-sprawdzian', 'zakonczone', ?, 120, 1,
		         'policz do trzech', 'raz, dwa, trzy')`,
		kanal, terazWMilisekundachAlertow()); err != nil {
		t.Fatalf("nie można założyć pierwowzoru: %v", err)
	}

	var powtorzenie shared.ProvenanceCallReplayResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandProvenanceCallReplay,
		shared.ProvenanceCallReplayRequest{CallId: "wywolanie-pierwowzor"}, &powtorzenie)

	if !powtorzenie.Replay.ContentAvailable {
		t.Fatal("powtórzenie melduje brak treści, choć pierwowzór ją niesie")
	}
	if powtorzenie.Replay.ReplayCallId == "" {
		t.Fatal("powtórzenie nie oddało identyfikatora nowego wywołania")
	}

	var rodzic string
	var stan string
	if err := baza.QueryRow(
		`SELECT COALESCE(rodzic_kod, ''), stan FROM prowenancja_wywolanie
		  WHERE kod = ?`, powtorzenie.Replay.ReplayCallId).
		Scan(&rodzic, &stan); err != nil {
		t.Fatalf("powtórzenie nie zostawiło wiersza śladu — wywołanie, którego "+
			"rozliczenie nie widzi: %v", err)
	}
	if rodzic != "wywolanie-pierwowzor" {
		t.Fatalf("wiersz powtórzenia wskazuje rodzica %q", rodzic)
	}
	if stan != shared.WartosciBazyModelCallStatus[shared.ModelCallStatusOk] {
		t.Fatalf("powtórzenie na kanale echa zakończyło się stanem %q", stan)
	}
}

// TestPowtorzenieBezZapisanejTresciNieUdajeWykonania wykazuje, że ślad
// prowadzony bez treści daje odpowiedź mówiącą o tym wprost, a nie wywołanie
// z pustym promptem.
func TestPowtorzenieBezZapisanejTresciNieUdajeWykonania(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianu(t, katalog)
	kanal := kanalEchoSprawdzianu(t, zmontowany, zycie, 8192)

	if _, err := baza.Exec(
		`INSERT INTO prowenancja_wywolanie
		     (kod, kanal_kod, stan, poczatek, tresc_zapisana)
		 VALUES ('wywolanie-bez-tresci', ?, 'zakonczone', ?, 0)`,
		kanal, terazWMilisekundachAlertow()); err != nil {
		t.Fatalf("nie można założyć pierwowzoru bez treści: %v", err)
	}

	var powtorzenie shared.ProvenanceCallReplayResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandProvenanceCallReplay,
		shared.ProvenanceCallReplayRequest{CallId: "wywolanie-bez-tresci"}, &powtorzenie)

	if powtorzenie.Replay.ContentAvailable {
		t.Fatal("rdzeń zameldował dostępność treści, której ślad nie zapisał")
	}
	if powtorzenie.Replay.ReplayCallId != "" {
		t.Fatal("rdzeń wykonał powtórzenie mimo braku treści do wysłania")
	}
	if wierszy(t, baza, `SELECT COUNT(*) FROM prowenancja_wywolanie`) != 1 {
		t.Fatal("powtórzenie bez treści zostawiło wiersz nowego wywołania")
	}
}
