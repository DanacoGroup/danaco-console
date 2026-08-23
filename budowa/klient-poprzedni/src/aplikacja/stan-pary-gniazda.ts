import { EventType } from '../../../shared/contract';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import { czyRolaWPetli, stanZKolejki, type GniazdoOkna } from '../okna-rownolegle/indeks';
import type { Kanal } from '../protokol/kanal';

/**
 * Doprowadza figurę pętli koordynator–wykonawca do nagłówka gniazda sceny.
 *
 * „Stan pary" ma dwa człony: skład pary (kto jest koordynatorem, kto wykonawcą)
 * oraz jej położenie (czy wykonawca pracuje, czy kolejka stoi). Oba przychodzą
 * z rdzenia zdarzeniami `window.changed` i `queue.changed`.
 *
 * Właścicielem pola `windowRole` jest rdzeń, nie scena: po odtworzeniu okna
 * z bazy, po `window.update` z innego urządzenia albo po `context.transfer`
 * może wrócić rola inna niż domyślna nadana przy tworzeniu gniazda. Dlatego
 * każde `window.changed` o tym oknie przestawia nagłówek; zmianę na tę samą
 * rolę gniazdo pomija, więc widok nie przebudowuje się bez powodu.
 *
 * Okno poza pętlą nie ma pary — jego nagłówek zostaje bez plakietki stanu
 * zamiast pokazywać stan cudzej kolejki.
 *
 * Każde gniazdo słucha całej szyny, więc zdarzenie dotyczące innego okna jest
 * pomijane. Kolejka bez wykazu okien dotyczy całej sesji i wchodzi do każdego
 * gniazda tej sceny.
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
