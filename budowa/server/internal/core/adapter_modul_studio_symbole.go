// Odpowiedzialność pliku: tablica znaków specjalnych — wykaz do wybrania wraz
// z wyszukaniem po nazwie i po kodzie, znaki ostatnio użyte, wstawienie znaku
// w miejsce kursora oraz zasady autozamiany skrótu na znak.
//
// ── Dlaczego tablica znaków stoi w rdzeniu, a nie w bazie ────────────────────
// Nazwy znaków, ich punkty kodowe i grupy są WIEDZĄ rdzenia — tak samo jak
// arkusz stylów fabryczny i wykaz nośników druku. Do bazy schodzi wyłącznie to,
// co Operator zmienił (jego zasady autozamiany) albo czym się posłużył (jego
// znaki ostatnio użyte). Wpisanie tablicy do migracji dałoby dwa wykazy, które
// rozjadą się przy pierwszym uzupełnieniu.
//
// ── Dlaczego wstawienie znaku idzie drogą zmiany treści ─────────────────────
// Znak specjalny jest ZNAKIEM, nie obiektem: § w treści pisma ma się liczyć do
// długości akapitu, znaleźć w wyszukiwaniu i przenieść przy zmianie formatu tak
// samo jak litera. Dlatego wstawienie idzie przez `postacZamienTresc`, tę samą
// drogę, którą idzie pisanie — a nie przez osobny byt, którego reszta modułu by
// nie widziała. Skutkiem jest zmiana śledzona autora `model`, gdy znak wstawił
// model, i wpis dziennika, którym da się to cofnąć pojedynczo.
//
// ── Dlaczego autozamiana nie zamienia niczego w rdzeniu ─────────────────────
// Zasada autozamiany jest NASTAWĄ, a nie czynnością na dokumencie: zamiana
// zachodzi w chwili pisania, czyli w oknie, na naciśnięcie klawisza. Rdzeń
// trzyma wykaz zasad, wystawia go oknu i pozwala go zmienić — i tyle ma robić.
// Gdyby rdzeń przepuszczał treść przez zasady przy zapisie, Operator, który
// napisał „(c)" świadomie, dostałby znak praw autorskich wbrew sobie i nie
// miałby czym tego cofnąć.
//
// Wykaz zasad — także fabrycznych — stoi przy tym W BAZIE, nie w tym pliku.
// Migracja 368 założyła tabelę `autozamiana_znaku_studio` wraz z kolumną
// `fabryczna` i wpisała zasady fabryczne wierszami, żeby Operator mógł je
// WYŁĄCZYĆ; migracja 371 dołożyła do tamtego wykazu znaki prawnicze i ułamki.
// Powtórzenie wykazu fabrycznego tutaj byłoby drugą prawdą o tym, co wchodzi
// w miejsce skrótu, i pierwsza poprawka by je rozjechała. To jest odwrotnie niż
// z tablicą znaków wyżej — i różnica jest z zamysłu: tablicy znaków Operator nie
// zmienia, a zasadę autozamiany zmienia i wyłącza.
package core

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// symbolIleOstatnich to ile znaków ostatnio użytych rdzeń trzyma pod ręką.
// Dwadzieścia cztery mieści się w tablicy okna bez przewijania.
const symbolIleOstatnich = 24

// symbolIleDomyslnie to ile znaków oddaje wykaz bez wskazania granicy.
const symbolIleDomyslnie = 200

// Grupy znaków. Nazwy są pełne i polskie, bo wchodzą wprost do okna wyboru —
// żadnego kodu ani skrótu.
const (
	symbolGrupaMatematyczne = "matematyczne"
	symbolGrupaWaluty       = "waluty"
	symbolGrupaGreckie      = "greckie"
	symbolGrupaStrzalki     = "strzałki"
	symbolGrupaPrawnicze    = "prawnicze"
	symbolGrupaDiakrytyczne = "diakrytyczne"
	symbolGrupaInterpunkcja = "interpunkcyjne"
)

// SymbolSkladnicaStudia jest kontraktem tabel tablicy znaków: zasad autozamiany
// i znaków ostatnio użytych (migracja 371).
//
// Obszar sięga po nie osobnym kontraktem, a nie po całe repozytorium — tak samo
// jak obszar kontroli pracy po swoje tabele. Brak tych tabel jest brakiem
// montażu rdzenia i mówi to wprost, zamiast udawać puste wykazy.
type SymbolSkladnicaStudia interface {
	ZapiszZasadeAutozamiany(ctx context.Context,
		zasada dane.ZasadaAutozamianyStudia) (dane.ZasadaAutozamianyStudia, error)
	ZasadaAutozamiany(ctx context.Context, skrot string) (dane.ZasadaAutozamianyStudia, error)
	ZasadyAutozamiany(ctx context.Context) ([]dane.ZasadaAutozamianyStudia, error)
	UsunZasadeAutozamiany(ctx context.Context, skrot string) (bool, error)
	OdnotujUzycieZnaku(ctx context.Context, kod, znak string) error
	ZnakiOstatnioUzyte(ctx context.Context, ile int) ([]dane.ZnakOstatnioUzytyStudia, error)
}

