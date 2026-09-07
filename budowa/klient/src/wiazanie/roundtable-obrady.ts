// Prowadzenie obrad w oknie Roundtable: skład debaty, tury, głosowanie
// i stanowisko końcowe. Okno wypisywało dotąd wyłącznie wykazy rdzenia.
import {
  Command,
  RoundtableFormat,
  RoundtableTranscriptFormat,
  RoundtableVoteMethod,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';

const NAGLOWEK = 'Debata';

/*
zwiazObrady podpina pas działań debaty i czynności wierszy.

Nasłuch stoi na korzeniu karty, bo panele są wymieniane w całości przy każdym
odświeżeniu wykazu; nasłuch na wierszu zniknąłby razem z nim.
*/
export function zwiazObrady(
  kanal: Kanal,
  korzen: Element,
  idOkna: () => string,
  odswiez: () => void,
  przy: AddEventListenerOptions,
): void {
  postawPasObrad(korzen);
  korzen.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const czynnosc = cel.closest<HTMLElement>('[data-obrady]')?.dataset.obrady;
    if (czynnosc === undefined) return;
    zdarzenie.stopPropagation();
    void wykonaj(kanal, korzen, czynnosc, idOkna(), odswiez);
  }, przy);
}

/* Pas stoi nad wykazem wypowiedzi: ciało panelu jest wymieniane przy każdym
   odświeżeniu, więc pas postawiony w nim znikałby razem z wykazem. */
function postawPasObrad(korzen: Element): void {
  const panel = korzen.querySelector('#panel-debate');
  const cialo = panel?.querySelector('.sta-okno-tresc');
  if (panel === null || panel === undefined || cialo === null || cialo === undefined) return;
  const pas = korzen.ownerDocument.createElement('div');
  pas.className = 'dn-pas-dzialan';
  pas.append(
    pole(korzen, 'obrady-pytanie', 'Pytanie do uczestników'),
    przycisk(korzen, 'rozpocznij', 'Rozpocznij debatę'),
    przycisk(korzen, 'dopytaj', 'Dopytaj'),
    przycisk(korzen, 'glosowanie', 'Otwórz głosowanie'),
    przycisk(korzen, 'stanowisko', 'Stanowisko końcowe'),
    przycisk(korzen, 'wydaj', 'Wydaj transkrypt'),
  );
  panel.insertBefore(pas, cialo);
}

function pole(korzen: Element, znacznik: string, opis: string): HTMLInputElement {
  const wezel = korzen.ownerDocument.createElement('input');
  wezel.type = 'text';
  wezel.className = 'dn-pole dn-pole--sm';
  wezel.placeholder = opis;
  wezel.setAttribute('aria-label', opis);
  wezel.dataset[znacznik === 'obrady-pytanie' ? 'obradyPytanie' : znacznik] = '';
  return wezel;
}

function przycisk(korzen: Element, czynnosc: string, etykieta: string): HTMLButtonElement {
  const wezel = korzen.ownerDocument.createElement('button');
  wezel.type = 'button';
  wezel.className = 'dn-btn dn-btn--duch dn-btn--sm';
  wezel.dataset.obrady = czynnosc;
  wezel.textContent = etykieta;
  return wezel;
}

async function wykonaj(
  kanal: Kanal,
  korzen: Element,
  czynnosc: string,
  idOkna: string,
  odswiez: () => void,
): Promise<void> {
  if (idOkna === '') {
    oglos(NAGLOWEK, 'Rdzeń nie dał okna debaty dla tej karty.', 'ostrzezenie');
    return;
  }
  if (czynnosc === 'rozpocznij') return rozpocznij(kanal, korzen, idOkna, odswiez);
  if (czynnosc === 'dopytaj') return dopytaj(kanal, korzen, idOkna, odswiez);
  if (czynnosc === 'glosowanie') return otworzGlosowanie(kanal, idOkna, odswiez);
  if (czynnosc === 'stanowisko') return opiszStanowisko(kanal, idOkna);
  if (czynnosc === 'wydaj') return wydajTranskrypt(kanal, idOkna);
}

