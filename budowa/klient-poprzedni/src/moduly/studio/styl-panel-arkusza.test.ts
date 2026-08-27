import { describe, expect, it } from 'vitest';

import {
  StudioBulletSource,
  StudioCaseTransform,
  StudioListKind,
  StudioListNumberFormat,
  StudioStyleKind,
  StudioTextAlign,
  StudioUnderlineStyle,
} from '../../../../shared/contract';
import { utworzStylPanelArkusza, type CzynnosciStyluPanelu } from './styl-panel-arkusza';
import { utworzStylPanelList, type CzynnosciListPanelu } from './styl-panel-list';

/** Zapis żądań złożonych przez panel gromadzi nazwę wywołanej czynności wraz z pełną treścią zgłoszonego zadania, zachowując kolejność wywołań. */
type Zapis = { nazwa: string; zadanie: Record<string, unknown> }[];

/** Kontrolka pola o wskazanej etykiecie odnajduje wiersz formularza dopasowany do etykiety i zwraca jego pole wejściowe gotowe do odczytu albo zapisu wartości. */
function pole(korzen: HTMLElement, etykieta: string): HTMLInputElement & HTMLSelectElement {
  for (const wiersz of Array.from(korzen.querySelectorAll<HTMLElement>('.dn-pole'))) {
    const napis = wiersz.querySelector('label');
    if (napis?.textContent !== etykieta) continue;
    const kontrolka = wiersz.querySelector<HTMLElement>('input, select, textarea');
    if (kontrolka !== null) return kontrolka as HTMLInputElement & HTMLSelectElement;
  }
  throw new Error(`Sprawdzian nie znalazł pola o etykiecie „${etykieta}"`);
}

function czynnosc(korzen: HTMLElement, kod: string): HTMLButtonElement {
  const przycisk = korzen.querySelector<HTMLButtonElement>(`[data-czynnosc='${kod}']`);
  if (przycisk === null) throw new Error(`Sprawdzian nie znalazł czynności „${kod}"`);
  return przycisk;
}

function odpowiedzPanelu(korzen: HTMLElement): string {
  return korzen.querySelector('.dm-odpowiedz')?.textContent ?? '';
}

function panelStylu(zapis: Zapis) {
  const czynnosci: CzynnosciStyluPanelu = {
    naZapisStylu: (zadanie) => zapis.push({ nazwa: 'zapisStylu', zadanie: { ...zadanie } }),
    naStosowanieStylu: (zadanie) => zapis.push({ nazwa: 'stosowanie', zadanie: { ...zadanie } }),
    naUsuniecieStylu: (zadanie) => zapis.push({ nazwa: 'usuniecie', zadanie: { ...zadanie } }),
    naStylZnaku: (zadanie) => zapis.push({ nazwa: 'stylZnaku', zadanie: { ...zadanie } }),
    naStylAkapitu: (zadanie) => zapis.push({ nazwa: 'stylAkapitu', zadanie: { ...zadanie } }),
    naCzyszczenie: (zadanie) => zapis.push({ nazwa: 'czyszczenie', zadanie: { ...zadanie } }),
    naWielkoscLiter: (zadanie) => zapis.push({ nazwa: 'wielkoscLiter', zadanie: { ...zadanie } }),
    naPobraniePostaci: (zadanie) => zapis.push({ nazwa: 'pobranie', zadanie: { ...zadanie } }),
    naNalozeniePostaci: (zadanie) => zapis.push({ nazwa: 'nalozenie', zadanie: { ...zadanie } }),
    naPodobne: (zadanie) => zapis.push({ nazwa: 'podobne', zadanie: { ...zadanie } }),
    naZamiane: (zadanie) => zapis.push({ nazwa: 'zamiana', zadanie: { ...zadanie } }),
    naTabulator: (zadanie) => zapis.push({ nazwa: 'tabulator', zadanie: { ...zadanie } }),
    naOdczyt: () => zapis.push({ nazwa: 'odczyt', zadanie: {} }),
  };
  return utworzStylPanelArkusza(czynnosci);
}

