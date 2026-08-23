import {
  StudioDocumentFormat,
  StudioOcrEngine,
  type StudioInputDevice,
  type StudioRecognitionSettings,
} from '../../../../shared/contract';
import { utworzDymekObjasnienia } from '../../komponenty/dymek';
import {
  poleLogiczne,
  poleTekstowe,
  poleWyboru,
  przycisk,
  ustawPozycje,
} from '../../modele/kontrolki-formularza';
import { DROGI_CYFRYZACJI, NASTAWY_CYFRYZACJI, SKLADNIKI_PAKIETU_SERWERA } from './braki-cyfryzacji';

/**
 * Kontrolki wejściowe narzędziowni cyfryzacji wraz z odczytem ich wartości.
 *
 * ── Pola odpowiadają temu, co przyjmuje rdzeń — wszystkiemu ─────────────────
 * `StudioRecognitionSettings` niesie silnik, zestaw języków, próg pewności,
 * prostowanie skosu, odszumianie, progowanie, przycięcie marginesów, odtwarzanie
 * układu i zakres stron. Wszystkie dziewięć stoi tutaj jako kontrolki
 * ze skutkiem. Poprzednia postać tego pliku miała jedno pole języka i wymuszenie
 * rozpoznania, bo tyle przyjmowała starsza droga `document.text.extract` —
 * i wypisywała obok siebie braki, które przestały być brakami.
 *
 * ── Czego tu nadal nie ma i dlaczego ────────────────────────────────────────
 * Wymuszenia rozpoznania mimo obecnej warstwy tekstowej: `ingest.recognize`
 * rozpoznaje pozycję kolejki bezwarunkowo, więc osobny przełącznik „wymuś" nie
 * miałby czego przestawiać. Wybierania pliku okienkiem systemowym: klient dysku
 * nie czyta i nie zapisuje — podaje wskazanie widziane przez rdzeń.
 */

/** Dwa źródła materiału odpowiadające dwóm polom żądania kontraktu. */
const ZRODLA_MATERIALU = [
  { wartosc: 'zasob', etykieta: 'Zasób magazynu rdzenia (assetIds)' },
  { wartosc: 'sciezka', etykieta: 'Ścieżka widziana przez rdzeń (sourcePaths)' },
];

/** Dwa silniki rozpoznawania z kontraktu. */
const SILNIKI = [
  { wartosc: '', etykieta: 'silnik z katalogu ustawień' },
  { wartosc: StudioOcrEngine.Lokalny, etykieta: 'lokalny — materiał nie opuszcza maszyny rdzenia' },
  { wartosc: StudioOcrEngine.Chmurowy, etykieta: 'chmurowy — usługa zewnętrzna z jawnym kluczem' },
];

/** Formaty dokumentu zakładanego z wyniku cyfryzacji. */
const FORMATY_DOKUMENTU = [
  { wartosc: '', etykieta: 'format domyślny rdzenia' },
  ...Object.values(StudioDocumentFormat).map((format) => ({
    wartosc: format,
    etykieta: format,
  })),
];

