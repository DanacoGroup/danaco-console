import {
  zadajArchiwizacjeSesji,
  zadajPrzywrocenieSesji,
  zadajWykazArchiwum,
} from '../protokol/archiwum-sesji';
import { zadajWznowienieSesji, zadajZatrzymanieSesji } from '../protokol/bieg-sesji';
import type { Kanal } from '../protokol/kanal';
import { zadajWyjecieZProjektu } from '../protokol/projekt-sesji';
import type { CzynnosciSesji, MeldunekCzynnosci } from './czynnosci-sesji';
import { nazwa, odmowa } from './meldunki-sesji';
import type { WynikArchiwum } from './strefa-archiwum';
import { czynnosciZNazwa } from './wykonanie-nazwy';

/**
 * Wykonanie operacji historii sesji — most między czynnością widoku a komendą
 * kontraktu.
 *
 * Składa komplet czynności i przeprowadza te, które idą do rdzenia bez pytania
 * Operatora o cokolwiek. Czynności pytające o nazwę stoją w `wykonanie-nazwy`;
 * dzieli je nie temat, tylko kształt — tam każda ma etap zbierania danych, tu
 * żadna go nie ma.
 *
 * Meldunek składany jest z wyniku, nie z żądania: komendy zbiorowe (`archive`,
 * `restore`, `project.clear`) oddają wykaz sesji faktycznie przeniesionych,
 * a ten bywa krótszy od żądania.
 *
 * Po każdej udanej zmianie woła się odświeżenie podane przez wpięcie: rdzeń
 * rozsyła `session.changed`, ale archiwum jest odpytywane osobno i tego
 * zdarzenia nie widzi.
 */

export interface ZaleznosciWykonania {
  kanal: Kanal;
  /** Ponowne odpytanie wykazów po udanej zmianie. */
  odswiez(): void;
}

export function utworzCzynnosciHistorii(zaleznosci: ZaleznosciWykonania): CzynnosciSesji {
  const { kanal, odswiez } = zaleznosci;

  return {
    ...czynnosciZNazwa(kanal, odswiez),

    async wyjmijZProjektu(wpis): Promise<MeldunekCzynnosci> {
      const wynik = await zadajWyjecieZProjektu(kanal, { sessionIds: [wpis.sesja.id] });
      if (!wynik.udany) return odmowa(wynik, 'wyjęcie sesji z projektu');
      odswiez();
      return (wynik.wynik?.clearedIds.length ?? 0) > 0
        ? `Sesja „${nazwa(wpis)}” nie należy już do projektu.`
        : 'Rdzeń nie wyjął żadnej sesji z projektu.';
    },

    async archiwizuj(wpis): Promise<MeldunekCzynnosci> {
      const wynik = await zadajArchiwizacjeSesji(kanal, { sessionIds: [wpis.sesja.id] });
      if (!wynik.udany) return odmowa(wynik, 'archiwizacja sesji');
      odswiez();
      return (wynik.wynik?.archivedIds.length ?? 0) > 0
        ? `Sesja „${nazwa(wpis)}” poszła do archiwum. Zapis został w całości.`
        : 'Rdzeń nie przeniósł żadnej sesji do archiwum.';
    },

    async przywroc(wpis): Promise<MeldunekCzynnosci> {
      const wynik = await zadajPrzywrocenieSesji(kanal, { sessionIds: [wpis.sesja.id] });
      if (!wynik.udany) return odmowa(wynik, 'przywrócenie sesji z archiwum');
      odswiez();
      return (wynik.wynik?.restoredIds.length ?? 0) > 0
        ? `Sesja „${nazwa(wpis)}” wróciła do historii bieżącej.`
        : 'Rdzeń nie przywrócił żadnej sesji.';
    },

    async wznow(wpis): Promise<MeldunekCzynnosci> {
      const wynik = await zadajWznowienieSesji(kanal, { sessionId: wpis.sesja.id });
      if (!wynik.udany) return odmowa(wynik, 'wznowienie sesji');
      odswiez();
      return `Sesja „${nazwa(wpis)}” wznowiona — otwarto ${wynik.wynik?.windows.length ?? 0} okien.`;
    },

    async zatrzymaj(wpis): Promise<MeldunekCzynnosci> {
      const wynik = await zadajZatrzymanieSesji(kanal, { sessionId: wpis.sesja.id });
      if (!wynik.udany) return odmowa(wynik, 'zatrzymanie tur sesji');
      odswiez();
      const zatrzymane = wynik.wynik?.stoppedWindowIds.length ?? 0;
      // Zero okien to wynik poprawny: nic nie biegło. Mówimy to wprost,
      // zamiast udawać, że coś przerwaliśmy.
      return zatrzymane > 0
        ? `Zatrzymano tury w ${zatrzymane} oknach sesji „${nazwa(wpis)}”.`
        : `W sesji „${nazwa(wpis)}” nie biegła żadna tura.`;
    },
  };
}

/** Odpytanie archiwum dla `strefa-archiwum` — odmowa wraca treścią, nie wyjątkiem. */
export function utworzOdpytanieArchiwum(kanal: Kanal): () => Promise<WynikArchiwum> {
  return async () => {
    const wynik = await zadajWykazArchiwum(kanal, {});
    if (!wynik.udany || wynik.wynik === undefined) {
      return {
        sesje: [],
        razem: 0,
        blad: wynik.blad?.message ?? 'Rdzeń nie oddał wykazu archiwum.',
      };
    }
    return { sesje: wynik.wynik.sessions, razem: wynik.wynik.total };
  };
}
