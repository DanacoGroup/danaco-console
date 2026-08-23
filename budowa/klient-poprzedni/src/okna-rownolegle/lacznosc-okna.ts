import type { StanPolaczenia } from '../polaczenie/stan-polaczenia';

/**
 * Łączność widziana z wnętrza okna roboczego — model i napisy.
 *
 * Pasek górny aplikacji ma własny wskaźnik (`aplikacja/wskaznik-lacznosci.ts`),
 * ale okno rozwinięte na pełny ekran go przykrywa. Okno robocze mówi więc samo,
 * że łącze padło, zanim polecenie zostanie wysłane w próżnię.
 *
 * Ten plik niczego nie liczy i niczego nie pamięta: składa napis z odczytu
 * podanego przez `polaczenie/gniazdo.ts`. Własny licznik kolejki albo własne
 * zliczanie prób byłyby drugim źródłem prawdy o tej samej rzeczy.
 *
 * Stan łączności jest informacją, nie bramą. Pole wypowiedzi zostaje czynne
 * przy rozłączeniu, bo kolejka wychodząca jest nieograniczona
 * (`polaczenie/kolejka-wychodzaca.ts`) — plakietka mówi, że polecenie stanie
 * w kolejce, nie że wysłać się nie da.
 */

/** Odczyt transportu — tyle, ile transport wystawia dziś każdemu pytającemu. */
export interface OdczytLacznosci {
  stan: StanPolaczenia;
  /** Ramki czekające w kolejce wychodzącej. */
  oczekujace: number;
}

/**
 * Dojście do przebiegu ponowienia.
 *
 * Interfejs `Transport` (`polaczenie/gniazdo.ts`) tego nie wystawia: klasa
 * `Gniazdo` trzyma `numerProby` i uchwyt `zaplanowane` jako pola prywatne.
 * Port jest więc opcjonalny — plakietka bez niego pokazuje stan i kolejkę,
 * a brak numeru próby nazywa w podpowiedzi, zamiast liczyć próby u siebie.
 */
export interface PortPonawiania {
  /** Numer próby, którą transport podejmie jako następną (od 1). */
  numerProby(): number;
  /** Milisekundy do zaplanowanej próby; wartość ujemna albo 0 = próba trwa. */
  zaPonowieniem(): number;
  /** Podejmuje próbę natychmiast, bez czekania na zaplanowany odstęp. */
  ponowTeraz(): void;
}

/** Pełny stan łączności okna: odczyt transportu wraz z przebiegiem ponowienia. */
export interface StanLacznosciOkna extends OdczytLacznosci {
  /** Numer próby; `null` znaczy „transport go nie wystawia". */
  numerProby: number | null;
  /** Milisekundy do następnej próby; `null` znaczy „transport ich nie wystawia". */
  zaPonowieniemMs: number | null;
}

/** Nazwa stanu widoczna w interfejsie — ta sama, którą niesie pasek górny. */
const NAZWY: Readonly<Record<StanPolaczenia, string>> = {
  rozlaczony: 'Rozłączony',
  laczenie: 'Łączenie',
  polaczony: 'Połączony',
  ponawianie: 'Ponawianie',
};

/** Odmiana kropki z `komponenty/plakietka.css`; bez barw własnych. */
const KROPKI: Readonly<Record<StanPolaczenia, string>> = {
  rozlaczony: 'dn-kropka--blad',
  laczenie: 'dn-kropka--tetno',
  polaczony: 'dn-kropka--sukces',
  ponawianie: 'dn-kropka--ostrzezenie',
};

/** Odmiana plakietki — stan nigdy nie idzie samą barwą. */
const PLAKIETKI: Readonly<Record<StanPolaczenia, string>> = {
  rozlaczony: 'dn-plakietka--blad',
  laczenie: 'dn-plakietka--informacja',
  polaczony: 'dn-plakietka--sukces',
  ponawianie: 'dn-plakietka--ostrzezenie',
};

/** Stany, w których łączność jest w toku — nośnikiem jest wtedy `.dn-spinner`. */
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

/** Czy stan jest w toku — spinner zastępuje wtedy kropkę, tak jak w pasku górnym. */
export function czyLacznoscWToku(stan: StanPolaczenia): boolean {
  return W_TOKU.has(stan);
}

/** Klasa kropki dla stanu. */
export function wariantKropkiLacznosci(stan: StanPolaczenia): string {
  return KROPKI[stan];
}

/** Klasa plakietki dla stanu. */
export function wariantPlakietkiLacznosci(stan: StanPolaczenia): string {
  return PLAKIETKI[stan];
}

/** Nazwa stanu — jedno źródło napisu dla nagłówka i dla podpowiedzi. */
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

/** „za 4 s" — sekundy zaokrąglone w górę; poniżej sekundy próba już biegnie. */
function sekundyDoProby(ms: number): string {
  if (ms <= 0) return 'teraz';
  return `${Math.ceil(ms / 1000)} s`;
}

/**
 * Opis polityki ponawiania dla podpowiedzi plakietki.
 *
 * Wartości pochodzą z `polaczenie/ponawianie.ts` i to tam są ustalane; ten
 * napis tylko je nazywa. Uzasadnienie doboru:
 *
 * - baza 500 ms — rdzeń wstający lokalnie wraca zwykle w pierwszej albo drugiej
 *   próbie, więc przerwa bywa niezauważalna;
 * - mnożnik dwa z pułapem 15 s — po rdzeniu, który nie wstał, klient nie dobija
 *   się co pół sekundy bez końca, a 15 s to jeszcze czekanie, po którym produkt
 *   nie wygląda na martwy;
 * - rozproszenie do 25 % — okna równoległe i karty sesji ponawiają niezależnie;
 *   bez rozproszenia trafiałyby w rdzeń równocześnie;
 * - brak górnego limitu prób — poddanie się po ustalonej liczbie prób zostawiłoby
 *   użytkownika z produktem wymagającym przeładowania strony. Przyspieszenie
 *   daje czynność „Ponów teraz"; rezygnacji nie ma.
 */
export const OPIS_PONAWIANIA =
  'Ponawianie jest wykładnicze: pierwsza próba po 500 ms, każda następna dwa '
  + 'razy dalej, nie dalej niż co 15 s, z rozproszeniem losowym do 25 %. '
  + 'Liczba prób nie jest ograniczona — klient nie poddaje się sam.';

/** Zdanie o kolejce — mówi, co się dzieje z poleceniem wpisanym teraz. */
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

/** Podpowiedź plakietki — pełne zdanie tam, gdzie w napisie mieści się skrót. */
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
