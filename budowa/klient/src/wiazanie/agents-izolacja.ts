// Izolacja w oknie Agents: co sesja dzieli z resztą aplikacji, a co ma
// osobne — historia, pamięć, katalog roboczy, sieć, żeton konta.
import {
  Command,
  ConfigScope,
  IsolationContextKind,
  IsolationLayer,
  IsolationTechnicalScope,
} from '../../../shared/contract.ts';
import type { IsolationProfile } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { dolozCzynnosciPanelu, zapytajWSzufladzie } from './czynnosci-okna.ts';
import { oglos } from './ogloszenie.ts';

const NAGLOWEK = 'Izolacja';

const ZAKRESY: ReadonlyArray<readonly [string, string]> = [
  [ConfigScope.Application, 'Cała aplikacja'],
  [ConfigScope.Environment, 'Środowisko'],
  [ConfigScope.Project, 'Projekt'],
  [ConfigScope.Session, 'Sesja'],
  [ConfigScope.Window, 'Okno'],
];

const WARSTWY: ReadonlyArray<readonly [string, string]> = [
  [IsolationLayer.Default, 'Nastawa domyślna'],
  [IsolationLayer.Session, 'Nastawa sesji'],
];

const RODZAJE: ReadonlyArray<readonly [string, string]> = [
  [IsolationContextKind.History, 'Dzieje rozmów'],
  [IsolationContextKind.Memory, 'Pamięć'],
  [IsolationContextKind.Context, 'Kontekst pracy'],
];

const ZAKRESY_TECHNICZNE: ReadonlyArray<readonly [string, string]> = [
  [IsolationTechnicalScope.WorkingDirectory, 'Katalog roboczy'],
  [IsolationTechnicalScope.ProcessEnvironment, 'Środowisko procesu'],
  [IsolationTechnicalScope.ModelDataDirectory, 'Katalog danych modelu'],
  [IsolationTechnicalScope.NetworkAccess, 'Dostęp do sieci'],
  [IsolationTechnicalScope.FileAccess, 'Dostęp do plików'],
  [IsolationTechnicalScope.AccountToken, 'Żeton konta'],
  [IsolationTechnicalScope.ProcessModel, 'Proces modelu'],
  [IsolationTechnicalScope.ExecutionServer, 'Serwer wykonawczy'],
];

interface Otoczenie {
  kanal: Kanal;
  korzen: Element;
  idSesji: () => string;
}

export function zwiazIzolacjeAgentow(
  kanal: Kanal,
  korzen: Element,
  idSesji: () => string,
  dolacz: (zdejmij: () => void) => void,
  przy: AddEventListenerOptions,
): void {
  const otoczenie: Otoczenie = { kanal, korzen, idSesji };
  const zdejmij = dolozCzynnosciPanelu(korzen, 'panel-skills', 'Czynności izolacji', [
    {
      naglowek: 'Wgląd',
      pozycje: [
        { kod: 'poziomy', nazwa: 'Poziomy izolacji' },
        { kod: 'obowiazuje', nazwa: 'Co obowiązuje tu i teraz…' },
        { kod: 'kontekst', nazwa: 'Odczytaj rozdział treści…' },
        { kod: 'techniczne', nazwa: 'Odczytaj rozdział techniczny…' },
      ],
    },
    {
      naglowek: 'Nastawy',
      pozycje: [
        { kod: 'kontekst-ustaw', nazwa: 'Rozdziel treść…' },
        { kod: 'techniczne-ustaw', nazwa: 'Rozdziel warunki techniczne…' },
        { kod: 'warstwa', nazwa: 'Przestaw warstwę sesji…' },
      ],
    },
    {
      naglowek: 'Wzorce izolacji',
      pozycje: [
        { kod: 'wzorce', nazwa: 'Wykaz wzorców' },
        { kod: 'wzorzec-zapisz', nazwa: 'Zapisz wzorzec…' },
        { kod: 'wzorzec-wczytaj', nazwa: 'Wczytaj wzorzec…' },
        { kod: 'wzorzec-przypisz', nazwa: 'Przypisz wzorzec…' },
        { kod: 'wzorzec-usun', nazwa: 'Usuń wzorzec…' },
      ],
    },
  ], (kod) => {
    void wykonaj(otoczenie, kod, 'panel-skills');
  }, przy);
  if (zdejmij !== null) dolacz(zdejmij);
}

