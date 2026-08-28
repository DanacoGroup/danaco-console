import type { RoundtableParticipant, RoundtableStatement } from '../../../../shared/contract';
import type { GlosBiezacy } from './glos-biezacy';
import { nazwaUczestnika, type StanDebaty } from './stan-debaty';

/**
 * Kod moderatora w zapisie tury, literał stały, nie domysł — moderator nie jest uczestnikiem
 * składu, więc dostaje własny rodzaj głosu i własny znacznik jednokolumnowy.
 */
export const PARTICIPANT_ID_MODERATORA = 'moderator';

/** Znacznik głosu moderatora — wersaliki, żeby nie mylił się z numerem uczestnika w kolejności wystąpienia w turze debaty. */
const ZNACZNIK_MODERATORA = 'MOD';

const ROLE_MOWCOW = ['pierwszy', 'drugi', 'trzeci'] as const;

/** Rodzaj głosu: uczestnik znany składowi debaty, moderator rozpoznany stałym kodem, albo mówca składowi dziś nieznany. */
export type RodzajGlosu = 'uczestnik' | 'moderator' | 'nieznany';

/**
 * Chwila wypowiedzi — albo ta, którą podał rdzeń, albo jawne zdanie o braku wiedzy, gdy rdzeń
 * przyniósł znacznik czasu zerowy.
 */
export function chwila(znacznik: number): string {
  if (!Number.isFinite(znacznik) || znacznik <= 0) return 'czas nieznany (rdzeń nie podał chwili)';
  return new Date(znacznik).toLocaleString('pl-PL');
}

/**
 * Rodzaj głosu. Moderatora rozpoznaje stały kod, nie nieudane odnalezienie w składzie —
 * uczestnik jeszcze nieznany oknu nie dostaje podpisu moderatora.
 */
export function rodzajGlosu(idUczestnika: string, uczestnik: RoundtableParticipant | null): RodzajGlosu {
  if (idUczestnika === PARTICIPANT_ID_MODERATORA) return 'moderator';
  return uczestnik === null ? 'nieznany' : 'uczestnik';
}

/**
 * Nazwa mówcy w nagłówku wypowiedzi. Mówca składowi nieznany dostaje zdanie
 * o braku wiedzy okna, nie podpis moderatora — wypowiedź zostaje widoczna
 * zamiast znikać albo kłamać rolą.
 */
export function opisMowcy(idUczestnika: string, uczestnik: RoundtableParticipant | null, stan: StanDebaty): string {
  if (idUczestnika === PARTICIPANT_ID_MODERATORA) return 'Moderator debaty';
  if (uczestnik !== null) return nazwaUczestnika(uczestnik, stan.opisKanalu(uczestnik.channelId));
  return `${idUczestnika} (skład tego uczestnika jeszcze nieznany temu oknu)`;
}

/**
 * Pamięć jednego rysowania listy: role wstęgi, znaczniki mówców i mówca
 * poprzedni. Kontekst żyje tyle, co jeden przebieg rysowania — dwa okna i dwie
 * tury nie dzielą numeracji.
 */
export interface KontekstCzytelnosci {
  role: Map<string, string>;
  znaczniki: Map<string, string>;
  poprzedniMowca: string;
}

/** Zakłada pamięć jednego rysowania listy wypowiedzi: role wstęgi, znaczniki mówców i mówcę poprzedniego, od zera. */
export function utworzKontekstCzytelnosci(): KontekstCzytelnosci {
  return { role: new Map(), znaczniki: new Map(), poprzedniMowca: '' };
}

/**
 * Znacznik mówcy — te same znaki, którymi skład podpisuje uczestnika, albo numer wystąpienia,
 * gdy skład go nie zna.
 */
export function znacznikMowcy(
  idUczestnika: string,
  stan: StanDebaty,
  kontekst: KontekstCzytelnosci,
): string {
  if (idUczestnika === PARTICIPANT_ID_MODERATORA) return ZNACZNIK_MODERATORA;
  const zapamietany = kontekst.znaczniki.get(idUczestnika);
  if (zapamietany !== undefined) return zapamietany;
  const miejsce = stan.uczestnicy().findIndex((uczestnik) => uczestnik.id === idUczestnika);
  const znacznik = miejsce === -1 ? `?${kontekst.znaczniki.size + 1}` : `U${miejsce + 1}`;
  kontekst.znaczniki.set(idUczestnika, znacznik);
  return znacznik;
}

/**
 * Rola wstęgi barwnej. Przydzielana po kolejności wystąpienia w tej turze, nie
 * po składzie: wykaz oddaje przebieg, a nie spis uczestników.
 */
function rolaMowcy(idUczestnika: string, kontekst: KontekstCzytelnosci): string {
  if (idUczestnika === PARTICIPANT_ID_MODERATORA) return 'moderator';
  const znana = kontekst.role.get(idUczestnika);
  if (znana !== undefined) return znana;
  const rola = ROLE_MOWCOW[kontekst.role.size] ?? 'dalszy';
  kontekst.role.set(idUczestnika, rola);
  return rola;
}

