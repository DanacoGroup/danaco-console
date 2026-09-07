// Czynności pracy na ustaleniach: scalenie, wątki tematyczne, kodowanie,
// macierz kod na źródło, ślad pochodzenia, rozstrzygnięcie sprzeczności
// oraz książka kodów badania.
import { Command, type ResearchCode, type ResearchFinding } from '../../../shared/contract.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import { niegotowyPanel, wykazPanelu } from './okno-modulu.ts';
import {
  brakOkna,
  nasluchCzynnosci,
  nieznanaCzynnosc,
  odmowa,
  potwierdzone,
  powiedz,
  wybierz,
  zapytaj,
  type KontekstBadania,
} from './research-czynnosci.ts';
import { wywolaj } from '../protokol/wywolanie.ts';

const ATRYBUT_CZYNNOSCI = 'data-czynnosc-syntezy';
const GNIAZDO_WYNIKU = '[data-synteza-wynik]';

export function zwiazCzynnosciSyntezy(kontekst: KontekstBadania): Odsubskrybuj {
  return nasluchCzynnosci(kontekst.korzen, ATRYBUT_CZYNNOSCI, async (kod) => {
    if (brakOkna(kontekst.idOkna())) return;
    await wykonaj(kontekst, kod);
  });
}

async function wykonaj(kontekst: KontekstBadania, kod: string): Promise<void> {
  if (kod === 'scal') return scal(kontekst);
  if (kod === 'watki') return watki(kontekst);
  if (kod === 'koduj') return koduj(kontekst);
  if (kod === 'macierz') return macierz(kontekst);
  if (kod === 'slad') return slad(kontekst);
  if (kod === 'rozstrzygnij') return rozstrzygnij(kontekst);
  if (kod === 'ksiazka') return ksiazka(kontekst);
  if (kod === 'dopisz-kod') return dopiszKod(kontekst);
  nieznanaCzynnosc(kod);
}

function gniazdo(kontekst: KontekstBadania): Element | null {
  return kontekst.korzen.querySelector(GNIAZDO_WYNIKU);
}

async function ustalenia(kontekst: KontekstBadania): Promise<ResearchFinding[]> {
  const wykaz = await wywolaj(kontekst.kanal, Command.ResearchFindingList, {
    windowId: kontekst.idOkna(),
  });
  if (!wykaz.udany || wykaz.wynik === undefined) {
    odmowa(wykaz.blad, 'Wykaz ustaleń nie doszedł, więc nie ma na czym pracować.');
    return [];
  }
  if (wykaz.wynik.findings.length === 0) {
    odmowa(undefined, 'Okno nie ma ustaleń, więc czynność nie ma przedmiotu.');
  }
  return wykaz.wynik.findings;
}

function wskaz(pytanie: string, wykaz: readonly ResearchFinding[]): ResearchFinding | undefined {
  if (wykaz.length === 0) return undefined;
  const numer = wybierz(pytanie, wykaz.map((jedno) => jedno.content));
  const wybrane = numer < 0 ? undefined : wykaz[numer];
  if (wybrane === undefined) {
    odmowa(undefined, 'Żadne ustalenie nie zostało wskazane, więc czynność nie poszła do rdzenia.');
  }
  return wybrane;
}

async function scal(kontekst: KontekstBadania): Promise<void> {
  const wykaz = await ustalenia(kontekst);
  if (wykaz.length < 2) {
    if (wykaz.length === 1) odmowa(undefined, 'Scalenie wymaga dwóch ustaleń, a okno ma jedno.');
    return;
  }
  const zachowane = wskaz('Ustalenie zachowywane', wykaz);
  if (zachowane === undefined) return;
  const pozostale = wykaz.filter((jedno) => jedno.id !== zachowane.id);
  const scalane = wskaz('Ustalenie scalane i usuwane', pozostale);
  if (scalane === undefined) return;
  if (!potwierdzone(`scalenie-ustalen:${zachowane.id}`,
    'Scalenie usuwa ustalenie scalane. Naciśnij drugi raz, żeby je wykonać.')) return;
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchFindingMerge, {
    targetFindingId: zachowane.id,
    mergedFindingIds: [scalane.id],
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    odmowa(wynik.blad, 'Rdzeń odmówił scalenia ustaleń.');
    return;
  }
  powiedz(`Scalono; przeniesionych powiązań ze źródłami: ${String(wynik.wynik.movedSources)}.`);
  await kontekst.odswiez();
}

