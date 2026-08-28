import { describe, expect, it } from 'vitest';

import {
  AutomationDependencyKind,
  AutomationExecutionStatus,
  AutomationStepKind,
  AutomationTriggerKind,
  QueueAction,
  type AutomationExecution,
  type AutomationStep,
} from '../../../../shared/contract';
import { RODZAJE_KROKU } from './edytor-krokow';
import { DZIALANIA } from './okno-queue-manager';
import { RODZAJE_WYZWALACZA } from './okno-scheduler';
import { zapisDot, zapisMermaid } from './eksport-grafu';
import { warstwyGrafu, type OpisGrafu } from './graf-krokow';
import { dniKalendarza, ZAKRESY_KALENDARZA } from './kalendarz-uruchomien';
import { miaryPrzebiegow, zawezonePrzebiegi, BEZ_ZAWEZENIA } from './przeglad-przebiegow';
import { scenariuszZKompletu } from './przyjecie-przekazania';
import { zastrzezeniaDefinicji } from './walidacja-definicji';
import {
  opisCyklicznosci,
  rozpoznajWzorzec,
  WZORCE_CYKLICZNOSCI,
  zapisWzorca,
  NASTAWY_WYJSCIOWE,
} from './wzorce-cyklicznosci';

/** Sprawdziany warstwy klienckiej modułu Automations obejmują wyłącznie byty rozstrzygalne bez rdzenia. */

/** Krok przykładowy w kształcie, w którym niesie go kontrakt, gotowy do nadpisania wybranych pól testu. */
function krok(id: string, dodatki: Partial<AutomationStep> = {}): AutomationStep {
  return { id, kind: AutomationStepKind.Command, command: 'przyklad.krok', ...dodatki };
}

/** Przebieg przykładowy, którego czas podawany jest w milisekundach liczonych od wspólnego punktu odniesienia. */
const POCZATEK = Date.UTC(2026, 7, 17, 7, 0, 0);

function przebieg(
  id: string,
  status: string,
  odsuniecie: number,
  trwanie?: number,
): AutomationExecution {
  const wpis: AutomationExecution = {
    id,
    workflowId: 'raport-tygodniowy',
    status: status as AutomationExecution['status'],
    startedAt: POCZATEK + odsuniecie,
  };
  if (trwanie !== undefined) wpis.finishedAt = POCZATEK + odsuniecie + trwanie;
  return wpis;
}

/**
 * Wykazy rozwijane modułu biorą się z wyliczeń kontraktu, nie z zapisu w oknie, zaporą na
 * powtórzenie usterki.
 */
describe('wykazy modułu odpowiadają wyliczeniom kontraktu', () => {
  it('rodzaje kroku pokrywają wyliczenie rodzajów kroku', () => {
    expect(RODZAJE_KROKU.map(([wartosc]) => wartosc)).toEqual(Object.values(AutomationStepKind));
  });

  it('rodzaje wyzwalacza pokrywają wyliczenie rodzajów wyzwalacza', () => {
    expect(RODZAJE_WYZWALACZA.map(([wartosc]) => wartosc)).toEqual(
      Object.values(AutomationTriggerKind),
    );
  });

  it('działania kolejki pokrywają wyliczenie działań silnika', () => {
    expect(DZIALANIA.map(([wartosc]) => wartosc)).toEqual(Object.values(QueueAction));
  });

  it('każda pozycja wykazu ma nazwę na ekranie, nie samą wartość kontraktu', () => {
    for (const [wartosc, nazwa] of [...RODZAJE_KROKU, ...RODZAJE_WYZWALACZA, ...DZIALANIA]) {
      expect(nazwa.trim(), `nazwa pozycji ${wartosc}`).not.toBe('');
    }
  });
});

