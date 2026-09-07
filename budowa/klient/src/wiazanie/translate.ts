// Wiązanie karty modułu Translate z rdzeniem: okno tłumaczenia, wykazy rdzenia
// i pasy czynności nad panelami. Węzły pochodzą ze znacznika Właściciela,
// treść z rdzenia.

import {
  Command,
  type ApprovalRecord,
  type GlossaryTerm,
  type QaProfile,
  type SegmentationRuleset,
  type TranslationMemoryEntry,
  type TranslationStep,
} from '../../../shared/contract.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import {
  cialoPanelu,
  opiszNaglowek,
  zapewnijOknoModulu,
  zdejmijSterowanieWspolne,
  wykazPanelu,
  zdejmijTrescWspolna,
} from './okno-modulu.ts';
import { zwiazKatalogModulu } from './katalog-modulu.ts';
import { czynnosciGlosariusza } from './translate-glosariusz.ts';
import { czynnosciJakosci } from './translate-jakosc.ts';
import { czynnosciNapisow } from './translate-napisy.ts';
import { czynnosciPamieci } from './translate-pamiec.ts';
import {
  NAGLOWEK,
  pasDzialan,
  przyjmijPanele,
  zapytaj,
  type CzynnosciTlumaczenia,
  type StanTlumaczenia,
} from './translate-wspolne.ts';
import { czynnosciWymiany } from './translate-wymiana.ts';
import { czynnosciZrodla } from './translate-zrodlo.ts';

const KOD_MODULU = 'translate';

interface WiazanieKarty {
  korzen: Element;
  odlaczenia: Odsubskrybuj[];
}

const WIAZANIA = new Map<string, WiazanieKarty>();

export function zwiazTlumaczenie(
  kanal: Kanal,
  nazwaSrodowiska: string,
  idOknaStojacego: string,
  wskazanieKorzenia: Element | string,
): boolean {
  const korzen = korzenKarty(wskazanieKorzenia);
  if (korzen === null) return false;
  const idKarty = korzen.getAttribute('data-karta') ?? '';
  if (idKarty === '') return false;
  if (WIAZANIA.get(idKarty)?.korzen === korzen) return false;
  zwolnijTlumaczenie(idKarty);

  const odlaczenia: Odsubskrybuj[] = [];
  const sterowanie = new AbortController();
  const przy = { signal: sterowanie.signal };
  odlaczenia.push(() => {
    sterowanie.abort();
  });
  WIAZANIA.set(idKarty, { korzen, odlaczenia });

  if (korzen instanceof HTMLElement) {
    zdejmijTrescWspolna(korzen);
    zdejmijSterowanieWspolne(korzen);
  }
  zdejmijTrescPrzykladowa(korzen);
  void opiszNaglowek(kanal, nazwaSrodowiska, korzen);

  let idOkna = idOknaStojacego;
  const odswiez = async (): Promise<void> => {
    if (idOkna === '') return;
    await Promise.all([
      wypelnijGlosariusz(kanal, korzen),
      wykaz(kanal, korzen, 'panel-zadania', Command.TranslateStepList,
        'Rdzeń nie podaje kroków tłumaczenia.',
        (o) => (o.steps as TranslationStep[]).map((s) => [s.name, s.kind] as const)),
      wykaz(kanal, korzen, 'panel-obszar-tlum', Command.TranslateMemoryList,
        'Pamięć tłumaczeń jest pusta.',
        (o) => (o.entries as TranslationMemoryEntry[]).map((e) => [e.sourceSegment, e.language] as const)),
      wykaz(kanal, korzen, 'panel-plan', Command.TranslateQaProfileList,
        'Rdzeń nie ma profilu kontroli jakości.',
        (o) => (o.profiles as QaProfile[]).map((q) => [q.name, q.scope ?? ''] as const)),
      wykaz(kanal, korzen, 'panel-kolejka', Command.TranslateApprovalList,
        'Żadna zgoda nie została jeszcze wydana.',
        (o) => (o.records as ApprovalRecord[]).map((r) => [r.stage, r.author ?? ''] as const)),
      wykaz(kanal, korzen, 'panel-pliki', Command.TranslateSegmentationRulesList,
        'Rdzeń nie ma zestawu reguł segmentacji.',
        (o) => (o.rulesets as SegmentationRuleset[]).map((r) => [r.name, r.language ?? ''] as const)),
    ]);
  };

  const stan: StanTlumaczenia = {
    kanal,
    korzen,
    idOkna: () => idOkna,
    odswiez,
    panele: [],
    idDokumentu: '',
    idZasobu: '',
    idPakietu: '',
  };
  zalozPasyDzialan(korzen);
  const czynnosci: CzynnosciTlumaczenia = {
    ...czynnosciZrodla(stan),
    ...czynnosciPamieci(stan),
    ...czynnosciGlosariusza(stan),
    ...czynnosciJakosci(stan),
    ...czynnosciNapisow(stan),
    ...czynnosciWymiany(stan),
  };

  void (async (): Promise<void> => {
    idOkna = await zapewnijOknoModulu(kanal, idKarty, KOD_MODULU, 'Translate', idOkna);
    if (idOkna === '') {
      niegotowe(panel(korzen, 'panel-glossary'), 'Rdzeń nie dał okna tłumaczenia dla tej karty.');
      return;
    }
    await odswiez();
  })();

  korzen.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    if (cel.closest('#panel-glossary .sta-okno-akcje .dn-btn-ikona') !== null) {
      zdarzenie.stopPropagation();
      void wniesPanel(stan).then(odswiez);
      return;
    }
    const przycisk = cel.closest<HTMLElement>('[data-czynnosc-tlumaczenia]');
    if (przycisk === null) return;
    zdarzenie.stopPropagation();
    const czynnosc = czynnosci[przycisk.dataset.czynnoscTlumaczenia ?? ''];
    if (czynnosc !== undefined) void czynnosc();
  }, przy);

  const katalog = zwiazKatalogModulu(kanal, idOkna, korzen, KOD_MODULU, 'Tłumaczenie');
  if (katalog !== null) odlaczenia.push(katalog);
  return true;
}

