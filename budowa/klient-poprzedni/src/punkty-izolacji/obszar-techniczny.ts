import {
  Command,
  IsolationTechnicalScope,
  type IsolationTechnicalSwitch,
  type RequestOf,
  type ResponseOf,
} from '../../../shared/contract';
import { utworzDymekObjasnienia } from '../komponenty/dymek';
import { utworzStanTresci } from '../komponenty/stan-tresci';
import { opisOdmowyBledu } from '../komponenty/odmowa';
import type { Kanal, Wynik } from '../protokol/kanal';
import type { ObszarIzolacji, ZaleznosciObszaru } from './obszary';
import { NAZWY_WARSTW } from './stan-warstwy';

/**
 * Obszar „Zakres techniczny" — osiem punktów izolacji technicznej
 * z `konfig/definicjeIzolacjiTechnicznej`
 * (`server/internal/konfig/definicje_izolacji.go`): katalog roboczy sesji,
 * środowisko procesu, katalog danych modelu, dostęp sieciowy, odczyt/zapis
 * plików, konto i token, model procesu, serwer wykonania. Wartości kluczy to
 * `wlaczony`/`wylaczony`, domyślnie zawsze `wylaczony`.
 *
 * Różnica wobec obszaru „Kontekst": tam wartość domyślna
 * (`odrebna`/`wspoldzielona`) opisuje zakres widoczności i obie wartości coś
 * znaczą. Tu `wylaczony` nie jest słabszym wariantem ochrony: egzekutor rdzenia
 * bierze wyłącznie klucze włączone i odrzuca wykonanie, które by je naruszyło.
 * Klucz wyłączony nie jest sprawdzany w ogóle — nikt nie pilnuje tego zakresu,
 * dopóki Operator świadomie go nie włączy. Stan wyjściowy platformy to pełna
 * swoboda operacyjna, a widok nazywa ten skutek przy każdym z ośmiu kluczy,
 * nie tylko raz na górze, żeby przewijanie wykazu nie zgubiło ostrzeżenia.
 *
 * `odswiez()` woła `isolation.technical.get`; przełączenie klucza woła
 * `isolation.technical.set` z pełnym zestawem ośmiu przełączników (zmieniony
 * jeden, reszta bez zmian) — kontrakt niesie cały zestaw naraz, nie pojedynczy
 * klucz.
 *
 * Wartość na ekranie zmienia się wyłącznie po potwierdzeniu z rdzenia:
 * `switches` w tym module to zawsze ostatnia odpowiedź `get`/`set`, nigdy
 * przewidywanie kliknięcia. Odmowa zapisu nie idzie do paska ogólnego obszaru —
 * ląduje w komórce „Zmiana" wiersza klucza, którego dotyczy, żeby Operator
 * wiedział, który z ośmiu się nie zapisał; wartość wyświetlona zostaje ta
 * sprzed próby, bo rdzeń jej nie potwierdził.
 */
