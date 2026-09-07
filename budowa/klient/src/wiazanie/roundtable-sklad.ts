// Skład debaty, ocena i stanowisko w oknie Roundtable: uczestnicy, głosy,
// analiza wypowiedzi i ranking. Okno wypisywało dotąd sam skład.
import {
  Command,
  RoundtableAnalysisKind,
  RoundtableLeaderboardScope,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';

const NAGLOWEK = 'Debata';

/* Skład i ocena stoją w panelu moderatora; pas nad ciałem panelu, bo ciało
   jest wymieniane przy każdym odświeżeniu wykazu uczestników. */
export function zwiazSklad(
  kanal: Kanal,
  korzen: Element,
  idOkna: () => string,
  odswiez: () => void,
  przy: AddEventListenerOptions,
): void {
  postawPasSkladu(korzen);
  korzen.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const czynnosc = cel.closest<HTMLElement>('[data-sklad]')?.dataset.sklad;
    if (czynnosc === undefined) return;
    zdarzenie.stopPropagation();
    void wykonaj(kanal, korzen, czynnosc, idOkna(), odswiez);
  }, przy);
}

function postawPasSkladu(korzen: Element): void {
  const panel = korzen.querySelector('#panel-moderator');
  const cialo = panel?.querySelector('.sta-okno-tresc');
  if (panel === null || panel === undefined || cialo === null || cialo === undefined) return;
  const pas = korzen.ownerDocument.createElement('div');
  pas.className = 'dn-pas-dzialan';
  const nazwa = korzen.ownerDocument.createElement('input');
  nazwa.type = 'text';
  nazwa.className = 'dn-pole dn-pole--sm';
  nazwa.placeholder = 'Nazwa tożsamości uczestnika';
  nazwa.setAttribute('aria-label', 'Nazwa tożsamości uczestnika');
  nazwa.dataset.skladNazwa = '';
  pas.append(
    nazwa,
    przycisk(korzen, 'dodaj', 'Dodaj uczestnika'),
    przycisk(korzen, 'przemianuj', 'Przemianuj pierwszego'),
    przycisk(korzen, 'usun', 'Usuń ostatniego'),
    przycisk(korzen, 'analiza', 'Analizuj wypowiedzi'),
    przycisk(korzen, 'ranking', 'Ranking uczestników'),
    przycisk(korzen, 'glos', 'Oddaj głos'),
    przycisk(korzen, 'zapisz-stanowisko', 'Zapisz stanowisko'),
  );
  panel.insertBefore(pas, cialo);
}

function przycisk(korzen: Element, czynnosc: string, etykieta: string): HTMLButtonElement {
  const wezel = korzen.ownerDocument.createElement('button');
  wezel.type = 'button';
  wezel.className = 'dn-btn dn-btn--duch dn-btn--sm';
  wezel.dataset.sklad = czynnosc;
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
  if (czynnosc === 'dodaj') return dodajUczestnika(kanal, korzen, idOkna, odswiez);
  if (czynnosc === 'przemianuj') return przemianujUczestnika(kanal, korzen, idOkna, odswiez);
  if (czynnosc === 'usun') return usunUczestnika(kanal, korzen, idOkna, odswiez);
  if (czynnosc === 'analiza') return uruchomAnalize(kanal, idOkna);
  if (czynnosc === 'ranking') return pokazRanking(kanal, idOkna);
  if (czynnosc === 'glos') return oddajGlos(kanal, idOkna);
  if (czynnosc === 'zapisz-stanowisko') return zapiszStanowisko(kanal, korzen, idOkna);
  /* Czynność spoza obsłużonych odmawia zamiast milczeć: przycisk bez gałęzi
     wyglądałby jak działający. */
  oglos(NAGLOWEK, `Czynność „${czynnosc}" nie jest prowadzona przez to okno.`, 'ostrzezenie');
}

function nazwaZPola(korzen: Element): string {
  const wezel = korzen.querySelector<HTMLInputElement>('[data-sklad-nazwa]');
  return (wezel?.value ?? '').trim();
}

async function uczestnicy(kanal: Kanal, idOkna: string): Promise<readonly { id: string }[]> {
  const wynik = await wywolaj(kanal, Command.RoundtableModelList, { windowId: idOkna });
  return wynik.wynik?.participants ?? [];
}

/* Uczestnik wchodzi na kanale modelu: bez kanału rdzeń odmawia, a okno nie ma
   pola wskazania kanału, więc bierze pierwszy czynny z wykazu rdzenia. */
