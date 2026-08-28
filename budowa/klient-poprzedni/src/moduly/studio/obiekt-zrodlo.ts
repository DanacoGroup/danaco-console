import {
  Command,
  type StudioObjectFormatSetRequest,
  type StudioObjectFormatSetResponse,
  type StudioObjectInsertRequest,
  type StudioObjectInsertResponse,
  type StudioObjectListRequest,
  type StudioObjectListResponse,
  type StudioObjectRemoveRequest,
  type StudioObjectRemoveResponse,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyLogiczna, czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolajUczciwie } from './odmowa-rdzenia';

/** Cztery komendy obiektu osadzonego w dokumencie: obraz, kształt, ikona, pole tekstowe i logo, bez dublowania rachunku kształtu ani biblioteki ikon. */
export interface ObiektZrodlo {
  wstaw(zadanie: StudioObjectInsertRequest): Promise<Wynik<StudioObjectInsertResponse>>;
  postac(zadanie: StudioObjectFormatSetRequest): Promise<Wynik<StudioObjectFormatSetResponse>>;
  wykaz(zadanie: StudioObjectListRequest): Promise<Wynik<StudioObjectListResponse>>;
  usun(zadanie: StudioObjectRemoveRequest): Promise<Wynik<StudioObjectRemoveResponse>>;
}

export function utworzObiektZrodlo(kanal: Kanal): ObiektZrodlo {
  return {
    async wstaw(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioObjectInsert, zadanie),
        Command.StudioObjectInsert,
        (tresc) => czyObiekt(tresc.object) && czyObiekt(tresc.balance),
      );
    },

    async postac(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioObjectFormatSet, zadanie),
        Command.StudioObjectFormatSet,
        (tresc) => czyObiekt(tresc.object) && czyObiekt(tresc.balance),
      );
    },

    async wykaz(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioObjectList, zadanie),
        Command.StudioObjectList,
        (tresc) => czyTablica(tresc.objects),
      );
    },

    async usun(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioObjectRemove, zadanie),
        Command.StudioObjectRemove,
        // Pole removed podlega sprawdzeniu wprost: nieusunięcie jest wynikiem, nie awarią, jako odmowa.
        (tresc) => czyLogiczna(tresc.removed) && czyObiekt(tresc.balance),
      );
    },
  };
}