/* Ustalenia poza wątkami są nazwane osobno: milczenie o nich wyglądałoby tak,
   jakby klastrowanie objęło całość materiału. */
async function watki(kontekst: KontekstBadania): Promise<void> {
  const cialo = gniazdo(kontekst);
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchFindingCluster, {
    windowId: kontekst.idOkna(),
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    niegotowyPanel(cialo, wynik.blad?.message ?? 'Wątki tematyczne nie doszły.');
    odmowa(wynik.blad, 'Rdzeń odmówił pogrupowania ustaleń w wątki.');
    return;
  }
  wykazPanelu(cialo, true, undefined, wynik.wynik.threads,
    'Rdzeń nie wydzielił żadnego wątku tematycznego.',
    (watek) => [watek.name, `${String(watek.findingIds.length)} ustaleń`] as const);
  const poza = wynik.wynik.ungroupedFindingIds ?? [];
  if (poza.length > 0) {
    odmowa(undefined, `Ustaleń poza wątkami: ${String(poza.length)}.`);
  }
}

async function kody(kontekst: KontekstBadania): Promise<ResearchCode[] | undefined> {
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchCodebookGet, {
    windowId: kontekst.idOkna(),
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    odmowa(wynik.blad, 'Rdzeń nie oddał książki kodów badania.');
    return undefined;
  }
  return wynik.wynik.codes;
}

/* Nazwy wpisane po przecinku rozchodzą się na kody już istniejące i kody
   zakładane przy okazji — inaczej każde kodowanie mnożyłoby kody o tej samej
   nazwie, a macierz przestałaby cokolwiek zliczać. */
async function koduj(kontekst: KontekstBadania): Promise<void> {
  const wybrane = wskaz('Kodowane ustalenie', await ustalenia(kontekst));
  if (wybrane === undefined) return;
  const ksiazkaKodow = await kody(kontekst);
  if (ksiazkaKodow === undefined) return;
  const wpis = zapytaj('Nazwy kodów tematycznych, po przecinku');
  if (wpis === '') return;
  const nazwy = wpis.split(',').map((czesc) => czesc.trim()).filter((czesc) => czesc !== '');
  if (nazwy.length === 0) {
    odmowa(undefined, 'Żaden kod nie został podany, więc kodowanie nie poszło do rdzenia.');
    return;
  }
  const znane: string[] = [];
  const nowe: string[] = [];
  for (const nazwa of nazwy) {
    const kod = ksiazkaKodow.find((jeden) => jeden.name === nazwa);
    if (kod === undefined) nowe.push(nazwa);
    else znane.push(kod.id);
  }
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchFindingCode, {
    findingId: wybrane.id,
    codeIds: znane.length === 0 ? undefined : znane,
    newCodeNames: nowe.length === 0 ? undefined : nowe,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    odmowa(wynik.blad, 'Rdzeń odmówił zakodowania ustalenia.');
    return;
  }
  powiedz(`Ustalenie nosi ${String(wynik.wynik.codes.length)} kodów.`);
  await kontekst.odswiez();
}

async function macierz(kontekst: KontekstBadania): Promise<void> {
  const cialo = gniazdo(kontekst);
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchFindingMatrix, {
    windowId: kontekst.idOkna(),
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    niegotowyPanel(cialo, wynik.blad?.message ?? 'Macierz kodowania nie doszła.');
    odmowa(wynik.blad, 'Rdzeń odmówił złożenia macierzy kodowania.');
    return;
  }
  const nazwy = new Map(wynik.wynik.codes.map((kod) => [kod.id, kod.name]));
  wykazPanelu(cialo, true, undefined, wynik.wynik.cells,
    'Żadne źródło nie zostało jeszcze zakodowane.',
    (komorka) => [`${nazwy.get(komorka.codeId) ?? komorka.codeId} — ${komorka.sourceId}`,
      `${String(komorka.count)} wystąpień`] as const);
}

