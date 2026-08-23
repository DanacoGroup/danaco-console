import {
  Command,
  ConfigScope,
  IsolationContextKind,
  IsolationTechnicalScope,
  NARZEDZIA_MODELU,
  type IsolationPolicy,
  type IsolationPolicyPreviewRequest,
} from '../../../shared/contract';
import { utworzDymekObjasnienia } from '../komponenty/dymek';
import { opisOdmowyBledu } from '../komponenty/odmowa';
import { utworzStanTresci } from '../komponenty/stan-tresci';
import { utworzPunktWidzenia } from './efektywna-punkt-widzenia';
import type { ObszarIzolacji, ZaleznosciObszaru } from './obszary';

/**
 * Obszar „Polityka efektywna" okna Punktów Izolacji: pokazuje, co obowiązuje
 * w tej chwili dla jedenastu punktów izolacji i skąd każda wartość pochodzi.
 *
 * Wartość zapisana świadomie i wartość domyślna muszą wyglądać inaczej,
 * ponieważ stan wyjściowy platformy to pełna swoboda operacyjna — egzekutor
 * odrzuca naruszenie tylko przy punkcie włączonym, a wartość wyłączona niczego
 * nie ogranicza. Pokazanie domyślnej tak samo jak zapisu Operatora
 * sugerowałoby ochronę, której nikt nie włączył.
 *
 * Pochodzenie czytamy z odpowiedzi rdzenia na `isolation.policy.preview`:
 * obecne `policy.origin` znaczy zapis na wskazanym poziomie, brak `origin`
 * znaczy wartość z rejestru definicji. Stan `nieznane` wchodzi wtedy, gdy
 * odpowiedź w ogóle nie niesie przełącznika danego punktu.
 *
 * Żądanie podglądu musi nieść punkt widzenia — sesję, warstwę, zasięg i okno.
 * Przy pustym żądaniu rdzeń schodzi kolejką do poziomu globalnego, więc okno
 * pokazałoby politykę całej platformy pod nazwą „efektywna".
 *
 * Czwarty stan `brak-odczytu` jest wyłącznie kliencki i nie wolno go zlewać
 * z wartością domyślną: oba prowadzą do przeciwnych wniosków o bezpieczeństwie.
 */

/** Trzy wartości pochodzenia z rdzenia (`Pochodzenie`) i czwarta, wyłącznie kliencka: brak odczytu. */
type StanPochodzenia = 'zapis' | 'domyslna' | 'nieznane' | 'brak-odczytu';

interface OpisPlakietki {
  etykieta: string;
  klasa: string;
}

/** Wygląd plakietki dla każdego z czterech stanów — każdy ma własny. */
const PLAKIETKI_POCHODZENIA: Readonly<Record<StanPochodzenia, OpisPlakietki>> = {
  zapis: { etykieta: 'Zapis Operatora', klasa: 'dn-plakietka--sukces' },
  domyslna: { etykieta: 'Wartość domyślna', klasa: 'dn-plakietka--informacja' },
  nieznane: { etykieta: 'Nieznane rejestrowi', klasa: 'dn-plakietka--ostrzezenie' },
  'brak-odczytu': { etykieta: 'Brak odczytu', klasa: 'dn-plakietka--blad' },
};

const POWOD_BRAKU_ODCZYTU =
  'Rdzeń nie odpowiedział na isolation.policy.preview — komenda jest w kontrakcie i została wysłana, ' +
  'ale odczyt się nie powiódł. Wartość i pochodzenie każdego punktu poniżej są nieznane — to nie jest ' +
  'wartość domyślna, to brak możliwości sprawdzenia.';

type GrupaPunktu = 'kontekst' | 'techniczny';

interface PunktPolityki {
  klucz: string;
  nazwa: string;
  grupa: GrupaPunktu;
  /** Objaśnienie statyczne z rejestru definicji — treść stabilna, niezależna od odczytu. */
  objasnienie: string;
  /** Klucz przełącznika w odpowiedzi rdzenia (`IsolationSwitch.kind` albo `IsolationTechnicalSwitch.scope`). */
  przelacznikKlucz: IsolationContextKind | IsolationTechnicalScope;
  /** Etykiety wartości włączonej (`isolated: true`) i wyłączonej, właściwe grupie punktu. */
  etykietaWlaczony: string;
  etykietaWylaczony: string;
}

