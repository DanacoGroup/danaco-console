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
 * Siedem komend warsztatu szablonów pism.
 *
 * Wykaz szablonów (`studio.template.list`) i założenie dokumentu z szablonu
 * (`studio.template.apply`) stały już w galerii, ale były wykazem FABRYCZNYM —
 * nie było czym szablonu założyć ani zmienić. Te siedem komend to domykają:
 * `template.save` zakłada szablon z bieżącego dokumentu wraz z arkuszem stylów,
 * nastawami strony, nagłówkiem, stopką, tabelami i blokadami wzorcowymi,
 * `template.field.set` i `.field.list` prowadzą pola do wypełnienia,
 * `template.fill` wypełnia je wartościami, `template.import` i `.export` wnoszą
 * i oddają plik Operatora (`dotx`, `ott`), a `template.delete` usuwa szablon
 * własny.
 *
 * ── Szablonu fabrycznego się nie usuwa ──────────────────────────────────────
 * `template.delete` odmawia usunięcia szablonu fabrycznego nazwanym powodem —
 * tak samo jak `studio.operation.delete` dla operacji fabrycznych. Okno czyta
 * pole `deleted` i przy odpowiedzi „nie usunąłem" nazywa powód, zamiast zdejmować
 * pozycję z galerii i pozwolić jej wrócić przy następnym odczycie.
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
        // Bilans wniesienia jest obowiązkowy: szablon z pliku Operatora przejmuje
        // arkusz stylów, nastawy strony i pola tylko na tyle, na ile plik je
        // niesie, i okno ma powiedzieć, ile tego było.
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
        // Pole `missingRequired` jest nieobowiązkowe, ale gdy przychodzi, jest
        // odpowiedzią najważniejszą: dokument powstał z dziurami i Operator ma
        // to zobaczyć, zamiast dostać „wypełniono".
        (tresc) => czyObiekt(tresc.document) && czyObiekt(tresc.balance),
      );
    },
  };
}
