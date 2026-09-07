// Czynności okna Apps rozłożone na panele wedle tego, czego dotykają: wyrób
// i plan w budowniczym, architektura we własnym panelu, wydanie we wdrożeniach.
import {
  AppArchitectureTemplate,
  AppDeployEnvironment,
  AppEndpointMethod,
  AppExportFormat,
  AppPackageFormat,
  AppWorkspaceLayer,
  Command,
  ExtensionKind,
} from '../../../shared/contract.ts';
import type {
  AppEnvironmentVariable,
  AppSchemaTable,
  AppTimelineEntry,
  AppValidationIssue,
  DeveloperFile,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { dolozCzynnosciPanelu, zapytajWSzufladzie } from './czynnosci-okna.ts';
import { oglos } from './ogloszenie.ts';

const NAGLOWEK = 'Aplikacje';

const SRODOWISKA: ReadonlyArray<readonly [string, string]> = [
  [AppDeployEnvironment.Dev, 'Robocze'],
  [AppDeployEnvironment.Staging, 'Przedprodukcyjne'],
  [AppDeployEnvironment.Production, 'Produkcyjne'],
];

const WARSTWY: ReadonlyArray<readonly [string, string]> = [
  [AppWorkspaceLayer.Frontend, 'Widok'],
  [AppWorkspaceLayer.Backend, 'Rdzeń'],
];

interface Otoczenie {
  kanal: Kanal;
  korzen: Element;
  idOkna: () => string;
  odswiez: () => void;
}

export function zwiazCzynnosciAplikacji(
  kanal: Kanal,
  korzen: Element,
  idOkna: () => string,
  odswiez: () => void,
  dolacz: (zdejmij: () => void) => void,
  przy: AddEventListenerOptions,
): void {
  const otoczenie: Otoczenie = { kanal, korzen, idOkna, odswiez };
  const zdejmowanie = [
    dolozCzynnosciPanelu(korzen, 'panel-builder', 'Czynności wyrobu', [
      {
        naglowek: 'Wyrób',
        pozycje: [
          { kod: 'wyrob-zapisz', nazwa: 'Zapisz wyrób…' },
          { kod: 'wyrob-odczytaj', nazwa: 'Odczytaj wyrób' },
          { kod: 'motyw-odczytaj', nazwa: 'Odczytaj motyw' },
          { kod: 'motyw-ustaw', nazwa: 'Ustaw motyw…' },
        ],
      },
      {
        naglowek: 'Plan budowy',
        pozycje: [
          { kod: 'etap', nazwa: 'Dołóż etap…' },
          { kod: 'kamien', nazwa: 'Dołóż kamień milowy…' },
          { kod: 'kamien-usun', nazwa: 'Usuń kamień milowy…' },
          { kod: 'dzieje', nazwa: 'Odczytaj dzieje okna' },
        ],
      },
      {
        naglowek: 'Warstwy i podgląd',
        pozycje: [
          { kod: 'warstwa-odczytaj', nazwa: 'Odczytaj pliki warstwy…' },
          { kod: 'warstwa-zapisz', nazwa: 'Zapisz plik warstwy…' },
          { kod: 'schemat', nazwa: 'Odczytaj schemat danych' },
          { kod: 'podglad', nazwa: 'Uruchom podgląd…' },
          { kod: 'podglad-stop', nazwa: 'Zatrzymaj podgląd' },
          { kod: 'wydajnosc', nazwa: 'Zbadaj wydajność…' },
          { kod: 'sonda', nazwa: 'Sonduj punkt styku…' },
        ],
      },
    ], (kod) => {
      void wykonaj(otoczenie, kod, 'panel-builder');
    }, przy),
    dolozCzynnosciPanelu(korzen, 'panel-architektura', 'Czynności architektury', [
      {
        naglowek: 'Architektura',
        pozycje: [
          { kod: 'arch-zaloz', nazwa: 'Załóż architekturę…' },
          { kod: 'arch-odczytaj', nazwa: 'Odczytaj architekturę' },
          { kod: 'arch-sprawdz', nazwa: 'Sprawdź architekturę' },
          { kod: 'arch-wydaj', nazwa: 'Wydaj architekturę…' },
          { kod: 'arch-przypis', nazwa: 'Dołóż przypis…' },
        ],
      },
    ], (kod) => {
      void wykonaj(otoczenie, kod, 'panel-architektura');
    }, przy),
    dolozCzynnosciPanelu(korzen, 'panel-deployment', 'Czynności wydania', [
      {
        naglowek: 'Pakiet',
        pozycje: [
          { kod: 'manifest', nazwa: 'Zapisz manifest…' },
          { kod: 'pakiet', nazwa: 'Zbuduj pakiet' },
          { kod: 'pakiet-sprawdz', nazwa: 'Sprawdź pakiet' },
          { kod: 'pakiet-podpisz', nazwa: 'Podpisz pakiet…' },
          { kod: 'pakiet-oglos', nazwa: 'Ogłoś pakiet…' },
        ],
      },
      {
        naglowek: 'Wdrożenie',
        pozycje: [
          { kod: 'wdroz', nazwa: 'Wdroż na środowisko…' },
          { kod: 'domena', nazwa: 'Ustaw domenę…' },
          { kod: 'skala', nazwa: 'Ustaw skalę…' },
          { kod: 'zdrowie', nazwa: 'Odczytaj zdrowie…' },
          { kod: 'dziennik-wdrozenia', nazwa: 'Dziennik ostatniego wdrożenia' },
          { kod: 'dziennik-uslugi', nazwa: 'Dziennik usługi' },
        ],
      },
      {
        naglowek: 'Zmienne środowiska',
        pozycje: [
          { kod: 'zmienne', nazwa: 'Odczytaj zmienne…' },
          { kod: 'zmienna-ustaw', nazwa: 'Ustaw zmienną…' },
        ],
      },
    ], (kod) => {
      void wykonaj(otoczenie, kod, 'panel-deployment');
    }, przy),
  ];
  for (const zdejmij of zdejmowanie) if (zdejmij !== null) dolacz(zdejmij);
}

/* Oznaczenie pakietu żyje w oknie, bo kontrakt nie ma wykazu pakietów: bierze
   się z budowy albo zapisu manifestu i służy podpisowi, kontroli i ogłoszeniu. */
const PAKIETY = new Map<string, string>();

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
  if (otoczenie.idOkna() === '') {
    oglos(NAGLOWEK, 'Rdzeń nie dał okna aplikacji dla tej karty.', 'ostrzezenie');
    return;
  }
  if (kod === 'wyrob-zapisz') return zapiszWyrob(otoczenie, panel);
  if (kod === 'wyrob-odczytaj') return odczytajWyrob(otoczenie);
  if (kod === 'motyw-odczytaj') return odczytajMotyw(otoczenie);
  if (kod === 'motyw-ustaw') return ustawMotyw(otoczenie, panel);
  if (kod === 'etap') return dolozEtap(otoczenie, panel);
  if (kod === 'kamien') return dolozKamien(otoczenie, panel);
  if (kod === 'kamien-usun') return usunKamien(otoczenie, panel);
  if (kod === 'dzieje') return odczytajDzieje(otoczenie);
  if (kod === 'warstwa-odczytaj') return odczytajWarstwe(otoczenie, panel);
  if (kod === 'warstwa-zapisz') return zapiszPlikWarstwy(otoczenie, panel);
  if (kod === 'schemat') return odczytajSchemat(otoczenie);
  if (kod === 'podglad') return uruchomPodglad(otoczenie, panel);
  if (kod === 'podglad-stop') return zatrzymajPodglad(otoczenie);
  if (kod === 'wydajnosc') return zbadajWydajnosc(otoczenie, panel);
  if (kod === 'sonda') return sondujPunktStyku(otoczenie, panel);
  if (kod === 'arch-zaloz') return zalozArchitekture(otoczenie, panel);
  if (kod === 'arch-odczytaj') return odczytajArchitekture(otoczenie);
  if (kod === 'arch-sprawdz') return sprawdzArchitekture(otoczenie);
  if (kod === 'arch-wydaj') return wydajArchitekture(otoczenie, panel);
  if (kod === 'arch-przypis') return dolozPrzypis(otoczenie, panel);
  if (kod === 'manifest') return zapiszManifest(otoczenie, panel);
  if (kod === 'pakiet') return zbudujPakiet(otoczenie);
  if (kod === 'pakiet-sprawdz') return sprawdzPakiet(otoczenie);
  if (kod === 'pakiet-podpisz') return podpiszPakiet(otoczenie, panel);
  if (kod === 'pakiet-oglos') return oglosPakiet(otoczenie, panel);
  if (kod === 'wdroz') return wdroz(otoczenie, panel);
  if (kod === 'domena') return ustawDomene(otoczenie, panel);
  if (kod === 'skala') return ustawSkale(otoczenie, panel);
  if (kod === 'zdrowie') return odczytajZdrowie(otoczenie, panel);
  if (kod === 'dziennik-wdrozenia') return dziennikWdrozenia(otoczenie);
  if (kod === 'dziennik-uslugi') return dziennikUslugi(otoczenie);
  if (kod === 'zmienne') return odczytajZmienne(otoczenie, panel);
  if (kod === 'zmienna-ustaw') return ustawZmienna(otoczenie, panel);
  /* Czynność spoza obsłużonych odmawia zamiast milczeć: pozycja bez gałęzi
     wyglądałaby jak działająca. */
  oglos(NAGLOWEK, `Czynność „${kod}" nie jest prowadzona przez to okno.`, 'ostrzezenie');
}

