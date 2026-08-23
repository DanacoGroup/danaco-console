import { beforeEach, describe, expect, it } from 'vitest';
import { WyciszenieCzasowe } from './progi-aod';
import { PowodDecyzji } from './rozpoznanie-decyzji';
import { WagaUjawnienia } from './rodzaje-sugestii';
import { TrybObecnosci, utworzStanObecnosci } from './tryb-obecnosci';
import {
  KLASA_POWODU,
  KlasaZdarzen,
  ZakresKontekstu,
  utworzStanWyciszen,
  wyciszenieObejmujace,
  zdanieWyciszenia,
  type MagazynWyciszen,
} from './wyciszenie-aod';
import { brakiCzynne, zdanieBraku, zdanieGranicyWyciszenia } from './wyciszenie-braki-kontraktu';

/**
 * Sprawdziany wyciszania Always On Display — czynność z rozdz. 3.5 opracowania
 * `docs/funkcje-globalne/always-on-display.md`, nie kształt pliku.
 *
 * Pilnowane są rzeczy, których zlecenie żąda wprost:
 *   1. trzy rodzaje wyciszenia wstrzymują to, co mają wstrzymywać, i NIE
 *      wstrzymują niczego więcej (kontekstowe — tylko wskazany byt, klasy —
 *      tylko wskazaną klasę);
 *   2. wyjątek wagi krytycznej przechodzi przez WSZYSTKIE rodzaje wyciszenia
 *      i przez tryb cichy — plakietką, bez dymka;
 *   3. wyciszenie i jego zniesienie idą jednym ruchem;
 *   4. magazyn jest podawany, nie brany z globalnej przestrzeni na sztywno;
 *   5. brak pozycji kontraktu jest nazwany wprost, nie zasłonięty.
 */

/** Magazyn na mapie — dowód, że stan nie sięga po `localStorage` na sztywno. */
function magazynPamieciowy(): MagazynWyciszen & { zapis: Map<string, string> } {
  const zapis = new Map<string, string>();
  return {
    zapis,
    getItem: (klucz) => zapis.get(klucz) ?? null,
    setItem: (klucz, wartosc) => void zapis.set(klucz, wartosc),
  };
}

const TERAZ = 1_700_000_000_000;

