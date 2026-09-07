// Czynności obszaru badania: pytania badawcze, notatka robocza, pokrycie pytań
// źródłami i świeżość materiału. Wyniki stają w panelu Research Workspace.
import { Command, type ResearchQuestion } from '../../../shared/contract.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { niegotowyPanel, wykazPanelu } from './okno-modulu.ts';
import {
  brakOkna,
  nasluchCzynnosci,
  nieznanaCzynnosc,
  odmowa,
  powiedz,
  zapytaj,
  type KontekstBadania,
} from './research-czynnosci.ts';

const ATRYBUT_CZYNNOSCI = 'data-czynnosc-obszaru';
const GNIAZDO_WYNIKU = '[data-obszar-wynik]';

export function zwiazCzynnosciObszaru(kontekst: KontekstBadania): Odsubskrybuj {
  return nasluchCzynnosci(kontekst.korzen, ATRYBUT_CZYNNOSCI, async (kod) => {
    if (brakOkna(kontekst.idOkna())) return;
    await wykonaj(kontekst, kod);
  });
}

async function wykonaj(kontekst: KontekstBadania, kod: string): Promise<void> {
  if (kod === 'pytanie') return dopiszPytanie(kontekst);
  if (kod === 'notatka') return zapiszNotatke(kontekst);
  if (kod === 'pokrycie') return pokrycie(kontekst);
  if (kod === 'swiezosc') return swiezosc(kontekst);
  nieznanaCzynnosc(kod);
}

function gniazdo(kontekst: KontekstBadania): Element | null {
  return kontekst.korzen.querySelector(GNIAZDO_WYNIKU);
}

/* Zapis pytań idzie całym zestawem, więc nowe pytanie dopisuje się do tego, co
   rdzeń już prowadzi — inaczej zapis skasowałby pytania wcześniejsze. */
async function dopiszPytanie(kontekst: KontekstBadania): Promise<void> {
  const tresc = zapytaj('Pytanie badawcze');
  if (tresc === '') return;
  const przestrzen = await wywolaj(kontekst.kanal, Command.ResearchWorkspaceGet, {});
  if (!przestrzen.udany) {
    odmowa(przestrzen.blad, 'Rdzeń nie oddał przestrzeni badania, więc pytanie nie zostało dopisane.');
    return;
  }
  const pytania: ResearchQuestion[] = [
    ...(przestrzen.wynik?.questions ?? []),
    { id: '', text: tresc },
  ];
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchWorkspaceQuestionSet, { questions: pytania });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił zapisania pytań badawczych.');
    return;
  }
  powiedz(`Badanie prowadzi ${String(wynik.wynik?.questions.length ?? pytania.length)} pytań.`);
  await pokrycie(kontekst);
}

async function zapiszNotatke(kontekst: KontekstBadania): Promise<void> {
  const przestrzen = await wywolaj(kontekst.kanal, Command.ResearchWorkspaceGet, {});
  const tresc = zapytaj('Notatka robocza badania', przestrzen.wynik?.note ?? '');
  if (tresc === '') return;
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchWorkspaceNoteSet, { content: tresc });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił zapisania notatki badania.');
    return;
  }
  powiedz('Notatka robocza zapisana.');
}

async function pokrycie(kontekst: KontekstBadania): Promise<void> {
  const cialo = gniazdo(kontekst);
  if (cialo === null) return;
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchWorkspaceCoverage, {
    windowId: kontekst.idOkna(),
  });
  wykazPanelu(cialo, wynik.udany, wynik.blad?.message, wynik.wynik?.questions,
    'Badanie nie prowadzi jeszcze żadnego pytania badawczego.',
    (pytanie) => [pytanie.text, `${String(pytanie.sourceIds?.length ?? 0)} źródeł`] as const);
  if (wynik.udany && wynik.wynik !== undefined) {
    powiedz(`Pytań bez źródła: ${String(wynik.wynik.uncoveredCount)}.`);
  }
}

async function swiezosc(kontekst: KontekstBadania): Promise<void> {
  const cialo = gniazdo(kontekst);
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchWorkspaceFreshness, {
    windowId: kontekst.idOkna(),
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    niegotowyPanel(cialo, wynik.blad?.message ?? 'Świeżość materiału nie doszła.');
    odmowa(wynik.blad, 'Rdzeń odmówił oceny świeżości materiału.');
    return;
  }
  const zdanie = `Od ostatniego pozyskania minęło ${String(wynik.wynik.staleDays)} dni; `
    + (wynik.wynik.refreshSuggested ? 'rdzeń sugeruje odświeżenie.' : 'rdzeń nie sugeruje odświeżenia.');
  niegotowyPanel(cialo, zdanie);
  powiedz(zdanie);
}
