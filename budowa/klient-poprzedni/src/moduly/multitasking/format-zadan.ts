/**
 * Zapis liczb i czasów panelu „Zadania w tle" — sam tekst, bez komend
 * i bez tworzenia elementów.
 *
 * Plik stoi osobno, bo zapis czasu trwania jest potrzebny w trzech miejscach
 * naraz (wiersz podsumowania przepływu, kolumna „Czas" tabeli agentów, licznik
 * zwiniętych) i trzy kopie tej samej reguły rozjechałyby się przy pierwszej
 * poprawce.
 *
 * Czas bierze się ze znaczników kontraktu: `Subagent.startedAt`
 * i `Subagent.finishedAt` to milisekundy epoki i oba są nieobowiązkowe.
 * Podagent bez znacznika startu dostaje brak, a nie czas zerowy — zero
 * znaczyłoby „ruszył i nic nie trwał".
 */

/** Zapis czasu trwania po polsku: `24 m 19 s`, `1 h 04 m`, `19 s`. */
export function zapiszCzas(milisekundy: number): string {
  const sekundyRazem = Math.max(0, Math.floor(milisekundy / 1000));
  const godziny = Math.floor(sekundyRazem / 3600);
  const minuty = Math.floor((sekundyRazem % 3600) / 60);
  const sekundy = sekundyRazem % 60;
  if (godziny > 0) return `${godziny} h ${String(minuty).padStart(2, '0')} m`;
  if (minuty > 0) return `${minuty} m ${String(sekundy).padStart(2, '0')} s`;
  return `${sekundy} s`;
}

/**
 * Czas trwania odcinka albo brak.
 *
 * `doKiedy` puste znaczy „praca trwa" — wtedy odcinek liczy się do chwili
 * podanej przez wołającego, a nie do `Date.now()` wziętego wewnątrz. Chwila
 * przychodzi z zewnątrz, żeby wszystkie wiersze jednego przerysowania mierzyły
 * się do tej samej sekundy.
 */
export function czasOdcinka(
  odKiedy: number | undefined,
  doKiedy: number | undefined,
  teraz: number,
): number | null {
  if (odKiedy === undefined || odKiedy <= 0) return null;
  const koniec = doKiedy === undefined || doKiedy <= 0 ? teraz : doKiedy;
  return Math.max(0, koniec - odKiedy);
}

/** Zapis czasu odcinka albo zdanie o braku znacznika — nigdy zero z domysłu. */
export function zapiszCzasOdcinka(czas: number | null): string {
  return czas === null ? 'bez znacznika czasu' : zapiszCzas(czas);
}

/**
 * Liczebnik rzeczownika „agent" w odmianie polskiej.
 *
 * Panel pisze „5 agentów", a nie „5 agent(ów)" — produkt jest polskojęzyczny
 * i odmienia rzeczownik zamiast dopisywać końcówkę w nawiasie.
 */
export function odmienAgentow(liczba: number): string {
  const jednosci = liczba % 10;
  const dziesiatki = liczba % 100;
  if (liczba === 1) return '1 agent';
  if (jednosci >= 2 && jednosci <= 4 && !(dziesiatki >= 12 && dziesiatki <= 14)) {
    return `${liczba} agenci`;
  }
  return `${liczba} agentów`;
}

/**
 * Wskaźnik kropkowy etapu jako tekst: `●●●●○`.
 *
 * Wypełnionych kropek jest tyle, ile wynosi numer etapu; reszta skali zostaje
 * pusta. Wartości spoza zakresu są przycinane, żeby wskaźnik nie wyszedł
 * dłuższy od skali.
 */
export function kropkiEtapu(numerEtapu: number, liczbaEtapow: number): string {
  const skala = Math.max(0, Math.floor(liczbaEtapow));
  if (skala === 0) return '';
  const pelne = Math.min(skala, Math.max(0, Math.floor(numerEtapu)));
  return '●'.repeat(pelne) + '○'.repeat(skala - pelne);
}
