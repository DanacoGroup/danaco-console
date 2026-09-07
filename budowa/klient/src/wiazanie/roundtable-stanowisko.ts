// Stanowisko końcowe w oknie Roundtable: zdanie odrębne uczestnika, wykaz
// wersji stanowiska i przekazanie ustaleń do innego modułu.
import { Command, RoundtableHandoffTarget } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import {
  odmowaCzynnosci,
  polePasa,
  postawPas,
  przyciskPasa,
  wartoscPola,
} from './roundtable-pas.ts';

const NAGLOWEK = 'Debata';

export function zwiazStanowisko(
  kanal: Kanal,
  korzen: Element,
  idOkna: () => string,
  przy: AddEventListenerOptions,
): void {
  postawPas(korzen, 'panel-consensus', [
    polePasa(korzen, 'stanowiskoTresc', 'Treść zdania odrębnego'),
    przyciskPasa(korzen, 'stanowisko', 'odrebne', 'Zdanie odrębne'),
    przyciskPasa(korzen, 'stanowisko', 'wersje', 'Wersje stanowiska'),
    przyciskPasa(korzen, 'stanowisko', 'przekaz', 'Przekaż do Studio'),
  ]);
  korzen.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const czynnosc = cel.closest<HTMLElement>('[data-stanowisko]')?.dataset.stanowisko;
    if (czynnosc === undefined) return;
    zdarzenie.stopPropagation();
    void wykonaj(kanal, korzen, czynnosc, idOkna());
  }, przy);
}

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
  if (czynnosc === 'odrebne') return zdanieOdrebne(kanal, korzen, idOkna);
  if (czynnosc === 'wersje') return wersjeStanowiska(kanal, idOkna);
  if (czynnosc === 'przekaz') return przekazUstalenia(kanal, idOkna);
  odmowaCzynnosci(NAGLOWEK, czynnosc);
}

async function idStanowiska(kanal: Kanal, idOkna: string): Promise<string> {
  const wynik = await wywolaj(kanal, Command.RoundtableConsensusGet, { windowId: idOkna });
  return wynik.wynik?.consensus.id ?? '';
}

/* Zdanie odrębne przypisuje się pierwszemu uczestnikowi składu: okno nie
   prowadzi wskazania autora, a komenda bez uczestnika odmawia. */
async function zdanieOdrebne(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const tresc = wartoscPola(korzen, '[data-stanowisko-tresc]');
  if (tresc === '') {
    oglos(NAGLOWEK, 'Zdanie odrębne musi mieć treść wpisaną w polu.', 'ostrzezenie');
    return;
  }
  const stanowisko = await idStanowiska(kanal, idOkna);
  if (stanowisko === '') {
    oglos(NAGLOWEK, 'Debata nie ma jeszcze stanowiska końcowego.', 'ostrzezenie');
    return;
  }
  const sklad = await wywolaj(kanal, Command.RoundtableModelList, { windowId: idOkna });
  const uczestnik = sklad.wynik?.participants[0]?.id ?? '';
  if (uczestnik === '') {
    oglos(NAGLOWEK, 'Debata nie ma uczestnika, któremu zdanie odrębne przypadłoby.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.RoundtableConsensusMinoritySet, {
    windowId: idOkna,
    consensusId: stanowisko,
    participantId: uczestnik,
    content: tresc,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu zdania odrębnego.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Zdanie odrębne zapisane.');
}

async function wersjeStanowiska(kanal: Kanal, idOkna: string): Promise<void> {
  const wynik = await wywolaj(kanal, Command.RoundtableConsensusVersionList, { windowId: idOkna });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń nie podał wersji stanowiska.', 'ostrzezenie');
    return;
  }
  const wersje = wynik.wynik.versions;
  oglos(NAGLOWEK, wersje.length === 0
    ? 'Stanowisko końcowe nie ma jeszcze żadnej wersji.'
    : `Stanowisko ma ${wersje.length} wersji; ostatnia nosi numer ${wersje[wersje.length - 1].version}.`);
}

/* Przekazanie wraca jako wytwór magazynu, więc okno nazywa jego oznaczenie
   zamiast udawać, że przeniosło treść ustaleń. */
async function przekazUstalenia(kanal: Kanal, idOkna: string): Promise<void> {
  const stanowisko = await idStanowiska(kanal, idOkna);
  if (stanowisko === '') {
    oglos(NAGLOWEK, 'Debata nie ma stanowiska, które dałoby się przekazać.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.RoundtableConsensusHandoff, {
    windowId: idOkna,
    consensusId: stanowisko,
    target: RoundtableHandoffTarget.Studio,
    includeTranscript: true,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił przekazania ustaleń.', 'ostrzezenie');
    return;
  }
  const wytwor = wynik.wynik.handoff.artifactId ?? '';
  oglos(NAGLOWEK, wytwor === ''
    ? 'Ustalenia przekazane do Studio.'
    : `Ustalenia przekazane do Studio jako wytwór ${wytwor}.`);
}
