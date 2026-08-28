import type { RoundtableParticipant } from '../../../../shared/contract';
import { przyciskAkcji as przycisk } from '../../modele/kontrolki-formularza-braki';
import { powodBrakuCzynnosci } from './braki-kontraktu';
import {
  utworzKontekstCzytelnosci,
  znacznikMowcy,
  type KontekstCzytelnosci,
} from './czytelnosc-glosow';
import { nazwaUczestnika, type StanDebaty } from './stan-debaty';
import type { StanTresci } from './stany-okna';

/** Skład debaty i tura bieżąca w Moderator Panelu: obsługa czynności wiersza uczestnika, wnoszona przez wytwórnię okna. */
export interface ObslugaSkladu {
  przelaczWyciszenie: (uczestnik: RoundtableParticipant) => void;
  przesunKolejnosc: (idUczestnika: string, kierunek: -1 | 1) => void;
}

/** Co o uczestniku wiadomo z wypowiedzi tej tury — zmierzone, nie założone: ostatni mówca, mówiący i liczba znanych wypowiedzi. */
interface GlosyTury {
  ostatni: string;
  mowiacy: Set<string>;
  znane: number;
}

/**
 * Rysuje turę bieżącą i skład debaty, albo stan pusty, gdy debata nie ma otwartej tury —
 * stan poprawny, nie błąd.
 */
export function rysujTuraISklad(
  stan: StanDebaty,
  tresc: StanTresci,
  obsluga: ObslugaSkladu,
): void {
  const tura = stan.definicjaTury();
  if (tura === null) {
    tresc.pusto('Debata nie ma otwartej tury — uruchom ją w Debate Panelu.');
    return;
  }

  const miejsce = tresc.tresc();
  miejsce.append(naglowekTury(tura));

  const uczestnicy = uczestnicyWKolejnosci(stan);
  if (uczestnicy.length === 0) {
    const brak = document.createElement('p');
    brak.className = 'dr-stan';
    brak.dataset['stan'] = 'pusto';
    brak.textContent = 'Debata nie ma jeszcze uczestników — dodaj ich w Model Panels.';
    miejsce.append(brak);
    return;
  }

  const glosy = glosyTury(stan);
  miejsce.append(notaKolejnosci(stan, uczestnicy), notaGlosuBiezacego(glosy));

  const kontekst = utworzKontekstCzytelnosci();
  const sklad = document.createElement('div');
  sklad.className = 'dr-sklad';
  uczestnicy.forEach((uczestnik, indeks) => {
    sklad.append(wierszUczestnika(uczestnik, indeks, stan, glosy, kontekst, obsluga));
  });
  miejsce.append(sklad);
}

/** Nagłówek tury bieżącej: numer, stan i — gdy jest — zagadnienie, wypisane jednym zdaniem opisowym tury. */
function naglowekTury(tura: { index: number; status: string; topic?: string }): HTMLElement {
  const element = document.createElement('div');
  element.className = 'dr-tura';
  element.dataset['stanTury'] = tura.status;
  const zagadnienie = tura.topic !== undefined && tura.topic !== '' ? ` · ${tura.topic}` : '';
  element.textContent = `Tura #${tura.index} · ${tura.status}${zagadnienie}`;
  return element;
}

/**
 * Skład debaty w kolejności głosu wskazanej przez rdzeń; uczestnicy bez wskazanej kolejności
 * schodzą na koniec.
 */
export function uczestnicyWKolejnosci(stan: StanDebaty): RoundtableParticipant[] {
  return [...stan.uczestnicy()].sort((a, b) => {
    if (a.order === undefined && b.order === undefined) return 0;
    if (a.order === undefined) return 1;
    if (b.order === undefined) return -1;
    return a.order - b.order;
  });
}

/** Głosy tej tury policzone z wypowiedzi, które okno widziało, nie z założenia o pełnym udziale składu. */
function glosyTury(stan: StanDebaty): GlosyTury {
  const wypowiedzi = stan.wypowiedzi();
  const mowiacy = new Set(wypowiedzi.map((wypowiedz) => wypowiedz.participantId));
  const ostatnia = wypowiedzi[wypowiedzi.length - 1];
  return {
    ostatni: ostatnia === undefined ? '' : ostatnia.participantId,
    mowiacy,
    znane: wypowiedzi.length,
  };
}

/**
 * Skąd bierze się porządek wykazu. Rdzeń albo podał `order` każdemu, albo wykaz
 * stoi na kolejności dodania — od tego zależy, czy strzałki przestawiają coś, co
 * rdzeń już zna, czy ustalają to pierwszy raz.
 */
function notaKolejnosci(stan: StanDebaty, uczestnicy: readonly RoundtableParticipant[]): HTMLElement {
  const zKolejnoscia = uczestnicy.filter((uczestnik) => uczestnik.order !== undefined).length;
  const nota = document.createElement('p');
  nota.className = 'dr-sklad__nota';
  if (zKolejnoscia === uczestnicy.length) {
    nota.textContent = `Kolejność głosu wskazał rdzeń — ma ją każdy z ${uczestnicy.length} uczestników składu znanego oknu.`;
    return nota;
  }
  if (zKolejnoscia === 0) {
    nota.textContent = `Rdzeń nie wskazał kolejności głosu nikomu ze składu (${uczestnicy.length}) — wykaz stoi na kolejności dodania uczestników.`;
    return nota;
  }
  nota.textContent =
    `Rdzeń wskazał kolejność głosu ${zKolejnoscia} z ${uczestnicy.length} uczestników — ` +
    `ci z kolejnością stoją wyżej, reszta w kolejności dodania. Okno debaty: ${stan.okno()}.`;
  return nota;
}

