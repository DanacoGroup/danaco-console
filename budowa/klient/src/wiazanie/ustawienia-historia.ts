// Sekcja konta w oknie Ustawień: historia rozmów okna, granica przechowywania
// oraz spis nastaw, które platforma w ogóle zna.
import { Command } from '../../../shared/contract.ts';
import type { HistoryEntry, SettingDefinition } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import { grupaKatalogu, komorka, komorkaCzynnosci, wypelnijWybor, zdanieWiersza }
  from './ustawienia-katalog.ts';
import { zapytajWOknieModalnym } from './pytanie-modalne.ts';

const NAGLOWEK = 'Historia i nastawy';

const ZASIEGI: ReadonlyArray<readonly [string, string]> = [
  ['window', 'Okno komunikacji'],
  ['session', 'Sesja'],
  ['application', 'Cała platforma'],
];

/*
zwiazHistorieUstawien dokłada do sekcji konta dwa katalogi: historię wskazanego
okna i spis nastaw platformy. Historia jest własnością okna, nie konta, więc
wykaz zaczyna się od wskazania okna — tak samo jak nadania dostępu.
*/
export function zwiazHistorieUstawien(kanal: Kanal, korzen: Element): void {
  const sekcja = korzen.querySelector('#sekcja-konto');
  if (sekcja === null) return;
  const dokument = sekcja.ownerDocument;
  sekcja.append(
    grupaKatalogu(dokument, {
      kod: 'us-historia',
      tytul: 'Historia rozmów',
      opis: 'Zapis wypowiedzi okna komunikacji. Granica przechowywania mówi, ile '
        + 'rdzeń zostawia; reszta odchodzi sama.',
      napisPrzycisku: 'Ustaw granicę',
      kolumny: ['Kto', 'Zapowiedź', 'Kiedy', 'Akcje'],
      poleWyboru: 'okno',
      napisWyboru: 'Pokaż historię',
    }),
    grupaKatalogu(dokument, {
      kod: 'us-nastawy',
      tytul: 'Spis nastaw platformy',
      opis: 'Wszystko, co platforma umie nastawić, wraz z zasięgami, w których '
        + 'wolno to zmienić.',
      napisPrzycisku: 'Odczytaj spis',
      kolumny: ['Nastawa', 'Rodzaj', 'Zasięgi', 'Dział'],
    }),
  );

  void wypelnijOkna(kanal, sekcja);

  sekcja.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    if (cel.closest('#us-historia-pokaz') !== null) {
      zdarzenie.stopPropagation();
      void wypelnijHistorie(kanal, sekcja);
      return;
    }
    if (cel.closest('#us-historia-dodaj') !== null) {
      zdarzenie.stopPropagation();
      void ustawGranice(kanal, sekcja);
      return;
    }
    if (cel.closest('#us-nastawy-dodaj') !== null) {
      zdarzenie.stopPropagation();
      void wypelnijNastawy(kanal, sekcja);
      return;
    }
    const wiersz = cel.closest<HTMLElement>('[data-wpis-historii]');
    const czynnosc = cel.closest<HTMLElement>('[data-historia-czynnosc]');
    if (wiersz === null || czynnosc === null) return;
    zdarzenie.stopPropagation();
    void zdejmijWpis(kanal, sekcja, wiersz.dataset.wpisHistorii ?? '');
  });
}

function wskazaneOkno(sekcja: Element): string {
  return sekcja.querySelector<HTMLSelectElement>('#us-historia-okno')?.value ?? '';
}

async function wypelnijOkna(kanal: Kanal, sekcja: Element): Promise<void> {
  const wybor = sekcja.querySelector<HTMLSelectElement>('#us-historia-okno');
  if (wybor === null) return;
  const wynik = await wywolaj(kanal, Command.WindowList, {});
  const okna = (wynik.wynik?.windows ?? []).map((okno) =>
    [okno.id, okno.title ?? okno.id] as const);
  wypelnijWybor(wybor, okna, 'Wskaż okno');
}

async function wypelnijHistorie(kanal: Kanal, sekcja: Element): Promise<void> {
  const cialo = sekcja.querySelector('#us-historia-cialo');
  if (cialo === null) return;
  const okno = wskazaneOkno(sekcja);
  if (okno === '') {
    cialo.replaceChildren(zdanieWiersza(cialo, 4,
      'Wskaż okno komunikacji, aby zobaczyć jego historię.'));
    return;
  }
  const wynik = await wywolaj(kanal, Command.HistoryLoad, { windowId: okno, limit: 50 });
  if (!wynik.udany || wynik.wynik === undefined) {
    cialo.replaceChildren(zdanieWiersza(cialo, 4,
      wynik.blad?.message ?? 'Historia nie doszła.'));
    return;
  }
  const wpisy = wynik.wynik.entries;
  if (wpisy.length === 0) {
    cialo.replaceChildren(zdanieWiersza(cialo, 4, 'To okno nie ma jeszcze historii.'));
    return;
  }
  cialo.replaceChildren(...wpisy.map((wpis) => wierszWpisu(cialo, wpis)));
  oglos(NAGLOWEK, `Wpisów w historii: ${String(wynik.wynik.total)}.`);
}

