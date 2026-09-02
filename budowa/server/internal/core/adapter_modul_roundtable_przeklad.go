// Przekład wierszy obszaru Roundtable na byty kontraktu i odmowy modułu z kodami
// kontraktu: brak danych, byt nieistniejący i usterkę rdzenia okno pokazuje inaczej.
package core

import (
	"errors"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// RoundtableParticipantRole jest wyliczeniem zamkniętym: rola pusta wychodzi jako pole nieobecne.
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

func uczestnicyKontraktu(uczestnicy []dane.UczestnikDebaty) []shared.RoundtableParticipant {
	wykaz := make([]shared.RoundtableParticipant, 0, len(uczestnicy))
	for _, uczestnik := range uczestnicy {
		wykaz = append(wykaz, uczestnikKontraktu(uczestnik))
	}
	return wykaz
}

func kodyUczestnikow(uczestnicy []dane.UczestnikDebaty) []string {
	kody := make([]string, 0, len(uczestnicy))
	for _, uczestnik := range uczestnicy {
		kody = append(kody, uczestnik.Kod)
	}
	return kody
}

// Granice zerowe wychodzą jako pole nieobecne: zero klient odczytałby jako zakaz mówienia.
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

// Pewność ujemna znaczy „nie deklarowano”; zero jest deklaracją odrębną od braku pytania.
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

func itoaRundy(liczba int) string {
	return strconv.Itoa(liczba)
}

func bladDebaty(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, dane.ErrKolizjaWiersza) {
		return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeConflict, err))
	}
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeNotFound, err))
	}
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeInternalError, err))
}

func bladWskazaniaDebaty(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Roundtable: "+powod))
}

// Przerwanie tury jest osobną decyzją moderatora, nie skutkiem ubocznym otwarcia następnej.
func odmowaTuryDebatyWBiegu(okno string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeConflict,
		"moduł Roundtable: okno "+okno+" prowadzi turę i nie otworzy kolejnej. "+
			"Otwarcie nowego zagadnienia nie przerywa wypowiedzi uczestników — "+
			"zamknij turę bieżącą w Moderator Panelu albo ukierunkuj dyskusję "+
			"(roundtable.moderator.direct), co przerywa ją jawnie."))
}

// Debata bez wykonawcy nie ma prawa udawać, że uczestnicy odpowiedzieli.
func bladBrakuKanalow() error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
		"moduł Roundtable: serwer nie ma rejestru kanałów modelu — uczestnicy nie mają czym odpowiedzieć"))
}

func bladNieznanegoKanalu(kod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
		"moduł Roundtable: kanał "+kod+" nie istnieje w rejestrze albo jest nieczynny"))
}

func bladNieznanejTury(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Roundtable: tura debaty nie istnieje: "+kod))
	}
	return bladDebaty(err)
}

func bladNieznanegoUczestnika(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Roundtable: uczestnik debaty nie istnieje: "+kod))
	}
	return bladDebaty(err)
}
