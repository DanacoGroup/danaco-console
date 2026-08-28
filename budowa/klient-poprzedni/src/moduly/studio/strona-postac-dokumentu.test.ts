import { beforeEach, describe, expect, it } from 'vitest';

import {
  StudioAuthor,
  StudioHeaderScope,
  StudioPageNumberFormat,
  StudioTabKind,
  StudioTabLeader,
  StudioWatermarkKind,
} from '../../../../shared/contract';
import { domyslnaStrona } from './nastawy-strony';
import { utworzMagazynWidokuRdzenia, naNastawyRdzenia, zNastawRdzenia } from './strona-magazyn-widoku';
import { utworzPostacDokumentu, type PostacDokumentu } from './strona-postac-dokumentu';
import { liczbaPola, opiszBilans, wyborPola } from './strona-pola-postaci';
import { domyslneNastawyWidoku } from './widok-nastawy-operatora';
import type { ZrodloPostaciStudio } from './zrodlo-postaci-studio';

// Sprawdziany mierzą treść żądania do rdzenia, nie samo to, że przycisk dał się nacisnąć.

/** Jedno wywołanie zapisane przez atrapę źródła: nazwa wywołanej komendy wraz z treścią żądania, jakie do niej trafiło. */
interface Wywolanie {
  nazwa: string;
  zadanie: Record<string, unknown>;
}

/** Bilans czynności bez pominięć — najczęstsza odpowiedź rdzenia, gdy operacja objęła całe zaznaczenie bez odmów. */
const BILANS_CZYSTY = { applied: 1, skippedCount: 0 };

/** Postać dokumentu w kształcie najmniejszym, jaki niesie kontrakt — sam identyfikator dokumentu, bez list i pól opcjonalnych. */
const POSTAC = { documentId: 'dokument-jeden' };

/**
 * Atrapa źródła komend.
 *
 * Zapisuje każde wywołanie wraz z treścią żądania i oddaje odpowiedź kanoniczną.
 * `odpowiedzi` pozwala podmienić wybraną — tak mierzy się odmowę i odpowiedź
 * „rdzeń nie usunął", która jest odpowiedzią POPRAWNĄ, a nie usterką.
 */
