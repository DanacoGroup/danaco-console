// Skład nazwany w oknie Roundtable: zapis i przywołanie zespołu, rubryka
// oceny, ocena wystawiona uczestnikowi i ocena sędziowska.
import { Command, RoundtableRatingKind } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import {
  odmowaCzynnosci,
  polePasa,
  postawPas,
  potwierdzone,
  przyciskPasa,
  wartoscPola,
} from './roundtable-pas.ts';

const NAGLOWEK = 'Debata';
const SELEKTOR_NAZWY = '[data-zespoly-nazwa]';

export function zwiazZespoly(
  kanal: Kanal,
  korzen: Element,
  idOkna: () => string,
  odswiez: () => void,
  przy: AddEventListenerOptions,
): void {
  postawPas(korzen, 'panel-plan', [
    polePasa(korzen, 'zespolyNazwa', 'Nazwa zespołu albo rubryki'),
    polePasa(korzen, 'zespolyGwiazdki', 'Ocena w gwiazdkach'),
    przyciskPasa(korzen, 'zespoly', 'zapisz', 'Zapisz zespół'),
    przyciskPasa(korzen, 'zespoly', 'wnies', 'Wnieś zespół'),
    przyciskPasa(korzen, 'zespoly', 'rubryka', 'Załóż rubrykę'),
    przyciskPasa(korzen, 'zespoly', 'ocena', 'Wystaw ocenę'),
    przyciskPasa(korzen, 'zespoly', 'sedzia', 'Ocena sędziowska'),
  ]);
  korzen.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const czynnosc = cel.closest<HTMLElement>('[data-zespoly]')?.dataset.zespoly;
    if (czynnosc === undefined) return;
    zdarzenie.stopPropagation();
    void wykonaj(kanal, korzen, czynnosc, idOkna(), odswiez);
  }, przy);
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
  if (czynnosc === 'zapisz') return zapiszZespol(kanal, korzen, idOkna, odswiez);
  if (czynnosc === 'wnies') return wniesZespol(kanal, korzen, idOkna, odswiez);
  if (czynnosc === 'rubryka') return zalozRubryke(kanal, korzen, idOkna, odswiez);
  if (czynnosc === 'ocena') return wystawOcene(kanal, korzen, idOkna);
  if (czynnosc === 'sedzia') return ocenaSedziowska(kanal, idOkna);
  odmowaCzynnosci(NAGLOWEK, czynnosc);
}

async function pierwszyUczestnik(kanal: Kanal, idOkna: string): Promise<string> {
  const wynik = await wywolaj(kanal, Command.RoundtableModelList, { windowId: idOkna });
  return wynik.wynik?.participants[0]?.id ?? '';
}

async function zapiszZespol(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  odswiez: () => void,
): Promise<void> {
  const nazwa = wartoscPola(korzen, SELEKTOR_NAZWY);
  if (nazwa === '') {
    oglos(NAGLOWEK, 'Zespół zapisywany musi mieć nazwę wpisaną w polu.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.RoundtableTeamSave, {
    windowId: idOkna,
    name: nazwa,
    includeFormat: true,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu zespołu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Zespół „${nazwa}” zapisany.`);
  odswiez();
}

/* Wniesienie zespołu zastępuje skład bieżący, więc pierwsze naciśnięcie
   uzbraja przycisk zamiast nadpisywać uczestników debaty. */
async function wniesZespol(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  odswiez: () => void,
): Promise<void> {
  if (!potwierdzone(korzen, '[data-zespoly="wnies"]', 'Potwierdź wniesienie', 'Wnieś zespół')) return;
  const nazwa = wartoscPola(korzen, SELEKTOR_NAZWY);
  const wykaz = await wywolaj(kanal, Command.RoundtableTeamList, nazwa === '' ? {} : { query: nazwa });
  const zespol = wykaz.wynik?.teams[0];
  if (zespol === undefined) {
    oglos(NAGLOWEK, 'Rdzeń nie ma zapisanego zespołu o tej nazwie.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.RoundtableTeamApply, {
    windowId: idOkna,
    teamId: zespol.id,
    replace: true,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wniesienia zespołu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Skład debaty wzięty z zespołu „${zespol.name}”.`);
  odswiez();
}

/* Rubryka zakładana z okna ma jedno kryterium o pełnej wadze: okno nie
   prowadzi edytora kryteriów, a zmyślona lista byłaby decyzją za Operatora. */
async function zalozRubryke(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  odswiez: () => void,
): Promise<void> {
  const nazwa = wartoscPola(korzen, SELEKTOR_NAZWY);
  if (nazwa === '') {
    oglos(NAGLOWEK, 'Rubryka oceny musi mieć nazwę wpisaną w polu.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.RoundtableRubricSet, {
    windowId: idOkna,
    name: nazwa,
    criteria: [{ name: nazwa, weight: 1 }],
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił założenia rubryki.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Rubryka „${nazwa}” założona.`);
  odswiez();
}

async function wystawOcene(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const gwiazdki = Number(wartoscPola(korzen, '[data-zespoly-gwiazdki]'));
  if (!Number.isFinite(gwiazdki) || gwiazdki <= 0) {
    oglos(NAGLOWEK, 'Ocena wymaga liczby gwiazdek większej od zera.', 'ostrzezenie');
    return;
  }
  const uczestnik = await pierwszyUczestnik(kanal, idOkna);
  if (uczestnik === '') {
    oglos(NAGLOWEK, 'Debata nie ma uczestnika, któremu ocena mogłaby przypaść.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.RoundtableRatingSet, {
    windowId: idOkna,
    kind: RoundtableRatingKind.Star,
    targetParticipantId: uczestnik,
    stars: gwiazdki,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu oceny.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Ocena zapisana.');
}

/* Sędzią jest pierwszy uczestnik składu: okno nie prowadzi wskazania sędziego,
   a komenda bez uczestnika oceniającego odmawia. */
async function ocenaSedziowska(kanal: Kanal, idOkna: string): Promise<void> {
  const rubryki = await wywolaj(kanal, Command.RoundtableRubricList, { windowId: idOkna });
  const rubryka = rubryki.wynik?.rubrics[0];
  if (rubryka === undefined) {
    oglos(NAGLOWEK, 'Rdzeń nie ma rubryki, według której dałoby się oceniać.', 'ostrzezenie');
    return;
  }
  const uczestnik = await pierwszyUczestnik(kanal, idOkna);
  if (uczestnik === '') {
    oglos(NAGLOWEK, 'Debata nie ma uczestnika, który mógłby oceniać.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.RoundtableJudgeRun, {
    windowId: idOkna,
    rubricId: rubryka.id,
    judgeParticipantIds: [uczestnik],
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił oceny sędziowskiej.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Ocena sędziowska wykonana: ${wynik.wynik.judgements.length} orzeczeń.`);
}
