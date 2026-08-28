import { describe, expect, it } from 'vitest';

import {
  StudioActionKind,
  StudioActionState,
  StudioAuthor,
  type StudioActionBalance,
  type StudioDocumentAction,
  type StudioFormDiffEntry,
  type StudioModelChangeSummary,
} from '../../../../shared/contract';
import {
  DZIENNIK_ROZNICA_DRZEW,
  dziennikCzyCofnieciePowstrzymane,
  dziennikNazwaRodzaju,
  dziennikNazwaWykonawcy,
  dziennikOpiszBilans,
  dziennikOpiszCofniecie,
  dziennikOpiszCofniecieModelu,
  dziennikOpiszRoznicePostaci,
  dziennikOpiszZaleznosci,
  dziennikOpiszZmianyModelu,
  dziennikPoKolejnosci,
} from './dziennik-czynnosci';

/** Buduje wpis dziennika czynności o wskazanym kodzie, kolejności, stanie i zależnościach blokujących go zadania. */
function wpis(
  kod: string,
  kolejnosc: number,
  stan: StudioActionState,
  stojaNaNiej: readonly string[] = [],
): StudioDocumentAction {
  return {
    id: kod,
    documentId: 'pismo',
    kind: StudioActionKind.FormatChange,
    author: StudioAuthor.Uzytkownik,
    description: '',
    state: stan,
    sequence: kolejnosc,
    createdAt: 1_700_000_000_000,
    ...(stojaNaNiej.length === 0 ? {} : { blockedBy: [...stojaNaNiej] }),
  };
}

/** Buduje bilans czynności z jednym pominięciem, zatrzymanym przez nazwaną blokadę fragmentu dokumentu. */
function bilansZBlokada(): StudioActionBalance {
  return {
    applied: 3,
    skippedCount: 1,
    skipped: [
      {
        reason: 'fragment pod blokadą',
        lockName: 'podstawa prawna',
        rangeStart: 10,
        rangeEnd: 40,
      },
    ],
  };
}

describe('dziennik czynności — cofanie i jego odmowa', () => {
  it('odmowę cofnięcia rozpoznaje po pustym wykazie cofniętych, nie po błędzie wywołania', () => {
    const ocena = dziennikOpiszCofniecie([], {
      applied: 0,
      skippedCount: 1,
      skipped: [{ reason: 'na tej czynności stoi czynność późniejsza' }],
    });
    expect(ocena.udane).toBe(false);
    expect(ocena.zdanie).toContain('nie cofnął ani jednej czynności');
    expect(ocena.zdanie).toContain('stoi czynność późniejsza');
  });

  it('każde zdanie o cofnięciu niesie zasadę różnicy drzew, także przy powodzeniu', () => {
    const udane = dziennikOpiszCofniecie(['czynnosc-1'], { applied: 1, skippedCount: 0 });
    expect(udane.udane).toBe(true);
    expect(udane.zdanie).toContain(DZIENNIK_ROZNICA_DRZEW);
    const odmowa = dziennikOpiszCofniecie([], { applied: 0, skippedCount: 0 });
    expect(odmowa.zdanie).toContain(DZIENNIK_ROZNICA_DRZEW);
  });

  it('zależność jest nazwana kodami czynności, a nie samą liczbą', () => {
    const zablokowany = wpis('czynnosc-2', 2, StudioActionState.Active, ['czynnosc-7', 'czynnosc-9']);
    expect(dziennikCzyCofnieciePowstrzymane(zablokowany)).toBe(true);
    const zdanie = dziennikOpiszZaleznosci(zablokowany);
    expect(zdanie).toContain('czynnosc-7');
    expect(zdanie).toContain('czynnosc-9');
    expect(zdanie).toContain('ODMÓWI');
  });

  it('czynność bez zależności mówi to wprost, zamiast milczeć', () => {
    const wolny = wpis('czynnosc-3', 3, StudioActionState.Active);
    expect(dziennikCzyCofnieciePowstrzymane(wolny)).toBe(false);
    expect(dziennikOpiszZaleznosci(wolny)).toContain('nie ma zależności');
  });

  it('wpisy idą od najświeższego, bo tak się cofa', () => {
    const kolejnosc = dziennikPoKolejnosci([
      wpis('pierwsza', 1, StudioActionState.Active),
      wpis('trzecia', 3, StudioActionState.Active),
      wpis('druga', 2, StudioActionState.Reverted),
    ]).map((pozycja) => pozycja.id);
    expect(kolejnosc).toStrictEqual(['trzecia', 'druga', 'pierwsza']);
  });
});