function atrapaZrodla(
  zapis: Wywolanie[],
  odpowiedzi: Record<string, unknown> = {},
): ZrodloPostaciStudio {
  const kanoniczne: Record<string, unknown> = {
    postac: { form: { ...POSTAC, lists: [] } },
    zapiszPostac: { document: { id: 'dokument-jeden' }, form: POSTAC, balance: BILANS_CZYSTY },
    tekst: { text: 'fragment' },
    zmienTekst: { form: POSTAC, balance: BILANS_CZYSTY },
    pochodzenia: { entries: [] },
    nastawyStrony: { pageSetup: { marginTop: 20 }, sections: [] },
    ustawStrone: { form: POSTAC, balance: BILANS_CZYSTY, pageSetup: { marginTop: 20 } },
    nosniki: { papers: [{ name: 'A4', kind: 'sheet', widthMm: 210, heightMm: 297 }] },
    ustawKoperte: { form: POSTAC, balance: BILANS_CZYSTY, envelope: {} },
    wstawPodzial: { form: POSTAC, balance: BILANS_CZYSTY },
    naglowki: { headersFooters: [] },
    ustawNaglowek: { form: POSTAC, balance: BILANS_CZYSTY, headersFooters: [] },
    ustawNumeracje: { form: POSTAC, balance: BILANS_CZYSTY, numbering: { enabled: true } },
    ustawZnakWodny: {
      form: POSTAC,
      balance: BILANS_CZYSTY,
      watermark: { kind: StudioWatermarkKind.Text },
    },
    sekcje: { sections: [] },
    zapiszSekcje: { form: POSTAC, balance: BILANS_CZYSTY, section: { id: 'sekcja-jeden' } },
    usunSekcje: { deleted: true, form: POSTAC, balance: BILANS_CZYSTY },
    ustawTabulator: { form: POSTAC, balance: BILANS_CZYSTY, tabStops: [] },
    style: { styles: [] },
    zapiszStyl: { style: { name: 'Cytat', kind: 'paragraph' }, form: POSTAC, balance: BILANS_CZYSTY },
    zastosujStyl: { form: POSTAC, balance: BILANS_CZYSTY },
    usunStyl: { deleted: true, form: POSTAC, balance: BILANS_CZYSTY },
    postacZnaku: { character: {} },
    ustawZnak: { form: POSTAC, balance: BILANS_CZYSTY },
    postacAkapitu: { paragraph: {} },
    ustawAkapit: { form: POSTAC, balance: BILANS_CZYSTY },
    wyczyscFormat: { form: POSTAC, balance: BILANS_CZYSTY },
    ustawWielkoscLiter: { form: POSTAC, balance: BILANS_CZYSTY },
    pobierzPostac: { clipId: 'uchwyt-jeden', character: { fontFamily: 'Source Serif' } },
    nalozPostac: { form: POSTAC, balance: BILANS_CZYSTY },
    podobnePostacia: { matches: [], count: 0 },
    zamienZPostacia: { form: POSTAC, balance: BILANS_CZYSTY, matches: 3, replaced: 3 },
    zastosujListe: { form: POSTAC, balance: BILANS_CZYSTY, list: { id: 'lista-jeden', kind: 'bullet' } },
    ustawPunktator: { form: POSTAC, balance: BILANS_CZYSTY, list: { id: 'lista-jeden', kind: 'bullet' } },
    ustawNumeracjeListy: {
      form: POSTAC,
      balance: BILANS_CZYSTY,
      list: { id: 'lista-jeden', kind: 'number' },
    },
    wznowNumeracje: { form: POSTAC, balance: BILANS_CZYSTY },
    przestawPoziom: { form: POSTAC, balance: BILANS_CZYSTY },
    znaki: { symbols: [], categories: [] },
    wstawZnak: { form: POSTAC, balance: BILANS_CZYSTY, symbol: { code: '00A7', character: '§', name: 'paragraf' } },
    autozamiany: { rules: [] },
    ustawAutozamiane: { rule: { shortcut: '(c)', replacement: '©', enabled: true } },
    widok: { settings: {} },
    ustawWidok: { settings: {} },
  };

  const zrodlo: Record<string, unknown> = {};
  for (const nazwa of Object.keys(kanoniczne)) {
    zrodlo[nazwa] = async (zadanie: Record<string, unknown>) => {
      zapis.push({ nazwa, zadanie });
      const podmiana = odpowiedzi[nazwa];
      if (podmiana !== undefined) return podmiana;
      return { udany: true, wynik: kanoniczne[nazwa] };
    };
  }
  return zrodlo as unknown as ZrodloPostaciStudio;
}

/** Kontrolka pola o wskazanej etykiecie; brak pola jest błędem sprawdzianu, nie stanem, który sprawdzian ma obsłużyć. */
function pole(korzen: HTMLElement, etykieta: string): HTMLInputElement & HTMLSelectElement {
  for (const wiersz of Array.from(korzen.querySelectorAll<HTMLElement>('.dn-pole'))) {
    const napis = wiersz.querySelector('label');
    if (napis?.textContent !== etykieta) continue;
    const kontrolka = wiersz.querySelector<HTMLElement>('input, select, textarea');
    if (kontrolka !== null) return kontrolka as HTMLInputElement & HTMLSelectElement;
  }
  throw new Error(`Sprawdzian nie znalazł pola o etykiecie „${etykieta}"`);
}

/** Przycisk czynności wskazany jej kodem; brak przycisku jest błędem sprawdzianu, nie stanem do obsłużenia. */
function czynnosc(korzen: HTMLElement, kod: string): HTMLButtonElement {
  const przycisk = korzen.querySelector<HTMLButtonElement>(`[data-czynnosc='${kod}']`);
  if (przycisk === null) throw new Error(`Sprawdzian nie znalazł czynności „${kod}"`);
  return przycisk;
}

/** Zbudowana warstwa postaci wraz z zapisem wywołań przekazanych do atrapy źródła i zdaniami pokazanymi w oknie. */
interface Stanowisko {
  postac: PostacDokumentu;
  zapis: Wywolanie[];
  zdania: { tresc: string; powodzenie: boolean }[];
}

function zbuduj(
  ustawienia: {
    idDokumentu?: string;
    zaznaczenie?: { poczatek: number; koniec: number } | null;
    odpowiedzi?: Record<string, unknown>;
  } = {},
): Stanowisko {
  const zapis: Wywolanie[] = [];
  const zdania: { tresc: string; powodzenie: boolean }[] = [];
  const postac = utworzPostacDokumentu({
    zrodlo: atrapaZrodla(zapis, ustawienia.odpowiedzi ?? {}),
    idDokumentu: () => ustawienia.idDokumentu ?? 'dokument-jeden',
    idOkna: () => 'okno-jeden',
    // Wartość null znaczy „nic nie zaznaczono”, więc nie wolno jej podmienić wartością domyślną.
    zaznaczenie: () =>
      'zaznaczenie' in ustawienia ? ustawienia.zaznaczenie ?? null : { poczatek: 10, koniec: 40 },
    miejsceKursora: () => 25,
    naPostac: () => undefined,
    naZdanie: (tresc, powodzenie) => zdania.push({ tresc, powodzenie }),
  });
  return { postac, zapis, zdania };
}