async function zapiszWyrob({ kanal, korzen, idOkna, odswiez }: Otoczenie, panel: string):
Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Wyrób okna',
    pola: [
      { klucz: 'nazwa', etykieta: 'Nazwa wyrobu', wymagane: true },
      { klucz: 'opis', etykieta: 'Opis', obszerne: true },
      { klucz: 'repozytorium', etykieta: 'Adres repozytorium' },
    ],
    wykonanie: 'Zapisz wyrób',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.AppsProductSave, {
    windowId: idOkna(),
    name: wartosci.nazwa ?? '',
    ...(wartosci.opis === '' ? {} : { description: wartosci.opis }),
    ...(wartosci.repozytorium === '' ? {} : { repositoryUrl: wartosci.repozytorium }),
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu wyrobu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Wyrób „${wynik.wynik.product.name}" zapisany.`);
  odswiez();
}

async function odczytajWyrob({ kanal, idOkna }: Otoczenie): Promise<void> {
  const wynik = await wywolaj(kanal, Command.AppsProductGet, { windowId: idOkna() });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu wyrobu.', 'ostrzezenie');
    return;
  }
  const wyrob = wynik.wynik.product;
  oglos(NAGLOWEK, wyrob === undefined
    ? 'To okno nie ma jeszcze zapisanego wyrobu.'
    : `Wyrób „${wyrob.name}": ${wyrob.description ?? 'bez opisu'}.`);
}

