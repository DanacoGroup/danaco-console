// Zestawy żetonów okna: wykaz, zapis żetonów motywu obowiązującego, wydanie
// zapisu, wczytanie cudzego zapisu i wydanie przewodnika do modułu.

import { Command, DesignTokenKind, DesignTokenTarget, type DesignToken } from '../../../shared/contract.ts';
import { oglos } from './ogloszenie.ts';
import { chwila, poproszony, wykonany, type Kontekst, type Wiersz } from './design-wspolne.ts';
import { dolozPlik, dopiszDoKorzenia, pokazWykaz } from './design-wykaz.ts';

export function zwiazZetony(kontekst: Kontekst): void {
  kontekst.stan.odswiezenia.set('zetony', () => {
    void pokazZestawy(kontekst);
  });
  dopiszDoKorzenia({
    nazwa: 'Żetony i system projektowy',
    opis: 'zestawy, wydanie, wczytanie',
    otworz: (biezacy) => {
      void pokazZestawy(biezacy);
    },
  });
}

export async function pokazZestawy(kontekst: Kontekst): Promise<void> {
  if (kontekst.stan.idOkna === '') return;
  const odpowiedz = await poproszony(kontekst.kanal, Command.DesignTokensetList, {
    windowId: kontekst.stan.idOkna,
  });
  const zestawy = odpowiedz?.tokenSets ?? [];
  const wiersze: Wiersz[] = zestawy.map((zestaw) => ({
    tekst: zestaw.name,
    meta: `${String(zestaw.tokenCount)} żetonów · ${chwila(zestaw.updatedAt)}`,
    plakietka: zestaw.theme ?? '',
    kropka: zestaw.id === kontekst.stan.idZestawuZetonow ? 'sygnal' : 'neutralna',
    naKlik: () => {
      kontekst.stan.idZestawuZetonow = zestaw.id;
      void pokazZeton(kontekst, zestaw.id, zestaw.name, zestaw.tokens ?? []);
    },
  }));
  wiersze.unshift({
    tekst: 'Zapisz żetony motywu obowiązującego jako zestaw',
    kropka: 'sygnal',
    naKlik: () => {
      void zapisz(kontekst);
    },
  });
  pokazWykaz(kontekst, 'Zestawy żetonów', wiersze, 'Okno nie ma jeszcze zestawu żetonów.');
}

async function pokazZeton(
  kontekst: Kontekst,
  idZestawu: string,
  nazwa: string,
  zebrane: DesignToken[],
): Promise<void> {
  let zetony = zebrane;
  if (zetony.length === 0) {
    const odpowiedz = await poproszony(kontekst.kanal, Command.DesignTokensetList, {
      windowId: kontekst.stan.idOkna,
      tokenSetId: idZestawu,
    });
    zetony = odpowiedz?.tokenSets[0]?.tokens ?? [];
  }
  const wiersze: Wiersz[] = [
    {
      tekst: '◂ Zestawy żetonów',
      naKlik: () => {
        void pokazZestawy(kontekst);
      },
    },
    {
      tekst: 'Wydaj zapis zmiennych CSS',
      kropka: 'sygnal',
      naKlik: () => {
        void wydaj(kontekst, idZestawu);
      },
    },
    {
      tekst: 'Wydaj przewodnik do modułu wskazanego zasobem',
      kropka: 'sygnal',
      naKlik: () => {
        void wydajPrzewodnik(kontekst, idZestawu);
      },
    },
    ...zetony.map((zeton) => ({
      tekst: zeton.name,
      meta: zeton.value,
      plakietka: zeton.kind,
    } satisfies Wiersz)),
  ];
  pokazWykaz(kontekst, nazwa, wiersze, 'Zestaw nie niesie ani jednego żetonu.');
}