function panelList(zapis: Zapis) {
  const czynnosci: CzynnosciListPanelu = {
    naListe: (zadanie) => zapis.push({ nazwa: 'lista', zadanie: { ...zadanie } }),
    naPunktator: (zadanie) => zapis.push({ nazwa: 'punktator', zadanie: { ...zadanie } }),
    naNumeracjeListy: (zadanie) => zapis.push({ nazwa: 'numeracjaListy', zadanie: { ...zadanie } }),
    naWznowienie: (zadanie) => zapis.push({ nazwa: 'wznowienie', zadanie: { ...zadanie } }),
    naPoziom: (zadanie) => zapis.push({ nazwa: 'poziom', zadanie: { ...zadanie } }),
    naZnak: (zadanie) => zapis.push({ nazwa: 'znak', zadanie: { ...zadanie } }),
    naSzukanieZnaku: (fraza, tylkoOstatnie) =>
      zapis.push({ nazwa: 'szukanie', zadanie: { fraza, tylkoOstatnie } }),
    naAutozamiane: (zadanie) => zapis.push({ nazwa: 'autozamiana', zadanie: { ...zadanie } }),
    naOdczyt: () => zapis.push({ nazwa: 'odczyt', zadanie: {} }),
  };
  return utworzStylPanelList(czynnosci);
}

describe('arkusz stylów nazwanych', () => {
  it('zapis stylu niesie nazwę, rodzaj, dziedziczenie i postać z pól', () => {
    const zapis: Zapis = [];
    const panel = panelStylu(zapis);
    pole(panel.element, 'Nazwa stylu').value = 'Cytat';
    pole(panel.element, 'Nazwa widoczna dla Operatora').value = 'Cytat blokowy';
    pole(panel.element, 'Rodzaj stylu').value = StudioStyleKind.Paragraph;
    pole(panel.element, 'Styl nadrzędny, po którym ten dziedziczy').value = 'Tekst zasadniczy';
    pole(panel.element, 'Krój pisma').value = 'Source Serif';
    pole(panel.element, 'Wyrównanie').value = StudioTextAlign.Justify;
    czynnosc(panel.element, 'zapisz-styl').click();

    const zadanie = zapis[0]?.zadanie ?? {};
    expect(zadanie['name']).toBe('Cytat');
    expect(zadanie['displayName']).toBe('Cytat blokowy');
    expect(zadanie['basedOn']).toBe('Tekst zasadniczy');
    expect(zadanie['character']).toMatchObject({ fontFamily: 'Source Serif' });
    expect(zadanie['paragraph']).toMatchObject({ align: StudioTextAlign.Justify });
  });

  it('styl bez nazwy jest odmawiany nazwanym powodem', () => {
    const zapis: Zapis = [];
    const panel = panelStylu(zapis);
    czynnosc(panel.element, 'zapisz-styl').click();
    expect(zapis).toHaveLength(0);
    expect(odpowiedzPanelu(panel.element)).toContain('nazwij go');
  });

  it('wykaz stylów pokazuje liczbę miejsc użycia i wskazuje styl w polach', () => {
    const zapis: Zapis = [];
    const panel = panelStylu(zapis);
    panel.pokazStyle([
      {
        name: 'Naglowek1',
        displayName: 'Nagłówek poziomu pierwszego',
        kind: StudioStyleKind.Paragraph,
        builtin: true,
        usageCount: 12,
      },
    ]);
    const pozycja = panel.element.querySelector<HTMLElement>("[data-styl='Naglowek1']");
    expect(pozycja?.textContent ?? '').toContain('miejsc użycia 12');
    expect(pozycja?.textContent ?? '').toContain('fabryczny');
    pozycja?.querySelector('button')?.click();
    expect(pole(panel.element, 'Nazwa stylu').value).toBe('Naglowek1');
  });

  it('wykaz pusty mówi wprost, że dokument nie ma arkusza — nie milczy', () => {
    const panel = panelStylu([]);
    panel.pokazStyle([]);
    expect(panel.element.textContent ?? '').toContain('ani jednego stylu');
  });
});