/** Ostatnie żądanie wskazanej komendy — kolejne wywołania tej samej komendy nadpisują poprzednie w sprawdzianie. */
function ostatnie(zapis: Wywolanie[], nazwa: string): Record<string, unknown> {
  const wybrane = zapis.filter((wpis) => wpis.nazwa === nazwa);
  const koniec = wybrane[wybrane.length - 1];
  if (koniec === undefined) throw new Error(`Komenda ${nazwa} nie została wywołana ani raz`);
  return koniec.zadanie;
}

describe('numeracja stron jest nastawą Operatora, nie przełącznikiem tak/nie', () => {
  let stanowisko: Stanowisko;

  beforeEach(() => {
    stanowisko = zbuduj();
  });

  it('wysyła styl, punkt startu, wznowienie, liczbę stron i umiejscowienie', () => {
    const panel = stanowisko.postac.element;
    pole(panel, 'Numeracja stron włączona').checked = true;
    pole(panel, 'Styl numeru').value = StudioPageNumberFormat.RomanUpper;
    pole(panel, 'Numeracja od numeru').value = '7';
    pole(panel, 'Numeracja wznawia się w tej sekcji').checked = true;
    pole(panel, 'Dopisz liczbę stron — „strona N z M"').checked = true;
    pole(panel, 'Umiejscowienie numeru').value = 'stopka-srodek';
    czynnosc(panel, 'zapisz-numeracje').click();

    const zadanie = ostatnie(stanowisko.zapis, 'ustawNumeracje');
    expect(zadanie['documentId']).toBe('dokument-jeden');
    expect(zadanie['enabled']).toBe(true);
    expect(zadanie['format']).toBe(StudioPageNumberFormat.RomanUpper);
    expect(zadanie['startAt']).toBe(7);
    expect(zadanie['restartInSection']).toBe(true);
    expect(zadanie['showTotal']).toBe(true);
    expect(zadanie['position']).toBe('stopka-srodek');
    // Autor jedzie jawnie — bez niego przełącznik zmian modelu nie odróżni Operatora od modelu.
    expect(zadanie['author']).toBe(StudioAuthor.Uzytkownik);
  });

  it('styl niewybrany nie jedzie wcale — pole puste znaczy „nie ruszaj"', () => {
    const panel = stanowisko.postac.element;
    czynnosc(panel, 'zapisz-numeracje').click();
    const zadanie = ostatnie(stanowisko.zapis, 'ustawNumeracje');
    expect('format' in zadanie).toBe(false);
    expect('startAt' in zadanie).toBe(false);
    expect('position' in zadanie).toBe(false);
    // Sam przełącznik jedzie zawsze — kontrakt ma go jako pole obowiązkowe.
    expect(zadanie['enabled']).toBe(false);
  });

  it('numeracja wyłączona jest odpowiedzią prawdziwą, nie usterką', async () => {
    const stanowiskoWylaczone = zbuduj({
      odpowiedzi: {
        ustawNumeracje: {
          udany: true,
          wynik: { form: POSTAC, balance: BILANS_CZYSTY, numbering: { enabled: false } },
        },
      },
    });
    czynnosc(stanowiskoWylaczone.postac.element, 'zapisz-numeracje').click();
    await Promise.resolve();
    await Promise.resolve();
    const panel = stanowiskoWylaczone.postac.element;
    const zdanie = panel.querySelector('.dm-odpowiedz')?.textContent ?? '';
    expect(zdanie).toContain('wyłączona');
  });

  it('numeracja sekcji jedzie z jej identyfikatorem, nie z dokumentem', () => {
    const panel = stanowisko.postac.element;
    const wybor = pole(panel, 'Sekcja, której nastawy dotyczą');
    const opcja = document.createElement('option');
    opcja.value = 'sekcja-zalacznika';
    wybor.append(opcja);
    wybor.value = 'sekcja-zalacznika';
    czynnosc(panel, 'zapisz-numeracje').click();
    expect(ostatnie(stanowisko.zapis, 'ustawNumeracje')['sectionId']).toBe('sekcja-zalacznika');
  });
});

