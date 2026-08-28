import { GRANICA_NIEPODANA, narzedzie, profil, type ProfilModulu } from './profil-modulu';

/**
 * Profile pięciu modułów pracy z wiedzą i przekazem — Library, Translate, Roundtable, Design i Assistant — obejmujące dwa moduły bez rozmowy w postaci okna.
 */
export const PROFILE_WIEDZY: readonly ProfilModulu[] = [
  profil(
    'library',
    'Library',
    'Repozytorium plików, wersji, etykiet i kolekcji',
    ['Library Explorer', 'Tags & Collections', 'File Preview', 'Versioning Panel'],
    [
      narzedzie(
        'library.pokaz',
        'Pokaż w Explorerze',
        'Wskazanie pliku w Library Explorer (library.file.list)',
        'Pokaż w Library Explorer plik, o którym mowa, i podaj jego moduł źródłowy.',
      ),
      narzedzie(
        'library.opis',
        'Opisz plik',
        'Generowanie opisu i metadanych pliku (library.file.preview)',
        'Zbuduj opis i metadane wskazanego pliku na podstawie jego zawartości.',
      ),
      narzedzie(
        'library.etykiety',
        'Etykiety i kolekcje',
        'Sugestia porządkująca z akcją zastosowania (library.tag.set)',
        'Zaproponuj etykiety i kolekcję dla wskazanych plików wraz z uzasadnieniem.',
      ),
      narzedzie(
        'library.wersje',
        'Wersje dokumentu',
        'Przegląd wersji i przywrócenie wcześniejszej (library.version.list)',
        'Pokaż wersje wskazanego dokumentu i wskaż, która nadaje się do przywrócenia.',
      ),
    ],
    {
      // Library nie prowadzi rozmowy — pokazuje eksplorator, bo puste okno czatu czytałoby się jak awaria.
      postacRozmowy: 'brak',
      granicaOkien: 1,
      pamiecSesyjna: false,
    },
  ),
  profil(
    'translate',
    'Translate',
    'Tłumaczenie równoległe z glosariuszem',
    ['Source Panel', 'Translation Panels', 'Glossary Manager'],
    [
      narzedzie(
        'translate.jezyk',
        'Dodaj język',
        'Nowy panel tłumaczenia dla kolejnego języka (translate.target.add)',
        'Dodaj panel tłumaczenia dla języka: ',
      ),
      narzedzie(
        'translate.zwrotne',
        'Tłumaczenie zwrotne',
        'Kontrola przekładu przez tłumaczenie wsteczne',
        'Wykonaj tłumaczenie zwrotne wskazanego panelu i pokaż rozbieżności wobec źródła.',
      ),
      narzedzie(
        'translate.jakosc',
        'Kontrola jakości',
        'Liczby, daty, waluty, placeholdery, długość, segmenty pominięte',
        'Sprawdź jakość przekładu: liczby, daty, waluty, placeholdery, długość segmentów ' +
          'i fragmenty pominięte.',
      ),
      narzedzie(
        'translate.glosariusz',
        'Glosariusz',
        'Ujednolicenie terminu w panelach (translate.glossary.set)',
        'Ujednolić termin we wszystkich panelach i dopisz go do glosariusza: ',
      ),
    ],
    {
      // Translate pracuje w układzie dwupanelowym: źródło i tłumaczenie, stąd granica dwóch okien.
      postacRozmowy: 'okno',
      granicaOkien: 2,
      pamiecSesyjna: true,
    },
  ),
  profil(
    'roundtable',
    'Roundtable',
    'Debata wielu modeli z moderacją',
    ['Model Panels', 'Debate Panel', 'Moderator Panel', 'Consensus Panel'],
    [
      narzedzie(
        'roundtable.pytanie',
        'Pytanie do wszystkich',
        'Jedno pytanie równolegle do wszystkich uczestników debaty',
        'Zadaj to pytanie jednocześnie wszystkim uczestnikom debaty: ',
      ),
      narzedzie(
        'roundtable.tura',
        'Kolejna tura',
        'Uruchomienie następnej tury debaty (roundtable.debate.start)',
        'Uruchom kolejną turę debaty i podaj kolejność głosu.',
      ),
      narzedzie(
        'roundtable.moderacja',
        'Interwencja moderatora',
        'Ukierunkowanie dyskusji i zamknięcie tury (roundtable.moderator.direct)',
        'Wykonaj interwencję moderującą: zawęź temat debaty do — ',
      ),
      narzedzie(
        'roundtable.konsensus',
        'Do Consensus Panel',
        'Zapis stanowiska końcowego debaty (roundtable.consensus.get)',
        'Złóż stanowisko końcowe z dotychczasowych tur i zapisz je w Consensus Panel.',
      ),
    ],
    {
      // Roundtable prowadzi wiele równoległych okien po uczestniku; wiąże je sufit sceny, nie profil.
      postacRozmowy: 'okno',
      granicaOkien: GRANICA_NIEPODANA,
      pamiecSesyjna: true,
    },
  ),
  profil(
    'design',
    'Design',
    'Praca wizualna: kompozycja, zasoby, generowanie',
    ['Design Board', 'Assets Panel', 'Prompt Builder', 'Preview Window'],
    [
      narzedzie(
        'design.generuj',
        'Generuj zasób',
        'Polecenie generujące wprost z okna komunikacji (design.asset.generate)',
        'Wygeneruj zasób wizualny — temat, styl, kompozycja, oświetlenie, paleta: ',
      ),
      narzedzie(
        'design.warianty',
        'Warianty',
        'Powtórzenie generowania z modyfikacją parametrów',
        'Powtórz generowanie ostatniego zasobu z modyfikacją: ',
      ),
      narzedzie(
        'design.szablon',
        'Zapisz jako szablon',
        'Zapis promptu w Prompt Builderze do ponownego użycia',
        'Zapisz ostatnie polecenie generujące jako szablon w Prompt Builderze.',
      ),
      narzedzie(
        'design.kompozycja',
        'Kompozycja',
        'Zestawienie zasobów na Design Board (design.board.update)',
        'Zestaw wskazane zasoby na Design Board i opisz układ kompozycji.',
      ),
    ],
    {
      // Okno rozmowy jest głównym kanałem współpracy z modelem, a obok stoi Design Board jako płótno.
      postacRozmowy: 'okno',
      granicaOkien: GRANICA_NIEPODANA,
      pamiecSesyjna: true,
      oknaObowiazkowe: ['Design Board'],
    },
  ),
  profil(
    'assistant',
    'Assistant',
    'Asystent głosowy z historią działań',
    ['Voice Console', 'Actions Monitor', 'Activity Feed'],
    [
      narzedzie(
        'assistant.polecenie',
        'Polecenie asystenta',
        'Zlecenie wieloetapowe wykonywane przez asystenta (assistant.voice.command)',
        'Wykonaj zlecenie wieloetapowe i pokaż jego postęp w Actions Monitor: ',
      ),
      narzedzie(
        'assistant.status',
        'Status zlecenia',
        'Stan realizacji zleceń na żywo (assistant.action.status)',
        'Podaj status zleceń biegnących w Actions Monitor wraz z etapem każdego z nich.',
      ),
      narzedzie(
        'assistant.dziennik',
        'Dziennik działań',
        'Chronologiczny zapis działań i odtworzenie przebiegu (assistant.activity.list)',
        'Pokaż dziennik działań z Activity Feed i odtwórz przebieg ostatniego zlecenia.',
      ),
      narzedzie(
        'assistant.popraw',
        'Popraw polecenie',
        'Poprawa polecenia po transkrypcji głosu',
        'Popraw treść ostatniego polecenia i wykonaj je ponownie w brzmieniu poprawionym.',
      ),
    ],
    {
      // Assistant nie ma okna rozmowy — mówi awatarem, jednym dymkiem, warstwą wsparcia platformy.
      postacRozmowy: 'dymek-glosowy',
      granicaOkien: 1,
      pamiecSesyjna: true,
    },
  ),
];
