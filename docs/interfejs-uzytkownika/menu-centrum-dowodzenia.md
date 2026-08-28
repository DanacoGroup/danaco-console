# KOMPLET DLA DEWELOPERA — `design/05-okna/przeplyw/centrum-dowodzenia.html` (1510 wierszy)

Drugi przebieg weryfikacji. Wszystkie liczby zmierzone na repozytorium `/home/ubuntu/budowa` w dniu 28.08.2026. Gdzie pomiaru nie wykonano, napisano „nie zmierzono”.

Skróty źródeł używane dalej:
- **WZ** = `design/05-okna/WZORZEC-STANOWISKA.html`
- **33 okna** = pozostałe pliki `design/05-okna/**/*.html` niosące dany składnik wspólny (bez ocenianego)
- **MAN** = `design/zasoby/ikony/manifest.json`
- **KAT** = `design/01-dokumentacja-md/07-icons-widgets.md`
- **HAND** = `design/01-dokumentacja-md/08-handoff-motion-dostepnosc.md`
- **KONTRAKT** = `budowa/shared/contract.json`

---

## 1. SUMY CAŁOŚCI

### 1.1. Menu

| Miara | Wartość |
|---|---|
| Menu w pliku (`role="menu"` z `id`) | 40 |
| Układów niepowtarzalnych | 22 (w tym 1 pusty, wypełniany skryptem) |
| Pozycji łącznie (`sta-menu-poz[role^=menuitem]`) | 358 |
| Pozycji niepowtarzalnych | 237 |
| Nagłówków grup | 32 |
| Pól skrótu (`dn-menu-skrot-pole`) | 4 |

Powtórzenia (nie liczone drugi raz): `menu-sesja-1…4` + `menu-drzewo-1,2,4,5` — 8 menu identycznych po 8 pozycji; `menu-sesja-5…7` + `menu-drzewo-3,6,7` — 6 menu identycznych po 8 pozycji; `menu-projekt-1/-2` — 2 × 5; `menu-wlasny-1…6` — 6 × 4.

**Ocena 237 pozycji niepowtarzalnych**

| | menu-aplikacji (przebieg 1) | pozostałe 21 menu (przebieg 2) | razem |
|---|---|---|---|
| Pozycji ocenionych | 67 | 170 | **237** |
| Potwierdzone | 32 | 129 | **161** |
| Do poprawy nazwy/skrótu | 10 | 30 | **40** |
| Do usunięcia | 1 | 5 | **6** |
| Bez pokrycia w źródle | 13 | 6 | **19** |
| Nierozliczone | 11 (7 wyzwalaczy podmenu + 4 nie zmierzono) | 0 | **11** |
| Do dopisania (poza 237) | 6 | 11 | **17** |

Osobno, poza wykazem pozycji:
- nagłówki do poprawy — **6**; nagłówki do usunięcia — **2**;
- pola skrótu do poprawy — **4**; separatory do usunięcia — **2**;
- ikony pozycji do poprawy — **13** rozstrzygnięte źródłem, **5** bez rozstrzygnięcia (do pytań).

### 1.2. Szyna nawigacji (w. 41–378)

| Miara | Wartość |
|---|---|
| Pozycji szyny łącznie | **53** |
| — sekcja „Praca” | 3 (w. 151, 154, 157) |
| — środowiska | 4 (w. 164, 197, 230, 260) |
| — moduły środowisk | 34 (9 + 9 + 8 + 8) |
| — szybki wybór | 6 (w. 294–309) + „Dodaj skrót” (w. 312) |
| — stałe (Menu aplikacji, Ustawienia, Dostosuj pasek, Motyw, Profil) | 5 |

| Ocena | Liczba |
|---|---|
| Potwierdzone | **42** |
| Do poprawy (nazwa/opis) | **6** |
| Do poprawy (ikona) | **3** |
| Bez pokrycia w źródle | **2** |
| Do usunięcia | **0** |
| Do dopisania | **0** |

### 1.3. Pomiar ikon — korekta ustaleń przebiegu pierwszego

Zmierzono sygnaturami kształtów (`d`, `cx/cy/r`, `x/y/w/h`, `points`) wobec 152 plików `design/zasoby/ikony/svg/`:

| Miara | Wartość |
|---|---|
| SVG w tym oknie | 573 |
| — zgodne co do znaku z plikiem zestawu | 416 (73 %) |
| — kształty spoza zestawu | 157 w 64 odmianach |
| SVG we wszystkich 36 plikach `design/05-okna/` | 10 517, z tego 6 989 spoza zestawu (66 %) |

**Wniosek:** kształty spoza zestawu są cechą całego produktu, także WZORCA (menu `menu-ustawienia` wzorca ma 17 z 18 ikon spoza zestawu). Zapis przebiegu pierwszego „szyna: 2 ikony spoza zestawu” nie oddaje stanu — to okno jest pod tym względem **najlepsze z 36**. Dlatego niżej **nie proponuje się wymiany kształtów na pliki zestawu**; proponuje się wyłącznie zmianę tam, gdzie jeden glif niesie w tym oknie dwie różne role, a źródło (MAN/KAT) nazywa rolę właściwą.

---

## 2. WYKAZ ZMIAN DO WYKONANIA

Zapis: **w.** = wiersz pliku. Przy pozycjach menu ikona stoi w tym samym wierszu co pozycja. Przy pozycjach szyny przycisk stoi w wierszu *N*, ikona w wierszu *N+1*.

### 2.1. `menu-ustawienia` (w. 322–353) — 22 pozycje, 5 grup; wzorzec ma 18 pozycji i 4 grupy

| w. | czynność | treść do wpisania | ikona | skrót | źródło |
|---|---|---|---|---|---|
| 324 | popraw | `Wszystkie ustawienia` → **`Ustawienia aplikacji`** | bez zmian | `Ctrl+,` bez zmian | WZ:219; brzmienie w 33 oknach, obecne tylko tu |
| 325 | popraw | `Język i formaty` → **`Język i format zapisu`** | bez zmian | — | WZ:219; 33 okna |
| 326 | popraw | `Wygląd` → **`Motyw i gęstość widoku`** | bez zmian | — | WZ:219; 33 okna; sekcja docelowa `platformowe/ustawienia.html:588` |
| 334 | popraw | `Modele i kanały` → **`Modele, kanały i poziom wysiłku`** | bez zmian | — | WZ:219; 33 okna |
| 335 | popraw | `Magistrala kontekstu` → **`Magistrala kontekstu i pamięć`** | bez zmian | — | WZ:219; 33 okna. Bez członu „i pamięć” pozycja zlewa się z narzędziem w. 109 i przyciskiem paska w. 474 |
| 345 | popraw | `Kopie zapasowe` → **`Kopie zapasowe i przywracanie`** | bez zmian | — | WZ:219; 33 okna |
| 346 | popraw | `Pamięć podręczna` → **`Pamięć podręczna i dane lokalne`** | bez zmian | — | WZ:219; 33 okna |
| 348 | **usuń** | separator `<div class="sta-menu-sep">` | — | — | konsekwencja usunięcia grupy niżej |
| 349 | **usuń** | nagłówek `Szybkie przełączniki` | — | — | grupy nie ma we WZ:219 ani w żadnym z 33 okien |
| 350 | **usuń** | pozycja `Motyw` | — | — | powielenie: przycisk szyny w. 358, pozycja w. 87, nastawa w pozycji w. 326 |
| 351 | **usuń** | pozycja `Gęstość` | — | — | powielenie: `Dostosuj paski` › Wygląd i zachowanie, w. 1452–1454; sekcja Ustawień `platformowe/ustawienia.html:588` |
| 352 | **usuń** | pozycja `Nie przeszkadzać` | — | — | brak pokrycia — patrz pytanie P-1 |
| 353 | **usuń** | pozycja `Always On Display` | — | `Ctrl+Shift+A` znika | powielenie czterokrotne: w. 309, 502, 515, 1236 |

