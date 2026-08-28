# Komplet dla dewelopera — pakiet pięciu okien

Stan na 28.08.2026. Gałąź `teren/centrum-poprawki` w `~/budowa`, 25 rewizji przed `main`.
Wszystkie liczby niżej **zmierzone** na tym drzewie. Gdzie pomiaru nie było, napisano „nie zmierzono".

Zakres nadany przez Właściciela: **tylko te pięć okien**. Pozostałe 31 prototypów w
`design/05-okna/` jest poza pakietem — nie zmieniano ich i nie należy ich zmieniać
razem z tym wydaniem.

---

## 1. Skład pakietu

| Okno | Plik | Wierszy | Menu | Pozycji menu | SVG |
|---|---|---:|---:|---:|---:|
| Instalator | `design/05-okna/platformowe/instalator.html` | 106 | — | — | 1 |
| Uruchomienie i rejestracja | `design/05-okna/przeplyw/przeplyw-wejscia.html` | 365 | — | — | 18 |
| Centrum dowodzenia | `design/05-okna/przeplyw/centrum-dowodzenia.html` | 1523 | 40 | 377 | 602 |
| TalkIn | `design/05-okna/srodowiska/talkin.html` | 1037 | 11 | 159 | 312 |
| Studio | `design/05-okna/moduly/studio.html` | 990 | 11 | 159 | 323 |

Pozycje liczone po `.sta-menu-poz`. W TalkIn i Studio każda pozycja niesie rolę
`menuitem*` (159 = 159). W Centrum rola pada 379 razy — o dwa więcej niż pozycji, bo dwa
elementy niosą ją poza rodziną `.sta-menu-poz`. Przy kontroli trzymać się jednego licznika.

**Instalator i Przepływ wejścia nie mają treści w pliku HTML** — okno składa JavaScript
ze składników w `design/zasoby/okna/instalator/` i `design/zasoby/okna/wejscie/`.
Kto szuka znacznika w pliku okna, znajdzie samo rusztowanie prototypu. To najczęstsza
pułapka przy tym pakiecie: skan klas po plikach `05-okna/**` pomija te dwa okna.

---

## 2. Co pakiet wczytuje

Warstwy są wspólne dla wszystkich pięciu: `zetony/fonty.css` → `zetony/zetony.css`
(wciąga `zetony/ruch.css` przez `@import`) → `css/fundament.css` → `css/komponenty.css`.
Dalej różnią się.

| Okno | Arkusze ponad warstwę wspólną | Skryptów |
|---|---|---:|
| Instalator | `kreator.css` | 37 |
| Uruchomienie | `wejscie.css` | 28 |
| Centrum | `menu`, `stany`, `rama`, `karty-okna`, `panel-sesji`, `pasek-okna`, `izolacja`, `okno-robocze`, `stanowisko`, `okna-modalne` + 4 arkusze `okna/centrum-*` + `okna/danaco-anim-3d` | 18 |
| TalkIn | `stanowisko`, `rama`, `okna-modalne`, `okna/talkin.css` | 5 |
| Studio | `stanowisko`, `rama`, `okna-modalne`, `karty-okna`, `okna/studio.css` | 8 |

`prototyp.css` i `prototyp.js` to **rusztowanie prototypu** — panel wyboru widoku,
uchwyt z nazwą pliku, noty projektowe. W produkcie ich nie ma; znacznik, który je
niesie, nie wchodzi do implementacji.

**Okno przed uwierzytelnieniem nie wczytuje `rama.css`.** Powłoka aplikacji dochodzi
dopiero po wydaniu tokenu — w prototypie `przeplyw-wejscia.html` doczepia ją skryptem
przy wejściu w etap 3. Dlatego rodzina `.dn-okno-wejsciowe` ma **własne** części belki
(znak, tytuł, uchwyt, sterowanie oknem), osobne od `.dn-belka-*` z ramy. To nie jest
duplikat do złożenia.

---

## 3. Stan zmierzony

