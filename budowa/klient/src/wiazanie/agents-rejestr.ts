// Konektory w oknie Agents: katalog rozszerzeń, ich narzędzia, uprawnienia,
// wersje, poświadczenia, zestawy, dzienniki i praca zbiorcza.
import {
  Command,
  ExtensionAuthKind,
  ExtensionBulkAction,
  ExtensionDefinitionFormat,
  ExtensionKind,
  ExtensionPermissionScope,
  ExtensionToolKind,
  ExtensionWebhookDirection,
  McpTransport,
} from '../../../shared/contract.ts';
import type { Extension, ExtensionToolEntry } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { dolozCzynnosciPanelu, zapytajWSzufladzie } from './czynnosci-okna.ts';
import { oglos } from './ogloszenie.ts';

const NAGLOWEK = 'Konektory';
const DOBA_MS = 86400000;

interface Otoczenie {
  kanal: Kanal;
  korzen: Element;
  odswiez: () => void;
}

export function zwiazRejestrKonektorow(
  kanal: Kanal,
  korzen: Element,
  odswiez: () => void,
  dolacz: (zdejmij: () => void) => void,
  przy: AddEventListenerOptions,
): void {
  const otoczenie: Otoczenie = { kanal, korzen, odswiez };
  const zdejmowanie = [
    dolozCzynnosciPanelu(korzen, 'panel-connectors', 'Czynności konektorów', [
      {
        naglowek: 'Katalog',
        pozycje: [
          { kod: 'katalog', nazwa: 'Wykaz konektorów…' },
          { kod: 'szukaj', nazwa: 'Szukaj konektora…' },
          { kod: 'szczegoly', nazwa: 'Szczegóły konektora…' },
          { kod: 'zdrowie', nazwa: 'Sprawdź dostępność' },
          { kod: 'aktualizacje', nazwa: 'Sprawdź nowe wydania' },
          { kod: 'zuzycie', nazwa: 'Zużycie konektorów…' },
        ],
      },
      {
        naglowek: 'Narzędzia',
        pozycje: [
          { kod: 'narzedzia', nazwa: 'Wykaz narzędzi…' },
          { kod: 'wywolaj', nazwa: 'Wywołaj narzędzie…' },
          { kod: 'proba', nazwa: 'Próba w piaskownicy…' },
          { kod: 'przeklad', nazwa: 'Zapisz przekład pól…' },
        ],
      },
      {
        naglowek: 'Wnoszenie',
        pozycje: [
          { kod: 'zestaw-wnies', nazwa: 'Wnieś zestaw z manifestu…' },
          { kod: 'opis-wnies', nazwa: 'Wnieś z opisu usługi…' },
          { kod: 'pakiet-wnies', nazwa: 'Wnieś pakiet…' },
        ],
      },
    ], (kod) => {
      void wykonaj(otoczenie, kod, 'panel-connectors');
    }, przy),
    dolozCzynnosciPanelu(korzen, 'panel-permissions', 'Zaufanie konektorów', [
      {
        naglowek: 'Uprawnienia',
        pozycje: [
          { kod: 'uprawnienia', nazwa: 'Wykaz uprawnień…' },
          { kod: 'uprawnienie-nadaj', nazwa: 'Nadaj uprawnienie…' },
        ],
      },
      {
        naglowek: 'Poświadczenia i tajemnice',
        pozycje: [
          { kod: 'poswiadczenie', nazwa: 'Zwiąż poświadczenie…' },
          { kod: 'tajemnice', nazwa: 'Wykaz tajemnic…' },
          { kod: 'tajemnica-udostepnij', nazwa: 'Udostępnij tajemnicę…' },
        ],
      },
      {
        naglowek: 'Wiarygodność',
        pozycje: [
          { kod: 'podpis', nazwa: 'Sprawdź podpis…' },
          { kod: 'manifest', nazwa: 'Przejrzyj manifest…' },
        ],
      },
    ], (kod) => {
      void wykonaj(otoczenie, kod, 'panel-permissions');
    }, przy),
    dolozCzynnosciPanelu(korzen, 'panel-artefakty', 'Wydania i dzienniki', [
      {
        naglowek: 'Wydania',
        pozycje: [
          { kod: 'wersja-przypnij', nazwa: 'Przypnij wydanie…' },
          { kod: 'wersja-cofnij', nazwa: 'Cofnij do wydania…' },
          { kod: 'dzieje', nazwa: 'Dzieje konektora…' },
        ],
      },
      {
        naglowek: 'Dzienniki',
        pozycje: [
          { kod: 'dziennik-protokolu', nazwa: 'Dziennik rozmowy…' },
          { kod: 'dziennik-czynnosci', nazwa: 'Dziennik czynności…' },
        ],
      },
      {
        naglowek: 'Zestawy i praca zbiorcza',
        pozycje: [
          { kod: 'zestawy', nazwa: 'Wykaz zestawów' },
          { kod: 'zestaw-zapisz', nazwa: 'Zapisz zestaw…' },
          { kod: 'zestaw-zastosuj', nazwa: 'Zastosuj zestaw…' },
          { kod: 'zbiorczo', nazwa: 'Przestaw wiele naraz…' },
        ],
      },
      {
        naglowek: 'Łącze i wejścia sieciowe',
        pozycje: [
          { kod: 'lacze', nazwa: 'Ustaw łącze…' },
          { kod: 'wejscia', nazwa: 'Wykaz wejść sieciowych…' },
          { kod: 'wejscie-zapisz', nazwa: 'Zapisz wejście sieciowe…' },
        ],
      },
    ], (kod) => {
      void wykonaj(otoczenie, kod, 'panel-artefakty');
    }, przy),
  ];
  for (const zdejmij of zdejmowanie) if (zdejmij !== null) dolacz(zdejmij);
}

