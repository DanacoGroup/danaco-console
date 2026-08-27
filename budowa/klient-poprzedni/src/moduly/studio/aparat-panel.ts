import {
  StudioApparatusKind,
  StudioFieldKind,
  type StudioActionBalance,
  type StudioApparatusItem,
  type StudioDocumentField,
} from '../../../../shared/contract';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import {
  poleLiczbowe,
  poleTekstowe,
  poleTresci,
  przycisk,
  utworzWierszOdpowiedzi,
  wybor,
} from '../../modele/kontrolki-formularza';
import type { AparatZrodlo } from './aparat-zrodlo';
import type { Wynik } from '../../protokol/kanal';
import type { StanStudio } from './stan-studio';
import { wstawieniaOpiszBilans } from './zrodlo-wstawien-studio';

/**
 * Aparat dokumentu i pola — nakładka na żądanie. Elementy aparatu i pola
 * dokumentu są wyliczane z treści i po jej zmianie stają się nieświeże, aż
 * odświeżenie policzy je od nowa. Numer przypisu i powołania nadaje rdzeń,
 * panel go wyłącznie pokazuje.
 */
export interface AparatPanel {
  element: HTMLElement;
  przestawWidocznosc(): void;
  odswiez(): void;
}

/** Trzynaście rodzajów elementów aparatu wraz z nazwą pełną, pokazywaną w oknie zamiast kodu kontraktu. */
const RODZAJE_APARATU: readonly (readonly [string, string])[] = [
  [StudioApparatusKind.Toc, 'spis treści'],
  [StudioApparatusKind.FigureIndex, 'spis ilustracji'],
  [StudioApparatusKind.TableIndex, 'spis tabel'],
  [StudioApparatusKind.Footnote, 'przypis dolny'],
  [StudioApparatusKind.Endnote, 'przypis końcowy'],
  [StudioApparatusKind.Caption, 'podpis pod ilustracją albo tabelą'],
  [StudioApparatusKind.Bookmark, 'zakładka'],
  [StudioApparatusKind.CrossReference, 'odwołanie wzajemne'],
  [StudioApparatusKind.Hyperlink, 'odsyłacz'],
  [StudioApparatusKind.Citation, 'powołanie bibliograficzne'],
  [StudioApparatusKind.Bibliography, 'bibliografia'],
  [StudioApparatusKind.IndexEntry, 'hasło indeksu'],
  [StudioApparatusKind.Index, 'indeks'],
];

/** Dziewięć rodzajów pól dokumentu wraz z ich nazwą pełną, pokazywaną w tym panelu zamiast kodu kontraktu. */
const RODZAJE_POL: readonly (readonly [string, string])[] = [
  [StudioFieldKind.PageNumber, 'numer strony'],
  [StudioFieldKind.PageCount, 'liczba stron'],
  [StudioFieldKind.Date, 'data'],
  [StudioFieldKind.Time, 'godzina'],
  [StudioFieldKind.DocumentTitle, 'tytuł dokumentu'],
  [StudioFieldKind.DocumentAuthor, 'autor dokumentu'],
  [StudioFieldKind.DocumentProperty, 'właściwość dokumentu'],
  [StudioFieldKind.Calculated, 'pole obliczane'],
  [StudioFieldKind.TemplateField, 'pole szablonu do wypełnienia'],
];

