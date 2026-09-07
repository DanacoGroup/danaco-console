// Orkiestracja i harmonogram w oknie Automations: zależności między krokami,
// grupy, bramy, wyrównania, sprawdzenie układu oraz okna czasu biegu.
import {
  AutomationDependencyKind,
  Command,
  OrchestrationGateRule,
  WindowRole,
} from '../../../shared/contract.ts';
import type { AutomationDependency, AutomationSchedule } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { dolozCzynnosciPanelu, zapytajWSzufladzie } from './czynnosci-okna.ts';
import { oglos } from './ogloszenie.ts';

const NAGLOWEK = 'Orkiestracja';
const MINUTA = 60000;

const RODZAJE: ReadonlyArray<readonly [string, string]> = [
  [AutomationDependencyKind.Sequential, 'Po kolei'],
  [AutomationDependencyKind.Parallel, 'Równolegle'],
  [AutomationDependencyKind.Conditional, 'Warunkowo'],
];

interface Otoczenie {
  kanal: Kanal;
  korzen: Element;
  odswiez: () => void;
}

export function zwiazOrkiestracjeAutomatyk(
  kanal: Kanal,
  korzen: Element,
  odswiez: () => void,
  dolacz: (zdejmij: () => void) => void,
  przy: AddEventListenerOptions,
): void {
  const otoczenie: Otoczenie = { kanal, korzen, odswiez };
  const zdejmij = dolozCzynnosciPanelu(korzen, 'panel-plan', 'Układ i harmonogram', [
    {
      naglowek: 'Układ kroków',
      pozycje: [
        { kod: 'zaleznosci', nazwa: 'Wykaz zależności…' },
        { kod: 'zaleznosc-zaloz', nazwa: 'Zwiąż dwa kroki…' },
        { kod: 'zaleznosc-zdejmij', nazwa: 'Rozwiąż dwa kroki…' },
        { kod: 'grupa', nazwa: 'Złóż kroki w grupę…' },
        { kod: 'brama', nazwa: 'Postaw bramę przed krokiem…' },
        { kod: 'wyrownanie', nazwa: 'Wskaż krok wyrównujący…' },
        { kod: 'sprawdz', nazwa: 'Sprawdź układ…' },
        { kod: 'wiazanie', nazwa: 'Przestaw wiązanie z orkiestracją…' },
      ],
    },
    {
      naglowek: 'Harmonogram',
      pozycje: [
        { kod: 'harmonogramy', nazwa: 'Wykaz harmonogramów' },
        { kod: 'okno-czasu', nazwa: 'Ustaw okno biegu…' },
        { kod: 'tetno', nazwa: 'Ustaw tolerancję tętna…' },
        { kod: 'nadrobienie', nazwa: 'Nadrób zaległe biegi…' },
        { kod: 'dzieje-wyzwolen', nazwa: 'Dzieje wyzwoleń' },
        { kod: 'wejscie-sieciowe', nazwa: 'Adres wejścia sieciowego…' },
      ],
    },
    {
      naglowek: 'Role okien',
      pozycje: [
        { kod: 'role', nazwa: 'Wykaz ról okien' },
        { kod: 'rola-nadaj', nazwa: 'Nadaj oknu rolę…' },
        { kod: 'rola-popraw', nazwa: 'Popraw rolę okna…' },
        { kod: 'rola-zdejmij', nazwa: 'Zdejmij rolę okna…' },
      ],
    },
  ], (kod) => {
    void wykonaj(otoczenie, kod, 'panel-plan');
  }, przy);
  if (zdejmij !== null) dolacz(zdejmij);
}

function wypelnij(korzen: Element, wiersze: string[]): void {
  const cialo = korzen.querySelector('#panel-plan .sta-okno-tresc');
  if (cialo === null) return;
  cialo.replaceChildren(...wiersze.map((tresc) => {
    const wiersz = cialo.ownerDocument.createElement('div');
    wiersz.className = 'dn-wykaz-modulu-poz';
    wiersz.textContent = tresc;
    return wiersz;
  }));
}