describe('dziennik czynności — bilans nie przemilcza pominięcia', () => {
  it('pominięcie wypisuje blokadę po NAZWIE i zakres, nie samą liczbę', () => {
    const zdanie = dziennikOpiszBilans(bilansZBlokada());
    expect(zdanie).toContain('podstawa prawna');
    expect(zdanie).toContain('znaki 10–40');
    expect(zdanie).toContain('stanęło w 1');
  });

  it('zmiany odłożone po spięciu wykonawców mówią, że brzmienie nie przepadło', () => {
    const zdanie = dziennikOpiszBilans({ applied: 2, skippedCount: 0, deferredCount: 1 });
    expect(zdanie).toContain('NIE przepadło');
  });
});

describe('dziennik czynności — zmiany modelu', () => {
  it('cofnięcie pracy modelu podaje liczbę zmian Operatora ZACHOWANYCH', () => {
    const zdanie = dziennikOpiszCofniecieModelu(4, 2, { applied: 4, skippedCount: 0 }, 'kopia-1');
    expect(zdanie).toContain('ZACHOWANO: 2');
    expect(zdanie).toContain('nie jest przywrócenie wersji sprzed');
    expect(zdanie).toContain('kopia-1');
  });

  it('brak kopii zapasowej przy cofaniu jest nazwany, a nie pominięty', () => {
    const zdanie = dziennikOpiszCofniecieModelu(1, 0, { applied: 1, skippedCount: 0 }, undefined);
    expect(zdanie).toContain('brak do sprawdzenia');
  });

  it('zestawienie rozbija zmiany po wykonawcy, więc dwóch agentów nie zlewa się w jednego', () => {
    const zestawienie: StudioModelChangeSummary = {
      documentId: 'pismo',
      total: 5,
      contentChanges: 3,
      formChanges: 2,
      openChanges: 4,
      byAgent: [
        { actor: { kind: StudioAuthor.Model, agentName: 'redaktor' }, total: 3, contentChanges: 3, formChanges: 0 },
        { actor: { kind: StudioAuthor.Model, agentName: 'zdun' }, total: 2, contentChanges: 0, formChanges: 2 },
      ],
    };
    const zdanie = dziennikOpiszZmianyModelu(zestawienie);
    expect(zdanie).toContain('redaktor — 3');
    expect(zdanie).toContain('zdun — 2');
    expect(zdanie).toContain('postaci 2');
  });

  it('brak zmian modelu ma własne zdanie, nie samo zero', () => {
    const zdanie = dziennikOpiszZmianyModelu({
      documentId: 'pismo',
      total: 0,
      contentChanges: 0,
      formChanges: 0,
      openChanges: 0,
    });
    expect(zdanie).toContain('ani jednej zmiany');
  });

  it('wykonawca nienazwany jest nazwany brakiem, a nie podstawiony Operatorem', () => {
    const czynnoscModelu: StudioDocumentAction = {
      ...wpis('czynnosc-4', 4, StudioActionState.Active),
      author: StudioAuthor.Model,
    };
    expect(dziennikNazwaWykonawcy(czynnoscModelu)).toContain('nienazwany');
    expect(dziennikNazwaWykonawcy(wpis('czynnosc-5', 5, StudioActionState.Active))).toBe('Operator');
  });
});

describe('dziennik czynności — nazwy i różnica postaci', () => {
  it('rodzaj czynności ma pełną nazwę, a nieznany oddaje swoją wartość', () => {
    expect(dziennikNazwaRodzaju(StudioActionKind.StyleChange)).toBe('zmiana stylu nazwanego');
    expect(dziennikNazwaRodzaju('rodzajNieznany')).toBe('rodzajNieznany');
  });

  it('brak różnicy postaci jest odpowiedzią, nie pustką', () => {
    expect(dziennikOpiszRoznicePostaci([], { added: 0, removed: 0, changed: 0 })).toContain(
      'ta sama',
    );
  });

  it('różnica postaci podaje liczby cech, bo zmiana kroju nie ma milczeć', () => {
    const wpisy: StudioFormDiffEntry[] = [
      { kind: 'changed', area: 'characterStyle', detail: 'krój pisma', before: 'Times', after: 'Arial' },
    ];
    const zdanie = dziennikOpiszRoznicePostaci(wpisy, { added: 0, removed: 0, changed: 1 });
    expect(zdanie).toContain('Różnic postaci: 1');
    expect(zdanie).toContain('zmieniło się 1');
  });
});
