/**
 * Zakres zaznaczenia w kanwie dokumentu liczony w znakach — kształt, którego
 * żądają komendy postaci rodziny `studio`. Przycisk paska narzędziowego
 * zabiera ognisko i zwija zaznaczenie, więc ostatni zakres wskazany w kanwie
 * zostaje zapamiętany przy niej.
 */

/** Fragment tekstu wskazany w kanwie; oba końce liczone w znakach od początku treści. */
export interface ZakresTresci {
  poczatek: number;
  koniec: number;
}

const ZAPAMIETANE = new WeakMap<HTMLElement, ZakresTresci>();

/** Zapamiętuje zakres wskazany w kanwie; nasłuch schodzi sterownikiem przerwania wołającego. */
export function sledzZaznaczenie(kanwa: HTMLElement, przy: AddEventListenerOptions): void {
  document.addEventListener('selectionchange', () => {
    const zakres = odczytajZakres(kanwa);
    if (zakres !== null) ZAPAMIETANE.set(kanwa, zakres);
  }, przy);
}

/** Zakres wskazany w kanwie albo ostatni zapamiętany; pustka znaczy kanwę bez zaznaczenia. */
export function zakresZaznaczenia(kanwa: HTMLElement): ZakresTresci | null {
  return odczytajZakres(kanwa) ?? ZAPAMIETANE.get(kanwa) ?? null;
}

/** Zakres niepusty; komendy postaci biorą brak zakresu jako cały dokument, więc zwinięty punkt odpada. */
export function fragmentZaznaczony(kanwa: HTMLElement): ZakresTresci | null {
  const zakres = zakresZaznaczenia(kanwa);
  return zakres === null || zakres.poczatek === zakres.koniec ? null : zakres;
}

/** Miejsce kursora w znakach; bez zaznaczenia wstawianie idzie na koniec treści. */
export function miejsceKursora(kanwa: HTMLElement): number {
  return zakresZaznaczenia(kanwa)?.koniec ?? (kanwa.textContent ?? '').length;
}

function odczytajZakres(kanwa: HTMLElement): ZakresTresci | null {
  const zaznaczenie = document.getSelection();
  if (zaznaczenie === null || zaznaczenie.rangeCount === 0) return null;
  const zakres = zaznaczenie.getRangeAt(0);
  if (!kanwa.contains(zakres.commonAncestorContainer)) return null;
  const poprzedzajacy = zakres.cloneRange();
  poprzedzajacy.selectNodeContents(kanwa);
  poprzedzajacy.setEnd(zakres.startContainer, zakres.startOffset);
  const poczatek = poprzedzajacy.toString().length;
  return { poczatek, koniec: poczatek + zakres.toString().length };
}
