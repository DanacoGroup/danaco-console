/**
 * Zapis liczb i czasów panelu „Zadania w tle" — sam tekst, bez komend i bez
 * tworzenia elementów. Czas bierze się ze znaczników kontraktu
 * `Subagent.startedAt` oraz `Subagent.finishedAt`, podanych w milisekundach
 * epoki i nieobowiązkowych.
 */

/**
 * Zapis czasu trwania po polsku w postaci `24 m 19 s`, `1 h 04 m` albo `19 s`.
 * Jednostka mniejsza jest dopełniana zerem do dwóch cyfr, a wartości ujemne
 * są przycinane do zera.
 */
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
 * Czas trwania odcinka albo brak. Puste `doKiedy` znaczy „praca trwa" —
 * odcinek liczy się wtedy do chwili podanej przez wołającego, nie do
 * `Date.now()` branego wewnątrz. Chwila przychodzi z zewnątrz, aby wiersze
 * przerysowania mierzyły się jednakowo.
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

/**
 * Zapis czasu odcinka albo zdanie o braku znacznika czasu. Brak nigdy nie
 * zamienia się w zero, ponieważ zero znaczyłoby, że praca ruszyła i nie
 * trwała ani chwili.
 */
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
