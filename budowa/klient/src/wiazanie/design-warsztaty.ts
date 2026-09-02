// Warsztaty Design: sześć grup czynności katalogu wraz z ich polami wchodzi
// w panel roboczy prototypu, a wynik wywołania wraca w ten sam panel.

import type { Wynik } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import {
  panel,
  rozmiar,
  tresc,
  type Kontekst,
  type Wiersz,
} from './design-wspolne.ts';
import { dolozPlik, pokazWykaz, pozycjeKorzenia } from './design-wykaz.ts';
import {
  CZYNNOSCI_DESIGNU,
  type CzynnoscDesignu,
  type GrupaWarsztatu,
  type PoleCzynnosci,
} from './design-czynnosci.ts';

const WARSZTATY: readonly { grupa: GrupaWarsztatu; nazwa: string }[] = [
  { grupa: 'fotografia', nazwa: 'Fotografia i barwa' },
  { grupa: 'wektor', nazwa: 'Wektor, ikony i kroje' },
  { grupa: 'druk', nazwa: 'Druk i wielki format' },
  { grupa: 'bazy', nazwa: 'Bazy zasobów' },
  { grupa: 'publikacja', nazwa: 'Publikacje i kampanie' },
  { grupa: 'makieta', nazwa: 'Makiety i komponenty' },
];

// Klucze, pod którymi rdzeń oddaje wykazy — jeden odczyt zamiast tylu gałęzi,
// ile jest rodzin odpowiedzi.
const WYKAZY = [
  'assets', 'boards', 'collections', 'versions', 'annotations', 'tokenSets', 'paths',
  'symbols', 'frames', 'components', 'templates', 'icons', 'presets', 'profiles',
  'sizes', 'edits', 'issues', 'colors', 'pairs', 'glyphs', 'tiles', 'prompts', 'links',
  'files', 'paperSizes', 'devicePresets', 'sets', 'providersQueried',
];

export function pokazWarsztaty(kontekst: Kontekst): void {
  const wiersze: Wiersz[] = WARSZTATY.map((warsztat) => {
    const ile = CZYNNOSCI_DESIGNU.filter((czynnosc) => czynnosc.grupa === warsztat.grupa).length;
    return {
      tekst: warsztat.nazwa,
      meta: `${String(ile)} czynności`,
      kropka: 'neutralna',
      naKlik: () => {
        pokazCzynnosci(kontekst, warsztat.grupa, warsztat.nazwa);
      },
    };
  });
  for (const pozycja of pozycjeKorzenia()) {
    wiersze.push({
      tekst: pozycja.nazwa,
      meta: pozycja.opis,
      kropka: 'neutralna',
      naKlik: () => {
        pozycja.otworz(kontekst);
      },
    });
  }
  pokazWykaz(kontekst, 'Warsztaty Design', wiersze);
}

function pokazCzynnosci(kontekst: Kontekst, grupa: GrupaWarsztatu, nazwa: string): void {
  const wiersze: Wiersz[] = [{
    tekst: '◂ Warsztaty',
    naKlik: () => {
      pokazWarsztaty(kontekst);
    },
  }];
  for (const czynnosc of CZYNNOSCI_DESIGNU.filter((pozycja) => pozycja.grupa === grupa)) {
    wiersze.push({
      tekst: czynnosc.nazwa,
      meta: String(czynnosc.komenda),
      naKlik: () => {
        otworzCzynnosc(kontekst, czynnosc, grupa, nazwa);
      },
    });
  }
  pokazWykaz(kontekst, nazwa, wiersze);
}

function otworzCzynnosc(
  kontekst: Kontekst,
  czynnosc: CzynnoscDesignu,
  grupa: GrupaWarsztatu,
  nazwaGrupy: string,
): void {
  const okno = panel(kontekst.korzen, 'panel-plan');
  const cialo = tresc(okno);
  if (okno === null || cialo === null) return;
  const naglowek = okno.querySelector('.sta-okno-tytul b');
  if (naglowek !== null) naglowek.textContent = czynnosc.nazwa;

  const dokument = cialo.ownerDocument;
  const powrot = dokument.createElement('div');
  powrot.className = 'dn-wykaz-modulu-poz';
  powrot.setAttribute('role', 'button');
  powrot.textContent = `◂ ${nazwaGrupy}`;
  powrot.addEventListener('click', () => {
    pokazCzynnosci(kontekst, grupa, nazwaGrupy);
  }, kontekst.przy);

  const kontrolki = new Map<string, HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>();
  const pola = czynnosc.pola.map((pole) => zbudujPole(dokument, czynnosc, pole, kontrolki, kontekst));

  const stopka = dokument.createElement('div');
  stopka.className = 'dn-wykaz-modulu-poz';
  const wykonaj = dokument.createElement('button');
  wykonaj.type = 'button';
  wykonaj.className = 'dn-btn dn-btn--sygnal dn-btn--sm';
  wykonaj.style.marginLeft = 'auto';
  wykonaj.textContent = 'Wykonaj';
  wykonaj.addEventListener('click', () => {
    void wykonajCzynnosc(kontekst, czynnosc, kontrolki);
  }, kontekst.przy);
  stopka.append(wykonaj);

  cialo.replaceChildren(powrot, ...pola, stopka);
}

