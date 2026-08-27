# Dług składników — lista robocza

Rozbiór **229 własnych składników** czterech okien wobec biblioteki 261 składników.
Podstawa: [NORMA-OKNA.md](NORMA-OKNA.md) — zasada zera.

| kubełek | ile | robota |
|---|---:|---|
| duplikat biblioteki | 92 | skreślić z okna, wpisać `.dn-*` |
| martwy | 12 | skreślić — nieużywany |
| czyste rozmieszczenie | 56 | zostaje w oknie |
| do biblioteki | 45 | 33 rodzin · 6 wariantów · 6 części |

**104 składników znika bez śladu w wyglądzie** — powielają bibliotekę albo są martwe.

---

## 0 · Role uwięzione w rodzinach — ZROBIONE

Przyczyna głębsza niż niedbalstwo. Biblioteka **miała** te role, ale każda była zamknięta
w jednej rodzinie i nie dało się jej użyć gdzie indziej. Okno, które potrzebowało zwykłej
noty pod tytułem, nie miało czego wziąć — więc pisało własną.

Sprawdzenie ręczne zawęziło listę z pięciu do **trzech**: `.dn-kafel-etykieta` jest na swoim
miejscu, a `.dn-tekst-ciagly` to osobna nowa rodzina, nie uogólnienie modalu.

| wyniesione do biblioteki | uwolnione z | rola |
|---|---|---|
| `.dn-meta` + `--mocna` `--na-ramie` `--w-wierszu` | `.dn-krok-meta`, `.dn-stan-miara` | nota maszynowa — metryka, miara, nota |
| `.dn-etyk-mono` + `--mocna` `--na-ramie` `--z-czynnoscia` | `.dn-aod-sekcja-glowa` | etykieta mono-wersaliki |
| `.dn-lid` + `--bez-odstepu` `--na-ramie` | `.dn-modal-lid` | zdanie wprowadzające pod tytułem |

Przy okazji uzupełniona rodzina `.dn-wybor`, która nie miała **żadnych** części — stąd własne
klasy wyboru wariantu w instalatorze:

| dopisane | rola |
|---|---|
| `.dn-wybor--blokowy` | wybór dwuwierszowy: kontrolka u góry, obok kolumna treści |
| `.dn-wybor-nazwa` | nazwa opcji, pismo półgrube |
| `.dn-wybor-opis` | zdanie objaśniające pod nazwą |

Klasy zastane renderują się bez zmiany — sprawdzone zrzutami czterech okien w dwóch
motywach: **różnica 1 piksela przy progu szumu 10**.

---

## 0b · Rodzina okna wejściowego — NIEZROBIONE

Bryła, z której składają się instalator, okno startowe, logowanie i przygotowanie
środowiska. Dotąd każde z tych okien miało własną kopię.

**Stan: rodzina dodana do biblioteki, migracja okien niewykonana.** Składniki
`.dn-*` z tabeli stoją w bibliotece, ale klasy z kolumny „zastępuje" **nie
zostały skreślone** — `.we-scena`, `.we-okno`, `.we-belka`, `.we-panel`,
`.we-akcje`, `.au-akcje`, `.we-marka` (z `-poz`, `-stopka`, `-godlo`), `.pg-znak`,
`.pd-tytul` i `.au-link` nadal mają definicje w `wejscie.css` i `przedsionek.css`,
obok swoich zamienników. Zamiana wejdzie przy pracy nad tymi oknami; do tego czasu
pozycje pozostają otwarte (por. §1).

| składnik | warianty | zastępuje |
|---|---|---|
| `.dn-scena` | `--gora` | `.we-scena` |
| `.dn-okno-wejsciowe` | `--z-belka`, `--jednokolumnowe` | `.we-okno`, `.in-okno` |
| `.dn-okno-wejsciowe-belka` | — | `.we-belka` |
| `.dn-okno-wejsciowe-panel` | `--scentrowany` | `.we-panel` |
| `.dn-akcje` | `--poboczne` | `.we-akcje`, `.we-akcje-poboczne`, `.au-akcje` |
| `.dn-pas-dzialan` | `--ciag` | `.we-stopka`, `.we-stopka--ciag` |
| `.dn-tozsamosc` `-poz` `-stopka` | — | `.we-marka`, `.we-marka-poz`, `.we-marka-stopka` |
| `.dn-godlo` | `--lg` `--ikonowy` `--lockup` `--tetno` `--na-plotnie` `--sygnal` | `.we-marka-godlo`, `.pg-znak`, `.cd-marka-lockup`, `.st-pasmo-znak` |
| `.dn-tytul-okna` | `--sm`, `--na-ramie` | `.cd-tytul`, `.pd-tytul` |
| `.dn-link` | `--drugorzedny`, `--na-ramie` | `.au-link` |