export function utworzAparatPanel(stan: StanStudio, zrodlo: AparatZrodlo): AparatPanel {
  const odpowiedz = utworzWierszOdpowiedzi();

  let elementy: readonly StudioApparatusItem[] = [];
  let pola: readonly StudioDocumentField[] = [];
  let wskazany = '';

  /* ── Założenie elementu aparatu ────────────────────────────────────────── */

  const rodzaj = wybor('Rodzaj elementu aparatu', RODZAJE_APARATU.map((pozycja) => pozycja));
  const brzmienie = poleTresci(
    'Brzmienie przypisu, podpisu albo hasła',
    3,
    'treść elementu widoczna w dokumencie',
  );
  const etykieta = poleTekstowe({
    etykieta: 'Etykieta widoczna w treści',
    podpowiedz: 'np. Spis treści',
  });
  const poziomOd = poleLiczbowe('Od którego poziomu nagłówków zbierać spis', 'np. 1');
  const poziomDo = poleLiczbowe('Do którego poziomu nagłówków zbierać spis', 'np. 3');
  const celElementu = poleTekstowe({
    etykieta: 'Element, do którego odwołanie prowadzi',
    podpowiedz: 'identyfikator zakładki albo podpisu',
  });
  const celAdresu = poleTekstowe({
    etykieta: 'Adres, do którego odsyłacz prowadzi',
    podpowiedz: 'adres strony',
  });
  const kluczPowolania = poleTekstowe({
    etykieta: 'Klucz powołania bibliograficznego',
    podpowiedz: 'czym powołanie odwołuje się do pozycji bibliografii',
  });
  const tytulZrodla = poleTekstowe({ etykieta: 'Tytuł źródła', podpowiedz: '' });
  const autorZrodla = poleTekstowe({ etykieta: 'Autor źródła', podpowiedz: '' });
  const rokZrodla = poleTekstowe({ etykieta: 'Rok źródła', podpowiedz: '' });
  const adresZrodla = poleTekstowe({ etykieta: 'Adres źródła', podpowiedz: '' });

  const zaloz = przycisk('Załóż element aparatu', 'dn-btn dn-btn--sm dn-btn--sygnal');
  zaloz.dataset['czynnosc'] = 'zaloz-aparat';
  zaloz.addEventListener('click', () => void zalozElement());

  /* ── Wykaz aparatu ─────────────────────────────────────────────────────── */

  const zawezenie = wybor('Zawężenie wykazu do rodzaju', [
    ['', 'wszystkie rodzaje'],
    ...RODZAJE_APARATU,
  ]);
  zawezenie.addEventListener('change', () => void odczytajAparat());

  const tylkoNieswieze = document.createElement('input');
  tylkoNieswieze.type = 'checkbox';
  tylkoNieswieze.className = 'dn-przelacznik';
  tylkoNieswieze.setAttribute('aria-label', 'Tylko elementy wymagające odświeżenia');
  tylkoNieswieze.addEventListener('change', () => void odczytajAparat());
  const etykietaNieswiezych = document.createElement('label');
  etykietaNieswiezych.className = 'ms-aparat__zawezenie';
  etykietaNieswiezych.append(
    tylkoNieswieze,
    document.createTextNode('tylko wymagające odświeżenia'),
  );

  const stanSwiezosci = document.createElement('p');
  stanSwiezosci.className = 'dn-pole-opis ms-aparat__swiezosc';

  const wykaz = document.createElement('ul');
  wykaz.className = 'ms-aparat__wykaz';
  wykaz.setAttribute('aria-label', 'Aparat dokumentu');

  const odswiezWszystko = przycisk('Odśwież cały aparat', 'dn-btn dn-btn--sm dn-btn--atrament');
  odswiezWszystko.dataset['czynnosc'] = 'odswiez-aparat';
  odswiezWszystko.title =
    'Spis treści zgadza się po odświeżeniu z nagłówkami dokumentu, przypisy są przenumerowane, ' +
    'a spisy ilustracji i tabel policzone od nowa. Numery nadaje rdzeń.';
  odswiezWszystko.addEventListener('click', () => void odswiezAparat(''));

  const odswiezWskazany = przycisk(
    'Odśwież element wskazany',
    'dn-btn dn-btn--sm dn-btn--zarys',
  );
  odswiezWskazany.dataset['czynnosc'] = 'odswiez-element';
  odswiezWskazany.addEventListener('click', () => {
    if (wskazany === '') {
      odpowiedz.pokaz(BRAK_WSKAZANEGO, false);
      return;
    }
    void odswiezAparat(wskazany);
  });

  const usun = przycisk('Usuń element wskazany', 'dn-btn dn-btn--sm dn-btn--duch');
  usun.dataset['czynnosc'] = 'usun-element';
  usun.title =
    'Po usunięciu przypisu przypisy pozostałe przenumerowują się — numer nadaje rdzeń, nie okno.';
  usun.addEventListener('click', () => void usunElement());

  /* ── Pola dokumentu ────────────────────────────────────────────────────── */

  const rodzajPola = wybor('Rodzaj pola dokumentu', RODZAJE_POL.map((pozycja) => pozycja));
  const formatPola = poleTekstowe({
    etykieta: 'Format wartości',
    podpowiedz: 'wzór daty albo format liczby',
  });
  const wyrazenie = poleTekstowe({
    etykieta: 'Wyrażenie pola obliczanego',
    podpowiedz: 'liczy je rdzeń przy odświeżeniu',
  });
  const nazwaWlasciwosci = poleTekstowe({
    etykieta: 'Nazwa właściwości dokumentu albo pola szablonu',
    podpowiedz: '',
  });

  const wstawPole = przycisk('Wstaw pole w miejsce kursora', 'dn-btn dn-btn--sm dn-btn--sygnal');
  wstawPole.dataset['czynnosc'] = 'wstaw-pole';
  wstawPole.addEventListener('click', () => void wstawPoleDokumentu());

  const odswiezPola = przycisk('Odśwież pola dokumentu', 'dn-btn dn-btn--sm dn-btn--atrament');
  odswiezPola.dataset['czynnosc'] = 'odswiez-pola';
  odswiezPola.addEventListener('click', () => void odswiezPolaDokumentu());

  const wykazPol = document.createElement('ul');
  wykazPol.className = 'ms-aparat__pola';
  wykazPol.setAttribute('aria-label', 'Pola dokumentu');

  /* ── Czynności ─────────────────────────────────────────────────────────── */

  function idDokumentu(): string | null {
    const dokument = stan.dokument();
    if (dokument === null) {
      odpowiedz.pokaz(BRAK_DOKUMENTU, false);
      return null;
    }
    return dokument.id;
  }

  function miejsceKursora(): number {
    const zakres = stan.zaznaczenie();
    return zakres === null ? stan.trescRobocza().length : zakres.koniec;
  }

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

  async function zalozElement(): Promise<void> {
    const dokument = idDokumentu();
    if (dokument === null) return;
    const zakres = stan.zaznaczenie();
    const wynik = await zrodlo.zaloz({
      documentId: dokument,
      kind: rodzaj.value as StudioApparatusKind,
      offset: miejsceKursora(),
      ...(zakres === null ? {} : { rangeStart: zakres.poczatek, rangeEnd: zakres.koniec }),
      ...(brzmienie.value.trim() === '' ? {} : { text: brzmienie.value.trim() }),
      ...(etykieta.kontrolka.value.trim() === ''
        ? {}
        : { label: etykieta.kontrolka.value.trim() }),
      ...(liczba(poziomOd.value) > 0 ? { levelsFrom: liczba(poziomOd.value) } : {}),
      ...(liczba(poziomDo.value) > 0 ? { levelsTo: liczba(poziomDo.value) } : {}),
      ...(celElementu.kontrolka.value.trim() === ''
        ? {}
        : { targetId: celElementu.kontrolka.value.trim() }),
      ...(celAdresu.kontrolka.value.trim() === ''
        ? {}
        : { targetUrl: celAdresu.kontrolka.value.trim() }),
      ...(kluczPowolania.kontrolka.value.trim() === ''
        ? {}
        : { citationKey: kluczPowolania.kontrolka.value.trim() }),
      ...(tytulZrodla.kontrolka.value.trim() === ''
        ? {}
        : { sourceTitle: tytulZrodla.kontrolka.value.trim() }),
      ...(autorZrodla.kontrolka.value.trim() === ''
        ? {}
        : { sourceAuthor: autorZrodla.kontrolka.value.trim() }),
      ...(rokZrodla.kontrolka.value.trim() === ''
        ? {}
        : { sourceYear: rokZrodla.kontrolka.value.trim() }),
      ...(adresZrodla.kontrolka.value.trim() === ''
        ? {}
        : { sourceUrl: adresZrodla.kontrolka.value.trim() }),
    });
    if (!przyjalSie('Założenie elementu aparatu', wynik)) return;
    if (wynik.wynik === undefined) return;
    wskazany = wynik.wynik.item.id;
    opiszSkutek(
      `Element aparatu założony: ${opiszRodzajAparatu(wynik.wynik.item.kind)} ${wynik.wynik.item.id}`,
      wynik.wynik.balance,
      wynik.wynik.actionId,
      opiszElement(wynik.wynik.item),
    );
    void odczytajAparat();
  }

  async function odswiezAparat(idElementu: string): Promise<void> {
    const dokument = idDokumentu();
    if (dokument === null) return;
    const wynik = await zrodlo.odswiez({
      documentId: dokument,
      ...(idElementu === '' ? {} : { itemId: idElementu }),
      ...(idElementu === '' && zawezenie.value !== ''
        ? { kind: zawezenie.value as StudioApparatusKind }
        : {}),
    });
    if (!przyjalSie('Odświeżenie aparatu', wynik)) return;
    if (wynik.wynik === undefined) return;
    elementy = wynik.wynik.items;
    opiszSkutek(
      idElementu === '' ? 'Aparat dokumentu odświeżony' : `Element ${idElementu} odświeżony`,
      wynik.wynik.balance,
      wynik.wynik.actionId,
      `Elementów po odświeżeniu: ${elementy.length}, wciąż nieświeżych: ${ileNieswiezych()}.`,
    );
    przerysujAparat();
  }

  async function usunElement(): Promise<void> {
    const dokument = idDokumentu();
    if (dokument === null) return;
    if (wskazany === '') {
      odpowiedz.pokaz(BRAK_WSKAZANEGO, false);
      return;
    }
    const wynik = await zrodlo.usun({ documentId: dokument, itemId: wskazany });
    if (!przyjalSie('Usunięcie elementu aparatu', wynik)) return;
    if (wynik.wynik === undefined) return;
    if (!wynik.wynik.removed) {
      odpowiedz.pokaz(
        `Rdzeń NIE usunął elementu ${wskazany}. ${wstawieniaOpiszBilans(wynik.wynik.balance)}`,
        false,
      );
      return;
    }
    const usuniety = wskazany;
    wskazany = '';
    opiszSkutek(
      `Element aparatu ${usuniety} usunięty`,
      wynik.wynik.balance,
      wynik.wynik.actionId,
      'Przypisy pozostałe przenumerowały się — odśwież aparat, żeby zobaczyć nowe numery.',
    );
    void odczytajAparat();
  }

  async function wstawPoleDokumentu(): Promise<void> {
    const dokument = idDokumentu();
    if (dokument === null) return;
    const wynik = await zrodlo.wstawPole({
      documentId: dokument,
      offset: miejsceKursora(),
      kind: rodzajPola.value as StudioFieldKind,
      ...(formatPola.kontrolka.value.trim() === ''
        ? {}
        : { format: formatPola.kontrolka.value.trim() }),
      ...(wyrazenie.kontrolka.value.trim() === ''
        ? {}
        : { expression: wyrazenie.kontrolka.value.trim() }),
      ...(nazwaWlasciwosci.kontrolka.value.trim() === ''
        ? {}
        : { propertyName: nazwaWlasciwosci.kontrolka.value.trim() }),
    });
    if (!przyjalSie('Wstawienie pola', wynik)) return;
    if (wynik.wynik === undefined) return;
    opiszSkutek(
      `Pole ${opiszRodzajPola(wynik.wynik.field.kind)} wstawione jako ${wynik.wynik.field.id}`,
      wynik.wynik.balance,
      wynik.wynik.actionId,
      opiszPole(wynik.wynik.field),
    );
    void odczytajPola();
  }

  async function odswiezPolaDokumentu(): Promise<void> {
    const dokument = idDokumentu();
    if (dokument === null) return;
    const wynik = await zrodlo.odswiezPola({ documentId: dokument });
    if (!przyjalSie('Odświeżenie pól', wynik)) return;
    if (wynik.wynik === undefined) return;
    pola = wynik.wynik.fields;
    opiszSkutek(
      'Pola dokumentu odświeżone',
      wynik.wynik.balance,
      wynik.wynik.actionId,
      `Pól po odświeżeniu: ${pola.length}, wciąż nieświeżych: ${
        pola.filter((pole) => pole.stale === true).length
      }.`,
    );
    przerysujPola();
  }

  async function odczytajAparat(): Promise<void> {
    const dokument = stan.dokument();
    if (dokument === null) {
      elementy = [];
      przerysujAparat();
      return;
    }
    const wynik = await zrodlo.wykaz({
      documentId: dokument.id,
      ...(zawezenie.value === '' ? {} : { kind: zawezenie.value as StudioApparatusKind }),
      ...(tylkoNieswieze.checked ? { staleOnly: true } : {}),
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu('Odczyt aparatu dokumentu', wynik.blad), false);
      return;
    }
    elementy = wynik.wynik.items;
    przerysujAparat();
  }

  async function odczytajPola(): Promise<void> {
    const dokument = stan.dokument();
    if (dokument === null) {
      pola = [];
      przerysujPola();
      return;
    }
    const wynik = await zrodlo.wykazPol({ documentId: dokument.id });
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu('Odczyt pól dokumentu', wynik.blad), false);
      return;
    }
    pola = wynik.wynik.fields;
    przerysujPola();
  }

  function przyjalSie(nazwa: string, wynik: Wynik<unknown>): boolean {
    if (wynik.udany && wynik.wynik !== undefined) return true;
    odpowiedz.pokaz(opisOdmowyBledu(nazwa, wynik.blad), false);
    return false;
  }

  function ileNieswiezych(): number {
    return (
      elementy.filter((pozycja) => pozycja.stale === true).length +
      pola.filter((pole) => pole.stale === true).length
    );
  }

  function przerysujAparat(): void {
    const nieswieze = ileNieswiezych();
    stanSwiezosci.textContent =
      nieswieze === 0
        ? `Elementów aparatu ${elementy.length}, pól ${pola.length}. Żaden nie wymaga odświeżenia.`
        : `Elementów aparatu ${elementy.length}, pól ${pola.length}. WYMAGA ODŚWIEŻENIA: ` +
          `${nieswieze} — spis treści albo numeracja rozjechały się z dokumentem, więc to, co ` +
          'widać na kartce, nie zgadza się z jego treścią.';
    stanSwiezosci.dataset['nieswieze'] = String(nieswieze);
    // Licznik na przycisku przelicza się przy przerysowaniu, bo wykaz przychodzi z rdzenia później.
    opiszWyzwalacz(nieswieze);

    if (elementy.length === 0) {
      const puste = document.createElement('li');
      puste.className = 'dn-pole-opis';
      puste.textContent =
        'Dokument nie ma ani jednego elementu aparatu tego rodzaju. To odpowiedź rdzenia, nie brak ' +
        'odczytu — załóż spis treści, przypis albo zakładkę wyżej.';
      wykaz.replaceChildren(puste);
      return;
    }
    wykaz.replaceChildren(
      ...elementy.map((pozycja) => {
        const glowa = document.createElement('p');
        glowa.className = 'ms-aparat__glowa';
        glowa.textContent =
          `${opiszRodzajAparatu(pozycja.kind)}` +
          (pozycja.number === undefined ? '' : ` · numer nadany przez rdzeń: ${pozycja.number}`) +
          (pozycja.label === undefined ? '' : ` · „${pozycja.label}"`) +
          (pozycja.stale === true ? ' · WYMAGA ODŚWIEŻENIA' : '');

        const szczegoly = document.createElement('p');
        szczegoly.className = 'dn-pole-opis';
        szczegoly.textContent = opiszElement(pozycja);

        const wskaz = przycisk(
          pozycja.id === wskazany ? 'Wskazany' : 'Wskaż',
          'dn-btn dn-btn--sm dn-btn--zarys',
        );
        wskaz.dataset['element'] = pozycja.id;
        wskaz.setAttribute('aria-pressed', String(pozycja.id === wskazany));
        wskaz.addEventListener('click', () => {
          wskazany = pozycja.id;
          przerysujAparat();
        });

        const wiersz = document.createElement('li');
        wiersz.dataset['element'] = pozycja.id;
        wiersz.dataset['rodzaj'] = pozycja.kind;
        wiersz.dataset['nieswiezy'] = String(pozycja.stale === true);
        wiersz.dataset['wskazany'] = String(pozycja.id === wskazany);
        wiersz.append(glowa, szczegoly, wskaz);
        return wiersz;
      }),
    );
  }

  function przerysujPola(): void {
    przerysujAparat();
    if (pola.length === 0) {
      const puste = document.createElement('li');
      puste.className = 'dn-pole-opis';
      puste.textContent =
        'Dokument nie ma ani jednego pola. Wstaw numer strony, datę albo pole obliczane wyżej.';
      wykazPol.replaceChildren(puste);
      return;
    }
    wykazPol.replaceChildren(
      ...pola.map((pole) => {
        const wiersz = document.createElement('li');
        wiersz.dataset['pole'] = pole.id;
        wiersz.dataset['rodzaj'] = pole.kind;
        wiersz.dataset['nieswieze'] = String(pole.stale === true);
        wiersz.textContent =
          `${opiszRodzajPola(pole.kind)} · ${opiszPole(pole)}` +
          (pole.stale === true ? ' WYMAGA ODŚWIEŻENIA.' : '');
        return wiersz;
      }),
    );
  }

  /* ── Nakładka ──────────────────────────────────────────────────────────── */

  const tresc = document.createElement('div');
  tresc.className = 'ms-aparat__tresc';
  tresc.append(
    czesc('Założenie elementu aparatu', [
      rodzaj,
      brzmienie,
      etykieta.element,
      poziomOd,
      poziomDo,
      celElementu.element,
      celAdresu.element,
      kluczPowolania.element,
      tytulZrodla.element,
      autorZrodla.element,
      rokZrodla.element,
      adresZrodla.element,
      zaloz,
    ]),
    czesc('Aparat dokumentu', [
      zawezenie,
      etykietaNieswiezych,
      stanSwiezosci,
      wykaz,
      odswiezWszystko,
      odswiezWskazany,
      usun,
    ]),
    czesc('Pola dokumentu', [
      rodzajPola,
      formatPola.element,
      wyrazenie.element,
      nazwaWlasciwosci.element,
      wstawPole,
      wykazPol,
      odswiezPola,
    ]),
    odpowiedz.element,
  );

  const nakladka = document.createElement('div');
  nakladka.className = 'ms-aparat__nakladka';
  nakladka.hidden = true;
  nakladka.append(tresc);

  const wyzwalacz = przycisk('Aparat dokumentu ▾', 'dn-btn dn-btn--sm dn-btn--zarys');
  wyzwalacz.dataset['czynnosc'] = 'warsztat-aparatu';
  wyzwalacz.setAttribute('aria-expanded', 'false');
  wyzwalacz.title =
    'Otwiera aparat dokumentu nakładką: spis treści z poziomów nagłówków, spisy ilustracji i tabel, ' +
    'przypisy dolne i końcowe, podpisy, zakładki, odsyłacze, odwołania wzajemne, powołania, ' +
    'bibliografia, indeks oraz pola dokumentu. Schodzi drugim naciśnięciem albo Escapem.';
  wyzwalacz.addEventListener('click', () => przestaw(nakladka.hidden));

  const element = document.createElement('section');
  element.className = 'ms-aparat';
  element.setAttribute('aria-label', 'Aparat dokumentu i pola');
  element.append(wyzwalacz, nakladka);
  element.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key === 'Escape' && !nakladka.hidden) przestaw(false);
  });

  /** Napis przycisku wraz z licznikiem elementów do odświeżenia. */
  function opiszWyzwalacz(nieswieze: number): void {
    const znak = nakladka.hidden ? '▾' : '▴';
    wyzwalacz.textContent =
      nieswieze === 0
        ? `Aparat dokumentu ${znak}`
        : `Aparat dokumentu: ${nieswieze} do odświeżenia ${znak}`;
  }

  function przestaw(otwarta: boolean): void {
    nakladka.hidden = !otwarta;
    wyzwalacz.setAttribute('aria-expanded', otwarta ? 'true' : 'false');
    opiszWyzwalacz(ileNieswiezych());
    if (otwarta) {
      void odczytajAparat();
      void odczytajPola();
    }
  }

  przerysujPola();

  return {
    element,
    przestawWidocznosc: () => przestaw(nakladka.hidden),
    odswiez: () => {
      if (nakladka.hidden) return;
      void odczytajAparat();
      void odczytajPola();
    },
  };
}

