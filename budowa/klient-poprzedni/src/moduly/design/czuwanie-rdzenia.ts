import type { Wynik } from '../../protokol/kanal';

// Czuwanie odróżnia rdzeń pracujący długo od kanału milczącego, pytając kanał o życie.

/** Jedno tanie pytanie kierowane do rdzenia o życie kanału; obietnica spełnia się dopiero, gdy rdzeń na nie odpowie. */
export type ProbaZycia = () => Promise<unknown>;

/** Określa, co okno robi z prawdą o czynności, dla której rdzeń nie zdążył jeszcze rozstrzygnąć skutku wywołania. */
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
  /** Prowadzi wywołanie rdzenia pod czuwaniem; zwraca odpowiedź albo `null`, gdy kanał zamilkł. */
  prowadz<T>(
    czynnosc: string,
    wywolanie: Promise<Wynik<T>>,
    nasluch: NasluchCzuwania,
  ): Promise<Wynik<T> | null>;
}

/** Po upływie tylu milisekund bez odpowiedzi na czynność czuwanie pyta kanał, czy jeszcze w ogóle odpowiada. */
const CZAS_CISZY_MS = 4000;
/** Tyle czasu czuwanie czeka na odpowiedź próby życia kanału; rdzeń lokalny odpowiada zwykle w milisekundach. */
const CZAS_PROBY_MS = 2500;

/** Znacznik wygranej zegara w wyścigu obietnic czuwania; nie jest wynikiem żadnej czynności ani próby życia. */
const ZEGAR = Symbol('zegar czuwania');

function poCzasie<T>(ms: number, wartosc: T): Promise<T> {
  return new Promise((rozstrzygnij) => {
    setTimeout(() => rozstrzygnij(wartosc), ms);
  });
}

/** Zdanie przedstawiane Operatorowi o czynności, która wciąż trwa, mimo że kanał w międzyczasie odpowiedział. */
function zdanieWToku(czynnosc: string, sekund: number): string {
  return (
    `Rdzeń odpowiada na inne pytania, więc łączność jest — ${czynnosc} trwa już ${sekund} s ` +
    'bez rozstrzygnięcia. Czekam dalej i nie orzekam o skutku.'
  );
}

/** Zdanie przedstawiane Operatorowi o kanale, który zamilkł; nie orzeka nic o skutku prowadzonej czynności. */
function zdanieCiszy(czynnosc: string): string {
  return (
    `Połączenie z rdzeniem zamilkło w trakcie czynności: ${czynnosc}. Odpowiedzi na to żądanie ` +
    'już nie będzie — była przypisana do gniazda, które padło, i ponowne połączenie jej nie ' +
    'wskrzesi. NIE WIADOMO, czy czynność się odbyła: żądanie mogło dojść do rdzenia i zostać ' +
    'wykonane, więc nie ogłaszam niepowodzenia. Gdy łączność wróci, odczytaj stan z rdzenia — ' +
    'tylko odczyt powie, co naprawdę zostało zapisane.'
  );
}

/** Zdanie przedstawiane Operatorowi o powrocie łączności, bez zmiany wcześniejszego orzeczenia o skutku czynności. */
function zdaniePowrotu(czynnosc: string): string {
  return (
    `Łączność z rdzeniem wróciła (rdzeń odpowiedział na próbę), ale odpowiedzi na czynność ` +
    `„${czynnosc}" już nie ma i nie będzie — jej skutek pozostaje NIEZNANY. Odczytaj zasoby, ` +
    'żeby zobaczyć stan, który naprawdę leży w rdzeniu.'
  );
}

/** Zdanie przedstawiane Operatorowi o odpowiedzi, która przyszła spóźniona; czuwanie poprawia własne orzeczenie. */
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
      // Obietnica wywołania jest jedna na cały czas czuwania; drugie `then` nie ponawia żądania.
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
        // Ta sama próba czeka w kolejce wychodzącej i jest czujnikiem powrotu za darmo.
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