export function utworzObszar(zaleznosci: ZaleznosciObszaru): ObszarIzolacji {
  const { kanal, warstwa, zasieg } = zaleznosci;
  const tresc = utworzStanTresci('pi');

  /** Ostatnia odpowiedź rdzenia — jedyne źródło wartości pokazywanych Operatorowi. */
  let switches: readonly IsolationTechnicalSwitch[] = [];
  /** Odmowa zapisu ostatniej próby, per klucz — czyszczona przy udanym odczycie i udanym zapisie tego klucza. */
  const bledyZapisu = new Map<string, string>();
  /** Klucz aktualnie w trakcie zapisu — blokuje ponowne kliknięcie tego samego wiersza, nie całej tabeli. */
  const wTrakcieZapisu = new Set<string>();

  async function odczytaj(): Promise<void> {
    tresc.ladowanie('Pytam rdzeń o wartości ośmiu kluczy zakresu technicznego (isolation.technical.get)…');

    const wynik = await posijKomende(kanal, Command.IsolationTechnicalGet, {
      scope: zasieg.zasieg(),
      scopeId: zasieg.bytDoZadania(),
      layer: warstwa.warstwa(),
    });

    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.blad(opisOdmowyBledu('Odczyt zakresu technicznego', wynik.blad), wynik.blad);
      return;
    }

    switches = wynik.wynik.switches;
    bledyZapisu.clear();
    renderuj();
  }

  function renderuj(): void {
    const miejsce = tresc.tresc();
    miejsce.append(wstepObszaru(zdanieZasiegu()), tabelaPunktow());
  }

  /** Jedno zdanie o tym, czego dotyczy odczyt i zapis — zasięg z lewego panelu, warstwa z pasa narzędzi. */
  function zdanieZasiegu(): string {
    return `Odczyt i zapis obejmują ${zasieg.opis()} oraz warstwę „${NAZWY_WARSTW[warstwa.warstwa()]}”.`;
  }

  async function ustawKlucz(punkt: PunktTechniczny, docelowa: boolean): Promise<void> {
    if (wTrakcieZapisu.has(punkt.klucz)) return;
    wTrakcieZapisu.add(punkt.klucz);
    bledyZapisu.delete(punkt.klucz);
    renderuj();

    const zestawDoZapisu = zbudujZestawZeZmiana(switches, punkt.scope, docelowa);

    const wynik = await posijKomende(kanal, Command.IsolationTechnicalSet, {
      scope: zasieg.zasieg(),
      scopeId: zasieg.bytDoZadania(),
      layer: warstwa.warstwa(),
      switches: zestawDoZapisu,
    });

    wTrakcieZapisu.delete(punkt.klucz);

    if (!wynik.udany || wynik.wynik === undefined) {
      // Wartość nie zmienia się na ekranie — `switches` zostaje sprzed próby.
      // Odmowa idzie do wiersza tego klucza, nie do paska ogólnego.
      bledyZapisu.set(
        punkt.klucz,
        opisOdmowyBledu(`Zapis klucza „${punkt.nazwa}” nie powiódł się`, wynik.blad),
      );
      renderuj();
      return;
    }

    switches = wynik.wynik.switches;
    bledyZapisu.delete(punkt.klucz);
    renderuj();
  }

  return {
    element: tresc.element,
    odswiez: () => void odczytaj(),
  };

  /** Wiersz jednego punktu: nazwa i klucz, dwa skutki, wartość bieżąca, przełączenie. */
  function wierszPunktu(punkt: PunktTechniczny): HTMLTableRowElement {
    const wiersz = document.createElement('tr');
    wiersz.dataset['klucz'] = punkt.klucz;
    wiersz.append(komorkaNazwy(punkt), komorkaSkutku(punkt), komorkaWartosci(punkt), komorkaZmiany(punkt));
    return wiersz;
  }

  /** Wartość obecna wprost z ostatniej odpowiedzi rdzenia — nigdy przewidywana. */
  function komorkaWartosci(punkt: PunktTechniczny): HTMLTableCellElement {
    const wpis = switches.find((s) => s.scope === punkt.scope);

    const plakietka = document.createElement('span');
    const wlaczony = wpis?.isolated === true;

    // Trzy stany, trzy wyglądy — nie dwa. „Nieznany" (rdzeń nie zwrócił tego
    // klucza) i „wyłączony" (rdzeń zwrócił: nikt nie pilnuje) prowadzą do
    // przeciwnych wniosków o bezpieczeństwie maszyny, więc nie mogą wyglądać
    // tak samo. Warianty zdefiniowane w `komponenty/plakietka.css`: sukces,
    // ostrzezenie, blad, informacja, sygnal, rola, atrament — spoza tego zbioru
    // plakietka nie dostaje żadnego wyróżnienia.
    plakietka.className = [
      'dn-plakietka',
      wpis === undefined
        ? 'dn-plakietka--informacja'
        : wlaczony
          ? 'dn-plakietka--sukces'
          : 'dn-plakietka--ostrzezenie',
    ].join(' ');
    plakietka.textContent = wpis === undefined ? 'nieznany' : wlaczony ? 'włączony' : 'wyłączony';

    const komorka = document.createElement('td');
    komorka.append(plakietka);

    if (wpis?.explanation !== undefined && wpis.explanation !== '') {
      const objasnienie = document.createElement('p');
      objasnienie.className = 'dn-pole-opis';
      objasnienie.textContent = wpis.explanation;
      komorka.append(objasnienie);
    }

    return komorka;
  }

  /** Para przycisków przełączenia oraz odmowa przy tym kluczu, gdy ostatni zapis się nie powiódł. */
  function komorkaZmiany(punkt: PunktTechniczny): HTMLTableCellElement {
    const zapisuje = wTrakcieZapisu.has(punkt.klucz);

    const przyciskWlacz = document.createElement('button');
    przyciskWlacz.type = 'button';
    przyciskWlacz.className = 'dn-btn dn-btn--zarys';
    przyciskWlacz.textContent = zapisuje ? 'Zapisuję…' : 'Ustaw: włączony';
    przyciskWlacz.addEventListener('click', () => void ustawKlucz(punkt, true));

    const przyciskWylacz = document.createElement('button');
    przyciskWylacz.type = 'button';
    przyciskWylacz.className = 'dn-btn dn-btn--zarys';
    przyciskWylacz.textContent = zapisuje ? 'Zapisuję…' : 'Ustaw: wyłączony';
    przyciskWylacz.addEventListener('click', () => void ustawKlucz(punkt, false));

    // Żaden przycisk nie niesie `disabled` — nawet w trakcie zapisu klik jest
    // przyjmowany (kolejne wywołanie po prostu nadpisze poprzednie w locie
    // rdzenia); etykieta „Zapisuję…" informuje, nie blokuje.

    const grupa = document.createElement('div');
    grupa.className = 'dn-przybornik';
    grupa.append(przyciskWlacz, przyciskWylacz);

    const komorka = document.createElement('td');
    komorka.append(grupa);

    const powodOdmowy = bledyZapisu.get(punkt.klucz);
    if (powodOdmowy !== undefined) {
      const odmowa = document.createElement('p');
      odmowa.className = 'dn-pole-opis';
      odmowa.setAttribute('role', 'alert');
      odmowa.textContent = powodOdmowy;
      komorka.append(odmowa);
    }

    return komorka;
  }

  /** Tabela ośmiu punktów — biblioteczny `dn-tabela`, bez własnej klasy CSS. */
  function tabelaPunktow(): HTMLElement {
    const podpis = document.createElement('caption');
    podpis.className = 'dn-sr-only';
    podpis.textContent =
      'Osiem punktów izolacji zakresu technicznego: nazwa, skutek włączenia i wyłączenia, wartość obecna, sposób zmiany.';

    const glowa = document.createElement('thead');
    const wierszGlowy = document.createElement('tr');
    wierszGlowy.append(
      naglowekKolumny('Punkt izolacji'),
      naglowekKolumny('Skutek'),
      naglowekKolumny('Wartość'),
      naglowekKolumny('Zmiana'),
    );
    glowa.append(wierszGlowy);

    const cialo = document.createElement('tbody');
    cialo.append(...PUNKTY_TECHNICZNE.map(wierszPunktu));

    const tabela = document.createElement('table');
    tabela.className = 'dn-tabela';
    tabela.setAttribute('aria-label', 'Osiem punktów izolacji zakresu technicznego');
    tabela.append(podpis, glowa, cialo);
    return tabela;
  }
}

