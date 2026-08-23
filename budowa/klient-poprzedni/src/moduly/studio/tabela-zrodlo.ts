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
 * Sześć komend tabeli dokumentu.
 *
 * Rodzina jest zamknięta i pokrywa cały zakres zamówiony: `table.insert` zakłada
 * tabelę o wskazanym rozmiarze albo ze stylu gotowego, `table.structure.edit`
 * wstawia i usuwa wiersz oraz kolumnę, scala i dzieli komórki,
 * `table.format.set` ustawia szerokości kolumn, obramowanie, cieniowanie,
 * wyrównanie w komórce, styl i powtarzanie wiersza nagłówkowego, `table.sort`
 * sortuje zawartość, `table.convert` zamienia tekst na tabelę i odwrotnie,
 * a `table.list` oddaje tabele wraz z komórkami i policzonymi szerokościami.
 *
 * Każda czynność zmieniająca oddaje `balance` i `actionId`: bilans mówi, co
 * pominęła blokada fragmentu, a wpis dziennika daje cofnięcie pojedyncze. Panel
 * czyta oba — bez tego scalenie komórek w zablokowanym fragmencie wyglądałoby
 * na wykonane.
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
        // Tabela w odpowiedzi jest nieobowiązkowa: zamiana tabeli NA TEKST tabeli
        // po sobie nie zostawia, więc jej brak jest wynikiem poprawnym.
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