/**
 * Jedenaście punktów izolacji z etykietami i objaśnieniami odpowiadającymi
 * definicjom rdzenia (`server/internal/konfig/definicje_izolacji.go`:
 * `definicjeIzolacjiKontekstu`, `definicjeIzolacjiTechnicznej`). Objaśnienie
 * jest opisem klucza — co znaczy i jaki jest stan wyjściowy platformy — a nie
 * wynikiem odczytu, dlatego stoi w pliku: odpowiedź `policy.preview` niesie
 * wartości, nie objaśnienia.
 */
const PUNKTY_POLITYKI: readonly PunktPolityki[] = [
  {
    klucz: 'izolacja_historia',
    nazwa: 'Historia wymiany',
    grupa: 'kontekst',
    objasnienie:
      'Zapis wymiany wiadomości. Odrębna: nowe okno zaczyna z pustą historią. Współdzielona: ten sam ' +
      'zapis widoczny w kilku oknach lub modułach. Stan wyjściowy platformy: odrębna.',
    przelacznikKlucz: IsolationContextKind.History,
    etykietaWlaczony: 'Odrębna',
    etykietaWylaczony: 'Współdzielona',
  },
  {
    klucz: 'izolacja_pamiec',
    nazwa: 'Pamięć długoterminowa',
    grupa: 'kontekst',
    objasnienie:
      'Pamięć długoterminowa zasięgu. Odrębna: pamięć jednego zasięgu niewidoczna w innym. ' +
      'Współdzielona: jeden zasób zasila kilka zasięgów. Stan wyjściowy platformy: odrębna.',
    przelacznikKlucz: IsolationContextKind.Memory,
    etykietaWlaczony: 'Odrębna',
    etykietaWylaczony: 'Współdzielona',
  },
  {
    klucz: 'izolacja_kontekst',
    nazwa: 'Bieżący stan roboczy',
    grupa: 'kontekst',
    objasnienie:
      'Aktywne pliki, projekt, załączniki, zmienne. Współdzielony przenosi się między oknami bez ' +
      'przeładowania. Stan wyjściowy platformy: odrębny.',
    przelacznikKlucz: IsolationContextKind.Context,
    etykietaWlaczony: 'Odrębny',
    etykietaWylaczony: 'Współdzielony',
  },
  {
    klucz: 'izolacja_katalog_roboczy_sesji',
    nazwa: 'Katalog roboczy sesji',
    grupa: 'techniczny',
    objasnienie:
      'Fizyczny katalog plików procesu. Włączony: własny katalog, niewidoczny dla innych sesji ' +
      'zasięgu. Stan wyjściowy platformy: wyłączony.',
    przelacznikKlucz: IsolationTechnicalScope.WorkingDirectory,
    etykietaWlaczony: 'Włączony',
    etykietaWylaczony: 'Wyłączony',
  },
  {
    klucz: 'izolacja_srodowisko_procesu',
    nazwa: 'Środowisko procesu',
    grupa: 'techniczny',
    objasnienie:
      'Zmienne środowiskowe i kontekst uruchomieniowy. Włączony: własny zestaw zmiennych zamiast ' +
      'wspólnego środowiska serwera. Stan wyjściowy platformy: wyłączony.',
    przelacznikKlucz: IsolationTechnicalScope.ProcessEnvironment,
    etykietaWlaczony: 'Włączony',
    etykietaWylaczony: 'Wyłączony',
  },
  {
    klucz: 'izolacja_katalog_danych_modelu',
    nazwa: 'Katalog danych modelu',
    grupa: 'techniczny',
    objasnienie:
      'Dane pomocnicze kanału modelu: konfiguracja, dane tymczasowe, ustawienia dostawcy. Włączony: ' +
      'własny katalog danych modelu. Stan wyjściowy platformy: wyłączony.',
    przelacznikKlucz: IsolationTechnicalScope.ModelDataDirectory,
    etykietaWlaczony: 'Włączony',
    etykietaWylaczony: 'Wyłączony',
  },
  {
    klucz: 'izolacja_dostep_sieciowy',
    nazwa: 'Dostęp sieciowy',
    grupa: 'techniczny',
    objasnienie:
      'Połączenia wychodzące: API, strony, serwery MCP, rozszerzenia. Włączony: odrębny, ograniczony ' +
      'dostęp sieciowy. Stan wyjściowy platformy: wyłączony.',
    przelacznikKlucz: IsolationTechnicalScope.NetworkAccess,
    etykietaWlaczony: 'Włączony',
    etykietaWylaczony: 'Wyłączony',
  },
  {
    klucz: 'izolacja_odczyt_zapis_plikow',
    nazwa: 'Odczyt/zapis plików',
    grupa: 'techniczny',
    objasnienie:
      'Uprawnienia do plików poza własnym katalogiem roboczym. Włączony: dostęp wyłącznie do ścieżek ' +
      'jawnie dozwolonych. Stan wyjściowy platformy: wyłączony.',
    przelacznikKlucz: IsolationTechnicalScope.FileAccess,
    etykietaWlaczony: 'Włączony',
    etykietaWylaczony: 'Wyłączony',
  },
  {
    klucz: 'izolacja_konto_i_token',
    nazwa: 'Konto i token',
    grupa: 'techniczny',
    objasnienie:
      'Token dostępu, klucz API, dane logowania kanału modelu. Włączony: własne dane dostępowe ' +
      'zamiast platformowych. Stan wyjściowy platformy: wyłączony.',
    przelacznikKlucz: IsolationTechnicalScope.AccountToken,
    etykietaWlaczony: 'Włączony',
    etykietaWylaczony: 'Wyłączony',
  },
  {
    klucz: 'izolacja_model_procesu',
    nazwa: 'Model procesu',
    grupa: 'techniczny',
    objasnienie:
      'Instancja procesu wykonawczego modelu. Włączony: własna, niezależna instancja zamiast wspólnej ' +
      'puli. Stan wyjściowy platformy: wyłączony.',
    przelacznikKlucz: IsolationTechnicalScope.ProcessModel,
    etykietaWlaczony: 'Włączony',
    etykietaWylaczony: 'Wyłączony',
  },
  {
    klucz: 'izolacja_serwer_wykonania',
    nazwa: 'Serwer wykonania',
    grupa: 'techniczny',
    objasnienie:
      'Serwer wykonania procesu — istotne przy kanale zdalnym. Włączony: serwer dedykowany zamiast ' +
      'współdzielonego. Stan wyjściowy platformy: wyłączony.',
    przelacznikKlucz: IsolationTechnicalScope.ExecutionServer,
    etykietaWlaczony: 'Włączony',
    etykietaWylaczony: 'Wyłączony',
  },
];

