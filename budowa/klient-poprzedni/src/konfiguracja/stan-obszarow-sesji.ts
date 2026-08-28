import {
  ErrorCode,
  SessionConfigArea,
  type ConfigSessionSetResponse,
  type ErrorInfo,
  type SessionConfigEffective,
  type SessionConfigFieldCapability,
  type SessionConfigHooks,
  type SessionConfigOrigin,
} from '../../../shared/contract';
import { opisOdmowyBledu } from '../komponenty/odmowa';
import type { Kanal, Wynik } from '../protokol/kanal';
import { adresPoczatkowy } from './adres-ustawienia';
import type { PunktWidzenia } from './rozstrzygniecie';
import type { FazaOdczytu } from './stan-konfiguracji';
import { utworzZrodloObszarowSesji, type ZrodloObszarowSesji } from './zrodlo-obszarow-sesji';

/**
 * Stan panelu obszarów konfiguracji sesji — jedno źródło prawdy dla wykazu
 * obszarów, plakietek pochodzenia i zapisu. Żadna ścieżka nie zatrzymuje
 * okna: odmowa odczytu zostawia wykaz pusty, odmowa zapisu wraca `Wynikiem`
 * z błędem.
 */
export interface StanObszarowSesji {
  /** Faza odczytu konfiguracji obowiązującej. */
  faza(): FazaOdczytu;
  /** Powód niepowodzenia odczytu; pusty, gdy odczyt się powiódł. */
  powodNiepowodzenia(): string;
  /** Konfiguracja obowiązująca; `null`, dopóki rdzeń jej nie oddał. */
  obowiazujaca(): SessionConfigEffective | null;
  /** Pochodzenie obszaru; `null`, gdy rdzeń o tym obszarze nie mówił. */
  pochodzenie(obszar: SessionConfigArea): SessionConfigOrigin | null;
  /** Obszary zapisane na poziomie punktu widzenia po ostatnim zapisie. */
  zapisane(): readonly SessionConfigArea[];
  /** Pola, których adapter dostawcy nie obsłuży; puste, gdy rdzeń nic nie zgłosił. */
  nieobsluzone(): readonly SessionConfigFieldCapability[];
  /** Punkt widzenia, dla którego liczona jest konfiguracja obowiązująca. */
  punkt(): PunktWidzenia;
  /** Zmienia punkt widzenia; nie czyta rdzenia — o to prosi się `odswiez`. */
  ustawPunkt(punkt: PunktWidzenia): void;
  /** `config.effective.get` — rozstrzygnięcie kompletu obszarów. */
  odswiez(): Promise<void>;
  /** `config.session.set` — utrwala wskazane obszary na poziomie punktu. */
  utrwal(obszary: readonly SessionConfigArea[]): Promise<Wynik<ConfigSessionSetResponse>>;
  /** `config.session.set` — zapisuje obszar `hooks` treścią złożoną przez Operatora w panelu zaczepów. */
  utrwalZaczepy(zaczepy: SessionConfigHooks): Promise<Wynik<ConfigSessionSetResponse>>;
  /** Subskrypcja przeliczenia stanu. */
  naZmiane(sluchacz: () => void): void;
}

/**
 * Zapis stanu trzymany osobno od wytwórni.
 *
 * Dzięki temu dwie czynności rdzeniowe — odczyt i zapis — są zwykłymi funkcjami
 * modułu, a nie ciałami domknięć w jednej długiej wytwórni.
 */
interface ZapisStanu {
  obowiazujaca: SessionConfigEffective | null;
  zapisane: SessionConfigArea[];
  nieobsluzone: SessionConfigFieldCapability[];
  punkt: PunktWidzenia;
  faza: FazaOdczytu;
  powod: string;
}

export function utworzStanObszarowSesji(
  kanal: Kanal,
  zrodlo: ZrodloObszarowSesji = utworzZrodloObszarowSesji(kanal),
): StanObszarowSesji {
  const zapis: ZapisStanu = {
    obowiazujaca: null,
    zapisane: [],
    nieobsluzone: [],
    punkt: adresPoczatkowy(),
    faza: 'spoczynek',
    powod: '',
  };

  const sluchacze: Array<() => void> = [];
  const oglos = (): void => {
    for (const sluchacz of [...sluchacze]) sluchacz();
  };

  return {
    faza: () => zapis.faza,
    powodNiepowodzenia: () => zapis.powod,
    obowiazujaca: () => zapis.obowiazujaca,
    pochodzenie: (obszar) =>
      zapis.obowiazujaca?.origins.find((wpis) => wpis.area === obszar) ?? null,
    zapisane: () => zapis.zapisane,
    nieobsluzone: () => zapis.nieobsluzone,
    punkt: () => zapis.punkt,

    ustawPunkt(nowy) {
      zapis.punkt = nowy;
      oglos();
    },

    odswiez: () => wczytaj(zapis, zrodlo, oglos),
    utrwal: (obszary) => utrwalObszary(zapis, zrodlo, obszary, oglos),
    utrwalZaczepy: (zaczepy) => utrwalZaczepy(zapis, zrodlo, zaczepy, oglos),
    naZmiane: (sluchacz) => void sluchacze.push(sluchacz),
  };
}