describe('styl znaku ma trzy stany cech logicznych', () => {
  it('cecha niewskazana nie jedzie wcale', () => {
    const zapis: Zapis = [];
    const panel = panelStylu(zapis);
    pole(panel.element, 'Pogrubienie').value = 'tak';
    pole(panel.element, 'Podkreślenie').value = StudioUnderlineStyle.Double;
    czynnosc(panel.element, 'nanies-styl-znaku').click();

    const zadanie = zapis[0]?.zadanie ?? {};
    expect(zadanie['bold']).toBe(true);
    expect(zadanie['underline']).toBe(StudioUnderlineStyle.Double);
    // Kursywy Operator nie tknął — nie wolno jej zdejmować przy okazji.
    expect('italic' in zadanie).toBe(false);
    expect('smallCaps' in zadanie).toBe(false);
  });

  it('cecha wyłączona jedzie jako fałsz, nie jako brak', () => {
    const zapis: Zapis = [];
    const panel = panelStylu(zapis);
    pole(panel.element, 'Kapitaliki').value = 'nie';
    czynnosc(panel.element, 'nanies-styl-znaku').click();
    expect(zapis[0]?.zadanie['smallCaps']).toBe(false);
  });

  it('żądanie bez ani jednego pola jest odmawiane, a nie wysyłane', () => {
    const zapis: Zapis = [];
    const panel = panelStylu(zapis);
    czynnosc(panel.element, 'nanies-styl-znaku').click();
    expect(zapis).toHaveLength(0);
    expect(odpowiedzPanelu(panel.element)).toContain('nie ma czego nanieść');
  });

  it('krok stopnia jedzie osobno od stopnia — działa bez znajomości zastanego', () => {
    const zapis: Zapis = [];
    const panel = panelStylu(zapis);
    pole(panel.element, 'Powiększ albo pomniejsz stopień o').value = '-1';
    czynnosc(panel.element, 'nanies-styl-znaku').click();
    const zadanie = zapis[0]?.zadanie ?? {};
    expect(zadanie['fontSizeStepPt']).toBe(-1);
    expect('fontSizePt' in zadanie).toBe(false);
  });

  it('postać znaku i akapitu odczytana z rdzenia wchodzi w pola wraz z niejednolitymi', () => {
    const panel = panelStylu([]);
    panel.pokazPostacZnaku({ fontFamily: 'Inter', bold: false, fontSizePt: 11 }, ['fontFamily']);
    expect(pole(panel.element, 'Krój pisma').value).toBe('Inter');
    expect(pole(panel.element, 'Pogrubienie').value).toBe('nie');
    expect(panel.element.textContent ?? '').toContain('NIEJEDNOLITE: fontFamily');
    panel.pokazPostacAkapitu({ align: StudioTextAlign.Center, indentLeftMm: 10 }, []);
    expect(pole(panel.element, 'Wcięcie lewe w milimetrach').value).toBe('10');
    expect(panel.element.textContent ?? '').toContain('akapitu jest we fragmencie jednolita');
  });
});