function wierszWpisu(cialo: Element, wpis: HistoryEntry): HTMLElement {
  const dokument = cialo.ownerDocument;
  const wiersz = dokument.createElement('tr');
  wiersz.dataset.wpisHistorii = wpis.id;
  wiersz.append(
    komorka(dokument, wpis.role),
    komorka(dokument, (wpis.preview ?? '').slice(0, 70)),
    komorka(dokument, new Date(wpis.createdAt).toLocaleString('pl')),
    komorkaCzynnosci(dokument, [['zdejmij', 'Zdejmij']], 'data-historia-czynnosc'),
  );
  return wiersz;
}

async function zdejmijWpis(kanal: Kanal, sekcja: Element, idWpisu: string): Promise<void> {
  const okno = wskazaneOkno(sekcja);
  if (okno === '') {
    oglos(NAGLOWEK, 'Bez wskazanego okna nie ma czego zdejmować.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWOknieModalnym({
    tytul: 'Zdjęcie wpisu z historii',
    pola: [],
    wykonanie: 'Zdejmij wpis',
    nieodwracalne: 'Wypowiedź znika z zapisu rozmowy bez możliwości odzyskania.',
  });
  if (wartosci === null) return;
  const wynik = await wywolaj(kanal, Command.HistoryDelete, {
    windowId: okno,
    entryIds: [idWpisu],
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zdjęcia wpisu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Wpis zdjęty z historii.');
  await wypelnijHistorie(kanal, sekcja);
}

/* Granica działa dwiema miarami naraz: dniami i liczbą wpisów. Pusta miara
   znaczy brak granicy z tej strony, nie zero. */
async function ustawGranice(kanal: Kanal, sekcja: Element): Promise<void> {
  const wartosci = await zapytajWOknieModalnym({
    tytul: 'Granica przechowywania',
    opis: 'Pusta miara znaczy brak granicy z tej strony.',
    pola: [
      { klucz: 'zasieg', etykieta: 'Czego dotyczy', wybor: ZASIEGI },
      { klucz: 'dni', etykieta: 'Ile dni trzymać' },
      { klucz: 'wpisy', etykieta: 'Ile wpisów trzymać' },
    ],
    wykonanie: 'Ustaw granicę',
  });
  if (wartosci === null) return;
  const dni = Number(wartosci.dni ?? '');
  const wpisy = Number(wartosci.wpisy ?? '');
  const okno = wskazaneOkno(sekcja);
  const wynik = await wywolaj(kanal, Command.RetentionSet, {
    scope: wartosci.zasieg ?? 'application',
    ...(wartosci.zasieg === 'window' && okno !== '' ? { scopeId: okno } : {}),
    ...(wartosci.dni === '' || !Number.isFinite(dni) ? {} : { keepDays: dni }),
    ...(wartosci.wpisy === '' || !Number.isFinite(wpisy) ? {} : { keepEntries: wpisy }),
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił ustawienia granicy.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Granica przechowywania zapisana.');
}

function opiszZasiegi(nastawa: SettingDefinition): string {
  return nastawa.allowedScopes.join(' · ');
}

/* Spis nastaw idzie w parze z działami: sama nastawa bez działu nie mówi
   Operatorowi, gdzie jej szukać w oknie. */
async function wypelnijNastawy(kanal: Kanal, sekcja: Element): Promise<void> {
  const cialo = sekcja.querySelector('#us-nastawy-cialo');
  if (cialo === null) return;
  const dzialy = await wywolaj(kanal, Command.SettingsCategoryList, {});
  const nazwyDzialow = new Map<string, string>();
  for (const dzial of dzialy.wynik?.categories ?? []) nazwyDzialow.set(dzial.id, dzial.name);
  const wynik = await wywolaj(kanal, Command.SettingsDefinitionList, {});
  if (!wynik.udany || wynik.wynik === undefined) {
    cialo.replaceChildren(zdanieWiersza(cialo, 4,
      wynik.blad?.message ?? 'Spis nastaw nie doszedł.'));
    return;
  }
  const nastawy = wynik.wynik.definitions;
  if (nastawy.length === 0) {
    cialo.replaceChildren(zdanieWiersza(cialo, 4, 'Rdzeń nie podaje żadnej nastawy.'));
    return;
  }
  const dokument = cialo.ownerDocument;
  cialo.replaceChildren(...nastawy.slice(0, 60).map((nastawa) => {
    const wiersz = dokument.createElement('tr');
    wiersz.append(
      komorka(dokument, nastawa.name),
      komorka(dokument, nastawa.valueType),
      komorka(dokument, opiszZasiegi(nastawa), 'pt-mono'),
      komorka(dokument, nazwyDzialow.get(nastawa.categoryId) ?? nastawa.categoryId),
    );
    return wiersz;
  }));
  oglos(NAGLOWEK, `Nastaw w spisie: ${String(wynik.wynik.total)}`
    + ` · działów: ${String(nazwyDzialow.size)}.`);
}
