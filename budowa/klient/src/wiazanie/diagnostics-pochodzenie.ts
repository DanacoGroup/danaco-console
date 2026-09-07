// Panel artefaktów w oknie Diagnostics: wywołania modeli — ich wykaz, wgląd,
// powtórzenie, ocena i wydanie śladu, obok rachunku zużycia, biegu procesów
// i pytania do doradcy.
import {
  Command,
  ModelCallQuality,
  TelemetryFormat,
  UsageDimension,
} from '../../../shared/contract.ts';
import type { ModelCallTrace, UsageAggregate } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { dolozCzynnosciPanelu, zapytajWSzufladzie } from './czynnosci-okna.ts';
import { oglos } from './ogloszenie.ts';

const NAGLOWEK = 'Wywołania modeli';
const PANEL = 'panel-artefakty';
const DOBA = 86400000;

const OCENY: ReadonlyArray<readonly [string, string]> = [
  [ModelCallQuality.Accurate, 'Trafne'],
  [ModelCallQuality.Partial, 'Częściowo trafne'],
  [ModelCallQuality.Inaccurate, 'Nietrafne'],
  [ModelCallQuality.Unrated, 'Bez oceny'],
];

const POSTACIE: ReadonlyArray<readonly [string, string]> = [
  [TelemetryFormat.Json, 'JSON'],
  [TelemetryFormat.Jsonl, 'JSON po wierszu'],
  [TelemetryFormat.Csv, 'CSV'],
  [TelemetryFormat.Otlp, 'OTLP'],
];

const WYMIARY: ReadonlyArray<readonly [string, string]> = [
  [UsageDimension.Channel, 'Kanał'],
  [UsageDimension.Provider, 'Dostawca'],
  [UsageDimension.Account, 'Konto'],
  [UsageDimension.Session, 'Sesja'],
  [UsageDimension.Project, 'Projekt'],
  [UsageDimension.Environment, 'Środowisko'],
  [UsageDimension.Window, 'Okno'],
];

interface Otoczenie {
  kanal: Kanal;
  korzen: Element;
  idOkna: () => string;
}

export function zwiazPochodzenieDiagnostyki(
  kanal: Kanal,
  korzen: Element,
  idOkna: () => string,
  dolacz: (zdejmij: () => void) => void,
  przy: AddEventListenerOptions,
): void {
  const otoczenie: Otoczenie = { kanal, korzen, idOkna };
  const zdejmij = dolozCzynnosciPanelu(korzen, PANEL, 'Wywołania i rachunek', [
    {
      naglowek: 'Wywołania modeli',
      pozycje: [
        { kod: 'wykaz', nazwa: 'Wykaz wywołań' },
        { kod: 'wglad', nazwa: 'Wgląd w wywołanie…' },
        { kod: 'powtorz', nazwa: 'Powtórz wywołanie…' },
        { kod: 'ocen', nazwa: 'Oceń wywołanie…' },
        { kod: 'wydaj', nazwa: 'Wydaj ślad wywołań…' },
      ],
    },
    {
      naglowek: 'Rachunek zużycia',
      pozycje: [
        { kod: 'zuzycie', nazwa: 'Rachunek zużycia…' },
        { kod: 'zestawienie', nazwa: 'Zbuduj zestawienie…' },
      ],
    },
    {
      naglowek: 'Bieg i doradca',
      pozycje: [
        { kod: 'bieg', nazwa: 'Stan biegu procesów' },
        { kod: 'obserwuj', nazwa: 'Obserwuj bieg okna' },
        { kod: 'doradca', nazwa: 'Zapytaj doradcę…' },
      ],
    },
  ], (kod) => {
    void wykonaj(otoczenie, kod);
  }, przy);
  if (zdejmij !== null) dolacz(zdejmij);
}

function wypelnij(korzen: Element, wiersze: string[]): void {
  const cialo = korzen.querySelector(`#${PANEL} .sta-okno-tresc`);
  if (cialo === null) return;
  cialo.replaceChildren(...wiersze.map((tresc) => {
    const wiersz = cialo.ownerDocument.createElement('div');
    wiersz.className = 'dn-wykaz-modulu-poz';
    wiersz.textContent = tresc;
    return wiersz;
  }));
}

