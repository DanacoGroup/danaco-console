import type { BrowserSource } from '../../../../shared/contract';
import { pokazKomunikat } from '../../aplikacja/komunikaty';
import { opisOdmowy } from '../../komponenty/odmowa';
import { POZYCJE_BEZ_OBSLUGI } from './etykiety-browser';
import { skutekPobrania, skutekPrzekazania, skutekZapisuZrodla } from './skutek-zapisu';
import type { StanPrzegladania } from './stan-przegladania';

/**
 * Panel akcji Sources Panel po stronie rdzenia: dodanie źródła, otwarcie,
 * podgląd migawki, oznaczenie kluczowym, usunięcie i przekazanie do Research.
 *
 * Jedna odpowiedzialność: rozmowa z rdzeniem w imieniu panelu źródeł. Panel
 * składa kontrolki i wykaz; tutaj mieszka to, co dzieje się po naciśnięciu.
 *
 * Usunięcie nie usuwa: `browser.source.remove` kontrakt niesie, ale ta czynność
 * jeszcze jej nie wywołuje. Powód bierze się z odczytu wykazu komend rdzenia,
 * a pozycja zostaje w wykazie — zniknięcie wiersza bez zapisu w rdzeniu byłoby
 * udawaniem wykonania.
 *
 * Adres w zdaniu potwierdzenia bierze się z odpowiedzi, nie z pola formularza
 * (`skutek-zapisu.ts`). Do rdzenia idzie adres przycięty, więc zdanie o skutku
 * musi mówić o tym, co wróciło, a nie o surowej treści pola.
 */
export interface CzynnosciZrodel {
  /** `browser.source.add`; zwraca `true`, gdy rdzeń przyjął źródło. */
  dodaj(url: string, tytul: string, kluczowe: boolean): Promise<boolean>;
  /** `browser.navigate` na adres źródła. */
  otworz(zrodlo: BrowserSource): Promise<void>;
  /** `context.transfer` zaznaczonych źródeł do modułu Research. */
  przekaz(identyfikatory: readonly string[]): Promise<void>;
  /** Migawka zapisana ze źródłem — bez komendy odczytu po identyfikatorze. */
  pokazMigawke(zrodlo: BrowserSource): void;
  /** Nazywa brak komendy usunięcia; wykaz zostaje nietknięty. */
  usun(zrodlo: BrowserSource): void;
}

export function utworzCzynnosciZrodel(
  stan: StanPrzegladania,
  powiedz: (tresc: string, powodzenie: boolean) => void,
): CzynnosciZrodel {
  return {
    async dodaj(url, tytul, kluczowe) {
      const idOkna = stan.idOkna();
      if (idOkna === '') {
        powiedz(stan.powod(), false);
        return false;
      }
      // Przycięcie robi się raz i dalej idzie już tylko wartość przycięta —
      // do rdzenia i do zdania o skutku trafia dokładnie ten sam adres.
      const adres = url.trim();
      const nazwa = tytul.trim();
      if (adres === '') {
        powiedz('Wskaż adres źródła — rdzeń odmówi dodania bez niego.', false);
        return false;
      }
      powiedz(`Dodawanie źródła ${adres}…`, true);
      const migawka = stan.migawka();
      const wynik = await stan.zrodlo.dodajZrodlo({
        windowId: idOkna,
        url: adres,
        ...(nazwa === '' ? {} : { title: nazwa }),
        ...(migawka === null ? {} : { snapshotId: migawka.id }),
        key: kluczowe,
      });
      if (!wynik.udany || wynik.wynik === undefined) {
        powiedz(opisOdmowy('Dodanie źródła', wynik.blad?.code, wynik.blad?.message), false);
        return false;
      }
      stan.zebrane.dopiszZrodlo(wynik.wynik.source);
      const skutek = skutekZapisuZrodla(wynik.wynik.source, adres);
      powiedz(skutek.zdanie, skutek.udany);
      // Wykaz odświeża się także przy rozbieżności: wiersz w rdzeniu istnieje,
      // więc obok zdania o rozjeździe ma stanąć jego prawdziwa postać.
      return true;
    },

    async otworz(zrodlo) {
      const idOkna = stan.idOkna();
      if (idOkna === '') {
        powiedz(stan.powod(), false);
        return;
      }
      powiedz(`Otwieranie źródła ${zrodlo.url}…`, true);
      const wynik = await stan.zrodlo.przejdz({ windowId: idOkna, url: zrodlo.url });
      if (!wynik.udany || wynik.wynik === undefined) {
        powiedz(opisOdmowy('Otwarcie źródła', wynik.blad?.code, wynik.blad?.message), false);
        return;
      }
      stan.wchlonMigawke(wynik.wynik.snapshot);
      const skutek = skutekPobrania(wynik.wynik.snapshot, zrodlo.url, 'źródło');
      powiedz(skutek.zdanie, skutek.udany);
    },

    async przekaz(identyfikatory) {
      const idOkna = stan.idOkna();
      if (idOkna === '' || identyfikatory.length === 0) {
        powiedz('Zaznacz źródła do przekazania — przekazanie idzie z wykazem, nie puste.', false);
        return;
      }
      powiedz('Przekazanie źródeł do modułu Research…', true);
      const wynik = await stan.zapisy.przekaz(idOkna, 'research', {
        knowledgeSourceIds: [...identyfikatory],
      });
      if (!wynik.udany || wynik.wynik === undefined) {
        powiedz(opisOdmowy('Przekazanie do Research', wynik.blad?.code, wynik.blad?.message), false);
        return;
      }
      // Liczba mówi o tym, co wysłano — `ContextTransferResponse` zawartości
      // kompletu nie oddaje, więc zdanie nie udaje, że rdzeń ją potwierdził.
      const skutek = skutekPrzekazania(
        wynik.wynik,
        'research',
        `Wysłano ${identyfikatory.length} źródeł`,
      );
      powiedz(skutek.zdanie, skutek.udany);
    },

    pokazMigawke(zrodlo) {
      const identyfikator = (zrodlo.snapshotId ?? '').trim();
      powiedz(
        identyfikator === ''
          ? 'Źródło zapisano bez migawki — rdzeń nie ma czego pokazać.'
          : `Migawka źródła: ${identyfikator}. Odczyt treści migawki po identyfikatorze nie ma komendy w kontrakcie.`,
        identyfikator !== '',
      );
    },

    usun(zrodlo) {
      const powod = stan.pokrycie.zdanie(
        POZYCJE_BEZ_OBSLUGI.usuniecieZrodla.komenda,
        POZYCJE_BEZ_OBSLUGI.usuniecieZrodla.czynnosc,
      );
      pokazKomunikat({ tytul: 'Usunięcie źródła', tresc: powod, waga: 'ostrz' });
      powiedz(`${powod} Pozycja ${zrodlo.url} zostaje w wykazie.`, false);
    },
  };
}