// symbolSkladnica oddaje tabele tablicy znaków.
func (a *adapterStudia) symbolSkladnica() (SymbolSkladnicaStudia, error) {
	if a == nil || a.repozytorium == nil {
		return nil, postacBladZaplecza(
			"repozytorium Studia nie zostało podane przy montażu rdzenia")
	}
	skladnica, jest := a.repozytorium.(SymbolSkladnicaStudia)
	if !jest {
		return nil, postacBladZaplecza("repozytorium Studia nie niesie tabel tablicy znaków " +
			"z migracji 371 — autozamiany i znaków ostatnio użytych nie ma gdzie zapisać")
	}
	return skladnica, nil
}

// symbolTablica wylicza znaki do wybrania.
//
// Wykaz jest tym, czego Właściciel wymienia wprost: znaki matematyczne, waluty,
// litery greckie, strzałki, znaki prawnicze, znaki diakrytyczne oraz znaki
// interpunkcyjne niedostępne z klawiatury — półpauza, pauza, cudzysłowy
// drukarskie, twarda spacja, twardy dywiz i znak podziału wyrazu.
func symbolTablica() []shared.StudioSymbol {
	wykaz := make([]shared.StudioSymbol, 0, 220)
	dodaj := func(grupa string, znak string, nazwa string) {
		wykaz = append(wykaz, shared.StudioSymbol{
			Code:      symbolKodZnaku(znak),
			Character: znak,
			Name:      nazwa,
			Category:  postacWskaznikTekstu(grupa),
		})
	}

	// Znaki interpunkcyjne niedostępne z klawiatury. Stoją pierwsze, bo są
	// najczęściej potrzebne w pismie polskim: półpauza w zakresach liczb,
	// cudzysłowy drukarskie w cytatach, twarda spacja przed jednostką.
	dodaj(symbolGrupaInterpunkcja, "–", "półpauza")
	dodaj(symbolGrupaInterpunkcja, "—", "pauza")
	dodaj(symbolGrupaInterpunkcja, "‐", "dywiz")
	dodaj(symbolGrupaInterpunkcja, "‑", "twardy dywiz")
	dodaj(symbolGrupaInterpunkcja, "­", "znak podziału wyrazu")
	dodaj(symbolGrupaInterpunkcja, " ", "twarda spacja")
	dodaj(symbolGrupaInterpunkcja, " ", "spacja cienka")
	dodaj(symbolGrupaInterpunkcja, "„", "cudzysłów drukarski otwierający")
	dodaj(symbolGrupaInterpunkcja, "”", "cudzysłów drukarski zamykający")
	dodaj(symbolGrupaInterpunkcja, "«", "cudzysłów ostrokątny otwierający")
	dodaj(symbolGrupaInterpunkcja, "»", "cudzysłów ostrokątny zamykający")
	dodaj(symbolGrupaInterpunkcja, "‘", "apostrof drukarski otwierający")
	dodaj(symbolGrupaInterpunkcja, "’", "apostrof drukarski zamykający")
	dodaj(symbolGrupaInterpunkcja, "…", "wielokropek")
	dodaj(symbolGrupaInterpunkcja, "•", "kropka wypunktowania")
	dodaj(symbolGrupaInterpunkcja, "·", "kropka środkowa")
	dodaj(symbolGrupaInterpunkcja, "‰", "promil")
	dodaj(symbolGrupaInterpunkcja, "†", "krzyżyk")
	dodaj(symbolGrupaInterpunkcja, "‡", "krzyżyk podwójny")

	// Znaki prawnicze — te, które Właściciel wymienia po nazwie.
	dodaj(symbolGrupaPrawnicze, "§", "paragraf")
	dodaj(symbolGrupaPrawnicze, "¶", "znak akapitu")
	dodaj(symbolGrupaPrawnicze, "©", "prawa autorskie")
	dodaj(symbolGrupaPrawnicze, "®", "znak zastrzeżony")
	dodaj(symbolGrupaPrawnicze, "™", "znak towarowy")
	dodaj(symbolGrupaPrawnicze, "℗", "prawa do nagrania")
	dodaj(symbolGrupaPrawnicze, "№", "numer")
	dodaj(symbolGrupaPrawnicze, "℘", "znak wagi")
	dodaj(symbolGrupaPrawnicze, "℮", "znak szacowanej ilości")

	// Waluty.
	dodaj(symbolGrupaWaluty, "€", "euro")
	dodaj(symbolGrupaWaluty, "$", "dolar")
	dodaj(symbolGrupaWaluty, "£", "funt")
	dodaj(symbolGrupaWaluty, "¥", "jen")
	dodaj(symbolGrupaWaluty, "₽", "rubel")
	dodaj(symbolGrupaWaluty, "₴", "hrywna")
	dodaj(symbolGrupaWaluty, "₹", "rupia")
	dodaj(symbolGrupaWaluty, "₺", "lira turecka")
	dodaj(symbolGrupaWaluty, "₣", "frank")
	dodaj(symbolGrupaWaluty, "¢", "cent")
	dodaj(symbolGrupaWaluty, "¤", "znak waluty")

	// Znaki matematyczne.
	dodaj(symbolGrupaMatematyczne, "±", "plus minus")
	dodaj(symbolGrupaMatematyczne, "×", "znak mnożenia")
	dodaj(symbolGrupaMatematyczne, "÷", "znak dzielenia")
	dodaj(symbolGrupaMatematyczne, "≠", "nie równa się")
	dodaj(symbolGrupaMatematyczne, "≈", "w przybliżeniu równa się")
	dodaj(symbolGrupaMatematyczne, "≡", "tożsamościowo równa się")
	dodaj(symbolGrupaMatematyczne, "≤", "mniejsze albo równe")
	dodaj(symbolGrupaMatematyczne, "≥", "większe albo równe")
	dodaj(symbolGrupaMatematyczne, "∞", "nieskończoność")
	dodaj(symbolGrupaMatematyczne, "√", "pierwiastek")
	dodaj(symbolGrupaMatematyczne, "∑", "suma")
	dodaj(symbolGrupaMatematyczne, "∏", "iloczyn")
	dodaj(symbolGrupaMatematyczne, "∫", "całka")
	dodaj(symbolGrupaMatematyczne, "∂", "pochodna cząstkowa")
	dodaj(symbolGrupaMatematyczne, "∆", "przyrost")
	dodaj(symbolGrupaMatematyczne, "∇", "nabla")
	dodaj(symbolGrupaMatematyczne, "∈", "należy do")
	dodaj(symbolGrupaMatematyczne, "∉", "nie należy do")
	dodaj(symbolGrupaMatematyczne, "⊂", "zawiera się w")
	dodaj(symbolGrupaMatematyczne, "⊃", "zawiera")
	dodaj(symbolGrupaMatematyczne, "∪", "suma zbiorów")
	dodaj(symbolGrupaMatematyczne, "∩", "iloczyn zbiorów")
	dodaj(symbolGrupaMatematyczne, "∅", "zbiór pusty")
	dodaj(symbolGrupaMatematyczne, "∀", "dla każdego")
	dodaj(symbolGrupaMatematyczne, "∃", "istnieje")
	dodaj(symbolGrupaMatematyczne, "¬", "negacja")
	dodaj(symbolGrupaMatematyczne, "∧", "koniunkcja")
	dodaj(symbolGrupaMatematyczne, "∨", "alternatywa")
	dodaj(symbolGrupaMatematyczne, "°", "stopień")
	dodaj(symbolGrupaMatematyczne, "′", "minuta kątowa")
	dodaj(symbolGrupaMatematyczne, "″", "sekunda kątowa")
	dodaj(symbolGrupaMatematyczne, "‱", "punkt bazowy")
	dodaj(symbolGrupaMatematyczne, "½", "jedna druga")
	dodaj(symbolGrupaMatematyczne, "⅓", "jedna trzecia")
	dodaj(symbolGrupaMatematyczne, "¼", "jedna czwarta")
	dodaj(symbolGrupaMatematyczne, "¾", "trzy czwarte")
	dodaj(symbolGrupaMatematyczne, "µ", "mikro")
	dodaj(symbolGrupaMatematyczne, "∝", "proporcjonalne do")
	dodaj(symbolGrupaMatematyczne, "∠", "kąt")
	dodaj(symbolGrupaMatematyczne, "⊥", "prostopadłe")
	dodaj(symbolGrupaMatematyczne, "∥", "równoległe")

	// Strzałki.
	dodaj(symbolGrupaStrzalki, "←", "strzałka w lewo")
	dodaj(symbolGrupaStrzalki, "→", "strzałka w prawo")
	dodaj(symbolGrupaStrzalki, "↑", "strzałka w górę")
	dodaj(symbolGrupaStrzalki, "↓", "strzałka w dół")
	dodaj(symbolGrupaStrzalki, "↔", "strzałka w lewo i w prawo")
	dodaj(symbolGrupaStrzalki, "↕", "strzałka w górę i w dół")
	dodaj(symbolGrupaStrzalki, "⇐", "strzałka podwójna w lewo")
	dodaj(symbolGrupaStrzalki, "⇒", "strzałka podwójna w prawo")
	dodaj(symbolGrupaStrzalki, "⇔", "strzałka podwójna w obie strony")
	dodaj(symbolGrupaStrzalki, "↵", "znak powrotu")
	dodaj(symbolGrupaStrzalki, "⇧", "strzałka wypełniona w górę")
	dodaj(symbolGrupaStrzalki, "⇨", "strzałka wypełniona w prawo")
	dodaj(symbolGrupaStrzalki, "↗", "strzałka w prawo w górę")
	dodaj(symbolGrupaStrzalki, "↘", "strzałka w prawo w dół")

	// Litery greckie — małe i wielkie, bo w piśmie technicznym potrzebne są oba.
	greckieMale := []struct{ znak, nazwa string }{
		{"α", "alfa"}, {"β", "beta"}, {"γ", "gamma"}, {"δ", "delta"},
		{"ε", "epsilon"}, {"ζ", "dzeta"}, {"η", "eta"}, {"θ", "teta"},
		{"ι", "jota"}, {"κ", "kappa"}, {"λ", "lambda"}, {"μ", "mi"},
		{"ν", "ni"}, {"ξ", "ksi"}, {"ο", "omikron"}, {"π", "pi"},
		{"ρ", "rho"}, {"σ", "sigma"}, {"τ", "tau"}, {"υ", "ypsilon"},
		{"φ", "fi"}, {"χ", "chi"}, {"ψ", "psi"}, {"ω", "omega"},
	}
	for _, litera := range greckieMale {
		dodaj(symbolGrupaGreckie, litera.znak, "litera grecka mała "+litera.nazwa)
	}
	greckieWielkie := []struct{ znak, nazwa string }{
		{"Α", "alfa"}, {"Β", "beta"}, {"Γ", "gamma"}, {"Δ", "delta"},
		{"Ε", "epsilon"}, {"Ζ", "dzeta"}, {"Η", "eta"}, {"Θ", "teta"},
		{"Ι", "jota"}, {"Κ", "kappa"}, {"Λ", "lambda"}, {"Μ", "mi"},
		{"Ν", "ni"}, {"Ξ", "ksi"}, {"Ο", "omikron"}, {"Π", "pi"},
		{"Ρ", "rho"}, {"Σ", "sigma"}, {"Τ", "tau"}, {"Υ", "ypsilon"},
		{"Φ", "fi"}, {"Χ", "chi"}, {"Ψ", "psi"}, {"Ω", "omega"},
	}
	for _, litera := range greckieWielkie {
		dodaj(symbolGrupaGreckie, litera.znak, "litera grecka wielka "+litera.nazwa)
	}

	// Znaki diakrytyczne — samodzielne i łączące. Samodzielne służą pismu obcemu,
	// łączące składaniu znaku, którego w wykazie nie ma.
	diakrytyczne := []struct{ znak, nazwa string }{
		{"´", "akcent ostry"}, {"`", "akcent słaby"}, {"ˆ", "daszek"},
		{"˜", "tylda"}, {"¨", "diereza"}, {"˚", "kółko"},
		{"˝", "akcent podwójny ostry"}, {"˛", "ogonek"}, {"ˇ", "haczek"},
		{"¯", "makron"}, {"˙", "kropka nad"}, {"¸", "cedylla"},
		{"́", "akcent ostry łączący"}, {"̀", "akcent słaby łączący"},
		{"̂", "daszek łączący"}, {"̃", "tylda łącząca"},
		{"̈", "diereza łącząca"}, {"̊", "kółko łączące"},
		{"̨", "ogonek łączący"}, {"̧", "cedylla łącząca"},
		{"̄", "makron łączący"}, {"̇", "kropka nad łącząca"},
		{"̌", "haczek łączący"}, {"̆", "brewis łączący"},
	}
	for _, znak := range diakrytyczne {
		dodaj(symbolGrupaDiakrytyczne, znak.znak, znak.nazwa)
	}
	return wykaz
}

