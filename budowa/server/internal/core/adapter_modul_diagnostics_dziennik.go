// Odpowiedzialność pliku: dziennik modułu Diagnostics — przyjęcie linii
// z dziennika rdzenia, zapis partią i obsługa komendy `diagnostics.log.query`
// zasilającej okno Logs Viewer.
//
// Adapter jest odbiorcą dziennika rdzenia, a nie drugim dziennikiem: rdzeń pisze
// o sobie jednym `*log.Logger`, a moduł wpina się w jego wyjście. Drugi,
// równoległy dziennik rozjechałby się z pierwszym co do treści i chwili.
//
// Nagłówek linii zostaje obcięty, bo `log.Logger` wkłada na jej początek
// przedrostek i znacznik czasu. Gdyby wchodziły do treści wpisu, każda linia
// byłaby niepowtarzalna i deduplikacja z licznikiem nie zgrupowałaby ani jednej
// pary. Obcięcie liczy się z flag i przedrostka tego samego dziennika, więc jest
// dokładne, a nie zgadywane.
package core

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// zrodloRdzenia znakuje wpisy pochodzące z dziennika technicznego rdzenia.
const zrodloRdzenia = "rdzeń"

// wielkoscPartii ogranicza jedną transakcję zapisu dziennika.
const wielkoscPartii = 64

// odstepZapisu wyznacza, jak długo pisarz czeka na dopełnienie partii, zanim
// zapisze to, co ma. Bez niego pojedynczy wpis czekałby na dopełnienie całej
// partii, a Logs Viewer pokazywałby dziennik z opóźnieniem.
const odstepZapisu = 250 * time.Millisecond

// PodepnijDziennik wpina moduł w wyjście dziennika rdzenia jako drugiego
// odbiorcę linii; przekazany `*log.Logger` pozostaje tym samym dziennikiem.
//
// Dziennik pusty zostawia moduł bez tego źródła; źródło drugie, czyli odmowy
// wykonania komend, działa bez zmian.
func (a *adapterDiagnostyki) PodepnijDziennik(dziennik *log.Logger) {
	if a == nil || dziennik == nil {
		return
	}
	a.przedrostek, a.flagi = dziennik.Prefix(), dziennik.Flags()
	dziennik.SetOutput(rozgalezienie{pierwsze: dziennik.Writer(), drugie: a})
}

// rozgalezienie oddaje tę samą linię dwóm odbiorcom. Niepowodzenie zapisu
// do dziennika modułu nie ma prawa zabrać linii odbiorcy pierwotnemu —
// diagnostyka nie może zepsuć tego, co diagnozuje.
type rozgalezienie struct {
	pierwsze interface{ Write([]byte) (int, error) }
	drugie   interface{ Write([]byte) (int, error) }
}

func (r rozgalezienie) Write(dane []byte) (int, error) {
	liczba, err := r.pierwsze.Write(dane)
	_, _ = r.drugie.Write(dane)
	return liczba, err
}

// Write przyjmuje linię dziennika rdzenia. Nie zapisuje jej do bazy sam:
// `log.Logger` trzyma przy zapisie własną blokadę, więc czekanie na dysk w tym
// miejscu wstrzymywałoby każdy wątek rdzenia, który chce coś odnotować.
func (a *adapterDiagnostyki) Write(linia []byte) (int, error) {
	tresc := obetnijNaglowekDziennika(string(linia), a.przedrostek, a.flagi)
	if tresc != "" {
		a.zakolejkuj(dane.WpisDiagnostyki{
			Kod:    nowyIdentyfikator(przedrostekWpisuDziennika),
			Chwila: time.Now().UnixMilli(),
			Poziom: shared.LogLevelInfo,
			Zrodlo: wskaznikTekstu(zrodloRdzenia),
			Tresc:  tresc,
			Odcisk: odciskWpisu(shared.LogLevelInfo, zrodloRdzenia, tresc),
		})
	}
	return len(linia), nil
}

// zakolejkuj oddaje wpis pisarzowi. Kolejka pełna oznacza stratę wpisu, a nie
// wstrzymanie rdzenia; strata jest liczona i wychodzi do podsumowania analizy.
func (a *adapterDiagnostyki) zakolejkuj(wpis dane.WpisDiagnostyki) {
	if a == nil || a.repozytorium == nil {
		return
	}
	select {
	case a.wpisy <- wpis:
	default:
		a.odrzucone.Add(1)
	}
}

// pisz zapisuje wpisy partiami do zamknięcia kolejki albo końca życia rdzenia.
func (a *adapterDiagnostyki) pisz(ctx context.Context) {
	defer close(a.koniec)
	zegar := time.NewTicker(odstepZapisu)
	defer zegar.Stop()

	partia := make([]dane.WpisDiagnostyki, 0, wielkoscPartii)
	for {
		select {
		case wpis, otwarta := <-a.wpisy:
			if !otwarta {
				a.zapiszPartie(context.WithoutCancel(ctx), partia)
				return
			}
			partia = append(partia, wpis)
			if len(partia) >= wielkoscPartii {
				a.zapiszPartie(ctx, partia)
				partia = partia[:0]
			}
		case <-zegar.C:
			if len(partia) > 0 {
				a.zapiszPartie(ctx, partia)
				partia = partia[:0]
			}
		case <-ctx.Done():
			a.zapiszPartie(context.WithoutCancel(ctx), partia)
			return
		}
	}
}