export function zwolnijTlumaczenie(idKarty: string): void {
  const wiazanie = WIAZANIA.get(idKarty);
  if (wiazanie === undefined) return;
  for (const odlacz of wiazanie.odlaczenia) odlacz();
  WIAZANIA.delete(idKarty);
}

/* Pas czynności stoi przy panelu, którego wykazu dotyczy: pamięć i źródło przy
   obszarze tłumaczenia, terminy przy glosariuszu, kontrola przy planie. */
function zalozPasyDzialan(korzen: Element): void {
  pasDzialan(korzen, 'panel-obszar-tlum', 'Źródło i segmenty', [
    ['zrodlo-wczytaj', 'Wczytaj źródło'],
    ['zrodlo-jezyk', 'Rozpoznaj język'],
    ['zrodlo-segmentuj', 'Segmentuj'],
    ['segment-scal', 'Scal segmenty'],
    ['segment-podziel', 'Podziel segment'],
    ['tlumaczenie-zapisz', 'Zapisz przekład'],
    ['panel-ton', 'Ton panelu'],
    ['reguly-segmentacji', 'Reguły segmentacji'],
  ]);
  pasDzialan(korzen, 'panel-obszar-tlum', 'Pamięć tłumaczeń', [
    ['pamiec-podpowiedz', 'Podpowiedzi'],
    ['pamiec-wstepne', 'Tłumacz wstępnie'],
    ['pamiec-zapisz', 'Zapisz parę'],
    ['pamiec-usun', 'Usuń parę'],
    ['pamiec-wciagnij', 'Wciągnij pamięć'],
    ['pamiec-wydaj', 'Wydaj pamięć'],
    ['pamiec-wyrownaj', 'Wyrównaj teksty'],
    ['pamiec-uporzadkuj', 'Uporządkuj pamięć'],
    ['pamiec-polityka-odczyt', 'Polityka pamięci'],
    ['pamiec-polityka-zapis', 'Zapisz politykę'],
  ]);
  pasDzialan(korzen, 'panel-glossary', 'Glosariusz', [
    ['glosariusz-zapisz', 'Zapisz termin'],
    ['glosariusz-wystapienia', 'Wystąpienia'],
    ['glosariusz-ujednolic', 'Ujednolić'],
    ['glosariusz-wciagnij', 'Wciągnij glosariusz'],
    ['glosariusz-wydaj', 'Wydaj glosariusz'],
    ['terminy-wydobadz', 'Wydobądź terminy'],
  ]);
  pasDzialan(korzen, 'panel-plan', 'Kontrola jakości', [
    ['jakosc-sprawdz', 'Sprawdź jakość'],
    ['spojnosc-sprawdz', 'Sprawdź spójność'],
    ['korekta-uruchom', 'Uruchom korektę'],
    ['korekta-zastosuj', 'Zastosuj korektę'],
    ['zwrotne-tlumaczenie', 'Tłumaczenie zwrotne'],
    ['zatwierdzenie-ustaw', 'Ustaw etap'],
    ['profil-jakosci-zapisz', 'Zapisz profil'],
    ['profil-jakosci-usun', 'Usuń profil'],
  ]);
  pasDzialan(korzen, 'panel-pliki', 'Dokumenty i zasoby', [
    ['dokument-wczytaj', 'Wczytaj dokument'],
    ['dokument-zloz', 'Złóż dokument'],
    ['dokument-uklad', 'Porównaj układ'],
    ['xliff-wciagnij', 'Wciągnij XLIFF'],
    ['zasob-wciagnij', 'Wciągnij zasób'],
    ['zasob-wydaj', 'Wydaj zasób'],
    ['zasob-formy-mnogie', 'Formy mnogie'],
    ['zasob-kontekst-klucza', 'Kontekst klucza'],
  ]);
  pasDzialan(korzen, 'panel-pliki', 'Napisy i dubbing', [
    ['napisy-wciagnij', 'Wciągnij napisy'],
    ['napisy-wydaj', 'Wydaj napisy'],
    ['napisy-taktowanie', 'Sprawdź taktowanie'],
    ['dubbing-skrypt', 'Skrypt dubbingu'],
  ]);
  pasDzialan(korzen, 'panel-kolejka', 'Przekład i wymiana', [
    ['profile-przekladu', 'Profile przekładu'],
    ['profil-przekladu-zapisz', 'Zapisz profil przekładu'],
    ['przeklady-porownaj', 'Porównaj przekłady'],
    ['pivot-odczyt', 'Język pośredni'],
    ['pivot-zapis', 'Zapisz język pośredni'],
    ['zlecenie-wsadowe', 'Zlecenie wsadowe'],
    ['panel-wydaj', 'Wydaj panel'],
    ['mowa-synteza', 'Synteza mowy'],
    ['artefakt-opublikuj', 'Opublikuj artefakt'],
    ['pakiet-zloz', 'Złóż pakiet'],
    ['pakiet-przyjmij', 'Przyjmij zwrot'],
    ['most-przyjmij-zrodlo', 'Przyjmij ze Studia'],
    ['most-odeslij-wynik', 'Odeślij do Studia'],
  ]);
}