async function wyborAutomatyki(
  kanal: Kanal,
): Promise<ReadonlyArray<readonly [string, string]>> {
  const wykaz = await wywolaj(kanal, Command.AutomationWorkflowList, {});
  return (wykaz.wynik?.workflows ?? []).map((przeplyw) => [przeplyw.id, przeplyw.name] as const);
}

async function wskazAutomatyke(
  otoczenie: Otoczenie,
  panel: string,
  tytul: string,
  wykonanie: string,
  dodatkowe: Parameters<typeof zapytajWSzufladzie>[2]['pola'] = [],
  opis?: string,
): Promise<Record<string, string> | null> {
  const wybor = await wyborAutomatyki(otoczenie.kanal);
  if (wybor.length === 0) {
    oglos(NAGLOWEK, 'To okno nie ma jeszcze żadnej automatyki.', 'ostrzezenie');
    return null;
  }
  return zapytajWSzufladzie(otoczenie.korzen, panel, {
    tytul,
    ...(opis === undefined ? {} : { opis }),
    pola: [{ klucz: 'przeplyw', etykieta: 'Automatyka', wybor }, ...dodatkowe],
    wykonanie,
  });
}

async function wykonaj(otoczenie: Otoczenie, kod: string, panel: string): Promise<void> {
  if (kod === 'zaleznosci') return wykazZaleznosci(otoczenie, panel);
  if (kod === 'zaleznosc-zaloz') return zwiazKroki(otoczenie, panel);
  if (kod === 'zaleznosc-zdejmij') return rozwiazKroki(otoczenie, panel);
  if (kod === 'grupa') return zlozGrupe(otoczenie, panel);
  if (kod === 'brama') return postawBrame(otoczenie, panel);
  if (kod === 'wyrownanie') return wskazWyrownanie(otoczenie, panel);
  if (kod === 'sprawdz') return sprawdzUklad(otoczenie, panel);
  if (kod === 'wiazanie') return przestawWiazanie(otoczenie, panel);
  if (kod === 'harmonogramy') return wykazHarmonogramow(otoczenie);
  if (kod === 'okno-czasu') return ustawOknoCzasu(otoczenie, panel);
  if (kod === 'tetno') return ustawTetno(otoczenie, panel);
  if (kod === 'nadrobienie') return nadrobBiegi(otoczenie, panel);
  if (kod === 'dzieje-wyzwolen') return dziejeWyzwolen(otoczenie);
  if (kod === 'wejscie-sieciowe') return wejscieSieciowe(otoczenie, panel);
  if (kod === 'role') return wykazRol(otoczenie);
  if (kod === 'rola-nadaj') return nadajRole(otoczenie, panel);
  if (kod === 'rola-popraw') return poprawRole(otoczenie, panel);
  if (kod === 'rola-zdejmij') return zdejmijRole(otoczenie, panel);
  /* Czynność spoza obsłużonych odmawia zamiast milczeć: pozycja bez gałęzi
     wyglądałaby jak działająca. */
  oglos(NAGLOWEK, `Czynność „${kod}" nie jest prowadzona przez to okno.`, 'ostrzezenie');
}

async function wykazZaleznosci(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazAutomatyke(otoczenie, panel, 'Zależności kroków',
    'Odczytaj zależności');
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.OrchestrationDependencyList, {
    workflowId: wartosci.przeplyw ?? '',
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu zależności.', 'ostrzezenie');
    return;
  }
  const zaleznosci = wynik.wynik.dependencies;
  wypelnij(otoczenie.korzen, zaleznosci.map((zaleznosc: AutomationDependency) =>
    `${zaleznosc.fromStepId} → ${zaleznosc.toStepId} · ${zaleznosc.kind}`));
  if (zaleznosci.length === 0) oglos(NAGLOWEK, 'Ta automatyka nie ma jeszcze zależności.');
}