Bez zmian (11): 327, 328, 329, 332, 333, 336, 337, 340, 341, 342, 347.
Ikona `diagnostyka` w w. 347 **zostaje** — po usunięciu w. 350–353 nie dzieli już glifu z niczym w tym menu (rola wg KAT:264 właściwa).

### 2.2. `menu-profil-szyna` (w. 366–375) — 6 pozycji; wzorzec ma 10

| w. | czynność | treść do wpisania | ikona | skrót | źródło |
|---|---|---|---|---|---|
| 370 | popraw | `Licencja i plan` → **`Plan i rozliczenia`** | bez zmian | — | WZ:232; 33 okna |
| 371 | popraw | `Urządzenia połączone` → **`Urządzenia i sesje logowania`** | bez zmian | — | WZ:232; 33 okna |
| 372 | popraw | `Bezpieczeństwo i logowanie` → **`Bezpieczeństwo i uwierzytelnianie dwuskładnikowe`** | bez zmian | — | WZ:232; 33 okna |
| po 372 | **dopisz** | **`Powiadomienia konta`** | `dzwonek` | — | WZ:232; 33 okna; MAN „Powiadomienia” |
| po 372, za separatorem | **dopisz** | **`Status: dostępny`** (`role="menuitemradio"`) | `konto` | — | WZ:232; 33 okna |
| jw. | **dopisz** | **`Status: zajęty`** | `minus-kolo` | — | WZ:232; MAN: `minus-kolo` = „Stan zajętości — dostępność wstrzymana” |
| jw. | **dopisz** | **`Status: niewidoczny`** | `oko-przekreslone` | — | WZ:232; MAN: „Wartość ukryta, podgląd wyłączony” |

Bez zmian: 369, 374, 375.

### 2.3. `menu-zrzut` (w. 412–434) — 14 pozycji + 2 pola skrótu

| w. | czynność | treść do wpisania | źródło |
|---|---|---|---|
| 417 | popraw | `Widok przewijany` → **`Widok przewijany w całości`** | WZ:280; 33 okna |
| 418 | popraw | `Z opóźnieniem` → **`Zrzut z opóźnieniem 5 s`** | WZ:280; 33 okna |
| 422 | popraw | `Zapisz do pliku` → **`Do pliku w katalogu zrzutów`** | WZ:280; 33 okna; `prowadzenie/decyzje.md:844` |
| 423 | popraw | `Zapisz w module Library` → **`Do modułu Library`** | WZ:280; 33 okna; `decyzje.md:844` |
| 424 | popraw | `Dołącz do bieżącej sesji` → **`Do Chat Window jako załącznik`** | WZ:280; 33 okna; `decyzje.md:844–850` (rozróżnienie załącznika od przytoczenia) |
| 427 | popraw | `Edytor adnotacji` → **`Otwórz w edytorze adnotacji po wykonaniu`** | WZ:280; 33 okna |
| 428 | popraw | `Ukrywaj dane wrażliwe` → **`Ukryj dane wrażliwe automatycznie`** | WZ:280; 33 okna |
| 431 | popraw | `Katalog i nazwa…` → **`Katalog docelowy i wzorzec nazwy…`** | WZ:280; 33 okna |
| 432 | popraw | etykieta `Zrzut obszaru` → **`Skrót zrzutu obszaru`**; `aria-label` → **`Skrót zrzutu obszaru — naciśnij kombinację, aby zmienić`**; wartość `Ctrl + PrtSc` bez zmian | WZ:280; 33 okna |
| 433 | popraw | etykieta `Zrzut okna` → **`Skrót zrzutu okna`**; `aria-label` → **`Skrót zrzutu okna — naciśnij kombinację, aby zmienić`**; wartość `Alt + PrtSc` bez zmian | WZ:280; 33 okna |

Bez zmian: 414, 415, 416, 421, 434.
**Bez pokrycia (nie usuwać bez rozstrzygnięcia):** w. 429 `Pokazuj kursor`, w. 430 `Format pliku` — po 1 wystąpieniu w produkcie, brak we WZ:280 i w 33 oknach. Patrz P-2.

### 2.4. `menu-schowek` (w. 439–455) — 10 pozycji + 2 pola skrótu

| w. | czynność | treść do wpisania | źródło |
|---|---|---|---|
| 443 | popraw | `Wklej bez formatu` → **`Wklej jako zwykły tekst`** | WZ:284; 33 okna |
| 449 | popraw | `Pomijaj hasła` → **`Pomijaj treść z pól haseł`** | WZ:284; 33 okna |
| 450 | popraw | `Współdziel między urządzeniami` → **`Współdziel schowek między urządzeniami`** | WZ:284; 33 okna |
| 451 | popraw | `Czyść przy zamknięciu` → **`Czyść schowek przy zamknięciu aplikacji`** | WZ:284; 33 okna |
| 452 | popraw | `Zakres pamiętania` → **`Liczba przechowywanych pozycji: 25`** | WZ:284; 33 okna |
| po 446 | **dopisz** | separator + nagłówek **`Ostatnie pozycje`** + 3 pozycje: **`„Sprawozdanie końcowe za trzeci kwartał” — fragment`**, **`raport-koncowy.md — ścieżka pliku`**, **`Zrzut obszaru — 12:04`** | WZ:284; 33 okna |
| 453 | popraw | `aria-label` → **`Otwarcie schowka — naciśnij kombinację, aby zmienić`** | WZ:284; 33 okna |
| 454 | popraw | etykieta `Wklejanie` → **`Wklejenie z listy`**; `aria-label` → **`Wklejenie z listy — naciśnij kombinację, aby zmienić`** | WZ:284; 33 okna |

Bez zmian: 441, 442, 444, 446, 455.
**Bez pokrycia:** w. 445 `Przypnij pozycję` — jedno wystąpienie w produkcie, brak we WZ:284. Patrz P-2.

### 2.5. `menu-wiecej` (w. 487–504) — 12 pozycji

| w. | czynność | treść do wpisania | źródło |
|---|---|---|---|
| 488 | popraw | nagłówek `Znajdź` → **`Znajdź w widoku`** | WZ:312; 33 okna |
| 490 | popraw | `Okno wiodące` → **`Okno robocze wiodące`** | WZ:312; 33 okna (jako pozycja menu skrócone brzmienie występuje tylko tu) |
| 491 | popraw | `Wszystkie okna karty` → **`Wszystkie okna bieżącej karty`** | WZ:312; 33 okna |

Bez zmian: 492–495, 498–503. Pozycje 493, 494, 501, 502, 503 są nadmiarem wobec wzorca **z pokryciem** — akapit w. 504 stanowi, że skład rozwinięcia zmienia okno „Dostosuj”.

### 2.6. `menu-aod` (w. 513–533) — 11 pozycji, 4 grupy

| w. | czynność | treść do wpisania | źródło |
|---|---|---|---|
| 515 | popraw | `Włącz Always On Display` → **`Powierzchnia interakcji`**; skrót `Ctrl+Shift+A` zostaje | WZ:319; 33 okna. **Zastrzeżenie:** wzorzec nie nadaje tej pozycji roli `menuitemcheckbox`; jeżeli ma zostać przełącznikiem, patrz P-3 |
| 524 | **usuń** | separator stojący między nagłówkiem `Wejścia` (w. 523) a pierwszą pozycją grupy (w. 525) | usterka struktury: separator nie może rozdzielać nagłówka od jego pierwszej pozycji — w tym samym pliku wszystkie 31 pozostałych grup są zbudowane odwrotnie |

