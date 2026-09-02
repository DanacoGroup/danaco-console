// Prompt Builder: pola strukturalnego promptu, zamówienie zasobu, szablony
// promptu oraz wykaz zleceń generowania w panelu zadań.

import { Command, type DesignPrompt } from '../../../shared/contract.ts';
import { oglos } from './ogloszenie.ts';
import {
  chip,
  chwila,
  nieGotowe,
  panel,
  poproszony,
  przycisk,
  tresc,
  wykonany,
  wypelnijWykaz,
  znacznik,
  type Kontekst,
  type Wiersz,
} from './design-wspolne.ts';
import { pokazWykaz } from './design-wykaz.ts';

interface PolaPromptu {
  temat: HTMLInputElement | null;
  styl: HTMLSelectElement | null;
  paleta: HTMLSelectElement | null;
  proporcje: HTMLSelectElement | null;
  warianty: HTMLSelectElement | null;
  wykluczenia: HTMLInputElement | null;
  kreatywnosc: HTMLInputElement | null;
  miara: HTMLElement | null;
}

const OSTATNIE = new Map<string, DesignPrompt>();

export function ostatniPrompt(idKarty: string): DesignPrompt | undefined {
  return OSTATNIE.get(idKarty);
}

export function zwiazPrompt(kontekst: Kontekst): void {
  const okno = panel(kontekst.korzen, 'panel-prompt');
  const cialo = tresc(okno);
  if (okno === null || cialo === null) return;
  const pola = zbierzPola(cialo);
  wyzeruj(pola);
  const zamow = przycisk(cialo, 'Generuj') ?? cialo.querySelector<HTMLElement>('.dn-btn--sygnal');

  const odswiez = (): void => {
    void wczytaj(kontekst, pola);
  };
  kontekst.stan.odswiezenia.set('prompt', odswiez);

  for (const pole of [pola.temat, pola.wykluczenia, pola.kreatywnosc]) {
    pole?.addEventListener('input', () => {
      opiszMiare(pola);
    }, kontekst.przy);
  }

  if (zamow !== null) {
    zamow.textContent = 'Zamów warianty';
    zamow.removeAttribute('data-komunikat');
    zamow.removeAttribute('data-komunikat-tytul');
    zamow.addEventListener('click', () => {
      void zamowWarianty(kontekst, pola);
    }, kontekst.przy);
  }

  const szablony = chip(okno, 'Szablony');
  if (szablony !== null) {
    szablony.setAttribute('role', 'button');
    szablony.addEventListener('click', () => {
      void pokazSzablony(kontekst, pola);
    }, kontekst.przy);
  }

  odswiez();
  zwiazZadania(kontekst);
}

// Wartości pól prototypu są przykładem Właściciela; do odpowiedzi rdzenia pole
// ma stać puste, a lista wyboru — bez ani jednej pozycji zmyślonej.
function wyzeruj(pola: PolaPromptu): void {
  for (const pole of [pola.temat, pola.wykluczenia]) {
    if (pole !== null) pole.value = '';
  }
  for (const wybor of [pola.styl, pola.paleta, pola.proporcje, pola.warianty]) {
    wypelnijWybor(wybor, [], '');
  }
  if (pola.miara !== null) pola.miara.textContent = '';
}

function zbierzPola(cialo: HTMLElement): PolaPromptu {
  return {
    temat: cialo.querySelector<HTMLInputElement>('#dg-temat'),
    styl: cialo.querySelector<HTMLSelectElement>('#dg-styl'),
    paleta: cialo.querySelector<HTMLSelectElement>('#dg-paleta-barw'),
    proporcje: cialo.querySelector<HTMLSelectElement>('#dg-proporcje'),
    warianty: cialo.querySelector<HTMLSelectElement>('#dg-warianty'),
    wykluczenia: cialo.querySelector<HTMLInputElement>('#dg-wykluczenia-negative-prompt'),
    kreatywnosc: cialo.querySelector<HTMLInputElement>('#dg-kreatywnosc-wiernosc-0-7'),
    miara: cialo.querySelector<HTMLElement>('.dn-meta'),
  };
}

