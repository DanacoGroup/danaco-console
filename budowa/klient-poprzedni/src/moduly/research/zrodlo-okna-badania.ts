import {
  Command,
  type Window,
  type WindowActionResponse,
  type WindowStateGetResponse,
} from '../../../../shared/contract';
import type { Kanal } from '../../protokol/kanal';
import { czyObiekt, czyTablica, czyTekst, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';
import { KOD_MODULU } from './kody-okien';
import { utworzNasluchOdmow, type NasluchOdmow, type OdpowiedzBadania } from './nasluch-odmow';

/**
 * Trzy komendy obszaru window w rdzeniu, na których stoi moduł Research: wykaz okien, stan
 * okna i akcja panelu, wszystkie wymagające identyfikatora okna badania.
 */
export interface ZrodloOknaBadania {
  /** Okna sesji należące do modułu Research, w kolejności rdzenia. */
  oknaBadania(idSesji: string): Promise<OdpowiedzBadania<Window[]>>;
  /** Stan okna badania wraz z konfiguracją efektywną. */
  stanOkna(idOkna: string): Promise<OdpowiedzBadania<WindowStateGetResponse>>;
  /** Akcja panelu akcji okna operacyjnego. */
  akcja(
    idOkna: string,
    idAkcji: string,
    parametry?: unknown,
  ): Promise<OdpowiedzBadania<WindowActionResponse>>;
  rozlacz(): void;
}

export function utworzZrodloOknaBadania(kanal: Kanal): ZrodloOknaBadania {
  const nasluch: NasluchOdmow = utworzNasluchOdmow(kanal);

  return {
    async oknaBadania(idSesji) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.WindowList, { sessionId: idSesji }),
        Command.WindowList,
        (tresc) => czyTablica(tresc.windows),
      );
      if (!wynik.udany || wynik.wynik === undefined) {
        return { udany: false, ...(wynik.blad === undefined ? {} : { blad: wynik.blad }) };
      }
      return {
        udany: true,
        wynik: wynik.wynik.windows.filter((okno) => okno.moduleId === KOD_MODULU),
      };
    },

    async stanOkna(idOkna) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.WindowStateGet, { windowId: idOkna, includeConfig: true }),
        Command.WindowStateGet,
        (tresc) => czyObiekt(tresc.window),
      );
    },

    async akcja(idOkna, idAkcji, parametry) {
      return sprawdzKsztaltAkcji(
        await nasluch.wyslij(Command.WindowAction, {
          windowId: idOkna,
          actionId: idAkcji,
          ...(parametry === undefined ? {} : { parameters: parametry }),
        }),
      );
    },

    rozlacz: () => nasluch.rozlacz(),
  };
}

/** Sprawdzian odpowiedzi akcji okna zachowujący nazwę nieznanego typu, żeby okno mogło wypisać ją Operatorowi wprost. */
function sprawdzKsztaltAkcji(
  odpowiedz: OdpowiedzBadania<WindowActionResponse>,
): OdpowiedzBadania<WindowActionResponse> {
  const wynik = sprawdzKsztalt(odpowiedz, Command.WindowAction, (tresc) =>
    czyTekst(tresc.actionId),
  );
  if (wynik.udany || odpowiedz.nieznanyTyp === undefined) return wynik;
  return { ...wynik, nieznanyTyp: odpowiedz.nieznanyTyp };
}
