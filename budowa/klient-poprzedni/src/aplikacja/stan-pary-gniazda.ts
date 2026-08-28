import { EventType } from '../../../shared/contract';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import { czyRolaWPetli, stanZKolejki, type GniazdoOkna } from '../okna-rownolegle/indeks';
import type { Kanal } from '../protokol/kanal';

/**
 * Doprowadza figurę pętli koordynator–wykonawca do nagłówka gniazda sceny.
 * Stan pary ma dwa człony: skład pary oraz jej położenie w kolejce, a oba
 * przychodzą z rdzenia zdarzeniami `window.changed` i `queue.changed`.
 * Okno poza pętlą pary nie ma.
 */
export function zwiazStanPary(kanal: Kanal, gniazdo: GniazdoOkna, idOkna: string): Odsubskrybuj {
  if (!czyRolaWPetli(gniazdo.rola())) {
    gniazdo.ustawStanPary(null);
  }

  const odRoli = kanal.naZdarzenie(EventType.WindowChanged, (tresc) => {
    if (tresc.window.id !== idOkna) return;
    gniazdo.ustawRole(tresc.window.windowRole);
    if (!czyRolaWPetli(gniazdo.rola())) gniazdo.ustawStanPary(null);
  });

  const odKolejki = kanal.naZdarzenie(EventType.QueueChanged, (tresc) => {
    if (!czyRolaWPetli(gniazdo.rola())) {
      gniazdo.ustawStanPary(null);
      return;
    }
    const okna = tresc.queue.windowIds;
    if (okna !== undefined && okna.length > 0 && !okna.includes(idOkna)) return;
    gniazdo.ustawStanPary(stanZKolejki(tresc.queue.status));
  });

  return () => {
    odRoli();
    odKolejki();
  };
}