// Wykazu stylów, palet i proporcji kontrakt nie oddaje — jedynym źródłem
// wartości są prompty, które w tym oknie już poszły do rdzenia.
function wypelnijWybor(wybor: HTMLSelectElement | null, wartosci: string[], wybrana: string): void {
  if (wybor === null) return;
  const zestaw = [...new Set(wartosci.filter((wartosc) => wartosc !== ''))];
  if (zestaw.length === 0) {
    const pusta = wybor.ownerDocument.createElement('option');
    pusta.textContent = 'brak wcześniejszych wartości';
    pusta.value = '';
    wybor.replaceChildren(pusta);
    wybor.disabled = true;
    return;
  }
  wybor.disabled = false;
  wybor.replaceChildren(...zestaw.map((wartosc) => {
    const pozycja = wybor.ownerDocument.createElement('option');
    pozycja.value = wartosc;
    pozycja.textContent = wartosc;
    return pozycja;
  }));
  if (wybrana !== '') wybor.value = wybrana;
}

async function wczytaj(kontekst: Kontekst, pola: PolaPromptu): Promise<void> {
  if (kontekst.stan.idOkna === '') return;
  const historia = await poproszony(kontekst.kanal, Command.DesignPromptHistoryList, {
    windowId: kontekst.stan.idOkna,
    limit: 50,
  });
  const zapisy = historia?.prompts ?? [];
  const prompty = zapisy.map((zapis) => zapis.prompt);
  const biezacy = prompty[0];

  if (pola.temat !== null) pola.temat.value = biezacy?.subject ?? '';
  if (pola.wykluczenia !== null) pola.wykluczenia.value = biezacy?.exclusions ?? '';
  if (pola.kreatywnosc !== null) {
    pola.kreatywnosc.value = String(Math.round((biezacy?.creativity ?? 0.7) * 10));
  }
  wypelnijWybor(pola.styl, prompty.map((prompt) => prompt.style ?? ''), biezacy?.style ?? '');
  wypelnijWybor(pola.paleta, prompty.map((prompt) => prompt.palette ?? ''), biezacy?.palette ?? '');
  wypelnijWybor(
    pola.proporcje,
    prompty.map((prompt) => prompt.aspectRatio ?? ''),
    biezacy?.aspectRatio ?? '',
  );
  wypelnijWybor(
    pola.warianty,
    prompty.map((prompt) => (prompt.variants === undefined ? '' : String(prompt.variants))),
    biezacy?.variants === undefined ? '' : String(biezacy.variants),
  );
  if (biezacy !== undefined) OSTATNIE.set(kontekst.idKarty, biezacy);
  opiszMiare(pola, biezacy?.seed);
  wypelnijZadania(kontekst, zapisy.map((zapis) => ({
    tekst: zapis.prompt.subject,
    meta: chwila(zapis.createdAt),
    plakietka: `${String((zapis.assetIds ?? []).length)} zasobów`,
    kropka: (zapis.assetIds ?? []).length > 0 ? 'sukces' : 'neutralna',
  } satisfies Wiersz)));
}

function opiszMiare(pola: PolaPromptu, ziarno?: number): void {
  if (pola.miara === null) return;
  const dlugosc = zlozony(pola).length;
  pola.miara.textContent = ziarno === undefined
    ? `≈ ${String(dlugosc)} znaków`
    : `≈ ${String(dlugosc)} znaków · ziarno ${String(ziarno)}`;
}

function zlozony(pola: PolaPromptu): string {
  return [
    pola.temat?.value ?? '',
    pola.styl?.value ?? '',
    pola.paleta?.value ?? '',
    pola.wykluczenia?.value ?? '',
  ].filter((czesc) => czesc !== '').join(', ');
}

function zbierzPrompt(pola: PolaPromptu): DesignPrompt {
  const warianty = Number.parseInt(pola.warianty?.value ?? '', 10);
  const kreatywnosc = Number.parseInt(pola.kreatywnosc?.value ?? '', 10);
  return {
    subject: pola.temat?.value.trim() ?? '',
    ...(pola.styl?.value === undefined || pola.styl.value === '' ? {} : { style: pola.styl.value }),
    ...(pola.paleta?.value === undefined || pola.paleta.value === ''
      ? {}
      : { palette: pola.paleta.value }),
    ...(pola.proporcje?.value === undefined || pola.proporcje.value === ''
      ? {}
      : { aspectRatio: pola.proporcje.value }),
    ...(pola.wykluczenia?.value === undefined || pola.wykluczenia.value === ''
      ? {}
      : { exclusions: pola.wykluczenia.value }),
    ...(Number.isNaN(warianty) ? {} : { variants: warianty }),
    ...(Number.isNaN(kreatywnosc) ? {} : { creativity: kreatywnosc / 10 }),
  };
}

