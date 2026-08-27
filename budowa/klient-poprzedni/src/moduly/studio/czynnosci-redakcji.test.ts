import { describe, expect, it } from 'vitest';

import { Command } from '../../../../shared/contract';
import { CZYNNOSCI_REDAKCJI, opiszSkutekRedakcji } from './czynnosci-redakcji';

/** Odnajduje czynność redakcji po komendzie kontraktu, rzucając błąd, gdy katalog czynności jej nie zawiera. */
function czynnosc(komenda: Command) {
  const znaleziona = CZYNNOSCI_REDAKCJI.find((pozycja) => pozycja.komenda === komenda);
  if (znaleziona === undefined) throw new Error(`brak czynności ${komenda} w katalogu`);
  return znaleziona;
}

const OTOCZENIE = { idOkna: 'okno-studio', idDokumentu: 'studio-dok-1' };
const BEZ_DOKUMENTU = { idOkna: 'okno-studio', idDokumentu: null };

describe('katalog czynności redakcji', () => {
  it('niesie wszystkie piętnaście komend, każdą raz', () => {
    const komendy = CZYNNOSCI_REDAKCJI.map((pozycja) => pozycja.komenda);
    expect(komendy).toHaveLength(15);
    expect(new Set(komendy).size).toBe(15);
    expect(komendy.every((komenda) => komenda.startsWith('studio.'))).toBe(true);
  });

  it('daje drogę z okna do każdej z piętnastu komend dobudowanych', () => {
    const oczekiwane = [
      Command.StudioAssetEmbed,
      Command.StudioBatchRun,
      Command.StudioDiffReportExport,
      Command.StudioDiffVisual,
      Command.StudioDiffSource,
      Command.StudioSearchSemantic,
      Command.StudioBranchCreate,
      Command.StudioBranchList,
      Command.StudioBranchMerge,
      Command.StudioRepositoryExport,
      Command.StudioPackageExport,
      Command.StudioVersionReferenceCreate,
      Command.StudioPreviewRender,
      Command.StudioIngestUrl,
      Command.StudioIngestDeviceScan,
    ];
    const komendy = new Set(CZYNNOSCI_REDAKCJI.map((pozycja) => pozycja.komenda));
    for (const komenda of oczekiwane) expect(komendy.has(komenda)).toBe(true);
  });

  it('pomija pola puste zamiast wysyłać je jako pusty napis', () => {
    const zlozenie = czynnosc(Command.StudioPreviewRender).zloz(
      {
        format: 'pdf',
        versionId: '',
        pageFrom: '',
        pageTo: '',
        profileId: '',
        watermark: '  ',
      },
      OTOCZENIE,
    );
    expect('zadanie' in zlozenie).toBe(true);
    if (!('zadanie' in zlozenie)) return;
    expect(zlozenie.zadanie).toEqual({
      documentId: 'studio-dok-1',
      format: 'pdf',
      windowId: 'okno-studio',
    });
  });

  it('nie wysyła fałszu tam, gdzie Operator zostawił wartość domyślną', () => {
    const zlozenie = czynnosc(Command.StudioPackageExport).zloz(
      { includeHistory: '', includeDiffReport: 'nie', includeAnnotations: 'tak' },
      OTOCZENIE,
    );
    expect('zadanie' in zlozenie).toBe(true);
    if (!('zadanie' in zlozenie)) return;
    expect(zlozenie.zadanie).toEqual({
      documentId: 'studio-dok-1',
      windowId: 'okno-studio',
      includeDiffReport: false,
      includeAnnotations: true,
    });
  });

  it('odmawia zdaniem, gdy czynność potrzebuje dokumentu, a edytor jest pusty', () => {
    const zlozenie = czynnosc(Command.StudioRepositoryExport).zloz({}, BEZ_DOKUMENTU);
    expect('odmowa' in zlozenie).toBe(true);
    if (!('odmowa' in zlozenie)) return;
    expect(zlozenie.odmowa).toContain('Studio Editorze');
  });

  it('składa wsad z wykazu dokumentów oddzielonych przecinkami', () => {
    const zlozenie = czynnosc(Command.StudioBatchRun).zloz(
      { documentIds: 'studio-dok-1, studio-dok-2 ,', actionId: 'skroc' },
      OTOCZENIE,
    );
    expect('zadanie' in zlozenie).toBe(true);
    if (!('zadanie' in zlozenie)) return;
    expect(zlozenie.zadanie).toEqual({
      windowId: 'okno-studio',
      documentIds: ['studio-dok-1', 'studio-dok-2'],
      actionId: 'skroc',
    });
  });

  it('dokłada rozstrzygnięcie konfliktu dopiero wtedy, gdy podano numer i stronę', () => {
    const bez = czynnosc(Command.StudioBranchMerge).zloz(
      { sourceBranchId: 'g1', targetBranchId: 'g2', resolutionIndex: '', resolutionSide: '' },
      OTOCZENIE,
    );
    expect('zadanie' in bez).toBe(true);
    if (!('zadanie' in bez)) return;
    expect(bez.zadanie).toEqual({ sourceBranchId: 'g1', targetBranchId: 'g2' });

    const z = czynnosc(Command.StudioBranchMerge).zloz(
      { sourceBranchId: 'g1', targetBranchId: 'g2', resolutionIndex: '1', resolutionSide: 'scalana' },
      OTOCZENIE,
    );
    expect('zadanie' in z).toBe(true);
    if (!('zadanie' in z)) return;
    expect(z.zadanie['resolutions']).toEqual([{ index: 1, side: 'scalana' }]);
  });

  it('skanowanie z urządzenia składa żądanie — odmowę stawia rdzeń, nie okno', () => {
    // Okno nie udaje wiedzy o braku po stronie serwera: wysyła żądanie i pokazuje odmowę rdzenia.
    const zlozenie = czynnosc(Command.StudioIngestDeviceScan).zloz({ deviceId: '' }, BEZ_DOKUMENTU);
    expect('zadanie' in zlozenie).toBe(true);
    if (!('zadanie' in zlozenie)) return;
    expect(zlozenie.zadanie).toEqual({ windowId: 'okno-studio' });
  });
});

describe('opis skutku redakcji', () => {
  it('nazywa drogę wyszukiwania znaczeniowego, a nie samą liczbę trafień', () => {
    const zdanie = opiszSkutekRedakcji({
      matches: [{ text: 'a' }, { text: 'b' }],
      mode: 'miara-rdzenia',
    });
    expect(zdanie).toContain('fragmentów: 2');
    expect(zdanie).toContain('droga: miara-rdzenia');
  });

  it('odróżnia scalenie dokonane od scalenia wstrzymanego konfliktem', () => {
    expect(opiszSkutekRedakcji({ merged: true })).toContain('scalone');
    const wstrzymane = opiszSkutekRedakcji({ merged: false, conflicts: [{ index: 1 }] });
    expect(wstrzymane).toContain('wstrzymane konfliktem');
    expect(wstrzymane).toContain('konfliktów do rozstrzygnięcia: 1');
  });

  it('mówi wprost, gdy materiału wejściowego nie udało się odczytać', () => {
    expect(opiszSkutekRedakcji({ hunks: [], sourceResolved: false })).toContain(
      'nie udało się odczytać',
    );
  });

  it('nie melduje pustego zdania, gdy odpowiedź nie ma żadnego znanego pola', () => {
    expect(opiszSkutekRedakcji({})).toBe('Rdzeń przyjął czynność.');
  });
});