| Miara | Wynik | Czym mierzone |
|---|---|---|
| axe WCAG 2.0/2.1/2.2 A+AA, oba motywy, 2560 px | **0 krytycznych i poważnych** we wszystkich pięciu | `axe-core` 4.x przez Playwright |
| Błędy skryptu w konsoli | **0** w 36 oknach drzewa | Playwright, `pageerror` + `console.error` |
| Odpowiedzi 4xx/5xx | **0** | Playwright, `response` |
| Przewijanie poziome przy 2560 px | **brak** we wszystkich pięciu | `scrollWidth > clientWidth` |
| Przyciski bez nazwy dostępnej | **0** | brak tekstu i `aria-label`/`aria-labelledby` |
| Pole dotyku < 24 px (WCAG 2.5.8) | **0** — 11 kontrolek ma rysunek mniejszy, wszystkie mają pole 24 px | trafienie `elementFromPoint` 11 px od środka |
| Zasada trzech stref (okna przedaplikacyjne) | **dotrzymana** we wszystkich krokach i widokach, oba motywy | `~/robocze/pomiar/strefy-okna.mjs` |
| Kolizje skrótów | Centrum **2**, TalkIn **3**, Studio **3** — wszystkie odziedziczone po WZORCU, wypisane w rozdz. 6 | grupowanie `<span class="skrot">` po treści |
| Bloki `<style>` odrzucone przez sprawdzenie | **0 z 5 okien pakietu** | `~/robocze/pomiar/sprawdzenie-skladnikow.sh` |
| Deklaracje wyglądu w arkuszach okien pakietu | `wejscie.css` **133**, `okna/studio.css` **75** | jw. |

Ogniskowalnych kontrolek: Instalator 25, Uruchomienie 130, Centrum 682, TalkIn 412, Studio 441.

---

## 4. Co zmieniono — 25 rewizji

Kolejność chronologiczna; pełne opisy w `git log main..teren/centrum-poprawki`.

**Dostępność i znacznik (1–14).** Semantyka zakładek, przełączników, list wyboru i pól
formularza: 127 naruszeń znacznika → 2. Kontrast obu motywów do progu AA: 691 naruszeń → 0
(jeden żeton `--dn-tekst-3` w motywie ciemnym niósł 381 z 389). Dwa wyjątki JavaScript
przerywające obsługę okien. Granice kontrolek — 39 przycisków brało `buttonface` przeglądarki.

**Biblioteka (15–20).** Wyniesienie wyglądu z okien do `komponenty.css`: belka okna
przedaplikacyjnego, kafel wyboru ze wstęgą, krok w bryle (`.dn-krok--pole`), skala pisma,
rytm panelu wejścia. Każde przeniesienie zmierzone odciskiem geometrii przed i po.

**Menu i wstążka (21–25).** Domknięcie 237 pozycji menu centrum wobec wzorca, 33 okien,
manifestu ikon i kontraktu. Rozstrzygnięcie skrótów `Ctrl+V` / `Ctrl+W`. Zestaw narzędzi
wstążki właściwy modułowi karty. Pole dotyku 24 px. Rola `menuitem` dla siedmiu wyzwalaczy
podmenu w TalkIn i Studio — stały w `role="menu"` jako zwykły przycisk.

### 4.1. Menu centrum — szczegółowo

- **47 brzmień poprawionych.** Każde poprawione stoi teraz w 34 plikach, każde zastąpione
  w zerze — zmierzone. Reguła kontroli: `grep -rl "<brzmienie>" design/05-okna --include='*.html' | wc -l` = 34.
- **5 pozycji usuniętych**: grupa „Szybkie przełączniki" z menu ustawień (Motyw, Gęstość,
  Nie przeszkadzać, Always On Display — powielenia przycisku szyny albo nastawy; grupy nie ma
  we WZORCU ani w żadnym z 33 okien) oraz „Archiwum sesji" z menu panelu.
- **22 pozycje dopisane**: cztery w menu profilu (Powiadomienia konta + trzy stany obecności),
  grupa „Ostatnie pozycje" schowka, **„Wyjmij z projektu" we wszystkich 14 menu wierszy sesji**
  (`session.project.clear` z kontraktu), trzy w menu pasma kart, „Samouczek" w oknie dostosowania.
- **16 ikon podmienionych** na właściwe wg manifestu.
- Pozycji menu 358 → 377.

**Uchwyty do podpięcia po stronie aplikacji** (dziś prowadzą do komunikatu prototypu):
`data-poz-akcja="wyjmij"` → `session.project.clear`; pozycje `Zamknij karty po prawej`,
`Przekaż na Mobile`, `Zapisz zestaw kart…` w `menu-karty-otwarte` nie mają jeszcze uchwytu.

### 4.2. Wstążka okna roboczego

`pasek-okna.js` deklarował w swoim kontrakcie „zestaw narzędzi właściwy karcie bieżącej",
a niósł jeden stały skład Studia — karta Research pokazywała narzędzia Studia pod własną nazwą.