Bez zmian: 516, 517, 520, 521, 525, 526, 528, 529, 532, 533.
Zgłoszenie porządkowe (bez zmiany treści): pozycje 528–529 („Wyciszenie sugestii”, „Reguły wyzwalania”) stoją pod nagłówkiem `Wejścia`, choć wejściami nie są. Rozstrzygnięcie wymaga źródła — P-3.

### 2.7. `menu-sesje-filtr` (w. 548–560) — 9 pozycji, 2 grupy

| w. | czynność | treść / ikona | źródło |
|---|---|---|---|
| 549 | popraw | nagłówek `Pokaż` → **`Stan`** | 20 plików / 21 wystąpień; `moduly/research.html:369`, `moduly/workspace.html:374`. „Pokaż” — 1 plik, 2 wystąpienia (549, 564) |
| 555 | popraw | nagłówek `Porządek` → **`Sortowanie`** | 20 plików; `research.html:373`, `workspace.html:378`. „Porządek” — 1 plik, 2 wystąpienia (555, 569) |
| 556 | popraw | `Ostatnio czynne` → **`Ostatnia aktywność`**; zmiana obejmuje `panel-sesji.js:41` (`czynnosc: 'ostatnio czynne'`) | 16 plików; `research.html:374`, `workspace.html:379` |
| 553 | popraw ikonę | `ptaszek` → **`ptaszek-kolo`** | KAT:180 i :224; MAN: `ptaszek` = zaznaczenie, `ptaszek-kolo` = „proces zakończony poprawnie”. Pozycja ma już znacznik wyboru `<span class="ptaszek">✓</span>` — dziś w jednym wierszu stoją dwa te same znaki |
| 558 | popraw ikonę | `wielkosc-liter` → **`sortowanie`** | MAN: `sortowanie` = „Zmiana porządku wykazu”; `wielkosc-liter` = „Rozróżnianie wielkości liter w wyszukiwaniu” — ten sam glif niesie w tym oknie pozycje w. 75 i 495 |
| 559 | popraw ikonę | `wielkosc-liter` → **`sortowanie`** | jw. |

Bez zmian: 550, 551, 552, 557, 560 (ikona w. 560 — patrz P-4).

### 2.8. `menu-projekty-filtr` (w. 563–571) — 5 pozycji, 2 grupy

| w. | czynność | treść / ikona | źródło |
|---|---|---|---|
| 564 | popraw | nagłówek `Pokaż` → **`Stan`** | jak w. 549; **oba menu stoją w jednym pasku panelu i muszą nieść jeden nagłówek** |
| 569 | popraw | nagłówek `Porządek` → **`Sortowanie`** | jak w. 555 |
| 566 | popraw | `Z pracą` → **`W toku`** | filtr działa na `data-stan="praca"` (`panel-sesji.css:110`) — dokładnie ten sam stan, który w sąsiednim menu tego samego paska nazywa się „W toku” (w. 551). Nazwa produktowa: `research.html:371` |
| 571 | popraw | `Ostatnio czynne` → **`Ostatnia aktywność`** | jak w. 556 |
| 570 | popraw ikonę | `wielkosc-liter` → **`sortowanie`** | jak w. 558 |

Bez zmian: 565, 567.

### 2.9. `menu-panel` (w. 576–584) — 6 pozycji

| w. | czynność | uzasadnienie i źródło |
|---|---|---|
| 578 | **usuń** pozycję `Archiwum sesji` | powielenie: przycisk ikonowy „Archiwum sesji” stoi w tym samym pasku panelu (w. 545) z tym samym uchwytem `data-otwarz-historie` i tą samą ikoną `archiwum`. Oba widoczne zawsze — `cd-wstazka-grupa` nie zwęża się responsywnie (`zasoby/okna/centrum-dowodzenia.css:95`) |
| 584 | zostaw, **zgłoś** | `Eksportuj wykaz…` — czynność zaimplementowana (`panel-sesji.js:263–267`), ale KONTRAKT nie ma komendy wydającej wykaz sesji ani projektów: rodzina `session.*` to `create, list, focus, bind, open, close, delete, rename, copy, project.set, project.clear, archive, restore, archive.list, resume, stop, tool.*`; wszystkie 32 komendy `*.export` dotyczą raportów, dokumentów, zasobów i glosariuszy. Wymaga wpisu w kontrakcie, nie usunięcia |

Bez zmian: 577, 580, 581, 583. Skrótu przy w. 583 nie nadawać — `Ctrl+B` jest w tym oknie zajęty (w. 81, 874).

### 2.10. Menu wierszy sesji — 14 menu

Warianty: „praca” — `menu-sesja-1…4` (w. 590, 603, 616, 629) i `menu-drzewo-1, 2, 4, 5` (w. 695, 708, 746, 759); „zakończona” — `menu-sesja-5…7` (w. 642, 655, 668) i `menu-drzewo-3, 6, 7` (w. 721, 772, 790).

| czynność | treść | ikona | miejsce | źródło |
|---|---|---|---|---|
| **dopisz** we **wszystkich 14 menu** | **`Wyjmij z projektu`** | `minus` (MAN: „Minimalizacja okna, **odjęcie pozycji z wykazu**”) — patrz P-5 | bezpośrednio po `Przenieś do projektu`, czyli po w. 593, 606, 619, 632, 645, 658, 671, 698, 711, 724, 749, 762, 775, 792 | KONTRAKT → `session.project.clear`: „Wyjmuje sesje z projektu; sesja zostaje w historii bez projektu”, pole `sessionIds string[] wymagane`. W pliku „Wyjmij” — 0 wystąpień, „z projektu” — 0 |

**Nie dopisywać ikony `minus-kolo`** — MAN opisuje ją jako znak *stanu* („Stan zajętości — dostępność wstrzymana”), nie czynności.

**Usterka prototypu do naprawy przy okazji (nie zmiana treści):** cztery z ośmiu pozycji każdego menu wiersza nie mają uchwytu `data-poz-akcja` i są martwe — `Kopiuj sesję` (592 i 13 dalszych), `Katalog roboczy` (595 …), `Kopiuj ścieżkę` (596 …), `Zatrzymaj pracę` / `Wznów sesję` (598 …, 650 …). Obsłużone są tylko `nazwa`, `przenies`, `archiwizuj`, `usun` (`panel-sesji.js:19, 184–198`). Wszystkie cztery mają pokrycie w kontrakcie (`session.copy`, `SessionConfigArea.workingDirectory`, `session.stop`, `session.resume`).

### 2.11. `menu-projekt-1` / `menu-projekt-2` (w. 686, 737) — 5 pozycji

Bez zmian treści. **Zgłoszenia bez pokrycia w kontrakcie:**
- w. 687 / 738 `Zmień nazwę` — brak komendy zmiany nazwy projektu (`workspace.project.*` zawiera wyłącznie `status.set`; `session.rename` dotyczy sesji).
- w. 692 / 743 `Usuń trwale` — brak komendy usunięcia projektu.
- w. 691 / 742 `Archiwizuj` — **pokryte**: `workspace.project.status.set` + `WorkspaceProjectStatus.archived`.

### 2.12. `menu-karty-otwarte` (w. 811–824) — 9 pozycji, 2 grupy

