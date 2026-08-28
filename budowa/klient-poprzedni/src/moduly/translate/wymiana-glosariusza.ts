import { opisOdmowy } from '../../komponenty/odmowa';
import { poleTekstowe, przycisk, type WierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import { OBJASNIENIA } from './etykiety-translate';
import { dopnijDymek } from './kontrolki-translate';
import type { ZrodloGlosariusza } from './zrodlo-glosariusza';

/**
 * Trzy czynności glosariusza wykonywane na całości, nie na jednym terminie: ujednolicenie
 * terminologii paneli, wczytanie glosariusza z pliku i zapisanie go do pliku.
 */
export interface WymianaGlosariusza {
  element: HTMLElement;
}

export function utworzWymianeGlosariusza(
  zrodlo: ZrodloGlosariusza,
  odpowiedz: WierszOdpowiedzi,
): WymianaGlosariusza {
  const sciezka = poleTekstowe({
    etykieta: 'Ścieżka pliku glosariusza (TBX/CSV)',
    podpowiedz: '/dane/glosariusz.csv',
  });
  dopnijDymek(sciezka.element, OBJASNIENIA.sciezkaGlosariusza);

  const ujednolic = przycisk('Ujednolić panele', 'dn-btn dn-btn--sm dn-btn--zarys');
  const importuj = przycisk('Importuj z pliku', 'dn-btn dn-btn--sm dn-btn--zarys');
  const eksportuj = przycisk('Eksportuj do pliku', 'dn-btn dn-btn--sm dn-btn--zarys');

  const pasek = document.createElement('div');
  pasek.className = 'mt-pasek';
  pasek.append(ujednolic, importuj, eksportuj);

  const element = document.createElement('div');
  element.className = 'mt-wymiana';
  element.append(sciezka.element, pasek);

  const wskazana = (): string => sciezka.kontrolka.value.trim();
  ujednolic.addEventListener('click', () => void ujednolicPanele(zrodlo, odpowiedz));
  importuj.addEventListener('click', () => void wczytajZPliku(zrodlo, odpowiedz, wskazana()));
  eksportuj.addEventListener('click', () => void zapiszDoPliku(zrodlo, odpowiedz, wskazana()));

  return { element };
}

/** Komenda ujednolicenia terminologii działa na komplecie paneli, bez zawężania czynności do jednego języka panelu. */
async function ujednolicPanele(
  zrodlo: ZrodloGlosariusza,
  odpowiedz: WierszOdpowiedzi,
): Promise<void> {
  odpowiedz.pokaz('Ujednolicanie terminologii we wszystkich panelach…', true);
  // Puste wskazanie panelu obejmuje komplet paneli — okno zarządcy działa na glosariuszu, nie na języku.
  const wynik = await zrodlo.ujednolic('');
  if (!wynik.udany || wynik.wynik === undefined) {
    odpowiedz.pokaz(
      opisOdmowy('Ujednolicenie terminologii', wynik.blad?.code, wynik.blad?.message),
      false,
    );
    return;
  }
  // Liczba jest świadkiem skutku: rdzeń liczy wystąpienia podmienione, powtórzone wywołanie oddaje zero.
  odpowiedz.pokaz(`Rdzeń ujednolicił ${wynik.wynik.changedCount} wystąpień.`, true);
}

/** Komenda wczytania glosariusza z pliku odmawia po stronie rdzenia, bo rdzeń nie czyta dysku operatora bezpośrednio. */
async function wczytajZPliku(
  zrodlo: ZrodloGlosariusza,
  odpowiedz: WierszOdpowiedzi,
  wskazana: string,
): Promise<void> {
  if (wskazana === '') {
    odpowiedz.pokaz('Podaj ścieżkę pliku — rdzeń odmówi wczytania bez niej.', false);
    return;
  }
  odpowiedz.pokaz(`Wczytywanie glosariusza z ${wskazana}…`, true);
  const wynik = await zrodlo.wczytajZPliku(wskazana);
  if (!wynik.udany || wynik.wynik === undefined) {
    // Powód rdzenia idzie bez zdania o kanale: kod niesie tu inną odmowę niż przy czynnościach modelowych.
    odpowiedz.pokaz(opisOdmowy('Import glosariusza', wynik.blad?.code, wynik.blad?.message), false);
    return;
  }
  odpowiedz.pokaz(
    `Rdzeń wczytał ${wynik.wynik.importedCount} terminów. Kontrakt nie ma komendy odczytu ` +
      'glosariusza, więc wczytane terminy nie pojawią się w wykazie poniżej.',
    true,
  );
}

/**
 * Eksport glosariusza liczy zawartość w chwili zlecenia, a nie potwierdza zapisu; liczba nie jest
 * świadkiem zapisu pliku.
 */
async function zapiszDoPliku(
  zrodlo: ZrodloGlosariusza,
  odpowiedz: WierszOdpowiedzi,
  wskazana: string,
): Promise<void> {
  if (wskazana === '') {
    odpowiedz.pokaz('Podaj ścieżkę pliku wyniku — rdzeń odmówi zapisu bez niej.', false);
    return;
  }
  odpowiedz.pokaz(`Zapisywanie glosariusza do ${wskazana}…`, true);
  const wynik = await zrodlo.zapiszDoPliku(wskazana);
  if (!wynik.udany || wynik.wynik === undefined) {
    odpowiedz.pokaz(opisOdmowy('Eksport glosariusza', wynik.blad?.code, wynik.blad?.message), false);
    return;
  }
  odpowiedz.pokaz(
    `Eksport glosariusza: rdzeń przyjął zlecenie i naliczył ${wynik.wynik.exportedCount} ` +
      'terminów, ale NIE potwierdził zapisu. Odpowiedź kontraktu niesie samą liczbę terminów ' +
      `— bez ścieżki wyniku i bez znaku powstania pliku — a ścieżkę ${wskazana} podałeś Ty, ` +
      'nie rdzeń. Okno nie ma czym sprawdzić, czy plik pod nią jest, więc go nie potwierdza.',
    false,
  );
}
