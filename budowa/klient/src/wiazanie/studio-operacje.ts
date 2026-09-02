// Katalog operacji Studia: nazwa komendy, etykieta i grupa. Widok nie idzie
// stąd — pozycje stają w panelu operacji, a żądanie składa renderer z kontekstu
// dokumentu i pól panelu. Komendy pochodzą z rodziny studio.* kontraktu.

import { Command } from '../../../shared/contract.ts';

export type GrupaOperacji =
  | 'strony'
  | 'styl'
  | 'tabela'
  | 'wstawki'
  | 'szablony'
  | 'znakowanie'
  | 'wersje'
  | 'pdf'
  | 'ochrona'
  | 'przyjmowanie'
  | 'praca';

export interface OperacjaStudia {
  komenda: Command;
  nazwa: string;
  grupa: GrupaOperacji;
}

export const OPERACJE_STUDIA: readonly OperacjaStudia[] = [
  { komenda: Command.StudioPageSetupGet, nazwa: 'Ustawienia strony — odczyt', grupa: 'strony' },
  { komenda: Command.StudioPageSetupSet, nazwa: 'Ustawienia strony — zapis', grupa: 'strony' },
  { komenda: Command.StudioPagePaperList, nazwa: 'Formaty papieru', grupa: 'strony' },
  { komenda: Command.StudioPageEnvelopeSet, nazwa: 'Koperta', grupa: 'strony' },
  { komenda: Command.StudioPageBreakInsert, nazwa: 'Wstaw podział strony', grupa: 'strony' },
  { komenda: Command.StudioPageHeaderfooterGet, nazwa: 'Nagłówek i stopka — odczyt', grupa: 'strony' },
  { komenda: Command.StudioPageHeaderfooterSet, nazwa: 'Nagłówek i stopka — zapis', grupa: 'strony' },
  { komenda: Command.StudioPageNumberingSet, nazwa: 'Numeracja stron', grupa: 'strony' },
  { komenda: Command.StudioPageWatermarkSet, nazwa: 'Znak wodny', grupa: 'strony' },
  { komenda: Command.StudioSectionList, nazwa: 'Sekcje — wykaz', grupa: 'strony' },
  { komenda: Command.StudioSectionSave, nazwa: 'Sekcja — zapis', grupa: 'strony' },
  { komenda: Command.StudioSectionDelete, nazwa: 'Sekcja — usuń', grupa: 'strony' },
  { komenda: Command.StudioRulerTabstopSet, nazwa: 'Tabulator linijki', grupa: 'strony' },

  { komenda: Command.StudioStyleSave, nazwa: 'Styl — zapis', grupa: 'styl' },
  { komenda: Command.StudioStyleDelete, nazwa: 'Styl — usuń', grupa: 'styl' },
  { komenda: Command.StudioFormatCaseSet, nazwa: 'Wielkość liter', grupa: 'styl' },
  { komenda: Command.StudioFormatClear, nazwa: 'Wyczyść formatowanie', grupa: 'styl' },
  { komenda: Command.StudioFormatPainterCopy, nazwa: 'Malarz formatu — pobierz', grupa: 'styl' },
  { komenda: Command.StudioFormatPainterApply, nazwa: 'Malarz formatu — nanieś', grupa: 'styl' },
  { komenda: Command.StudioFormatReplace, nazwa: 'Zamień formatowanie', grupa: 'styl' },
  { komenda: Command.StudioFormatSimilarSelect, nazwa: 'Zaznacz podobne', grupa: 'styl' },
  { komenda: Command.StudioListBulletSet, nazwa: 'Punktory listy', grupa: 'styl' },
  { komenda: Command.StudioListNumberingSet, nazwa: 'Numeracja listy', grupa: 'styl' },
  { komenda: Command.StudioListRestart, nazwa: 'Restart numeracji listy', grupa: 'styl' },
  { komenda: Command.StudioListLevelIndent, nazwa: 'Poziom wcięcia listy', grupa: 'styl' },

  { komenda: Command.StudioTableConvert, nazwa: 'Tekst na tabelę', grupa: 'tabela' },
  { komenda: Command.StudioTableFormatSet, nazwa: 'Format tabeli', grupa: 'tabela' },
  { komenda: Command.StudioTableList, nazwa: 'Tabele — wykaz', grupa: 'tabela' },
  { komenda: Command.StudioTableSort, nazwa: 'Sortuj tabelę', grupa: 'tabela' },
  { komenda: Command.StudioTableStructureEdit, nazwa: 'Struktura tabeli', grupa: 'tabela' },

  { komenda: Command.StudioObjectInsert, nazwa: 'Wstaw obiekt', grupa: 'wstawki' },
  { komenda: Command.StudioObjectList, nazwa: 'Obiekty — wykaz', grupa: 'wstawki' },
  { komenda: Command.StudioObjectRemove, nazwa: 'Usuń obiekt', grupa: 'wstawki' },
  { komenda: Command.StudioObjectFormatSet, nazwa: 'Format obiektu', grupa: 'wstawki' },
  { komenda: Command.StudioFieldInsert, nazwa: 'Wstaw pole', grupa: 'wstawki' },
  { komenda: Command.StudioFieldList, nazwa: 'Pola — wykaz', grupa: 'wstawki' },
  { komenda: Command.StudioFieldRefresh, nazwa: 'Odśwież pola', grupa: 'wstawki' },
  { komenda: Command.StudioSymbolInsert, nazwa: 'Wstaw symbol', grupa: 'wstawki' },
  { komenda: Command.StudioSymbolList, nazwa: 'Symbole — wykaz', grupa: 'wstawki' },
  { komenda: Command.StudioSymbolAutoreplaceList, nazwa: 'Autozamiana symboli — wykaz', grupa: 'wstawki' },
  { komenda: Command.StudioSymbolAutoreplaceSet, nazwa: 'Autozamiana symboli — zapis', grupa: 'wstawki' },
  { komenda: Command.StudioApparatusInsert, nazwa: 'Wstaw aparat naukowy', grupa: 'wstawki' },
  { komenda: Command.StudioApparatusList, nazwa: 'Aparat — wykaz', grupa: 'wstawki' },
  { komenda: Command.StudioApparatusRefresh, nazwa: 'Odśwież aparat', grupa: 'wstawki' },
  { komenda: Command.StudioApparatusRemove, nazwa: 'Usuń aparat', grupa: 'wstawki' },
  { komenda: Command.StudioAssetEmbed, nazwa: 'Osadź zasób', grupa: 'wstawki' },
  { komenda: Command.StudioInsertFromLibrary, nazwa: 'Wstaw z biblioteki', grupa: 'wstawki' },
  { komenda: Command.StudioInsertFromWeb, nazwa: 'Wstaw z sieci', grupa: 'wstawki' },
  { komenda: Command.StudioClipboardCopy, nazwa: 'Kopiuj do schowka', grupa: 'wstawki' },
  { komenda: Command.StudioClipboardPaste, nazwa: 'Wklej ze schowka', grupa: 'wstawki' },
  { komenda: Command.StudioTextEdit, nazwa: 'Edytuj tekst', grupa: 'wstawki' },

  { komenda: Command.StudioTemplateApply, nazwa: 'Zastosuj szablon', grupa: 'szablony' },
  { komenda: Command.StudioTemplateDelete, nazwa: 'Usuń szablon', grupa: 'szablony' },
  { komenda: Command.StudioTemplateExport, nazwa: 'Eksport szablonu', grupa: 'szablony' },
  { komenda: Command.StudioTemplateImport, nazwa: 'Import szablonu', grupa: 'szablony' },
  { komenda: Command.StudioTemplateList, nazwa: 'Szablony — wykaz', grupa: 'szablony' },
  { komenda: Command.StudioTemplateSave, nazwa: 'Zapis szablonu', grupa: 'szablony' },
  { komenda: Command.StudioTemplateFill, nazwa: 'Wypełnij szablon', grupa: 'szablony' },
  { komenda: Command.StudioTemplateFieldList, nazwa: 'Pola szablonu — wykaz', grupa: 'szablony' },
  { komenda: Command.StudioTemplateFieldSet, nazwa: 'Pola szablonu — zapis', grupa: 'szablony' },
  { komenda: Command.StudioExportProfileSave, nazwa: 'Profil eksportu — zapis', grupa: 'szablony' },

  { komenda: Command.StudioMarkupAdd, nazwa: 'Znakowanie — dodaj', grupa: 'znakowanie' },
  { komenda: Command.StudioMarkupDecide, nazwa: 'Znakowanie — rozstrzygnij', grupa: 'znakowanie' },
  { komenda: Command.StudioMarkupList, nazwa: 'Znakowanie — wykaz', grupa: 'znakowanie' },
  { komenda: Command.StudioMarkupRemove, nazwa: 'Znakowanie — usuń', grupa: 'znakowanie' },
  { komenda: Command.StudioMarkupTypeList, nazwa: 'Rodzaje znakowania — wykaz', grupa: 'znakowanie' },
  { komenda: Command.StudioMarkupTypeSave, nazwa: 'Rodzaj znakowania — zapis', grupa: 'znakowanie' },
  { komenda: Command.StudioMarkupTypeDelete, nazwa: 'Rodzaj znakowania — usuń', grupa: 'znakowanie' },
  { komenda: Command.StudioCommentAdd, nazwa: 'Komentarz — dodaj', grupa: 'znakowanie' },
  { komenda: Command.StudioCommentList, nazwa: 'Komentarze — wykaz', grupa: 'znakowanie' },
  { komenda: Command.StudioCommentResolve, nazwa: 'Komentarz — domknij', grupa: 'znakowanie' },
  { komenda: Command.StudioAnnotationAdd, nazwa: 'Adnotacja — dodaj', grupa: 'znakowanie' },
  { komenda: Command.StudioAnnotationList, nazwa: 'Adnotacje — wykaz', grupa: 'znakowanie' },
  { komenda: Command.StudioLockAdd, nazwa: 'Blokada — dodaj', grupa: 'znakowanie' },
  { komenda: Command.StudioLockList, nazwa: 'Blokady — wykaz', grupa: 'znakowanie' },
  { komenda: Command.StudioLockRemove, nazwa: 'Blokada — usuń', grupa: 'znakowanie' },

  { komenda: Command.StudioBranchCreate, nazwa: 'Gałąź — utwórz', grupa: 'wersje' },
  { komenda: Command.StudioBranchList, nazwa: 'Gałęzie — wykaz', grupa: 'wersje' },
  { komenda: Command.StudioBranchMerge, nazwa: 'Gałąź — scal', grupa: 'wersje' },
  { komenda: Command.StudioBackupCreate, nazwa: 'Kopia — utwórz', grupa: 'wersje' },
  { komenda: Command.StudioBackupList, nazwa: 'Kopie — wykaz', grupa: 'wersje' },
  { komenda: Command.StudioBackupRestore, nazwa: 'Kopia — przywróć', grupa: 'wersje' },
  { komenda: Command.StudioJournalList, nazwa: 'Dziennik — wykaz', grupa: 'wersje' },
  { komenda: Command.StudioJournalRedo, nazwa: 'Dziennik — ponów', grupa: 'wersje' },
  { komenda: Command.StudioJournalRevert, nazwa: 'Dziennik — cofnij', grupa: 'wersje' },
  { komenda: Command.StudioVersionReferenceCreate, nazwa: 'Wersja odniesienia — utwórz', grupa: 'wersje' },
  { komenda: Command.StudioVersionRestoreInitial, nazwa: 'Przywróć wersję pierwotną', grupa: 'wersje' },
  { komenda: Command.StudioModelChangesList, nazwa: 'Zmiany modelu — wykaz', grupa: 'wersje' },
  { komenda: Command.StudioModelChangesNavigate, nazwa: 'Zmiany modelu — przejdź', grupa: 'wersje' },
  { komenda: Command.StudioModelChangesRevert, nazwa: 'Zmiany modelu — cofnij', grupa: 'wersje' },
  { komenda: Command.StudioDiffSource, nazwa: 'Różnica treści źródłowej', grupa: 'wersje' },
  { komenda: Command.StudioDiffVisual, nazwa: 'Różnica wizualna', grupa: 'wersje' },
  { komenda: Command.StudioDiffReportExport, nazwa: 'Eksport raportu różnic', grupa: 'wersje' },
  { komenda: Command.StudioProposalDecide, nazwa: 'Propozycja — rozstrzygnij', grupa: 'wersje' },
  { komenda: Command.StudioProvenanceList, nazwa: 'Prowenancja — wykaz', grupa: 'wersje' },

  { komenda: Command.StudioPdfMerge, nazwa: 'PDF — scal', grupa: 'pdf' },
  { komenda: Command.StudioPdfSplit, nazwa: 'PDF — podziel', grupa: 'pdf' },
  { komenda: Command.StudioPdfExtract, nazwa: 'PDF — wydobądź strony', grupa: 'pdf' },
  { komenda: Command.StudioPdfPagesReorder, nazwa: 'PDF — przestaw strony', grupa: 'pdf' },
  { komenda: Command.StudioPdfOptimize, nazwa: 'PDF — optymalizuj', grupa: 'pdf' },
  { komenda: Command.StudioPdfBookmarksSet, nazwa: 'PDF — zakładki', grupa: 'pdf' },
  { komenda: Command.StudioPdfBates, nazwa: 'PDF — numeracja Batesa', grupa: 'pdf' },
  { komenda: Command.StudioPdfStamp, nazwa: 'PDF — stempel', grupa: 'pdf' },
  { komenda: Command.StudioPdfFormFill, nazwa: 'PDF — wypełnij formularz', grupa: 'pdf' },

  { komenda: Command.StudioSecurityEncrypt, nazwa: 'Szyfruj dokument', grupa: 'ochrona' },
  { komenda: Command.StudioSecurityMetadataStrip, nazwa: 'Zdejmij metadane', grupa: 'ochrona' },
  { komenda: Command.StudioSecurityRedact, nazwa: 'Wymaż fragment', grupa: 'ochrona' },
  { komenda: Command.StudioSecuritySensitiveDetect, nazwa: 'Wykryj dane wrażliwe', grupa: 'ochrona' },
  { komenda: Command.StudioSecuritySign, nazwa: 'Podpisz dokument', grupa: 'ochrona' },
  { komenda: Command.StudioSecuritySignVerify, nazwa: 'Sprawdź podpis', grupa: 'ochrona' },

  { komenda: Command.StudioIngestQueueAdd, nazwa: 'Przyjmowanie — dodaj do kolejki', grupa: 'przyjmowanie' },
  { komenda: Command.StudioIngestQueueList, nazwa: 'Przyjmowanie — kolejka', grupa: 'przyjmowanie' },
  { komenda: Command.StudioIngestUrl, nazwa: 'Przyjmowanie — z adresu', grupa: 'przyjmowanie' },
  { komenda: Command.StudioIngestRecognize, nazwa: 'Przyjmowanie — rozpoznaj', grupa: 'przyjmowanie' },
  { komenda: Command.StudioIngestCorrectionSet, nazwa: 'Przyjmowanie — korekta', grupa: 'przyjmowanie' },
  { komenda: Command.StudioIngestItemAccept, nazwa: 'Przyjmowanie — przyjmij pozycję', grupa: 'przyjmowanie' },
  { komenda: Command.StudioIngestDeviceList, nazwa: 'Przyjmowanie — urządzenia', grupa: 'przyjmowanie' },
  { komenda: Command.StudioIngestDeviceScan, nazwa: 'Przyjmowanie — skanuj', grupa: 'przyjmowanie' },

  { komenda: Command.StudioOperationSave, nazwa: 'Operacja własna — zapis', grupa: 'praca' },
  { komenda: Command.StudioOperationDelete, nazwa: 'Operacja własna — usuń', grupa: 'praca' },
  { komenda: Command.StudioChainList, nazwa: 'Łańcuchy operacji — wykaz', grupa: 'praca' },
  { komenda: Command.StudioChainSave, nazwa: 'Łańcuch operacji — zapis', grupa: 'praca' },
  { komenda: Command.StudioChainRun, nazwa: 'Łańcuch operacji — uruchom', grupa: 'praca' },
  { komenda: Command.StudioPlanCreate, nazwa: 'Plan pracy — utwórz', grupa: 'praca' },
  { komenda: Command.StudioPlanRun, nazwa: 'Plan pracy — uruchom', grupa: 'praca' },
  { komenda: Command.StudioPlanStop, nazwa: 'Plan pracy — zatrzymaj', grupa: 'praca' },
  { komenda: Command.StudioBatchRun, nazwa: 'Przebieg wsadowy', grupa: 'praca' },
  { komenda: Command.StudioPackageExport, nazwa: 'Eksport pakietu', grupa: 'praca' },
  { komenda: Command.StudioSearchSemantic, nazwa: 'Wyszukiwanie znaczeniowe', grupa: 'praca' },
  { komenda: Command.StudioDocumentCopy, nazwa: 'Dokument — kopiuj', grupa: 'praca' },
  { komenda: Command.StudioDocumentSaveAs, nazwa: 'Dokument — zapisz jako', grupa: 'praca' },
  { komenda: Command.StudioDocumentExportBatch, nazwa: 'Dokument — eksport wsadowy', grupa: 'praca' },
  { komenda: Command.StudioDocumentImportFile, nazwa: 'Dokument — import pliku', grupa: 'praca' },
  { komenda: Command.StudioDocumentImportPdf, nazwa: 'Dokument — import PDF', grupa: 'praca' },
  { komenda: Command.StudioDocumentImageImport, nazwa: 'Dokument — import obrazu', grupa: 'praca' },
  { komenda: Command.StudioAgentsSlotsList, nazwa: 'Sloty ekspertów — wykaz', grupa: 'praca' },
  { komenda: Command.StudioAgentsClaim, nazwa: 'Slot eksperta — zajmij', grupa: 'praca' },
  { komenda: Command.StudioAgentsRelease, nazwa: 'Slot eksperta — zwolnij', grupa: 'praca' },
  { komenda: Command.StudioAgentsConflictsList, nazwa: 'Konflikty ekspertów — wykaz', grupa: 'praca' },
  { komenda: Command.StudioAgentsSettingsGet, nazwa: 'Nastawy ekspertów — odczyt', grupa: 'praca' },
  { komenda: Command.StudioAgentsSettingsSet, nazwa: 'Nastawy ekspertów — zapis', grupa: 'praca' },
];

export const NAZWY_GRUP: Readonly<Record<GrupaOperacji, string>> = {
  strony: 'Strony i sekcje',
  styl: 'Styl i formatowanie',
  tabela: 'Tabele',
  wstawki: 'Wstawki',
  szablony: 'Szablony',
  znakowanie: 'Znakowanie i komentarze',
  wersje: 'Wersje i różnice',
  pdf: 'PDF',
  ochrona: 'Ochrona',
  przyjmowanie: 'Przyjmowanie treści',
  praca: 'Praca i eksport',
};