function wypelnij(korzen: Element, wiersze: string[]): void {
  const cialo = korzen.querySelector('#panel-skills .sta-okno-tresc');
  if (cialo === null) return;
  cialo.replaceChildren(...wiersze.map((tresc) => {
    const wiersz = cialo.ownerDocument.createElement('div');
    wiersz.className = 'dn-wykaz-modulu-poz';
    wiersz.textContent = tresc;
    return wiersz;
  }));
}

async function wykonaj(otoczenie: Otoczenie, kod: string, panel: string): Promise<void> {
  if (kod === 'poziomy') return poziomyIzolacji(otoczenie);
  if (kod === 'obowiazuje') return coObowiazuje(otoczenie, panel);
  if (kod === 'kontekst') return odczytajRozdzialTresci(otoczenie, panel);
  if (kod === 'techniczne') return odczytajRozdzialTechniczny(otoczenie, panel);
  if (kod === 'kontekst-ustaw') return rozdzielTresc(otoczenie, panel);
  if (kod === 'techniczne-ustaw') return rozdzielTechniczne(otoczenie, panel);
  if (kod === 'warstwa') return przestawWarstwe(otoczenie, panel);
  if (kod === 'wzorce') return wykazWzorcow(otoczenie);
  if (kod === 'wzorzec-zapisz') return zapiszWzorzec(otoczenie, panel);
  if (kod === 'wzorzec-wczytaj') return wczytajWzorzec(otoczenie, panel);
  if (kod === 'wzorzec-przypisz') return przypiszWzorzec(otoczenie, panel);
  if (kod === 'wzorzec-usun') return usunWzorzec(otoczenie, panel);
  /* Czynność spoza obsłużonych odmawia zamiast milczeć: pozycja bez gałęzi
     wyglądałaby jak działająca. */
  oglos(NAGLOWEK, `Czynność „${kod}" nie jest prowadzona przez to okno.`, 'ostrzezenie');
}

/* Poziomy idą od najszerszego do najwęższego: nastawa węższa bierze górę nad
   szerszą, więc kolejność jest tu treścią, nie ozdobą. */
async function poziomyIzolacji({ kanal, korzen }: Otoczenie): Promise<void> {
  const wynik = await wywolaj(kanal, Command.IsolationScopeList, {});
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu poziomów.', 'ostrzezenie');
    return;
  }
  const poziomy = [...wynik.wynik.scopes].sort((jeden, drugi) => jeden.order - drugi.order);
  wypelnij(korzen, poziomy.map((poziom) =>
    `${poziom.label}${poziom.narrowest === true ? ' · najwęższy' : ''}`
    + `${poziom.description === undefined ? '' : ` · ${poziom.description}`}`));
  if (poziomy.length === 0) oglos(NAGLOWEK, 'Rdzeń nie podaje poziomów izolacji.');
}

async function coObowiazuje({ kanal, korzen, idSesji }: Otoczenie, panel: string):
Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Izolacja obowiązująca',
    opis: 'Rdzeń złoży nastawy wszystkich poziomów i pokaże, co z nich wychodzi.',
    pola: [{ klucz: 'warstwa', etykieta: 'Warstwa', wybor: WARSTWY }],
    wykonanie: 'Pokaż, co obowiązuje',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.IsolationPolicyPreview, {
    ...(idSesji() === '' ? {} : { sessionId: idSesji() }),
    layer: (wartosci.warstwa ?? IsolationLayer.Session) as IsolationLayer,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił złożenia nastaw.', 'ostrzezenie');
    return;
  }
  const nastawa = wynik.wynik.policy;
  wypelnij(korzen, [
    `poziom rozstrzygający: ${nastawa.origin ?? nastawa.scope}`,
    ...nastawa.contextSwitches.map((przelacznik) =>
      `${przelacznik.kind}: ${przelacznik.isolated ? 'osobne' : 'wspólne'}`),
    ...nastawa.technicalSwitches.map((przelacznik) =>
      `${przelacznik.scope}: ${przelacznik.isolated ? 'osobne' : 'wspólne'}`),
  ]);
}

