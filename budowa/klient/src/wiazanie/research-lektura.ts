// Czynności lektury i odkryć: układanie zapytań, graf cytowań, wczytanie treści
// źródła, adnotacje, pytanie do korpusu i zakres przestrzeni badania.
import {
  Command,
  ResearchAnchorKind,
  ResearchAnnotationKind,
  ResearchDiscoveryMode,
  ResearchSnowballDirection,
  type ResearchSource,
} from '../../../shared/contract.ts';
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

const ATRYBUT_CZYNNOSCI = 'data-czynnosc-lektury';
const GNIAZDO_WYNIKU = '[data-lektura-wynik]';

const TRYBY_SZUKANIA: readonly (readonly [string, ResearchDiscoveryMode])[] = [
  ['Sieć', ResearchDiscoveryMode.Web],
  ['Bazy publikacji naukowych', ResearchDiscoveryMode.Scholarly],
];

const KIERUNKI_GRAFU: readonly (readonly [string, ResearchSnowballDirection])[] = [
  ['Prace cytowane', ResearchSnowballDirection.Backward],
  ['Prace cytujące', ResearchSnowballDirection.Forward],
  ['Oba kierunki', ResearchSnowballDirection.Both],
];

const RODZAJE_ADNOTACJI: readonly (readonly [string, ResearchAnnotationKind])[] = [
  ['Podświetlenie', ResearchAnnotationKind.Highlight],
  ['Notatka', ResearchAnnotationKind.Note],
  ['Zakładka', ResearchAnnotationKind.Bookmark],
];

export function zwiazCzynnosciLektury(kontekst: KontekstBadania): Odsubskrybuj {
  return nasluchCzynnosci(kontekst.korzen, ATRYBUT_CZYNNOSCI, async (kod) => {
    if (brakOkna(kontekst.idOkna())) return;
    await wykonaj(kontekst, kod);
  });
}

async function wykonaj(kontekst: KontekstBadania, kod: string): Promise<void> {
  if (kod === 'zapytania') return zapytania(kontekst);
  if (kod === 'cytowania') return grafCytowan(kontekst);
  if (kod === 'czytaj') return czytaj(kontekst);
  if (kod === 'adnotuj') return adnotuj(kontekst);
  if (kod === 'adnotacje') return adnotacje(kontekst);
  if (kod === 'usun-adnotacje') return usunAdnotacje(kontekst);
  if (kod === 'pytanie') return pytanieDoKorpusu(kontekst);
  if (kod === 'zakres') return zakres(kontekst);
  nieznanaCzynnosc(kod);
}

function gniazdo(kontekst: KontekstBadania): Element | null {
  return kontekst.korzen.querySelector(GNIAZDO_WYNIKU);
}

async function wskazZrodlo(kontekst: KontekstBadania): Promise<ResearchSource | undefined> {
  const wykaz = await wywolaj(kontekst.kanal, Command.ResearchSourceList, {
    windowId: kontekst.idOkna(),
  });
  if (!wykaz.udany || wykaz.wynik === undefined) {
    odmowa(wykaz.blad, 'Wykaz źródeł nie doszedł, więc nie ma czego czytać.');
    return undefined;
  }
  const zrodla = wykaz.wynik.sources;
  if (zrodla.length === 0) {
    odmowa(undefined, 'Okno nie ma źródeł, więc czynność nie ma przedmiotu.');
    return undefined;
  }
  const numer = wybierz('Źródło', zrodla.map((zrodlo) => zrodlo.title));
  const wybrane = numer < 0 ? undefined : zrodla[numer];
  if (wybrane === undefined) {
    odmowa(undefined, 'Żadne źródło nie zostało wskazane, więc czynność nie poszła do rdzenia.');
  }
  return wybrane;
}

async function zapytania(kontekst: KontekstBadania): Promise<void> {
  const pytanie = zapytaj('Pytanie badawcze do przełożenia na zapytania');
  if (pytanie === '') return;
  const numer = wybierz('Tryb wyszukiwania', TRYBY_SZUKANIA.map((pozycja) => pozycja[0]));
  const tryb = numer < 0 ? undefined : TRYBY_SZUKANIA[numer];
  if (tryb === undefined) {
    odmowa(undefined, 'Tryb wyszukiwania nie został wskazany, więc pytanie nie poszło do rdzenia.');
    return;
  }
  const cialo = gniazdo(kontekst);
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchDiscoveryAssist, {
    question: pytanie,
    mode: tryb[1],
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    niegotowyPanel(cialo, wynik.blad?.message ?? 'Propozycje zapytań nie doszły.');
    odmowa(wynik.blad, 'Rdzeń odmówił ułożenia zapytań wyszukiwawczych.');
    return;
  }
  wykazPanelu(cialo, true, undefined, wynik.wynik.queries,
    'Rdzeń nie zaproponował żadnego zapytania.', (zapytanie) => [zapytanie, ''] as const);
  const uzasadnienie = wynik.wynik.rationale;
  if (uzasadnienie !== undefined) powiedz(uzasadnienie);
}

