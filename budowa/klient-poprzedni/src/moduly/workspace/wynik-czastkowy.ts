/**
 * Stałe dziedzinowe modułu Workspace używane przez oba źródła modułu.
 *
 * Wycinanie pola z odpowiedzi rdzenia (`przenies`) mieszka
 * w `protokol/wynik-czastkowy.ts`, bo jest pomocnikiem nad `Wynik`,
 * a nie wiedzą modułu.
 */

/** Kod modułu Workspace, który znakuje pliki wgrywane z okien należących do tego modułu podczas pracy nad projektem. */
export const KOD_MODULU = 'workspace';

/**
 * Klucz ustawienia niosącego instrukcje systemowe projektu. Ta sama wartość
 * stoi po stronie rdzenia w `KluczInstrukcjiProjektu`; kontrakt nie ma dla niej
 * stałej, bo klucze ustawień są danymi katalogu, nie wyliczeniem.
 */
export const KLUCZ_INSTRUKCJI = 'workspace.instrukcje';