describe('marginesy są zmienialne i dojeżdżają do rdzenia', () => {
  it('nastawa gotowa wpisuje cztery liczby i oprawę w pola', () => {
    const { postac } = zbuduj();
    const panel = postac.element;
    const nastawa = pole(panel, 'Nastawa gotowa marginesów');
    nastawa.value = 'oprawa';
    nastawa.dispatchEvent(new Event('change'));
    expect(pole(panel, 'Margines górny w milimetrach').value).toBe('20');
    expect(pole(panel, 'Margines na oprawę w milimetrach').value).toBe('12');
    expect(pole(panel, 'Marginesy odbicia dla druku dwustronnego').checked).toBe(true);
  });

  it('zapis wysyła cztery marginesy, oprawę i marginesy odbicia', () => {
    const { postac, zapis } = zbuduj();
    const panel = postac.element;
    pole(panel, 'Margines górny w milimetrach').value = '15';
    pole(panel, 'Margines dolny w milimetrach').value = '16';
    pole(panel, 'Margines lewy w milimetrach').value = '30';
    pole(panel, 'Margines prawy w milimetrach').value = '17';
    pole(panel, 'Margines na oprawę w milimetrach').value = '8';
    pole(panel, 'Marginesy odbicia dla druku dwustronnego').checked = true;
    czynnosc(panel, 'zapisz-nastawy-strony').click();

    const zadanie = ostatnie(zapis, 'ustawStrone');
    expect(zadanie['marginTopMm']).toBe(15);
    expect(zadanie['marginBottomMm']).toBe(16);
    expect(zadanie['marginLeftMm']).toBe(30);
    expect(zadanie['marginRightMm']).toBe(17);
    expect(zadanie['gutterMm']).toBe(8);
    expect(zadanie['mirrorMargins']).toBe(true);
  });

  it('margines zerowy jedzie, a margines niewpisany nie — to nie to samo', () => {
    const { postac, zapis } = zbuduj();
    const panel = postac.element;
    pole(panel, 'Margines lewy w milimetrach').value = '0';
    czynnosc(panel, 'zapisz-nastawy-strony').click();
    const zadanie = ostatnie(zapis, 'ustawStrone');
    expect(zadanie['marginLeftMm']).toBe(0);
    expect('marginTopMm' in zadanie).toBe(false);
  });

  it('format własny podany jedną liczbą jest odmawiany, nie zgadywany', () => {
    const { postac, zapis } = zbuduj();
    const panel = postac.element;
    pole(panel, 'Format własny — szerokość w milimetrach').value = '250';
    czynnosc(panel, 'zapisz-nastawy-strony').click();
    expect(zapis.filter((wpis) => wpis.nazwa === 'ustawStrone')).toHaveLength(0);
    expect(panel.querySelector('.dm-odpowiedz')?.textContent ?? '').toContain('DWIEMA');
  });

  it('chwyt marginesu z linijki jedzie tą samą komendą co pole panelu', () => {
    const { postac, zapis } = zbuduj();
    postac.zglosMargines('lewy', 34.5);
    const zadanie = ostatnie(zapis, 'ustawStrone');
    expect(zadanie['marginLeftMm']).toBe(34.5);
    // Chwyt przestawia JEDEN margines i tylko on jedzie: pozostałe zostają
    // takie, jakie stoją w rdzeniu.
    expect('marginRightMm' in zadanie).toBe(false);
  });

  it('całe nastawy strony z powierzchni jadą jednym żądaniem', () => {
    const { postac, zapis } = zbuduj();
    const strona = { ...domyslnaStrona(), marginesOprawyMm: 12, marginesyOdbicia: true };
    postac.zglosStrone(strona);
    const zadanie = ostatnie(zapis, 'ustawStrone');
    expect(zadanie['paperName']).toBe('A4');
    expect(zadanie['gutterMm']).toBe(12);
    expect(zadanie['mirrorMargins']).toBe(true);
    expect(zadanie['marginTopMm']).toBe(20);
  });

  it('nośnik własny jedzie wymiarami, nie samym opisem', () => {
    const { postac, zapis } = zbuduj();
    postac.zglosStrone({
      ...domyslnaStrona(),
      nosnik: { oznaczenie: 'własny 250×350 mm', szerokoscMm: 250, wysokoscMm: 350, rodzaj: 'wlasny' },
    });
    const zadanie = ostatnie(zapis, 'ustawStrone');
    expect(zadanie['widthMm']).toBe(250);
    expect(zadanie['heightMm']).toBe(350);
  });
});