// zapiszPartie utrwala partię. Niepowodzenie nie idzie do dziennika rdzenia, bo
// zapis dziennika wywołany niepowodzeniem zapisu dziennika kręciłby się w kółko;
// zamiast tego rośnie licznik widoczny w podsumowaniu analizy.
func (a *adapterDiagnostyki) zapiszPartie(ctx context.Context, partia []dane.WpisDiagnostyki) {
	if len(partia) == 0 || a.repozytorium == nil {
		return
	}
	if err := a.repozytorium.DopiszWpisy(ctx, partia); err != nil {
		a.niezapisane.Add(int64(len(partia)))
	}
}

// PrzeszukajDziennik obsługuje `diagnostics.log.query` — okno Logs Viewer.
//
// Wzorzec zwykły zawęża w bazie, regularny w rdzeniu: SQLite bez rozszerzenia
// nie zna operatora REGEXP. Wzorzec regularny przechodzi więc przez okno wpisów
// odczytane pozostałymi zawężeniami, a wynik przycięty tą drogą wraca oznaczony
// polem `truncated`.
func (a *adapterDiagnostyki) PrzeszukajDziennik(ctx context.Context,
	z shared.DiagnosticsLogQueryRequest) (shared.DiagnosticsLogQueryResponse, error) {

	if a.repozytorium == nil {
		return shared.DiagnosticsLogQueryResponse{}, bladBrakuTrwalosci("przeszukanie dziennika")
	}
	granica := wartoscLiczby(z.Limit)
	filtr := dane.FiltrDziennika{
		Poziom:  wartoscPoziomu(z.Level),
		Zrodlo:  strings.TrimSpace(wartoscTekstu(z.Source)),
		Od:      wartoscChwili(z.FromTime),
		Do:      wartoscChwili(z.ToTime),
		Scal:    wartoscPrawdy(z.Deduplicate),
		Granica: granica,
	}

	wzorzec := strings.TrimSpace(wartoscTekstu(z.Pattern))
	wyrazenie, err := wyrazenieWzorca(wzorzec, wartoscPrawdy(z.Regex))
	if err != nil {
		return shared.DiagnosticsLogQueryResponse{}, err
	}
	if wyrazenie == nil {
		filtr.Wzorzec = wzorzec
	} else {
		// Zawężenie treścią wykona rdzeń, więc baza oddaje okno bez niego.
		filtr.Granica = 0
	}

	wiersze, razem, err := a.repozytorium.Wpisy(ctx, filtr)
	if err != nil {
		return shared.DiagnosticsLogQueryResponse{}, bladDiagnostyki(err)
	}
	przyciete := razem > len(wiersze)

	if wyrazenie != nil {
		dopasowane := make([]dane.WpisDiagnostyki, 0, len(wiersze))
		for _, wiersz := range wiersze {
			if wyrazenie.MatchString(wiersz.Tresc) {
				dopasowane = append(dopasowane, wiersz)
			}
		}
		razem = len(dopasowane)
		if granica > 0 && len(dopasowane) > granica {
			dopasowane, przyciete = dopasowane[:granica], true
		}
		wiersze = dopasowane
	}

	return shared.DiagnosticsLogQueryResponse{
		Entries:   wpisyDziennikaKontraktu(wiersze),
		Total:     wskaznikLiczby(razem),
		Truncated: wskaznikPrawdy(przyciete),
	}, nil
}

// odciskWpisu grupuje wpisy identyczne co do poziomu, źródła i treści.
func odciskWpisu(poziom shared.LogLevel, zrodlo, tresc string) string {
	suma := sha256.Sum256([]byte(string(poziom) + "\x00" + zrodlo + "\x00" + tresc))
	return hex.EncodeToString(suma[:16])
}

// obetnijNaglowekDziennika zdejmuje z linii to, co dołożył `log.Logger`:
// przedrostek, datę, godzinę i wskazanie pliku. Kolejność jest kolejnością
// pakietu `log`, a długości są stałe wyznaczone jego formatem — obcięcie jest
// więc dokładne, nie odgadywane.
func obetnijNaglowekDziennika(linia, przedrostek string, flagi int) string {
	linia = strings.TrimRight(linia, "\r\n")
	if flagi&log.Lmsgprefix == 0 {
		linia = strings.TrimPrefix(linia, przedrostek)
	}
	if flagi&log.Ldate != 0 {
		linia = obetnijPole(linia, len("2006/01/02"))
	}
	if flagi&(log.Ltime|log.Lmicroseconds) != 0 {
		dlugosc := len("15:04:05")
		if flagi&log.Lmicroseconds != 0 {
			dlugosc += len(".000000")
		}
		linia = obetnijPole(linia, dlugosc)
	}
	if flagi&(log.Lshortfile|log.Llongfile) != 0 {
		if koniec := strings.Index(linia, ": "); koniec >= 0 {
			linia = linia[koniec+2:]
		}
	}
	if flagi&log.Lmsgprefix != 0 {
		linia = strings.TrimPrefix(linia, przedrostek)
	}
	return strings.TrimSpace(linia)
}

// obetnijPole zdejmuje pole stałej długości wraz z odstępem po nim.
func obetnijPole(linia string, dlugosc int) string {
	if len(linia) < dlugosc {
		return linia
	}
	return strings.TrimPrefix(linia[dlugosc:], " ")
}

// wyrazenieWzorca kompiluje wzorzec regularny. Wzorzec niepoprawny jest błędem
// wywołującego: wraca nazwa usterki wzorca, a nie wynik pusty, nieodróżnialny
// od pustego dziennika.
func wyrazenieWzorca(wzorzec string, regularny bool) (*regexp.Regexp, error) {
	if wzorzec == "" || !regularny {
		return nil, nil
	}
	skompilowane, err := regexp.Compile(wzorzec)
	if err != nil {
		return nil, bladWskazaniaDiagnostyki("wzorzec regularny " +
			strconv.Quote(wzorzec) + " jest niepoprawny: " + err.Error())
	}
	return skompilowane, nil
}