export interface PolaCyfryzacji {
  /** Kontrolki wsadu w kolejności osadzenia. */
  wsad: readonly HTMLElement[];
  /** Kontrolki nastaw rozpoznawania w kolejności osadzenia. */
  nastawy: readonly HTMLElement[];
  /** Kontrolki przyjęcia wyniku do edytora. */
  przyjecie: readonly HTMLElement[];
  /** Przycisk „Dołóż do kolejki". */
  doloz: HTMLButtonElement;
  /** Przycisk „Rozpakuj archiwum do kolejki rdzenia". */
  rozpakuj: HTMLButtonElement;
  /** Przycisk „Odczytaj urządzenia wejściowe". */
  odczytajUrzadzenia: HTMLButtonElement;
  /** Wskazanie materiału wpisane przez Operatora; puste znaczy brak wskazania. */
  wskazanie(): string;
  /** Czy wskazanie jest zasobem magazynu; `false` znaczy ścieżkę rdzenia. */
  czyZasob(): boolean;
  /** Ścieżka archiwum wsadu; pusta znaczy brak wskazania. */
  archiwum(): string;
  /** Czyści pole wskazania po dołożeniu pozycji do kolejki. */
  wyczyscWskazanie(): void;
  /** Nastawy rozpoznawania wspólne pozycjom wsadu. */
  ustawienia(): StudioRecognitionSettings;
  /** Tytuł dokumentu zakładanego z wyniku; pusty zostawia rdzeniowi wybór. */
  tytul(): string;
  /** Format dokumentu zakładanego z wyniku; pusty zostawia rdzeniowi wybór. */
  format(): string;
  /** Wstawia urządzenia oddane przez rdzeń wraz ze zdaniem o pustym wykazie. */
  ustawUrzadzenia(urzadzenia: readonly StudioInputDevice[]): void;
}