describe('nagłówek, stopka, znak wodny i koperta', () => {
  it('nagłówek jedzie z zasięgiem i pozwala opróżnić treść', () => {
    const { postac, zapis } = zbuduj();
    const panel = postac.element;
    pole(panel, 'Zasięg nagłówka i stopki').value = StudioHeaderScope.FirstPage;
    pole(panel, 'Treść nagłówka').value = '';
    pole(panel, 'Treść stopki').value = 'Strona';
    czynnosc(panel, 'zapisz-naglowek').click();
    const zadanie = ostatnie(zapis, 'ustawNaglowek');
    expect(zadanie['scope']).toBe(StudioHeaderScope.FirstPage);
    // Treść pusta jedzie jawnie: opróżnienie nagłówka jest czynnością, nie brakiem wskazania.
    expect(zadanie['headerText']).toBe('');
    expect(zadanie['footerText']).toBe('Strona');
  });

  it('znak wodny jedzie rodzajem, napisem i kryciem', () => {
    const { postac, zapis } = zbuduj();
    const panel = postac.element;
    pole(panel, 'Znak wodny').value = StudioWatermarkKind.Text;
    pole(panel, 'Napis znaku wodnego').value = 'PROJEKT';
    pole(panel, 'Krycie od 0 do 1').value = '0.3';
    czynnosc(panel, 'zapisz-znak-wodny').click();
    const zadanie = ostatnie(zapis, 'ustawZnakWodny');
    expect(zadanie['kind']).toBe(StudioWatermarkKind.Text);
    expect(zadanie['text']).toBe('PROJEKT');
    expect(zadanie['opacity']).toBe(0.3);
  });

  it('koperta jedzie adresami i ich położeniem', () => {
    const { postac, zapis } = zbuduj();
    const panel = postac.element;
    pole(panel, 'Adres adresata').value = 'Sąd Okręgowy';
    pole(panel, 'Adresat — od lewej krawędzi w milimetrach').value = '100';
    pole(panel, 'Nadrukuj nadawcę').checked = true;
    czynnosc(panel, 'zapisz-koperte').click();
    const zadanie = ostatnie(zapis, 'ustawKoperte');
    expect(zadanie['recipient']).toBe('Sąd Okręgowy');
    expect(zadanie['recipientXMm']).toBe(100);
    expect(zadanie['includeSender']).toBe(true);
  });

  it('podział wstawia się w miejscu kursora, nie na końcu treści', () => {
    const { postac, zapis } = zbuduj();
    czynnosc(postac.element, 'wstaw-podzial').click();
    expect(ostatnie(zapis, 'wstawPodzial')['offset']).toBe(25);
  });

  it('usunięcie sekcji bez wskazanej sekcji jest odmawiane nazwanym powodem', () => {
    const { postac, zapis } = zbuduj();
    czynnosc(postac.element, 'usun-sekcje').click();
    expect(zapis.filter((wpis) => wpis.nazwa === 'usunSekcje')).toHaveLength(0);
    expect(postac.element.querySelector('.dm-odpowiedz')?.textContent ?? '').toContain(
      'nie jest sekcją',
    );
  });
});