async function odczytajRozdzialTresci({ kanal, korzen }: Otoczenie, panel: string):
Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Rozdział treści',
    pola: [
      { klucz: 'zakres', etykieta: 'Poziom', wybor: ZAKRESY },
      { klucz: 'warstwa', etykieta: 'Warstwa', wybor: WARSTWY },
    ],
    wykonanie: 'Odczytaj nastawę',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.IsolationContextGet, {
    scope: (wartosci.zakres ?? ConfigScope.Session) as ConfigScope,
    layer: (wartosci.warstwa ?? IsolationLayer.Session) as IsolationLayer,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu nastawy.', 'ostrzezenie');
    return;
  }
  wypelnij(korzen, wynik.wynik.switches.map((przelacznik) =>
    `${przelacznik.kind}: ${przelacznik.isolated ? 'osobne' : 'wspólne'}`));
  if (wynik.wynik.switches.length === 0) oglos(NAGLOWEK, 'Ten poziom nie ma własnej nastawy.');
}

async function odczytajRozdzialTechniczny({ kanal, korzen }: Otoczenie, panel: string):
Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Rozdział techniczny',
    pola: [
      { klucz: 'zakres', etykieta: 'Poziom', wybor: ZAKRESY },
      { klucz: 'warstwa', etykieta: 'Warstwa', wybor: WARSTWY },
    ],
    wykonanie: 'Odczytaj nastawę',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.IsolationTechnicalGet, {
    scope: (wartosci.zakres ?? ConfigScope.Session) as ConfigScope,
    layer: (wartosci.warstwa ?? IsolationLayer.Session) as IsolationLayer,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu nastawy.', 'ostrzezenie');
    return;
  }
  wypelnij(korzen, wynik.wynik.switches.map((przelacznik) =>
    `${przelacznik.scope}: ${przelacznik.isolated ? 'osobne' : 'wspólne'}`));
  if (wynik.wynik.switches.length === 0) oglos(NAGLOWEK, 'Ten poziom nie ma własnej nastawy.');
}

/* Nastawa idzie jednym rodzajem naraz: komenda przyjmuje wykaz, a okno nie
   prowadzi zaznaczania wielu przełączników — mówi to wprost w opisie. */
async function rozdzielTresc({ kanal, korzen }: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Rozdzielenie treści',
    opis: 'Nastawa obejmuje jeden rodzaj treści; pozostałe zostają, jak stały.',
    pola: [
      { klucz: 'zakres', etykieta: 'Poziom', wybor: ZAKRESY },
      { klucz: 'warstwa', etykieta: 'Warstwa', wybor: WARSTWY },
      { klucz: 'rodzaj', etykieta: 'Rodzaj treści', wybor: RODZAJE },
      {
        klucz: 'stan',
        etykieta: 'Treść tego rodzaju',
        wybor: [['tak', 'Osobna dla tego poziomu'], ['nie', 'Wspólna z szerszym poziomem']],
      },
      { klucz: 'powod', etykieta: 'Powód nastawy' },
    ],
    wykonanie: 'Ustaw rozdział',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.IsolationContextSet, {
    scope: (wartosci.zakres ?? ConfigScope.Session) as ConfigScope,
    layer: (wartosci.warstwa ?? IsolationLayer.Session) as IsolationLayer,
    switches: [{
      kind: (wartosci.rodzaj ?? IsolationContextKind.Memory) as IsolationContextKind,
      isolated: wartosci.stan === 'tak',
      ...(wartosci.powod === '' ? {} : { explanation: wartosci.powod }),
    }],
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił ustawienia rozdziału.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, wartosci.stan === 'tak'
    ? 'Ten rodzaj treści jest teraz osobny dla wskazanego poziomu.'
    : 'Ten rodzaj treści wraca do wspólnego.');
}

