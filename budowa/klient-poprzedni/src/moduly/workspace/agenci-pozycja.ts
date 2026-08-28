import { AgentPermissionGroup, type Agent } from '../../../../shared/contract';
import { pozycjaWykazu, przyciskAkcji as przycisk } from '../../modele/kontrolki-formularza';

/** Cztery grupy zakresu uprawnień eksperta wraz z opisem każdej z nich, wyświetlanym w oknie dla Operatora. */
const GRUPY: ReadonlyArray<[AgentPermissionGroup, string]> = [
  [AgentPermissionGroup.Files, 'pliki — odczyt i zapis'],
  [AgentPermissionGroup.Network, 'sieć — połączenia wychodzące'],
  [AgentPermissionGroup.Processes, 'procesy — uruchamianie'],
  [AgentPermissionGroup.Integrations, 'integracje — konektory i MCP'],
];

/** Czynność okna wywoływana naciśnięciem przełącznika uprawnienia jednej z czterech grup uprawnień eksperta. */
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
    // Naciśnięcie żąda stanu przeciwnego do podanego przez rdzeń; przycisk stanu sobie nie wpisuje.
    przelacznik.addEventListener('click', () =>
      czynnosci.zmienUprawnienie(grupa, opis, stan !== true),
    );
    akcje.append(przelacznik);
  }
  return element;
}

/** Buduje zdanie o przypisaniu eksperta do projektu; brak odczytu przypisania jest nazwany wprost jako niewiedza. */
function opisPrzypisania(przypisany: boolean | null): string {
  if (przypisany === null) return 'przypisanie do projektu nieodczytane';
  return przypisany ? 'przypisany do projektu' : 'w bibliotece, nieprzypisany do projektu';
}

/** Buduje etykietę przełącznika, niosącą stan uprawnienia wedle rdzenia i czynność, którą wykona naciśnięcie. */
function etykietaUprawnienia(opis: string, stan: boolean | null): string {
  if (stan === null) return `Uprawnienie: ${opis} — rdzeń nie podał (przyznaj)`;
  return stan
    ? `Uprawnienie: ${opis} — przyznane (odbierz)`
    : `Uprawnienie: ${opis} — odebrane (przyznaj)`;
}