function korzenKarty(wskazanie: Element | string): Element | null {
  if (typeof wskazanie !== 'string') return wskazanie;
  return document.querySelector(`.cd-tresc--modul[data-karta="${wskazanie}"]`);
}

function panel(korzen: Element, identyfikator: string): Element | null {
  return korzen.querySelector(`#${identyfikator} .sta-okno-tresc`);
}

function niegotowe(cialo: Element | null, zdanie: string): void {
  if (cialo === null) return;
  const napis = cialo.ownerDocument.createElement('div');
  napis.className = 'dn-meta';
  napis.textContent = zdanie;
  cialo.replaceChildren(napis);
}

/* Znacznik niesie źródła, ustalenia i miary wpisane wprost. Schodzą przed
   pierwszym pytaniem rdzenia, żeby okno nie pokazywało cudzego badania. */
function zdejmijTrescPrzykladowa(korzen: Element): void {
  niegotowe(panel(korzen, 'panel-glossary'), 'Wykaz czeka na odpowiedź rdzenia.');
  /* Znacznik niesie tekst źródłowy, gotowe tłumaczenia i miary segmentów
     wpisane wprost. Rdzeń ich nie odtworzy, więc schodzą razem z panelami. */
  for (const segment of korzen.querySelectorAll('.tr-seg, .tr-kol-tresc > *')) segment.remove();
  for (const miara of korzen.querySelectorAll('.tr-glowa span:not(.tr-kod)')) miara.remove();
  // Panele języków i wykaz plików prototyp wpisuje wprost; rdzeń poda swoje.
  for (const kolumna of korzen.querySelectorAll('.tr-kol')) kolumna.remove();
  niegotowe(panel(korzen, 'panel-pliki'), 'Rdzeń nie podaje plików tłumaczenia dla tego okna.');
  niegotowe(panel(korzen, 'panel-obszar-tlum'), 'Panele tłumaczenia czekają na wskazanie języka.');
  niegotowe(panel(korzen, 'panel-plan'), 'Plan tłumaczenia czeka na pierwszy panel języka.');
  niegotowe(panel(korzen, 'panel-kolejka'), 'Rdzeń nie podaje zadań w tle dla tego okna.');
  niegotowe(panel(korzen, 'panel-zadania'), 'Rdzeń nie podaje zadań tłumaczenia dla tego okna.');
}

