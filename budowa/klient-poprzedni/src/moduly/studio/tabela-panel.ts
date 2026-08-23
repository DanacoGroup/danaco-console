import {
  StudioTableConvert,
  StudioTableStructureOp,
  StudioTextAlign,
  StudioVerticalAlign,
  type StudioActionBalance,
  type StudioDocumentTable,
} from '../../../../shared/contract';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import {
  poleLiczbowe,
  poleTekstowe,
  przycisk,
  utworzWierszOdpowiedzi,
  wybor,
} from '../../modele/kontrolki-formularza';
import type { Wynik } from '../../protokol/kanal';
import type { StanStudio } from './stan-studio';
import type { TabelaZrodlo } from './tabela-zrodlo';
import { wstawieniaOpiszBilans } from './zrodlo-wstawien-studio';

/**
 * Warsztat tabel — nakładka na żądanie, nie stała kolumna.
 *
 * ── Wejście przez wskazanie rozmiaru siatką ─────────────────────────────────
 * Zamówione wprost: tabela wstawia się przez wskazanie rozmiaru siatką, a nie
 * przez wpisanie dwóch liczb. Siatka jest tu pierwsza, a pola liczbowe stoją
 * obok jako droga druga — dla tabeli większej niż siatka i dla pracy
 * z klawiatury. Oba wejścia prowadzą do tej samej komendy
 * `studio.table.insert`, więc nie ma dwóch zachowań.
 *
 * ── Bilans zamiast ciszy ────────────────────────────────────────────────────
 * Każda czynność tabeli oddaje bilans i panel go wypisuje ZAWSZE, nie tylko przy
 * pominięciu. Scalenie komórek w zablokowanym fragmencie wraca odpowiedzią
 * pomyślną z pominięciem w bilansie — bez wypisania bilansu wyglądałoby to na
 * scalenie wykonane. Wpis dziennika (`actionId`) idzie do tego samego zdania, bo
 * bez niego Operator nie wie, co ma cofnąć.
 *
 * ── Czego panel nie liczy sam ───────────────────────────────────────────────
 * Ani szerokości kolumn, ani wyniku sortowania, ani zamiany tekstu na tabelę.
 * Wszystko to robi rdzeń; panel składa żądanie, czyta tabelę z odpowiedzi
 * i pokazuje szerokości POLICZONE, a nie założone.
 */
export interface TabelaPanel {
  element: HTMLElement;
  /** Otwiera albo zamyka warsztat — dla przycisku spoza tego pliku. */
  przestawWidocznosc(): void;
  /** Odczytuje tabele dokumentu czynnego; woła to okno po zmianie dokumentu. */
  odswiez(): void;
}

/** Największy rozmiar wskazywany siatką; większe tabele idą polami liczbowymi. */
const SIATKA_WIERSZY = 8;
const SIATKA_KOLUMN = 10;

/** Sześć czynności na budowie tabeli wraz z nazwą widoczną dla Operatora. */
const CZYNNOSCI_BUDOWY: readonly { operacja: StudioTableStructureOp; nazwa: string }[] = [
  { operacja: StudioTableStructureOp.InsertRow, nazwa: 'Wstaw wiersz' },
  { operacja: StudioTableStructureOp.DeleteRow, nazwa: 'Usuń wiersz' },
  { operacja: StudioTableStructureOp.InsertColumn, nazwa: 'Wstaw kolumnę' },
  { operacja: StudioTableStructureOp.DeleteColumn, nazwa: 'Usuń kolumnę' },
  { operacja: StudioTableStructureOp.MergeCells, nazwa: 'Scal komórki' },
  { operacja: StudioTableStructureOp.SplitCell, nazwa: 'Podziel komórkę' },
];

