import {
  AccessPointStatus,
  ExtensionKind,
  type AccessPoint,
  type Extension,
} from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import {
  poleTekstowe,
  poleWielowierszowe,
  przycisk,
  utworzWierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import { utworzWykazBrakow } from './braki-kontraktu';
import { opiszPole } from './dymek-objasnienia';
import {
  BEZ_ODCZYTU_KATALOGU,
  BRAKI_INTEGRATIONS_HUB,
  KODY_OKIEN,
  NAZWY_OKIEN,
  OKNO_SPOZA_KATALOGU,
} from './etykiety-apps';
import { utworzRameApps } from './rama-okna';
import { narzedziaIntegracji } from './narzedzia-rozszerzen';
import { utworzPrzybornikApps } from './przybornik-apps';
import type { StanRozszerzen } from './stan-rozszerzen';
import { utworzWyborZMenu, wierszWyboru } from './wybor-z-menu';

/**
 * Integrations Hub — okno wiodące integracji zewnętrznych: konektory i serwery
 * MCP w jednym widoku.
 *
 * Widok obejmuje dwa rodzaje pozycji katalogu — serwer MCP i integrację API —
 * bo tak brzmi cel okna. Wtyczki i umiejętności do niego nie należą: są
 * rozszerzeniami powłoki i sposobami wykonania zadań, nie drogami do usług
 * zewnętrznych, i mają swoje miejsce w App Catalogu.
 *
 * „Test połączenia" i „monitor zdrowia" stoją na `access.point.check` i okno
 * mówi wprost, co jest sprawdzane: MOST, którym się do serwera dochodzi, a nie
 * usługa za mostem. Sprawdzenie oddaje stan, czas i korzenie potwierdzone przez
 * punkt; czasu odpowiedzi samej integracji ani liczby jej narzędzi kontrakt nie
 * niesie i okno tego nie zmyśla.
 *
 * Integracja bez punktu dostępu nie jest usterką: pozycja rodzaju `api` bywa
 * dostępna wprost adresem, a `AccessPoint` opisuje most do maszyny. Okno
 * odróżnia „bez mostu" od „most nieznany rdzeniowi", bo to dwa różne stany
 * i drugi jest niespójnością rejestru.
 *
 * Transport MCP da się dziś wpisać wyłącznie do nieprzezroczystej konfiguracji.
 * Okno pozwala to zrobić i nazywa ryzyko: kontrakt tego pola nie nazywa, więc
 * rdzeń i okno mogą rozumieć je inaczej.
 */
export interface OknoIntegrationsHub {
  element: HTMLElement;
  odswiez(): void;
}

/** Rodzaje pozycji, które są drogą do usługi zewnętrznej. */
const RODZAJE_INTEGRACJI: readonly ExtensionKind[] = [ExtensionKind.Mcp, ExtensionKind.Api];

/** Odmiana plakietki stanu punktu dostępu; nazwy z biblioteki `komponenty/`. */
const PLAKIETKI_PUNKTU: Readonly<Record<string, string>> = {
  [AccessPointStatus.Reachable]: 'dn-plakietka dn-plakietka--sukces',
  [AccessPointStatus.Unreachable]: 'dn-plakietka dn-plakietka--blad',
  [AccessPointStatus.Unknown]: 'dn-plakietka',
};

export function utworzOknoIntegrationsHub(stan: StanRozszerzen): OknoIntegrationsHub {
  const kod = KODY_OKIEN.IntegrationsHub;
  const rama = utworzRameApps(kod, NAZWY_OKIEN[kod] ?? kod, 'wiodące');

  const zdrowie = document.createElement('p');
  zdrowie.className = 'mp-zdrowie';

  const wykaz = document.createElement('ul');
  wykaz.className = 'mp-integracje';

  const kodPozycji = poleTekstowe({
    etykieta: 'Kod integracji',
    podpowiedz: 'kod stały między wydaniami',
  });
  const rodzaj = utworzWyborZMenu('Rodzaj integracji', [
    { wartosc: ExtensionKind.Mcp, etykieta: 'serwer MCP' },
    { wartosc: ExtensionKind.Api, etykieta: 'integracja API' },
  ]);
  const adres = poleTekstowe({ etykieta: 'Adres serwera', podpowiedz: 'adres albo ścieżka procesu' });
  const most = utworzWyborZMenu('Most z katalogu punktów dostępu');
  const konfiguracja = poleWielowierszowe(
    { etykieta: 'Konfiguracja integracji (JSON)', podpowiedz: '{"transport":"stdio"}' },
    3,
  );
  const dodaj = przycisk('Dodaj integrację', 'dn-btn dn-btn--sm dn-btn--atrament');

  const odczytaj = przycisk('Odczytaj integracje i mosty', 'dn-btn dn-btn--sm dn-btn--zarys');
  const odpowiedz = utworzWierszOdpowiedzi();

  const granica = document.createElement('p');
  granica.className = 'dn-pole-opis mp-granica';
  granica.textContent = OKNO_SPOZA_KATALOGU;

  const pasek = document.createElement('div');
  pasek.className = 'mp-pasek';
  pasek.append(dodaj);

  rama.akcje.append(
    utworzPrzybornikApps('Transport, protokół, webhooki i zdrowie', narzedziaIntegracji(stan)).element,
    utworzWykazBrakow('Bez drogi w kontrakcie', BRAKI_INTEGRATIONS_HUB, 'extension'),
  );
  rama.tresc.append(
    odczytaj,
    zdrowie,
    wykaz,
    naglowekCzesci('Dodanie integracji'),
    kodPozycji.element,
    wierszWyboru('Rodzaj integracji', rodzaj),
    opiszPole(
      adres.element,
      'Napis podawany rdzeniowi jako źródło pozycji. Dla serwera uruchamianego lokalnie ' +
        'jest to ścieżka procesu, dla usługi sieciowej — jej adres.',
    ),
    opiszPole(
      wierszWyboru('Most z katalogu punktów dostępu', most),
      'Most opisuje, DO CZEGO platforma ma wgląd, i to jego dotyczy sprawdzenie połączenia. ' +
        'Integracja bez mostu jest poprawna, jeżeli usługa stoi wprost pod adresem.',
    ),
    opiszPole(
      konfiguracja.element,
      'Konfiguracja jest dla kontraktu nieprzezroczystym JSON-em. Transport MCP wpisany ' +
        'tutaj zadziała tylko wtedy, gdy rdzeń szuka go pod tą samą nazwą pola — kontrakt ' +
        'tego pola nie nazywa.',
    ),
    pasek,
    odpowiedz.element,
    granica,
  );

  odczytaj.addEventListener('click', () => void odczytajCalosc());
  dodaj.addEventListener('click', () => void wyslijDodanie());

  function integracje(): readonly Extension[] {
    return stan.rozszerzeniaRodzaju(RODZAJE_INTEGRACJI);
  }

  async function odczytajCalosc(): Promise<void> {
    rama.ladowanie('Odczyt integracji i katalogu mostów…');
    odpowiedz.pokaz('Odczyt integracji: żądania wysłane do rdzenia…', true);
    // Dwa odczyty idą równolegle, bo są od siebie niezależne: odmowa jednego
    // nie zabiera drugiego, a integracje da się pokazać także bez katalogu
    // mostów — wtedy wiersz mówi, że mostu nie ma czym rozwiązać.
    await Promise.all([stan.odczytajKatalog(), stan.odczytajPunkty()]);
    const powodKatalogu = stan.powodOdczytu('katalog');
    const powodPunktow = stan.powodOdczytu('punkty');
    odswiez();
    if (powodKatalogu !== '') {
      rama.blad(powodKatalogu);
      odpowiedz.pokaz(powodKatalogu, false);
      return;
    }
    rama.gotowe();
    odpowiedz.pokaz(
      `Odczyt integracji: rdzeń oddał ${integracje().length} integracji ` +
        (powodPunktow === ''
          ? `i ${stan.punkty().length} mostów.`
          : `; katalog mostów odmówił: ${powodPunktow}`),
      powodPunktow === '',
    );
  }

  async function wyslijDodanie(): Promise<void> {
    if (kodPozycji.kontrolka.value === '') {
      odpowiedz.pokaz('Okno nie wysyła pustego kodu integracji — wpisz go.', false);
      return;
    }
    let tresc: unknown;
    const wpisana = konfiguracja.kontrolka.value;
    try {
      tresc = wpisana === '' ? undefined : JSON.parse(wpisana);
    } catch (blad) {
      odpowiedz.pokaz(
        `Konfiguracja nie jest poprawnym zapisem JSON: ${(blad as Error).message}. ` +
          'Żądanie nie zostało wysłane.',
        false,
      );
      return;
    }
    rama.ladowanie('Dodawanie integracji w toku…');
    const wynik = await stan.zrodlo.zainstaluj({
      kod: kodPozycji.kontrolka.value,
      rodzaj: rodzaj.wartosc() as ExtensionKind,
      zrodlo: adres.kontrolka.value,
      idPunktuDostepu: most.wartosc(),
      ...(tresc === undefined ? {} : { konfiguracja: tresc }),
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      const zdanie = opisOdmowy('Dodanie integracji', wynik.blad?.code, wynik.blad?.message);
      rama.blad(zdanie);
      odpowiedz.pokaz(zdanie, false);
      return;
    }
    stan.wchlonPozycje(wynik.wynik.extension);
    rama.gotowe();
    odswiez();
    odpowiedz.pokaz(
      `Rdzeń zarejestrował integrację ${wynik.wynik.extension.name} (kod ` +
        `${wynik.wynik.extension.code}). Pochodzenia nie zamawiano, więc rdzeń przyjął ` +
        'wartość domyślną kontraktu — pozycja Operatora jest wyłączona do świadomego włączenia.',
      true,
    );
  }

  /**
   * Sprawdzenie mostu obsługującego integrację.
   *
   * Wynik wchodzi do wspólnego zbioru punktów, a nie tylko do wiersza: ten sam
   * most bywa obsługą wielu integracji, więc sprawdzenie zrobione raz ma być
   * widoczne przy każdej z nich.
   */
  async function sprawdzMost(pozycja: Extension, punkt: AccessPoint): Promise<void> {
    rama.ladowanie(`Sprawdzanie mostu ${punkt.name}…`);
    const wynik = await stan.punktyDostepu.sprawdz(punkt.id);
    if (!wynik.udany || wynik.wynik === undefined) {
      const zdanie = opisOdmowy('Sprawdzenie mostu', wynik.blad?.code, wynik.blad?.message);
      rama.blad(zdanie);
      odpowiedz.pokaz(zdanie, false);
      return;
    }
    const oddany = wynik.wynik;
    stan.wchlonPunkt({ ...punkt, status: oddany.status, checkedAt: oddany.checkedAt });
    rama.gotowe();
    odswiez();
    const korzenie = oddany.roots ?? [];
    const oKorzeniach =
      korzenie.length === 0
        ? ' Punkt nie potwierdził żadnego korzenia.'
        : ` Potwierdzone korzenie: ${korzenie.join(', ')}.`;
    const oSzczegole =
      oddany.detail === undefined || oddany.detail === '' ? '' : ` Rdzeń podał: ${oddany.detail}.`;
    odpowiedz.pokaz(
      `Sprawdzenie mostu ${punkt.name} obsługującego integrację ${pozycja.name}: stan ` +
        `${oddany.status}.${oKorzeniach}${oSzczegole} Sprawdzeniu podlega most, nie usługa za nim.`,
      oddany.status === AccessPointStatus.Reachable,
    );
  }

  async function przelacz(pozycja: Extension): Promise<void> {
    const czynnosc = pozycja.enabled ? 'Wyłączenie' : 'Włączenie';
    rama.ladowanie(`${czynnosc} integracji ${pozycja.name} w toku…`);
    const wynik = await stan.zrodlo.przelacz(pozycja.id, !pozycja.enabled);
    if (!wynik.udany || wynik.wynik === undefined) {
      const zdanie = opisOdmowy(czynnosc, wynik.blad?.code, wynik.blad?.message);
      rama.blad(zdanie);
      odpowiedz.pokaz(zdanie, false);
      return;
    }
    stan.wchlonPozycje(wynik.wynik.extension);
    rama.gotowe();
    odswiez();
    odpowiedz.pokaz(
      `${czynnosc} integracji ${wynik.wynik.extension.name}: rdzeń oddał stan ` +
        `${wynik.wynik.extension.enabled ? 'włączona' : 'wyłączona'}.`,
      true,
    );
  }

  /**
   * Zbiorczy stan integracji — składany wyłącznie z pól, które rdzeń oddał.
   *
   * Zdanie nie mówi o opóźnieniach ani o błędach integracji, bo tych kontrakt
   * nie niesie. Mówi o tym, co wiadomo: ile integracji jest włączonych i jak
   * odpowiadają ich mosty.
   */
  function zdanieZdrowia(pozycje: readonly Extension[]): string {
    if (pozycje.length === 0) return '';
    const wlaczone = pozycje.filter((pozycja) => pozycja.enabled).length;
    const mosty = new Set(
      pozycje
        .map((pozycja) => pozycja.accessPointId ?? '')
        .filter((identyfikator) => identyfikator !== ''),
    );
    const odpowiadajace = [...mosty].filter(
      (identyfikator) => stan.punkt(identyfikator)?.status === AccessPointStatus.Reachable,
    ).length;
    const oMostach =
      mosty.size === 0
        ? 'żadna nie wskazuje mostu'
        : `mostów ${mosty.size}, odpowiadających przy ostatnim sprawdzeniu ${odpowiadajace}`;
    return (
      `Integracji ${pozycje.length}, włączonych ${wlaczone}; ${oMostach}. ` +
      'Czasu odpowiedzi i liczby błędów integracji kontrakt nie niesie — nie ma ich tu skąd wziąć.'
    );
  }

  function odswiez(): void {
    const pozycje = integracje();
    wykaz.replaceChildren(...pozycje.map(wierszIntegracji));
    zdrowie.textContent = zdanieZdrowia(pozycje);
    zdrowie.hidden = zdrowie.textContent === '';
    most.ustawPozycje([
      { wartosc: '', etykieta: 'bez mostu' },
      ...stan.punkty().map((punkt) => ({
        wartosc: punkt.id,
        etykieta: punkt.name,
        opis: punkt.kind,
      })),
    ]);
    if (rama.faza() === 'blad' || rama.faza() === 'ladowanie') return;
    if (pozycje.length > 0) {
      rama.gotowe();
      return;
    }
    if (!stan.czyKatalogCzytany()) {
      rama.puste(BEZ_ODCZYTU_KATALOGU);
      return;
    }
    rama.puste(
      'Rdzeń odpowiedział na odczyt i nie zna ani jednego serwera MCP ani integracji API. ' +
        'To jest pustka POTWIERDZONA — dodaj pierwszą integrację poniżej.',
    );
  }

  /** Jeden wiersz integracji: rodzaj, stan, most i czynności. */
  function wierszIntegracji(pozycja: Extension): HTMLElement {
    const nazwa = document.createElement('span');
    nazwa.className = 'mp-integracje__nazwa';
    nazwa.textContent = pozycja.name;

    const rodzajPozycji = document.createElement('span');
    rodzajPozycji.className = 'dn-plakietka dn-plakietka--rola';
    rodzajPozycji.textContent = pozycja.kind;

    const stanPozycji = document.createElement('span');
    stanPozycji.className = pozycja.enabled
      ? 'dn-plakietka dn-plakietka--sukces'
      : 'dn-plakietka dn-plakietka--ostrzezenie';
    stanPozycji.textContent = pozycja.enabled ? 'włączona' : 'wyłączona';

    const glowa = document.createElement('div');
    glowa.className = 'mp-integracje__glowa';
    glowa.append(nazwa, rodzajPozycji, stanPozycji);

    const element = document.createElement('li');
    element.className = 'mp-integracje__wiersz';
    element.dataset['integracja'] = pozycja.id;
    element.append(glowa, wierszMostu(pozycja));

    const czynnosci = document.createElement('div');
    czynnosci.className = 'mp-integracje__czynnosci';
    czynnosci.append(
      przyciskWiersza(pozycja.enabled ? 'Wyłącz' : 'Włącz', 'dn-btn dn-btn--zarys dn-btn--sm', () =>
        void przelacz(pozycja),
      ),
    );
    const punkt = pozycja.accessPointId === undefined ? null : stan.punkt(pozycja.accessPointId);
    if (punkt !== null) {
      czynnosci.append(
        przyciskWiersza('Testuj most', 'dn-btn dn-btn--zarys dn-btn--sm', () =>
          void sprawdzMost(pozycja, punkt),
        ),
      );
    }
    czynnosci.append(
      przyciskWiersza('Konsola i inspektor', 'dn-btn dn-btn--duch dn-btn--sm', () =>
        stan.wybierz(pozycja.id),
      ),
      przyciskWiersza('Uprawnienia i zaufanie', 'dn-btn dn-btn--duch dn-btn--sm', () =>
        stan.wybierz(pozycja.id),
      ),
    );
    element.append(czynnosci);
    return element;
  }

  /** Wiersz mostu: trzy rozróżnione stany, nie jeden. */
  function wierszMostu(pozycja: Extension): HTMLElement {
    const element = document.createElement('p');
    element.className = 'mp-integracje__most';
    const identyfikator = pozycja.accessPointId ?? '';
    if (identyfikator === '') {
      element.textContent =
        'Bez mostu — integracja stoi wprost pod swoim adresem. To poprawny stan, nie brak.';
      return element;
    }
    const punkt = stan.punkt(identyfikator);
    if (punkt === null) {
      element.textContent = stan.czyPunktyCzytane()
        ? `Most ${identyfikator} wskazany przez integrację NIE STOI w katalogu punktów ` +
          'dostępu — to niespójność rejestru, nie stan usługi.'
        : `Most ${identyfikator} — katalog punktów dostępu nie był jeszcze czytany, ` +
          'więc nie ma czym rozwiązać tego wskazania.';
      return element;
    }
    const plakietka = document.createElement('span');
    plakietka.className = PLAKIETKI_PUNKTU[punkt.status] ?? 'dn-plakietka';
    plakietka.textContent = punkt.status;
    const opis = document.createElement('span');
    opis.className = 'mp-integracje__opis-mostu';
    opis.textContent =
      `most ${punkt.name}` +
      (punkt.roots.length === 0 ? ', bez korzeni' : `, korzenie: ${punkt.roots.join(', ')}`);
    element.append(plakietka, opis);
    return element;
  }

  odswiez();
  return { element: rama.element, odswiez };
}

function przyciskWiersza(etykieta: string, klasa: string, naNacisniecie: () => void): HTMLElement {
  const element = document.createElement('button');
  element.type = 'button';
  element.className = klasa;
  element.textContent = etykieta;
  element.addEventListener('click', naNacisniecie);
  return element;
}

/** Nagłówek części okna. */
function naglowekCzesci(tresc: string): HTMLElement {
  const element = document.createElement('p');
  element.className = 'mp-czesc__tytul';
  element.textContent = tresc;
  return element;
}
