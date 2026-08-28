import {
  StudioAnchorKind,
  StudioObjectKind,
  StudioObjectSource,
  StudioShapeKind,
  StudioTextWrap,
  type StudioActionBalance,
  type StudioDocumentObject,
} from '../../../../shared/contract';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import {
  poleLiczbowe,
  poleTekstowe,
  przycisk,
  przyciskBezKomendy,
  utworzWierszOdpowiedzi,
  wybor,
} from '../../modele/kontrolki-formularza';
import type { Wynik } from '../../protokol/kanal';
import type { ObiektZrodlo } from './obiekt-zrodlo';
import type { StanStudio } from './stan-studio';
import { wstawieniaOpiszBilans } from './zrodlo-wstawien-studio';

/**
 * Warsztat obiektów — obrazy, kształty, ikony, pola tekstowe, logo. Kształt
 * i ikona przychodzą z modułu Design; wykres rdzeń odmawia i odmowa jest
 * widoczna przed próbą wstawienia.
 */
export interface ObiektPanel {
  element: HTMLElement;
  przestawWidocznosc(): void;
  odswiez(): void;
}

/**
 * Pięć rodzajów obiektu, które rdzeń naprawdę wstawia, wraz z nazwą widoczną
 * dla Operatora w wykazie obiektów dokumentu.
 */
const RODZAJE: readonly (readonly [string, string])[] = [
  [StudioObjectKind.Image, 'obraz'],
  [StudioObjectKind.Shape, 'kształt'],
  [StudioObjectKind.Icon, 'ikona'],
  [StudioObjectKind.Textbox, 'pole tekstowe'],
  [StudioObjectKind.Logo, 'logo'],
];

/**
 * Skąd obiekt pochodzi — siedem dróg z kontraktu, od pliku wskazanego przez
 * Operatora po obiekt narysowany wprost w dokumencie.
 */
const ZRODLA: readonly (readonly [string, string])[] = [
  ['', 'bez wskazania źródła'],
  [StudioObjectSource.File, 'plik wskazany przez Operatora'],
  [StudioObjectSource.CoreAsset, 'magazyn zasobów rdzenia'],
  [StudioObjectSource.DesignModule, 'moduł Design'],
  [StudioObjectSource.PhotoBank, 'baza zdjęciowa'],
  [StudioObjectSource.LibraryFile, 'plik Biblioteki Library'],
  [StudioObjectSource.Web, 'strona z sieci'],
  [StudioObjectSource.Drawn, 'narysowany w dokumencie'],
];

/**
 * Osiem kształtów z kontraktu; rysuje je moduł Design, a Studio jedynie
 * wskazuje który z nich osadzić w dokumencie.
 */
const KSZTALTY: readonly (readonly [string, string])[] = [
  ['', 'bez kształtu'],
  [StudioShapeKind.Rectangle, 'prostokąt'],
  [StudioShapeKind.RoundedRectangle, 'prostokąt zaokrąglony'],
  [StudioShapeKind.Ellipse, 'elipsa'],
  [StudioShapeKind.Line, 'linia'],
  [StudioShapeKind.Arrow, 'strzałka'],
  [StudioShapeKind.Polygon, 'wielokąt'],
  [StudioShapeKind.Star, 'gwiazda'],
  [StudioShapeKind.Callout, 'dymek'],
];

/**
 * Siedem sposobów opływania tekstem wokół obiektu, od wstawienia w wierszu
 * tekstu po ukrycie obiektu za tekstem.
 */
const OPLYWANIE: readonly (readonly [string, string])[] = [
  ['', 'bez zmiany opływania'],
  [StudioTextWrap.Inline, 'w wierszu tekstu'],
  [StudioTextWrap.Square, 'prostokątem'],
  [StudioTextWrap.Tight, 'ściśle wokół kształtu'],
  [StudioTextWrap.Through, 'przez kształt'],
  [StudioTextWrap.TopAndBottom, 'góra i dół'],
  [StudioTextWrap.BehindText, 'za tekstem'],
  [StudioTextWrap.InFrontOfText, 'przed tekstem'],
];