// Żetony bierze się z motywu obowiązującego w chwili naciśnięcia, nie z kopii
// zdjętej przy otwarciu okna — motyw jest własnością powłoki i bywa zmieniany.
function zetonyMotywu(dokument: Document): DesignToken[] {
  const zebrane = new Map<string, string>();
  const arkusze = dokument.styleSheets;
  for (let a = 0; a < arkusze.length; a += 1) {
    let reguly: CSSRuleList;
    try {
      // Arkusz z obcego źródła odmawia wglądu w reguły; taki pomijamy.
      reguly = arkusze[a].cssRules;
    } catch {
      continue;
    }
    for (let r = 0; r < reguly.length; r += 1) {
      const regula = reguly[r];
      if (!(regula instanceof CSSStyleRule)) continue;
      if (!regula.selectorText.includes(':root')) continue;
      for (let p = 0; p < regula.style.length; p += 1) {
        const nazwa = regula.style[p];
        if (!nazwa.startsWith('--dn-')) continue;
        zebrane.set(nazwa, regula.style.getPropertyValue(nazwa).trim());
      }
    }
  }
  return [...zebrane].map(([nazwa, wartosc]) => ({
    name: nazwa,
    kind: rodzajZetonu(nazwa),
    value: wartosc,
  }));
}

function rodzajZetonu(nazwa: string): DesignTokenKind {
  if (nazwa.startsWith('--dn-ff')) return DesignTokenKind.FontFamily;
  if (nazwa.startsWith('--dn-fw')) return DesignTokenKind.FontWeight;
  if (nazwa.startsWith('--dn-cien')) return DesignTokenKind.Shadow;
  if (nazwa.startsWith('--dn-czas') || nazwa.startsWith('--dn-ruch')) {
    return DesignTokenKind.Duration;
  }
  if (nazwa.startsWith('--dn-od') || nazwa.startsWith('--dn-r-') || nazwa.startsWith('--dn-fs')) {
    return DesignTokenKind.Dimension;
  }
  return DesignTokenKind.Color;
}

async function zapisz(kontekst: Kontekst): Promise<void> {
  const zetony = zetonyMotywu(kontekst.korzen.ownerDocument);
  if (zetony.length === 0) {
    oglos(
      'Design',
      'Motyw obowiązujący nie oddał ani jednego żetonu — zestaw pusty rdzeń odmówi.',
      'ostrzezenie',
    );
    return;
  }
  const odpowiedz = await wykonany(kontekst.kanal, Command.DesignTokensetSave, {
    windowId: kontekst.stan.idOkna,
    name: `Motyw okna ${chwila(Date.now())}`,
    tokens: zetony,
  }, 'Zapis zestawu żetonów');
  if (odpowiedz !== null) kontekst.stan.idZestawuZetonow = odpowiedz.tokenSet.id;
  await pokazZestawy(kontekst);
}

async function wydaj(kontekst: Kontekst, idZestawu: string): Promise<void> {
  const odpowiedz = await wykonany(kontekst.kanal, Command.DesignTokensetExport, {
    tokenSetId: idZestawu,
    target: DesignTokenTarget.CssVariables,
  }, 'Wydanie żetonów');
  if (odpowiedz === null) return;
  dolozPlik(kontekst, {
    fileName: odpowiedz.fileName,
    mediaType: 'text/css',
    contentBase64: btoa(unescape(encodeURIComponent(odpowiedz.content))),
  });
}

export async function wczytajZapis(kontekst: Kontekst, zapis: string, nazwa: string): Promise<void> {
  if (zapis.trim() === '') {
    oglos('Design', 'Pusty zapis nie niesie ani jednej roli.', 'ostrzezenie');
    return;
  }
  const odpowiedz = await wykonany(kontekst.kanal, Command.DesignTokensetImport, {
    windowId: kontekst.stan.idOkna,
    name: nazwa === '' ? 'Zestaw wczytany' : nazwa,
    contentBase64: btoa(unescape(encodeURIComponent(zapis))),
  }, 'Wczytanie zestawu żetonów');
  if (odpowiedz === null) return;
  const nieznane = odpowiedz.unknownNames ?? [];
  if (nieznane.length > 0) {
    oglos('Design', `Role spoza systemu produktu: ${nieznane.join(', ')}.`, 'ostrzezenie');
  }
  await pokazZestawy(kontekst);
}

async function wydajPrzewodnik(kontekst: Kontekst, idZestawu: string): Promise<void> {
  if (kontekst.stan.idOkna === '') return;
  const odpowiedz = await wykonany(kontekst.kanal, Command.DesignStyleguidePublish, {
    tokenSetId: idZestawu,
    targetModuleId: kontekst.stan.idModulu,
  }, 'Wydanie przewodnika');
  if (odpowiedz === null) return;
  kontekst.stan.idZasobu = odpowiedz.assetId;
  kontekst.stan.odswiezenia.get('zasoby')?.();
}
