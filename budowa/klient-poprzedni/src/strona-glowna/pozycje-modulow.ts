import type { Module } from '../../../shared/contract';
import { czyNazwaIkony, type NazwaIkony } from '../ikony/ikony';
import type { PozycjaKomponentu } from './pozycje-komponentow';

/** Przekład wykazu modułów rdzenia na kafle dwóch stref strony głównej stoi w jednym pliku, bo obie odpowiedzi biorą się z tego samego odczytu i dwóch różnych kolumn rdzenia. */
export interface PozycjaModuluStrony {
  /** Kod modułu z rdzenia, po którym powłoka otwiera moduł. */
  kod: string;
  /** Nazwa własna modułu, tak jak podaje ją rdzeń. */
  nazwa: string;
  /** Wezwanie do działania na kaflu. */
  wezwanie: string;
  /** Opis przeznaczenia z rdzenia; pusty, gdy rdzeń go nie podał. */
  opis: string;
  /** Ikona kafla; kontrakt nie niesie rysunku, więc dokłada ją klient sam. */
  ikona: NazwaIkony;
  /** Ile okien operacyjnych moduł otwiera wraz z sobą (katalog rdzenia). */
  okien: number;
}

/** Ikona kafla, gdy kod modułu jest rdzeniowi znany, a klientowi jego rysunek nie jest znany z zestawu. */
const IKONA_ZASTEPCZA: NazwaIkony = 'karta-okna';

/** Kod modułu tłumaczy się na ikonę zestawu, gdy rdzeń nie podał własnej; mapa działa tylko przy pustej kolumnie ikony. */
const IKONY: Readonly<Record<string, NazwaIkony>> = {
  automations: 'automatyzacja',
  terminal: 'terminal',
};

/** Rozstrzyga, czy moduł nie stoi w bocznej nawigacji żadnego środowiska produktu, wracając wtedy wartość prawdy. */
function bezNawigacji(modul: Module): boolean {
  return (modul.environmentCodes ?? []).length === 0;
}

/** Rozstrzyga, czy moduł nastawia się w strefie kafli komponentów na stronie głównej, wyłącznie na podstawie odpowiedzi rdzenia. */
function nastawianyNaStronie(modul: Module): boolean {
  return modul.configuredOnHome;
}

/** Dobiera ikonę kafla z rdzenia, gdy rdzeń ją podał i klient zna taki rysunek w swoim całym zestawie ikon produktu. */
function ikonaModulu(modul: Module): NazwaIkony {
  const nazwa = modul.icon ?? '';
  if (nazwa !== '' && czyNazwaIkony(nazwa)) return nazwa;
  return IKONY[modul.code] ?? IKONA_ZASTEPCZA;
}

/** Przekłada wykaz modułów rdzenia na kafle strefy modułów poza nawigacją, zachowując kolejność oddaną przez rdzeń. */
export function pozycjeModulowBezNawigacji(
  moduly: readonly Module[],
): readonly PozycjaModuluStrony[] {
  return moduly
    .filter((modul) => bezNawigacji(modul) && !nastawianyNaStronie(modul))
    .map((modul) => ({
      kod: modul.code,
      nazwa: modul.name,
      wezwanie: 'Otwórz moduł',
      opis: modul.description ?? '',
      ikona: ikonaModulu(modul),
      okien: modul.operationalWindowCodes.length,
    }));
}

/** Przekłada wykaz modułów rdzenia na kafle strefy kafli komponentów, biorąc nazwę, opis i ikonę bezpośrednio z rdzenia. */
export function pozycjeModulowStrefyDrugiej(
  moduly: readonly Module[],
  wezwanieZastane: (kod: string) => string | undefined = () => undefined,
): readonly PozycjaKomponentu[] {
  return moduly.filter(nastawianyNaStronie).map((modul) => {
    const opis = modul.description ?? '';
    return {
      kod: modul.code,
      nazwa: modul.name,
      wezwanie: opis !== '' ? opis : (wezwanieZastane(modul.code) ?? 'Otwórz moduł'),
      ikona: ikonaModulu(modul),
    };
  });
}