function wypelnij(korzen: Element, panelKod: string, wiersze: string[]): void {
  const cialo = korzen.querySelector(`#${panelKod} .sta-okno-tresc`);
  if (cialo === null) return;
  cialo.replaceChildren(...wiersze.map((tresc) => {
    const wiersz = cialo.ownerDocument.createElement('div');
    wiersz.className = 'dn-wykaz-modulu-poz';
    wiersz.textContent = tresc;
    return wiersz;
  }));
}

/* Konektor wskazuje wykaz rejestru: okno nie każe Operatorowi znać oznaczeń,
   które rdzeń i tak zna. */
async function wyborKonektorow(
  kanal: Kanal,
): Promise<ReadonlyArray<readonly [string, string]>> {
  const wykaz = await wywolaj(kanal, Command.ExtensionRegistryList, {});
  return (wykaz.wynik?.extensions ?? []).map((rozszerzenie: Extension) =>
    [rozszerzenie.id, `${rozszerzenie.name} · ${rozszerzenie.kind}`
      + `${rozszerzenie.installed ? '' : ' · niezainstalowany'}`] as const);
}

async function wskazKonektor(
  otoczenie: Otoczenie,
  panel: string,
  tytul: string,
  wykonanie: string,
  dodatkowe: Parameters<typeof zapytajWSzufladzie>[2]['pola'] = [],
  nieodwracalne?: string,
): Promise<Record<string, string> | null> {
  const wybor = await wyborKonektorow(otoczenie.kanal);
  if (wybor.length === 0) {
    oglos(NAGLOWEK, 'Rejestr nie ma jeszcze żadnego konektora.', 'ostrzezenie');
    return null;
  }
  return zapytajWSzufladzie(otoczenie.korzen, panel, {
    tytul,
    pola: [{ klucz: 'konektor', etykieta: 'Konektor', wybor }, ...dodatkowe],
    wykonanie,
    ...(nieodwracalne === undefined ? {} : { nieodwracalne }),
  });
}

async function wykonaj(otoczenie: Otoczenie, kod: string, panel: string): Promise<void> {
  if (kod === 'katalog') return wykazKonektorow(otoczenie, panel);
  if (kod === 'szukaj') return szukajKonektora(otoczenie, panel);
  if (kod === 'szczegoly') return szczegolyKonektora(otoczenie, panel);
  if (kod === 'zdrowie') return sprawdzDostepnosc(otoczenie);
  if (kod === 'aktualizacje') return sprawdzWydania(otoczenie);
  if (kod === 'zuzycie') return zuzycieKonektorow(otoczenie, panel);
  if (kod === 'narzedzia') return wykazNarzedzi(otoczenie, panel);
  if (kod === 'wywolaj') return wywolajNarzedzie(otoczenie, panel);
  if (kod === 'proba') return probaWPiaskownicy(otoczenie, panel);
  if (kod === 'przeklad') return zapiszPrzeklad(otoczenie, panel);
  if (kod === 'zestaw-wnies') return wniesZestaw(otoczenie, panel);
  if (kod === 'opis-wnies') return wniesZOpisu(otoczenie, panel);
  if (kod === 'pakiet-wnies') return wniesPakiet(otoczenie, panel);
  if (kod === 'uprawnienia') return wykazUprawnien(otoczenie, panel);
  if (kod === 'uprawnienie-nadaj') return nadajUprawnienie(otoczenie, panel);
  if (kod === 'poswiadczenie') return zwiazPoswiadczenie(otoczenie, panel);
  if (kod === 'tajemnice') return wykazTajemnic(otoczenie, panel);
  if (kod === 'tajemnica-udostepnij') return udostepnijTajemnice(otoczenie, panel);
  if (kod === 'podpis') return sprawdzPodpis(otoczenie, panel);
  if (kod === 'manifest') return przejrzyjManifest(otoczenie, panel);
  if (kod === 'wersja-przypnij') return przypnijWydanie(otoczenie, panel);
  if (kod === 'wersja-cofnij') return cofnijWydanie(otoczenie, panel);
  if (kod === 'dzieje') return dziejeKonektora(otoczenie, panel);
  if (kod === 'dziennik-protokolu') return dziennikProtokolu(otoczenie, panel);
  if (kod === 'dziennik-czynnosci') return dziennikCzynnosci(otoczenie, panel);
  if (kod === 'zestawy') return wykazZestawow(otoczenie);
  if (kod === 'zestaw-zapisz') return zapiszZestaw(otoczenie, panel);
  if (kod === 'zestaw-zastosuj') return zastosujZestaw(otoczenie, panel);
  if (kod === 'zbiorczo') return przestawWieleNaraz(otoczenie, panel);
  if (kod === 'lacze') return ustawLacze(otoczenie, panel);
  if (kod === 'wejscia') return wykazWejsc(otoczenie, panel);
  if (kod === 'wejscie-zapisz') return zapiszWejscie(otoczenie, panel);
  /* Czynność spoza obsłużonych odmawia zamiast milczeć: pozycja bez gałęzi
     wyglądałaby jak działająca. */
  oglos(NAGLOWEK, `Czynność „${kod}" nie jest prowadzona przez to okno.`, 'ostrzezenie');
}

const RODZAJE: ReadonlyArray<readonly [string, string]> = [
  ['', 'Wszystkie rodzaje'],
  [ExtensionKind.Mcp, 'Most MCP'],
  [ExtensionKind.Plugin, 'Wtyczka'],
  [ExtensionKind.Api, 'Usługa'],
  [ExtensionKind.Skill, 'Umiejętność'],
];