/**
 * Zdanie o powtórzonym kanale — liczone, nigdy domyślane. `wystapieniaKanalu`
 * liczy uczestników tego kanału w składzie znanym oknu; przy jednym
 * wystąpieniu zdania nie ma, bo nie byłoby o czym mówić.
 */
function notaPowtorzonegoKanalu(uczestnik: RoundtableParticipant, stan: StanDebaty): string {
  const ile = stan.wystapieniaKanalu(uczestnik.channelId);
  if (ile <= 1) return '';
  const kanal = stan.opisKanalu(uczestnik.channelId) === '' ? uczestnik.channelId : stan.opisKanalu(uczestnik.channelId);
  return `Ten sam kanał modelu (${kanal}) niesie w tej debacie ${ile} odrębne tożsamości — to osobny głos, nie powtórzenie poprzedniego.`;
}

/**
 * Buduje blok jednej wypowiedzi wraz z widoczną tożsamością mówcy, przyjmując opcjonalnie głos
 * bieżący dla wypowiedzi otwartej.
 */
export function utworzWypowiedz(
  wypowiedz: RoundtableStatement,
  stan: StanDebaty,
  kontekst: KontekstCzytelnosci,
  biezacy?: GlosBiezacy,
): HTMLElement {
  const uczestnik = stan.uczestnik(wypowiedz.participantId);
  const rodzaj = rodzajGlosu(wypowiedz.participantId, uczestnik);

  const element = document.createElement('li');
  element.className = 'dr-wypowiedz';
  element.dataset['rolaMowcy'] = rolaMowcy(wypowiedz.participantId, kontekst);
  element.dataset['rodzajGlosu'] = rodzaj;
  element.dataset['ciag'] =
    kontekst.poprzedniMowca === wypowiedz.participantId ? 'ten-sam-mowca' : 'nowy-mowca';
  kontekst.poprzedniMowca = wypowiedz.participantId;

  const znacznik = document.createElement('span');
  znacznik.className = 'dr-wypowiedz__znacznik';
  znacznik.textContent = znacznikMowcy(wypowiedz.participantId, stan, kontekst);

  const glowa = document.createElement('div');
  glowa.className = 'dr-wypowiedz__glowa';
  const mowca = document.createElement('span');
  mowca.className = 'dr-wypowiedz__mowca';
  mowca.textContent = opisMowcy(wypowiedz.participantId, uczestnik, stan);
  const czas = document.createElement('span');
  czas.className = 'dr-wypowiedz__czas';
  czas.textContent = chwila(wypowiedz.createdAt);
  glowa.append(mowca, czas);

  const tresc = document.createElement('p');
  tresc.className = 'dr-wypowiedz__tresc';
  // Treść bierze się z głosu bieżącego, gdy podany — rdzeń inaczej niesie ją pustą do domknięcia.
  tresc.textContent = biezacy === undefined ? wypowiedz.content : biezacy.tekst;

  element.append(znacznik, glowa, tresc);

  if (biezacy !== undefined) element.append(...dopiskiGlosu(biezacy));

  const nota = uczestnik === null ? '' : notaPowtorzonegoKanalu(uczestnik, stan);
  if (nota !== '') {
    const zdanie = document.createElement('p');
    zdanie.className = 'dr-wypowiedz__nota';
    zdanie.textContent = nota;
    element.append(zdanie);
  }
  return element;
}

/**
 * Dopiski wypowiedzi rosnącej — stan głosu jako napis w treści węzła, tok rozumowania osobno
 * oraz zdanie o zerwaniu strumienia.
 */
function dopiskiGlosu(biezacy: GlosBiezacy): HTMLElement[] {
  const dopiski: HTMLElement[] = [];

  const stanGlosu = document.createElement('p');
  stanGlosu.className = 'dr-wypowiedz__stan';
  stanGlosu.dataset['stanGlosu'] = biezacy.stan;
  stanGlosu.textContent = biezacy.zeStrumienia
    ? `${biezacy.stan} — treść rośnie ze strumienia, rdzeń jeszcze jej nie utrwalił`
    : biezacy.stan;
  dopiski.push(stanGlosu);

  if (biezacy.rozumowanie !== '') {
    const rozumowanie = document.createElement('p');
    rozumowanie.className = 'dr-wypowiedz__rozumowanie';
    rozumowanie.textContent = biezacy.rozumowanie;
    dopiski.push(rozumowanie);
  }
  if (biezacy.przyczyna !== '') {
    const zerwanie = document.createElement('p');
    zerwanie.className = 'dr-wypowiedz__nota';
    zerwanie.setAttribute('role', 'note');
    zerwanie.textContent = `Strumień zerwany: ${biezacy.przyczyna}`;
    dopiski.push(zerwanie);
  }
  return dopiski;
}
