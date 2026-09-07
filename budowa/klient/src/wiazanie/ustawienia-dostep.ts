// Sekcja urządzeń w oknie Ustawień: katalog punktów dostępu i nadania okien —
// dołożenie, poprawa, sprawdzenie punktu oraz nadanie i zdjęcie dostępu.
import { AccessMode, AccessPointKind, Command } from '../../../shared/contract.ts';
import type { AccessGrant, AccessPoint } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import { zapytajWOknieModalnym } from './pytanie-modalne.ts';

const NAGLOWEK = 'Dostęp do maszyn';

const RODZAJE: ReadonlyArray<readonly [string, string]> = [
  [AccessPointKind.LocalDirectory, 'Katalog na urządzeniu'],
  [AccessPointKind.McpBridge, 'Maszyna przez most'],
];

const TRYBY: ReadonlyArray<readonly [string, string]> = [
  [AccessMode.Read, 'Odczyt'],
  [AccessMode.Write, 'Zapis'],
];

const STANY: Readonly<Record<string, string>> = {
  unknown: 'niesprawdzony',
  reachable: 'odpowiada',
  unreachable: 'nie odpowiada',
};

/*
zwiazDostepUstawien dokłada do sekcji urządzeń dwa katalogi: punkty dostępu
i nadania okien. Prototyp postawił w tej sekcji sam wykaz urządzeń, więc oba
katalogi powstają z tego samego słownika, którym opisana jest reszta sekcji:
pozycja z nagłówkiem i przyciskiem, pod nią ramka z tabelą.
*/
export function zwiazDostepUstawien(kanal: Kanal, korzen: Element): void {
  const sekcja = korzen.querySelector('#sekcja-urzadzenia');
  if (sekcja === null) return;
  const dokument = sekcja.ownerDocument;

  const punkty = grupa(dokument, 'us-punkty', 'Punkty dostępu',
    'Katalog maszyn i katalogów, po które sięgają okna komunikacji. Punkt sam '
    + 'z siebie niczego nie otwiera — dopiero nadanie wiąże go z oknem.',
    'Dołóż punkt', ['Punkt', 'Rodzaj', 'Korzenie', 'Stan', 'Akcje'], false);
  const nadania = grupa(dokument, 'us-nadania', 'Nadania dostępu',
    'Nadania rdzeń prowadzi osobno dla każdego okna komunikacji, więc wykaz '
    + 'zaczyna się od wskazania okna. Nadanie główne rozstrzyga, gdzie okno '
    + 'szuka najpierw.',
    'Nadaj dostęp', ['Okno', 'Punkt', 'Tryb', 'Kolejność', 'Akcje'], true);
  sekcja.append(punkty, nadania);

  void odswiez(kanal, sekcja);

  sekcja.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    if (cel.closest('#us-punkty-dodaj') !== null) {
      zdarzenie.stopPropagation();
      void dolozPunkt(kanal, sekcja);
      return;
    }
    if (cel.closest('#us-nadania-pokaz') !== null) {
      zdarzenie.stopPropagation();
      void wypelnijNadania(kanal, sekcja);
      return;
    }
    if (cel.closest('#us-nadania-dodaj') !== null) {
      zdarzenie.stopPropagation();
      void dolozNadanie(kanal, sekcja);
      return;
    }
    const wiersz = cel.closest<HTMLElement>('[data-dostep-id]');
    const czynnosc = cel.closest<HTMLElement>('[data-dostep-czynnosc]');
    if (wiersz === null || czynnosc === null) return;
    zdarzenie.stopPropagation();
    void wykonaj(kanal, sekcja, czynnosc.dataset.dostepCzynnosc ?? '',
      wiersz.dataset.dostepId ?? '');
  });
}

