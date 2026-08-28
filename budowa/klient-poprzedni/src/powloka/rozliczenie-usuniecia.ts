import type { RozliczenieUsuniecia } from './usuniecie-sesji';

// Zdania rozliczenia usunięcia: co rdzeń skasował, a czego nie znalazł spośród wskazanych sesji.

/** Odczyt tytułu karty sesji dla wskazanego identyfikatora sesji; pusty napis znaczy, że tytułu nie znamy. */
export type TytulSesji = (idSesji: string) => string;

/** Wykaz nazw sesji zapisany w cudzysłowie drukarskim, kolejne pozycje rozdzielone przecinkiem i spacją. */
function nazwy(idSesji: readonly string[], tytul: TytulSesji): string {
  return idSesji
    .map((id) => {
      const nazwa = tytul(id).trim();
      return nazwa === '' ? id : `„${nazwa}"`;
    })
    .join(', ');
}

/** Odmiana rzeczownika sesja przez liczbę wskazanych sesji, tak żeby zdanie nigdy nie powiedziało „1 sesje". */
function ileSesji(liczba: number): string {
  if (liczba === 1) return '1 sesję';
  const dziesiatki = liczba % 100;
  const jednosci = liczba % 10;
  const mnoga = jednosci >= 2 && jednosci <= 4 && (dziesiatki < 12 || dziesiatki > 14);
  return `${liczba} ${mnoga ? 'sesje' : 'sesji'}`;
}

/**
 * Zdanie o tym, co rdzeń faktycznie skasował.
 *
 * Licznik bierzemy z `liczba` (pole `deletedCount` rdzenia), a nie z długości
 * wykazu: rozbieżność obu jest wtedy widoczna, a nie zamaskowana.
 */
export function zdanieUsunietych(
  rozliczenie: RozliczenieUsuniecia,
  tytul: TytulSesji,
): string {
  if (rozliczenie.usuniete.length === 0) {
    return 'Rdzeń nie usunął ani jednej sesji — żaden zapis nie zginął.';
  }
  return `Usunięto trwale ${ileSesji(rozliczenie.liczba)}: ${nazwy(rozliczenie.usuniete, tytul)}. Zapisu nie da się przywrócić.`;
}

/**
 * Zdanie o wskazaniach bez odpowiednika w historii albo `null`, gdy takich
 * wskazań nie było. `null`, nie napis pusty — wywołujący ma odróżnić „nie ma
 * o czym mówić" od „mam zdanie puste".
 */
export function zdanieNieznalezionych(
  rozliczenie: RozliczenieUsuniecia,
  tytul: TytulSesji,
): string | null {
  if (rozliczenie.nieznalezione.length === 0) return null;
  return `Bez odpowiednika w historii, więc nic tu nie zginęło (rdzeń nie uznaje tego za błąd): ${nazwy(rozliczenie.nieznalezione, tytul)}.`;
}

/**
 * Oba zdania złożone w jedną treść dla pasa kart, gdzie miejsca jest na jeden
 * napis. Kolejność jest stała: najpierw skutek, potem wskazania pominięte.
 */
export function trescRozliczenia(
  rozliczenie: RozliczenieUsuniecia,
  tytul: TytulSesji,
): string {
  const pominiete = zdanieNieznalezionych(rozliczenie, tytul);
  const pierwsze = zdanieUsunietych(rozliczenie, tytul);
  return pominiete === null ? pierwsze : `${pierwsze} ${pominiete}`;
}

/** Czy rozliczenie usunięcia nadaje się na komunikat o wadze błędu — sprawdza, że nic w nim nie zginęło. */
export function czyRozliczeniePuste(rozliczenie: RozliczenieUsuniecia): boolean {
  return rozliczenie.usuniete.length === 0;
}