async function zwiazKroki(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazAutomatyke(otoczenie, panel, 'Zależność kroków', 'Zwiąż kroki', [
    { klucz: 'od', etykieta: 'Krok wcześniejszy', wymagane: true },
    { klucz: 'do', etykieta: 'Krok późniejszy', wymagane: true },
    { klucz: 'rodzaj', etykieta: 'Rodzaj zależności', wybor: RODZAJE },
    { klucz: 'warunek', etykieta: 'Warunek przejścia', podpowiedz: 'Tylko dla zależności warunkowej' },
  ]);
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.OrchestrationDependencySet, {
    workflowId: wartosci.przeplyw ?? '',
    dependency: {
      fromStepId: wartosci.od ?? '',
      toStepId: wartosci.do ?? '',
      kind: (wartosci.rodzaj ?? AutomationDependencyKind.Sequential) as AutomationDependencyKind,
      ...(wartosci.warunek === '' ? {} : { condition: wartosci.warunek }),
    },
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił związania kroków.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Krok ${wartosci.do ?? ''} idzie po kroku ${wartosci.od ?? ''}.`);
  otoczenie.odswiez();
}

async function rozwiazKroki(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazAutomatyke(otoczenie, panel, 'Zdjęcie zależności', 'Rozwiąż kroki', [
    { klucz: 'od', etykieta: 'Krok wcześniejszy', wymagane: true },
    { klucz: 'do', etykieta: 'Krok późniejszy', wymagane: true },
  ]);
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.OrchestrationDependencyRemove, {
    workflowId: wartosci.przeplyw ?? '',
    fromStepId: wartosci.od ?? '',
    toStepId: wartosci.do ?? '',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zdjęcia zależności.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Zależność zdjęta.');
  otoczenie.odswiez();
}

/* Kroki grupy podaje się rozdzielone spacją: okno nie prowadzi zaznaczania
   wielu kroków naraz, a oznaczenia widać w wykazie zależności. */
async function zlozGrupe(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazAutomatyke(otoczenie, panel, 'Grupa kroków', 'Złóż grupę', [
    { klucz: 'nazwa', etykieta: 'Nazwa grupy', wymagane: true },
    { klucz: 'kroki', etykieta: 'Kroki grupy', podpowiedz: 'Rozdzielone spacją', wymagane: true },
    { klucz: 'rodzaj', etykieta: 'Kroki grupy biegną', wybor: RODZAJE },
  ]);
  if (wartosci === null) return;
  const kroki = (wartosci.kroki ?? '').split(/\s+/).filter((krok) => krok !== '');
  const wynik = await wywolaj(otoczenie.kanal, Command.OrchestrationGroupSet, {
    workflowId: wartosci.przeplyw ?? '',
    name: wartosci.nazwa ?? '',
    stepIds: kroki,
    kind: (wartosci.rodzaj ?? AutomationDependencyKind.Parallel) as AutomationDependencyKind,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił złożenia grupy.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Grupa „${wartosci.nazwa ?? ''}" obejmuje ${kroki.length} kroków.`);
  otoczenie.odswiez();
}

/* Brama rozstrzyga, ile poprzedników musi się skończyć, zanim krok ruszy;
   liczba ma znaczenie wyłącznie przy regule liczbowej. */
