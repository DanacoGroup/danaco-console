/**
 * Formy rzeczownika po liczebniku. Polszczyzna wymaga trzech i nie da się ich
 * zapisać jednym wzorcem w katalogu treści, więc wybór formy stoi w kodzie,
 * a same słowa — w katalogu panelu, tam gdzie tłumacz je zobaczy.
 *
 * Reguła: 1 bierze formę pojedynczą; 2–4 bierze mnogą, ale nie w drugiej
 * dziesiątce (12–14 idzie z resztą); pozostałe biorą dopełniacz mnogi.
 */

export interface FormyLiczebnika {
  /** Forma po „1” — np. „wersja”, „znak”, „słowo”. */
  jedna: string;
  /** Forma po „2”, „3”, „4” — np. „wersje”, „znaki”, „słowa”. */
  kilka: string;
  /** Forma po „5” i dalej oraz po 12–14 — np. „wersji”, „znaków”, „słów”. */
  wiele: string;
}

/** Wybiera formę rzeczownika właściwą dla podanej liczby. */
export function formaPo(liczba: number, formy: FormyLiczebnika): string {
  const bezwzgledna = Math.abs(Math.trunc(liczba));
  if (bezwzgledna === 1) return formy.jedna;
  const ostatniaCyfra = bezwzgledna % 10;
  const drugaDziesiatka = bezwzgledna % 100 >= 12 && bezwzgledna % 100 <= 14;
  return ostatniaCyfra >= 2 && ostatniaCyfra <= 4 && !drugaDziesiatka ? formy.kilka : formy.wiele;
}

/** Liczba wraz z formą rzeczownika właściwą dla niej — np. „34 słowa”. */
export function zLiczba(liczba: number, formy: FormyLiczebnika): string {
  return `${liczba} ${formaPo(liczba, formy)}`;
}