async function wykonaj(otoczenie: Otoczenie, kod: string): Promise<void> {
  if (kod === 'wykaz') return wykazWywolan(otoczenie);
  if (kod === 'wglad') return wgladWWywolanie(otoczenie);
  if (kod === 'powtorz') return powtorzWywolanie(otoczenie);
  if (kod === 'ocen') return ocenWywolanie(otoczenie);
  if (kod === 'wydaj') return wydajSlad(otoczenie);
  if (kod === 'zuzycie') return rachunekZuzycia(otoczenie);
  if (kod === 'zestawienie') return zbudujZestawienie(otoczenie);
  if (kod === 'bieg') return stanBiegu(otoczenie);
  if (kod === 'obserwuj') return obserwujBieg(otoczenie);
  if (kod === 'doradca') return zapytajDoradce(otoczenie);
  /* Czynność spoza obsłużonych odmawia zamiast milczeć. */
  oglos(NAGLOWEK, `Czynność „${kod}" nie jest prowadzona przez to okno.`, 'ostrzezenie');
}

function opiszWywolanie(wywolanie: ModelCallTrace): string {
  const model = wywolanie.model ?? wywolanie.provider ?? 'bez wskazania modelu';
  const czas = wywolanie.latencyMs === undefined ? '' : ` · ${String(wywolanie.latencyMs)} ms`;
  const zetony = wywolanie.totalTokens === undefined
    ? '' : ` · ${String(wywolanie.totalTokens)} żetonów`;
  return `${model} · ${wywolanie.status}${czas}${zetony}`;
}

async function wywolania(kanal: Kanal): Promise<ModelCallTrace[]> {
  const wynik = await wywolaj(kanal, Command.ProvenanceCallList, { limit: 50 });
  return wynik.wynik?.calls ?? [];
}

async function wyborWywolan(kanal: Kanal): Promise<ReadonlyArray<readonly [string, string]>> {
  return (await wywolania(kanal)).map((wywolanie) =>
    [wywolanie.id, opiszWywolanie(wywolanie)] as const);
}

async function wykazWywolan(otoczenie: Otoczenie): Promise<void> {
  const wynik = await wywolaj(otoczenie.kanal, Command.ProvenanceCallList, { limit: 50 });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wykazu wywołań.', 'ostrzezenie');
    return;
  }
  const wykaz = wynik.wynik.calls;
  wypelnij(otoczenie.korzen, wykaz.length === 0
    ? ['Rdzeń nie zapisał jeszcze żadnego wywołania modelu.']
    : wykaz.map(opiszWywolanie));
  oglos(NAGLOWEK, `Wywołań w wykazie: ${String(wykaz.length)}.`);
}

/* Wskazanie wywołania idzie przez wybór z wykazu, bo identyfikatora wywołania
   Operator nigdzie nie widzi. */
async function wskazWywolanie(
  otoczenie: Otoczenie,
  tytul: string,
  wykonanie: string,
  dodatkowe: Parameters<typeof zapytajWSzufladzie>[2]['pola'] = [],
): Promise<Record<string, string> | null> {
  const wybor = await wyborWywolan(otoczenie.kanal);
  if (wybor.length === 0) {
    oglos(NAGLOWEK, 'Rdzeń nie zapisał jeszcze żadnego wywołania modelu.', 'ostrzezenie');
    return null;
  }
  return zapytajWSzufladzie(otoczenie.korzen, PANEL, {
    tytul,
    pola: [{ klucz: 'wywolanie', etykieta: 'Wywołanie', wybor }, ...dodatkowe],
    wykonanie,
  });
}

async function wgladWWywolanie(otoczenie: Otoczenie): Promise<void> {
  const wartosci = await wskazWywolanie(otoczenie, 'Wgląd w wywołanie', 'Odczytaj');
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.ProvenanceCallGet, {
    callId: wartosci.wywolanie ?? '',
    includeContent: true,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wglądu w wywołanie.', 'ostrzezenie');
    return;
  }
  const wywolanie = wynik.wynik.call;
  wypelnij(otoczenie.korzen, [
    opiszWywolanie(wywolanie),
    `kanał: ${wywolanie.channelId ?? 'bez wskazania'}`,
    `okno: ${wywolanie.windowId ?? 'bez wskazania'}`,
    `koszt: ${wywolanie.cost === undefined ? 'bez wyceny' : String(wywolanie.cost)}`,
    `ocena: ${wywolanie.quality ?? 'bez oceny'}`,
  ]);
  oglos(NAGLOWEK, 'Wywołanie odczytane.');
}