describe('tabulatory z linijki idą różnicą, nie przepisaniem wykazu', () => {
  it('nowy tabulator jedzie jednym żądaniem założenia', () => {
    const { postac, zapis } = zbuduj();
    postac.zglosTabulatory([{ milimetry: 40, rodzaj: 'dziesietny', znakWiodacy: 'kropka' }]);
    const zadanie = ostatnie(zapis, 'ustawTabulator');
    expect(zadanie['positionMm']).toBe(40);
    expect(zadanie['kind']).toBe(StudioTabKind.Decimal);
    expect(zadanie['leader']).toBe(StudioTabLeader.Dot);
    expect(zadanie['remove']).toBe(false);
  });

  it('przesunięcie tabulatora daje zdjęcie starego i założenie nowego', async () => {
    const { postac, zapis } = zbuduj();
    postac.zglosTabulatory([{ milimetry: 40, rodzaj: 'lewy', znakWiodacy: 'brak' }]);
    await Promise.resolve();
    zapis.length = 0;
    postac.zglosTabulatory([{ milimetry: 55, rodzaj: 'lewy', znakWiodacy: 'brak' }]);
    // Pierwsze żądanie wychodzi natychmiast, drugie po jego rozstrzygnięciu — kolejność ma znaczenie.
    await Promise.resolve();
    await Promise.resolve();
    await Promise.resolve();
    const wywolania = zapis.filter((wpis) => wpis.nazwa === 'ustawTabulator');
    expect(wywolania[0]?.zadanie['positionMm']).toBe(40);
    expect(wywolania[0]?.zadanie['remove']).toBe(true);
    expect(wywolania[1]?.zadanie['positionMm']).toBe(55);
    expect(wywolania[1]?.zadanie['remove']).toBe(false);
  });

  it('wykaz bez zmian nie woła rdzenia ani raz', () => {
    const { postac, zapis } = zbuduj();
    postac.zglosTabulatory([{ milimetry: 40, rodzaj: 'lewy', znakWiodacy: 'brak' }]);
    zapis.length = 0;
    postac.zglosTabulatory([{ milimetry: 40, rodzaj: 'lewy', znakWiodacy: 'brak' }]);
    expect(zapis).toHaveLength(0);
  });

  it('wcięcia z linijki idą stylem akapitu na zaznaczonym fragmencie', () => {
    const { postac, zapis } = zbuduj();
    postac.zglosWciecia({ leweMm: 12.5, pierwszyWierszMm: -5 });
    const zadanie = ostatnie(zapis, 'ustawAkapit');
    expect(zadanie['indentLeftMm']).toBe(12.5);
    expect(zadanie['firstLineIndentMm']).toBe(-5);
    expect(zadanie['rangeStart']).toBe(10);
    expect(zadanie['rangeEnd']).toBe(40);
  });
});

describe('brak dokumentu jest odmową nazwaną, nie ciszą', () => {
  it('nastawy strony bez dokumentu nie jadą i mówią, czego brakuje', () => {
    const { postac, zapis } = zbuduj({ idDokumentu: '' });
    czynnosc(postac.element, 'zapisz-nastawy-strony').click();
    expect(zapis.filter((wpis) => wpis.nazwa === 'ustawStrone')).toHaveLength(0);
    expect(postac.element.querySelector('.dm-odpowiedz')?.textContent ?? '').toContain(
      'po stronie okna',
    );
  });

  it('zmiana brzmienia bez zaznaczenia jest odmawiana — komenda wymaga fragmentu', async () => {
    const { postac, zapis, zdania } = zbuduj({ zaznaczenie: null });
    const udane = await postac.zmienTekst('nowe brzmienie', true);
    expect(udane).toBe(false);
    expect(zapis.filter((wpis) => wpis.nazwa === 'zmienTekst')).toHaveLength(0);
    expect(zdania[zdania.length - 1]?.tresc ?? '').toContain('FRAGMENTU');
  });
});

describe('zapis postaci i uczciwość niepowodzenia', () => {
  it('zapis niesie postać, treść i zakłada wersję', async () => {
    const { postac, zapis } = zbuduj();
    const udane = await postac.zapiszPostac({ documentId: 'dokument-jeden' }, 'treść', 'Pismo');
    expect(udane).toBe(true);
    const zadanie = ostatnie(zapis, 'zapiszPostac');
    expect(zadanie['content']).toBe('treść');
    expect(zadanie['title']).toBe('Pismo');
    expect(zadanie['createVersion']).toBe(true);
    expect(zadanie['form']).toEqual({ documentId: 'dokument-jeden' });
  });

  it('nieudany zapis oddaje fałsz i nazywa odmowę — nigdy „zapisano"', async () => {
    const { postac, zdania } = zbuduj({
      odpowiedzi: {
        zapiszPostac: {
          udany: false,
          blad: { code: 'validation_failed', message: 'postać niezgodna', retryable: false },
        },
      },
    });
    const udane = await postac.zapiszPostac({ documentId: 'dokument-jeden' }, 'treść', '');
    expect(udane).toBe(false);
    const zdanie = zdania[zdania.length - 1];
    expect(zdanie?.powodzenie).toBe(false);
    expect(zdanie?.tresc ?? '').toContain('postać niezgodna');
  });

  it('bilans z pominięciem nazywa blokadę, która czynność zatrzymała', () => {
    const zdanie = opiszBilans({
      applied: 4,
      skippedCount: 1,
      skipped: [
        {
          reason: 'fragment pod blokadą',
          detail: 'podstawa prawna',
          rangeStart: 100,
          rangeEnd: 180,
          lockName: 'podstawa prawna — nie zmieniać',
        },
      ],
      note: 'Zamiana wykonana poza blokadą.',
    });
    expect(zdanie).toContain('miejsc zmienionych 4');
    expect(zdanie).toContain('pominiętych 1');
    expect(zdanie).toContain('podstawa prawna — nie zmieniać');
    expect(zdanie).toContain('100–180');
  });
});