const ZAKOTWICZENIE: readonly (readonly [string, string])[] = [
  ['', 'bez zmiany zakotwiczenia'],
  [StudioAnchorKind.Character, 'do znaku'],
  [StudioAnchorKind.Paragraph, 'do akapitu'],
  [StudioAnchorKind.Page, 'do strony'],
];

const ODMOWA_WYKRESU =
  'Wykresu rdzeń NIE UMIE wstawić do dokumentu Studia: rachunku wykresu po stronie Studia nie ma ' +
  'i komenda studio.object.insert odmawia rodzaju „wykres" nazwanym powodem. Brak jest po stronie ' +
  'rdzenia, nie okna. Droga, która działa: złóż wykres w module Design, a potem osadź go tutaj jako ' +
  'obiekt ze źródłem „moduł Design" i wskazaniem węzła — panel przyjmie go jak każdy inny obraz.';

export function utworzObiektPanel(stan: StanStudio, zrodlo: ObiektZrodlo): ObiektPanel {
  const odpowiedz = utworzWierszOdpowiedzi();

  let wskazany = '';
  let obiekty: readonly StudioDocumentObject[] = [];

  /* ── Wstawienie ────────────────────────────────────────────────────────── */

  const rodzaj = wybor('Rodzaj obiektu', RODZAJE.map((pozycja) => pozycja));
  const zrodloObiektu = wybor('Skąd obiekt pochodzi', ZRODLA.map((pozycja) => pozycja));
  const ksztalt = wybor('Rodzaj kształtu', KSZTALTY.map((pozycja) => pozycja));

  const wskazanie = poleTekstowe({
    etykieta: 'Wskazanie źródła',
    podpowiedz: 'ścieżka, zasób magazynu, węzeł Designu albo plik Biblioteki',
    opis:
      'Pole idzie do rdzenia wedle wybranego źródła: ścieżka jako path, zasób jako assetId, węzeł ' +
      'Designu jako designNodeId, plik Biblioteki jako libraryFileId, adres strony jako sourceUrl. ' +
      'Klient dysku nie czyta — podaje wskazanie, a treść wciąga rdzeń.',
  });
  const nazwaIkony = poleTekstowe({
    etykieta: 'Nazwa ikony w bibliotece modułu Design',
    podpowiedz: 'szuka jej design.icon.library.search',
  });
  const szerokosc = poleLiczbowe('Szerokość w milimetrach', '');
  const wysokosc = poleLiczbowe('Wysokość w milimetrach', '');
  const tekstZastepczy = poleTekstowe({
    etykieta: 'Tekst zastępczy',
    podpowiedz: 'co obraz przedstawia — dla czytnika ekranu i dla wydania bez obrazów',
  });
  const tekstWewnatrz = poleTekstowe({
    etykieta: 'Tekst wewnątrz kształtu albo pola tekstowego',
    podpowiedz: '',
  });
  const podpis = poleTekstowe({ etykieta: 'Podpis pod obiektem', podpowiedz: '' });
  const wypelnienie = poleTekstowe({
    etykieta: 'Wypełnienie zapisem szesnastkowym',
    podpowiedz: 'np. #f4f4f4',
  });
  const obrys = poleTekstowe({
    etykieta: 'Barwa obrysu zapisem szesnastkowym',
    podpowiedz: 'np. #333333',
  });
  const oplywanieWstawienia = wybor('Opływanie tekstem', OPLYWANIE.map((pozycja) => pozycja));
  const zakotwiczenieWstawienia = wybor(
    'Zakotwiczenie',
    ZAKOTWICZENIE.map((pozycja) => pozycja),
  );

  const wstaw = przycisk('Wstaw obiekt w miejsce kursora', 'dn-btn dn-btn--sm dn-btn--sygnal');
  wstaw.dataset['czynnosc'] = 'wstaw-obiekt';
  wstaw.addEventListener('click', () => void wstawObiekt());

  /* ── Wykaz obiektów ────────────────────────────────────────────────────── */

  const zawezenie = wybor('Zawężenie wykazu do rodzaju', [
    ['', 'wszystkie rodzaje'],
    ...RODZAJE,
  ]);
  zawezenie.addEventListener('change', () => void odczytajObiekty());

  const wykaz = document.createElement('ul');
  wykaz.className = 'ms-obiekt__wykaz';
  wykaz.setAttribute('aria-label', 'Obiekty osadzone w dokumencie');

  /* ── Postać obiektu wskazanego ─────────────────────────────────────────── */

  const zachowajProporcje = document.createElement('input');
  zachowajProporcje.type = 'checkbox';
  zachowajProporcje.className = 'dn-przelacznik';
  zachowajProporcje.checked = true;
  zachowajProporcje.setAttribute('aria-label', 'Zachowaj proporcje przy zmianie rozmiaru');
  const etykietaProporcji = document.createElement('label');
  etykietaProporcji.className = 'ms-obiekt__zawezenie';
  etykietaProporcji.append(zachowajProporcje, document.createTextNode('zachowaj proporcje'));

  const polozenieX = poleLiczbowe('Położenie od lewej krawędzi w milimetrach', '');
  const polozenieY = poleLiczbowe('Położenie od górnej krawędzi w milimetrach', '');
  const warstwa = poleLiczbowe('Warstwa — im większa, tym wyżej', '');
  const obrot = poleLiczbowe('Obrót w stopniach', '');
  const grubosc = poleLiczbowe('Grubość obrysu w punktach', '');
  const cieniowanie = document.createElement('input');
  cieniowanie.type = 'checkbox';
  cieniowanie.className = 'dn-przelacznik';
  cieniowanie.setAttribute('aria-label', 'Cieniowanie obiektu');
  const etykietaCienia = document.createElement('label');
  etykietaCienia.className = 'ms-obiekt__zawezenie';
  etykietaCienia.append(cieniowanie, document.createTextNode('cieniowanie obiektu'));

  const ustawPostac = przycisk('Ustaw postać obiektu', 'dn-btn dn-btn--sm dn-btn--atrament');
  ustawPostac.dataset['czynnosc'] = 'postac-obiektu';
  ustawPostac.addEventListener('click', () => void przestawPostac());

  const usun = przycisk('Usuń obiekt wskazany', 'dn-btn dn-btn--sm dn-btn--duch');
  usun.dataset['czynnosc'] = 'usun-obiekt';
  usun.addEventListener('click', () => void usunObiekt());

  /* ── Czynności ─────────────────────────────────────────────────────────── */

  function idDokumentu(): string | null {
    const dokument = stan.dokument();
    if (dokument === null) {
      odpowiedz.pokaz(BRAK_DOKUMENTU, false);
      return null;
    }
    return dokument.id;
  }

  function idObiektu(): string | null {
    if (wskazany === '') {
      odpowiedz.pokaz(BRAK_WSKAZANEGO, false);
      return null;
    }
    return wskazany;
  }

  function miejsceKursora(): number {
    const zakres = stan.zaznaczenie();
    return zakres === null ? stan.trescRobocza().length : zakres.koniec;
  }

  /** Wskazanie źródła w polu, które rdzeń dla tego źródła czyta. */
  function poleWskazania(): Record<string, string> {
    const wartosc = wskazanie.kontrolka.value.trim();
    if (wartosc === '') return {};
    switch (zrodloObiektu.value) {
      case StudioObjectSource.CoreAsset:
        return { assetId: wartosc };
      case StudioObjectSource.DesignModule:
        return { designNodeId: wartosc };
      case StudioObjectSource.LibraryFile:
        return { libraryFileId: wartosc };
      case StudioObjectSource.Web:
        return { sourceUrl: wartosc };
      default:
        return { path: wartosc };
    }
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

  async function wstawObiekt(): Promise<void> {
    const dokument = idDokumentu();
    if (dokument === null) return;
    const wynik = await zrodlo.wstaw({
      documentId: dokument,
      offset: miejsceKursora(),
      kind: rodzaj.value as StudioObjectKind,
      ...(zrodloObiektu.value === ''
        ? {}
        : { source: zrodloObiektu.value as StudioObjectSource }),
      ...poleWskazania(),
      ...(ksztalt.value === '' ? {} : { shapeKind: ksztalt.value as StudioShapeKind }),
      ...(nazwaIkony.kontrolka.value.trim() === ''
        ? {}
        : { iconName: nazwaIkony.kontrolka.value.trim() }),
      ...(liczba(szerokosc.value) > 0 ? { widthMm: liczba(szerokosc.value) } : {}),
      ...(liczba(wysokosc.value) > 0 ? { heightMm: liczba(wysokosc.value) } : {}),
      ...(tekstZastepczy.kontrolka.value.trim() === ''
        ? {}
        : { altText: tekstZastepczy.kontrolka.value.trim() }),
      ...(tekstWewnatrz.kontrolka.value.trim() === ''
        ? {}
        : { innerText: tekstWewnatrz.kontrolka.value.trim() }),
      ...(podpis.kontrolka.value.trim() === '' ? {} : { caption: podpis.kontrolka.value.trim() }),
      ...(oplywanieWstawienia.value === ''
        ? {}
        : { wrap: oplywanieWstawienia.value as StudioTextWrap }),
      ...(zakotwiczenieWstawienia.value === ''
        ? {}
        : { anchor: zakotwiczenieWstawienia.value as StudioAnchorKind }),
      ...(wypelnienie.kontrolka.value.trim() === ''
        ? {}
        : { fillColor: wypelnienie.kontrolka.value.trim() }),
      ...(obrys.kontrolka.value.trim() === ''
        ? {}
        : { strokeColor: obrys.kontrolka.value.trim() }),
    });
    if (!przyjalSie('Wstawienie obiektu', wynik)) return;
    if (wynik.wynik === undefined) return;
    wskazany = wynik.wynik.object.id;
    opiszSkutek(
      `Obiekt ${opiszRodzaj(wynik.wynik.object.kind)} wstawiony jako ${wynik.wynik.object.id}`,
      wynik.wynik.balance,
      wynik.wynik.actionId,
      opiszObiekt(wynik.wynik.object),
    );
    void odczytajObiekty();
  }

  async function przestawPostac(): Promise<void> {
    const dokument = idDokumentu();
    const obiekt = idObiektu();
    if (dokument === null || obiekt === null) return;
    const wynik = await zrodlo.postac({
      documentId: dokument,
      objectId: obiekt,
      ...(liczba(szerokosc.value) > 0 ? { widthMm: liczba(szerokosc.value) } : {}),
      ...(liczba(wysokosc.value) > 0 ? { heightMm: liczba(wysokosc.value) } : {}),
      keepAspect: zachowajProporcje.checked,
      ...(oplywanieWstawienia.value === ''
        ? {}
        : { wrap: oplywanieWstawienia.value as StudioTextWrap }),
      ...(zakotwiczenieWstawienia.value === ''
        ? {}
        : { anchor: zakotwiczenieWstawienia.value as StudioAnchorKind }),
      ...(polozenieX.value === '' ? {} : { positionXMm: liczba(polozenieX.value) }),
      ...(polozenieY.value === '' ? {} : { positionYMm: liczba(polozenieY.value) }),
      ...(warstwa.value === '' ? {} : { zOrder: liczba(warstwa.value) }),
      ...(obrot.value === '' ? {} : { rotationDeg: liczba(obrot.value) }),
      ...(liczba(grubosc.value) > 0 ? { strokeWidthPt: liczba(grubosc.value) } : {}),
      ...(cieniowanie.checked ? { shadow: true } : {}),
      ...(tekstZastepczy.kontrolka.value.trim() === ''
        ? {}
        : { altText: tekstZastepczy.kontrolka.value.trim() }),
      ...(tekstWewnatrz.kontrolka.value.trim() === ''
        ? {}
        : { innerText: tekstWewnatrz.kontrolka.value.trim() }),
      ...(podpis.kontrolka.value.trim() === '' ? {} : { caption: podpis.kontrolka.value.trim() }),
      ...(wypelnienie.kontrolka.value.trim() === ''
        ? {}
        : { fillColor: wypelnienie.kontrolka.value.trim() }),
      ...(obrys.kontrolka.value.trim() === ''
        ? {}
        : { strokeColor: obrys.kontrolka.value.trim() }),
    });
    if (!przyjalSie('Postać obiektu', wynik)) return;
    if (wynik.wynik === undefined) return;
    opiszSkutek(
      'Postać obiektu przestawiona',
      wynik.wynik.balance,
      wynik.wynik.actionId,
      opiszObiekt(wynik.wynik.object),
    );
    void odczytajObiekty();
  }

  async function usunObiekt(): Promise<void> {
    const dokument = idDokumentu();
    const obiekt = idObiektu();
    if (dokument === null || obiekt === null) return;
    const wynik = await zrodlo.usun({ documentId: dokument, objectId: obiekt });
    if (!przyjalSie('Usunięcie obiektu', wynik)) return;
    if (wynik.wynik === undefined) return;
    if (!wynik.wynik.removed) {
      // Odpowiedź „nie usunąłem" to wynik, nie awaria; obiekt zostaje wskazany, bo nadal jest.
      odpowiedz.pokaz(
        `Rdzeń NIE usunął obiektu ${obiekt}. ${wstawieniaOpiszBilans(wynik.wynik.balance)}`,
        false,
      );
      return;
    }
    wskazany = '';
    opiszSkutek(
      `Obiekt ${obiekt} usunięty`,
      wynik.wynik.balance,
      wynik.wynik.actionId,
      '',
    );
    void odczytajObiekty();
  }

  async function odczytajObiekty(): Promise<void> {
    const dokument = stan.dokument();
    if (dokument === null) {
      obiekty = [];
      przerysujWykaz();
      return;
    }
    const wynik = await zrodlo.wykaz({
      documentId: dokument.id,
      ...(zawezenie.value === '' ? {} : { kind: zawezenie.value as StudioObjectKind }),
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu('Odczyt obiektów dokumentu', wynik.blad), false);
      return;
    }
    obiekty = wynik.wynik.objects;
    przerysujWykaz();
  }

  function przyjalSie(nazwa: string, wynik: Wynik<unknown>): boolean {
    if (wynik.udany && wynik.wynik !== undefined) return true;
    odpowiedz.pokaz(opisOdmowyBledu(nazwa, wynik.blad), false);
    return false;
  }

  function przerysujWykaz(): void {
    if (obiekty.length === 0) {
      const puste = document.createElement('li');
      puste.className = 'dn-pole-opis';
      puste.textContent =
        'Dokument nie ma ani jednego obiektu tego rodzaju. To odpowiedź rdzenia, nie brak odczytu.';
      wykaz.replaceChildren(puste);
      return;
    }
    wykaz.replaceChildren(
      ...obiekty.map((obiekt) => {
        const glowa = document.createElement('p');
        glowa.className = 'ms-obiekt__glowa';
        glowa.textContent =
          `${opiszRodzaj(obiekt.kind)} · ${obiekt.id}` +
          (obiekt.locked === true ? ' · POD BLOKADĄ FRAGMENTU' : '');

        const miary = document.createElement('p');
        miary.className = 'dn-pole-opis';
        miary.textContent = opiszObiekt(obiekt);

        const wskaz = przycisk(
          obiekt.id === wskazany ? 'Wskazany' : 'Wskaż',
          'dn-btn dn-btn--sm dn-btn--zarys',
        );
        wskaz.dataset['obiekt'] = obiekt.id;
        wskaz.setAttribute('aria-pressed', String(obiekt.id === wskazany));
        wskaz.addEventListener('click', () => {
          wskazany = obiekt.id;
          przerysujWykaz();
        });

        const pozycja = document.createElement('li');
        pozycja.dataset['obiekt'] = obiekt.id;
        pozycja.dataset['rodzaj'] = obiekt.kind;
        pozycja.dataset['wskazany'] = String(obiekt.id === wskazany);
        pozycja.append(glowa, miary, wskaz);
        return pozycja;
      }),
    );
  }

  /* ── Nakładka ──────────────────────────────────────────────────────────── */

  const oDesignie = document.createElement('p');
  oDesignie.className = 'dn-pole-opis';
  oDesignie.textContent =
    'Kształt i ikona mają zaplecze w module Design: kształt zakłada się tam od razu jako węzły ' +
    'ścieżki, więc da się go dalej edytować, a ikonę szuka biblioteka Designu. Studio wskazuje, co ' +
    'osadzić, i nie prowadzi drugiego rachunku kształtu.';

  const tresc = document.createElement('div');
  tresc.className = 'ms-obiekt__tresc';
  tresc.append(
    czesc('Wstawienie obiektu', [
      rodzaj,
      zrodloObiektu,
      wskazanie.element,
      ksztalt,
      nazwaIkony.element,
      oDesignie,
      szerokosc,
      wysokosc,
      tekstZastepczy.element,
      tekstWewnatrz.element,
      podpis.element,
      wypelnienie.element,
      obrys.element,
      oplywanieWstawienia,
      zakotwiczenieWstawienia,
      wstaw,
      przyciskBezKomendy('Wstaw wykres', ODMOWA_WYKRESU),
    ]),
    czesc('Obiekty dokumentu', [zawezenie, wykaz]),
    czesc('Postać obiektu wskazanego', [
      etykietaProporcji,
      polozenieX,
      polozenieY,
      warstwa,
      obrot,
      grubosc,
      etykietaCienia,
      ustawPostac,
      usun,
    ]),
    odpowiedz.element,
  );

  const nakladka = document.createElement('div');
  nakladka.className = 'ms-obiekt__nakladka';
  nakladka.hidden = true;
  nakladka.append(tresc);

  const wyzwalacz = przycisk('Obiekty ▾', 'dn-btn dn-btn--sm dn-btn--zarys');
  wyzwalacz.dataset['czynnosc'] = 'warsztat-obiektow';
  wyzwalacz.setAttribute('aria-expanded', 'false');
  wyzwalacz.title =
    'Otwiera warsztat obiektów nakładką: obraz z pliku, magazynu, Designu, bazy zdjęciowej albo ' +
    'Biblioteki, kształt, ikona, pole tekstowe i logo, wraz z rozmiarem, opływaniem, warstwą ' +
    'i obrotem. Schodzi drugim naciśnięciem albo Escapem.';
  wyzwalacz.addEventListener('click', () => przestaw(nakladka.hidden));

  const element = document.createElement('section');
  element.className = 'ms-obiekt';
  element.setAttribute('aria-label', 'Warsztat obiektów dokumentu');
  element.append(wyzwalacz, nakladka);
  element.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key === 'Escape' && !nakladka.hidden) przestaw(false);
  });

  function przestaw(otwarta: boolean): void {
    nakladka.hidden = !otwarta;
    wyzwalacz.setAttribute('aria-expanded', otwarta ? 'true' : 'false');
    wyzwalacz.textContent =
      obiekty.length === 0
        ? `Obiekty ${otwarta ? '▴' : '▾'}`
        : `Obiekty: ${obiekty.length} ${otwarta ? '▴' : '▾'}`;
    if (otwarta) void odczytajObiekty();
  }

  przerysujWykaz();

  return {
    element,
    przestawWidocznosc: () => przestaw(nakladka.hidden),
    odswiez: () => {
      if (!nakladka.hidden) void odczytajObiekty();
    },
  };
}