/** Nazwa rodzaju elementu pełnym słowem, czytana z wykazu trzynastu rodzajów aparatu zamiast z kodu kontraktu. */
function opiszRodzajAparatu(rodzaj: StudioApparatusKind): string {
  const znaleziony = RODZAJE_APARATU.find((pozycja) => pozycja[0] === rodzaj);
  return znaleziony === undefined ? rodzaj : znaleziony[1];
}

/** Nazwa rodzaju pola pełnym słowem, czytana z wykazu dziewięciu rodzajów pól zamiast z kodu kontraktu. */
function opiszRodzajPola(rodzaj: StudioFieldKind): string {
  const znaleziony = RODZAJE_POL.find((pozycja) => pozycja[0] === rodzaj);
  return znaleziony === undefined ? rodzaj : znaleziony[1];
}

/** Zdanie o elemencie aparatu — zakotwiczenie, treść, cel i pozycje zebrane, złożone z pól, które rdzeń oddał. */
function opiszElement(pozycja: StudioApparatusItem): string {
  const czesci: string[] = [];
  if (pozycja.anchorStart !== undefined) {
    czesci.push(
      pozycja.anchorEnd === undefined
        ? `na znaku ${pozycja.anchorStart}`
        : `znaki ${pozycja.anchorStart}–${pozycja.anchorEnd}`,
    );
  }
  if (pozycja.text !== undefined && pozycja.text !== '') {
    czesci.push(`treść: ${pozycja.text.slice(0, 120)}`);
  }
  if (pozycja.levelsFrom !== undefined || pozycja.levelsTo !== undefined) {
    czesci.push(`poziomy nagłówków ${pozycja.levelsFrom ?? 1}–${pozycja.levelsTo ?? 9}`);
  }
  if (pozycja.entries !== undefined) {
    czesci.push(`pozycji zebranych ${pozycja.entries.length}`);
  }
  if (pozycja.targetId !== undefined) czesci.push(`prowadzi do ${pozycja.targetId}`);
  if (pozycja.targetUrl !== undefined) czesci.push(`prowadzi na ${pozycja.targetUrl}`);
  if (pozycja.citationKey !== undefined) czesci.push(`klucz powołania ${pozycja.citationKey}`);
  if (pozycja.sourceTitle !== undefined) czesci.push(`źródło „${pozycja.sourceTitle}"`);
  if (pozycja.sourceAuthor !== undefined) czesci.push(pozycja.sourceAuthor);
  if (pozycja.sourceYear !== undefined) czesci.push(pozycja.sourceYear);
  return czesci.length === 0 ? 'Rdzeń nie oddał szczegółów tego elementu.' : `${czesci.join(' · ')}.`;
}