async function postawBrame(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazAutomatyke(otoczenie, panel, 'Brama przed krokiem', 'Postaw bramę', [
    { klucz: 'krok', etykieta: 'Krok za bramą', wymagane: true },
    {
      klucz: 'regula',
      etykieta: 'Krok rusza, gdy skończą się',
      wybor: [
        [OrchestrationGateRule.All, 'wszyscy poprzednicy'],
        [OrchestrationGateRule.Any, 'dowolny poprzednik'],
        [OrchestrationGateRule.Count, 'wskazana liczba poprzedników'],
      ],
    },
    { klucz: 'liczba', etykieta: 'Liczba poprzedników', podpowiedz: 'Tylko dla reguły liczbowej' },
  ]);
  if (wartosci === null) return;
  const liczba = Number.parseInt(wartosci.liczba ?? '', 10);
  const wynik = await wywolaj(otoczenie.kanal, Command.OrchestrationGateSet, {
    workflowId: wartosci.przeplyw ?? '',
    stepId: wartosci.krok ?? '',
    rule: (wartosci.regula ?? OrchestrationGateRule.All) as OrchestrationGateRule,
    ...(Number.isFinite(liczba) ? { count: liczba } : {}),
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił postawienia bramy.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Brama stoi przed krokiem ${wartosci.krok ?? ''}.`);
  otoczenie.odswiez();
}

/* Krok wyrównujący biegnie, gdy wskazany krok się nie powiedzie — po to, żeby
   cofnąć to, co zdążył zrobić. Puste wskazanie zdejmuje wyrównanie. */
async function wskazWyrownanie(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazAutomatyke(otoczenie, panel, 'Krok wyrównujący',
    'Wskaż wyrównanie', [
      { klucz: 'krok', etykieta: 'Krok, który może zawieść', wymagane: true },
      {
        klucz: 'wyrownanie',
        etykieta: 'Krok wyrównujący',
        podpowiedz: 'Puste zdejmuje wyrównanie',
      },
    ], 'Krok wyrównujący biegnie po niepowodzeniu wskazanego kroku.');
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.OrchestrationCompensationSet, {
    workflowId: wartosci.przeplyw ?? '',
    stepId: wartosci.krok ?? '',
    ...(wartosci.wyrownanie === '' ? {} : { compensationStepId: wartosci.wyrownanie }),
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wskazania wyrównania.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, wartosci.wyrownanie === ''
    ? `Wyrównanie kroku ${wartosci.krok ?? ''} zdjęte.`
    : `Krok ${wartosci.wyrownanie} wyrówna niepowodzenie kroku ${wartosci.krok ?? ''}.`);
  otoczenie.odswiez();
}

async function sprawdzUklad(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazAutomatyke(otoczenie, panel, 'Sprawdzenie układu', 'Sprawdź układ');
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.OrchestrationValidate, {
    workflowId: wartosci.przeplyw ?? '',
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił sprawdzenia układu.', 'ostrzezenie');
    return;
  }
  const sciezka = wynik.wynik.criticalPathStepIds ?? [];
  wypelnij(otoczenie.korzen, [
    ...(wynik.wynik.issues ?? []).map((uwaga) => `uwaga: ${uwaga}`),
    ...(sciezka.length === 0 ? [] : [`ścieżka rozstrzygająca: ${sciezka.join(' → ')}`]),
  ]);
  oglos(NAGLOWEK, wynik.wynik.valid
    ? 'Układ kroków bez zastrzeżeń.'
    : `Układ ma ${wynik.wynik.issues?.length ?? 0} zastrzeżeń.`,
  wynik.wynik.valid ? 'informacja' : 'ostrzezenie');
}

async function przestawWiazanie(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazAutomatyke(otoczenie, panel, 'Wiązanie z orkiestracją',
    'Przestaw wiązanie', [{
      klucz: 'stan',
      etykieta: 'Automatyka w orkiestracji',
      wybor: [['tak', 'Związana'], ['nie', 'Rozwiązana']],
    }]);
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.OrchestrationMultitaskingLink, {
    workflowId: wartosci.przeplyw ?? '',
    linked: wartosci.stan === 'tak',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zmiany wiązania.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, wartosci.stan === 'tak'
    ? 'Automatyka związana z orkiestracją.'
    : 'Automatyka rozwiązana z orkiestracji.');
  otoczenie.odswiez();
}

async function wyborHarmonogramow(
  kanal: Kanal,
): Promise<ReadonlyArray<readonly [string, string]>> {
  const wykaz = await wywolaj(kanal, Command.ScheduleGet, {});
  return (wykaz.wynik?.schedules ?? []).map((harmonogram: AutomationSchedule) =>
    [harmonogram.id, `${harmonogram.cron ?? 'bez zapisu czasu'}`
      + `${harmonogram.enabled ? '' : ' · wstrzymany'}`] as const);
}

async function wykazHarmonogramow({ kanal, korzen }: Otoczenie): Promise<void> {
  const wynik = await wywolaj(kanal, Command.ScheduleGet, {});
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu harmonogramów.', 'ostrzezenie');
    return;
  }
  const harmonogramy = wynik.wynik.schedules;
  wypelnij(korzen, harmonogramy.map((harmonogram: AutomationSchedule) =>
    `${harmonogram.enabled ? 'czynny' : 'wstrzymany'} · ${harmonogram.cron ?? 'bez zapisu czasu'}`
    + `${harmonogram.nextRunAt === undefined ? '' : ' · następny bieg zapowiedziany'}`));
  if (harmonogramy.length === 0) oglos(NAGLOWEK, 'Żadna automatyka nie ma harmonogramu.');
}

/* Okno biegu podaje się godzinami doby, bo tak Operator o nim myśli; rdzeń
   liczy minutami od północy i okno przelicza to za niego. */
async function ustawOknoCzasu(otoczenie: Otoczenie, panel: string): Promise<void> {
  const harmonogramy = await wyborHarmonogramow(otoczenie.kanal);
  if (harmonogramy.length === 0) {
    oglos(NAGLOWEK, 'Żadna automatyka nie ma harmonogramu.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWSzufladzie(otoczenie.korzen, panel, {
    tytul: 'Okno biegu',
    opis: 'Poza tym oknem harmonogram nie uruchomi automatyki.',
    pola: [
      { klucz: 'harmonogram', etykieta: 'Harmonogram', wybor: harmonogramy },
      { klucz: 'od', etykieta: 'Od godziny', wartosc: '8', wymagane: true },
      { klucz: 'do', etykieta: 'Do godziny', wartosc: '18', wymagane: true },
    ],
    wykonanie: 'Ustaw okno biegu',
  });
  if (wartosci === null) return;
  const od = Number.parseInt(wartosci.od ?? '', 10);
  const dokad = Number.parseInt(wartosci.do ?? '', 10);
  if (!Number.isFinite(od) || !Number.isFinite(dokad)) {
    oglos(NAGLOWEK, 'Godziny muszą być liczbami.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(otoczenie.kanal, Command.ScheduleWindowSet, {
    scheduleId: wartosci.harmonogram ?? '',
    windows: [{ fromMinuteOfDay: od * 60, toMinuteOfDay: dokad * 60 }],
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił ustawienia okna biegu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Automatyka rusza między ${od} a ${dokad} godziną.`);
  otoczenie.odswiez();
}

/* Tolerancja tętna mówi, po jakim milczeniu rdzeń uzna, że bieg nie ruszył —
   bez czujki nikt się o tym nie dowie, więc okno nazywa to wprost. */
async function ustawTetno(otoczenie: Otoczenie, panel: string): Promise<void> {
  const harmonogramy = await wyborHarmonogramow(otoczenie.kanal);
  if (harmonogramy.length === 0) {
    oglos(NAGLOWEK, 'Żadna automatyka nie ma harmonogramu.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWSzufladzie(otoczenie.korzen, panel, {
    tytul: 'Tolerancja tętna',
    opis: 'O przekroczeniu tolerancji melduje czujka; bez niej nikt się nie dowie.',
    pola: [
      { klucz: 'harmonogram', etykieta: 'Harmonogram', wybor: harmonogramy },
      { klucz: 'sekundy', etykieta: 'Tolerancja w sekundach', wartosc: '300', wymagane: true },
    ],
    wykonanie: 'Ustaw tolerancję',
  });
  if (wartosci === null) return;
  const sekundy = Number.parseInt(wartosci.sekundy ?? '', 10);
  if (!Number.isFinite(sekundy)) {
    oglos(NAGLOWEK, 'Tolerancja musi być liczbą sekund.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(otoczenie.kanal, Command.ScheduleHeartbeatSet, {
    scheduleId: wartosci.harmonogram ?? '',
    toleranceSeconds: sekundy,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił ustawienia tolerancji.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Tolerancja tętna to ${sekundy} sekund.`);
}

/* Nadrobienie uruchamia biegi za czas miniony, więc szuflada pyta o granicę
   liczby biegów — bez niej długi zakres zalałby kolejkę. */
async function nadrobBiegi(otoczenie: Otoczenie, panel: string): Promise<void> {
  const harmonogramy = await wyborHarmonogramow(otoczenie.kanal);
  if (harmonogramy.length === 0) {
    oglos(NAGLOWEK, 'Żadna automatyka nie ma harmonogramu.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWSzufladzie(otoczenie.korzen, panel, {
    tytul: 'Nadrobienie zaległych biegów',
    pola: [
      { klucz: 'harmonogram', etykieta: 'Harmonogram', wybor: harmonogramy },
      { klucz: 'godziny', etykieta: 'Ile godzin wstecz', wartosc: '24', wymagane: true },
      { klucz: 'granica', etykieta: 'Najwyżej biegów', wartosc: '10', wymagane: true },
    ],
    wykonanie: 'Nadrób biegi',
  });
  if (wartosci === null) return;
  const godziny = Number.parseInt(wartosci.godziny ?? '', 10);
  const granica = Number.parseInt(wartosci.granica ?? '', 10);
  if (!Number.isFinite(godziny) || !Number.isFinite(granica)) {
    oglos(NAGLOWEK, 'Zakres i granica muszą być liczbami.', 'ostrzezenie');
    return;
  }
  const teraz = Date.now();
  const wynik = await wywolaj(otoczenie.kanal, Command.ScheduleBackfillRun, {
    scheduleId: wartosci.harmonogram ?? '',
    fromAt: teraz - godziny * 60 * MINUTA,
    toAt: teraz,
    maxRuns: granica,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił nadrobienia biegów.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Nadrabianie ruszyło dla ${godziny} godzin wstecz.`);
  otoczenie.odswiez();
}

async function dziejeWyzwolen({ kanal, korzen }: Otoczenie): Promise<void> {
  const wynik = await wywolaj(kanal, Command.ScheduleTriggerHistory, { limit: 100 });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu dziejów wyzwoleń.',
      'ostrzezenie');
    return;
  }
  const wpisy = wynik.wynik.entries;
  wypelnij(korzen, wpisy.map((wpis) => `${wpis.cause}`
    + `${wpis.executionId === undefined ? ' · bez przebiegu' : " · przebieg " + wpis.executionId}`));
  if (wpisy.length === 0) oglos(NAGLOWEK, 'Żadne wyzwolenie nie zostało jeszcze odnotowane.');
}