/** Opakowuje `kanal.wyslij` w Promise — kanał sam daje wyłącznie wersję z wywołaniem zwrotnym. */
function posijKomende<K extends Command>(
  kanal: Kanal,
  komenda: K,
  zadanie: RequestOf<K>,
): Promise<Wynik<ResponseOf<K>>> {
  return new Promise((rozstrzygnij) => {
    kanal.wyslij(komenda, zadanie, (wynik) => rozstrzygnij(wynik));
  });
}

/** Zestaw ośmiu przełączników z jedną zmianą naniesioną — kontrakt `isolation.technical.set` niesie cały zestaw naraz. */
function zbudujZestawZeZmiana(
  biezace: readonly IsolationTechnicalSwitch[],
  scope: IsolationTechnicalScope,
  isolated: boolean,
): IsolationTechnicalSwitch[] {
  let trafiono = false;
  const zestaw = biezace.map((s) => {
    if (s.scope !== scope) return s;
    trafiono = true;
    return { ...s, isolated };
  });
  if (!trafiono) zestaw.push({ scope, isolated });
  return zestaw;
}

/** Jeden z ośmiu punktów izolacji zakresu technicznego, wedle `definicje_izolacji.go`. */
interface PunktTechniczny {
  /** Klucz kolumny `regula_izolacji_technicznej.zakres` (`KluczIzolacja…`), do etykiet i identyfikacji wiersza. */
  klucz: string;
  /** Wartość kontraktu (`IsolationTechnicalScope`) niesiona w `isolation.technical.get`/`.set`. */
  scope: IsolationTechnicalScope;
  /** Nazwa czytelna dla Operatora. */
  nazwa: string;
  /** Co dokładnie zyskuje zasięg, gdy klucz jest `wlaczony`. */
  skutekWlaczony: string;
  /** Co dokładnie nie jest pilnowane, gdy klucz jest `wylaczony` — nazwane wprost. */
  skutekWylaczony: string;
}