Wymiary okna i skale godła wyniesione do żetonów: `--dn-wym-okno-wejsciowe-szer`,
`-wys`, `-kolumna`, `--dn-wym-godlo`, `-lg`, `-lockup`, `-ikonowy`.

**Biblioteka: 261 → 303 składników.**

### Usterka znaleziona przy próbie składania

Rozdzielenie stref w motywie **ciemnym** praktycznie nie istniało: kolumna tożsamości
`#131313` przy panelu `#181818` to różnica **2,4 stopnia L\***. W jasnym motywie ta sama
para dawała 94. Panel przeniesiony na żeton `--dn-panel` — różnica **6,9 L\*** w ciemnym,
92,4 w jasnym. Rozdzielenie nadal idzie samą powierzchnią, bez kreski.

### Próba składania

[`_samodzielne/proba-zera.html`](_samodzielne/proba-zera.html) — okno logowania złożone
**wyłącznie** z biblioteki: zero klas własnych, zero bloku `<style>`, zero arkusza okna.


---

## 0c · Braki odsłonięte przez kroki 3, 5 i 6 instalatora — ZROBIONE

Treść finalna trzech kroków wymagała ról, których biblioteka nie miała. Wszystkie
weszły do warstwy wspólnej **przed** użyciem w oknie; okno instalatora nie ma
bloku `<style>` ani własnego arkusza dla tych kroków.

| dopisane | arkusz | rola |
|---|---|---|
| `[hidden]` | `css/fundament.css` | atrybut `hidden` wygrywa z `display` rodziny — dotąd każda rodzina łatała to osobno (`.dn-modal-tlo[hidden]`, `.dn-pasek-okna[hidden]`, dwa `:not([hidden])`) |
| `.dn-krok-stan` | `css/komponenty.css` | słowo stanu w kolejce kroków — pismo bezszeryfowe, bo status to tekst interfejsu, nie miara jak `.dn-krok-meta` |
| `.dn-postep-tor--nieokreslony` | `css/komponenty.css` | tor bez miary — wycofywanie zmian, oczekiwanie bez procentu |
| `@keyframes dn-postep-przebieg` | `zetony/ruch.css` | przebieg toru bez miary |
| `.dn-tytul-ze-znakiem` `.dn-tytul-znak` `--sukces` `--ostrzezenie` | `css/komponenty.css` | tytuł ze znakiem wyniku w wierszu — `svg` jest w pakiecie blokowy, więc rząd musi zbudować rodzina |
| `button.dn-link` | `css/komponenty.css` | odsyłacz jako czynność, nie adres — rodzina zdejmuje bryłę przycisku |
| `.dn-bryla` + `zasoby/bryla.js` | `css/komponenty.css` | bryła warstw platformy — scena przestrzenna w płótnie, barwy z żetonów, postęp podawany przez okno |

Reguła `[hidden]` czyni zbędnymi cztery łatki rodzinowe. Zdjęta została jedna
(`.dn-modal-tlo[hidden]`); pozostałe stoją w `pasek-okna.css` i `panel-sesji.css`
i są do zdjęcia przy pracy nad tamtymi oknami.

---

## 0d · Rodzina okna kreatora — ZROBIONE

Instalator nie miał w bibliotece swojej bryły. Zbieżność nazw w raporcie strażnika
myliła: `.dn-belka` z `rama.css` to belka **ramy aplikacji**, której okna sprzed
uwierzytelnienia celowo nie mają, a `.dn-okno-wejsciowe` ma dwie kolumny i nie ma
pasa działań na dnie. Okno kreatora to trzy pasma, nie dwie kolumny.

Cała rama instalatora wyszła z arkusza okna do biblioteki, do nowego arkusza
rodziny [`zasoby/kreator.css`](../zasoby/kreator.css) — 48 reguł:

```
.dn-kreator              .dn-kreator-korpus       .dn-kreator-nadtytul
.dn-kreator-belka        .dn-kreator-szyna        .dn-kreator-tytul
.dn-kreator-belka-znak   .dn-kreator-kroki        .dn-kreator-lid
.dn-kreator-belka-tytul  .dn-kreator-krok         .dn-kreator-dokument
.dn-kreator-belka-btn    .dn-kreator-krok-nr      .dn-kreator-wybory
.dn-kreator-belka-sterowanie                      .dn-kreator-bryla
.dn-kreator-plotno       .dn-kreator-ekran
.dn-kreator-pas          .dn-kreator-pas-tresc    .dn-kreator-pas-odstep
```

Deklaracje przeniesione **bez zmiany** — zmieniły się wyłącznie nazwy selektorów.
Dowód: porównanie zrzutów płótna treści dla pięciu kroków w dwóch motywach daje
**0 pikseli różnicy**; rozbieżność w całym oknie pochodzi z losowych cząstek
bryły i mieści się w szumie dwóch uruchomień tego samego kodu.

