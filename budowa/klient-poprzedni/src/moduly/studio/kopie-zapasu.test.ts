import { describe, expect, it } from 'vitest';

import {
  StudioBackupReason,
  StudioVersionSeries,
  type StudioAutosaveSettings,
  type StudioDocumentBackup,
} from '../../../../shared/contract';
import {
  kopieNazwaPowodu,
  kopieNazwaSzeregu,
  kopieOpiszKopie,
  kopieOpiszNastawy,
  kopieOpiszStanZapisu,
  kopieOpiszSzeregi,
  kopieStanZapisu,
  kopieZgloszeniePoZamknieciu,
} from './kopie-zapasu';

/**
 * Sprawdziany autozapisu i kopii — mierzą UCZCIWOŚĆ ZAPISU.
 *
 * Wskaźnik „zapisano" pokazany po nieudanym zapisie jest najgorszym możliwym
 * błędem tego modułu: Operator zamknie okno i straci pracę. Sprawdzian mierzy to
 * na zapisie NIEUDANYM, nie tylko na udanym — tak, jak żąda zlecenie.
 */

/** Nastawy autozapisu z niepowodzeniem ostatniego zapisu. */
function nastawyPoNiepowodzeniu(): StudioAutosaveSettings {
  return {
    enabled: true,
    intervalSeconds: 120,
    lastSaveAt: 1_700_000_000_000,
    lastSaveFailed: true,
    lastFailureReason: 'magazyn odmówił zapisu',
  };
}

/** Kopia zapasowa o wskazanym czasie i stanie zmian niezapisanych. */
function kopia(
  kod: string,
  czas: number,
  niezapisane: boolean,
  udana = true,
): StudioDocumentBackup {
  return {
    id: kod,
    documentId: 'pismo',
    reason: StudioBackupReason.Interval,
    bytes: 2048,
    succeeded: udana,
    unsavedChanges: niezapisane,
    createdAt: czas,
  };
}

describe('autozapis — uczciwość wskaźnika zapisu', () => {
  it('niepowodzenie ostatniego zapisu WYGRYWA nad wszystkim, także nad brakiem zmian', () => {
    expect(kopieStanZapisu(nastawyPoNiepowodzeniu(), false, false)).toBe('nieudany');
    expect(kopieStanZapisu(nastawyPoNiepowodzeniu(), false, true)).toBe('nieudany');
  });

  it('zdanie o nieudanym zapisie niesie powód i wskazuje kopię, a nie samo „błąd"', () => {
    const zdanie = kopieOpiszStanZapisu('nieudany', nastawyPoNiepowodzeniu());
    expect(zdanie).toContain('magazyn odmówił zapisu');
    expect(zdanie).toContain('kopii zapasowej');
    expect(zdanie).not.toContain('zapisano ·');
  });

  it('trzy pozostałe stany są rozłączne i wynikają ze zmian oraz zapisu w toku', () => {
    const nastawy: StudioAutosaveSettings = { enabled: true, lastSaveAt: 1 };
    expect(kopieStanZapisu(nastawy, false, false)).toBe('zapisano');
    expect(kopieStanZapisu(nastawy, true, false)).toBe('niezapisane');
    expect(kopieStanZapisu(nastawy, true, true)).toBe('zapisywanie');
  });

  it('brak jakiegokolwiek udanego zapisu jest nazwany, a nie podany jako czas zerowy', () => {
    expect(kopieOpiszStanZapisu('zapisano', null)).toContain('nie zgłosił jeszcze żadnego');
  });

  it('autozapis wyłączony mówi, że jest ustawieniem Operatora, a nie awarią', () => {
    const zdanie = kopieOpiszNastawy({ enabled: false });
    expect(zdanie).toContain('WYŁĄCZONY');
    expect(zdanie).toContain('ustawieniem');
  });

  it('nastawy czynne wypisują odstęp, zdarzenia i zasadę wygasania', () => {
    const zdanie = kopieOpiszNastawy({
      enabled: true,
      intervalSeconds: 60,
      onBlur: true,
      onClose: true,
      backupRetentionCount: 10,
      backupRetentionHours: 72,
    });
    expect(zdanie).toContain('co 60 sekund');
    expect(zdanie).toContain('odejście od okna');
    expect(zdanie).toContain('zachowywanych kopii: 10');
    expect(zdanie).toContain('72 godzinach');
    expect(zdanie).toContain('osobnym szeregiem');
  });
});

describe('kopie zapasowe — zgłoszenie po nagłym zamknięciu', () => {
  it('zgłasza kopię NAJŚWIEŻSZĄ niosącą zmiany niezapisane i mówi, ile ich jest', () => {
    const zgloszenie = kopieZgloszeniePoZamknieciu([
      kopia('starsza', 1_000, true),
      kopia('najswiezsza', 9_000, true),
      kopia('zapisana', 5_000, false),
    ]);
    expect(zgloszenie?.kopia.id).toBe('najswiezsza');
    expect(zgloszenie?.zdanie).toContain('Przywrócić?');
    expect(zgloszenie?.zdanie).toContain('DO NOWEGO DOKUMENTU');
    expect(zgloszenie?.zdanie).toContain('jest 2');
  });

  it('kopia, której zapis się NIE udał, nie jest podstawą zgłoszenia', () => {
    expect(kopieZgloszeniePoZamknieciu([kopia('polamana', 9_000, true, false)])).toBeNull();
  });

  it('brak kopii z niezapisanymi zmianami znaczy „nie ma czego zgłaszać"', () => {
    expect(kopieZgloszeniePoZamknieciu([kopia('zapisana', 1_000, false)])).toBeNull();
  });

  it('kopia nieudana jest w wykazie oznaczona powodem, nie ukryta', () => {
    const polamana: StudioDocumentBackup = {
      ...kopia('polamana', 1_000, true, false),
      failureReason: 'brak miejsca w magazynie',
    };
    const zdanie = kopieOpiszKopie(polamana);
    expect(zdanie).toContain('ZAPIS KOPII SIĘ NIE UDAŁ');
    expect(zdanie).toContain('brak miejsca w magazynie');
    expect(zdanie).toContain('NIESIE ZMIANY NIEZAPISANE');
  });

  it('powód założenia kopii ma pełną nazwę, a nieznany oddaje swoją wartość', () => {
    expect(kopieNazwaPowodu(StudioBackupReason.BeforeIrreversible)).toBe(
      'przed czynnością nieodwracalną',
    );
    expect(kopieNazwaPowodu('powodNieznany')).toBe('powodNieznany');
  });
});

describe('szeregi wersji — autozapis idzie osobno', () => {
  it('zdanie o szeregach podaje oba liczniki i wskazuje wersję założycielską', () => {
    const zdanie = kopieOpiszSzeregi([], 4, 11, 'wersja-pierwsza');
    expect(zdanie).toContain('Szereg Operatora: 4');
    expect(zdanie).toContain('szereg autozapisu: 11');
    expect(zdanie).toContain('wersja-pierwsza');
  });

  it('brak wersji założycielskiej jest nazwany, a nie podstawiony pierwszą z wykazu', () => {
    expect(kopieOpiszSzeregi([], 0, 0, undefined)).toContain('nie wskazał');
  });

  it('wiersz wykazu zbiorczego nie udaje, że zna swój szereg', () => {
    expect(kopieNazwaSzeregu(StudioVersionSeries.Autosave)).toBe('zapis samoczynny');
    expect(kopieNazwaSzeregu(undefined)).toContain('sam szeregu nie niesie');
  });
});
