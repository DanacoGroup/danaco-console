// Miary debaty w oknie Roundtable: punkty zgody, zbieżność tur, sedno sporu,
// zmiana stanowisk, trafność deklaracji i skupiska opinii.
import { Command } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import { cialoPanelu, wykazPanelu } from './okno-modulu.ts';
import { odmowaCzynnosci, postawPas, przyciskPasa } from './roundtable-pas.ts';

const NAGLOWEK = 'Debata';
const PANEL = 'panel-consensus';

export function zwiazMiary(
  kanal: Kanal,
  korzen: Element,
  idOkna: () => string,
  przy: AddEventListenerOptions,
): void {
  postawPas(korzen, PANEL, [
    przyciskPasa(korzen, 'miary', 'zgoda', 'Punkty zgody'),
    przyciskPasa(korzen, 'miary', 'zbieznosc', 'Zbieżność tur'),
    przyciskPasa(korzen, 'miary', 'sedno', 'Sedno sporu'),
    przyciskPasa(korzen, 'miary', 'zmiana', 'Zmiana stanowisk'),
    przyciskPasa(korzen, 'miary', 'trafnosc', 'Trafność deklaracji'),
    przyciskPasa(korzen, 'miary', 'skupiska', 'Skupiska opinii'),
  ]);
  korzen.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const czynnosc = cel.closest<HTMLElement>('[data-miary]')?.dataset.miary;
    if (czynnosc === undefined) return;
    zdarzenie.stopPropagation();
    void wykonaj(kanal, korzen, czynnosc, idOkna());
  }, przy);
}

/* Miara wchodzi w ciało panelu zgody zamiast w ogłoszenie: wykaz punktów nie
   mieści się w jednym zdaniu, a panel i tak wraca do dowodów przy odświeżeniu. */
async function wykonaj(
  kanal: Kanal,
  korzen: Element,
  czynnosc: string,
  idOkna: string,
): Promise<void> {
  if (idOkna === '') {
    oglos(NAGLOWEK, 'Rdzeń nie dał okna debaty dla tej karty.', 'ostrzezenie');
    return;
  }
  if (czynnosc === 'zgoda') return punktyZgody(kanal, korzen, idOkna);
  if (czynnosc === 'zbieznosc') return zbieznoscTur(kanal, korzen, idOkna);
  if (czynnosc === 'sedno') return sednoSporu(kanal, korzen, idOkna);
  if (czynnosc === 'zmiana') return zmianaStanowisk(kanal, korzen, idOkna);
  if (czynnosc === 'trafnosc') return trafnoscDeklaracji(kanal, korzen, idOkna);
  if (czynnosc === 'skupiska') return skupiskaOpinii(kanal, korzen, idOkna);
  odmowaCzynnosci(NAGLOWEK, czynnosc);
}

function wypisz(
  korzen: Element,
  udany: boolean,
  blad: string | undefined,
  pozycje: readonly (readonly [string, string])[] | undefined,
  pusty: string,
): void {
  wykazPanelu(cialoPanelu(korzen, PANEL), udany, blad,
    pozycje === undefined ? undefined : [...pozycje], pusty, (x) => x);
}

async function punktyZgody(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const wynik = await wywolaj(kanal, Command.RoundtableAgreementGet, { windowId: idOkna });
  wypisz(korzen, wynik.udany, wynik.blad?.message,
    wynik.wynik?.points.map((p) => [p.text, p.agreed ? 'zgoda' : 'spór'] as const),
    'Rdzeń nie wskazał ani punktu zgody, ani punktu spornego.');
}

async function zbieznoscTur(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const wynik = await wywolaj(kanal, Command.RoundtableConvergenceGet, { windowId: idOkna });
  wypisz(korzen, wynik.udany, wynik.blad?.message,
    wynik.wynik?.points.map((p) => [`Tura ${p.turnIndex}`, p.convergence.toFixed(2)] as const),
    'Zbieżności nie ma czym policzyć — debata nie ma zamkniętych tur.');
}

async function sednoSporu(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const wynik = await wywolaj(kanal, Command.RoundtableCruxGet, { windowId: idOkna });
  const sedno = wynik.wynik?.crux;
  wypisz(korzen, wynik.udany, wynik.blad?.message,
    sedno === undefined ? [] : [[sedno.text, sedno.rationale] as const],
    'Rdzeń nie wskazał sedna sporu.');
}

async function zmianaStanowisk(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const wynik = await wywolaj(kanal, Command.RoundtableDriftGet, { windowId: idOkna });
  wypisz(korzen, wynik.udany, wynik.blad?.message,
    wynik.wynik?.entries.map((e) => [e.participantId, `${e.fromPosition} → ${e.toPosition}`] as const),
    'Żaden uczestnik nie zmienił stanowiska.');
}

async function trafnoscDeklaracji(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const wynik = await wywolaj(kanal, Command.RoundtableCalibrationGet, { windowId: idOkna });
  wypisz(korzen, wynik.udany, wynik.blad?.message,
    wynik.wynik?.calibrations.map((k) => [k.participantId,
      `deklaracja ${k.declaredConfidence.toFixed(2)} wobec trafności ${k.measuredAccuracy.toFixed(2)}`,
    ] as const),
    'Rdzeń nie ma czym zestawić deklaracji z trafnością.');
}

async function skupiskaOpinii(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const wynik = await wywolaj(kanal, Command.RoundtableClusterGet, { windowId: idOkna });
  wypisz(korzen, wynik.udany, wynik.blad?.message,
    wynik.wynik?.clusters.map((s) => [s.label, s.summary ?? ''] as const),
    'Rdzeń nie wydzielił skupisk opinii.');
}
