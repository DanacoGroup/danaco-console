# Danaco Console — Okno historii sesji

| | |
|---|---|
| **Produkt** | Danaco Console |
| **Rodzaj** | Platforma AI Workspace OS |
| **Opis** | Platforma jest wielośrodowiskowym systemem operacyjnym dla sztucznej inteligencji, integrującym komunikację, zarządzanie wiedzą, tworzenie treści, projektowanie, automatyzacje procesów oraz rozwój oprogramowania. |
| **Producent** | Danaco Holding Group Sp. z o.o. |
| **Twórca** | Dariusz Naharnowicz |
| **Wersja** | v2.0 |
| **Status** | Deweloperski |
| **Data** | 2026-08-20 |

**Informacje szczegółowe dokumentu:**

| | |
|---|---|
| **Tytuł** | Okno historii sesji |
| **Klasa dokumentu** | Specyfikacja docelowa |
| **Odbiorcy** | projektant — buduje układ i zachowanie okna · deweloper — wiąże okno z kontraktem i modelem danych |
| **Przeznaczenie** | Ustala kompletną budowę, zawartość i zachowanie okna historii sesji — nakładki okna ustawień otwartej na sekcji rejestru sesji — tak, aby deweloper zbudował je bez rozstrzygania czegokolwiek samodzielnie. |
| **Zakres** | wywołanie okna, układ, trzy zbiory sesji (czynne, zakończone, archiwalne) i ich kolumny, filtry po środowisku, wyszukiwanie, operacje właściwe każdemu zbiorowi, stany kontrolek, skróty, punkty łamania, dostępność, komendy i struktury kontraktu wraz z komendami niewyeksponowanymi w widoku, kryteria odbioru |
| **Poza zakresem** | mechanika samego okna ustawień i pozostałych ośmiu jego sekcji — [Okno ustawień](ustawienia.md); definicja encji sesji i okna — [Model danych](../architektura/model-danych.md); model kart sesji na wstążce — [Rama okna](rama-okna.md) |
| **Dokument nadrzędny** | [Elementy okien](elementy-okien.md) |
| **Dokumenty powiązane** | [Okno ustawień](ustawienia.md) · [Rama okna](rama-okna.md) · [Strona główna i nawigacja](strona-glowna-i-nawigacja.md) · [Model danych](../architektura/model-danych.md) · [Okno nowego projektu](okno-nowego-projektu.md) |
| **Prototypy odniesienia** | `design/05-okna/platformowe/historia-sesji.html` |
| **Źródła normatywne** | `design/05-okna/platformowe/historia-sesji.html` (układ, trzy zbiory, kolumny, przyciski per wiersz, treści dosłowne) · `budowa/shared/contract.json` (obszary `session`, `history`; wyliczenie `SessionStatus`; struktury `Session`, `SessionPresence`, `HistoryEntry`, `Window`, `ErrorInfo`) · `design/zasoby/css/komponenty.css` · `design/zasoby/prototyp.css` (animacja `.pt-tetno`) |
| **Zasada nadrzędna** | Rozłączenie klienta nie kończy sesji ani jej procesów. Sesja ma dokładnie cztery stany wyliczenia `SessionStatus` — `active`, `paused`, `finished`, `archived` — lecz widok zweryfikowany w prototypie grupuje wykaz w trzy zbiory, nie cztery; rozdz. 1.3 rozstrzyga tę rozbieżność wprost, nie milczeniem. |

---

## Spis treści