Wprowadzono `.dn-pasek-okna-zestaw[data-zestaw-modul]`: pasek nosi tyle zestawów, ile modułów
ma karta w oknie, widoczny jest ten, którego moduł prowadzi kartę bieżącą. Lustra menu nadmiaru
przebudowują się przy zmianie zestawu.

| Zestaw | Narzędzia |
|---|---|
| Studio (10) | Zapisz · Cofnij · Ponów ‖ Rewizje · Walidacja · Przypisy ‖ Terminal · Różnice ‖ Udostępnij · Drukuj |
| Research (8) | Wyszukiwanie źródeł · Dodaj źródło · Widok lektury ‖ Nowe ustalenie · Budowa raportu ‖ Eksportuj raport · Udostępnij · Drukuj |

Zestaw Research wyprowadzony z sześciu paneli `docs/moduly/research.md` (Discovery Panel,
Sources Manager, Reading View, Findings Panel, Report Builder, Export Panel). Trzy narzędzia
mają skrót nadany w `08-handoff-motion-dostepnosc.md` rozdz. 18.2: dodanie źródła `Ctrl+K`,
nowe ustalenie `Ctrl+N`, eksport `Ctrl+E`. **Skróty nie są jeszcze wpisane w znacznik** —
wstążka pokazuje same znaki z etykietką.

Kolejne moduły dokłada się jednym blokiem `.dn-pasek-okna-zestaw` — skrypt nie wymaga zmian.

---

## 5. Rozstrzygnięcia Właściciela wiążące dla tego pakietu

| Rzecz | Rozstrzygnięcie | Skutek |
|---|---|---|
| `Ctrl+V` | należy do **schowka** | `Wklej ostatnią pozycję` trzyma `Ctrl+V`; zwykłe `Wklej` traci skrót powłoki (obsługuje je system) |
| `Ctrl+W` | należy do **okna** | `Zamknij okno` = `Ctrl+W`; zamknięcie karty przeszło na `Ctrl+F4` |
| Zestaw 152 ikon | **nie wiąże znacznika** | Kształty wpisane wprost w okna są dopuszczone. Pozycja zdjęta z długu, nie mierzy się jej w audytach. Zestaw zostaje biblioteką do użycia, nie normą |
| Zakres | **tylko pięć okien** | Pozostałe 31 prototypów bez zmian, także tam, gdzie odbiegają od tego pakietu |

---

## 6. Otwarte — z liczbami, bez rozstrzygnięcia

### 6.1. Kolizje skrótów odziedziczone po WZORCU

Nie są usterkami tych okien — pochodzą ze składnika wspólnego i stoją tak samo we wszystkich 34.

| Okno | Skrót | Dwie czynności |
|---|---|---|
| Centrum | `Ctrl+F4` | `Zamknij kartę` / `Zamknij kartę sesji` — jedna czynność, dwie nazwy; WZORZEC nie zna żadnej |
| Centrum | `Ctrl+Shift+A` | `Always On Display` (menu Więcej) / `Powierzchnia interakcji` (menu AOD) — rozeszły się po poprawieniu nazwy wg WZ:319 |
| TalkIn, Studio | `Ctrl+O` | `Otwórz projekt…` / `Otwórz projekt istniejący` |
| TalkIn, Studio | `Ctrl+Shift+N` | `Nowe okno aplikacji` / `Powtórz ostatnią sesję` — obie nadane przez WZORZEC |
| TalkIn, Studio | `F5` | `Odśwież widok` / `Ponowne wczytanie widoku` — **obie w tym samym menu** |

### 6.2. Wyglądy wciąż w arkuszach okien

`design/zasoby/wejscie.css` — 133 deklaracje, `design/zasoby/okna/studio.css` — 75.
Rodziny do wyniesienia: `we-marka*` (22), `pg-warstwa*` (17), `au-krok*` (15), `au-kod*` (13),
`au-sila*` (9), `au-haslo*` (7), `we-tytul`/`we-nadtytul` (10). Wzór postępowania — trzy
rewizje 17–20: rodzinę przenieść do `komponenty.css` jako odmianę, arkusz okna odchudzić,
zmianę zmierzyć odciskiem geometrii 22 widoków w obu motywach.

`design/zasoby/okna/danaco-anim-3d.css` (89 deklaracji) to **sankcjonowany wyjątek** —
rejestr terenów w. 28–29.

### 6.3. Rozbieżności wobec dokumentacji — świadomie nietknięte

Wymagają decyzji obejmującej więcej niż ten pakiet.

