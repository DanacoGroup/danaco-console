import { describe, expect, it } from 'vitest';

import {
  StudioAgentSlotState,
  StudioAuthor,
  StudioLockScope,
  type StudioActionBalance,
  type StudioAgentSlot,
  type StudioFragmentLock,
} from '../../../../shared/contract';
import {
  blokadaBilansPominiec,
  blokadaNazwaZasiegu,
  blokadaOpisz,
  blokadaOpiszOdmoweZajecia,
  blokadaOpiszWykaz,
  blokadaOpiszZajecie,
  blokadaZakresPoprawny,
} from './blokada-fragmentow';

/** Buduje blokadę fragmentu o wskazanym zasięgu, z ustaloną nazwą i zakresem, do użycia w sprawdzianach opisu blokad. */
function blokada(zasieg: StudioLockScope): StudioFragmentLock {
  return {
    id: 'blokada-1',
    documentId: 'pismo',
    name: 'podstawa prawna',
    reason: 'nie zmieniać',
    rangeStart: 10,
    rangeEnd: 40,
    scope: zasieg,
    createdBy: StudioAuthor.Uzytkownik,
    createdAt: 1_700_000_000_000,
  };
}

/** Buduje zajęcie fragmentu dokumentu przez nazwanego wykonawcę modelowego, do użycia w sprawdzianach opisu zajęć. */
function zajecie(): StudioAgentSlot {
  return {
    id: 'zajecie-1',
    documentId: 'pismo',
    actor: { kind: StudioAuthor.Model, agentId: 'ekspert-1', agentName: 'redaktor' },
    state: StudioAgentSlotState.Working,
    rangeStart: 5,
    rangeEnd: 25,
    expiresAt: 1_700_000_300_000,
  };
}

describe('blokada fragmentu', () => {
  it('zasięg obejmujący Operatora jest opisany jako ustawienie JAWNE, nie domyślne', () => {
    expect(blokadaNazwaZasiegu(StudioLockScope.Everyone)).toContain('jawne');
    expect(blokadaNazwaZasiegu(StudioLockScope.Model)).toContain('bez przeszkód');
  });

  it('opis blokady niesie nazwę, powód, zakres i zasięg naraz', () => {
    const zdanie = blokadaOpisz(blokada(StudioLockScope.Model));
    expect(zdanie).toContain('podstawa prawna');
    expect(zdanie).toContain('nie zmieniać');
    expect(zdanie).toContain('znaki 10–40');
    expect(zdanie).toContain('wyłącznie model');
  });

  it('blokada wzorcowa bez wskazania szablonu mówi, że rdzeń go nie podał', () => {
    const zSzablonu: StudioFragmentLock = {
      ...blokada(StudioLockScope.Model),
      fromTemplate: true,
    };
    expect(blokadaOpisz(zSzablonu)).toContain('KTÓREGO, rdzeń nie podał');
  });

  it('pusty wykaz blokad znaczy „model może zmieniać wszystko", a nie „blokady wyłączone"', () => {
    const zdanie = blokadaOpiszWykaz([]);
    expect(zdanie).toContain('może zmieniać wszystko');
    expect(zdanie).not.toContain('wyłączone.');
  });

  it('wykaz liczy osobno blokady obejmujące Operatora i wzorcowe z szablonu', () => {
    const zdanie = blokadaOpiszWykaz([
      blokada(StudioLockScope.Everyone),
      { ...blokada(StudioLockScope.Model), id: 'blokada-2', fromTemplate: true },
    ]);
    expect(zdanie).toContain('Blokad: 2');
    expect(zdanie).toContain('także Operatora: 1');
    expect(zdanie).toContain('wzorcowych z szablonu: 1');
  });

  it('zakres pusty nie nadaje się na blokadę — fragment o zerowej długości nie jest fragmentem', () => {
    expect(blokadaZakresPoprawny(null)).toBe(false);
    expect(blokadaZakresPoprawny({ poczatek: 7, koniec: 7 })).toBe(false);
    expect(blokadaZakresPoprawny({ poczatek: 7, koniec: 8 })).toBe(true);
  });
});

describe('zmiana obejmująca blokadę częściowo', () => {
  it('oddaje bilans wraz z NAZWĄ blokady i zakresem, zamiast przemilczeć pominięcie', () => {
    const bilans: StudioActionBalance = {
      applied: 12,
      skippedCount: 2,
      skipped: [
        { reason: 'fragment pod blokadą', lockName: 'podstawa prawna', rangeStart: 10, rangeEnd: 40 },
        { reason: 'fragment pod blokadą', lockId: 'blokada-7' },
      ],
    };
    const zdanie = blokadaBilansPominiec(bilans);
    expect(zdanie).toContain('POZA blokadami w 12 miejscach');
    expect(zdanie).toContain('podstawa prawna');
    expect(zdanie).toContain('blokada-7');
  });

  it('bilans bez pominięć przez blokadę nie wymyśla zdania o blokadzie', () => {
    expect(blokadaBilansPominiec({ applied: 3, skippedCount: 0 })).toBeNull();
    expect(
      blokadaBilansPominiec({
        applied: 3,
        skippedCount: 1,
        skipped: [{ reason: 'fragment poza dokumentem' }],
      }),
    ).toBeNull();
  });
});

describe('zajęcie fragmentu przez wykonawcę', () => {
  it('odmowa zajęcia NAZYWA wykonawcę, zakres i czas wygaśnięcia', () => {
    const zdanie = blokadaOpiszOdmoweZajecia(zajecie(), 'fragment zajęty');
    expect(zdanie).toContain('NIE zajęto');
    expect(zdanie).toContain('redaktor');
    expect(zdanie).toContain('znakach 5–25');
  });

  it('odmowa bez wskazania wykonawcy jest zgłoszona jako BRAK, nie przemilczana', () => {
    const zdanie = blokadaOpiszOdmoweZajecia(undefined, '');
    expect(zdanie).toContain('brak do zgłoszenia');
  });

  it('opis zajęcia podaje czas wygaśnięcia, bo wykonawca ubity w pół pracy nie trzyma fragmentu na zawsze', () => {
    expect(blokadaOpiszZajecie(zajecie())).toContain('wygasa');
    const bezCzasu: StudioAgentSlot = { ...zajecie(), expiresAt: undefined };
    expect(blokadaOpiszZajecie(bezCzasu)).toContain('nie podał');
  });
});
