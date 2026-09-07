// Pamięć w oknie Agents: konteksty pamięci, wpisy, wyciszenia i zasady
// przechowania — to, co model pamięta między rozmowami.
import { Command, ConfigScope, MemoryLevel } from '../../../shared/contract.ts';
import type { MemoryContext } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { dolozCzynnosciPanelu, zapytajWSzufladzie } from './czynnosci-okna.ts';
import { oglos } from './ogloszenie.ts';

const NAGLOWEK = 'Pamięć';

const POZIOMY: ReadonlyArray<readonly [string, string]> = [
  [ConfigScope.Global, 'Cała aplikacja'],
  [ConfigScope.Project, 'Projekt'],
  [ConfigScope.Session, 'Sesja'],
  [ConfigScope.Environment, 'Środowisko'],
];

interface Otoczenie {
  kanal: Kanal;
  korzen: Element;
  idSesji: () => string;
}

export function zwiazPamiecAgentow(
  kanal: Kanal,
  korzen: Element,
  idSesji: () => string,
  dolacz: (zdejmij: () => void) => void,
  przy: AddEventListenerOptions,
): void {
  const otoczenie: Otoczenie = { kanal, korzen, idSesji };
  const zdejmij = dolozCzynnosciPanelu(korzen, 'panel-model', 'Czynności pamięci', [
    {
      naglowek: 'Konteksty pamięci',
      pozycje: [
        { kod: 'konteksty', nazwa: 'Wykaz kontekstów' },
        { kod: 'kontekst-zapisz', nazwa: 'Zapisz kontekst…' },
        { kod: 'kontekst-wlacz', nazwa: 'Włącz kontekst w sesji…' },
        { kod: 'kontekst-usun', nazwa: 'Usuń kontekst…' },
      ],
    },
    {
      naglowek: 'Wpisy pamięci',
      pozycje: [
        { kod: 'wpis-zapisz', nazwa: 'Zapisz wpis…' },
        { kod: 'wpis-odepnij', nazwa: 'Odepnij wpis…' },
        { kod: 'przelacz', nazwa: 'Przestaw pamięć sesji…' },
      ],
    },
    {
      naglowek: 'Wyciszenia i przechowanie',
      pozycje: [
        { kod: 'wyciszenia', nazwa: 'Wykaz wyciszeń' },
        { kod: 'wycisz', nazwa: 'Wycisz poziom pamięci…' },
        { kod: 'przechowanie', nazwa: 'Zasady przechowania' },
        { kod: 'przechowanie-ustaw', nazwa: 'Ustaw przechowanie…' },
      ],
    },
  ], (kod) => {
    void wykonaj(otoczenie, kod, 'panel-model');
  }, przy);
  if (zdejmij !== null) dolacz(zdejmij);
}

function wypelnij(korzen: Element, wiersze: string[]): void {
  const cialo = korzen.querySelector('#panel-model .sta-okno-tresc');
  if (cialo === null) return;
  cialo.replaceChildren(...wiersze.map((tresc) => {
    const wiersz = cialo.ownerDocument.createElement('div');
    wiersz.className = 'dn-wykaz-modulu-poz';
    wiersz.textContent = tresc;
    return wiersz;
  }));
}

async function wykonaj(otoczenie: Otoczenie, kod: string, panel: string): Promise<void> {
  if (kod === 'konteksty') return wykazKontekstow(otoczenie);
  if (kod === 'kontekst-zapisz') return zapiszKontekst(otoczenie, panel);
  if (kod === 'kontekst-wlacz') return wlaczKontekst(otoczenie, panel);
  if (kod === 'kontekst-usun') return usunKontekst(otoczenie, panel);
  if (kod === 'wpis-zapisz') return zapiszWpis(otoczenie, panel);
  if (kod === 'wpis-odepnij') return odepnijWpis(otoczenie, panel);
  if (kod === 'przelacz') return przestawPamiec(otoczenie, panel);
  if (kod === 'wyciszenia') return wykazWyciszen(otoczenie);
  if (kod === 'wycisz') return wyciszPoziom(otoczenie, panel);
  if (kod === 'przechowanie') return zasadyPrzechowania(otoczenie);
  if (kod === 'przechowanie-ustaw') return ustawPrzechowanie(otoczenie, panel);
  /* Czynność spoza obsłużonych odmawia zamiast milczeć: pozycja bez gałęzi
     wyglądałaby jak działająca. */
  oglos(NAGLOWEK, `Czynność „${kod}" nie jest prowadzona przez to okno.`, 'ostrzezenie');
}