/** Zdanie o głosie bieżącym: mierzalny jest ostatni zapisany, nie trwający, którego kontrakt nie niesie. */
function notaGlosuBiezacego(glosy: GlosyTury): HTMLElement {
  const nota = document.createElement('p');
  nota.className = 'dr-sklad__nota';
  nota.dataset['brakSygnalu'] = 'tak';
  const zdanie =
    glosy.znane === 0
      ? 'Żadna wypowiedź tej tury nie dotarła jeszcze do tego okna, więc ostatniego głosu nie ma jak wskazać.'
      : `Ostatni głos w turze jest oznaczony w wykazie; okno widziało ${glosy.znane} wypowiedzi od otwarcia.`;
  nota.textContent = powodBrakuCzynnosci(
    `${zdanie} Głosu TRWAJĄCEGO nie da się pokazać: ani uczestnik, ani tura, ani wypowiedź nie mają w kontrakcie pola o mówieniu w toku — mierzalna jest wyłącznie wypowiedź już zapisana.`,
  );
  return nota;
}

/** Wiersz uczestnika: znacznik, tożsamość, stany słowem i czynności moderatora — wyciszenie i przesunięcie kolejności. */
function wierszUczestnika(
  uczestnik: RoundtableParticipant,
  indeks: number,
  stan: StanDebaty,
  glosy: GlosyTury,
  kontekst: KontekstCzytelnosci,
  obsluga: ObslugaSkladu,
): HTMLElement {
  const opisKanalu = stan.opisKanalu(uczestnik.channelId);
  const wyciszony = uczestnik.muted === true;
  const powtorzony = stan.wystapieniaKanalu(uczestnik.channelId) > 1;
  const ostatni = glosy.ostatni === uczestnik.id;

  const kolejnosc = document.createElement('span');
  kolejnosc.className = 'dr-uczestnik__kolejnosc';
  kolejnosc.textContent = znacznikMowcy(uczestnik.id, stan, kontekst);
  kolejnosc.title = 'Ten sam znacznik stoi nad wypowiedziami tego uczestnika w Debate Panelu.';

  const tozsamosc = document.createElement('div');
  tozsamosc.className = 'dr-tozsamosc';
  const nazwa = document.createElement('strong');
  nazwa.textContent = nazwaUczestnika(uczestnik, opisKanalu);
  tozsamosc.append(nazwa);

  const pola = document.createElement('div');
  pola.className = 'dr-uczestnik__pola';
  pola.append(tozsamosc, stanyUczestnika(indeks, wyciszony, powtorzony, ostatni, glosy.mowiacy.has(uczestnik.id)));

  const akcje = document.createElement('div');
  akcje.className = 'dr-uczestnik__akcje';

  // Strzałki stoją bez wygaszenia na skrajnych wierszach; przesunięcie samo odrzuca ruch poza wykaz.
  const gora = przycisk('↑', 'dn-btn dn-btn--zarys');
  gora.setAttribute('aria-label', `Przesuń ${nazwa.textContent} wyżej w kolejności głosu`);
  gora.addEventListener('click', () => obsluga.przesunKolejnosc(uczestnik.id, -1));

  const dol = przycisk('↓', 'dn-btn dn-btn--zarys');
  dol.setAttribute('aria-label', `Przesuń ${nazwa.textContent} niżej w kolejności głosu`);
  dol.addEventListener('click', () => obsluga.przesunKolejnosc(uczestnik.id, 1));

  const wyciszenie = przycisk(wyciszony ? 'Odcisz' : 'Wycisz', 'dn-btn dn-btn--zarys');
  wyciszenie.setAttribute(
    'aria-label',
    `${wyciszony ? 'Przywróć głos' : 'Wycisz w turze'}: ${nazwa.textContent}`,
  );
  wyciszenie.addEventListener('click', () => obsluga.przelaczWyciszenie(uczestnik));

  akcje.append(gora, dol, wyciszenie);

  const panel = document.createElement('div');
  panel.className = 'dr-uczestnik';
  if (wyciszony) panel.dataset['wyciszony'] = 'tak';
  if (powtorzony) panel.dataset['kanalPowtorzony'] = 'tak';
  if (ostatni) panel.dataset['ostatniGlos'] = 'tak';
  panel.append(kolejnosc, pola, akcje);
  return panel;
}

/**
 * Stany uczestnika wypisane słowem, nie samym atrybutem danych czy krojem pisma.
 * Wyciszenie mówi, czego dotyczy: uczestnik zostaje w składzie, pytanie tury do
 * niego nie idzie (kolumna `wyciszony` w `migracja_044_roundtable.sql`).
 */
function stanyUczestnika(
  indeks: number,
  wyciszony: boolean,
  powtorzony: boolean,
  ostatni: boolean,
  mowil: boolean,
): HTMLElement {
  const stany = document.createElement('ul');
  stany.className = 'dr-uczestnik__stany';

  const opisy = [
    `Miejsce #${indeks + 1} w kolejności głosu`,
    wyciszony ? 'Wyciszony w turze — zostaje w składzie, pytanie do niego nie idzie' : 'Słyszany w turze',
    ostatni ? 'Ostatni głos w tej turze' : mowil ? 'Zabierał głos w tej turze' : 'Bez wypowiedzi widzianej przez to okno',
  ];
  if (powtorzony) opisy.push('Ten sam kanał modelu pod inną tożsamością');

  for (const opis of opisy) {
    const pozycja = document.createElement('li');
    pozycja.className = 'dr-uczestnik__stan';
    pozycja.textContent = opis;
    stany.append(pozycja);
  }
  return stany;
}