describe('czyszczenie, wielkość liter, malarz i podobieństwo', () => {
  it('czyszczenie bez wskazania, co czyścić, jest odmawiane', () => {
    const zapis: Zapis = [];
    const panel = panelStylu(zapis);
    pole(panel.element, 'Czyść styl znaku').checked = false;
    pole(panel.element, 'Czyść styl akapitu').checked = false;
    czynnosc(panel.element, 'wyczysc-format').click();
    expect(zapis).toHaveLength(0);
    expect(odpowiedzPanelu(panel.element)).toContain('Wskaż, co czyścić');
  });

  it('wielkość liter jedzie wskazanym przestawieniem', () => {
    const zapis: Zapis = [];
    const panel = panelStylu(zapis);
    pole(panel.element, 'Wielkość liter zaznaczenia').value = StudioCaseTransform.Sentence;
    czynnosc(panel.element, 'przestaw-litery').click();
    expect(zapis[0]?.zadanie['transform']).toBe(StudioCaseTransform.Sentence);
  });

  it('malarz bez pobranej postaci ma przycisk naniesienia zablokowany', () => {
    const panel = panelStylu([]);
    panel.pokazMalarza(false, '');
    expect(czynnosc(panel.element, 'malarz-naloz').disabled).toBe(true);
    panel.pokazMalarza(true, 'uchwyt-jeden');
    expect(czynnosc(panel.element, 'malarz-naloz').disabled).toBe(false);
    expect(panel.element.textContent ?? '').toContain('uchwyt-jeden');
  });

  it('zaznaczenie wedle podobieństwa niesie wzór i zakres dopasowania', () => {
    const zapis: Zapis = [];
    const panel = panelStylu(zapis);
    pole(panel.element, 'Wzór podobieństwa — styl nazwany').value = 'Cytat';
    pole(panel.element, 'Dopasuj styl akapitu').checked = true;
    czynnosc(panel.element, 'podobne-formatowanie').click();
    const zadanie = zapis[0]?.zadanie ?? {};
    expect(zadanie['styleName']).toBe('Cytat');
    expect(zadanie['matchParagraph']).toBe(true);
    expect(zadanie['matchCharacter']).toBe(true);
  });
});

describe('znajdź i zamień wraz z postacią', () => {
  it('zamiana bez wzoru jest odmawiana — nie objęłaby całego dokumentu po cichu', () => {
    const zapis: Zapis = [];
    const panel = panelStylu(zapis);
    czynnosc(panel.element, 'zamien-z-postacia').click();
    expect(zapis).toHaveLength(0);
    expect(odpowiedzPanelu(panel.element)).toContain('Nie ma czego szukać');
  });

  it('zamiana niesie treść, styl i postać nadawaną trafieniom', () => {
    const zapis: Zapis = [];
    const panel = panelStylu(zapis);
    pole(panel.element, 'Szukana treść').value = 'ustawa';
    pole(panel.element, 'Treść wstawiana').value = 'Ustawa';
    pole(panel.element, 'Szukany styl nazwany').value = 'Tekst zasadniczy';
    pole(panel.element, 'Styl nazwany nadawany trafieniom').value = 'Cytat';
    pole(panel.element, 'Odróżniaj wielkość liter').checked = true;
    pole(panel.element, 'Barwa tekstu').value = '#800000';
    czynnosc(panel.element, 'zamien-z-postacia').click();

    const zadanie = zapis[0]?.zadanie ?? {};
    expect(zadanie['findText']).toBe('ustawa');
    expect(zadanie['replaceText']).toBe('Ustawa');
    expect(zadanie['findStyleName']).toBe('Tekst zasadniczy');
    expect(zadanie['replaceStyleName']).toBe('Cytat');
    expect(zadanie['matchCase']).toBe(true);
    expect(zadanie['replaceFormat']).toMatchObject({ color: '#800000' });
  });

  it('wybór „tylko w zaznaczeniu" jest czytany osobno, nie zgadywany', () => {
    const panel = panelStylu([]);
    expect(panel.zamianaTylkoWZaznaczeniu()).toBe(false);
    pole(panel.element, 'Tylko w zaznaczeniu').checked = true;
    expect(panel.zamianaTylkoWZaznaczeniu()).toBe(true);
  });
});

