import type { RoundtableParticipant, RoundtableStatement } from '../../../../shared/contract';
import type { GlosBiezacy } from './glos-biezacy';
import { nazwaUczestnika, type StanDebaty } from './stan-debaty';

/**
 * Czytelność kilku głosów obok siebie — jedna wypowiedź, jeden blok, jedna
 * widoczna tożsamość.
 *
 * Dwa głosy tego samego modelu to dwa głosy, nie jeden: schemat debaty nie ma
 * więzu `UNIQUE(okno, kanał)` (`store/migracja_044_roundtable.sql`), a kontrakt
 * powtarza to w opisie `roundtable.model.add`. Kluczem wypowiedzi jest więc
 * `participantId`, nigdy `channelId` — i dwie tożsamości jednego kanału muszą
 * różnić się na ekranie, bo inaczej dwa głosy czyta się jako jeden. Stąd
 * znacznik mówcy przy każdej wypowiedzi i zdanie o powtórzonym kanale.
 *
 * Rozróżnienie idzie układem i gęstością, nie samą barwą. Wstęga barwna
 * (`data-rola-mowcy`) zostaje, ale sama nic nie niesie temu, kto barw nie
 * rozróżnia. Nośnikiem jest znacznik mówcy w osobnej kolumnie (te same znaki
 * co w składzie), nagłówek z pełną tożsamością przy każdej wypowiedzi oraz
 * odstęp — nowy mówca dostaje przerwę, ciąg tego samego mówcy zostaje ciasny
 * (`data-ciag`).
 *
 * Moderator nie jest uczestnikiem: rdzeń zapisuje jego interwencję kodem
 * stałym, poza wykazem uczestników, bo więz obcy do `debata_uczestnik`
 * odciąłby ją od zapisu. Jego głos dostaje więc własny rodzaj, własny znacznik
 * i układ jednokolumnowy.
 */

/**
 * Kod moderatora w zapisie tury — literał stały, nie domysł. Moderator nie ma
 * kanału ani persony, więc rdzeń podpisuje go stałą `kodModeratora`
 * (`core/adapter_modul_roundtable.go`); tę samą wartość niesie interwencja
 * moderatora (`core/adapter_modul_roundtable_moderator.go`). Uczestnicy
 * dostają kod z przedrostkiem `uczest-`, więc kolizja nie zachodzi.
 */
export const PARTICIPANT_ID_MODERATORA = 'moderator';

/** Znacznik głosu moderatora — wersaliki, żeby nie mylił się z numerem uczestnika. */
const ZNACZNIK_MODERATORA = 'MOD';

const ROLE_MOWCOW = ['pierwszy', 'drugi', 'trzeci'] as const;

/** Rodzaj głosu: uczestnik składu, moderator albo mówca składowi nieznany. */
export type RodzajGlosu = 'uczestnik' | 'moderator' | 'nieznany';

/**
 * Chwila wypowiedzi — albo ta, którą podał rdzeń, albo jawne „nie wiadomo”.
 *
 * `roundtable.debate.changed` potrafi przynieść `statement.createdAt` równe
 * zeru. Brak wiedzy nie jest treścią, więc zero i wartość niebędąca liczbą
 * dostają zdanie o braku, a nie datę z początku epoki.
 */
export function chwila(znacznik: number): string {
  if (!Number.isFinite(znacznik) || znacznik <= 0) return 'czas nieznany (rdzeń nie podał chwili)';
  return new Date(znacznik).toLocaleString('pl-PL');
}

/**
 * Rodzaj głosu. Moderatora rozpoznaje stała, a nie nieudane odnalezienie
 * w składzie: okno otwarte w trakcie debaty nie zna jeszcze całego składu —
 * odczytu `roundtable.model.list` żadne okno modułu dziś nie wywołuje — więc
 * prawdziwy uczestnik, o którym okno jeszcze nie wie, nie może dostać podpisu
 * moderatora.
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

/** Zakłada pamięć jednego rysowania listy wypowiedzi. */
export function utworzKontekstCzytelnosci(): KontekstCzytelnosci {
  return { role: new Map(), znaczniki: new Map(), poprzedniMowca: '' };
}

/**
 * Znacznik mówcy — te same znaki, którymi skład podpisuje uczestnika.
 *
 * Uczestnik znany składowi dostaje `U` i swoje miejsce w kolejności głosu, więc
 * znacznik w przebiegu debaty i znacznik w składzie to ta sama etykieta.
 * Uczestnik składowi nieznany dostaje `?` i numer wystąpienia w tej turze —
 * miejsca w składzie nie ma skąd wziąć, a zmyślenie go rozjechałoby obie listy.
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
 * Buduje blok jednej wypowiedzi wraz z widoczną tożsamością mówcy.
 *
 * @param biezacy głos tego mówcy złożony z obu dróg rdzenia — zapisu i
 *   strumienia (`glos-biezacy.ts`). Podawany wyłącznie przy wypowiedzi otwartej
 *   ostatnio przez tego mówcę: treść rosnąca należy do niej i tylko do niej,
 *   bo doklejona do wcześniejszych powtórzyłaby te same słowa. Pominięty daje
 *   blok wypowiedzi utrwalonej.
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
  // Treść bierze się z głosu bieżącego, gdy jest podany: rdzeń rozgłasza
  // wypowiedź `created` z treścią pustą i dopiero po domknięciu strumienia
  // `updated` z pełną, więc `wypowiedz.content` jest przez cały czas mówienia
  // modelu pustym napisem, a treść jedzie tymczasem drogą `stream.chunk`.
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
 * Dopiski wypowiedzi rosnącej — stan głosu, tok rozumowania, zerwanie.
 *
 * Stan głosu jest słowem, nie samą barwą: „Mówi teraz" odróżnione wyłącznie
 * odcieniem znacznika byłoby dla nierozróżniającego barw tym samym, czym
 * wypowiedź domknięta — napis stoi więc w treści węzła, a `data-stan-glosu`
 * jest dla arkusza, nie zamiast napisu.
 *
 * Tok rozumowania stoi osobno, bo rdzeń nie liczy go do treści wypowiedzi
 * (`core/adapter_modul_roundtable_glos.go` sumuje wyłącznie fragmenty `text`).
 * Doklejony do zdania dałby na żywo wypowiedź inną niż ta, którą za chwilę
 * utrwali rdzeń.
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
