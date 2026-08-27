import {
  StudioBulletSource,
  StudioListKind,
  StudioListNumberFormat,
  StudioTextAlign,
  type StudioAutoReplaceRule,
  type StudioListApplyRequest,
  type StudioListBulletSetRequest,
  type StudioListDefinition,
  type StudioListLevelIndentRequest,
  type StudioListNumberingSetRequest,
  type StudioListRestartRequest,
  type StudioSymbol,
  type StudioSymbolAutoreplaceSetRequest,
  type StudioSymbolInsertRequest,
} from '../../../../shared/contract';
import {
  poleLogiczne,
  poleTekstowe,
  poleWyboru,
  przycisk,
  ustawPozycje,
  utworzWierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import { BEZ_ZMIANY, liczbaPola, poleLiczbowe, tekstPola, wyborPola } from './strona-pola-postaci';

/** Panel list i znaków: wypunktowanie, numeracja, wielopoziomowość i tablica znaków z autozamianą. */

/** Treść żądania bez pola dokumentu — dokument dokłada dopiero warstwa wołająca rdzeń komendą właściwą temu żądaniu. */
type BezDokumentu<T> = Omit<T, 'documentId'>;
/** Treść żądania bez pola dokumentu i bez pól zakresu fragmentu — obie wartości dokłada warstwa wołająca rdzeń. */
type BezZakresu<T> = Omit<T, 'documentId' | 'rangeStart' | 'rangeEnd'>;

/** Czynności panelu list i znaków zlecane oknu: zastosowanie listy, ustawienie punktatora i numeracji poziomu, wznowienie numeracji, zmiana poziomu, wstawienie znaku oraz zapis reguły autozamiany. */
export interface CzynnosciListPanelu {
  naListe(zadanie: BezZakresu<StudioListApplyRequest>): void;
  naPunktator(zadanie: BezDokumentu<StudioListBulletSetRequest>): void;
  naNumeracjeListy(zadanie: BezDokumentu<StudioListNumberingSetRequest>): void;
  naWznowienie(zadanie: BezDokumentu<StudioListRestartRequest>): void;
  naPoziom(zadanie: BezZakresu<StudioListLevelIndentRequest>): void;
  naZnak(zadanie: BezDokumentu<Omit<StudioSymbolInsertRequest, 'offset'>>): void;
  naSzukanieZnaku(fraza: string, tylkoOstatnie: boolean): void;
  naAutozamiane(zadanie: StudioSymbolAutoreplaceSetRequest): void;
  /** Ponowny odczyt definicji list, tablicy znaków i zasad autozamiany. */
  naOdczyt(): void;
}

/** Panel list i znaków wraz z jego sterowaniem: widoczność panelu, pokazanie list dokumentu, tablicy znaków, reguł autozamiany oraz odpowiedzi na żądanie. */
export interface StylPanelList {
  element: HTMLElement;
  przestawWidocznosc(): void;
  widoczny(): boolean;
  /** Definicje list dokumentu — wykaz do wskazania listy i poziomu. */
  pokazListy(listy: readonly StudioListDefinition[]): void;
  /** Tablica znaków wraz z grupami. */
  pokazZnaki(znaki: readonly StudioSymbol[], grupy: readonly string[]): void;
  /** Zasady autozamiany skrótów na znaki. */
  pokazAutozamiany(zasady: readonly StudioAutoReplaceRule[]): void;
  pokazOdpowiedz(tresc: string, powodzenie: boolean): void;
}

export function utworzStylPanelList(czynnosci: CzynnosciListPanelu): StylPanelList {
  const odpowiedz = utworzWierszOdpowiedzi();

  /* ── Lista na fragmencie ─────────────────────────────────────────────────── */

  const rodzajListy = poleWyboru(
    {
      etykieta: 'Rodzaj listy',
      opis: '„Bez listy" zdejmuje listę z zaznaczonego fragmentu.',
    },
    [
      { wartosc: StudioListKind.Bullet, etykieta: 'Wypunktowanie' },
      { wartosc: StudioListKind.Number, etykieta: 'Numeracja' },
      { wartosc: StudioListKind.Multilevel, etykieta: 'Lista wielopoziomowa' },
      { wartosc: StudioListKind.None, etykieta: 'Bez listy — zdejmij' },
    ],
  );

  const listaZastana = poleWyboru(
    {
      etykieta: 'Lista zastana, do której fragment dołączyć',
      opis: 'Puste zakłada listę nową. Dołączenie utrzymuje ciągłość numeracji.',
    },
    [BEZ_ZMIANY],
  );

  const poziomListy = poleLiczbowe({ etykieta: 'Poziom listy' }, { dolna: 1, gorna: 9 });
  const formatNumeracji = poleWyboru({ etykieta: 'Format numeracji' }, FORMATY_NUMERACJI);
  const startNumeracji = poleLiczbowe({ etykieta: 'Punkt startu numeracji' }, { dolna: 0 });

  const zastosujListe = przycisk('Zastosuj listę do zaznaczenia', 'dn-btn dn-btn--sm dn-btn--sygnal');
  zastosujListe.dataset['czynnosc'] = 'zastosuj-liste';
  zastosujListe.title = 'Idzie komendą studio.list.apply na zaznaczonym fragmencie.';
  zastosujListe.addEventListener('click', () => {
    const zadanie: BezZakresu<StudioListApplyRequest> = {
      kind: rodzajListy.kontrolka.value as StudioListKind,
    };
    const zastana = wyborPola(listaZastana.kontrolka);
    if (zastana !== undefined) zadanie.listId = zastana;
    const poziom = liczbaPola(poziomListy.kontrolka);
    if (poziom !== undefined) zadanie.level = poziom;
    const format = wyborPola(formatNumeracji.kontrolka);
    if (format !== undefined) zadanie.numberFormat = format as StudioListNumberFormat;
    const start = liczbaPola(startNumeracji.kontrolka);
    if (start !== undefined) zadanie.startAt = start;
    czynnosci.naListe(zadanie);
  });

  /* ── Punktator poziomu ───────────────────────────────────────────────────── */

  const zrodloZnaku = poleWyboru({ etykieta: 'Źródło znaku wypunktowania' }, [
    BEZ_ZMIANY,
    { wartosc: StudioBulletSource.Character, etykieta: 'Znak z klawiatury' },
    { wartosc: StudioBulletSource.Symbol, etykieta: 'Symbol specjalny' },
    { wartosc: StudioBulletSource.Icon, etykieta: 'Ikona z biblioteki modułu Design' },
    { wartosc: StudioBulletSource.Image, etykieta: 'Obraz własny' },
  ]);
  const znakPunktatora = poleTekstowe({ etykieta: 'Znak wypunktowania', podpowiedz: '•' });
  const zasobPunktatora = poleTekstowe({
    etykieta: 'Zasób ikony albo obrazu wypunktowania',
    opis: 'Identyfikator zasobu magazynu rdzenia — dla źródła „ikona" albo „obraz własny".',
  });
  const wcieciePoziomu = poleLiczbowe(
    { etykieta: 'Wcięcie poziomu w milimetrach' },
    { dolna: 0, gorna: 200, krok: 0.5 },
  );
  const odstepZnaku = poleLiczbowe(
    { etykieta: 'Odstęp znaku od tekstu w milimetrach' },
    { dolna: 0, gorna: 100, krok: 0.5 },
  );
  const wyrownanieZnaku = poleWyboru({ etykieta: 'Wyrównanie znaku albo numeru' }, [
    BEZ_ZMIANY,
    { wartosc: StudioTextAlign.Left, etykieta: 'Do lewej' },
    { wartosc: StudioTextAlign.Center, etykieta: 'Wyśrodkowanie' },
    { wartosc: StudioTextAlign.Right, etykieta: 'Do prawej' },
    { wartosc: StudioTextAlign.Justify, etykieta: 'Wyjustowanie' },
  ]);

  /** Lista i poziom, których dotyczą nastawy punktatora i numeracji poziomu. */
  function listaIPoziom(): { listId: string; level: number } | null {
    const lista = wyborPola(listaZastana.kontrolka);
    const poziom = liczbaPola(poziomListy.kontrolka);
    if (lista === undefined) {
      odpowiedz.pokaz(
        'Nastawa poziomu dotyczy LISTY ZASTANEJ — wskaż ją w wykazie. Lista jeszcze niezałożona ' +
          'nie ma poziomu, który dałoby się przestawić; najpierw zastosuj listę do fragmentu.',
        false,
      );
      return null;
    }
    if (poziom === undefined) {
      odpowiedz.pokaz('Podaj poziom listy — nastawa dotyczy jednego poziomu, nie całej listy.', false);
      return null;
    }
    return { listId: lista, level: poziom };
  }

  const zapiszPunktator = przycisk('Zapisz znak wypunktowania poziomu', 'dn-btn dn-btn--sm dn-btn--sygnal');
  zapiszPunktator.dataset['czynnosc'] = 'zapisz-punktator';
  zapiszPunktator.title = 'Idzie komendą studio.list.bullet.set.';
  zapiszPunktator.addEventListener('click', () => {
    const wskazanie = listaIPoziom();
    if (wskazanie === null) return;
    const zadanie: BezDokumentu<StudioListBulletSetRequest> = { ...wskazanie };
    const zrodlo = wyborPola(zrodloZnaku.kontrolka);
    if (zrodlo !== undefined) zadanie.bulletSource = zrodlo as StudioBulletSource;
    const znak = tekstPola(znakPunktatora.kontrolka);
    if (znak !== undefined) zadanie.bulletCharacter = znak;
    const zasob = tekstPola(zasobPunktatora.kontrolka);
    if (zasob !== undefined) zadanie.bulletAssetId = zasob;
    const wciecie = liczbaPola(wcieciePoziomu.kontrolka);
    if (wciecie !== undefined) zadanie.indentMm = wciecie;
    const odstep = liczbaPola(odstepZnaku.kontrolka);
    if (odstep !== undefined) zadanie.hangingMm = odstep;
    const wyrownanie = wyborPola(wyrownanieZnaku.kontrolka);
    if (wyrownanie !== undefined) zadanie.align = wyrownanie as StudioTextAlign;
    czynnosci.naPunktator(zadanie);
  });

  /* ── Numeracja poziomu ───────────────────────────────────────────────────── */

  const wzorNumeru = poleTekstowe({
    etykieta: 'Wzór numeru poziomu',
    podpowiedz: '%1.%2.',
    opis:
      'Wzór składa numer z poziomów: „%1.%2." daje 1.1, 1.2, 2.1. Tym się buduje numerację ' +
      'prawniczą wielopoziomową typu 1.1.2.',
  });

  const zapiszNumeracjeListy = przycisk(
    'Zapisz format numeracji poziomu',
    'dn-btn dn-btn--sm dn-btn--sygnal',
  );
  zapiszNumeracjeListy.dataset['czynnosc'] = 'zapisz-numeracje-listy';
  zapiszNumeracjeListy.title = 'Idzie komendą studio.list.numbering.set.';
  zapiszNumeracjeListy.addEventListener('click', () => {
    const wskazanie = listaIPoziom();
    if (wskazanie === null) return;
    const zadanie: BezDokumentu<StudioListNumberingSetRequest> = { ...wskazanie };
    const format = wyborPola(formatNumeracji.kontrolka);
    if (format !== undefined) zadanie.numberFormat = format as StudioListNumberFormat;
    const wzor = tekstPola(wzorNumeru.kontrolka);
    if (wzor !== undefined) zadanie.pattern = wzor;
    const start = liczbaPola(startNumeracji.kontrolka);
    if (start !== undefined) zadanie.startAt = start;
    czynnosci.naNumeracjeListy(zadanie);
  });

  /* ── Wznowienie i poziom ─────────────────────────────────────────────────── */

  const wznow = przycisk('Wznów numerację w miejscu kursora', 'dn-btn dn-btn--sm dn-btn--duch');
  wznow.dataset['czynnosc'] = 'wznow-numeracje';
  wznow.title = 'Idzie komendą studio.list.restart z miejscem kursora jako punktem wznowienia.';
  wznow.addEventListener('click', () => {
    const lista = wyborPola(listaZastana.kontrolka);
    if (lista === undefined) {
      odpowiedz.pokaz('Wznowienie dotyczy wskazanej listy — wskaż ją w wykazie.', false);
      return;
    }
    const zadanie: BezDokumentu<StudioListRestartRequest> = { listId: lista, offset: 0 };
    const start = liczbaPola(startNumeracji.kontrolka);
    if (start !== undefined) zadanie.startAt = start;
    czynnosci.naWznowienie(zadanie);
  });

  const wglab = przycisk('Zwiększ poziom listy', 'dn-btn dn-btn--sm dn-btn--duch');
  wglab.dataset['czynnosc'] = 'poziom-glebiej';
  wglab.title = 'Idzie komendą studio.list.level.indent — wcięcie i odstęp poziomu idą za nim.';
  wglab.addEventListener('click', () => czynnosci.naPoziom({ step: 1 }));

  const wyzej = przycisk('Zmniejsz poziom listy', 'dn-btn dn-btn--sm dn-btn--duch');
  wyzej.dataset['czynnosc'] = 'poziom-wyzej';
  wyzej.addEventListener('click', () => czynnosci.naPoziom({ step: -1 }));

  /* ── Znaki specjalne ─────────────────────────────────────────────────────── */

  const szukanyZnak = poleTekstowe({
    etykieta: 'Szukaj znaku po nazwie albo kodzie',
    podpowiedz: 'paragraf, półpauza, 00A0',
  });
  const tylkoOstatnie = poleLogiczne({ etykieta: 'Tylko znaki ostatnio użyte' });

  const szukajZnaku = przycisk('Szukaj', 'dn-btn dn-btn--sm dn-btn--zarys');
  szukajZnaku.dataset['czynnosc'] = 'szukaj-znaku';
  szukajZnaku.title = 'Idzie komendą studio.symbol.list.';
  szukajZnaku.addEventListener('click', () => {
    czynnosci.naSzukanieZnaku(szukanyZnak.kontrolka.value, tylkoOstatnie.kontrolka.checked);
  });

  const grupyZnakow = document.createElement('p');
  grupyZnakow.className = 'dn-pole-opis';

  const tablicaZnakow = document.createElement('div');
  tablicaZnakow.className = 'ms-postac__znaki';

  const znakWprost = poleTekstowe({
    etykieta: 'Znak wstawiany wprost',
    opis: 'Alternatywa dla tablicy: sam znak, jego punkt kodowy albo nazwa.',
  });
  const kodZnaku = poleTekstowe({ etykieta: 'Punkt kodowy zapisem szesnastkowym', podpowiedz: '00A7' });
  const nazwaZnaku = poleTekstowe({ etykieta: 'Nazwa znaku', podpowiedz: 'znak paragrafu' });

  const wstawZnak = przycisk('Wstaw znak w miejscu kursora', 'dn-btn dn-btn--sm dn-btn--sygnal');
  wstawZnak.dataset['czynnosc'] = 'wstaw-znak';
  wstawZnak.title = 'Idzie komendą studio.symbol.insert.';
  wstawZnak.addEventListener('click', () => {
    const zadanie: BezDokumentu<Omit<StudioSymbolInsertRequest, 'offset'>> = {};
    const znak = tekstPola(znakWprost.kontrolka);
    if (znak !== undefined) zadanie.character = znak;
    const kod = tekstPola(kodZnaku.kontrolka);
    if (kod !== undefined) zadanie.code = kod;
    const nazwa = tekstPola(nazwaZnaku.kontrolka);
    if (nazwa !== undefined) zadanie.name = nazwa;
    if (Object.keys(zadanie).length === 0) {
      odpowiedz.pokaz(
        'Wskaż znak: samym znakiem, punktem kodowym albo nazwą. Wybranie z tablicy wyżej robi to ' +
          'jednym naciśnięciem.',
        false,
      );
      return;
    }
    czynnosci.naZnak(zadanie);
  });

  /* ── Autozamiana ─────────────────────────────────────────────────────────── */

  const wykazZasad = document.createElement('ul');
  wykazZasad.className = 'ms-postac__wykaz';

  const skrot = poleTekstowe({ etykieta: 'Skrót wpisywany przez Operatora', podpowiedz: '(c)' });
  const zastapienie = poleTekstowe({
    etykieta: 'Co wchodzi w miejsce skrótu',
    opis: 'Puste USUWA zasadę — to jest droga odwrócenia nastawy, nie osobna komenda.',
    podpowiedz: '©',
  });
  const zasadaCzynna = poleLogiczne({ etykieta: 'Zasada obowiązuje' });
  zasadaCzynna.kontrolka.checked = true;

  const zapiszZasade = przycisk('Zapisz zasadę autozamiany', 'dn-btn dn-btn--sm dn-btn--sygnal');
  zapiszZasade.dataset['czynnosc'] = 'zapisz-autozamiane';
  zapiszZasade.title = 'Idzie komendą studio.symbol.autoreplace.set.';
  zapiszZasade.addEventListener('click', () => {
    const wpisywany = tekstPola(skrot.kontrolka);
    if (wpisywany === undefined) {
      odpowiedz.pokaz('Zasada bez skrótu nie ma na co reagować — podaj skrót.', false);
      return;
    }
    const zadanie: StudioSymbolAutoreplaceSetRequest = {
      shortcut: wpisywany,
      enabled: zasadaCzynna.kontrolka.checked,
    };
    const wchodzi = tekstPola(zastapienie.kontrolka);
    if (wchodzi !== undefined) zadanie.replacement = wchodzi;
    czynnosci.naAutozamiane(zadanie);
  });

  const odczytaj = przycisk('Odczytaj listy, znaki i autozamiany', 'dn-btn dn-btn--sm dn-btn--zarys');
  odczytaj.dataset['czynnosc'] = 'odczytaj-listy';
  odczytaj.addEventListener('click', () => czynnosci.naOdczyt());

  /* ── Układ panelu ────────────────────────────────────────────────────────── */

  const element = document.createElement('section');
  element.className = 'ms-postac ms-postac--listy';
  element.dataset['panel'] = 'listy-i-znaki';
  element.hidden = true;
  element.setAttribute('aria-label', 'Listy, punktatory i znaki specjalne');

  element.append(
    grupa('Lista na zaznaczeniu', [
      odczytaj,
      rodzajListy.element,
      listaZastana.element,
      poziomListy.element,
      formatNumeracji.element,
      startNumeracji.element,
      zastosujListe,
      wglab,
      wyzej,
      wznow,
    ]),
    grupa('Znak wypunktowania poziomu', [
      zrodloZnaku.element,
      znakPunktatora.element,
      zasobPunktatora.element,
      wcieciePoziomu.element,
      odstepZnaku.element,
      wyrownanieZnaku.element,
      zapiszPunktator,
    ]),
    grupa('Numeracja poziomu', [wzorNumeru.element, zapiszNumeracjeListy]),
    grupa('Znaki specjalne', [
      szukanyZnak.element,
      tylkoOstatnie.element,
      szukajZnaku,
      grupyZnakow,
      tablicaZnakow,
      znakWprost.element,
      kodZnaku.element,
      nazwaZnaku.element,
      wstawZnak,
    ]),
    grupa('Autozamiana skrótów', [
      wykazZasad,
      skrot.element,
      zastapienie.element,
      zasadaCzynna.element,
      zapiszZasade,
    ]),
    odpowiedz.element,
  );

  return {
    element,

    przestawWidocznosc() {
      element.hidden = !element.hidden;
    },

    widoczny: () => !element.hidden,

    pokazListy(listy) {
      ustawPozycje(listaZastana.kontrolka, [
        BEZ_ZMIANY,
        ...listy.map((lista) => ({
          wartosc: lista.id,
          etykieta:
            `${lista.id} · ${lista.kind} · poziomów ${lista.levels?.length ?? 0}` +
            `${lista.startAt === undefined ? '' : ` · od ${lista.startAt}`}`,
        })),
      ]);
    },

    pokazZnaki(znaki, grupy) {
      grupyZnakow.textContent =
        grupy.length === 0
          ? 'Rdzeń nie podał grup znaków.'
          : `Grupy znaków: ${grupy.join(', ')}.`;
      tablicaZnakow.replaceChildren(
        ...znaki.map((znak) => {
          const kafel = przycisk(znak.character, 'dn-btn dn-btn--sm dn-btn--zarys ms-postac__znak');
          kafel.dataset['znak'] = znak.code;
          kafel.title =
            `${znak.name} · kod ${znak.code}` +
            `${znak.category === undefined ? '' : ` · ${znak.category}`}` +
            `${znak.recentlyUsed === true ? ' · ostatnio użyty' : ''}`;
          kafel.setAttribute('aria-label', `Wstaw znak ${znak.name}`);
          kafel.addEventListener('click', () => {
            // Wybranie z tablicy jedzie kodem — jest jednoznaczny nawet dla znaków wyglądających identycznie.
            czynnosci.naZnak({ code: znak.code, character: znak.character });
          });
          return kafel;
        }),
      );
      if (znaki.length === 0) {
        const puste = document.createElement('p');
        puste.className = 'dn-pole-opis';
        puste.textContent = 'Rdzeń nie oddał ani jednego znaku dla tego szukania.';
        tablicaZnakow.append(puste);
      }
    },

    pokazAutozamiany(zasady) {
      wykazZasad.replaceChildren(
        ...zasady.map((zasada) => {
          const pozycja = document.createElement('li');
          pozycja.className = 'ms-postac__zasada';
          pozycja.dataset['skrot'] = zasada.shortcut;
          pozycja.textContent =
            `„${zasada.shortcut}" → „${zasada.replacement}" · ` +
            `${zasada.enabled ? 'obowiązuje' : 'wyłączona'}` +
            `${zasada.builtin === true ? ' · fabryczna' : ' · własna'}`;
          const wybierz = przycisk('Wskaż', 'dn-btn dn-btn--sm dn-btn--zarys');
          wybierz.addEventListener('click', () => {
            skrot.kontrolka.value = zasada.shortcut;
            zastapienie.kontrolka.value = zasada.replacement;
            zasadaCzynna.kontrolka.checked = zasada.enabled;
          });
          pozycja.append(wybierz);
          return pozycja;
        }),
      );
    },

    pokazOdpowiedz: (tresc, powodzenie) => odpowiedz.pokaz(tresc, powodzenie),
  };
}

/** Formaty numeracji poziomu wraz z prawniczą wielopoziomową — wykaz pozycji do wyboru w polu formatu numeracji panelu list. */
const FORMATY_NUMERACJI = [
  BEZ_ZMIANY,
  { wartosc: StudioListNumberFormat.Arabic, etykieta: 'Cyfry arabskie — 1, 2, 3' },
  { wartosc: StudioListNumberFormat.RomanUpper, etykieta: 'Cyfry rzymskie wielkie — I, II, III' },
  { wartosc: StudioListNumberFormat.RomanLower, etykieta: 'Cyfry rzymskie małe — i, ii, iii' },
  { wartosc: StudioListNumberFormat.LetterUpper, etykieta: 'Litery wielkie — A, B, C' },
  { wartosc: StudioListNumberFormat.LetterLower, etykieta: 'Litery małe — a, b, c' },
  { wartosc: StudioListNumberFormat.Legal, etykieta: 'Prawnicza wielopoziomowa — 1.1.2' },
];

/** Grupa pól panelu złożona z nagłówka tytułowego oraz przekazanej zawartości, ułożona w jednym bloku pionowym. */
function grupa(tytul: string, zawartosc: readonly HTMLElement[]): HTMLElement {
  const naglowek = document.createElement('h4');
  naglowek.className = 'ms-postac__tytul';
  naglowek.textContent = tytul;

  const element = document.createElement('div');
  element.className = 'ms-postac__grupa';
  element.append(naglowek, ...zawartosc);
  return element;
}
