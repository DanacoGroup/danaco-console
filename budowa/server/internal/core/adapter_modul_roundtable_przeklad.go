// Odpowiedzialność pliku: przekład wierszy obszaru Roundtable na byty kontraktu
// i odmowy modułu opatrzone kodami błędu kontraktu.
//
// Kody błędu rozróżniają trzy przypadki, bo każdy z nich okno pokazuje inaczej:
// braku danych w żądaniu (`validation_failed`), bytu, którego nie ma
// (`not_found`) i usterki rdzenia (`internal_error`). Okno debaty pokazujące na
// wszystko jedno „nie udało się” nie mówiłoby Operatorowi niczego.
package core

import (
	"errors"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// uczestnikKontraktu przekłada wiersz składu na uczestnika kontraktu.
//
// Rola pusta nie wychodzi jako pusty napis: `RoundtableParticipantRole` jest
// wyliczeniem o zamkniętym zbiorze wartości, a napis pusty żadną z nich nie
// jest. Uczestnik bez wskazanej roli oddaje pole nieobecne.
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

// uczestnicyKontraktu przekłada cały skład debaty.
func uczestnicyKontraktu(uczestnicy []dane.UczestnikDebaty) []shared.RoundtableParticipant {
	wykaz := make([]shared.RoundtableParticipant, 0, len(uczestnicy))
	for _, uczestnik := range uczestnicy {
		wykaz = append(wykaz, uczestnikKontraktu(uczestnik))
	}
	return wykaz
}

// kodyUczestnikow wyciąga identyfikatory adresatów tury.
func kodyUczestnikow(uczestnicy []dane.UczestnikDebaty) []string {
	kody := make([]string, 0, len(uczestnicy))
	for _, uczestnik := range uczestnicy {
		kody = append(kody, uczestnik.Kod)
	}
	return kody
}

// turaKontraktu przekłada wiersz tury na turę kontraktu.
//
// Granice zerowe nie wychodzą jako zero, tylko jako pole nieobecne: zero
// w kolumnie znaczy „bez granicy”, a zero w polu kontraktu klient odczytałby
// jako granicę zerową, czyli zakaz mówienia.
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
//
// Pewność ujemna znaczy „nie deklarowano” i wtedy pole nie wychodzi. Zero jest
// deklaracją — uczestnik, który powiedział, że nie jest pewien niczego, mówi coś
// innego niż uczestnik, którego o pewność nie pytano.
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

// stanowiskoKontraktu przekłada wiersz stanowiska wraz z turami, które objęło.
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

// rozdzielWiersze rozbija wykaz zapisany w jednej kolumnie na pozycje.
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

// itoaRundy zapisuje numer tury dziesiętnie.
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
// uczestnicy jeszcze mówią.
//
// Odmowa jest trójczęściowa:
//
//	co odmówiło  — otwarcie kolejnej tury debaty w tym oknie,
//	dlaczego     — okno prowadzi turę, a jej przerwanie jest osobną decyzją
//	               moderatora, nie skutkiem ubocznym otwarcia następnej,
//	czym zmienić — zamknięciem tury bieżącej w Moderator Panelu albo
//	               ukierunkowaniem dyskusji (`roundtable.moderator.direct`),
//	               które przerywa ją jawnie.
//
// Kod `conflict`: żądanie jest poprawne, odmawia mu stan okna.
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
		"moduł Roundtable: rdzeń nie ma rejestru kanałów modelu — uczestnicy nie mają czym odpowiedzieć"))
}

// bladNieznanegoKanalu odmawia dodania uczestnika na kanale spoza rejestru.
func bladNieznanegoKanalu(kod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
		"moduł Roundtable: kanał "+kod+" nie istnieje w rejestrze albo jest nieczynny"))
}

// bladNieznanejTury odróżnia „tury nie ma” od „odczyt się nie powiódł”.
func bladNieznanejTury(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Roundtable: tura debaty nie istnieje: "+kod))
	}
	return bladDebaty(err)
}

// bladNieznanegoUczestnika odróżnia „uczestnika nie ma” od usterki odczytu.
func bladNieznanegoUczestnika(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Roundtable: uczestnik debaty nie istnieje: "+kod))
	}
	return bladDebaty(err)
}
