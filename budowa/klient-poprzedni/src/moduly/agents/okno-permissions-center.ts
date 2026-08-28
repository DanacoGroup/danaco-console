import { AgentPermissionGroup, type Agent } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import { poleTekstowe, przycisk, utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import type { Kanal } from '../../protokol/kanal';
import { naglowek } from './biblioteka-ekspertow';
import { utworzPanelZakresowNarzedzi, type PanelZakresowNarzedzi } from './panel-zakresow-narzedzi';
import { utworzPanelZakresuEksperta, type PanelZakresuEksperta } from './panel-zakresu-eksperta';
import { utworzZrodloZakresuEksperta, type ZrodloZakresuEksperta } from './zrodlo-zakresu-eksperta';
import { utworzZasiegOkien, type ZasiegOkien } from './zasieg-okien';
import { utworzStanOkna, type StanOkna } from './stan-okna';
import type { StanAgentow } from './stan-agentow';
import type { ZrodloZaplecza } from './zrodlo-zaplecza';

/**
 * Permissions Center jest oknem zarządcą modułu Agents: konfiguruje
 * możliwości eksperta w grupach zakresu i tryb uprawnień okna komunikacji,
 * nie kontrolę dostępu do aplikacji.
 */
export interface OknoPermissionsCenter {
  element: HTMLElement;
  odswiez(): void;
  /** Odczytuje okna komunikacji wskazanej sesji. */
  wczytajOkna(idSesji: string): Promise<void>;
  /** Odpina wykaz pokrycia od wspólnego bytu. Wołane z `rozlacz()` modułu. */
  zamknij(): void;
}

/** Grupy zakresu wraz z ich znaczeniem dla Operatora, wyświetlane w kolejności ustalonej dla tego okna. */
const OPISY_GRUP: Record<AgentPermissionGroup, string> = {
  [AgentPermissionGroup.Files]: 'odczyt i zapis plików',
  [AgentPermissionGroup.Network]: 'dostęp sieciowy',
  [AgentPermissionGroup.Processes]: 'uruchamianie procesów',
  [AgentPermissionGroup.Integrations]: 'konektory i serwery MCP',
  [AgentPermissionGroup.Modules]: 'moduły i zasoby platformy',
};

export function utworzOknoPermissionsCenter(
  stan: StanAgentow,
  zaplecze: ZrodloZaplecza,
  kanal: Kanal,
): OknoPermissionsCenter {
  const okno: StanOkna = utworzStanOkna();
  const zasieg: ZasiegOkien = utworzZasiegOkien(zaplecze);
  const zrodloZakresu: ZrodloZakresuEksperta = utworzZrodloZakresuEksperta(kanal);
  // Trzy grupy zakresu, których wiersz przełącznika unieść nie potrafi, mają własny panel.
  const zakresDzialania: PanelZakresuEksperta = utworzPanelZakresuEksperta(stan, zrodloZakresu);
  const zakresyNarzedzi: PanelZakresowNarzedzi = utworzPanelZakresowNarzedzi(zrodloZakresu);

  // Baner jest elementem obowiązkowym okna: bez niego zawężenie zakresu czyta się jak odebranie dostępu.
  const baner = document.createElement('div');
  baner.className = 'da-baner da-baner--wyjsciowy';
  baner.setAttribute('role', 'note');

  const trescBaneru = document.createElement('p');
  trescBaneru.className = 'da-baner__tresc';
  trescBaneru.textContent =
    'Stan wyjściowy: PEŁNY DOSTĘP OPERACYJNY. Centrum uprawnień nie kontroluje dostępu ' +
    'do aplikacji — konfiguruje zakres działania TEGO WYKONAWCY w imieniu już ' +
    'zalogowanego Operatora. Brak wpisu znaczy wartość domyślną, nigdy blokadę.';

  const reset = przycisk('Resetuj do pełnego dostępu', 'dn-btn dn-btn--sm dn-btn--zarys');
  baner.append(trescBaneru, reset);

  const grupy = document.createElement('ul');
  grupy.className = 'da-uprawnienia';

  const zakres = poleTekstowe({
    etykieta: 'Zakres szczegółowy w obrębie grupy',
    podpowiedz: 'puste znaczy cała grupa',
    opis:
      'W grupie modułów zakresem szczegółowym jest kod modułu platformy — tędy zawęża ' +
      'się dostęp eksperta do pojedynczego modułu, bez powielania wykazu modułów w tym oknie.',
  });

  const odpowiedz = utworzWierszOdpowiedzi();

  okno.tresc.append(
    baner,
    grupy,
    zakres.element,
    odpowiedz.element,
    zakresDzialania.element,
    zakresyNarzedzi.element,
    zasieg.element,
  );

  const element = document.createElement('section');
  element.className = 'da-okno da-okno--zarzadca';
  element.dataset['okno'] = 'permissions-center';
  element.append(naglowek('Permissions Center'), okno.element);

  async function ustaw(grupa: AgentPermissionGroup, przyznane: boolean): Promise<void> {
    const ekspert = stan.wybrany();
    if (ekspert === null) return;
    const wskazanyZakres = zakres.kontrolka.value.trim();
    const czego = wskazanyZakres === '' ? `grupy ${grupa}` : `zakresu ${wskazanyZakres} w grupie ${grupa}`;
    odpowiedz.pokaz(`Zapis uprawnienia ${czego}…`, true);
    const wynik = await stan.zrodlo.ustawUprawnienie({
      idEksperta: ekspert.id,
      grupa,
      przyznane,
      zakres: wskazanyZakres,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      // Wiersz wraca do stanu rdzenia: przestawienie przeglądarki nie jest świadectwem zapisu.
      odswiez();
      odpowiedz.pokaz(opisOdmowy('Zapis uprawnienia', wynik.blad?.code, wynik.blad?.message), false);
      return;
    }
    const oddane = wynik.wynik.permissions;
    // Wchłonięcie idzie przed rozstrzygnięciem: wiersze mają pokazać stan rdzenia niezależnie od zapisu.
    stan.wchlon({ ...ekspert, permissions: oddane });
    // Świadectwem zapisu jest wpis, który wrócił — o tej grupie, o tym zakresie
    // i o tej wartości.
    const wpis = oddane.find(
      (kandydat) => kandydat.group === grupa && (kandydat.scope ?? '') === wskazanyZakres,
    );
    if (wpis === undefined) {
      odpowiedz.pokaz(
        `Rdzeń przyjął wywołanie, ale w wykazie, który oddał, nie ma wpisu ${czego} — ` +
          'zapis się nie odbył.',
        false,
      );
      return;
    }
    if (wpis.granted !== przyznane) {
      odpowiedz.pokaz(
        `Rdzeń przyjął wywołanie, ale wpis ${czego} stoi w wykazie, który oddał, jako ` +
          `${wpis.granted ? 'przyznany' : 'odebrany'} — zapis żądanej wartości się nie odbył.`,
        false,
      );
      return;
    }
    odpowiedz.pokaz(
      `Uprawnienie ${czego} ${wpis.granted ? 'przyznane' : 'odebrane'} — ekspert ma teraz ` +
        `${oddane.length} wpisów zakresu.` +
        (wskazanyZakres === ''
          ? ''
          : ' Przełącznik całej grupy zostaje bez zmian: rdzeń zapisał wyłącznie ten zakres.'),
      true,
    );
  }

  function wierszGrupy(grupa: AgentPermissionGroup, ekspert: Agent): HTMLElement {
    const wpisy = (ekspert.permissions ?? []).filter((wpis) => wpis.group === grupa);
    const calaGrupa = wpisy.find((wpis) => (wpis.scope ?? '') === '');
    // Brak wiersza znaczy „nie rozstrzygnięto”; rozstrzygnięciem wyjściowym jest pełny dostęp operacyjny.
    const przyznane = calaGrupa?.granted ?? true;

    const przelacznik = document.createElement('input');
    przelacznik.type = 'checkbox';
    przelacznik.className = 'dn-przelacznik';
    przelacznik.checked = przyznane;
    przelacznik.addEventListener('change', () => void ustaw(grupa, przelacznik.checked));

    const etykieta = document.createElement('span');
    etykieta.className = 'da-uprawnienia__etykieta';
    etykieta.textContent = `${grupa} — ${OPISY_GRUP[grupa]}`;

    const szczegolowe = document.createElement('span');
    szczegolowe.className = 'da-uprawnienia__zakresy';
    const zawezenia = wpisy.filter((wpis) => (wpis.scope ?? '') !== '');
    szczegolowe.textContent =
      zawezenia.length === 0
        ? 'bez zakresów szczegółowych'
        : zawezenia
            .map((wpis) => `${wpis.scope ?? ''}: ${wpis.granted ? 'tak' : 'nie'}`)
            .join(' · ');

    const element = document.createElement('li');
    element.className = 'da-uprawnienia__wiersz';
    element.dataset['grupa'] = grupa;
    element.append(przelacznik, etykieta, szczegolowe);
    return element;
  }

  /** Przywrócenie pełnego dostępu jednym wywołaniem: reset zdejmuje wpisy, a nie przyznaje je z powrotem. */
  async function resetuj(): Promise<void> {
    const ekspert = stan.wybrany();
    if (ekspert === null) {
      odpowiedz.pokaz('Wybierz eksperta w Agent Builderze — reset dotyczy jednego.', false);
      return;
    }
    odpowiedz.pokaz('Przywracanie pełnego dostępu…', true);
    const wynik = await zrodloZakresu.zdejmijUprawnienie({ idEksperta: ekspert.id });
    if (!wynik.udany || wynik.wynik === undefined) {
      odswiez();
      odpowiedz.pokaz(
        opisOdmowy('Przywrócenie pełnego dostępu', wynik.blad?.code, wynik.blad?.message),
        false,
      );
      return;
    }
    stan.wchlon({ ...ekspert, permissions: wynik.wynik.permissions });
    odswiez();
    odpowiedz.pokaz(
      wynik.wynik.removed === 0
        ? 'Już pełny dostęp, brak zmian — ekspert nie miał ani jednego zawężenia.'
        : `Przywrócono pełny dostęp — zdjęto ${wynik.wynik.removed} wpisów. Wiersze zakresu ` +
            'znikają w całości, więc wraca stan „brak ustawienia = wartość domyślna”.',
      true,
    );
  }

  reset.addEventListener('click', () => void resetuj());

  /** Ekspert, którego wiersze stoją w oknie — po nim poznajemy zmianę wyboru. */
  let pokazany = '';

  function odswiez(): void {
    const ekspert = stan.wybrany();
    // Zakres wpisany dla jednego eksperta nie może czekać w polu przy drugim ekspercie.
    if ((ekspert?.id ?? '') !== pokazany) {
      pokazany = ekspert?.id ?? '';
      zakres.kontrolka.value = '';
      odpowiedz.wyczysc();
    }
    if (ekspert === null) {
      grupy.replaceChildren();
      okno.puste('Wybierz eksperta w Agent Builderze, aby ustalić jego uprawnienia.');
      return;
    }
    grupy.replaceChildren(
      ...Object.values(AgentPermissionGroup).map((grupa) => wierszGrupy(grupa, ekspert)),
    );
    zakresDzialania.odswiez();
    okno.gotowe();
  }

  return {
    element,

    odswiez,

    async wczytajOkna(idSesji) {
      void zakresyNarzedzi.wczytaj();
      await zasieg.wczytaj(idSesji);
    },

    // Okno nie trzyma już wykazu pokrycia: wszystkie cztery grupy zakresu mają drogę do rdzenia.
    zamknij: () => undefined,
  };
}