// symbolGrupy oddaje grupy znaków w kolejności, w jakiej stoją w tablicy.
func symbolGrupy() []string {
	grupy := make([]string, 0, 7)
	znane := map[string]bool{}
	for _, znak := range symbolTablica() {
		if znak.Category == nil || znane[*znak.Category] {
			continue
		}
		znane[*znak.Category] = true
		grupy = append(grupy, *znak.Category)
	}
	return grupy
}

// symbolKodZnaku zapisuje punkt kodowy znaku szesnastkowo, w postaci `U+00A7`.
// Znak złożony z kilku punktów kodowych oddaje je rozdzielone spacją — inaczej
// kod byłby nieprawdą o znaku, który Operator widzi jako jeden.
func symbolKodZnaku(znak string) string {
	czesci := make([]string, 0, 2)
	for _, punkt := range znak {
		czesci = append(czesci, "U+"+strings.ToUpper(
			symbolDopelnijDoCzterech(strconv.FormatInt(int64(punkt), 16))))
	}
	return strings.Join(czesci, " ")
}

// symbolDopelnijDoCzterech dopełnia zapis szesnastkowy zerami do czterech
// znaków, bo tak zapisuje się punkty kodowe i tak Operator ich szuka.
func symbolDopelnijDoCzterech(zapis string) string {
	for len(zapis) < 4 {
		zapis = "0" + zapis
	}
	return zapis
}