/* Odświeżenie tajemnicy podpisu unieważnia poprzednią, więc wejście sieciowe
   pyta o to osobno i mówi, co się stanie z tym, co już wysyła. */
async function wejscieSieciowe(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazAutomatyke(otoczenie, panel, 'Wejście sieciowe',
    'Odczytaj adres', [{
      klucz: 'odswiez',
      etykieta: 'Tajemnica podpisu',
      wybor: [
        ['nie', 'Zostaw stojącą'],
        ['tak', 'Odśwież — dotychczasowi nadawcy przestaną być przyjmowani'],
      ],
    }]);
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.ScheduleWebhookEndpointGet, {
    workflowId: wartosci.przeplyw ?? '',
    rotateSecret: wartosci.odswiez === 'tak',
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił podania wejścia.', 'ostrzezenie');
    return;
  }
  wypelnij(otoczenie.korzen, [
    `adres wejścia: ${wynik.wynik.endpointUrl}`,
    `tajemnica podpisu: ${wynik.wynik.signatureSecretRef}`,
    `okno odrzucania powtórzeń: ${wynik.wynik.deduplicationWindowSeconds} sekund`,
  ]);
  oglos(NAGLOWEK, 'Adres wejścia sieciowego stoi w planie.');
}


const ROLE: ReadonlyArray<readonly [string, string]> = [
  [WindowRole.Standalone, 'Osobne'],
  [WindowRole.Coordinator, 'Prowadzące'],
  [WindowRole.Executor, 'Wykonawcze'],
];

