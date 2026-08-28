import type { WindowRoleAssignment } from '../../../../shared/contract';
import {
  przyciskAkcji as przycisk,
  przelacznikWidoku,
  przestaw,
  pozycjaWykazu,
  wykaz,
} from '../../modele/kontrolki-formularza';
import { PRZEDROSTEK } from './kontrolki';
import { powod } from './nadanie-rol';
import { utworzStanTresci } from './stany-okna';
import type { StanMultitaskingu } from './stan-multitaskingu';
import type { ZrodloOkien } from './zrodlo-okien';

/**
 * Wykaz nadań ról jest jedynym czytelnikiem nadań i jedynym miejscem, gdzie
 * widać wcielenie z rdzenia; zawężenie do wykonawców koordynatora wykonuje
 * rdzeń, a zdjęcie roli działa bez bramek, bo skutek nie jest utratą zapisu.
 */
export interface WykazNadanRol {
  element: HTMLElement;
  /** Ponowny odczyt nadań z rdzenia. */
  odswiez(): void;
}

export interface OpcjeWykazuNadan {
  zrodlo: ZrodloOkien;
  stan: StanMultitaskingu;
  /** Wywoływane po skutecznym zdjęciu roli — scena musi przeczytać okna od nowa. */
  poZdjeciuRoli(): void;
}

export function utworzWykazNadanRol(opcje: OpcjeWykazuNadan): WykazNadanRol {
  const { zrodlo, stan, poZdjeciuRoli } = opcje;

  const tresci = utworzStanTresci();
  const pozycje = wykaz('Nadania ról widziane przez rdzeń', `${PRZEDROSTEK}-wykaz`);

  const zawez = przelacznikWidoku('Tylko wykonawcy tego koordynatora', false);
  zawez.addEventListener('click', () => {
    przestaw(zawez);
    void odczytaj();
  });

  const odczytajPonownie = przycisk('Odczytaj nadania', 'dn-btn dn-btn--sm');
  odczytajPonownie.addEventListener('click', () => {
    void odczytaj();
  });

  const pasek = document.createElement('div');
  pasek.className = `${PRZEDROSTEK}-nadania__pasek`;
  pasek.append(zawez, odczytajPonownie);

  const element = document.createElement('div');
  element.className = `${PRZEDROSTEK}-nadania`;
  element.append(pasek, pozycje, tresci.element);

  // Odczyt nadań; zawężenie idzie do rdzenia tylko wtedy, gdy obsada ma koordynatora.
  async function odczytaj(): Promise<void> {
    const koordynator = stan.obsada().koordynator;
    const zadane = zawez.dataset['wlaczony'] === 'true';
    const zawezone = zadane && koordynator !== null;
    // Zawężenie zadane, a niewykonane, nie ma prawa zniknąć w ciszy przed przeglądającym wykaz.
    const oZawezeniu =
      zadane && !zawezone
        ? ' Zawężenie do wykonawców koordynatora NIE poszło — obsada nie ma koordynatora, więc nie ma adresu zawężenia.'
        : '';
    tresci.ladowanie('Odczyt nadań ról…');

    const wynik = await zrodlo.nadaniaRol(
      zawezone && koordynator !== null ? { coordinatorWindowId: koordynator.id } : {},
    );
    if (!wynik.udany || wynik.wynik === undefined) {
      pozycje.replaceChildren();
      tresci.blad('Rdzeń odmówił wykazu nadań ról.', wynik.blad);
      return;
    }
    const nadania = wynik.wynik;
    pozycje.replaceChildren(
      ...nadania.map((nadanie) =>
        pozycjaNadania(nadanie, () => {
          void zdejmij(nadanie);
        }),
      ),
    );
    tresci.pusto(
      `${
        nadania.length === 0
          ? 'Rdzeń nie ma ani jednego nadania roli w tym zawężeniu.'
          : `Rdzeń oddał ${nadania.length} nadań ról.`
      }${oZawezeniu}`,
    );
  }

  /** Zdjęcie roli — bez potwierdzania; zdanie powstaje z pola `removed`. */
  async function zdejmij(nadanie: WindowRoleAssignment): Promise<void> {
    const wynik = await zrodlo.zdejmijRole({ windowId: nadanie.windowId });
    if (!wynik.udany || wynik.wynik === undefined) {
      tresci.potwierdzenie(
        `Rdzeń odmówił zdjęcia roli z okna ${nadanie.windowId}. ${powod(wynik.blad)}`,
        false,
      );
      return;
    }
    const oddane = wynik.wynik;
    if (!oddane.removed) {
      tresci.potwierdzenie(
        `Rdzeń przyjął wywołanie, ale oddał „nie zdjęto" dla okna ${oddane.windowId} — rola została na miejscu.`,
        false,
      );
      await odczytaj();
      return;
    }
    tresci.potwierdzenie(`Rdzeń zdjął rolę z okna ${oddane.windowId}.`, true);
    // Scena stoi na wykazie okien, więc po zdjęciu roli musi przeczytać okna od nowa.
    poZdjeciuRoli();
  }

  return {
    element,
    odswiez: () => {
      void odczytaj();
    },
  };
}

/**
 * Jedna pozycja wykazu nadań.
 *
 * Czysta konstrukcja: zdjęcie roli oddaje wywołaniem zwrotnym, więc nie zna ani
 * źródła rdzenia, ani stanu sceny. Opis niesie wcielenie — po to ten wykaz jest.
 */
function pozycjaNadania(nadanie: WindowRoleAssignment, naZdjecie: () => void): HTMLElement {
  const opis = [
    `rola ${nadanie.role}`,
    `wcielenie ${nadanie.persona === undefined || nadanie.persona === '' ? 'nie nadane' : nadanie.persona}`,
    `koordynator ${nadanie.coordinatorWindowId ?? 'brak'}`,
  ].join(' · ');

  const { element, akcje } = pozycjaWykazu(nadanie.windowId, opis, PRZEDROSTEK);
  const zdejmij = przycisk('Zdejmij rolę', 'dn-btn dn-btn--sm dn-btn--zarys');
  zdejmij.title = `Zdejmij rolę ${nadanie.role} z okna ${nadanie.windowId} — okno zostaje, traci wyłącznie rolę`;
  zdejmij.addEventListener('click', naZdjecie);
  akcje.append(zdejmij);
  return element;
}