// symbolZnakZKodu odczytuje znak z zapisu punktu kodowego. Przyjmuje `U+00A7`,
// `u+00a7`, `0x00A7` i sam zapis szesnastkowy — Operator wkleja kod w takiej
// postaci, w jakiej go znalazł, a nie w takiej, jakiej rdzeń by sobie życzył.
func symbolZnakZKodu(kod string) (string, bool) {
	oczyszczony := strings.TrimSpace(kod)
	if oczyszczony == "" {
		return "", false
	}
	czesci := strings.Fields(oczyszczony)
	var budowa strings.Builder
	for _, czesc := range czesci {
		zapis := strings.ToLower(strings.TrimSpace(czesc))
		zapis = strings.TrimPrefix(zapis, "u+")
		zapis = strings.TrimPrefix(zapis, "0x")
		zapis = strings.TrimPrefix(zapis, "\\u")
		punkt, err := strconv.ParseInt(zapis, 16, 32)
		if err != nil || punkt <= 0 || punkt > 0x10FFFF {
			return "", false
		}
		budowa.WriteRune(rune(punkt))
	}
	if budowa.Len() == 0 {
		return "", false
	}
	return budowa.String(), true
}

// ── Czynności ───────────────────────────────────────────────────────────────

// TabliceZnakow oddaje znaki do wybrania (`studio.symbol.list`).
//
// Wyszukiwanie idzie po NAZWIE i po KODZIE naraz — Właściciel wymienia oba —
// a znaki ostatnio użyte stoją na wierzchu i są oznaczone, żeby okno mogło je
// pokazać osobno. Zawężenie, które nie trafia w ani jeden znak, jest odmową
// nazwaną: pusty wykaz z odpowiedzią „ok" znaczyłby dla okna, że tablica znaków
// jest pusta.
func (a *adapterStudia) TabliceZnakow(ctx context.Context,
	z shared.StudioSymbolListRequest) (shared.StudioSymbolListResponse, error) {

	skladnica, err := a.symbolSkladnica()
	if err != nil {
		return shared.StudioSymbolListResponse{}, err
	}
	ostatnie, err := skladnica.ZnakiOstatnioUzyte(ctx, symbolIleOstatnich)
	if err != nil {
		return shared.StudioSymbolListResponse{}, bladStudio(err)
	}
	kolejnoscUzycia := make(map[string]int, len(ostatnie))
	for i, wiersz := range ostatnie {
		kolejnoscUzycia[strings.ToUpper(wiersz.Kod)] = i + 1
	}

	grupa := ""
	if z.Category != nil {
		grupa = strings.ToLower(strings.TrimSpace(*z.Category))
	}
	if grupa != "" {
		znana := false
		for _, nazwa := range symbolGrupy() {
			if strings.ToLower(nazwa) == grupa {
				znana = true
				break
			}
		}
		if !znana {
			return shared.StudioSymbolListResponse{}, bladWskazaniaStudio(
				"grupy znaków „" + strings.TrimSpace(*z.Category) + "” rdzeń nie zna; grupy: " +
					strings.Join(symbolGrupy(), ", "))
		}
	}
	szukane := ""
	if z.Query != nil {
		szukane = strings.ToLower(strings.TrimSpace(*z.Query))
	}
	tylkoOstatnie := z.RecentOnly != nil && *z.RecentOnly

	wybrane := make([]shared.StudioSymbol, 0, symbolIleDomyslnie)
	for _, znak := range symbolTablica() {
		_, uzyty := kolejnoscUzycia[strings.ToUpper(znak.Code)]
		if tylkoOstatnie && !uzyty {
			continue
		}
		if grupa != "" && (znak.Category == nil ||
			strings.ToLower(*znak.Category) != grupa) {
			continue
		}
		if szukane != "" && !symbolPasuje(znak, szukane) {
			continue
		}
		// Kolejność użycia niesie porządek wykazu, a nie pole znaku: kontrakt ma
		// na to samo `recentlyUsed`, a wymyślanie drugiego pola nie wchodzi
		// w rachubę.
		znak.RecentlyUsed = postacWskaznikPrawdy(uzyty)
		wybrane = append(wybrane, znak)
	}

	// Znak wskazany kodem, którego tablica nie niesie, i tak jest znakiem —
	// Unicode ma ich więcej, niż zmieści się w wykazie. Odmowa w tym miejscu
	// byłaby odmową wobec znaku istniejącego, więc znak wchodzi do wyniku
	// z nazwą mówiącą, że pochodzi z kodu.
	if szukane != "" && len(wybrane) == 0 {
		if znak, jest := symbolZnakZKodu(szukane); jest {
			wybrane = append(wybrane, shared.StudioSymbol{
				Code:      symbolKodZnaku(znak),
				Character: znak,
				Name:      "znak z punktu kodowego " + symbolKodZnaku(znak),
			})
		}
	}
	if len(wybrane) == 0 {
		powod := "tablica znaków nie niesie ani jednego znaku spełniającego zawężenie"
		if szukane != "" {
			powod += ": szukane „" + strings.TrimSpace(*z.Query) + "”"
		}
		if grupa != "" {
			powod += ", grupa „" + strings.TrimSpace(*z.Category) + "”"
		}
		if tylkoOstatnie {
			powod += ", wyłącznie znaki ostatnio użyte — a Operator nie posłużył się " +
				"jeszcze żadnym"
		}
		return shared.StudioSymbolListResponse{}, bladWskazaniaStudio(powod +
			". Grupy znane: " + strings.Join(symbolGrupy(), ", "))
	}

	// Znaki ostatnio użyte idą na wierzch, w kolejności od najbliższego ręce —
	// tak stanowi wymaganie „znaki ostatnio użyte pod ręką".
	sort.SliceStable(wybrane, func(i, j int) bool {
		pierwszy := kolejnoscUzycia[strings.ToUpper(wybrane[i].Code)]
		drugi := kolejnoscUzycia[strings.ToUpper(wybrane[j].Code)]
		if pierwszy == drugi {
			return false
		}
		if pierwszy == 0 {
			return false
		}
		if drugi == 0 {
			return true
		}
		return pierwszy < drugi
	})

	granica := symbolIleDomyslnie
	if z.Limit != nil {
		if *z.Limit <= 0 {
			return shared.StudioSymbolListResponse{}, bladWskazaniaStudio(
				"granica wykazu znaków musi być większa od zera")
		}
		granica = *z.Limit
	}
	if len(wybrane) > granica {
		wybrane = wybrane[:granica]
	}
	return shared.StudioSymbolListResponse{
		Symbols: wybrane, Categories: symbolGrupy(),
	}, nil
}

