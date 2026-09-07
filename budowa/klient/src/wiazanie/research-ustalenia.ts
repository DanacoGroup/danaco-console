// Czynności ustaleń: dodanie, zmiana, usunięcie, powiązanie ze źródłem,
// sprawdzenie twierdzenia i wykaz sprzeczności między ustaleniami okna.
import { Command, type ResearchFinding } from '../../../shared/contract.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import {
  brakOkna,
  nasluchCzynnosci,
  nieznanaCzynnosc,
  odmowa,
  potwierdzone,
  powiedz,
  wierszCzynnosci,
  wskazanie,
  wybierz,
  zapytaj,
  type CzynnoscWiersza,
  type KontekstBadania,
} from './research-czynnosci.ts';

const ATRYBUT_USTALENIA = 'data-ustalenie-badania';
const ATRYBUT_CZYNNOSCI = 'data-czynnosc-ustalenia';

const CZYNNOSCI_WIERSZA: readonly CzynnoscWiersza[] = [
  { kod: 'zmien', etykieta: 'Zmień' },
  { kod: 'powiaz', etykieta: 'Powiąż źródło' },
  { kod: 'sprawdz', etykieta: 'Sprawdź' },
  { kod: 'usun', etykieta: 'Usuń' },
];

export function wierszUstalenia(cialo: Element, ustalenie: ResearchFinding): HTMLElement {
  return wierszCzynnosci(cialo, ustalenie.content, String(ustalenie.status),
    ATRYBUT_USTALENIA, ustalenie.id, ATRYBUT_CZYNNOSCI, CZYNNOSCI_WIERSZA);
}

export function zwiazCzynnosciUstalen(kontekst: KontekstBadania): Odsubskrybuj {
  return nasluchCzynnosci(kontekst.korzen, ATRYBUT_CZYNNOSCI, async (kod, przycisk) => {
    if (brakOkna(kontekst.idOkna())) return;
    await wykonaj(kontekst, kod, wskazanie(przycisk, ATRYBUT_USTALENIA));
  });
}

async function wykonaj(kontekst: KontekstBadania, kod: string, idUstalenia: string): Promise<void> {
  if (kod === 'dodaj') return dodaj(kontekst);
  if (kod === 'sprzecznosci') return sprzecznosci(kontekst);
  if (kod === 'zmien') return zmien(kontekst, idUstalenia);
  if (kod === 'powiaz') return powiaz(kontekst, idUstalenia);
  if (kod === 'sprawdz') return sprawdz(kontekst, idUstalenia);
  if (kod === 'usun') return usun(kontekst, idUstalenia);
  nieznanaCzynnosc(kod);
}

function brakUstalenia(idUstalenia: string): boolean {
  if (idUstalenia !== '') return false;
  odmowa(undefined, 'Przycisk nie stoi przy żadnym ustaleniu, więc nic nie zostało zmienione.');
  return true;
}

async function dodaj(kontekst: KontekstBadania): Promise<void> {
  const tresc = zapytaj('Treść ustalenia');
  if (tresc === '') return;
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchFindingAdd, {
    windowId: kontekst.idOkna(),
    content: tresc,
  });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił zapisania ustalenia.');
    return;
  }
  powiedz('Ustalenie zapisane.');
  await kontekst.odswiez();
}

async function zmien(kontekst: KontekstBadania, idUstalenia: string): Promise<void> {
  if (brakUstalenia(idUstalenia)) return;
  const biezace = await ustalenie(kontekst.kanal, kontekst.idOkna(), idUstalenia);
  const tresc = zapytaj('Nowa treść ustalenia', biezace?.content ?? '');
  if (tresc === '') return;
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchFindingUpdate, {
    findingId: idUstalenia,
    content: tresc,
  });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił zmiany ustalenia.');
    return;
  }
  powiedz('Ustalenie zmienione.');
  await kontekst.odswiez();
}

