import { Command } from '../../../shared/contract';
import type { Kanal } from '../protokol/kanal';
import { utworzMacierzModulow, type MacierzModulow } from './macierz-modulow';
import { POZYCJE_KOMPONENTOW } from './pozycje-komponentow';
import { pozycjeModulowStrefyDrugiej } from './pozycje-modulow';
import { pozycjeZKomponentow } from './pozycje-personalizowane';
import type { KodSrodowiska } from './pozycje-srodowisk';
import type { StrefaKomponentow } from './strefa-komponentow';

/** Wpięcie macierzy widoczności w kafle strefy drugiej łączy odpowiedź rdzenia o środowiskach i modułach z kaflami, kierując przejście do środowiska początkowego, gdy macierz jeszcze nie przyszła. */
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

/** Wezwanie kafla, gdy rdzeń nie opisał modułu, sięga po wykaz zastany; kod nieznany czwórce zwraca brak wartości ostatecznej. */
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

  // Kafle strefy drugiej wynikają z kolumny widoczności na stronie głównej, nie z wyliczenia rodzajów.
  zaleznosci.kanal.wyslij(Command.ModuleList, {}, (wynik) => {
    const moduly = wynik.wynik?.modules;
    if (!wynik.udany || moduly === undefined) return;
    zaleznosci.strefa.ustawRodzaje(
      pozycjeModulowStrefyDrugiej(moduly, wezwanieZastane),
    );
  });

  // Kafle personalizowane: każdy komponent zbudowany przez operatora dostaje kafel pod nadaną nazwą.
  function odczytajKomponenty(): void {
    zaleznosci.kanal.wyslij(Command.ComponentList, {}, (wynik) => {
      const komponenty = wynik.wynik?.components;
      if (!wynik.udany || komponenty === undefined) return;
      zaleznosci.strefa.ustawPersonalizowane(pozycjeZKomponentow(komponenty));
    });
  }

  odczytajKomponenty();

  return {
    // Ponowny odczyt po założeniu, nie doklejenie kafla z odpowiedzi, bo o wykazie rozstrzyga rdzeń.
    odswiez: odczytajKomponenty,

    wybierz(kodModulu) {
      // Wskazanie modułu idzie zawsze, także gdy macierz go nie zna; moduł bez pozycji na liście ma okna.
      const srodowisko = macierz.srodowiskoModulu(kodModulu) ?? macierz.srodowiskoPoczatkowe();
      zaleznosci.naPrzejscie(srodowisko, kodModulu);
    },
  };
}
