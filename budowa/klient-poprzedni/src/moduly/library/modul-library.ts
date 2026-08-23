import './library.css';

import type { Kanal } from '../../protokol/kanal';
import { utworzNarzedziaMaterialu } from './material-narzedzia';
import { utworzOknoExplorer } from './okno-library-explorer';
import { utworzOknoMetadanych } from './okno-metadata-archive';
import { utworzOknoEtykiet } from './okno-tags-collections';
import { utworzOknoPodgladu } from './okno-file-preview';
import { utworzOknoWersji } from './okno-versioning-panel';
import { utworzPasekKontekstu } from './pasek-kontekstu';
import { utworzStanBiblioteki } from './stan-biblioteki';
import { utworzStrazOdmow } from './straz-odmow';
import { utworzZrodloBiblioteki } from './zrodlo-biblioteki';
import { utworzZrodloOtoczenia } from './zrodlo-otoczenia';
import { utworzZrodloZnaczenia } from './zrodlo-znaczenia';

/**
 * Moduł Library — pięć okien operacyjnych własnych modułu w jednym układzie.
 *
 * Układ wynika z roli okna. Library Explorer jest oknem wiodącym i jedynym
 * z własnym wejściem, więc stoi w pasie pierwszym na całą szerokość. File
 * Preview (pomocnicze) oraz Versioning Panel, Tags & Collections i Metadata &
 * Archive Panel (zarządcy) stoją w pasie drugim: wszystkie odnoszą się do
 * wykazu i do pliku wskazanego wyżej.
 *
 * Jeden plik czynny i jedno zaznaczenie na cały moduł: wskazanie w Explorerze
 * przestawia pozostałe cztery okna naraz, bo stan jest jeden.
 *
 * Przekazanie z innego modułu ma odbiorcę w oknie wiodącym
 * (`odbior-przekazania.ts`). `context.transfer` rozgłasza `window.changed`,
 * a nie `library.file.changed`, i nie zakłada pliku w repozytorium — treść
 * przekazanego kompletu czyta się osobno komendą `aod.context.get`. Nasłuch
 * wyłącznie na `library.file.changed` przekazania by nie zobaczył.
 *
 * Dziesięć komend obszaru `library.*` ma uchwyt w rdzeniu; wpina je
 * `zarejestrujBiblioteke` (`server/internal/core/adapter_modul_library_uchwyty.go`),
 * a port `Biblioteka` wypełnia w `montaz_porty.go` adapter stojący nad
 * repozytorium i katalogiem danych. Cztery z nich — wgranie, dołożenie
 * i przywrócenie wersji oraz ustawienie etykiet — rozgłaszają po udanym
 * wykonaniu `library.file.changed`, i to zdarzenie jest jedyną drogą
 * odświeżenia okien poza ich własnym działaniem. Moduł wysyła komendy
 * naprawdę i pokazuje odpowiedź rdzenia, a odmowę merytoryczną — jako stan
 * błędu okna wraz z powodem.
 *
 * Poza obszarem `library.*` moduł sięga po jedną rodzinę więcej:
 * `knowledge.index` i `knowledge.search` w zakresie biblioteki
 * (`zrodlo-znaczenia.ts`). Czytelnikiem tego zakresu jest repozytorium modułu
 * Library, a trafienie niesie identyfikator pliku biblioteki, więc rodzina
 * opisuje ten sam zbiór — wnosi wyszukiwanie po ZNACZENIU tam, gdzie
 * `library.file.search` dopasowuje słowa.
 *
 * Rodziny arsenału moduł bierze WYŁĄCZNIE drogą ścieżki na dysku Operatora,
 * i to jest rozstrzygnięcie mierzone, nie ostrożność. Wszystkie trzy —
 * `document.*`, `media.*`, `archive.*` — rozwiązują `assetId` przez repozytorium
 * zasobów modułu Design (`adapter_narzedzia_dokument.go`,
 * `adapter_narzedzia_media.go`, `adapter_narzedzia_archiwum.go`), a plik
 * biblioteki leży w innym rejestrze: jego identyfikator wraca stamtąd odmową
 * `not_found`. Przycisk wysyłający identyfikator zasobu biblioteki zawodziłby
 * zawsze, więc go nie ma.
 *
 * Kontrakt niesie jednak drugą drogę źródła — `sourcePath`, treść wciąganą
 * do magazynu pod sumą kontrolną — i tą drogą trzy czynności stoją w obszarze
 * Archiwum panelu Metadata & Archive Panel: rozpoznanie materiału
 * (`media.inspect`), jego przetworzenie (`media.transcode`) i spakowanie
 * archiwum (`archive.pack`). Ich wynik jest zasobem magazynu Designu, nie
 * zasobem biblioteki, i widok mówi to Operatorowi zdaniem, żeby nikt nie wziął
 * jednego za drugie (`material-narzedzia.ts`). Rodzina `document.*` zostaje
 * nadal poza modułem: jej wynikiem byłby dokument, a nad dokumentami moduł
 * Library nie prowadzi ani jednej czynności własnej.
 *
 * Straż odmów (`straz-odmow.ts`) i odmowa `library.unknown` zostają na
 * miejscu, bo mierzą stan rdzenia w chwili wywołania: gdyby uchwyt wypadł,
 * okno powie to z pomiaru, a nie z komentarza.
 */