function grupa(
  dokument: Document,
  kod: string,
  tytul: string,
  opis: string,
  napis: string,
  kolumny: readonly string[],
  zOknem: boolean,
): HTMLElement {
  const pozycja = dokument.createElement('div');
  pozycja.className = 'us-pozycja';

  const glowa = dokument.createElement('div');
  glowa.className = 'us-pozycja-glowa';
  const etykieta = dokument.createElement('span');
  etykieta.className = 'us-etykieta';
  etykieta.textContent = tytul;
  const prawa = dokument.createElement('span');
  prawa.className = 'us-prawa';
  if (zOknem) {
    const okno = dokument.createElement('select');
    okno.className = 'dn-pole-kontrolka';
    okno.id = `${kod}-okno`;
    const pokaz = dokument.createElement('button');
    pokaz.type = 'button';
    pokaz.className = 'dn-btn dn-btn--duch dn-btn--sm';
    pokaz.id = `${kod}-pokaz`;
    pokaz.textContent = 'Pokaż nadania';
    prawa.append(okno, pokaz);
  }
  const przyciskDodania = dokument.createElement('button');
  przyciskDodania.type = 'button';
  przyciskDodania.className = 'dn-btn dn-btn--atrament dn-btn--sm';
  przyciskDodania.id = `${kod}-dodaj`;
  przyciskDodania.textContent = napis;
  prawa.appendChild(przyciskDodania);
  glowa.append(etykieta, prawa);

  const zdanie = dokument.createElement('p');
  zdanie.className = 'us-opis';
  zdanie.textContent = opis;

  const ramka = dokument.createElement('div');
  ramka.className = 'us-tabela-ramka';
  const tabela = dokument.createElement('table');
  tabela.className = 'dn-tabela us-tabela';
  const glowica = dokument.createElement('thead');
  const wierszGlowicy = dokument.createElement('tr');
  for (const nazwa of kolumny) {
    const komorkaGlowicy = dokument.createElement('th');
    komorkaGlowicy.scope = 'col';
    komorkaGlowicy.textContent = nazwa;
    wierszGlowicy.appendChild(komorkaGlowicy);
  }
  glowica.appendChild(wierszGlowicy);
  const cialo = dokument.createElement('tbody');
  cialo.id = `${kod}-cialo`;
  tabela.append(glowica, cialo);
  ramka.appendChild(tabela);

  pozycja.append(glowa, zdanie, ramka);
  return pozycja;
}

async function odswiez(kanal: Kanal, sekcja: Element): Promise<void> {
  await wypelnijPunkty(kanal, sekcja);
  await wypelnijOkna(kanal, sekcja);
  await wypelnijNadania(kanal, sekcja);
}

async function wypelnijPunkty(kanal: Kanal, sekcja: Element): Promise<void> {
  const cialo = sekcja.querySelector('#us-punkty-cialo');
  if (cialo === null) return;
  const wynik = await wywolaj(kanal, Command.AccessPointList, {});
  if (!wynik.udany || wynik.wynik === undefined) {
    cialo.replaceChildren(zdanieWiersza(cialo, 5,
      wynik.blad?.message ?? 'Katalog punktów dostępu nie doszedł.'));
    return;
  }
  const punkty = wynik.wynik.points;
  if (punkty.length === 0) {
    cialo.replaceChildren(zdanieWiersza(cialo, 5,
      'Katalog nie ma jeszcze żadnego punktu dostępu.'));
    return;
  }
  cialo.replaceChildren(...punkty.map((punkt) => wierszPunktu(cialo, punkt)));
}

function wskazaneOkno(sekcja: Element): string {
  return sekcja.querySelector<HTMLSelectElement>('#us-nadania-okno')?.value ?? '';
}

/* Identyfikatora okna komunikacji Operator nigdzie nie widzi, więc wybór idzie
   po nazwie okna, a identyfikator zostaje wartością pozycji. */
async function nazwyOkien(kanal: Kanal): Promise<Array<readonly [string, string]>> {
  const wynik = await wywolaj(kanal, Command.WindowList, {});
  return (wynik.wynik?.windows ?? []).map((okno) =>
    [okno.id, okno.title ?? okno.id] as const);
}

async function nazwyPunktow(kanal: Kanal): Promise<Array<readonly [string, string]>> {
  const wynik = await wywolaj(kanal, Command.AccessPointList, {});
  return (wynik.wynik?.points ?? []).map((punkt) => [punkt.id, punkt.name] as const);
}