async function dodajUczestnika(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  odswiez: () => void,
): Promise<void> {
  const kanaly = await wywolaj(kanal, Command.ChannelList, { enabledOnly: true });
  const idKanalu = kanaly.wynik?.channels[0]?.id ?? '';
  if (idKanalu === '') {
    oglos(NAGLOWEK, 'Rdzeń nie prowadzi żadnego czynnego kanału modelu.', 'ostrzezenie');
    return;
  }
  const nazwa = nazwaZPola(korzen);
  const wynik = await wywolaj(kanal, Command.RoundtableModelAdd, {
    windowId: idOkna,
    channelId: idKanalu,
    ...(nazwa === '' ? {} : { personaName: nazwa }),
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił dodania uczestnika.', 'ostrzezenie');
    return;
  }
  odswiez();
}

async function przemianujUczestnika(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  odswiez: () => void,
): Promise<void> {
  const nazwa = nazwaZPola(korzen);
  if (nazwa === '') {
    oglos(NAGLOWEK, 'Nowa nazwa tożsamości musi mieć treść.', 'ostrzezenie');
    return;
  }
  const sklad = await uczestnicy(kanal, idOkna);
  if (sklad.length === 0) {
    oglos(NAGLOWEK, 'Debata nie ma jeszcze uczestników.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.RoundtableModelUpdate, {
    windowId: idOkna,
    participantId: sklad[0].id,
    personaName: nazwa,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zmiany tożsamości.', 'ostrzezenie');
    return;
  }
  odswiez();
}

/* Usunięcie uczestnika jest nieodwracalne, więc pierwsze naciśnięcie uzbraja. */
async function usunUczestnika(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  odswiez: () => void,
): Promise<void> {
  const przycisk = korzen.querySelector<HTMLElement>('[data-sklad="usun"]');
  if (przycisk !== null && przycisk.dataset.uzbrojone !== 'tak') {
    przycisk.dataset.uzbrojone = 'tak';
    przycisk.textContent = 'Potwierdź usunięcie';
    return;
  }
  const sklad = await uczestnicy(kanal, idOkna);
  const ostatni = sklad[sklad.length - 1]?.id ?? '';
  if (ostatni === '') {
    oglos(NAGLOWEK, 'Debata nie ma uczestnika do usunięcia.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.RoundtableModelRemove, {
    windowId: idOkna,
    participantId: ostatni,
  });
  if (przycisk !== null) {
    delete przycisk.dataset.uzbrojone;
    przycisk.textContent = 'Usuń ostatniego';
  }
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił usunięcia uczestnika.', 'ostrzezenie');
    return;
  }
  odswiez();
}

async function uruchomAnalize(kanal: Kanal, idOkna: string): Promise<void> {
  const wynik = await wywolaj(kanal, Command.RoundtableAnalysisRun, {
    windowId: idOkna,
    kind: RoundtableAnalysisKind.ArgumentMining,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił analizy wypowiedzi.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Analiza wypowiedzi wykonana; graf argumentów odświeżony.');
}

async function pokazRanking(kanal: Kanal, idOkna: string): Promise<void> {
  const wynik = await wywolaj(kanal, Command.RoundtableLeaderboardGet, {
    scope: RoundtableLeaderboardScope.Window,
    windowId: idOkna,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń nie podał rankingu.', 'ostrzezenie');
    return;
  }
  const pozycje = wynik.wynik.entries;
  oglos(NAGLOWEK, pozycje.length === 0
    ? 'Ranking jest pusty — debata nie ma jeszcze ocen.'
    : `Ranking prowadzi ${pozycje[0].displayName} na ${pozycje.length} ocenianych.`);
}

/* Głos oddaje Operator, nie uczestnik: głosowanie wskazuje rdzeń, a wariantem
   jest wypowiedź pierwsza z wykazu głosowania. */
async function oddajGlos(kanal: Kanal, idOkna: string): Promise<void> {
  const stan = await wywolaj(kanal, Command.RoundtableVoteGet, { windowId: idOkna });
  if (!stan.udany || stan.wynik === undefined) {
    oglos(NAGLOWEK, stan.blad?.message ?? 'Żadne głosowanie nie stoi otwarte.', 'ostrzezenie');
    return;
  }
  const glosowanie = stan.wynik.vote;
  const wariant = glosowanie.options?.[0];
  if (wariant === undefined) {
    oglos(NAGLOWEK, 'Głosowanie nie ma wariantów.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.RoundtableVoteCast, {
    windowId: idOkna,
    voteId: glosowanie.id,
    voterId: 'operator',
    approvals: [wariant.id],
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił przyjęcia głosu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Głos oddany.');
}

async function zapiszStanowisko(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const tresc = nazwaZPola(korzen);
  if (tresc === '') {
    oglos(NAGLOWEK, 'Stanowisko końcowe musi mieć treść wpisaną w polu.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.RoundtableConsensusSet, {
    windowId: idOkna,
    content: tresc,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu stanowiska.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Stanowisko końcowe zapisane.');
}