async function zamowWarianty(kontekst: Kontekst, pola: PolaPromptu): Promise<void> {
  const prompt = zbierzPrompt(pola);
  if (prompt.subject === '' || kontekst.stan.idOkna === '') {
    oglos('Design', 'Temat promptu nie może zostać pusty.', 'ostrzezenie');
    return;
  }
  OSTATNIE.set(kontekst.idKarty, prompt);
  const odpowiedz = await wykonany(kontekst.kanal, Command.DesignAssetGenerate, {
    windowId: kontekst.stan.idOkna,
    prompt,
    ...(kontekst.stan.idZasobu === '' ? {} : { referenceAssetId: kontekst.stan.idZasobu }),
  }, 'Zamówienie wariantów');
  if (odpowiedz === null) return;
  kontekst.stan.idZasobu = odpowiedz.assets[0]?.id ?? kontekst.stan.idZasobu;
  kontekst.stan.odswiezenia.get('zasoby')?.();
  kontekst.stan.odswiezenia.get('podglad')?.();
  kontekst.stan.odswiezenia.get('prompt')?.();
}

async function pokazSzablony(kontekst: Kontekst, pola: PolaPromptu): Promise<void> {
  if (kontekst.stan.idOkna === '') return;
  const odpowiedz = await poproszony(kontekst.kanal, Command.DesignPromptTemplateList, {
    windowId: kontekst.stan.idOkna,
  });
  const wiersze: Wiersz[] = (odpowiedz?.templates ?? []).map((szablon) => ({
    tekst: szablon.name,
    meta: szablon.prompt.subject,
    naKlik: () => {
      nanies(pola, szablon.prompt);
    },
  }));
  wiersze.unshift({
    tekst: 'Zapisz obecny prompt jako szablon',
    kropka: 'sygnal',
    naKlik: () => {
      void zapiszSzablon(kontekst, pola);
    },
  });
  pokazWykaz(kontekst, 'Szablony promptu', wiersze);
}

function nanies(pola: PolaPromptu, prompt: DesignPrompt): void {
  if (pola.temat !== null) pola.temat.value = prompt.subject;
  if (pola.wykluczenia !== null) pola.wykluczenia.value = prompt.exclusions ?? '';
  if (pola.styl !== null && prompt.style !== undefined) pola.styl.value = prompt.style;
  if (pola.paleta !== null && prompt.palette !== undefined) pola.paleta.value = prompt.palette;
  if (pola.proporcje !== null && prompt.aspectRatio !== undefined) {
    pola.proporcje.value = prompt.aspectRatio;
  }
  opiszMiare(pola, prompt.seed);
}

async function zapiszSzablon(kontekst: Kontekst, pola: PolaPromptu): Promise<void> {
  const prompt = zbierzPrompt(pola);
  if (prompt.subject === '') return;
  await wykonany(kontekst.kanal, Command.DesignPromptTemplateSave, {
    windowId: kontekst.stan.idOkna,
    name: prompt.subject.slice(0, 60),
    prompt,
  }, 'Zapis szablonu promptu');
  await pokazSzablony(kontekst, pola);
}

function zwiazZadania(kontekst: Kontekst): void {
  const okno = panel(kontekst.korzen, 'panel-zadania');
  const cialo = tresc(okno);
  if (cialo === null) return;
  // Postępu pojedynczego zlecenia kontrakt nie oddaje, więc pasek schodzi.
  cialo.querySelector('.dn-postep')?.remove();
  nieGotowe(cialo, 'Żadne zlecenie generowania nie poszło jeszcze z tego okna.');
}

function wypelnijZadania(kontekst: Kontekst, wiersze: Wiersz[]): void {
  const okno = panel(kontekst.korzen, 'panel-zadania');
  const cialo = tresc(okno);
  if (cialo === null) return;
  znacznik(okno, String(wiersze.length));
  wypelnijWykaz(
    kontekst,
    cialo,
    wiersze,
    'Żadne zlecenie generowania nie poszło jeszcze z tego okna.',
  );
}