1. [Czym jest okno historii sesji](#1-czym-jest-okno-historii-sesji)
   - [1.1 Rola okna w platformie](#11-rola-okna-w-platformie)
   - [1.2 Dlaczego jest sekcją okna ustawień](#12-dlaczego-jest-sekcją-okna-ustawień)
   - [1.3 Cztery stany sesji, trzy zbiory widoku](#13-cztery-stany-sesji-trzy-zbiory-widoku)
2. [Umiejscowienie i sposób wywołania](#2-umiejscowienie-i-sposób-wywołania)
3. [Układ okna](#3-układ-okna)
   - [3.1 Szkic układu](#31-szkic-układu)
   - [3.2 Dziewięć sekcji okna ustawień](#32-dziewięć-sekcji-okna-ustawień)
   - [3.3 Strefy treści sekcji](#33-strefy-treści-sekcji)
   - [3.4 Żetony i komponenty](#34-żetony-i-komponenty)
   - [3.5 Punkty łamania](#35-punkty-łamania)
4. [Trzy zbiory sesji](#4-trzy-zbiory-sesji)
   - [4.1 Sesje czynne](#41-sesje-czynne)
   - [4.2 Sesje zakończone](#42-sesje-zakończone)
   - [4.3 Sesje archiwalne](#43-sesje-archiwalne)
   - [4.4 Kolumny wspólne i ich źródła](#44-kolumny-wspólne-i-ich-źródła)
   - [4.5 Reguła podtytułu i stronicowania trzech zbiorów](#45-reguła-podtytułu-i-stronicowania-trzech-zbiorów)
   - [4.6 Nagłówek kolumn jest opisowy, nie interaktywny](#46-nagłówek-kolumn-jest-opisowy-nie-interaktywny)
5. [Filtry i wyszukiwanie](#5-filtry-i-wyszukiwanie)
   - [5.1 Zakres środowisk](#51-zakres-środowisk)
   - [5.2 Wyszukiwanie](#52-wyszukiwanie)
   - [5.3 Zakres dat i sortowanie](#53-zakres-dat-i-sortowanie)
6. [Operacje na sesji](#6-operacje-na-sesji)
   - [6.1 Wróć / Zakończ — sesje czynne](#61-wróć--zakończ--sesje-czynne)
   - [6.2 Otwórz / Archiwizuj — sesje zakończone](#62-otwórz--archiwizuj--sesje-zakończone)
   - [6.3 Przywróć / Pobierz — sesje archiwalne](#63-przywróć--pobierz--sesje-archiwalne)
   - [6.4 Komendy obszaru `session` bez odpowiednika w widoku](#64-komendy-obszaru-session-bez-odpowiednika-w-widoku)
7. [Stany kontrolek](#7-stany-kontrolek)
8. [Skróty klawiszowe](#8-skróty-klawiszowe)
9. [Dostępność](#9-dostępność)
10. [Komendy i struktury kontraktu](#10-komendy-i-struktury-kontraktu)
   - [10.1 Pełny wykaz komend obszaru `session`](#101-pełny-wykaz-komend-obszaru-session)
   - [10.2 Komendy obszaru `history`](#102-komendy-obszaru-history)
   - [10.3 Struktury i wyliczenia](#103-struktury-i-wyliczenia)
   - [10.4 Zdarzenia zasilające widok](#104-zdarzenia-zasilające-widok)
11. [Okno komunikacji jako nośnik historii](#11-okno-komunikacji-jako-nośnik-historii)
   - [11.1 Struktura Window](#111-struktura-window)
   - [11.2 Wyliczenia okna](#112-wyliczenia-okna)
   - [11.3 Historia per okno a widok per sesja](#113-historia-per-okno-a-widok-per-sesja)
   - [11.4 Narzędzia dołożone do sesji](#114-narzędzia-dołożone-do-sesji)
12. [Kody błędów i zachowanie okna](#12-kody-błędów-i-zachowanie-okna)
13. [Przebiegi operacji](#13-przebiegi-operacji)
   - [13.1 Wróć — sesja czynna](#131-wróć--sesja-czynna)
   - [13.2 Otwórz — sesja zakończona](#132-otwórz--sesja-zakończona)
   - [13.3 Archiwizuj i Przywróć](#133-archiwizuj-i-przywróć)
   - [13.4 Pobierz — sesja archiwalna](#134-pobierz--sesja-archiwalna)
   - [13.5 Zakończ — sesja czynna](#135-zakończ--sesja-czynna)
14. [Scenariusze Operatora](#14-scenariusze-operatora)
15. [Powiązania z projektami i pamięcią](#15-powiązania-z-projektami-i-pamięcią)
16. [Terminologia okna](#16-terminologia-okna)
17. [Kryteria odbioru](#17-kryteria-odbioru)
18. [Załącznik A. Tabela zbiorcza — kontrolka, komenda, źródło](#załącznik-a-tabela-zbiorcza--kontrolka-komenda-źródło)
19. [Załącznik B. Słowniczek pojęć okna](#załącznik-b-słowniczek-pojęć-okna)

---

## 1. Czym jest okno historii sesji

### 1.1 Rola okna w platformie

Okno historii sesji jest jedynym miejscem platformy, w którym Operator widzi **pełny rejestr sesji
swojego konta** — niezależnie od stanowiska, na którym powstały, środowiska, w którym są
prowadzone, i tego, czy proces sesji nadal trwa na rdzeniu. Prototyp
`design/05-okna/platformowe/historia-sesji.html` deklaruje siebie atrybutem `data-prototyp-okno`
jako „Historia sesji — nakładka okna ustawień z otwartą sekcją rejestru”.

Sesja w Danaco Console jest bytem trwałym po stronie rdzenia (struktura `Session` w kontrakcie).
Klient — okno aplikacji na stanowisku — jest wyłącznie oknem na sesję. Rozłączenie klienta jej nie
kończy. Z tej zasady wynika racja istnienia okna: skoro sesje żyją dłużej niż połączenie klienta,
Operator potrzebuje widoku, który je zbiera, pozwala odnaleźć sesję sprzed tygodnia, wrócić do niej,
zakończyć ją, otworzyć zakończoną, zarchiwizować albo przywrócić z archiwum.

Okno nie prowadzi pracy — jest **rejestrem i pulpitem operacji na sesjach**. Praca dzieje się
w oknach operacyjnych modułów; tutaj Operator odnajduje sesję i decyduje, co z nią zrobić. Panel
wprowadzający okna (`.hs-tlo-karta`, widoczny zanim Operator otworzy właściwą nakładkę) niesie tę
samą myśl wprost, cytatem dosłownym:

> „Rejestr sesji otwiera się jako nakładka nad bieżącym widokiem — Operator zagląda do niego bez
> opuszczania miejsca pracy i wraca jednym zamknięciem.”

Panel wprowadzający jest elementem wyłącznie prototypowym — dema składu okna przed pokazaniem
właściwej nakładki, obecny w pliku `historia-sesji.html` jako `.pt-okno` osobny od dialogu
`.hs-okno`. Nie jest częścią budowanego produktu: w gotowej aplikacji cztery drogi wejścia z rozdz. 2
otwierają nakładkę bezpośrednio, bez pośredniego ekranu wyjaśniającego. Cytat powyżej jest jednak
przywołany jako dosłowne źródło opisu roli okna, nie jako element interfejsu do zbudowania.

### 1.2 Dlaczego jest sekcją okna ustawień

Historia sesji nie jest osobnym typem okna. Jest **sekcją okna ustawień**, otwartą bezpośrednio na
rejestrze. Nakłada się na powłokę jako okno nakładkowe (`.dn-modal`, lokalnie rozszerzone klasą
`.hs-okno`), z tytułem nagłówka nakładki „Ustawienia i konfiguracja” i sekcją „Historia sesji”
zaznaczoną atrybutem `aria-current="true"` na lewej liście sekcji. Nakładka dziedziczy mechanikę
okna ustawień: przeciąganie za nagłówek, zmianę rozmiaru, zwinięcie do nagłówka, maksymalizację oraz
lewą kolumnę nawigacji sekcji (rozdz. 3.2). Kontrakt potwierdza to umiejscowienie: komenda
`session.archive.list` opisana jest w `contract.json` jako „wgląd z okna ustawień”.

Konsekwencja projektowa: **okno historii sesji nigdy nie zajmuje całego obszaru roboczego** i nie
zastępuje okna operacyjnego. Otwiera się nad nim, a po zamknięciu przywraca widok, z którego zostało
wywołane.

Bycie sekcją, nie osobnym oknem, ma też konsekwencję dla dokumentacji: pełny opis mechaniki wspólnej
wszystkim dziewięciu sekcjom (rozdz. 3.2) — przeciąganie, zmiana rozmiaru, zwinięcie, maksymalizacja
— należy do [Okna ustawień](ustawienia.md), nie do niniejszego dokumentu. Niniejszy dokument opisuje
wyłącznie to, co jest właściwe sekcji „Historia sesji” i niepowtarzalne w żadnej z pozostałych ośmiu:
trzy zbiory sesji, filtry środowiska, sześć operacji jednowierszowych i ich mapowanie na komendy
kontraktu.

### 1.3 Cztery stany sesji, trzy zbiory widoku

Stan sesji niesie wyliczenie `SessionStatus` z kontraktu. Ma **cztery wartości** — `active`,
`paused`, `finished`, `archived`.

```
 session.create
        │
        ▼
   ┌──────────┐  bezczynność / rozłączenie   ┌──────────┐
   │  active  │ ────────────────────────────►│  paused  │
   │ czynna   │◄──────── session.bind ───────│ wstrzym. │
   │          │                              └────┬─────┘
   └────┬─────┘                                   │
        │ session.close / session.stop            │ session.resume
        ▼                                         ▼
   ┌──────────────┐   session.resume / open  ┌──────────┐
   │  finished    │◄─────────────────────────│  active  │
   │  zakończona  │                          └──────────┘
   └──────┬───────┘
          │ session.archive           ┌────────────┐
          └──────────────────────────►│  archived  │
                session.restore ◄─────│ archiwum   │
                                      └────────────┘
```

Legenda: strzałki podpisane są nazwami komend kontraktu wykonujących przejście. Wartości `active`,
`paused`, `finished`, `archived` pochodzą wprost z wyliczenia `SessionStatus`.

**Rozbieżność zweryfikowana ze źródłem prototypu.** Wykaz sesji zweryfikowany w prototypie
(rozdz. 4) nie grupuje sesji w cztery zbiory odpowiadające czterem wartościom `SessionStatus` —
grupuje je w **trzy zbiory** (`.hs-zbior`, atrybut `data-stan` na wierszu): „Sesje czynne”
(`data-stan="czynna"`), „Sesje zakończone” (`data-stan="zakonczona"`), „Sesje archiwalne”
(`data-stan="archiwalna"`). Żaden wiersz przykładowy nie niesie `data-stan` odpowiadającego stanowi
`paused`, i arkusz stylów lokalnych (`.hs-wiersz[data-stan=...]`) definiuje reguły wyłącznie dla
dwóch wartości: `czynna` i `archiwalna` — `zakonczona` nie ma własnej reguły barwy, bo dziedziczy
wygląd domyślny wiersza.

| `SessionStatus` (kontrakt) | Zbiór widoku (prototyp) | Status weryfikacji |
|---|---|---|
| `active` | „Sesje czynne” | zweryfikowane — trzy wiersze przykładowe |
| `paused` | brak zbioru dedykowanego | **[DO DECYZJI OPERATORA]** — czy sesje `paused` wchodzą do zbioru „Sesje czynne” (skoro proces bywa nieaktywny, lecz zapis kompletny), do zbioru „Sesje zakończone”, czy otrzymują własny, czwarty zbiór nieobecny w zweryfikowanym zrzucie |
| `finished` | „Sesje zakończone” | zweryfikowane — cztery wiersze przykładowe |
| `archived` | „Sesje archiwalne” | zweryfikowane — trzy wiersze przykładowe |

Zasada, której poprzednia redakcja tego dokumentu się trzymała i którą niniejsza redakcja
utrzymuje jako wymóg projektowy niezależnie od rozbieżności powyżej: **sesja `paused` nie może być
przedstawiona jako `finished`.** Wstrzymanie jest odwracalne — Operator, który zamknął okno
aplikacji, po ponownym otwarciu zastaje swoje sesje wstrzymane, gotowe do wznowienia, nie
zakończone. Rozstrzygnięcie, do którego z trzech zbiorów zweryfikowanych w prototypie sesja
`paused` ma trafiać, pozostaje otwarte (wiersz tabeli powyżej) — nie wolno jej jednak
utożsamiać automatycznie ze zbiorem „Sesje zakończone”.

---

## 2. Umiejscowienie i sposób wywołania

| Droga | Miejsce | Dosłowne brzmienie kontrolki |
|---|---|---|
| Menu aplikacji | Przejdź › Historia sesji | „Historia sesji” |
| Szyna nawigacji, strefa pracy | przycisk osobny, obok „Nowa sesja” i „Nowy projekt” | „Historia sesji”, etykieta wskazania kursorem: „Pełny rejestr sesji konta: czynne, zakończone i archiwalne.” |
| Przedsionek środowiska | szyna sesji przedsionka | „Pełna historia sesji” |
| Okno ustawień już otwarte | sekcja „Historia sesji” listy sekcji | „Historia sesji” |

Trzy z czterech dróg wejścia wykazanych powyżej — szyna nawigacji, przedsionek środowiska
i sekcja otwartego okna ustawień — są potwierdzone dosłownie w prototypie, cytatem z panelu
wprowadzającego okna; czwarta, pozycja menu aplikacji, pochodzi z wykazu poleceń menu powłoki:

> „Wywołują go: pozycja „Historia sesji” w strefie pracy szyny nawigacji, przycisk „Pełna historia
> sesji” w szynie sesji przedsionka oraz sekcja „Historia sesji” w oknie ustawień i konfiguracji.”

| Warunek wywołania | Zachowanie |
|---|---|
| Wywołanie z dowolnej z czterech dróg powyżej | nakładka otwiera się z sekcją „Historia sesji” już wybraną (`aria-current="true"`) |
| Okno ustawień już otwarte na innej sekcji | przełączenie lewej kolumny na sekcję „Historia sesji” bez ponownego otwierania nakładki |

Otwarcie rejestru jest odczytem (`session.list`) — nie zmienia stanu żadnej sesji. Dopiero jawna
operacja właściwa jednemu z trzech zbiorów (rozdz. 6) zmienia stan.

Szyna nawigacji niesie pozycję „Historia sesji” w tej samej strefie pracy co „Nowa sesja” i „Nowy
projekt” (rozdz. 2.1–2.2 opracowania [Okno nowego projektu](okno-nowego-projektu.md)) — trzy pozycje
równorzędne strefy pracy, z których wyłącznie „Historia sesji” otwiera nakładkę zamiast pełnego
okna Workspace albo przedsionka.

Trzy pozycje strefy pracy różnią się więc nie tylko celem, ale i klasą okna, które otwierają:
„Nowa sesja” prowadzi do przedsionka (byt bez ramy nakładkowej), „Nowy projekt” do pełnego okna
Workspace w trybie zakładania ([Okno nowego projektu](okno-nowego-projektu.md) rozdz. 1.1), a
„Historia sesji” — jako jedyna z trzech — do nakładki (`.dn-modal`) opisanej w niniejszym rozdziale.
Ten kontrast jest wart odnotowania przy dalszej budowie szyny nawigacji: trzy pozycje sąsiadujące
wizualnie w tej samej strefie prowadzą do trzech odmiennych klas okien platformy.

---

## 3. Układ okna

### 3.1 Szkic układu

```
┌───────────────────────────────────────────────────────────────────────────────┐
│  Ustawienia i konfiguracja                                          [_][□][×] │  .dn-modal-naglowek
├────────────────────────┬──────────────────────────────────────────────────────┤
│  ZAKRES KONTA          │  Historia sesji                                      │
│  ○ Konto Operatora     │  Rejestr obejmuje sesje ze wszystkich czterech       │
│  ○ Uwierzytelnianie    │  środowisk pracy — także te, które trwają w tle…     │
│  ○ Wygląd i język      │                                                      │
│  ○ Urządzenia          │  🔍 Szukaj w nazwach, projektach i modułach…          │
│  ○ Powiadomienia       │  [Wszystkie][TalkIn][WorkSpace][CodeStudio][Multi…]  │
│  ● Historia sesji      │                                    [Zakres dat][Sort]│
│  ○ Konta modeli        │  ────────────────────────────────────────────────    │
│  ○ Izolacja kontekstu  │   Sesja          Środ.  Moduł  Ostatnia  Zapis       │
│  ○ Pamięć i dane lok.  │  Sesje czynne · 3 sesje · praca trwa po stronie serw.│
│                        │  ● Raport końcowy…  WorkSp. Studio  12:04   7 wersji │
│                        │    …                                    [Wróć][Zakoń]│
│                        │  Sesje zakończone · 4 sesje · ostatnie 7 dni         │
│                        │  ✓ Streszczenie zarz. WorkSp. Studio  dziś10:31 42min│
│                        │    …                                  [Otwórz][Arch.]│
│                        │  Sesje archiwalne · 3 sesje · archiwum konta         │
│                        │  🗃 Kary umowne…    WorkSp. Studio  12 sierp archiw.  │
│                        │    …                              [Przywróć][Pobierz]│
└────────────────────────┴──────────────────────────────────────────────────────┘
```

Legenda: lewa kolumna to nawigacja dziewięciu sekcji okna ustawień (`●` — sekcja zaznaczona,
`○` — pozostałe; rozdz. 3.2). Prawa kolumna niesie nagłówek sekcji, pasek filtrów i trzy zbiory
sesji, każdy z własnym nagłówkiem grupy i parą przycisków akcji przy każdym wierszu — nie jeden
wspólny pasek operacji na dole, jak zakładała wcześniejsza redakcja tego dokumentu.

### 3.2 Dziewięć sekcji okna ustawień

Lewa kolumna (`.hs-sekcje`, szerokość 260 px) niesie nagłówek grupy „Zakres konta”
(`.hs-sekcje-glowa`) i dziewięć pozycji (`.hs-sekcja`), każda z ikoną, nazwą pogrubioną i
jednowierszowym opisem (`.hs-meta`). Pełny wykaz dziewięciu pozycji, w kolejności zweryfikowanej w
prototypie:

| # | Nazwa sekcji | Opis (`.hs-meta`) |
|---|---|---|
| 1 | Konto Operatora | „Login, e-mail, hasło” |
| 2 | Uwierzytelnianie | „Metody logowania” |
| 3 | Wygląd i język | „Motyw, język, gęstość” |
| 4 | Urządzenia | „Wykaz i parowanie” |
| 5 | Powiadomienia | „Kategorie i kanały” |
| 6 | **Historia sesji** | „Rejestr sesji konta” — sekcja bieżąca (`aria-current="true"`), przedmiot niniejszego dokumentu |
| 7 | Konta modeli | „Kanał CLI i profile” |
| 8 | Izolacja kontekstu | „Profile zasięgu sesji” |
| 9 | Pamięć i dane lokalne | „Pamięć podręczna, kopie” |

Osiem sekcji poza „Historią sesji” jest poza zakresem niniejszego dokumentu — pełny opis okna ustawień
i wszystkich jego sekcji należy do [Okna ustawień](ustawienia.md). Wykaz powyższy jest tu
przywołany wyłącznie jako mapa nawigacyjna kontekstu, w którym sekcja „Historia sesji” się znajduje.

### 3.3 Strefy treści sekcji

| Strefa | Klasa | Zawartość |
|---|---|---|
| Nagłówek sekcji | `.hs-tytul`, `.hs-lid` | tytuł „Historia sesji” i wprowadzenie (cytat pełny w rozdz. 4) |
| Pasek filtrów | `.hs-filtry` | pole wyszukiwania, zakładki zakresu środowisk, przyciski „Zakres dat” i „Sortowanie” (rozdz. 5) |
| Rejestr | `.hs-rejestr` | nagłówek kolumn (dekoracyjny, rozdz. 4.6) i trzy zbiory sesji (rozdz. 4) |
| Nota zamykająca | `.hs-nota` | wyjaśnienie skutków zakończenia i archiwizacji (cytat pełny w rozdz. 15) |

Cała treść sekcji (`.hs-tresc`) przewija się jako jedna kolumna pionowa — w odróżnieniu od
poprzedniej redakcji tego dokumentu, prototyp nie wyodrębnia paska filtrów i paska operacji jako
elementów nieruchomych względem przewijanego wykazu; `.hs-tresc` ma `overflow-y: auto` na całej
swojej wysokości.

### 3.4 Żetony i komponenty

Nazwy klas potwierdzone w `design/05-okna/platformowe/historia-sesji.html`,
`design/zasoby/css/komponenty.css` i `design/zasoby/prototyp.css`.

| Element | Klasa / żeton | Rola |
|---|---|---|
| Nakładka okna | `.dn-modal.hs-okno` | okno nakładkowe, rozszerzone lokalnym prefiksem `.hs-*` |
| Nagłówek nakładki | `.dn-modal-naglowek` | pas przeciągania i sterowania oknem, tytuł „Ustawienia i konfiguracja” |
| Ciało nakładki | `.dn-modal-cialo` | obszar dwukolumnowy `.hs-korpus` |
| Siatka korpusu | `.hs-korpus` | `grid-template-columns: 260px minmax(0, 1fr)` |
| Kolumna sekcji | `.hs-sekcje` | nawigacja dziewięciu sekcji (rozdz. 3.2) |
| Pozycja sekcji | `.hs-sekcja` | jedna z dziewięciu; `aria-current='true'` na bieżącej |
| Pole wyszukiwania | `.dn-szukaj` z `.dn-pole-kontrolka` | zawężanie wykazu (rozdz. 5.2) |
| Zakładki zakresu środowisk | `.dn-zakladki.dn-zakladki--pigulki`, `.dn-zakladka`, `role="tablist"`/`role="tab"` | pięć pozycji (rozdz. 5.1) |
| Przyciski „Zakres dat” / „Sortowanie” | `.dn-btn.dn-btn--duch.dn-btn--sm` | rozdz. 5.3 |
| Siatka wiersza i nagłówka | zmienna CSS `--hs-siatka` (`18px minmax(0,1fr) 104px 92px 96px 84px 152px`) | siedem kolumn wspólnych nagłówkowi i każdemu wierszowi |
| Nagłówek grupy zbioru | `.hs-zbior-glowa`, `.hs-zbior-tytul`, `.hs-zbior-liczba`, `.hs-zbior-kreska` | tytuł, licznik, kreska wypełniająca |
| Wiersz sesji | `.hs-wiersz`, atrybut `data-stan` | `czynna` / `zakonczona` / `archiwalna` (rozdz. 1.3, 4) |
| Znak wiersza | `.hs-znak`; sesja czynna dodatkowo `.pt-tetno` | ikona stanu; tętnienie wyłącznie dla `czynna` |
| Nazwa sesji | `.hs-nazwa` z `<b>` i `<small>` | tytuł sesji i podtytuł (projekt/plik/zakres) |
| Metadana kolumnowa | `.hs-meta`, wariant `.hs-ukryj-waskie` | wartości kolumn ukrywane poniżej 900 px |
| Przyciski akcji wiersza | `.hs-akcje` z `.dn-btn.dn-btn--duch.dn-btn--sm` | dwa przyciski, różne wg zbioru (rozdz. 6) |
| Animacja tętna | `.pt-tetno` (`design/zasoby/prototyp.css`), niezależna od `.dn-kropka--tetno` (`komponenty.css`) | znak sesji czynnej; obie klasy dają wizualnie zbliżony efekt, lecz są zdefiniowane osobno — wiersz sesji korzysta z pierwszej |

**Wymiary i tło nakładki bazowej.** Reguła `.dn-modal` w `komponenty.css` (niedziedziczona lokalnie
przez `.hs-okno` poza tym, co rozdz. 3.5 opisuje jako przycięcie na `900px`) ustala: szerokość
`min(var(--dn-wym-modal), calc(100vw - var(--dn-od-8)))`, wysokość maksymalną
`min(80dvh, 720px)`, obrys `--dn-obrys`, zaokrąglenie `--dn-r-lg`, tło `--dn-panel` (ten sam żeton
zweryfikowany dla panelu podsumowania opracowania [Okno nowego projektu](okno-nowego-projektu.md) rozdz. 6 —
jaśniejszy od bieli w motywie jasnym, jaśniejszy od tła kart w motywie ciemnym) i cień
`--dn-cien-lg`. Warstwa pod nakładką (`::backdrop`) korzysta z żetonu `--dn-nakladka` z rozmyciem
`blur(2px)`. Wejście nakładki animuje się przejściem `dn-wejscie` w czasie `--dn-czas-3`: przesunięcie
pionowe o `--dn-od-2` i skalowanie z `0.98` do `1`, jednocześnie z pojawieniem przezroczystości od
`0`.

### 3.5 Punkty łamania

| Próg szerokości | Zachowanie |
|---|---|
| > 900 px | układ dwukolumnowy pełny; siedem kolumn wiersza (`--hs-siatka: 18px minmax(0,1fr) 104px 92px 96px 84px 152px`); lewa kolumna sekcji pionowa |
| ≤ 900 px | `.hs-korpus` przechodzi na jedną kolumnę; lewa kolumna sekcji staje się poziomym paskiem przewijanym (`overflow-x: auto`); siatka wiersza zwęża się do czterech kolumn (`--hs-siatka: 20px minmax(0,1fr) 104px 168px`); kolumny „Środowisko”, „Moduł” i „Zapis” znikają (`.hs-ukryj-waskie { display: none }`) |

Próg 900 px nie odpowiada żadnemu żetonowi `--dn-bp-*` nazwanemu w `zetony.css` — najbliższy niższy
jest `--dn-bp-w1` (640 px), najbliższy wyższy `--dn-bp-w2` (960 px). Wartość jest zapisana w arkuszu
lokalnym wprost, jako liczba pikseli, analogicznie do progów zweryfikowanych w
[Oknie instalatora](okno-instalatora.md) (1000 px) i [Oknie nowego projektu](okno-nowego-projektu.md)
(1240 px) — trzeci przykład tego samego wzorca w zbiorze: progi lokalne okien nie pokrywają się z
progami nazwanymi `--dn-bp-*` ramy aplikacji.

---

## 4. Trzy zbiory sesji

Wprowadzenie sekcji (`.hs-lid`), cytat dosłowny:

> „Rejestr obejmuje sesje ze wszystkich czterech środowisk pracy — także te, które trwają w tle po
> zamknięciu aplikacji. Sesja czynna wraca do pracy w stanie, w jakim została pozostawiona; sesja
> zakończona otwiera się jako nowa sesja z odtworzonym kontekstem; sesja archiwalna wymaga
> wcześniejszego przywrócenia.”

Trzy zdania tego wprowadzenia odpowiadają dokładnie trzem zbiorom poniżej — każdy niesie regułę
odrębną co do tego, co znaczy „wrócić” do sesji tego zbioru.

### 4.1 Sesje czynne

| Cecha | Wartość |
|---|---|
| Nagłówek grupy | „Sesje czynne” |
| Licznik grupy (zrzut) | „3 sesje · praca trwa po stronie serwera” |
| Znak wiersza | `.hs-znak.pt-tetno` — tętniący, barwa `--dn-sygnal` |
| Przyciski wiersza | „Wróć”, „Zakończ” (rozdz. 6.1) |

Trzy wiersze przykładowe:

| Nazwa (pogrubiona) | Podtytuł | Środowisko | Moduł | Ostatnia praca | Zapis |
|---|---|---|---|---|---|
| Raport końcowy — redakcja | Sprawozdawczość Q3 · raport-koncowy.md | WorkSpace | Studio | 12:04 | 7 wersji |
| Klasyfikacja umów według rodzaju świadczenia | Rejestr umów serwisowych · 240 pozycji | TalkIn | Research | 11:52 | w toku |
| Migracja bazy — dokumentacja zmian | Wdrożenie 2.1 · gałąź release/2.1 | CodeStudio | Developer | 11:20 | w toku |

### 4.2 Sesje zakończone

| Cecha | Wartość |
|---|---|
| Nagłówek grupy | „Sesje zakończone” |
| Licznik grupy (zrzut) | „4 sesje · ostatnie 7 dni” |
| Znak wiersza | fajka (ikona potwierdzenia), bez tętnienia |
| Przyciski wiersza | „Otwórz”, „Archiwizuj” (rozdz. 6.2) |

Cztery wiersze przykładowe:

| Nazwa (pogrubiona) | Podtytuł | Środowisko | Moduł | Ostatnia praca | Zapis |
|---|---|---|---|---|---|
| Streszczenie zarządcze | Sprawozdawczość Q3 · 4 strony | WorkSpace | Studio | dziś 10:31 | 42 min |
| Przekład noty technicznej DE → PL | Materiały regulacyjne · 11 stron | TalkIn | Translate | wczoraj 16:08 | 1 g 12 min |
| Przegląd zmian w prawie energetycznym | Materiały regulacyjne · 18 źródeł | TalkIn | Browser | wczoraj 09:47 | 55 min |
| Konfiguracja zespołu modeli — narada | Wdrożenie 2.1 · 4 role | MultitaskingAI | Zespoły | 2 dni temu | 2 g 05 min |

Wartość kolumny „Moduł” dla ostatniego wiersza — „Zespoły” — jest jedyną wartością tej kolumny w
całym zrzucie niebędącą nazwą modułu z katalogu modułów opisanego w
[Stronie głównej i nawigacji](strona-glowna-i-nawigacja.md); w środowisku MultitaskingAI boczna
nawigacja jest panelem orkiestracji, nie listą modułów (`NavigationKind.orchestration`), więc
„Zespoły” oznacza tu sekcję panelu orkiestracji, nie modułu w rozumieniu struktury `Module`.

### 4.3 Sesje archiwalne

| Cecha | Wartość |
|---|---|
| Nagłówek grupy | „Sesje archiwalne” |
| Licznik grupy (zrzut) | „3 sesje · archiwum konta” |
| Znak wiersza | ikona archiwum (teczka), bez tętnienia |
| Kolor nazwy | `.hs-wiersz[data-stan='archiwalna'] .hs-nazwa b` — kolor `--dn-tekst-2` zamiast `--dn-tekst` domyślnego, jedyny zbiór z przygaszoną nazwą |
| Przyciski wiersza | „Przywróć”, „Pobierz” (rozdz. 6.3) |

Trzy wiersze przykładowe:

| Nazwa (pogrubiona) | Podtytuł | Środowisko | Moduł | Ostatnia praca | Zapis |
|---|---|---|---|---|---|
| Kary umowne — przegląd zdarzeń | Rejestr umów serwisowych | WorkSpace | Studio | 12 sierpnia | archiwum |
| Opis produktu — zmiana tonu | Marketing | TalkIn | Studio | 4 sierpnia | archiwum |
| Regulamin usługi — korekta | Dział prawny | WorkSpace | Studio | 29 lipca | archiwum |

Kolumna „Zapis” tych trzech wierszy niesie stałą wartość „archiwum” zamiast wartości liczbowej
(liczby wersji albo czasu trwania jak w dwóch pozostałych zbiorach) — dla sesji zarchiwizowanej
metryka pracy przestaje być istotna; istotne jest wyłącznie miejsce, w którym zapis spoczywa. Trzy
wiersze archiwalne mają też wspólną cechę tematyczną w danych przykładowych — wszystkie trzy tytuły
niosą charakter dokumentów zamkniętych i rozliczonych („przegląd zdarzeń”, „zmiana tonu”,
„korekta”), spójnie z naturą archiwum jako miejsca dla pracy definitywnie skończonej.

### 4.4 Kolumny wspólne i ich źródła

Siedem kolumn wspólnych trzem zbiorom (licząc kolumnę znaku i kolumnę akcji). Źródła pól wskazane
wprost — struktury `Session` i `SessionPresence` z kontraktu.

| Kolumna | Nagłówek dekoracyjny | Pole źródłowe | Struktura |
|---|---|---|---|
| Znak | *(brak nagłówka tekstowego)* | `status` (rodzaj ikony) + `live` (tętnienie) | `Session` + `SessionPresence` |
| Sesja | „Sesja” | `title` (pogrubione) + podtytuł swobodny (projekt/plik/zakres) | `Session.title` + kontekst pochodny |
| Środowisko | „Środowisko” | `environmentCode` | `SessionPresence` |
| Moduł | „Moduł” | `moduleCode` | `SessionPresence` |
| Ostatnia praca | „Ostatnia praca” | `updatedAt` / `lastActivityAt` | `Session` / `SessionPresence` |
| Zapis | „Zapis” | metryka pochodna (liczba wersji, czas trwania albo „archiwum”) | pochodna, bez pola 1:1 w kontrakcie — rozdz. 6.4 |
| *(akcje)* | *(brak nagłówka tekstowego)* | — | — |

Kolumna „Ostatnia praca” prezentuje czas względnie („12:04” dla dziś, „wczoraj 16:08”, „2 dni temu”)
do progu kilku dni; powyżej — datą przybliżoną („12 sierpnia”). Znaczniki `createdAt`/`updatedAt` są
typu `int64`; format prezentacji względnej nie ma osobnej struktury w kontrakcie — jest regułą
prezentacji klienta.

**Progi zweryfikowane w dziesięciu wierszach przykładowych.** Cztery formy odrębne występują w
zrzucie, każda przypisana konsekwentnie do zakresu czasu:

| Forma | Przykład ze zrzutu | Zakres odwzorowany |
|---|---|---|
| Sama godzina | „12:04”, „11:52”, „11:20” | dzisiaj — bez dopisku dnia |
| „dziś” + godzina | „dziś 10:31” | wariant jawny formy „sama godzina”, użyty w zbiorze „Sesje zakończone” — **[DO DECYZJI OPERATORA]**, czy „dziś” poprzedza godzinę także w zbiorze „Sesje czynne”, gdzie zrzut pokazuje wyłącznie samą godzinę bez tego dopisku |
| „wczoraj” + godzina | „wczoraj 16:08”, „wczoraj 09:47” | dzień poprzedzający |
| „N dni temu” | „2 dni temu” | od dwóch dni wzwyż, do progu nieustalonego |
| Data przybliżona | „12 sierpnia”, „4 sierpnia”, „29 lipca” | zbiór „Sesje archiwalne” wyłącznie — żaden wiersz spoza tego zbioru nie używa tej formy w zrzucie |

Rozbieżność formy „dziś 10:31” (zbiór zakończone) wobec „12:04” bez dopisku (zbiór czynne) dla tego
samego dnia kalendarzowego nie ma wyjaśnienia w prototypie — oba wiersze pochodzą z tego samego dnia
zrzutu, lecz różnią się obecnością słowa „dziś”. **[DO DECYZJI OPERATORA]**, czy reguła prezentacji
różni się między zbiorami, czy jest to niespójność samego zestawu danych przykładowych, którą
przyszła redakcja tego dokumentu powinna ujednolicić po konsultacji z prototypem zaktualizowanym.

**Pole „Zapis” nie ma odpowiednika 1:1 w kontrakcie.** Wartości „7 wersji”, „w toku”, „42 min”,
„1 g 12 min”, „archiwum” są metrykami różnej natury w zależności od zbioru — licznik wersji, znacznik
stanu, czas trwania, stała tekstowa. Żadna ze struktur `Session`, `SessionPresence` ani
`WorkspaceDashboard` (zweryfikowana przy okazji opracowania [Okno nowego projektu](okno-nowego-projektu.md)
rozdz. 12.7) nie niesie pola nazwanego wprost tak ogólnie. **[DO DECYZJI OPERATORA]** — czy kolumna
„Zapis” czerpie z pola istniejącego pod inną nazwą, na przykład `SessionPresence.openWindowCount`
dla wartości „w toku”, czy wymaga nowego pola po stronie kontraktu.

### 4.5 Reguła podtytułu i stronicowania trzech zbiorów

**Reguła podtytułu.** Dziesięć wierszy przykładowych (rozdz. 4.1–4.3) niosą podtytuł
(`.hs-nazwa small`) o postaci zmiennej, lecz wzorzec powtarza się konsekwentnie: pierwszy człon
nazywa kontekst nadrzędny (nazwa projektu, np. „Sprawozdawczość Q3”, „Wdrożenie 2.1”, „Rejestr umów
serwisowych”, albo dziedzina, np. „Marketing”, „Dział prawny”, „Materiały regulacyjne”), drugi człon
— po znaku `·` — niesie szczegół specyficzny dla sesji (nazwa pliku, liczba pozycji, liczba stron,
gałąź kontroli wersji, liczba ról). Dwa wiersze („Kary umowne — przegląd zdarzeń”, „Opis produktu —
zmiana tonu”) niosą wyłącznie pierwszy człon, bez drugiego — obie należą do zbioru „Sesje
archiwalne”, co sugeruje, że szczegół drugiego członu traci znaczenie dla sesji dawno
zarchiwizowanych. **[DO DECYZJI OPERATORA]** — czy brak drugiego członu jest regułą zbioru
archiwalnego, czy przypadkiem doboru danych przykładowych; trzeci wiersz archiwalny („Regulamin
usługi — korekta” / „Dział prawny”) ma tylko pierwszy człon również, co wzmacnia pierwszą
interpretację, lecz żadna z trzech nie ma zrzutu potwierdzającego regułę wprost.

**Stronicowanie.** `session.list` (rozdz. 10.1) niesie pola `limit` i `offset`; `session.archive.list`
niesie własną parę `offset`/`limit`, odrębną od pierwszej. Trzy zbiory widoku (rozdz. 4.1–4.3) są
więc stronicowane niezależnie od siebie: przewinięcie zbioru „Sesje zakończone” do kolejnej strony
nie wywołuje ponownego odczytu zbioru „Sesje czynne” ani „Sesje archiwalne”, bo każdy pochodzi z
osobnego wywołania kontraktu (`session.list` z filtrem `status` dla dwóch pierwszych zbiorów,
`session.archive.list` dla trzeciego). Liczniki nagłówków grup (rozdz. 4.1–4.3: „3 sesje”, „4
sesje”, „3 sesje”) odpowiadają polu `total` zwracanemu przez wywołanie właściwe każdemu zbiorowi, nie
liczbie wierszy faktycznie wyrenderowanych — przy wykazie dłuższym niż limit stronicowania nagłówek
grupy pokazuje liczbę całkowitą, wykaz poniżej pokazuje podzbiór pierwszej strony.

### 4.6 Nagłówek kolumn jest opisowy, nie interaktywny

Wiersz nagłówka kolumn (`.hs-naglowek-kolumn`) niesie atrybut `aria-hidden="true"` w prototypie —
jest elementem czysto dekoracyjnym, nieinteraktywnym, ukrytym przed technologią wspomagającą.
**Poprzednia redakcja tego dokumentu twierdziła, że „kliknięcie nagłówka kolumny sortuje wykaz” —
zweryfikowany prototyp temu przeczy:** nagłówek nie niesie znaczników interaktywności (`role`,
`tabindex`, uchwytu zdarzeń) i jest jawnie wyłączony z drzewa dostępności. Sortowanie w zweryfikowanym
widoku jest funkcją osobnego przycisku „Sortowanie” (rozdz. 5.3), nie nagłówka kolumny.

Ta korekta ma konsekwencję dla rozdz. 8 (skróty klawiszowe) i rozdz. 9 (dostępność): żaden skrót
klawiszowy tego okna nie jest przypisany do nagłówka kolumny, bo nagłówek nie przyjmuje ogniska —
kolejność `Tab` (rozdz. 8) pomija go całkowicie, przechodząc z paska filtrów wprost do pierwszego
wiersza pierwszego zbioru.

---

## 5. Filtry i wyszukiwanie

### 5.1 Zakres środowisk

Pasek filtrów niesie zakładki pigułkowe (`.dn-zakladki--pigulki`, `role="tablist"`) filtrujące wykaz
po **środowisku**, nie po stanie sesji — poprzednia redakcja tego dokumentu opisywała zakładki
„Czynne / Wszystkie / Zakończone” filtrujące po `SessionStatus`; zweryfikowany prototyp niesie pięć
zakładek odmiennej natury:

| Zakładka | `role` | `aria-selected` domyślnie |
|---|---|---|
| Wszystkie | `tab` | `true` (wybrana domyślnie) |
| TalkIn | `tab` | `false` |
| WorkSpace | `tab` | `false` |
| CodeStudio | `tab` | `false` |
| MultitaskingAI | `tab` | `false` |

Cztery nazwy środowisk odpowiadają dokładnie czterem środowiskom platformy z
[Strony głównej i nawigacji](strona-glowna-i-nawigacja.md). Filtr środowiska działa **w poprzek**
trzech zbiorów jednocześnie (rozdz. 4) — zawężenie do „TalkIn” pokazuje sesje TalkIn należące do
wszystkich trzech zbiorów naraz, nie zastępuje ich jednym płaskim wykazem.

### 5.2 Wyszukiwanie

| Cecha | Wartość |
|---|---|
| Komponent | `.dn-szukaj` z `.dn-pole-kontrolka` |
| Tekst podpowiedzi | „Szukaj w nazwach, projektach i modułach…” |
| Etykieta dostępności | `aria-label="Szukaj w historii sesji"` |
| Zasięg zawężenia | tytuł sesji, nazwa projektu, nazwa modułu — trzy pola przywołane wprost w tekście podpowiedzi |

Pole wyszukiwania zawęża wykaz w obrębie zakładki środowiska bieżącej (rozdz. 5.1) i działa w
poprzek wszystkich trzech zbiorów jednocześnie, analogicznie do filtra środowiska.

### 5.3 Zakres dat i sortowanie

| Przycisk | Klasa | Dosłowne brzmienie |
|---|---|---|
| Zakres dat | `.dn-btn.dn-btn--duch.dn-btn--sm` | „Zakres dat” |
| Sortowanie | `.dn-btn.dn-btn--duch.dn-btn--sm` | „Sortowanie” |

Oba przyciski są obecne w zrzucie statycznym bez rozwiniętego stanu — prototyp nie niesie zrzutu
otwartego panelu żadnego z nich. **[DO DECYZJI OPERATORA]** — dokładna postać kontrolki „Zakres
dat” (para pól daty, lista gotowych zakresów typu „ostatnie 7 dni”, czy kalendarz) i kontrolki
„Sortowanie” (lista pól sortowania i kierunku) nie ma zrzutu potwierdzającego. Jedyne, co
zweryfikowane, to obecność obu przycisków, ich klasa i dosłowne brzmienie.

Położenie obu przycisków na paska filtrów — po prawej stronie, za rozpychaczem `.hs-rozpychacz`
(rozdz. 3.4) — je odróżnia od pola wyszukiwania i zakładek zakresu środowisk, umieszczonych po lewej.
Ten podział przestrzenny sugeruje odmienną rolę: wyszukiwanie i zakres środowisk są filtrami
zawężającymi treść trzech zbiorów jednocześnie (rozdz. 5.1–5.2), podczas gdy „Zakres dat” i
„Sortowanie” prawdopodobnie działają na porządek prezentacji w obrębie zbioru już zawężonego, a nie
na to, które sesje w ogóle wchodzą do wykazu — rozróżnienie analogiczne do różnicy między filtrem a
sortowaniem w każdym innym rejestrze platformy, choć bez zrzutu jednoznacznie to potwierdzającego dla
tego konkretnego okna.

---

## 6. Operacje na sesji

W odróżnieniu od modelu „zaznacz wiele wierszy, wykonaj operację zbiorczą z paska na dole” opisanego
w poprzedniej redakcji tego dokumentu, zweryfikowany prototyp niesie wyłącznie **operacje
jednowierszowe**: każdy wiersz sesji ma dokładnie dwa przyciski akcji (`.hs-akcje`), różne zależnie
od zbioru, do którego wiersz należy. Nie istnieje w zrzucie żaden mechanizm zaznaczenia wielu
wierszy naraz, żaden licznik zaznaczenia, żaden pasek operacji zbiorczych.

### 6.1 Wróć / Zakończ — sesje czynne

| Przycisk | Komenda kandydująca | Pola żądania | Uzasadnienie mapowania |
|---|---|---|---|
| „Wróć” | `session.bind`, następnie `session.focus` | `sessionId`, `clientId` (`session.bind`); `sessionId`, `clientId` (`session.focus`) | sesja `active` trwa na rdzeniu — powrót wymaga związania połączenia klienta z sesją trwającą, potem ogniskowania karty, dokładnie jak w przebiegu 13.1 |
| „Zakończ” | `session.close` | `sessionId` | „Zakończ” jest czasownikiem dokonanym odpowiadającym `session.close` („zamyka sesję”); alternatywą jest `session.stop` („zatrzymuje tury”), lecz ta komenda nie kończy sesji, wyłącznie przerywa tury biegnące — rozdz. 6.4 |

### 6.2 Otwórz / Archiwizuj — sesje zakończone

| Przycisk | Komenda kandydująca | Pola żądania | Uzasadnienie mapowania |
|---|---|---|---|
| „Otwórz” | `session.resume` (preferowana) albo `session.open` | `sessionId` | obie komendy pasują semantycznie; `session.resume` „wznawia zamkniętą sesję wraz z jej oknami” — dosłownie opisuje sesję `finished`; `session.open` „otwiera sesję wraz z oknami” bez założenia o stanie poprzednim |
| „Archiwizuj” | `session.archive` | `sessionIds: string[]` | dosłowne dopasowanie nazwy |

### 6.3 Przywróć / Pobierz — sesje archiwalne

| Przycisk | Komenda kandydująca | Pola żądania | Uzasadnienie mapowania |
|---|---|---|---|
| „Przywróć” | `session.restore` | `sessionIds: string[]` | dosłowne dopasowanie nazwy; lid sekcji (rozdz. 4) potwierdza: „sesja archiwalna wymaga wcześniejszego przywrócenia” |
| „Pobierz” | `history.load` (per `windowId` sesji), złożone w plik | `windowId` (wym), `limit` (opc), `before:int64` (opc) | „Pobierz” różni się nazwą od „Eksportuj…” opisywanej w poprzedniej redakcji, lecz mechanizm źródłowy jest ten sam — sesja archiwalna nie otwiera się bezpośrednio (wymaga przywrócenia najpierw), więc jedyną drogą wglądu bez przywrócenia jest pobranie zapisu |

Historia jest **per okno komunikacji**, nie per sesja: `HistoryEntry` niesie pola `windowId` i
`sessionId`. Sesja ma wiele okien (`Session.windowIds`), więc „Pobierz” odczytuje `history.load` dla
każdego okna sesji i składa wynik w jeden zapis (rozdz. 13.4). Pobranie nie zmienia stanu sesji.

### 6.4 Komendy obszaru `session` bez odpowiednika w widoku

Obszar `session` liczy dziewiętnaście komend (rozdz. 10.1). Sześć przycisków opisanych w rozdz.
6.1–6.3 wyczerpuje to, co zweryfikowany prototyp eksponuje. Pozostałe komendy obszaru **nie mają
kontrolki w tym oknie**, choć są semantycznie bliskie rejestrowi sesji:

| Komenda | Przeznaczenie | Status w tym oknie |
|---|---|---|
| `session.delete` | trwałe usunięcie sesji wraz z zapisem | **nie eksponowana** — brak przycisku „Usuń” w żadnym z trzech zbiorów zweryfikowanego prototypu; **[DO DECYZJI OPERATORA]**, czy i w którym zbiorze (najpewniej „Sesje archiwalne”, obok „Przywróć”/„Pobierz”) taka kontrolka ma powstać |
| `session.stop` | zatrzymanie tur biegnących w oknach sesji, bez zamykania jej | **nie eksponowana** osobno — rozdz. 6.1 rozważa ją jako alternatywę dla „Zakończ”, lecz zatrzymanie tur bez zamknięcia sesji nie ma własnego przycisku |
| `session.rename` | zmiana nazwy sesji | **nie eksponowana** w tym oknie |
| `session.copy` | kopia sesji wraz z zapisem | **nie eksponowana** w tym oknie |
| `session.project.set` / `session.project.clear` | przypisanie/wyjęcie z projektu | **nie eksponowane** w tym oknie; przypisanie odbywa się przy zakładaniu projektu ([Okno nowego projektu](okno-nowego-projektu.md) rozdz. 4.1) albo z poziomu Workspace, poza zakresem tego dokumentu |
| `session.tool.attach` / `.detach` / `.list` | narzędzia dołożone do sesji | **nie eksponowane** jako kontrolki edycji w tym oknie; `session.tool.list` pozostaje właściwe do podglądu (rozdz. 11.4), gdyby podgląd sesji miał go pokazywać — brak zrzutu takiego podglądu |
| `session.create` | zakłada sesję | poza zakresem tego okna — zakładanie należy do przedsionka środowiska i do opracowania [Okno nowego projektu](okno-nowego-projektu.md) |

Ten podział — sześć komend eksponowanych, osiem nieeksponowanych, jedna (`session.create`) spoza
zakresu z definicji — jest rozstrzygnięciem zweryfikowanym, nie domysłem: każdy wiersz tabeli
odzwierciedla nieobecność potwierdzoną brakiem odpowiedniego przycisku w znaczniku HTML prototypu,
sprawdzonym wprost (rozdz. 4 wyżej, `grep` po nazwach przycisków w źródle).

Trzy z ośmiu komend nieeksponowanych (`session.rename`, `session.copy`, `session.project.set`/
`.clear`) są komendami zmiany tożsamości albo przynależności sesji, nie zmiany jej stanu w
wyliczeniu `SessionStatus` (rozdz. 1.3) — ich brak w tym oknie jest spójny z charakterem rejestru
jako miejsca **przeglądu i przejść stanu**, nie edycji metadanych sesji. Dwie kolejne
(`session.tool.attach`/`.detach`) dotyczą zestawu narzędzi dołożonych do sesji trwającej — również
poza naturalnym zakresem rejestru historycznego. Jedyna komenda nieeksponowana, której brak
odbiega od tej reguły spójności, jest `session.delete` (rozdz. 6.4 wyżej, scenariusz E w rozdz. 14)
— usunięcie trwałe jest operacją zmiany stanu (a właściwie zniesienia bytu) analogiczną do
`session.archive`, która swoją kontrolkę ma, więc brak analogicznej kontrolki dla usunięcia jest
rozbieżnością wartą uwagi Operatora przy dalszej budowie tego okna, nie jedynie faktem
odnotowanym na równi z pozostałymi siedmioma.

---

## 7. Stany kontrolek

Wykaz w postaci ośmiu stanów wymaganych przez rozdz. 4.2 standardu redakcyjnego zbioru: spoczynek,
wskazanie kursorem, wciśnięcie, ognisko, nieaktywny, ładowanie, pusty, błąd.

| Kontrolka | Spoczynek | Wskazanie kursorem | Wciśnięcie | Ognisko | Nieaktywny | Ładowanie | Pusty | Błąd |
|---|---|---|---|---|---|---|---|---|
| Wiersz sesji | tło `--dn-powierzchnia` | obrys `--dn-obrys`, tło `--dn-powierzchnia-2` | — | pierścień `--dn-fokus` | brak (zasada zero blokad) | szkielet wiersza | zbiór pusty pomija swój nagłówek (rozdz. 4) — **[DO DECYZJI OPERATORA]**, brak zrzutu stanu pustego jednego zbioru | — |
| Zakładka zakresu środowisk | `.dn-zakladka` | tekst `--dn-tekst` | — | pierścień `--dn-fokus` | brak | — | — | — |
| Zakładka wybrana | `aria-selected="true"` | rozjaśnienie | — | pierścień `--dn-fokus` | brak | — | — | — |
| Pole wyszukiwania | obrys `--dn-obrys-mocny` | obrys `--dn-tekst-3` | — | obrys `--dn-fokus` | brak | — | tekst podpowiedzi „Szukaj w nazwach, projektach i modułach…” | — |
| „Wróć” | `.dn-btn--duch.dn-btn--sm` | tło `--dn-hover` | wciśnięcie | pierścień `--dn-fokus` | brak (zasada zero blokad) | `aria-busy` z wirnikiem `--dn-wym-spinner` — **[DO DECYZJI OPERATORA]**, brak zrzutu | nie dotyczy | dymek błędu (rozdz. 12) |
| „Zakończ” | `.dn-btn--duch.dn-btn--sm` | tło `--dn-hover` | wciśnięcie | pierścień `--dn-fokus` | brak | — | nie dotyczy | dymek błędu |
| „Otwórz” | `.dn-btn--duch.dn-btn--sm` | tło `--dn-hover` | wciśnięcie | pierścień `--dn-fokus` | brak | jak „Wróć” | nie dotyczy | dymek błędu |
| „Archiwizuj” | `.dn-btn--duch.dn-btn--sm` | tło `--dn-hover` | wciśnięcie | pierścień `--dn-fokus` | brak | — | nie dotyczy | dymek błędu |
| „Przywróć” | `.dn-btn--duch.dn-btn--sm` | tło `--dn-hover` | wciśnięcie | pierścień `--dn-fokus` | brak | — | nie dotyczy | dymek błędu |
| „Pobierz” | `.dn-btn--duch.dn-btn--sm` | tło `--dn-hover` | wciśnięcie | pierścień `--dn-fokus` | brak | pasek postępu pobierania — **[DO DECYZJI OPERATORA]** | nie dotyczy | dymek błędu |
| Cały rejestr | — | — | — | — | — | szkielet trzech zbiorów | „Brak sesji w wybranym zakresie” — **[DO DECYZJI OPERATORA]**, dosłowna treść nieobecna w zrzucie, przyjęta analogicznie do innych okien zbioru | „Nie udało się wczytać rejestru — spróbuj ponownie” — jak wyżej |

Wszystkie sześć przycisków akcji (rozdz. 6.1–6.3) niesie tę samą klasę wizualną
`.dn-btn.dn-btn--duch.dn-btn--sm` niezależnie od tego, czy działanie jest odwracalne
(„Wróć”, „Otwórz”, „Archiwizuj”, „Przywróć”) czy potencjalnie kończące pracę w danym miejscu
(„Zakończ”, „Pobierz”) — żaden z sześciu nie niesie wariantu `--sygnal` ani `--niebezpieczny`
widocznego w poprzedniej redakcji tego dokumentu dla „Wznów” i „Usuń”. Zgodnie z zasadą zero blokad
(`komponenty.css`, w. 7) żaden przycisk nie jest wyłączany technicznie — niegotowość (np. sesja już
w toku innej operacji) komunikuje się stanem `aria-busy` albo komunikatem, nie odjęciem klikalności.

---

## 8. Skróty klawiszowe

| Skrót | Działanie | Zasięg |
|---|---|---|
| `Ctrl+F` | ognisko w polu wyszukiwania | okno |
| `↑` / `↓` | przejście między wierszami w obrębie zbioru bieżącego | wykaz |
| `Enter` | aktywacja pierwszego przycisku wiersza z ogniskiem — „Wróć” dla sesji czynnej, „Otwórz” dla zakończonej, „Przywróć” dla archiwalnej | wykaz |
| `Esc` | zamknięcie nakładki | okno |
| `Tab` | przejście między strefami: sekcje ustawień → pasek filtrów → wiersze rejestru → przyciski wiersza | okno |

Skrót `Delete` obecny w poprzedniej redakcji tego dokumentu jest usunięty z niniejszego wykazu —
wiązał się z przyciskiem „Usuń”, którego zweryfikowany prototyp nie zawiera (rozdz. 6.4). Jeżeli
kontrolka usunięcia trwałego powstanie w przyszłej redakcji tego okna, skrót `Delete` jest
kandydatem naturalnym do przywrócenia w tym miejscu.

---

## 9. Dostępność

| Wymóg | Realizacja |
|---|---|
| Nawigacja klawiaturą | pełna; każdy element osiągalny `Tab` (rozdz. 8) |
| Ognisko widoczne | pierścień `--dn-fokus` grubości `--dn-wym-fokus` |
| Stan czynny | znak `.pt-tetno` towarzyszy nazwie sesji w kolumnie „Sesja”, nie jest jedynym nośnikiem — nazwa i podtytuł tekstowe niosą znaczenie niezależnie od animacji |
| Ruch | tętnienie znaku (`.pt-tetno`) wyciszane przy `prefers-reduced-motion: reduce`, zgodnie z regułą zapisaną w `prototyp.css` przy tej klasie |
| Trzy zbiory jako regiony nazwane | każda z trzech sekcji `.hs-zbior` niesie `aria-label` własny („Sesje czynne”, „Sesje zakończone”, „Sesje archiwalne”) — czytnik ekranu ogłasza przynależność wiersza do zbioru bez polegania wyłącznie na położeniu wizualnym |
| Nagłówek kolumn | `aria-hidden="true"` (rozdz. 4.6) — informacja o znaczeniu kolumny jest niesiona inaczej: `.hs-meta` w każdym wierszu stoi w tej samej kolejności, a etykieta `aria-label` pola wyszukiwania i zakładek środowisk wystarcza do orientacji bez odczytu nagłówka wizualnego |
| Zakładki zakresu środowisk | `role="tablist"`/`role="tab"`, `aria-selected` — wzorzec standardowy nawigacji kartami |
| Kontrast | znaki stanu i tekst przygaszony (`--dn-tekst-2` dla nazw archiwalnych, rozdz. 4.3) spełniają progi `budowa/shared/kontrasty-progi.json` |
| Komunikaty | stan pusty i błąd czytane przez `role="status"` |
| Dwa przyciski akcji per wiersz rozróżnialne przez tekst | „Wróć”/„Zakończ”, „Otwórz”/„Archiwizuj”, „Przywróć”/„Pobierz” (rozdz. 6) różnią się treścią tekstową, nie tylko położeniem — technologia wspomagająca czytająca sam przycisk bez kontekstu wiersza nadal rozróżnia działanie |
| Lewa nawigacja pozioma poniżej 900 px | `.hs-sekcje` zachowuje `role` i strukturę listy przy przejściu na układ poziomy (rozdz. 3.5) — zmienia się wyłącznie kierunek wizualny (`flex-direction: row`), nie semantyka |
| Podtytuł jako informacja dodatkowa, nie zastępcza | `.hs-nazwa small` (rozdz. 4.5) uzupełnia tytuł sesji (`<b>`), nie zastępuje go — czytnik ekranu odczytujący wyłącznie pierwszy poziom nadal otrzymuje nazwę sesji sensowną samodzielnie |

---

## 10. Komendy i struktury kontraktu

Nazwy komend, pól, struktur i wyliczeń — wprost z `budowa/shared/contract.json`. Deweloper nie
wprowadza nazw poza tym wykazem. Kolumna „Eksponowana w oknie” rozstrzyga, które z komend mają
kontrolkę zweryfikowaną w prototypie (rozdz. 6) — wykaz pełny nie jest tożsamy z powierzchnią
widoczną, co rozdz. 6.4 wykazuje wprost.

### 10.1 Pełny wykaz komend obszaru `session`

Obszar `session` liczy 19 komend.

| Komenda | Przeznaczenie | Pola żądania | Pola wyniku | Eksponowana w oknie |
|---|---|---|---|---|
| `session.create` | zakłada sesję | `title` (opc), `projectId` (opc), `metadata:json` (opc) | `session:Session` | nie — poza zakresem (rozdz. 6.4) |
| `session.list` | zwraca listę sesji | `status:SessionStatus` (opc), `limit` (opc), `offset` (opc), `includePresence:bool` (opc) | `sessions:Session[]`, `total:int`, `presence:SessionPresence[]` (opc) | tak — zasila trzy zbiory (rozdz. 4) |
| `session.open` | otwiera sesję wraz z oknami | `sessionId` | `session`, `windows:Window[]` | częściowo — kandydat dla „Otwórz” (rozdz. 6.2) |
| `session.bind` | wiąże połączenie z sesją trwającą na rdzeniu | `sessionId`, `clientId`, `windowIds` (opc) | `session`, `windows`, `bound:bool`, `resumed:bool` | tak — „Wróć” (rozdz. 6.1) |
| `session.focus` | przenosi ognisko na kartę sesji | `sessionId`, `clientId`, `windowId` (opc), `targetClientId` (opc) | `sessionId`, `windowId` (opc), `previousSessionId` (opc), `focusedAt:int64` | tak — „Wróć”, drugi krok (rozdz. 6.1, 13.1) |
| `session.resume` | wznawia zamkniętą sesję wraz z oknami | `sessionId` | `session`, `windows` | tak — „Otwórz” (rozdz. 6.2) |
| `session.close` | zamyka sesję | `sessionId` | `session` | tak — „Zakończ” (rozdz. 6.1) |
| `session.stop` | zatrzymuje tury biegnące w oknach sesji | `sessionId` | `sessionId`, `stoppedWindowIds` | nie (rozdz. 6.4) |
| `session.delete` | trwale usuwa sesje wraz z zapisem | `sessionIds:string[]`, `confirm:bool` | `deletedIds`, `deletedCount:int`, `missingIds` (opc) | nie (rozdz. 6.4) |
| `session.rename` | zmienia nazwę sesji | `sessionId`, `title` | `session` | nie (rozdz. 6.4) |
| `session.copy` | kopiuje sesję wraz z zapisem | `sessionId`, `title` (opc) | `session`, `copiedMessages:int` | nie (rozdz. 6.4) |
| `session.project.set` | przenosi sesje do projektu | `sessionIds`, `projectId` (opc), `projectName` (opc) | `projectId`, `movedIds` | nie (rozdz. 6.4) |
| `session.project.clear` | wyjmuje sesje z projektu | `sessionIds` | `clearedIds` | nie (rozdz. 6.4) |
| `session.archive` | przenosi sesje do archiwum | `sessionIds` | `archivedIds` | tak — „Archiwizuj” (rozdz. 6.2) |
| `session.restore` | przywraca sesje z archiwum | `sessionIds` | `restoredIds` | tak — „Przywróć” (rozdz. 6.3) |
| `session.archive.list` | zwraca sesje archiwum (wgląd z okna ustawień) | `offset` (opc), `limit` (opc) | `sessions`, `total:int` | tak — zasila zbiór „Sesje archiwalne” (rozdz. 4.3) |
| `session.tool.attach` | dokłada narzędzie/skill do sesji | `sessionId`, `toolName`, `source:SessionToolSource` (opc) | `tool`, `tools`, `alreadyAttached:bool` | nie (rozdz. 6.4) |
| `session.tool.detach` | zdejmuje dołożenie z sesji | `sessionId`, `toolName` | `detached:bool`, `tools` | nie (rozdz. 6.4) |
| `session.tool.list` | zwraca narzędzia dołożone do sesji | `sessionId` | `tools:SessionTool[]`, `total:int` | nie (rozdz. 6.4, 11.4) |

Sześć komend eksponowanych (`session.list`, `session.bind`, `session.focus`, `session.resume`,
`session.close`, `session.archive`, `session.restore`, `session.archive.list` — osiem licząc
odczyty zasilające widok) i jedenaście nieeksponowanych sumują się do dziewiętnastu — komplet
obszaru.

### 10.2 Komendy obszaru `history`

| Komenda | Przeznaczenie | Pola żądania | Pola wyniku | Eksponowana w oknie |
|---|---|---|---|---|
| `history.load` | zwraca historię rozmowy okna od pozycji najnowszej | `windowId` (wym), `limit` (opc), `before:int64` (opc) | `windowId`, `entries:HistoryEntry[]`, `total:int` | tak — „Pobierz” (rozdz. 6.3, 13.4) |
| `history.delete` | usuwa pozycje historii okna; bez wskazania czyści całą historię okna | `windowId` (wym), `entryIds:string[]` (opc) | `windowId`, `deleted:int` | nie — brak kontrolki w tym oknie, analogicznie do `session.delete` |

### 10.3 Struktury i wyliczenia

**`Session`** — sesja jako byt trwały rdzenia:

| Pole | Typ |
|---|---|
| `id` | `string` |
| `title` | `string` |
| `projectId` | `string` |
| `status` | `SessionStatus` |
| `windowIds` | `string[]` |
| `createdAt` / `updatedAt` | `int64` |
| `metadata` | `json` |

**`SessionPresence`** — żywy stan sesji trwającej na rdzeniu:

| Pole | Typ |
|---|---|
| `sessionId` | `string` |
| `live` | `bool` |
| `environmentCode` / `moduleCode` | `string` |
| `focusedWindowId` | `string` |
| `openWindowCount` / `streamingWindowCount` | `int` |
| `accessGrantIds` | `string[]` |
| `loop` | `LoopState` |
| `lastActivityAt` | `int64` |

**`HistoryEntry`** — pozycja historii rozmowy okna:

| Pole | Typ |
|---|---|
| `id` | `string` |
| `windowId` / `sessionId` | `string` |
| `role` | `string` |
| `preview` | `string` |
| `createdAt` | `int64` |

**`SessionStatus`** — wyliczenie stanu sesji: `active`, `paused`, `finished`, `archived` (rozdz. 1.3).

**`ErrorInfo`** — treść błędu w odpowiedzi rdzenia; struktura wspólna całemu kontraktowi,
przywoływana w rozdz. 12 przy rozróżnieniu kodów ponawialnych i nieponawialnych:

| Pole | Typ |
|---|---|
| `code` | `ErrorCode` |
| `message` | `string` |
| `details` | `json` |
| `retryable` | `bool` |

**Trzy struktury w postaci wymuszonej**, zgodnie z rozdz. 4.2 wymogu normatywnej szczegółowości
standardu redakcyjnego zbioru — pole po polu, z wymagalnością, tak jak zwraca je `session.list` z
`includePresence=true` dla wiersza „Raport końcowy — redakcja” (rozdz. 4.1):

```json
{
  "struktura": "Session",
  "przyklad": "wiersz «Raport końcowy — redakcja», zbiór Sesje czynne",
  "pola": {
    "id":         "string — wymagane",
    "title":      "string — opcjonalne, tu: „Raport końcowy — redakcja”",
    "projectId":  "string — opcjonalne, tu: identyfikator projektu „Sprawozdawczość Q3”",
    "status":     "SessionStatus — wymagane, tu: active",
    "windowIds":  "string[] — opcjonalne, tu: co najmniej jedno okno modułu Studio",
    "createdAt":  "int64 — wymagane, milisekundy epoki",
    "updatedAt":  "int64 — wymagane, tu odpowiada wartości prezentowanej jako „12:04”",
    "metadata":   "json — opcjonalne"
  }
}
```

```json
{
  "struktura": "SessionPresence",
  "przyklad": "ten sam wiersz — obecność żywa towarzysząca sesji active",
  "pola": {
    "sessionId":            "string — wymagane",
    "live":                 "bool — wymagane, tu: true (znak .pt-tetno, rozdz. 3.4)",
    "environmentCode":      "string — opcjonalne, tu: workspace",
    "moduleCode":           "string — opcjonalne, tu: studio",
    "focusedWindowId":      "string — opcjonalne",
    "openWindowCount":      "int — wymagane",
    "streamingWindowCount": "int — wymagane",
    "accessGrantIds":       "string[] — opcjonalne",
    "loop":                 "LoopState — opcjonalne, puste dla okna samodzielne",
    "lastActivityAt":       "int64 — wymagane"
  }
}
```

```json
{
  "struktura": "HistoryEntry",
  "przyklad": "jedna pozycja zwrócona przez history.load przy «Pobierz» (rozdz. 6.3, 13.4)",
  "pola": {
    "id":        "string — wymagane",
    "windowId":  "string — wymagane",
    "sessionId": "string — opcjonalne",
    "role":      "string — wymagane, tu: operator albo model",
    "preview":   "string — opcjonalne, skrót treści na potrzeby wykazu",
    "createdAt": "int64 — wymagane"
  }
}
```

Trzy struktury powyżej różnią się liczbą pól wymaganych — `Session` niesie dwa pola opcjonalne
kluczowe dla prezentacji (`title`, `projectId`) obok pięciu wymaganych; `HistoryEntry` niesie
`sessionId` jako opcjonalne mimo że okno historii sesji odczytuje ją zawsze w kontekście znanej
sesji — pole istnieje osobno, bo `history.load` przyjmuje wyłącznie `windowId`, nie `sessionId`
(rozdz. 11.3), więc pozycja historii niesie własny zapis przynależności do sesji niezależnie od
parametru wywołania, które ją zwróciło.

### 10.4 Zdarzenia zasilające widok

| Zdarzenie | Kiedy | Skutek w oknie |
|---|---|---|
| `session.changed` | zmiana sesji na rdzeniu | aktualizacja wiersza bez pełnego odczytu; przeniesienie wiersza między zbiorami, gdy zmiana `status` tego wymaga |
| `session.focus.changed` | zmiana ogniska karty sesji | odświeżenie znacznika sesji bieżącej |
| `history.changed` | zmiana historii rozmowy okna | odświeżenie kolumny „Zapis” (rozdz. 4.4), gdy pochodzi z liczby wpisów |
| `session.tool.attached` / `session.tool.detached` | dołożenie/zdjęcie narzędzia sesji | bez skutku widocznego w tym oknie — rejestr nie pokazuje zestawu narzędzi (rozdz. 6.4) |

---

## 11. Okno komunikacji jako nośnik historii

Historia rozmowy jest w Danaco Console przypisana do **okna komunikacji** (struktura `Window`), nie
wprost do sesji. Sesja skupia okna (`Session.windowIds`), a każde okno niesie własny przebieg
rozmowy. Dlatego komendy obszaru `history` operują na `windowId`, a nie `sessionId`. Zrozumienie
struktury `Window` jest warunkiem poprawnego działania przycisku „Pobierz” (rozdz. 6.3) z poziomu
okna historii sesji.

### 11.1 Struktura Window

Struktura `Window` z kontraktu — komplet pól:

| Pole | Typ | Rola |
|---|---|---|
| `id` | `string` | identyfikator okna |
| `sessionId` | `string` | sesja, do której okno należy |
| `moduleId` | `string` | moduł, w którym okno pracuje |
| `modelChannelId` | `string` | kanał modelu przypisany do okna |
| `agentId` | `string` | agent obsługujący okno (jeśli przypisany) |
| `workingDirs` | `string[]` | katalogi robocze okna |
| `executionEnv` | `ExecutionEnv` | zasięg wykonania modelu |
| `permissionMode` | `PermissionMode` | tryb uprawnień okna |
| `windowRole` | `WindowRole` | rola okna w pętli koordynator–wykonawca |
| `coordinatorWindowId` | `string` | okno koordynatora (dla okna wykonawcy) |
| `title` | `string` | tytuł okna |
| `status` | `WindowStatus` | stan otwarcia okna |
| `createdAt` / `updatedAt` | `int64` | znaczniki czasu |

Okno historii sesji korzysta z pól `id` (do `history.load`), `sessionId` (do wiązania okien
z sesją) i `status` (do rozróżnienia okien otwartych i zamkniętych przy składaniu pobrania).

### 11.2 Wyliczenia okna

**`WindowStatus`** — stan okna komunikacji:

| Wartość | Znaczenie |
|---|---|
| `open` | okno otwarte |
| `closed` | okno zamknięte |

**`WindowRole`** — rola okna w pętli koordynator–wykonawca:

| Wartość | Znaczenie |
|---|---|
| `executor` | okno wykonawcy; zakończenie tury wybudza koordynatora |
| `coordinator` | okno koordynatora; przekazuje zlecenie przez narzędzie |
| `standalone` | okno samodzielne, poza pętlą |

**`ExecutionEnv`** — zasięg wykonania modelu:

| Wartość | Znaczenie |
|---|---|
| `local` | urządzenie użytkownika, przez agenta lokalnego |
| `core` | host rdzenia serwera |
| `remote` | host zdalny wskazany w konfiguracji okna |

**`PermissionMode`** — tryb uprawnień okna komunikacji:

| Wartość | Znaczenie |
|---|---|
| `manual` | pytanie o zgodę przed każdą zmianą |
| `acceptEdits` | automatyczna zgoda na zmiany plików |
| `plan` | praca planistyczna bez zmian w systemie |
| `auto` | decyzje o uprawnieniach podejmuje model |
| `dontAsk` | bez zapytań, z zachowaniem ograniczeń |
| `bypassPermissions` | pominięcie kontroli uprawnień |

Okno historii nie zmienia tych ustawień — nie pokazuje ich nawet jako kontekst; zweryfikowany
prototyp nie niesie zrzutu podglądu okien sesji z poziomu rejestru. Tryb uprawnień ustala się przy
zakładaniu okna komunikacji komendą `window.create` i zmienia komendą `window.update` — obie
przenoszą pole `permissionMode` i obie opisuje [Rama okna](rama-okna.md) w wykazie komend obszaru
`window`. Zmiana trybu uprawnień należy zatem do okna operacyjnego, nie do rejestru.

### 11.3 Historia per okno a widok per sesja

```
   Sesja  (Session)
   ├── windowId: w-1  ──►  history.load(windowId=w-1)  ──►  HistoryEntry[]
   ├── windowId: w-2  ──►  history.load(windowId=w-2)  ──►  HistoryEntry[]
   └── windowId: w-3  ──►  history.load(windowId=w-3)  ──►  HistoryEntry[]
                                     │
                                     ▼
                „Pobierz” = suma przebiegów wszystkich okien sesji
```

Legenda: „Pobierz” (rozdz. 6.3) odczytuje `history.load` dla każdego `windowId` z
`Session.windowIds` i składa wynik w jeden zapis. `HistoryEntry` niesie oba pola — `windowId` i
`sessionId` — więc pozycje da się przypisać zarówno do okna, jak i do sesji.

**Dlaczego „Pobierz” nie wymaga otwarcia okna.** Sesja archiwalna (rozdz. 4.3) „wymaga wcześniejszego
przywrócenia”, by wrócić do pracy — lid sekcji (rozdz. 4) to potwierdza. „Pobierz” jest jednak
operacją odczytu czystego: `history.load` przyjmuje `windowId`, pole zachowane w strukturze `Window`
niezależnie od tego, czy sesja nadrzędna jest `active`, `finished` czy `archived` (rozdz. 1.3) —
`WindowStatus` (rozdz. 11.2) rozstrzyga wyłącznie, czy okno jest `open` czy `closed`, nie czy sesja
jest zarchiwizowana. Dlatego pobranie zapisu sesji archiwalnej nie wymaga uprzedniego
`session.restore` — identyfikatory okien pozostają czytelne, nawet gdy sama sesja nie jest w stanie
gotowym do wznowienia pracy. To rozróżnienie — między stanem sesji a stanem jej okien — jest
kluczowe dla zrozumienia, czemu „Przywróć” i „Pobierz” (rozdz. 6.3) są dwiema niezależnymi drogami z
tego samego wiersza, nie jedną drogą warunkującą drugą.

### 11.4 Narzędzia dołożone do sesji

Obszar `session` niesie trzy komendy zarządzania narzędziami dołożonymi (`session.tool.list`,
`.attach`, `.detach`) — rozdz. 6.4 stwierdza wprost, że żadna z nich nie ma kontrolki w zweryfikowanym
prototypie tego okna. Wykaz poniżej jest zachowany jako kompletny opis obszaru kontraktowego, nie
jako opis czegoś widocznego w tym oknie.

| Komenda | Przeznaczenie | Pola żądania | Pola wyniku |
|---|---|---|---|
| `session.tool.list` | zwraca narzędzia dołożone do sesji | `sessionId` | `tools: SessionTool[]`, `total: int` |
| `session.tool.attach` | dokłada narzędzie albo skill do sesji | `sessionId`, `toolName`, `source: SessionToolSource` (opc) | `tool`, `tools`, `alreadyAttached: bool` |
| `session.tool.detach` | zdejmuje dołożenie; zestaw wraca do podstawy z definicji eksperta | `sessionId`, `toolName` | `detached: bool`, `tools` |

**`SessionToolSource`** — czyja ręka dołożyła narzędzie:

| Wartość | Znaczenie |
|---|---|
| `slashCommand` | Operator wpisał komendę po ukośniku z klawiatury |
| `assistant` | asystent działający za Operatora; Operator widzi dołożenie zdarzeniem `session.tool.attached` |

---

## 12. Kody błędów i zachowanie okna

Każda komenda obszarów `session` i `history` może zakończyć się jednym z dziewięciu kodów błędów
kontraktu. Okno reaguje na kod, nie na komunikat — komunikat jest opisem dla Operatora, kod
rozstrzyga o zachowaniu.

Dziewięć kodów jest wspólnych całemu kontraktowi — żaden nie jest właściwy wyłącznie obszarowi
`session` albo `history`. Pole `kodyBledow` w `budowa/shared/contract.json` zawiera dziś osiem
pierwszych pozycji; dopisanie dziewiątego kodu `command_not_understood` jest osobnym zadaniem
w kodzie. To, co odróżnia zachowanie
tego okna od zachowania innego okna platformy przy tym samym kodzie, nie jest samym kodem, lecz
miejscem, w którym komunikat się pojawia (rozdz. 7: dymek błędu przy przycisku wiersza, nie przy
formularzu ani przy paśmie akcji zbiorczej — bo taka nie istnieje, rozdz. 6) i tym, co z sesją dzieje
się po stronie widoku (przeniesienie między zbiorami przy sukcesie operacji zmieniającej stan, brak
przeniesienia przy błędzie).

| Kod | Ponawialny | Znaczenie | Zachowanie okna |
|---|---|---|---|
| `validation_failed` | nie | treść żądania niezgodna z kontraktem | komunikat błędu przy przycisku wiersza; operacja niewykonana |
| `not_found` | nie | wskazany byt nie istnieje | odświeżenie wykazu (`session.list`); sesja mogła zniknąć spod wiersza |
| `not_authenticated` | nie | brak uwierzytelnienia | przejście do okna logowania |
| `permission_denied` | nie | brak uprawnienia do czynności | komunikat „Brak uprawnienia do tej operacji” |
| `conflict` | nie | stan bytu wyklucza czynność | komunikat wyjaśniający stan — przy próbie „Zakończ” na sesji już zamykanej gdzie indziej brzmi „Sesja jest właśnie zamykana na innym urządzeniu” |
| `channel_unavailable` | tak | kanał modelu niedostępny | ponowienie dozwolone; dotyczy tylko bieżącego wywołania |
| `rate_limited` | tak | ograniczenie tempa | ponowienie po odczekaniu |
| `internal_error` | tak | błąd wewnętrzny rdzenia | ponowienie dozwolone; przy powtórce — zgłoszenie diagnostyki |

Tabela wymienia osiem kodów, ponieważ dziewiąty kod kontraktu — `command_not_understood` —
powstaje wyłącznie przy niedopasowaniu treści polecenia języka naturalnego do funkcji
(`layer.function.invoke`, [Kontrakty komunikacji](../architektura/kontrakty-komunikacji.md)
rozdz. 7.2). Okno historii sesji nie przyjmuje poleceń języka naturalnego, więc ten kod nie
występuje w żadnym z jego wywołań. Wykaz pełny dziewięciu kodów niesie
[Architektura techniczna](../architektura/architektura.md) rozdz. 19.

Mapowanie sześciu operacji eksponowanych (rozdz. 6.1–6.3) na najczęstsze kody:

| Operacja | Typowy kod przy niepowodzeniu |
|---|---|
| „Wróć” (`session.bind`/`focus`) | `not_found` (sesja zniknęła), `channel_unavailable` (kanał modelu) |
| „Zakończ” (`session.close`) | `conflict` (sesja już zamykana), `not_found` |
| „Otwórz” (`session.resume`) | `not_found`, `channel_unavailable` |
| „Archiwizuj” (`session.archive`) | `not_found`, `conflict` |
| „Przywróć” (`session.restore`) | `not_found`, `conflict` |
| „Pobierz” (`history.load`) | `not_found` (okno usunięte) |

Kody ponawialne (`channel_unavailable`, `rate_limited`, `internal_error`) niosą pole
`retryable=true` w strukturze `ErrorInfo` — okno pokazuje przy nich przycisk ponowienia obok
komunikatu.

---

## 13. Przebiegi operacji

Każdy przebieg jako diagram sekwencji między Operatorem, powłoką (klientem) a rdzeniem, nazwany
dosłownym brzmieniem przycisku zweryfikowanego w prototypie (rozdz. 6). Nazwy komend i pól — z
kontraktu.

### 13.1 Wróć — sesja czynna

```
Operator            Powłoka                         Rdzeń
   │  „Wróć”            │                              │
   ├───────────────────►│  session.bind                │
   │                    │  {sessionId, clientId}       │
   │                    ├─────────────────────────────►│
   │                    │  {session, windows,          │
   │                    │   bound, resumed}            │
   │                    │◄─────────────────────────────┤
   │                    │  session.focus               │
   │                    │  {sessionId, clientId}       │
   │                    ├─────────────────────────────►│
   │                    │  {focusedAt, windowId}       │
   │                    │◄─────────────────────────────┤
   │  obszar roboczy    │                              │
   │◄───────────────────┤  (nakładka zamknięta)        │
```

Legenda: dla sesji `active` najpierw `session.bind` wiąże bieżące połączenie klienta z trwającą
sesją i odtwarza jej okna, następnie `session.focus` ogniskuje kartę. Okno historii zamyka się po
ogniskowaniu.

### 13.2 Otwórz — sesja zakończona

```
Operator            Powłoka                         Rdzeń
   │  „Otwórz”          │                              │
   ├───────────────────►│  session.resume              │
   │                    │  {sessionId}                 │
   │                    ├─────────────────────────────►│
   │                    │  {session, windows}          │
   │                    │◄─────────────────────────────┤
   │  obszar roboczy    │                              │
   │◄───────────────────┤                              │
```

Legenda: dla sesji `finished` komenda `session.resume` wznawia zamkniętą sesję wraz z jej oknami —
nie potrzeba `session.bind`, bo sesja nie trwała na rdzeniu.

### 13.3 Archiwizuj i Przywróć

```
Operator            Powłoka                         Rdzeń
   │  „Archiwizuj”      │  session.archive             │
   ├───────────────────►├─────────────────────────────►│
   │                    │  {archivedIds}               │
   │                    │◄─────────────────────────────┤
   │  (sesja przenosi się do zbioru „Sesje archiwalne”)│
   │                    │                              │
   │  „Przywróć”        │  session.restore             │
   ├───────────────────►├─────────────────────────────►│
   │                    │  {restoredIds}               │
   │                    │◄─────────────────────────────┤
   │  (sesja wraca do zbioru „Sesje zakończone”)       │
```

Legenda: archiwizacja zachowuje zapis w całości i przenosi sesję ze zbioru „Sesje zakończone” do
„Sesje archiwalne” (rozdz. 4); przywrócenie wykonuje przejście odwrotne. Obie operacje działają na
rejestrze bieżącym — `session.archive.list` (rozdz. 10.1) jest wyłącznie odczytem zasilającym zbiór
„Sesje archiwalne”, nie osobnym oknem.

### 13.4 Pobierz — sesja archiwalna

```
Operator            Powłoka                         Rdzeń
   │  „Pobierz”         │  dla każdego windowId sesji: │
   ├───────────────────►│  history.load                │
   │                    │  {windowId, limit, before}   │
   │                    ├─────────────────────────────►│
   │                    │  {entries:HistoryEntry[],    │
   │                    │   total}                     │
   │                    │◄─────────────────────────────┤
   │  wybór pliku       │  (złożenie przebiegu)        │
   ├───────────────────►│                              │
   │  plik zapisany     │                              │
   │◄───────────────────┤                              │
```

Legenda: „Pobierz” odczytuje `history.load` dla każdego okna sesji (pole `before` pozwala
stronicować od pozycji najnowszej wstecz), po czym składa przebieg w jeden plik. Operacja nie zmienia
stanu sesji i nie wymaga uprzedniego przywrócenia — różni się tym od „Otwórz” (rozdz. 13.2), które
dla sesji archiwalnej nie jest dostępne wprost (rozdz. 4.3: „sesja archiwalna wymaga wcześniejszego
przywrócenia”).

### 13.5 Zakończ — sesja czynna

```
Operator            Powłoka                         Rdzeń
   │  „Zakończ”         │  session.close               │
   ├───────────────────►├─────────────────────────────►│
   │                    │  {session}                   │
   │                    │◄─────────────────────────────┤
   │  (sesja przenosi się do zbioru „Sesje zakończone”)│
```

Legenda: „Zakończ” zamyka sesję czynną — po powodzeniu sesja opuszcza zbiór „Sesje czynne” i wchodzi
do zbioru „Sesje zakończone” (rozdz. 4), skąd „Otwórz” (rozdz. 13.2) prowadzi z powrotem do pracy.

---

## 14. Scenariusze Operatora

**Scenariusz A — powrót do wczorajszej pracy.** Operator otwiera historię z menu aplikacji. W
zbiorze „Sesje czynne” nie ma wczorajszej sesji — trwa ona jako `paused`, a rozdz. 1.3 pozostawia
otwartym, do którego zbioru taka sesja trafia. Operator odnajduje ją po tytule w wyszukiwarce
(rozdz. 5.2), niezależnie od tego, w którym zbiorze się znajduje, i klika przycisk pierwszy wiersza
(„Wróć” albo „Otwórz”, zależnie od zbioru). Obszar roboczy wraca do stanu sprzed przerwy.

**Scenariusz B — porządki w rejestrze.** Operator odnajduje w zbiorze „Sesje zakończone” kilka
sesji projektu, którego już nie prowadzi. Dla każdej klika „Archiwizuj” — zapis zostaje w całości
(`session.archive`), a sesja przenosi się do zbioru „Sesje archiwalne” (rozdz. 13.3). Gdyby okazały
się potrzebne, odzyska je stamtąd przyciskiem „Przywróć”.

**Scenariusz C — pobranie zapisu bez przywracania.** Operator potrzebuje treści archiwalnej sesji
wyłącznie do wglądu, bez wznawiania pracy w niej. W zbiorze „Sesje archiwalne” klika „Pobierz”
zamiast „Przywróć” — powłoka składa przebieg z `history.load` (rozdz. 13.4) i zapisuje plik. Stan
sesji pozostaje `archived`, nieporuszony.

**Scenariusz D — zakończenie sesji czynnej z poziomu rejestru.** Operator zauważa w zbiorze „Sesje
czynne” sesję, którą uznaje za skończoną, choć nikt jej formalnie nie zamknął. Klika „Zakończ” —
`session.close` przenosi ją do zbioru „Sesje zakończone” (rozdz. 13.5), bez opuszczania okna
historii sesji.

**Scenariusz E — próba trwałego usunięcia, niedostępna w tym oknie.** Operator szuka sposobu na
trwałe usunięcie sesji zawierającej dane, które mają zniknąć bez śladu. Żaden z trzech zbiorów
zweryfikowanego prototypu nie niesie przycisku „Usuń” (rozdz. 6.4) — komenda `session.delete`
istnieje w kontrakcie, lecz nie ma tu kontrolki. Operator nie ma w tym oknie drogi do trwałego
usunięcia; dokument odnotowuje to jako lukę otwartą, nie jako możliwość ukrytą gdzie indziej w tym
samym oknie.

**Scenariusz F — zawężenie po środowisku przed odnalezieniem sesji.** Operator pracuje jednocześnie
w czterech środowiskach i chce przejrzeć wyłącznie sesje CodeStudio, niezależnie od ich stanu.
Klika zakładkę „CodeStudio” (rozdz. 5.1) — wszystkie trzy zbiory (rozdz. 4.1–4.3) przeliczają się
jednocześnie, każdy pokazując wyłącznie wiersze o `SessionPresence.environmentCode = codestudio`.
Nagłówki grup aktualizują liczniki („1 sesja” zamiast „3 sesje” dla zbioru czynnych, w zestawie
przykładowym — jedynym wierszem CodeStudio jest „Migracja bazy — dokumentacja zmian”). Operator
odnajduje sesję szybciej niż przewijając wszystkie trzy zbiory w pełnej liczbie.

**Scenariusz G — sesja wstrzymana, nieoznaczony zbiór.** Operator zamyka aplikację w trakcie pracy
nad sesją, która przechodzi w stan `paused` (rozdz. 1.3) — proces na rdzeniu przestaje trwać, lecz
zapis jest kompletny i sesja pozostaje odzyskiwalna. Po ponownym otwarciu okna historii sesji
Operator szuka tej sesji. Zgodnie z rozstrzygnięciem otwartym w rozdz. 1.3, nie jest ustalone, w
którym z trzech zbiorów zweryfikowanego widoku taka sesja się znajdzie — dokument nie prowadzi
Operatora scenariusza dalej niż do tego punktu, bo dalszy przebieg (klika w zbiorze „Sesje czynne”
czy w zbiorze „Sesje zakończone”) zależy od rozstrzygnięcia, które ma zapaść przed budową tego
fragmentu okna.

---

## 15. Powiązania z projektami i pamięcią

| Powiązanie | Zachowanie |
|---|---|
| Projekt (`Session.projectId`) | podtytuł wiersza (`.hs-nazwa small`, rozdz. 4) niesie nazwę projektu, gdy sesja do niego należy |
| Przypisanie do projektu | odbywa się przy zakładaniu ([Okno nowego projektu](okno-nowego-projektu.md) rozdz. 4.1, `workspace.agent.assign` i sąsiednie komendy) albo z poziomu Workspace — nie z tego okna (rozdz. 6.4) |
| Wznowienie sesji (rozdz. 13.1, 13.2) | przywraca dostęp do pamięci jej projektu |
| Archiwizacja (rozdz. 13.3) | zachowuje zapis w całości; sesja zmienia zbiór, projekt i jego pamięć pozostają nietknięte |

Nota zamykająca sekcję okna (`.hs-nota`), cytat dosłowny:

> „Zakończenie sesji nie usuwa jej wytworów: dokumenty, wersje i wpisy pamięci pozostają w
> projekcie, do którego sesja należała. Archiwizacja przenosi zapis sesji do archiwum konta — wraca
> stamtąd przyciskiem „Przywróć”, wraz z pełnym kontekstem rozmowy i listą plików.”

Ta nota jest jedynym miejscem prototypu, w którym „Usuń” (rozdz. 6.4) mogłoby się pojawić
naturalnie jako przeciwieństwo trwałości opisanej w pierwszym zdaniu — nota go jednak nie
przywołuje, spójnie z brakiem kontrolki usunięcia w całym oknie.

**Trwałość niezależna od stanu widoku.** Trzy zbiory widoku (rozdz. 4) opisują wyłącznie to, gdzie
Operator znajduje sesję na liście — żaden z trzech nie zmienia przynależności projektowej sesji ani
nie dotyka jej pamięci. Ten sam projekt może w danej chwili mieć sesje rozproszone po wszystkich
trzech zbiorach jednocześnie: jedną czynną, dwie zakończone, jedną archiwalną — wszystkie cztery
dzielące ten sam `Session.projectId` i tę samą pamięć projektu. Rejestr historii sesji nie grupuje
wykazu po projekcie (w odróżnieniu od opracowania [Okno nowego projektu](okno-nowego-projektu.md), gdzie
projekt jest jednostką organizującą sam formularz) — grupuje go po stanie sesji (rozdz. 4) i filtruje
po środowisku (rozdz. 5.1); powiązanie projektowe pozostaje widoczne wyłącznie jako fragment
podtytułu wiersza, nie jako oś grupowania.

---

## 16. Terminologia okna

Rozdział zbiera cztery pary pojęć, które łatwo pomylić przy dalszej budowie tego okna, w miejscu
pierwszego pełnego zestawienia w niniejszym dokumencie, zgodnie z rozdz. 5.3 standardu redakcyjnego
zbioru.

| Para pojęć | Rozróżnienie |
|---|---|
| Zbiór a zakładka | „Sesje czynne / zakończone / archiwalne” (rozdz. 4) są **zbiorami** widocznymi jednocześnie, jeden pod drugim; „Wszystkie / TalkIn / WorkSpace / CodeStudio / MultitaskingAI” (rozdz. 5.1) są **zakładkami** — wybór jednej ukrywa pozostałe. Poprzednia redakcja tego dokumentu myliła te dwa mechanizmy, opisując zbiory jako zakładki stanu |
| Wznowienie a otwarcie a powrót | „Wróć” (rozdz. 6.1) dotyczy sesji trwającej na rdzeniu (`session.bind`); „Otwórz” (rozdz. 6.2) dotyczy sesji zamkniętej (`session.resume`) — oba przywracają dostęp do pracy, lecz różnią się komendą źródłową i stanem wyjściowym sesji |
| Zakończenie a usunięcie | „Zakończ” (`session.close`) zmienia stan sesji na `finished`, zapis pozostaje w całości i sesja jest odzyskiwalna przez „Otwórz”; usunięcie (`session.delete`, rozdz. 6.4) niszczy zapis trwale i nieodwracalnie — okno eksponuje wyłącznie pierwsze |
| Archiwizacja a pobranie | „Archiwizuj” (`session.archive`) przenosi sesję między zbiorami widoku, zapis pozostaje w systemie, odzyskiwalny przez „Przywróć”; „Pobierz” (`history.load`, rozdz. 6.3) kopiuje zapis na zewnątrz systemu jako plik, nie zmieniając niczego w rejestrze |
| Sesja a okno komunikacji | Sesja (`Session`) jest jednostką organizacyjną trwałą, widoczną jako wiersz rejestru; okno komunikacji (`Window`) jest jednostką techniczną wewnątrz sesji, niewidoczną wprost w tym oknie poza tym, że historia (rozdz. 11) i „Pobierz” operują na jego identyfikatorze — Operator nigdy nie wybiera okna wprost, wybiera sesję |
| Eksponowana a nieeksponowana komenda | „Eksponowana” (rozdz. 10.1–10.2) znaczy: ma kontrolkę zweryfikowaną w prototypie tego okna; „nieeksponowana” znaczy: istnieje w kontrakcie, lecz to okno jej nie udostępnia — nie znaczy to, że komenda jest niedostępna z platformy w ogóle, wyłącznie że nie stąd |

---

## 17. Kryteria odbioru

| Warunek | Sposób sprawdzenia |
|---|---|
| Okno otwiera się jako nakładka z zaznaczoną sekcją „Historia sesji” spośród dziewięciu sekcji ustawień | wywołanie z menu aplikacji › Przejdź › Historia sesji |
| Wykaz grupuje sesje w trzy zbiory — czynne, zakończone, archiwalne — nie w zakładki stanu | otwarcie okna i sprawdzenie obecności trzech nagłówków `.hs-zbior-tytul` jednocześnie |
| Zakładki paska filtrów filtrują po środowisku (pięć pozycji), nie po stanie sesji | odczyt `role="tab"` pięciu zakładek i ich etykiet |
| Nagłówek kolumn jest dekoracyjny — `aria-hidden="true"`, bez reakcji na kliknięcie | próba kliknięcia nagłówka kolumny i sprawdzenie braku efektu |
| Sesja `paused` nie jest przedstawiana jako `finished` | sprawdzenie zbioru, do którego trafia sesja wstrzymana, po rozstrzygnięciu rozdz. 1.3 |
| Każdy wiersz niesie dokładnie dwa przyciski akcji, różne wg zbioru (rozdz. 6.1–6.3) | przegląd wiersza w każdym z trzech zbiorów |
| Brak mechanizmu zaznaczenia wielu wierszy i paska operacji zbiorczych | przegląd znaczników wiersza — brak `.dn-check`, brak licznika zaznaczenia |
| „Wróć” wywołuje `session.bind` + `session.focus`; „Otwórz” wywołuje `session.resume` | wznowienie sesji z obu zbiorów i podgląd wywołań kontraktu |
| „Zakończ” wywołuje `session.close` i przenosi sesję do zbioru „Sesje zakończone” | zakończenie sesji czynnej i obserwacja przejścia między zbiorami |
| „Archiwizuj” i „Przywróć” przenoszą sesję między zbiorami „Sesje zakończone” i „Sesje archiwalne” | wykonanie obu operacji i obserwacja |
| „Pobierz” czyta `history.load` po `windowId` dla okien sesji, bez zmiany stanu sesji | pobranie zapisu i podgląd wywołań oraz stanu sesji po operacji |
| Brak kontrolki „Usuń” w żadnym z trzech zbiorów | przegląd wszystkich przycisków akcji w rejestrze |
| Poniżej 900 px kolumny „Środowisko”, „Moduł” i „Zapis” znikają, lewa nawigacja staje się poziomym paskiem | zwężenie nakładki poniżej progu |
| Animacja tętna wiersza czynnego korzysta z klasy `.pt-tetno`, wyłączana przy `prefers-reduced-motion` | przegląd znacznika `.hs-znak` sesji czynnej i test preferencji ruchu |
| Usunięcie/zakończenie sesji nie usuwa jej projektu ani pamięci | wykonanie operacji na sesji projektu i sprawdzenie trwałości projektu |
| Filtr środowiska działa jednocześnie na wszystkich trzech zbiorach, nie zastępuje ich jednym płaskim wykazem | wybór zakładki środowiska i sprawdzenie, że trzy nagłówki grup pozostają widoczne z przeliczonymi licznikami |
| Trzy struktury `Session`, `SessionPresence`, `HistoryEntry` niosą pola zgodne z wykazem rozdz. 10.3 | podgląd surowej odpowiedzi `session.list` z `includePresence=true` i porównanie pól |
| Dziewiętnaście komend obszaru `session` są rozliczone jako eksponowane albo nieeksponowane, bez pominięć | porównanie wykazu rozdz. 10.1 z pełną listą komend obszaru w `contract.json` |
| Dziewięć sekcji lewej nawigacji okna ustawień jest wymienionych w kolejności zgodnej z prototypem | porównanie rozdz. 3.2 z `design/05-okna/platformowe/historia-sesji.html` |

---

## Załącznik A. Tabela zbiorcza — kontrolka, komenda, źródło

Zestawienie wszystkich kontrolek okna z komendą kontraktu, którą wywołują (albo kandydatką, gdy
mapowanie nie ma zrzutu jednoznacznego), oraz polem źródłowym prezentowanej wartości.

| Kontrolka | Dosłowne brzmienie | Komenda kontraktu | Pole / struktura źródłowa |
|---|---|---|---|
| Zakładka „Wszystkie” | „Wszystkie” | `session.list` (bez filtra środowiska) | — |
| Zakładka środowiska | „TalkIn” / „WorkSpace” / „CodeStudio” / „MultitaskingAI” | `session.list` (filtr klienta po `SessionPresence.environmentCode`) | `SessionPresence.environmentCode` |
| Pole wyszukiwania | „Szukaj w nazwach, projektach i modułach…” | `session.list` (filtr klienta) | `Session.title`, nazwa projektu, `SessionPresence.moduleCode` |
| Przycisk „Zakres dat” | „Zakres dat” | **[DO DECYZJI OPERATORA]** (rozdz. 5.3) | `Session.updatedAt` / `SessionPresence.lastActivityAt` |
| Przycisk „Sortowanie” | „Sortowanie” | **[DO DECYZJI OPERATORA]** (rozdz. 5.3) | — |
| Nagłówek zbioru „Sesje czynne” | „Sesje czynne” | `session.list` (`status=active`, `includePresence=true`) | `Session.status`, `SessionPresence.live` |
| Nagłówek zbioru „Sesje zakończone” | „Sesje zakończone” | `session.list` (`status=finished`) | `Session.status` |
| Nagłówek zbioru „Sesje archiwalne” | „Sesje archiwalne” | `session.archive.list` | `Session[]` |
| Kolumna Sesja | „Sesja” | `session.list` | `Session.title` |
| Kolumna Środowisko | „Środowisko” | `session.list` (`includePresence`) | `SessionPresence.environmentCode` |
| Kolumna Moduł | „Moduł” | `session.list` (`includePresence`) | `SessionPresence.moduleCode` |
| Kolumna Ostatnia praca | „Ostatnia praca” | `session.list` | `Session.updatedAt` / `SessionPresence.lastActivityAt` |
| Kolumna Zapis | „Zapis” | brak 1:1 (rozdz. 4.4) | **[DO DECYZJI OPERATORA]** |
| Znak tętniący | — | zdarzenie `session.changed` / `progress.changed` | `SessionPresence.live` |
| „Wróć” | „Wróć” | `session.bind` → `session.focus` | `Session`, `Window[]` |
| „Zakończ” | „Zakończ” | `session.close` | `Session` |
| „Otwórz” | „Otwórz” | `session.resume` | `Session`, `Window[]` |
| „Archiwizuj” | „Archiwizuj” | `session.archive` | `archivedIds` |
| „Przywróć” | „Przywróć” | `session.restore` | `restoredIds` |
| „Pobierz” | „Pobierz” | `history.load` (per `windowId`) | `HistoryEntry[]` |
| *(nieobecne)* Usuń trwale | — | `session.delete` — nieeksponowana w żadnym z trzech zbiorów (rozdz. 6.4) | `deletedIds`, `deletedCount` |

---

## Załącznik B. Słowniczek pojęć okna

| Pojęcie | Znaczenie w kontekście okna |
|---|---|
| Sesja | byt trwały rdzenia (struktura `Session`, rozdz. 10.3); żyje dłużej niż połączenie klienta |
| Stan sesji | jedna z czterech wartości `SessionStatus`: `active`, `paused`, `finished`, `archived` (rozdz. 1.3) |
| Zbiór widoku | jedna z trzech grup wykazu: „Sesje czynne”, „Sesje zakończone”, „Sesje archiwalne” (rozdz. 4) — nie tożsame z czterema wartościami `SessionStatus` |
| Zakładka zakresu | jedna z pięciu pozycji filtra środowiska: „Wszystkie” i cztery nazwy środowisk (rozdz. 5.1) |
| Obecność sesji | żywy stan sesji trwającej (`SessionPresence`); pole `live` rozstrzyga, czy proces trwa |
| Okno komunikacji | struktura `Window` (rozdz. 11.1); sesja skupia okna (`Session.windowIds`), historia jest per okno |
| Pozycja historii | struktura `HistoryEntry`; niesie `windowId`, `sessionId`, `role`, `preview` |
| Archiwum | zbiór sesji `archived`; zapis w całości, poza zbiorem „Sesje zakończone”; wgląd `session.archive.list` (rozdz. 4.3) |
| Projekt | jednostka grupująca sesje (`Session.projectId`); trwa niezależnie od pojedynczej sesji (rozdz. 15) |
| Zdarzenie | powiadomienie rdzenia o zmianie: `session.changed`, `session.focus.changed`, `history.changed`, `session.tool.attached`/`detached` |
| Kod błędu | jedna z ośmiu wartości katalogu błędów kontraktu, wspólnych całemu `contract.json`; rozstrzyga o zachowaniu okna (rozdz. 12) |
| Komenda nieeksponowana | komenda kontraktu istniejąca w obszarze `session`/`history`, bez kontrolki w zweryfikowanym prototypie tego okna (rozdz. 6.4, 10.1–10.2) |

---

*Koniec dokumentu. Okno historii sesji — Specyfikacja docelowa, wersja 2.0, 2026-08-20.*

---
*Danaco Console — Platforma AI Workspace OS · v2.0 · status Deweloperski*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — Dariusz Naharnowicz.*
*Warunki korzystania: [Licencja produktu](../LICENSE.md). Kontakt: support@danaco-group.pl*
