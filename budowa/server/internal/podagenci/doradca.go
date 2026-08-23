// Pojęcie doradcy — konsultacja u innego modelu.
//
// Doradca jest kanałem modelu, którego agent pytający woła po radę, a nie po
// wykonanie pracy. Pojęcie trzyma się na trzech własnościach:
//
//	jawność rady — `Rada` niesie treść wraz z kanałem i modelem doradcy,
//	    a `Jawnie` składa z nich blok do pokazania. Rada nigdy nie wraca samym
//	    napisem treści, bo wpleciona w wynik agenta wyglądałaby na jego własną
//	    odpowiedź.
//	ślad w prowenancji — jedzie istniejącą drogą rdzenia (fragment
//	    `provenance`), a dziennik (`dziennik_doradcy.go`) zapisuje ten sam opis
//	    wywołania obok pytania i rady. Drugiej prowenancji nie ma.
//	konsultacja, nie delegacja — `Wiazaca()` oddaje zawsze fałsz, a rama pytania
//	    mówi to doradcy wprost: odpowiedzialność za wynik zostaje przy agencie
//	    pytającym.
//
// Kto może być doradcą, rozstrzygają dane, nie kod: czynny wiersz rejestru
// kanałów z parametrem `doradca` i liczbową `sila`. Jak daleko wolno sięgnąć,
// rozstrzyga sufit siły opisany w `doradca_wybor.go`; model z własnej
// inicjatywy w górę nie sięga. Tu jest samo pojęcie: pytanie, rada i jawność.
package podagenci

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"danacoconsole/server/internal/models"
)

// Kandydat jest kanałem dopuszczonym do roli doradcy — wycinkiem wiersza
// rejestru, którego potrzebuje konsultacja. Model jedzie do prowenancji
// i dziennika, żeby ślad mówił, kto radził, a nie tylko którym kanałem.
type Kandydat struct {
	Kanal string
	Nazwa string
	Model string
	// Sila jest porównywalną miarą wiersza; wyższa znaczy silniejszy. To ona
	// jest sufitem doboru — patrz `doradca_wybor.go`.
	Sila int
	// SilaZnana mówi, czy `Sila` w ogóle coś znaczy. Wiersz bez czytelnego
	// parametru `sila` nie ma siły zerowej — nie ma jej wcale, a zero byłoby
	// tu orzeczeniem „najsłabszy" wyciągniętym z braku danych.
	SilaZnana bool
}

// Pytanie jest zapytaniem agenta do doradcy. Niesie sprawę i kontekst osobno,
// bo doradca ma odpowiedzieć na sprawę, a kontekst tylko przeczytać.
type Pytanie struct {
	// Okno i Wiadomosc wiążą konsultację ze strumieniem, w którym padła — tym
	// samym, którym idzie odpowiedź agenta pytającego.
	Okno      string
	Wiadomosc string
	// Pytajacy to kod kanału agenta zadającego pytanie; rozstrzyga próg sufitu.
	Pytajacy string
	// Pola „WskazanyDoradca" tu nie ma — patrz nagłówek `doradca_wybor.go`.
	// Wskazania Operatora produkt nie ma dziś czym przyjąć, więc takie pole
	// znosiłoby sufit przy każdej konsultacji i kłamało o powodzie doboru.
	//
	// ZadanyDoradca to prośba samego wołającego modelu. Wolno jej wybór zawęzić,
	// nie wolno jej podnieść sufitu — prośba o silniejszego bez wskazania
	// Operatora kończy się odmową, nie podmianą.
	ZadanyDoradca string
	// Sprawa jest pytaniem wprost; Kontekst materiałem do przeczytania (bywa pusty).
	Sprawa   string
	Kontekst string
}

// Sprawdz odrzuca pytanie, którego nie da się zadać; pusty kontekst przechodzi.
func (p Pytanie) Sprawdz() error {
	if strings.TrimSpace(p.Sprawa) == "" {
		return errors.New("podagenci: pytanie do doradcy bez treści sprawy")
	}
	return nil
}

// Doradca rozstrzyga, kogo zapytać, przykładając sufit siły do tego pytania.
// Pytanie zna oba wskazania (pytający i prośba modelu), więc to ono, a nie
// wołający, zestawia je w jedno wywołanie doboru.
func (p Pytanie) Doradca(wykaz []models.Definicja) (Kandydat, error) {
	return WybierzDoradce(wykaz, p.Pytajacy, p.ZadanyDoradca)
}