// symbolPasuje mówi, czy znak spełnia szukane — po nazwie, po kodzie i po samym
// znaku. Operator wkleja czasem znak, żeby dowiedzieć się, co to jest.
func symbolPasuje(znak shared.StudioSymbol, szukane string) bool {
	if strings.Contains(strings.ToLower(znak.Name), szukane) {
		return true
	}
	if strings.Contains(strings.ToLower(znak.Code), szukane) {
		return true
	}
	if znak.Character == szukane {
		return true
	}
	// Kod bez przedrostka: Operator wpisuje „00a7" tak samo często jak „U+00A7".
	if zapis, jest := symbolZnakZKodu(szukane); jest && zapis == znak.Character {
		return true
	}
	return false
}

// WstawZnak wstawia znak specjalny w miejsce kursora (`studio.symbol.insert`).
//
// Znak da się wskazać na trzy sposoby: samym znakiem, punktem kodowym albo nazwą
// z tablicy — Właściciel wymaga wszystkich trzech, bo „wstaw tu paragraf" jest
// poleceniem modelu, a nie kliknięciem w tablicę. Znak wchodzi drogą zmiany
// treści, więc zmiana śledzona autora `model` i wpis dziennika odkładają się
// same, a blokada fragmentu zatrzymuje wstawienie odmową nazwaną.
func (a *adapterStudia) WstawZnak(ctx context.Context,
	z shared.StudioSymbolInsertRequest) (shared.StudioSymbolInsertResponse, error) {

	znak, err := symbolRozstrzygnij(z)
	if err != nil {
		return shared.StudioSymbolInsertResponse{}, err
	}
	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioSymbolInsertResponse{}, err
	}
	skladnica, err := a.symbolSkladnica()
	if err != nil {
		return shared.StudioSymbolInsertResponse{}, err
	}
	autor := postacAutor(z.Author)
	miejsce, _ := postacZakres(&z.Offset, &z.Offset, postacDlugosc(&stan.forma))

	for _, blokada := range postacBlokadyZakresu(&stan.forma, miejsce, miejsce) {
		if blokada.Scope == shared.StudioLockScopeEveryone || autor == shared.StudioAuthorModel {
			return shared.StudioSymbolInsertResponse{}, bladWskazaniaStudio(
				"znaku „" + znak.Character + "” nie da się wstawić " +
					postacZapisZakresu(miejsce, miejsce) + ": blokada „" + blokada.Name +
					"” nie przepuszcza tej czynności")
		}
	}

	// Postać znaku przejmuje się z miejsca wstawienia: paragraf wstawiony
	// w wytłuszczony nagłówek ma być wytłuszczony, a nie wrócić do kroju
	// domyślnego.
	postacRozetnij(&stan.forma, miejsce)
	var przejeta *shared.StudioCharacterFormat
	if wskazania := postacFragmentyZakresu(&stan.forma, miejsce, miejsce); len(wskazania) > 0 {
		przejeta = postacKopiaZnaku(stan.forma.Blocks[wskazania[0][0]].Runs[wskazania[0][1]].Format)
	}
	autorWpisu := autor
	postacZamienTresc(&stan.forma, miejsce, miejsce, znak.Character, przejeta, &autorWpisu)

	if err := skladnica.OdnotujUzycieZnaku(ctx, znak.Code, znak.Character); err != nil {
		return shared.StudioSymbolInsertResponse{}, bladStudio(err)
	}
	znak.RecentlyUsed = postacWskaznikPrawdy(true)

	stan.opisCzynnosci = "znak „" + znak.Character + "” (" + znak.Name + ") wstawiony " +
		postacZapisZakresu(miejsce, miejsce)
	bilans := shared.StudioActionBalance{
		Applied: 1,
		Skipped: []shared.StudioSkippedItem{},
		Note:    postacWskaznikTekstu(stan.opisCzynnosci),
	}
	forma, bilansGotowy, zmiana, err := a.postacZakoncz(ctx, stan, autor,
		shared.StudioChangeKindWstawienie, shared.StudioActionKindTextEdit,
		miejsce, miejsce+1, bilans)
	if err != nil {
		return shared.StudioSymbolInsertResponse{}, err
	}
	return shared.StudioSymbolInsertResponse{
		Form: forma, Balance: bilansGotowy, Change: zmiana,
		ActionId: stan.czynnosc, Symbol: znak,
	}, nil
}