const NAGLOWKI_GRUP: Readonly<Record<GrupaPunktu, string>> = {
  kontekst: 'Kontekst — trzy klucze (odrębna / współdzielona)',
  techniczny: 'Zakres techniczny — osiem kluczy (włączony / wyłączony)',
};

const ETYKIETY_SCOPE: Readonly<Record<ConfigScope, string>> = {
  // Słownik jest zupełny wobec kontraktu: brak etykiety dla nowego zasięgu
  // zatrzymuje kompilację.
  [ConfigScope.Application]: 'aplikacja',
  [ConfigScope.Global]: 'Globalnie',
  [ConfigScope.Environment]: 'Środowisko',
  [ConfigScope.Module]: 'Moduł',
  [ConfigScope.ModulePair]: 'Para modułów',
  [ConfigScope.Project]: 'Projekt',
  [ConfigScope.Session]: 'Karta sesji',
  [ConfigScope.Role]: 'Rola',
  [ConfigScope.Window]: 'Okno komunikacji',
};

export function utworzObszar(zaleznosci: ZaleznosciObszaru): ObszarIzolacji {
  const tresc = utworzStanTresci('pi');
  const punktWidzenia = utworzPunktWidzenia(
    zaleznosci.kanal,
    zaleznosci.warstwa,
    zaleznosci.zasieg,
    () => renderuj(),
  );

  async function odczytaj(): Promise<void> {
    const miejsce = tresc.tresc();
    miejsce.replaceChildren();

    const wynik = await pobierzPolityke(zaleznosci, punktWidzenia.zadanie());

    if (!wynik.udana) {
      miejsce.append(
        punktWidzenia.element,
        zbudujWstep(null, wynik.powod),
        zbudujTabele(null),
        // Granica asystenta stoi także przy nieudanym odczycie: wynika
        // z kontraktu, nie z wyniku tego wywołania.
        zbudujGraniceAsystenta(),
      );
      return;
    }

    const polityka = wynik.polityka ?? null;
    miejsce.append(
      punktWidzenia.element,
      zdanieZakresu(punktWidzenia.opis()),
      zbudujWstep(polityka, undefined),
      zbudujTabele(polityka),
      zbudujGraniceAsystenta(),
    );
  }

  function renderuj(): void {
    void odczytaj();
  }

  renderuj();

  return {
    element: tresc.element,
    odswiez: renderuj,
  };
}