async function grafCytowan(kontekst: KontekstBadania): Promise<void> {
  const zrodlo = await wskazZrodlo(kontekst);
  if (zrodlo === undefined) return;
  const numer = wybierz('Kierunek rozwijania', KIERUNKI_GRAFU.map((pozycja) => pozycja[0]));
  const kierunek = numer < 0 ? undefined : KIERUNKI_GRAFU[numer];
  if (kierunek === undefined) {
    odmowa(undefined, 'Kierunek nie został wskazany, więc graf cytowań nie poszedł do rdzenia.');
    return;
  }
  const cialo = gniazdo(kontekst);
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchDiscoverySnowball, {
    sourceId: zrodlo.id,
    direction: kierunek[1],
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    niegotowyPanel(cialo, wynik.blad?.message ?? 'Graf cytowań nie doszedł.');
    odmowa(wynik.blad, 'Rdzeń odmówił rozwinięcia grafu cytowań.');
    return;
  }
  const pozycje = [...(wynik.wynik.referenced ?? []), ...(wynik.wynik.citing ?? [])];
  wykazPanelu(cialo, true, undefined, pozycje,
    'Rdzeń nie oddał żadnej pracy powiązanej cytowaniem.',
    (pozycja) => [pozycja.title, pozycja.provider] as const);
}

/* Skan bez warstwy tekstowej i treść skrócona są nazwane, bo w obu wypadkach
   w panelu stoi mniej, niż źródło niesie — a wygląda to tak samo. */
async function czytaj(kontekst: KontekstBadania): Promise<void> {
  const zrodlo = await wskazZrodlo(kontekst);
  if (zrodlo === undefined) return;
  const wpis = zapytaj('Numer strony; puste znaczy stronę pierwszą');
  const strona = wpis === '' ? undefined : Number.parseInt(wpis, 10);
  if (strona !== undefined && (Number.isNaN(strona) || strona < 1)) {
    odmowa(undefined, 'Numer strony musi być liczbą od jednego, więc lektura nie poszła do rdzenia.');
    return;
  }
  const cialo = gniazdo(kontekst);
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchReadingOpen, {
    sourceId: zrodlo.id,
    page: strona,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    niegotowyPanel(cialo, wynik.blad?.message ?? 'Treść źródła nie doszła.');
    odmowa(wynik.blad, 'Rdzeń odmówił wczytania treści źródła.');
    return;
  }
  const tresc = wynik.wynik.content;
  niegotowyPanel(cialo, tresc.text ?? 'Rdzeń oddał stronę bez treści tekstowej.');
  powiedz(`Strona ${String(tresc.page ?? 1)} z ${String(tresc.pageCount ?? 1)}.`);
  if (tresc.hasTextLayer === false) {
    odmowa(undefined, 'Źródło nie ma warstwy tekstowej — najpierw rozpoznaj w nim pismo.');
  }
  if (tresc.truncated === true) {
    odmowa(undefined, 'Rdzeń oddał treść skróconą, więc w panelu stoi niepełna strona.');
  }
}

async function adnotuj(kontekst: KontekstBadania): Promise<void> {
  const zrodlo = await wskazZrodlo(kontekst);
  if (zrodlo === undefined) return;
  const numer = wybierz('Rodzaj adnotacji', RODZAJE_ADNOTACJI.map((pozycja) => pozycja[0]));
  const rodzaj = numer < 0 ? undefined : RODZAJE_ADNOTACJI[numer];
  if (rodzaj === undefined) {
    odmowa(undefined, 'Rodzaj adnotacji nie został wskazany, więc zapis nie poszedł do rdzenia.');
    return;
  }
  const wpis = zapytaj('Numer strony, na której stoi fragment', '1');
  const strona = Number.parseInt(wpis, 10);
  if (Number.isNaN(strona) || strona < 1) {
    odmowa(undefined, 'Adnotacja bez strony nie da się zakotwiczyć, więc nie poszła do rdzenia.');
    return;
  }
  const cytat = zapytaj('Cytowany fragment');
  const komentarz = zapytaj('Komentarz roboczy');
  if (cytat === '' && komentarz === '') {
    odmowa(undefined, 'Adnotacja bez cytatu i bez komentarza nic nie odnotowuje.');
    return;
  }
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchAnnotationAdd, {
    sourceId: zrodlo.id,
    kind: rodzaj[1],
    anchor: { kind: ResearchAnchorKind.Page, page: strona },
    quote: cytat === '' ? undefined : cytat,
    comment: komentarz === '' ? undefined : komentarz,
  });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił zapisania adnotacji.');
    return;
  }
  powiedz(`Adnotacja zapisana przy źródle „${zrodlo.title}".`);
  await kontekst.odswiez();
}

