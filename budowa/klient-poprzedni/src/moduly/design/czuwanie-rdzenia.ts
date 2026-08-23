import type { Wynik } from '../../protokol/kanal';

/**
 * Czuwanie nad czynnością, której rdzeń nie rozstrzygnął: mówi prawdę
 * o wywołaniu bez odpowiedzi i odróżnia rdzeń pracujący od kanału milczącego.
 *
 * Gdy gniazdo padnie w trakcie oczekiwania, okno samo się nie odnajdzie:
 * obietnica `protokol/wywolanie.ts` nigdy nie jest odrzucana, a odbiorca
 * odpowiedzi żyje w rejestrze korelacji przypisanym do gniazda, które padło —
 * odpowiedź nie przyjdzie już nigdy, także po ponownym połączeniu.
 *
 * Zwykły limit czasu nie odróżnia rdzenia, który pracuje długo, od kanału,
 * który zamilkł, a komendy modułu bywają wolne i długie oczekiwanie jest tu
 * stanem poprawnym. Czuwanie pyta więc kanał o życie: wysyła jedno tanie
 * żądanie kontraktu tą samą drogą.
 *   · próba odpowiedziała → kanał żyje, czynność nadal trwa
 *   · próba zamilkła      → kanał zamilkł, skutek pozostaje nieznany
 *
 * Ta sama próba jest zarazem czujnikiem powrotu i nie kosztuje dodatkowej
 * ramki: żądanie wysłane przy rozłączeniu czeka w kolejce wychodzącej
 * (`polaczenie/kolejka-wychodzaca.ts`), a rdzeń odpowiada na nie w chwili
 * ponownego połączenia. Jedna obietnica próby mówi więc najpierw „kanał
 * milczy", a potem „łączność wróciła"; pętli odpytującej tu nie ma.
 *
 * Czuwanie nie ogłasza niepowodzenia czynności i nie zgaduje jej skutku —
 * rdzeń mógł żądanie odebrać i wykonać, zanim gniazdo padło. Jedyną drogą do
 * prawdy jest odczyt po powrocie łączności i tak brzmi zdanie dla Operatora.
 */

/** Jedno tanie pytanie do rdzenia; obietnica spełnia się, gdy rdzeń odpowie. */
export type ProbaZycia = () => Promise<unknown>;

/** Co okno robi z prawdą o czynności bez rozstrzygnięcia. */
export interface NasluchCzuwania {
  /** Kanał odpowiada, czynność nadal trwa — rdzeń pracuje, nie milczy. */
  wToku?(zdanie: string): void;
  /** Kanał zamilkł: czynność bez rozstrzygnięcia, skutek nieznany. */
  cisza(zdanie: string): void;
  /** Łączność wróciła, a odpowiedzi na tę czynność już nie będzie. */
  powrot?(zdanie: string): void;
  /** Odpowiedź przyszła po uznaniu kanału za milczący — poprawka zdania. */
  spozniona?(zdanie: string): void;
}

export interface CzuwanieRdzenia {
  /**
   * Prowadzi wywołanie rdzenia pod czuwaniem.
   *
   * Zwraca odpowiedź rdzenia albo `null`, gdy kanał zamilkł. `null` nie jest
   * odmową — znaczy „bez rozstrzygnięcia".
   *
   * @param czynnosc nazwa czynności w mianowniku, np. „generowanie zasobu"
   */
  prowadz<T>(
    czynnosc: string,
    wywolanie: Promise<Wynik<T>>,
    nasluch: NasluchCzuwania,
  ): Promise<Wynik<T> | null>;
}

/** Po tyle bez odpowiedzi pytamy kanał, czy w ogóle żyje. */
const CZAS_CISZY_MS = 4000;
/** Tyle czekamy na odpowiedź próby życia; rdzeń lokalny odpowiada w milisekundach. */
const CZAS_PROBY_MS = 2500;

/** Znacznik wygranej zegara w wyścigu — nie do pomylenia z żadnym wynikiem. */
const ZEGAR = Symbol('zegar czuwania');

function poCzasie<T>(ms: number, wartosc: T): Promise<T> {
  return new Promise((rozstrzygnij) => {
    setTimeout(() => rozstrzygnij(wartosc), ms);
  });
}

/** Zdanie o czynności, która trwa, choć kanał odpowiada. */
function zdanieWToku(czynnosc: string, sekund: number): string {
  return (
    `Rdzeń odpowiada na inne pytania, więc łączność jest — ${czynnosc} trwa już ${sekund} s ` +
    'bez rozstrzygnięcia. Czekam dalej i nie orzekam o skutku.'
  );
}