export function utworzTabelaPanel(stan: StanStudio, zrodlo: TabelaZrodlo): TabelaPanel {
  const odpowiedz = utworzWierszOdpowiedzi();

  /** Tabela wskazana — przedmiot budowy, postaci i sortowania. */
  let wskazana = '';
  let tabele: readonly StudioDocumentTable[] = [];

  /* ── Wskazanie rozmiaru siatką ─────────────────────────────────────────── */

  const siatka = document.createElement('div');
  siatka.className = 'ms-tabela__siatka';
  siatka.setAttribute('role', 'grid');
  siatka.setAttribute('aria-label', 'Wskazanie rozmiaru tabeli siatką');

  const podpisSiatki = document.createElement('p');
  podpisSiatki.className = 'dn-pole-opis';
  podpisSiatki.textContent =
    'Wskaż rozmiar tabeli w siatce albo podaj go liczbami niżej — obie drogi wołają tę samą ' +
    'komendę studio.table.insert.';

  for (let wiersz = 1; wiersz <= SIATKA_WIERSZY; wiersz += 1) {
    for (let kolumna = 1; kolumna <= SIATKA_KOLUMN; kolumna += 1) {
      const komorka = document.createElement('button');
      komorka.type = 'button';
      komorka.className = 'ms-tabela__komorka';
      komorka.dataset['wiersze'] = String(wiersz);
      komorka.dataset['kolumny'] = String(kolumna);
      komorka.setAttribute('aria-label', `Tabela ${wiersz} na ${kolumna}`);
      komorka.title = `${wiersz} wierszy na ${kolumna} kolumn`;
      komorka.addEventListener('mouseenter', () => oznaczSiatke(wiersz, kolumna));
      komorka.addEventListener('focus', () => oznaczSiatke(wiersz, kolumna));
      komorka.addEventListener('click', () => void wstawTabele(wiersz, kolumna));
      siatka.append(komorka);
    }
  }

  /** Podświetla prostokąt siatki do wskazanej komórki — podgląd przed wstawieniem. */
  function oznaczSiatke(wiersze: number, kolumny: number): void {
    for (const komorka of siatka.children) {
      if (!(komorka instanceof HTMLElement)) continue;
      const wiersz = Number(komorka.dataset['wiersze'] ?? '0');
      const kolumna = Number(komorka.dataset['kolumny'] ?? '0');
      komorka.dataset['objeta'] = wiersz <= wiersze && kolumna <= kolumny ? 'tak' : 'nie';
    }
    podpisSiatki.textContent = `Tabela ${wiersze} × ${kolumny} — naciśnij, żeby ją wstawić.`;
  }

  /* ── Nastawy wstawienia ────────────────────────────────────────────────── */

  const wierszeRecznie = poleLiczbowe('Liczba wierszy', 'np. 12');
  const kolumnyRecznie = poleLiczbowe('Liczba kolumn', 'np. 4');
  const styl = poleTekstowe({
    etykieta: 'Styl tabeli nazwany albo szablon gotowy',
    podpowiedz: 'puste zostawia rdzeniowi styl domyślny',
  });
  const wierszeNaglowka = poleLiczbowe('Ile wierszy początkowych jest nagłówkiem', 'np. 1');
  const powtarzajNaglowek = document.createElement('input');
  powtarzajNaglowek.type = 'checkbox';
  powtarzajNaglowek.className = 'dn-przelacznik';
  powtarzajNaglowek.setAttribute(
    'aria-label',
    'Powtarzaj wiersz nagłówkowy na kolejnych stronach',
  );
  const etykietaPowtarzania = document.createElement('label');
  etykietaPowtarzania.className = 'ms-tabela__zawezenie';
  etykietaPowtarzania.append(
    powtarzajNaglowek,
    document.createTextNode('powtarzaj wiersz nagłówkowy na kolejnych stronach'),
  );

  const wstawRecznie = przycisk('Wstaw tabelę o podanym rozmiarze', 'dn-btn dn-btn--sm dn-btn--sygnal');
  wstawRecznie.dataset['czynnosc'] = 'wstaw-tabele';
  wstawRecznie.addEventListener('click', () => {
    void wstawTabele(liczba(wierszeRecznie.value), liczba(kolumnyRecznie.value));
  });

  /* ── Wykaz tabel dokumentu ─────────────────────────────────────────────── */

  const wykaz = document.createElement('ul');
  wykaz.className = 'ms-tabela__wykaz';
  wykaz.setAttribute('aria-label', 'Tabele dokumentu');

  /* ── Budowa tabeli wskazanej ───────────────────────────────────────────── */

  const wiersz = poleLiczbowe('Wiersz liczony od zera', '0');
  const kolumna = poleLiczbowe('Kolumna liczona od zera', '0');
  const ile = poleLiczbowe('Ile wierszy albo kolumn objąć', 'puste znaczy jeden');
  const obejmijWierszy = poleLiczbowe('Ile wierszy objąć scaleniem albo na ile podzielić', '');
  const obejmijKolumn = poleLiczbowe('Ile kolumn objąć scaleniem albo na ile podzielić', '');
  const przedWskazanym = document.createElement('input');
  przedWskazanym.type = 'checkbox';
  przedWskazanym.className = 'dn-przelacznik';
  przedWskazanym.setAttribute('aria-label', 'Wstaw przed wskazanym wierszem albo kolumną');
  const etykietaPrzed = document.createElement('label');
  etykietaPrzed.className = 'ms-tabela__zawezenie';
  etykietaPrzed.append(przedWskazanym, document.createTextNode('wstaw PRZED wskazanym'));

  const rzadBudowy = document.createElement('div');
  rzadBudowy.className = 'ms-tabela__rzad';
  for (const czynnosc of CZYNNOSCI_BUDOWY) {
    const kontrolka = przycisk(czynnosc.nazwa, 'dn-btn dn-btn--sm dn-btn--atrament');
    kontrolka.dataset['budowa'] = czynnosc.operacja;
    kontrolka.addEventListener('click', () => void zmienBudowe(czynnosc.operacja, czynnosc.nazwa));
    rzadBudowy.append(kontrolka);
  }

  /* ── Postać tabeli ─────────────────────────────────────────────────────── */

  const szerokosciKolumn = poleTekstowe({
    etykieta: 'Szerokości kolumn w milimetrach',
    podpowiedz: 'np. 40, 30, 30 — puste zostawia szerokości policzone przez rdzeń',
  });
  const szerokoscTabeli = poleLiczbowe('Szerokość całej tabeli w milimetrach', '');
  const cieniowanie = poleTekstowe({
    etykieta: 'Cieniowanie zapisem szesnastkowym',
    podpowiedz: 'np. #eeeeee',
  });
  const podpisTabeli = poleTekstowe({ etykieta: 'Podpis tabeli', podpowiedz: '' });
  const wyrownaniePionowe = wybor('Wyrównanie pionowe w komórce', [
    ['', 'bez zmiany'],
    [StudioVerticalAlign.Top, 'do góry'],
    [StudioVerticalAlign.Middle, 'do środka'],
    [StudioVerticalAlign.Bottom, 'do dołu'],
  ]);
  const wyrownaniePoziome = wybor('Wyrównanie w komórce albo tabeli na stronie', [
    ['', 'bez zmiany'],
    [StudioTextAlign.Left, 'do lewej'],
    [StudioTextAlign.Center, 'do środka'],
    [StudioTextAlign.Right, 'do prawej'],
    [StudioTextAlign.Justify, 'justowanie'],
  ]);

  const ustawPostac = przycisk('Ustaw postać tabeli', 'dn-btn dn-btn--sm dn-btn--atrament');
  ustawPostac.dataset['czynnosc'] = 'postac-tabeli';
  ustawPostac.addEventListener('click', () => void przestawPostac());

  /* ── Sortowanie ────────────────────────────────────────────────────────── */

  const kolumnaSortowania = poleLiczbowe('Kolumna sortowania liczona od zera', '0');
  const malejaco = document.createElement('input');
  malejaco.type = 'checkbox';
  malejaco.className = 'dn-przelacznik';
  malejaco.setAttribute('aria-label', 'Sortuj malejąco');
  const jakLiczby = document.createElement('input');
  jakLiczby.type = 'checkbox';
  jakLiczby.className = 'dn-przelacznik';
  jakLiczby.setAttribute('aria-label', 'Porównuj jako liczby');
  const etykietaMalejaco = document.createElement('label');
  etykietaMalejaco.className = 'ms-tabela__zawezenie';
  etykietaMalejaco.append(malejaco, document.createTextNode('malejąco'));
  const etykietaLiczb = document.createElement('label');
  etykietaLiczb.className = 'ms-tabela__zawezenie';
  etykietaLiczb.append(jakLiczby, document.createTextNode('porównuj jako liczby'));

  const sortuj = przycisk('Sortuj zawartość tabeli', 'dn-btn dn-btn--sm dn-btn--atrament');
  sortuj.dataset['czynnosc'] = 'sortuj-tabele';
  sortuj.addEventListener('click', () => void sortujTabele());

  /* ── Zamiana tekstu i tabeli ───────────────────────────────────────────── */

  const rozdzielnik = poleTekstowe({
    etykieta: 'Znak rozdzielający kolumny',
    podpowiedz: 'puste znaczy tabulator',
  });

  const naTabele = przycisk('Zamień zaznaczony tekst na tabelę', 'dn-btn dn-btn--sm dn-btn--zarys');
  naTabele.dataset['czynnosc'] = 'tekst-na-tabele';
  naTabele.addEventListener('click', () => void zamienTekstNaTabele());

  const naTekst = przycisk('Zamień wskazaną tabelę na tekst', 'dn-btn dn-btn--sm dn-btn--zarys');
  naTekst.dataset['czynnosc'] = 'tabela-na-tekst';
  naTekst.addEventListener('click', () => void zamienTabeleNaTekst());

  /* ── Czynności ─────────────────────────────────────────────────────────── */

  /** Dokument czynny albo odmowa nazwana; tabela bez dokumentu nie ma gdzie stanąć. */
  function idDokumentu(): string | null {
    const dokument = stan.dokument();
    if (dokument === null) {
      odpowiedz.pokaz(BRAK_DOKUMENTU, false);
      return null;
    }
    return dokument.id;
  }

  /** Miejsce wstawienia: koniec zaznaczenia albo koniec treści roboczej. */
  function miejsceKursora(): number {
    const zakres = stan.zaznaczenie();
    return zakres === null ? stan.trescRobocza().length : zakres.koniec;
  }

  /** Wskazuje tabelę albo mówi, że wskazania brakuje. */
  function idTabeli(): string | null {
    if (wskazana === '') {
      odpowiedz.pokaz(BRAK_WSKAZANEJ, false);
      return null;
    }
    return wskazana;
  }

  /**
   * Jedno miejsce, w którym odpowiedź rdzenia zamienia się w zdanie.
   *
   * Bilans wypisuje się bezwarunkowo, także przy pełnym powodzeniu: „zmienionych
   * miejsc 1, pominiętych 0" jest zdaniem prawdziwym, a wypisywanie bilansu tylko
   * przy pominięciu uczyłoby Operatora, że brak bilansu znaczy „wszystko weszło".
   */
  function opiszSkutek(
    nazwa: string,
    bilans: StudioActionBalance,
    idCzynnosci: string | undefined,
    dopisek: string,
  ): void {
    const cofniecie =
      idCzynnosci === undefined || idCzynnosci === ''
        ? 'Rdzeń nie oddał wpisu dziennika — tej czynności nie da się cofnąć pojedynczo.'
        : `Wpis dziennika do cofnięcia pojedynczego: ${idCzynnosci}.`;
    odpowiedz.pokaz(
      `${nazwa}. ${wstawieniaOpiszBilans(bilans)} ${dopisek} ${cofniecie}`.replace(/\s+/g, ' '),
      bilans.applied > 0,
    );
  }

  async function wstawTabele(wiersze: number, kolumny: number): Promise<void> {
    const dokument = idDokumentu();
    if (dokument === null) return;
    if (wiersze <= 0 || kolumny <= 0) {
      odpowiedz.pokaz(BRAK_ROZMIARU, false);
      return;
    }
    const wynik = await zrodlo.wstaw({
      documentId: dokument,
      offset: miejsceKursora(),
      rows: wiersze,
      columns: kolumny,
      ...(styl.kontrolka.value.trim() === '' ? {} : { styleName: styl.kontrolka.value.trim() }),
      ...(liczba(wierszeNaglowka.value) > 0 ? { headerRows: liczba(wierszeNaglowka.value) } : {}),
      ...(powtarzajNaglowek.checked ? { repeatHeader: true } : {}),
    });
    if (!przyjalSie('Wstawienie tabeli', wynik)) return;
    if (wynik.wynik === undefined) return;
    wskazana = wynik.wynik.table.id;
    opiszSkutek(
      `Tabela ${wiersze} × ${kolumny} wstawiona jako ${wynik.wynik.table.id}`,
      wynik.wynik.balance,
      wynik.wynik.actionId,
      opiszSzerokosci(wynik.wynik.table),
    );
    void odczytajTabele();
  }

  async function zmienBudowe(operacja: StudioTableStructureOp, nazwa: string): Promise<void> {
    const dokument = idDokumentu();
    const tabela = idTabeli();
    if (dokument === null || tabela === null) return;
    const wynik = await zrodlo.budowa({
      documentId: dokument,
      tableId: tabela,
      operation: operacja,
      ...(liczba(wiersz.value) >= 0 && wiersz.value !== '' ? { row: liczba(wiersz.value) } : {}),
      ...(kolumna.value === '' ? {} : { column: liczba(kolumna.value) }),
      ...(liczba(ile.value) > 0 ? { count: liczba(ile.value) } : {}),
      ...(liczba(obejmijWierszy.value) > 0 ? { rowSpan: liczba(obejmijWierszy.value) } : {}),
      ...(liczba(obejmijKolumn.value) > 0 ? { columnSpan: liczba(obejmijKolumn.value) } : {}),
      ...(przedWskazanym.checked ? { before: true } : {}),
    });
    if (!przyjalSie(nazwa, wynik)) return;
    if (wynik.wynik === undefined) return;
    opiszSkutek(nazwa, wynik.wynik.balance, wynik.wynik.actionId, opiszSzerokosci(wynik.wynik.table));
    void odczytajTabele();
  }

  async function przestawPostac(): Promise<void> {
    const dokument = idDokumentu();
    const tabela = idTabeli();
    if (dokument === null || tabela === null) return;
    const szerokosci = rozbijSzerokosci(szerokosciKolumn.kontrolka.value);
    const wynik = await zrodlo.postac({
      documentId: dokument,
      tableId: tabela,
      ...(wiersz.value === '' ? {} : { row: liczba(wiersz.value) }),
      ...(kolumna.value === '' ? {} : { column: liczba(kolumna.value) }),
      ...(szerokosci.length === 0 ? {} : { columnWidthsMm: szerokosci }),
      ...(liczba(szerokoscTabeli.value) > 0 ? { widthMm: liczba(szerokoscTabeli.value) } : {}),
      ...(styl.kontrolka.value.trim() === '' ? {} : { styleName: styl.kontrolka.value.trim() }),
      ...(cieniowanie.kontrolka.value.trim() === ''
        ? {}
        : { shadingColor: cieniowanie.kontrolka.value.trim() }),
      ...(wyrownaniePionowe.value === ''
        ? {}
        : { verticalAlign: wyrownaniePionowe.value as StudioVerticalAlign }),
      ...(wyrownaniePoziome.value === ''
        ? {}
        : { align: wyrownaniePoziome.value as StudioTextAlign }),
      ...(liczba(wierszeNaglowka.value) > 0 ? { headerRows: liczba(wierszeNaglowka.value) } : {}),
      ...(powtarzajNaglowek.checked ? { repeatHeader: true } : {}),
      ...(podpisTabeli.kontrolka.value.trim() === ''
        ? {}
        : { caption: podpisTabeli.kontrolka.value.trim() }),
    });
    if (!przyjalSie('Postać tabeli', wynik)) return;
    if (wynik.wynik === undefined) return;
    opiszSkutek(
      'Postać tabeli przestawiona',
      wynik.wynik.balance,
      wynik.wynik.actionId,
      opiszSzerokosci(wynik.wynik.table),
    );
    void odczytajTabele();
  }

  async function sortujTabele(): Promise<void> {
    const dokument = idDokumentu();
    const tabela = idTabeli();
    if (dokument === null || tabela === null) return;
    const wynik = await zrodlo.sortuj({
      documentId: dokument,
      tableId: tabela,
      column: liczba(kolumnaSortowania.value),
      ...(malejaco.checked ? { descending: true } : {}),
      ...(jakLiczby.checked ? { numeric: true } : {}),
    });
    if (!przyjalSie('Sortowanie tabeli', wynik)) return;
    if (wynik.wynik === undefined) return;
    opiszSkutek(
      `Tabela posortowana po kolumnie ${liczba(kolumnaSortowania.value)}`,
      wynik.wynik.balance,
      wynik.wynik.actionId,
      'Wiersz nagłówkowy zostaje na miejscu, o ile tabela go ma.',
    );
    void odczytajTabele();
  }

  async function zamienTekstNaTabele(): Promise<void> {
    const dokument = idDokumentu();
    if (dokument === null) return;
    const zakres = stan.zaznaczenie();
    if (zakres === null) {
      odpowiedz.pokaz(BRAK_ZAZNACZENIA, false);
      return;
    }
    const wynik = await zrodlo.zamien({
      documentId: dokument,
      direction: StudioTableConvert.TextToTable,
      rangeStart: zakres.poczatek,
      rangeEnd: zakres.koniec,
      ...(rozdzielnik.kontrolka.value === ''
        ? {}
        : { separator: rozdzielnik.kontrolka.value }),
    });
    if (!przyjalSie('Zamiana tekstu na tabelę', wynik)) return;
    if (wynik.wynik === undefined) return;
    const tabela = wynik.wynik.table;
    if (tabela !== undefined) wskazana = tabela.id;
    opiszSkutek(
      'Zaznaczony tekst zamieniony na tabelę',
      wynik.wynik.balance,
      wynik.wynik.actionId,
      tabela === undefined ? '' : opiszSzerokosci(tabela),
    );
    void odczytajTabele();
  }

  async function zamienTabeleNaTekst(): Promise<void> {
    const dokument = idDokumentu();
    const tabela = idTabeli();
    if (dokument === null || tabela === null) return;
    const wynik = await zrodlo.zamien({
      documentId: dokument,
      direction: StudioTableConvert.TableToText,
      tableId: tabela,
      ...(rozdzielnik.kontrolka.value === '' ? {} : { separator: rozdzielnik.kontrolka.value }),
    });
    if (!przyjalSie('Zamiana tabeli na tekst', wynik)) return;
    if (wynik.wynik === undefined) return;
    wskazana = '';
    opiszSkutek(
      'Tabela zamieniona na tekst',
      wynik.wynik.balance,
      wynik.wynik.actionId,
      'Obramowanie, cieniowanie i scalenia komórek w tekście nie mają czego nieść — bilans wyżej ' +
        'mówi, co przez to odpadło.',
    );
    void odczytajTabele();
  }

  async function odczytajTabele(): Promise<void> {
    const dokument = stan.dokument();
    if (dokument === null) {
      tabele = [];
      przerysujWykaz();
      return;
    }
    const wynik = await zrodlo.wykaz({ documentId: dokument.id });
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu('Odczyt tabel dokumentu', wynik.blad), false);
      return;
    }
    tabele = wynik.wynik.tables;
    przerysujWykaz();
  }

  /** Wynik wywołania: `true`, gdy rdzeń odpowiedział; odmowę wypisuje sam. */
  function przyjalSie(nazwa: string, wynik: Wynik<unknown>): boolean {
    if (wynik.udany && wynik.wynik !== undefined) return true;
    odpowiedz.pokaz(opisOdmowyBledu(nazwa, wynik.blad), false);
    return false;
  }

  function przerysujWykaz(): void {
    if (tabele.length === 0) {
      const puste = document.createElement('li');
      puste.className = 'dn-pole-opis';
      puste.textContent =
        'Dokument nie ma ani jednej tabeli. To odpowiedź rdzenia, nie brak odczytu — wstaw tabelę ' +
        'siatką wyżej.';
      wykaz.replaceChildren(puste);
      return;
    }
    wykaz.replaceChildren(
      ...tabele.map((tabela) => {
        const glowa = document.createElement('p');
        glowa.className = 'ms-tabela__glowa';
        glowa.textContent =
          `${tabela.rows} × ${tabela.columns}` +
          (tabela.styleName === undefined ? '' : ` · styl ${tabela.styleName}`) +
          (tabela.headerRows === undefined ? '' : ` · wierszy nagłówka ${tabela.headerRows}`) +
          (tabela.repeatHeader === true ? ' · nagłówek powtarzany' : '') +
          (tabela.caption === undefined ? '' : ` · podpis „${tabela.caption}"`);

        const miary = document.createElement('p');
        miary.className = 'dn-pole-opis';
        miary.textContent = opiszSzerokosci(tabela);

        const wskaz = przycisk(
          tabela.id === wskazana ? 'Wskazana' : 'Wskaż',
          'dn-btn dn-btn--sm dn-btn--zarys',
        );
        wskaz.dataset['tabela'] = tabela.id;
        wskaz.setAttribute('aria-pressed', String(tabela.id === wskazana));
        wskaz.addEventListener('click', () => {
          wskazana = tabela.id;
          przerysujWykaz();
        });

        const pozycja = document.createElement('li');
        pozycja.dataset['tabela'] = tabela.id;
        pozycja.dataset['wskazana'] = String(tabela.id === wskazana);
        pozycja.append(glowa, miary, wskaz);
        return pozycja;
      }),
    );
  }

  /* ── Nakładka ──────────────────────────────────────────────────────────── */

  const tresc = document.createElement('div');
  tresc.className = 'ms-tabela__tresc';
  tresc.append(
    czesc('Wstawienie tabeli', [
      podpisSiatki,
      siatka,
      wierszeRecznie,
      kolumnyRecznie,
      styl.element,
      wierszeNaglowka,
      etykietaPowtarzania,
      wstawRecznie,
    ]),
    czesc('Tabele dokumentu', [wykaz]),
    czesc('Budowa tabeli wskazanej', [
      wiersz,
      kolumna,
      ile,
      obejmijWierszy,
      obejmijKolumn,
      etykietaPrzed,
      rzadBudowy,
    ]),
    czesc('Postać tabeli i komórek', [
      szerokosciKolumn.element,
      szerokoscTabeli,
      cieniowanie.element,
      podpisTabeli.element,
      wyrownaniePionowe,
      wyrownaniePoziome,
      ustawPostac,
    ]),
    czesc('Sortowanie', [kolumnaSortowania, etykietaMalejaco, etykietaLiczb, sortuj]),
    czesc('Zamiana tekstu i tabeli', [rozdzielnik.element, naTabele, naTekst]),
    odpowiedz.element,
  );

  const nakladka = document.createElement('div');
  nakladka.className = 'ms-tabela__nakladka';
  nakladka.hidden = true;
  nakladka.append(tresc);

  const wyzwalacz = przycisk('Tabele ▾', 'dn-btn dn-btn--sm dn-btn--zarys');
  wyzwalacz.dataset['czynnosc'] = 'warsztat-tabel';
  wyzwalacz.setAttribute('aria-expanded', 'false');
  wyzwalacz.title =
    'Otwiera warsztat tabel nakładką: wskazanie rozmiaru siatką, budowa tabeli, postać komórek, ' +
    'sortowanie i zamiana tekstu na tabelę. Powierzchnia należy do dokumentu, więc warsztat ' +
    'schodzi drugim naciśnięciem albo Escapem.';
  wyzwalacz.addEventListener('click', () => przestaw(nakladka.hidden));

  const element = document.createElement('section');
  element.className = 'ms-tabela';
  element.setAttribute('aria-label', 'Warsztat tabel dokumentu');
  element.append(wyzwalacz, nakladka);
  element.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key === 'Escape' && !nakladka.hidden) przestaw(false);
  });

  function przestaw(otwarta: boolean): void {
    nakladka.hidden = !otwarta;
    wyzwalacz.setAttribute('aria-expanded', otwarta ? 'true' : 'false');
    wyzwalacz.textContent =
      tabele.length === 0 ? `Tabele ${otwarta ? '▴' : '▾'}` : `Tabele: ${tabele.length} ${otwarta ? '▴' : '▾'}`;
    if (otwarta) void odczytajTabele();
  }

  przerysujWykaz();

  return {
    element,
    przestawWidocznosc: () => przestaw(nakladka.hidden),
    odswiez: () => {
      if (!nakladka.hidden) void odczytajTabele();
    },
  };
}