- **`menu-aod`**: WZORZEC ma 1 grupę i 7 pozycji, Centrum 4 grupy i 11.
- **Rozwinięcie MultitaskingAI**: 8 pozycji wobec 6 sekcji z `docs/srodowiska/multitaskingai.md:636–653`.
- **Kolejność modułów WorkSpace w szynie** różna od `docs/srodowiska/workspace.md:145–153`.
  Oba stany identyczne we wszystkich 34 oknach.
- **Cztery pozycje szyny, w których Centrum jest zgodne z katalogiem, a 33 okna odbiegają**:
  Browser `karta-okna` wobec `globus` ×66, Workspace `warstwy` wobec `folder` ×132,
  Automations `automatyzacja`, pisownia „Always On Display" (844 wystąpienia wobec 136).
- **`Zmień nazwę` i `Usuń trwale` w menu projektu** nie mają komendy w kontrakcie
  (`workspace.project.*` to wyłącznie `status.set`).
- **`Eksportuj wykaz…`** działa w prototypie, ale kontrakt nie ma komendy wydającej wykaz sesji.
- **Stany sesji**: `SessionStatus` kontraktu ma `active/paused/finished/archived`, filtr okna
  ma `Wszystkie/W toku/Czekające/Zakończone`. `paused` nie ma odpowiednika w oknie, „Czekające"
  (`data-stan="reakcja"`) nie ma odpowiednika w kontrakcie.

Pełny wykaz 16 pytań z uzasadnieniem i pomiarem: `~/robocze/prowadzenie/komplet-menu-centrum.md`, rozdz. 6.

---

## 7. Jak to mierzyć

Serwer statyczny na `127.0.0.1:8611` z korzeniem w `~/budowa`. Wszystkie polecenia z `~/budowa`.

```bash
bash ~/robocze/pomiar/sprawdzenie-skladnikow.sh
```

```bash
PLAYWRIGHT_BROWSERS_PATH=/opt/ms-playwright node ~/robocze/pomiar/strefy-okna.mjs
```

Miary niezmiennicze centrum — muszą dać tę samą wartość przed i po każdej zmianie:
menu i wyzwalacze **40 / 40 bez sierot**, bilans `<button>` **566**, bilans `<div>` **314**,
pozycji `.sta-menu-poz` **377**, elementów z rolą `menuitem*` **379** (dwa niosą rolę
poza rodziną pozycji menu).

### Pułapki zmierzone w tej pracy

1. **Sprawdzenie obcina wykaz na 25 pozycjach** („… i N dalszych"). Pełną listę daje kopia
   skryptu z `head -25` podniesionym.
2. **Bilans `<div>` trzeba sprawdzać osobno.** Usunięcie wiersza menu zdjęło `</div>`
   domykające menu, bo stało na końcu ostatniej pozycji, nie w osobnym wierszu. Bilans
   `<button>` i wykaz menu tego nie wykryły — dopiero pomiar geometrii stopki szyny (209 → 40 px).
3. **Porównanie zrzutów jest bezużyteczne** — okna mają animacje wejścia, dwa zrzuty tego
   samego stanu różnią się. Mierzyć geometrię, nie piksele. Sondy odciskowe uruchamiać
   z `reducedMotion: 'reduce'`.
4. **`::after` przycisku ikonowego nosi dymek etykietki** (`rama.css:467`). Nakładka pola
   dotyku musi iść przez `::before`.
5. **Skan klas musi obejmować `design/zasoby/okna/**` rekurencyjnie** — Instalator
   i Przepływ wejścia składa JavaScript, w plikach HTML klas nie ma.
6. **Po zmianie CSS/JS podbić `?w=` w odsyłaczach okna.** Konwencji używa wyłącznie
   `centrum-dowodzenia.html` (38 odsyłaczy) — Właściciel ogląda prosto z dysku, przez `file://`.
7. **Szerokość odniesienia to ~2560 px**, nie 1600.

---

## 8. Dokumenty towarzyszące

| Plik | Co zawiera |
|---|---|
| `~/robocze/prowadzenie/komplet-menu-centrum.md` | 483 wiersze: pozycja po pozycji, z wierszem pliku i źródłem; 16 pytań; czego nie robić i dlaczego |
| `~/robocze/przekazanie-centrum-i-motyw-jasny.md` | historia kontroli warstwy wspólnej, motyw jasny, kontrast |
| `~/robocze/audyt-design-korekty-przekazanie.md` | poprzednia tura — 72 ustalenia audytu |
| `design/05-okna/NORMA-OKNA.md` | zasada zera wyglądu w oknie, zasada trzech stref |
| `design/zasoby/ARCHITEKTURA.md` | warstwy arkuszy i kolejność wczytywania |