describe('tabulator wpisywany liczbą', () => {
  it('bez położenia nie jedzie — tabulator bez miejsca nie istnieje', () => {
    const zapis: Zapis = [];
    const panel = panelStylu(zapis);
    czynnosc(panel.element, 'ustaw-tabulator').click();
    expect(zapis).toHaveLength(0);
    expect(odpowiedzPanelu(panel.element)).toContain('bez położenia');
  });

  it('niesie położenie, rodzaj, znak wiodący i zdjęcie', () => {
    const zapis: Zapis = [];
    const panel = panelStylu(zapis);
    pole(panel.element, 'Tabulator — położenie od lewego marginesu w milimetrach').value = '62.5';
    pole(panel.element, 'Rodzaj tabulatora').value = 'decimal';
    pole(panel.element, 'Znak wiodący').value = 'dot';
    pole(panel.element, 'Zdejmij tabulator z tego położenia').checked = true;
    czynnosc(panel.element, 'ustaw-tabulator').click();
    const zadanie = zapis[0]?.zadanie ?? {};
    expect(zadanie['positionMm']).toBe(62.5);
    expect(zadanie['kind']).toBe('decimal');
    expect(zadanie['leader']).toBe('dot');
    expect(zadanie['remove']).toBe(true);
  });
});

describe('listy, punktatory i numeracja poziomu', () => {
  it('lista niesie rodzaj, poziom i format numeracji', () => {
    const zapis: Zapis = [];
    const panel = panelList(zapis);
    pole(panel.element, 'Rodzaj listy').value = StudioListKind.Multilevel;
    pole(panel.element, 'Poziom listy').value = '2';
    pole(panel.element, 'Format numeracji').value = StudioListNumberFormat.Legal;
    pole(panel.element, 'Punkt startu numeracji').value = '3';
    czynnosc(panel.element, 'zastosuj-liste').click();

    const zadanie = zapis[0]?.zadanie ?? {};
    expect(zadanie['kind']).toBe(StudioListKind.Multilevel);
    expect(zadanie['level']).toBe(2);
    expect(zadanie['numberFormat']).toBe(StudioListNumberFormat.Legal);
    expect(zadanie['startAt']).toBe(3);
  });

  it('nastawa poziomu bez wskazanej listy jest odmawiana nazwanym powodem', () => {
    const zapis: Zapis = [];
    const panel = panelList(zapis);
    pole(panel.element, 'Poziom listy').value = '1';
    czynnosc(panel.element, 'zapisz-punktator').click();
    expect(zapis).toHaveLength(0);
    expect(odpowiedzPanelu(panel.element)).toContain('LISTY ZASTANEJ');
  });

  it('punktator niesie źródło znaku, znak, wcięcie i wyrównanie', () => {
    const zapis: Zapis = [];
    const panel = panelList(zapis);
    panel.pokazListy([{ id: 'lista-jeden', kind: StudioListKind.Bullet, levels: [] }]);
    pole(panel.element, 'Lista zastana, do której fragment dołączyć').value = 'lista-jeden';
    pole(panel.element, 'Poziom listy').value = '1';
    pole(panel.element, 'Źródło znaku wypunktowania').value = StudioBulletSource.Symbol;
    pole(panel.element, 'Znak wypunktowania').value = '§';
    pole(panel.element, 'Wcięcie poziomu w milimetrach').value = '7.5';
    pole(panel.element, 'Wyrównanie znaku albo numeru').value = StudioTextAlign.Right;
    czynnosc(panel.element, 'zapisz-punktator').click();

    const zadanie = zapis[0]?.zadanie ?? {};
    expect(zadanie['listId']).toBe('lista-jeden');
    expect(zadanie['level']).toBe(1);
    expect(zadanie['bulletSource']).toBe(StudioBulletSource.Symbol);
    expect(zadanie['bulletCharacter']).toBe('§');
    expect(zadanie['indentMm']).toBe(7.5);
    expect(zadanie['align']).toBe(StudioTextAlign.Right);
  });

  it('numeracja poziomu niesie wzór numeru — tym buduje się 1.1.2', () => {
    const zapis: Zapis = [];
    const panel = panelList(zapis);
    panel.pokazListy([{ id: 'lista-jeden', kind: StudioListKind.Multilevel }]);
    pole(panel.element, 'Lista zastana, do której fragment dołączyć').value = 'lista-jeden';
    pole(panel.element, 'Poziom listy').value = '3';
    pole(panel.element, 'Wzór numeru poziomu').value = '%1.%2.%3';
    czynnosc(panel.element, 'zapisz-numeracje-listy').click();
    expect(zapis[0]?.zadanie['pattern']).toBe('%1.%2.%3');
  });

  it('poziom listy przestawia się w dwie strony jedną komendą', () => {
    const zapis: Zapis = [];
    const panel = panelList(zapis);
    czynnosc(panel.element, 'poziom-glebiej').click();
    czynnosc(panel.element, 'poziom-wyzej').click();
    expect(zapis[0]?.zadanie['step']).toBe(1);
    expect(zapis[1]?.zadanie['step']).toBe(-1);
  });
});