/**
 * Zdanie o szerokościach kolumn — miara policzona, nie założona.
 *
 * Sprawdzian skutku zlecenia mierzy właśnie to: tabela po scaleniu komórek ma
 * szerokości POLICZONE, nie zerowe. Panel wypisuje je wprost, żeby zero było
 * widoczne w oknie, a nie tylko w bazie.
 */
function opiszSzerokosci(tabela: StudioDocumentTable): string {
  const szerokosci = tabela.columnWidthsMm ?? [];
  if (szerokosci.length === 0) {
    return 'Rdzeń nie oddał szerokości kolumn — tabela stoi bez policzonych miar.';
  }
  const zerowe = szerokosci.filter((miara) => miara <= 0).length;
  const razem = tabela.widthMm === undefined ? '' : ` · szerokość tabeli ${tabela.widthMm} mm`;
  return (
    `Szerokości kolumn w milimetrach: ${szerokosci.join(', ')}${razem}.` +
    (zerowe === 0 ? '' : ` UWAGA: kolumn o szerokości zerowej ${zerowe} — miara nie policzyła się.`)
  );
}

/** Rozbija „40, 30, 30" na liczby; wartości niepoprawne odpadają. */
function rozbijSzerokosci(wartosc: string): number[] {
  return wartosc
    .split(',')
    .map((czesc) => Number.parseFloat(czesc.trim()))
    .filter((miara) => Number.isFinite(miara) && miara > 0);
}

