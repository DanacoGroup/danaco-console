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
 * Wykaz nadań ról — jedyny czytelnik `role.list` i jedyny wołacz `role.remove`
 * w kliencie.
 *
 * Wykaz bierze się z rdzenia, a nie z okien sesji. Panel obsady składa scenę
 * z `window.list`: bierze okna i odczytuje z nich pole `windowRole`. To wystarcza
 * do pokazania sceny, ale nie jest rejestrem nadań. Rejestrem jest rdzeń, a jego
 * odpowiedź (`WindowRoleAssignment`) niesie wcielenie — pole, którego `Window`
 * nie ma w ogóle. Ten wykaz jest jedynym miejscem, gdzie widać, czy wcielenie
 * nadane przez `role.update` w rdzeniu stoi.
 *
 * Zawężenie jest czynnością, nie filtrem miejscowym: `role.list` przyjmuje
 * `coordinatorWindowId`, więc zawężenie do wykonawców tej obsady wykonuje rdzeń.
 * Przesiewanie odpowiedzi w kliencie dałoby ten sam obraz, ale kłamałoby przy
 * wykazie uciętym po stronie rdzenia — i byłoby drugą regułą przynależności
 * do pętli obok tej z kontraktu.
 *
 * „Zdejmij rolę" wykonuje się bez bramek: bez pytania „czy na pewno?", bez
 * wygaszania i bez uprawnienia per okno. Skutek nie jest utratą zapisu — okno
 * zostaje, traci wyłącznie rolę, a nadać ją z powrotem można paskiem obsady
 * stojącym wyżej w tym samym panelu.
 *
 * Zgoda nie jest dowodem: `role.remove` odpowiada polem `removed`, a odpowiedź
 * udana z `removed: false` znaczy „nie było czego zdjąć". Zdanie dla Operatora
 * powstaje z tego pola, nie z faktu, że wywołanie nie zwróciło błędu.
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

  /**
   * Odczyt nadań. Zawężenie idzie do rdzenia tylko wtedy, gdy obsada ma
   * koordynatora — pole `coordinatorWindowId` z pustym napisem byłoby zawężeniem
   * do okna, którego nie ma, i oddałoby wykaz pusty bez powodu.
   */
  async function odczytaj(): Promise<void> {
    const koordynator = stan.obsada().koordynator;
    const zadane = zawez.dataset['wlaczony'] === 'true';
    const zawezone = zadane && koordynator !== null;
    // Zawężenie zadane, a niewykonane, nie ma prawa zniknąć w ciszy:
    // wykaz byłby wtedy szerszy, niż Operator prosił, i nikt by tego nie wiedział.
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
    // Scena stoi na `window.list`, więc po zdjęciu roli musi przeczytać okna
    // od nowa — inaczej pokazywałaby obsadę, której w rdzeniu już nie ma.
    // Ten odczyt zawraca też tutaj (panel odświeża wykaz nadań po odczycie
    // okien), więc drugiego wywołania stąd nie ma: byłoby tym samym pytaniem
    // zadanym dwa razy pod rząd.
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