interface OdczytPolityki {
  udana: boolean;
  polityka?: IsolationPolicy;
  powod?: string;
}

/**
 * Woła `isolation.policy.preview` z punktem widzenia wskazanym przez Operatora.
 * Żądania nie wolno wysłać pustego: rdzeń sprowadza je wtedy do poziomu
 * globalnego i okno pokazałoby politykę całej platformy pod nazwą „efektywna".
 */
function pobierzPolityke(
  zaleznosci: ZaleznosciObszaru,
  zadanie: IsolationPolicyPreviewRequest,
): Promise<OdczytPolityki> {
  return new Promise((rozstrzygnij) => {
    zaleznosci.kanal.wyslij(Command.IsolationPolicyPreview, zadanie, (wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        rozstrzygnij({ udana: false, powod: opisOdmowyBledu('Odczyt polityki efektywnej', wynik.blad) });
        return;
      }
      rozstrzygnij({ udana: true, polityka: wynik.wynik.policy });
    });
  });
}

/** Wstęp nad tabelą: co pokazuje, jak czytać plakietki i co znaczy nieudany odczyt. */
function zbudujWstep(polityka: IsolationPolicy | null, powodBraku: string | undefined): HTMLElement {
  const sekcja = document.createElement('div');

  const opis = document.createElement('p');
  opis.textContent =
    'Rozstrzygnięcie jedenastu punktów izolacji dla bieżącego kontekstu: wartość obowiązująca i jej ' +
    'pochodzenie — czy ktoś ją zapisał (i na jakim poziomie), czy obowiązuje wartość wbudowana, bo ' +
    'nikt niczego nie zapisał.';

  const ostrzezenie = document.createElement('p');
  const plakietkaOstrzezenia = document.createElement('strong');
  plakietkaOstrzezenia.className = 'dn-plakietka dn-plakietka--blad';
  plakietkaOstrzezenia.textContent = 'Bez potwierdzenia';
  ostrzezenie.append(
    plakietkaOstrzezenia,
    document.createTextNode(
      ' — stan wyjściowy platformy to pełna swoboda operacyjna; egzekutor odrzuca tylko naruszenie ' +
        'wartości WŁĄCZONEJ. Ta tabela jest jedynym miejscem sprawdzenia, czy ochrona naprawdę działa. ' +
        'Gdy odczyt poniżej zawodzi, żaden punkt nie jest potwierdzony jako chroniony — nawet jeśli ' +
        'tak było wcześniej.',
    ),
  );

  const legenda = document.createElement('p');
  legenda.textContent =
    '„Zapis Operatora" (zielono) — świadomy zapis na wskazanym poziomie. „Wartość domyślna" ' +
    '(niebiesko) — nikt nie zapisał, obowiązuje rejestr definicji. „Nieznane rejestrowi" (żółto) — ' +
    'klucz zapisany, lecz spoza jedenastki. „Brak odczytu" (czerwono) — rdzeń nie odpowiedział; nie ' +
    'wolno mylić tego z „Wartość domyślna", bo prowadzą do przeciwnych wniosków o bezpieczeństwie.';

  sekcja.append(opis, ostrzezenie, legenda);

  if (!polityka) {
    // Odmowa mówi trzy rzeczy: co się nie udało, dlaczego to nie jest wartość
    // domyślna i czym Operator to zmieni. `isolation.policy.preview` stoi
    // w kontrakcie i została wysłana, więc brak wyniku jest odmową odczytu,
    // nie brakiem komendy.
    const plakietka = document.createElement('span');
    plakietka.className = 'dn-plakietka dn-plakietka--blad';
    plakietka.textContent = 'Odczyt nieudany';

    const zdanie = document.createElement('p');
    zdanie.textContent =
      `${powodBraku ?? POWOD_BRAKU_ODCZYTU} Popraw punkt widzenia powyżej — poziom zasięgu, ` +
      'identyfikator bytu albo okno komunikacji — i naciśnij „Odczytaj podgląd dla tego punktu ' +
      'widzenia” ponownie; warstwę zmienisz na pasie narzędzi okna.';

    const odmowa = document.createElement('div');
    odmowa.append(plakietka, zdanie);
    sekcja.append(odmowa);
  }

  return sekcja;
}

