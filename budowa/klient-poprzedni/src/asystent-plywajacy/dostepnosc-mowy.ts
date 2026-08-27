/**
 * Czym rozporządza dziś kanał głosowy Asystenta. Dymek prowadzi rozmowę tekstem,
 * jednym wywołaniem komendy `assistant.voice.command` z polem transkrypcji,
 * a o braku głosu mówi wprost, zamiast go udawać.
 */

import { pokazKomunikat } from '../aplikacja/komunikaty';

/**
 * Ogniwo drogi głosowej wraz z odpowiedzią, czy jest zbudowane: nazwa widoczna
 * dla Operatora i w zgłoszeniu do rdzenia, wskaźnik istnienia oraz dowód, czyli
 * komenda, pole kontraktu albo plik rdzenia.
 */
export interface OgniwoGlosu {
  /** Nazwa ogniwa widoczna dla Operatora i w zgłoszeniu do rdzenia. */
  nazwa: string;
  /** Czy ogniwo istnieje dziś w produkcie. */
  jest: boolean;
  /** Czym zmierzone — komenda, pole kontraktu albo plik rdzenia. */
  dowod: string;
}

/**
 * Ogniwa, z których składa się rozmowa głosowa, wraz ze stanem każdego z nich.
 * Wykaz jest jednym miejscem, w którym stoi odpowiedź o stanie kanału głosowego,
 * i jednym źródłem zdania, którym przycisk głosu odmawia.
 */
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

/**
 * Czy Operator może dziś powiedzieć polecenie zamiast je napisać. Odpowiedź
 * twierdząca wymaga wszystkich ogniw drogi głosowej naraz, ponieważ brak
 * jednego przerywa drogę w całości.
 */
export function czyGlosDziala(): boolean {
  return OGNIWA_GLOSU.every((ogniwo) => ogniwo.jest);
}

/**
 * Ogniwa, których brakuje, czyli wykaz do zbudowania w rdzeniu i w kontrakcie.
 * Wykaz powstaje z odsiania ogniw już zbudowanych, więc nie rozjeżdża się
 * z zapisem stanu drogi głosowej.
 */
export function brakujaceOgniwa(): readonly OgniwoGlosu[] {
  return OGNIWA_GLOSU.filter((ogniwo) => !ogniwo.jest);
}

/**
 * Zdanie stojące w dymku od chwili otwarcia, żeby o braku kanału głosowego
 * wiadomo było przed naciśnięciem mikrofonu, a nie dopiero po nim.
 */
export const ZAPOWIEDZ_BRAKU_GLOSU =
  'Ten dymek prowadzi dziś rozmowę TEKSTEM. Rdzeń nie rozpoznaje mowy i nie odczytuje ' +
  'odpowiedzi, a kontrakt nie ma komendy przesyłu nagrania — kanału głosowego nie ma, ' +
  'choć moduł jest do niego przeznaczony. Wpisane polecenie jedzie naprawdę, ' +
  'komendą assistant.voice.command.';

/**
 * Powód odmowy przycisku głosu, będący rozwinięciem zapowiedzi o wykaz ogniw,
 * których brakuje, wraz z dowodem zmierzenia każdego z nich. Wykaz nadaje się
 * wprost na zgłoszenie braku do rdzenia.
 */
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
 * Odpowiedź na naciśnięcie przycisku głosu. Przycisk zostaje klikalny i za
 * każdym razem mówi to samo zdanie, ponieważ przycisk wygaszony kazałby zgadywać,
 * czy głos jest niedostępny chwilowo, czy niezbudowany w ogóle.
 */
export function odmowMowy(): void {
  pokazKomunikat({
    tytul: 'Głosu jeszcze nie ma — i nie udajemy, że jest',
    tresc: powodOdmowyGlosu(),
    waga: 'ostrz',
  });
}