/** `config.effective.get` — rozstrzygnięcie obszarów konfiguracji sesji dla wskazanego punktu widzenia okna. */
async function wczytaj(
  zapis: ZapisStanu,
  zrodlo: ZrodloObszarowSesji,
  oglos: () => void,
): Promise<void> {
  zapis.faza = 'odczyt';
  zapis.powod = '';
  oglos();

  const wynik = await zrodlo.obowiazujaca(zapis.punkt);
  if (!wynik.udany || wynik.wynik === undefined) {
    // Wykaz zostaje pusty, a nie udawany: pusty obszar i nieudany odczyt to różne stany.
    zapis.obowiazujaca = null;
    zapis.powod = opisOdmowyBledu('Konfiguracja obowiązująca nie dotarła', wynik.blad);
    zapis.faza = 'blad';
    oglos();
    return;
  }

  zapis.obowiazujaca = wynik.wynik.effective;
  zapis.nieobsluzone = [...(wynik.wynik.effective.unsupportedFields ?? [])];
  zapis.faza = 'gotowe';
  oglos();
}

/** `config.session.set` — utrwalenie wskazanych obszarów sesji na poziomie wybranego punktu widzenia okna. */
async function utrwalObszary(
  zapis: ZapisStanu,
  zrodlo: ZrodloObszarowSesji,
  obszary: readonly SessionConfigArea[],
  oglos: () => void,
): Promise<Wynik<ConfigSessionSetResponse>> {
  const zrodlowa = zapis.obowiazujaca;
  if (zrodlowa === null) return { udany: false, blad: bladBezOdczytu() };

  const wynik = await zrodlo.zapiszObszary(zapis.punkt, obszary, zrodlowa.config);
  if (wynik.udany && wynik.wynik !== undefined) {
    zapis.zapisane = [...wynik.wynik.storedAreas];
    zapis.nieobsluzone = [...(wynik.wynik.unsupportedFields ?? zapis.nieobsluzone)];
  }
  oglos();
  return wynik;
}

/**
 * Zapis obszaru `hooks` treścią z panelu zaczepów. Klient podmienia wyłącznie
 * obszar, który Operator złożył; rdzeń zapisuje tylko obszary z wykazu
 * `areas`, więc pozostałe pola są kontekstem, nie zapisem.
 */
async function utrwalZaczepy(
  zapis: ZapisStanu,
  zrodlo: ZrodloObszarowSesji,
  zaczepy: SessionConfigHooks,
  oglos: () => void,
): Promise<Wynik<ConfigSessionSetResponse>> {
  const zrodlowa = zapis.obowiazujaca;
  if (zrodlowa === null) return { udany: false, blad: bladBezOdczytu() };

  const wynik = await zrodlo.zapiszObszary(zapis.punkt, [SessionConfigArea.Hooks], {
    ...zrodlowa.config,
    hooks: zaczepy,
  });
  if (wynik.udany && wynik.wynik !== undefined) {
    zapis.zapisane = [...wynik.wynik.storedAreas];
  }
  oglos();
  return wynik;
}

/**
 * Zapis bez wcześniejszego odczytu nie idzie do rdzenia. Bez treści
 * konfiguracji obowiązującej klient musiałby wysłać obszar pusty, co
 * w kontrakcie znaczy usunięcie zapisu, więc czynność odmawia z powodem,
 * zamiast zgadywać.
 */
function bladBezOdczytu(): ErrorInfo {
  return {
    code: ErrorCode.ValidationFailed,
    message:
      'Nie ma czego utrwalić: konfiguracja obowiązująca jeszcze nie dotarła z rdzenia. ' +
      'Zapis pustego obszaru usunąłby zapis z poziomu, a nie utrwalił wartości.',
    retryable: true,
  };
}
