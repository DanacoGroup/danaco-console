import './przelacznik-srodowiska.css';

import { ExecutionEnv } from '../../../shared/contract';
import { utworzMenuDrzewo, type PozycjaMenu } from '../komponenty/menu-drzewo';

// Urządzenie to pierwszy komponent paska zlecenia, wybierający jednym ruchem maszynę wykonania modelu.

/** Migawka stanu, z której komponent czerpie całą swoją treść: zasięg wykonania oraz nazwę hosta zdalnego. */
export interface StanPrzelacznika {
  /** Zasięg wykonania okna — wartość wyliczenia kontraktu. */
  srodowisko: ExecutionEnv;
  /** Nazwa hosta z ustawienia `host_wykonania`; pusta, gdy nie wskazano. */
  host: string;
}

/** Zależności komponentu, wstrzykiwane przez montaż okna: migawka stanu, subskrypcja zmian oraz zastosowanie wyboru. */
export interface ZaleznosciPrzelacznika {
  /** Bieżąca migawka stanu okna. */
  migawka(): StanPrzelacznika;
  /** Subskrypcja zmian stanu (window.changed, config.changed). */
  naZmiane(sluchacz: () => void): void;
  /** Zastosowanie wyboru komendą window.update; odrzucenie niesie zdanie odmowy od rdzenia. */
  zastosuj(srodowisko: ExecutionEnv): Promise<void>;
}

/** Komponent „urządzenie” paska zlecenia, pokazujący maszynę wykonania i pozwalający ją zmienić jednym ruchem. */
export interface PrzelacznikSrodowiska {
  element: HTMLElement;
  /** Przerysowuje etykietę uchwytu i drzewo z bieżącej migawki. */
  odswiez(): void;
}

/**
 * Zdanie odmowy dla wyboru hosta zdalnego bez wskazanej maszyny, eksportowane, bo sprawdza je test jednostkowy.
 */
export const ZDANIE_BRAKU_HOSTA =
  'Okno nie zostało przełączone na host zdalny. ' +
  'Pole „Host wykonania” jest puste, a przełącznik wskazuje maszyny, nie nazwy wartości — ' +
  'nie ma więc dokąd przełączyć. ' +
  'Wpisz nazwę hosta w polu „Host wykonania” panelu sterowania i ponów wybór.';

/** Etykieta maszyny dla braku wskazania hosta zdalnego — mówi wprost o braku, niczego nie udaje operatorowi. */
export const ETYKIETA_BRAKU_HOSTA = 'Host zdalny — nie wskazano';

/**
 * Objaśnienia pozycji menu mówią o skutku wyboru, nie o nazwie wartości — maszynę wybiera się po tym, co się na niej stanie.
 */
const OPISY: Record<ExecutionEnv, string> = {
  [ExecutionEnv.Local]:
    'Model pracuje na tym komputerze i sięga po jego pliki, procesy i sieć.',
  [ExecutionEnv.Core]:
    'Model pracuje na maszynie, na której stoi rdzeń — bez dostępu do tego komputera.',
  [ExecutionEnv.Remote]:
    'Model pracuje na maszynie wskazanej w polu „Host wykonania” panelu sterowania.',
};

export function utworzPrzelacznikSrodowiska(
  zaleznosci: ZaleznosciPrzelacznika,
): PrzelacznikSrodowiska {
  const element = document.createElement('div');
  element.className = 'dc-przelacznik-srodowiska';

  const odmowa = document.createElement('p');
  odmowa.className = 'dn-plakietka dn-plakietka--blad';
  odmowa.setAttribute('role', 'alert');
  odmowa.hidden = true;

  const menu = utworzMenuDrzewo({
    nastawa: 'Urządzenie',
    ikona: 'cpu',
    naWybor: (klucz) => {
      void wybierz(klucz as ExecutionEnv);
    },
  });

  element.append(menu.element, odmowa);

  /** Nazwa MASZYNY pokazywana dla danego zasięgu. */
  function nazwaMaszyny(srodowisko: ExecutionEnv, host: string): string {
    switch (srodowisko) {
      case ExecutionEnv.Local:
        return 'To urządzenie';
      case ExecutionEnv.Core:
        return 'Maszyna rdzenia';
      case ExecutionEnv.Remote:
        return host !== '' ? host : ETYKIETA_BRAKU_HOSTA;
    }
  }

  function odswiez(): void {
    const stan = zaleznosci.migawka();
    const host = stan.host.trim();
    const drzewo: PozycjaMenu[] = [
      ExecutionEnv.Local,
      ExecutionEnv.Core,
      ExecutionEnv.Remote,
    ].map((srodowisko) => ({
      rodzaj: 'wybor',
      klucz: srodowisko,
      nazwa: nazwaMaszyny(srodowisko, host),
      opis: OPISY[srodowisko],
      wybrany: srodowisko === stan.srodowisko,
    }));
    // Uchwyt niesie wartość — nazwę maszyny, na której okno pracuje w tej
    // chwili.
    menu.ustaw(nazwaMaszyny(stan.srodowisko, host), drzewo);
  }

  function pokazOdmowe(zdanie: string): void {
    odmowa.textContent = zdanie;
    odmowa.hidden = false;
  }

  async function wybierz(srodowisko: ExecutionEnv): Promise<void> {
    const stan = zaleznosci.migawka();
    if (srodowisko === stan.srodowisko) return;
    if (srodowisko === ExecutionEnv.Remote && stan.host.trim() === '') {
      pokazOdmowe(ZDANIE_BRAKU_HOSTA);
      return;
    }
    odmowa.hidden = true;
    element.setAttribute('aria-busy', 'true');
    try {
      await zaleznosci.zastosuj(srodowisko);
    } catch (blad) {
      // Odmowa rdzenia idzie do operatora dosłownie, bez parafrazy i bez przełączenia wyróżnienia.
      pokazOdmowe(blad instanceof Error ? blad.message : String(blad));
    } finally {
      element.removeAttribute('aria-busy');
      odswiez();
    }
  }

  zaleznosci.naZmiane(odswiez);
  odswiez();

  return { element, odswiez };
}
