// Odpowiedzialność pliku: przekład wierszy obszaru Roundtable na byty kontraktu
// i odmowy modułu opatrzone kodami błędu kontraktu. Kody rozróżniają brak
// danych, byt nieistniejący i usterkę rdzenia — każdy okno pokazuje inaczej.
package core

import (
	"errors"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// uczestnikKontraktu przekłada wiersz składu na uczestnika kontraktu. Rola
// pusta nie wychodzi jako pusty napis, bo RoundtableParticipantRole jest
// wyliczeniem zamkniętym — uczestnik bez wskazanej roli oddaje pole nieobecne.
func uczestnikKontraktu(u dane.UczestnikDebaty) shared.RoundtableParticipant {
	wyciszony, kolejnosc, kluczowy, waga := u.Wyciszony, u.Kolejnosc, u.Kluczowy, u.Waga
	uczestnik := shared.RoundtableParticipant{
		Id:              u.Kod,
		WindowId:        u.Okno,
		ChannelId:       u.KanalModelu,
		PersonaName:     u.NazwaTozsamosci,
		SystemPrompt:    u.PromptSystemowy,
		Muted:           &wyciszony,
		Order:           &kolejnosc,
		Key:             &kluczowy,
		Weight:          &waga,
		AgentId:         u.Agent,
		Avatar:          u.Awatar,
		RoleDescription: u.OpisRoli,
	}
	if rola := strings.TrimSpace(u.Rola); rola != "" {
		wartosc := shared.RoundtableParticipantRole(rola)
		uczestnik.Role = &wartosc
	}
	if u.LiczbaProbek > 0 {
		probki := u.LiczbaProbek
		uczestnik.SampleCount = &probki
	}
	return uczestnik
}

// uczestnicyKontraktu przekłada cały skład debaty na wykaz uczestników kontraktu, wiersz po wierszu składu.
func uczestnicyKontraktu(uczestnicy []dane.UczestnikDebaty) []shared.RoundtableParticipant {
	wykaz := make([]shared.RoundtableParticipant, 0, len(uczestnicy))
	for _, uczestnik := range uczestnicy {
		wykaz = append(wykaz, uczestnikKontraktu(uczestnik))
	}
	return wykaz
}

// kodyUczestnikow wyciąga identyfikatory adresatów tury z wykazu uczestników debaty złożonych w oknie.
func kodyUczestnikow(uczestnicy []dane.UczestnikDebaty) []string {
	kody := make([]string, 0, len(uczestnicy))
	for _, uczestnik := range uczestnicy {
		kody = append(kody, uczestnik.Kod)
	}
	return kody
}

// turaKontraktu przekłada wiersz tury na turę kontraktu. Granice zerowe nie
// wychodzą jako zero, tylko jako pole nieobecne, bo zero w polu kontraktu
// klient odczytałby jako zakaz mówienia.
func turaKontraktu(t dane.TuraDebaty) shared.RoundtableTurn {
	pytanie, anonimowa, granica := t.Pytanie, t.Anonimowa, t.GranicaTur
	tura := shared.RoundtableTurn{
		Id:        t.Kod,
		WindowId:  t.Okno,
		Index:     t.Numer,
		Topic:     t.Zagadnienie,
		Format:    shared.RoundtableFormat(t.Format),
		Status:    shared.RoundtableTurnStatus(t.Stan),
		StartedAt: chwilaBazy(t.Rozpoczeto),
		Question:  &pytanie,
		Anonymous: &anonimowa,
	}
	if t.Zamknieto != nil {
		zamknieto := chwilaBazy(*t.Zamknieto)
		tura.ClosedAt = &zamknieto
	}
	if nadrzedna := strings.TrimSpace(t.TuraNadrzedna); nadrzedna != "" {
		tura.ParentTurnId = &nadrzedna
	}
	if t.GranicaCzasuMs > 0 {
		czas := t.GranicaCzasuMs
		tura.TimeLimitMs = &czas
	}
	if t.GranicaZnakow > 0 {
		znaki := t.GranicaZnakow
		tura.MaxStatementChars = &znaki
	}
	if granica > 0 {
		tura.TurnLimit = &granica
	}
	return tura
}

// wypowiedzKontraktu przekłada wiersz wypowiedzi na wypowiedź kontraktu.
// Pewność ujemna znaczy nie deklarowano i wtedy pole nie wychodzi, bo zero
// jest deklaracją odrębną od braku pytania o pewność.
func wypowiedzKontraktu(w dane.WypowiedzDebaty) shared.RoundtableStatement {
	redakcja := w.Redakcja
	wypowiedz := shared.RoundtableStatement{
		Id: w.Kod, TurnId: w.TuraKod, ParticipantId: w.Uczestnik,
		Content: w.Tresc, CreatedAt: chwilaBazy(w.Utworzono),
		Revision: &redakcja,
	}
	if odpowiedz := strings.TrimSpace(w.OdpowiedzNa); odpowiedz != "" {
		wypowiedz.ReplyToId = &odpowiedz
	}
	if akt := strings.TrimSpace(w.AktMowy); akt != "" {
		wartosc := shared.RoundtableSpeechAct(akt)
		wypowiedz.SpeechAct = &wartosc
	}
	if w.Pewnosc >= 0 {
		pewnosc := w.Pewnosc
		wypowiedz.Confidence = &pewnosc
	}
	return wypowiedz
}

// stanowiskoKontraktu przekłada wiersz stanowiska wraz z turami, które objęło, na byt kontraktu konsensusu.
func stanowiskoKontraktu(s dane.StanowiskoDebaty, kody []string) shared.RoundtableConsensus {
	wersja, zaakceptowane := s.Wersja, s.Zaakceptowane
	if len(kody) == 0 {
		kody = rozdzielWiersze(s.Tury)
	}
	return shared.RoundtableConsensus{
		Id: s.Kod, WindowId: s.Okno, Content: s.Tresc, TurnIds: kody,
		Version: &wersja, UpdatedAt: chwilaBazy(s.Zaktualizowano),
		Accepted: &zaakceptowane, Context: s.Kontekst, Options: s.Warianty,
		Consequences: s.Konsekwencje,
	}
}

// rozdzielWiersze rozbija wykaz zapisany w jednej kolumnie tekstu na osobne pozycje wykazu wynikowego.
func rozdzielWiersze(zapis string) []string {
	if strings.TrimSpace(zapis) == "" {
		return nil
	}
	czesci := strings.Split(zapis, "\n")
	wykaz := make([]string, 0, len(czesci))
	for _, czesc := range czesci {
		if przyciete := strings.TrimSpace(czesc); przyciete != "" {
			wykaz = append(wykaz, przyciete)
		}
	}
	return wykaz
}

// itoaRundy zapisuje numer tury dziesiętnie, wygodnym zapisem tekstowym do treści komunikatu odmowy okna.
func itoaRundy(liczba int) string {
	return strconv.Itoa(liczba)
}

// bladDebaty znakuje usterkę kodem kontraktu. Błąd, któremu kod już nadano,
// przechodzi tędy bez zmiany — dopiero usterka bez kodu staje się usterką
// wewnętrzną rdzenia.
func bladDebaty(err error) error {
	if err == nil {
		return nil
	}
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeInternalError, err))
}

