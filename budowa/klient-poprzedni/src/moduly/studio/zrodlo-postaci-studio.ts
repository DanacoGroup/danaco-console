import {
  Command,
  type StudioDocumentFormGetRequest,
  type StudioDocumentFormGetResponse,
  type StudioDocumentFormSaveRequest,
  type StudioDocumentFormSaveResponse,
  type StudioFormatCaseSetRequest,
  type StudioFormatCaseSetResponse,
  type StudioFormatCharacterGetRequest,
  type StudioFormatCharacterGetResponse,
  type StudioFormatCharacterSetRequest,
  type StudioFormatCharacterSetResponse,
  type StudioFormatClearRequest,
  type StudioFormatClearResponse,
  type StudioFormatPainterApplyRequest,
  type StudioFormatPainterApplyResponse,
  type StudioFormatPainterCopyRequest,
  type StudioFormatPainterCopyResponse,
  type StudioFormatParagraphGetRequest,
  type StudioFormatParagraphGetResponse,
  type StudioFormatParagraphSetRequest,
  type StudioFormatParagraphSetResponse,
  type StudioFormatReplaceRequest,
  type StudioFormatReplaceResponse,
  type StudioFormatSimilarSelectRequest,
  type StudioFormatSimilarSelectResponse,
  type StudioListApplyRequest,
  type StudioListApplyResponse,
  type StudioListBulletSetRequest,
  type StudioListBulletSetResponse,
  type StudioListLevelIndentRequest,
  type StudioListLevelIndentResponse,
  type StudioListNumberingSetRequest,
  type StudioListNumberingSetResponse,
  type StudioListRestartRequest,
  type StudioListRestartResponse,
  type StudioPageBreakInsertRequest,
  type StudioPageBreakInsertResponse,
  type StudioPageEnvelopeSetRequest,
  type StudioPageEnvelopeSetResponse,
  type StudioPageHeaderfooterGetRequest,
  type StudioPageHeaderfooterGetResponse,
  type StudioPageHeaderfooterSetRequest,
  type StudioPageHeaderfooterSetResponse,
  type StudioPageNumberingSetRequest,
  type StudioPageNumberingSetResponse,
  type StudioPagePaperListRequest,
  type StudioPagePaperListResponse,
  type StudioPageSetupGetRequest,
  type StudioPageSetupGetResponse,
  type StudioPageSetupSetRequest,
  type StudioPageSetupSetResponse,
  type StudioPageWatermarkSetRequest,
  type StudioPageWatermarkSetResponse,
  type StudioProvenanceListRequest,
  type StudioProvenanceListResponse,
  type StudioRulerTabstopSetRequest,
  type StudioRulerTabstopSetResponse,
  type StudioSectionDeleteRequest,
  type StudioSectionDeleteResponse,
  type StudioSectionListRequest,
  type StudioSectionListResponse,
  type StudioSectionSaveRequest,
  type StudioSectionSaveResponse,
  type StudioStyleApplyRequest,
  type StudioStyleApplyResponse,
  type StudioStyleDeleteRequest,
  type StudioStyleDeleteResponse,
  type StudioStyleListRequest,
  type StudioStyleListResponse,
  type StudioStyleSaveRequest,
  type StudioStyleSaveResponse,
  type StudioSymbolAutoreplaceListRequest,
  type StudioSymbolAutoreplaceListResponse,
  type StudioSymbolAutoreplaceSetRequest,
  type StudioSymbolAutoreplaceSetResponse,
  type StudioSymbolInsertRequest,
  type StudioSymbolInsertResponse,
  type StudioSymbolListRequest,
  type StudioSymbolListResponse,
  type StudioTextEditRequest,
  type StudioTextEditResponse,
  type StudioTextGetRequest,
  type StudioTextGetResponse,
  type StudioViewGetRequest,
  type StudioViewGetResponse,
  type StudioViewSetRequest,
  type StudioViewSetResponse,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import {
  czyLiczba,
  czyLogiczna,
  czyObiekt,
  czyTablica,
  czyTekst,
  sprawdzKsztalt,
} from '../../protokol/ksztalt-odpowiedzi';
import { wywolajUczciwie } from './odmowa-rdzenia';

/**
 * Postać dokumentu — droga okna do rdzenia dla nastaw strony, stylów,
 * formatowania, list, znaków, widoku i pochodzenia. Źródło nie ma stanu i nie
 * buduje elementu — sprawdza kształt odpowiedzi i oddaje ją oknu.
 */
export interface ZrodloPostaciStudio {
  /* ── Postać w całości i treść fragmentu ──────────────────────────────────── */

  /** `studio.document.form.get` — arkusz stylów, nastawy, sekcje, bloki, pola. */
  postac(zadanie: StudioDocumentFormGetRequest): Promise<Wynik<StudioDocumentFormGetResponse>>;
  /** `studio.document.form.save` — zapis postaci wraz z treścią; postać przestaje ginąć. */
  zapiszPostac(
    zadanie: StudioDocumentFormSaveRequest,
  ): Promise<Wynik<StudioDocumentFormSaveResponse>>;
  /** `studio.text.get` — treść fragmentu wraz z jego postacią i blokadami. */
  tekst(zadanie: StudioTextGetRequest): Promise<Wynik<StudioTextGetResponse>>;
  /** `studio.text.edit` — brzmienie fragmentu bez przepisywania całości. */
  zmienTekst(zadanie: StudioTextEditRequest): Promise<Wynik<StudioTextEditResponse>>;
  /** `studio.provenance.list` — skąd pochodzi każdy wniesiony fragment. */
  pochodzenia(
    zadanie: StudioProvenanceListRequest,
  ): Promise<Wynik<StudioProvenanceListResponse>>;

  /* ── Nastawy strony, sekcje, nagłówki ────────────────────────────────────── */

  /** `studio.page.setup.get` — nastawy strony dokumentu albo sekcji. */
  nastawyStrony(zadanie: StudioPageSetupGetRequest): Promise<Wynik<StudioPageSetupGetResponse>>;
  /** `studio.page.setup.set` — nośnik, orientacja, cztery marginesy, oprawa, kolumny. */
  ustawStrone(zadanie: StudioPageSetupSetRequest): Promise<Wynik<StudioPageSetupSetResponse>>;
  /** `studio.page.paper.list` — formaty nośnika wraz z kopertami. */
  nosniki(zadanie: StudioPagePaperListRequest): Promise<Wynik<StudioPagePaperListResponse>>;
  /** `studio.page.envelope.set` — nadruk koperty: adresat, nadawca, położenia. */
  ustawKoperte(
    zadanie: StudioPageEnvelopeSetRequest,
  ): Promise<Wynik<StudioPageEnvelopeSetResponse>>;
  /** `studio.page.break.insert` — podział strony, kolumny, sekcji albo wiersza. */
  wstawPodzial(
    zadanie: StudioPageBreakInsertRequest,
  ): Promise<Wynik<StudioPageBreakInsertResponse>>;
  /** `studio.page.headerfooter.get` — nagłówki i stopki wedle zasięgu. */
  naglowki(
    zadanie: StudioPageHeaderfooterGetRequest,
  ): Promise<Wynik<StudioPageHeaderfooterGetResponse>>;
  /** `studio.page.headerfooter.set` — osobno dla sekcji, pierwszej strony i parzystych. */
  ustawNaglowek(
    zadanie: StudioPageHeaderfooterSetRequest,
  ): Promise<Wynik<StudioPageHeaderfooterSetResponse>>;
  /** `studio.page.numbering.set` — format, punkt startu, wznowienie i położenie numeru. */
  ustawNumeracje(
    zadanie: StudioPageNumberingSetRequest,
  ): Promise<Wynik<StudioPageNumberingSetResponse>>;
  /** `studio.page.watermark.set` — napis albo obraz wraz z kryciem i obrotem. */
  ustawZnakWodny(
    zadanie: StudioPageWatermarkSetRequest,
  ): Promise<Wynik<StudioPageWatermarkSetResponse>>;
  /** `studio.section.list` — sekcje wraz z nastawami, nagłówkami i numeracją. */
  sekcje(zadanie: StudioSectionListRequest): Promise<Wynik<StudioSectionListResponse>>;
  /** `studio.section.save` — sekcja o własnych nastawach; brak sekcji zakłada nową. */
  zapiszSekcje(zadanie: StudioSectionSaveRequest): Promise<Wynik<StudioSectionSaveResponse>>;
  /** `studio.section.delete` — treść przechodzi do sekcji poprzedniej. */
  usunSekcje(zadanie: StudioSectionDeleteRequest): Promise<Wynik<StudioSectionDeleteResponse>>;

  /* ── Linijka ─────────────────────────────────────────────────────────────── */

  /** `studio.ruler.tabstop.set` — tabulator wraz z rodzajem i znakiem wiodącym. */
  ustawTabulator(
    zadanie: StudioRulerTabstopSetRequest,
  ): Promise<Wynik<StudioRulerTabstopSetResponse>>;

  /* ── Style nazwane ───────────────────────────────────────────────────────── */

  /** `studio.style.list` — arkusz stylów wraz z dziedziczeniem i liczbą użyć. */
  style(zadanie: StudioStyleListRequest): Promise<Wynik<StudioStyleListResponse>>;
  /** `studio.style.save` — zmiana stylu przestawia WSZYSTKIE miejsca jego użycia. */
  zapiszStyl(zadanie: StudioStyleSaveRequest): Promise<Wynik<StudioStyleSaveResponse>>;
  /** `studio.style.apply` — styl nazwany na wskazanym fragmencie. */
  zastosujStyl(zadanie: StudioStyleApplyRequest): Promise<Wynik<StudioStyleApplyResponse>>;
  /** `studio.style.delete` — styl własny; fabrycznego rdzeń nie usuwa. */
  usunStyl(zadanie: StudioStyleDeleteRequest): Promise<Wynik<StudioStyleDeleteResponse>>;

  /* ── Formatowanie fragmentu ──────────────────────────────────────────────── */

  /** `studio.format.character.get` — postać znaku wraz z tym, co niejednolite. */
  postacZnaku(
    zadanie: StudioFormatCharacterGetRequest,
  ): Promise<Wynik<StudioFormatCharacterGetResponse>>;
  /** `studio.format.character.set` — krój, stopień, barwa, wyróżnienie, kapitaliki. */
  ustawZnak(
    zadanie: StudioFormatCharacterSetRequest,
  ): Promise<Wynik<StudioFormatCharacterSetResponse>>;
  /** `studio.format.paragraph.get` — postać akapitu obowiązująca na fragmencie. */
  postacAkapitu(
    zadanie: StudioFormatParagraphGetRequest,
  ): Promise<Wynik<StudioFormatParagraphGetResponse>>;
  /** `studio.format.paragraph.set` — wyrównanie, wcięcia, odstępy, interlinia. */
  ustawAkapit(
    zadanie: StudioFormatParagraphSetRequest,
  ): Promise<Wynik<StudioFormatParagraphSetResponse>>;
  /** `studio.format.clear` — postać wraca do stylu nazwanego albo domyślnej. */
  wyczyscFormat(zadanie: StudioFormatClearRequest): Promise<Wynik<StudioFormatClearResponse>>;
  /** `studio.format.case.set` — wielkość liter fragmentu. */
  ustawWielkoscLiter(
    zadanie: StudioFormatCaseSetRequest,
  ): Promise<Wynik<StudioFormatCaseSetResponse>>;
  /** `studio.format.painter.copy` — malarz formatów pobiera postać, nie treść. */
  pobierzPostac(
    zadanie: StudioFormatPainterCopyRequest,
  ): Promise<Wynik<StudioFormatPainterCopyResponse>>;
  /** `studio.format.painter.apply` — malarz formatów nanosi pobraną postać. */
  nalozPostac(
    zadanie: StudioFormatPainterApplyRequest,
  ): Promise<Wynik<StudioFormatPainterApplyResponse>>;
  /** `studio.format.similar.select` — zaznacz wedle podobnego formatowania. */
  podobnePostacia(
    zadanie: StudioFormatSimilarSelectRequest,
  ): Promise<Wynik<StudioFormatSimilarSelectResponse>>;
  /** `studio.format.replace` — znajdź i zamień wraz z postacią; blokady wracają bilansem. */
  zamienZPostacia(
    zadanie: StudioFormatReplaceRequest,
  ): Promise<Wynik<StudioFormatReplaceResponse>>;

  /* ── Listy ───────────────────────────────────────────────────────────────── */

  /** `studio.list.apply` — wypunktowanie, numeracja albo lista wielopoziomowa. */
  zastosujListe(zadanie: StudioListApplyRequest): Promise<Wynik<StudioListApplyResponse>>;
  /** `studio.list.bullet.set` — znak wypunktowania poziomu wraz z wcięciem. */
  ustawPunktator(
    zadanie: StudioListBulletSetRequest,
  ): Promise<Wynik<StudioListBulletSetResponse>>;
  /** `studio.list.numbering.set` — format numeracji poziomu wraz ze wzorem numeru. */
  ustawNumeracjeListy(
    zadanie: StudioListNumberingSetRequest,
  ): Promise<Wynik<StudioListNumberingSetResponse>>;
  /** `studio.list.restart` — wznowienie numeracji od wskazanego miejsca. */
  wznowNumeracje(zadanie: StudioListRestartRequest): Promise<Wynik<StudioListRestartResponse>>;
  /** `studio.list.level.indent` — poziom listy fragmentu w górę albo w dół. */
  przestawPoziom(
    zadanie: StudioListLevelIndentRequest,
  ): Promise<Wynik<StudioListLevelIndentResponse>>;

  /* ── Znaki specjalne ─────────────────────────────────────────────────────── */

  /** `studio.symbol.list` — tablica znaków wraz z szukaniem i ostatnio użytymi. */
  znaki(zadanie: StudioSymbolListRequest): Promise<Wynik<StudioSymbolListResponse>>;
  /** `studio.symbol.insert` — znak w miejsce kursora. */
  wstawZnak(zadanie: StudioSymbolInsertRequest): Promise<Wynik<StudioSymbolInsertResponse>>;
  /** `studio.symbol.autoreplace.list` — zasady autozamiany skrótów. */
  autozamiany(
    zadanie: StudioSymbolAutoreplaceListRequest,
  ): Promise<Wynik<StudioSymbolAutoreplaceListResponse>>;
  /** `studio.symbol.autoreplace.set` — zasada autozamiany; puste zastąpienie ją usuwa. */
  ustawAutozamiane(
    zadanie: StudioSymbolAutoreplaceSetRequest,
  ): Promise<Wynik<StudioSymbolAutoreplaceSetResponse>>;

  /* ── Widok ───────────────────────────────────────────────────────────────── */

  /** `studio.view.get` — nastawy widoku pamiętane przy dokumencie. */
  widok(zadanie: StudioViewGetRequest): Promise<Wynik<StudioViewGetResponse>>;
  /** `studio.view.set` — tryb powierzchni, skala, linijki, układ stron, przewijanie. */
  ustawWidok(zadanie: StudioViewSetRequest): Promise<Wynik<StudioViewSetResponse>>;
}

/**
 * Sprawdzian odpowiedzi czynności zmieniającej postać. Każda z nich oddaje
 * `form` i `balance` — postać po zmianie oraz bilans tego, co zmienione i co
 * pominięte przez blokadę; bilans jest polem obowiązkowym kontraktu.
 */
function czyPostacIBilans(tresc: { form?: unknown; balance?: unknown }): boolean {
  return (
    czyObiekt(tresc.form) &&
    czyObiekt(tresc.balance) &&
    czyLiczba((tresc.balance as { applied?: unknown }).applied) &&
    czyLiczba((tresc.balance as { skippedCount?: unknown }).skippedCount)
  );
}

export function utworzZrodloPostaciStudio(kanal: Kanal): ZrodloPostaciStudio {
  return {
    async postac(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioDocumentFormGet, zadanie),
        Command.StudioDocumentFormGet,
        (tresc) => czyObiekt(tresc.form),
      );
    },

    async zapiszPostac(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioDocumentFormSave, zadanie),
        Command.StudioDocumentFormSave,
        (tresc) => czyObiekt(tresc.document) && czyPostacIBilans(tresc),
      );
    },

    async tekst(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioTextGet, zadanie),
        Command.StudioTextGet,
        // Treść pusta jest odpowiedzią prawdziwą — sprawdza się obecność pola, nie jego długość.
        (tresc) => czyTekst(tresc.text),
      );
    },

    async zmienTekst(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioTextEdit, zadanie),
        Command.StudioTextEdit,
        czyPostacIBilans,
      );
    },

    async pochodzenia(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioProvenanceList, zadanie),
        Command.StudioProvenanceList,
        (tresc) => czyTablica(tresc.entries),
      );
    },

    async nastawyStrony(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioPageSetupGet, zadanie),
        Command.StudioPageSetupGet,
        (tresc) => czyObiekt(tresc.pageSetup),
      );
    },

    async ustawStrone(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioPageSetupSet, zadanie),
        Command.StudioPageSetupSet,
        (tresc) => czyObiekt(tresc.pageSetup) && czyPostacIBilans(tresc),
      );
    },

    async nosniki(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioPagePaperList, zadanie),
        Command.StudioPagePaperList,
        (tresc) => czyTablica(tresc.papers),
      );
    },

    async ustawKoperte(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioPageEnvelopeSet, zadanie),
        Command.StudioPageEnvelopeSet,
        (tresc) => czyObiekt(tresc.envelope) && czyPostacIBilans(tresc),
      );
    },

    async wstawPodzial(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioPageBreakInsert, zadanie),
        Command.StudioPageBreakInsert,
        czyPostacIBilans,
      );
    },

    async naglowki(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioPageHeaderfooterGet, zadanie),
        Command.StudioPageHeaderfooterGet,
        (tresc) => czyTablica(tresc.headersFooters),
      );
    },

    async ustawNaglowek(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioPageHeaderfooterSet, zadanie),
        Command.StudioPageHeaderfooterSet,
        (tresc) => czyTablica(tresc.headersFooters) && czyPostacIBilans(tresc),
      );
    },

    async ustawNumeracje(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioPageNumberingSet, zadanie),
        Command.StudioPageNumberingSet,
        // enabled: false jest odpowiedzią poprawną — pilnuje się obecności pola, nie jego wartości.
        (tresc) =>
          czyObiekt(tresc.numbering) &&
          czyLogiczna((tresc.numbering as { enabled?: unknown }).enabled) &&
          czyPostacIBilans(tresc),
      );
    },

    async ustawZnakWodny(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioPageWatermarkSet, zadanie),
        Command.StudioPageWatermarkSet,
        (tresc) => czyObiekt(tresc.watermark) && czyPostacIBilans(tresc),
      );
    },

    async sekcje(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioSectionList, zadanie),
        Command.StudioSectionList,
        (tresc) => czyTablica(tresc.sections),
      );
    },

    async zapiszSekcje(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioSectionSave, zadanie),
        Command.StudioSectionSave,
        (tresc) => czyObiekt(tresc.section) && czyPostacIBilans(tresc),
      );
    },

    async usunSekcje(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioSectionDelete, zadanie),
        Command.StudioSectionDelete,
        (tresc) => czyLogiczna(tresc.deleted) && czyPostacIBilans(tresc),
      );
    },

    async ustawTabulator(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioRulerTabstopSet, zadanie),
        Command.StudioRulerTabstopSet,
        (tresc) => czyTablica(tresc.tabStops) && czyPostacIBilans(tresc),
      );
    },

    async style(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioStyleList, zadanie),
        Command.StudioStyleList,
        (tresc) => czyTablica(tresc.styles),
      );
    },

    async zapiszStyl(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioStyleSave, zadanie),
        Command.StudioStyleSave,
        (tresc) => czyObiekt(tresc.style) && czyPostacIBilans(tresc),
      );
    },

    async zastosujStyl(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioStyleApply, zadanie),
        Command.StudioStyleApply,
        czyPostacIBilans,
      );
    },

    async usunStyl(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioStyleDelete, zadanie),
        Command.StudioStyleDelete,
        // deleted: false z odmową nazywającą powód jest odpowiedzią poprawną dla stylu fabrycznego.
        (tresc) => czyLogiczna(tresc.deleted) && czyPostacIBilans(tresc),
      );
    },

    async postacZnaku(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioFormatCharacterGet, zadanie),
        Command.StudioFormatCharacterGet,
        (tresc) => czyObiekt(tresc.character),
      );
    },

    async ustawZnak(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioFormatCharacterSet, zadanie),
        Command.StudioFormatCharacterSet,
        czyPostacIBilans,
      );
    },

    async postacAkapitu(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioFormatParagraphGet, zadanie),
        Command.StudioFormatParagraphGet,
        (tresc) => czyObiekt(tresc.paragraph),
      );
    },

    async ustawAkapit(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioFormatParagraphSet, zadanie),
        Command.StudioFormatParagraphSet,
        czyPostacIBilans,
      );
    },

    async wyczyscFormat(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioFormatClear, zadanie),
        Command.StudioFormatClear,
        czyPostacIBilans,
      );
    },

    async ustawWielkoscLiter(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioFormatCaseSet, zadanie),
        Command.StudioFormatCaseSet,
        czyPostacIBilans,
      );
    },

    async pobierzPostac(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioFormatPainterCopy, zadanie),
        Command.StudioFormatPainterCopy,
        // Uchwyt pusty byłby malarzem bez czego nanieść — druga połowa czynności nie ma czym zadziałać.
        (tresc) => czyTekst(tresc.clipId) && tresc.clipId !== '',
      );
    },

    async nalozPostac(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioFormatPainterApply, zadanie),
        Command.StudioFormatPainterApply,
        czyPostacIBilans,
      );
    },

    async podobnePostacia(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioFormatSimilarSelect, zadanie),
        Command.StudioFormatSimilarSelect,
        (tresc) => czyTablica(tresc.matches) && czyLiczba(tresc.count),
      );
    },

    async zamienZPostacia(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioFormatReplace, zadanie),
        Command.StudioFormatReplace,
        (tresc) => czyLiczba(tresc.matches) && czyLiczba(tresc.replaced) && czyPostacIBilans(tresc),
      );
    },

    async zastosujListe(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioListApply, zadanie),
        Command.StudioListApply,
        (tresc) => czyObiekt(tresc.list) && czyPostacIBilans(tresc),
      );
    },

    async ustawPunktator(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioListBulletSet, zadanie),
        Command.StudioListBulletSet,
        (tresc) => czyObiekt(tresc.list) && czyPostacIBilans(tresc),
      );
    },

    async ustawNumeracjeListy(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioListNumberingSet, zadanie),
        Command.StudioListNumberingSet,
        (tresc) => czyObiekt(tresc.list) && czyPostacIBilans(tresc),
      );
    },

    async wznowNumeracje(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioListRestart, zadanie),
        Command.StudioListRestart,
        czyPostacIBilans,
      );
    },

    async przestawPoziom(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioListLevelIndent, zadanie),
        Command.StudioListLevelIndent,
        czyPostacIBilans,
      );
    },

    async znaki(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioSymbolList, zadanie),
        Command.StudioSymbolList,
        (tresc) => czyTablica(tresc.symbols),
      );
    },

    async wstawZnak(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioSymbolInsert, zadanie),
        Command.StudioSymbolInsert,
        (tresc) => czyObiekt(tresc.symbol) && czyPostacIBilans(tresc),
      );
    },

    async autozamiany(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioSymbolAutoreplaceList, zadanie),
        Command.StudioSymbolAutoreplaceList,
        (tresc) => czyTablica(tresc.rules),
      );
    },

    async ustawAutozamiane(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioSymbolAutoreplaceSet, zadanie),
        Command.StudioSymbolAutoreplaceSet,
        // Kontrakt oddaje albo zasadę po zapisie, albo oznaczenie usunięcia — obie są prawdziwe.
        (tresc) => czyObiekt(tresc.rule) || czyLogiczna(tresc.removed),
      );
    },

    async widok(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioViewGet, zadanie),
        Command.StudioViewGet,
        (tresc) => czyObiekt(tresc.settings),
      );
    },

    async ustawWidok(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioViewSet, zadanie),
        Command.StudioViewSet,
        (tresc) => czyObiekt(tresc.settings),
      );
    },
  };
}
