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

/**
 * Cztery komendy obiektu osadzonego w dokumencie — obraz, kształt, ikona, pole
 * tekstowe, logo.
 *
 * ── Czego to źródło NIE dubluje ─────────────────────────────────────────────
 * Rachunku kształtu ani biblioteki ikon. Kształt jest zapleczem modułu Design
 * (`design.vector.shape.add` zakłada go od razu jako węzły ścieżki), a ikona
 * jego biblioteką (`design.icon.library.search`). `studio.object.insert`
 * przyjmuje `designNodeId` i `iconName` właśnie dlatego: Studio wskazuje, co
 * osadzić, a rysuje to Design. Drugiego rachunku kształtu tu nie ma i nie
 * będzie.
 *
 * ── Wykresu rdzeń nie ma i okno tego nie udaje ──────────────────────────────
 * `StudioObjectKind` niesie wartość „wykres", ale rdzeń odmawia jej nazwanym
 * powodem: rachunku wykresu po stronie Studia nie ma. Odmowa stoi w oknie PRZED
 * próbą (`obiekt-panel.ts`) i kieruje do modułu Design, bo pokazanie kontrolki
 * kończącej się odmową rdzenia byłoby obietnicą bez pokrycia.
 */
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
        // Pole `removed` sprawdzamy wprost: odpowiedź „nie usunąłem" jest
        // wynikiem, nie awarią, i okno ma ją przeczytać jako odmowę nazwaną,
        // a nie zameldować usunięcie, którego nie było.
        (tresc) => czyLogiczna(tresc.removed) && czyObiekt(tresc.balance),
      );
    },
  };
}
