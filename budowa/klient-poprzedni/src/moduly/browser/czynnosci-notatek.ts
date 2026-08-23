import type { BrowserNote } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import { skutekPobrania, skutekPrzekazania, skutekZapisuNotatki } from './skutek-zapisu';
import type { StanPrzegladania } from './stan-przegladania';

/**
 * Panel akcji Notes Panel po stronie rdzenia: zapis notatki, otwarcie źródła
 * powiązanego i przekazanie notatek do Research albo Library.
 *
 * Jedna odpowiedzialność: rozmowa z rdzeniem w imieniu panelu notatek. Panel
 * składa formularz i wykaz; tutaj mieszka to, co dzieje się po naciśnięciu.
 *
 * Zdanie końcowe każdej czynności powstaje z odpowiedzi rdzenia
 * (`skutek-zapisu.ts`), nie z treści żądania: zapis notatki opisuje jej postać
 * po zapisie, a `context.transfer` oddaje okno docelowe wraz z jego modułem
 * i znacznikiem `transferred`.
 */
export interface TrescNotatki {
  tresc: string;
  idZrodla: string;
  cytat: string;
}

export interface CzynnosciNotatek {
  /** `browser.note.add`; zwraca `true`, gdy rdzeń przyjął notatkę. */
  dodaj(zapis: TrescNotatki): Promise<boolean>;
  /** `context.transfer` wykazu notatek do wskazanego modułu. */
  przekaz(kodModulu: string): Promise<void>;
  /** `browser.navigate` na adres źródła powiązanego z notatką. */
  otworzZrodlo(notatka: BrowserNote): Promise<void>;
}

export function utworzCzynnosciNotatek(
  stan: StanPrzegladania,
  powiedz: (tresc: string, powodzenie: boolean) => void,
): CzynnosciNotatek {
  return {
    async dodaj(zapis) {
      const idOkna = stan.idOkna();
      if (idOkna === '') {
        powiedz(stan.powod(), false);
        return false;
      }
      powiedz('Zapis notatki w rdzeniu…', true);
      const wynik = await stan.zrodlo.dodajNotatke({
        windowId: idOkna,
        content: zapis.tresc,
        ...(zapis.idZrodla === '' ? {} : { sourceId: zapis.idZrodla }),
        ...(zapis.cytat === '' ? {} : { quote: zapis.cytat }),
      });
      if (!wynik.udany || wynik.wynik === undefined) {
        powiedz(opisOdmowy('Zapis notatki', wynik.blad?.code, wynik.blad?.message), false);
        return false;
      }
      // Wykaz dopisujemy także wtedy, gdy postać zapisana rozjeżdża się
      // z wysłaną: notatka w rdzeniu jest i widok ma pokazać jej prawdziwą
      // postać obok zdania o rozjeździe.
      stan.zebrane.dopiszNotatke(wynik.wynik.note);
      const skutek = skutekZapisuNotatki(wynik.wynik.note, zapis);
      powiedz(skutek.zdanie, skutek.udany);
      return true;
    },

    async przekaz(kodModulu) {
      const idOkna = stan.idOkna();
      const notatki = stan.zebrane.notatki();
      if (idOkna === '' || notatki.length === 0) {
        powiedz('Nie ma czego przekazać — wykaz notatek jest pusty.', false);
        return;
      }
      powiedz(`Przekazanie notatek do modułu ${kodModulu}…`, true);
      const wynik = await stan.zapisy.przekaz(idOkna, kodModulu, {
        prompt: notatki.map((notatka) => notatka.content).join('\n\n'),
        knowledgeSourceIds: zrodlaNotatek(notatki),
      });
      if (!wynik.udany || wynik.wynik === undefined) {
        powiedz(opisOdmowy('Przekazanie notatek', wynik.blad?.code, wynik.blad?.message), false);
        return;
      }
      const skutek = skutekPrzekazania(
        wynik.wynik,
        kodModulu,
        `Wysłano ${notatki.length} notatek`,
      );
      powiedz(skutek.zdanie, skutek.udany);
    },

    async otworzZrodlo(notatka) {
      const wpis = stan.zebrane.zrodlo(notatka.sourceId ?? '');
      const idOkna = stan.idOkna();
      if (wpis === null || idOkna === '') {
        powiedz('Notatka nie wskazuje źródła znanego temu panelowi.', false);
        return;
      }
      powiedz(`Otwieranie źródła ${wpis.url}…`, true);
      const wynik = await stan.zrodlo.przejdz({ windowId: idOkna, url: wpis.url });
      if (!wynik.udany || wynik.wynik === undefined) {
        powiedz(opisOdmowy('Otwarcie źródła', wynik.blad?.code, wynik.blad?.message), false);
        return;
      }
      stan.wchlonMigawke(wynik.wynik.snapshot);
      const skutek = skutekPobrania(wynik.wynik.snapshot, wpis.url, 'źródło');
      powiedz(skutek.zdanie, skutek.udany);
    },
  };
}

/** Źródła powiązane z notatkami, bez powtórzeń — ładunek `knowledgeSourceIds`. */
function zrodlaNotatek(notatki: readonly BrowserNote[]): string[] {
  const zebrane = new Set<string>();
  for (const notatka of notatki) {
    if (notatka.sourceId !== undefined && notatka.sourceId !== '') zebrane.add(notatka.sourceId);
  }
  return [...zebrane];
}