async function wypelnijOkna(kanal: Kanal, sekcja: Element): Promise<void> {
  const wybor = sekcja.querySelector<HTMLSelectElement>('#us-nadania-okno');
  if (wybor === null) return;
  const okna = await nazwyOkien(kanal);
  const stojace = wybor.value;
  wybor.replaceChildren();
  for (const [wartosc, nazwa] of [['', 'Wskaż okno'] as const, ...okna]) {
    const pozycja = sekcja.ownerDocument.createElement('option');
    pozycja.value = wartosc;
    pozycja.textContent = nazwa;
    wybor.appendChild(pozycja);
  }
  wybor.value = okna.some((pozycja) => pozycja[0] === stojace) ? stojace : '';
}

/* Rdzeń prowadzi nadania osobno dla każdego okna i bez wskazania okna oddaje
   pusty zbiór, więc wykaz bez wskazania mówi wprost, czego brakuje. */
async function wypelnijNadania(kanal: Kanal, sekcja: Element): Promise<void> {
  const cialo = sekcja.querySelector('#us-nadania-cialo');
  if (cialo === null) return;
  const okno = wskazaneOkno(sekcja);
  if (okno === '') {
    cialo.replaceChildren(zdanieWiersza(cialo, 5,
      'Wskaż okno komunikacji, aby zobaczyć jego nadania.'));
    return;
  }
  const wynik = await wywolaj(kanal, Command.AccessGrantList, { windowId: okno });
  if (!wynik.udany || wynik.wynik === undefined) {
    cialo.replaceChildren(zdanieWiersza(cialo, 5,
      wynik.blad?.message ?? 'Wykaz nadań nie doszedł.'));
    return;
  }
  const nadania = wynik.wynik.grants;
  if (nadania.length === 0) {
    cialo.replaceChildren(zdanieWiersza(cialo, 5,
      'To okno nie ma jeszcze nadanego dostępu.'));
    return;
  }
  cialo.replaceChildren(...nadania.map((nadanie) => wierszNadania(cialo, nadanie)));
}

function zdanieWiersza(cialo: Element, ile: number, zdanie: string): HTMLElement {
  const wiersz = cialo.ownerDocument.createElement('tr');
  const komorkaZdania = cialo.ownerDocument.createElement('td');
  komorkaZdania.colSpan = ile;
  komorkaZdania.className = 'us-opis--drobny';
  komorkaZdania.textContent = zdanie;
  wiersz.appendChild(komorkaZdania);
  return wiersz;
}

function wierszPunktu(cialo: Element, punkt: AccessPoint): HTMLElement {
  const dokument = cialo.ownerDocument;
  const wiersz = dokument.createElement('tr');
  wiersz.dataset.dostepId = punkt.id;
  wiersz.dataset.dostepRodzaj = 'punkt';
  const rodzaj = RODZAJE.find((pozycja) => pozycja[0] === punkt.kind);
  wiersz.append(
    komorka(dokument, punkt.name),
    komorka(dokument, rodzaj === undefined ? punkt.kind : rodzaj[1]),
    komorka(dokument, punkt.roots.join(' · '), 'pt-mono'),
    plakietkaPunktu(dokument, punkt),
    komorkaCzynnosci(dokument, [
      ['popraw-punkt', 'Popraw'],
      ['sprawdz-punkt', 'Sprawdź'],
      ['zdejmij-punkt', 'Zdejmij'],
    ]),
  );
  return wiersz;
}

function wierszNadania(cialo: Element, nadanie: AccessGrant): HTMLElement {
  const dokument = cialo.ownerDocument;
  const wiersz = dokument.createElement('tr');
  wiersz.dataset.dostepId = nadanie.id;
  wiersz.dataset.dostepRodzaj = 'nadanie';
  const tryb = TRYBY.find((pozycja) => pozycja[0] === nadanie.mode);
  wiersz.append(
    komorka(dokument, nadanie.windowId, 'pt-mono'),
    komorka(dokument, nadanie.accessPointId, 'pt-mono'),
    komorka(dokument, tryb === undefined ? nadanie.mode : tryb[1]),
    komorka(dokument, `${String(nadanie.order)}${nadanie.primary ? ' · główne' : ''}`),
    komorkaCzynnosci(dokument, [
      ['popraw-nadanie', 'Popraw'],
      ['zdejmij-nadanie', 'Zdejmij'],
    ]),
  );
  return wiersz;
}