async function odczytajMotyw({ kanal, idOkna }: Otoczenie): Promise<void> {
  const wynik = await wywolaj(kanal, Command.AppsThemeGet, { windowId: idOkna() });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu motywu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, wynik.wynik.theme === undefined
    ? 'To okno nie ma jeszcze ustawionego motywu.'
    : JSON.stringify(wynik.wynik.theme).slice(0, 300));
}

/* Motyw jest dla rdzenia treścią nieprzejrzystą, więc pole podaje go zapisem
   JSON; zapis niepoprawny odmawia zamiast wysyłać rdzeniowi tekst. */
async function ustawMotyw({ kanal, korzen, idOkna }: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Motyw wyrobu',
    pola: [{
      klucz: 'motyw',
      etykieta: 'Zapis JSON motywu',
      podpowiedz: '{"tlo":"#101014"}',
      obszerne: true,
      wymagane: true,
    }],
    wykonanie: 'Ustaw motyw',
  });
  if (wartosci === null) return;
  let motyw: unknown;
  try {
    motyw = JSON.parse(wartosci.motyw ?? '');
  } catch {
    oglos(NAGLOWEK, 'Wpisana treść nie jest poprawnym zapisem JSON.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.AppsThemeSet, { windowId: idOkna(), theme: motyw });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił ustawienia motywu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Motyw ustawiony.');
}

async function dolozEtap({ kanal, korzen, idOkna, odswiez }: Otoczenie, panel: string):
Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Etap planu',
    pola: [{ klucz: 'nazwa', etykieta: 'Nazwa etapu', wymagane: true }],
    wykonanie: 'Dołóż etap',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.AppsStageSave, {
    windowId: idOkna(),
    name: wartosci.nazwa ?? '',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu etapu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Etap „${wartosci.nazwa ?? ''}" stoi w planie.`);
  odswiez();
}

async function dolozKamien({ kanal, korzen, idOkna, odswiez }: Otoczenie, panel: string):
Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Kamień milowy',
    pola: [{ klucz: 'nazwa', etykieta: 'Nazwa kamienia', wymagane: true }],
    wykonanie: 'Dołóż kamień',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.AppsMilestoneSave, {
    windowId: idOkna(),
    name: wartosci.nazwa ?? '',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu kamienia.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Kamień „${wartosci.nazwa ?? ''}" stoi w planie.`);
  odswiez();
}