/** Zdanie o kanale, który zamilkł. Nie orzeka o skutku czynności. */
function zdanieCiszy(czynnosc: string): string {
  return (
    `Połączenie z rdzeniem zamilkło w trakcie czynności: ${czynnosc}. Odpowiedzi na to żądanie ` +
    'już nie będzie — była przypisana do gniazda, które padło, i ponowne połączenie jej nie ' +
    'wskrzesi. NIE WIADOMO, czy czynność się odbyła: żądanie mogło dojść do rdzenia i zostać ' +
    'wykonane, więc nie ogłaszam niepowodzenia. Gdy łączność wróci, odczytaj stan z rdzenia — ' +
    'tylko odczyt powie, co naprawdę zostało zapisane.'
  );
}

/** Zdanie o powrocie łączności — bez zmiany orzeczenia o skutku. */
function zdaniePowrotu(czynnosc: string): string {
  return (
    `Łączność z rdzeniem wróciła (rdzeń odpowiedział na próbę), ale odpowiedzi na czynność ` +
    `„${czynnosc}" już nie ma i nie będzie — jej skutek pozostaje NIEZNANY. Odczytaj zasoby, ` +
    'żeby zobaczyć stan, który naprawdę leży w rdzeniu.'
  );
}

/** Zdanie o odpowiedzi spóźnionej — czuwanie poprawia własne orzeczenie. */
function zdanieSpoznione(czynnosc: string, wynik: Wynik<unknown>): string {
  const co = wynik.udany
    ? 'rdzeń przyjął żądanie'
    : `rdzeń odmówił: ${wynik.blad?.message ?? 'bez powodu w odpowiedzi'}`;
  return (
    `Odpowiedź na czynność „${czynnosc}" przyszła PO tym, jak uznałem kanał za milczący — ${co}. ` +
    'Skutek nie jest już nieznany; odczytaj zasoby, żeby zobaczyć zapisany stan.'
  );
}

export function utworzCzuwanieRdzenia(probaZycia: ProbaZycia): CzuwanieRdzenia {
  return {
    async prowadz(czynnosc, wywolanie, nasluch) {
      const poczatek = Date.now();
      // Obietnica wywołania jest jedna na cały czas czuwania: drugie `then`
      // na tej samej obietnicy nie wysyła drugiego żądania do rdzenia.
      const odpowiedz = wywolanie.then((wynik) => ({ rodzaj: 'odpowiedz' as const, wynik }));

      for (;;) {
        const pierwszy = await Promise.race([odpowiedz, poCzasie(CZAS_CISZY_MS, ZEGAR)]);
        if (pierwszy !== ZEGAR) return pierwszy.wynik;

        // Jedno żądanie kontraktu tą samą drogą — pytanie „czy kanał żyje".
        const proba = probaZycia().then(() => 'zyje' as const);
        const drugi = await Promise.race([odpowiedz, proba, poCzasie(CZAS_PROBY_MS, ZEGAR)]);
        if (drugi !== ZEGAR && drugi !== 'zyje') return drugi.wynik;
        if (drugi === 'zyje') {
          nasluch.wToku?.(zdanieWToku(czynnosc, Math.round((Date.now() - poczatek) / 1000)));
          continue;
        }

        nasluch.cisza(zdanieCiszy(czynnosc));
        // Ta sama próba czeka teraz w kolejce wychodzącej i odpowie w chwili
        // ponownego połączenia — jest więc czujnikiem powrotu za darmo.
        void pilnujPowrotu(czynnosc, odpowiedz, proba, nasluch);
        return null;
      }
    },
  };
}

/**
 * Po ogłoszeniu ciszy zostaje jedno pytanie: co odezwie się pierwsze —
 * odpowiedź na czynność (wtedy poprawiamy orzeczenie) czy próba życia
 * (wtedy wiemy, że łączność wróciła, a odpowiedzi nie będzie).
 */
async function pilnujPowrotu<T>(
  czynnosc: string,
  odpowiedz: Promise<{ rodzaj: 'odpowiedz'; wynik: Wynik<T> }>,
  proba: Promise<'zyje'>,
  nasluch: NasluchCzuwania,
): Promise<void> {
  const co = await Promise.race([odpowiedz, proba]);
  if (co === 'zyje') {
    nasluch.powrot?.(zdaniePowrotu(czynnosc));
    return;
  }
  nasluch.spozniona?.(zdanieSpoznione(czynnosc, co.wynik));
}