Instalator nie ma już arkusza okna. Wpina żetony, fundament, komponenty,
`kreator.css` i `prototyp.css`; jego mechanika stoi w `okna/instalator.js`.
W znaczniku zostało **227 klas `dn-*` i zero klas własnych**.

Do przeniesienia zostaje reszta `wejscie.css` — rodziny `.we-*`, `.au-*`, `.pg-*`
(okno startowe, logowanie, przygotowanie środowiska). Przy okazji: belka okna
systemowego stoi teraz w dwóch kopiach — `.we-belka-btn` i `.dn-kreator-belka-btn`
— bo reguła była wspólna i podział ją rozdwoił. To jeden składnik i docelowo ma
mieć jedną nazwę.

---

## 1 · Do skreślenia — duplikaty biblioteki

| własna klasa | zastąpić przez | użyć razy |
|---|---|---:|
| `.au-etykieta` | `.dn-menu-naglowek` | 2 |
| `.au-haslo-oko` | `.dn-btn-ikona--sm` | 6 |
| `.au-kod-adres` | `.dn-meta` ⚠ | 1 |
| `.au-krok-nr` | `.dn-krok-znak` | 6 |
| `.au-metoda` | `.dn-kafel` | 6 |
| `.au-metoda-stan` | `.dn-kafel-opis` | 6 |
| `.au-odliczanie` | `.dn-meta` ⚠ | 1 |
| `.au-opcja` | `.dn-check-etyk` | 5 |
| `.au-sila-opis` | `.dn-pole-opis` | 1 |
| ~~`.cd-ik`~~ | `.dn-ik` — ZROBIONE | 5 |
| ~~`.cd-listwa`~~ | `.dn-listwa` — ZROBIONE | 1 |
| ~~`.cd-listwa-poz`~~ | `.dn-listwa-pozycja--pulpit` — ZROBIONE; podniesienie i stan `aria-pressed` weszły do biblioteki jako wariant | 3 |
| ~~`.cd-listwa-sep`~~ | `.dn-separator--pionowy` — ZROBIONE. Nie `.dn-stan-sep`: tamten ma sztywne 14 px pod pasek stanu, a separator listwy rozciąga się na wysokość pasa | 2 |
| `.cd-modul-sesja` | `.dn-meta` ⚠ | 1 |
| `.cd-nauka-kroki` | `.dn-kolejka` | 1 |
| `.cd-nauka-nota` | `.dn-dost-nota` | 1 |
| `.cd-nauka-numer` | `.dn-krok-znak` | 4 |
| `.cd-nauka-wstep` | `.dn-tekst-ciagly` ⚠ | 1 |
| `.cd-plotno` | `.dn-obszar-tresc` | 1 |
| `.cd-start` | `.dn-alert` | 1 |
| `.cd-tresc` | `.dn-wejscie-kaskada` | 2 |
| `.cd-wlasne-pusto` | `.dn-pusty-stan` | 1 |
| `.cd-wlasny-nazwa` | `.dn-obszar-pozycja-nazwa` | 6 |
| `.cd-wlasny-poz` | `.dn-kafel` | 6 |
| `.cd-wlasny-rodzaj` | `.dn-kafel-opis` | 6 |
| `.cd-wlasny-znak` | `.dn-kafel-ikona` | 6 |
| `.in-dzialania` | `.dn-modal-stopka` | 1 |
| `.in-faza` | `.dn-krok` | 1 |
| `.in-faza-stan` | `.dn-meta` ⚠ | 1 |
| `.in-fazy` | `.dn-kolejka` | 1 |
| `.in-krok` | `.dn-krok` | 6 |
| `.in-krok-numer` | `.dn-krok-znak` | 6 |
| `.in-kroki` | `.dn-kolejka` | 1 |
| `.in-odsylacze` | `.dn-tekst-ciagly` ⚠ | 1 |
| `.in-okno` | `.dn-modal` | 1 |
| `.in-pole` | `.dn-pole` | 1 |
| `.in-pole-etykieta` | `.dn-pole-etykieta` | 1 |
| `.in-postep-stan` | `.dn-postep-wartosc` ⚠ | 1 |
| `.in-postep-tor` | `.dn-postep-tor` | 1 |
| `.in-powitanie` | `.dn-tekst-ciagly` ⚠ | 1 |
| `.in-sciezka` | `.dn-pole-kontrolka` | 1 |
| `.in-tresc` | `.dn-modal-cialo` | 1 |
| `.in-wariant` | `.dn-karta` | 2 |
| `.in-wariant-nazwa` | `.dn-wybor-nazwa` ⚠ | 2 |
| `.in-wariant-opis` | `.dn-kafel-opis` | 4 |
| `.in-wybor` | `.dn-wybor` | 5 |
| `.in-wybor-opis` | `.dn-pole-opis` | 3 |
| `.in-wynik-zdanie` | `.dn-tekst-ciagly` ⚠ | 1 |
| `.pg-postep-wiersz` | `.dn-postep-glowa` ⚠ | 1 |
| `.pg-warstwa-etykieta` | `.dn-menu-naglowek` | 4 |
| `.st-cialo` | `.dn-obszar` | 1 |
| `.st-edytor-grupa` | `.dn-narzedzia-grupa` | 3 |
| `.st-edytor-pasek` | `.dn-obszar-pasek` | 1 |
| `.st-karta` | `.dn-karta-sesji` | 7 |
| `.st-karta--robocza` | `.dn-karta-widoku--glowna` | 2 |
| `.st-kartka` | `.dn-karta` | 1 |
| `.st-kartka-numer` | `.dn-pole-opis` | 1 |
| `.st-kartka-tytul` | `.dn-karta-tytul` | 1 |
| `.st-karty` | `.dn-karty-lista` | 2 |
| `.st-meta` | `.dn-meta` ⚠ | 19 |
| `.st-miara` | `.dn-meta` ⚠ | 1 |
| `.st-nota` | `.dn-dost-nota` | 2 |
| `.st-okno-robocze` | `.dn-obszar-panel--glowny` | 2 |
| `.st-panel-lista` | `.dn-pole-grupa` | 6 |
| `.st-pasmo-sterowanie` | `.dn-karty-narzedzia` | 1 |
| `.st-status` | `.dn-stan` | 1 |
| `.st-status-prawa` | `.dn-stan-poz--prawa` | 1 |
| `.st-szyna` | `.dn-boczna` | 1 |
| `.st-szyna-grupa` | `.dn-obszar-grupa` | 4 |
| `.st-szyna-lista` | `.dn-obszar-tresc` | 1 |
| `.st-szyna-naglowek` | `.dn-obszar-glowa` | 1 |
| `.st-wstazka` | `.dn-obszar-pasek` | 1 |
| `.st-wstazka-grupa` | `.dn-narzedzia-grupa` | 3 |
| `.st-wstazka-grupa--drugorzedna` | `.dn-narzedzia-grupa--proba-1400` | 1 |
| `.st-wstazka-stan` | `.dn-stan-poz` | 1 |
| `.we-alarm` | `.dn-alert` | 6 |
| `.we-alarm--informacja` | `.dn-alert--info` | 2 |
| `.we-alarm--ostrzezenie` | `.dn-alert--ostrzezenie` | 1 |
| `.we-belka` | `.dn-belka` | 3 |
| `.we-belka-tytul` | `.dn-belka-tytul` | 3 |
| `.we-belka-zamknij` | `.dn-belka-btn--zamknij` | 3 |
| `.we-belka-znak` | `.dn-belka-marka` | 3 |
| `.we-krok` | `.dn-krok` | 17 |
| `.we-krok-meta` | `.dn-meta` ⚠ | 17 |
| `.we-krok-znak` | `.dn-krok-znak` | 17 |
| `.we-kroki` | `.dn-kolejka` | 4 |
| `.we-lid` | `.dn-tekst-ciagly` ⚠ | 6 |
| `.we-marka-motto` | `.dn-menu-naglowek` | 4 |
| `.we-marka-nazwa` | `.dn-karta-srodowiska-tytul` | 4 |
| `.we-nadtytul` | `.dn-menu-naglowek` | 4 |
| `.we-stopka` | `.dn-modal-stopka` | 10 |
| `.we-tytul` | `.dn-modal-tytul` | 16 |