async function usunKamien({ kanal, korzen, idOkna, odswiez }: Otoczenie, panel: string):
Promise<void> {
  const wykaz = await wywolaj(kanal, Command.AppsMilestoneList, { windowId: idOkna() });
  const kamienie = wykaz.wynik?.milestones ?? [];
  if (kamienie.length === 0) {
    oglos(NAGLOWEK, 'Plan nie ma jeszcze kamieni milowych.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Usunięcie kamienia',
    pola: [{
      klucz: 'kamien',
      etykieta: 'Kamień milowy',
      wybor: kamienie.map((kamien) => [kamien.id, `${kamien.name} · ${kamien.status}`] as const),
    }],
    wykonanie: 'Usuń kamień',
    nieodwracalne: 'Kamień znika z planu wraz z powiązaniem z etapami.',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.AppsMilestoneDelete, {
    windowId: idOkna(),
    milestoneId: wartosci.kamien ?? '',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił usunięcia kamienia.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Kamień zdjęty z planu.');
  odswiez();
}

async function odczytajDzieje({ kanal, korzen, idOkna }: Otoczenie): Promise<void> {
  const wynik = await wywolaj(kanal, Command.AppsTimelineList, { windowId: idOkna(), limit: 100 });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu dziejów.', 'ostrzezenie');
    return;
  }
  const wpisy = wynik.wynik.entries;
  wypelnij(korzen, 'panel-plan', wpisy.map((wpis: AppTimelineEntry) =>
    `${wpis.kind} · ${wpis.summary}`));
  if (wpisy.length === 0) oglos(NAGLOWEK, 'Dzieje tego okna są jeszcze puste.');
}

async function odczytajWarstwe({ kanal, korzen, idOkna }: Otoczenie, panel: string):
Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Pliki warstwy',
    pola: [{ klucz: 'warstwa', etykieta: 'Warstwa', wybor: WARSTWY }],
    wykonanie: 'Odczytaj pliki',
  });
  if (wartosci === null) return;
  const warstwa = (wartosci.warstwa ?? AppWorkspaceLayer.Frontend) as AppWorkspaceLayer;
  const wynik = await wywolaj(kanal, Command.AppsWorkspaceList, {
    windowId: idOkna(),
    layer: warstwa,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu warstwy.', 'ostrzezenie');
    return;
  }
  const pliki = wynik.wynik.files;
  wypelnij(korzen, warstwa === AppWorkspaceLayer.Frontend ? 'panel-frontend' : 'panel-backend',
    pliki.map((plik: DeveloperFile) => plik.path));
  if (pliki.length === 0) oglos(NAGLOWEK, 'Ta warstwa nie ma jeszcze plików.');
}

async function zapiszPlikWarstwy({ kanal, korzen, idOkna }: Otoczenie, panel: string):
Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Plik warstwy',
    pola: [
      { klucz: 'warstwa', etykieta: 'Warstwa', wybor: WARSTWY },
      { klucz: 'sciezka', etykieta: 'Ścieżka pliku', podpowiedz: 'src/widok.ts', wymagane: true },
      { klucz: 'tresc', etykieta: 'Treść', obszerne: true, wymagane: true },
    ],
    wykonanie: 'Zapisz plik',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.AppsWorkspaceUpdate, {
    windowId: idOkna(),
    layer: (wartosci.warstwa ?? AppWorkspaceLayer.Frontend) as AppWorkspaceLayer,
    path: wartosci.sciezka ?? '',
    content: wartosci.tresc ?? '',
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu pliku.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Plik ${wynik.wynik.file.path} zapisany.`);
}

async function odczytajSchemat({ kanal, korzen, idOkna }: Otoczenie): Promise<void> {
  const wynik = await wywolaj(kanal, Command.AppsSchemaGet, { windowId: idOkna() });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu schematu.', 'ostrzezenie');
    return;
  }
  const tablice = wynik.wynik.schema?.tables ?? [];
  if (tablice.length === 0) {
    oglos(NAGLOWEK, 'Rdzeń nie zna tablic tej aplikacji.', 'ostrzezenie');
    return;
  }
  wypelnij(korzen, 'panel-backend', tablice.map((tablica: AppSchemaTable) =>
    `${tablica.name} · ${tablica.columns.length} kolumn`));
}

async function uruchomPodglad({ kanal, korzen, idOkna }: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Podgląd warstwy',
    pola: [{ klucz: 'warstwa', etykieta: 'Warstwa', wybor: WARSTWY }],
    wykonanie: 'Uruchom podgląd',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.AppsPreviewStart, {
    windowId: idOkna(),
    layer: (wartosci.warstwa ?? AppWorkspaceLayer.Frontend) as AppWorkspaceLayer,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił uruchomienia podglądu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Podgląd stoi pod ${wynik.wynik.previewUrl}.`);
}

async function zatrzymajPodglad({ kanal, idOkna }: Otoczenie): Promise<void> {
  const wynik = await wywolaj(kanal, Command.AppsPreviewStop, { windowId: idOkna() });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zatrzymania podglądu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Podgląd zatrzymany.');
}

async function zbadajWydajnosc({ kanal, korzen, idOkna }: Otoczenie, panel: string):
Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Badanie wydajności',
    pola: [{
      klucz: 'adres',
      etykieta: 'Adres strony',
      podpowiedz: 'http://127.0.0.1:3000',
      wymagane: true,
    }],
    wykonanie: 'Zbadaj wydajność',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.AppsPerformanceAudit, {
    windowId: idOkna(),
    url: wartosci.adres ?? '',
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił badania wydajności.', 'ostrzezenie');
    return;
  }
  const badanie = wynik.wynik.audit;
  oglos(NAGLOWEK,
    `Wydajność ${badanie.url}: ocena ${badanie.performanceScore}, miar ${badanie.metrics.length}.`);
}