function nazwaRoli(rola: string): string {
  return ROLE.find((pozycja) => pozycja[0] === rola)?.[1] ?? rola;
}

/* Rola mówi, czy okno pracuje samo, prowadzi inne, czy wykonuje zlecone —
   wykaz bierze wszystkie okna, nie tylko podległe jednemu prowadzącemu. */
async function wykazRol(otoczenie: Otoczenie): Promise<void> {
  const wynik = await wywolaj(otoczenie.kanal, Command.RoleList, {});
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu ról.', 'ostrzezenie');
    return;
  }
  const nadania = wynik.wynik.assignments;
  wypelnij(otoczenie.korzen, nadania.length === 0
    ? ['Żadne okno nie ma jeszcze nadanej roli.']
    : nadania.map((nadanie) => `${nadanie.windowId} · ${nazwaRoli(nadanie.role)}`
      + (nadanie.persona === undefined || nadanie.persona === ''
        ? '' : ` · ${nadanie.persona}`)
      + (nadanie.coordinatorWindowId === undefined
        ? '' : ` · prowadzi ${nadanie.coordinatorWindowId}`)));
  oglos(NAGLOWEK, `Ról nadanych: ${String(wynik.wynik.total)}.`);
}

async function wyborOkien(kanal: Kanal): Promise<ReadonlyArray<readonly [string, string]>> {
  const wykaz = await wywolaj(kanal, Command.WindowList, {});
  return (wykaz.wynik?.windows ?? []).map((okno) => [okno.id, okno.title ?? okno.id] as const);
}

