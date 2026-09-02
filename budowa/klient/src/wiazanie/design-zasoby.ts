// Assets Panel: wykaz zasobów rdzenia, ich znaczniki, ulubione, wniesienie,
// wyniesienie i szukanie w bibliotekach zewnętrznych.

import {
  AssetContentDisposition,
  ChangeKind,
  Command,
  DesignAssetKind,
  EventType,
  type DesignAsset,
} from '../../../shared/contract.ts';
import { zglosUchwyt } from '../polaczenie/rozdzielacz-zdarzen.ts';
import { oglos } from './ogloszenie.ts';
import {
  chwila,
  nieGotowe,
  panel,
  podpisRodzaju,
  poproszony,
  przycisk,
  rozmiar,
  tresc,
  wykonany,
  zdejmij,
  znacznik,
  type Kontekst,
} from './design-wspolne.ts';
import { dolozZasobDoPlanszy } from './design-plansza.ts';

const RODZAJE: string[] = [
  '',
  DesignAssetKind.Image,
  DesignAssetKind.Vector,
  DesignAssetKind.Composition,
  DesignAssetKind.Document,
  DesignAssetKind.Audio,
  DesignAssetKind.Video,
  DesignAssetKind.Archive,
];

// Rdzeń oddaje ścieżkę magazynu, spod której przeglądarka nie wczyta bajtów;
// miniatura musi więc iść osobnym wywołaniem treści, z zaciąganiem na żądanie.
const PROG_MINIATURY = 512 * 1024;

let rodzajFiltru = 0;
let zasoby: DesignAsset[] = [];

export function zwiazZasoby(kontekst: Kontekst): (() => void)[] {
  const okno = panel(kontekst.korzen, 'panel-assets');
  const cialo = tresc(okno);
  if (okno === null || cialo === null) return [];

  const siatka = cialo.querySelector<HTMLElement>('.dg-siatka');
  const wzor = przygotujWzor(siatka);
  const szukaj = cialo.querySelector<HTMLInputElement>('input[type="search"]');
  const filtr = cialo.querySelector<HTMLElement>('.dn-plakietka--sygnal');
  const stopka = [...cialo.querySelectorAll<HTMLElement>('.dg-narzedzia')].at(-1) ?? null;

  const doPlanszy = przycisk(stopka, '→ Design Board');
  // Przejścia do Studia i do Apps nie mają w kontrakcie komendy przenoszącej
  // zasób między modułami, więc oba przyciski schodzą zamiast kłamać.
  zdejmij(przycisk(stopka, '→ Studio'));
  zdejmij(przycisk(stopka, '→ Apps'));

  const odswiez = (): void => {
    void wczytaj(kontekst, okno, siatka, wzor);
  };
  kontekst.stan.odswiezenia.set('zasoby', odswiez);

  if (szukaj !== null) {
    szukaj.setAttribute(
      'title',
      'Nazwa albo znacznik. Prefiks + nadaje znacznik wskazanemu zasobowi, '
      + 'prefiks ? szuka w bibliotekach zewnętrznych.',
    );
    szukaj.addEventListener('input', () => {
      kontekst.stan.szukanie = szukaj.value.trim();
      if (!kontekst.stan.szukanie.startsWith('+') && !kontekst.stan.szukanie.startsWith('?')) {
        odswiez();
      }
    }, kontekst.przy);
    szukaj.addEventListener('keydown', (zdarzenie) => {
      if (zdarzenie.key !== 'Enter') return;
      zdarzenie.preventDefault();
      const wpis = szukaj.value.trim();
      if (wpis.startsWith('+')) void nadajZnacznik(kontekst, wpis.slice(1).trim(), odswiez);
      else if (wpis.startsWith('?')) void szukajZewnetrznie(kontekst, wpis.slice(1).trim(), siatka, wzor);
      else odswiez();
    }, kontekst.przy);
  }

  if (filtr !== null) {
    filtr.setAttribute('role', 'button');
    filtr.setAttribute('title', 'Przełącza rodzaj zasobu w wykazie.');
    filtr.textContent = 'wszystkie rodzaje';
    filtr.addEventListener('click', () => {
      rodzajFiltru = (rodzajFiltru + 1) % RODZAJE.length;
      const wybrany = RODZAJE[rodzajFiltru];
      filtr.textContent = wybrany === '' ? 'wszystkie rodzaje' : podpisRodzaju(wybrany);
      odswiez();
    }, kontekst.przy);
  }

  if (doPlanszy !== null) {
    doPlanszy.addEventListener('click', () => {
      void dolozZasobDoPlanszy(kontekst, kontekst.stan.idZasobu);
    }, kontekst.przy);
  }

  if (siatka !== null) {
    siatka.addEventListener('dragover', (zdarzenie) => {
      zdarzenie.preventDefault();
    }, kontekst.przy);
    siatka.addEventListener('drop', (zdarzenie) => {
      zdarzenie.preventDefault();
      void wniesUpuszczone(kontekst, zdarzenie, odswiez);
    }, kontekst.przy);
  }

  const odlaczenia = [
    zglosUchwyt(EventType.DesignAssetChanged, (zdarzenie) => {
      if (zdarzenie.asset.windowId !== kontekst.stan.idOkna) return;
      if (zdarzenie.change === ChangeKind.Deleted && zdarzenie.asset.id === kontekst.stan.idZasobu) {
        kontekst.stan.idZasobu = '';
      }
      odswiez();
    }),
  ];

  odswiez();
  return odlaczenia;
}

