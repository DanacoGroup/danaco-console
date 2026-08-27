import { ErrorCode } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import type { ZebraneWSesji } from './zebrane-w-sesji';
import type { ZrodloBrowser } from './zrodlo-browser';

/**
 * Zaciągnięcie wykazu źródeł i notatek okna z rdzenia: dwa odczyty i przełożenie
 * ich wyniku na zbiór zebranego. Bez elementów widoku, bo okna pytają o wynik,
 * a nie o sposób. Odmowa „okna nie ma" kończy rozpoznanie stanem nietkniętym,
 * nie błędem.
 */
export type StanZaciagniecia = 'gotowe' | 'nietkniete' | 'odczyt' | 'blad';

export interface WynikZaciagniecia {
  stan: StanZaciagniecia;
  /** Zdanie dla Operatora; puste, gdy wykaz zaciągnął się bez uwag. */
  powod: string;
}

export interface WykazyZebranego {
  /** Czyta oba wykazy okna i podmienia nimi zbiór zebranego. */
  zaciagnij(idOkna: string): Promise<WynikZaciagniecia>;
}

export function utworzWykazyZebranego(
  zrodlo: ZrodloBrowser,
  zebrane: ZebraneWSesji,
): WykazyZebranego {
  return {
    async zaciagnij(idOkna) {
      if (idOkna === '') {
        return { stan: 'nietkniete', powod: 'Wykaz nie ma o co pytać: okno nie jest ustalone.' };
      }

      const [wykazZrodel, wykazNotatek] = await Promise.all([
        zrodlo.wykazZrodel({ windowId: idOkna }),
        zrodlo.wykazNotatek({ windowId: idOkna }),
      ]);

      // Okno nieznane rdzeniowi zgłasza się przy obu wykazach naraz, więc mówi
      // się o nim raz.
      if (czyOknoNieznane(wykazZrodel.blad?.code) && czyOknoNieznane(wykazNotatek.blad?.code)) {
        zebrane.zastapZrodla([]);
        zebrane.zastapNotatki([]);
        return {
          stan: 'nietkniete',
          powod:
            'Rdzeń nie zna jeszcze tego okna przeglądania — nic w nim nie zebrano. ' +
            'Wykaz zapełni się po pierwszym przejściu, źródle albo notatce.',
        };
      }

      if (!wykazZrodel.udany || wykazZrodel.wynik === undefined) {
        return {
          stan: 'blad',
          powod: opisOdmowy('Odczyt wykazu źródeł', wykazZrodel.blad?.code, wykazZrodel.blad?.message),
        };
      }
      if (!wykazNotatek.udany || wykazNotatek.wynik === undefined) {
        return {
          stan: 'blad',
          powod: opisOdmowy(
            'Odczyt wykazu notatek',
            wykazNotatek.blad?.code,
            wykazNotatek.blad?.message,
          ),
        };
      }

      zebrane.zastapZrodla(wykazZrodel.wynik.sources);
      zebrane.zastapNotatki(wykazNotatek.wynik.notes);
      return { stan: 'gotowe', powod: '' };
    },
  };
}

/**
 * Czy odmowa mówi „takiego okna przeglądania nie ma", a nie o usterce. Rozstrzyga
 * wyłącznie kod odmowy, więc treść komunikatu rdzenia nie wpływa na rozpoznanie.
 */
function czyOknoNieznane(kod: string | undefined): boolean {
  return kod === ErrorCode.NotFound;
}