/** Zdanie nad tabelą: czego dotyczy bieżący odczyt. */
function zdanieZakresu(opis: string): HTMLElement {
  const element = document.createElement('p');
  element.className = 'pi-widzenie__zakres';
  element.textContent = opis;
  return element;
}

function zbudujTabele(polityka: IsolationPolicy | null): HTMLElement {
  const tabela = document.createElement('table');
  tabela.className = 'dn-tabela';
  tabela.setAttribute(
    'aria-label',
    'Jedenaście punktów izolacji: wartość obowiązująca i pochodzenie dla każdego klucza',
  );

  const thead = document.createElement('thead');
  const wierszNaglowka = document.createElement('tr');
  for (const etykieta of ['Punkt izolacji', 'Wartość obowiązująca', 'Pochodzenie', 'Poziom']) {
    const th = document.createElement('th');
    th.scope = 'col';
    th.textContent = etykieta;
    wierszNaglowka.append(th);
  }
  thead.append(wierszNaglowka);

  const tbody = document.createElement('tbody');
  let grupaPoprzednia: GrupaPunktu | null = null;
  for (const punkt of PUNKTY_POLITYKI) {
    if (punkt.grupa !== grupaPoprzednia) {
      tbody.append(zbudujWierszGrupy(NAGLOWKI_GRUP[punkt.grupa]));
      grupaPoprzednia = punkt.grupa;
    }
    tbody.append(zbudujWierszPunktu(punkt, polityka));
  }

  tabela.append(thead, tbody);
  return tabela;
}

function zbudujWierszGrupy(naglowek: string): HTMLElement {
  const wiersz = document.createElement('tr');
  const td = document.createElement('td');
  td.colSpan = 4;
  td.className = 'dn-dane';
  const etykieta = document.createElement('strong');
  etykieta.textContent = naglowek;
  td.append(etykieta);
  wiersz.append(td);
  return wiersz;
}

/** Odnajduje przełącznik punktu w odpowiedzi rdzenia — `isolated: boolean` albo `undefined`, gdy klucz nie wrócił (traktowane jak „nieznane"). */
function znajdzIsolated(punkt: PunktPolityki, polityka: IsolationPolicy): boolean | undefined {
  if (punkt.grupa === 'kontekst') {
    return polityka.contextSwitches.find((p) => p.kind === punkt.przelacznikKlucz)?.isolated;
  }
  return polityka.technicalSwitches.find((p) => p.scope === punkt.przelacznikKlucz)?.isolated;
}

/**
 * Jeden wiersz tabeli punktów izolacji. Bez `polityka` — gdy rdzeń nie
 * odpowiedział — każda komórka niesie stan `brak-odczytu`, odróżniony od
 * `domyslna`, bo oba prowadzą do przeciwnych wniosków o bezpieczeństwie.
 * Z `polityka` wartość pochodzi z `isolated`, a pochodzenie z `policy.origin`:
 * obecny `origin` znaczy zapis na wskazanym poziomie, brak `origin` znaczy
 * wartość z rejestru definicji.
 */