function pytanie(korzen: Element): string {
  const wezel = korzen.querySelector<HTMLInputElement>('[data-obrady-pytanie]');
  return (wezel?.value ?? '').trim();
}

/* Format i liczba tur zostają domyślne rdzenia: okno nie ma pola, którym
   Operator mógłby je wskazać, a zmyślanie wartości byłoby decyzją za niego. */
async function rozpocznij(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  odswiez: () => void,
): Promise<void> {
  const tresc = pytanie(korzen);
  if (tresc === '') {
    oglos(NAGLOWEK, 'Debata potrzebuje pytania kierowanego do uczestników.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.RoundtableDebateStart, {
    windowId: idOkna,
    question: tresc,
    format: RoundtableFormat.Structured,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił rozpoczęcia debaty.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Debata ruszyła.');
  odswiez();
}

/* Pytanie doprecyzowujące idzie do uczestnika pierwszego składu: okno nie
   prowadzi wskazania mówcy, a bez uczestnika komenda odmawia. */
async function dopytaj(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  odswiez: () => void,
): Promise<void> {
  const tresc = pytanie(korzen);
  if (tresc === '') {
    oglos(NAGLOWEK, 'Pytanie doprecyzowujące musi mieć treść.', 'ostrzezenie');
    return;
  }
  const sklad = await wywolaj(kanal, Command.RoundtableModelList, { windowId: idOkna });
  const uczestnik = sklad.wynik?.participants[0]?.id ?? '';
  if (uczestnik === '') {
    oglos(NAGLOWEK, 'Debata nie ma jeszcze uczestnika, do którego pytanie mogłoby pójść.',
      'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.RoundtableDebateFollowup, {
    windowId: idOkna,
    participantId: uczestnik,
    question: tresc,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił pytania doprecyzowującego.', 'ostrzezenie');
    return;
  }
  odswiez();
}

/* Warianty głosowania biorą się z wypowiedzi debaty: głosuje się nad tym, co
   padło, a nie nad treścią wpisaną obok. */
async function otworzGlosowanie(
  kanal: Kanal,
  idOkna: string,
  odswiez: () => void,
): Promise<void> {
  const stan = await wywolaj(kanal, Command.RoundtableDebateGet, { windowId: idOkna });
  const wypowiedzi = stan.wynik?.snapshot.statements ?? [];
  const warianty = wypowiedzi.slice(0, 5).map((w) => w.id);
  if (warianty.length === 0) {
    oglos(NAGLOWEK, 'Nie ma jeszcze wypowiedzi, nad którymi dałoby się głosować.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.RoundtableVoteStart, {
    windowId: idOkna,
    method: RoundtableVoteMethod.Approval,
    options: warianty,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił otwarcia głosowania.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Głosowanie otwarte nad ${warianty.length} wypowiedziami.`);
  odswiez();
}

async function opiszStanowisko(kanal: Kanal, idOkna: string): Promise<void> {
  const wynik = await wywolaj(kanal, Command.RoundtableConsensusGet, { windowId: idOkna });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń nie podał stanowiska końcowego.', 'ostrzezenie');
    return;
  }
  const stanowisko = wynik.wynik.consensus;
  oglos(NAGLOWEK, stanowisko.content ?? 'Stanowisko końcowe nie ma jeszcze treści.');
}

/* Transkrypt wraca jako artefakt magazynu rdzenia, nie jako treść, więc okno
   nazywa jego oznaczenie zamiast przenosić zapis do schowka. */
async function wydajTranskrypt(kanal: Kanal, idOkna: string): Promise<void> {
  const wynik = await wywolaj(kanal, Command.RoundtableTranscriptExport, {
    windowId: idOkna,
    format: RoundtableTranscriptFormat.Markdown,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wydania transkryptu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Transkrypt zapisany jako wytwór ${wynik.wynik.artifactId}.`);
}