function przygotujWzor(siatka: HTMLElement | null): HTMLElement | null {
  const pierwsza = siatka?.querySelector<HTMLElement>('.dg-miniatura') ?? null;
  if (pierwsza === null) return null;
  const wzor = pierwsza.cloneNode(true) as HTMLElement;
  wzor.querySelector('.dg-nowy')?.remove();
  const podpis = wzor.querySelector('.podpis');
  if (podpis !== null) podpis.textContent = '';
  siatka?.replaceChildren();
  return wzor;
}

async function wczytaj(
  kontekst: Kontekst,
  okno: HTMLElement,
  siatka: HTMLElement | null,
  wzor: HTMLElement | null,
): Promise<void> {
  if (kontekst.stan.idOkna === '') return;
  const rodzaj = RODZAJE[rodzajFiltru];
  const odpowiedz = await poproszony(kontekst.kanal, Command.DesignAssetList, {
    windowId: kontekst.stan.idOkna,
    ...(rodzaj === '' ? {} : { kind: rodzaj as DesignAssetKind }),
    limit: 120,
  });
  if (odpowiedz === null) {
    znacznik(okno, '');
    nieGotowe(siatka, 'Rdzeń nie oddał wykazu zasobów.');
    return;
  }
  zasoby = przesiej(odpowiedz.assets, kontekst.stan.szukanie);
  znacznik(okno, `${String(odpowiedz.total ?? zasoby.length)} zasobów`);
  postawMiniatury(kontekst, siatka, wzor, zasoby);
}

function przesiej(wykaz: DesignAsset[], szukanie: string): DesignAsset[] {
  const fraza = szukanie.toLowerCase();
  if (fraza === '' || fraza.startsWith('+') || fraza.startsWith('?')) return wykaz;
  return wykaz.filter((zasob) => {
    const nazwa = (zasob.name ?? zasob.id).toLowerCase();
    const znaczniki = (zasob.tags ?? []).join(' ').toLowerCase();
    return nazwa.includes(fraza) || znaczniki.includes(fraza);
  });
}

function postawMiniatury(
  kontekst: Kontekst,
  siatka: HTMLElement | null,
  wzor: HTMLElement | null,
  wykaz: DesignAsset[],
): void {
  if (siatka === null || wzor === null) return;
  if (wykaz.length === 0) {
    nieGotowe(siatka, 'Okno nie ma jeszcze żadnego zasobu.');
    return;
  }
  const kafle = wykaz.map((zasob) => zbudujKafel(kontekst, wzor, zasob));
  siatka.replaceChildren(...kafle);
}