async function sondujPunktStyku({ kanal, korzen, idOkna }: Otoczenie, panel: string):
Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Sonda punktu styku',
    pola: [
      {
        klucz: 'metoda',
        etykieta: 'Metoda',
        wybor: [
          [AppEndpointMethod.Get, 'GET'],
          [AppEndpointMethod.Post, 'POST'],
          [AppEndpointMethod.Put, 'PUT'],
          [AppEndpointMethod.Patch, 'PATCH'],
          [AppEndpointMethod.Delete, 'DELETE'],
        ],
      },
      { klucz: 'sciezka', etykieta: 'Ścieżka', podpowiedz: '/zdrowie', wymagane: true },
      { klucz: 'srodowisko', etykieta: 'Środowisko', wybor: SRODOWISKA },
    ],
    wykonanie: 'Wyślij sondę',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.AppsEndpointProbe, {
    windowId: idOkna(),
    method: (wartosci.metoda ?? AppEndpointMethod.Get) as AppEndpointMethod,
    path: wartosci.sciezka ?? '',
    environment: (wartosci.srodowisko ?? AppDeployEnvironment.Dev) as AppDeployEnvironment,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wysłania sondy.', 'ostrzezenie');
    return;
  }
  const odpowiedz = wynik.wynik.result;
  oglos(NAGLOWEK, `Sonda: kod ${odpowiedz.statusCode}, czas ${odpowiedz.durationMs} ms.`);
}

async function zalozArchitekture({ kanal, korzen, idOkna, odswiez }: Otoczenie, panel: string):
Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Architektura aplikacji',
    pola: [
      { klucz: 'nazwa', etykieta: 'Nazwa architektury', wymagane: true },
      {
        klucz: 'wzorzec',
        etykieta: 'Wzorzec',
        wybor: [
          [AppArchitectureTemplate.Monolith, 'Monolit'],
          [AppArchitectureTemplate.Microservices, 'Usługi rozłączne'],
          [AppArchitectureTemplate.Serverless, 'Bez własnego serwera'],
        ],
      },
    ],
    wykonanie: 'Załóż architekturę',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.AppsArchitectureDefine, {
    windowId: idOkna(),
    name: wartosci.nazwa ?? '',
    template: (wartosci.wzorzec ?? AppArchitectureTemplate.Monolith) as AppArchitectureTemplate,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił założenia architektury.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Architektura „${wartosci.nazwa ?? ''}" założona.`);
  odswiez();
}

async function odczytajArchitekture({ kanal, korzen, idOkna }: Otoczenie): Promise<void> {
  const wynik = await wywolaj(kanal, Command.AppsArchitectureGet, { windowId: idOkna() });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu architektury.', 'ostrzezenie');
    return;
  }
  const architektura = wynik.wynik.architecture;
  if (architektura === undefined) {
    oglos(NAGLOWEK, 'To okno nie ma jeszcze założonej architektury.', 'ostrzezenie');
    return;
  }
  wypelnij(korzen, 'panel-architektura',
    (architektura.components ?? []).map((skladnik) => `${skladnik.name} · ${skladnik.kind}`));
  oglos(NAGLOWEK,
    `Architektura „${architektura.name ?? 'bez nazwy'}" na wzorcu ${architektura.template}.`);
}

async function sprawdzArchitekture({ kanal, korzen, idOkna }: Otoczenie): Promise<void> {
  const wynik = await wywolaj(kanal, Command.AppsArchitectureValidate, { windowId: idOkna() });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił sprawdzenia architektury.', 'ostrzezenie');
    return;
  }
  const uwagi = wynik.wynik.issues;
  if (uwagi.length === 0) {
    oglos(NAGLOWEK, 'Architektura bez uwag.');
    return;
  }
  wypelnij(korzen, 'panel-architektura',
    uwagi.map((uwaga: AppValidationIssue) => `${uwaga.severity} · ${uwaga.message}`));
  oglos(NAGLOWEK, `Architektura ma ${uwagi.length} uwag.`, 'ostrzezenie');
}

/* Wydanie wraca oznaczeniem wytworu magazynu, więc okno nazywa oznaczenie
   zamiast udawać zapis pliku, którego nie ma gdzie zostawić. */
