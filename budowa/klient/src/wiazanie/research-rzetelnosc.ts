// Czynności rzetelności badania: kompletność cytowań, skład cytatu, sygnały
// wycofania prac, liczniki przesiewu PRISMA, graf dowodów i luki badawcze.
import { Command, ResearchCitationMode } from '../../../shared/contract.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import { niegotowyPanel, wykazPanelu } from './okno-modulu.ts';
import {
  brakOkna,
  doSchowka,
  nasluchCzynnosci,
  nieznanaCzynnosc,
  odmowa,
  powiedz,
  wybierz,
  type KontekstBadania,
} from './research-czynnosci.ts';
import { wywolaj } from '../protokol/wywolanie.ts';

const ATRYBUT_CZYNNOSCI = 'data-czynnosc-rzetelnosci';
const GNIAZDO_WYNIKU = '[data-rzetelnosc-wynik]';

export function zwiazCzynnosciRzetelnosci(kontekst: KontekstBadania): Odsubskrybuj {
  return nasluchCzynnosci(kontekst.korzen, ATRYBUT_CZYNNOSCI, async (kod) => {
    if (brakOkna(kontekst.idOkna())) return;
    await wykonaj(kontekst, kod);
  });
}

async function wykonaj(kontekst: KontekstBadania, kod: string): Promise<void> {
  if (kod === 'cytowania') return cytowania(kontekst);
  if (kod === 'zloz-cytat') return zlozCytat(kontekst);
  if (kod === 'wycofania') return wycofania(kontekst);
  if (kod === 'przesiew') return przesiew(kontekst);
  if (kod === 'graf') return graf(kontekst);
  if (kod === 'luki') return luki(kontekst);
  nieznanaCzynnosc(kod);
}

function gniazdo(kontekst: KontekstBadania): Element | null {
  return kontekst.korzen.querySelector(GNIAZDO_WYNIKU);
}

async function identyfikatoryZrodel(kontekst: KontekstBadania): Promise<string[]> {
  const wykaz = await wywolaj(kontekst.kanal, Command.ResearchSourceList, {
    windowId: kontekst.idOkna(),
  });
  if (!wykaz.udany || wykaz.wynik === undefined) {
    odmowa(wykaz.blad, 'Wykaz źródeł nie doszedł, więc nie ma czego sprawdzić.');
    return [];
  }
  const zrodla = wykaz.wynik.sources;
  if (zrodla.length === 0) odmowa(undefined, 'Okno nie ma źródeł, więc czynność nie ma przedmiotu.');
  return zrodla.map((zrodlo) => zrodlo.id);
}

async function cytowania(kontekst: KontekstBadania): Promise<void> {
  const cialo = gniazdo(kontekst);
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchCitationCheck, {
    windowId: kontekst.idOkna(),
  });
  wykazPanelu(cialo, wynik.udany, wynik.blad?.message, wynik.wynik?.issues,
    'Rdzeń nie wykrył braków w metadanych ani w zgodności cytowań.',
    (brak) => [`${brak.kind}${brak.field === undefined ? '' : ` — ${brak.field}`}`,
      brak.fixableFrom ?? brak.sourceId ?? ''] as const);
  if (!wynik.udany) odmowa(wynik.blad, 'Rdzeń odmówił sprawdzenia cytowań.');
}

/* Źródła o metadanych niepełnych są nazwane osobno: skład, który wyszedł
   z brakami, wygląda w schowku tak samo jak skład kompletny. */
