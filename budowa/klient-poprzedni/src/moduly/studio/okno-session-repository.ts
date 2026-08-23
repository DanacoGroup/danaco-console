import { Command, type StudioVersion } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import {
  poleTekstowe,
  poleWyboru,
  przycisk,
  przyciskBezKomendy,
  utworzWierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import {
  BRAK_FILTRU_AUTORA,
  POZYCJE_FILTRU,
  czyZapisSamoczynny,
  opiszZawezenie,
  pasujeDoAutora,
  przefiltrujHistorie,
  type FiltrHistorii,
} from './filtr-historii';
import { utworzOknoStudio } from './okno-studio';
import type { StanStudio } from './stan-studio';
import { utworzWierszWersji } from './wiersz-wersji';

/** Nagłówki kolumn tabeli wersji — zostają widoczne także w stanie pustym. */
const KOLUMNY = ['Wersja', 'Opis zmiany', 'Czas', 'Autor', 'Stan', 'Czynność'];

/**
 * Historia wersji — pod przyciskiem, nie w stałej kolumnie.
 *
 * ── Powierzchnia należy do dokumentu ────────────────────────────────────────
 * Rozstrzygnięcie Właściciela, podjęte dwukrotnie: panele wchodzą na żądanie
 * i schodzą, gdy nie są używane; stała kolumna może być wyłącznie trybem do
 * wyboru. Historia wersji stała dotąd wąską kolumną na dole modułu, więc zabierała
 * miejsce także wtedy, gdy nikt do niej nie zaglądał. Teraz `element` jest
 * WĄSKIM PASKIEM ze znacznikiem wersji i przyciskiem; cała tabela otwiera się
 * nakładką nad treścią i schodzi naciśnięciem albo `Escape`.
 *
 * ── Co zostaje bez zmiany ───────────────────────────────────────────────────
 * Repozytorium narasta i **niczego nie usuwa samo**: `studio.repository.restore`
 * przywraca wersję bez usuwania wersji nowszych, więc samo cofnięcie jest
 * odwracalne. „Usuń z wykazu" jest czynnością MIEJSCOWĄ i jawną — schowaniem
 * pozycji w tym oknie, nie usunięciem wersji w rdzeniu; okno mówi to wprost, żeby
 * Operator nie sądził, że stracił wersję.
 *
 * ── Autozapis idzie osobnym szeregiem ───────────────────────────────────────
 * Zapisy samoczynne są domyślnie UKRYTE, z przełącznikiem „pokaż także zapisy
 * samoczynne". Gdyby wchodziły do jednego wykazu z wersjami nazwanymi, historia
 * zasypałaby się w kilka minut. Rozróżnienie bierze się z pola `milestone`
 * (`studio.version.label.set`) i z etykiety — drugiego pojęcia okno nie zakłada.
 *
 * ── Gałęzie ─────────────────────────────────────────────────────────────────
 * Czynności „Rozgałęź" tu nie ma decyzją Właściciela. Komendy `studio.branch.*`
 * są zbudowane i pracują dalej — po prostu nie dostają okna w tej turze.
 */
export interface OknoSessionRepository {
  element: HTMLElement;
  /** Odczytuje historię wersji dokumentu czynnego. */
  wczytaj(): Promise<void>;
  odswiez(): void;
  /** Otwiera albo zamyka nakładkę historii — dla przycisku spoza tego pliku. */
  przestawWidocznosc(): void;
}

/**
 * Zaplecze czynności, których `ZrodloStudio` nie niesie.
 *
 * `ZrodloStudio` ma sześć komend i nie ma wśród nich etykietowania wersji,
 * eksportu historii ani paczki przekazania. Zaplecze jest **nieobowiązkowe**, bo
 * wołacz tego okna (`modul-studio.ts`) go dziś nie podaje, a plik ten nie należy
 * do tego odcinka prac. Bez zaplecza trzy czynności stoją jako brak NAZWANY,
 * wraz z komendą, która czeka w rdzeniu — zamiast jako przycisk, który milczy.
 */
export interface ZapleczeHistorii {
  wykonaj(komenda: Command, zadanie: Record<string, unknown>): Promise<{
    udany: boolean;
    blad?: { code?: string; message?: string };
  }>;
}

export function utworzOknoSessionRepository(
  stan: StanStudio,
  zaplecze?: ZapleczeHistorii,
): OknoSessionRepository {
  const rama = utworzOknoStudio({
    kod: 'studio.session-repository',
    tytul: 'Historia wersji',
    rola: 'zarządca',
    objasnienie:
      'Wersje narastają w toku sesji przy każdym zapisie. Przywrócenie wersji przestawia treść ' +
      'dokumentu i NIE usuwa wersji nowszych (studio.repository.restore), więc samo cofnięcie ' +
      'jest odwracalne. Historia otwiera się na żądanie i schodzi — powierzchnia należy do ' +
      'dokumentu, nie do paneli.',
  });

  const cialo = document.createElement('tbody');

  const glowa = document.createElement('thead');
  const wierszGlowy = document.createElement('tr');
  for (const kolumna of KOLUMNY) {
    const komorka = document.createElement('th');
    komorka.scope = 'col';
    komorka.textContent = kolumna;
    wierszGlowy.append(komorka);
  }
  glowa.append(wierszGlowy);

  const tabela = document.createElement('table');
  tabela.className = 'dn-tabela ms-repozytorium';
  tabela.append(glowa, cialo);

  const filtr = poleWyboru(
    { etykieta: 'Filtr historii', opis: BRAK_FILTRU_AUTORA },
    POZYCJE_FILTRU.map((pozycja) => ({ wartosc: pozycja.wartosc, etykieta: pozycja.etykieta })),
  );

  const filtrAutora = poleTekstowe({
    etykieta: 'Filtr autora',
    opis:
      'Zawęża wykaz do wersji autora, którego nazwa zawiera wpisany ciąg. Rdzeń oddaje autora ' +
      'w polu StudioVersion.author (uzytkownik albo model); wersje założone przed wprowadzeniem ' +
      'tego pola autora nie niosą i zostają widoczne, zamiast wypadać z wykazu po cichu.',
  });

  const zapisySamoczynne = document.createElement('input');
  zapisySamoczynne.type = 'checkbox';
  zapisySamoczynne.className = 'dn-przelacznik';
  zapisySamoczynne.setAttribute('aria-label', 'Pokaż także zapisy samoczynne');
  const etykietaSamoczynnych = document.createElement('label');
  etykietaSamoczynnych.className = 'ms-repozytorium__przelacznik';
  const napisSamoczynnych = document.createElement('span');
  napisSamoczynnych.textContent = 'Pokaż także zapisy samoczynne';
  napisSamoczynnych.title =
    'Autozapis idzie osobnym szeregiem: wersje samoczynne są domyślnie ukryte, żeby historia ' +
    'nie zasypała wersji nazwanych przez Operatora. Wersja kluczowa (milestone) i wersja ' +
    'z etykietą własną nigdy nie są liczone jako samoczynne.';
  etykietaSamoczynnych.append(zapisySamoczynne, napisSamoczynnych);

  const zawezenie = document.createElement('p');
  zawezenie.className = 'dn-pole-opis ms-repozytorium__zawezenie';

  const odswiezWykaz = przycisk('Odczytaj historię', 'dn-btn dn-btn--sm dn-btn--zarys');
  // Znacznik czynności jak w pozostałych oknach modułu: nazywa ster po roli,
  // dzięki czemu zmiana etykiety widocznej dla użytkownika niczego nie psuje.
  odswiezWykaz.dataset['czynnosc'] = 'odczytaj';
  const odpowiedz = utworzWierszOdpowiedzi();

  /* ── Czynności zbiorcze ──────────────────────────────────────────────────── */

  const pasZbiorczy = document.createElement('div');
  pasZbiorczy.className = 'ms-repozytorium__zbiorcze';
  pasZbiorczy.append(
    czynnoscZaplecza(
      'Eksportuj całą historię',
      'eksport-historii',
      Command.StudioRepositoryExport,
      () => ({ documentId: stan.dokument()?.id ?? '' }),
      'Zakłada archiwum całej historii wraz z manifestem — komenda studio.repository.export.',
    ),
    czynnoscZaplecza(
      'Paczka redakcyjna przekazania',
      'paczka-redakcyjna',
      Command.StudioPackageExport,
      () => ({ documentId: stan.dokument()?.id ?? '' }),
      'Zakłada paczkę przekazania: dokument, historia, raport zmian i adnotacje — komenda ' +
        'studio.package.export.',
    ),
  );

  /**
   * Przycisk czynności, która potrzebuje zaplecza.
   *
   * Bez zaplecza oddaje brak NAZWANY: mówi, że komenda w rdzeniu jest, i czego
   * brakuje po stronie klienta. Przycisk, który po naciśnięciu milczy, byłby tu
   * gorszy od braku — Operator sądziłby, że paczka powstała.
   */
  function czynnoscZaplecza(
    nazwa: string,
    kod: string,
    komenda: Command,
    zadanie: () => Record<string, unknown>,
    objasnienie: string,
  ): HTMLElement {
    if (zaplecze === undefined) {
      return przyciskBezKomendy(
        nazwa,
        `${objasnienie} Komenda stoi w rdzeniu i jest zbudowana, ale okno historii nie ma dziś ` +
          'drogi do niej: ZrodloStudio niesie sześć komend i tej wśród nich nie ma, a wołacz okna ' +
          '(modul-studio.ts) nie podaje zaplecza. Brak jest po stronie klienta, nie w rdzeniu — ' +
          'wpisany w sprawozdanie jako jedno wywołanie do domknięcia.',
      );
    }
    const element = przycisk(nazwa, 'dn-btn dn-btn--sm dn-btn--zarys');
    element.dataset['czynnosc'] = kod;
    element.title = objasnienie;
    element.addEventListener('click', () => void wykonajZapleczem(nazwa, komenda, zadanie()));
    return element;
  }

  async function wykonajZapleczem(
    nazwa: string,
    komenda: Command,
    zadanie: Record<string, unknown>,
  ): Promise<void> {
    if (zaplecze === undefined) return;
    if (zadanie['documentId'] === '') {
      odpowiedz.pokaz(`${nazwa}: najpierw wczytaj dokument — czynność dotyczy jego historii.`, false);
      return;
    }
    odpowiedz.pokaz(`${nazwa} w toku…`, true);
    const wynik = await zaplecze.wykonaj(komenda, zadanie);
    odpowiedz.pokaz(
      wynik.udany
        ? `${nazwa}: rdzeń przyjął czynność.`
        : opisOdmowy(nazwa, wynik.blad?.code, wynik.blad?.message),
      wynik.udany,
    );
    if (wynik.udany) await wczytaj();
  }

  rama.pasek.append(odswiezWykaz);
  rama.stan.tresc.append(
    filtr.element,
    filtrAutora.element,
    etykietaSamoczynnych,
    zawezenie,
    tabela,
    pasZbiorczy,
    odpowiedz.element,
  );

  /* ── Nakładka i jej wyzwalacz ────────────────────────────────────────────── */

  const nakladka = document.createElement('div');
  nakladka.className = 'ms-historia__nakladka';
  nakladka.hidden = true;
  nakladka.append(rama.element);

  const wyzwalacz = document.createElement('button');
  wyzwalacz.type = 'button';
  wyzwalacz.className = 'dn-btn dn-btn--sm dn-btn--zarys ms-historia__wyzwalacz';
  wyzwalacz.dataset['czynnosc'] = 'historia';
  wyzwalacz.setAttribute('aria-expanded', 'false');
  wyzwalacz.title =
    'Otwiera historię wersji nakładką nad treścią. Zamyka ją drugie naciśnięcie albo Escape — ' +
    'historia nie zajmuje miejsca, kiedy nie jest używana.';
  wyzwalacz.addEventListener('click', () => przestawWidocznosc());

  const element = document.createElement('div');
  element.className = 'ms-historia';
  element.append(wyzwalacz, nakladka);

  element.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key === 'Escape' && !nakladka.hidden) przestawWidocznosc();
  });

  function przestawWidocznosc(): void {
    const otwarta = nakladka.hidden;
    nakladka.hidden = !otwarta;
    wyzwalacz.setAttribute('aria-expanded', otwarta ? 'true' : 'false');
    opiszWyzwalacz();
    if (otwarta && znacznikDokumentu() !== znacznikOdczytu) void wczytaj();
  }

  function opiszWyzwalacz(): void {
    const ile = stan.wersje().length;
    const otwarta = !nakladka.hidden;
    wyzwalacz.textContent =
      ile === 0
        ? `Historia wersji ${otwarta ? '▴' : '▾'}`
        : `Historia wersji: ${ile} ${otwarta ? '▴' : '▾'}`;
  }

  filtr.kontrolka.addEventListener('change', () => odswiez());
  filtrAutora.kontrolka.addEventListener('input', () => przerysujWiersze());
  zapisySamoczynne.addEventListener('change', () => przerysujWiersze());

  /* ── Czynności na pozycji ────────────────────────────────────────────────── */

  async function przywroc(wersja: StudioVersion, oznacz: (powod: string) => void): Promise<void> {
    const dokument = stan.dokument();
    if (dokument === null) return;
    odpowiedz.pokaz(`Przywracanie wersji ${wersja.id}…`, true);
    const wynik = await stan.zrodlo.przywroc({ documentId: dokument.id, versionId: wersja.id });
    if (!wynik.udany || wynik.wynik === undefined) {
      const powod = opisOdmowy('Przywrócenie wersji', wynik.blad?.code, wynik.blad?.message);
      oznacz(powod);
      odpowiedz.pokaz(powod, false);
      return;
    }
    stan.wchlon(wynik.wynik.document);
    odpowiedz.pokaz(
      `Wersja ${wersja.id} przywrócona jako stan bieżący. Wersje nowsze POZOSTAŁY w historii, ` +
        'więc to cofnięcie samo da się cofnąć.',
      true,
    );
    await wczytaj();
  }

  /**
   * Podgląd wersji — treść wersji obok bieżącej, bez jej przywracania.
   *
   * Idzie tą samą drogą co porównanie (`studio.diff.compare` przez parę
   * porównania), bo komendy „oddaj treść wersji" kontrakt nie niesie: różnica
   * wobec bieżącej JEST podglądem tego, co ta wersja niosła, i nie wymaga
   * zakładania drugiej drogi do rdzenia.
   */
  function podejrzyj(wersja: StudioVersion): void {
    stan.ustawParePorownania({ odniesienie: wersja.id, porownywana: '' });
    odpowiedz.pokaz(
      `Wersja ${wersja.id} wskazana do podglądu. Okno pracy pokaże ją trybem „Różnica na treści" — ` +
        'komendy oddającej samą treść wersji kontrakt nie ma, a różnica wobec bieżącej mówi ' +
        'dokładnie to samo i więcej.',
      true,
    );
  }

  /** Pozycje schowane miejscowo — decyzją Operatora, nie regułą okna. */
  const schowane = new Set<string>();

  function przerysujWiersze(): void {
    const wszystkie = stan.wersje();
    const wybrany = filtr.kontrolka.value as FiltrHistorii;
    const poFiltrze = przefiltrujHistorie(wszystkie, wybrany);
    const szukanyAutor = filtrAutora.kontrolka.value.trim();
    const widoczne = poFiltrze.filter(
      (wersja) =>
        !schowane.has(wersja.id) &&
        pasujeDoAutora(wersja, szukanyAutor) &&
        (zapisySamoczynne.checked || !czyZapisSamoczynny(wersja)),
    );
    zawezenie.textContent =
      `${opiszZawezenie(wszystkie.length, widoczne.length, wybrany)}` +
      (schowane.size === 0 ? '' : ` Schowanych miejscowo: ${schowane.size}.`) +
      (zapisySamoczynne.checked ? '' : ' Zapisy samoczynne są ukryte — przełącznik wyżej je pokazuje.');

    cialo.replaceChildren(
      ...widoczne.map((wersja) => {
        const wiersz = utworzWierszWersji(wersja, {
          biezaca: stan.dokument()?.versionId === wersja.id,
          zapisSamoczynny: czyZapisSamoczynny(wersja),
        });
        wiersz.podejrzyj.addEventListener('click', () => podejrzyj(wersja));
        wiersz.przywroc.addEventListener('click', () => {
          void przywroc(wersja, wiersz.oznaczNiedostepna);
        });
        wiersz.porownaj.addEventListener('click', () => wskazDoPorownania(wersja));
        wiersz.etykieta.addEventListener('click', () => void nadajEtykiete(wersja));
        wiersz.eksportuj.addEventListener('click', () =>
          void wykonajZapleczem('Eksport wersji', Command.StudioRepositoryExport, {
            documentId: wersja.documentId,
            versionId: wersja.id,
          }),
        );
        wiersz.odwolanie.addEventListener('click', () => odwolajSie(wersja));
        wiersz.usun.addEventListener('click', () => {
          schowane.add(wersja.id);
          odpowiedz.pokaz(
            `Pozycja ${wersja.id} schowana MIEJSCOWO — w tym oknie i tylko tutaj. Wersja stoi ` +
              'dalej w rdzeniu i wróci do wykazu po ponownym odczycie historii; komendy usuwającej ' +
              'wersję kontrakt nie niesie i historia nic nie usuwa bez decyzji Operatora.',
            true,
          );
          przerysujWiersze();
        });
        return wiersz.element;
      }),
    );
    opiszWyzwalacz();
  }

  /** Nadaje wersji nazwę własną i oznaczenie kluczowej — `studio.version.label.set`. */
  async function nadajEtykiete(wersja: StudioVersion): Promise<void> {
    if (zaplecze === undefined) {
      odpowiedz.pokaz(
        'Etykieta wersji idzie komendą studio.version.label.set, która jest w rdzeniu zbudowana. ' +
          'Okno historii nie ma dziś do niej drogi: ZrodloStudio jej nie niesie, a zaplecze nie ' +
          'jest podane przez wołacza okna. Brak po stronie klienta, wpisany w sprawozdanie.',
        false,
      );
      return;
    }
    const nazwa = globalThis.prompt?.(
      `Nazwa własna wersji ${wersja.id} (pusta zdejmuje etykietę):`,
      wersja.label ?? '',
    );
    if (nazwa === null || nazwa === undefined) return;
    await wykonajZapleczem('Etykieta wersji', Command.StudioVersionLabelSet, {
      versionId: wersja.id,
      label: nazwa,
      milestone: nazwa !== '',
    });
  }

  /**
   * Odwołanie do wersji — jej identyfikator do wklejenia w pismo albo w rozmowę.
   *
   * Bez schowka systemowego: odwołanie idzie do wiersza odpowiedzi, skąd da się je
   * przeczytać i zaznaczyć. Schowek ma własną rodzinę komend i własnego wykonawcę,
   * więc drugiej drogi do niego to okno nie zakłada.
   */
  function odwolajSie(wersja: StudioVersion): void {
    const czas = new Date(wersja.createdAt).toLocaleString('pl');
    odpowiedz.pokaz(
      `Odwołanie do wersji: ${wersja.id} · dokument ${wersja.documentId} · ${czas}` +
        `${wersja.label === undefined || wersja.label === '' ? '' : ` · „${wersja.label}"`}`,
      true,
    );
  }

  /**
   * Wskazuje wersję jako odniesienie porównania.
   *
   * Strona druga zostaje pusta z zamysłem: `studio.diff.compare` czyta brak
   * wersji porównywanej jako zgodę na porównanie z propozycją zmiany albo
   * z wersją bieżącą. Wpisanie tam czegokolwiek na siłę odbierałoby Operatorowi
   * najczęstszy przypadek — „ta wersja wobec tego, co mam teraz".
   */
  function wskazDoPorownania(wersja: StudioVersion): void {
    stan.ustawParePorownania({ odniesienie: wersja.id, porownywana: '' });
    odpowiedz.pokaz(
      `Wersja ${wersja.id} wskazana jako odniesienie porównania. Różnicę liczy okno pracy ` +
        'z dokumentem — czeka z wypełnionymi polami, żeby dało się je jeszcze poprawić.',
      true,
    );
  }

  /**
   * Znacznik odczytu: dokument i jego wersja bieżąca w chwili ostatniego
   * czytania historii. Pusty łańcuch znaczy „nie czytano jeszcze nic".
   */
  let znacznikOdczytu = '';

  /** Dokument i jego wersja bieżąca — para, po której poznaje się zdezaktualizowanie wykazu. */
  function znacznikDokumentu(): string {
    const dokument = stan.dokument();
    if (dokument === null) return '';
    return `${dokument.id}:${dokument.versionId ?? ''}`;
  }

  async function wczytaj(): Promise<void> {
    const dokument = stan.dokument();
    if (dokument === null) {
      rama.stan.puste('Historia bez wskazanego dokumentu', BEZ_DOKUMENTU);
      return;
    }
    // Znacznik przestawiany przed wywołaniem: odmowa nie może wprawić odświeżania
    // w pętlę ponowień czegoś, co rdzeń przed chwilą odrzucił.
    znacznikOdczytu = znacznikDokumentu();
    rama.stan.ladowanie('Odczyt historii wersji z repozytorium sesji…');
    const wynik = await stan.zrodlo.wersje({ documentId: dokument.id });
    if (!wynik.udany || wynik.wynik === undefined) {
      rama.stan.blad(opisOdmowy('Odczyt historii wersji', wynik.blad?.code, wynik.blad?.message));
      return;
    }
    stan.ustawWersje(wynik.wynik.versions);
    // `ustawWersje` odświeżyło okna w fazie „ładowanie", więc stan pusty jeszcze
    // nie zapadł. Zdejmujemy fazę i odświeżamy raz jeszcze — inaczej pusta
    // historia kończyłaby się tabelą bez wiersza i bez jednego zdania powodu.
    rama.stan.gotowe();
    odswiez();
  }

  odswiezWykaz.addEventListener('click', () => void wczytaj());

  function odswiez(): void {
    przerysujWiersze();
    if (rama.stan.faza() === 'ladowanie' || rama.stan.faza() === 'blad') return;
    if (stan.dokument() === null) {
      rama.stan.puste('Historia bez wskazanego dokumentu', BEZ_DOKUMENTU);
      return;
    }
    // Odczyt tylko przy otwartej nakładce: historia zamknięta nie jest powodem,
    // żeby wołać rdzeń przy każdym naciśnięciu klawisza w dokumencie.
    if (!nakladka.hidden && znacznikDokumentu() !== znacznikOdczytu) {
      void wczytaj();
      return;
    }
    if (stan.wersje().length === 0) {
      rama.stan.puste('Dokument jeszcze bez wersji', BEZ_WERSJI);
      return;
    }
    rama.stan.gotowe();
  }

  opiszWyzwalacz();

  return { element, wczytaj, odswiez, przestawWidocznosc };
}

const BEZ_DOKUMENTU =
  'Historia prowadzi wersje JEDNEGO dokumentu sesji i pozwala wrócić do wcześniejszej bez ' +
  'usuwania nowszych. Wczytaj dokument w oknie pracy — wykaz weźmie jego wersje z komendy ' +
  'studio.repository.list.';

const BEZ_WERSJI =
  'Rdzeń oddał wykaz pusty: ten dokument nie ma jeszcze ani jednej wersji. Wersje narastają ' +
  'przy każdym zapisie — pierwszy zapis zakłada wersję pierwszą i od niej zaczyna się historia, ' +
  'do której można wrócić.';