function zbudujPole(
  dokument: Document,
  czynnosc: CzynnoscDesignu,
  pole: PoleCzynnosci,
  kontrolki: Map<string, HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>,
  kontekst: Kontekst,
): HTMLElement {
  const opakowanie = dokument.createElement('div');
  opakowanie.className = 'dg-pole';
  const identyfikator = `dg-${String(czynnosc.komenda).replaceAll('.', '-')}-${pole.kod}`;
  const etykieta = dokument.createElement('label');
  etykieta.setAttribute('for', identyfikator);
  etykieta.textContent = pole.wymagane === true ? `${pole.etykieta} *` : pole.etykieta;

  let kontrolka: HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement;
  if (pole.rodzaj === 'wybor') {
    const wybor = dokument.createElement('select');
    wybor.className = 'dn-wybor';
    for (const pozycja of pole.pozycje ?? []) {
      const opcja = dokument.createElement('option');
      opcja.value = pozycja.wartosc;
      opcja.textContent = pozycja.etykieta;
      wybor.append(opcja);
    }
    kontrolka = wybor;
  } else if (pole.rodzaj === 'wielowiersz') {
    const obszar = dokument.createElement('textarea');
    obszar.className = 'dn-pole-kontrolka';
    obszar.rows = 3;
    kontrolka = obszar;
  } else if (pole.rodzaj === 'logiczne') {
    const przelacznik = dokument.createElement('input');
    przelacznik.type = 'checkbox';
    przelacznik.className = 'dn-check';
    kontrolka = przelacznik;
  } else {
    const wejscie = dokument.createElement('input');
    wejscie.className = 'dn-pole-kontrolka';
    wejscie.type = pole.rodzaj === 'liczba' ? 'number' : 'text';
    if (pole.rodzaj === 'zasob') wejscie.value = kontekst.stan.idZasobu;
    kontrolka = wejscie;
  }
  kontrolka.id = identyfikator;
  kontrolki.set(pole.kod, kontrolka);
  opakowanie.append(etykieta, kontrolka);
  return opakowanie;
}

function wartoscPola(
  pole: PoleCzynnosci,
  kontrolka: HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement,
): unknown {
  if (pole.postac === 'logiczne') {
    return kontrolka instanceof HTMLInputElement && kontrolka.checked ? true : undefined;
  }
  const surowa = kontrolka.value.trim();
  if (surowa === '') return undefined;
  if (pole.postac === 'liczba') {
    const liczba = Number(surowa);
    return Number.isFinite(liczba) ? liczba : undefined;
  }
  if (pole.postac === 'lista') {
    const czesci = surowa.split(',').map((czesc) => czesc.trim()).filter((czesc) => czesc !== '');
    return czesci.length === 0 ? undefined : czesci;
  }
  if (pole.postac === 'listaLiczb') {
    const czesci = surowa.split(',').map((czesc) => Number(czesc.trim()))
      .filter((czesc) => Number.isFinite(czesc));
    return czesci.length === 0 ? undefined : czesci;
  }
  if (pole.postac === 'obszary') {
    const obszary = surowa.split('\n').map((wiersz) => wiersz.split(',')
      .map((czesc) => Number(czesc.trim())).filter((czesc) => Number.isFinite(czesc)))
      .filter((liczby) => liczby.length >= 4)
      .map((liczby) => ({ x: liczby[0], y: liczby[1], width: liczby[2], height: liczby[3] }));
    return obszary.length === 0 ? undefined : obszary;
  }
  if (pole.postac === 'json') {
    try {
      return JSON.parse(surowa) as unknown;
    } catch {
      return undefined;
    }
  }
  return surowa;
}