async function wydajArchitekture({ kanal, korzen, idOkna }: Otoczenie, panel: string):
Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Wydanie architektury',
    pola: [{
      klucz: 'postac',
      etykieta: 'Postać',
      wybor: [
        [AppExportFormat.Mermaid, 'Mermaid'],
        [AppExportFormat.Markdown, 'Markdown'],
        [AppExportFormat.Svg, 'SVG'],
        [AppExportFormat.Png, 'PNG'],
      ],
    }],
    wykonanie: 'Wydaj architekturę',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.AppsArchitectureExport, {
    windowId: idOkna(),
    format: (wartosci.postac ?? AppExportFormat.Mermaid) as AppExportFormat,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wydania architektury.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Architektura wydana jako wytwór ${wynik.wynik.artifactRef}.`);
}

async function dolozPrzypis({ kanal, korzen, idOkna }: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Przypis architektury',
    pola: [
      { klucz: 'tresc', etykieta: 'Treść przypisu', obszerne: true, wymagane: true },
      { klucz: 'skladnik', etykieta: 'Oznaczenie składnika', podpowiedz: 'Puste znaczy całość' },
    ],
    wykonanie: 'Dołóż przypis',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.AppsArchitectureAnnotationSave, {
    windowId: idOkna(),
    text: wartosci.tresc ?? '',
    ...(wartosci.skladnik === '' ? {} : { componentId: wartosci.skladnik }),
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu przypisu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Przypis dołożony do architektury.');
}

async function zapiszManifest({ kanal, korzen, idOkna }: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Manifest pakietu',
    pola: [
      { klucz: 'oznaczenie', etykieta: 'Oznaczenie', podpowiedz: 'sklep', wymagane: true },
      { klucz: 'nazwa', etykieta: 'Nazwa', wymagane: true },
      { klucz: 'wersja', etykieta: 'Wersja', podpowiedz: '1.0.0', wymagane: true },
      {
        klucz: 'rodzaj',
        etykieta: 'Rodzaj',
        wybor: [
          [ExtensionKind.Plugin, 'Wtyczka'],
          [ExtensionKind.Mcp, 'Most MCP'],
          [ExtensionKind.Api, 'Usługa'],
          [ExtensionKind.Skill, 'Umiejętność'],
        ],
      },
    ],
    wykonanie: 'Zapisz manifest',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.AppsPackageManifestSave, {
    windowId: idOkna(),
    manifest: {
      identifier: wartosci.oznaczenie ?? '',
      name: wartosci.nazwa ?? '',
      version: wartosci.wersja ?? '',
      kind: (wartosci.rodzaj ?? ExtensionKind.Plugin) as ExtensionKind,
    },
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu manifestu.', 'ostrzezenie');
    return;
  }
  PAKIETY.set(idOkna(), wynik.wynik.package.id);
  oglos(NAGLOWEK, `Manifest zapisany; pakiet ${wynik.wynik.package.id}.`);
}

async function zbudujPakiet({ kanal, idOkna }: Otoczenie): Promise<void> {
  const wynik = await wywolaj(kanal, Command.AppsPackageBuild, {
    windowId: idOkna(),
    format: AppPackageFormat.Zip,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zbudowania pakietu.', 'ostrzezenie');
    return;
  }
  PAKIETY.set(idOkna(), wynik.wynik.package.id);
  oglos(NAGLOWEK, `Pakiet ${wynik.wynik.package.id} zbudowany jako wytwór `
    + `${wynik.wynik.package.artifactRef ?? 'bez oznaczenia'}.`);
}

function pakiet(idOkna: string): string {
  const zapamietany = PAKIETY.get(idOkna) ?? '';
  if (zapamietany === '') {
    oglos(NAGLOWEK, 'Okno nie zna jeszcze pakietu — zbuduj go albo zapisz manifest.', 'ostrzezenie');
  }
  return zapamietany;
}

async function sprawdzPakiet({ kanal, korzen, idOkna }: Otoczenie): Promise<void> {
  const cel = pakiet(idOkna());
  if (cel === '') return;
  const wynik = await wywolaj(kanal, Command.AppsPackageValidate, {
    windowId: idOkna(),
    packageId: cel,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił sprawdzenia pakietu.', 'ostrzezenie');
    return;
  }
  const uwagi = wynik.wynik.issues;
  if (uwagi.length === 0) {
    oglos(NAGLOWEK, 'Pakiet bez uwag.');
    return;
  }
  wypelnij(korzen, 'panel-deployment',
    uwagi.map((uwaga: AppValidationIssue) => `${uwaga.severity} · ${uwaga.message}`));
  oglos(NAGLOWEK, `Pakiet ma ${uwagi.length} uwag.`, 'ostrzezenie');
}

/* Podpis idzie odwołaniem do klucza w sejfie rdzenia; materiał klucza nie
   przechodzi przez okno ani przez ogłoszenie. */
async function podpiszPakiet({ kanal, korzen, idOkna }: Otoczenie, panel: string): Promise<void> {
  const cel = pakiet(idOkna());
  if (cel === '') return;
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Podpis pakietu',
    opis: 'Klucz zostaje w sejfie rdzenia; okno podaje odwołanie, nie materiał.',
    pola: [{ klucz: 'klucz', etykieta: 'Odwołanie do klucza w sejfie', wymagane: true }],
    wykonanie: 'Podpisz pakiet',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.AppsPackageSign, {
    windowId: idOkna(),
    packageId: cel,
    signingKeyRef: wartosci.klucz ?? '',
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił podpisania pakietu.', 'ostrzezenie');
    return;
  }
  const podpis = wynik.wynik.signature;
  oglos(NAGLOWEK, podpis.signed
    ? `Pakiet podpisany na poziomie zaufania ${podpis.trustLevel}.`
    : `Rdzeń nie podpisał pakietu: ${podpis.detail ?? 'bez wyjaśnienia'}.`,
  podpis.signed ? 'informacja' : 'ostrzezenie');
}

async function oglosPakiet({ kanal, korzen, idOkna }: Otoczenie, panel: string): Promise<void> {
  const cel = pakiet(idOkna());
  if (cel === '') return;
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Ogłoszenie pakietu',
    pola: [{ klucz: 'noty', etykieta: 'Noty wydania', obszerne: true }],
    wykonanie: 'Ogłoś pakiet',
    nieodwracalne: 'Pakiet staje się widoczny dla innych kont; ogłoszenia nie da się cofnąć.',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.AppsPackagePublish, {
    windowId: idOkna(),
    packageId: cel,
    ...(wartosci.noty === '' ? {} : { releaseNotes: wartosci.noty }),
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił ogłoszenia pakietu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Pakiet ogłoszony jako rozszerzenie „${wynik.wynik.extension.name}".`);
}