/** Zdanie o polu — miejsce, format i wartość policzona, złożone z pól odpowiedzi, które oddał rdzeń panelowi. */
function opiszPole(pole: StudioDocumentField): string {
  const czesci: string[] = [pole.id];
  if (pole.anchorOffset !== undefined) czesci.push(`na znaku ${pole.anchorOffset}`);
  if (pole.format !== undefined && pole.format !== '') czesci.push(`format ${pole.format}`);
  if (pole.expression !== undefined && pole.expression !== '') {
    czesci.push(`wyrażenie ${pole.expression}`);
  }
  if (pole.propertyName !== undefined && pole.propertyName !== '') {
    czesci.push(`właściwość ${pole.propertyName}`);
  }
  czesci.push(
    pole.value === undefined || pole.value === ''
      ? 'wartość NIEPOLICZONA — pole nie było jeszcze odświeżone'
      : `wartość „${pole.value}"`,
  );
  return `${czesci.join(' · ')}.`;
}

function liczba(wartosc: string): number {
  const odczytana = Number.parseInt(wartosc, 10);
  return Number.isFinite(odczytana) ? odczytana : 0;
}

function czesc(tytul: string, elementy: readonly HTMLElement[]): HTMLElement {
  const naglowek = document.createElement('p');
  naglowek.className = 'ms-aparat__tytul';
  naglowek.textContent = tytul;
  const sekcja = document.createElement('section');
  sekcja.className = 'ms-aparat__czesc';
  sekcja.append(naglowek, ...elementy);
  return sekcja;
}

const BRAK_DOKUMENTU =
  'Aparat należy do dokumentu — wczytaj go albo załóż nowy. Komendy studio.apparatus.* ' +
  'i studio.field.* przyjmują identyfikator dokumentu jako pole obowiązkowe.';

const BRAK_WSKAZANEGO =
  'Ta czynność dotyczy elementu wskazanego — naciśnij „Wskaż" przy elemencie w wykazie wyżej.';