/** Liczba całkowita z pola; wartość pusta i niepoprawna znaczą zero. */
function liczba(wartosc: string): number {
  const odczytana = Number.parseInt(wartosc, 10);
  return Number.isFinite(odczytana) ? odczytana : 0;
}

/** Część nakładki wraz z jej tytułem. */
function czesc(tytul: string, elementy: readonly HTMLElement[]): HTMLElement {
  const naglowek = document.createElement('p');
  naglowek.className = 'ms-tabela__tytul';
  naglowek.textContent = tytul;
  const sekcja = document.createElement('section');
  sekcja.className = 'ms-tabela__czesc';
  sekcja.append(naglowek, ...elementy);
  return sekcja;
}

const BRAK_DOKUMENTU =
  'Tabela wchodzi do dokumentu — wczytaj go albo załóż nowy. Komenda studio.table.insert przyjmuje ' +
  'identyfikator dokumentu jako pole obowiązkowe.';

const BRAK_WSKAZANEJ =
  'Ta czynność dotyczy tabeli wskazanej — naciśnij „Wskaż" przy tabeli w wykazie wyżej. Rdzeń ' +
  'przyjmuje identyfikator tabeli jako pole obowiązkowe i bez niego nie ma na czym pracować.';

const BRAK_ROZMIARU =
  'Tabela potrzebuje liczby wierszy i kolumn większej od zera — wskaż rozmiar siatką albo podaj go ' +
  'liczbami.';

const BRAK_ZAZNACZENIA =
  'Zamiana tekstu na tabelę pracuje na zaznaczonym fragmencie — zaznacz tekst w dokumencie. Bez ' +
  'zakresu rdzeń nie wie, co zamienić, a zamiana całego dokumentu nie jest tym, o co Operator prosi.';