async function wykazKontekstow({ kanal, korzen }: Otoczenie): Promise<void> {
  const wynik = await wywolaj(kanal, Command.MemoryContextList, { includeDisabled: true });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu kontekstów.', 'ostrzezenie');
    return;
  }
  const czynny = wynik.wynik.activeId;
  const konteksty = wynik.wynik.contexts;
  wypelnij(korzen, konteksty.map((kontekst: MemoryContext) =>
    `${kontekst.id === czynny ? 'czynny' : (kontekst.enabled ? 'gotowy' : 'wyłączony')}`
    + ` · ${kontekst.name} · wpisów ${kontekst.entryIds?.length ?? 0}`));
  if (konteksty.length === 0) oglos(NAGLOWEK, 'Nie ma jeszcze żadnego kontekstu pamięci.');
}

async function wyborKontekstow(kanal: Kanal): Promise<ReadonlyArray<readonly [string, string]>> {
  const wykaz = await wywolaj(kanal, Command.MemoryContextList, { includeDisabled: true });
  return (wykaz.wynik?.contexts ?? []).map((kontekst: MemoryContext) =>
    [kontekst.id, kontekst.name] as const);
}

/* Zapowiedź systemowa kontekstu idzie do modelu przed każdą rozmową, więc
   szuflada pyta o nią wprost, a nie chowa jej za ustawieniami. */