describe('nastawy widoku prowadzi rdzeń, bez zmiany wołaczy', () => {
  it('przełożenie na kontrakt i z powrotem zachowuje wybory Operatora', () => {
    const nastawy = {
      ...domyslneNastawyWidoku(),
      trybDokumentow: 'podzial' as const,
      kierunekPodzialu: 'poziomy' as const,
      przewijanie: 'strona-po-stronie' as const,
      ukladKartek: 'obok' as const,
      kartekWRzedzie: 3,
      jednostka: 'cal' as const,
      linijkiWidoczne: false,
      graniceMarginesow: false,
      skala: 145,
    };
    const wrocone = zNastawRdzenia({ ...naNastawyRdzenia(nastawy) });
    expect(wrocone.trybDokumentow).toBe('podzial');
    expect(wrocone.kierunekPodzialu).toBe('poziomy');
    expect(wrocone.przewijanie).toBe('strona-po-stronie');
    expect(wrocone.ukladKartek).toBe('obok');
    expect(wrocone.kartekWRzedzie).toBe(3);
    expect(wrocone.jednostka).toBe('cal');
    expect(wrocone.linijkiWidoczne).toBe(false);
    expect(wrocone.graniceMarginesow).toBe(false);
    expect(wrocone.skala).toBe(145);
  });

  it('rozkładówka i liczba kartek w rzędzie nie mieszają się', () => {
    const wrocone = zNastawRdzenia({ spreadView: true, pagesPerRow: 4 });
    expect(wrocone.ukladKartek).toBe('rozkladowka');
  });

  it('zapis magazynu jedzie komendą widoku rdzenia', async () => {
    const zapis: Wywolanie[] = [];
    const magazyn = utworzMagazynWidokuRdzenia({
      zrodlo: atrapaZrodla(zapis),
      idDokumentu: () => 'dokument-jeden',
      idOkna: () => 'okno-jeden',
      naZdanie: () => undefined,
      odbicie: null,
    });
    magazyn.setItem(
      'dn.studio.widok',
      JSON.stringify({ ...domyslneNastawyWidoku(), skala: 175, jednostka: 'cal' }),
    );
    await Promise.resolve();
    const zadanie = ostatnie(zapis, 'ustawWidok');
    expect(zadanie['zoomPercent']).toBe(175);
    expect(zadanie['rulerUnit']).toBe('inch');
    expect(zadanie['documentId']).toBe('dokument-jeden');
  });

  it('odczyt z rdzenia napełnia magazyn, więc wołacze nie zmieniają drogi', async () => {
    const zapis: Wywolanie[] = [];
    const magazyn = utworzMagazynWidokuRdzenia({
      zrodlo: atrapaZrodla(zapis, {
        widok: { udany: true, wynik: { settings: { zoomPercent: 133, rulersVisible: false } } },
      }),
      idDokumentu: () => 'dokument-jeden',
      idOkna: () => '',
      naZdanie: () => undefined,
      odbicie: null,
    });
    expect(magazyn.getItem('dn.studio.widok')).toBeNull();
    expect(await magazyn.wczytaj()).toBe(true);
    const zapisany = magazyn.getItem('dn.studio.widok') ?? '{}';
    expect(JSON.parse(zapisany)).toMatchObject({ skala: 133, linijkiWidoczne: false });
    // Okno puste nie jedzie w żądaniu: kontrakt czyta brak jako „bez zawężenia do okna”.
    expect('windowId' in ostatnie(zapis, 'widok')).toBe(false);
  });

  it('odmowa odczytu nie przerywa pracy, ale nazywa brak trwałości', async () => {
    const zdania: string[] = [];
    const magazyn = utworzMagazynWidokuRdzenia({
      zrodlo: atrapaZrodla([], {
        widok: { udany: false, blad: { code: 'not_found', message: 'brak uchwytu' } },
      }),
      idDokumentu: () => 'dokument-jeden',
      idOkna: () => 'okno-jeden',
      naZdanie: (tresc) => zdania.push(tresc),
      odbicie: null,
    });
    expect(await magazyn.wczytaj()).toBe(false);
    expect(zdania[0] ?? '').toContain('nie przeniosą się na inne');
  });
});