describe('wykaz wyciszeń czynnych', () => {
  it('wycisza czasowo i znosi jednym ruchem', () => {
    const wyciszenia = utworzStanWyciszen(magazynPamieciowy());

    wyciszenia.wyciszCzasem(WyciszenieCzasowe.Godzina, TERAZ);
    const czynne = wyciszenia.czasowe(TERAZ);
    expect(czynne?.doChwili).toBe(TERAZ + 3_600_000);

    wyciszenia.znies('czasowe');
    expect(wyciszenia.czynne(TERAZ)).toHaveLength(0);
  });

  it('drugi czas zastępuje pierwszy — Operator ma jedno „do kiedy"', () => {
    const wyciszenia = utworzStanWyciszen(magazynPamieciowy());
    wyciszenia.wyciszCzasem(WyciszenieCzasowe.Kwadrans, TERAZ);
    wyciszenia.wyciszCzasem(WyciszenieCzasowe.Godzina, TERAZ);

    expect(wyciszenia.czynne(TERAZ)).toHaveLength(1);
    expect(wyciszenia.czasowe(TERAZ)?.czas).toBe(WyciszenieCzasowe.Godzina);
  });

  it('wyciszenie czasowe przeterminowane znosi się samo', () => {
    const wyciszenia = utworzStanWyciszen(magazynPamieciowy());
    wyciszenia.wyciszCzasem(WyciszenieCzasowe.Kwadrans, TERAZ);

    expect(wyciszenia.czynne(TERAZ + 16 * 60_000)).toHaveLength(0);
  });

  it('trzy rodzaje stoją obok siebie i każdy znosi się osobno', () => {
    const wyciszenia = utworzStanWyciszen(magazynPamieciowy());
    wyciszenia.wyciszCzasem(WyciszenieCzasowe.Godzina, TERAZ);
    wyciszenia.przelaczKontekst(ZakresKontekstu.Modul, 'mod_studio', 'Studio');
    wyciszenia.przelaczKlase(KlasaZdarzen.StanKolejkiZadan);

    expect(wyciszenia.czynne(TERAZ)).toHaveLength(3);

    wyciszenia.znies('kontekstowe:modul:mod_studio');
    expect(wyciszenia.czyKontekstWyciszony(ZakresKontekstu.Modul, 'mod_studio')).toBe(false);
    expect(wyciszenia.czyKlasaWyciszona(KlasaZdarzen.StanKolejkiZadan)).toBe(true);

    wyciszenia.zniesWszystkie();
    expect(wyciszenia.czynne(TERAZ)).toHaveLength(0);
  });

  it('zapisuje się w PODANYM magazynie i odtwarza z niego', () => {
    const magazyn = magazynPamieciowy();
    utworzStanWyciszen(magazyn).przelaczKlase(KlasaZdarzen.Harmonogram);

    expect([...magazyn.zapis.keys()]).toEqual(['danaco.aod.wyciszenia']);
    expect(utworzStanWyciszen(magazyn).czyKlasaWyciszona(KlasaZdarzen.Harmonogram)).toBe(true);
  });

  it('magazyn `null` znaczy „bez zapisu", a nie awarię', () => {
    const wyciszenia = utworzStanWyciszen(null);
    wyciszenia.przelaczKlase(KlasaZdarzen.Harmonogram);
    expect(wyciszenia.czyKlasaWyciszona(KlasaZdarzen.Harmonogram)).toBe(true);
    expect(utworzStanWyciszen(null).czyKlasaWyciszona(KlasaZdarzen.Harmonogram)).toBe(false);
  });

  it('zapis uszkodzony daje wykaz pusty, nie ciszę bez podstawy', () => {
    const magazyn = magazynPamieciowy();
    magazyn.zapis.set('danaco.aod.wyciszenia', '{to nie jest wykaz}');
    expect(utworzStanWyciszen(magazyn).czynne(TERAZ)).toHaveLength(0);
  });

  it('pozycja nieznanego rodzaju nie wycisza niczego', () => {
    const magazyn = magazynPamieciowy();
    magazyn.zapis.set(
      'danaco.aod.wyciszenia',
      JSON.stringify([{ rodzaj: 'wymyslone' }, { rodzaj: 'klasa-zdarzen', klasa: 'nieznana' }]),
    );
    expect(utworzStanWyciszen(magazyn).czynne(TERAZ)).toHaveLength(0);
  });
});

describe('zakres każdego rodzaju wyciszenia', () => {
  it('wyciszenie czasowe obejmuje wszystko', () => {
    const wyciszenia = utworzStanWyciszen(null);
    wyciszenia.wyciszCzasem(WyciszenieCzasowe.Godzina, TERAZ);

    expect(
      wyciszenieObejmujace(wyciszenia.czynne(TERAZ), { klasa: KlasaZdarzen.Harmonogram }),
    ).not.toBeNull();
  });

  it('wyciszenie kontekstowe obejmuje wskazany moduł, a nie sąsiedni', () => {
    const wyciszenia = utworzStanWyciszen(null);
    wyciszenia.przelaczKontekst(ZakresKontekstu.Modul, 'mod_studio', 'Studio');

    expect(wyciszenieObejmujace(wyciszenia.czynne(TERAZ), { idModulu: 'mod_studio' })).not.toBeNull();
    expect(wyciszenieObejmujace(wyciszenia.czynne(TERAZ), { idModulu: 'mod_terminal' })).toBeNull();
  });

  it('wyciszenie karty sesji nie obejmuje innej karty', () => {
    const wyciszenia = utworzStanWyciszen(null);
    wyciszenia.przelaczKontekst(ZakresKontekstu.KartaSesji, 'ses_1', 'Praca nad ofertą');

    expect(wyciszenieObejmujace(wyciszenia.czynne(TERAZ), { idSesji: 'ses_1' })).not.toBeNull();
    expect(wyciszenieObejmujace(wyciszenia.czynne(TERAZ), { idSesji: 'ses_2' })).toBeNull();
  });

  it('sugestia bez rozpoznanego modułu NIE wpada w wyciszenie modułu', () => {
    const wyciszenia = utworzStanWyciszen(null);
    wyciszenia.przelaczKontekst(ZakresKontekstu.Modul, 'mod_studio', 'Studio');

    expect(wyciszenieObejmujace(wyciszenia.czynne(TERAZ), { idSesji: 'ses_1' })).toBeNull();
  });

  it('wyciszenie klasy wstrzymuje jedną klasę, pozostałe pracują', () => {
    const wyciszenia = utworzStanWyciszen(null);
    wyciszenia.przelaczKlase(KlasaZdarzen.StanKolejkiZadan);

    expect(
      wyciszenieObejmujace(wyciszenia.czynne(TERAZ), { klasa: KlasaZdarzen.StanKolejkiZadan }),
    ).not.toBeNull();
    expect(
      wyciszenieObejmujace(wyciszenia.czynne(TERAZ), { klasa: KlasaZdarzen.StanPetliWykonawczej }),
    ).toBeNull();
  });

  it('każdy powód rozpoznania ma klasę zdarzeń z rozdz. 3.2', () => {
    for (const powod of Object.values(PowodDecyzji)) {
      expect(Object.values(KlasaZdarzen)).toContain(KLASA_POWODU[powod]);
    }
    expect(KLASA_POWODU[PowodDecyzji.Wstrzymany]).toBe(KlasaZdarzen.StanPetliWykonawczej);
    expect(KLASA_POWODU[PowodDecyzji.BezRuchu]).toBe(KlasaZdarzen.StanKolejkiZadan);
  });

  it('zdanie wyciszenia mówi CO jest wyciszone i DO KIEDY', () => {
    expect(
      zdanieWyciszenia({ rodzaj: 'czasowe', czas: WyciszenieCzasowe.Godzina, doChwili: TERAZ }),
    ).toContain(new Date(TERAZ).toLocaleTimeString());
    expect(
      zdanieWyciszenia({
        rodzaj: 'kontekstowe',
        zakres: ZakresKontekstu.Modul,
        wartosc: 'mod_studio',
        nazwa: 'Studio',
      }),
    ).toContain('Studio');
    expect(zdanieWyciszenia({ rodzaj: 'klasa-zdarzen', klasa: KlasaZdarzen.Harmonogram })).toContain(
      'harmonogram',
    );
  });
});

