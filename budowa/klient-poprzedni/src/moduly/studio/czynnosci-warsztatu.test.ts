import { describe, expect, it } from 'vitest';

import { Command } from '../../../../shared/contract';
import { CZYNNOSCI_WARSZTATU, opiszSkutek } from './czynnosci-warsztatu';

/** Odnajduje czynność warsztatu po komendzie kontraktu, rzucając błąd, gdy katalog czynności jej nie zawiera. */
function czynnosc(komenda: Command) {
  const znaleziona = CZYNNOSCI_WARSZTATU.find((pozycja) => pozycja.komenda === komenda);
  if (znaleziona === undefined) throw new Error(`brak czynności ${komenda} w katalogu`);
  return znaleziona;
}

const OTOCZENIE = { idOkna: 'okno-studio', idDokumentu: 'dokument-1' };

describe('katalog czynności warsztatu', () => {
  it('niesie wszystkie piętnaście komend obu rodzin, każdą raz', () => {
    const komendy = CZYNNOSCI_WARSZTATU.map((pozycja) => pozycja.komenda);
    expect(komendy).toHaveLength(15);
    expect(new Set(komendy).size).toBe(15);
    expect(komendy.every((komenda) => komenda.startsWith('studio.'))).toBe(true);
  });

  it('pomija pola puste zamiast wysyłać je jako pusty napis', () => {
    const zlozenie = czynnosc(Command.StudioPdfBates).zloz(
      { assetId: 'zasob-1', prefix: '', startNumber: '5', digits: '', header: '', footer: '' },
      OTOCZENIE,
    );
    expect('zadanie' in zlozenie).toBe(true);
    if (!('zadanie' in zlozenie)) return;
    expect(zlozenie.zadanie).toEqual({
      assetId: 'zasob-1',
      windowId: 'okno-studio',
      startNumber: 5,
    });
  });

  it('odmawia scalania jednego dokumentu, zanim żądanie pójdzie do rdzenia', () => {
    const zlozenie = czynnosc(Command.StudioPdfMerge).zloz({ assetIds: 'zasob-1' }, OTOCZENIE);
    expect(zlozenie).toEqual({ odmowa: 'Scalanie wymaga co najmniej dwóch dokumentów.' });
  });

  it('rozbiera zakładki na tytuł i numer strony', () => {
    const zlozenie = czynnosc(Command.StudioPdfBookmarksSet).zloz(
      { assetId: 'zasob-1', bookmarks: 'Wstęp | 1\nRozdział pierwszy | 4' },
      OTOCZENIE,
    );
    if (!('zadanie' in zlozenie)) throw new Error('złożenie miało się udać');
    expect(zlozenie.zadanie['bookmarks']).toEqual([
      { title: 'Wstęp', page: 1 },
      { title: 'Rozdział pierwszy', page: 4 },
    ]);
  });

  it('nazywa wiersz obszaru redakcji, którego nie da się odczytać', () => {
    const zlozenie = czynnosc(Command.StudioSecurityRedact).zloz(
      { assetId: 'zasob-1', regions: '1 | 72 | 640' },
      OTOCZENIE,
    );
    expect(zlozenie).toHaveProperty('odmowa');
    if (!('odmowa' in zlozenie)) return;
    expect(zlozenie.odmowa).toContain('pięciu liczb');
  });

  it('mówi wprost, że rozpoznanie danych wrażliwych potrzebuje dokumentu', () => {
    const zlozenie = czynnosc(Command.StudioSecuritySensitiveDetect).zloz(
      {},
      { idOkna: 'okno-studio', idDokumentu: null },
    );
    expect(zlozenie).toHaveProperty('odmowa');
  });
});

describe('opis skutku czynności', () => {
  it('opisuje wynik liczbą, a nie słowem „gotowe"', () => {
    expect(opiszSkutek({ asset: { id: 'zasob-9' }, pages: 5 })).toBe('nowy zasób zasob-9, stron: 5.');
    expect(opiszSkutek({ parts: 3, assetIds: ['a', 'b', 'c'] })).toBe('części: 3, zasobów: 3.');
  });

  it('nie udaje wiedzy o odpowiedzi, której nie rozumie', () => {
    expect(opiszSkutek(null)).toBe('Rdzeń przyjął czynność.');
    expect(opiszSkutek({})).toBe('Rdzeń przyjął czynność.');
  });
});