async function powtorzWywolanie(otoczenie: Otoczenie): Promise<void> {
  const wartosci = await wskazWywolanie(otoczenie, 'Powtórzenie wywołania', 'Powtórz', [
    { klucz: 'model', etykieta: 'Model zastępczy' },
  ]);
  if (wartosci === null) return;
  oglos(NAGLOWEK, 'Zlecam powtórzenie wywołania…');
  const wynik = await wywolaj(otoczenie.kanal, Command.ProvenanceCallReplay, {
    callId: wartosci.wywolanie ?? '',
    ...(wartosci.model === '' ? {} : { model: wartosci.model }),
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił powtórzenia wywołania.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Wywołanie powtórzone.');
}

async function ocenWywolanie(otoczenie: Otoczenie): Promise<void> {
  const wartosci = await wskazWywolanie(otoczenie, 'Ocena wywołania', 'Zapisz ocenę', [
    { klucz: 'ocena', etykieta: 'Ocena', wybor: OCENY },
    { klucz: 'uwaga', etykieta: 'Uwaga', obszerne: true },
  ]);
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.ProvenanceCallRate, {
    callId: wartosci.wywolanie ?? '',
    quality: (wartosci.ocena ?? ModelCallQuality.Unrated) as ModelCallQuality,
    ...(wartosci.uwaga === '' ? {} : { note: wartosci.uwaga }),
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu oceny.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Ocena wywołania zapisana.');
}

/* Ślad wychodzi do schowka: okno nie prowadzi powłoki wskazania katalogu. */
async function wydajSlad(otoczenie: Otoczenie): Promise<void> {
  const wartosci = await zapytajWSzufladzie(otoczenie.korzen, PANEL, {
    tytul: 'Wydanie śladu wywołań',
    opis: 'Ślad idzie do schowka Operatora.',
    pola: [
      { klucz: 'postac', etykieta: 'Postać zapisu', wybor: POSTACIE },
      { klucz: 'dni', etykieta: 'Ile dni wstecz', wartosc: '7' },
    ],
    wykonanie: 'Wydaj ślad',
  });
  if (wartosci === null) return;
  const dni = Number(wartosci.dni ?? '7');
  const wynik = await wywolaj(otoczenie.kanal, Command.ProvenanceTraceExport, {
    format: (wartosci.postac ?? TelemetryFormat.Json) as TelemetryFormat,
    fromTime: Date.now() - (Number.isFinite(dni) ? dni : 7) * DOBA,
    toTime: Date.now(),
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wydania śladu.', 'ostrzezenie');
    return;
  }
  try {
    await navigator.clipboard.writeText(wynik.wynik.content);
    oglos(NAGLOWEK, `Ślad ${String(wynik.wynik.callCount)} wywołań jest w schowku.`);
  } catch {
    oglos(NAGLOWEK, 'Przeglądarka nie dała dostępu do schowka, więc ślad nie został przeniesiony.',
      'ostrzezenie');
  }
}

function opiszRachunek(pozycja: UsageAggregate): string {
  const koszt = pozycja.cost === undefined
    ? 'bez wyceny'
    : `${String(pozycja.cost)} ${pozycja.currency ?? ''}`.trim();
  return `${pozycja.dimensionLabel ?? pozycja.dimensionId} · ${String(pozycja.requests)} wywołań`
    + ` · ${String(pozycja.totalTokens)} żetonów · ${koszt}`;
}

async function rachunekZuzycia(otoczenie: Otoczenie): Promise<void> {
  const wartosci = await zapytajWSzufladzie(otoczenie.korzen, PANEL, {
    tytul: 'Rachunek zużycia',
    pola: [
      { klucz: 'wymiar', etykieta: 'W podziale na', wybor: WYMIARY },
      { klucz: 'dni', etykieta: 'Ile dni wstecz', wartosc: '30' },
    ],
    wykonanie: 'Policz',
  });
  if (wartosci === null) return;
  const dni = Number(wartosci.dni ?? '30');
  const wynik = await wywolaj(otoczenie.kanal, Command.UsageSummaryGet, {
    dimension: (wartosci.wymiar ?? UsageDimension.Channel) as UsageDimension,
    fromTime: Date.now() - (Number.isFinite(dni) ? dni : 30) * DOBA,
    toTime: Date.now(),
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił rachunku zużycia.', 'ostrzezenie');
    return;
  }
  const pozycje = wynik.wynik.aggregates;
  wypelnij(otoczenie.korzen, pozycje.length === 0
    ? ['W tym okresie rdzeń nie zapisał żadnego zużycia.']
    : pozycje.map(opiszRachunek));
  oglos(NAGLOWEK, `Pozycji rachunku: ${String(pozycje.length)}.`);
}

async function zbudujZestawienie(otoczenie: Otoczenie): Promise<void> {
  const wartosci = await zapytajWSzufladzie(otoczenie.korzen, PANEL, {
    tytul: 'Zestawienie zużycia',
    opis: 'Zestawienie idzie do schowka Operatora.',
    pola: [
      { klucz: 'postac', etykieta: 'Postać zapisu', wybor: POSTACIE },
      { klucz: 'dni', etykieta: 'Ile dni wstecz', wartosc: '30' },
    ],
    wykonanie: 'Zbuduj',
  });
  if (wartosci === null) return;
  const dni = Number(wartosci.dni ?? '30');
  const wynik = await wywolaj(otoczenie.kanal, Command.UsageReportBuild, {
    format: (wartosci.postac ?? TelemetryFormat.Csv) as TelemetryFormat,
    fromTime: Date.now() - (Number.isFinite(dni) ? dni : 30) * DOBA,
    toTime: Date.now(),
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zbudowania zestawienia.', 'ostrzezenie');
    return;
  }
  try {
    await navigator.clipboard.writeText(wynik.wynik.content);
    oglos(NAGLOWEK, `Zestawienie ${String(wynik.wynik.requests)} wywołań jest w schowku.`);
  } catch {
    oglos(NAGLOWEK, 'Przeglądarka nie dała dostępu do schowka.', 'ostrzezenie');
  }
}

async function stanBiegu(otoczenie: Otoczenie): Promise<void> {
  const wynik = await wywolaj(otoczenie.kanal, Command.MonitorStatus, {});
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił stanu biegu.', 'ostrzezenie');
    return;
  }
  const stany = wynik.wynik.statuses;
  wypelnij(otoczenie.korzen, stany.length === 0
    ? ['Żaden proces nie biegnie.']
    : stany.map((stan) => `${stan.label ?? stan.processId} · ${stan.status}`
      + (stan.stage === undefined ? '' : ` · ${stan.stage}`)));
  oglos(NAGLOWEK, `Procesów w biegu: ${String(stany.length)}.`);
}

/* Obserwacja biegu zapisuje okno u rdzenia: od tej chwili postępy przychodzą
   same, bez pytania. */
async function obserwujBieg(otoczenie: Otoczenie): Promise<void> {
  const wynik = await wywolaj(otoczenie.kanal, Command.MonitorSubscribe, {
    ...(otoczenie.idOkna() === '' ? {} : { windowId: otoczenie.idOkna() }),
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił obserwacji biegu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Rdzeń będzie odtąd donosił o biegu tego okna.');
}

async function zapytajDoradce(otoczenie: Otoczenie): Promise<void> {
  if (otoczenie.idOkna() === '') {
    oglos(NAGLOWEK, 'To okno nie stoi jeszcze przy sesji.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWSzufladzie(otoczenie.korzen, PANEL, {
    tytul: 'Pytanie do doradcy',
    pola: [
      { klucz: 'pytanie', etykieta: 'Pytanie', obszerne: true, wymagane: true },
      { klucz: 'tlo', etykieta: 'Tło pytania', obszerne: true },
    ],
    wykonanie: 'Zapytaj',
  });
  if (wartosci === null) return;
  oglos(NAGLOWEK, 'Pytam doradcę…');
  const wynik = await wywolaj(otoczenie.kanal, Command.AdvisorConsult, {
    windowId: otoczenie.idOkna(),
    question: wartosci.pytanie ?? '',
    ...(wartosci.tlo === '' ? {} : { context: wartosci.tlo }),
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił pytania do doradcy.', 'ostrzezenie');
    return;
  }
  wypelnij(otoczenie.korzen, [
    `doradca: ${wynik.wynik.advisorChannel}`,
    wynik.wynik.advice,
  ]);
  oglos(NAGLOWEK, 'Doradca odpowiedział.');
}