| w. | czynność | treść / ikona | źródło |
|---|---|---|---|
| 812 | popraw | nagłówek `Otwarte karty (3)` → **`Otwarte karty (4)`** | pod nagłówkiem stoją cztery pozycje (813–816); pasmo `role="tablist"` (w. 843) niesie trzy zakładki, bo trzecia jest grupą dwóch kart — `menu-grupa-umowy` w. 888–893 („Klasyfikacja umów”, „Migracja rejestru umów”) |
| 821 | popraw ikonę | `dom` → **`zamknij`** | MAN: `zamknij` = „Zamknięcie karty, modala, powiadomienia”; `dom` = „Centrum dowodzenia — powrót do strony głównej”. Glif `dom` niosą dziś w tym oknie pozycje 813, 821, 822, 823, 824 |
| 823 | popraw ikonę | `dom` → **`pinezka`** | MAN: `pinezka` = „Przypięcie wniosku, źródła, **karty sesji**” |
| po 822 | **dopisz** | **`Zamknij karty po prawej`**, ikona `zamknij` | `docs/srodowiska/workspace.md:361` — „Menu karty \| Zamknij inne, **zamknij po prawej**, przypnij, dodaj do grupy, zmień nazwę, duplikuj”. W pliku „po prawej” pada raz (w. 583, w innej roli). **Zastrzeżenie:** źródło przypisuje zestaw menu kontekstowemu karty, którego to okno nie ma — patrz P-6 |
| po dopisanej pozycji | **dopisz** | **`Przekaż na Mobile`**, ikona `telefon` | `docs/srodowiska/talkin.md:459` („Menu karty \| … przekaż na Mobile …”) i `:759` („Akcja «Przekaż na Mobile» oznacza kartę do monitorowania z urządzenia mobilnego”). MAN: `telefon` = „Widok mobilny — pozycja «Mobile»”; KAT:266. W pliku „Przekaż na Mobile” — 0 wystąpień; kafel „Mobile” stoi w w. 1231 |
| na końcu grupy | **dopisz** | **`Zapisz zestaw kart…`**, ikona `zapisz` | `docs/srodowiska/codestudio.md:419` — „Zestaw kart jako obszar pracy \| Zapisany zestaw otwartych kart wywoływany jednym poleceniem \| … \| Paleta poleceń; **menu ☰ paska kart**”; nastawa `shell.session.saved` w `docs/srodowiska/workspace.md:612`. `menu-karty-otwarte` jest właśnie menu ☰ pasma kart (wyzwalacz w. 810, `aria-label="Menu okna roboczego"`). W pliku nie zapisuje zestawu ani „Nowa karta z szablonu…” (w. 864), ani „Nowy z szablonu” (w. 53), ani „Zapisz sesję” (w. 55) |

Bez zmian: 813–816, 818, 822, 824 (ikony 822 i 824 — patrz P-4).

### 2.13. `menu-powloki` (w. 826–841) — 10 pozycji, 2 grupy

| w. | czynność | treść / ikona | źródło |
|---|---|---|---|
| 828 | popraw | `<span class="skrot">kart: 3</span>` → **`kart: 4`** | Okno 1 jest oknem bieżącym (`aria-current="true"`); jego pasmo niesie 4 karty (813–816) |
| 829 | popraw ikonę | `oko-przekreslone` → **`karta-okna`** | MAN: `oko-przekreslone` = „Wartość ukryta, podgląd wyłączony” — rola niezwiązana z oknem roboczym; `karta-okna` = „Okno, karta przeglądarki”, tę ikonę niesie już w. 828 |
| 830 | popraw ikonę | `rozwidlenie` → **`karta-okna`** | MAN: `rozwidlenie` = „Rozwidlenie gałęzi repozytorium” |
| 840 | popraw ikonę | `okno-plywajace` → **`aktywnosc`** | MAN: `okno-plywajace` = „Always On Display — awatar pływający…”; `aktywnosc` = „**Praca w tle**, telemetria procesu”. Pozycja stoi pod nagłówkiem „Sesje w tle” |
| 841 | popraw ikonę | `okno-plywajace` → **`aktywnosc`** | jw. |

Bez zmian: 832–837 (ikony 832–835 — patrz P-4).
**Sprzeczność do rozstrzygnięcia:** pozycje 840 i 841 nazywają te same dwie karty, które stoją w grupie „Rejestr umów” (w. 846, 891, 893) jako karty otwarte. Patrz P-7.

### 2.14. `menu-nowa-karta` (w. 849–865) — 12 pozycji, 2 nagłówki

| w. | czynność | treść / ikona | źródło |
|---|---|---|---|
| 850 | **usuń** albo przenieś | nagłówek `Moduły środowisk` nie ma pod sobą żadnej pozycji — bezpośrednio po nim stoi drugi nagłówek `Najczęstsze` (w. 851), a moduły (855–862) stoją niżej, za separatorem 854, bez nagłówka. Zalecenie: **przenieść nagłówek `Moduły środowisk` z w. 850 na miejsce między w. 854 a 855** | usterka struktury; w pozostałych 31 grupach tego pliku po nagłówku zawsze stoi pozycja |
| 852 | popraw ikonę | `dokument` → **`rozmowa-nowa`** | MAN: `rozmowa-nowa` = „**Otwarcie nowej sesji rozmowy**”. Glif `dokument` niosą dziś w tym oknie pozycje 814, 815, 852, 853, 855, 864, 865, 891 |
| 853 | popraw ikonę | `dokument` → **`plus`** | MAN: `plus` = „**Utworzenie bytu: sesja, okno, kanał**” |
| 865 | popraw ikonę | `dokument` → **`kopiuj`** | MAN: `kopiuj` = „Kopiowanie treści, identyfikatora”; KAT:222. Ten sam glif niosą już pozycje „Duplikuj” kafli (1119, 1136, 1153, 1170, 1187, 1204) |

Bez zmian: 855–862, 864 (ikony modułów 855–862 zgodne z KAT:328–341).
**Nie zmieniać nazw** w. 852 i 865 — brzmienia są właściwe (obalone w kontroli jako wpisy źle zakwalifikowane).

### 2.15. Szyna nawigacji

| w. | czynność | treść do wpisania | źródło |
|---|---|---|---|
| 152 | popraw opis | `Otwarcie sesji. Środowisko wskazujesz w następnym kroku.` → **`Otwarcie nowej sesji — najpierw wskazanie środowiska, w którym ma powstać.`** | WZ:42; 33 okna |
| 158 | popraw opis | `Założenie projektu. Wskazanie środowiska otwiera jego Workspace.` → **`Założenie projektu — wskazanie środowiska otwiera od razu jego Workspace.`** | WZ:51; 33 okna |
| 172 → ikona w 173 | popraw ikonę | `ksiazka` → **`biblioteka`** | KAT:330 („Library \| `biblioteka`”); MAN: `biblioteka` = „Moduł Library — repozytorium wiedzy”, `ksiazka` = „Instrukcja użytkowania, słownik pojęć”. Pomiar: `biblioteka` przy Library — 66 wystąpień w 33 oknach; `ksiazka` — 2 wystąpienia, oba w tym pliku |
| 217 → ikona w 218 | popraw ikonę | `ksiazka` → **`biblioteka`** | jw. |
| 235 → ikona w 236 | popraw ikonę | `plik-kodu` → **`kod`** | KAT:339 („Developer \| `kod`”); MAN: `kod` = „Fragment kodu, moduł Developer”, `plik-kodu` = „Plik źródłowy w drzewie projektu”. Pomiar: `kod` przy Developer — 33 okna; `plik-kodu` — 1, ten plik. **W tym samym oknie** pozycja „Developer” w `menu-nowa-karta` (w. 860) niesie już `kod` |
| 277 | popraw etykietę | `data-etykietka="Harmonogram"` i `aria-label="Harmonogram"` → **`Harmonogram i automatyki`** | `docs/srodowiska/multitaskingai.md:650` (sekcja 6.2) i `:129`; 33 okna niosą `<b>Harmonogram i automatyki</b>` — to jedyny plik z formą skróconą. `data-modul` już brzmi „Harmonogram i automatyki” |
| 278 | popraw nazwę | `<b>Harmonogram</b>` → **`<b>Harmonogram i automatyki</b>`** | jw. |
| 313 | popraw opis | `Dodanie komponentu do strefy szybkiego wyboru.` → **`Dodanie pozycji szybkiego wyboru — ulubionego albo najczęściej używanego komponentu.`** | WZ:210; 33 okna |
| 319 | popraw etykietę | `data-etykietka="Ustawienia"` i `aria-label="Ustawienia"` → **`Ustawienia i konfiguracja`** | WZ:217; 33 okna |
| 320 | popraw nazwę | `<b>Ustawienia</b>` → **`<b>Ustawienia i konfiguracja</b>`** | jw. Opis „Ustawienia aplikacji, środowiska, sesji i modeli.” bez zmian — zgodny |
| 358 | popraw etykietę | `data-etykietka="Motyw"` i `aria-label="Motyw"` → **`Motyw jasny i ciemny`** | WZ:225; 33 okna |
| 359 | popraw nazwę | `<b>Motyw</b>` → **`<b>Motyw jasny i ciemny</b>`** | jw. |

