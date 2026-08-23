import { KOMENDY, type Command } from '../../../../shared/contract';

/**
 * Powody „bez obsługi” i „bez komendy” składane z wykazu komend kontraktu,
 * a nie z napisów wpisanych na sztywno.
 *
 * Napis orzekający na stałe, czego kontrakt nie ma, przestaje być prawdą
 * w dniu powstania komendy i nikt go wtedy nie zdejmuje. Ten moduł przeszedł
 * dokładnie przez taki dzień: obszar `roundtable.*` niósł cztery komendy,
 * a po scaleniu kontraktu niesie ich czterdzieści sześć. Zdania mówiące
 * „nie wykonuje tego żadna komenda obszaru” stały się wtedy nieprawdą co do
 * czterdziestu dwóch czynności naraz.
 *
 * Stąd podział na dwa powody, nie jeden:
 *
 *   • `powodBrakuObslugi` — komenda JEST w kontrakcie, a okno jej nie wywołuje.
 *     Nazwa komendy przychodzi wartością `Command`, nie napisem, więc zmiana
 *     nazwy w kontrakcie przerywa kompilację zamiast zostawiać w dymku napis
 *     o komendzie, której już nie ma.
 *   • `powodBrakuCzynnosci` — czynności nie da się wskazać jedną komendą, bo
 *     kontrakt nie ma jej w ogóle albo brakuje pola, nie komendy.
 *
 * Czego ten plik nie wie: czy rdzeń ma obsługiwacz danej komendy. Kontrakt nie
 * daje takiego odczytu — nie ma komendy wyliczającej komendy obsługiwane
 * (`module.list` mówi o oknach modułu, nie o obsługiwaczach). Rdzeń odpowiadający
 * `roundtable.unknown` na komendę obecną w wykazie jest z tego miejsca
 * nierozpoznawalny; rozpoznaje go dopiero odmowa po naciśnięciu. Dlatego zdania
 * mówią o obsłudze niezbudowanej, a nie orzekają, gdzie dokładnie jej brakuje.
 */

/** Przedrostek obszaru komend modułu. */
const OBSZAR = 'roundtable.';

/** Komendy obszaru obecne w wykazie kontraktu. */
function komendyObszaru(): string[] {
  return KOMENDY.map((komenda) => String(komenda)).filter((komenda) => komenda.startsWith(OBSZAR));
}

/**
 * Zdanie o zasobności obszaru — ta sama treść pod każdą pozycją bez obsługi.
 *
 * Wykaz podaje liczbę, a nie nazwy. Przy czterdziestu sześciu komendach
 * wypisanie ich wszystkich w dymku przy przycisku zasłoniłoby powód, dla którego
 * dymek się otworzył; liczba mówi to samo, co trzeba: obszar nie jest ubogi,
 * uboga jest obsługa.
 */
function zasobnoscObszaru(): string {
  const ile = komendyObszaru().length;
  if (ile === 0) {
    return 'Wykaz kontraktu nie niesie dziś ani jednej komendy obszaru roundtable.';
  }
  return `Obszar roundtable niesie dziś w kontrakcie ${ile} komend.`;
}

/**
 * Powód dla czynności, której komenda jest w kontrakcie, a okno jej nie wywołuje.
 *
 * Zdanie nie mówi „kontrakt tego nie przewiduje”, bo przewiduje. Mówi, że
 * czynności nie zbudowano — i nazywa komendę, po której Operator odnajdzie ją
 * w kontrakcie. Komenda spoza obszaru jest usterką wywołania, nie stanem
 * produktu, więc zdanie mówi to wprost zamiast ją przemilczeć.
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
 * Powód dla czynności, której nie da się wskazać jedną komendą.
 *
 * Dotyczy dwóch przypadków: kontrakt nie ma komendy w ogóle albo ma komendę,
 * lecz brakuje pola, którym czynność jechałaby w żądaniu. Zdanie autora nazywa,
 * który to przypadek, i jedzie wraz z liczbą komend obszaru, żeby nie czytało
 * się jako zarzut wobec kontraktu ubogiego.
 */
export function powodBrakuCzynnosci(zdanie: string): string {
  return `${zdanie} ${zasobnoscObszaru()}`;
}
