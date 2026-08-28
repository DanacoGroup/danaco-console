import {
  Command,
  type StudioTableConvertRequest,
  type StudioTableConvertResponse,
  type StudioTableFormatSetRequest,
  type StudioTableFormatSetResponse,
  type StudioTableInsertRequest,
  type StudioTableInsertResponse,
  type StudioTableListRequest,
  type StudioTableListResponse,
  type StudioTableSortRequest,
  type StudioTableSortResponse,
  type StudioTableStructureEditRequest,
  type StudioTableStructureEditResponse,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolajUczciwie } from './odmowa-rdzenia';

/**
 * Interfejs TabelaZrodlo obejmuje sześć komend tabeli dokumentu: wstawianie, edycję struktury, formatowanie, sortowanie, zamianę tekstu na tabelę i odwrotnie oraz odczyt wykazu tabel.
 */
export interface TabelaZrodlo {
  wstaw(zadanie: StudioTableInsertRequest): Promise<Wynik<StudioTableInsertResponse>>;
  budowa(
    zadanie: StudioTableStructureEditRequest,
  ): Promise<Wynik<StudioTableStructureEditResponse>>;
  postac(zadanie: StudioTableFormatSetRequest): Promise<Wynik<StudioTableFormatSetResponse>>;
  sortuj(zadanie: StudioTableSortRequest): Promise<Wynik<StudioTableSortResponse>>;
  zamien(zadanie: StudioTableConvertRequest): Promise<Wynik<StudioTableConvertResponse>>;
  wykaz(zadanie: StudioTableListRequest): Promise<Wynik<StudioTableListResponse>>;
}

export function utworzTabeleZrodlo(kanal: Kanal): TabelaZrodlo {
  return {
    async wstaw(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioTableInsert, zadanie),
        Command.StudioTableInsert,
        (tresc) => czyObiekt(tresc.table) && czyObiekt(tresc.balance),
      );
    },

    async budowa(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioTableStructureEdit, zadanie),
        Command.StudioTableStructureEdit,
        (tresc) => czyObiekt(tresc.table) && czyObiekt(tresc.balance),
      );
    },

    async postac(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioTableFormatSet, zadanie),
        Command.StudioTableFormatSet,
        (tresc) => czyObiekt(tresc.table) && czyObiekt(tresc.balance),
      );
    },

    async sortuj(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioTableSort, zadanie),
        Command.StudioTableSort,
        (tresc) => czyObiekt(tresc.table) && czyObiekt(tresc.balance),
      );
    },

    async zamien(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioTableConvert, zadanie),
        Command.StudioTableConvert,
        // Tabela w odpowiedzi jest nieobowiązkowa, ponieważ zamiana tabeli na tekst nie zostawia tabeli.
        (tresc) => czyObiekt(tresc.balance) && czyObiekt(tresc.form),
      );
    },

    async wykaz(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioTableList, zadanie),
        Command.StudioTableList,
        (tresc) => czyTablica(tresc.tables),
      );
    },
  };
}
