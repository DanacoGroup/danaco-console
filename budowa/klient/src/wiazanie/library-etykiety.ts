// Panele „Tags & Collections" i „Kolejka" wobec słownika etykiet rdzenia.
// Prototyp nie ma pól tekstowych, więc nazwę nowej etykiety niesie edycja
// wpisana w istniejącą plakietkę, a przypisanie — upuszczenie wiersza pliku.
import {
  Command,
  LibraryDuplicateKind,
  LibraryThesaurusRelation,
  type LibraryCollection,
  type LibraryRule,
  type LibrarySuggestion,
  type LibraryTag,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import {
  etykietaGrupy,
  panel,
  plakietka,
  pozycja,
  pustka,
  wartosc,
  zdejmijDzieci,
  type Kontekst,
} from './library-wspolne.ts';

const GRANICA = 200;
const RODZAJ_PLIKU = 'text/dn-library-plik';
const RODZAJ_ETYKIETY = 'text/dn-library-etykieta';

export function zwiazEtykiety(
  kanal: Kanal,
  korzen: ParentNode,
  kontekst: Kontekst,
  przy: AddEventListenerOptions,
): void {
  const trescMozeMoze = panel(korzen, 'panel-tags');
  if (trescMozeMoze === null) return;
  const trescMoze = trescMozeMoze;
  const tresc = trescMoze;
  const dokument = tresc.ownerDocument;
  const przyciskSugestii = tresc.querySelector('.dn-btn');
  const podpisSugestii = przyciskSugestii?.textContent?.trim() ?? '';
  const kolekcje = new Map<HTMLElement, LibraryCollection>();
  let reguly: LibraryRule[] = [];

  async function odswiez(): Promise<void> {
    const [odpowiedzKolekcji, odpowiedzEtykiet, odpowiedzRegul, odpowiedzSugestii] =
      await Promise.all([
        wywolaj(kanal, Command.LibraryCollectionList, { limit: GRANICA }),
        wywolaj(kanal, Command.LibraryTagList, { limit: GRANICA }),
        wywolaj(kanal, Command.LibraryRuleList, { enabledOnly: false }),
        wywolaj(kanal, Command.LibrarySuggestionList, { limit: GRANICA }),
      ]);
    reguly = wartosc(odpowiedzRegul)?.rules ?? [];
    const wykazKolekcji = wartosc(odpowiedzKolekcji)?.collections ?? [];
    const wykazEtykiet = wartosc(odpowiedzEtykiet)?.tags ?? [];
    const liczbaSugestii = wartosc(odpowiedzSugestii)?.total ?? 0;
    kolekcje.clear();
    zdejmijDzieci(tresc);
    tresc.appendChild(etykietaGrupy(dokument, 'Kolekcje inteligentne'));
    if (wykazKolekcji.length === 0) {
      tresc.appendChild(nota(dokument, 'Rdzeń nie zwrócił kolekcji.'));
    }
    for (const kolekcja of wykazKolekcji) {
      const wiersz = zbudujKolekcje(dokument, kolekcja);
      kolekcje.set(wiersz, kolekcja);
      tresc.appendChild(wiersz);
    }
    tresc.appendChild(etykietaGrupy(dokument, 'Słownik etykiet'));
    tresc.appendChild(zbudujSlownik(dokument, wykazEtykiet));
    if (przyciskSugestii !== null) {
      przyciskSugestii.textContent = podpisPrzycisku(podpisSugestii, liczbaSugestii);
      tresc.appendChild(przyciskSugestii);
    }
  }

  tresc.addEventListener(
    'click',
    (zdarzenie) => {
      const cel = zdarzenie.target;
      if (!(cel instanceof Element)) return;
      if (cel.closest('.dn-btn') !== null) {
        void sklasyfikuj(kanal, kontekst).then(odswiez);
        return;
      }
      const wiersz = cel.closest('.dn-wykaz-modulu-poz');
      const kolekcja = wiersz instanceof HTMLElement ? kolekcje.get(wiersz) : undefined;
      if (kolekcja === undefined) return;
      if (zdarzenie.shiftKey) {
        void zdejmijRegule(kanal, kolekcja, reguly).then(odswiez);
        return;
      }
      if (zdarzenie.altKey) {
        void przestawRegule(kanal, kolekcja, reguly).then(odswiez);
        return;
      }
      kontekst.ustawKolekcje(kolekcja.id);
    },
    przy,
  );

  tresc.addEventListener(
    'dblclick',
    (zdarzenie) => {
      const cel = zdarzenie.target;
      if (!(cel instanceof HTMLElement) || !cel.classList.contains('dn-plakietka')) return;
      cel.contentEditable = 'true';
      cel.dataset.poprzednia = nazwaEtykiety(cel);
      cel.focus();
    },
    przy,
  );

  tresc.addEventListener(
    'focusout',
    (zdarzenie) => {
      const cel = zdarzenie.target;
      if (!(cel instanceof HTMLElement) || cel.contentEditable !== 'true') return;
      cel.contentEditable = 'false';
      const poprzednia = cel.dataset.poprzednia ?? '';
      const nowa = nazwaEtykiety(cel);
      if (poprzednia === '' || nowa === '' || poprzednia === nowa) return;
      void przemianujEtykiete(kanal, poprzednia, nowa).then(odswiez);
    },
    przy,
  );

  tresc.addEventListener(
    'keydown',
    (zdarzenie) => {
      const cel = zdarzenie.target;
      if (zdarzenie.key !== 'Delete' || !(cel instanceof HTMLElement)) return;
      if (!cel.classList.contains('dn-plakietka') || cel.contentEditable === 'true') return;
      zdarzenie.preventDefault();
      void usunEtykiete(kanal, nazwaEtykiety(cel)).then(odswiez);
    },
    przy,
  );

  tresc.addEventListener(
    'dragstart',
    (zdarzenie) => {
      const cel = zdarzenie.target;
      if (!(cel instanceof HTMLElement) || !cel.classList.contains('dn-plakietka')) return;
      zdarzenie.dataTransfer?.setData(RODZAJ_ETYKIETY, nazwaEtykiety(cel));
    },
    przy,
  );

  tresc.addEventListener(
    'dragover',
    (zdarzenie) => {
      const rodzaje = zdarzenie.dataTransfer?.types ?? [];
      if (rodzaje.includes(RODZAJ_PLIKU) || rodzaje.includes(RODZAJ_ETYKIETY)) {
        zdarzenie.preventDefault();
      }
    },
    przy,
  );

  tresc.addEventListener(
    'drop',
    (zdarzenie) => {
      zdarzenie.preventDefault();
      const dane = zdarzenie.dataTransfer;
      const cel = zdarzenie.target;
      if (dane === null || !(cel instanceof HTMLElement)) return;
      const idPliku = dane.getData(RODZAJ_PLIKU);
      const zrodlowa = dane.getData(RODZAJ_ETYKIETY);
      if (idPliku !== '') {
        void przypiszUpuszczony(kanal, cel, idPliku, kolekcje).then(odswiez);
        return;
      }
      if (zrodlowa === '') return;
      void polaczEtykiety(kanal, cel, zrodlowa, zdarzenie.altKey).then(odswiez);
    },
    przy,
  );

  zwiazKolejke(kanal, korzen, kontekst, przy, odswiez);
  kontekst.zglosOdswiezenie(() => {
    void odswiez();
  });
  void odswiez();
}

function zwiazKolejke(
  kanal: Kanal,
  korzen: ParentNode,
  kontekst: Kontekst,
  przy: AddEventListenerOptions,
  odswiezEtykiety: () => Promise<void>,
): void {
  const trescMozeMoze = panel(korzen, 'panel-kolejka');
  if (trescMozeMoze === null) return;
  const trescMoze = trescMozeMoze;
  const tresc = trescMoze;
  const sugestie = new Map<HTMLElement, LibrarySuggestion>();
  const duplikaty = new Map<HTMLElement, string[]>();

  async function odswiez(): Promise<void> {
    const [odpowiedzSugestii, odpowiedzDuplikatow] = await Promise.all([
      wywolaj(kanal, Command.LibrarySuggestionList, { limit: GRANICA }),
      wywolaj(kanal, Command.LibraryDuplicateScan, {
        ...(kontekst.kolekcja() === '' ? {} : { collectionId: kontekst.kolekcja() }),
        kinds: [LibraryDuplicateKind.Exact],
        limit: GRANICA,
      }),
    ]);
    const wykazSugestii = wartosc(odpowiedzSugestii)?.suggestions ?? [];
    const grupy = wartosc(odpowiedzDuplikatow)?.groups ?? [];
    sugestie.clear();
    duplikaty.clear();
    if (wykazSugestii.length === 0 && grupy.length === 0) {
      pustka(tresc, 'Rdzeń nie zwrócił pozycji do zatwierdzenia.');
      return;
    }
    zdejmijDzieci(tresc);
    for (const sugestia of wykazSugestii) {
      const wiersz = pozycja(tresc.ownerDocument, sugestia.reason, sugestia.value ?? '');
      sugestie.set(wiersz, sugestia);
      tresc.appendChild(wiersz);
    }
    for (const grupa of grupy) {
      const wiersz = pozycja(
        tresc.ownerDocument,
        `Duplikaty: ${String(grupa.fileIds.length)} plików`,
        grupa.checksum ?? '',
      );
      duplikaty.set(wiersz, grupa.fileIds);
      tresc.appendChild(wiersz);
    }
  }

  tresc.addEventListener(
    'click',
    (zdarzenie) => {
      const cel = zdarzenie.target;
      if (!(cel instanceof Element)) return;
      const wiersz = cel.closest('.dn-wykaz-modulu-poz');
      if (!(wiersz instanceof HTMLElement)) return;
      const sugestia = sugestie.get(wiersz);
      if (sugestia !== undefined) {
        void zastosujSugestie(kanal, sugestia.id, !zdarzenie.shiftKey)
          .then(odswiez)
          .then(odswiezEtykiety);
        return;
      }
      const grupa = duplikaty.get(wiersz);
      if (grupa === undefined || grupa.length < 2) return;
      void scalDuplikaty(kanal, grupa).then(odswiez);
    },
    przy,
  );

  kontekst.przyKolekcji(() => {
    void odswiez();
  });
  kontekst.zglosOdswiezenie(() => {
    void odswiez();
  });
  void odswiez();
}

function zbudujKolekcje(dokument: Document, kolekcja: LibraryCollection): HTMLElement {
  const wiersz = pozycja(dokument, kolekcja.name, String(kolekcja.fileCount));
  const znak = dokument.createElement('span');
  znak.className =
    kolekcja.ruleId === undefined ? 'dn-kropka dn-kropka--neutralna' : 'pt-tetno';
  znak.setAttribute('aria-hidden', 'true');
  wiersz.prepend(znak);
  wiersz.tabIndex = 0;
  return wiersz;
}

function zbudujSlownik(dokument: Document, etykiety: LibraryTag[]): HTMLElement {
  const pojemnik = dokument.createElement('div');
  pojemnik.style.display = 'flex';
  pojemnik.style.gap = 'var(--dn-od-1)';
  pojemnik.style.flexWrap = 'wrap';
  if (etykiety.length === 0) {
    pojemnik.appendChild(nota(dokument, 'Rdzeń nie zwrócił etykiet.'));
    return pojemnik;
  }
  for (const etykieta of etykiety) {
    const znak = plakietka(dokument, `${etykieta.name} · ${String(etykieta.fileCount)}`);
    znak.setAttribute('draggable', 'true');
    znak.tabIndex = 0;
    znak.dataset.etykieta = etykieta.name;
    pojemnik.appendChild(znak);
  }
  return pojemnik;
}

function nota(dokument: Document, zdanie: string): HTMLElement {
  const wezel = dokument.createElement('div');
  wezel.className = 'dn-meta';
  wezel.textContent = zdanie;
  return wezel;
}

function nazwaEtykiety(znak: HTMLElement): string {
  return znak.dataset.etykieta ?? (znak.textContent ?? '').split('·')[0].trim();
}

function podpisPrzycisku(podpis: string, liczba: number): string {
  const rdzen = podpis.replace(/\s*\(\d+\)\s*$/, '');
  return `${rdzen} (${String(liczba)})`;
}

async function przypiszUpuszczony(
  kanal: Kanal,
  cel: HTMLElement,
  idPliku: string,
  kolekcje: Map<HTMLElement, LibraryCollection>,
): Promise<void> {
  if (cel.classList.contains('dn-plakietka')) {
    await wywolaj(kanal, Command.LibraryTagSet, {
      fileId: idPliku,
      tags: [nazwaEtykiety(cel)],
    });
    return;
  }
  const wiersz = cel.closest('.dn-wykaz-modulu-poz');
  const kolekcja = wiersz instanceof HTMLElement ? kolekcje.get(wiersz) : undefined;
  if (kolekcja === undefined) return;
  await wywolaj(kanal, Command.LibraryCollectionAssign, {
    collectionId: kolekcja.id,
    fileIds: [idPliku],
  });
}

async function polaczEtykiety(
  kanal: Kanal,
  cel: HTMLElement,
  zrodlowa: string,
  szerzej: boolean,
): Promise<void> {
  if (cel.classList.contains('pt-etykieta')) {
    await wywolaj(kanal, Command.LibraryCollectionCreate, { name: zrodlowa });
    return;
  }
  if (!cel.classList.contains('dn-plakietka')) return;
  const docelowa = nazwaEtykiety(cel);
  if (docelowa === '' || docelowa === zrodlowa) return;
  if (szerzej) {
    await wywolaj(kanal, Command.LibraryThesaurusRelate, {
      sourceName: zrodlowa,
      targetName: docelowa,
      relation: LibraryThesaurusRelation.Related,
      remove: false,
    });
    return;
  }
  const odpowiedz = await wywolaj(kanal, Command.LibraryTagMerge, {
    sourceNames: [zrodlowa],
    targetName: docelowa,
  });
  const wynikMozeMoze = wartosc(odpowiedz);
  if (wynikMozeMoze === null) return;
  const wynikMoze = wynikMozeMoze;
  const wynik = wynikMoze;
  oglos('Library', `Scalono etykiety w „${docelowa}": ${String(wynik.affectedFiles)} plików.`);
}

async function przemianujEtykiete(kanal: Kanal, poprzednia: string, nowa: string): Promise<void> {
  await wywolaj(kanal, Command.LibraryTagUpdate, { name: poprzednia, newName: nowa });
}

async function usunEtykiete(kanal: Kanal, nazwa: string): Promise<void> {
  if (nazwa === '') return;
  await wywolaj(kanal, Command.LibraryTagRemove, { name: nazwa, confirm: true });
}

async function przestawRegule(
  kanal: Kanal,
  kolekcja: LibraryCollection,
  reguly: LibraryRule[],
): Promise<void> {
  const regula = reguly.find((kandydat) => kandydat.id === kolekcja.ruleId);
  if (regula === undefined) return;
  await wywolaj(kanal, Command.LibraryRuleSet, {
    rule: { ...regula, enabled: !regula.enabled },
    runNow: false,
  });
}

async function zdejmijRegule(
  kanal: Kanal,
  kolekcja: LibraryCollection,
  reguly: LibraryRule[],
): Promise<void> {
  const regula = reguly.find((kandydat) => kandydat.id === kolekcja.ruleId);
  if (regula === undefined) return;
  await wywolaj(kanal, Command.LibraryRuleRemove, { ruleId: regula.id, detachFiles: false });
}

async function sklasyfikuj(kanal: Kanal, kontekst: Kontekst): Promise<void> {
  const odpowiedz = await wywolaj(kanal, Command.LibraryClassifyRun, {
    ...(kontekst.kolekcja() === '' ? {} : { collectionId: kontekst.kolekcja() }),
    untaggedOnly: true,
    apply: false,
  });
  const wynikMozeMoze = wartosc(odpowiedz);
  if (wynikMozeMoze === null) return;
  const wynikMoze = wynikMozeMoze;
  const wynik = wynikMoze;
  oglos('Library', `Sugestii etykiet: ${String(wynik.suggestions.length)}.`);
}

async function zastosujSugestie(kanal: Kanal, id: string, przyjmij: boolean): Promise<void> {
  await wywolaj(kanal, Command.LibrarySuggestionApply, {
    suggestionIds: [id],
    accept: przyjmij,
  });
}

async function scalDuplikaty(kanal: Kanal, grupa: string[]): Promise<void> {
  const odpowiedz = await wywolaj(kanal, Command.LibraryDuplicateMerge, {
    targetFileId: grupa[0],
    sourceFileIds: grupa.slice(1),
    keepVersions: true,
  });
  const wynikMozeMoze = wartosc(odpowiedz);
  if (wynikMozeMoze === null) return;
  const wynikMoze = wynikMozeMoze;
  const wynik = wynikMoze;
  oglos('Library', `Scalono duplikatów: ${String(wynik.mergedCount)}.`);
}