async function adnotacje(kontekst: KontekstBadania): Promise<void> {
  const cialo = gniazdo(kontekst);
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchAnnotationList, {
    windowId: kontekst.idOkna(),
  });
  wykazPanelu(cialo, wynik.udany, wynik.blad?.message, wynik.wynik?.annotations,
    'Żadne źródło tego okna nie ma adnotacji.',
    (adnotacja) => [adnotacja.quote ?? adnotacja.comment ?? adnotacja.id,
      String(adnotacja.kind)] as const);
  if (!wynik.udany) odmowa(wynik.blad, 'Rdzeń odmówił wydania adnotacji.');
}

async function usunAdnotacje(kontekst: KontekstBadania): Promise<void> {
  const wykaz = await wywolaj(kontekst.kanal, Command.ResearchAnnotationList, {
    windowId: kontekst.idOkna(),
  });
  if (!wykaz.udany || wykaz.wynik === undefined) {
    odmowa(wykaz.blad, 'Wykaz adnotacji nie doszedł, więc nie ma czego usunąć.');
    return;
  }
  const lista = wykaz.wynik.annotations;
  if (lista.length === 0) {
    odmowa(undefined, 'Żadne źródło tego okna nie ma adnotacji, więc nic nie zostało usunięte.');
    return;
  }
  const numer = wybierz('Adnotacja do usunięcia',
    lista.map((jedna) => jedna.quote ?? jedna.comment ?? jedna.id));
  const adnotacja = numer < 0 ? undefined : lista[numer];
  if (adnotacja === undefined) {
    odmowa(undefined, 'Żadna adnotacja nie została wskazana, więc nic nie zostało usunięte.');
    return;
  }
  if (!potwierdzone(`adnotacja:${adnotacja.id}`,
    'Usunięcie adnotacji jest nieodwracalne. Naciśnij drugi raz, żeby je wykonać.')) return;
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchAnnotationRemove, {
    annotationId: adnotacja.id,
  });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił usunięcia adnotacji.');
    return;
  }
  powiedz('Adnotacja usunięta.');
  await adnotacje(kontekst);
}

/* Zakotwiczenie w cytatach jest wymagane, a odpowiedź bez pokrycia nazwana:
   inaczej wywód modelu stanąłby w panelu jako ustalenie ze źródeł. */
async function pytanieDoKorpusu(kontekst: KontekstBadania): Promise<void> {
  const pytanie = zapytaj('Pytanie do korpusu źródeł');
  if (pytanie === '') return;
  const cialo = gniazdo(kontekst);
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchCorpusAsk, {
    windowId: kontekst.idOkna(),
    question: pytanie,
    requireGrounding: true,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    niegotowyPanel(cialo, wynik.blad?.message ?? 'Odpowiedź korpusu nie doszła.');
    odmowa(wynik.blad, 'Rdzeń odmówił odpowiedzi z korpusu źródeł.');
    return;
  }
  const odpowiedz = wynik.wynik.answer;
  niegotowyPanel(cialo, odpowiedz.text);
  powiedz(`Cytatów, w których odpowiedź stoi: ${String(odpowiedz.citations.length)}.`);
  if (!odpowiedz.grounded) {
    odmowa(undefined, 'Nie każde zdanie odpowiedzi ma pokrycie w cytacie ze źródła.');
  }
}

async function zakres(kontekst: KontekstBadania): Promise<void> {
  const biezaca = await wywolaj(kontekst.kanal, Command.ResearchWorkspaceGet, {});
  if (!biezaca.udany || biezaca.wynik === undefined) {
    odmowa(biezaca.blad, 'Rdzeń nie oddał przestrzeni badania, więc zakres nie ma od czego wyjść.');
    return;
  }
  const opis = zapytaj('Zakres badania', biezaca.wynik.scope);
  if (opis === '') return;
  const odbiorca = zapytaj('Odbiorca raportu', biezaca.wynik.audience ?? '');
  const protokol = zapytaj('Protokół badania; wpis prisma włącza przegląd systematyczny',
    biezaca.wynik.protocol ?? '');
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchWorkspaceSet, {
    scope: opis,
    stages: biezaca.wynik.stages,
    questions: biezaca.wynik.questions,
    audience: odbiorca === '' ? undefined : odbiorca,
    protocol: protokol === '' ? undefined : protokol,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    odmowa(wynik.blad, 'Rdzeń odmówił zapisania zakresu badania.');
    return;
  }
  powiedz(`Zakres badania zapisany; etapów: ${String(wynik.wynik.stages.length)}.`);
  await kontekst.odswiez();
}
