// Czuwanie w oknie Diagnostics: sondy zdrowia w panelu zadań oraz czujki
// progowe wraz z ich zapłonami w centrum diagnostyki.
import {
  AlertChannel,
  AlertComparison,
  AlertMetric,
  AlertRuleKind,
  AlertTriggerStatus,
  Command,
  DiagnosticPriority,
  HealthProbeKind,
} from '../../../shared/contract.ts';
import type { AlertRule, AlertTrigger, HealthProbe } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { dolozCzynnosciPanelu, zapytajWSzufladzie } from './czynnosci-okna.ts';
import { oglos } from './ogloszenie.ts';

const NAGLOWEK = 'Czuwanie';
const GODZINA_MS = 3600000;

const MIARY: ReadonlyArray<readonly [string, string]> = [
  [AlertMetric.ErrorCount, 'Liczba błędów'],
  [AlertMetric.ErrorRate, 'Odsetek błędów'],
  [AlertMetric.Cost, 'Koszt'],
  [AlertMetric.Tokens, 'Zużycie żetonów'],
  [AlertMetric.BudgetPercent, 'Odsetek budżetu'],
  [AlertMetric.CallLatency, 'Czas odpowiedzi'],
  [AlertMetric.ProbeFailure, 'Niepowodzenie sondy'],
  [AlertMetric.ProcessFailure, 'Niepowodzenie procesu'],
];

const WAGI: ReadonlyArray<readonly [string, string]> = [
  [DiagnosticPriority.Low, 'Niska'],
  [DiagnosticPriority.Medium, 'Średnia'],
  [DiagnosticPriority.High, 'Wysoka'],
  [DiagnosticPriority.Critical, 'Krytyczna'],
];

interface Otoczenie {
  kanal: Kanal;
  korzen: Element;
}

