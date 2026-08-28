import {
  Command,
  type StudioTemplateDeleteRequest,
  type StudioTemplateDeleteResponse,
  type StudioTemplateExportRequest,
  type StudioTemplateExportResponse,
  type StudioTemplateFieldListRequest,
  type StudioTemplateFieldListResponse,
  type StudioTemplateFieldSetRequest,
  type StudioTemplateFieldSetResponse,
  type StudioTemplateFillRequest,
  type StudioTemplateFillResponse,
  type StudioTemplateImportRequest,
  type StudioTemplateImportResponse,
  type StudioTemplateSaveRequest,
  type StudioTemplateSaveResponse,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyLogiczna, czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolajUczciwie } from './odmowa-rdzenia';

/**
 * Interfejs SzablonZrodlo obejmuje siedem komend warsztatu szablonów pism: zapis, usunięcie, ustawienie i odczyt pól, wniesienie i eksport pliku szablonu oraz wypełnienie pól wartościami.
 */
export interface SzablonZrodlo {
  zapisz(zadanie: StudioTemplateSaveRequest): Promise<Wynik<StudioTemplateSaveResponse>>;
  usun(zadanie: StudioTemplateDeleteRequest): Promise<Wynik<StudioTemplateDeleteResponse>>;
  ustawPole(
    zadanie: StudioTemplateFieldSetRequest,
  ): Promise<Wynik<StudioTemplateFieldSetResponse>>;
  pola(zadanie: StudioTemplateFieldListRequest): Promise<Wynik<StudioTemplateFieldListResponse>>;
  wnies(zadanie: StudioTemplateImportRequest): Promise<Wynik<StudioTemplateImportResponse>>;
  oddaj(zadanie: StudioTemplateExportRequest): Promise<Wynik<StudioTemplateExportResponse>>;
  wypelnij(zadanie: StudioTemplateFillRequest): Promise<Wynik<StudioTemplateFillResponse>>;
}

export function utworzSzablonZrodlo(kanal: Kanal): SzablonZrodlo {
  return {
    async zapisz(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioTemplateSave, zadanie),
        Command.StudioTemplateSave,
        (tresc) => czyObiekt(tresc.template),
      );
    },

    async usun(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioTemplateDelete, zadanie),
        Command.StudioTemplateDelete,
        (tresc) => czyLogiczna(tresc.deleted),
      );
    },

    async ustawPole(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioTemplateFieldSet, zadanie),
        Command.StudioTemplateFieldSet,
        (tresc) => czyObiekt(tresc.template) && czyTablica(tresc.fields),
      );
    },

    async pola(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioTemplateFieldList, zadanie),
        Command.StudioTemplateFieldList,
        (tresc) => czyTablica(tresc.fields),
      );
    },

    async wnies(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioTemplateImport, zadanie),
        Command.StudioTemplateImport,
        // Bilans wniesienia jest obowiązkowy i pokazuje, ile arkusza stylów, nastaw strony i pól przejął plik.
        (tresc) => czyObiekt(tresc.template) && czyObiekt(tresc.balance),
      );
    },

    async oddaj(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioTemplateExport, zadanie),
        Command.StudioTemplateExport,
        (tresc) => czyObiekt(tresc.result),
      );
    },

    async wypelnij(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioTemplateFill, zadanie),
        Command.StudioTemplateFill,
        // Pole missingRequired oznacza, gdy występuje, dokument z niewypełnionymi polami wymaganymi.
        (tresc) => czyObiekt(tresc.document) && czyObiekt(tresc.balance),
      );
    },
  };
}
