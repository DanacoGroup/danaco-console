// Doradca jest kanałem modelu, którego agent pytający woła po radę, a nie po
// wykonanie pracy; pojęcie trzyma się na jawności rady, śladzie w prowenancji
// i braku wiążącej mocy.
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
	// Sila jest porównywalną miarą wiersza, wyższa znaczy silniejszy; to ona
	// jest sufitem doboru.
	Sila int
	// SilaZnana mówi, czy Sila w ogóle coś znaczy; wiersz bez czytelnego
	// parametru nie ma siły zerowej.
	SilaZnana bool
}

// Pytanie jest zapytaniem agenta do doradcy. Niesie sprawę i kontekst osobno,
// bo doradca ma odpowiedzieć na sprawę, a kontekst tylko przeczytać.
type Pytanie struct {
	// Okno i Wiadomosc wiążą konsultację ze strumieniem, którym idzie
	// odpowiedź agenta pytającego.
	Okno      string
	Wiadomosc string
	// Pytajacy to kod kanału agenta zadającego pytanie; rozstrzyga próg sufitu.
	Pytajacy string
	// ZadanyDoradca jest prośbą wołającego modelu: wolno jej zawężać wybór, nie
	// podnosić sufitu.
	ZadanyDoradca string
	// Sprawa jest pytaniem wprost; Kontekst materiałem do przeczytania (bywa pusty).
	Sprawa   string
	Kontekst string
}

// Sprawdz odrzuca każde pytanie do doradcy, którego zadać się nie da; pusty
// kontekst przechodzi bez odmowy.
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

// Tresc składa pełną treść wysyłaną doradcy: ramę konsultacji, sprawę do
// rozważenia i kontekst dodatkowy.
func (p Pytanie) Tresc() string {
	czesci := []string{ramaKonsultacji, "PYTANIE:\n" + strings.TrimSpace(p.Sprawa)}
	if kontekst := strings.TrimSpace(p.Kontekst); kontekst != "" {
		czesci = append(czesci, "KONTEKST:\n"+kontekst)
	}
	return strings.Join(czesci, "\n\n")
}

// Rada jest odpowiedzią doradcy; niesie zawsze kanał i model doradcy, bo rada
// bez wskazania, kto ją dał, wyglądałaby jak odpowiedź agenta.
type Rada struct {
	Doradca Kandydat
	Tresc   string
	Skrot   string
	Chwila  time.Time
}

// ZlozRade wiąże radę z doradcą, danym pytaniem i chwilą jej otrzymania,
// licząc przy tym skrót konsultacji.
func ZlozRade(doradca Kandydat, pytanie Pytanie, tresc string) Rada {
	chwila := time.Now().UTC()
	return Rada{
		Doradca: doradca,
		Tresc:   strings.TrimSpace(tresc),
		Skrot:   skrotKonsultacji(doradca, pytanie, tresc, chwila),
		Chwila:  chwila,
	}
}

// skrotKonsultacji liczy skrót jednej konsultacji, nie samej pary
// pytanie-rada, żeby dwie konsultacje tego samego pytania u różnych doradców
// miały skrót różny.
func skrotKonsultacji(doradca Kandydat, pytanie Pytanie, tresc string, chwila time.Time) string {
	return models.SkrotTresci(pytanie.Tresc(), tresc, doradca.Kanal, doradca.Model,
		strconv.FormatInt(chwila.UnixNano(), 10))
}

// Wiazaca oddaje zawsze fałsz (własność c) — kod, który chciałby radę wykonać,
// musi najpierw napisać, że robi to wbrew pojęciu.
func (r Rada) Wiazaca() bool { return false }

// Jawnie składa blok pokazywany Operatorowi: kto radził, o co był pytany, co
// odpowiedział i skąd wziął się ten doradca.
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
