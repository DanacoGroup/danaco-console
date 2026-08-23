import { utworzDymekObjasnienia } from '../../komponenty/dymek';
import { opisOdmowy } from '../../komponenty/odmowa';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import { poleTekstowe, przycisk, utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import { utworzCzynnosciZbiorcze } from './czynnosci-zbiorcze';
import { utworzNarzedziaExplorera } from './narzedzia-explorera';
import { utworzOdbiorPrzekazania, type OdbiorPrzekazania } from './odbior-przekazania';
import { utworzPanelAkcji, type PanelAkcji } from './panel-akcji';
import type { StanBiblioteki, WidokWykazu } from './stan-biblioteki';
import { utworzStanOkna } from './stan-okna';
import { utworzWgraniePliku } from './wgranie-pliku';
import { WIDOKI } from './widoki-wykazu';
import { utworzWykazPlikow } from './wykaz-plikow';
import { BEZ_KOMENDY_EXPLORER } from './etykiety-biblioteki';
import type { ZrodloOtoczenia } from './zrodlo-otoczenia';

/**
 * Library Explorer — okno wiodące modułu (`library.library-explorer`).
 *
 * Cztery funkcje Operatora z wiersza wykazu, każda z własną drogą do rdzenia:
 * nawigacja po strukturze (`library.file.list`), wyszukiwanie w trzech trybach
 * (`library.file.search` po słowach, `knowledge.search` po znaczeniu, hybryda
 * złożona z obu w oknie), otwarcie zasobu (wskazanie pliku przestawia cztery
 * pozostałe okna, a otwarcie w module źródłowym idzie `context.transfer`) oraz
 * wgranie pliku (`library.file.upload`).
 *
 * Prezentacja wykazu ma pięć postaci — siatka miniatur jako domyślna oraz
 * lista, galeria, oś czasu i mapa spod przełącznika widoku. Wszystkie liczą się
 * z tej samej odpowiedzi rdzenia i nie wysyłają ani jednej komendy
 * (`widoki-wykazu.ts`).
 *
 * Okno jest źródłem zaznaczenia dla całego modułu: File Preview, Versioning
 * Panel, Tags & Collections oraz Metadata & Archive Panel nie mają własnego
 * wejścia — biorą plik czynny i zaznaczenie stąd, przez wspólny stan.
 */
export interface OknoExplorer {
  element: HTMLElement;
  odswiez(): void;
  /** Pierwszy odczyt okna: wykaz plików, katalog akcji i katalog modułów. */
  wczytaj(): Promise<void>;
  /** Odpina nasłuch przybycia przekazania (`window.changed`). */
  rozlacz(): void;
}

export function utworzOknoExplorer(
  stan: StanBiblioteki,
  otoczenie: ZrodloOtoczenia,
): OknoExplorer {
  const rama = utworzRameOkna({
    kod: 'library-explorer',
    tytul: 'Library Explorer',
    rola: 'wiodące',
    modul: 'Library',
    dodatkiNaglowka: [
      utworzDymekObjasnienia(
        'Repozytorium wiedzy platformy: struktura plików, wyszukiwanie po treści i znaczeniu, ' +
          'wgranie pliku oraz odbiór artefaktów wytworzonych w innych modułach.',
        { powloka: 'ml-dymek', znak: 'ml-dymek__znak' },
      ),
    ],
  });
  const okno = utworzStanOkna();
  const wykaz = utworzWykazPlikow(stan);
  const narzedzia = utworzNarzedziaExplorera(stan);
  const czynnosci = utworzCzynnosciZbiorcze(stan, otoczenie);
  const panel: PanelAkcji = utworzPanelAkcji(otoczenie, stan, BEZ_KOMENDY_EXPLORER);
  const odpowiedz = utworzWierszOdpowiedzi();
  const wgranie = utworzWgraniePliku(stan, odpowiedz.pokaz);
  // Odbiór przekazania stoi w oknie wiodącym, bo to ono prowadzi wykaz plików
  // i zaznaczenie — komplet przybyły z innego modułu trafia tam, gdzie Operator
  // może z nim cokolwiek zrobić.
  const odbior: OdbiorPrzekazania = utworzOdbiorPrzekazania(stan, otoczenie);

  const fraza = poleTekstowe({
    etykieta: 'Fraza',
    podpowiedz: 'fragment nazwy albo treści',
    opis: 'Odczyt wykazu zawęża frazą; wyszukiwanie sięga treści i znaczenia po stronie rdzenia.',
  });
  const etykieta = poleTekstowe({
    etykieta: 'Etykieta zawężająca',
    podpowiedz: 'np. umowy',
  });

  /**
   * Mówi, że okno wysłało co innego, niż Operator napisał.
   *
   * Okno przycina pola przed wysłaniem, a rdzeń Library dopasowuje etykietę
   * dosłownie — nie przycina ani zapisu, ani filtru. `trim()` przeglądarki
   * zdejmuje przy tym znaki, które rdzeń zostawia (U+00A0, U+FEFF), więc
   * etykieta zapisana z takim znakiem na brzegu jest z tego pola nieosiągalna,
   * a pusty wykaz wyglądałby jak zdanie o repozytorium.
   */
  function ostrzezOPrzycieciu(): void {
    const zmienione: string[] = [];
    if (fraza.kontrolka.value.trim() !== fraza.kontrolka.value) {
      zmienione.push(`frazę wysłano jako „${fraza.kontrolka.value.trim()}"`);
    }
    if (etykieta.kontrolka.value.trim() !== etykieta.kontrolka.value) {
      zmienione.push(`etykietę wysłano jako „${etykieta.kontrolka.value.trim()}"`);
    }
    if (zmienione.length === 0) return;
    odpowiedz.pokaz(
      `Okno przycięło pola przed wysłaniem: ${zmienione.join('; ')}. Rdzeń dopasowuje ` +
        'etykietę dosłownie, więc zapis ze znakiem na brzegu tym polem się nie znajdzie.',
      false,
    );
  }

  const odswiezWykaz = przycisk('Odśwież wykaz', 'dn-btn dn-btn--sm dn-btn--zarys');
  odswiezWykaz.dataset['czynnosc'] = 'wykaz';
  odswiezWykaz.addEventListener('click', () => {
    ostrzezOPrzycieciu();
    void stan.odczytaj(fraza.kontrolka.value, etykieta.kontrolka.value);
  });

  const szukaj = przycisk('Szukaj', 'dn-btn dn-btn--sm dn-btn--zarys');
  szukaj.dataset['czynnosc'] = 'szukaj';
  szukaj.addEventListener('click', () => void wykonajSzukanie());

  async function wykonajSzukanie(): Promise<void> {
    if (fraza.kontrolka.value.trim() === '') {
      odpowiedz.pokaz(
        'Fraza jest pusta po przycięciu przez okno — nic nie zostało wysłane.',
        false,
      );
      return;
    }
    ostrzezOPrzycieciu();
    // Zdanie o drodze wyniku jest tu obowiązkowe: trzy tryby dają wyniki
    // nieporównywalne, a pusty wykaz w trybie semantycznym może znaczyć „nie ma
    // wskaźnika znaczenia", a nie „nie ma takich plików". O powodzeniu orzeka
    // faza wykazu, nie treść zdania — zdanie mówi to samo o odmowie i o wyniku.
    const zdanie = await stan.szukaj(fraza.kontrolka.value);
    odpowiedz.pokaz(zdanie, stan.faza() !== 'blad');
  }

  const pasekWyszukiwania = document.createElement('div');
  pasekWyszukiwania.className = 'ml-explorer__szukanie';
  pasekWyszukiwania.append(fraza.element, etykieta.element, odswiezWykaz, szukaj, wgranie.element);

  okno.tresc.append(wykaz.element, czynnosci.element);
  okno.pierwszaAkcja('Dodaj plik', () => wgranie.wskazPlik());
  rama.narzedzia.append(narzedzia.element);
  rama.cialo.append(
    pasekWyszukiwania,
    odbior.element,
    okno.element,
    odpowiedz.element,
    panel.element,
  );

  /** Zdanie widoku, który nie ma czego pokazać mimo niepustego wykazu. */
  let brakWidoku = '';

  function odswiez(): void {
    brakWidoku = wykaz.odswiez();
    narzedzia.odswiez();
    czynnosci.odswiez();
    ustawStan();
  }

  function ustawStan(): void {
    if (stan.faza() === 'odczyt') {
      okno.ladowanie('Rdzeń odczytuje wykaz plików repozytorium.');
      return;
    }
    if (stan.faza() === 'blad') {
      okno.blad(stan.powod());
      return;
    }
    if (stan.faza() === 'spoczynek') {
      okno.puste(
        'Repozytorium nieodpytane',
        'Okno nie pytało jeszcze rdzenia o pliki. Naciśnij „Odśwież wykaz", żeby ' +
          'zobaczyć zawartość repozytorium wiedzy.',
      );
      return;
    }
    if (stan.pliki().length === 0) {
      // Pusty wykaz zawężony to nie puste repozytorium. Zdanie „rdzeń nie ma
      // ani jednego pliku" należy się wyłącznie odczytowi bez zawężenia; przy
      // frazie albo etykiecie w polu rdzeń odpowiedział o tym, o co go pytano,
      // a nie o całym zbiorze.
      const zawezenia: string[] = [];
      if (fraza.kontrolka.value.trim() !== '') zawezenia.push(`fraza „${fraza.kontrolka.value.trim()}"`);
      if (etykieta.kontrolka.value.trim() !== '') {
        zawezenia.push(`etykieta „${etykieta.kontrolka.value.trim()}"`);
      }
      if (zawezenia.length > 0) {
        okno.puste(
          'Bez trafień w zawężeniu',
          `Rdzeń nie zwrócił pliku dla zawężenia: ${zawezenia.join(', ')}. O reszcie ` +
            'repozytorium ta odpowiedź nie mówi — wyczyść pola i odśwież wykaz.',
        );
        return;
      }
      okno.puste(
        'Repozytorium jest puste',
        'Rdzeń nie ma ani jednego pliku. Dodaj pierwszy albo poczekaj na artefakt z modułu.',
      );
      return;
    }
    // Wykaz niepusty, a widocznych pozycji brak — zawęziło je okno, nie rdzeń.
    // Zdanie nazywa oba możliwe zawężenia, bo składają się koniunkcyjnie
    // i Operator ma wiedzieć, które z nich zdjąć.
    if (stan.widoczne().length === 0) {
      const zawezenia: string[] = [];
      const wskazane = stan.zawezenie();
      if (wskazane !== null) zawezenia.push(`zbiór wskazany (${wskazane.opis})`);
      if (stan.katalog() !== '') zawezenia.push(`katalog „${stan.katalog()}"`);
      okno.puste(
        'Zawężenie okna nie zostawiło ani jednej pozycji',
        `Rdzeń oddał plików: ${stan.pliki().length}, a po zawężeniu — ${zawezenia.join(' oraz ')} ` +
          '— nie widać żadnego. Zdejmij zawężenie w pasku narzędzi albo wróć do całego zbioru.',
      );
      return;
    }
    // Widok bez czego pokazać to nie wykaz bez plików: galeria bez obrazów
    // i mapa bez współrzędnych orzekają o widoku, nie o repozytorium.
    if (brakWidoku !== '') {
      okno.puste(`Widok „${nazwaWidoku(stan.widok())}" bez pozycji`, brakWidoku);
      return;
    }
    okno.gotowe();
  }

  return {
    element: rama.element,
    odswiez,

    async wczytaj() {
      // Trzy odczyty niezależne, więc idą razem, a odmowa jednego nie gasi
      // pozostałych. Katalog modułów obsadza ster modułu docelowego
      // w tym oknie i w File Preview — jedna nastawa, więc jeden odczyt.
      await Promise.all([
        stan.odczytaj(fraza.kontrolka.value, etykieta.kontrolka.value),
        panel.wczytaj(),
        wczytajModuly(),
      ]);
    },

    rozlacz: () => odbior.rozlacz(),
  };

  /** Nazwa widoku z przełącznika — do zdania o widoku bez pozycji. */
  function nazwaWidoku(kod: WidokWykazu): string {
    return WIDOKI.find((pozycja) => pozycja.kod === kod)?.nazwa ?? kod;
  }

  /** Katalog modułów do stera; odmowa zostaje nazwana przy sterze, nie zgubiona. */
  async function wczytajModuly(): Promise<void> {
    const wynik = await otoczenie.moduly();
    if (!wynik.udany || wynik.wynik === undefined) {
      stan.ustawModuly(
        [],
        `Katalog modułów nieodczytany — ster ma samą pozycję wytwórcy. ` +
          `${opisOdmowy('Odczyt katalogu modułów', wynik.blad?.code, wynik.blad?.message)}`,
      );
      return;
    }
    stan.ustawModuly(
      wynik.wynik.modules,
      wynik.wynik.modules.length === 0
        ? 'Rdzeń oddał katalog modułów bez ani jednej pozycji — ster ma samą pozycję wytwórcy.'
        : '',
    );
  }
}
