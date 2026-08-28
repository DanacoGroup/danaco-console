import { Command, WindowStatus, type Window } from '../../../../shared/contract';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import type { Kanal } from '../../protokol/kanal';
import { czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';
import { KOD_MODULU } from './wynik-czastkowy';

/**
 * Odnajduje okno rozmowy modułu Workspace w karcie sesji jako nośnik
 * przeniesienia kontekstu, bez zgadywania w przypadku odmowy rdzenia.
 */
export interface OknoRozmowyModulu {
  /** Identyfikator okna; pusty znaczy „nie odnaleziono” i wtedy `powod` mówi, czemu. */
  okno: string;
  /** Powód braku okna, zdaniem gotowym dla Operatora; pusty, gdy okno jest. */
  powod: string;
}

/** Powód sprzed pierwszego szukania okna — moduł jeszcze nie pytał rdzenia o okna należące do karty sesji. */
export const POWOD_PRZED_WCZYTANIEM =
  'Moduł nie zna jeszcze okna rozmowy — powłoka nie wczytała go do żadnej karty sesji.';

export async function odnajdzOknoRozmowy(
  kanal: Kanal,
  idSesji: string,
): Promise<OknoRozmowyModulu> {
  if (idSesji === '') {
    return { okno: '', powod: POWOD_PRZED_WCZYTANIEM };
  }

  const wynik = sprawdzKsztalt(
    await wywolaj(kanal, Command.WindowList, { sessionId: idSesji }),
    Command.WindowList,
    (tresc) => czyTablica(tresc.windows),
  );
  if (!wynik.udany || wynik.wynik === undefined) {
    return {
      okno: '',
      powod: opisOdmowyBledu('Odczyt okien karty sesji nie udał się', wynik.blad),
    };
  }

  const okna = wynik.wynik.windows;
  const moje = okna.filter((okno: Window) => okno.moduleId === KOD_MODULU);
  // Okno otwarte ma pierwszeństwo — przeniesienie do zamkniętego nie ma dokąd dojść.
  const otwarte = moje.find((okno: Window) => okno.status === WindowStatus.Open);
  if (otwarte !== undefined) return { okno: otwarte.id, powod: '' };

  if (moje.length > 0) {
    return {
      okno: '',
      powod:
        `Rdzeń oddał ${moje.length} okno/okien modułu ${KOD_MODULU} w tej karcie sesji, ` +
        'ale wszystkie zamknięte — przeniesienie kontekstu nie ma źródła.',
    };
  }
  return {
    okno: '',
    powod:
      okna.length === 0
        ? 'Ta karta sesji nie ma jeszcze ani jednego okna — okno rozmowy modułu powstanie razem z pierwszym z nich.'
        : `Rdzeń oddał ${okna.length} okno/okien tej karty sesji, ale żadne nie należy do modułu ${KOD_MODULU}.`,
  };
}
