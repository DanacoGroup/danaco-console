/* Polszczyzna odmienia rzeczownik po liczbie w trzech postaciach; miara podana
   samą liczbą zostawiłaby w oknie zapis, którego Operator nie czyta zdaniem. */

/** Prawda dla liczby biorącej mianownik liczby mnogiej: 2-4 poza nastoma. */
function mnoga(ile: number): boolean {
  const reszta = ile % 10;
  const setka = ile % 100;
  return reszta >= 2 && reszta <= 4 && (setka < 12 || setka > 14);
}

/** Liczba sesji wraz z odmianą rzeczownika. */
export function miaraSesji(ile: number): string {
  if (ile === 1) return '1 sesja';
  return `${ile} ${mnoga(ile) ? 'sesje' : 'sesji'}`;
}

/** Liczba modułów wraz z odmianą rzeczownika. */
export function miaraModulow(ile: number): string {
  if (ile === 1) return '1 moduł';
  return `${ile} ${mnoga(ile) ? 'moduły' : 'modułów'}`;
}
