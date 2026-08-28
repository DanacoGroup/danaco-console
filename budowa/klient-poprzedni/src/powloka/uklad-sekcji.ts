import type { PanelSection } from '../../../shared/contract';

// Reguła układu sekcji panelu: kto stoi przed kim, co jest zwinięte, co zdjęte z widoku panelu.

/** Sekcja panelu w postaci, którą podaje moduł: identyfikator sekcji, tytuł oraz gotowa treść tej sekcji. */
export interface OpisSekcji {
  /** Identyfikator sekcji w obrębie panelu — trafia w `PanelSection.id`. */
  id: string;
  /** Nagłówek sekcji widoczny dla Operatora. */
  tytul: string;
  /** Gotowy element treści; sekcja go tylko obudowuje, nigdy nie buduje sama. */
  tresc: HTMLElement;
}

/** Przeliczenie kolejności sekcji na liczby od jeden do n, według bieżącej kolejności w tej samej tablicy. */
function ponumeruj(uklad: readonly PanelSection[]): PanelSection[] {
  return uklad.map((sekcja, indeks) => ({ ...sekcja, order: indeks + 1 }));
}

/** Układ domyślny — kolejność zgłoszenia sekcji przez moduł, bez ani jednej sekcji tutaj zwiniętej wcale. */
export function ukladDomyslny(sekcje: readonly OpisSekcji[]): PanelSection[] {
  return sekcje.map((opis, indeks) => ({ id: opis.id, order: indeks + 1, collapsed: false }));
}

// Zestawienie układu z rdzenia z sekcjami, które ten panel faktycznie buduje po stronie tego samego klienta.
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

/** Przesunięcie sekcji o jedno miejsce w podanym kierunku, liczone tylko względem sekcji widocznych okiem. */
export function przestawSekcje(
  uklad: readonly PanelSection[],
  id: string,
  kierunek: -1 | 1,
): PanelSection[] {
  const wynik = [...uklad];
  const skad = wynik.findIndex((sekcja) => sekcja.id === id);
  if (skad === -1) return wynik;

  // Sąsiad to najbliższa sekcja widoczna w danym kierunku, nie sąsiad sąsiadujący w samej tablicy.
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

/** Zwinięcie albo rozwinięcie jednej wskazanej sekcji panelu; reszta całego układu zostaje bez żadnej zmiany. */
export function przelaczZwiniecie(uklad: readonly PanelSection[], id: string): PanelSection[] {
  return uklad.map((sekcja) =>
    sekcja.id === id ? { ...sekcja, collapsed: !sekcja.collapsed } : sekcja,
  );
}

/** Zdjęcie sekcji z widoku albo jej przywrócenie, bez usuwania samej tej sekcji z całego układu panelu. */
export function przelaczZdjecie(uklad: readonly PanelSection[], id: string): PanelSection[] {
  return uklad.map((sekcja) =>
    sekcja.id === id ? { ...sekcja, hidden: sekcja.hidden !== true } : sekcja,
  );
}

/** Czy dwa układy sekcji są ze sobą tożsame — zapis bez żadnego skutku nie ma po co w ogóle iść do rdzenia. */
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