**Bez pokrycia (nie zmieniać w tym pliku):** w. 283 `Workspace` i w. 286 `Agents` w rozwinięciu MultitaskingAI. `docs/srodowiska/multitaskingai.md:636–653` stanowi, że panel orkiestracji ma **sześć** sekcji o nazwach Zespoły · Role · Kolejki · Orkiestracja · Harmonogram i automatyki · Monitor procesu, a boczna nawigacja tego środowiska „nie zawiera listy modułów”. Obie nadmiarowe pozycje stoją identycznie we wszystkich 34 oknach — patrz P-8.

**Opis w. 262 „sześć sekcji sterowania” jest prawidłowy** — zgodny z `multitaskingai.md:68` i `:95`. Zapis przebiegu pierwszego („6 sekcji w opisie wobec 8 w rozwinięciu”) wskazywał właściwą rozbieżność, ale po niewłaściwej stronie: poprawić należy skład, nie opis.

### 2.16. Okno „Dostosuj paski” (w. 1319–1470)

| w. | czynność | treść | źródło |
|---|---|---|---|
| 1357 | popraw | `Zwiń szynę sesji` → **`Zwiń panel`** | wykaz ma nazywać składniki tak, jak brzmią na pasku; przycisk faktyczny to w. 396 `data-etykietka="Zwiń panel"` |
| po 1374 | **dopisz** | pozycja **`Samouczek`**, rodzina `widok`, ikona jak w. 509 | przycisk „Samouczek” stoi na pasku (w. 509), a nie występuje w żadnej z trzech kolumn (w. 1336–1348, 1356–1375, 1385–1395); nota o składzie stałym (w. 1404) wymienia wyłącznie cztery karty środowisk, menu aplikacji, ustawienia i profil |

Zgłoszenia bez zmiany (rozbieżności nazw wykazu wobec pasków, do jednej decyzji): w. 1339 `Eksportuj sesję` ↔ w. 56 `Eksportuj sesję…`; w. 1342/1343 `Powiększ widok` / `Pomniejsz widok` ↔ w. 89/90 `Powiększ` / `Pomniejsz`; w. 1365 `Pasmo kart sesji` ↔ w. 82 `Pasmo kart` ↔ WZ:34 `Pas kart sesji`. Przycisk „Always On Display” (w. 512) także nie występuje w wykazie ani w nocie składu stałego.

---

## 3. KOLIZJE SKRÓTÓW

Zmierzono wszystkie 51 kombinacji zapisanych w `<span class="skrot">`. Trzynaście kombinacji pada więcej niż raz. **Siedem** z nich to jedna czynność wywoływana z dwóch miejsc — regułą `HAND:1070` („Skrót nie może być jedyną drogą do funkcji”) to nie jest kolizja. **Sześć** to kolizje rzeczywiste.

### 3.1. Kolizje rzeczywiste — propozycja rozstrzygnięcia

| Skrót | Czynności (w.) | Ustępuje | Uzasadnienie i źródło |
|---|---|---|---|
| **Ctrl+Shift+T** | w. 87 `Motyw jasny/ciemny` · w. 443 `Wklej bez formatu` · w. 824 `Przywróć ostatnią kartę` | **w. 87 i w. 824 — skrót zdjąć bez zastąpienia** | WZ:284 nadaje `Ctrl+Shift+T` pozycji „Wklej jako zwykły tekst” — to jedyne z trzech użyć poparte źródłem. WZ:34 (podmenu Widok) nie ma pozycji motywu w ogóle, a jej skrótu nie nadaje żadne źródło. Dla „Przywróć ostatnią kartę” źródła nie ma; `docs/srodowiska/workspace.md:366` przypisuje przywracanie „Menu kebab (⋮) paska kart”, bez skrótu. Zostaje **w. 443**. |
| **Ctrl+Shift+N** | w. 51 `Nowe okno aplikacji` · w. 853 `Nowa sesja` | **w. 853** → **`Ctrl+N`** | WZ:34 nadaje `Ctrl+Shift+N` pozycji „Nowe okno aplikacji”, a `Ctrl+N` pozycji „Nowa sesja”. W tym pliku `Ctrl+N` niesie tylko w. 50 — tę samą czynność. |
| **Ctrl+B** | w. 81 `Okno pętli wykonawczej` · w. 874 `Panel boczny` | **w. 81 — skrót zdjąć** | WZ:34 nadaje `Ctrl+B` pozycji „Szyna sesji”, czyli panelowi bocznemu sesji — bliżej w. 874. Nazwy „Okno pętli wykonawczej” wzorzec nie zna. Zmiana dotyczy `menu-aplikacji`, ocenionego w przebiegu pierwszym — wykonać po jego rozstrzygnięciu. |
| **Ctrl+1** | w. 813 `Centrum dowodzenia` · w. 869 `Widok pojedynczy` | **w. 869 — skrót zdjąć bez zastąpienia** | Rozstrzyga nastawa w tym samym pliku: „Dostosuj paski” › Karty sesji › **`Numer karty — Numer przy podpisie; przejście do karty skrótem Ctrl + cyfra`** (w. 1434). Przestrzeń `Ctrl + cyfra` należy do przełączania kart (w. 813–816: Ctrl+1…Ctrl+4). Źródło nie nadaje układom okna żadnego skrótu — zastąpienia nie proponuje się. |
| **Ctrl+V** | w. 71 `Wklej` · w. 442 `Wklej ostatnią pozycję` | **żadna — nie pogłębiać** | Kolizja jest we wzorcu: WZ:34 nadaje `Ctrl+V` pozycji „Wklej”, WZ:284 tej samej kombinacji pozycji „Wklej ostatnią pozycję”. Rozstrzygnięcie należy do składnika wspólnego — P-9. |
| **Ctrl+W** | w. 59 `Zamknij kartę sesji` · w. 821 `Zamknij kartę` | **żadna — to jedna czynność pod dwiema nazwami** | Do ujednolicenia zostaje brzmienie, nie skrót. WZ:34 nadaje `Ctrl+W` pozycji **„Zamknij okno”**, a `HAND:1048` tej samej kombinacji „Zamknięcie bieżącej karty”. Źródła są sprzeczne — P-9. |

### 3.2. Powtórzenia bez kolizji (jedna czynność, dwa wejścia) — zamknięte, nic nie zmieniać

| Skrót | Wiersze | Czynność |
|---|---|---|
| Ctrl+Shift+A | 353, 502, 515 | Always On Display (po usunięciu w. 353 zostają dwa) |
| Ctrl+, | 115, 324 | Ustawienia aplikacji |
| Ctrl+/ | 138, 499 | Skróty klawiszowe |
| Ctrl+H | 75, 495 | Znajdź i zamień |
| F1 | 137, 498 | Dokumentacja platformy |
| F5 | 92, 904 | Odśwież widok |
| Ctrl+K | 102, 501 | Paleta poleceń |

### 3.3. Nakładanie skrótów powłoki na skróty modułów — do rozstrzygnięcia, nie do zmiany

