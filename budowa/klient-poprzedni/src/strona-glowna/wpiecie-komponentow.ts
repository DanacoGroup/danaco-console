import { Command } from '../../../shared/contract';
import type { Kanal } from '../protokol/kanal';
import { utworzMacierzModulow, type MacierzModulow } from './macierz-modulow';
import { POZYCJE_KOMPONENTOW } from './pozycje-komponentow';
import { pozycjeModulowStrefyDrugiej } from './pozycje-modulow';
import { pozycjeZKomponentow } from './pozycje-personalizowane';
import type { KodSrodowiska } from './pozycje-srodowisk';
import type { StrefaKomponentow } from './strefa-komponentow';

/**
 * Wpięcie macierzy widoczności w kafle strefy drugiej.
 *
 * Docelowe środowisko kafla wynika z macierzy `srodowisko_modul`, którą rdzeń
 * podaje w `environment.list`. Kopia tej macierzy po stronie klienta mogłaby
 * rozjechać się z bazą po cichu, więc jej tutaj nie ma.
 *
 * Zapytanie o środowiska pada tu osobno, z `includeModules: true`, bo
 * `wpiecie-srodowisk.ts` pyta bez tego pola — kartom strefy pierwszej moduły
 * nie są potrzebne. Podobnie `module.list` wołają dwa wpięcia: to po kafle
 * strefy drugiej, a `wpiecie-modulow.ts` po kafle modułów poza nawigacją.
 * Scalenie wołań wymaga zmiany w `aplikacja/widok-strony-glownej.ts`, poza tym
 * pakietem; strona wstaje raz na wejście, nie w pętli, więc dwa wołania są
 * ceną za niezależność obu wpięć.
 *
 * Kafel jest zawsze klikalny i zawsze prowadzi do pracy. Gdy macierz jeszcze
 * nie przyszła albo moduł nie stoi w żadnej bocznej nawigacji, przejście idzie
 * do środowiska początkowego — nadal ze wskazaniem modułu — zamiast pokazać
 * odmowę lub kafel wyszarzony.
 */

export interface ZaleznosciWpieciaKomponentow {
  /** Kanał kontraktu — źródło macierzy widoczności i wykazu komponentów. */
  kanal: Kanal;
  /** Strefa druga — przyjmuje kafle komponentów zbudowanych przez Operatora. */
  strefa: StrefaKomponentow;
  /** Przejście do pracy; drugi argument wskazuje pozycję bocznej nawigacji. */
  naPrzejscie(kod: KodSrodowiska, modul?: string): void;
}

export interface WpiecieKomponentow {
  /** Odczytuje wykaz komponentów z rdzenia na nowo. */
  odswiez(): void;
  /** Skutek wyboru kafla o wskazanym kodzie modułu. */
  wybierz(kodModulu: string): void;
}

/**
 * Wezwanie kafla, gdy rdzeń nie opisał modułu — z wykazu zastanego.
 *
 * Kod nieznany zastanej czwórce zwraca `undefined`, a przekład sięga wtedy po
 * własną wartość ostateczną.
 */
function wezwanieZastane(kod: string): string | undefined {
  return POZYCJE_KOMPONENTOW.find((pozycja) => pozycja.kod === kod)?.wezwanie;
}

export function wepnijKomponenty(
  zaleznosci: ZaleznosciWpieciaKomponentow,
): WpiecieKomponentow {
  const macierz: MacierzModulow = utworzMacierzModulow();

  zaleznosci.kanal.wyslij(Command.EnvironmentList, { includeModules: true }, (wynik) => {
    const srodowiska = wynik.wynik?.environments;
    if (!wynik.udany || srodowiska === undefined) return;
    macierz.ustawZeSrodowisk(srodowiska);
  });

  // Kafle strefy drugiej wynikają z kolumny
  // `modul.konfigurowany_na_stronie_glownej`, którą kontrakt wystawia jako
  // `Module.configuredOnHome` — nie z zamkniętego wyliczenia `ComponentKind`.
  // Dzięki temu oznaczenie kolejnego modułu w bazie zmienia ekran bez zmiany
  // kodu. Odmowa rdzenia zostawia czwórkę zastaną, więc strefa nie gaśnie.
  zaleznosci.kanal.wyslij(Command.ModuleList, {}, (wynik) => {
    const moduly = wynik.wynik?.modules;
    if (!wynik.udany || moduly === undefined) return;
    zaleznosci.strefa.ustawRodzaje(
      pozycjeModulowStrefyDrugiej(moduly, wezwanieZastane),
    );
  });

  // Kafle personalizowane: każdy komponent zbudowany przez Operatora dostaje
  // własny kafel pod nadaną mu nazwą. Wykaz przychodzi bez `includeDisabled`,
  // więc komponent wyłączony skraca listę zamiast dawać kafel wyszarzony.
  // Odmowa rdzenia zostawia same kafle rodzajów — te wynikają z kontraktu,
  // nie z tej odpowiedzi.
  function odczytajKomponenty(): void {
    zaleznosci.kanal.wyslij(Command.ComponentList, {}, (wynik) => {
      const komponenty = wynik.wynik?.components;
      if (!wynik.udany || komponenty === undefined) return;
      zaleznosci.strefa.ustawPersonalizowane(pozycjeZKomponentow(komponenty));
    });
  }

  odczytajKomponenty();

  return {
    // Ponowny odczyt po założeniu, nie doklejenie kafla z odpowiedzi.
    // `component.create` oddaje komponent, ale o wykazie rozstrzyga rdzeń:
    // kafel doklejony po stronie widoku pokazywałby stan, którego drugi
    // odczyt nie potwierdził.
    odswiez: odczytajKomponenty,

    wybierz(kodModulu) {
      // Wskazanie modułu idzie zawsze, także gdy macierz go nie zna. Moduł
      // niewidoczny w żadnej bocznej nawigacji nie ma pozycji na liście, ale
      // ma okna; powłoka otwiera go wtedy drogą bezpośrednią, z katalogu
      // `module.list` z pominięciem wykazu. Zgubienie wskazania odbierałoby
      // kaflowi jedyne wejście do tych okien.
      const srodowisko = macierz.srodowiskoModulu(kodModulu) ?? macierz.srodowiskoPoczatkowe();
      zaleznosci.naPrzejscie(srodowisko, kodModulu);
    },
  };
}