describe('wzorce cykliczności', () => {
  it('składa zapis cron z wzorca i nastaw', () => {
    expect(
      zapisWzorca(WZORCE_CYKLICZNOSCI.coTydzien, { ...NASTAWY_WYJSCIOWE, godzina: 7, minuta: 0, dzienTygodnia: 1 }),
    ).toBe('0 7 * * 1');
    expect(zapisWzorca(WZORCE_CYKLICZNOSCI.coMinute, NASTAWY_WYJSCIOWE)).toBe('* * * * *');
    expect(
      zapisWzorca(WZORCE_CYKLICZNOSCI.interwalMinutowy, { ...NASTAWY_WYJSCIOWE, krok: 15 }),
    ).toBe('*/15 * * * *');
    expect(
      zapisWzorca(WZORCE_CYKLICZNOSCI.kwartalnie, { ...NASTAWY_WYJSCIOWE, dzienMiesiaca: 1 }),
    ).toBe('0 7 1 1,4,7,10 *');
  });

  it('wtłacza nastawę spoza zakresu pola w jego granice, zamiast wysyłać ją dalej', () => {
    expect(zapisWzorca(WZORCE_CYKLICZNOSCI.codziennie, { ...NASTAWY_WYJSCIOWE, godzina: 99 })).toBe(
      '0 23 * * *',
    );
  });

  it('rozpoznaje wzorzec w zapisie odczytanym z rdzenia', () => {
    expect(rozpoznajWzorzec('0 7 * * 1')).toEqual({
      wzorzec: WZORCE_CYKLICZNOSCI.coTydzien,
      nastawy: { ...NASTAWY_WYJSCIOWE, minuta: 0, godzina: 7, dzienTygodnia: 1, dzienMiesiaca: 1 },
    });
    expect(rozpoznajWzorzec('*/15 * * * *').wzorzec).toBe(WZORCE_CYKLICZNOSCI.interwalMinutowy);
    expect(rozpoznajWzorzec('0 7 1 1,4,7,10 *').wzorzec).toBe(WZORCE_CYKLICZNOSCI.kwartalnie);
  });

  it('zapis, którego żaden wzorzec nie obejmuje, jest wzorcem własnym, nie usterką', () => {
    expect(rozpoznajWzorzec('0 6 * * 1-5').wzorzec).toBe(WZORCE_CYKLICZNOSCI.wlasny);
    expect(rozpoznajWzorzec('nie jest zapisem').wzorzec).toBe(WZORCE_CYKLICZNOSCI.wlasny);
  });

  it('przekład w obie strony nie gubi zapisu', () => {
    for (const zapis of ['* * * * *', '30 * * * *', '0 7 * * *', '0 7 * * 3', '0 7 12 * *']) {
      const rozpoznany = rozpoznajWzorzec(zapis);
      expect(zapisWzorca(rozpoznany.wzorzec, rozpoznany.nastawy), zapis).toBe(zapis);
    }
  });

  it('opisuje cykliczność mową Operatora', () => {
    expect(opisCyklicznosci('0 7 * * 1')).toContain('poniedziałek');
    expect(opisCyklicznosci('')).toContain('Bez cykliczności');
  });
});

describe('walidacja definicji', () => {
  it('definicja bez zastrzeżeń oddaje wykaz pusty', () => {
    expect(zastrzezeniaDefinicji([krok('zbierz'), krok('wyslij')])).toEqual([]);
  });

  it('nazywa krok bez identyfikatora i identyfikator powielony', () => {
    const zastrzezenia = zastrzezeniaDefinicji([krok(''), krok('zbierz'), krok('zbierz')]);
    expect(zastrzezenia.filter((wpis) => wpis.waga === 'blad')).toHaveLength(3);
  });

  it('nazywa krok komendy bez komendy i krok rozstrzygający bez warunku', () => {
    const zastrzezenia = zastrzezeniaDefinicji([
      krok('zbierz', { command: '' }),
      krok('sprawdz', { kind: AutomationStepKind.Condition }),
    ]);
    expect(zastrzezenia.filter((wpis) => wpis.waga === 'ostrzezenie')).toHaveLength(2);
  });

  it('nazywa poprzednika, którego w definicji nie ma', () => {
    const zastrzezenia = zastrzezeniaDefinicji([krok('wyslij', { dependsOn: ['zbierz'] })]);
    expect(zastrzezenia).toHaveLength(1);
    expect(zastrzezenia[0]?.zdanie).toContain('zbierz');
  });

  it('wskazuje wszystkie kroki cyklu, nie jeden wybrany', () => {
    const zastrzezenia = zastrzezeniaDefinicji([
      krok('pierwszy', { dependsOn: ['drugi'] }),
      krok('drugi', { dependsOn: ['pierwszy'] }),
    ]);
    const wCyklu = zastrzezenia.filter((wpis) => wpis.zdanie.includes('cyklu'));
    expect(wCyklu.map((wpis) => wpis.idKroku).sort()).toEqual(['drugi', 'pierwszy']);
  });

  it('nie zgłasza kroku bez połączenia, gdy definicja nie używa zależności wcale', () => {
    const zastrzezenia = zastrzezeniaDefinicji([krok('zbierz'), krok('wyslij')]);
    expect(zastrzezenia.filter((wpis) => wpis.zdanie.includes('bez połączenia'))).toEqual([]);
  });
});

