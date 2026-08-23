import {
  Command,
  type StudioApparatusInsertRequest,
  type StudioApparatusInsertResponse,
  type StudioApparatusListRequest,
  type StudioApparatusListResponse,
  type StudioApparatusRefreshRequest,
  type StudioApparatusRefreshResponse,
  type StudioApparatusRemoveRequest,
  type StudioApparatusRemoveResponse,
  type StudioFieldInsertRequest,
  type StudioFieldInsertResponse,
  type StudioFieldListRequest,
  type StudioFieldListResponse,
  type StudioFieldRefreshRequest,
  type StudioFieldRefreshResponse,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyLogiczna, czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolajUczciwie } from './odmowa-rdzenia';

/**
 * Siedem komend aparatu dokumentu i pól — dwie rodziny o jednej naturze.
 *
 * Aparat (`apparatus.*`) niesie trzynaście rodzajów elementów: spis treści, spis
 * ilustracji i tabel, przypis dolny i końcowy, podpis, zakładkę, odwołanie
 * wzajemne, odsyłacz, powołanie, bibliografię, hasło indeksu i indeks. Pola
 * (`field.*`) niosą numer strony, liczbę stron, datę, godzinę, tytuł i autora
 * dokumentu, właściwość i pole obliczane.
 *
 * ── Dlaczego razem ──────────────────────────────────────────────────────────
 * Obie rodziny mają tę samą oś: element WYLICZANY z dokumentu, który po zmianie
 * treści staje się nieświeży i wymaga odświeżenia. `StudioApparatusItem.stale`
 * i `StudioDocumentField.stale` są tym samym pojęciem, a panel prowadzi jeden
 * wykaz „do odświeżenia" dla obu. Rozdzielenie tego na dwa źródła kazałoby
 * pisać dwa razy to samo zdanie o nieświeżości.
 */
export interface AparatZrodlo {
  /* ── Aparat dokumentu ──────────────────────────────────────────────────── */
  zaloz(zadanie: StudioApparatusInsertRequest): Promise<Wynik<StudioApparatusInsertResponse>>;
  wykaz(zadanie: StudioApparatusListRequest): Promise<Wynik<StudioApparatusListResponse>>;
  odswiez(zadanie: StudioApparatusRefreshRequest): Promise<Wynik<StudioApparatusRefreshResponse>>;
  usun(zadanie: StudioApparatusRemoveRequest): Promise<Wynik<StudioApparatusRemoveResponse>>;

  /* ── Pola dokumentu ────────────────────────────────────────────────────── */
  wstawPole(zadanie: StudioFieldInsertRequest): Promise<Wynik<StudioFieldInsertResponse>>;
  wykazPol(zadanie: StudioFieldListRequest): Promise<Wynik<StudioFieldListResponse>>;
  odswiezPola(zadanie: StudioFieldRefreshRequest): Promise<Wynik<StudioFieldRefreshResponse>>;
}

export function utworzAparatZrodlo(kanal: Kanal): AparatZrodlo {
  return {
    async zaloz(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioApparatusInsert, zadanie),
        Command.StudioApparatusInsert,
        (tresc) => czyObiekt(tresc.item) && czyObiekt(tresc.balance),
      );
    },

    async wykaz(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioApparatusList, zadanie),
        Command.StudioApparatusList,
        (tresc) => czyTablica(tresc.items),
      );
    },

    async odswiez(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioApparatusRefresh, zadanie),
        Command.StudioApparatusRefresh,
        (tresc) => czyTablica(tresc.items) && czyObiekt(tresc.balance),
      );
    },

    async usun(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioApparatusRemove, zadanie),
        Command.StudioApparatusRemove,
        (tresc) => czyLogiczna(tresc.removed) && czyObiekt(tresc.balance),
      );
    },

    async wstawPole(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioFieldInsert, zadanie),
        Command.StudioFieldInsert,
        (tresc) => czyObiekt(tresc.field) && czyObiekt(tresc.balance),
      );
    },

    async wykazPol(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioFieldList, zadanie),
        Command.StudioFieldList,
        (tresc) => czyTablica(tresc.fields),
      );
    },

    async odswiezPola(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioFieldRefresh, zadanie),
        Command.StudioFieldRefresh,
        (tresc) => czyTablica(tresc.fields) && czyObiekt(tresc.balance),
      );
    },
  };
}
