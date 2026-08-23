import type { Module } from '../../../shared/contract';
import { czyNazwaIkony, type NazwaIkony } from '../ikony/ikony';
import type { PozycjaKomponentu } from './pozycje-komponentow';

/**
 * Przekład wykazu `module.list` na kafle dwóch stref strony głównej.
 *
 * Obie odpowiedzi biorą się z tego samego odczytu i z dwóch różnych kolumn
 * rdzenia, więc stoją w jednym pliku — inaczej rozstrzygnięcie, który moduł
 * gdzie stoi, rozjechałoby się między dwoma miejscami:
 *   · `pozycjeModulowStrefyDrugiej` — kolumna
 *     `modul.konfigurowany_na_stronie_glownej`,
 *   · `pozycjeModulowBezNawigacji`  — pusty `Module.environmentCodes`.
 *
 * Kafel modułu bez nawigacji jest jedyną drogą do modułu, którego rdzeń nie
 * pokazuje w bocznej nawigacji żadnego środowiska: `module.list` oddaje takie
 * moduły z pustym wykazem środowisk, a bez kafla ich okna operacyjne byłyby
 * nieosiągalne.
 *
 * Kryterium pochodzi z kontraktu, nie z wykazu nazw: `Module.environmentCodes`
 * jest opisane wprost jako „puste znaczy modul bez okna modulowego, dostepny
 * wylacznie ze strony glownej". Plik pyta dokładnie o to i nie wpisuje ani
 * jednego kodu modułu, więc gdy rdzeń przypnie moduł do środowiska, kafel
 * znika sam.
 *
 * To nie jest druga macierz widoczności: macierz mówi, gdzie moduł stoi, a ten
 * wykaz wyłącznie to, że nie stoi nigdzie.
 */

/** Pozycja kafla modułu bez pozycji w nawigacji. */
export interface PozycjaModuluStrony {
  /** Kod modułu z rdzenia (`modul.kod`) — po nim powłoka otwiera moduł. */
  kod: string;
  /** Nazwa własna modułu, tak jak podaje ją rdzeń. */
  nazwa: string;
  /** Wezwanie do działania na kaflu. */
  wezwanie: string;
  /** Opis przeznaczenia z rdzenia; pusty, gdy rdzeń go nie podał. */
  opis: string;
  /** Ikona kafla — kontrakt nie niesie rysunku, więc dokłada ją klient. */
  ikona: NazwaIkony;
  /** Ile okien operacyjnych moduł otwiera wraz z sobą (katalog rdzenia). */
  okien: number;
}

/** Ikona kafla, gdy kod modułu jest rdzeniowi znany, a klientowi nie. */
const IKONA_ZASTEPCZA: NazwaIkony = 'karta-okna';

/**
 * Kod modułu → ikona zestawu, gdy rdzeń ikony nie podał.
 *
 * Źródłem pierwszym jest `Module.icon` (kolumna `modul.ikona`) — patrz
 * `ikonaModulu`. Ta mapa wchodzi wyłącznie przy pustej kolumnie; kod spoza niej
 * nie jest błędem i dostaje ikonę zastępczą.
 */
const IKONY: Readonly<Record<string, NazwaIkony>> = {
  automations: 'automatyzacja',
  terminal: 'terminal',
};

/** Czy moduł nie stoi w bocznej nawigacji ani jednego środowiska. */
function bezNawigacji(modul: Module): boolean {
  return (modul.environmentCodes ?? []).length === 0;
}

/**
 * Czy moduł nastawia się w strefie kafli komponentów strony głównej.
 *
 * Jedno pytanie do rdzenia i ani jednego kodu modułu po tej stronie:
 * `Module.configuredOnHome` odpowiada kolumnie
 * `modul.konfigurowany_na_stronie_glownej`.
 */
function nastawianyNaStronie(modul: Module): boolean {
  return modul.configuredOnHome;
}

/**
 * Ikona kafla z rdzenia, gdy rdzeń ją podał i klient taki rysunek zna.
 *
 * Kolumna `modul.ikona` niesie nazwę z tego samego zestawu, którym posługuje się
 * klient, więc plik sprawdza przynależność zamiast trzymać drugą mapę kod
 * modułu → rysunek. Nazwa nieznana nie gasi kafla — dostaje rysunek zastępczy.
 */
function ikonaModulu(modul: Module): NazwaIkony {
  const nazwa = modul.icon ?? '';
  if (nazwa !== '' && czyNazwaIkony(nazwa)) return nazwa;
  return IKONY[modul.code] ?? IKONA_ZASTEPCZA;
}

/**
 * Przekład wykazu modułów rdzenia na kafle strefy modułów poza nawigacją.
 *
 * Kolejność jest kolejnością, w jakiej oddał je rdzeń — klient nie sortuje po
 * swojemu, bo `modul` niesie własną kolumnę kolejności.
 *
 * Moduł nastawiany na stronie głównej tu nie wraca. Moduł potrafi spełniać oba
 * warunki naraz — nie stać w żadnej bocznej nawigacji i zarazem nastawiać się
 * na stronie głównej — a wtedy rysowałby się na jednym ekranie dwa razy,
 * w dwóch sąsiednich strefach, choć oba kafle wołałyby ten sam kod. Pierwszeństwo
 * ma strefa kafli komponentów. Warunek pyta o dane rdzenia, nie o kod modułu,
 * więc rozstrzygnięcie przenosi się samo wraz ze zmianą oznaczenia w bazie.
 */
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

/**
 * Przekład wykazu modułów rdzenia na kafle strefy kafli komponentów.
 *
 * Wybór kafli idzie z kolumny `modul.konfigurowany_na_stronie_glownej`, a nie
 * z wyliczenia `ComponentKind`: oznaczenie kolejnego modułu w bazie ma dołożyć
 * kafel bez zmiany w kodzie klienta.
 *
 * Nazwa, opis i ikona idą z rdzenia. Wezwanie kafla to `modul.opis` — zdanie,
 * którym rdzeń sam opisuje moduł; pusty opis wraca do wezwania zastanego z tej
 * samej klatki, co `POZYCJE_KOMPONENTOW`, tak samo jak w kaflu personalizowanym
 * (`pozycje-personalizowane.ts`).
 *
 * Wezwanie opisuje moduł, a nie obiecuje czynności, której kafel nie wykonuje:
 * jedynym jego skutkiem jest przejście do modułu. Kafel zakładający komponent
 * stoi osobno, w przyborniku listwy ustawień, i woła `component.create`.
 */
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