function komorka(dokument: Document, tresc: string, klasa = ''): HTMLElement {
  const wezel = dokument.createElement('td');
  if (klasa !== '') wezel.className = klasa;
  wezel.textContent = tresc;
  return wezel;
}

/* Stan punktu mówi o ostatnim sprawdzeniu, a nie o tym, czy punkt jest czynny —
   punkt wyłączony też może odpowiadać, więc obie rzeczy stoją obok siebie. */
function plakietkaPunktu(dokument: Document, punkt: AccessPoint): HTMLElement {
  const wezel = dokument.createElement('td');
  const znak = dokument.createElement('span');
  znak.className = punkt.status === 'reachable'
    ? 'dn-plakietka dn-plakietka--sukces'
    : 'dn-plakietka';
  znak.textContent = STANY[punkt.status] ?? punkt.status;
  wezel.appendChild(znak);
  if (!punkt.enabled) {
    const wylaczony = dokument.createElement('span');
    wylaczony.className = 'dn-plakietka';
    wylaczony.textContent = 'wyłączony';
    wezel.appendChild(wylaczony);
  }
  return wezel;
}

function komorkaCzynnosci(
  dokument: Document,
  czynnosci: ReadonlyArray<readonly [string, string]>,
): HTMLElement {
  const wezel = dokument.createElement('td');
  const gniazdo = dokument.createElement('span');
  gniazdo.className = 'us-akcje-komorka';
  for (const [kod, napis] of czynnosci) {
    const przyciskCzynnosci = dokument.createElement('button');
    przyciskCzynnosci.type = 'button';
    przyciskCzynnosci.className = 'dn-btn dn-btn--duch dn-btn--sm';
    przyciskCzynnosci.dataset.dostepCzynnosc = kod;
    przyciskCzynnosci.textContent = napis;
    gniazdo.appendChild(przyciskCzynnosci);
  }
  wezel.appendChild(gniazdo);
  return wezel;
}

async function wykonaj(
  kanal: Kanal,
  sekcja: Element,
  czynnosc: string,
  identyfikator: string,
): Promise<void> {
  if (czynnosc === 'popraw-punkt') return poprawPunkt(kanal, sekcja, identyfikator);
  if (czynnosc === 'sprawdz-punkt') return sprawdzPunkt(kanal, sekcja, identyfikator);
  if (czynnosc === 'zdejmij-punkt') return zdejmijPunkt(kanal, sekcja, identyfikator);
  if (czynnosc === 'popraw-nadanie') return poprawNadanie(kanal, sekcja, identyfikator);
  if (czynnosc === 'zdejmij-nadanie') return zdejmijNadanie(kanal, sekcja, identyfikator);
  /* Czynność spoza obsłużonych odmawia zamiast milczeć. */
  oglos(NAGLOWEK, `Czynność „${czynnosc}" nie jest prowadzona przez to okno.`, 'ostrzezenie');
}

async function pobierzPunkt(kanal: Kanal, identyfikator: string): Promise<AccessPoint | undefined> {
  const wykaz = await wywolaj(kanal, Command.AccessPointList, {});
  return wykaz.wynik?.points.find((pozycja) => pozycja.id === identyfikator);
}

async function pobierzNadanie(
  kanal: Kanal,
  sekcja: Element,
  identyfikator: string,
): Promise<AccessGrant | undefined> {
  const wykaz = await wywolaj(kanal, Command.AccessGrantList, { windowId: wskazaneOkno(sekcja) });
  return wykaz.wynik?.grants.find((pozycja) => pozycja.id === identyfikator);
}