// ramaKonsultacji mówi doradcy wprost, czym jest to wywołanie. Bez niej
// model odpowiada wykonaniem zadania i konsultacja staje się delegacją
// (złamanie własności c).
const ramaKonsultacji = `Jesteś doradcą. To jest KONSULTACJA, nie zlecenie.
Odpowiedz radą: co rozważyć, czego uniknąć, co sprawdzić.
Nie wykonuj zadania za pytającego i nie udawaj, że je wykonałeś.
Odpowiedzialność za wynik zostaje przy agencie, który pyta — Twoja odpowiedź
zostanie mu pokazana JAWNIE, wraz z tym pytaniem.`

// Tresc składa treść wysyłaną doradcy: rama, sprawa i kontekst.
func (p Pytanie) Tresc() string {
	czesci := []string{ramaKonsultacji, "PYTANIE:\n" + strings.TrimSpace(p.Sprawa)}
	if kontekst := strings.TrimSpace(p.Kontekst); kontekst != "" {
		czesci = append(czesci, "KONTEKST:\n"+kontekst)
	}
	return strings.Join(czesci, "\n\n")
}

// Rada jest odpowiedzią doradcy. Niesie zawsze kanał i model doradcy, bo rada
// bez wskazania, kto ją dał, dałaby się podać za odpowiedź własną agenta
// (własność a). Skrót z pytania i rady służy diagnostyce i zestawieniu wpisu
// dziennika ze strumieniem — niczego nie dopuszcza.
type Rada struct {
	Doradca Kandydat
	Tresc   string
	Skrot   string
	Chwila  time.Time
}

// ZlozRade wiąże radę z doradcą i chwilą otrzymania.
func ZlozRade(doradca Kandydat, pytanie Pytanie, tresc string) Rada {
	chwila := time.Now().UTC()
	return Rada{
		Doradca: doradca,
		Tresc:   strings.TrimSpace(tresc),
		Skrot:   skrotKonsultacji(doradca, pytanie, tresc, chwila),
		Chwila:  chwila,
	}
}

// skrotKonsultacji liczy skrót jednej konsultacji, a nie samej pary
// pytanie-rada. Doradca i chwila wchodzą do skrótu, bo `adviceDigest`
// zdarzenia `advisor.consulted` ma dać się zestawić z jednym wierszem
// dziennika `konsultacja_doradcy`. Bez nich dwie konsultacje o tym samym
// pytaniu — u dwóch różnych doradców, dwa różne wiersze dziennika — miałyby
// skrót identyczny, więc jedyna wartość niosąca treść rady w zdarzeniu nie
// wskazywałaby niczego jednoznacznie.
func skrotKonsultacji(doradca Kandydat, pytanie Pytanie, tresc string, chwila time.Time) string {
	return models.SkrotTresci(pytanie.Tresc(), tresc, doradca.Kanal, doradca.Model,
		strconv.FormatInt(chwila.UnixNano(), 10))
}

// Wiazaca oddaje zawsze fałsz (własność c) — kod, który chciałby radę wykonać,
// musi najpierw napisać, że robi to wbrew pojęciu.
func (r Rada) Wiazaca() bool { return false }

// Jawnie składa blok pokazywany Operatorowi: kto radził, o co był pytany i co
// odpowiedział. Skrócony do samej rady przestaje być jawny (własność a).
//
// Mówi też, skąd wziął się ten doradca: w górę model z własnej inicjatywy nie
// sięga, a Operator ma widzieć podstawę doboru w oknie, nie dopiero w dzienniku.
func (r Rada) Jawnie(p Pytanie) string {
	return strings.Join([]string{
		fmt.Sprintf("KONSULTACJA U DORADCY — %s (%s, model %s), %s",
			r.Doradca.Nazwa, r.Doradca.Kanal, r.Doradca.Model, zrodloDoboru(p)),
		"ZAPYTANO: " + strings.TrimSpace(p.Sprawa),
		"RADA: " + r.Tresc,
		"Rada nie jest wiążąca — odpowiedzialność za wynik zostaje przy agencie pytającym.",
	}, "\n")
}

// zrodloDoboru nazywa podstawę wyboru doradcy jednym zdaniem do bloku jawności.
// Podstawa jest jedna — sufit siły — i zdanie mówi dokładnie to. Blok jawności
// ma opisywać zdarzenie, a nie nadawać mu powagi, której nie ma.
func zrodloDoboru(p Pytanie) string {
	if strings.TrimSpace(p.ZadanyDoradca) != "" {
		return "doradca poproszony przez model i przepuszczony pod sufitem siły kanału pytającego"
	}
	return "doradca dobrany pod sufitem siły kanału pytającego (model równy albo słabszy)"
}