`HAND:1071` stanowi: „Kolizja … rozstrzygana **kontekstem okna**; okna nie współdzielą przestrzeni skrótów”. Reguła nie obejmuje przypadku, w którym oknem modułu jest karta wewnątrz tej samej powłoki. Zmierzone nakładania:

| Skrót | W tym oknie | W module (HAND, rozdz. 18.2) |
|---|---|---|
| Ctrl+T | w. 852 `Nowe okno komunikacji` | :1047 „Nowa karta powłoki / nowa karta przeglądania — Terminal Tabs · Browser Window” |
| Ctrl+D | w. 865 `Duplikuj bieżącą kartę` | :1045 „Zaznacz kolejne wystąpienie — Code Editor”; :1052 „Dodanie strony do zakładek — Browser Window” |
| Ctrl+Shift+S | w. 864 `Nowa karta z szablonu…` | :1053 „Zrzut ekranu — Browser Window” |
| Ctrl+Shift+A | w. 502, 515 `Always On Display` | :1054 „Uruchomienie pełnej analizy — Diagnostics Center” |
| Ctrl+Shift+F | w. 85 `Tryb skupienia` | :1040 „Grep — znajdź w plikach — Code Editor, Output Console” |
| Ctrl+E | w. 98 `Moduły środowiska` | :1060 „Otwarcie Export Panel — Report Builder” |
| Ctrl+K | w. 102, 501 `Paleta poleceń` | :1058 „Szybkie dodanie źródła — Sources Manager” |
| Ctrl+N | proponowany dla w. 853 | :1059 „Nowe ustalenie — Findings Panel” |

Patrz P-10.

---

## 4. NARZĘDZIOWNIA — poprawki ikon i etykiet szyny

### 4.1. Do wykonania w tym pliku (9 zmian)

Wypisane w punkcie 2.15: opisy w. 152, 158, 313; nazwy w. 277/278, 319/320, 358/359; ikony w. 172/173, 217/218, 235/236.

### 4.2. Zmierzone pozycje szyny, w których **to okno jest zgodne z dokumentacją, a 33 pozostałe odbiegają** — NIE zmieniać

| Pozycja | Tutaj | W 33 oknach | Rozstrzygające źródło |
|---|---|---|---|
| Browser (w. 175/176, 220/221) | `karta-okna` | `globus` × 66 | KAT:331 — „Browser \| `karta-okna` \| Ikona z grupy nawigacja”; MAN: `karta-okna` = „Okno, karta przeglądarki, **moduł Browser**” |
| Workspace (w. 190, 202, 250, 283, 303 — ikony w wierszach +1) | `warstwy` | `folder` × 132 | KAT:334 — „Workspace \| `warstwy`”; MAN: „Moduł Workspace — przestrzeń projektowa” |
| Automations (w. 300/301) | `automatyzacja` | kształt spoza zestawu × 33 | KAT:335 — „Automations \| `automatyzacja`” |
| Always On Display (w. 309/310) | `<b>Always On Display</b>` | `<b>Always on Display</b>` × 33 | pisownia produktowa: 844 wystąpień „Always On Display” wobec 136 „Always on Display” w `design` + `docs` + `prowadzenie`; w `docs` i `prowadzenie` forma z małym „on” nie występuje ani razu |

Te cztery pozycje należy zgłosić Prowadzącemu jako poprawkę **składnika wspólnego w 33 plikach**, nie tutaj. Patrz P-11.

### 4.3. Pozostałe pozycje szyny — zgodne, bez uwag (42)

Sekcja „Praca” (3), emblematy czterech środowisk (`srodowisko-talkin/-workspace/-codestudio/-multitaskingai`), moduły TalkIn (9, kolejność zgodna z `docs/srodowiska/talkin.md:248–256`), moduły WorkSpace (9) i CodeStudio (8) — kolejność identyczna we wszystkich 34 oknach, więc nie jest odstępstwem tego pliku (patrz P-8), szybki wybór (6) i „Dodaj skrót”.

### 4.4. Pas narzędzi i pasek okna

- **`menu-pasek-nadmiar` (w. 899–901) — nie zmieniać.** Menu ma nagłówek „Poza paskiem” i zero pozycji, ponieważ jego zawartość powstaje w czasie działania z pozycji paska: `zasoby/pasek-okna.js:7, 29–30, 36, 85–95` (`nadmiar.hidden = zwiniete.length === 0`). Cały pojemnik `dn-pasek-okna` (w. 898) jest ukryty na karcie Centrum dowodzenia. Zapis przebiegu pierwszego („nagłówek i ZERO pozycji”) opisuje stan spoczynkowy, nie usterkę.
- Wyzwalacze i cele menu: sprawdzono komplet — **40 identyfikatorów `id="menu-*"`, 40 wyzwalaczy `data-menu`, zero sierot po obu stronach.**

---

## 5. CZEGO NIE ROBIĆ

| Nie wykonywać | Zarzut |
|---|---|
| Dopisania `Nazwa Z–A` do `menu-projekty-filtr` | Wskazane źródło (w. 559 tego samego pliku) nie jest źródłem — okno nie dowodzi samo sobie. Żadne inne okno nie daje drzewu projektów porządku odwrotnego (`research.html:373–375`, `workspace.html:378–380` niosą „Ostatnia aktywność” i „Nazwa”, bez kierunku). Nadto `panel-sesji.js:153–172` obsługuje dla drzewa wyłącznie `nazwa` i gałąź `else` — dopisanie wymaga zmiany skryptu, nie samego znacznika. Asymetria jest szersza, niż zgłoszenie przyznawało: wykaz sesji ma pięć pozycji porządku, drzewo dwie; brakuje też „Najstarsze” i „Środowisko”. → P-12 |
| Dopisania `Zarchiwizowane` do `menu-projekty-filtr` | Brzmienie występuje w 2 plikach, przy czym `workspace.html:377` dotyczy filtra **projektów modułu Workspace**, a nie drzewa projektów w panelu powłoki. Kontrakt zna `WorkspaceProjectStatus.archived`, ale nie ma komendy wydającej wykaz projektów archiwalnych. → P-12 |
| Zmiany nazwy w. 852 `Nowe okno komunikacji` i w. 865 `Duplikuj bieżącą kartę` | Zgłoszono jako „do-poprawy-nazwy” z nazwą właściwą **identyczną z obecną**. Poprawiona ma być wyłącznie ikona (punkt 2.14). |
| Dopisania `Assistant` do `menu-nowa-karta` | Teza „jedyny moduł spoza menu, który zakłada kartę sesji” jest nieprawdziwa: `data-tworzy-sesje="tak"` niosą także Zespoły (265), Role (268), Kolejki (271), Orkiestracja (274), Harmonogram i automatyki (277), Monitor procesu (280). Szyna ma 14 różnych modułów z `tworzy-sesje="tak"`, blok menu miałby 8 — twierdzenie o pokryciu zbiorów jest fałszywe. |
| Ikony `minus-kolo` dla `Wyjmij z projektu` | MAN opisuje `minus-kolo` jako „Stan zajętości — dostępność wstrzymana” — znak stanu, nie polecenia. |
| Wymiany 157 kształtów spoza zestawu na pliki z `design/zasoby/ikony/svg/` | Kształty spoza zestawu niesie 66 % SVG w `design/05-okna/` i sam WZORZEC. To okno ma najniższy udział z 36 plików (27 %). Wymiana w jednym oknie rozjedzie je z resztą produktu. → P-13 |
| Uzupełniania `menu-pasek-nadmiar` (w. 899–901) | Menu wypełnia `zasoby/pasek-okna.js` przy zwężeniu okna; puste w spoczynku jest stanem zamierzonym. |
| Zmiany opisu w. 262 („sześć sekcji sterowania”) | Opis jest zgodny z `docs/srodowiska/multitaskingai.md:68, 95, 636–653`. Rozbieżny jest skład rozwinięcia (8 pozycji), nie opis. |
| Zmiany ikon Browser / Workspace / Automations i pisowni „Always On Display” w szynie | To okno jest tu zgodne z KAT i z pisownią produktową; odbiegają 33 pozostałe pliki. |
| Dodawania skrótu do w. 327 `Skróty klawiszowe` | WZ:219 skrótu nie nadaje, a `Ctrl+/` w tym oknie obsadzają dwie pozycje o tej samej nazwie (w. 138, 499). |
| Dodawania skrótu do w. 583 `Panel po prawej` | `Ctrl+B` jest zajęty (w. 81, 874) i sam wymaga rozstrzygnięcia. |
| Usuwania w. 429, 430, 445 | Brak pokrycia we wzorcu nie jest dowodem, że czynność jest zbędna. Do rozstrzygnięcia (P-2), nie do usunięcia. |