/**
 * Osiem punktów w kolejności `definicjeIzolacjiTechnicznej`.
 * Treść `skutekWlaczony` parafrazuje `Objasnienie` z `definicje_izolacji.go`; `skutekWylaczony`
 * dopisuje to, czego `Objasnienie` nie mówi wprost — że brak włączenia oznacza
 * brak jakiejkolwiek kontroli tego zakresu, nie kontrolę słabszą.
 */
const PUNKTY_TECHNICZNE: readonly PunktTechniczny[] = [
  {
    klucz: 'izolacja_katalog_roboczy_sesji',
    scope: IsolationTechnicalScope.WorkingDirectory,
    nazwa: 'Katalog roboczy sesji',
    skutekWlaczony:
      'Sesja dostaje własny, fizyczny katalog plików procesu, niewidoczny dla innych sesji tego zasięgu.',
    skutekWylaczony:
      'Sesje tego zasięgu współdzielą jeden katalog roboczy — inna sesja może odczytać albo nadpisać jej pliki robocze, i nikt tego nie pilnuje.',
  },
  {
    klucz: 'izolacja_srodowisko_procesu',
    scope: IsolationTechnicalScope.ProcessEnvironment,
    nazwa: 'Środowisko procesu',
    skutekWlaczony:
      'Proces dostaje własny zestaw zmiennych środowiskowych zamiast wspólnego środowiska serwera.',
    skutekWylaczony:
      'Proces dziedziczy wspólne środowisko serwera — zmienne jednej sesji, także poufne, są widoczne innym sesjom tego zasięgu, i nikt tego nie pilnuje.',
  },
  {
    klucz: 'izolacja_katalog_danych_modelu',
    scope: IsolationTechnicalScope.ModelDataDirectory,
    nazwa: 'Katalog danych modelu',
    skutekWlaczony:
      'Kanał modelu dostaje własny katalog danych pomocniczych: konfigurację, dane tymczasowe, ustawienia dostawcy.',
    skutekWylaczony:
      'Sesje tego zasięgu współdzielą jeden katalog danych modelu — dane tymczasowe i ustawienia dostawcy jednej sesji są dostępne innym, i nikt tego nie pilnuje.',
  },
  {
    klucz: 'izolacja_dostep_sieciowy',
    scope: IsolationTechnicalScope.NetworkAccess,
    nazwa: 'Dostęp sieciowy',
    skutekWlaczony:
      'Sesja dostaje odrębny, ograniczony dostęp sieciowy do API, stron, serwerów MCP i rozszerzeń.',
    skutekWylaczony:
      'Sesja może otworzyć dowolne połączenie wychodzące, bez ograniczenia właściwego temu zasięgowi, i nikt tego nie pilnuje.',
  },
  {
    klucz: 'izolacja_odczyt_zapis_plikow',
    scope: IsolationTechnicalScope.FileAccess,
    nazwa: 'Odczyt/zapis plików',
    skutekWlaczony:
      'Poza własnym katalogiem roboczym sesja dostaje dostęp wyłącznie do ścieżek, które Operator jawnie dopuścił.',
    skutekWylaczony:
      'Sesja może czytać i zapisywać pliki poza swoim katalogiem roboczym bez ograniczenia do listy dozwolonych ścieżek, i nikt tego nie pilnuje.',
  },
  {
    klucz: 'izolacja_konto_i_token',
    scope: IsolationTechnicalScope.AccountToken,
    nazwa: 'Konto i token',
    skutekWlaczony:
      'Kanał modelu używa własnego tokenu dostępu, klucza API i danych logowania zamiast platformowych; platforma trzyma odwołania do nich, nie treść sekretu.',
    skutekWylaczony:
      'Kanał modelu korzysta z danych dostępowych platformy, współdzielonych z innymi sesjami tego zasięgu, i nikt nie pilnuje, kto z nich korzysta.',
  },
  {
    klucz: 'izolacja_model_procesu',
    scope: IsolationTechnicalScope.ProcessModel,
    nazwa: 'Model procesu',
    skutekWlaczony: 'Sesja dostaje własną, niezależną instancję procesu wykonawczego modelu.',
    skutekWylaczony:
      'Sesja korzysta ze wspólnej puli instancji procesu razem z innymi sesjami tego zasięgu, i nikt tego nie pilnuje.',
  },
  {
    klucz: 'izolacja_serwer_wykonania',
    scope: IsolationTechnicalScope.ExecutionServer,
    nazwa: 'Serwer wykonania',
    skutekWlaczony:
      'Sesja dostaje dedykowany serwer wykonania procesu — istotne przy kanale zdalnym.',
    skutekWylaczony:
      'Sesja wykonuje się na serwerze współdzielonym z innymi sesjami tego zasięgu, i nikt tego nie pilnuje.',
  },
];