describe('reguła ujawniania wobec trzech rodzajów wyciszenia', () => {
  let stan: ReturnType<typeof utworzStanObecnosci>;

  beforeEach(() => {
    stan = utworzStanObecnosci({ magazyn: null });
  });

  const wysoka = { waga: WagaUjawnienia.Wysoka } as const;

  it('bez wyciszenia sugestia wagi wysokiej otwiera dymek', () => {
    expect(stan.czyOtworzycDymek({ ...wysoka }, TERAZ)).toBe(true);
  });

  it('wyciszenie kontekstowe wstrzymuje sugestię wskazanego modułu', () => {
    stan.wyciszenia.przelaczKontekst(ZakresKontekstu.Modul, 'mod_studio', 'Studio');

    expect(stan.czyOtworzycDymek({ ...wysoka, idModulu: 'mod_studio' }, TERAZ)).toBe(false);
    expect(stan.powodMilczenia({ ...wysoka, idModulu: 'mod_studio' }, TERAZ)).toContain('Studio');
    expect(stan.czyOtworzycDymek({ ...wysoka, idModulu: 'mod_terminal' }, TERAZ)).toBe(true);
  });

  it('wyciszenie klasy wstrzymuje jedną klasę, pozostałe ujawniają się', () => {
    stan.wyciszenia.przelaczKlase(KlasaZdarzen.StanKolejkiZadan);

    expect(stan.czyOtworzycDymek({ ...wysoka, klasa: KlasaZdarzen.StanKolejkiZadan }, TERAZ)).toBe(
      false,
    );
    expect(
      stan.czyOtworzycDymek({ ...wysoka, klasa: KlasaZdarzen.StanPetliWykonawczej }, TERAZ),
    ).toBe(true);
  });

  it('WYJĄTEK WAGI KRYTYCZNEJ przechodzi przez każdy rodzaj wyciszenia — plakietką, bez dymka', () => {
    const krytyczna = {
      waga: WagaUjawnienia.Wysoka,
      krytyczna: true,
      klasa: KlasaZdarzen.StanPetliWykonawczej,
      idModulu: 'mod_studio',
      idSesji: 'ses_1',
    } as const;

    for (const zaloz of [
      () => stan.wyciszenia.wyciszCzasem(WyciszenieCzasowe.Godzina, TERAZ),
      () => stan.wyciszenia.przelaczKontekst(ZakresKontekstu.Modul, 'mod_studio', 'Studio'),
      () => stan.wyciszenia.przelaczKontekst(ZakresKontekstu.KartaSesji, 'ses_1', 'Karta'),
      () => stan.wyciszenia.przelaczKlase(KlasaZdarzen.StanPetliWykonawczej),
      () => stan.przelaczTrybCichy(),
    ]) {
      stan = utworzStanObecnosci({ magazyn: null });
      zaloz();

      // Dymek się nie otwiera…
      expect(stan.czyOtworzycDymek(krytyczna, TERAZ)).toBe(false);
      // …ale powód mówi wprost, że sugestia ujawnia się plakietką…
      expect(stan.powodMilczenia(krytyczna, TERAZ)).toContain('plakietką, bez dymka');
      // …i sugestia krytyczna nie jest liczona jako wstrzymana.
      expect(stan.czySugestiaWstrzymana(krytyczna, TERAZ)).toBe(false);
    }
  });

  it('tryb cichy wyłącza syntezę mowy i zachowuje wyjątek', () => {
    stan.przelaczTrybCichy();
    expect(stan.tryb()).toBe(TrybObecnosci.Cichy);
    expect(stan.czySyntezaMowy()).toBe(false);
    expect(stan.czyAwatarWidoczny()).toBe(true);
  });

  it('plakietkę chowa wyłącznie wyciszenie czasowe', () => {
    stan.wyciszenia.przelaczKlase(KlasaZdarzen.StanKolejkiZadan);
    expect(stan.czyPlakietkaWidoczna(TERAZ)).toBe(true);

    stan.wycisz(WyciszenieCzasowe.Godzina, TERAZ);
    expect(stan.czyPlakietkaWidoczna(TERAZ)).toBe(false);
  });

  it('stan „Wyciszony" zachodzi przy każdym z trzech rodzajów', () => {
    expect(stan.czyJakiekolwiekWyciszenie(TERAZ)).toBe(false);
    stan.wyciszenia.przelaczKontekst(ZakresKontekstu.KartaSesji, 'ses_1', 'Karta');
    expect(stan.czyJakiekolwiekWyciszenie(TERAZ)).toBe(true);
  });

  it('skrót kwadransowy przełącza wyciszenie w obie strony', () => {
    stan.przelaczWyciszenieKwadransem(TERAZ);
    expect(stan.wyciszenie(TERAZ)?.doChwili).toBe(TERAZ + 15 * 60_000);
    stan.przelaczWyciszenieKwadransem(TERAZ);
    expect(stan.wyciszenie(TERAZ)).toBeNull();
  });

  it('zmiana wykazu wyciszeń rozgłasza się obserwatorom stanu obecności', () => {
    let powiadomien = 0;
    stan.obserwuj(() => (powiadomien += 1));
    stan.wyciszenia.przelaczKlase(KlasaZdarzen.Harmonogram);
    expect(powiadomien).toBe(1);
  });
});

