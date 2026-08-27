import type { Channel } from '../../../../shared/contract';
import type { StanDebaty } from './stan-debaty';

/** Stan jednego kanału niosącego uczestników tej debaty, zestawiony z rejestrem kanałów rdzenia jako warstwa czwarta Model Panels. */
export interface StanKanaluDebaty {
  idKanalu: string;
  /** Wiersz rejestru; `null`, gdy rejestr tego kanału nie zna. */
  wpis: Channel | null;
  /** Ilu uczestników debaty jedzie tym kanałem. */
  uczestnikow: number;
}

/**
 * Kanały niosące skład debaty, w kolejności pierwszego wystąpienia w składzie.
 *
 * Wykaz idzie po składzie, nie po rejestrze: rejestr niesie wszystkie kanały
 * instalacji, a diagnostyka dotyczy wyłącznie tych, które w debacie mówią.
 */
export function kanalyDebaty(stan: StanDebaty): StanKanaluDebaty[] {
  const wykaz: StanKanaluDebaty[] = [];
  const rejestr = stan.kanaly();
  for (const uczestnik of stan.uczestnicy()) {
    const znany = wykaz.find((pozycja) => pozycja.idKanalu === uczestnik.channelId);
    if (znany !== undefined) {
      znany.uczestnikow += 1;
      continue;
    }
    wykaz.push({
      idKanalu: uczestnik.channelId,
      wpis: rejestr.find((kanal) => kanal.id === uczestnik.channelId) ?? null,
      uczestnikow: 1,
    });
  }
  return wykaz;
}

/** Zdania diagnostyki jednego kanału — po jednym na wiersz wykazu, nazywające wprost brak wpisu w rejestrze albo stan kanału czynnego. */
export function zdaniaKanalu(kanal: StanKanaluDebaty, kanalyOdczytane: boolean): string[] {
  if (kanal.wpis === null) {
    return [
      kanalyOdczytane
        ? `Kanał ${kanal.idKanalu} nie ma wiersza w rejestrze rdzenia — uczestnik zostaje w debacie, a jego panel podpisuje się identyfikatorem kanału zamiast nazwą.`
        : `Kanał ${kanal.idKanalu}: rejestr kanałów nie został jeszcze odczytany, więc o tym kanale nie wiadomo nic ponad identyfikator.`,
      `Uczestników na tym kanale: ${kanal.uczestnikow}.`,
    ];
  }
  const wpis = kanal.wpis;
  const model = wpis.model === undefined || wpis.model === '' ? 'model nienazwany w rejestrze' : `model ${wpis.model}`;
  const konto =
    wpis.accountId === undefined || wpis.accountId === ''
      ? 'bez konta poświadczeń w rejestrze'
      : `konto poświadczeń ${wpis.accountId}`;
  return [
    `${wpis.name} — rodzaj ${wpis.kind}, ${model}, ${konto}.`,
    wpis.enabled
      ? 'Kanał czynny w rejestrze rdzenia.'
      : 'Kanał NIECZYNNY w rejestrze rdzenia — uczestnik na nim nie odpowie, choć zostaje w składzie.',
    `Uczestników na tym kanale: ${kanal.uczestnikow}${kanal.uczestnikow > 1 ? ' — to odrębne tożsamości jednego kanału.' : '.'}`,
  ];
}

/**
 * Zdanie zamykające diagnostykę: czego kontrakt nie mierzy.
 *
 * Bez niego wykaz kanałów wyglądałby na pełny obraz zdrowia kanału, którym nie
 * jest — rejestr mówi o konfiguracji kanału, nie o jego zachowaniu w debacie.
 */
export const ZDANIE_GRANICY_DIAGNOSTYKI =
  'Diagnostyka czyta wyłącznie rejestr kanałów rdzenia. Czasu odpowiedzi kanału, liczby odmów ani stanu połączenia nie oddaje żadna komenda obszaru roundtable, więc okno ich nie pokazuje.';
