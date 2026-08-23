import { AgentPermissionGroup, type Agent } from '../../../../shared/contract';
import { pozycjaWykazu, przyciskAkcji as przycisk } from '../../modele/kontrolki-formularza';

/**
 * Jedna pozycja wykazu ekspertów wraz z czterema przełącznikami uprawnień.
 * Osobny plik od okna, bo to inna odpowiedzialność: okno prowadzi odczyt
 * i przypisanie, pozycja rysuje jednego eksperta.
 *
 * Pozycja nie trzyma własnego stanu — rysuje komplet uprawnień z pola
 * `Agent.permissions`, tak jak podał go rdzeń.
 *
 * Uprawnienie ma trzy stany, nie dwa: grupa, o której rdzeń nie powiedział nic,
 * nie jest ani przyznana, ani odebrana i tak też jest wypisana.
 */

/** Cztery grupy zakresu uprawnień eksperta. */
const GRUPY: ReadonlyArray<[AgentPermissionGroup, string]> = [
  [AgentPermissionGroup.Files, 'pliki — odczyt i zapis'],
  [AgentPermissionGroup.Network, 'sieć — połączenia wychodzące'],
  [AgentPermissionGroup.Processes, 'procesy — uruchamianie'],
  [AgentPermissionGroup.Integrations, 'integracje — konektory i MCP'],
];

/** Czynność okna wywoływana przełącznikiem uprawnienia. */
export interface CzynnosciAgenta {
  /** @param chciane stan żądany od rdzenia; wykaz pokaże to, co rdzeń odpowie. */
  zmienUprawnienie(grupa: AgentPermissionGroup, opis: string, chciane: boolean): void;
}

/**
 * Uprawnienie grupy wedle rdzenia. `null` znaczy „rdzeń o tej grupie nie
 * powiedział nic” i nie wolno tego mylić z odebraniem uprawnienia.
 */
export function uprawnienieGrupy(agent: Agent, grupa: AgentPermissionGroup): boolean | null {
  const wiersz = (agent.permissions ?? []).find((pozycja) => pozycja.group === grupa);
  return wiersz === undefined ? null : wiersz.granted;
}

/**
 * @param przypisany `true`/`false` wedle zestawienia projektu z rdzenia,
 *   `null` — gdy rdzeń o przypisaniach tego projektu nie był pytany.
 */
export function pozycjaAgenta(
  agent: Agent,
  przypisany: boolean | null,
  czynnosci: CzynnosciAgenta,
): HTMLElement {
  const { element, akcje } = pozycjaWykazu(
    agent.name,
    [
      opisPrzypisania(przypisany),
      agent.model ?? 'model bazowy nieustalony',
      agent.enabled ? 'czynny' : 'wyłączony',
    ].join(' · '),
    'dw',
  );
  element.dataset['przypisany'] = przypisany === null ? 'nieznane' : String(przypisany);

  for (const [grupa, opis] of GRUPY) {
    const stan = uprawnienieGrupy(agent, grupa);
    const przelacznik = przycisk(etykietaUprawnienia(opis, stan), 'dn-btn dn-btn--zarys');
    przelacznik.dataset['przyznane'] = stan === null ? 'nieznane' : stan ? 'tak' : 'nie';
    // Naciśnięcie ŻĄDA stanu przeciwnego do podanego przez rdzeń; grupa bez
    // odpowiedzi rdzenia daje żądanie przyznania. Nowego stanu przycisk sobie
    // nie wpisuje — wykaz odrysowuje okno z odpowiedzi `agent.permission.set`.
    przelacznik.addEventListener('click', () =>
      czynnosci.zmienUprawnienie(grupa, opis, stan !== true),
    );
    akcje.append(przelacznik);
  }
  return element;
}

/** Zdanie o przypisaniu do projektu; niewiedza nazywa się niewiedzą. */
function opisPrzypisania(przypisany: boolean | null): string {
  if (przypisany === null) return 'przypisanie do projektu nieodczytane';
  return przypisany ? 'przypisany do projektu' : 'w bibliotece, nieprzypisany do projektu';
}

/** Etykieta przełącznika: stan wedle rdzenia i czynność, którą wykona naciśnięcie. */
function etykietaUprawnienia(opis: string, stan: boolean | null): string {
  if (stan === null) return `Uprawnienie: ${opis} — rdzeń nie podał (przyznaj)`;
  return stan
    ? `Uprawnienie: ${opis} — przyznane (odbierz)`
    : `Uprawnienie: ${opis} — odebrane (przyznaj)`;
}