/* Wdrożenie na produkcję zmienia to, co widzą odbiorcy usługi, więc szuflada
   nazywa następstwo i wymaga drugiego naciśnięcia. */
async function wdroz({ kanal, korzen, idOkna, odswiez }: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Wdrożenie',
    pola: [
      { klucz: 'srodowisko', etykieta: 'Środowisko', wybor: SRODOWISKA },
      { klucz: 'wersja', etykieta: 'Wersja', podpowiedz: 'Puste bierze wersję wyrobu' },
      { klucz: 'noty', etykieta: 'Noty wydania', obszerne: true },
    ],
    wykonanie: 'Wdroż',
    nieodwracalne: 'Wdrożenie zmienia to, co widzą odbiorcy wskazanego środowiska.',
  });
  if (wartosci === null) return;
  const cel = (wartosci.srodowisko ?? AppDeployEnvironment.Dev) as AppDeployEnvironment;
  const wynik = await wywolaj(kanal, Command.AppsDeploymentRun, {
    windowId: idOkna(),
    environment: cel,
    ...(wartosci.wersja === '' ? {} : { version: wartosci.wersja }),
    ...(wartosci.noty === '' ? {} : { releaseNotes: wartosci.noty }),
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wdrożenia.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Wdrożenie na ${cel} w stanie ${wynik.wynik.deployment.status}.`);
  odswiez();
}

async function ustawDomene({ kanal, korzen, idOkna }: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Domena środowiska',
    pola: [
      { klucz: 'srodowisko', etykieta: 'Środowisko', wybor: SRODOWISKA },
      { klucz: 'domena', etykieta: 'Adres domeny', podpowiedz: 'sklep.example.com', wymagane: true },
    ],
    wykonanie: 'Ustaw domenę',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.AppsDeploymentDomainSet, {
    windowId: idOkna(),
    environment: (wartosci.srodowisko ?? AppDeployEnvironment.Dev) as AppDeployEnvironment,
    domain: wartosci.domena ?? '',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił ustawienia domeny.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Domena ${wartosci.domena ?? ''} ustawiona.`);
}