async function rozdzielTechniczne({ kanal, korzen }: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Rozdzielenie warunków technicznych',
    opis: 'Nastawa obejmuje jeden warunek; pozostałe zostają, jak stały.',
    pola: [
      { klucz: 'zakres', etykieta: 'Poziom', wybor: ZAKRESY },
      { klucz: 'warstwa', etykieta: 'Warstwa', wybor: WARSTWY },
      { klucz: 'warunek', etykieta: 'Warunek', wybor: ZAKRESY_TECHNICZNE },
      {
        klucz: 'stan',
        etykieta: 'Ten warunek',
        wybor: [['tak', 'Osobny dla tego poziomu'], ['nie', 'Wspólny z szerszym poziomem']],
      },
      { klucz: 'powod', etykieta: 'Powód nastawy' },
    ],
    wykonanie: 'Ustaw rozdział',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.IsolationTechnicalSet, {
    scope: (wartosci.zakres ?? ConfigScope.Session) as ConfigScope,
    layer: (wartosci.warstwa ?? IsolationLayer.Session) as IsolationLayer,
    switches: [{
      scope: (wartosci.warunek
        ?? IsolationTechnicalScope.WorkingDirectory) as IsolationTechnicalScope,
      isolated: wartosci.stan === 'tak',
      ...(wartosci.powod === '' ? {} : { explanation: wartosci.powod }),
    }],
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił ustawienia rozdziału.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, wartosci.stan === 'tak'
    ? 'Ten warunek jest teraz osobny dla wskazanego poziomu.'
    : 'Ten warunek wraca do wspólnego.');
}

