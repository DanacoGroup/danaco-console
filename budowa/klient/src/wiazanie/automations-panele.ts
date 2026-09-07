// Czynności okna Automations: przepływ i jego wersje w panelu orkiestracji,
// przebiegi w monitorze, tajemnice w harmonogramie, kolejka we własnym panelu.
import {
  AutomationAlertTrigger,
  Command,
  QueueAction,
  QueueItemStatus,
} from '../../../shared/contract.ts';
import type {
  AutomationCheckpoint,
  AutomationExecution,
  AutomationExecutionStep,
  AutomationLogEntry,
  AutomationVersionChange,
  AutomationWorkflowVersion,
  QueueItem,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { dolozCzynnosciPanelu, zapytajWSzufladzie } from './czynnosci-okna.ts';
import { oglos } from './ogloszenie.ts';

const NAGLOWEK = 'Automatyki';

/* Przebieg wskazany w monitorze żyje w oknie: kontrakt nie ma wykazu
   przebiegów, a zapis na zdarzenia oddaje ich komplet przy zakładaniu. */
const PRZEBIEGI = new Map<string, string>();

interface Otoczenie {
  kanal: Kanal;
  korzen: Element;
  idOkna: () => string;
  idSesji: () => string;
  odswiez: () => void;
}

export function zwiazPaneleAutomatyk(
  kanal: Kanal,
  korzen: Element,
  idOkna: () => string,
  idSesji: () => string,
  odswiez: () => void,
  dolacz: (zdejmij: () => void) => void,
  przy: AddEventListenerOptions,
): void {
  const otoczenie: Otoczenie = { kanal, korzen, idOkna, idSesji, odswiez };
  const zdejmowanie = [
    dolozCzynnosciPanelu(korzen, 'panel-orch', 'Czynności przepływu', [
      {
        naglowek: 'Wersje',
        pozycje: [
          { kod: 'wersje', nazwa: 'Wykaz wersji…' },
          { kod: 'roznica', nazwa: 'Różnica wersji…' },
          { kod: 'przywroc', nazwa: 'Przywróć wersję…' },
        ],
      },
      {
        naglowek: 'Wydanie',
        pozycje: [
          { kod: 'oglos', nazwa: 'Ogłoś automatykę…' },
          { kod: 'udostepnij', nazwa: 'Przestaw udostępnienie…' },
          { kod: 'proba', nazwa: 'Przeprowadź próbę…' },
        ],
      },
      {
        naglowek: 'Ustawienia przepływu',
        pozycje: [
          { kod: 'znaczniki', nazwa: 'Ustaw znaczniki…' },
          { kod: 'zmienna', nazwa: 'Ustaw zmienną…' },
          { kod: 'harmonogram', nazwa: 'Ustaw harmonogram…' },
          { kod: 'zaleznosc', nazwa: 'Ustaw zależność kroków…' },
          { kod: 'polozenie', nazwa: 'Ustaw położenie kroku…' },
          { kod: 'przypis', nazwa: 'Dołóż przypis kroku…' },
          { kod: 'granice', nazwa: 'Ustaw granice czasu…' },
        ],
      },
    ], (kod) => {
      void wykonaj(otoczenie, kod, 'panel-orch');
    }, przy),
    dolozCzynnosciPanelu(korzen, 'panel-builder', 'Czynności wzorców', [
      {
        naglowek: 'Wzorce',
        pozycje: [
          { kod: 'wzorzec-zapisz', nazwa: 'Zapisz wzorzec…' },
          { kod: 'wzorzec-zaloz', nazwa: 'Załóż z wzorca…' },
        ],
      },
    ], (kod) => {
      void wykonaj(otoczenie, kod, 'panel-builder');
    }, przy),
    dolozCzynnosciPanelu(korzen, 'panel-monitor', 'Czynności przebiegów', [
      {
        naglowek: 'Wgląd w przebieg',
        pozycje: [
          { kod: 'sledz', nazwa: 'Zapisz okno na przebiegi' },
          { kod: 'dziennik', nazwa: 'Dziennik przebiegu…' },
          { kod: 'kroki', nazwa: 'Kroki przebiegu' },
          { kod: 'ladunek', nazwa: 'Ładunek kroku…' },
          { kod: 'zaczepienia', nazwa: 'Punkty zaczepienia' },
        ],
      },
      {
        naglowek: 'Prowadzenie przebiegu',
        pozycje: [
          { kod: 'powtorz', nazwa: 'Powtórz przebieg…' },
          { kod: 'wznow', nazwa: 'Wznów przebieg' },
        ],
      },
      {
        naglowek: 'Czujki',
        pozycje: [{ kod: 'czujka', nazwa: 'Ustaw czujkę…' }],
      },
    ], (kod) => {
      void wykonaj(otoczenie, kod, 'panel-monitor');
    }, przy),
    dolozCzynnosciPanelu(korzen, 'panel-scheduler', 'Czynności sejfu', [
      {
        naglowek: 'Tajemnice',
        pozycje: [
          { kod: 'tajemnica', nazwa: 'Zapisz tajemnicę…' },
          { kod: 'tajemnica-usun', nazwa: 'Usuń tajemnicę…' },
        ],
      },
    ], (kod) => {
      void wykonaj(otoczenie, kod, 'panel-scheduler');
    }, przy),
    dolozCzynnosciPanelu(korzen, 'panel-queue', 'Czynności kolejki', [
      {
        naglowek: 'Kolejka',
        pozycje: [
          { kod: 'kolejka-zaloz', nazwa: 'Załóż kolejkę…' },
          { kod: 'kolejka-przestaw', nazwa: 'Przestaw kolejkę…' },
          { kod: 'kolejka-zasady', nazwa: 'Ustaw zasady…' },
          { kod: 'kolejka-wiaz', nazwa: 'Zwiąż kolejkę z oknem' },
          { kod: 'kolejka-glebokosc', nazwa: 'Odczytaj obciążenie' },
        ],
      },
      {
        naglowek: 'Pozycje kolejki',
        pozycje: [
          { kod: 'pozycje', nazwa: 'Wykaz pozycji…' },
          { kod: 'pozycja-dodaj', nazwa: 'Dołóż pozycję…' },
          { kod: 'pozycja-zdejmij', nazwa: 'Zdejmij pozycję…' },
          { kod: 'pozycja-odlóż', nazwa: 'Odłóż pozycję…' },
          { kod: 'pozycja-skieruj', nazwa: 'Skieruj pozycję…' },
          { kod: 'pozycja-warunek', nazwa: 'Ustaw warunek pozycji…' },
          { kod: 'pozycja-rozdziel', nazwa: 'Rozdziel pozycję…' },
          { kod: 'pozycja-zloz', nazwa: 'Złóż pozycje…' },
          { kod: 'pozycja-rozgalez', nazwa: 'Rozgałęź pozycję…' },
          { kod: 'martwe', nazwa: 'Wykaz pozycji martwych' },
        ],
      },
    ], (kod) => {
      void wykonaj(otoczenie, kod, 'panel-queue');
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

/* Wybór automatyki idzie wykazem rdzenia, nie wpisem z ręki: okno pokazuje
   nazwy, które Operator ma przed sobą w panelu orkiestracji. */
async function wyborAutomatyki(
  kanal: Kanal,
): Promise<ReadonlyArray<readonly [string, string]>> {
  const wykaz = await wywolaj(kanal, Command.AutomationWorkflowList, {});
  return (wykaz.wynik?.workflows ?? []).map((przeplyw) =>
    [przeplyw.id, `${przeplyw.name}${przeplyw.enabled ? '' : ' · wyłączona'}`] as const);
}

async function wykonaj(otoczenie: Otoczenie, kod: string, panel: string): Promise<void> {
  if (kod === 'wersje') return wykazWersji(otoczenie, panel);
  if (kod === 'roznica') return roznicaWersji(otoczenie, panel);
  if (kod === 'przywroc') return przywrocWersje(otoczenie, panel);
  if (kod === 'oglos') return oglosAutomatyke(otoczenie, panel);
  if (kod === 'udostepnij') return przestawUdostepnienie(otoczenie, panel);
  if (kod === 'proba') return przeprowadzProbe(otoczenie, panel);
  if (kod === 'znaczniki') return ustawZnaczniki(otoczenie, panel);
  if (kod === 'zmienna') return ustawZmienna(otoczenie, panel);
  if (kod === 'harmonogram') return ustawHarmonogram(otoczenie, panel);
  if (kod === 'zaleznosc') return ustawZaleznosc(otoczenie, panel);
  if (kod === 'polozenie') return ustawPolozenie(otoczenie, panel);
  if (kod === 'przypis') return dolozPrzypis(otoczenie, panel);
  if (kod === 'granice') return ustawGranice(otoczenie, panel);
  if (kod === 'wzorzec-zapisz') return zapiszWzorzec(otoczenie, panel);
  if (kod === 'wzorzec-zaloz') return zalozZWzorca(otoczenie, panel);
  if (kod === 'sledz') return sledzPrzebiegi(otoczenie);
  if (kod === 'dziennik') return dziennikPrzebiegu(otoczenie, panel);
  if (kod === 'kroki') return krokiPrzebiegu(otoczenie);
  if (kod === 'ladunek') return ladunekKroku(otoczenie, panel);
  if (kod === 'zaczepienia') return punktyZaczepienia(otoczenie);
  if (kod === 'powtorz') return powtorzPrzebieg(otoczenie, panel);
  if (kod === 'wznow') return wznowPrzebieg(otoczenie);
  if (kod === 'czujka') return ustawCzujke(otoczenie, panel);
  if (kod === 'tajemnica') return zapiszTajemnice(otoczenie, panel);
  if (kod === 'tajemnica-usun') return usunTajemnice(otoczenie, panel);
  if (kod === 'kolejka-zaloz') return zalozKolejke(otoczenie, panel);
  if (kod === 'kolejka-przestaw') return przestawKolejke(otoczenie, panel);
  if (kod === 'kolejka-zasady') return ustawZasady(otoczenie, panel);
  if (kod === 'kolejka-wiaz') return zwiazKolejke(otoczenie);
  if (kod === 'kolejka-glebokosc') return odczytajObciazenie(otoczenie);
  if (kod === 'pozycje') return wykazPozycji(otoczenie, panel);
  if (kod === 'pozycja-dodaj') return dolozPozycje(otoczenie, panel);
  if (kod === 'pozycja-zdejmij') return zdejmijPozycje(otoczenie, panel);
  if (kod === 'pozycja-odlóż') return odlozPozycje(otoczenie, panel);
  if (kod === 'pozycja-skieruj') return skierujPozycje(otoczenie, panel);
  if (kod === 'pozycja-warunek') return warunekPozycji(otoczenie, panel);
  if (kod === 'pozycja-rozdziel') return rozdzielPozycje(otoczenie, panel);
  if (kod === 'pozycja-zloz') return zlozPozycje(otoczenie, panel);
  if (kod === 'pozycja-rozgalez') return rozgalezPozycje(otoczenie, panel);
  if (kod === 'martwe') return wykazMartwych(otoczenie);
  /* Czynność spoza obsłużonych odmawia zamiast milczeć: pozycja bez gałęzi
     wyglądałaby jak działająca. */
  oglos(NAGLOWEK, `Czynność „${kod}" nie jest prowadzona przez to okno.`, 'ostrzezenie');
}

async function wskazAutomatyke(
  otoczenie: Otoczenie,
  panel: string,
  tytul: string,
  wykonanie: string,
  dodatkowe: Parameters<typeof zapytajWSzufladzie>[2]['pola'] = [],
  nieodwracalne?: string,
): Promise<Record<string, string> | null> {
  const wybor = await wyborAutomatyki(otoczenie.kanal);
  if (wybor.length === 0) {
    oglos(NAGLOWEK, 'To okno nie ma jeszcze żadnej automatyki.', 'ostrzezenie');
    return null;
  }
  return zapytajWSzufladzie(otoczenie.korzen, panel, {
    tytul,
    pola: [{ klucz: 'przeplyw', etykieta: 'Automatyka', wybor }, ...dodatkowe],
    wykonanie,
    ...(nieodwracalne === undefined ? {} : { nieodwracalne }),
  });
}

async function wykazWersji(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazAutomatyke(otoczenie, panel, 'Wersje automatyki', 'Odczytaj wersje');
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.AutomationWorkflowVersionList, {
    workflowId: wartosci.przeplyw ?? '',
    limit: 50,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu wersji.', 'ostrzezenie');
    return;
  }
  const wersje = wynik.wynik.versions;
  wypelnij(otoczenie.korzen, 'panel-plan', wersje.map((wersja: AutomationWorkflowVersion) =>
    `wersja ${wersja.version} · kroków ${wersja.steps.length}`
    + `${wersja.published ? ' · ogłoszona' : ''}`));
  if (wersje.length === 0) oglos(NAGLOWEK, 'Ta automatyka nie ma jeszcze wersji.');
}

async function roznicaWersji(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazAutomatyke(otoczenie, panel, 'Różnica wersji', 'Porównaj wersje', [
    { klucz: 'od', etykieta: 'Wersja wcześniejsza', wymagane: true },
    { klucz: 'do', etykieta: 'Wersja późniejsza', wymagane: true },
  ]);
  if (wartosci === null) return;
  const od = Number.parseInt(wartosci.od ?? '', 10);
  const dokad = Number.parseInt(wartosci.do ?? '', 10);
  if (!Number.isFinite(od) || !Number.isFinite(dokad)) {
    oglos(NAGLOWEK, 'Numery wersji muszą być liczbami.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(otoczenie.kanal, Command.AutomationWorkflowVersionDiff, {
    workflowId: wartosci.przeplyw ?? '',
    fromVersion: od,
    toVersion: dokad,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił porównania wersji.', 'ostrzezenie');
    return;
  }
  const zmiany = wynik.wynik.changes;
  wypelnij(otoczenie.korzen, 'panel-plan', zmiany.map((zmiana: AutomationVersionChange) =>
    `${zmiana.change} ${zmiana.stepId}${zmiana.field === undefined ? '' : ` · ${zmiana.field}`}`));
  if (zmiany.length === 0) oglos(NAGLOWEK, 'Te wersje nie różnią się niczym.');
}

async function przywrocWersje(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazAutomatyke(otoczenie, panel, 'Przywrócenie wersji',
    'Przywróć wersję',
    [{ klucz: 'wersja', etykieta: 'Numer wersji', wymagane: true }],
    'Kroki stojące zostaną zastąpione krokami wskazanej wersji.');
  if (wartosci === null) return;
  const numer = Number.parseInt(wartosci.wersja ?? '', 10);
  if (!Number.isFinite(numer)) {
    oglos(NAGLOWEK, 'Numer wersji musi być liczbą.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(otoczenie.kanal, Command.AutomationWorkflowVersionRestore, {
    workflowId: wartosci.przeplyw ?? '',
    version: numer,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił przywrócenia wersji.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Kroki wzięte z wersji ${numer}.`);
  otoczenie.odswiez();
}

async function oglosAutomatyke(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazAutomatyke(otoczenie, panel, 'Ogłoszenie automatyki',
    'Ogłoś automatykę', [],
    'Ogłoszona automatyka rusza wedle swojego harmonogramu i wyzwalaczy.');
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.AutomationWorkflowPublish, {
    workflowId: wartosci.przeplyw ?? '',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił ogłoszenia automatyki.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Automatyka ogłoszona.');
  otoczenie.odswiez();
}

async function przestawUdostepnienie(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazAutomatyke(otoczenie, panel, 'Udostępnienie automatyki',
    'Przestaw udostępnienie', [{
      klucz: 'stan',
      etykieta: 'Udostępnienie',
      wybor: [['tak', 'Udostępniona innym kontom'], ['nie', 'Tylko dla tego konta']],
    }]);
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.AutomationWorkflowShare, {
    workflowId: wartosci.przeplyw ?? '',
    shared: wartosci.stan === 'tak',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zmiany udostępnienia.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, wartosci.stan === 'tak'
    ? 'Automatyka udostępniona innym kontom.'
    : 'Udostępnienie zdjęte.');
}

/* Próba nie rusza kroków naprawdę: rdzeń oddaje ich wynik bez wykonania. */
async function przeprowadzProbe(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazAutomatyke(otoczenie, panel, 'Próba automatyki',
    'Przeprowadź próbę');
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.AutomationWorkflowSimulate, {
    workflowId: wartosci.przeplyw ?? '',
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił próby.', 'ostrzezenie');
    return;
  }
  wypelnij(otoczenie.korzen, 'panel-plan', wynik.wynik.results.map((krok) =>
    `${krok.status} ${krok.stepId}`
    + `${krok.errorMessage === undefined ? '' : ` · ${krok.errorMessage}`}`));
  oglos(NAGLOWEK, wynik.wynik.succeeded
    ? 'Próba przeszła bez zastrzeżeń.'
    : `Próba zgłasza ${wynik.wynik.issues?.length ?? 0} zastrzeżeń.`,
  wynik.wynik.succeeded ? 'informacja' : 'ostrzezenie');
}

async function ustawZnaczniki(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazAutomatyke(otoczenie, panel, 'Znaczniki automatyki',
    'Ustaw znaczniki', [{
      klucz: 'znaczniki',
      etykieta: 'Znaczniki',
      podpowiedz: 'Rozdzielone spacją; puste zdejmuje wszystkie',
    }]);
  if (wartosci === null) return;
  const znaczniki = (wartosci.znaczniki ?? '').split(/\s+/).filter((znak) => znak !== '');
  const wynik = await wywolaj(otoczenie.kanal, Command.AutomationWorkflowTagSet, {
    workflowId: wartosci.przeplyw ?? '',
    tags: znaczniki,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił ustawienia znaczników.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, znaczniki.length === 0
    ? 'Znaczniki zdjęte z automatyki.'
    : `Automatyka ma ${znaczniki.length} znaczników.`);
}

/* Komenda przyjmuje komplet zmiennych naraz, więc zapis jednej zdejmuje
   poprzednie; okno mówi to wprost w opisie szuflady. */
async function ustawZmienna(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wybor = await wyborAutomatyki(otoczenie.kanal);
  if (wybor.length === 0) {
    oglos(NAGLOWEK, 'To okno nie ma jeszcze żadnej automatyki.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWSzufladzie(otoczenie.korzen, panel, {
    tytul: 'Zmienna automatyki',
    opis: 'Rdzeń przyjmuje komplet zmiennych naraz — ta zmienna zastąpi wszystkie stojące.',
    pola: [
      { klucz: 'przeplyw', etykieta: 'Automatyka', wybor },
      { klucz: 'nazwa', etykieta: 'Nazwa zmiennej', wymagane: true },
      { klucz: 'rodzaj', etykieta: 'Rodzaj wartości', wartosc: 'tekst' },
    ],
    wykonanie: 'Ustaw zmienną',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.AutomationWorkflowVariablesSet, {
    workflowId: wartosci.przeplyw ?? '',
    variables: [{ name: wartosci.nazwa ?? '', kind: wartosci.rodzaj ?? 'tekst' }],
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił ustawienia zmiennej.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Automatyka ma jedną zmienną „${wartosci.nazwa ?? ''}".`);
}

async function ustawHarmonogram(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazAutomatyke(otoczenie, panel, 'Harmonogram automatyki',
    'Ustaw harmonogram', [
      { klucz: 'zapis', etykieta: 'Zapis czasu', podpowiedz: '0 6 * * *', wymagane: true },
      { klucz: 'strefa', etykieta: 'Strefa czasu', podpowiedz: 'Europe/Warsaw' },
      {
        klucz: 'czynny',
        etykieta: 'Harmonogram',
        wybor: [['tak', 'Czynny'], ['nie', 'Wstrzymany']],
      },
    ]);
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.AutomationScheduleSet, {
    workflowId: wartosci.przeplyw ?? '',
    cron: wartosci.zapis ?? '',
    ...(wartosci.strefa === '' ? {} : { timeZone: wartosci.strefa }),
    enabled: wartosci.czynny !== 'nie',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił ustawienia harmonogramu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Automatyka rusza wedle zapisu „${wartosci.zapis ?? ''}".`);
  otoczenie.odswiez();
}

async function ustawZaleznosc(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazAutomatyke(otoczenie, panel, 'Zależność kroków',
    'Ustaw zależność', [
      { klucz: 'od', etykieta: 'Krok wcześniejszy', wymagane: true },
      { klucz: 'do', etykieta: 'Krok późniejszy', wymagane: true },
    ]);
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.AutomationOrchestratorDefine, {
    workflowId: wartosci.przeplyw ?? '',
    dependencies: [{
      fromStepId: wartosci.od ?? '',
      toStepId: wartosci.do ?? '',
      kind: 'sequential',
    }],
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił ustawienia zależności.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Krok ${wartosci.do ?? ''} idzie po kroku ${wartosci.od ?? ''}.`);
}

async function ustawPolozenie(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazAutomatyke(otoczenie, panel, 'Położenie kroku',
    'Ustaw położenie', [
      { klucz: 'krok', etykieta: 'Oznaczenie kroku', wymagane: true },
      { klucz: 'poziomo', etykieta: 'Współrzędna pozioma', wymagane: true },
      { klucz: 'pionowo', etykieta: 'Współrzędna pionowa', wymagane: true },
    ]);
  if (wartosci === null) return;
  const x = Number.parseInt(wartosci.poziomo ?? '', 10);
  const y = Number.parseInt(wartosci.pionowo ?? '', 10);
  if (!Number.isFinite(x) || !Number.isFinite(y)) {
    oglos(NAGLOWEK, 'Współrzędne muszą być liczbami.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(otoczenie.kanal, Command.AutomationStepLayoutSet, {
    workflowId: wartosci.przeplyw ?? '',
    positions: [{ stepId: wartosci.krok ?? '', x, y }],
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił ustawienia położenia.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Krok ${wartosci.krok ?? ''} stoi na ${x}, ${y}.`);
}

async function dolozPrzypis(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazAutomatyke(otoczenie, panel, 'Przypis kroku', 'Dołóż przypis', [
    { klucz: 'krok', etykieta: 'Oznaczenie kroku', wymagane: true },
    { klucz: 'tresc', etykieta: 'Treść przypisu', obszerne: true, wymagane: true },
  ]);
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.AutomationStepNoteSet, {
    workflowId: wartosci.przeplyw ?? '',
    stepId: wartosci.krok ?? '',
    note: wartosci.tresc ?? '',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu przypisu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Przypis dołożony do kroku ${wartosci.krok ?? ''}.`);
}

async function ustawGranice(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazAutomatyke(otoczenie, panel, 'Granice czasu biegu',
    'Ustaw granice', [
      { klucz: 'bieg', etykieta: 'Sekund na cały bieg', wymagane: true },
      { klucz: 'krok', etykieta: 'Sekund na krok' },
    ]);
  if (wartosci === null) return;
  const bieg = Number.parseInt(wartosci.bieg ?? '', 10);
  const krok = Number.parseInt(wartosci.krok ?? '', 10);
  if (!Number.isFinite(bieg)) {
    oglos(NAGLOWEK, 'Granica biegu musi być liczbą sekund.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(otoczenie.kanal, Command.AutomationExecutionBudgetSet, {
    workflowId: wartosci.przeplyw ?? '',
    runBudgetSeconds: bieg,
    ...(Number.isFinite(krok) ? { stepBudgetSeconds: krok } : {}),
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił ustawienia granic.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Bieg ma ${bieg} sekund granicy.`);
}

async function zapiszWzorzec(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazAutomatyke(otoczenie, panel, 'Wzorzec z automatyki',
    'Zapisz wzorzec', [
      { klucz: 'nazwa', etykieta: 'Nazwa wzorca', wymagane: true },
      { klucz: 'opis', etykieta: 'Opis', obszerne: true },
    ]);
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.AutomationTemplateSave, {
    workflowId: wartosci.przeplyw ?? '',
    name: wartosci.nazwa ?? '',
    ...(wartosci.opis === '' ? {} : { description: wartosci.opis }),
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu wzorca.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Wzorzec „${wartosci.nazwa ?? ''}" zapisany.`);
  otoczenie.odswiez();
}

async function zalozZWzorca(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wykaz = await wywolaj(otoczenie.kanal, Command.AutomationTemplateList, {});
  const wzorce = wykaz.wynik?.templates ?? [];
  if (wzorce.length === 0) {
    oglos(NAGLOWEK, 'Nie ma jeszcze żadnego wzorca.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWSzufladzie(otoczenie.korzen, panel, {
    tytul: 'Automatyka z wzorca',
    pola: [
      {
        klucz: 'wzorzec',
        etykieta: 'Wzorzec',
        wybor: wzorce.map((wzorzec) => [wzorzec.id, wzorzec.name] as const),
      },
      { klucz: 'nazwa', etykieta: 'Nazwa nowej automatyki', wymagane: true },
    ],
    wykonanie: 'Załóż automatykę',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.AutomationTemplateApply, {
    templateId: wartosci.wzorzec ?? '',
    name: wartosci.nazwa ?? '',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił założenia z wzorca.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Automatyka „${wartosci.nazwa ?? ''}" założona z wzorca.`);
  otoczenie.odswiez();
}

/* Zapis okna na przebiegi oddaje ich komplet, więc jest też jedyną drogą
   wskazania przebiegu dla pozostałych czynności monitora. */
async function sledzPrzebiegi(otoczenie: Otoczenie): Promise<void> {
  const wynik = await wywolaj(otoczenie.kanal, Command.AutomationExecutionSubscribe, {
    windowId: otoczenie.idOkna(),
    limit: 50,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu na przebiegi.', 'ostrzezenie');
    return;
  }
  const przebiegi = wynik.wynik.executions;
  wypelnij(otoczenie.korzen, 'panel-monitor', przebiegi.map((przebieg: AutomationExecution) =>
    `${przebieg.status} ${przebieg.id}`
    + `${przebieg.errorMessage === undefined ? '' : ` · ${przebieg.errorMessage}`}`));
  const pierwszy = przebiegi[0]?.id ?? '';
  if (pierwszy !== '') PRZEBIEGI.set(otoczenie.idOkna(), pierwszy);
  oglos(NAGLOWEK, przebiegi.length === 0
    ? 'Okno śledzi przebiegi; żaden jeszcze nie ruszył.'
    : `Okno śledzi przebiegi; w wykazie ${przebiegi.length}.`);
}

function przebieg(idOkna: string): string {
  const zapamietany = PRZEBIEGI.get(idOkna) ?? '';
  if (zapamietany === '') {
    oglos(NAGLOWEK, 'Okno nie zna przebiegu — zapisz je najpierw na przebiegi.', 'ostrzezenie');
  }
  return zapamietany;
}

async function dziennikPrzebiegu(otoczenie: Otoczenie, panel: string): Promise<void> {
  const cel = przebieg(otoczenie.idOkna());
  if (cel === '') return;
  const wartosci = await zapytajWSzufladzie(otoczenie.korzen, panel, {
    tytul: 'Dziennik przebiegu',
    pola: [
      { klucz: 'krok', etykieta: 'Krok', podpowiedz: 'Puste bierze cały przebieg' },
      { klucz: 'granica', etykieta: 'Najwyżej wpisów', wartosc: '200' },
    ],
    wykonanie: 'Odczytaj dziennik',
  });
  if (wartosci === null) return;
  const granica = Number.parseInt(wartosci.granica ?? '', 10);
  const wynik = await wywolaj(otoczenie.kanal, Command.AutomationExecutionLog, {
    executionId: cel,
    ...(wartosci.krok === '' ? {} : { stepId: wartosci.krok }),
    limit: Number.isFinite(granica) ? granica : 200,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu dziennika.', 'ostrzezenie');
    return;
  }
  wypelnij(otoczenie.korzen, 'panel-monitor',
    wynik.wynik.entries.map((wpis: AutomationLogEntry) => `${wpis.level} ${wpis.message}`));
  if (wynik.wynik.truncated) {
    oglos(NAGLOWEK, 'Dziennik przycięty granicą odczytu — wpisów jest więcej.', 'ostrzezenie');
  }
}

async function krokiPrzebiegu(otoczenie: Otoczenie): Promise<void> {
  const cel = przebieg(otoczenie.idOkna());
  if (cel === '') return;
  const wynik = await wywolaj(otoczenie.kanal, Command.AutomationExecutionSteps, {
    executionId: cel,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu kroków.', 'ostrzezenie');
    return;
  }
  const kroki = wynik.wynik.steps;
  wypelnij(otoczenie.korzen, 'panel-plan', kroki.map((krok: AutomationExecutionStep) =>
    `${krok.status} ${krok.stepId}`
    + `${krok.errorMessage === undefined ? '' : ` · ${krok.errorMessage}`}`));
  if (kroki.length === 0) oglos(NAGLOWEK, 'Ten przebieg nie ma jeszcze kroków.');
}

/* Rdzeń zaciera pola objęte tajemnicą i nazywa je z osobna; okno powtarza tę
   listę, żeby nikt nie brał zatartego ładunku za pełny. */
async function ladunekKroku(otoczenie: Otoczenie, panel: string): Promise<void> {
  const cel = przebieg(otoczenie.idOkna());
  if (cel === '') return;
  const wartosci = await zapytajWSzufladzie(otoczenie.korzen, panel, {
    tytul: 'Ładunek kroku',
    pola: [{ klucz: 'krok', etykieta: 'Oznaczenie kroku', wymagane: true }],
    wykonanie: 'Odczytaj ładunek',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.AutomationExecutionPayloadGet, {
    executionId: cel,
    stepId: wartosci.krok ?? '',
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu ładunku.', 'ostrzezenie');
    return;
  }
  wypelnij(otoczenie.korzen, 'panel-artefakty', [
    `wejście: ${JSON.stringify(wynik.wynik.input ?? null).slice(0, 200)}`,
    `wyjście: ${JSON.stringify(wynik.wynik.output ?? null).slice(0, 200)}`,
  ]);
  const zatarte = wynik.wynik.redactedFields ?? [];
  if (zatarte.length > 0) {
    oglos(NAGLOWEK, `Rdzeń zatarł pola objęte tajemnicą: ${zatarte.join(', ')}.`, 'ostrzezenie');
  }
}

async function punktyZaczepienia(otoczenie: Otoczenie): Promise<void> {
  const cel = przebieg(otoczenie.idOkna());
  if (cel === '') return;
  const wynik = await wywolaj(otoczenie.kanal, Command.AutomationExecutionCheckpointList, {
    executionId: cel,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu punktów zaczepienia.',
      'ostrzezenie');
    return;
  }
  const punkty = wynik.wynik.checkpoints;
  wypelnij(otoczenie.korzen, 'panel-zadania', punkty.map((punkt: AutomationCheckpoint) =>
    `${punkt.id} · krok ${punkt.stepId} · kroków za sobą ${punkt.completedStepIds.length}`));
  if (punkty.length === 0) oglos(NAGLOWEK, 'Ten przebieg nie ma punktów zaczepienia.');
}

async function powtorzPrzebieg(otoczenie: Otoczenie, panel: string): Promise<void> {
  const cel = przebieg(otoczenie.idOkna());
  if (cel === '') return;
  const wartosci = await zapytajWSzufladzie(otoczenie.korzen, panel, {
    tytul: 'Powtórka przebiegu',
    pola: [{ klucz: 'krok', etykieta: 'Od kroku', podpowiedz: 'Puste zaczyna od początku' }],
    wykonanie: 'Powtórz przebieg',
    nieodwracalne: 'Kroki wykonają się na nowo wraz ze wszystkim, co robią na zewnątrz.',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.AutomationExecutionReplay, {
    executionId: cel,
    ...(wartosci.krok === '' ? {} : { fromStepId: wartosci.krok }),
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił powtórki.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, wartosci.krok === ''
    ? 'Przebieg powtórzony od początku.'
    : `Przebieg powtórzony od kroku ${wartosci.krok}.`);
  otoczenie.odswiez();
}

async function wznowPrzebieg(otoczenie: Otoczenie): Promise<void> {
  const cel = przebieg(otoczenie.idOkna());
  if (cel === '') return;
  const wynik = await wywolaj(otoczenie.kanal, Command.AutomationExecutionResume, {
    executionId: cel,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wznowienia przebiegu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Przebieg wznowiony od ostatniego punktu zaczepienia.');
  otoczenie.odswiez();
}

async function ustawCzujke(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazAutomatyke(otoczenie, panel, 'Czujka automatyki', 'Ustaw czujkę', [
    {
      klucz: 'powod',
      etykieta: 'Kiedy melduje',
      wybor: [
        [AutomationAlertTrigger.Failure, 'Po niepowodzeniu'],
        [AutomationAlertTrigger.Timeout, 'Po przekroczeniu czasu'],
        [AutomationAlertTrigger.MissingRun, 'Gdy bieg nie ruszył'],
        [AutomationAlertTrigger.SuccessRateDrop, 'Gdy spada odsetek powodzeń'],
      ],
    },
    { klucz: 'kanaly', etykieta: 'Kanały powiadomienia', wymagane: true },
  ]);
  if (wartosci === null) return;
  const kanaly = (wartosci.kanaly ?? '').split(/\s+/).filter((nazwa) => nazwa !== '');
  const wynik = await wywolaj(otoczenie.kanal, Command.AutomationAlertRuleSet, {
    workflowId: wartosci.przeplyw ?? '',
    trigger: (wartosci.powod ?? AutomationAlertTrigger.Failure) as AutomationAlertTrigger,
    channels: kanaly,
    enabled: true,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił ustawienia czujki.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Czujka melduje kanałami: ${kanaly.join(', ')}.`);
  otoczenie.odswiez();
}

/* Wartość tajemnicy idzie do sejfu rdzenia i wraca odwołaniem; szuflada
   znika zaraz po wysłaniu, więc treść nie zostaje na widoku. */
async function zapiszTajemnice(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(otoczenie.korzen, panel, {
    tytul: 'Tajemnica automatyk',
    opis: 'Wartość idzie do sejfu rdzenia; okno pozna wyłącznie odwołanie do niej.',
    pola: [
      { klucz: 'nazwa', etykieta: 'Nazwa tajemnicy', wymagane: true },
      { klucz: 'wartosc', etykieta: 'Wartość', wymagane: true },
    ],
    wykonanie: 'Zapisz w sejfie',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.AutomationSecretSet, {
    name: wartosci.nazwa ?? '',
    value: wartosci.wartosc ?? '',
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu tajemnicy.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Tajemnica stoi w sejfie pod odwołaniem ${wynik.wynik.secret.ref}.`);
  otoczenie.odswiez();
}

async function usunTajemnice(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wykaz = await wywolaj(otoczenie.kanal, Command.AutomationSecretList, {});
  const tajemnice = wykaz.wynik?.secrets ?? [];
  if (tajemnice.length === 0) {
    oglos(NAGLOWEK, 'Sejf automatyk jest pusty.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWSzufladzie(otoczenie.korzen, panel, {
    tytul: 'Usunięcie tajemnicy',
    pola: [{
      klucz: 'tajemnica',
      etykieta: 'Tajemnica',
      wybor: tajemnice.map((tajemnica) => [tajemnica.ref, tajemnica.name] as const),
    }],
    wykonanie: 'Usuń tajemnicę',
    nieodwracalne: 'Kroki, które sięgają po to odwołanie, przestaną działać.',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.AutomationSecretRemove, {
    secretRef: wartosci.tajemnica ?? '',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił usunięcia tajemnicy.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Tajemnica zdjęta z sejfu.');
  otoczenie.odswiez();
}

async function wyborKolejki(kanal: Kanal): Promise<ReadonlyArray<readonly [string, string]>> {
  const wykaz = await wywolaj(kanal, Command.QueueList, {});
  return (wykaz.wynik?.queues ?? []).map((kolejka) =>
    [kolejka.id, `${kolejka.name ?? kolejka.id} · ${kolejka.status}`] as const);
}

async function wskazKolejke(
  otoczenie: Otoczenie,
  panel: string,
  tytul: string,
  wykonanie: string,
  dodatkowe: Parameters<typeof zapytajWSzufladzie>[2]['pola'] = [],
  nieodwracalne?: string,
): Promise<Record<string, string> | null> {
  const wybor = await wyborKolejki(otoczenie.kanal);
  if (wybor.length === 0) {
    oglos(NAGLOWEK, 'Nie ma jeszcze żadnej kolejki.', 'ostrzezenie');
    return null;
  }
  return zapytajWSzufladzie(otoczenie.korzen, panel, {
    tytul,
    pola: [{ klucz: 'kolejka', etykieta: 'Kolejka', wybor }, ...dodatkowe],
    wykonanie,
    ...(nieodwracalne === undefined ? {} : { nieodwracalne }),
  });
}

async function zalozKolejke(otoczenie: Otoczenie, panel: string): Promise<void> {
  if (otoczenie.idSesji() === '') {
    oglos(NAGLOWEK, 'Rdzeń nie dał sesji dla tej karty.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWSzufladzie(otoczenie.korzen, panel, {
    tytul: 'Nowa kolejka',
    pola: [{ klucz: 'nazwa', etykieta: 'Nazwa kolejki', wymagane: true }],
    wykonanie: 'Załóż kolejkę',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.QueueCreate, {
    sessionId: otoczenie.idSesji(),
    name: wartosci.nazwa ?? '',
    windowIds: [otoczenie.idOkna()],
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił założenia kolejki.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Kolejka „${wynik.wynik.queue.name ?? ''}" założona.`);
  otoczenie.odswiez();
}

async function przestawKolejke(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazKolejke(otoczenie, panel, 'Przestawienie kolejki',
    'Przestaw kolejkę', [{
      klucz: 'czynnosc',
      etykieta: 'Czynność',
      wybor: [
        [QueueAction.Start, 'Uruchom'],
        [QueueAction.Pause, 'Wstrzymaj'],
        [QueueAction.Resume, 'Wznów'],
        [QueueAction.Stop, 'Zatrzymaj'],
        [QueueAction.Retry, 'Ponów nieudane'],
        [QueueAction.Clear, 'Wyczyść'],
      ],
    }], 'Wyczyszczenie zdejmuje z kolejki wszystkie pozycje czekające.');
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.QueueAction, {
    queueId: wartosci.kolejka ?? '',
    action: (wartosci.czynnosc ?? QueueAction.Start) as QueueAction,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił przestawienia kolejki.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Kolejka przestawiona.');
  otoczenie.odswiez();
}

async function ustawZasady(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazKolejke(otoczenie, panel, 'Zasady kolejki', 'Ustaw zasady', [
    { klucz: 'rownolegle', etykieta: 'Najwyżej równocześnie', wartosc: '1' },
    { klucz: 'naMinute', etykieta: 'Najwyżej na minutę' },
    { klucz: 'prob', etykieta: 'Najwyżej prób pozycji', wartosc: '3' },
  ]);
  if (wartosci === null) return;
  const rownolegle = Number.parseInt(wartosci.rownolegle ?? '', 10);
  const naMinute = Number.parseInt(wartosci.naMinute ?? '', 10);
  const prob = Number.parseInt(wartosci.prob ?? '', 10);
  const wynik = await wywolaj(otoczenie.kanal, Command.QueuePolicySet, {
    queueId: wartosci.kolejka ?? '',
    policy: {
      ...(Number.isFinite(rownolegle) ? { maxConcurrent: rownolegle } : {}),
      ...(Number.isFinite(naMinute) ? { ratePerMinute: naMinute } : {}),
      ...(Number.isFinite(prob) ? { maxAttempts: prob } : {}),
    },
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił ustawienia zasad.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Zasady kolejki ustawione.');
  otoczenie.odswiez();
}

async function zwiazKolejke(otoczenie: Otoczenie): Promise<void> {
  const wybor = await wyborKolejki(otoczenie.kanal);
  if (wybor.length === 0) {
    oglos(NAGLOWEK, 'Nie ma jeszcze żadnej kolejki.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(otoczenie.kanal, Command.QueueLink, {
    queueId: wybor[0]?.[0] ?? '',
    windowIds: [otoczenie.idOkna()],
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił związania kolejki.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Kolejka związana z tym oknem.');
  otoczenie.odswiez();
}

/* Obciążenie liczone jest od godziny wstecz w krokach minutowych: krótszy
   zakres nie pokazuje narastania, dłuższy nie mieści się w panelu. */
async function odczytajObciazenie(otoczenie: Otoczenie): Promise<void> {
  const wybor = await wyborKolejki(otoczenie.kanal);
  if (wybor.length === 0) {
    oglos(NAGLOWEK, 'Nie ma jeszcze żadnej kolejki.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(otoczenie.kanal, Command.QueueDepthGet, {
    queueId: wybor[0]?.[0] ?? '',
    fromAt: Date.now() - 3600000,
    bucketSeconds: 60,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu obciążenia.', 'ostrzezenie');
    return;
  }
  const punkty = wynik.wynik.points;
  wypelnij(otoczenie.korzen, 'panel-queue', punkty.map((punkt) =>
    `czekających ${punkt.pending} · biegnących ${punkt.running}`));
  if (punkty.length === 0) oglos(NAGLOWEK, 'Kolejka nie miała ruchu w ostatniej godzinie.');
}

async function wykazPozycji(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazKolejke(otoczenie, panel, 'Pozycje kolejki', 'Odczytaj pozycje', [{
    klucz: 'stan',
    etykieta: 'Stan pozycji',
    wybor: [
      ['', 'Wszystkie'],
      [QueueItemStatus.Pending, 'Czekające'],
      [QueueItemStatus.Running, 'Biegnące'],
      [QueueItemStatus.Failed, 'Nieudane'],
      [QueueItemStatus.Succeeded, 'Udane'],
    ],
  }]);
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.QueueItemList, {
    queueId: wartosci.kolejka ?? '',
    ...(wartosci.stan === '' ? {} : { status: wartosci.stan as QueueItemStatus }),
    limit: 200,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu pozycji.', 'ostrzezenie');
    return;
  }
  const pozycje = wynik.wynik.items;
  wypelnij(otoczenie.korzen, 'panel-queue', pozycje.map((pozycja: QueueItem) =>
    `${pozycja.status} ${pozycja.id}${pozycja.attempts === undefined
      ? '' : ` · prób ${pozycja.attempts}`}`));
  if (pozycje.length === 0) oglos(NAGLOWEK, 'Ta kolejka nie ma pozycji o tym stanie.');
}

async function dolozPozycje(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazKolejke(otoczenie, panel, 'Nowa pozycja kolejki', 'Dołóż pozycję', [
    { klucz: 'ladunek', etykieta: 'Ładunek pozycji', obszerne: true, wymagane: true },
    { klucz: 'pierwszenstwo', etykieta: 'Pierwszeństwo', wartosc: '0' },
    { klucz: 'klucz', etykieta: 'Klucz jednokrotności', podpowiedz: 'Chroni przed dublem' },
  ]);
  if (wartosci === null) return;
  const pierwszenstwo = Number.parseInt(wartosci.pierwszenstwo ?? '', 10);
  const wynik = await wywolaj(otoczenie.kanal, Command.QueueItemEnqueue, {
    queueId: wartosci.kolejka ?? '',
    payload: wartosci.ladunek ?? '',
    ...(Number.isFinite(pierwszenstwo) ? { priority: pierwszenstwo } : {}),
    ...(wartosci.klucz === '' ? {} : { idempotencyKey: wartosci.klucz }),
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił dołożenia pozycji.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, wynik.wynik.duplicate
    ? 'Pozycja o tym kluczu już stoi w kolejce — rdzeń nie dołożył drugiej.'
    : `Pozycja ${wynik.wynik.item.id} dołożona.`,
  wynik.wynik.duplicate ? 'ostrzezenie' : 'informacja');
  otoczenie.odswiez();
}

/* Pozycję wskazuje wykaz kolejki: okno pobiera go w chwili otwarcia szuflady,
   więc Operator wybiera z tego, co naprawdę stoi. */
async function wyborPozycji(
  kanal: Kanal,
  idKolejki: string,
): Promise<ReadonlyArray<readonly [string, string]>> {
  const wykaz = await wywolaj(kanal, Command.QueueItemList, { queueId: idKolejki, limit: 200 });
  return (wykaz.wynik?.items ?? []).map((pozycja) =>
    [pozycja.id, `${pozycja.status} ${pozycja.id}`] as const);
}

async function czynnoscNaPozycji(
  otoczenie: Otoczenie,
  panel: string,
  tytul: string,
  wykonanie: string,
  dodatkowe: Parameters<typeof zapytajWSzufladzie>[2]['pola'] = [],
  nieodwracalne?: string,
): Promise<Record<string, string> | null> {
  const kolejki = await wyborKolejki(otoczenie.kanal);
  const idKolejki = kolejki[0]?.[0] ?? '';
  if (idKolejki === '') {
    oglos(NAGLOWEK, 'Nie ma jeszcze żadnej kolejki.', 'ostrzezenie');
    return null;
  }
  const pozycje = await wyborPozycji(otoczenie.kanal, idKolejki);
  if (pozycje.length === 0) {
    oglos(NAGLOWEK, 'Ta kolejka nie ma pozycji.', 'ostrzezenie');
    return null;
  }
  return zapytajWSzufladzie(otoczenie.korzen, panel, {
    tytul,
    pola: [
      { klucz: 'kolejka', etykieta: 'Kolejka', wybor: kolejki },
      { klucz: 'pozycja', etykieta: 'Pozycja', wybor: pozycje },
      ...dodatkowe,
    ],
    wykonanie,
    ...(nieodwracalne === undefined ? {} : { nieodwracalne }),
  });
}

async function zdejmijPozycje(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await czynnoscNaPozycji(otoczenie, panel, 'Zdjęcie pozycji', 'Zdejmij pozycję',
    [], 'Pozycja schodzi z kolejki i nie zostanie wykonana.');
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.QueueItemDequeue, {
    queueId: wartosci.kolejka ?? '',
    itemId: wartosci.pozycja ?? '',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zdjęcia pozycji.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Pozycja zdjęta z kolejki.');
  otoczenie.odswiez();
}

async function odlozPozycje(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await czynnoscNaPozycji(otoczenie, panel, 'Odłożenie pozycji', 'Odłóż pozycję',
    [{ klucz: 'sekundy', etykieta: 'O ile sekund', wartosc: '60', wymagane: true }]);
  if (wartosci === null) return;
  const sekundy = Number.parseInt(wartosci.sekundy ?? '', 10);
  if (!Number.isFinite(sekundy)) {
    oglos(NAGLOWEK, 'Odłożenie potrzebuje liczby sekund.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(otoczenie.kanal, Command.QueueItemDelay, {
    queueId: wartosci.kolejka ?? '',
    itemId: wartosci.pozycja ?? '',
    delaySeconds: sekundy,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odłożenia pozycji.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Pozycja odłożona o ${sekundy} sekund.`);
  otoczenie.odswiez();
}

async function skierujPozycje(otoczenie: Otoczenie, panel: string): Promise<void> {
  const kolejki = await wyborKolejki(otoczenie.kanal);
  const wartosci = await czynnoscNaPozycji(otoczenie, panel, 'Skierowanie pozycji',
    'Skieruj pozycję', [{ klucz: 'cel', etykieta: 'Kolejka docelowa', wybor: kolejki }]);
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.QueueItemRoute, {
    queueId: wartosci.kolejka ?? '',
    itemId: wartosci.pozycja ?? '',
    targetQueueId: wartosci.cel ?? '',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił skierowania pozycji.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Pozycja skierowana do wskazanej kolejki.');
  otoczenie.odswiez();
}

async function warunekPozycji(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await czynnoscNaPozycji(otoczenie, panel, 'Warunek pozycji', 'Ustaw warunek',
    [{ klucz: 'warunek', etykieta: 'Warunek wykonania', obszerne: true, wymagane: true }]);
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.QueueItemCondition, {
    queueId: wartosci.kolejka ?? '',
    itemId: wartosci.pozycja ?? '',
    condition: wartosci.warunek ?? '',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił ustawienia warunku.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Warunek pozycji ustawiony.');
  otoczenie.odswiez();
}

async function rozdzielPozycje(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await czynnoscNaPozycji(otoczenie, panel, 'Rozdzielenie pozycji',
    'Rozdziel pozycję', [{
      klucz: 'ladunki',
      etykieta: 'Ładunki części',
      podpowiedz: 'Każdy w osobnym wierszu',
      obszerne: true,
      wymagane: true,
    }]);
  if (wartosci === null) return;
  const ladunki = (wartosci.ladunki ?? '').split('\n').filter((wiersz) => wiersz.trim() !== '');
  const wynik = await wywolaj(otoczenie.kanal, Command.QueueItemSplit, {
    queueId: wartosci.kolejka ?? '',
    itemId: wartosci.pozycja ?? '',
    payloads: ladunki,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił rozdzielenia pozycji.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Pozycja rozdzielona na ${ladunki.length} części.`);
  otoczenie.odswiez();
}

/* Złożenie bierze wszystkie pozycje czekające: okno nie prowadzi wyboru wielu
   wierszy naraz, a stan „czekająca" jest jedynym, w którym złożenie ma sens. */
async function zlozPozycje(otoczenie: Otoczenie, panel: string): Promise<void> {
  const kolejki = await wyborKolejki(otoczenie.kanal);
  const idKolejki = kolejki[0]?.[0] ?? '';
  if (idKolejki === '') {
    oglos(NAGLOWEK, 'Nie ma jeszcze żadnej kolejki.', 'ostrzezenie');
    return;
  }
  const wykaz = await wywolaj(otoczenie.kanal, Command.QueueItemList, {
    queueId: idKolejki,
    status: QueueItemStatus.Pending,
    limit: 200,
  });
  const czekajace = (wykaz.wynik?.items ?? []).map((pozycja) => pozycja.id);
  if (czekajace.length < 2) {
    oglos(NAGLOWEK, 'Do złożenia potrzeba co najmniej dwóch pozycji czekających.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWSzufladzie(otoczenie.korzen, panel, {
    tytul: 'Złożenie pozycji',
    opis: `Złożenie obejmie ${czekajace.length} pozycji czekających.`,
    pola: [
      { klucz: 'kolejka', etykieta: 'Kolejka', wybor: kolejki },
      { klucz: 'ladunek', etykieta: 'Ładunek pozycji złożonej', obszerne: true },
    ],
    wykonanie: 'Złóż pozycje',
    nieodwracalne: 'Pozycje składowe zostaną zastąpione jedną.',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.QueueItemMerge, {
    queueId: wartosci.kolejka ?? '',
    itemIds: czekajace,
    ...(wartosci.ladunek === '' ? {} : { payload: wartosci.ladunek }),
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił złożenia pozycji.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Złożono ${czekajace.length} pozycji w jedną.`);
  otoczenie.odswiez();
}

async function rozgalezPozycje(otoczenie: Otoczenie, panel: string): Promise<void> {
  const kolejki = await wyborKolejki(otoczenie.kanal);
  const wartosci = await czynnoscNaPozycji(otoczenie, panel, 'Rozgałęzienie pozycji',
    'Rozgałęź pozycję', [
      { klucz: 'cel', etykieta: 'Kolejka gałęzi', wybor: kolejki },
      { klucz: 'warunek', etykieta: 'Warunek gałęzi', obszerne: true },
    ]);
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.QueueItemBranch, {
    queueId: wartosci.kolejka ?? '',
    itemId: wartosci.pozycja ?? '',
    branches: [{
      targetQueueId: wartosci.cel ?? '',
      ...(wartosci.warunek === '' ? {} : { condition: wartosci.warunek }),
    }],
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił rozgałęzienia pozycji.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Gałąź pozycji założona.');
  otoczenie.odswiez();
}

async function wykazMartwych(otoczenie: Otoczenie): Promise<void> {
  const wynik = await wywolaj(otoczenie.kanal, Command.QueueDeadList, { limit: 200 });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu pozycji martwych.',
      'ostrzezenie');
    return;
  }
  const pozycje = wynik.wynik.items;
  wypelnij(otoczenie.korzen, 'panel-queue', pozycje.map((pozycja: QueueItem) =>
    `${pozycja.queueId} · ${pozycja.id} · prób ${pozycja.attempts ?? 0}`));
  if (pozycje.length === 0) oglos(NAGLOWEK, 'Żadna pozycja nie wpadła do wykazu martwych.');
}