function zbudujKafel(kontekst: Kontekst, wzor: HTMLElement, zasob: DesignAsset): HTMLElement {
  const kafel = wzor.cloneNode(true) as HTMLElement;
  kafel.dataset.zasob = zasob.id;
  kafel.tabIndex = 0;
  kafel.setAttribute('aria-selected', String(zasob.id === kontekst.stan.idZasobu));
  kafel.setAttribute(
    'title',
    'Klik wskazuje zasób, podwójny klik przełącza ulubione, Delete usuwa, E wynosi plik.',
  );

  const podpis = kafel.querySelector('.podpis');
  if (podpis !== null) {
    podpis.textContent = zasob.name ?? zasob.id;
    const plakietka = kafel.ownerDocument.createElement('span');
    plakietka.className = 'dn-plakietka dn-na-koniec';
    plakietka.textContent = podpisRodzaju(zasob.kind);
    podpis.append(' ', plakietka);
  }
  if (zasob.favorite === true) {
    const wyroznienie = kafel.ownerDocument.createElement('span');
    wyroznienie.className = 'dn-plakietka dn-plakietka--sygnal dg-nowy';
    wyroznienie.textContent = 'ulubione';
    kafel.prepend(wyroznienie);
  }

  void wstawPodglad(kontekst, kafel, zasob);

  kafel.addEventListener('click', () => {
    kontekst.stan.idZasobu = zasob.id;
    const rodzenstwo = kafel.parentElement;
    if (rodzenstwo !== null) {
      for (let i = 0; i < rodzenstwo.children.length; i += 1) {
        rodzenstwo.children[i].setAttribute('aria-selected', String(rodzenstwo.children[i] === kafel));
      }
    }
    kontekst.stan.odswiezenia.get('podglad')?.();
    kontekst.stan.odswiezenia.get('fotografia')?.();
  }, kontekst.przy);

  kafel.addEventListener('dblclick', () => {
    void wykonany(kontekst.kanal, Command.DesignAssetFavoriteSet, {
      assetId: zasob.id,
      favorite: zasob.favorite !== true,
    }, 'Ulubione');
  }, kontekst.przy);

  kafel.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key === 'Delete') {
      zdarzenie.preventDefault();
      void wykonany(kontekst.kanal, Command.DesignAssetRemove, { assetId: zasob.id }, 'Usunięcie zasobu');
      return;
    }
    if (zdarzenie.key.toLowerCase() !== 'e') return;
    zdarzenie.preventDefault();
    if (zdarzenie.shiftKey) void wyniesWszystkie(kontekst);
    else void wynies(kontekst, zasob);
  }, kontekst.przy);

  return kafel;
}

async function wstawPodglad(
  kontekst: Kontekst,
  kafel: HTMLElement,
  zasob: DesignAsset,
): Promise<void> {
  const rysunek = kafel.querySelector('svg');
  if (rysunek === null) return;
  if (zasob.kind !== DesignAssetKind.Image && zasob.kind !== DesignAssetKind.Vector) {
    rysunek.remove();
    return;
  }
  const odpowiedz = await poproszony(kontekst.kanal, Command.DesignAssetContentGet, {
    assetId: zasob.id,
    disposition: AssetContentDisposition.Inline,
    maxBytes: PROG_MINIATURY,
  });
  if (odpowiedz?.contentBase64 === undefined || odpowiedz.contentBase64 === '') {
    rysunek.remove();
    return;
  }
  const obraz = kafel.ownerDocument.createElement('img');
  obraz.src = `data:${odpowiedz.mediaType};base64,${odpowiedz.contentBase64}`;
  obraz.alt = zasob.name ?? zasob.id;
  obraz.loading = 'lazy';
  rysunek.replaceWith(obraz);
}

async function wynies(kontekst: Kontekst, zasob: DesignAsset): Promise<void> {
  const odpowiedz = await wykonany(kontekst.kanal, Command.DesignAssetExport, {
    assetId: zasob.id,
    format: zasob.format ?? 'png',
  }, 'Wyniesienie zasobu');
  if (odpowiedz === null) return;
  oglos('Design', `${odpowiedz.fileName} · ${rozmiar(odpowiedz.sizeBytes)}`, 'informacja');
}

async function wyniesWszystkie(kontekst: Kontekst): Promise<void> {
  if (zasoby.length === 0) return;
  await wykonany(kontekst.kanal, Command.DesignAssetExportBatch, {
    assetIds: zasoby.map((zasob) => zasob.id),
    format: 'png',
    archive: true,
  }, 'Wyniesienie zestawu');
}

async function nadajZnacznik(
  kontekst: Kontekst,
  slowo: string,
  odswiez: () => void,
): Promise<void> {
  if (slowo === '' || kontekst.stan.idZasobu === '') return;
  const zasob = zasoby.find((kandydat) => kandydat.id === kontekst.stan.idZasobu);
  const znaczniki = [...new Set([...(zasob?.tags ?? []), slowo])];
  await wykonany(kontekst.kanal, Command.DesignAssetTagSet, {
    assetId: kontekst.stan.idZasobu,
    tags: znaczniki,
  }, 'Znacznik zasobu');
  odswiez();
}