export function zwiazCzuwanieDiagnostyki(
  kanal: Kanal,
  korzen: Element,
  dolacz: (zdejmij: () => void) => void,
  przy: AddEventListenerOptions,
): void {
  const otoczenie: Otoczenie = { kanal, korzen };
  const zdejmowanie = [
    dolozCzynnosciPanelu(korzen, 'panel-zadania', 'Czynności sond', [
      {
        naglowek: 'Sondy zdrowia',
        pozycje: [
          { kod: 'sondy', nazwa: 'Wykaz sond' },
          { kod: 'sonda-zapisz', nazwa: 'Zapisz sondę…' },
          { kod: 'sonda-uruchom', nazwa: 'Uruchom sondę…' },
          { kod: 'sonda-usun', nazwa: 'Usuń sondę…' },
        ],
      },
      {
        naglowek: 'Wyniki',
        pozycje: [
          { kod: 'wyniki', nazwa: 'Wyniki sond…' },
          { kod: 'dostepnosc', nazwa: 'Dostępność sond…' },
        ],
      },
    ], (kod) => {
      void wykonaj(otoczenie, kod, 'panel-zadania');
    }, przy),
    dolozCzynnosciPanelu(korzen, 'panel-centrum', 'Czynności czujek', [
      {
        naglowek: 'Czujki progowe',
        pozycje: [
          { kod: 'czujki', nazwa: 'Wykaz czujek' },
          { kod: 'czujka-zapisz', nazwa: 'Zapisz czujkę…' },
          { kod: 'czujka-usun', nazwa: 'Usuń czujkę…' },
        ],
      },
      {
        naglowek: 'Zapłony',
        pozycje: [
          { kod: 'zaplony', nazwa: 'Wykaz zapłonów…' },
          { kod: 'zaplon-przyjmij', nazwa: 'Przyjmij zapłon…' },
        ],
      },
    ], (kod) => {
      void wykonaj(otoczenie, kod, 'panel-centrum');
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

async function wykonaj(otoczenie: Otoczenie, kod: string, panel: string): Promise<void> {
  if (kod === 'sondy') return wykazSond(otoczenie);
  if (kod === 'sonda-zapisz') return zapiszSonde(otoczenie, panel);
  if (kod === 'sonda-uruchom') return uruchomSonde(otoczenie, panel);
  if (kod === 'sonda-usun') return usunSonde(otoczenie, panel);
  if (kod === 'wyniki') return wynikiSond(otoczenie, panel);
  if (kod === 'dostepnosc') return dostepnoscSond(otoczenie, panel);
  if (kod === 'czujki') return wykazCzujek(otoczenie);
  if (kod === 'czujka-zapisz') return zapiszCzujke(otoczenie, panel);
  if (kod === 'czujka-usun') return usunCzujke(otoczenie, panel);
  if (kod === 'zaplony') return wykazZaplonow(otoczenie, panel);
  if (kod === 'zaplon-przyjmij') return przyjmijZaplon(otoczenie, panel);
  /* Czynność spoza obsłużonych odmawia zamiast milczeć: pozycja bez gałęzi
     wyglądałaby jak działająca. */
  oglos(NAGLOWEK, `Czynność „${kod}" nie jest prowadzona przez to okno.`, 'ostrzezenie');
}

async function wyborSond(kanal: Kanal): Promise<ReadonlyArray<readonly [string, string]>> {
  const wykaz = await wywolaj(kanal, Command.HealthProbeList, {});
  return (wykaz.wynik?.probes ?? []).map((sonda: HealthProbe) =>
    [sonda.id, `${sonda.name} · ${sonda.lastStatus ?? 'bez biegu'}`] as const);
}

async function wykazSond({ kanal, korzen }: Otoczenie): Promise<void> {
  const wynik = await wywolaj(kanal, Command.HealthProbeList, {});
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu sond.', 'ostrzezenie');
    return;
  }
  const sondy = wynik.wynik.probes;
  wypelnij(korzen, 'panel-zadania', sondy.map((sonda: HealthProbe) =>
    `${sonda.lastStatus ?? 'bez biegu'} · ${sonda.name} · ${sonda.kind} ${sonda.target}`));
  if (sondy.length === 0) oglos(NAGLOWEK, 'Rdzeń nie ma jeszcze żadnej sondy.');
}

/* Odstęp podawany jest w sekundach, bo w takich jednostkach Operator myśli o
   sondowaniu; rdzeń liczy w milisekundach i okno przelicza to za niego. */
async function zapiszSonde({ kanal, korzen }: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Sonda zdrowia',
    pola: [
      { klucz: 'nazwa', etykieta: 'Nazwa sondy', wymagane: true },
      {
        klucz: 'rodzaj',
        etykieta: 'Rodzaj sondy',
        wybor: [
          [HealthProbeKind.Http, 'Zapytanie sieciowe'],
          [HealthProbeKind.Tcp, 'Połączenie do portu'],
          [HealthProbeKind.Command, 'Polecenie powłoki'],
          [HealthProbeKind.ModelCall, 'Wywołanie modelu'],
          [HealthProbeKind.Internal, 'Sprawdzenie wewnętrzne'],
        ],
      },
      { klucz: 'cel', etykieta: 'Cel sondy', podpowiedz: 'http://127.0.0.1:8080/zdrowie',
        wymagane: true },
      { klucz: 'odstep', etykieta: 'Odstęp w sekundach', wartosc: '60', wymagane: true },
      { klucz: 'granica', etykieta: 'Granica czasu w sekundach', wartosc: '10' },
    ],
    wykonanie: 'Zapisz sondę',
  });
  if (wartosci === null) return;
  const odstep = Number.parseInt(wartosci.odstep ?? '', 10);
  const granica = Number.parseInt(wartosci.granica ?? '', 10);
  if (!Number.isFinite(odstep) || odstep < 1) {
    oglos(NAGLOWEK, 'Odstęp musi być liczbą sekund nie mniejszą niż jeden.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.HealthProbeSave, {
    name: wartosci.nazwa ?? '',
    kind: (wartosci.rodzaj ?? HealthProbeKind.Http) as HealthProbeKind,
    target: wartosci.cel ?? '',
    intervalMs: odstep * 1000,
    ...(Number.isFinite(granica) ? { timeoutMs: granica * 1000 } : {}),
    enabled: true,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu sondy.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Sonda „${wartosci.nazwa ?? ''}" zapisana.`);
  await wykazSond({ kanal, korzen });
}

async function uruchomSonde({ kanal, korzen }: Otoczenie, panel: string): Promise<void> {
  const sondy = await wyborSond(kanal);
  if (sondy.length === 0) {
    oglos(NAGLOWEK, 'Rdzeń nie ma jeszcze żadnej sondy.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Bieg sondy',
    pola: [{ klucz: 'sonda', etykieta: 'Sonda', wybor: sondy }],
    wykonanie: 'Uruchom sondę',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.HealthProbeRun, { probeId: wartosci.sonda ?? '' });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił biegu sondy.', 'ostrzezenie');
    return;
  }
  const bieg = wynik.wynik.result;
  oglos(NAGLOWEK, `Sonda „${wynik.wynik.probe.name}": ${bieg.status}`
    + `${bieg.latencyMs === undefined ? '' : `, ${bieg.latencyMs} ms`}`
    + `${bieg.detail === undefined ? '' : ` · ${bieg.detail}`}`);
  await wykazSond({ kanal, korzen });
}

async function usunSonde({ kanal, korzen }: Otoczenie, panel: string): Promise<void> {
  const sondy = await wyborSond(kanal);
  if (sondy.length === 0) {
    oglos(NAGLOWEK, 'Rdzeń nie ma jeszcze żadnej sondy.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Usunięcie sondy',
    pola: [{ klucz: 'sonda', etykieta: 'Sonda', wybor: sondy }],
    wykonanie: 'Usuń sondę',
    nieodwracalne: 'Sonda znika wraz ze swoimi wynikami i wyliczeniem dostępności.',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.HealthProbeRemove, { probeId: wartosci.sonda ?? '' });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił usunięcia sondy.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Sonda zdjęta z wykazu.');
  await wykazSond({ kanal, korzen });
}

/* Zakres liczony jest wstecz od chwili odczytu, bo okno nie prowadzi wyboru
   dat, a Operator pyta o to, co działo się przez ostatnie godziny. */
async function wynikiSond({ kanal, korzen }: Otoczenie, panel: string): Promise<void> {
  const sondy = await wyborSond(kanal);
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Wyniki sond',
    pola: [
      { klucz: 'sonda', etykieta: 'Sonda', wybor: [['', 'Wszystkie sondy'], ...sondy] },
      { klucz: 'godziny', etykieta: 'Ostatnich godzin', wartosc: '24' },
    ],
    wykonanie: 'Odczytaj wyniki',
  });
  if (wartosci === null) return;
  const godziny = Number.parseInt(wartosci.godziny ?? '', 10);
  const wynik = await wywolaj(kanal, Command.HealthResultList, {
    ...(wartosci.sonda === '' ? {} : { probeId: wartosci.sonda }),
    fromTime: Date.now() - (Number.isFinite(godziny) ? godziny : 24) * GODZINA_MS,
    limit: 200,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu wyników.', 'ostrzezenie');
    return;
  }
  const wyniki = wynik.wynik.results;
  wypelnij(korzen, 'panel-zadania', wyniki.map((bieg) =>
    `${bieg.status}${bieg.latencyMs === undefined ? '' : ` · ${bieg.latencyMs} ms`}`
    + `${bieg.detail === undefined ? '' : ` · ${bieg.detail}`}`));
  if (wyniki.length === 0) {
    oglos(NAGLOWEK, 'Żadna sonda nie biegła w tym zakresie.');
    return;
  }
  if (wynik.wynik.truncated === true) {
    oglos(NAGLOWEK, 'Wykaz przycięty granicą odczytu — wyników jest więcej.', 'ostrzezenie');
  }
}

async function dostepnoscSond({ kanal, korzen }: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Dostępność sond',
    pola: [{ klucz: 'godziny', etykieta: 'Ostatnich godzin', wartosc: '24' }],
    wykonanie: 'Wylicz dostępność',
  });
  if (wartosci === null) return;
  const godziny = Number.parseInt(wartosci.godziny ?? '', 10);
  const wynik = await wywolaj(kanal, Command.HealthUptimeGet, {
    fromTime: Date.now() - (Number.isFinite(godziny) ? godziny : 24) * GODZINA_MS,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wyliczenia dostępności.', 'ostrzezenie');
    return;
  }
  const wyliczenia = wynik.wynik.uptimes;
  wypelnij(korzen, 'panel-zadania', wyliczenia.map((wyliczenie) =>
    `${wyliczenie.probeName ?? wyliczenie.probeId} · `
    + `${wyliczenie.uptimePercent ?? 'brak pomiaru'}% z ${wyliczenie.samples} pomiarów`));
  if (wyliczenia.length === 0) oglos(NAGLOWEK, 'Żadna sonda nie ma pomiarów w tym zakresie.');
}

async function wykazCzujek({ kanal, korzen }: Otoczenie): Promise<void> {
  const wynik = await wywolaj(kanal, Command.AlertRuleList, {});
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu czujek.', 'ostrzezenie');
    return;
  }
  const czujki = wynik.wynik.rules;
  wypelnij(korzen, 'panel-centrum', czujki.map((czujka: AlertRule) =>
    `${czujka.enabled ? 'czynna' : 'wyłączona'} · ${czujka.name} · ${czujka.metric}`
    + `${czujka.threshold === undefined ? '' : ` ${czujka.comparison ?? ''} ${czujka.threshold}`}`));
  if (czujki.length === 0) oglos(NAGLOWEK, 'Rdzeń nie ma jeszcze żadnej czujki.');
}

/* Okno prowadzi czujkę progową: rozpoznanie odchyleń i nowych odcisków rdzeń
   liczy sam, więc nie potrzebuje od Operatora progu ani porównania. */
async function zapiszCzujke({ kanal, korzen }: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Czujka progowa',
    pola: [
      { klucz: 'nazwa', etykieta: 'Nazwa czujki', wymagane: true },
      { klucz: 'miara', etykieta: 'Miara', wybor: MIARY },
      {
        klucz: 'porownanie',
        etykieta: 'Zapłon następuje, gdy miara jest',
        wybor: [
          [AlertComparison.GreaterThan, 'większa niż próg'],
          [AlertComparison.GreaterOrEqual, 'nie mniejsza niż próg'],
          [AlertComparison.LessThan, 'mniejsza niż próg'],
          [AlertComparison.LessOrEqual, 'nie większa niż próg'],
        ],
      },
      { klucz: 'prog', etykieta: 'Próg', wymagane: true },
      { klucz: 'okno', etykieta: 'Okno pomiaru w minutach', wartosc: '15', wymagane: true },
      { klucz: 'waga', etykieta: 'Waga zapłonu', wybor: WAGI },
      {
        klucz: 'kanal',
        etykieta: 'Kanał powiadomienia',
        wybor: [
          [AlertChannel.App, 'W aplikacji'],
          [AlertChannel.AlwaysOnDisplay, 'Na awatarze'],
          [AlertChannel.Mail, 'Pocztą'],
        ],
      },
    ],
    wykonanie: 'Zapisz czujkę',
  });
  if (wartosci === null) return;
  const prog = Number.parseFloat(wartosci.prog ?? '');
  const okno = Number.parseInt(wartosci.okno ?? '', 10);
  if (!Number.isFinite(prog) || !Number.isFinite(okno) || okno < 1) {
    oglos(NAGLOWEK, 'Próg musi być liczbą, a okno pomiaru liczbą minut.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.AlertRuleSave, {
    name: wartosci.nazwa ?? '',
    kind: AlertRuleKind.Threshold,
    metric: (wartosci.miara ?? AlertMetric.ErrorCount) as AlertMetric,
    comparison: (wartosci.porownanie ?? AlertComparison.GreaterThan) as AlertComparison,
    threshold: prog,
    windowMs: okno * 60000,
    severity: (wartosci.waga ?? DiagnosticPriority.Medium) as DiagnosticPriority,
    channels: [(wartosci.kanal ?? AlertChannel.App) as AlertChannel],
    enabled: true,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu czujki.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Czujka „${wartosci.nazwa ?? ''}" czuwa.`);
  await wykazCzujek({ kanal, korzen });
}

async function usunCzujke({ kanal, korzen }: Otoczenie, panel: string): Promise<void> {
  const wykaz = await wywolaj(kanal, Command.AlertRuleList, {});
  const czujki = wykaz.wynik?.rules ?? [];
  if (czujki.length === 0) {
    oglos(NAGLOWEK, 'Rdzeń nie ma jeszcze żadnej czujki.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Usunięcie czujki',
    pola: [{
      klucz: 'czujka',
      etykieta: 'Czujka',
      wybor: czujki.map((czujka) => [czujka.id, czujka.name] as const),
    }],
    wykonanie: 'Usuń czujkę',
    nieodwracalne: 'Czujka znika wraz ze swoimi zapłonami.',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.AlertRuleRemove, { ruleId: wartosci.czujka ?? '' });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił usunięcia czujki.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Czujka zdjęta z wykazu.');
  await wykazCzujek({ kanal, korzen });
}

async function wykazZaplonow({ kanal, korzen }: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Zapłony czujek',
    pola: [{
      klucz: 'stan',
      etykieta: 'Stan zapłonu',
      wybor: [
        ['', 'Wszystkie'],
        [AlertTriggerStatus.Firing, 'Palące się'],
        [AlertTriggerStatus.Acknowledged, 'Przyjęte'],
        [AlertTriggerStatus.Resolved, 'Wygaszone'],
        [AlertTriggerStatus.Escalated, 'Podniesione wyżej'],
      ],
    }],
    wykonanie: 'Odczytaj zapłony',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.AlertTriggerList, {
    ...(wartosci.stan === '' ? {} : { status: wartosci.stan as AlertTriggerStatus }),
    limit: 200,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu zapłonów.', 'ostrzezenie');
    return;
  }
  const zaplony = wynik.wynik.triggers;
  wypelnij(korzen, 'panel-centrum', zaplony.map((zaplon: AlertTrigger) =>
    `${zaplon.status} · ${zaplon.ruleName ?? zaplon.ruleId} · ${zaplon.message}`));
  if (zaplony.length === 0) oglos(NAGLOWEK, 'Żadna czujka nie zapaliła się w tym stanie.');
}

/* Przyjęcie zapłonu nie wygasza czujki: mówi tylko, że ktoś go widział, więc
   okno nazywa to wprost i pyta o notatkę dla następnej osoby. */
async function przyjmijZaplon({ kanal, korzen }: Otoczenie, panel: string): Promise<void> {
  const wykaz = await wywolaj(kanal, Command.AlertTriggerList, {
    status: AlertTriggerStatus.Firing,
    limit: 100,
  });
  const zaplony = wykaz.wynik?.triggers ?? [];
  if (zaplony.length === 0) {
    oglos(NAGLOWEK, 'Żaden zapłon się teraz nie pali.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Przyjęcie zapłonu',
    opis: 'Przyjęcie odnotowuje, że zapłon jest widziany; czujki nie wygasza.',
    pola: [
      {
        klucz: 'zaplon',
        etykieta: 'Zapłon',
        wybor: zaplony.map((zaplon) =>
          [zaplon.id, `${zaplon.ruleName ?? zaplon.ruleId} · ${zaplon.message}`] as const),
      },
      { klucz: 'notatka', etykieta: 'Notatka', obszerne: true },
    ],
    wykonanie: 'Przyjmij zapłon',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.AlertTriggerAcknowledge, {
    triggerId: wartosci.zaplon ?? '',
    ...(wartosci.notatka === '' ? {} : { note: wartosci.notatka }),
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił przyjęcia zapłonu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Zapłon przyjęty.');
}