async function slad(kontekst: KontekstBadania): Promise<void> {
  const wybrane = wskaz('Ustalenie, którego ślad ma stanąć', await ustalenia(kontekst));
  if (wybrane === undefined) return;
  const cialo = gniazdo(kontekst);
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchFindingProvenance, {
    findingId: wybrane.id,
  });
  wykazPanelu(cialo, wynik.udany, wynik.blad?.message, wynik.wynik?.entries,
    'Rdzeń nie prowadzi śladu tego ustalenia.',
    (wpis) => [wpis.action, `${String(wpis.actor)} · ${wpis.sourceId ?? 'bez źródła'}`] as const);
  if (!wynik.udany) odmowa(wynik.blad, 'Rdzeń odmówił wydania śladu pochodzenia.');
}

async function rozstrzygnij(kontekst: KontekstBadania): Promise<void> {
  const wykryte = await wywolaj(kontekst.kanal, Command.ResearchFindingContradictions, {
    windowId: kontekst.idOkna(),
  });
  if (!wykryte.udany || wykryte.wynik === undefined) {
    odmowa(wykryte.blad, 'Rdzeń odmówił wykrycia sprzeczności.');
    return;
  }
  const otwarte = wykryte.wynik.contradictions.filter((jedna) => !jedna.resolved);
  if (otwarte.length === 0) {
    powiedz('Żadna sprzeczność tego okna nie czeka na rozstrzygnięcie.');
    return;
  }
  const numer = wybierz('Sprzeczność do rozstrzygnięcia', otwarte.map((jedna) => jedna.summary));
  const sprzecznosc = numer < 0 ? undefined : otwarte[numer];
  if (sprzecznosc === undefined) {
    odmowa(undefined, 'Żadna sprzeczność nie została wskazana, więc nic nie poszło do rdzenia.');
    return;
  }
  const uzasadnienie = zapytaj('Uzasadnienie rozstrzygnięcia');
  if (uzasadnienie === '') {
    odmowa(undefined, 'Rozstrzygnięcie bez uzasadnienia nic nie tłumaczy, więc nie poszło do rdzenia.');
    return;
  }
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchContradictionResolve, {
    contradictionId: sprzecznosc.id,
    rationale: uzasadnienie,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    odmowa(wynik.blad, 'Rdzeń odmówił rozstrzygnięcia sprzeczności.');
    return;
  }
  if (!wynik.wynik.contradiction.resolved) {
    odmowa(undefined, 'Rdzeń przyjął uzasadnienie, ale sprzeczność nadal stoi jako otwarta.');
    return;
  }
  powiedz('Sprzeczność rozstrzygnięta.');
  await kontekst.odswiez();
}

async function ksiazka(kontekst: KontekstBadania): Promise<void> {
  const cialo = gniazdo(kontekst);
  const ksiazkaKodow = await kody(kontekst);
  if (ksiazkaKodow === undefined) {
    niegotowyPanel(cialo, 'Książka kodów nie doszła.');
    return;
  }
  wykazPanelu(cialo, true, undefined, ksiazkaKodow,
    'Badanie nie prowadzi jeszcze żadnego kodu tematycznego.',
    (kod) => [kod.name, `${String(kod.occurrences ?? 0)} wystąpień`] as const);
}

/* Zapis książki kodów idzie całym zestawem, więc nowy kod dopisuje się do tego,
   co rdzeń już prowadzi — inaczej zapis skasowałby kody wcześniejsze. */
async function dopiszKod(kontekst: KontekstBadania): Promise<void> {
  const ksiazkaKodow = await kody(kontekst);
  if (ksiazkaKodow === undefined) return;
  const nazwa = zapytaj('Nazwa kodu tematycznego');
  if (nazwa === '') return;
  if (ksiazkaKodow.some((kod) => kod.name === nazwa)) {
    odmowa(undefined, `Kod „${nazwa}" już stoi w książce kodów, więc nie został dopisany.`);
    return;
  }
  const definicja = zapytaj('Definicja kodu w książce kodów');
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchCodebookSet, {
    windowId: kontekst.idOkna(),
    codes: [...ksiazkaKodow, {
      id: '',
      name: nazwa,
      description: definicja === '' ? undefined : definicja,
    }],
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    odmowa(wynik.blad, 'Rdzeń odmówił zapisania książki kodów.');
    return;
  }
  powiedz(`Książka kodów prowadzi ${String(wynik.wynik.codes.length)} kodów.`);
  await ksiazka(kontekst);
}
