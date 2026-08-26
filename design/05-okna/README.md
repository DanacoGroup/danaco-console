# Prototypy okien — Danaco Console

Ten katalog gromadzi prototypy wszystkich okien platformy w postaci plików HTML.
Każdy prototyp jest samodzielny: renderuje się w przeglądarce bez budowania,
korzysta wyłącznie z arkuszy `design/zasoby/` i żetonów systemu projektowego
`--dn-*`, działa w obu motywach (jasnym i ciemnym) i niesie pełną ramę okna
aplikacji (szyna nawigacji, belka tytułowa, wstążka pozioma, pasek stanu) tam,
gdzie okno należy do wnętrza aplikacji.

Prototypy są **wzorcem odniesienia** dla prac deweloperskich — nie makietą
poglądową. Treści stałe (etykiety, nagłówki, podpowiedzi) są w języku formalnym;
treści robocze imitują realną pracę. Warstwa `.pt-*` (plakietka „PROTOTYP",
noty „O opracowaniu") jest wyłącznie rusztowaniem prototypu i nie wchodzi do
produktu.

## Kolejność czytania

Okna układają się w przepływ od instalacji po pracę w module. Zalecana kolejność:

### 1 · Przed uruchomieniem
| Plik | Okno | Uwagi |
|---|---|---|
| `platformowe/instalator.html` | Instalator aplikacji | Okno przedaplikacyjne — bez ramy aplikacji; wygląd okien wejściowych. Kroki: licencja → katalog → składniki → skróty → instalacja → zakończenie |

### 2 · Przepływ uruchomienia (`przeplyw/`)
Osiem etapów drogi od uruchomienia do okna operacyjnego. Etapy 1–3 mieści jeden
plik z przełącznikiem etapów w rusztowaniu prototypu. Pliki bez numerów —
kolejność wynika z przepływu, nie z nazwy.
| Plik | Etap | Okno |
|---|---:|---|
| `przeplyw/przeplyw-wejscia.html` | 1–3 | Przepływ wejścia — okno startowe (trzy warianty stanu), rejestracja i logowanie (sześć stanów), przygotowanie środowiska pracy |
| `przeplyw/centrum-dowodzenia.html` | 4 | Strona główna (Centrum dowodzenia) — trzy strefy |

Etapy 5–8 (wybór środowiska, przedsionek, wybór modułu, okno operacyjne)
realizują pliki w `srodowiska/` i `moduly/`.

### 3 · Środowiska i przedsionki (`srodowiska/`)
Każde środowisko ma **przedsionek** (widok po wejściu, przed wyborem modułu —
kafle modułów i szyna sesji) oraz **przestrzeń roboczą** (po wyborze modułu).
| Plik | Okno |
|---|---|
| `srodowiska/talkin-przedsionek.html` | TalkIn — przedsionek |
| `srodowiska/talkin.html` | TalkIn — przestrzeń robocza |
| `srodowiska/workspace-przedsionek.html` | WorkSpace — przedsionek |
| `srodowiska/workspace.html` | WorkSpace — przestrzeń robocza |
| `srodowiska/codestudio-przedsionek.html` | CodeStudio — przedsionek |
| `srodowiska/codestudio.html` | CodeStudio — przestrzeń robocza |
| `srodowiska/multitaskingai-przedsionek.html` | MultitaskingAI — przedsionek |
| `srodowiska/multitaskingai.html` | MultitaskingAI — panel orkiestracji |
| `srodowiska/mtai-okna-rol.html` | MultitaskingAI — okna ról zespołu |
| `srodowiska/mtai-izolacja-i-zespoly.html` | MultitaskingAI — izolacja i zespoły |

### 4 · Moduły (`moduly/`)
Piętnaście modułów w przestrzeni roboczej. Chat Window (lewa kolumna, stały)
i Execution Loop Window są wspólne dla każdego modułu.
| Plik | Moduł | Środowisko wiodące |
|---|---|---|
| `moduly/studio.html` | Studio — redakcja treści długiej | TalkIn / WorkSpace |
| `moduly/library.html` | Library — repozytorium wiedzy | TalkIn |
| `moduly/browser.html` | Browser — praca ze źródłami sieci | TalkIn |
| `moduly/research.html` | Research — badanie zagadnienia | TalkIn |
| `moduly/translate.html` | Translate — przekład | TalkIn |
| `moduly/roundtable.html` | Roundtable — narada modeli | TalkIn |
| `moduly/assistant.html` | Assistant — dialog nad treścią | TalkIn / WorkSpace |
| `moduly/workspace.html` | Workspace — projekty i zadania | WorkSpace |
| `moduly/agents.html` | Agents — definicje agentów | WorkSpace / MultitaskingAI |
| `moduly/automations.html` | Automations — automatyki | WorkSpace |
| `moduly/apps.html` | Apps — budowa produktu | WorkSpace |
| `moduly/design.html` | Design — makiety i grafika | WorkSpace |
| `moduly/developer.html` | Developer — kod i wersje | CodeStudio |
| `moduly/terminal.html` | Terminal — powłoki i procesy | CodeStudio |
| `moduly/diagnostics.html` | Diagnostics — logi i błędy | CodeStudio |

### 5 · Okna platformowe (`platformowe/`)
Okna dostępne z każdego miejsca — konfiguracja, ustawienia, funkcje globalne,
okna kontekstowe strefy pracy szyny.
| Plik | Okno | Wywołanie |
|---|---|---|
| `platformowe/ustawienia.html` | Ustawienia i konfiguracja (poziom aplikacji) | Stopka szyny · menu aplikacji |
| `platformowe/konfiguracja.html` | Konfiguracja — trzynaście zakresów, nakładka i wywołanie modelu | Menu Ustawienia |
| `platformowe/historia-sesji.html` | Historia sesji — nakładka okna ustawień (czynne, zakończone, archiwalne) | Strefa pracy szyny · przedsionek · sekcja ustawień |
| `platformowe/nowy-projekt.html` | Nowy projekt — Workspace środowiska w trybie zakładania projektu | Strefa pracy szyny · listwa przedsionka |
| `platformowe/instrukcja-uzytkowania.html` | Instrukcja użytkowania — nakładka wywoływana z menu Pomoc | Menu aplikacji → Pomoc |
| `platformowe/mobile.html` | Mobile — widok mobilny i interwencja zdalna | Listwa ustawień · szybki wybór szyny |
| `platformowe/always-on-display.html` | Always On Display — awatar pływający | Listwa ustawień · szybki wybór szyny |

### Wzorzec
| Plik | Rola |
|---|---|
| `WZORZEC-STANOWISKA.html` | Wzorzec przestrzeni roboczej (stanowiska) — punkt odniesienia dla układu kolumn każdego modułu |

## Kluczowe uwagi

- **Rama okna jest jedna dla wszystkich okien wnętrza aplikacji.** Szyna
  nawigacji biegnie od górnej krawędzi; belka tytułowa i wstążka mieszczą się na
  prawo od niej; pasek stanu zamyka okno u dołu. Pełna specyfikacja:
  `docs/interfejs-uzytkownika/rama-okna.md`.
- **Szyna ma trzy strefy** rozdzielone delikatną kreską: praca (Nowa sesja,
  Historia sesji, Nowy projekt), środowiska, szybki wybór. Głowa (menu
  aplikacji) i stopka (ustawienia, dostosowanie paska, motyw, profil) stoją
  nieruchomo.
- **Wstążka pozioma jest belką nakładaną**, personalizowaną w oknie „Dostosuj"
  (zakładki Składniki, Karty sesji, Wygląd i zachowanie, Personalizacja widoku).
  Ten sam składnik można przenieść na szynę pionową i odwrotnie.
- **Okna nakładkowe** (Dostosuj, Ustawienia, Historia, Instrukcja) mają ramę
  okna: przeciąganie, zmianę rozmiaru, minimalizację i maksymalizację
  (`design/zasoby/okna-modalne.js`).
- **Trwałość pracy w tle**: proces sesji żyje po stronie serwera; zamknięcie
  okna nie przerywa pracy. Karty sesji i historia odtwarzają pełny stan.
- **Instalator** naśladuje wygląd okien wejściowych, ale stoi przed uruchomieniem
  aplikacji, więc nie ma ramy aplikacji — ma własny minimalny pasek tytułu.
- **Renderowanie do weryfikacji**: headless Chromium, np.
  `chromium-browser --headless --screenshot=out.png "file://…/moduly/studio.html"`.

## Zależności

Wszystkie prototypy zależą od `design/zasoby/`:
`zetony/` (żetony i fonty), `css/fundament.css`, `css/komponenty.css`,
`prototyp.css` (+ `.js`), `rama.css` (+ `.js`), `okna-modalne.css` (+ `.js`),
oraz — okna wejściowe — `wejscie.css`. Prototypy przestrzeni roboczej dodają
`stanowisko.css` (+ `.js`), przedsionki — `przedsionek.css`.
Instalator korzysta w pełni z `wejscie.css`: jego style `.in-*` stoją w tym
arkuszu, a nie w samym pliku okna.