---

## 6. PYTANIA DO WŁAŚCICIELA

**P-1. „Nie przeszkadzać” (w. 352).** Etykieta nie występuje w kontrakcie, w `docs/`, w `design/01-dokumentacja-md/`, we WZORCU ani w żadnym z 33 pozostałych okien — jedyne wystąpienie w produkcie to ten wiersz. W zestawie jest za to `dzwonek-wyciszony` („Powiadomienia wyciszone”). *Czy powłoka ma mieć przełącznik wyciszenia powiadomień — a jeśli tak, pod jakim brzmieniem i w którym menu, skoro grupa „Szybkie przełączniki” zostaje usunięta?*

**P-2. Pozycje bez pokrycia we wzorcu, po jednym wystąpieniu w produkcie: w. 429 `Pokazuj kursor`, w. 430 `Format pliku`, w. 445 `Przypnij pozycję`.** *Wchodzą do wzorca menu zrzutu i schowka dla wszystkich 34 okien, czy znikają z tego okna?*

**P-3. `menu-aod` (w. 513–533).** Wzorzec (WZ:319) ma jedną grupę i siedem pozycji; to okno ma cztery grupy i jedenaście, w tym trzy własne (`Tryb obecności`, `Położenie awatara`, `Środowiska`, `Moduły`). Wzorcowa pozycja `Powierzchnia interakcji` nie jest przełącznikiem, tutejsza `Włącz Always On Display` jest. *Który układ jest wiążący — wzorcowy, czy rozbudowany z tego okna? Jeżeli rozbudowany, czy przenieść go do wzorca i do 33 okien?*

**P-4. Glify bez roli w źródle.** Dla pięciu pozycji ani MAN, ani KAT nie nazywa roli właściwej: w. 560 `Środowisko` (dziś `warstwy` = moduł Workspace; kandydat `monitor` = „Stanowisko, **środowisko wykonania**”), w. 822 `Zamknij pozostałe` i w. 824 `Przywróć ostatnią kartę` (dziś oba `dom`), w. 832–835 `Nowe okno` / `Wydziel kartę` / `Przenieś kartę` / `Scal okna` (dziś wszystkie `okno-plywajace`), ikona pozycji `Wyjmij z projektu`. *Czy przyjąć wskazane kandydatury, czy zestaw ma zostać rozszerzony o brakujące kształty?*

**P-5. Zestaw nie ma pary kierunkowej dla porządku alfabetycznego.** Po zmianie na `sortowanie` pozycje `Nazwa A–Z` (w. 558) i `Nazwa Z–A` (w. 559) pozostaną nierozróżnialne glifem — odróżni je wyłącznie etykieta i znacznik wyboru. *Czy dodać do zestawu parę `sortowanie-rosnaco` / `sortowanie-malejaco`?*

**P-6. Menu kontekstowe karty.** `docs/srodowiska/workspace.md:361–363` i `docs/srodowiska/talkin.md:459` przypisują sześć czynności („zamknij inne, zamknij po prawej, przypnij, dodaj do grupy, zmień nazwę, duplikuj, przekaż na Mobile”) **menu kontekstowemu karty**, wywoływanemu prawym klawiszem albo kebabem karty. Takiego menu w oknie nie ma; najbliższa grupa to „Czynności na karcie” w `menu-karty-otwarte` (w. 820–824). Nadto `Dodaj do grupy` nie występuje w pliku, ale istnieje `Grupuj karty` (w. 878) o zakresie pasma, nie karty. *Budować menu kontekstowe karty, czy wszystkie te czynności zostają w menu ☰ pasma?*

**P-7. Karty grupy a „Sesje w tle”.** Karty „Klasyfikacja umów” i „Migracja rejestru umów” stoją jednocześnie jako karty otwarte grupy „Rejestr umów” (w. 846, 891, 893) i jako pozycje pod nagłówkiem „Sesje w tle” (w. 840, 841). Dopóki tak zostaje, poprawiona liczba „Otwarte karty (4)” spiera się z tamtym miejscem. *Czy karta należąca do grupy zwiniętej liczy się jako otwarta, czy jako sesja w tle?*

**P-8. Skład rozwinięcia MultitaskingAI i kolejność modułów.** Rozwinięcie MultitaskingAI ma 8 pozycji wobec sześciu sekcji z `multitaskingai.md:636–653` — nadmiarowe są `Workspace` (w. 283) i `Agents` (w. 286). Osobno: kolejność modułów WorkSpace w szynie (Workspace · Studio · Design · Apps · Research · Library · Browser · Roundtable · Agents) różni się od kolejności z `docs/srodowiska/workspace.md:145–153` (Studio · Workspace · Browser · Research · Library · Roundtable · Design · Apps · Agents). **Oba stany są identyczne we wszystkich 34 oknach.** *Czy poprawiać dokumentację, czy 34 pliki naraz?*

**P-9. Kolizje odziedziczone po wzorcu.** (a) `Ctrl+V` nadany dwóm różnym czynnościom w samym wzorcu — WZ:34 „Wklej”, WZ:284 „Wklej ostatnią pozycję”. (b) `Ctrl+W` — WZ:34 przypisuje go pozycji „Zamknij okno”, `HAND:1048` „Zamknięciu bieżącej karty”; w tym oknie „Zamknij okno” ma `Ctrl+Shift+W` (w. 60). *Którą stroną rozstrzygnąć obie?*

**P-10. Przestrzeń skrótów powłoki wobec modułów.** Reguła `HAND:1071` rozstrzyga kolizje „kontekstem okna, bo okna nie współdzielą przestrzeni skrótów”. W tej powłoce okno modułu jest kartą wewnątrz tego samego okna systemowego, więc przestrzeń jest wspólna. Zmierzono osiem nakładań (tabela 3.3), w tym `Ctrl+Shift+A` powłoki wobec „Uruchomienia pełnej analizy” w Diagnostics Center (`HAND:1054`). *Czy skróty powłoki mają pierwszeństwo przed skrótami modułu w karcie z ogniskiem, czy odwrotnie?*

**P-11. Cztery odstępstwa 33 okien od tego (punkt 4.2).** *Poprawiać 33 pliki, czy zrównać to okno z większością wbrew KAT i pisowni produktowej?*

**P-12. Symetria dwóch menu filtrujących w jednym pasku panelu.** `menu-sesje-filtr` daje pięć porządków (Ostatnia aktywność, Najstarsze, Nazwa A–Z, Nazwa Z–A, Środowisko), `menu-projekty-filtr` dwa. Nadto w 19 oknach wymiar „Środowisko” stoi pod nagłówkiem **Grupowanie**, nie „Sortowanie”, a menu filtrujące modułów mają trzecią grupę „Grupowanie” (`research.html:376`), której tu nie ma. *Czy oba menu panelu mają dostać pełny, symetryczny zestaw z grupą „Grupowanie”?*

