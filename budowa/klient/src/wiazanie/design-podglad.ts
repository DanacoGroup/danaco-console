// Preview Window: treść wskazanego zasobu, porównanie z jego pierwowzorem
// i trzy rozstrzygnięcia Operatora nad wariantem.

import { AssetContentDisposition, Command, type DesignAsset } from '../../../shared/contract.ts';
import {
  chip,
  nieGotowe,
  panel,
  poproszony,
  przycisk,
  tresc,
  tytul,
  wykonany,
  type Kontekst,
} from './design-wspolne.ts';
import { opisZasobu, wykazZasobow } from './design-zasoby.ts';
import { ostatniPrompt } from './design-prompt.ts';

const PROG_PODGLADU = 4 * 1024 * 1024;

export function zwiazPodglad(kontekst: Kontekst): void {
  const okno = panel(kontekst.korzen, 'panel-preview');
  const cialo = tresc(okno);
  if (okno === null || cialo === null) return;
  const plotno = cialo.querySelector<HTMLElement>('.dg-podglad');
  const suwak = cialo.querySelector<HTMLInputElement>('input[type="range"]');
  const stanChip = chip(cialo, 'oczekuje decyzji');

  const odswiez = (): void => {
    void wczytaj(kontekst, okno, plotno, suwak, stanChip);
  };
  kontekst.stan.odswiezenia.set('podglad', odswiez);

  przycisk(cialo, 'Akceptuj')?.addEventListener('click', () => {
    void rozstrzygnij(kontekst, 'akceptacja');
  }, kontekst.przy);
  przycisk(cialo, 'Regeneruj')?.addEventListener('click', () => {
    void rozstrzygnij(kontekst, 'regeneracja');
  }, kontekst.przy);
  przycisk(cialo, 'Odrzuć')?.addEventListener('click', () => {
    void rozstrzygnij(kontekst, 'odrzucenie');
  }, kontekst.przy);

  odswiez();
}

async function wczytaj(
  kontekst: Kontekst,
  okno: HTMLElement,
  plotno: HTMLElement | null,
  suwak: HTMLInputElement | null,
  stanChip: HTMLElement | null,
): Promise<void> {
  const zasob = wykazZasobow().find((kandydat) => kandydat.id === kontekst.stan.idZasobu);
  if (zasob === undefined) {
    tytul(okno, 'Preview Window');
    if (stanChip !== null) stanChip.textContent = 'brak wskazanego zasobu';
    nieGotowe(plotno, 'Wskaż zasób w Assets Panel, aby zobaczyć jego treść.');
    if (suwak !== null) suwak.disabled = true;
    return;
  }
  tytul(okno, `Preview Window — ${zasob.name ?? zasob.id}`);
  if (stanChip !== null) stanChip.textContent = opisZasobu(zasob);
  await postawWarstwy(kontekst, plotno, suwak, zasob);
}

async function postawWarstwy(
  kontekst: Kontekst,
  plotno: HTMLElement | null,
  suwak: HTMLInputElement | null,
  zasob: DesignAsset,
): Promise<void> {
  if (plotno === null) return;
  const biezacy = await tresc64(kontekst, zasob.id);
  if (biezacy === '') {
    nieGotowe(plotno, 'Rdzeń nie oddał treści tego zasobu.');
    if (suwak !== null) suwak.disabled = true;
    return;
  }
  const dokument = plotno.ownerDocument;
  const obraz = dokument.createElement('img');
  obraz.src = biezacy;
  obraz.alt = zasob.name ?? zasob.id;
  plotno.replaceChildren(obraz);

  // Porównanie „przed/po" ma sens tylko dla wariantu mającego pierwowzór.
  const pierwowzor = zasob.variantOfAssetId;
  if (pierwowzor === undefined) {
    if (suwak !== null) suwak.disabled = true;
    return;
  }
  const przed = await tresc64(kontekst, pierwowzor);
  if (przed === '') {
    if (suwak !== null) suwak.disabled = true;
    return;
  }
  const spod = dokument.createElement('img');
  spod.src = przed;
  spod.alt = `${zasob.name ?? zasob.id} — pierwowzór`;
  spod.style.position = 'absolute';
  spod.style.inset = '0';
  plotno.style.position = 'relative';
  plotno.append(spod);
  if (suwak === null) return;
  suwak.disabled = false;
  const nanies = (): void => {
    spod.style.opacity = String(1 - Number(suwak.value) / 100);
  };
  nanies();
  suwak.addEventListener('input', nanies, kontekst.przy);
}

async function tresc64(kontekst: Kontekst, idZasobu: string): Promise<string> {
  const odpowiedz = await poproszony(kontekst.kanal, Command.DesignAssetContentGet, {
    assetId: idZasobu,
    disposition: AssetContentDisposition.Inline,
    maxBytes: PROG_PODGLADU,
  });
  if (odpowiedz?.contentBase64 === undefined || odpowiedz.contentBase64 === '') return '';
  return `data:${odpowiedz.mediaType};base64,${odpowiedz.contentBase64}`;
}

async function rozstrzygnij(
  kontekst: Kontekst,
  rozstrzygniecie: 'akceptacja' | 'regeneracja' | 'odrzucenie',
): Promise<void> {
  const idZasobu = kontekst.stan.idZasobu;
  if (idZasobu === '') return;
  if (rozstrzygniecie === 'akceptacja') {
    await wykonany(kontekst.kanal, Command.DesignAssetFavoriteSet, {
      assetId: idZasobu,
      favorite: true,
    }, 'Akceptacja wariantu');
  } else if (rozstrzygniecie === 'odrzucenie') {
    await wykonany(kontekst.kanal, Command.DesignAssetRemove, {
      assetId: idZasobu,
    }, 'Odrzucenie wariantu');
    kontekst.stan.idZasobu = '';
  } else {
    const prompt = ostatniPrompt(kontekst.idKarty);
    if (prompt === undefined) return;
    await wykonany(kontekst.kanal, Command.DesignAssetGenerate, {
      windowId: kontekst.stan.idOkna,
      prompt,
      referenceAssetId: idZasobu,
    }, 'Ponowne generowanie');
  }
  kontekst.stan.odswiezenia.get('zasoby')?.();
  kontekst.stan.odswiezenia.get('podglad')?.();
}