async function zlozCytat(kontekst: KontekstBadania): Promise<void> {
  const zrodla = await identyfikatoryZrodel(kontekst);
  if (zrodla.length === 0) return;
  const style = await wywolaj(kontekst.kanal, Command.ResearchCitationStyles, {});
  const wykaz = style.wynik?.styles ?? [];
  if (!style.udany || wykaz.length === 0) {
    odmowa(style.blad, 'Rdzeń nie oddał żadnego stylu cytowania, więc skład nie ma reguły.');
    return;
  }
  const numer = wybierz('Styl cytowania', wykaz.map((styl) => styl.name));
  const styl = numer < 0 ? undefined : wykaz[numer];
  if (styl === undefined) {
    odmowa(undefined, 'Styl cytowania nie został wskazany, więc skład nie poszedł do rdzenia.');
    return;
  }
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchCitationRender, {
    sourceIds: zrodla,
    styleId: styl.id,
    mode: ResearchCitationMode.Both,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    odmowa(wynik.blad, 'Rdzeń odmówił złożenia cytatów.');
    return;
  }
  const niepelne = wynik.wynik.incompleteSourceIds ?? [];
  if (niepelne.length > 0) {
    odmowa(undefined, `Źródeł o metadanych niepełnych: ${String(niepelne.length)}.`);
  }
  const tresc = wynik.wynik.citations
    .map((cytat) => cytat.bibliography ?? cytat.inText ?? '')
    .filter((wiersz) => wiersz !== '')
    .join('\n');
  if (tresc === '') {
    odmowa(undefined, 'Rdzeń przyjął skład, ale nie oddał ani odnośnika, ani pozycji bibliograficznej.');
    return;
  }
  await doSchowka(tresc, `cytaty-${styl.id}`);
}

async function wycofania(kontekst: KontekstBadania): Promise<void> {
  const zrodla = await identyfikatoryZrodel(kontekst);
  if (zrodla.length === 0) return;
  const cialo = gniazdo(kontekst);
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchRetractionCheck, {
    sourceIds: zrodla,
  });
  wykazPanelu(cialo, wynik.udany, wynik.blad?.message, wynik.wynik?.flags,
    'Rdzeń nie oddał sygnału wycofania dla żadnego źródła.',
    (sygnal) => [sygnal.sourceId, sygnal.status] as const);
  if (!wynik.udany) odmowa(wynik.blad, 'Rdzeń odmówił sprawdzenia wycofań.');
}

async function przesiew(kontekst: KontekstBadania): Promise<void> {
  const cialo = gniazdo(kontekst);
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchPrismaGet, {
    windowId: kontekst.idOkna(),
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    niegotowyPanel(cialo, wynik.blad?.message ?? 'Liczniki przesiewu nie doszły.');
    odmowa(wynik.blad, 'Rdzeń odmówił wydania liczników przesiewu.');
    return;
  }
  const liczniki = wynik.wynik.counts;
  const etapy: readonly (readonly [string, number])[] = [
    ['Zidentyfikowane', liczniki.identified],
    ['Duplikaty usunięte', liczniki.duplicatesRemoved],
    ['Przesiane', liczniki.screened],
    ['Odrzucone', liczniki.excluded],
    ['Włączone', liczniki.included],
  ];
  wykazPanelu(cialo, true, undefined, [...etapy], 'Przesiew nie ma jeszcze żadnej pozycji.',
    (etap) => [etap[0], String(etap[1])] as const);
  const powody = liczniki.exclusionReasons ?? [];
  if (powody.length > 0) powiedz(`Uzasadnienia odrzuceń: ${powody.join('; ')}.`);
}

async function graf(kontekst: KontekstBadania): Promise<void> {
  const cialo = gniazdo(kontekst);
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchEvidenceGraph, {
    windowId: kontekst.idOkna(),
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    niegotowyPanel(cialo, wynik.blad?.message ?? 'Graf dowodów nie doszedł.');
    odmowa(wynik.blad, 'Rdzeń odmówił złożenia grafu dowodów.');
    return;
  }
  const nazwy = new Map(wynik.wynik.nodes.map((wezel) => [wezel.id, wezel.label]));
  wykazPanelu(cialo, true, undefined, wynik.wynik.edges,
    'Graf dowodów nie ma jeszcze żadnej relacji.',
    (krawedz) => [`${nazwy.get(krawedz.fromId) ?? krawedz.fromId} → `
      + `${nazwy.get(krawedz.toId) ?? krawedz.toId}`, krawedz.relation] as const);
  powiedz(`Węzłów grafu: ${String(wynik.wynik.nodes.length)}.`);
}

async function luki(kontekst: KontekstBadania): Promise<void> {
  const cialo = gniazdo(kontekst);
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchGapFind, {
    windowId: kontekst.idOkna(),
  });
  wykazPanelu(cialo, wynik.udany, wynik.blad?.message, wynik.wynik?.gaps,
    'Rdzeń nie wskazał luki badawczej.',
    (luka) => [luka.summary, luka.suggestedQuery ?? ''] as const);
  if (!wynik.udany) odmowa(wynik.blad, 'Rdzeń odmówił wskazania luk badawczych.');
}