/**
 * Nazwa rodzaju obiektu pełnym słowem — bez kodu i bez skrótu — widoczna dla
 * Operatora w wykazie obiektów.
 */
function opiszRodzaj(rodzaj: StudioObjectKind): string {
  const znaleziony = RODZAJE.find((pozycja) => pozycja[0] === rodzaj);
  return znaleziony === undefined ? rodzaj : znaleziony[1];
}

/**
 * Zdanie o obiekcie: miary w milimetrach, sposób opływania tekstem, warstwa
 * dokumentu i zapis pochodzenia obiektu.
 */
function opiszObiekt(obiekt: StudioDocumentObject): string {
  const czesci: string[] = [];
  if (obiekt.widthMm !== undefined || obiekt.heightMm !== undefined) {
    czesci.push(`${obiekt.widthMm ?? '—'} × ${obiekt.heightMm ?? '—'} mm`);
  }
  if (obiekt.wrap !== undefined) czesci.push(`opływanie ${obiekt.wrap}`);
  if (obiekt.anchor !== undefined) czesci.push(`zakotwiczenie ${obiekt.anchor}`);
  if (obiekt.anchorOffset !== undefined) czesci.push(`na znaku ${obiekt.anchorOffset}`);
  if (obiekt.zOrder !== undefined) czesci.push(`warstwa ${obiekt.zOrder}`);
  if (obiekt.rotationDeg !== undefined) czesci.push(`obrót ${obiekt.rotationDeg}°`);
  if (obiekt.shapeKind !== undefined) czesci.push(`kształt ${obiekt.shapeKind}`);
  if (obiekt.source !== undefined) czesci.push(`źródło ${obiekt.source}`);
  if (obiekt.designNodeId !== undefined) czesci.push(`węzeł Designu ${obiekt.designNodeId}`);
  if (obiekt.libraryFileId !== undefined) czesci.push(`plik Biblioteki ${obiekt.libraryFileId}`);
  if (obiekt.sourceUrl !== undefined) czesci.push(`adres ${obiekt.sourceUrl}`);
  if (obiekt.altText === undefined || obiekt.altText === '') {
    czesci.push('BEZ tekstu zastępczego — wydanie bez obrazów nie powie, co tu jest');
  }
  return czesci.length === 0 ? 'Rdzeń nie oddał żadnej miary tego obiektu.' : `${czesci.join(' · ')}.`;
}

function liczba(wartosc: string): number {
  const odczytana = Number.parseFloat(wartosc);
  return Number.isFinite(odczytana) ? odczytana : 0;
}

function czesc(tytul: string, elementy: readonly HTMLElement[]): HTMLElement {
  const naglowek = document.createElement('p');
  naglowek.className = 'ms-obiekt__tytul';
  naglowek.textContent = tytul;
  const sekcja = document.createElement('section');
  sekcja.className = 'ms-obiekt__czesc';
  sekcja.append(naglowek, ...elementy);
  return sekcja;
}

const BRAK_DOKUMENTU =
  'Obiekt wchodzi do dokumentu — wczytaj go albo załóż nowy. Komenda studio.object.insert przyjmuje ' +
  'identyfikator dokumentu jako pole obowiązkowe.';

const BRAK_WSKAZANEGO =
  'Ta czynność dotyczy obiektu wskazanego — naciśnij „Wskaż" przy obiekcie w wykazie wyżej.';