// symbolRozstrzygnij ustala, który znak Operator albo model wskazał.
//
// Kolejność jest z zamysłu: sam znak, potem punkt kodowy, potem nazwa. Znak
// podany wprost jest najmniej dwuznaczny, a nazwa najbardziej — nazwa niepełna
// dopasowuje się po zawarciu i wtedy odmowa mówi, ile znaków ją spełnia, zamiast
// wstawić pierwszy napotkany.
func symbolRozstrzygnij(z shared.StudioSymbolInsertRequest) (shared.StudioSymbol, error) {
	tablica := symbolTablica()

	if z.Character != nil && *z.Character != "" {
		znak := *z.Character
		if len([]rune(znak)) > 4 {
			return shared.StudioSymbol{}, bladWskazaniaStudio(
				"pole character niesie znak, nie napis — do wstawienia tekstu służy " +
					"komenda studio.text.edit")
		}
		for _, zTablicy := range tablica {
			if zTablicy.Character == znak {
				return zTablicy, nil
			}
		}
		// Znak spoza tablicy jest znakiem prawdziwym: Unicode ma ich więcej, niż
		// zmieści się w wykazie okna. Odmowa byłaby tu odmową wobec znaku,
		// który istnieje.
		return shared.StudioSymbol{
			Code:      symbolKodZnaku(znak),
			Character: znak,
			Name:      "znak z punktu kodowego " + symbolKodZnaku(znak),
		}, nil
	}

	if z.Code != nil && strings.TrimSpace(*z.Code) != "" {
		znak, jest := symbolZnakZKodu(*z.Code)
		if !jest {
			return shared.StudioSymbol{}, bladWskazaniaStudio(
				"punktu kodowego „" + strings.TrimSpace(*z.Code) +
					"” nie da się odczytać. Zapis: U+00A7, 0x00A7 albo samo 00A7; " +
					"znak złożony podaje się kilkoma punktami rozdzielonymi spacją")
		}
		for _, zTablicy := range tablica {
			if zTablicy.Character == znak {
				return zTablicy, nil
			}
		}
		return shared.StudioSymbol{
			Code:      symbolKodZnaku(znak),
			Character: znak,
			Name:      "znak z punktu kodowego " + symbolKodZnaku(znak),
		}, nil
	}

	if z.Name != nil && strings.TrimSpace(*z.Name) != "" {
		szukana := strings.ToLower(strings.TrimSpace(*z.Name))
		dokladne := make([]shared.StudioSymbol, 0, 2)
		czesciowe := make([]shared.StudioSymbol, 0, 8)
		for _, zTablicy := range tablica {
			nazwa := strings.ToLower(zTablicy.Name)
			if nazwa == szukana {
				dokladne = append(dokladne, zTablicy)
				continue
			}
			if strings.Contains(nazwa, szukana) {
				czesciowe = append(czesciowe, zTablicy)
			}
		}
		if len(dokladne) == 1 {
			return dokladne[0], nil
		}
		if len(dokladne) == 0 && len(czesciowe) == 1 {
			return czesciowe[0], nil
		}
		if len(dokladne) == 0 && len(czesciowe) == 0 {
			return shared.StudioSymbol{}, bladWskazaniaStudio(
				"znaku o nazwie „" + strings.TrimSpace(*z.Name) +
					"” tablica znaków nie niesie. Wykaz nazw oddaje komenda " +
					"studio.symbol.list; znak spoza tablicy wstawia się polem code")
		}
		kandydaci := dokladne
		if len(kandydaci) == 0 {
			kandydaci = czesciowe
		}
		nazwy := make([]string, 0, len(kandydaci))
		for _, kandydat := range kandydaci {
			nazwy = append(nazwy, "„"+kandydat.Name+"” ("+kandydat.Character+")")
		}
		if len(nazwy) > 8 {
			nazwy = nazwy[:8]
		}
		return shared.StudioSymbol{}, bladWskazaniaStudio(
			"nazwa „" + strings.TrimSpace(*z.Name) + "” wskazuje więcej niż jeden znak: " +
				strings.Join(nazwy, ", ") + ". Podaj nazwę pełną albo punkt kodowy")
	}

	return shared.StudioSymbol{}, bladWskazaniaStudio(
		"wstawienie znaku bez wskazania znaku — podaj pole character, code albo name")
}