async function zapiszKontekst({ kanal, korzen }: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Kontekst pamięci',
    pola: [
      { klucz: 'nazwa', etykieta: 'Nazwa kontekstu', wymagane: true },
      { klucz: 'opis', etykieta: 'Opis', obszerne: true },
      { klucz: 'poziom', etykieta: 'Poziom pamięci', wybor: POZIOMY },
      { klucz: 'zapowiedz', etykieta: 'Zapowiedź dla modelu', obszerne: true },
    ],
    wykonanie: 'Zapisz kontekst',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.MemoryContextSave, {
    name: wartosci.nazwa ?? '',
    ...(wartosci.opis === '' ? {} : { description: wartosci.opis }),
    levels: [(wartosci.poziom ?? ConfigScope.Session) as ConfigScope],
    ...(wartosci.zapowiedz === '' ? {} : { systemPrompt: wartosci.zapowiedz }),
    enabled: true,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu kontekstu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Kontekst „${wartosci.nazwa ?? ''}" zapisany.`);
  await wykazKontekstow({ kanal, korzen, idSesji: () => '' });
}

async function wlaczKontekst({ kanal, korzen, idSesji }: Otoczenie, panel: string):
Promise<void> {
  if (idSesji() === '') {
    oglos(NAGLOWEK, 'Rdzeń nie dał sesji dla tej karty.', 'ostrzezenie');
    return;
  }
  const konteksty = await wyborKontekstow(kanal);
  if (konteksty.length === 0) {
    oglos(NAGLOWEK, 'Nie ma jeszcze żadnego kontekstu pamięci.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Kontekst w sesji',
    opis: 'Czynny kontekst wchodzi do każdej rozmowy tej sesji.',
    pola: [{ klucz: 'kontekst', etykieta: 'Kontekst', wybor: konteksty }],
    wykonanie: 'Włącz kontekst',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.MemoryContextActivate, {
    contextId: wartosci.kontekst ?? '',
    sessionId: idSesji(),
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił włączenia kontekstu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Kontekst czynny w tej sesji.');
  await wykazKontekstow({ kanal, korzen, idSesji });
}

async function usunKontekst({ kanal, korzen, idSesji }: Otoczenie, panel: string):
Promise<void> {
  const konteksty = await wyborKontekstow(kanal);
  if (konteksty.length === 0) {
    oglos(NAGLOWEK, 'Nie ma jeszcze żadnego kontekstu pamięci.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Usunięcie kontekstu',
    pola: [{ klucz: 'kontekst', etykieta: 'Kontekst', wybor: konteksty }],
    wykonanie: 'Usuń kontekst',
    nieodwracalne: 'Kontekst znika; wpisy pamięci, które zbierał, zostają.',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.MemoryContextDelete, {
    contextId: wartosci.kontekst ?? '',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił usunięcia kontekstu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Kontekst usunięty.');
  await wykazKontekstow({ kanal, korzen, idSesji });
}

/* Wpis oznaczony jako wrażliwy rdzeń trzyma inaczej i nie podaje go dalej bez
   potrzeby, więc okno pyta o to wprost. */
async function zapiszWpis({ kanal, korzen }: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Wpis pamięci',
    pola: [
      { klucz: 'tresc', etykieta: 'Treść wpisu', obszerne: true, wymagane: true },
      { klucz: 'poziom', etykieta: 'Poziom pamięci', wybor: POZIOMY },
      { klucz: 'dni', etykieta: 'Ważny przez dni', podpowiedz: 'Puste znaczy bez terminu' },
      {
        klucz: 'wrazliwy',
        etykieta: 'Treść wrażliwa',
        wybor: [['nie', 'Zwykła'], ['tak', 'Wrażliwa']],
      },
      {
        klucz: 'przypiety',
        etykieta: 'Wpis przypięty',
        wybor: [['nie', 'Zwykły'], ['tak', 'Przypięty do każdej rozmowy']],
      },
    ],
    wykonanie: 'Zapisz wpis',
  });
  if (wartosci === null) return;
  const dni = Number.parseInt(wartosci.dni ?? '', 10);
  const wynik = await wywolaj(kanal, Command.MemorySet, {
    content: wartosci.tresc ?? '',
    scope: (wartosci.poziom ?? ConfigScope.Session) as ConfigScope,
    ...(Number.isFinite(dni) ? { ttlDays: dni } : {}),
    sensitive: wartosci.wrazliwy === 'tak',
    pinned: wartosci.przypiety === 'tak',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu wpisu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Wpis pamięci zapisany.');
}

async function odepnijWpis({ kanal, korzen, idSesji }: Otoczenie, panel: string):
Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Odpięcie wpisu',
    opis: 'Wpis zostaje w pamięci, ale przestaje wchodzić do tej sesji.',
    pola: [{ klucz: 'wpis', etykieta: 'Oznaczenie wpisu', wymagane: true }],
    wykonanie: 'Odepnij wpis',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.MemoryDetach, {
    entryId: wartosci.wpis ?? '',
    ...(idSesji() === '' ? {} : { sessionId: idSesji() }),
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odpięcia wpisu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Wpis odpięty od tej sesji.');
}

async function przestawPamiec({ kanal, korzen, idSesji }: Otoczenie, panel: string):
Promise<void> {
  if (idSesji() === '') {
    oglos(NAGLOWEK, 'Rdzeń nie dał sesji dla tej karty.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Pamięć tej sesji',
    pola: [
      { klucz: 'poziom', etykieta: 'Poziom pamięci', wybor: POZIOMY },
      {
        klucz: 'zapis',
        etykieta: 'Zapisywanie nowych wpisów',
        wybor: [['tak', 'Włączone'], ['nie', 'Wyłączone']],
      },
    ],
    wykonanie: 'Przestaw pamięć',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.MemoryToggle, {
    sessionId: idSesji(),
    levels: [(wartosci.poziom ?? ConfigScope.Session) as ConfigScope],
    writeEnabled: wartosci.zapis !== 'nie',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił przestawienia pamięci.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, wartosci.zapis === 'nie'
    ? 'Ta sesja nie zapisuje już nowych wpisów.'
    : 'Ta sesja zapisuje nowe wpisy.');
}

async function wykazWyciszen({ kanal, korzen }: Otoczenie): Promise<void> {
  const wynik = await wywolaj(kanal, Command.MemoryDisableList, {});
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu wyciszeń.', 'ostrzezenie');
    return;
  }
  const wyciszenia = wynik.wynik.disables;
  wypelnij(korzen, wyciszenia.map((wyciszenie) =>
    `${wyciszenie.scope}${wyciszenie.level === undefined ? '' : ` · ${wyciszenie.level}`}`
    + `${wyciszenie.entryId === undefined ? '' : ` · wpis ${wyciszenie.entryId}`}`));
  if (wyciszenia.length === 0) oglos(NAGLOWEK, 'Nic nie jest wyciszone.');
}

/* Wyciszenie nie kasuje pamięci: wpisy zostają, ale przestają wchodzić do
   rozmowy — okno nazywa to wprost, żeby nikt nie brał go za usunięcie. */
async function wyciszPoziom({ kanal, korzen }: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Wyciszenie pamięci',
    opis: 'Wyciszenie nie kasuje wpisów — przestają tylko wchodzić do rozmowy.',
    pola: [
      {
        klucz: 'poziom',
        etykieta: 'Poziom',
        wybor: [
          [MemoryLevel.Global, 'Cała aplikacja'],
          [MemoryLevel.Project, 'Projekt'],
          [MemoryLevel.Session, 'Sesja'],
          [MemoryLevel.Environment, 'Środowisko'],
        ],
      },
      { klucz: 'zakres', etykieta: 'Zakres nastawy', wybor: POZIOMY },
      {
        klucz: 'stan',
        etykieta: 'Poziom pamięci',
        wybor: [['tak', 'Wyciszony'], ['nie', 'Słyszalny']],
      },
    ],
    wykonanie: 'Przestaw wyciszenie',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.MemoryDisableSet, {
    level: (wartosci.poziom ?? MemoryLevel.Session) as MemoryLevel,
    scope: (wartosci.zakres ?? ConfigScope.Session) as ConfigScope,
    disabled: wartosci.stan === 'tak',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił przestawienia wyciszenia.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, wartosci.stan === 'tak'
    ? 'Poziom pamięci wyciszony.'
    : 'Poziom pamięci znów wchodzi do rozmowy.');
}

async function zasadyPrzechowania({ kanal, korzen }: Otoczenie): Promise<void> {
  const wynik = await wywolaj(kanal, Command.MemoryRetentionGet, {});
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu zasad.', 'ostrzezenie');
    return;
  }
  const zasady = wynik.wynik.policies;
  wypelnij(korzen, zasady.map((zasada) =>
    `${zasada.scope} · dni ${zasada.ttlDays}`
    + `${zasada.sensitiveDefault ? ' · domyślnie wrażliwe' : ''}`
    + `${zasada.enabled ? '' : ' · wyłączona'}`));
  if (zasady.length === 0) oglos(NAGLOWEK, 'Rdzeń nie ma jeszcze zasad przechowania.');
}

/* Wzorce treści nigdy niezapisywanej chronią przed wpisaniem do pamięci tego,
   co nie powinno w niej stanąć — hasła, kluczy, numerów. */
async function ustawPrzechowanie({ kanal, korzen }: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Zasady przechowania pamięci',
    pola: [
      { klucz: 'zakres', etykieta: 'Zakres', wybor: POZIOMY },
      { klucz: 'dni', etykieta: 'Wpisy ważne przez dni', wartosc: '90', wymagane: true },
      {
        klucz: 'wrazliwe',
        etykieta: 'Nowe wpisy',
        wybor: [['nie', 'Zwykłe'], ['tak', 'Domyślnie wrażliwe']],
      },
      {
        klucz: 'wzorce',
        etykieta: 'Nigdy nie zapisuj treści pasującej do',
        podpowiedz: 'Wzorce rozdzielone spacją',
      },
    ],
    wykonanie: 'Ustaw zasady',
  });
  if (wartosci === null) return;
  const dni = Number.parseInt(wartosci.dni ?? '', 10);
  if (!Number.isFinite(dni)) {
    oglos(NAGLOWEK, 'Termin ważności musi być liczbą dni.', 'ostrzezenie');
    return;
  }
  const wzorce = (wartosci.wzorce ?? '').split(/\s+/).filter((wzorzec) => wzorzec !== '');
  const wynik = await wywolaj(kanal, Command.MemoryRetentionSet, {
    scope: (wartosci.zakres ?? ConfigScope.Session) as ConfigScope,
    ttlDays: dni,
    sensitiveDefault: wartosci.wrazliwe === 'tak',
    ...(wzorce.length === 0 ? {} : { neverStorePatterns: wzorce }),
    enabled: true,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił ustawienia zasad.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Wpisy tego zakresu żyją ${dni} dni.`);
}