function zbudujWierszPunktu(punkt: PunktPolityki, polityka: IsolationPolicy | null): HTMLElement {
  const wiersz = document.createElement('tr');

  const tdPunkt = document.createElement('td');
  const nazwa = document.createElement('strong');
  nazwa.textContent = punkt.nazwa;
  const dymek = utworzDymekObjasnienia(punkt.objasnienie);
  const kod = document.createElement('code');
  kod.textContent = punkt.klucz;
  const wierszKodu = document.createElement('div');
  wierszKodu.className = 'dn-dane';
  wierszKodu.append(kod);
  tdPunkt.append(nazwa, dymek, wierszKodu);

  const tdWartosc = document.createElement('td');
  const tdPochodzenie = document.createElement('td');
  const tdPoziom = document.createElement('td');
  tdPoziom.className = 'dn-dane';

  if (!polityka) {
    tdWartosc.append(zbudujPlakietke('brak-odczytu', 'Wartość nieznana — ' + POWOD_BRAKU_ODCZYTU));
    tdPochodzenie.append(zbudujPlakietke('brak-odczytu', 'Pochodzenie nieznane — ' + POWOD_BRAKU_ODCZYTU));
    tdPoziom.textContent = '—';
    tdPoziom.title = 'Poziom nieznany — ' + POWOD_BRAKU_ODCZYTU;
  } else {
    const isolated = znajdzIsolated(punkt, polityka);
    if (isolated === undefined) {
      const powod =
        'Rdzeń nie zwrócił tego przełącznika w odpowiedzi — klucz nieznany rejestrowi tego odczytu.';
      tdWartosc.append(zbudujPlakietke('nieznane', powod));
      tdPochodzenie.append(zbudujPlakietke('nieznane', powod));
      tdPoziom.textContent = '—';
      tdPoziom.title = powod;
    } else {
      const etykietaWartosci = isolated ? punkt.etykietaWlaczony : punkt.etykietaWylaczony;
      const stanWartosci = isolated ? 'dn-plakietka--sukces' : 'dn-plakietka--informacja';
      const spanWartosc = document.createElement('span');
      spanWartosc.className = `dn-plakietka ${stanWartosci}`;
      spanWartosc.textContent = etykietaWartosci;
      tdWartosc.append(spanWartosc);

      if (polityka.origin) {
        const poziomEtykieta = ETYKIETY_SCOPE[polityka.origin] ?? polityka.origin;
        const opis = `Zapisano na poziomie: ${poziomEtykieta}${polityka.scopeId ? ` (${polityka.scopeId})` : ''}.`;
        tdPochodzenie.append(zbudujPlakietke('zapis', opis));
        tdPoziom.textContent = poziomEtykieta;
        tdPoziom.title = opis;
      } else {
        const opis = 'Nikt nie zapisał na żadnym z ośmiu poziomów — obowiązuje wartość z rejestru definicji.';
        tdPochodzenie.append(zbudujPlakietke('domyslna', opis));
        tdPoziom.textContent = '—';
        tdPoziom.title = opis;
      }
    }
  }

  wiersz.append(tdPunkt, tdWartosc, tdPochodzenie, tdPoziom);
  return wiersz;
}

/**
 * Granica asystenta wobec izolacji: czy asystent, który obsługuje aplikację
 * za Operatora, może przestawić sobie te punkty.
 *
 * Odpowiedź powstaje z wykazu narzędzi modelu (`NARZEDZIA_MODELU`) przy każdym
 * odczycie, a nie z wpisu na stałe: dołożenie do narzędzi komendy zapisu
 * izolacji zmienia treść wiersza samo.
 *
 * Miejscem styku jest `config.session.set` z roli klawiatury: pisze przedmioty
 * tych samych punktów — katalog roboczy, katalogi dodatkowe, środowisko,
 * uruchomienie, narzędzia, uprawnienia i konto. Przełączników izolacji nie
 * rusza i egzekutor zostaje na miejscu, więc granicy nie przekracza.
 */
const ZAPIS_IZOLACJI: readonly Command[] = [
  Command.IsolationContextSet,
  Command.IsolationTechnicalSet,
  Command.IsolationProfileSave,
  Command.IsolationProfileAssign,
  Command.IsolationProfileDelete,
  Command.IsolationLayerSet,
];