export function utworzPolaCyfryzacji(): PolaCyfryzacji {
  const zrodlo = poleWyboru(
    {
      etykieta: 'Źródło materiału',
      opis:
        'Kolejka rdzenia przyjmuje ścieżki i zasoby osobnymi polami; wybór ustala, które pole ' +
        'poleci do rdzenia. Wiele wskazań rozdziela się przecinkiem — wsad wchodzi jednym żądaniem.',
    },
    ZRODLA_MATERIALU,
  );
  zrodlo.element.append(
    utworzDymekObjasnienia(`${DROGI_CYFRYZACJI.zalaczniki} ${DROGI_CYFRYZACJI.library}`, {
      powloka: 'ms-dymek',
      znak: 'ms-dymek__znak',
    }),
  );

  const wskazanie = poleTekstowe({
    etykieta: 'Wskazanie materiału',
    podpowiedz: 'identyfikator zasobu albo ścieżka; wiele — po przecinku',
    opis: 'Klient dysku nie czyta ani nie zapisuje — podaje wskazanie, a treść wciąga rdzeń.',
  });

  const doloz = przycisk('Dołóż do kolejki rdzenia', 'dn-btn dn-btn--sm dn-btn--atrament');
  doloz.dataset['czynnosc'] = 'doloz';
  doloz.title =
    'studio.ingest.queue.add — kolejka stoi po stronie rdzenia i przeżywa odświeżenie okna.';

  const archiwum = poleTekstowe({
    etykieta: 'Archiwum wsadu (ścieżka rdzenia)',
    podpowiedz: 'ścieżka archiwum do rozpakowania',
    opis: SKLADNIKI_PAKIETU_SERWERA.archiwum,
  });

  const rozpakuj = przycisk('Rozpakuj archiwum do kolejki', 'dn-btn dn-btn--sm dn-btn--zarys');
  rozpakuj.dataset['czynnosc'] = 'rozpakuj';

  const urzadzenia = poleWyboru(
    {
      etykieta: 'Urządzenie wejściowe maszyny rdzenia',
      opis: SKLADNIKI_PAKIETU_SERWERA.urzadzenie,
    },
    [{ wartosc: '', etykieta: 'wykaz nieodczytany' }],
  );

  const odczytajUrzadzenia = przycisk(
    'Odczytaj urządzenia wejściowe',
    'dn-btn dn-btn--sm dn-btn--zarys',
  );
  odczytajUrzadzenia.dataset['czynnosc'] = 'urzadzenia';
  odczytajUrzadzenia.title =
    'studio.ingest.device.list — skanery i kamery widziane przez rdzeń. Bez tego wykazu wybór ' +
    'urządzenia nie ma z czego powstać.';

  const oUrzadzeniach = document.createElement('p');
  oUrzadzeniach.className = 'dn-pole-opis';
  oUrzadzeniach.textContent =
    'Wykaz urządzeń nie był jeszcze odczytany — naciśnij „Odczytaj urządzenia wejściowe". Wykaz ' +
    'pusty znaczy maszynę rdzenia bez skanera i bez kamery, a nie odmowę.';

  /* ── Nastawy rozpoznawania ─────────────────────────────────────────────── */

  const silnik = poleWyboru(
    { etykieta: 'Silnik rozpoznawania', opis: NASTAWY_CYFRYZACJI.silnik },
    SILNIKI,
  );

  const jezyki = poleTekstowe({
    etykieta: 'Zestaw języków rozpoznawania',
    podpowiedz: 'np. pol, eng — rozdzielone przecinkiem',
    opis: NASTAWY_CYFRYZACJI.jezyki,
  });

  const prog = poleLiczbowa('Próg pewności rozpoznania', 'np. 0,8 — puste zostawia próg rdzenia');
  prog.kontrolka.step = '0.01';
  prog.element.append(
    utworzDymekObjasnienia(NASTAWY_CYFRYZACJI.prog, {
      powloka: 'ms-dymek',
      znak: 'ms-dymek__znak',
    }),
  );

  const skos = poleLogiczne({ etykieta: 'Prostuj skos przed rozpoznaniem' });
  const szum = poleLogiczne({ etykieta: 'Odszumiaj obraz przed rozpoznaniem' });
  const progowanie = poleLogiczne({ etykieta: 'Proguj obraz do dwóch poziomów' });
  const marginesy = poleLogiczne({ etykieta: 'Przycinaj marginesy przed rozpoznaniem' });
  const uklad = poleLogiczne({
    etykieta: 'Odtwarzaj układ: kolumny, tabele, nagłówki, stopki, przypisy',
    opis: NASTAWY_CYFRYZACJI.uklad,
  });

  const stronaOd = poleLiczbowa('Pierwsza strona zakresu', 'puste = od początku');
  const stronaDo = poleLiczbowa('Ostatnia strona zakresu', 'puste = do końca');

  const oCzyszczeniu = document.createElement('p');
  oCzyszczeniu.className = 'dn-pole-opis';
  oCzyszczeniu.textContent = NASTAWY_CYFRYZACJI.czyszczenie;

  /* ── Przyjęcie wyniku ──────────────────────────────────────────────────── */

  const tytul = poleTekstowe({
    etykieta: 'Tytuł dokumentu zakładanego z wyniku',
    podpowiedz: 'puste zostawia rdzeniowi wybór',
    opis: NASTAWY_CYFRYZACJI.przyjecie,
  });

  const format = poleWyboru(
    { etykieta: 'Format dokumentu zakładanego z wyniku' },
    FORMATY_DOKUMENTU,
  );

  return {
    wsad: [
      zrodlo.element,
      wskazanie.element,
      doloz,
      archiwum.element,
      rozpakuj,
      urzadzenia.element,
      odczytajUrzadzenia,
      oUrzadzeniach,
    ],
    nastawy: [
      silnik.element,
      jezyki.element,
      prog.element,
      skos.element,
      szum.element,
      progowanie.element,
      marginesy.element,
      uklad.element,
      stronaOd.element,
      stronaDo.element,
      oCzyszczeniu,
    ],
    przyjecie: [tytul.element, format.element],
    doloz,
    rozpakuj,
    odczytajUrzadzenia,
    wskazanie: () => wskazanie.kontrolka.value.trim(),
    czyZasob: () => zrodlo.kontrolka.value === 'zasob',
    archiwum: () => archiwum.kontrolka.value.trim(),
    tytul: () => tytul.kontrolka.value.trim(),
    format: () => format.kontrolka.value,

    wyczyscWskazanie() {
      wskazanie.kontrolka.value = '';
    },

    /**
     * Nastawy w kształcie kontraktu — pola puste NIE wchodzą.
     *
     * Brak pola znaczy „weź wartość z katalogu ustawień", a to rozstrzygnięcie
     * rdzenia, nie okna. Pusty napis w polu języków byłby wskazaniem zestawu
     * pustego i zmieniłby znaczenie żądania.
     */
    ustawienia() {
      const zestaw = jezyki.kontrolka.value
        .split(',')
        .map((czesc) => czesc.trim())
        .filter((czesc) => czesc !== '');
      const progPewnosci = ulamek(prog.kontrolka.value);
      const nastawy: StudioRecognitionSettings = {};
      if (silnik.kontrolka.value !== '') {
        nastawy.engine = silnik.kontrolka.value as StudioRecognitionSettings['engine'];
      }
      if (zestaw.length > 0) nastawy.languages = zestaw;
      if (progPewnosci > 0) nastawy.minConfidence = progPewnosci;
      if (skos.kontrolka.checked) nastawy.deskew = true;
      if (szum.kontrolka.checked) nastawy.denoise = true;
      if (progowanie.kontrolka.checked) nastawy.binarize = true;
      if (marginesy.kontrolka.checked) nastawy.trimMargins = true;
      if (uklad.kontrolka.checked) nastawy.detectLayout = true;
      if (liczba(stronaOd.kontrolka.value) > 0) nastawy.pageFrom = liczba(stronaOd.kontrolka.value);
      if (liczba(stronaDo.kontrolka.value) > 0) nastawy.pageTo = liczba(stronaDo.kontrolka.value);
      return nastawy;
    },

    ustawUrzadzenia(wykaz) {
      if (wykaz.length === 0) {
        ustawPozycje(urzadzenia.kontrolka, [
          { wartosc: '', etykieta: 'brak urządzenia na maszynie rdzenia' },
        ]);
        oUrzadzeniach.textContent =
          'Rdzeń nie widzi ani jednego skanera i ani jednej kamery. To odpowiedź, nie odmowa: ' +
          'urządzenie wejściowe należy do maszyny rdzenia i tam musi być podłączone. Materiał ' +
          'wskaż ścieżką albo zasobem magazynu.';
        return;
      }
      ustawPozycje(urzadzenia.kontrolka, [
        { wartosc: '', etykieta: 'urządzenie domyślne rdzenia' },
        ...wykaz.map((urzadzenie) => ({
          wartosc: urzadzenie.id,
          etykieta:
            `${urzadzenie.name} · ${urzadzenie.kind}` +
            (urzadzenie.hasFeeder === true ? ' · z podajnikiem' : '') +
            (urzadzenie.resolutions === undefined
              ? ''
              : ` · ${urzadzenie.resolutions.join('/')} punktów na cal`),
        })),
      ]);
      oUrzadzeniach.textContent =
        `Rdzeń widzi urządzeń wejściowych: ${wykaz.length}. Wybór idzie do pobrania obrazu; ` +
        'urządzenie stoi na maszynie rdzenia, nie na maszynie Operatora.';
    },
  };
}

/** Pole liczbowe z etykietą — wiersz pola z kontrolką przestawioną na liczbę. */
function poleLiczbowa(
  etykieta: string,
  podpowiedz: string,
): { element: HTMLElement; kontrolka: HTMLInputElement } {
  const pole = poleTekstowe({ etykieta, podpowiedz });
  pole.kontrolka.type = 'number';
  return pole;
}

/** Liczba całkowita z pola; wartość niepoprawna i pusta znaczą „bez ograniczenia". */
function liczba(wartosc: string): number {
  const odczytana = Number.parseInt(wartosc, 10);
  return Number.isFinite(odczytana) ? odczytana : 0;
}

/** Ułamek z pola; przecinek dziesiętny przyjmowany, bo tak pisze Operator. */
function ulamek(wartosc: string): number {
  const odczytana = Number.parseFloat(wartosc.replace(',', '.'));
  return Number.isFinite(odczytana) ? odczytana : 0;
}