describe('panel fragmentu, postaci i pochodzenia', () => {
  it('odczyt fragmentu wpisuje treść i nazywa blokady, które go obejmują', async () => {
    const stanowisko = zbuduj({
      odpowiedzi: {
        tekst: {
          udany: true,
          wynik: {
            text: 'Na podstawie art. 10',
            runs: [{ text: 'Na podstawie art. 10' }],
            locks: [
              {
                id: 'blokada-jedna',
                documentId: 'dokument-jeden',
                name: 'podstawa prawna',
                rangeStart: 0,
                rangeEnd: 20,
                scope: 'model',
              },
            ],
          },
        },
      },
    });
    const panel = stanowisko.postac.element;
    czynnosc(panel, 'odczytaj-fragment').click();
    await Promise.resolve();
    await Promise.resolve();
    expect(pole(panel, 'Brzmienie zaznaczonego fragmentu').value).toBe('Na podstawie art. 10');
    expect(panel.textContent ?? '').toContain('podstawa prawna');
    // Zdanie o blokadzie mówi też, gdzie stoi sprawdzenie — wyłącznie w oknie byłoby pozorne.
    expect(panel.textContent ?? '').toContain('Sprawdzenie stoi w RDZENIU');
  });

  it('utrwalenie postaci odczytuje ją najpierw i NIE wysyła treści', async () => {
    const stanowisko = zbuduj();
    czynnosc(stanowisko.postac.element, 'utrwal-postac').click();
    await Promise.resolve();
    await Promise.resolve();
    await Promise.resolve();
    const zadanie = ostatnie(stanowisko.zapis, 'zapiszPostac');
    // Brak pola treści znaczy „bez zmiany treści"; treść pusta skasowałaby dokument.
    expect('content' in zadanie).toBe(false);
    expect(zadanie['createVersion']).toBe(true);
    // Odczyt wyprzedza zapis: postać sprzed cudzej zmiany cofnęłaby ją bez słowa.
    const kolejnosc = stanowisko.zapis.map((wpis) => wpis.nazwa);
    expect(kolejnosc.indexOf('postac')).toBeLessThan(kolejnosc.indexOf('zapiszPostac'));
  });

  it('wykaz pochodzeń pusty mówi, że to stan prawidłowy, a nie brak', async () => {
    const stanowisko = zbuduj();
    czynnosc(stanowisko.postac.element, 'odczytaj-pochodzenia').click();
    await Promise.resolve();
    await Promise.resolve();
    expect(stanowisko.postac.element.textContent ?? '').toContain('stan prawidłowy');
  });

  it('zestawienie postaci liczy style, sekcje i blokady', async () => {
    const stanowisko = zbuduj({
      odpowiedzi: {
        postac: {
          udany: true,
          wynik: {
            form: {
              documentId: 'dokument-jeden',
              styles: [{ name: 'Cytat', kind: 'paragraph' }],
              sections: [{ id: 'a', index: 0, rangeStart: 0, rangeEnd: 10 }],
              revision: 12,
            },
          },
        },
      },
    });
    czynnosc(stanowisko.postac.element, 'odczytaj-postac').click();
    await Promise.resolve();
    await Promise.resolve();
    expect(stanowisko.postac.element.textContent ?? '').toContain('stylów 1 · sekcji 1');
    expect(stanowisko.postac.element.textContent ?? '').toContain('numer porządkowy postaci 12');
    expect(stanowisko.postac.postacZapamietana()?.revision).toBe(12);
  });
});

describe('pola postaci czytają wartości ostrożnie', () => {
  it('pole puste oddaje brak, a nie zero', () => {
    const kontrolka = document.createElement('input');
    kontrolka.value = '';
    expect(liczbaPola(kontrolka)).toBeUndefined();
    kontrolka.value = '0';
    expect(liczbaPola(kontrolka)).toBe(0);
  });

  it('pozycja „bez zmiany" oddaje brak, a nie łańcuch pusty', () => {
    const kontrolka = document.createElement('select');
    const bezZmiany = document.createElement('option');
    bezZmiany.value = '';
    const wybrana = document.createElement('option');
    wybrana.value = 'arabic';
    kontrolka.append(bezZmiany, wybrana);
    expect(wyborPola(kontrolka)).toBeUndefined();
    kontrolka.value = 'arabic';
    expect(wyborPola(kontrolka)).toBe('arabic');
  });
});