async function usun(kontekst: KontekstBadania, idUstalenia: string): Promise<void> {
  if (brakUstalenia(idUstalenia)) return;
  if (!potwierdzone(`ustalenie:${idUstalenia}`,
    'Usunięcie ustalenia jest nieodwracalne. Naciśnij drugi raz, żeby je wykonać.')) return;
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchFindingRemove, {
    findingId: idUstalenia,
  });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił usunięcia ustalenia.');
    return;
  }
  powiedz(`Ustalenie usunięte; sekcji raportu, które je utraciły: ${String(wynik.wynik?.detachedSections ?? 0)}.`);
  await kontekst.odswiez();
}

/* Kontrakt nie ma komendy dowiązania samego źródła: powiązania idą zapisem
   ustalenia, w którym treść jest wymagana — dlatego bierze się ją z wykazu. */
async function powiaz(kontekst: KontekstBadania, idUstalenia: string): Promise<void> {
  if (brakUstalenia(idUstalenia)) return;
  const biezace = await ustalenie(kontekst.kanal, kontekst.idOkna(), idUstalenia);
  if (biezace === undefined) {
    odmowa(undefined, 'Rdzeń nie oddał treści ustalenia, a zapis bez niej nie przejdzie.');
    return;
  }
  const zrodla = await wywolaj(kontekst.kanal, Command.ResearchSourceList, {
    windowId: kontekst.idOkna(),
  });
  if (!zrodla.udany || zrodla.wynik === undefined) {
    odmowa(zrodla.blad, 'Wykaz źródeł nie doszedł, więc nie ma czym powiązać ustalenia.');
    return;
  }
  const wykaz = zrodla.wynik.sources;
  const numer = wybierz('Źródło do powiązania', wykaz.map((zrodlo) => zrodlo.title));
  const wybrane = numer < 0 ? undefined : wykaz[numer];
  if (wybrane === undefined) {
    odmowa(undefined, 'Żadne źródło nie zostało wskazane, więc powiązanie nie poszło do rdzenia.');
    return;
  }
  const powiazane = new Set(biezace.sourceIds ?? []);
  powiazane.add(wybrane.id);
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchFindingAdd, {
    windowId: kontekst.idOkna(),
    findingId: idUstalenia,
    content: biezace.content,
    sourceIds: [...powiazane],
  });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił powiązania ustalenia ze źródłem.');
    return;
  }
  powiedz(`Ustalenie opiera się teraz na ${String(powiazane.size)} źródłach.`);
  await kontekst.odswiez();
}

async function sprawdz(kontekst: KontekstBadania, idUstalenia: string): Promise<void> {
  if (brakUstalenia(idUstalenia)) return;
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchFindingFactCheck, {
    findingId: idUstalenia,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    odmowa(wynik.blad, 'Rdzeń odmówił weryfikacji ustalenia.');
    return;
  }
  powiedz(`${String(wynik.wynik.result.verdict)}: ${wynik.wynik.result.rationale}`);
}

async function sprzecznosci(kontekst: KontekstBadania): Promise<void> {
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchFindingContradictions, {
    windowId: kontekst.idOkna(),
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    odmowa(wynik.blad, 'Rdzeń odmówił wykrycia sprzeczności.');
    return;
  }
  const wykryte = wynik.wynik.contradictions;
  if (wykryte.length === 0) {
    powiedz('Rdzeń nie wykrył sprzeczności między ustaleniami tego okna.');
    return;
  }
  powiedz(`Sprzeczności: ${String(wykryte.length)}. ${wykryte.map((jedna) => jedna.summary).join(' | ')}`);
}

async function ustalenie(
  kanal: Kanal,
  idOkna: string,
  idUstalenia: string,
): Promise<ResearchFinding | undefined> {
  const wynik = await wywolaj(kanal, Command.ResearchFindingList, { windowId: idOkna });
  if (!wynik.udany) return undefined;
  return wynik.wynik?.findings.find((jedno) => jedno.id === idUstalenia);
}