async function przestawWarstwe({ kanal, korzen, idSesji }: Otoczenie, panel: string):
Promise<void> {
  if (idSesji() === '') {
    oglos(NAGLOWEK, 'Rdzeń nie dał sesji dla tej karty.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Warstwa izolacji sesji',
    opis: 'Warstwa sesji bierze górę nad nastawą domyślną, dopóki sesja trwa.',
    pola: [{ klucz: 'warstwa', etykieta: 'Warstwa', wybor: WARSTWY }],
    wykonanie: 'Przestaw warstwę',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.IsolationLayerSet, {
    layer: (wartosci.warstwa ?? IsolationLayer.Session) as IsolationLayer,
    sessionId: idSesji(),
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił przestawienia warstwy.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Warstwa izolacji przestawiona.');
}

async function wyborWzorcow(kanal: Kanal): Promise<ReadonlyArray<readonly [string, string]>> {
  const wykaz = await wywolaj(kanal, Command.IsolationProfileList, {});
  return (wykaz.wynik?.profiles ?? []).map((wzorzec: IsolationProfile) =>
    [wzorzec.id, wzorzec.name] as const);
}

async function wykazWzorcow({ kanal, korzen }: Otoczenie): Promise<void> {
  const wynik = await wywolaj(kanal, Command.IsolationProfileList, {});
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu wzorców.', 'ostrzezenie');
    return;
  }
  const wzorce = wynik.wynik.profiles;
  wypelnij(korzen, wzorce.map((wzorzec: IsolationProfile) =>
    `${wzorzec.name} · treści ${wzorzec.contextSwitches.length}`
    + ` · warunków ${wzorzec.technicalSwitches.length}`));
  if (wzorce.length === 0) oglos(NAGLOWEK, 'Nie ma jeszcze żadnego wzorca izolacji.');
}

/* Wzorzec bierze nastawę stojącą na wskazanym poziomie, bo składanie go z
   pojedynczych przełączników w szufladzie byłoby dłuższe niż sama praca. */
async function zapiszWzorzec({ kanal, korzen }: Otoczenie, panel: string): Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Wzorzec izolacji',
    opis: 'Wzorzec bierze nastawę stojącą na wskazanym poziomie.',
    pola: [
      { klucz: 'nazwa', etykieta: 'Nazwa wzorca', wymagane: true },
      { klucz: 'opis', etykieta: 'Opis', obszerne: true },
      { klucz: 'zakres', etykieta: 'Poziom, z którego brać nastawę', wybor: ZAKRESY },
    ],
    wykonanie: 'Zapisz wzorzec',
  });
  if (wartosci === null) return;
  const zakres = (wartosci.zakres ?? ConfigScope.Session) as ConfigScope;
  const tresci = await wywolaj(kanal, Command.IsolationContextGet, { scope: zakres });
  const techniczne = await wywolaj(kanal, Command.IsolationTechnicalGet, { scope: zakres });
  const wynik = await wywolaj(kanal, Command.IsolationProfileSave, {
    name: wartosci.nazwa ?? '',
    ...(wartosci.opis === '' ? {} : { description: wartosci.opis }),
    contextSwitches: tresci.wynik?.switches ?? [],
    technicalSwitches: techniczne.wynik?.switches ?? [],
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu wzorca.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Wzorzec „${wartosci.nazwa ?? ''}" zapisany.`);
  await wykazWzorcow({ kanal, korzen, idSesji: () => '' });
}

async function wczytajWzorzec({ kanal, korzen }: Otoczenie, panel: string): Promise<void> {
  const wzorce = await wyborWzorcow(kanal);
  if (wzorce.length === 0) {
    oglos(NAGLOWEK, 'Nie ma jeszcze żadnego wzorca izolacji.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Wczytanie wzorca',
    pola: [{ klucz: 'wzorzec', etykieta: 'Wzorzec', wybor: wzorce }],
    wykonanie: 'Wczytaj wzorzec',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.IsolationProfileLoad, {
    profileId: wartosci.wzorzec ?? '',
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wczytania wzorca.', 'ostrzezenie');
    return;
  }
  const wzorzec = wynik.wynik.profile;
  wypelnij(korzen, [
    `wzorzec: ${wzorzec.name}`,
    ...wzorzec.contextSwitches.map((przelacznik) =>
      `${przelacznik.kind}: ${przelacznik.isolated ? 'osobne' : 'wspólne'}`),
    ...wzorzec.technicalSwitches.map((przelacznik) =>
      `${przelacznik.scope}: ${przelacznik.isolated ? 'osobne' : 'wspólne'}`),
  ]);
}

async function przypiszWzorzec({ kanal, korzen }: Otoczenie, panel: string): Promise<void> {
  const wzorce = await wyborWzorcow(kanal);
  if (wzorce.length === 0) {
    oglos(NAGLOWEK, 'Nie ma jeszcze żadnego wzorca izolacji.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Przypisanie wzorca',
    pola: [
      { klucz: 'wzorzec', etykieta: 'Wzorzec', wybor: wzorce },
      { klucz: 'zakres', etykieta: 'Poziom', wybor: ZAKRESY },
      { klucz: 'warstwa', etykieta: 'Warstwa', wybor: WARSTWY },
    ],
    wykonanie: 'Przypisz wzorzec',
    nieodwracalne: 'Nastawy tego poziomu zostaną zastąpione nastawami wzorca.',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.IsolationProfileAssign, {
    profileId: wartosci.wzorzec ?? '',
    scope: (wartosci.zakres ?? ConfigScope.Session) as ConfigScope,
    layer: (wartosci.warstwa ?? IsolationLayer.Session) as IsolationLayer,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił przypisania wzorca.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Wzorzec przypisany do wskazanego poziomu.');
}

async function usunWzorzec({ kanal, korzen, idSesji }: Otoczenie, panel: string): Promise<void> {
  const wzorce = await wyborWzorcow(kanal);
  if (wzorce.length === 0) {
    oglos(NAGLOWEK, 'Nie ma jeszcze żadnego wzorca izolacji.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWSzufladzie(korzen, panel, {
    tytul: 'Usunięcie wzorca',
    pola: [{ klucz: 'wzorzec', etykieta: 'Wzorzec', wybor: wzorce }],
    wykonanie: 'Usuń wzorzec',
    nieodwracalne: 'Poziomy, którym wzorzec przypisano, zachowują swoje nastawy.',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.IsolationProfileDelete, {
    profileId: wartosci.wzorzec ?? '',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił usunięcia wzorca.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Wzorzec zdjęty z wykazu.');
  await wykazWzorcow({ kanal, korzen, idSesji });
}
