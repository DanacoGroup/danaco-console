/**
 * Plik niesie kontrakt strony dystrybucji modułu Apps: pięć komend rodziny
 * extension oraz jej zdarzenie, żadna z nich nie wymaga okna modułu.
 */

import {
  Command,
  EventType,
  type Extension,
  type ExtensionChangedEvent,
  type ExtensionConfigureRequest,
  type ExtensionInstallRequest,
  type ExtensionKind,
  type ExtensionListRequest,
  type ExtensionOrigin,
  type ExtensionListResponse,
  type ExtensionUninstallResponse,
} from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyLogiczna, czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { utworzWywolanieApps, type WywolanieApps } from './odmowa-rdzenia';

export interface ZawezenieKatalogu {
  /** Rodzaj zawężający wykaz; pominięty zwraca wszystkie rodzaje. */
  rodzaj?: ExtensionKind;
  /** Czy zwrócić wyłącznie pozycje zainstalowane. */
  tylkoZainstalowane?: boolean;
}

/**
 * Zlecenie instalacji pozycji katalogu w pełnym kształcie żądania kontraktu;
 * pole pominięte znaczy brak wskazania i zostawia rozstrzygnięcie rdzeniowi.
 */
export interface ZlecenieInstalacji {
  kod: string;
  rodzaj: ExtensionKind;
  /** Źródło paczki albo adres serwera; pusty łańcuch znaczy „bez wskazania”. */
  zrodlo: string;
  /** Punkt dostępu obsługujący rozszerzenie; pusty łańcuch znaczy „bez wskazania”. */
  idPunktuDostepu: string;
  /** Pochodzenie pozycji; pominięte zostawia rozstrzygnięcie rdzeniowi. */
  pochodzenie?: ExtensionOrigin;
  /** Konfiguracja początkowa; pominięta znaczy „bez konfiguracji”. */
  konfiguracja?: unknown;
}

export interface ZrodloRozszerzenApps {
  /** `extension.list` — katalog w kolejności wyświetlania. */
  katalog(zawezenie: ZawezenieKatalogu): Promise<Wynik<ExtensionListResponse>>;
  /** `extension.install` — instalacja pozycji wraz z jej pochodzeniem. */
  zainstaluj(z: ZlecenieInstalacji): Promise<Wynik<{ extension: Extension }>>;
  /** `extension.configure` — zapis konfiguracji rozszerzenia. */
  skonfiguruj(
    idRozszerzenia: string,
    konfiguracja: unknown,
    idPunktuDostepu: string,
  ): Promise<Wynik<{ extension: Extension }>>;
  /** `extension.toggle` — włączenie albo wyłączenie bez odinstalowania. */
  przelacz(idRozszerzenia: string, wlaczone: boolean): Promise<Wynik<{ extension: Extension }>>;
  /** `extension.uninstall` — odinstalowanie pozycji katalogu. */
  odinstaluj(idRozszerzenia: string): Promise<Wynik<ExtensionUninstallResponse>>;
  /** Subskrypcja `extension.changed` — zmiana katalogu rozszerzeń. */
  naZmianeKatalogu(sluchacz: (tresc: ExtensionChangedEvent) => void): Odsubskrybuj;
}

export function utworzZrodloRozszerzenApps(kanal: Kanal): ZrodloRozszerzenApps {
  const wywolaj: WywolanieApps = utworzWywolanieApps(kanal);

  return {
    async katalog(zawezenie) {
      const zadanie: ExtensionListRequest = {};
      if (zawezenie.rodzaj !== undefined) zadanie.kind = zawezenie.rodzaj;
      if (zawezenie.tylkoZainstalowane === true) zadanie.installedOnly = true;
      return sprawdzKsztalt(
        await wywolaj(Command.ExtensionList, zadanie),
        Command.ExtensionList,
        // Sprawdzany jest rodzaj tablicy, nie jej długość: katalog pusty jest poprawny.
        (tresc) => czyTablica(tresc.extensions),
      );
    },

    async zainstaluj(z) {
      const zadanie: ExtensionInstallRequest = { code: z.kod, kind: z.rodzaj };
      if (z.zrodlo !== '') zadanie.source = z.zrodlo;
      if (z.idPunktuDostepu !== '') zadanie.accessPointId = z.idPunktuDostepu;
      if (z.pochodzenie !== undefined) zadanie.origin = z.pochodzenie;
      if (z.konfiguracja !== undefined) zadanie.config = z.konfiguracja;
      return sprawdzKsztalt(
        await wywolaj(Command.ExtensionInstall, zadanie),
        Command.ExtensionInstall,
        (tresc) => czyObiekt(tresc.extension),
      );
    },

    async skonfiguruj(idRozszerzenia, konfiguracja, idPunktuDostepu) {
      const zadanie: ExtensionConfigureRequest = {
        extensionId: idRozszerzenia,
        config: konfiguracja,
      };
      if (idPunktuDostepu !== '') zadanie.accessPointId = idPunktuDostepu;
      return sprawdzKsztalt(
        await wywolaj(Command.ExtensionConfigure, zadanie),
        Command.ExtensionConfigure,
        (tresc) => czyObiekt(tresc.extension),
      );
    },

    async przelacz(idRozszerzenia, wlaczone) {
      return sprawdzKsztalt(
        await wywolaj(Command.ExtensionToggle, {
          extensionId: idRozszerzenia,
          enabled: wlaczone,
        }),
        Command.ExtensionToggle,
        (tresc) => czyObiekt(tresc.extension),
      );
    },

    async odinstaluj(idRozszerzenia) {
      return sprawdzKsztalt(
        await wywolaj(Command.ExtensionUninstall, { extensionId: idRozszerzenia }),
        Command.ExtensionUninstall,
        // Sprawdzany jest rodzaj pola uninstalled, nie jego prawdziwość.
        (tresc) => czyLogiczna(tresc.uninstalled),
      );
    },

    naZmianeKatalogu(sluchacz) {
      return kanal.naZdarzenie(EventType.ExtensionChanged, (tresc) => sluchacz(tresc));
    },
  };
}