describe('przegląd przebiegów', () => {
  const wykaz = [
    przebieg('przebieg-1', AutomationExecutionStatus.Succeeded, 0, 60_000),
    przebieg('przebieg-2', AutomationExecutionStatus.Failed, 86_400_000, 20_000),
    przebieg('przebieg-3', AutomationExecutionStatus.Succeeded, 172_800_000, 40_000),
    przebieg('przebieg-4', AutomationExecutionStatus.Running, 259_200_000),
  ];

  it('liczy miary z pól, które kontrakt naprawdę niesie', () => {
    const miary = miaryPrzebiegow(wykaz);
    expect(miary.liczba).toBe(4);
    expect(miary.udane).toBe(2);
    expect(miary.nieudane).toBe(1);
    expect(miary.wToku).toBe(1);
    expect(miary.wskaznikPowodzenia).toBe(67);
    expect(miary.sredniCzasMs).toBe(40_000);
    expect(miary.opoznienieP95Ms).toBe(60_000);
  });

  it('wykaz pusty daje miary puste, a nie zera udające pomiar', () => {
    const miary = miaryPrzebiegow([]);
    expect(miary.wskaznikPowodzenia).toBeNull();
    expect(miary.sredniCzasMs).toBeNull();
    expect(miary.opoznienieP95Ms).toBeNull();
  });

  it('zawęża wykaz po stanie i po napisie szukanym', () => {
    expect(
      zawezonePrzebiegi(wykaz, { ...BEZ_ZAWEZENIA, stan: AutomationExecutionStatus.Succeeded }),
    ).toHaveLength(2);
    expect(zawezonePrzebiegi(wykaz, { ...BEZ_ZAWEZENIA, szukane: 'przebieg-2' })).toHaveLength(1);
  });

  it('zawężenie puste oddaje wykaz w całości', () => {
    expect(zawezonePrzebiegi(wykaz, BEZ_ZAWEZENIA)).toHaveLength(wykaz.length);
  });
});

