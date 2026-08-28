import type { StanPolaczenia } from '../polaczenie/stan-polaczenia';

/**
 * Stan łączności okna: odczyt transportu.
 */

/**
 * Odczyt transportu łączności okna: tyle, ile transport wystawia dziś każdemu pytającemu okna o jego stan.
 */
export interface OdczytLacznosci {
  stan: StanPolaczenia;
  /** Ramki czekające w kolejce wychodzącej. */
  oczekujace: number;
}

/**
 * Dojście do przebiegu ponowienia; port jest opcjonalny, bo interfejs transportu go nie wystawia wprost.
 */
export interface PortPonawiania {
  /** Numer próby, którą transport podejmie jako następną (od 1). */
  numerProby(): number;
  /** Milisekundy do zaplanowanej próby; wartość ujemna albo 0 = próba trwa. */
  zaPonowieniem(): number;
  /** Podejmuje próbę natychmiast, bez czekania na zaplanowany odstęp. */
  ponowTeraz(): void;
}

/**
 * Pełny stan łączności okna: odczyt transportu wraz z całym przebiegiem ponowienia próby połączenia z rdzeniem.
 */
export interface StanLacznosciOkna extends OdczytLacznosci {
  /** Numer próby; `null` znaczy „transport go nie wystawia". */
  numerProby: number | null;
  /** Milisekundy do następnej próby; `null` znaczy „transport ich nie wystawia". */
  zaPonowieniemMs: number | null;
}

/**
 * Nazwa stanu łączności widoczna w interfejsie okna — ta sama, którą niesie pasek górny całej aplikacji.
 */
const NAZWY: Readonly<Record<StanPolaczenia, string>> = {
  rozlaczony: 'Rozłączony',
  laczenie: 'Łączenie',
  polaczony: 'Połączony',
  ponawianie: 'Ponawianie',
};

/**
 * Odmiana kropki plakietki łączności pochodząca ze wspólnego arkusza stylu aplikacji; bez barw własnych.
 */
const KROPKI: Readonly<Record<StanPolaczenia, string>> = {
  rozlaczony: 'dn-kropka--blad',
  laczenie: 'dn-kropka--tetno',
  polaczony: 'dn-kropka--sukces',
  ponawianie: 'dn-kropka--ostrzezenie',
};

/**
 * Odmiana plakietki łączności okna — stan łączności nigdy nie idzie samą barwą bez tekstu opisowego stanu.
 */
const PLAKIETKI: Readonly<Record<StanPolaczenia, string>> = {
  rozlaczony: 'dn-plakietka--blad',
  laczenie: 'dn-plakietka--informacja',
  polaczony: 'dn-plakietka--sukces',
  ponawianie: 'dn-plakietka--ostrzezenie',
};

/**
 * Stany łączności okna, w których połączenie jest w toku — nośnikiem stanu jest wtedy wirujący spinner.
 */
const W_TOKU: ReadonlySet<StanPolaczenia> = new Set<StanPolaczenia>(['laczenie', 'ponawianie']);

/**
 * Czy plakietka ma być widoczna.
 *
 * Połączenie bez kolejki nie dokłada niczego do nagłówka: cztery zielone kropki
 * obok siebie na scenie czterech okien to szum, w którym jedna czerwona ginie.
 * Plakietka pokazuje się wtedy, gdy ma co powiedzieć.
 */
export function czyWidocznaLacznosc(stan: StanLacznosciOkna): boolean {
  return stan.stan !== 'polaczony' || stan.oczekujace > 0;
}

/**
 * Czy stan łączności okna jest w toku — spinner zastępuje wtedy kropkę, tak jak w pasku górnym aplikacji.
 */
export function czyLacznoscWToku(stan: StanPolaczenia): boolean {
  return W_TOKU.has(stan);
}

/**
 * Klasa arkusza stylu dla samej kropki tego stanu łączności widocznej w nagłówku okna komunikacji rdzenia.
 */