export interface ModulLibrary {
  /** Element osadzany w obszarze roboczym powłoki. */
  element: HTMLElement;
  /** Zleca odczyt kontekstu okna, wykazu plików i katalogu akcji. */
  wczytaj(idSesji: string): Promise<void>;
  /** Odłącza subskrypcje zdarzeń rdzenia. */
  rozlacz(): void;
}

export function utworzModulLibrary(kanal: Kanal): ModulLibrary {
  const straz = utworzStrazOdmow(kanal);
  const stan = utworzStanBiblioteki(
    utworzZrodloBiblioteki(kanal, straz),
    utworzZrodloZnaczenia(straz),
  );
  const otoczenie = utworzZrodloOtoczenia(kanal, straz);

  const kontekst = utworzPasekKontekstu(stan, otoczenie);
  const explorer = utworzOknoExplorer(stan, otoczenie);
  const podglad = utworzOknoPodgladu(stan, otoczenie);
  const wersje = utworzOknoWersji(stan, otoczenie);
  const etykiety = utworzOknoEtykiet(stan, otoczenie);
  const metadane = utworzOknoMetadanych(stan, otoczenie, utworzNarzedziaMaterialu(kanal));

  const pasZarzadcow = document.createElement('div');
  pasZarzadcow.className = 'ml-modul__pas ml-modul__pas--zarzadcy';
  pasZarzadcow.append(podglad.element, wersje.element, etykiety.element, metadane.element);

  const element = document.createElement('div');
  element.className = 'ml-modul';
  element.dataset['modul'] = 'library';
  element.setAttribute('aria-label', 'Moduł Library — okna operacyjne');
  element.append(kontekst.element, explorer.element, pasZarzadcow);

  function odswiezWszystkie(): void {
    explorer.odswiez();
    podglad.odswiez();
    wersje.odswiez();
    etykiety.odswiez();
    metadane.odswiez();
  }

  const odsubskrybuj = stan.obserwuj(odswiezWszystkie);
  odswiezWszystkie();

  return {
    element,

    async wczytaj(idSesji) {
      // Kontekst okna idzie PIERWSZY: bez okna komunikacji przenoszenie
      // kontekstu nie ma okna źródłowego, a panel akcji nie ma czym wykonać
      // pozycji katalogu. Reszta odczytów jest niezależna i idzie równolegle —
      // odmowa jednego zostaje w jego oknie i nie gasi pozostałych.
      await kontekst.wczytaj(idSesji);
      // Wykaz plików i katalog akcji panelu metadanych są od siebie niezależne,
      // więc idą razem; odmowa jednego zostaje w jego oknie.
      await Promise.all([explorer.wczytaj(), metadane.wczytaj()]);
    },

    rozlacz() {
      odsubskrybuj();
      explorer.rozlacz();
      stan.rozlacz();
      straz.rozlacz();
    },
  };
}