/** Zdanie otwierające obszar — ramuje ośmiokrotnie powtórzony fakt, zanim Operator dojdzie do wykazu. */
function wstepObszaru(zdanieZasiegu: string): HTMLElement {
  const element = document.createElement('p');
  element.className = 'dn-pole-opis';
  element.textContent =
    `${zdanieZasiegu} ` +
    'Stan wyjściowy platformy dla wszystkich ośmiu punktów poniżej to pełna swoboda operacyjna: ' +
    'każdy zaczyna wyłączony. Egzekutor rdzenia sprawdza wyłącznie klucze włączone i odrzuca ' +
    'wykonanie, które by je naruszyło — klucz wyłączony nie jest słabszą ochroną, jest jej brakiem. ' +
    'Wartość każdego wiersza pochodzi z ostatniej odpowiedzi rdzenia (isolation.technical.get/.set) — ' +
    'nigdy z kliknięcia samego w sobie; nieudany zapis zostawia dawną wartość i nazywa odmowę przy ' +
    'kluczu, którego dotyczy.';
  return element;
}

function naglowekKolumny(tekst: string): HTMLTableCellElement {
  const th = document.createElement('th');
  th.scope = 'col';
  th.textContent = tekst;
  return th;
}

function komorkaNazwy(punkt: PunktTechniczny): HTMLTableCellElement {
  const nazwa = document.createElement('strong');
  nazwa.textContent = punkt.nazwa;

  const klucz = document.createElement('code');
  klucz.textContent = punkt.klucz;

  const komorka = document.createElement('td');
  komorka.append(
    nazwa,
    utworzDymekObjasnienia(objasnienieKontekstowe(punkt)),
    document.createElement('br'),
    klucz,
  );
  return komorka;
}

/**
 * Objaśnienie kontekstowe [?] punktu — co ustawienie robi i jaki ma wpływ na
 * działanie aplikacji (rozdz. 3.3 Modelu konfiguracji). Oba skutki stoją
 * w jednym zdaniu, bo różnica między „włączony" a „wyłączony" jest tu różnicą
 * między pilnowaniem zakresu a jego brakiem, nie między ochroną mocniejszą
 * a słabszą.
 */
function objasnienieKontekstowe(punkt: PunktTechniczny): string {
  return (
    `Włączony: ${punkt.skutekWlaczony} Wyłączony: ${punkt.skutekWylaczony} ` +
    'Zapis obejmuje zasięg wskazany w lewym panelu i warstwę z pasa narzędzi okna. ' +
    'Stan wyjściowy platformy to wyłączony — okno udostępnia izolację jako możliwość, nigdy jej ' +
    'nie wymusza.'
  );
}

function komorkaSkutku(punkt: PunktTechniczny): HTMLTableCellElement {
  const wlaczony = document.createElement('p');
  wlaczony.className = 'dn-pole-opis';
  wlaczony.append(silny('Włączony: '), document.createTextNode(punkt.skutekWlaczony));

  const wylaczony = document.createElement('p');
  wylaczony.className = 'dn-pole-opis';
  wylaczony.append(silny('Wyłączony: '), document.createTextNode(punkt.skutekWylaczony));

  const komorka = document.createElement('td');
  komorka.append(wlaczony, wylaczony);
  return komorka;
}

function silny(tekst: string): HTMLElement {
  const element = document.createElement('strong');
  element.textContent = tekst;
  return element;
}