describe('granica wobec kontraktu', () => {
  /**
   * Sprawdzian pilnuje tego samego, co przed dobudową wyciszenia w rdzeniu:
   * granica ma być NAZWANA, a nie zasłonięta. Zmieniła się wyłącznie strona
   * braku. Kontrakt niesie już wyciszenie nakładki wraz z odczytem, zapisem
   * i rozgłoszeniem, więc zdanie mówiące „brak po stronie kontraktu" byłoby
   * dziś nieprawdą — zostaje brak po stronie nakładki, której wołacze piszą
   * jeszcze do magazynu stanowiska.
   */
  it('granica wyciszenia jest nazwana wprost wraz ze stroną braku', () => {
    expect(brakiCzynne().length).toBeGreaterThan(0);
    const zdanie = zdanieGranicyWyciszenia();
    expect(zdanie).toContain('stanem TEGO okna');
    expect(zdanie).toContain('po stronie nakładki');
    // Zdanie nie ma wskazywać kontraktu jako winnego, bo pozycje w nim stoją.
    expect(zdanie).not.toContain('Brak jest po stronie kontraktu');
  });

  it('każdy brak czynny mówi, po czyjej stronie leży', () => {
    for (const brak of brakiCzynne()) {
      expect(zdanieBraku(brak)).toMatch(/po stronie (nakładki|kontraktu)/);
    }
  });
});