describe('graf kroków', () => {
  const opis: OpisGrafu = {
    wezly: [
      { kod: 'zbierz', nazwa: 'Zbierz dane' },
      { kod: 'sprawdz', nazwa: 'Waliduj wynik' },
      { kod: 'wyslij', nazwa: 'Wyślij raport' },
    ],
    krawedzie: [
      { odKroku: 'zbierz', doKroku: 'sprawdz', podpis: AutomationDependencyKind.Sequential },
      { odKroku: 'sprawdz', doKroku: 'wyslij', podpis: AutomationDependencyKind.Sequential },
    ],
    sciezkaKrytyczna: new Set(['zbierz', 'sprawdz', 'wyslij']),
  };

  it('układa węzły w warstwy według najdłuższej drogi od kroku bez poprzednika', () => {
    expect(warstwyGrafu(opis)).toEqual([['zbierz'], ['sprawdz'], ['wyslij']]);
  });

  it('układ z cyklem nadal daje warstwy, bo cykl jest zastrzeżeniem, nie powodem pustki', () => {
    const zCyklem: OpisGrafu = {
      wezly: [{ kod: 'pierwszy', nazwa: '' }, { kod: 'drugi', nazwa: '' }],
      krawedzie: [
        { odKroku: 'pierwszy', doKroku: 'drugi', podpis: '' },
        { odKroku: 'drugi', doKroku: 'pierwszy', podpis: '' },
      ],
      sciezkaKrytyczna: new Set(),
    };
    expect(warstwyGrafu(zCyklem).flat().sort()).toEqual(['drugi', 'pierwszy']);
  });

  it('składa zapis DOT wraz z podpisami krawędzi', () => {
    const dot = zapisDot(opis);
    expect(dot).toContain('digraph');
    expect(dot).toContain('"zbierz" -> "sprawdz"');
    expect(dot).toContain('label="Zbierz dane"');
  });

  it('składa zapis Mermaid z zastępnikami identyfikatorów', () => {
    const mermaid = zapisMermaid(opis);
    expect(mermaid).toContain('flowchart LR');
    expect(mermaid).toContain('krok0["Zbierz dane"]');
    expect(mermaid).toContain('krok0 -- "sequential" --> krok1');
  });
});

describe('kalendarz uruchomień', () => {
  it('rozkłada terminy zaplanowane i przebiegi odbyte na dni', () => {
    const od = new Date(2026, 7, 17, 6, 0, 0);
    const dni = dniKalendarza(
      '0 7 * * *',
      [{ ...przebieg('przebieg-1', AutomationExecutionStatus.Succeeded, 0, 1000), startedAt: od.getTime() }],
      ZAKRESY_KALENDARZA.tydzien,
      od,
    );
    expect(dni).toHaveLength(7);
    expect(dni.every((dzien) => dzien.zaplanowane === 1)).toBe(true);
    expect(dni[0]?.odbyte).toBe(1);
  });

  it('cykliczność nieczytelna daje kalendarz bez terminów, a nie brak kalendarza', () => {
    const dni = dniKalendarza('nie jest zapisem', [], ZAKRESY_KALENDARZA.miesiac);
    expect(dni).toHaveLength(28);
    expect(dni.every((dzien) => dzien.zaplanowane === 0)).toBe(true);
  });
});

describe('przyjęcie przekazania z innego modułu', () => {
  const komplet = {
    name: 'Scenariusz przeglądania PRZYKŁAD',
    description: 'zebranie kursów z witryny dostawcy',
    enabled: true,
    steps: [{ id: 'otworz', kind: AutomationStepKind.Command }],
  };

  it('wyjmuje scenariusz z parametrów wykonania przeniesionego kompletu', () => {
    const scenariusz = scenariuszZKompletu('okno-przykladowe', 'polecenie wyjściowe', komplet);
    expect(scenariusz?.nazwa).toBe('Scenariusz przeglądania PRZYKŁAD');
    expect(scenariusz?.kroki).toHaveLength(1);
    expect(scenariusz?.czynny).toBe(true);
  });

  it('pomija komplet bez nazwy, bez kroków i o kształcie obcym', () => {
    expect(scenariuszZKompletu('okno-przykladowe', '', { ...komplet, name: '' })).toBeNull();
    expect(scenariuszZKompletu('okno-przykladowe', '', { ...komplet, steps: [] })).toBeNull();
    expect(scenariuszZKompletu('okno-przykladowe', '', 'napis zamiast kompletu')).toBeNull();
    expect(scenariuszZKompletu('okno-przykladowe', '', undefined)).toBeNull();
  });

  it('pomija pozycje wykazu kroków bez identyfikatora albo bez rodzaju', () => {
    const scenariusz = scenariuszZKompletu('okno-przykladowe', '', {
      ...komplet,
      steps: [{ id: '', kind: AutomationStepKind.Command }, { id: 'otworz' }, komplet.steps[0]],
    });
    expect(scenariusz?.kroki).toHaveLength(1);
  });
});