const ODCZYT_IZOLACJI: readonly Command[] = [
  Command.IsolationScopeList,
  Command.IsolationContextGet,
  Command.IsolationTechnicalGet,
  Command.IsolationProfileList,
  Command.IsolationProfileLoad,
  Command.IsolationPolicyPreview,
];

/** Które z podanych komend stoją w wykazie narzędzi modelu. */
function narzedziamiModelu(komendy: readonly Command[]): readonly Command[] {
  const wykaz = new Set(NARZEDZIA_MODELU.map((narzedzie) => String(narzedzie.command)));
  return komendy.filter((komenda) => wykaz.has(String(komenda)));
}

function zbudujGraniceAsystenta(): HTMLElement {
  const sekcja = document.createElement('section');
  sekcja.className = 'pi-granica';

  const naglowek = document.createElement('h3');
  naglowek.textContent = 'Czy asystent może przestawić te punkty?';

  const zapisemNarzedzia = narzedziamiModelu(ZAPIS_IZOLACJI);
  const odczytemNarzedzia = narzedziamiModelu(ODCZYT_IZOLACJI);

  const werdykt = document.createElement('p');
  const plakietka = document.createElement('strong');
  plakietka.className = `dn-plakietka ${
    zapisemNarzedzia.length === 0 ? 'dn-plakietka--sukces' : 'dn-plakietka--blad'
  }`;
  plakietka.textContent = zapisemNarzedzia.length === 0 ? 'Nie ma czym' : 'MOŻE — sprawdź to';
  werdykt.append(plakietka, document.createTextNode(' '));

  if (zapisemNarzedzia.length === 0) {
    werdykt.append(
      document.createTextNode(
        `— z dwunastu komend izolacji ${odczytemNarzedzia.length} ODCZYTOWE są narzędziami modelu ` +
          `(${odczytemNarzedzia.join(', ')}), a ŻADNA z ${ZAPIS_IZOLACJI.length} zapisowych nie jest ` +
          `(${ZAPIS_IZOLACJI.join(', ')}). Asystent widzi swoje granice i nie ma czym ich przestawić. ` +
          'Zmierzone na wykazie narzędzi kontraktu przy tym odczycie, nie przepisane z dokumentu.',
      ),
    );
  } else {
    // Gałąź ostrzegawcza: wykaz narzędzi niesie komendę zapisu izolacji.
    werdykt.append(
      document.createTextNode(
        `— wykaz narzędzi modelu niesie dziś komendy ZAPISU izolacji: ${zapisemNarzedzia.join(', ')}. ` +
          'Asystent może przestawić punkt izolacji sam. Jeżeli to nie jest zamierzone, ' +
          'poprawka należy do wykazu narzędzi (server/internal/narzedzia), nie do tego okna.',
      ),
    );
  }

  const furtka = document.createElement('p');
  furtka.className = 'dn-pole-opis';
  furtka.textContent =
    'Jedno miejsce styku, którego nie przemilczamy: rola klawiatury daje asystentowi komendę ' +
    'config.session.set ponad wykaz narzędzi. Pisze nią obszary workingDirectory, ' +
    'additionalDirectories, environment, runtime, tools, permissions i account na OKNIE DOCELOWYM — ' +
    'czyli przedmioty tych samych punktów izolacji. Przełącznika izolacji nie rusza i egzekutor ' +
    'zostaje na miejscu, więc granicy nie przekracza; ale to jest jedyna furtka w jej pobliżu ' +
    'i to ona rodzi zdarzenie config.changed, które sprawcy nie niesie.';

  sekcja.append(naglowek, werdykt, furtka);
  return sekcja;
}

function zbudujPlakietke(stan: StanPochodzenia, opisDostepnosci: string): HTMLElement {
  const opisPlakietki = PLAKIETKI_POCHODZENIA[stan];
  const span = document.createElement('span');
  span.className = `dn-plakietka ${opisPlakietki.klasa}`;
  span.textContent = opisPlakietki.etykieta;
  span.setAttribute('aria-label', opisDostepnosci);
  span.title = opisDostepnosci;
  return span;
}