// ZasadyAutozamiany oddaje wykaz zasad autozamiany
// (`studio.symbol.autoreplace.list`).
//
// Wykaz jest zszyty z dwóch miejsc: zasady fabryczne wylicza rdzeń, zasady własne
// i wyłączenia fabrycznych stoją w bazie. Zasada Operatora o tym samym skrócie
// PRZEBIJA fabryczną — łącznie z jej wyłączeniem.
func (a *adapterStudia) ZasadyAutozamiany(ctx context.Context,
	z shared.StudioSymbolAutoreplaceListRequest) (shared.StudioSymbolAutoreplaceListResponse, error) {

	skladnica, err := a.symbolSkladnica()
	if err != nil {
		return shared.StudioSymbolAutoreplaceListResponse{}, err
	}
	wiersze, err := skladnica.ZasadyAutozamiany(ctx)
	if err != nil {
		return shared.StudioSymbolAutoreplaceListResponse{}, bladStudio(err)
	}
	zFabrycznymi := z.IncludeBuiltin == nil || *z.IncludeBuiltin

	wykaz := make([]shared.StudioAutoReplaceRule, 0, len(wiersze))
	for _, wiersz := range wiersze {
		if wiersz.Fabryczna && !zFabrycznymi {
			continue
		}
		wykaz = append(wykaz, shared.StudioAutoReplaceRule{
			Shortcut:    wiersz.Skrot,
			Replacement: wiersz.Zamiennik,
			Enabled:     wiersz.Czynna,
			Builtin:     postacWskaznikPrawdy(wiersz.Fabryczna),
		})
	}
	sort.SliceStable(wykaz, func(i, j int) bool {
		return wykaz[i].Shortcut < wykaz[j].Shortcut
	})
	if len(wykaz) == 0 {
		if !zFabrycznymi {
			return shared.StudioSymbolAutoreplaceListResponse{}, bladWskazaniaStudio(
				"wykaz zasad autozamiany bez zasad fabrycznych jest pusty: Operator nie " +
					"założył jeszcze ani jednej zasady własnej. Zasadę zakłada komenda " +
					"studio.symbol.autoreplace.set")
		}
		return shared.StudioSymbolAutoreplaceListResponse{}, postacBladZaplecza(
			"tabela zasad autozamiany jest pusta, a migracje 368 i 371 wpisują do niej " +
				"zasady fabryczne — baza tego wdrożenia jest niekompletna")
	}
	return shared.StudioSymbolAutoreplaceListResponse{Rules: wykaz}, nil
}