async function wypelnijGlosariusz(kanal: Kanal, korzen: Element): Promise<void> {
  const cialo = panel(korzen, 'panel-glossary');
  if (cialo === null) return;
  const odpowiedz = await wywolaj(kanal, Command.TranslateGlossaryList, {});
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    niegotowe(cialo, odpowiedz.blad?.message ?? 'Wykaz glosariusza nie doszedł.');
    return;
  }
  const wpisy: GlossaryTerm[] = odpowiedz.wynik.terms;
  if (wpisy.length === 0) {
    niegotowe(cialo, 'Glosariusz jest pusty — żadnego terminu jeszcze nie wniesiono.');
    return;
  }
  cialo.replaceChildren(...wpisy.map((w) => pozycja(cialo, w.source, w.target ?? w.language)));
}

/* Panel wykazu bez własnego kształtu: komenda bez pól wymaganych, pozycje
   z odpowiedzi rdzenia, zdanie przy pustce. */
async function wykaz(
  kanal: Kanal,
  korzen: Element,
  panelId: string,
  komenda: Parameters<typeof wywolaj>[1],
  pusty: string,
  mapuj: (wynik: Record<string, unknown>) => readonly (readonly [string, string])[],
): Promise<void> {
  const cialo = cialoPanelu(korzen, panelId);
  if (cialo === null) return;
  const odpowiedz = await wywolaj(kanal, komenda as never, {} as never);
  const wynik = odpowiedz.wynik as Record<string, unknown> | undefined;
  wykazPanelu(cialo, odpowiedz.udany, odpowiedz.blad?.message,
    wynik === undefined ? undefined : [...mapuj(wynik)], pusty, (p) => p);
}

function pozycja(cialo: Element, tytul: string, podpis: string): HTMLElement {
  const wiersz = cialo.ownerDocument.createElement('div');
  wiersz.className = 'dn-wykaz-modulu-poz';
  const nazwa = cialo.ownerDocument.createElement('span');
  nazwa.textContent = tytul;
  wiersz.append(nazwa);
  if (podpis !== '') {
    const meta = cialo.ownerDocument.createElement('span');
    meta.className = 'dn-meta';
    meta.textContent = podpis;
    wiersz.append(meta);
  }
  return wiersz;
}

/* Okna zakładania źródła wydanie nie niesie: tytuł wchodzi zapytaniem w oknie,
   tak samo jak nazwa komponentu w Centrum. */
async function wniesPanel(stan: StanTlumaczenia): Promise<void> {
  const idOkna = stan.idOkna();
  if (idOkna === '') {
    oglos(NAGLOWEK, 'Okno tłumaczenia jeszcze nie powstało, więc nie ma do czego dołożyć panelu.');
    return;
  }
  const jezyk = zapytaj('Język panelu (kod, np. DE)');
  if (jezyk === '') return;
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateTargetAdd, {
    windowId: idOkna,
    language: jezyk,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    oglos(NAGLOWEK, odpowiedz.blad?.message ?? 'Rdzeń odmówił dołożenia panelu.', 'blad');
    return;
  }
  przyjmijPanele(stan, [odpowiedz.wynik.panel]);
  oglos(NAGLOWEK, `Panel języka ${jezyk} dołożony.`);
}