describe('znaki specjalne i autozamiana', () => {
  it('kafel znaku wstawia go kodem, nie samym znakiem', () => {
    const zapis: Zapis = [];
    const panel = panelList(zapis);
    panel.pokazZnaki(
      [{ code: '00A0', character: ' ', name: 'twarda spacja', category: 'interpunkcyjne' }],
      ['interpunkcyjne'],
    );
    panel.element.querySelector<HTMLButtonElement>("[data-znak='00A0']")?.click();
    expect(zapis[0]?.zadanie['code']).toBe('00A0');
    expect(panel.element.textContent ?? '').toContain('Grupy znaków: interpunkcyjne');
  });

  it('szukanie znaku niesie frazę i zawężenie do ostatnio użytych', () => {
    const zapis: Zapis = [];
    const panel = panelList(zapis);
    pole(panel.element, 'Szukaj znaku po nazwie albo kodzie').value = 'paragraf';
    pole(panel.element, 'Tylko znaki ostatnio użyte').checked = true;
    czynnosc(panel.element, 'szukaj-znaku').click();
    expect(zapis[0]?.zadanie).toEqual({ fraza: 'paragraf', tylkoOstatnie: true });
  });

  it('tablica pusta mówi wprost, że rdzeń nic nie oddał', () => {
    const panel = panelList([]);
    panel.pokazZnaki([], []);
    expect(panel.element.textContent ?? '').toContain('ani jednego znaku');
  });

  it('zasada autozamiany bez skrótu jest odmawiana', () => {
    const zapis: Zapis = [];
    const panel = panelList(zapis);
    czynnosc(panel.element, 'zapisz-autozamiane').click();
    expect(zapis).toHaveLength(0);
    expect(odpowiedzPanelu(panel.element)).toContain('bez skrótu');
  });

  it('zastąpienie puste jedzie jako usunięcie zasady, nie jako pusty napis', () => {
    const zapis: Zapis = [];
    const panel = panelList(zapis);
    pole(panel.element, 'Skrót wpisywany przez Operatora').value = '(c)';
    czynnosc(panel.element, 'zapisz-autozamiane').click();
    const zadanie = zapis[0]?.zadanie ?? {};
    expect(zadanie['shortcut']).toBe('(c)');
    expect('replacement' in zadanie).toBe(false);
  });

  it('wykaz zasad wskazuje zasadę w polach', () => {
    const panel = panelList([]);
    panel.pokazAutozamiany([{ shortcut: '(c)', replacement: '©', enabled: true, builtin: true }]);
    const pozycja = panel.element.querySelector<HTMLElement>("[data-skrot='(c)']");
    expect(pozycja?.textContent ?? '').toContain('fabryczna');
    pozycja?.querySelector('button')?.click();
    expect(pole(panel.element, 'Co wchodzi w miejsce skrótu').value).toBe('©');
  });
});