export function wariantKropkiLacznosci(stan: StanPolaczenia): string {
  return KROPKI[stan];
}

/**
 * Klasa arkusza stylu dla całej plakietki tego stanu łączności widocznej w nagłówku danego okna komunikacji.
 */
export function wariantPlakietkiLacznosci(stan: StanPolaczenia): string {
  return PLAKIETKI[stan];
}

/**
 * Nazwa stanu łączności — jedno źródło napisu dla nagłówka okna oraz dla podpowiedzi tej plakietki stanu.
 */
export function nazwaStanuLacznosci(stan: StanPolaczenia): string {
  return NAZWY[stan];
}

/**
 * Napis plakietki: stan, przebieg ponowienia i długość kolejki.
 *
 * Każdy człon staje tylko wtedy, gdy jest prawdziwy. Człony „próba N" i „za N s"
 * znikają w całości, gdy transport ich nie wystawia — zero udawałoby wtedy
 * odczyt, którego nie ma.
 */
export function napisLacznosci(stan: StanLacznosciOkna): string {
  const czlony = [NAZWY[stan.stan]];

  if (stan.stan === 'ponawianie' && stan.numerProby !== null) {
    czlony.push(`próba ${stan.numerProby}`);
  }
  if (stan.stan === 'ponawianie' && stan.zaPonowieniemMs !== null) {
    czlony.push(`za ${sekundyDoProby(stan.zaPonowieniemMs)}`);
  }
  if (stan.oczekujace > 0) {
    czlony.push(`w kolejce ${stan.oczekujace}`);
  }
  return czlony.join(' · ');
}

/**
 * Zapis odliczania czasu w sekundach zaokrąglonych w górę; poniżej sekundy próba połączenia już biegnie.
 */
function sekundyDoProby(ms: number): string {
  if (ms <= 0) return 'teraz';
  return `${Math.ceil(ms / 1000)} s`;
}

/**
 * Opis polityki ponawiania dla podpowiedzi plakietki łączności; wartości pochodzą z modułu ponawiania.
 */
export const OPIS_PONAWIANIA =
  'Ponawianie jest wykładnicze: pierwsza próba po 500 ms, każda następna dwa '
  + 'razy dalej, nie dalej niż co 15 s, z rozproszeniem losowym do 25 %. '
  + 'Liczba prób nie jest ograniczona — klient nie poddaje się sam.';

/**
 * Zdanie o kolejce łączności okna — mówi, co się dzieje z poleceniem wpisanym w tej właśnie chwili pracy.
 */
export const ZDANIE_KOLEJKI =
  'Polecenie wpisane teraz nie ginie: staje w kolejce wychodzącej i pojedzie '
  + 'do rdzenia po powrocie łączności. Pole wypowiedzi zostaje czynne.';

/**
 * Zdanie o tym, czego transport nie wystawia — brak jest nazwany, żeby przy
 * samym słowie „Ponawianie" było wiadomo, dlaczego nie widać numeru próby.
 */
export const ZDANIE_BRAKU_PRZEBIEGU =
  'Numeru próby ani czasu do następnej ta plakietka nie pokazuje: transport '
  + 'trzyma je prywatnie i nie wystawia.';

/**
 * Podpowiedź plakietki łączności — pełne zdanie tam, gdzie w samym napisie mieści się tylko sam skrót stanu.
 */
export function podpowiedzLacznosci(stan: StanLacznosciOkna): string {
  const zdania = [`Łączność z rdzeniem: ${NAZWY[stan.stan]}.`];

  if (stan.oczekujace > 0) {
    zdania.push(`W kolejce wychodzącej czeka ${stan.oczekujace} ramek.`);
  }
  if (stan.stan !== 'polaczony') {
    zdania.push(ZDANIE_KOLEJKI);
  }
  if (stan.stan === 'ponawianie') {
    zdania.push(OPIS_PONAWIANIA);
    if (stan.numerProby === null) zdania.push(ZDANIE_BRAKU_PRZEBIEGU);
  }
  return zdania.join('\n');
}