// bladWskazaniaDebaty nazywa brak danych w żądaniu — to błąd Operatora, nie
// rdzenia, więc kod odmowy jest inny niż przy usterce.
func bladWskazaniaDebaty(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Roundtable: "+powod))
}

// odmowaTuryDebatyWBiegu odmawia otwarcia kolejnej tury oknu, w którym
// uczestnicy jeszcze mówią. Przerwanie tury jest osobną decyzją moderatora,
// nie skutkiem ubocznym otwarcia następnej.
func odmowaTuryDebatyWBiegu(okno string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeConflict,
		"moduł Roundtable: okno "+okno+" prowadzi turę i nie otworzy kolejnej. "+
			"Otwarcie nowego zagadnienia nie przerywa wypowiedzi uczestników — "+
			"zamknij turę bieżącą w Moderator Panelu albo ukierunkuj dyskusję "+
			"(roundtable.moderator.direct), co przerywa ją jawnie."))
}

// bladBrakuKanalow odmawia uruchomienia debaty bez rejestru kanałów. Debata bez
// wykonawcy nie ma prawa udawać, że uczestnicy odpowiedzieli.
func bladBrakuKanalow() error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
		"moduł Roundtable: serwer nie ma rejestru kanałów modelu — uczestnicy nie mają czym odpowiedzieć"))
}

// bladNieznanegoKanalu odmawia dodania uczestnika na kanale spoza rejestru kanałów modelu rdzenia platformy.
func bladNieznanegoKanalu(kod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
		"moduł Roundtable: kanał "+kod+" nie istnieje w rejestrze albo jest nieczynny"))
}

// bladNieznanejTury odróżnia stan tury nieistniejącej od stanu odczytu nieudanego z bazy danych rdzenia.
func bladNieznanejTury(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Roundtable: tura debaty nie istnieje: "+kod))
	}
	return bladDebaty(err)
}

// bladNieznanegoUczestnika odróżnia stan uczestnika nieistniejącego od usterki odczytu wiersza z bazy.
func bladNieznanegoUczestnika(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Roundtable: uczestnik debaty nie istnieje: "+kod))
	}
	return bladDebaty(err)
}
