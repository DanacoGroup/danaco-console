import { opisOdmowy } from '../../komponenty/odmowa';
import { poleTekstowe, przycisk, type WierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import { OBJASNIENIA } from './etykiety-translate';
import { dopnijDymek } from './kontrolki-translate';
import type { ZrodloGlosariusza } from './zrodlo-glosariusza';

/**
 * Trzy czynności glosariusza wykonywane na całości, nie na jednym terminie:
 * ujednolicenie terminologii paneli, wczytanie glosariusza z pliku i zapisanie
 * go do pliku (TBX/CSV).
 *
 * Wydzielone z okna, bo formularz terminu opisuje jeden termin, a to są
 * czynności zbiorcze — jedna odpowiedzialność na plik. Ścieżka pliku jest
 * ścieżką po stronie rdzenia: klient plików nie czyta i nie zapisuje, więc
 * kontrolką jest pole tekstowe, a nie okno wyboru pliku przeglądarki, które
 * sugerowałoby przesył nieprzewidziany kontraktem.
 *
 * Trzy czynności stoją poza wytwórnią elementów: wytwórnia składa pole i pasek
 * przycisków, a każda czynność jest osobną funkcją modułu — mówi o czym innym
 * i daje się sprawdzić bez klikania w przycisk.
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

/** `translate.glossary.apply` na komplecie paneli — bez zawężania do języka. */
async function ujednolicPanele(
  zrodlo: ZrodloGlosariusza,
  odpowiedz: WierszOdpowiedzi,
): Promise<void> {
  odpowiedz.pokaz('Ujednolicanie terminologii we wszystkich panelach…', true);
  // Puste `panelId` obejmuje komplet paneli — okno zarządcy działa na
  // glosariuszu, więc nie zawęża czynności do jednego języka.
  const wynik = await zrodlo.ujednolic('');
  if (!wynik.udany || wynik.wynik === undefined) {
    odpowiedz.pokaz(
      opisOdmowy('Ujednolicenie terminologii', wynik.blad?.code, wynik.blad?.message),
      false,
    );
    return;
  }
  // Liczba jest tu świadkiem skutku, w odróżnieniu od eksportu niżej: rdzeń
  // liczy wystąpienia faktycznie podmienione i zapisane (`zastosujWPanelu`),
  // więc wywołanie powtórzone oddaje zero — nie ma już czego podmieniać.
  odpowiedz.pokaz(`Rdzeń ujednolicił ${wynik.wynik.changedCount} wystąpień.`, true);
}

/** `translate.glossary.import` — rdzeń odmawia, bo nie czyta dysku Operatora. */
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
    // Powód rdzenia idzie w całości i bez zdania o kanale modelu: ten sam kod
    // `channel_unavailable` niesie tu zupełnie inną odmowę niż przy czynnościach
    // modelowych (`odmowa-translate.ts` — dopisanie tam zdania o kanale byłoby
    // zdaniem nieprawdziwym).
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
 * `translate.glossary.export`.
 *
 * Liczba nie jest świadkiem zapisu, więc wynik nie jest powodzeniem.
 * `TranslateGlossaryExportResponse` niesie wyłącznie `exportedCount`: ani
 * ścieżki wyniku, ani znaku, że plik powstał — odpowiedź wygląda tak samo także
 * przy ścieżce do nieistniejącego katalogu i przy napisie, który ścieżką nie
 * jest. Liczba mówi o zawartości glosariusza w chwili zlecenia (rdzeń liczy
 * zastane terminy — `EksportujSlownik`), a nie o zapisie, dlatego wiersz
 * odpowiedzi ma wydźwięk odmowy.
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
