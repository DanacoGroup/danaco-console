// Prowadzenie obrad w oknie Roundtable: interwencja moderatora, zamknięcie
// tury, wariant tury, powtórzenie wypowiedzi i szablony moderowania.
import { Command, ModeratorAction } from '../../../shared/contract.ts';
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
const SELEKTOR_TRESCI = '[data-prowadzenie-tresc]';

export function zwiazProwadzenie(
  kanal: Kanal,
  korzen: Element,
  idOkna: () => string,
  odswiez: () => void,
  przy: AddEventListenerOptions,
): void {
  postawPas(korzen, 'panel-moderator', [
    polePasa(korzen, 'prowadzenieTresc', 'Treść interwencji albo nazwa'),
    przyciskPasa(korzen, 'prowadzenie', 'ukierunkuj', 'Ukierunkuj debatę'),
    przyciskPasa(korzen, 'prowadzenie', 'zamknij', 'Zamknij turę'),
    przyciskPasa(korzen, 'prowadzenie', 'wariant', 'Załóż wariant tury'),
    przyciskPasa(korzen, 'prowadzenie', 'powtorz', 'Powtórz wypowiedź'),
    przyciskPasa(korzen, 'prowadzenie', 'szablony', 'Szablony moderowania'),
    przyciskPasa(korzen, 'prowadzenie', 'zapisz-szablon', 'Zapisz szablon'),
  ]);
  korzen.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const czynnosc = cel.closest<HTMLElement>('[data-prowadzenie]')?.dataset.prowadzenie;
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
  if (czynnosc === 'ukierunkuj') return ukierunkuj(kanal, korzen, idOkna, odswiez);
  if (czynnosc === 'zamknij') return zamknijTure(kanal, idOkna, odswiez);
  if (czynnosc === 'wariant') return zalozWariant(kanal, korzen, idOkna, odswiez);
  if (czynnosc === 'powtorz') return powtorzWypowiedz(kanal, korzen, idOkna, odswiez);
  if (czynnosc === 'szablony') return wypiszSzablony(kanal);
  if (czynnosc === 'zapisz-szablon') return zapiszSzablon(kanal, korzen, idOkna);
  odmowaCzynnosci(NAGLOWEK, czynnosc);
}

/* Czynności moderatora dotyczą tury ostatniej: okno nie prowadzi wyboru tury,
   a komendy bez wskazania tury odnoszą się do stanu, którego już nie widać. */
async function ostatniaTura(kanal: Kanal, idOkna: string): Promise<string> {
  const stan = await wywolaj(kanal, Command.RoundtableDebateGet, { windowId: idOkna });
  const tury = stan.wynik?.snapshot.turns ?? [];
  return tury[tury.length - 1]?.id ?? '';
}

async function ostatniaWypowiedz(kanal: Kanal, idOkna: string): Promise<string> {
  const stan = await wywolaj(kanal, Command.RoundtableDebateGet, { windowId: idOkna });
  const wypowiedzi = stan.wynik?.snapshot.statements ?? [];
  return wypowiedzi[wypowiedzi.length - 1]?.id ?? '';
}

async function ukierunkuj(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  odswiez: () => void,
): Promise<void> {
  const tresc = wartoscPola(korzen, SELEKTOR_TRESCI);
  if (tresc === '') {
    oglos(NAGLOWEK, 'Interwencja moderująca musi mieć treść.', 'ostrzezenie');
    return;
  }
  const tura = await ostatniaTura(kanal, idOkna);
  const wynik = await wywolaj(kanal, Command.RoundtableModeratorDirect, {
    windowId: idOkna,
    action: ModeratorAction.Direct,
    message: tresc,
    ...(tura === '' ? {} : { turnId: tura }),
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił interwencji moderującej.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Interwencja moderująca przyjęta.');
  odswiez();
}

async function zamknijTure(kanal: Kanal, idOkna: string, odswiez: () => void): Promise<void> {
  const tura = await ostatniaTura(kanal, idOkna);
  if (tura === '') {
    oglos(NAGLOWEK, 'Debata nie ma tury, którą dałoby się zamknąć.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.RoundtableModeratorDirect, {
    windowId: idOkna,
    action: ModeratorAction.CloseTurn,
    turnId: tura,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zamknięcia tury.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Tura zamknięta.');
  odswiez();
}

async function zalozWariant(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  odswiez: () => void,
): Promise<void> {
  const tura = await ostatniaTura(kanal, idOkna);
  if (tura === '') {
    oglos(NAGLOWEK, 'Debata nie ma tury, którą dałoby się rozgałęzić.', 'ostrzezenie');
    return;
  }
  const nazwa = wartoscPola(korzen, SELEKTOR_TRESCI);
  const wynik = await wywolaj(kanal, Command.RoundtableDebateBranch, {
    windowId: idOkna,
    turnId: tura,
    ...(nazwa === '' ? {} : { label: nazwa }),
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił założenia wariantu tury.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Wariant tury założony; obie gałęzie zostają w zapisie debaty.');
  odswiez();
}

/* Powtórzenie zastępuje wypowiedź odpowiedzią nową, więc pierwsze naciśnięcie
   uzbraja przycisk zamiast kasować to, co padło. */
async function powtorzWypowiedz(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  odswiez: () => void,
): Promise<void> {
  if (!potwierdzone(korzen, '[data-prowadzenie="powtorz"]', 'Potwierdź powtórzenie',
    'Powtórz wypowiedź')) return;
  const wypowiedz = await ostatniaWypowiedz(kanal, idOkna);
  if (wypowiedz === '') {
    oglos(NAGLOWEK, 'Debata nie ma wypowiedzi do powtórzenia.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.RoundtableStatementRegenerate, {
    windowId: idOkna,
    statementId: wypowiedz,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił powtórzenia wypowiedzi.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Wypowiedź zastąpiona odpowiedzią nową.');
  odswiez();
}

async function wypiszSzablony(kanal: Kanal): Promise<void> {
  const wynik = await wywolaj(kanal, Command.RoundtableModerationTemplateList, {});
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń nie podał szablonów moderowania.', 'ostrzezenie');
    return;
  }
  const szablony = wynik.wynik.templates;
  oglos(NAGLOWEK, szablony.length === 0
    ? 'Rdzeń nie ma zapisanego szablonu moderowania.'
    : `Szablony moderowania: ${szablony.map((s) => s.name).join(', ')}.`);
}

async function zapiszSzablon(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const nazwa = wartoscPola(korzen, SELEKTOR_TRESCI);
  if (nazwa === '') {
    oglos(NAGLOWEK, 'Szablon moderowania musi mieć nazwę wpisaną w polu.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.RoundtableModerationTemplateSave, {
    windowId: idOkna,
    name: nazwa,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu szablonu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Szablon moderowania „${nazwa}” zapisany.`);
}
