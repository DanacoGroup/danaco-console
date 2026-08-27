import { KOMENDY, type Command } from '../../../../shared/contract';

/** Przedrostek obszaru komend modułu roundtable, użyty do wydzielenia z wykazu kontraktu komend należących do tego obszaru. */
const OBSZAR = 'roundtable.';

/** Komendy obszaru roundtable obecne dziś w wykazie kontraktu, wyliczone dynamicznie zamiast wpisane na sztywno liczbą. */
function komendyObszaru(): string[] {
  return KOMENDY.map((komenda) => String(komenda)).filter((komenda) => komenda.startsWith(OBSZAR));
}

/**
 * Zdanie o zasobności obszaru roundtable — ta sama treść pod każdą pozycją bez obsługi,
 * podająca liczbę komend, nie ich nazwy.
 */
function zasobnoscObszaru(): string {
  const ile = komendyObszaru().length;
  if (ile === 0) {
    return 'Wykaz kontraktu nie niesie dziś ani jednej komendy obszaru roundtable.';
  }
  return `Obszar roundtable niesie dziś w kontrakcie ${ile} komend.`;
}

/**
 * Powód dla czynności, której komenda jest w kontrakcie, a okno jej nie wywołuje jeszcze —
 * nazywa komendę i mówi, że obsługi nie zbudowano.
 */
export function powodBrakuObslugi(komenda: Command, zdanie: string): string {
  const nazwa = String(komenda);
  if (!nazwa.startsWith(OBSZAR)) {
    return (
      `UWAGA: pod tą pozycją stoi komenda ${nazwa} spoza obszaru roundtable — ` +
      `to usterka okna, nie stan produktu. ${zasobnoscObszaru()}`
    );
  }
  return (
    `Komenda ${nazwa} jest w kontrakcie, ale obsługi tej czynności jeszcze nie zbudowano. ` +
    `${zdanie} ${zasobnoscObszaru()}`
  );
}

/**
 * Powód dla czynności, której nie da się wskazać jedną komendą kontraktu, wraz z liczbą
 * komend obszaru dla kontekstu.
 */
export function powodBrakuCzynnosci(zdanie: string): string {
  return `${zdanie} ${zasobnoscObszaru()}`;
}