async function ustawSkale({ kanal, korzen, idOkna }: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Skala środowiska',
    pola: [
      { klucz: 'srodowisko', etykieta: 'Środowisko', wybor: SRODOWISKA },
      { klucz: 'wystapienia', etykieta: 'Liczba wystąpień', wartosc: '1', wymagane: true },
    ],
    wykonanie: 'Ustaw skalę',
  });
  if (wartosci === null) return;
  const liczba = Number.parseInt(wartosci.wystapienia ?? '', 10);
  if (!Number.isFinite(liczba) || liczba < 1) {
    oglos(NAGLOWEK, 'Liczba wystąpień musi być liczbą nie mniejszą niż jeden.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.AppsDeploymentScaleSet, {
    windowId: idOkna(),
    environment: (wartosci.srodowisko ?? AppDeployEnvironment.Dev) as AppDeployEnvironment,
    instances: liczba,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił ustawienia skali.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Skala ${liczba} ustawiona.`);
}

async function odczytajZdrowie({ kanal, korzen, idOkna }: Otoczenie, panel: string):
Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Zdrowie środowiska',
    pola: [{ klucz: 'srodowisko', etykieta: 'Środowisko', wybor: SRODOWISKA }],
    wykonanie: 'Odczytaj zdrowie',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.AppsDeploymentHealthGet, {
    windowId: idOkna(),
    environment: (wartosci.srodowisko ?? AppDeployEnvironment.Dev) as AppDeployEnvironment,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu zdrowia.', 'ostrzezenie');
    return;
  }
  const zdrowie = wynik.wynik.health;
  oglos(NAGLOWEK, `Środowisko ${zdrowie.environment}: `
    + `${zdrowie.available ? 'dostępne' : 'niedostępne'}, `
    + `dostępność ${zdrowie.availabilityPercent ?? 'nieznana'}.`,
  zdrowie.available ? 'informacja' : 'ostrzezenie');
}

/* Dziennik czytany jest z wdrożenia stojącego w wykazie najwyżej: okno nie
   prowadzi wskazania wdrożenia, a wykaz podaje ich kolejność. */
async function dziennikWdrozenia({ kanal, korzen, idOkna }: Otoczenie): Promise<void> {
  const wykaz = await wywolaj(kanal, Command.AppsDeploymentList, { windowId: idOkna() });
  const cel = wykaz.wynik?.deployments[0]?.id ?? '';
  if (cel === '') {
    oglos(NAGLOWEK, 'To okno nie ma jeszcze żadnego wdrożenia.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.AppsDeploymentLogRead, {
    windowId: idOkna(),
    deploymentId: cel,
    limit: 200,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu dziennika.', 'ostrzezenie');
    return;
  }
  wypelnij(korzen, 'panel-deployment', wynik.wynik.lines);
  if (wynik.wynik.lines.length === 0) oglos(NAGLOWEK, 'Dziennik tego wdrożenia jest pusty.');
}

async function dziennikUslugi({ kanal, korzen, idOkna }: Otoczenie): Promise<void> {
  const wynik = await wywolaj(kanal, Command.AppsServiceLogRead, {
    windowId: idOkna(),
    limit: 200,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu dziennika usługi.', 'ostrzezenie');
    return;
  }
  wypelnij(korzen, 'panel-backend', wynik.wynik.lines);
  if (wynik.wynik.lines.length === 0) oglos(NAGLOWEK, 'Dziennik usługi jest pusty.');
}

/* Wartości tajne rdzeń oddaje odwołaniem do sejfu, nie treścią; wykaz nazywa
   odwołanie, żeby nie udawał, że zna wartość. */
async function odczytajZmienne({ kanal, korzen, idOkna }: Otoczenie, panel: string):
Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Zmienne środowiska',
    pola: [{ klucz: 'srodowisko', etykieta: 'Środowisko', wybor: SRODOWISKA }],
    wykonanie: 'Odczytaj zmienne',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.AppsEnvironmentVariableList, {
    windowId: idOkna(),
    environment: (wartosci.srodowisko ?? AppDeployEnvironment.Dev) as AppDeployEnvironment,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu zmiennych.', 'ostrzezenie');
    return;
  }
  const zmienne = wynik.wynik.variables;
  wypelnij(korzen, 'panel-terminal', zmienne.map((zmienna: AppEnvironmentVariable) =>
    `${zmienna.name} = ${zmienna.value ?? `sejf ${zmienna.secretRef ?? 'bez odwołania'}`}`));
  if (zmienne.length === 0) oglos(NAGLOWEK, 'To środowisko nie ma jeszcze zmiennych.');
}

async function ustawZmienna({ kanal, korzen, idOkna }: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Zmienna środowiska',
    pola: [
      { klucz: 'srodowisko', etykieta: 'Środowisko', wybor: SRODOWISKA },
      { klucz: 'nazwa', etykieta: 'Nazwa zmiennej', podpowiedz: 'PORT', wymagane: true },
      { klucz: 'wartosc', etykieta: 'Wartość jawna' },
      { klucz: 'sejf', etykieta: 'Odwołanie do sejfu', podpowiedz: 'Zamiast wartości jawnej' },
    ],
    wykonanie: 'Ustaw zmienną',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.AppsEnvironmentVariableSet, {
    windowId: idOkna(),
    environment: (wartosci.srodowisko ?? AppDeployEnvironment.Dev) as AppDeployEnvironment,
    name: wartosci.nazwa ?? '',
    ...(wartosci.sejf === '' ? { value: wartosci.wartosc } : { secretRef: wartosci.sejf }),
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił ustawienia zmiennej.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Zmienna ${wartosci.nazwa ?? ''} ustawiona.`);
}