**P-13. Zestaw ikon jako źródło.** 66 % kształtów w `design/05-okna/` nie odpowiada żadnemu z 152 plików zestawu; dotyczy to również WZORCA. *Czy zestaw ma być wiążący dla znacznika okien — a jeśli tak, kiedy i w jakiej partii przeprowadzić zrównanie?*

**P-14. Menu wierszy projektu (`menu-projekt-1/-2`).** Kontrakt nie ma komend zmiany nazwy ani usunięcia projektu (`workspace.project.*` to wyłącznie `status.set`). *Dopisać komendy do kontraktu, czy zdjąć obie pozycje z menu?*

**P-15. Rozgraniczenie „Ustawienia ↔ Konfiguracja”.** `platformowe/ustawienia.html:556–568` dzieli nastawy na 8 sekcji okna Ustawień i 13 zakresów okna Konfiguracji — 21 celów. Menu wspólne ma po poprawkach 18 pozycji i nie odwzorowuje podziału jeden do jednego (np. „Wtyczki i konektory” pokrywa zakresy 5.7 Rozszerzenia i 5.8 Integracje). Rozbieżność tkwi we wzorcu. *Rozstrzygnąć na poziomie składnika wspólnego przed albo po tej partii?*

**P-16. Stany sesji.** `SessionStatus` kontraktu ma cztery wartości (`active`, `paused`, `finished`, `archived`); filtr wykazu ma cztery pozycje, ale inne: Wszystkie · W toku · Czekające · Zakończone. `paused` („sesja wstrzymana; praca może być wznowiona”) nie ma odpowiednika w oknie, a „Czekające” (`data-stan="reakcja"`, `panel-sesji.js:38–39`) nie ma odpowiednika w kontrakcie. *Który zestaw jest wiążący?*

---

## 7. POMIAR PO KAŻDEJ PARTII

Przed pierwszą zmianą wykonać kopię pliku. Po każdej partii uruchomić komplet z tabeli — wynik oczekiwany podano obok. Wszystkie polecenia z katalogu `/home/ubuntu/budowa`.

### 7.1. Miary niezmiennicze — muszą dać tę samą wartość przed i po każdej partii

| Co | Polecenie | Wartość przed |
|---|---|---|
| Liczba wierszy | `wc -l design/05-okna/przeplyw/centrum-dowodzenia.html` | 1510 (zmienia się tylko o liczbę dopisanych/usuniętych wierszy danej partii) |
| Menu i wyzwalacze bez sierot | `python3 -c "import re,sys;s=open('design/05-okna/przeplyw/centrum-dowodzenia.html',encoding='utf-8').read();i=set(re.findall(r'id=\"(menu-[^\"]+)\"',s));r=set(re.findall(r'data-menu=\"([^\"]+)\"',s));print(len(i),len(r),sorted(i^r))"` | `40 40 []` |
| Bilans znaczników `button` | `grep -o '<button' … \| wc -l` wobec `grep -o '</button>' … \| wc -l` | równe |
| Poprawność XML/HTML pozycji menu | `python3 -c "…"` — liczba `sta-menu-poz` z `role^=menuitem` | 358 przed partią 2; po partii 2.1 → 354; po 2.2 → 358; po 2.4 → 361; po 2.9 → 360; po 2.10 → 374; po 2.12 → 377 |
| Brak nowych kolizji skrótów | skrypt z punktu 3 (grupowanie `<span class="skrot">` po treści) | 13 powtórzeń przed; po partii skrótów → 7, wszystkie z tabeli 3.2 |

### 7.2. Kontrola po każdej partii nazewniczej

```
cd design/05-okna
grep -c "Wszystkie ustawienia\|Język i formaty\|Wygląd<\|Modele i kanały\|Kopie zapasowe<" przeplyw/centrum-dowodzenia.html   # oczekiwane 0
for s in "Ustawienia aplikacji" "Język i format zapisu" "Motyw i gęstość widoku" \
         "Modele, kanały i poziom wysiłku" "Magistrala kontekstu i pamięć" \
         "Kopie zapasowe i przywracanie" "Pamięć podręczna i dane lokalne"; do
  printf '%-40s %s\n' "$s" "$(grep -rl "$s" --include='*.html' . | wc -l)"; done   # każde ma dać 34
```

Reguła ogólna: **każde brzmienie poprawione w tym oknie musi po zmianie występować w 34 plikach, a brzmienie zastąpione w 0.** Ta sama kontrola dla partii `menu-zrzut`, `menu-schowek`, `menu-profil-szyna`, `menu-wiecej`, `menu-aod` i szyny.

### 7.3. Kontrola ikon po partii ikonowej

Uruchomić mapowanie sygnatur kształtów (skrypt użyty w tym przebiegu) i sprawdzić:
- `sortowanie` przy w. 558, 559, 570; `wielkosc-liter` już tylko przy w. 75 i 495;
- `ptaszek-kolo` przy w. 553; `ptaszek` już tylko jako znacznik wyboru `<span class="ptaszek">`;
- `biblioteka` przy Library (2 wystąpienia) → w produkcie 68; `ksiazka` przy Library → 0;
- `kod` przy Developer w szynie → w produkcie 34; `plik-kodu` przy Developer → 0;
- glif `dom` niesie już tylko pozycje 97 i 813 (Centrum dowodzenia), nie 821–824;
- glif `okno-plywajace` niesie już tylko pozycje Always On Display, nie 840/841.

### 7.4. Kontrola zachowania okna (prototyp)

Otworzyć `http://127.0.0.1:8611/design/05-okna/przeplyw/centrum-dowodzenia.html` i sprawdzić:

1. Otwiera się każde z 40 menu; żadne nie jest puste poza `menu-pasek-nadmiar` przy szerokim oknie.
2. Konsola przeglądarki — **zero błędów skryptu** (kontrolne: `menu.js`, `panel-sesji.js`, `karty-okna.js`, `pasek-okna.js`).
3. Filtry panelu: przełączenie każdej z 4 pozycji „Stan” i 5 pozycji „Sortowanie” zmienia wykaz i wypowiada komunikat czytnika; po zmianie w. 556 komunikat brzmi „Wykaz sesji: ostatnia aktywność.” — **wymaga równoległej zmiany `zasoby/panel-sesji.js:41`**.
4. `Wyjmij z projektu` po dopisaniu: pozycja ma być widoczna we wszystkich 14 menu wierszy; jeżeli dostanie uchwyt `data-poz-akcja`, uzupełnić słownik komunikatów `panel-sesji.js:183–188`.
5. Po usunięciu w. 348–353: menu Ustawień ma 4 grupy i 18 pozycji, kończy się na „Diagnostyka aplikacji”; przycisk „Motyw” w stopce szyny (w. 358) nadal przełącza motyw.
6. Po usunięciu w. 578: przycisk „Archiwum sesji” (w. 545) nadal otwiera historię (`data-otwarz-historie`).
7. Zwężenie okna do ~1200 px i poniżej: grupa `dn-narzedzia-grupa--proba-1200` zwija się pod „Więcej”, `menu-pasek-nadmiar` zapełnia się na karcie modułu.
8. Motyw jasny i ciemny — obejrzeć oba; sprawdzić kontrast nowych pozycji menu profilu i grupy „Ostatnie pozycje” schowka.
9. Szerokość odniesienia przy oglądzie geometrii: **~2560 px**, nie 1600.

### 7.5. Bramka końcowa

```
python3 narzedzia/style_guard.py design/05-okna/przeplyw/centrum-dowodzenia.html   # jeżeli walidator obejmuje HTML okien
```
Jeżeli walidator nie obejmuje tego rozszerzenia — **nie zmierzono**; wtedy bramką jest komplet 7.1–7.4.