// UstawZasadeAutozamiany zakłada, zmienia, wyłącza albo usuwa zasadę
// autozamiany (`studio.symbol.autoreplace.set`).
//
// Zasady fabrycznej nie da się usunąć — odmowa nazywa powód i wskazuje drogę
// wyjścia (wyłączenie polem `enabled`), wzorem `studio.operation.delete`, który
// robi to samo dla operacji fabrycznych. Nastawa jest jawna i odwracalna: zasada
// wyłączona zostaje w wykazie, więc Operator widzi, że ją wyłączył, a nie że
// zniknęła.
func (a *adapterStudia) UstawZasadeAutozamiany(ctx context.Context,
	z shared.StudioSymbolAutoreplaceSetRequest) (shared.StudioSymbolAutoreplaceSetResponse, error) {

	skrot := strings.TrimSpace(z.Shortcut)
	if skrot == "" {
		return shared.StudioSymbolAutoreplaceSetResponse{}, bladWskazaniaStudio(
			"zasada autozamiany bez skrótu")
	}
	skladnica, err := a.symbolSkladnica()
	if err != nil {
		return shared.StudioSymbolAutoreplaceSetResponse{}, err
	}

	// Czy zasada jest fabryczna, wie BAZA — kolumna `fabryczna` z migracji 368.
	// Rdzeń tego nie zgaduje z własnego wykazu, bo własnego wykazu nie ma.
	zastana, err := skladnica.ZasadaAutozamiany(ctx, skrot)
	zastanaJest := err == nil
	if err != nil && !errors.Is(err, dane.ErrBrakWiersza) {
		return shared.StudioSymbolAutoreplaceSetResponse{}, bladStudio(err)
	}

	// Zamiennik pusty usuwa zasadę — tak mówi kontrakt tego pola.
	if z.Replacement != nil && strings.TrimSpace(*z.Replacement) == "" {
		if !zastanaJest {
			return shared.StudioSymbolAutoreplaceSetResponse{}, bladWskazaniaStudio(
				"zasady „" + skrot + "” nie ma czego usunąć: nie stoi ani w wykazie " +
					"fabrycznym, ani wśród zasad własnych Operatora")
		}
		if zastana.Fabryczna {
			return shared.StudioSymbolAutoreplaceSetResponse{}, bladWskazaniaStudio(
				"zasady „" + skrot + "” nie da się usunąć: jest fabryczna, a wykaz " +
					"fabryczny przyszedł wraz z produktem, nie od Operatora. Wyłącz ją " +
					"polem enabled — wyłączenie jest odwracalne i zostaje widoczne w wykazie")
		}
		usunieta, err := skladnica.UsunZasadeAutozamiany(ctx, skrot)
		if err != nil {
			return shared.StudioSymbolAutoreplaceSetResponse{}, bladStudio(err)
		}
		if !usunieta {
			return shared.StudioSymbolAutoreplaceSetResponse{}, postacBladZaplecza(
				"zasada „" + skrot + "” nie dała się usunąć, choć nie jest fabryczna")
		}
		return shared.StudioSymbolAutoreplaceSetResponse{
			Removed: postacWskaznikPrawdy(true),
		}, nil
	}

	// Zasada wskazana samym skrótem bez zamiennika jest przestawieniem tego, co
	// już stoi — najczęściej wyłączeniem. Zamiennik bierze się wtedy z zasady
	// zastanej; jej brak jest odmową nazwaną, bo zasada bez zamiennika nie ma
	// czego wstawić.
	zamiennik := ""
	switch {
	case z.Replacement != nil:
		zamiennik = *z.Replacement
	case zastanaJest:
		zamiennik = zastana.Zamiennik
	default:
		return shared.StudioSymbolAutoreplaceSetResponse{}, bladWskazaniaStudio(
			"zasada „" + skrot + "” jest nowa i nie niesie zamiennika — podaj pole " +
				"replacement, czyli to, co ma wejść w miejsce skrótu")
	}
	czynna := true
	if z.Enabled != nil {
		czynna = *z.Enabled
	}

	zapisana, err := skladnica.ZapiszZasadeAutozamiany(ctx, dane.ZasadaAutozamianyStudia{
		Skrot: skrot, Zamiennik: zamiennik, Czynna: czynna,
	})
	if err != nil {
		return shared.StudioSymbolAutoreplaceSetResponse{}, bladStudio(err)
	}
	zasada := shared.StudioAutoReplaceRule{
		Shortcut:    zapisana.Skrot,
		Replacement: zapisana.Zamiennik,
		Enabled:     zapisana.Czynna,
		Builtin:     postacWskaznikPrawdy(zapisana.Fabryczna),
	}
	return shared.StudioSymbolAutoreplaceSetResponse{
		Rule: &zasada, Removed: postacWskaznikPrawdy(false),
	}, nil
}
