import type { PanelSection } from '../../../shared/contract';

/**
 * Reguła układu sekcji panelu: kto stoi przed kim, co jest zwinięte, co zdjęte
 * z widoku. Rozmowę z rdzeniem prowadzi `zrodlo-sekcji-paneli`, kształt widoku
 * buduje `czesci-sekcji`. Każda funkcja tutaj jest czysta — bierze układ, oddaje
 * nowy, niczego nie modyfikuje w miejscu.
 *
 * Kolejność liczy się od 1 (`PanelSection.order` kontraktu) i jest przeliczana
 * po każdej zmianie: układ z dziurami w numeracji rdzeń przyjmie, ale kolejne
 * przestawienie policzy błędnie. Wpis rdzenia o identyfikatorze, jakiego ten
 * panel nie buduje, odpada w `scalUklad` — nie da się go narysować.
 */

/** Sekcja panelu w postaci, którą podaje moduł: identyfikator, tytuł i treść. */
export interface OpisSekcji {
  /** Identyfikator sekcji w obrębie panelu — trafia w `PanelSection.id`. */
  id: string;
  /** Nagłówek sekcji widoczny dla Operatora. */
  tytul: string;
  /** Gotowy element treści; sekcja go tylko obudowuje, nigdy nie buduje sama. */
  tresc: HTMLElement;
}

/** Przeliczenie `order` na 1…n według kolejności w tablicy. */
function ponumeruj(uklad: readonly PanelSection[]): PanelSection[] {
  return uklad.map((sekcja, indeks) => ({ ...sekcja, order: indeks + 1 }));
}

/** Układ domyślny — kolejność zgłoszenia sekcji przez moduł, nic zwiniętego. */
export function ukladDomyslny(sekcje: readonly OpisSekcji[]): PanelSection[] {
  return sekcje.map((opis, indeks) => ({ id: opis.id, order: indeks + 1, collapsed: false }));
}

/**
 * Zestawienie układu z rdzenia z sekcjami, które panel faktycznie buduje.
 *
 * Sekcje znane rdzeniowi idą w jego kolejności, sekcje panelu nieobecne
 * w odpowiedzi dochodzą na koniec (panel urósł od czasu zapisu), a wpisy rdzenia
 * bez odpowiednika w panelu odpadają (panel się skurczył). Numeracja wychodzi
 * ciągła niezależnie od tego, który z trzech przypadków zaszedł.
 */
export function scalUklad(
  oddane: readonly PanelSection[],
  sekcje: readonly OpisSekcji[],
): PanelSection[] {
  const znane = new Set(sekcje.map((opis) => opis.id));
  const wedlugRdzenia = [...oddane]
    .filter((sekcja) => znane.has(sekcja.id))
    .sort((a, b) => a.order - b.order);
  const juzUlozone = new Set(wedlugRdzenia.map((sekcja) => sekcja.id));
  const dopisane = sekcje
    .filter((opis) => !juzUlozone.has(opis.id))
    .map((opis) => ({ id: opis.id, order: 0, collapsed: false }));
  return ponumeruj([...wedlugRdzenia, ...dopisane]);
}

/**
 * Przesunięcie sekcji o jedno miejsce w podanym kierunku.
 *
 * Przesunięcie liczy się wyłącznie względem sekcji widocznych: sekcja zdjęta
 * z widoku nie zajmuje miejsca, więc ruch ponad nią nic by na ekranie nie
 * zmienił. Ruch poza wykaz oddaje układ bez zmiany, a wywołujący porówna go
 * z poprzednim i nie wyśle zapisu bez skutku.
 */
export function przestawSekcje(
  uklad: readonly PanelSection[],
  id: string,
  kierunek: -1 | 1,
): PanelSection[] {
  const wynik = [...uklad];
  const skad = wynik.findIndex((sekcja) => sekcja.id === id);
  if (skad === -1) return wynik;

  // Sąsiad to najbliższa sekcja widoczna w danym kierunku, nie sąsiad w tablicy.
  let dokad = skad + kierunek;
  while (dokad >= 0 && dokad < wynik.length && wynik[dokad]?.hidden === true) dokad += kierunek;
  if (dokad < 0 || dokad >= wynik.length) return wynik;

  const przesuwana = wynik[skad];
  const sasiad = wynik[dokad];
  if (przesuwana === undefined || sasiad === undefined) return wynik;
  wynik[skad] = sasiad;
  wynik[dokad] = przesuwana;
  return ponumeruj(wynik);
}

/** Zwinięcie albo rozwinięcie jednej sekcji; reszta układu bez zmiany. */
export function przelaczZwiniecie(uklad: readonly PanelSection[], id: string): PanelSection[] {
  return uklad.map((sekcja) =>
    sekcja.id === id ? { ...sekcja, collapsed: !sekcja.collapsed } : sekcja,
  );
}

/**
 * Zdjęcie sekcji z widoku albo jej przywrócenie.
 *
 * Zdjęta sekcja zostaje w układzie — traci widoczność, nie miejsce. Gdyby
 * wypadała z tablicy, przywrócenie musiałoby zgadywać, gdzie stała.
 */
export function przelaczZdjecie(uklad: readonly PanelSection[], id: string): PanelSection[] {
  return uklad.map((sekcja) =>
    sekcja.id === id ? { ...sekcja, hidden: sekcja.hidden !== true } : sekcja,
  );
}

/** Czy dwa układy są tożsame — zapis bez skutku nie ma po co iść do rdzenia. */
export function tenSamUklad(a: readonly PanelSection[], b: readonly PanelSection[]): boolean {
  if (a.length !== b.length) return false;
  return a.every((sekcja, indeks) => {
    const druga = b[indeks];
    return (
      druga !== undefined &&
      sekcja.id === druga.id &&
      sekcja.order === druga.order &&
      sekcja.collapsed === druga.collapsed &&
      (sekcja.hidden === true) === (druga.hidden === true)
    );
  });
}