async function wskazOkno(
  otoczenie: Otoczenie,
  panel: string,
  tytul: string,
  wykonanie: string,
  dodatkowe: Parameters<typeof zapytajWSzufladzie>[2]['pola'] = [],
  opis?: string,
): Promise<Record<string, string> | null> {
  const wybor = await wyborOkien(otoczenie.kanal);
  if (wybor.length === 0) {
    oglos(NAGLOWEK, 'Nie ma jeszcze żadnego okna komunikacji.', 'ostrzezenie');
    return null;
  }
  return zapytajWSzufladzie(otoczenie.korzen, panel, {
    tytul,
    ...(opis === undefined ? {} : { opis }),
    pola: [{ klucz: 'okno', etykieta: 'Okno komunikacji', wybor }, ...dodatkowe],
    wykonanie,
  });
}

async function nadajRole(otoczenie: Otoczenie, panel: string): Promise<void> {
  const prowadzace = await wyborOkien(otoczenie.kanal);
  const wartosci = await wskazOkno(otoczenie, panel, 'Nadanie roli oknu', 'Nadaj rolę', [
    { klucz: 'rola', etykieta: 'Rola', wybor: ROLE },
    {
      klucz: 'prowadzace',
      etykieta: 'Okno prowadzące',
      wybor: [['', 'bez prowadzącego'] as const, ...prowadzace],
    },
  ], 'Okno wykonawcze pracuje pod prowadzącym; osobne nie ma prowadzącego.');
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.RoleAssign, {
    windowId: wartosci.okno ?? '',
    role: (wartosci.rola ?? WindowRole.Standalone) as WindowRole,
    ...(wartosci.prowadzace === '' ? {} : { coordinatorWindowId: wartosci.prowadzace }),
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił nadania roli.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Okno ma teraz rolę: ${nazwaRoli(wartosci.rola ?? '')}.`);
  otoczenie.odswiez();
}

async function poprawRole(otoczenie: Otoczenie, panel: string): Promise<void> {
  const prowadzace = await wyborOkien(otoczenie.kanal);
  const wartosci = await wskazOkno(otoczenie, panel, 'Poprawa roli okna', 'Zapisz poprawki', [
    { klucz: 'rola', etykieta: 'Rola', wybor: ROLE },
    { klucz: 'postac', etykieta: 'Postać okna' },
    {
      klucz: 'prowadzace',
      etykieta: 'Okno prowadzące',
      wybor: [['', 'bez prowadzącego'] as const, ...prowadzace],
    },
  ]);
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.RoleUpdate, {
    windowId: wartosci.okno ?? '',
    role: (wartosci.rola ?? WindowRole.Standalone) as WindowRole,
    persona: wartosci.postac ?? '',
    ...(wartosci.prowadzace === '' ? {} : { coordinatorWindowId: wartosci.prowadzace }),
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił poprawy roli.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Rola okna poprawiona.');
  otoczenie.odswiez();
}

/* Zdjęcie roli zostawia okno bez przydziału: przestaje być prowadzone i samo
   przestaje prowadzić, więc czynność pyta o potwierdzenie. */
async function zdejmijRole(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazOkno(otoczenie, panel, 'Zdjęcie roli okna', 'Zdejmij rolę', [],
    'Okno przestanie być prowadzone i samo przestanie prowadzić inne okna.');
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.RoleRemove, {
    windowId: wartosci.okno ?? '',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zdjęcia roli.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Rola okna zdjęta.');
  otoczenie.odswiez();
}
