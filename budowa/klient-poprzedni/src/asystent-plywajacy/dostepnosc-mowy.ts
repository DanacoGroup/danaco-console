/**
 * Czym rozporządza dziś kanał głosowy Asystenta — jedno miejsce, w którym stoi
 * odpowiedź, i jedno źródło zdania, którym przycisk głosu odmawia.
 *
 * Stan drogi głosowej:
 *
 * - Droga polecenia do rdzenia jest: `assistant.voice.command` stoi
 *   w kontrakcie i ma uchwyt w rdzeniu — polecenie dojeżdża, zakłada zlecenie
 *   i wraca z jego stanem.
 * - Rozpoznania mowy nie ma. Żądanie z samym `audioRef`, bez `transcript`,
 *   zostawia zlecenie w stanie `queued`: rdzeń czeka na tekst i sam go nie
 *   wytworzy.
 * - Nie ma czym wytworzyć `audioRef` — żadna komenda kontraktu nie przyjmuje
 *   nagrania z przeglądarki, więc klient nie ma dokąd wysłać tego, co by
 *   nagrał. Nagrywanie do pamięci i wyrzucanie nagrania byłoby atrapą
 *   mikrofonu.
 * - Syntezy odpowiedzi nie ma: `speechRef` odpowiedzi wraca puste, a pole
 *   `speak` żądania nie ma kolumny w schemacie.
 *   `translate.speech.synthesize` syntezuje treść panelu tłumaczenia (żąda
 *   `panelId`) i odpowiedzi asystenta nie odczyta.
 *
 * Dymek prowadzi więc rozmowę tekstem — jednym wywołaniem
 * `assistant.voice.command` z polem `transcript` — a o braku głosu mówi
 * wprost, zamiast go udawać.
 */

import { pokazKomunikat } from '../aplikacja/komunikaty';

/** Ogniwo drogi głosowej wraz z odpowiedzią, czy jest zbudowane. */
export interface OgniwoGlosu {
  /** Nazwa ogniwa widoczna dla Operatora i w zgłoszeniu do rdzenia. */
  nazwa: string;
  /** Czy ogniwo istnieje dziś w produkcie. */
  jest: boolean;
  /** Czym zmierzone — komenda, pole kontraktu albo plik rdzenia. */
  dowod: string;
}

/** Ogniwa, z których składa się rozmowa głosowa, wraz ze stanem każdego. */
export const OGNIWA_GLOSU: readonly OgniwoGlosu[] = [
  {
    nazwa: 'Droga polecenia do rdzenia',
    jest: true,
    dowod: 'assistant.voice.command — komenda kontraktu z uchwytem w rdzeniu',
  },
  {
    nazwa: 'Przesył nagrania z przeglądarki',
    jest: false,
    dowod: 'żadna komenda kontraktu nie przyjmuje nagrania; pola audioRef nie ma czym wypełnić',
  },
  {
    nazwa: 'Rozpoznanie mowy (transkrypcja)',
    jest: false,
    dowod: 'rdzeń o sobie: „RDZEŃ NIE ROZPOZNAJE MOWY"',
  },
  {
    nazwa: 'Odczytanie odpowiedzi syntezą',
    jest: false,
    dowod: 'speechRef wraca puste — rdzeń nie ma silnika syntezy mowy; pole speak bez kolumny',
  },
];

/** Czy Operator może dziś powiedzieć polecenie zamiast je napisać. */
export function czyGlosDziala(): boolean {
  return OGNIWA_GLOSU.every((ogniwo) => ogniwo.jest);
}

/** Ogniwa, których brakuje — wykaz do zbudowania w rdzeniu i w kontrakcie. */
export function brakujaceOgniwa(): readonly OgniwoGlosu[] {
  return OGNIWA_GLOSU.filter((ogniwo) => !ogniwo.jest);
}

/**
 * Zdanie stojące w dymku od chwili otwarcia — o braku wiadomo przed
 * naciśnięciem mikrofonu, a nie dopiero po nim.
 */
export const ZAPOWIEDZ_BRAKU_GLOSU =
  'Ten dymek prowadzi dziś rozmowę TEKSTEM. Rdzeń nie rozpoznaje mowy i nie odczytuje ' +
  'odpowiedzi, a kontrakt nie ma komendy przesyłu nagrania — kanału głosowego nie ma, ' +
  'choć moduł jest do niego przeznaczony. Wpisane polecenie jedzie naprawdę, ' +
  'komendą assistant.voice.command.';

/** Powód odmowy przycisku głosu — rozwinięcie zapowiedzi o wykaz braków. */
export function powodOdmowyGlosu(): string {
  const braki = brakujaceOgniwa()
    .map((ogniwo) => `• ${ogniwo.nazwa} — ${ogniwo.dowod}`)
    .join('\n');
  return (
    'Rdzeń nie niesie jeszcze mowy, więc mikrofonu nie stawiamy — nagrywanie, ' +
    'którego nie ma dokąd wysłać, byłoby atrapą.\n\n' +
    `Brakuje ${brakujaceOgniwa().length} z ${OGNIWA_GLOSU.length} ogniw drogi głosowej:\n${braki}\n\n` +
    'Polecenie wydasz polem tekstowym obok — jedzie tą samą komendą, ' +
    'którą pojedzie kiedyś głos.'
  );
}

/**
 * Odpowiedź na naciśnięcie przycisku głosu.
 *
 * Przycisk zostaje klikalny i za każdym razem mówi to samo zdanie. Wygaszony
 * przycisk kazałby zgadywać, czy głos jest niedostępny chwilowo, czy
 * niezbudowany w ogóle.
 */
export function odmowMowy(): void {
  pokazKomunikat({
    tytul: 'Głosu jeszcze nie ma — i nie udajemy, że jest',
    tresc: powodOdmowyGlosu(),
    waga: 'ostrz',
  });
}
