import type { Agent, AgentPermissionGroup } from '../../../../shared/contract';
import type { StanTresci } from './stany-okna';
import type { ZrodloWorkspace } from './zrodlo-workspace';

/**
 * Czynności okna Agent Manager sięgające poza sam wykaz: zmiana uprawnienia
 * eksperta i zdanie o przypisaniach projektu.
 *
 * Osobny plik od okna, bo to inna odpowiedzialność: okno składa kontrolki
 * i prowadzi odczyt, tu leży przebieg czynności wraz z tym, co po niej widzi
 * operator. Wzór ten sam co w `biblioteka-czynnosci.ts`.
 *
 * Potwierdzenie mówi to, co zapisał rdzeń, nie to, czego żądało okno: rdzeń
 * oddaje po zapisie komplet uprawnień eksperta i ten komplet rozstrzyga.
 * Po odpowiedzi okno odrysowuje wykaz, bo stan ładowania go opróżnia.
 */

/** Zależności czynności okna. */
export interface ZaleznosciAgentow {
  zrodlo: ZrodloWorkspace;
  tresc: StanTresci;
  /** Odrysowuje wykaz ekspertów z danych, które okno ma po odpowiedzi rdzenia. */
  odrysuj(): void;
}

/**
 * Zmienia uprawnienie eksperta i podaje wynik wedle rdzenia. Oddany komplet
 * uprawnień wchodzi do wykazu, żeby wykaz i zdanie mówiły to samo.
 */
export function zmienUprawnienieEksperta(
  z: ZaleznosciAgentow,
  agent: Agent,
  grupa: AgentPermissionGroup,
  opis: string,
  chciane: boolean,
): void {
  z.tresc.ladowanie(`Zmiana uprawnienia „${opis}” eksperta ${agent.name}…`);
  void z.zrodlo.ustawUprawnienie(agent.id, grupa, chciane).then((wynik) => {
    if (!wynik.udany || wynik.wynik === undefined) {
      z.tresc.blad(`Rdzeń nie zmienił uprawnienia eksperta ${agent.name}.`, wynik.blad);
      return;
    }
    agent.permissions = wynik.wynik;
    z.odrysuj();
    const poZapisie = wynik.wynik.find((pozycja) => pozycja.group === grupa);
    if (poZapisie === undefined) {
      z.tresc.potwierdzenie(
        `Rdzeń przyjął wywołanie, ale w odpowiedzi nie ma grupy „${opis}” — ` +
          'nie wiadomo, co obowiązuje.',
        false,
      );
      return;
    }
    z.tresc.potwierdzenie(
      `Rdzeń zapisał uprawnienie „${opis}” eksperta ${agent.name} jako ` +
        `${poZapisie.granted ? 'przyznane' : 'odebrane'}` +
        `${poZapisie.granted === chciane ? '' : ' — INACZEJ, niż żądało okno'}.`,
      poZapisie.granted === chciane,
    );
  });
}

/**
 * Zdanie o przypisaniach projektu — osobne dla braku odczytu i dla odczytanego
 * zera. `null` znaczy „rdzeń nie był pytany” i nie wolno tego mylić z brakiem
 * przypisanych ekspertów.
 */
export function zdanieOPrzypisaniach(przypisani: ReadonlySet<string> | null): HTMLElement {
  const zdanie = document.createElement('p');
  zdanie.className = 'dn-pole-opis';
  if (przypisani === null) {
    zdanie.textContent =
      'Przypisań do projektu nie odczytano — rdzeń nie był o nie pytany. ' +
      'Wskaż projekt na pulpicie i odczytaj zestawienie.';
    return zdanie;
  }
  zdanie.textContent =
    przypisani.size === 0
      ? 'Rdzeń nie podał w zestawieniu projektu ani jednego przypisanego eksperta.'
      : `Rdzeń podał ${przypisani.size} eksperta/ów przypisanych do projektu.`;
  return zdanie;
}