async function wykazKonektorow(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(otoczenie.korzen, panel, {
    tytul: 'Wykaz konektorów',
    pola: [{ klucz: 'rodzaj', etykieta: 'Rodzaj', wybor: RODZAJE }],
    wykonanie: 'Odczytaj wykaz',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.ExtensionRegistryList, {
    ...(wartosci.rodzaj === '' ? {} : { kind: wartosci.rodzaj as ExtensionKind }),
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu rejestru.', 'ostrzezenie');
    return;
  }
  const konektory = wynik.wynik.extensions;
  wypelnij(otoczenie.korzen, 'panel-connectors', konektory.map((konektor: Extension) =>
    `${konektor.installed ? 'zainstalowany' : 'w rejestrze'} · ${konektor.name} · `
    + `${konektor.kind}${konektor.version === undefined ? '' : ` · ${konektor.version}`}`));
  if (konektory.length === 0) oglos(NAGLOWEK, 'Rejestr nie ma konektorów tego rodzaju.');
}

async function szukajKonektora(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(otoczenie.korzen, panel, {
    tytul: 'Szukanie konektora',
    pola: [
      { klucz: 'wzorzec', etykieta: 'Szukana treść', wymagane: true },
      { klucz: 'rodzaj', etykieta: 'Rodzaj', wybor: RODZAJE },
      {
        klucz: 'zakres',
        etykieta: 'Zakres szukania',
        wybor: [['wszystkie', 'Cały rejestr'], ['zainstalowane', 'Tylko zainstalowane']],
      },
    ],
    wykonanie: 'Szukaj',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.ExtensionSearch, {
    query: wartosci.wzorzec ?? '',
    ...(wartosci.rodzaj === '' ? {} : { kind: wartosci.rodzaj as ExtensionKind }),
    installedOnly: wartosci.zakres === 'zainstalowane',
    limit: 100,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił szukania.', 'ostrzezenie');
    return;
  }
  const znalezione = wynik.wynik.extensions;
  wypelnij(otoczenie.korzen, 'panel-connectors', znalezione.map((konektor: Extension) =>
    `${konektor.name} · ${konektor.kind}`));
  if (znalezione.length === 0) oglos(NAGLOWEK, 'Nic nie odpowiada tej treści.');
}

async function szczegolyKonektora(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazKonektor(otoczenie, panel, 'Szczegóły konektora', 'Odczytaj');
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.ExtensionDetailGet, {
    extensionId: wartosci.konektor ?? '',
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu szczegółów.', 'ostrzezenie');
    return;
  }
  const szczegol = wynik.wynik.detail;
  wypelnij(otoczenie.korzen, 'panel-connectors', [
    `konektor: ${szczegol.extensionId}`,
    `wydawca: ${szczegol.publisher ?? 'bez wskazania'}`,
    `narzędzi: ${szczegol.tools?.length ?? 0}`,
    `uprawnień: ${szczegol.permissions?.length ?? 0}`,
    `strona: ${szczegol.homepageUrl ?? 'bez adresu'}`,
  ]);
}

async function sprawdzDostepnosc({ kanal, korzen }: Otoczenie): Promise<void> {
  const wynik = await wywolaj(kanal, Command.ExtensionHealthCheck, {});
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił sprawdzenia dostępności.', 'ostrzezenie');
    return;
  }
  wypelnij(korzen, 'panel-connectors', wynik.wynik.results.map((stan) =>
    `${stan.status} · ${stan.extensionId}`
    + `${stan.handshakeError === undefined ? '' : ` · ${stan.handshakeError}`}`));
  if (wynik.wynik.results.length === 0) oglos(NAGLOWEK, 'Żaden konektor nie jest zainstalowany.');
}

async function sprawdzWydania({ kanal, korzen }: Otoczenie): Promise<void> {
  const wynik = await wywolaj(kanal, Command.ExtensionUpdateCheck, {});
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił sprawdzenia wydań.', 'ostrzezenie');
    return;
  }
  const wydania = wynik.wynik.updates;
  wypelnij(korzen, 'panel-connectors', wydania.map((wydanie) =>
    `${wydanie.extensionId} · ${wydanie.currentVersion} → ${wydanie.availableVersion}`
    + `${wydanie.breaking === true ? ' · zmiana łamiąca' : ''}`));
  if (wydania.length === 0) oglos(NAGLOWEK, 'Żaden konektor nie ma nowszego wydania.');
}

/* Zakres zużycia liczony jest wstecz od chwili odczytu, bo okno nie prowadzi
   wyboru dat, a pytanie brzmi „ile zużyły przez ostatnie dni". */
async function zuzycieKonektorow(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(otoczenie.korzen, panel, {
    tytul: 'Zużycie konektorów',
    pola: [{ klucz: 'dni', etykieta: 'Ostatnich dni', wartosc: '7' }],
    wykonanie: 'Odczytaj zużycie',
  });
  if (wartosci === null) return;
  const dni = Number.parseInt(wartosci.dni ?? '', 10);
  const wynik = await wywolaj(otoczenie.kanal, Command.ExtensionUsageGet, {
    since: Date.now() - (Number.isFinite(dni) ? dni : 7) * DOBA_MS,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu zużycia.', 'ostrzezenie');
    return;
  }
  wypelnij(otoczenie.korzen, 'panel-connectors', wynik.wynik.usage.map((pozycja) =>
    `${pozycja.extensionId} · wywołań ${pozycja.calls ?? 0}`));
  if (wynik.wynik.usage.length === 0) oglos(NAGLOWEK, 'Żaden konektor nie był wołany w tym czasie.');
}

async function wykazNarzedzi(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazKonektor(otoczenie, panel, 'Narzędzia konektora', 'Odczytaj wykaz', [
    {
      klucz: 'rodzaj',
      etykieta: 'Rodzaj pozycji',
      wybor: [
        ['', 'Wszystkie'],
        [ExtensionToolKind.Tool, 'Narzędzia'],
        [ExtensionToolKind.Resource, 'Zasoby'],
        [ExtensionToolKind.Prompt, 'Polecenia'],
      ],
    },
    {
      klucz: 'odswiez',
      etykieta: 'Wykaz',
      wybor: [['nie', 'Z pamięci rdzenia'], ['tak', 'Pobrany od konektora na nowo']],
    },
  ]);
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.ExtensionToolList, {
    extensionId: wartosci.konektor ?? '',
    ...(wartosci.rodzaj === '' ? {} : { kind: wartosci.rodzaj as ExtensionToolKind }),
    refresh: wartosci.odswiez === 'tak',
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu narzędzi.', 'ostrzezenie');
    return;
  }
  const narzedzia = wynik.wynik.entries;
  wypelnij(otoczenie.korzen, 'panel-connectors', narzedzia.map((narzedzie: ExtensionToolEntry) =>
    `${narzedzie.kind} · ${narzedzie.name}`
    + `${narzedzie.description === undefined ? '' : ` · ${narzedzie.description}`}`));
  if (narzedzia.length === 0) oglos(NAGLOWEK, 'Ten konektor nie podaje narzędzi.');
}

/* Argumenty podaje się zapisem JSON, bo rdzeń przekazuje je konektorowi
   nietknięte; zapis niepoprawny odmawia zamiast wysyłać tekst. */
async function wywolajNarzedzie(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazKonektor(otoczenie, panel, 'Wywołanie narzędzia',
    'Wywołaj narzędzie', [
      { klucz: 'narzedzie', etykieta: 'Nazwa narzędzia', wymagane: true },
      { klucz: 'argumenty', etykieta: 'Argumenty w zapisie JSON', wartosc: '{}', obszerne: true },
    ],
    'Konektor wykona to, co narzędzie robi — także poza tą maszyną.');
  if (wartosci === null) return;
  let argumenty: unknown;
  try {
    argumenty = JSON.parse(wartosci.argumenty === '' ? '{}' : (wartosci.argumenty ?? '{}'));
  } catch {
    oglos(NAGLOWEK, 'Argumenty nie są poprawnym zapisem JSON.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(otoczenie.kanal, Command.ExtensionToolCall, {
    extensionId: wartosci.konektor ?? '',
    toolName: wartosci.narzedzie ?? '',
    arguments: argumenty,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wywołania narzędzia.', 'ostrzezenie');
    return;
  }
  wypelnij(otoczenie.korzen, 'panel-connectors',
    (wynik.wynik.text ?? JSON.stringify(wynik.wynik.raw ?? null, null, 2)).split('\n').slice(0, 60));
  oglos(NAGLOWEK, wynik.wynik.ok
    ? `Narzędzie odpowiedziało w ${wynik.wynik.durationMs} ms.`
    : `Narzędzie zawiodło: ${wynik.wynik.errorDetail ?? 'bez wyjaśnienia'}.`,
  wynik.wynik.ok ? 'informacja' : 'ostrzezenie');
}

async function probaWPiaskownicy(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazKonektor(otoczenie, panel, 'Próba w piaskownicy',
    'Przeprowadź próbę', [
      { klucz: 'wejscie', etykieta: 'Wejście w zapisie JSON', wartosc: '{}', obszerne: true },
    ]);
  if (wartosci === null) return;
  let wejscie: unknown;
  try {
    wejscie = JSON.parse(wartosci.wejscie === '' ? '{}' : (wartosci.wejscie ?? '{}'));
  } catch {
    oglos(NAGLOWEK, 'Wejście nie jest poprawnym zapisem JSON.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(otoczenie.kanal, Command.ExtensionSandboxRun, {
    extensionId: wartosci.konektor ?? '',
    input: wejscie,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił próby.', 'ostrzezenie');
    return;
  }
  wypelnij(otoczenie.korzen, 'panel-connectors',
    JSON.stringify(wynik.wynik.output ?? null, null, 2).split('\n').slice(0, 60));
  oglos(NAGLOWEK, wynik.wynik.ok
    ? `Próba wykonana w ${wynik.wynik.durationMs} ms.`
    : 'Próba w piaskownicy nie powiodła się.',
  wynik.wynik.ok ? 'informacja' : 'ostrzezenie');
}

async function zapiszPrzeklad(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazKonektor(otoczenie, panel, 'Przekład pól', 'Zapisz przekład', [
    { klucz: 'nazwa', etykieta: 'Nazwa przekładu', wymagane: true },
    { klucz: 'zasady', etykieta: 'Zasady w zapisie JSON', wartosc: '{}', obszerne: true },
  ]);
  if (wartosci === null) return;
  let zasady: unknown;
  try {
    zasady = JSON.parse(wartosci.zasady === '' ? '{}' : (wartosci.zasady ?? '{}'));
  } catch {
    oglos(NAGLOWEK, 'Zasady nie są poprawnym zapisem JSON.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(otoczenie.kanal, Command.ExtensionMappingSave, {
    extensionId: wartosci.konektor ?? '',
    name: wartosci.nazwa ?? '',
    rules: zasady,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu przekładu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Przekład „${wartosci.nazwa ?? ''}" zapisany.`);
}

/* Wniesienie zestawu wprowadza konektory z cudzego manifestu; szuflada mówi to
   wprost, bo wniesione konektory zaczynają działać na tej maszynie. */
async function wniesZestaw(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(otoczenie.korzen, panel, {
    tytul: 'Zestaw z manifestu',
    pola: [
      { klucz: 'manifest', etykieta: 'Manifest w zapisie JSON', obszerne: true, wymagane: true },
      {
        klucz: 'wlacz',
        etykieta: 'Konektory zestawu',
        wybor: [['nie', 'Wniesione wyłączone'], ['tak', 'Wniesione i włączone']],
      },
    ],
    wykonanie: 'Wnieś zestaw',
    nieodwracalne: 'Konektory z manifestu zaczną działać na maszynie rdzenia.',
  });
  if (wartosci === null) return;
  let manifest: unknown;
  try {
    manifest = JSON.parse(wartosci.manifest ?? '');
  } catch {
    oglos(NAGLOWEK, 'Manifest nie jest poprawnym zapisem JSON.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(otoczenie.kanal, Command.ExtensionBundleInstall, {
    manifest,
    enable: wartosci.wlacz === 'tak',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wniesienia zestawu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Zestaw wniesiony.');
  otoczenie.odswiez();
}

async function wniesZOpisu(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(otoczenie.korzen, panel, {
    tytul: 'Konektor z opisu usługi',
    pola: [
      { klucz: 'kod', etykieta: 'Oznaczenie konektora', wymagane: true },
      {
        klucz: 'postac',
        etykieta: 'Postać opisu',
        wybor: [
          [ExtensionDefinitionFormat.Openapi3, 'OpenAPI 3'],
          [ExtensionDefinitionFormat.Graphql, 'GraphQL'],
        ],
      },
      { klucz: 'zrodlo', etykieta: 'Adres opisu', podpowiedz: 'Albo wklej treść niżej' },
      { klucz: 'tresc', etykieta: 'Treść opisu', obszerne: true },
    ],
    wykonanie: 'Wnieś konektor',
  });
  if (wartosci === null) return;
  if (wartosci.zrodlo === '' && wartosci.tresc === '') {
    oglos(NAGLOWEK, 'Opis potrzebuje adresu albo wklejonej treści.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(otoczenie.kanal, Command.ExtensionDefinitionImport, {
    code: wartosci.kod ?? '',
    format: (wartosci.postac ?? ExtensionDefinitionFormat.Openapi3) as ExtensionDefinitionFormat,
    ...(wartosci.zrodlo === '' ? {} : { source: wartosci.zrodlo }),
    ...(wartosci.tresc === '' ? {} : { content: wartosci.tresc }),
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wniesienia z opisu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Konektor „${wartosci.kod ?? ''}" wniesiony z opisu usługi.`);
  otoczenie.odswiez();
}

async function wniesPakiet(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(otoczenie.korzen, panel, {
    tytul: 'Pakiet konektora',
    opis: 'Treść pakietu podaje się zapisem base64; rdzeń sprawdzi jego sumę kontrolną.',
    pola: [
      { klucz: 'nazwa', etykieta: 'Nazwa pliku', wymagane: true },
      { klucz: 'tresc', etykieta: 'Pakiet w zapisie base64', obszerne: true, wymagane: true },
      { klucz: 'suma', etykieta: 'Suma kontrolna SHA-256' },
    ],
    wykonanie: 'Wnieś pakiet',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.ExtensionPackageUpload, {
    fileName: wartosci.nazwa ?? '',
    contentBase64: wartosci.tresc ?? '',
    ...(wartosci.suma === '' ? {} : { checksumSha256: wartosci.suma }),
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił przyjęcia pakietu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Pakiet ${wartosci.nazwa ?? ''} przyjęty.`);
  otoczenie.odswiez();
}

async function wykazUprawnien(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazKonektor(otoczenie, panel, 'Uprawnienia konektora', 'Odczytaj');
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.ExtensionPermissionList, {
    extensionId: wartosci.konektor ?? '',
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu uprawnień.', 'ostrzezenie');
    return;
  }
  const nadane = new Set(wynik.wynik.granted.map((uprawnienie) => uprawnienie.scope));
  const deklarowane = wynik.wynik.declared;
  wypelnij(otoczenie.korzen, 'panel-permissions', deklarowane.map((uprawnienie) =>
    `${uprawnienie.scope}${uprawnienie.target === undefined ? '' : ` · ${uprawnienie.target}`}`
    + `${nadane.has(uprawnienie.scope) ? ' · nadane' : ' · nienadane'}`));
  const nadmiarowe = wynik.wynik.excessive ?? [];
  if (nadmiarowe.length > 0) {
    oglos(NAGLOWEK, `Rdzeń uznaje za nadmiarowe: ${nadmiarowe.join(', ')}.`, 'ostrzezenie');
    return;
  }
  if (deklarowane.length === 0) oglos(NAGLOWEK, 'Ten konektor nie prosi o żadne uprawnienie.');
}

/* Nadanie uprawnienia wpuszcza konektor tam, gdzie wcześniej nie sięgał, więc
   szuflada nazywa zakres słowami i wymaga drugiego naciśnięcia. */
async function nadajUprawnienie(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazKonektor(otoczenie, panel, 'Nadanie uprawnienia',
    'Nadaj uprawnienie', [
      {
        klucz: 'zakres',
        etykieta: 'Zakres',
        wybor: [
          [ExtensionPermissionScope.Network, 'Sieć'],
          [ExtensionPermissionScope.FileRead, 'Odczyt plików'],
          [ExtensionPermissionScope.FileWrite, 'Zapis plików'],
          [ExtensionPermissionScope.ProcessSpawn, 'Uruchamianie procesów'],
          [ExtensionPermissionScope.SecretRead, 'Odczyt tajemnic'],
          [ExtensionPermissionScope.ModelCall, 'Wywołania modelu'],
        ],
      },
      { klucz: 'cel', etykieta: 'Zawężenie', podpowiedz: 'Ścieżka albo adres; puste znaczy całość' },
      { klucz: 'powod', etykieta: 'Powód nadania' },
    ],
    'Konektor sięgnie tam, gdzie wcześniej nie sięgał.');
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.ExtensionPermissionGrant, {
    extensionId: wartosci.konektor ?? '',
    permissions: [{
      scope: (wartosci.zakres ?? ExtensionPermissionScope.Network) as ExtensionPermissionScope,
      ...(wartosci.cel === '' ? {} : { target: wartosci.cel }),
      ...(wartosci.powod === '' ? {} : { explanation: wartosci.powod }),
    }],
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił nadania uprawnienia.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Uprawnienie nadane.');
}

async function zwiazPoswiadczenie(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazKonektor(otoczenie, panel, 'Poświadczenie konektora',
    'Zwiąż poświadczenie', [
      {
        klucz: 'sposob',
        etykieta: 'Sposób uwierzytelnienia',
        wybor: [
          [ExtensionAuthKind.ApiKey, 'Klucz usługi'],
          [ExtensionAuthKind.Token, 'Żeton'],
          [ExtensionAuthKind.Oauth2, 'OAuth 2'],
          [ExtensionAuthKind.Basic, 'Nazwa i hasło'],
          [ExtensionAuthKind.None, 'Bez uwierzytelnienia'],
        ],
      },
      { klucz: 'odwolanie', etykieta: 'Odwołanie do sejfu', wymagane: true },
      { klucz: 'zakresy', etykieta: 'Zakresy', podpowiedz: 'Rozdzielone spacją' },
    ]);
  if (wartosci === null) return;
  const zakresy = (wartosci.zakresy ?? '').split(/\s+/).filter((zakres) => zakres !== '');
  const wynik = await wywolaj(otoczenie.kanal, Command.ExtensionCredentialBind, {
    extensionId: wartosci.konektor ?? '',
    authKind: (wartosci.sposob ?? ExtensionAuthKind.ApiKey) as ExtensionAuthKind,
    credentialRef: wartosci.odwolanie ?? '',
    ...(zakresy.length === 0 ? {} : { scopes: zakresy }),
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił związania poświadczenia.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Poświadczenie związane z konektorem.');
}

async function wykazTajemnic(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(otoczenie.korzen, panel, {
    tytul: 'Tajemnice konektorów',
    pola: [{ klucz: 'dni', etykieta: 'Wygasające w ciągu dni', wartosc: '30' }],
    wykonanie: 'Odczytaj tajemnice',
  });
  if (wartosci === null) return;
  const dni = Number.parseInt(wartosci.dni ?? '', 10);
  const wynik = await wywolaj(otoczenie.kanal, Command.ExtensionSecretList, {
    ...(Number.isFinite(dni) ? { expiringWithinDays: dni } : {}),
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu tajemnic.', 'ostrzezenie');
    return;
  }
  wypelnij(otoczenie.korzen, 'panel-permissions', wynik.wynik.secrets.map((tajemnica) =>
    `${tajemnica.label ?? tajemnica.ref}`
    + `${tajemnica.expiresAt === undefined ? '' : ' · z terminem'}`
    + `${tajemnica.sharedWithExtensionIds === undefined
      ? '' : ` · dzielona z ${tajemnica.sharedWithExtensionIds.length}`}`));
  if (wynik.wynik.secrets.length === 0) {
    oglos(NAGLOWEK, 'Żadna tajemnica nie wygasa w tym czasie.');
  }
}

/* Udostępnienie tajemnicy daje wskazanym konektorom dostęp do tej samej
   wartości, więc szuflada nazywa następstwo i wymaga potwierdzenia. */
async function udostepnijTajemnice(otoczenie: Otoczenie, panel: string): Promise<void> {
  const konektory = await wyborKonektorow(otoczenie.kanal);
  if (konektory.length === 0) {
    oglos(NAGLOWEK, 'Rejestr nie ma jeszcze żadnego konektora.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWSzufladzie(otoczenie.korzen, panel, {
    tytul: 'Udostępnienie tajemnicy',
    pola: [
      { klucz: 'odwolanie', etykieta: 'Odwołanie do sejfu', wymagane: true },
      { klucz: 'konektor', etykieta: 'Konektor', wybor: konektory },
    ],
    wykonanie: 'Udostępnij tajemnicę',
    nieodwracalne: 'Wskazany konektor sięgnie po tę samą wartość z sejfu.',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.ExtensionSecretShare, {
    secretRef: wartosci.odwolanie ?? '',
    extensionIds: [wartosci.konektor ?? ''],
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił udostępnienia tajemnicy.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Tajemnica udostępniona konektorowi.');
}

async function sprawdzPodpis(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazKonektor(otoczenie, panel, 'Podpis konektora', 'Sprawdź podpis');
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.ExtensionSignatureVerify, {
    extensionId: wartosci.konektor ?? '',
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił sprawdzenia podpisu.', 'ostrzezenie');
    return;
  }
  const podpis = wynik.wynik.signature;
  oglos(NAGLOWEK, podpis.verified
    ? `Podpis potwierdzony; zaufanie ${podpis.trustLevel}`
    + `${podpis.publisher === undefined ? '' : `, wydawca ${podpis.publisher}`}.`
    : `Podpisu nie potwierdzono: ${podpis.detail ?? 'bez wyjaśnienia'}.`,
  podpis.verified ? 'informacja' : 'ostrzezenie');
}

async function przejrzyjManifest(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazKonektor(otoczenie, panel, 'Przegląd manifestu', 'Przejrzyj');
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.ExtensionManifestScan, {
    extensionId: wartosci.konektor ?? '',
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił przeglądu manifestu.', 'ostrzezenie');
    return;
  }
  const uwagi = wynik.wynik.findings;
  wypelnij(otoczenie.korzen, 'panel-permissions', uwagi.map((uwaga) =>
    `${uwaga.severity} · ${uwaga.message}`));
  oglos(NAGLOWEK, uwagi.length === 0
    ? 'Manifest bez uwag.'
    : `Manifest ma ${uwagi.length} uwag.`,
  uwagi.length === 0 ? 'informacja' : 'ostrzezenie');
}

async function przypnijWydanie(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazKonektor(otoczenie, panel, 'Przypięcie wydania', 'Przypnij', [
    { klucz: 'wydanie', etykieta: 'Wydanie', podpowiedz: 'Puste zdejmuje przypięcie' },
  ]);
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.ExtensionVersionPin, {
    extensionId: wartosci.konektor ?? '',
    ...(wartosci.wydanie === '' ? {} : { version: wartosci.wydanie }),
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił przypięcia wydania.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, wartosci.wydanie === ''
    ? 'Przypięcie zdjęte; konektor pójdzie za nowszymi wydaniami.'
    : `Konektor stoi na wydaniu ${wartosci.wydanie}.`);
}

async function cofnijWydanie(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazKonektor(otoczenie, panel, 'Cofnięcie wydania', 'Cofnij wydanie', [
    { klucz: 'wydanie', etykieta: 'Wydanie docelowe', wymagane: true },
  ],
  'Konektor wróci do wskazanego wydania wraz z jego zachowaniem.');
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.ExtensionVersionRollback, {
    extensionId: wartosci.konektor ?? '',
    targetVersion: wartosci.wydanie ?? '',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił cofnięcia wydania.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Konektor cofnięty do wydania ${wartosci.wydanie ?? ''}.`);
  otoczenie.odswiez();
}

async function dziejeKonektora(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazKonektor(otoczenie, panel, 'Dzieje konektora', 'Odczytaj dzieje');
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.ExtensionHistoryList, {
    extensionId: wartosci.konektor ?? '',
    limit: 100,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu dziejów.', 'ostrzezenie');
    return;
  }
  wypelnij(otoczenie.korzen, 'panel-artefakty', wynik.wynik.entries.map((wpis) =>
    `${wpis.action}${wpis.toVersion === undefined ? '' : ` → ${wpis.toVersion}`}`));
  if (wynik.wynik.entries.length === 0) oglos(NAGLOWEK, 'Ten konektor nie ma jeszcze dziejów.');
}

async function dziennikProtokolu(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazKonektor(otoczenie, panel, 'Dziennik rozmowy z konektorem',
    'Odczytaj dziennik');
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.ExtensionProtocolLogList, {
    extensionId: wartosci.konektor ?? '',
    limit: 100,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu dziennika.', 'ostrzezenie');
    return;
  }
  wypelnij(otoczenie.korzen, 'panel-artefakty', wynik.wynik.frames.map((ramka) =>
    `${ramka.direction} · ${ramka.method ?? 'bez nazwy'}`));
  if (wynik.wynik.frames.length === 0) oglos(NAGLOWEK, 'Rdzeń nie rozmawiał jeszcze z tym konektorem.');
}

async function dziennikCzynnosci(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(otoczenie.korzen, panel, {
    tytul: 'Dziennik czynności konektorów',
    pola: [{ klucz: 'dni', etykieta: 'Ostatnich dni', wartosc: '7' }],
    wykonanie: 'Odczytaj dziennik',
  });
  if (wartosci === null) return;
  const dni = Number.parseInt(wartosci.dni ?? '', 10);
  const wynik = await wywolaj(otoczenie.kanal, Command.ExtensionAuditList, {
    since: Date.now() - (Number.isFinite(dni) ? dni : 7) * DOBA_MS,
    limit: 200,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu dziennika.', 'ostrzezenie');
    return;
  }
  wypelnij(otoczenie.korzen, 'panel-artefakty', wynik.wynik.entries.map((wpis) =>
    `${wpis.extensionId}${wpis.toolName === undefined ? '' : ` · ${wpis.toolName}`}`
    + `${wpis.agentId === undefined ? '' : ` · agent ${wpis.agentId}`}`));
  if (wynik.wynik.entries.length === 0) oglos(NAGLOWEK, 'Dziennik czynności jest pusty.');
}

async function wykazZestawow({ kanal, korzen }: Otoczenie): Promise<void> {
  const wynik = await wywolaj(kanal, Command.ExtensionCollectionList, {});
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu zestawów.', 'ostrzezenie');
    return;
  }
  const zestawy = wynik.wynik.collections;
  wypelnij(korzen, 'panel-artefakty', zestawy.map((zestaw) =>
    `${zestaw.name} · konektorów ${zestaw.extensionIds?.length ?? 0}`));
  if (zestawy.length === 0) oglos(NAGLOWEK, 'Nie ma jeszcze żadnego zestawu.');
}

async function zapiszZestaw(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazKonektor(otoczenie, panel, 'Zestaw konektorów', 'Zapisz zestaw', [
    { klucz: 'nazwa', etykieta: 'Nazwa zestawu', wymagane: true },
    { klucz: 'opis', etykieta: 'Opis', obszerne: true },
  ]);
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.ExtensionCollectionSave, {
    name: wartosci.nazwa ?? '',
    ...(wartosci.opis === '' ? {} : { description: wartosci.opis }),
    extensionIds: [wartosci.konektor ?? ''],
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu zestawu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Zestaw „${wartosci.nazwa ?? ''}" zapisany.`);
  otoczenie.odswiez();
}

async function zastosujZestaw(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wykaz = await wywolaj(otoczenie.kanal, Command.ExtensionCollectionList, {});
  const zestawy = wykaz.wynik?.collections ?? [];
  if (zestawy.length === 0) {
    oglos(NAGLOWEK, 'Nie ma jeszcze żadnego zestawu.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWSzufladzie(otoczenie.korzen, panel, {
    tytul: 'Zastosowanie zestawu',
    pola: [
      {
        klucz: 'zestaw',
        etykieta: 'Zestaw',
        wybor: zestawy.map((zestaw) => [zestaw.id, zestaw.name] as const),
      },
      {
        klucz: 'stan',
        etykieta: 'Konektory zestawu',
        wybor: [['tak', 'Włącz'], ['nie', 'Wyłącz']],
      },
    ],
    wykonanie: 'Zastosuj zestaw',
    nieodwracalne: 'Stan wszystkich konektorów zestawu zmieni się naraz.',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.ExtensionCollectionApply, {
    collectionId: wartosci.zestaw ?? '',
    enable: wartosci.stan === 'tak',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zastosowania zestawu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Zestaw zastosowany.');
  otoczenie.odswiez();
}

/* Praca zbiorcza obejmuje jeden konektor naraz, bo okno nie prowadzi
   zaznaczania wielu wierszy — komenda przyjmuje wykaz, okno podaje jeden. */
async function przestawWieleNaraz(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazKonektor(otoczenie, panel, 'Przestawienie konektora',
    'Przestaw konektor', [{
      klucz: 'czynnosc',
      etykieta: 'Czynność',
      wybor: [
        [ExtensionBulkAction.Enable, 'Włącz'],
        [ExtensionBulkAction.Disable, 'Wyłącz'],
        [ExtensionBulkAction.Uninstall, 'Odinstaluj'],
      ],
    }],
    'Odinstalowanie zdejmuje konektor wraz z jego ustawieniami.');
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.ExtensionAdminBulk, {
    extensionIds: [wartosci.konektor ?? ''],
    action: (wartosci.czynnosc ?? ExtensionBulkAction.Disable) as ExtensionBulkAction,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił przestawienia konektora.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Konektor przestawiony.');
  otoczenie.odswiez();
}

async function ustawLacze(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazKonektor(otoczenie, panel, 'Łącze konektora', 'Ustaw łącze', [
    {
      klucz: 'rodzaj',
      etykieta: 'Rodzaj łącza',
      wybor: [
        [McpTransport.Stdio, 'Proces na maszynie rdzenia'],
        [McpTransport.Http, 'Zapytania sieciowe'],
        [McpTransport.Sse, 'Strumień zdarzeń'],
        [McpTransport.StreamableHttp, 'Strumień po sieci'],
      ],
    },
    { klucz: 'adres', etykieta: 'Adres usługi', podpowiedz: 'Dla łączy sieciowych' },
    { klucz: 'polecenie', etykieta: 'Polecenie', podpowiedz: 'Dla procesu na maszynie rdzenia' },
    {
      klucz: 'sonda',
      etykieta: 'Po zapisie',
      wybor: [['tak', 'Sprawdź łącze'], ['nie', 'Bez sprawdzenia']],
    },
  ]);
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.ExtensionTransportSet, {
    extensionId: wartosci.konektor ?? '',
    transport: (wartosci.rodzaj ?? McpTransport.Stdio) as McpTransport,
    ...(wartosci.adres === '' ? {} : { endpoint: wartosci.adres }),
    ...(wartosci.polecenie === '' ? {} : { command: wartosci.polecenie }),
    probe: wartosci.sonda !== 'nie',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił ustawienia łącza.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Łącze konektora ustawione.');
}

async function wykazWejsc(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazKonektor(otoczenie, panel, 'Wejścia sieciowe konektora',
    'Odczytaj wejścia');
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.ExtensionWebhookList, {
    extensionId: wartosci.konektor ?? '',
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu wejść.', 'ostrzezenie');
    return;
  }
  wypelnij(otoczenie.korzen, 'panel-artefakty', wynik.wynik.webhooks.map((wejscie) =>
    `${wejscie.direction} · ${wejscie.url ?? 'bez adresu'}`
    + `${wejscie.enabled ? '' : ' · wyłączone'}`));
  if (wynik.wynik.webhooks.length === 0) oglos(NAGLOWEK, 'Ten konektor nie ma wejść sieciowych.');
}

async function zapiszWejscie(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazKonektor(otoczenie, panel, 'Wejście sieciowe', 'Zapisz wejście', [
    {
      klucz: 'kierunek',
      etykieta: 'Kierunek',
      wybor: [
        [ExtensionWebhookDirection.Inbound, 'Do nas'],
        [ExtensionWebhookDirection.Outbound, 'Od nas'],
      ],
    },
    { klucz: 'adres', etykieta: 'Adres', podpowiedz: 'Wymagany dla kierunku od nas' },
    { klucz: 'zdarzenia', etykieta: 'Rodzaje zdarzeń', podpowiedz: 'Rozdzielone spacją' },
    { klucz: 'sejf', etykieta: 'Odwołanie do sejfu z podpisem' },
  ]);
  if (wartosci === null) return;
  const zdarzenia = (wartosci.zdarzenia ?? '').split(/\s+/).filter((rodzaj) => rodzaj !== '');
  const wynik = await wywolaj(otoczenie.kanal, Command.ExtensionWebhookSave, {
    extensionId: wartosci.konektor ?? '',
    direction: (wartosci.kierunek
      ?? ExtensionWebhookDirection.Inbound) as ExtensionWebhookDirection,
    ...(wartosci.adres === '' ? {} : { url: wartosci.adres }),
    ...(zdarzenia.length === 0 ? {} : { eventTypes: zdarzenia }),
    ...(wartosci.sejf === '' ? {} : { secretRef: wartosci.sejf }),
    enabled: true,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu wejścia.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Wejście sieciowe zapisane.');
}