function korzenie(wpis: string): string[] {
  return wpis.split('\n').map((wiersz) => wiersz.trim()).filter((wiersz) => wiersz !== '');
}

/* Poświadczenie punktu idzie do sejfu rdzenia i nie wraca do okna. */
async function dolozPunkt(kanal: Kanal, sekcja: Element): Promise<void> {
  const wartosci = await zapytajWOknieModalnym({
    tytul: 'Nowy punkt dostępu',
    opis: 'Każdy korzeń w osobnym wierszu. Poświadczenie idzie do sejfu rdzenia.',
    pola: [
      { klucz: 'nazwa', etykieta: 'Nazwa punktu', wymagane: true },
      { klucz: 'rodzaj', etykieta: 'Rodzaj punktu', wybor: RODZAJE },
      { klucz: 'korzenie', etykieta: 'Korzenie', obszerne: true, wymagane: true },
      { klucz: 'tryb', etykieta: 'Tryb domyślny', wybor: TRYBY },
      { klucz: 'opis', etykieta: 'Opis' },
      { klucz: 'urzadzenie', etykieta: 'Urządzenie' },
      { klucz: 'gospodarz', etykieta: 'Gospodarz' },
      { klucz: 'adres', etykieta: 'Adres mostu' },
      { klucz: 'most', etykieta: 'Nazwa mostu' },
      { klucz: 'poswiadczenie', etykieta: 'Poświadczenie' },
    ],
    wykonanie: 'Dołóż punkt',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.AccessPointAdd, {
    name: wartosci.nazwa ?? '',
    kind: (wartosci.rodzaj ?? AccessPointKind.LocalDirectory) as AccessPointKind,
    roots: korzenie(wartosci.korzenie ?? ''),
    defaultMode: (wartosci.tryb ?? AccessMode.Read) as AccessMode,
    ...(wartosci.opis === '' ? {} : { description: wartosci.opis }),
    ...(wartosci.urzadzenie === '' ? {} : { deviceId: wartosci.urzadzenie }),
    ...(wartosci.gospodarz === '' ? {} : { host: wartosci.gospodarz }),
    ...(wartosci.adres === '' ? {} : { endpoint: wartosci.adres }),
    ...(wartosci.most === '' ? {} : { bridgeName: wartosci.most }),
    ...(wartosci.poswiadczenie === '' ? {} : { credential: wartosci.poswiadczenie }),
    enabled: true,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił dołożenia punktu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Punkt „${wartosci.nazwa ?? ''}" stoi w katalogu.`);
  await odswiez(kanal, sekcja);
}

async function poprawPunkt(kanal: Kanal, sekcja: Element, identyfikator: string): Promise<void> {
  const punkt = await pobierzPunkt(kanal, identyfikator);
  if (punkt === undefined) {
    oglos(NAGLOWEK, 'Katalog nie ma już tego punktu.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWOknieModalnym({
    tytul: `Poprawa punktu „${punkt.name}"`,
    opis: 'Puste pole poświadczenia zostawia poświadczenie stojące nietknięte.',
    pola: [
      { klucz: 'nazwa', etykieta: 'Nazwa punktu', wartosc: punkt.name },
      { klucz: 'korzenie', etykieta: 'Korzenie', obszerne: true, wartosc: punkt.roots.join('\n') },
      { klucz: 'tryb', etykieta: 'Tryb domyślny', wybor: TRYBY, wartosc: punkt.defaultMode },
      { klucz: 'opis', etykieta: 'Opis', wartosc: punkt.description ?? '' },
      { klucz: 'urzadzenie', etykieta: 'Urządzenie', wartosc: punkt.deviceId ?? '' },
      { klucz: 'gospodarz', etykieta: 'Gospodarz', wartosc: punkt.host ?? '' },
      { klucz: 'adres', etykieta: 'Adres mostu', wartosc: punkt.endpoint ?? '' },
      { klucz: 'most', etykieta: 'Nazwa mostu', wartosc: punkt.bridgeName ?? '' },
      { klucz: 'poswiadczenie', etykieta: 'Nowe poświadczenie' },
      {
        klucz: 'czynny',
        etykieta: 'Stan punktu',
        wybor: [['tak', 'Czynny'], ['nie', 'Wyłączony']],
        wartosc: punkt.enabled ? 'tak' : 'nie',
      },
    ],
    wykonanie: 'Zapisz poprawki',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.AccessPointUpdate, {
    accessPointId: identyfikator,
    name: wartosci.nazwa ?? '',
    roots: korzenie(wartosci.korzenie ?? ''),
    defaultMode: (wartosci.tryb ?? punkt.defaultMode) as AccessMode,
    description: wartosci.opis ?? '',
    deviceId: wartosci.urzadzenie ?? '',
    host: wartosci.gospodarz ?? '',
    endpoint: wartosci.adres ?? '',
    bridgeName: wartosci.most ?? '',
    ...(wartosci.poswiadczenie === '' ? {} : { credential: wartosci.poswiadczenie }),
    enabled: wartosci.czynny !== 'nie',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił poprawy punktu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Punkt „${wartosci.nazwa ?? ''}" poprawiony.`);
  await odswiez(kanal, sekcja);
}

async function sprawdzPunkt(kanal: Kanal, sekcja: Element, identyfikator: string): Promise<void> {
  const punkt = await pobierzPunkt(kanal, identyfikator);
  if (punkt === undefined) {
    oglos(NAGLOWEK, 'Katalog nie ma już tego punktu.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.AccessPointCheck, { accessPointId: identyfikator });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił sprawdzenia punktu.', 'ostrzezenie');
    return;
  }
  const stan = STANY[wynik.wynik.status] ?? wynik.wynik.status;
  const potwierdzone = wynik.wynik.roots ?? [];
  const szczegol = wynik.wynik.detail ?? '';
  oglos(NAGLOWEK, `Punkt „${punkt.name}": ${stan}`
    + (potwierdzone.length === 0 ? '' : ` · korzenie ${potwierdzone.join(' · ')}`)
    + (szczegol === '' ? '' : ` · ${szczegol}`));
  await wypelnijPunkty(kanal, sekcja);
}

async function zdejmijPunkt(kanal: Kanal, sekcja: Element, identyfikator: string): Promise<void> {
  const punkt = await pobierzPunkt(kanal, identyfikator);
  if (punkt === undefined) {
    oglos(NAGLOWEK, 'Katalog nie ma już tego punktu.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWOknieModalnym({
    tytul: `Zdjęcie punktu „${punkt.name}"`,
    pola: [],
    wykonanie: 'Zdejmij punkt',
    nieodwracalne: 'Punkt znika z katalogu wraz z nadaniami okien, które po niego sięgały.',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.AccessPointRemove, { accessPointId: identyfikator });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zdjęcia punktu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Punkt „${punkt.name}" zdjęty z katalogu.`);
  await odswiez(kanal, sekcja);
}

/* Korzenie nadania zawężają korzenie punktu; puste pole zostawia zasięg punktu. */
async function dolozNadanie(kanal: Kanal, sekcja: Element): Promise<void> {
  const okna = await nazwyOkien(kanal);
  const punkty = await nazwyPunktow(kanal);
  if (okna.length === 0 || punkty.length === 0) {
    oglos(NAGLOWEK, 'Nadanie wiąże okno z punktem — bez obu stron nie ma czego nadać.',
      'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWOknieModalnym({
    tytul: 'Nowe nadanie dostępu',
    opis: 'Puste korzenie zostawiają nadaniu cały zasięg punktu.',
    pola: [
      { klucz: 'okno', etykieta: 'Okno komunikacji', wybor: okna, wartosc: wskazaneOkno(sekcja) },
      { klucz: 'punkt', etykieta: 'Punkt dostępu', wybor: punkty },
      { klucz: 'tryb', etykieta: 'Tryb', wybor: TRYBY },
      { klucz: 'korzenie', etykieta: 'Korzenie zawężające', obszerne: true },
      { klucz: 'kolejnosc', etykieta: 'Kolejność', wartosc: '0' },
      { klucz: 'glowne', etykieta: 'Nadanie główne', wybor: [['nie', 'Nie'], ['tak', 'Tak']] },
    ],
    wykonanie: 'Nadaj dostęp',
  });
  if (wartosci === null) return;
  const zawezenie = korzenie(wartosci.korzenie ?? '');
  const wynik = await wywolaj(kanal, Command.AccessGrantAdd, {
    windowId: wartosci.okno ?? '',
    accessPointId: wartosci.punkt ?? '',
    mode: (wartosci.tryb ?? AccessMode.Read) as AccessMode,
    ...(zawezenie.length === 0 ? {} : { roots: zawezenie }),
    order: Number(wartosci.kolejnosc ?? '0'),
    primary: wartosci.glowne === 'tak',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił nadania dostępu.', 'ostrzezenie');
    return;
  }
  const wybraneOkno = sekcja.querySelector<HTMLSelectElement>('#us-nadania-okno');
  if (wybraneOkno !== null) wybraneOkno.value = wartosci.okno ?? '';
  oglos(NAGLOWEK, 'Okno sięga teraz po wskazany punkt dostępu.');
  await wypelnijNadania(kanal, sekcja);
}

async function poprawNadanie(kanal: Kanal, sekcja: Element, identyfikator: string): Promise<void> {
  const nadanie = await pobierzNadanie(kanal, sekcja, identyfikator);
  if (nadanie === undefined) {
    oglos(NAGLOWEK, 'Wykaz nie ma już tego nadania.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWOknieModalnym({
    tytul: 'Poprawa nadania dostępu',
    opis: `Okno ${nadanie.windowId} · punkt ${nadanie.accessPointId}`,
    pola: [
      { klucz: 'tryb', etykieta: 'Tryb', wybor: TRYBY, wartosc: nadanie.mode },
      {
        klucz: 'korzenie',
        etykieta: 'Korzenie zawężające',
        obszerne: true,
        wartosc: (nadanie.roots ?? []).join('\n'),
      },
      { klucz: 'kolejnosc', etykieta: 'Kolejność', wartosc: String(nadanie.order) },
      {
        klucz: 'glowne',
        etykieta: 'Nadanie główne',
        wybor: [['nie', 'Nie'], ['tak', 'Tak']],
        wartosc: nadanie.primary ? 'tak' : 'nie',
      },
      {
        klucz: 'czynne',
        etykieta: 'Stan nadania',
        wybor: [['tak', 'Czynne'], ['nie', 'Wyłączone']],
        wartosc: nadanie.enabled ? 'tak' : 'nie',
      },
    ],
    wykonanie: 'Zapisz poprawki',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.AccessGrantUpdate, {
    grantId: identyfikator,
    mode: (wartosci.tryb ?? nadanie.mode) as AccessMode,
    roots: korzenie(wartosci.korzenie ?? ''),
    order: Number(wartosci.kolejnosc ?? '0'),
    primary: wartosci.glowne === 'tak',
    enabled: wartosci.czynne !== 'nie',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił poprawy nadania.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Nadanie poprawione.');
  await wypelnijNadania(kanal, sekcja);
}

async function zdejmijNadanie(kanal: Kanal, sekcja: Element, identyfikator: string): Promise<void> {
  const nadanie = await pobierzNadanie(kanal, sekcja, identyfikator);
  if (nadanie === undefined) {
    oglos(NAGLOWEK, 'Wykaz nie ma już tego nadania.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWOknieModalnym({
    tytul: 'Zdjęcie nadania dostępu',
    opis: `Okno ${nadanie.windowId} · punkt ${nadanie.accessPointId}`,
    pola: [],
    wykonanie: 'Zdejmij nadanie',
    nieodwracalne: 'Okno przestaje sięgać po ten punkt natychmiast po zdjęciu nadania.',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.AccessGrantRemove, { grantId: identyfikator });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zdjęcia nadania.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Nadanie zdjęte.');
  await wypelnijNadania(kanal, sekcja);
}