⚠ — orzeczenie poprawione ręcznie: agent wskazał klasę uwięzioną w cudzej rodzinie (punkt 0).

## 2 · Do skreślenia — martwe

| klasa | dowód |
|---|---|
| `.cd-archiwum` | Zero wystapien w .html i .js. Caly blok (border:1px dashed var(--dn-obrys); border-radius: |
| `.cd-archiwum-tresc` | Zero wystapien w .html i .js. Blok (padding, color:var(--dn-tekst-2), fs-sm, lh-luzny, ani |
| `.cd-karta-menu` | Zero wystapien w atrybutach class w .html. Jedyne odwolanie w calym drzewie to selektor ob |
| ~~`.cd-listwa-menu`~~ | ZROBIONE — skreślona z centrum-dowodzenia.css i centrum-obszar.css |
| `.cd-nadtytul` | Zero wystapien w atrybutach class w .html i zero w .js w calym /home/ubuntu/robocze/protot |
| `.cd-naglowek-akcje` | Zero wystapien w .html i .js. Blok display:flex; align-items:center; gap:var(--dn-od-2) ni |
| `.cd-podtytul` | Zero wystapien w .html i .js. (Blok: margin, max-width:78ch, color:var(--dn-tekst-2), fs-b |
| `.cd-sekcja-meta` | Zero wystapien w .html i .js. Blok margin-left:auto; ff-mono; fs-xs; color:var(--dn-tekst- |
| `.cd-sekcja-opis` | Zero wystapien w .html i .js. Blok max-width:74ch; color:var(--dn-tekst-2); font-size:var( |
| `.cd-zaslona` | Zero wystapien w .html i .js. Blok display:grid; place-content:center; padding:var(--dn-od |
| `.cd-zaslona-tresc` | Zero wystapien w .html i .js. Blok (flex, gap, max-width:64ch, color:var(--dn-tekst-3), fs |
| `.in-sciezka-stala` | grep po .html i .js w /home/ubuntu/robocze/prototypy/design nie znajduje ani jednego wysta |

---

## 3 · Do dopisania do biblioteki
### 3a · Nowe rodziny

Do `zasoby/css/komponenty.css`, z pełnym kompletem stanów, wyłącznie na żetonach `--dn-*`.

**`.dn-akcje-tekst`** — Slot tekstu ciągłego w pasie działań — objaśnienie albo odsyłacz stojący obok przycisków, wypychany do lewej krawędzi pasa.
  - zastępuje: `.we-lacznik`
  - warianty: `--blad`, `--drugorzedny`
  - stany: spoczynek, z odsyłaczem .dn-link, stan błędu (tekst w barwie błędu), wąski pas (tekst schodzi do drugiego wiersza)

**`.dn-bryla`** — ZROBIONE, patrz punkt 0c. Rodzina wyszła jednoskładnikowa: scena,
płyta i siatka nie są osobnymi składnikami CSS, bo bryłę rysuje płótno
(`zasoby/bryla.js`), a nie stos elementów. Pozycje `.dn-bryla-plyta`,
`.dn-bryla-scena` i `.dn-bryla-siatka` są bezprzedmiotowe i zostały skreślone;
klasy, które miały zastąpić (`.in-bryla`, `.in-plyta*`, `.in-bryla-pole`,
`.in-warstwa*`), znikły z `wejscie.css` wraz z wdrożeniem. Do przeniesienia
zostaje `.pg-bryla` i `.pg-warstwa*` z przepływu wejścia — to samo pole,
ten sam składnik.

**`.dn-etyk-mono`** — Etykieta w mono-wersalikach: nadtytuł, nazwa modułu, głowa kolumny, zapowiedź grupy kontrolek.
  - zastępuje: `.in-etykieta`, `.cd-modul-nazwa`, `.cd-nadtytul`, `.pd-nadtytul`
  - warianty: `--mocna (fs-sm, --dn-tekst)`, `--na-ramie`, `--z-kreska`
  - stany: spoczynek, na podłożu atramentowym, w rzędzie z czynnością po prawej, obcięcie długiego tekstu (ellipsis), forced-colors

**`.dn-godlo`** — Samodzielny znak marki albo modułu — nieinteraktywny svg w jednej skali z księgi znaku, w barwie z żetonu.
  - zastępuje: `.we-marka-godlo`, `.pg-znak`, `.cd-marka-lockup`, `.st-pasmo-znak`
  - warianty: `--lg (skala bryły)`, `--lockup (poziomy znak z logotypem)`, `--modul (skala ikonowa, barwa sygnału)`, `--tetno`
  - stany: spoczynek, z pulsującą kropką sygnału (--tetno), prefers-reduced-motion (puls zatrzymany), na podłożu atramentowym i na jasnym, forced-colors

**`.dn-kod-pole`** — Kratka jednego znaku kodu potwierdzającego — stała bryła, pismo maszynowe środkowane, własny stan skupienia.
  - zastępuje: `.au-kod-pole`
  - warianty: `--sm`, `--blad`
  - stany: pusta, wypełniona, fokus (obrys i cień sygnałowy), błąd (kod odrzucony), nieczynna (w trakcie sprawdzania), wklejenie całego kodu naraz, autouzupełnienie z wiadomości

**`.dn-kod-stopka`** — Wiersz pod kratkami kodu: adres wysyłki, odliczanie i ponowne wysłanie.
  - zastępuje: `.au-kod-stopka`
  - stany: odliczanie trwa (ponowienie nieczynne), ponowienie dostępne, wysłano ponownie (potwierdzenie), błąd wysyłki, wąski widok (zawinięcie do dwóch wierszy)

**`.dn-kryteria`** — Zawijany wykaz kryteriów do spełnienia — wymagania hasła, warunki wysyłki, lista kontrolna pola.
  - zastępuje: `.au-sila-warunki`
  - warianty: `--kolumna`, `--zwarte`
  - stany: spoczynek, wszystkie spełnione, część spełniona, po nieudanej próbie (podniesiona wyrazistość), ogłaszany asystująco przy zmianie

**`.dn-kryterium`** — Pojedyncze kryterium z ikoną i stanem spełnienia niesionym atrybutem.
  - zastępuje: `.au-sila-warunek`
  - stany: niespełnione, spełnione (data-spelniony), nieokreślone (jeszcze nie sprawdzane), forced-colors (znak, nie sama barwa — stan nie może zależeć wyłącznie od koloru)

**`.dn-kurtyna`** — Pełnopowierzchniowa warstwa ujęcia otwierającego okno, gasnąca przy wejściu w treść.
  - zastępuje: `.in-otwarcie`
  - warianty: `--atrament`, `--natychmiast`
  - stany: widoczna, gaśnie (data-stan='gasnie'), zgaszona (pointer-events: none, aria-hidden), prefers-reduced-motion (znika bez przejścia)

**`.dn-link`** — Odsyłacz tekstowy interfejsu — barwa sygnału i podkreślenie wymagane przez WCAG 1.4.1 tam, gdzie odsyłacz stoi w toku zdania.
  - zastępuje: `.au-link`
  - warianty: `--cichy (bez podkreślenia, wyłącznie poza tokiem zdania)`, `--na-ramie`, `--niebezpieczny`, `--zewnetrzny (znak wyjścia)`
  - stany: spoczynek, najechanie, fokus klawiaturą (:focus-visible, obrys --dn-fokus), wciśnięty (:active), odwiedzony, nieczynny (aria-disabled), w toku zdania vs samodzielny, forced-colors

**`.dn-meta`** — Drobna nota maszynowa drugorzędna — metryka pod tytułem, miara przy pozycji, nota na dnie kolumny. Jeden przepis pisma dla całej aplikacji.
  - zastępuje: `.in-wydanie`, `.in-wariant-miara`, `.in-suma-wolne`, `.in-podsumowanie`, `.we-marka-stopka`
  - warianty: `--na-ramie`, `--mocna`, `--prawa`
  - stany: spoczynek, na podłożu atramentowym (--na-ramie), wyróżnienie b/strong wewnątrz, stan błędu dziedziczony z rodzica (data-stan='blad'), wysoki kontrast / forced-colors

**`.dn-okno-wejsciowe`** — Bryła wolnostojącego okna sprzed uwierzytelnienia — stały wymiar, dwie kolumny, wspólny blok zmiennych całej rodziny (instalator, okno startowe, logowanie, przygotowanie).
  - zastępuje: `.we-okno`
  - warianty: `--jednokolumnowe`, `--waskie`, `--wysokie`
  - stany: spoczynek, fokus wewnątrz (pierwsza kontrolka), stan błędu przepływu, wąski widok — kolumna tożsamości chowana, prefers-reduced-motion przy wejściu

**`.dn-okno-wejsciowe-panel`** — Kolumna stanu albo formularza okna wejściowego: jasna powierzchnia, własny rytm bloków, przewijanie i pas działań opadający na dno.
  - zastępuje: `.we-panel`
  - warianty: `--scentrowany`, `--zwarty`
  - stany: spoczynek, z pasem działań (pas na dnie), bez pasa działań (stopka bez odstępu), treść dłuższa niż panel (przewijanie), stan pracy (kontrolki nieczynne), stan błędu

**`.dn-przewijane`** — Pole treści dłuższej niż jego okno, samo pokazujące zapas: gasnące krawędzie górna i dolna oraz własna szyna przewijania.
  - zastępuje: `.in-przewijane`
  - warianty: `--bez-szyny`, `--poziome`
  - stany: zapas u góry (data-poczatek), zapas u dołu (data-koniec), bez zapasu (obie krawędzie zgaszone), w trakcie przewijania (szyna widoczna), fokus wewnątrz, prefers-reduced-motion (bez przejścia krycia)

**`.dn-przewijane-suwak`** — Suwak w szynie pola przewijanego, prowadzony miarą przewinięcia.
  - zastępuje: `.in-przewijane-suwak`
  - stany: spoczynek, najechanie, ciągnięcie (:active), fokus klawiaturą, forced-colors

**`.dn-przewijane-szyna`** — Tor własnego paska przewijania pola — miejsce, po którym idzie suwak.
  - zastępuje: `.in-przewijane-szyna`
  - warianty: `--stala (zawsze widoczna)`
  - stany: spoczynek (przygaszona), najechanie na pole (rozjaśniona), w trakcie ciągnięcia, ukryta, gdy treść mieści się w całości

**`.dn-roznica`** — Wtręt wykazu różnic w toku zdania — fragment dodany, usunięty albo o zmienionym formatowaniu przy porównaniu wersji.
  - zastępuje: `.st-diff-dod`, `.st-diff-usu`, `.st-diff-fmt`
  - warianty: `--dodane`, `--usuniete`, `--format`
  - stany: dodane, usunięte, format zmieniony, wskazane (najechanie z podpowiedzią autora i czasu), wybrane do przywrócenia, forced-colors (kontur i przekreślenie muszą przeżyć zdjęcie barw)

**`.dn-scena`** — Podłoże pełnej wysokości, które środkuje wolnostojące okno sprzed uwierzytelnienia.
  - zastępuje: `.we-scena`
  - warianty: `--atrament`
  - stany: spoczynek, okno wyższe niż widok (przewijanie sceny), wąski widok (okno na całą szerokość), motyw jasny i ciemny

**`.dn-sekcja-tytul`** — Tytuł strefy na płótnie okna z kreską sygnałową niosącą wagę strefy w trzech stopniach krycia.
  - zastępuje: `.cd-sekcja-tytul`
  - warianty: `--bez-kreski`, `--z-czynnoscia`
  - stany: waga 1 / 2 / 3 (data-waga), z czynnością po prawej w wierszu, zwinięta / rozwinięta sekcja (aria-expanded), fokus na kontrolce zwijania, forced-colors (kreska nie może zniknąć)

**`.dn-sila`** — Odcinkowa miara jakości prowadzona stopniem (siła hasła, jakość wpisu) — odcinki, opis stopnia i barwa progu.
  - zastępuje: `.au-sila`
  - warianty: `--zwarta (bez opisu)`
  - stany: stopień 0 (pusto), 1 (słaba), 2 (dostateczna), 3 (dobra), 4 (mocna), nieczynna (pole zablokowane), zmiana stopnia ogłaszana asystująco (aria-live), prefers-reduced-motion

**`.dn-sila-odcinek`** — Pojedynczy odcinek miary odcinkowej.
  - zastępuje: `.au-sila-odc`
  - stany: nieosiągnięty, osiągnięty (barwa progu), prefers-reduced-motion (bez przejścia barwy), forced-colors

**`.dn-tekst-ciagly`** — Płótno tekstu ciągłego: miara wiersza, interlinia i skala nagłówków oraz akapitów wewnątrz okna.
  - zastępuje: `.in-licencja`, `.st-kanwa`, `.cd-nauka`
  - warianty: `--dokument (78ch, lh-luzny, skala h1/h2/p)`, `--osadzony (własna powierzchnia, obrys, przewijanie)`, `--drobny (fs-sm — samouczek, nota prawna)`
  - stany: spoczynek, ogniskowalny klawiaturą, gdy sam się przewija (tabindex=0, :focus-visible), przewinięty (współpraca z .dn-przewijane), zaznaczenie tekstu (::selection), tylko do odczytu vs redagowany (contenteditable), forced-colors

**`.dn-tor-krokow`** — Poziomy tor kroków procesu wskazujący miejsce w przepływie, czytany w jednym wierszu nad treścią.
  - zastępuje: `.au-kroki`
  - warianty: `--zwarty`, `--na-ramie`
  - stany: spoczynek, krok pierwszy / środkowy / ostatni, wąski widok (zwinięty do „krok 2 z 3”), ogłaszany asystująco przy przejściu, forced-colors

**`.dn-tor-krokow-poz`** — Pozycja poziomego toru kroków z numerem i stanem miejsca w przepływie; rozdzielacz rysowany przez ::after pozycji.
  - zastępuje: `.au-krok`, `.au-krok-strzalka`, `.au-krok-nr`
  - warianty: `--klikalna`
  - stany: przyszły, bieżący (data-biezacy), zrobiony (znak zamiast numeru), niedostępny, klikalny (gdy wolno wrócić): najechanie, fokus, wciśnięty, aria-current='step'

**`.dn-tozsamosc`** — Atramentowa kolumna tożsamości okna wejściowego — godło, nazwa, motto, wykaz cech i nota na dnie.
  - zastępuje: `.we-marka`
  - warianty: `--pas (pozioma na wąskim widoku)`
  - stany: spoczynek, wąski widok (zwinięta do pasa nad panelem), treść dłuższa niż kolumna, motyw jasny i ciemny (atrament stały)

**`.dn-tozsamosc-poz`** — Pozycja wykazu cech w kolumnie tożsamości: ikona, wytłuszczony lead w osobnym wierszu i zdanie opisu, na podłożu atramentowym.
  - zastępuje: `.we-marka-poz`
  - warianty: `--zwarta`
  - stany: spoczynek, bez ikony, lead bez opisu, zawijanie długiego zdania, forced-colors

**`.dn-tytul-okna`** — Tytuł okna (strony) w stopniu ekspozycyjnym — najwyższy poziom nagłówka na płótnie treści.
  - zastępuje: `.cd-tytul`, `.pd-tytul`
  - warianty: `--na-ramie`, `--sm (fs-3xl na wąskim oknie)`
  - stany: spoczynek, z nadtytułem .dn-etyk-mono nad sobą, z plakietką stanu w wierszu, zawijanie i obcięcie, na podłożu atramentowym

**`.dn-wersja`** — Pozycja repozytorium wersji: wielowierszowy wpis (wiersz stanu z metadanymi, nota zmiany, rząd czynności) odcięty kreską od następnego.
  - zastępuje: `.st-wersja`
  - warianty: `--biezaca`, `--wybrana`
  - stany: spoczynek, najechanie, fokus klawiaturą, bieżąca, wybrana do porównania, przywracana (praca w toku), ostatnia w wykazie (bez kreski)

**`.dn-zaznaczenie`** — Zaznaczony fragment tekstu w toku treści: podświetlenie sygnałowe z konturem, jednocześnie kontener odniesienia dla przybornika unoszonego nad nim.
  - zastępuje: `.st-zazn`
  - warianty: `--slad (zaznaczenie nieczynne, po odejściu kursora)`
  - stany: spoczynek, czynne (przybornik otwarty), zaznaczenie utracone (panel boczny zachowuje ślad), zaznaczenie wielowierszowe, forced-colors

**`.dn-zestawienie`** — Wiersz podsumowania liczbowego na własnej powierzchni: nazwa z lewej, miara z prawej.
  - zastępuje: `.in-suma`
  - warianty: `--mocne (suma końcowa)`, `--ostrzezenie`
  - stany: spoczynek, miara ostrzegawcza (za mało miejsca), miara błędna, stan pracy (miara jeszcze liczona), wiele wierszy w bloku

### 3b · Warianty istniejących składników

| modyfikator | rodzina | zastępuje |
|---|---|---|
| `--sukces` | `.dn-alert` | `.we-alarm--sukces`, `.we-alarm--blad` |
| `--cichy` | `.dn-btn-ikona` | `.cd-sekcja-menu`, `.cd-karta-menu`, `.cd-listwa-menu`, `.cd-wlasny-menu`, `prywatny mechanizm w .sta-menu-wiersz (komponenty.css:1392)` |
| `--zwarta` | `.dn-karta-srodowiska` | `.cd-karta-akcent` |
| `--w-krawedzi` | `.dn-karty-pasmo` | `.st-pasmo` |
| `--taca-ramy` | `.dn-narzedzia-pas` | `.st-obudowa` |
| `--plywajacy` | `.dn-przybornik` | `.st-plyw` |

### 3c · Brakujące części istniejących rodzin

| część | dołącza do | zastępuje |
|---|---|---|
| `.dn-alert-tresc` | `.dn-alert` | `.cd-start-tresc` |
| `.dn-alert-znak` | `.dn-alert` | `.cd-start-znak` |
| `.dn-kafel-tresc` | `.dn-kafel` | `.cd-kafel-tresc` |
| `.dn-karta-srodowiska-stopka` | `.dn-karta-srodowiska` | `.cd-karta-meta` |
| `.dn-obszar-pasek-tytul` | `.dn-obszar-pasek` | `.st-wstazka-sesja` |
| `.dn-postep-glowa` | `.dn-postep` | `.in-postep-wiersz` |

---

## 4 · Zostaje w oknie — czyste rozmieszczenie

```
  .we-okno--z-belka               .we-belka-uchwyt                .we-marka-lista                 .we-glowa
  .we-komunikaty                  .we-akcje                       .we-akcje-poboczne              .we-stopka--ciag
  .au-pola                        .au-para                        .au-haslo                       .au-sila-tor
  .au-opcje                       .au-akcje                       .au-metody                      .au-kod
  .pg-bryla-pole                  .pg-postep                      .in-korpus                      .in-otwarcie-znak
  .in-ekran                       .in-warianty                    .in-grupa                       .in-pola
  .in-wybory                      .in-przebieg                    .in-postep                      .in-prawa
  .in-zakonczenie                 .in-wynik                       .in-dzialania-odstep            .cd-tresc--modul
  .cd-naglowek                    .cd-strefa                      .cd-sekcja-glowa                .cd-siatka-srodowisk
  .cd-siatka-kafli                .cd-siatka-wlasnych             .cd-karta-godlo                 .cd-metryka-czlon
  .cd-wejdz                       .cd-kafel                       .cd-kafel-blok                  .cd-dodaj-pas
  .cd-wlasny                      .cd-wlasny-tresc                .cd-modul-glowa                 .cd-modul-odnosnik
  .st-pasmo-przywroc              .st-wstazka-odstep              .st-panel-wiersz                .st-pole-rozciagniete
  .st-edytor-odstep               .st-wersja-akcje                .st-wersja-stopka               .st-podglad
```