async function wniesUpuszczone(
  kontekst: Kontekst,
  zdarzenie: DragEvent,
  odswiez: () => void,
): Promise<void> {
  const przyniesione = zdarzenie.dataTransfer?.files;
  const pliki = przyniesione === undefined ? [] : Array.from(przyniesione);
  if (pliki.length === 0 || kontekst.stan.idOkna === '') return;
  for (const plik of pliki) {
    const bufor = await plik.arrayBuffer();
    await wykonany(kontekst.kanal, Command.DesignAssetUpload, {
      windowId: kontekst.stan.idOkna,
      name: plik.name,
      kind: rodzajPliku(plik.type),
      contentBase64: naBase64(bufor),
      format: plik.name.split('.').at(-1) ?? '',
    }, `Wniesienie ${plik.name}`);
  }
  odswiez();
}

function rodzajPliku(typ: string): DesignAssetKind {
  if (typ === 'image/svg+xml') return DesignAssetKind.Vector;
  if (typ.startsWith('image/')) return DesignAssetKind.Image;
  if (typ.startsWith('audio/')) return DesignAssetKind.Audio;
  if (typ.startsWith('video/')) return DesignAssetKind.Video;
  if (typ === 'application/zip') return DesignAssetKind.Archive;
  return DesignAssetKind.Document;
}

function naBase64(bufor: ArrayBuffer): string {
  const bajty = new Uint8Array(bufor);
  let napis = '';
  for (let i = 0; i < bajty.length; i += 1) napis += String.fromCharCode(bajty[i]);
  return btoa(napis);
}

async function szukajZewnetrznie(
  kontekst: Kontekst,
  fraza: string,
  siatka: HTMLElement | null,
  wzor: HTMLElement | null,
): Promise<void> {
  if (fraza === '' || siatka === null || wzor === null || kontekst.stan.idOkna === '') return;
  const odpowiedz = await poproszony(kontekst.kanal, Command.DesignStockSearch, {
    query: fraza,
    windowId: kontekst.stan.idOkna,
    limit: 24,
  });
  if (odpowiedz === null || odpowiedz.assets.length === 0) {
    nieGotowe(siatka, 'Biblioteki zewnętrzne nie oddały wyniku dla tej frazy.');
    return;
  }
  const kafle = odpowiedz.assets.map((pozycja) => {
    const kafel = wzor.cloneNode(true) as HTMLElement;
    kafel.tabIndex = 0;
    kafel.setAttribute('title', 'Klik wnosi pozycję do okna.');
    const rysunek = kafel.querySelector('svg');
    if (pozycja.previewUrl !== undefined && rysunek !== null) {
      const obraz = kafel.ownerDocument.createElement('img');
      obraz.src = pozycja.previewUrl;
      obraz.alt = pozycja.title ?? pozycja.id;
      obraz.loading = 'lazy';
      rysunek.replaceWith(obraz);
    } else rysunek?.remove();
    const podpis = kafel.querySelector('.podpis');
    if (podpis !== null) {
      podpis.textContent = pozycja.title ?? pozycja.id;
      const plakietka = kafel.ownerDocument.createElement('span');
      plakietka.className = 'dn-plakietka dn-na-koniec';
      plakietka.textContent = pozycja.provider;
      podpis.append(' ', plakietka);
    }
    kafel.addEventListener('click', () => {
      void wykonany(kontekst.kanal, Command.DesignStockImport, {
        provider: pozycja.provider,
        externalId: pozycja.id,
        windowId: kontekst.stan.idOkna,
      }, 'Wniesienie z biblioteki');
    }, kontekst.przy);
    return kafel;
  });
  siatka.replaceChildren(...kafle);
  const nieudane = odpowiedz.providersFailed ?? [];
  if (nieudane.length > 0) {
    oglos('Design', `Biblioteki bez odpowiedzi: ${nieudane.join(', ')}.`, 'ostrzezenie');
  }
}

export function wykazZasobow(): DesignAsset[] {
  return zasoby;
}

export function opisZasobu(zasob: DesignAsset): string {
  const wymiar = zasob.width !== undefined && zasob.height !== undefined
    ? `${String(zasob.width)}×${String(zasob.height)}`
    : '';
  return [podpisRodzaju(zasob.kind), wymiar, chwila(zasob.createdAt)]
    .filter((czesc) => czesc !== '')
    .join(' · ');
}