async function wykonajCzynnosc(
  kontekst: Kontekst,
  czynnosc: CzynnoscDesignu,
  kontrolki: Map<string, HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>,
): Promise<void> {
  const zadanie: Record<string, unknown> = {};
  const gniazda = new Map<string, Record<string, unknown>>();
  for (const pole of czynnosc.pola) {
    const kontrolka = kontrolki.get(pole.kod);
    if (kontrolka === undefined) continue;
    const wartosc = wartoscPola(pole, kontrolka);
    if (wartosc === undefined) {
      if (pole.wymagane === true) {
        oglos('Design', `${czynnosc.nazwa}: pole „${pole.etykieta}" zostaje puste.`, 'ostrzezenie');
        return;
      }
      continue;
    }
    if (pole.gniazdo === undefined) {
      zadanie[pole.kod] = wartosc;
      continue;
    }
    const gniazdo = gniazda.get(pole.gniazdo) ?? {};
    gniazdo[pole.kod] = wartosc;
    gniazda.set(pole.gniazdo, gniazdo);
  }
  for (const [nazwa, zawartosc] of gniazda) zadanie[nazwa] = zawartosc;
  if (czynnosc.okno) {
    if (kontekst.stan.idOkna === '') {
      oglos('Design', 'Czynność wymaga okna modułu, którego rdzeń nie założył.', 'ostrzezenie');
      return;
    }
    zadanie['windowId'] = kontekst.stan.idOkna;
  }

  // Komenda jest wskazana danymi katalogu, więc jej żądanie nie da się zawęzić
  // do jednego typu w miejscu wywołania.
  const wynik = await wywolaj(
    kontekst.kanal,
    czynnosc.komenda,
    zadanie as never,
  ) as unknown as Wynik<unknown>;
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos('Design', `${czynnosc.nazwa}: ${wynik.blad?.message ?? 'rdzeń odmówił bez opisu.'}`, 'blad');
    return;
  }
  oglos('Design', `${czynnosc.nazwa} — wykonane.`, 'informacja');
  opiszOdpowiedz(kontekst, czynnosc, wynik.wynik);
}

function opiszOdpowiedz(kontekst: Kontekst, czynnosc: CzynnoscDesignu, odpowiedz: unknown): void {
  if (typeof odpowiedz !== 'object' || odpowiedz === null) return;
  const zapis = odpowiedz as Record<string, unknown>;
  zbierzPliki(kontekst, zapis);
  const klucz = WYKAZY.find((nazwa) => Array.isArray(zapis[nazwa]));
  if (klucz === undefined) {
    kontekst.stan.odswiezenia.get('zasoby')?.();
    kontekst.stan.odswiezenia.get('plansza')?.();
    return;
  }
  const pozycje = zapis[klucz] as unknown[];
  const wiersze: Wiersz[] = pozycje.slice(0, 60).map((pozycja) => ({
    tekst: opisPozycji(pozycja),
    meta: '',
  }));
  wiersze.unshift({
    tekst: `◂ ${czynnosc.nazwa}`,
    naKlik: () => {
      pokazWarsztaty(kontekst);
    },
  });
  pokazWykaz(kontekst, `${czynnosc.nazwa} — wynik`, wiersze, 'Rdzeń oddał wykaz pusty.');
}

function zbierzPliki(kontekst: Kontekst, zapis: Record<string, unknown>): void {
  const pliki = Array.isArray(zapis['files']) ? (zapis['files'] as unknown[]) : [];
  for (const plik of pliki) {
    if (typeof plik !== 'object' || plik === null) continue;
    const wpis = plik as Record<string, unknown>;
    if (typeof wpis['contentBase64'] !== 'string' || typeof wpis['fileName'] !== 'string') continue;
    dolozPlik(kontekst, {
      fileName: wpis['fileName'],
      mediaType: typeof wpis['mediaType'] === 'string' ? wpis['mediaType'] : 'application/octet-stream',
      contentBase64: wpis['contentBase64'],
      ...(typeof wpis['sizeBytes'] === 'number' ? { sizeBytes: wpis['sizeBytes'] } : {}),
    });
  }
  if (typeof zapis['contentBase64'] !== 'string' || typeof zapis['fileName'] !== 'string') return;
  dolozPlik(kontekst, {
    fileName: zapis['fileName'],
    mediaType: typeof zapis['mediaType'] === 'string' ? zapis['mediaType'] : 'application/octet-stream',
    contentBase64: zapis['contentBase64'],
    ...(typeof zapis['sizeBytes'] === 'number' ? { sizeBytes: zapis['sizeBytes'] } : {}),
  });
}

function opisPozycji(pozycja: unknown): string {
  if (typeof pozycja === 'string') return pozycja;
  if (typeof pozycja !== 'object' || pozycja === null) return String(pozycja);
  const zapis = pozycja as Record<string, unknown>;
  for (const klucz of ['name', 'title', 'message', 'fileName', 'text', 'label', 'hex', 'id']) {
    const wartosc = zapis[klucz];
    if (typeof wartosc === 'string' && wartosc !== '') {
      const rozmiarPliku = zapis['sizeBytes'];
      return typeof rozmiarPliku === 'number'
        ? `${wartosc} · ${rozmiar(rozmiarPliku)}`
        : wartosc;
    }
  }
  return JSON.stringify(zapis).slice(0, 120);
}
